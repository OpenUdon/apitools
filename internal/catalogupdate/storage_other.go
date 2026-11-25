//go:build !linux

package catalogupdate

import (
	"context"
	"errors"
	"os"
)

func checkStateFilesystem(_ *os.File) error {
	return errors.New("catalog maintenance SQLite state requires Linux")
}

func syncArtifactFilesystem(_ *os.Root) error {
	return errors.New("catalog maintenance cache sync requires Linux")
}

func openObjectRead(root *os.Root, name string) (*os.File, error) {
	return root.Open(name)
}

func openStatePath(_ *os.Root, _ string) (*os.File, error) {
	return nil, errors.New("catalog maintenance SQLite state requires Linux")
}

func lockCatalogCache(_ context.Context, _ *os.Root) (*os.File, error) {
	return nil, errors.New("catalog maintenance cache locking requires Linux")
}

func lockCatalogFile(_ context.Context, _ *os.Root, _ string, _ bool) (*os.File, error) {
	return nil, errors.New("catalog maintenance run locking requires Linux")
}

func checkStateWritable(_ string) error {
	return errors.New("catalog maintenance SQLite state requires Linux")
}

func checkArtifactFilesystem(_ *os.Root, expectedUUID string, requiredBytes uint64) error {
	if expectedUUID != "" || requiredBytes != 0 {
		return errors.New("artifact cache volume and capacity checks require Linux")
	}
	return nil
}
