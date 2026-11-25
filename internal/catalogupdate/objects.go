package catalogupdate

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
)

const maxObjectBytes = 1_000_000_000

// ErrObjectIntegrity identifies evidence confined to one cached object.
var ErrObjectIntegrity = errors.New("catalog artifact integrity check failed")

func (o Object) filename() (string, error) {
	if !validDigest(o.SHA256) || o.Bytes <= 0 || o.Bytes > maxObjectBytes {
		return "", errors.New("invalid catalog artifact digest or byte count")
	}
	return "sha256-" + o.SHA256, nil
}

// PutObject makes a byte-identical artifact durable before registering it.
// Publishing the file first permits safe orphan reuse if the process stops
// between filesystem and database commits.
func (s *Store) PutObject(ctx context.Context, object Object, source io.Reader) error {
	name, err := object.filename()
	if err != nil {
		return err
	}
	lock, err := lockCatalogCache(ctx, s.cache)
	if err != nil {
		return err
	}
	defer lock.Close()
	const stage = ".catalog-update-stage"
	if err := s.cache.Remove(stage); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	existing, err := s.OpenObject(ctx, object)
	if err == nil {
		if err := errors.Join(existing.Sync(), existing.Close()); err != nil {
			return err
		}
		return s.registerObject(ctx, object)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		if !errors.Is(err, ErrObjectIntegrity) {
			return err
		}
		info, statErr := s.cache.Lstat(name)
		if statErr != nil && !errors.Is(statErr, fs.ErrNotExist) {
			return statErr
		}
		// A regular digest file can be restored by verified replacement
		// bytes. Leave foreign entries such as symlinks and directories alone.
		if statErr == nil && !info.Mode().IsRegular() {
			return err
		}
	}
	if err := checkArtifactFilesystem(s.cache, "", uint64(object.Bytes)); err != nil {
		return err
	}
	file, err := s.cache.OpenFile(stage, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	defer s.cache.Remove(stage)
	copyErr := verifyObject(ctx, source, file, object)
	if copyErr == nil {
		copyErr = file.Sync()
	}
	if err := errors.Join(copyErr, file.Close()); err != nil {
		return err
	}
	if err := s.cache.Rename(stage, name); err != nil {
		return err
	}
	return s.registerObject(ctx, object)
}

func (s *Store) registerObject(ctx context.Context, object Object) error {
	// A crash before the directory sync leaves no SQLite reference. Sync even
	// when reusing bytes so any colliding writer's link is made durable too.
	if err := syncDirectory(s.cache); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	result, err := s.db.ExecContext(ctx, `
INSERT INTO update_object VALUES (?, ?)
ON CONFLICT(digest) DO UPDATE SET bytes = excluded.bytes
WHERE update_object.bytes = excluded.bytes`, object.SHA256, object.Bytes)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("catalog artifact size conflicts with registry")
	}
	return nil
}

// OpenObject verifies stored bytes and positions the opened file at the
// beginning. A missing, substituted, or corrupt file cannot be used offline.
func (s *Store) OpenObject(ctx context.Context, object Object) (*os.File, error) {
	name, err := object.filename()
	if err != nil {
		return nil, err
	}
	info, err := s.cache.Lstat(name)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() != object.Bytes {
		return nil, fmt.Errorf("%w: invalid file type or size", ErrObjectIntegrity)
	}
	file, err := openObjectRead(s.cache, name)
	if err != nil {
		return nil, err
	}
	opened, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, err
	}
	if !os.SameFile(info, opened) {
		_ = file.Close()
		return nil, fmt.Errorf("%w: artifact changed during open", ErrObjectIntegrity)
	}
	if err := verifyObject(ctx, file, nil, object); err != nil {
		_ = file.Close()
		return nil, err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		_ = file.Close()
		return nil, err
	}
	return file, nil
}

func verifyObject(ctx context.Context, source io.Reader, destination io.Writer, object Object) error {
	digest := sha256.New()
	var dst io.Writer = digest
	if destination != nil {
		dst = io.MultiWriter(destination, digest)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if _, err := io.CopyN(dst, source, object.Bytes); err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return fmt.Errorf("%w: artifact is shorter than declared", ErrObjectIntegrity)
		}
		return fmt.Errorf("read catalog artifact: %w", err)
	}
	var extra [1]byte
	if n, err := io.ReadFull(source, extra[:]); n != 0 {
		return fmt.Errorf("%w: artifact exceeds its declared size", ErrObjectIntegrity)
	} else if err != nil && err != io.EOF {
		return err
	}
	if hex.EncodeToString(digest.Sum(nil)) != object.SHA256 {
		return fmt.Errorf("%w: digest mismatch", ErrObjectIntegrity)
	}
	return ctx.Err()
}
