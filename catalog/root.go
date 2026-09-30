package catalog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultRegistryPath       = "cache.sqlite"
	DefaultOperationIndexPath = "operations.v1.json"
)

// RootOptions is installation configuration, never a discovery request field.
// RegistryPath and IndexPath are relative to Directory. An empty Directory
// means no root was configured; callers must return metadata-only leads.
type RootOptions struct {
	Directory    string `json:"-"`
	RegistryPath string `json:"-"`
	IndexPath    string `json:"-"`
}

// RootPaths contains resolved local paths. It is not part of an evidence wire.
// Missing files are allowed so readers can report missing evidence without
// creating a root, registry or index.
type RootPaths struct {
	Directory    string `json:"-"`
	RegistryPath string `json:"-"`
	IndexPath    string `json:"-"`
}

// ResolveRoot validates an explicit existing root without writing anything.
// Symlinked ancestors above the selected root are resolved, but the root itself
// and its registry/index descendants cannot be symlinks or special files.
func ResolveRoot(options RootOptions) (RootPaths, error) {
	directory := strings.TrimSpace(options.Directory)
	if directory == "" {
		if options.RegistryPath != "" || options.IndexPath != "" {
			return RootPaths{}, fmt.Errorf("registry/index configuration requires an explicit catalog root")
		}
		return RootPaths{}, nil
	}
	abs, err := filepath.Abs(directory)
	if err != nil {
		return RootPaths{}, err
	}
	if filepath.Dir(abs) == abs {
		return RootPaths{}, fmt.Errorf("catalog root must not be a filesystem root")
	}
	parent, err := filepath.EvalSymlinks(filepath.Dir(abs))
	if err != nil {
		return RootPaths{}, err
	}
	directory = filepath.Join(parent, filepath.Base(abs))
	info, err := os.Lstat(directory)
	if err != nil {
		return RootPaths{}, err
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return RootPaths{}, fmt.Errorf("catalog root must be a regular directory")
	}
	root, err := os.OpenRoot(directory)
	if err != nil {
		return RootPaths{}, err
	}
	defer root.Close()
	registry, err := rootFilePath(root, options.RegistryPath, DefaultRegistryPath)
	if err != nil {
		return RootPaths{}, fmt.Errorf("catalog registry: %w", err)
	}
	index, err := rootFilePath(root, options.IndexPath, DefaultOperationIndexPath)
	if err != nil {
		return RootPaths{}, fmt.Errorf("catalog index: %w", err)
	}
	if registry == index || strings.HasPrefix(registry, index+string(filepath.Separator)) || strings.HasPrefix(index, registry+string(filepath.Separator)) {
		return RootPaths{}, fmt.Errorf("catalog registry and index paths overlap")
	}
	return RootPaths{Directory: directory, RegistryPath: filepath.Join(directory, registry), IndexPath: filepath.Join(directory, index)}, nil
}

func rootFilePath(root *os.Root, value, fallback string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		value = fallback
	}
	value = filepath.FromSlash(value)
	if !filepath.IsLocal(value) || filepath.Clean(value) == "." {
		return "", fmt.Errorf("path must be a non-empty relative file path")
	}
	parts := strings.Split(value, string(filepath.Separator))
	for _, part := range parts {
		if part == ".." || part == "" || part == "." {
			return "", fmt.Errorf("path contains an invalid component")
		}
	}
	for i := range parts {
		path := filepath.Join(parts[:i+1]...)
		info, err := root.Lstat(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return "", err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("symlink component is not allowed")
		}
		if i < len(parts)-1 && !info.IsDir() || i == len(parts)-1 && !info.Mode().IsRegular() {
			return "", fmt.Errorf("path has a non-regular component")
		}
	}
	return value, nil
}
