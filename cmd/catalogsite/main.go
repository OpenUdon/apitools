package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"
)

func main() {
	if err := run(os.Args[1:], os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "catalogsite:", err)
		os.Exit(1)
	}
}

func run(args []string, stderr io.Writer) error {
	flags := flag.NewFlagSet("catalogsite", flag.ContinueOnError)
	flags.SetOutput(stderr)
	baseURL := flags.String("base-url", "", "public HTTPS directory URL")
	commit := flags.String("source-commit", "", "full clean source commit used to build this binary")
	output := flags.String("out", "", "new output directory with an existing parent")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if flags.NArg() != 0 || *output == "" || !shaPattern.MatchString(*commit) {
		return errors.New("provide -base-url, -source-commit (40 lowercase hex characters), and -out")
	}
	base, err := parseBaseURL(*baseURL)
	if err != nil {
		return err
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return errors.New("missing Go build information")
	}
	if err := checkRevision(*commit, info.Settings); err != nil {
		return err
	}
	data, err := builtInSiteCatalog(*commit, base)
	if err != nil {
		return err
	}
	files, err := buildFiles(data)
	if err != nil {
		return err
	}
	return writeOutput(*output, files)
}

func parseBaseURL(raw string) (*url.URL, error) {
	base, err := url.Parse(raw)
	if err != nil || !publicHTTPS(base) || base.RawPath != "" || base.RawQuery != "" || base.ForceQuery || base.Fragment != "" {
		return nil, errors.New("-base-url must be an HTTPS directory URL without credentials, a port, query or fragment")
	}
	paths := strings.Trim(base.Path, "/")
	if paths != "" {
		for _, part := range strings.Split(paths, "/") {
			if !slugPattern.MatchString(part) {
				return nil, errors.New("-base-url paths must contain lowercase URL slugs")
			}
		}
		base.Path = "/" + paths + "/"
	} else {
		base.Path = "/"
	}
	return base, nil
}

func checkRevision(commit string, settings []debug.BuildSetting) error {
	values := make(map[string]string, len(settings))
	for _, setting := range settings {
		values[setting.Key] = setting.Value
	}
	if values["vcs"] != "git" || values["vcs.revision"] != commit || values["vcs.modified"] != "false" {
		return errors.New("build the generator with -buildvcs=true from a clean Git checkout of -source-commit")
	}
	return nil
}

func buildFiles(data siteCatalog) (map[string][]byte, error) {
	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}
	return map[string][]byte{"catalog.json": append(encoded, '\n')}, nil
}

func writeOutput(output string, files map[string][]byte) (err error) {
	names := make([]string, 0, len(files))
	for name := range files {
		if name == "SHA256SUMS" || name == "." || !fs.ValidPath(name) || strings.ContainsAny(name, " \\:\r\n\t") {
			return fmt.Errorf("invalid output file name %q", name)
		}
		names = append(names, name)
	}
	sort.Strings(names)
	if err := os.Mkdir(output, 0o755); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, os.RemoveAll(output))
		}
	}()
	var manifest strings.Builder
	for _, name := range names {
		contents := files[name]
		target := filepath.Join(output, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(target, contents, 0o644); err != nil {
			return err
		}
		fmt.Fprintf(&manifest, "%x  %s\n", sha256.Sum256(contents), name)
	}
	return os.WriteFile(filepath.Join(output, "SHA256SUMS"), []byte(manifest.String()), 0o644)
}
