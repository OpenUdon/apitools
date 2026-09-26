package apitools

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func resolvedLocalMaxBytes(maxBytes int64) int64 {
	if maxBytes <= 0 {
		return DefaultMaxBytes
	}
	return maxBytes
}

func validateInlineSpecContent(content []byte, maxBytes int64, label string) error {
	maxBytes = resolvedLocalMaxBytes(maxBytes)
	if int64(len(content)) > maxBytes {
		return fmt.Errorf("inline API source %q is larger than %d bytes", label, maxBytes)
	}
	return nil
}

func resolveLocalScanRoot(path string) (string, os.FileInfo, error) {
	resolved, err := resolveLocalPath(path)
	if err != nil {
		return "", nil, err
	}
	info, err := os.Lstat(resolved)
	if err != nil {
		return "", nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return "", nil, fmt.Errorf("local OpenAPI path %q is a symlink", path)
	}
	if !info.IsDir() {
		return "", nil, fmt.Errorf("local OpenAPI directory %q is not a directory", path)
	}
	return resolved, info, nil
}

func readLocalSpecFile(path string, maxBytes int64) ([]byte, error) {
	maxBytes = resolvedLocalMaxBytes(maxBytes)
	resolved, err := resolveLocalPath(path)
	if err != nil {
		return nil, err
	}
	info, err := os.Lstat(resolved)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("local OpenAPI document %q is a symlink", path)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("local OpenAPI document %q is a directory", path)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("local OpenAPI document %q is not a regular file", path)
	}
	if info.Size() > maxBytes {
		return nil, fmt.Errorf("local OpenAPI document %q is larger than %d bytes", path, maxBytes)
	}
	file, err := os.Open(resolved)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !os.SameFile(info, openedInfo) {
		return nil, fmt.Errorf("local OpenAPI document %q changed while opening", path)
	}
	if openedInfo.IsDir() {
		return nil, fmt.Errorf("local OpenAPI document %q is a directory", path)
	}
	if !openedInfo.Mode().IsRegular() {
		return nil, fmt.Errorf("local OpenAPI document %q is not a regular file", path)
	}
	if openedInfo.Size() > maxBytes {
		return nil, fmt.Errorf("local OpenAPI document %q is larger than %d bytes", path, maxBytes)
	}
	content, err := io.ReadAll(io.LimitReader(file, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(content)) > maxBytes {
		return nil, fmt.Errorf("local OpenAPI document %q is larger than %d bytes", path, maxBytes)
	}
	return content, nil
}

func lstatLocalLeaf(path string) (os.FileInfo, error) {
	resolved, err := resolveLocalPath(path)
	if err != nil {
		return nil, err
	}
	info, err := os.Lstat(resolved)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("local OpenAPI path %q is a symlink", path)
	}
	return info, nil
}

// resolveLocalPath resolves symlinked ancestors once while leaving the caller's
// final path component un-followed, so a selected file or scan root can still
// be checked for being a symlink itself.
func resolveLocalPath(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("local OpenAPI path is required")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	clean := filepath.Clean(abs)
	parent, err := filepath.EvalSymlinks(filepath.Dir(clean))
	if err != nil {
		return "", err
	}
	return filepath.Join(parent, filepath.Base(clean)), nil
}
