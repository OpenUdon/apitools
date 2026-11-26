package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
	"time"
)

//go:embed changes/*.json
var monthlyRecords embed.FS

type monthlyRecord struct {
	Period           string   `json:"period"`
	FromCommit       string   `json:"from_commit"`
	ToCommit         string   `json:"to_commit"`
	FirstPublication bool     `json:"first_publication"`
	Added            []string `json:"added"`
	Changed          []string `json:"changed"`
	Removed          []string `json:"removed"`
}

func runCompare(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("catalogsite compare", flag.ContinueOnError)
	flags.SetOutput(stderr)
	from := flags.String("from", "", "catalog.json from the earlier reviewed revision")
	to := flags.String("to", "", "catalog.json from the later reviewed revision")
	period := flags.String("period", "", "review month (YYYY-MM)")
	first := flags.Bool("first-publication", false, "there is no previous site publication")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if flags.NArg() != 0 || *from == "" || *to == "" || !validPeriod(*period) {
		return errors.New("provide -from, -to and -period (YYYY-MM)")
	}
	before, err := readComparisonCatalog(*from)
	if err != nil {
		return err
	}
	after, err := readComparisonCatalog(*to)
	if err != nil {
		return err
	}
	if before.SourceCommit == after.SourceCommit {
		return errors.New("the two catalogs must identify different source commits")
	}
	record := compareCatalogs(*period, *first, before, after)
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(record)
}

func readComparisonCatalog(name string) (siteCatalog, error) {
	info, err := os.Stat(name)
	if err != nil {
		return siteCatalog{}, err
	}
	if !info.Mode().IsRegular() || info.Size() > 16<<20 {
		return siteCatalog{}, fmt.Errorf("%s must be a catalog file no larger than 16 MiB", name)
	}
	file, err := os.Open(name)
	if err != nil {
		return siteCatalog{}, err
	}
	defer file.Close()
	var data siteCatalog
	if err := decodeDocument(io.LimitReader(file, 16<<20), &data); err != nil {
		return siteCatalog{}, fmt.Errorf("read %s: %w", name, err)
	}
	if data.SchemaVersion != 1 || !shaPattern.MatchString(data.SourceCommit) || len(data.Providers) == 0 {
		return siteCatalog{}, fmt.Errorf("%s is missing the version 1 catalog identity or providers", name)
	}
	previous := ""
	for _, provider := range data.Providers {
		if !slugPattern.MatchString(provider.ID) || provider.ID <= previous {
			return siteCatalog{}, fmt.Errorf("%s provider IDs must be unique sorted slugs", name)
		}
		previous = provider.ID
	}
	return data, nil
}

func compareCatalogs(period string, first bool, before, after siteCatalog) monthlyRecord {
	record := monthlyRecord{
		Period: period, FromCommit: before.SourceCommit, ToCommit: after.SourceCommit, FirstPublication: first,
		Added: []string{}, Changed: []string{}, Removed: []string{},
	}
	previous := make(map[string]siteProvider, len(before.Providers))
	for _, provider := range before.Providers {
		previous[provider.ID] = provider
	}
	for _, provider := range after.Providers {
		old, found := previous[provider.ID]
		delete(previous, provider.ID)
		if !found {
			record.Added = append(record.Added, provider.ID)
			continue
		}
		// A base URL change alone does not change reviewed provider metadata.
		old.URL, provider.URL = "", ""
		if !reflect.DeepEqual(old, provider) {
			record.Changed = append(record.Changed, provider.ID)
		}
	}
	for id := range previous {
		record.Removed = append(record.Removed, id)
	}
	slices.Sort(record.Removed)
	return record
}

func loadChangeHistory() ([]monthlyRecord, error) {
	entries, err := monthlyRecords.ReadDir("changes")
	if err != nil {
		return nil, err
	}
	history := make([]monthlyRecord, 0, len(entries))
	for i, entry := range entries {
		contents, err := monthlyRecords.ReadFile("changes/" + entry.Name())
		if err != nil {
			return nil, err
		}
		var record monthlyRecord
		if err := decodeDocument(bytes.NewReader(contents), &record); err != nil {
			return nil, fmt.Errorf("read monthly record %s: %w", entry.Name(), err)
		}
		if !validPeriod(record.Period) || record.Period+".json" != entry.Name() || record.FirstPublication != (i == 0) ||
			!shaPattern.MatchString(record.FromCommit) || !shaPattern.MatchString(record.ToCommit) || record.FromCommit == record.ToCommit {
			return nil, fmt.Errorf("invalid monthly record identity in %s", entry.Name())
		}
		seen := make(map[string]bool)
		for _, list := range [][]string{record.Added, record.Changed, record.Removed} {
			if list == nil || !slices.IsSorted(list) {
				return nil, fmt.Errorf("invalid provider list in %s", entry.Name())
			}
			for _, id := range list {
				if !slugPattern.MatchString(id) || seen[id] {
					return nil, fmt.Errorf("invalid or repeated provider ID in %s", entry.Name())
				}
				seen[id] = true
			}
		}
		history = append(history, record)
	}
	slices.Reverse(history)
	return history, nil
}

func validPeriod(value string) bool {
	period, err := time.Parse("2006-01", value)
	return err == nil && period.Format("2006-01") == value
}

func decodeDocument(reader io.Reader, destination any) error {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	var trailing json.RawMessage
	if err := decoder.Decode(&trailing); err != io.EOF {
		return errors.New("expected one complete JSON record")
	}
	return nil
}
