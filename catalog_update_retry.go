package apitools

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/OpenUdon/apitools/internal/catalogupdate"
)

const catalogUpdatePublisherAttempts = 3
const catalogUpdateRetryWindow = 24 * time.Hour
const catalogUpdateLastRepresentableUnix = 253402300799

type catalogUpdateRetry struct {
	Reason        string    `json:"reason"`
	Origin        string    `json:"origin"`
	HTTPStatus    int       `json:"http_status,omitempty"`
	Attempt       int       `json:"attempt"`
	NextAttemptAt time.Time `json:"next_attempt_at"`
	Shared        bool      `json:"shared,omitempty"`
	OutsideWindow bool      `json:"outside_window,omitempty"`
	HintSaturated bool      `json:"hint_saturated,omitempty"`
}

type catalogUpdateOriginBlock struct {
	Until         time.Time `json:"until"`
	OutsideWindow bool      `json:"outside_window,omitempty"`
}

type catalogUpdatePublisherError struct {
	cause  error
	retry  catalogUpdateRetry
	notify time.Time
}

func (e *catalogUpdatePublisherError) Error() string { return "retryable publisher " + e.retry.Reason }
func (e *catalogUpdatePublisherError) Unwrap() error { return e.cause }

func catalogUpdateRetryFailure(raw string, err error, status int, after, remaining, reset string) error {
	if err == nil {
		return nil
	}
	u, parseErr := url.Parse(raw)
	if parseErr != nil || u.Hostname() == "" {
		return err
	}
	now := time.Now().UTC()
	retryAfter, saturated := catalogUpdateRetryAfter(after, now)
	retry := catalogUpdateRetry{Origin: versionOrigin(u), HTTPStatus: status, HintSaturated: saturated}
	var timeout net.Error
	switch {
	case status == http.StatusTooManyRequests || status == http.StatusForbidden && (remaining == "0" || !retryAfter.IsZero()):
		retry.Reason, retry.Shared = "rate-limit", true
		// GitHub's metadata endpoint supplies an absolute seconds value.
		// Other publishers can supply the standard Retry-After header.
		if retry.Origin == "https://api.github.com:443" {
			epoch, resetSaturated := catalogUpdateRetryNumber(reset, catalogUpdateLastRepresentableUnix)
			if t := time.Unix(epoch, 0).UTC(); !resetSaturated && t.After(now) && t.After(retryAfter) {
				retryAfter = t
				retry.HintSaturated = false
			}
		}
	case status >= 500 && status <= 599:
		retry.Reason, retry.Shared = "publisher-server-error", !retryAfter.IsZero()
	case status == http.StatusRequestTimeout || errors.Is(err, context.DeadlineExceeded) || errors.As(err, &timeout) && timeout.Timeout():
		retry.Reason = "timeout"
	default:
		return err
	}
	return &catalogUpdatePublisherError{cause: err, retry: retry, notify: retryAfter}
}

func catalogUpdateRetryAfter(raw string, now time.Time) (time.Time, bool) {
	if n, saturated := catalogUpdateRetryNumber(raw, catalogUpdateLastRepresentableUnix-now.Unix()); n >= 0 {
		return time.Unix(now.Unix()+n, int64(now.Nanosecond())).UTC(), saturated
	}
	if date, err := http.ParseTime(raw); err == nil {
		if date.Before(now) {
			return now, false
		}
		return date.UTC(), false
	}
	return time.Time{}, false
}

func catalogUpdateRetryNumber(raw string, ceiling int64) (int64, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.Trim(raw, "0123456789") != "" {
		return -1, false
	}
	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || value > uint64(ceiling) {
		return ceiling, true
	}
	return int64(value), false
}

func (e *catalogUpdatePublisherError) nextAttempt(prior *catalogupdate.Outcome) catalogUpdateRetry {
	retry := e.retry
	retry.Attempt = 1
	if prior != nil {
		retry.Attempt = readCatalogUpdateRetry(*prior).Attempt + 1
	}
	retry.NextAttemptAt = e.notify
	retry.OutsideWindow = retry.NextAttemptAt.After(time.Now().UTC().Add(catalogUpdateRetryWindow))
	if retry.NextAttemptAt.IsZero() {
		delay := 10 * time.Second
		if retry.Shared {
			delay = time.Minute
		}
		retry.NextAttemptAt = time.Now().UTC().Add(delay * time.Duration(retry.Attempt))
	}
	return retry
}

func readCatalogUpdateRetry(out catalogupdate.Outcome) catalogUpdateRetry {
	var evidence struct {
		Retry catalogUpdateRetry `json:"publisher_retry"`
	}
	_ = json.Unmarshal(out.Evidence, &evidence)
	return evidence.Retry
}

func readCatalogUpdateRetryBlocks(out catalogupdate.Outcome) map[string]catalogUpdateOriginBlock {
	var evidence struct {
		Blocks map[string]catalogUpdateOriginBlock `json:"publisher_retry_blocks"`
	}
	_ = json.Unmarshal(out.Evidence, &evidence)
	return evidence.Blocks
}

func mergeCatalogUpdateRetryBlocks(blockedUntil map[string]catalogUpdateOriginBlock, out catalogupdate.Outcome, currentPass bool) {
	for origin, block := range readCatalogUpdateRetryBlocks(out) {
		if currentPass && !block.OutsideWindow {
			// Yield the origin for the remainder of this serial pass.
			block.Until = time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)
		}
		known := blockedUntil[origin]
		if block.OutsideWindow || !known.OutsideWindow && block.Until.After(known.Until) {
			blockedUntil[origin] = block
		}
	}
}

func catalogUpdateArtifactOrigin(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Hostname() == "" {
		return ""
	}
	return versionOrigin(u)
}

func catalogUpdateBlockForSource(source catalogUpdateSource, blockedUntil map[string]catalogUpdateOriginBlock) (string, catalogUpdateOriginBlock) {
	u, err := url.Parse(source.URL)
	if err != nil || u.Hostname() == "" {
		return "", catalogUpdateOriginBlock{}
	}
	origin := versionOrigin(u)
	block := blockedUntil[origin]
	if u.Hostname() == "raw.githubusercontent.com" {
		const github = "https://api.github.com:443"
		if metadata := blockedUntil[github]; metadata.OutsideWindow || (!block.OutsideWindow && metadata.Until.After(block.Until)) {
			return github, metadata
		}
	}
	return origin, block
}
