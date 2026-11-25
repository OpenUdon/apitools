package catalogupdate

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestArtifactCacheCapacityAndUnknownUUID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "objects")
	for _, tc := range []struct {
		uuid string
		size uint64
		want string
	}{
		{size: ^uint64(0), want: "insufficient capacity"},
		{uuid: "apitools-fixture-unallocated-filesystem", want: "expected filesystem"},
	} {
		root, err := OpenArtifactCache(path, tc.uuid, tc.size)
		if root != nil || err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("cache = %v, error = %v, want %q", root, err, tc.want)
		}
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("created cache directory despite failed readiness check: %v", err)
		}
	}
}

func TestArtifactCacheUsesFilesystemAllocationUnits(t *testing.T) {
	stat := syscall.Statfs_t{Bsize: 4096, Frsize: 512, Bavail: 2}
	if err := checkArtifactCapacity(stat, 1024); err != nil {
		t.Fatal(err)
	}
	if err := checkArtifactCapacity(stat, 1025); err == nil {
		t.Fatal("capacity should use fragment size when the filesystem reports one")
	}
	stat.Frsize = 0
	if err := checkArtifactCapacity(stat, 8192); err != nil {
		t.Fatalf("capacity should use block size when fragment size is absent: %v", err)
	}
}
