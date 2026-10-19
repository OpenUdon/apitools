package apitools

import (
	"context"
	"net/url"
)

// FetchedDocument is the result of a guarded fetch.
type FetchedDocument struct {
	// Content is the response body, decoded if it was gzip, within the byte limit.
	Content []byte
	// FinalURL is the URL that answered, after any redirects.
	FinalURL *url.URL
	// ContentType is the response Content-Type header, unparsed.
	ContentType string
}

// FetchDocument downloads rawURL with the client's download guard and returns
// the bytes, the final URL and the Content-Type. It applies the same host and
// port checks at dial time and on every redirect, the same timeout and the same
// wire and decoded byte limits as Import, and honours AllowUnsafeHosts and
// AllowedPorts. A response outside 2xx is an HTTPStatusError.
//
// Unlike Import it does not check that the body is an OpenAPI or Swagger
// document, does not use the cache, and writes nothing: it is for a consumer
// that stores or inspects companion files, documentation or other API source
// kinds itself and must not bypass the guard to fetch them. The result is
// untrusted data.
func (c *Client) FetchDocument(ctx context.Context, rawURL string) (FetchedDocument, error) {
	content, finalURL, contentType, err := c.downloadBoundedWithAccept(ctx, rawURL, "")
	if err != nil {
		return FetchedDocument{}, err
	}
	return FetchedDocument{Content: content, FinalURL: finalURL, ContentType: contentType}, nil
}
