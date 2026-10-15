package catalog

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

const (
	DefaultMaxRegistryBytes       int64 = 64 << 20 // 64 MiB
	DefaultMaxIndexBytes          int64 = 64 << 20 // 64 MiB
	DefaultMaxMetadataBytes       int64 = 2 << 20  // 2 MiB
	DefaultMaxSelectedSourceBytes int64 = 8 << 20  // 8 MiB
)

var (
	ErrActiveJournal = errors.New("cannot snapshot catalog with active SQLite journal or WAL file")
	ErrRootOverlap   = errors.New("source catalog root and snapshot destination overlap")
	ErrSnapshotDrift = errors.New("catalog source files drifted or changed during snapshot")
)

// SnapshotOptions specifies limits and optional metadata or selected source paths.
type SnapshotOptions struct {
	// MetadataSource is the path to an external metadata file.
	MetadataSource string
	// MetadataDestination is the relative destination path inside targetDir.
	MetadataDestination string
	// SelectedSourcePath is a relative specification path inside the source root.
	SelectedSourcePath string

	// Maximum allowed sizes for individual copied files.
	MaxRegistryBytes       int64
	MaxIndexBytes          int64
	MaxMetadataBytes       int64
	MaxSelectedSourceBytes int64

	// NormalizeJournal converts destination SQLite databases to PRAGMA journal_mode=DELETE.
	NormalizeJournal bool
}

// SnapshotRoot creates an isolated, read-only-safe snapshot of the catalog in targetDir.
// It verifies that no active journals exist, copies files preserving their relative paths,
// enforces ceilings, verifies that no source files drifted during copying, and normalizes
// destination SQLite databases to DELETE journal mode.
//
// Failures clean up any created files in targetDir and leave no producer mutation.
func SnapshotRoot(paths RootPaths, targetDir string, options SnapshotOptions) error {
	dir := strings.TrimSpace(paths.Directory)
	if dir == "" {
		return errors.New("cannot snapshot an unconfigured catalog root")
	}
	info, err := os.Lstat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("catalog root must be a regular directory")
	}

	targetDir = strings.TrimSpace(targetDir)
	if targetDir == "" {
		return errors.New("snapshot destination directory must not be empty")
	}
	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return err
	}
	if filepath.Dir(absTarget) == absTarget {
		return errors.New("snapshot destination cannot be filesystem root")
	}
	if rootsOverlap(dir, absTarget) {
		return ErrRootOverlap
	}

	// Resolve relative registry and index paths from paths
	relReg, err := filepath.Rel(dir, paths.RegistryPath)
	if err != nil || !isSafeRelativePath(relReg) {
		return fmt.Errorf("invalid registry relative path: %w", err)
	}
	relIdx, err := filepath.Rel(dir, paths.IndexPath)
	if err != nil || !isSafeRelativePath(relIdx) {
		return fmt.Errorf("invalid index relative path: %w", err)
	}

	// Apply ceiling defaults
	maxReg := options.MaxRegistryBytes
	if maxReg <= 0 {
		maxReg = DefaultMaxRegistryBytes
	}
	maxIdx := options.MaxIndexBytes
	if maxIdx <= 0 {
		maxIdx = DefaultMaxIndexBytes
	}
	maxMeta := options.MaxMetadataBytes
	if maxMeta <= 0 {
		maxMeta = DefaultMaxMetadataBytes
	}
	maxSel := options.MaxSelectedSourceBytes
	if maxSel <= 0 {
		maxSel = DefaultMaxSelectedSourceBytes
	}

	// Validate options
	if options.SelectedSourcePath != "" && !isSafeRelativePath(options.SelectedSourcePath) {
		return fmt.Errorf("selected source path must be a safe relative path: %q", options.SelectedSourcePath)
	}
	if options.MetadataDestination != "" {
		if !isSafeRelativePath(options.MetadataDestination) {
			return fmt.Errorf("metadata destination must be a safe relative path: %q", options.MetadataDestination)
		}
		if options.MetadataSource == "" {
			return errors.New("metadata destination requires a configured metadata source")
		}
	} else if options.MetadataSource != "" {
		return errors.New("metadata source requires a configured metadata destination")
	}

	// Destination collision and sidecar overlap checks
	destPaths := []string{relReg, relIdx}
	if options.SelectedSourcePath != "" {
		destPaths = append(destPaths, options.SelectedSourcePath)
	}
	if options.MetadataDestination != "" {
		destPaths = append(destPaths, options.MetadataDestination)
	}
	for i := 0; i < len(destPaths); i++ {
		for j := i + 1; j < len(destPaths); j++ {
			p1 := filepath.Clean(destPaths[i])
			p2 := filepath.Clean(destPaths[j])
			if p1 == p2 {
				return fmt.Errorf("destination paths collide: %q and %q", p1, p2)
			}
			if strings.HasPrefix(p1, p2+string(filepath.Separator)) || strings.HasPrefix(p2, p1+string(filepath.Separator)) {
				return fmt.Errorf("destination paths overlap hierarchically: %q and %q", p1, p2)
			}
		}
		// Sidecar overlap checks against registry
		for _, suffix := range []string{"-wal", "-shm", "-journal"} {
			if destPaths[i] != relReg && filepath.Clean(destPaths[i]) == filepath.Clean(relReg+suffix) {
				return fmt.Errorf("destination path %q overlaps registry sidecar", destPaths[i])
			}
		}
	}

	// Verify no active WAL or journal sidecars exist at source
	if err := checkActiveJournals(paths.RegistryPath); err != nil {
		return err
	}

	// Open source root with os.OpenRoot to confine file reads
	srcRoot, err := os.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer srcRoot.Close()

	// Ensure destination directory exists and is a regular directory
	if targetInfo, err := os.Lstat(absTarget); err == nil {
		if !targetInfo.IsDir() || targetInfo.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("snapshot destination %q must be a regular directory", absTarget)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(absTarget, 0o700); err != nil {
			return err
		}
	} else {
		return err
	}

	var createdFiles []string
	rollback := func() {
		for _, f := range createdFiles {
			_ = os.Remove(f)
		}
	}

	// 1. Copy registry if present (optional evidence)
	var regData []byte
	regExisted := false
	regData, err = readConfinedFile(srcRoot, relReg, maxReg)
	if err == nil {
		regExisted = true
		dest := filepath.Join(absTarget, relReg)
		if err := writeDestinationFile(dest, regData); err != nil {
			rollback()
			return err
		}
		createdFiles = append(createdFiles, dest)
	} else if !errors.Is(err, os.ErrNotExist) {
		rollback()
		return fmt.Errorf("reading registry: %w", err)
	}

	// 2. Copy index if present (optional evidence)
	var idxData []byte
	idxExisted := false
	idxData, err = readConfinedFile(srcRoot, relIdx, maxIdx)
	if err == nil {
		idxExisted = true
		dest := filepath.Join(absTarget, relIdx)
		if err := writeDestinationFile(dest, idxData); err != nil {
			rollback()
			return err
		}
		createdFiles = append(createdFiles, dest)
	} else if !errors.Is(err, os.ErrNotExist) {
		rollback()
		return fmt.Errorf("reading index: %w", err)
	}

	// 3. Copy selected source if configured (required evidence)
	var selData []byte
	if options.SelectedSourcePath != "" {
		selData, err = readConfinedFile(srcRoot, options.SelectedSourcePath, maxSel)
		if err != nil {
			rollback()
			return fmt.Errorf("reading selected source %q: %w", options.SelectedSourcePath, err)
		}
		dest := filepath.Join(absTarget, options.SelectedSourcePath)
		if err := writeDestinationFile(dest, selData); err != nil {
			rollback()
			return err
		}
		createdFiles = append(createdFiles, dest)
	}

	// 4. Copy metadata if configured (required evidence)
	var metaData []byte
	if options.MetadataSource != "" {
		metaData, err = readExternalFile(options.MetadataSource, maxMeta)
		if err != nil {
			rollback()
			return fmt.Errorf("reading metadata %q: %w", options.MetadataSource, err)
		}
		dest := filepath.Join(absTarget, options.MetadataDestination)
		if err := writeDestinationFile(dest, metaData); err != nil {
			rollback()
			return err
		}
		createdFiles = append(createdFiles, dest)
	}

	// 5. Full-set drift check
	// Recheck registry
	if regExisted {
		after, err := readConfinedFile(srcRoot, relReg, maxReg)
		if err != nil || !bytes.Equal(regData, after) {
			rollback()
			return ErrSnapshotDrift
		}
	} else {
		if _, err := srcRoot.Lstat(relReg); !errors.Is(err, os.ErrNotExist) {
			rollback()
			return ErrSnapshotDrift
		}
	}
	// Recheck index
	if idxExisted {
		after, err := readConfinedFile(srcRoot, relIdx, maxIdx)
		if err != nil || !bytes.Equal(idxData, after) {
			rollback()
			return ErrSnapshotDrift
		}
	} else {
		if _, err := srcRoot.Lstat(relIdx); !errors.Is(err, os.ErrNotExist) {
			rollback()
			return ErrSnapshotDrift
		}
	}
	// Recheck selected source
	if options.SelectedSourcePath != "" {
		after, err := readConfinedFile(srcRoot, options.SelectedSourcePath, maxSel)
		if err != nil || !bytes.Equal(selData, after) {
			rollback()
			return ErrSnapshotDrift
		}
	}
	// Recheck metadata
	if options.MetadataSource != "" {
		after, err := readExternalFile(options.MetadataSource, maxMeta)
		if err != nil || !bytes.Equal(metaData, after) {
			rollback()
			return ErrSnapshotDrift
		}
	}
	// Recheck active journals
	if err := checkActiveJournals(paths.RegistryPath); err != nil {
		rollback()
		return ErrSnapshotDrift
	}

	// 6. Destination-only SQLite journal mode normalization
	if regExisted {
		destReg := filepath.Join(absTarget, relReg)
		if shouldNormalizeSQLite(destReg, options.NormalizeJournal) {
			if err := normalizeJournalMode(destReg); err != nil {
				rollback()
				return fmt.Errorf("normalizing destination sqlite journal mode: %w", err)
			}
		}
	}

	return nil
}

func checkActiveJournals(registryPath string) error {
	for _, suffix := range []string{"-wal", "-journal"} {
		info, err := os.Lstat(registryPath + suffix)
		if err == nil && info.Size() > 0 {
			return fmt.Errorf("%w: active journal sidecar %s", ErrActiveJournal, suffix)
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}

func readConfinedFile(root *os.Root, relPath string, maxBytes int64) ([]byte, error) {
	f, err := root.Open(relPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("file %q is not a regular file", relPath)
	}
	if info.Size() > maxBytes {
		return nil, fmt.Errorf("file %q size %d exceeds maximum of %d bytes", relPath, info.Size(), maxBytes)
	}
	data, err := io.ReadAll(io.LimitReader(f, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxBytes {
		return nil, fmt.Errorf("file %q exceeds maximum of %d bytes", relPath, maxBytes)
	}
	return data, nil
}

func readExternalFile(path string, maxBytes int64) ([]byte, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	info, err := os.Lstat(abs)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("file %q is not a regular file", path)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("file %q cannot be a symlink", path)
	}
	if info.Size() > maxBytes {
		return nil, fmt.Errorf("file %q size %d exceeds maximum of %d bytes", path, info.Size(), maxBytes)
	}
	f, err := os.Open(abs)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxBytes {
		return nil, fmt.Errorf("file %q exceeds maximum of %d bytes", path, maxBytes)
	}
	return data, nil
}

func writeDestinationFile(destPath string, data []byte) error {
	parent := filepath.Dir(destPath)
	if err := os.MkdirAll(parent, 0o700); err != nil {
		return err
	}
	return os.WriteFile(destPath, data, 0o600)
}

func shouldNormalizeSQLite(path string, forceNormalize bool) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	header := make([]byte, 20)
	n, err := io.ReadFull(f, header)
	if err != nil || n < 20 {
		return false
	}
	if !bytes.Equal(header[:16], []byte("SQLite format 3\x00")) {
		return false
	}
	// Check if WAL mode (write version 2 or read version 2) or forced
	return forceNormalize || header[18] == 2 || header[19] == 2
}

func normalizeJournalMode(dbPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u := url.URL{Scheme: "file", Path: filepath.ToSlash(dbPath), RawQuery: "mode=rw"}
	db, err := sql.Open("sqlite", u.String())
	if err != nil {
		return err
	}
	defer db.Close()

	var mode string
	if err := db.QueryRowContext(ctx, "PRAGMA journal_mode=DELETE").Scan(&mode); err != nil {
		return err
	}
	if strings.ToLower(mode) != "delete" {
		return fmt.Errorf("expected delete journal mode, got %q", mode)
	}
	return nil
}

func rootsOverlap(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	a = filepath.Clean(a)
	b = filepath.Clean(b)
	return a == b || strings.HasPrefix(a, b+string(filepath.Separator)) || strings.HasPrefix(b, a+string(filepath.Separator))
}

func isSafeRelativePath(p string) bool {
	if p == "" || !filepath.IsLocal(p) || filepath.Clean(p) == "." {
		return false
	}
	if strings.ContainsAny(p, "\\\x00\r\n\t,") || strings.HasPrefix(p, "-") {
		return false
	}
	parts := strings.Split(p, string(filepath.Separator))
	for _, part := range parts {
		if part == ".." || part == "." || part == "" {
			return false
		}
	}
	return true
}
