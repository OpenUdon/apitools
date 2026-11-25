package catalogupdate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

func checkArtifactFilesystem(root *os.Root, expectedUUID string, requiredBytes uint64) error {
	dir, err := root.Open(".")
	if err != nil {
		return err
	}
	defer dir.Close()
	if expectedUUID != "" {
		var device, directory syscall.Stat_t
		if err := syscall.Stat(filepath.Join("/dev/disk/by-uuid", expectedUUID), &device); err != nil {
			return fmt.Errorf("find expected filesystem %q: %w", expectedUUID, err)
		}
		if err := syscall.Fstat(int(dir.Fd()), &directory); err != nil {
			return err
		}
		if device.Mode&syscall.S_IFMT != syscall.S_IFBLK || directory.Dev != device.Rdev {
			return fmt.Errorf("cache is outside expected filesystem %q", expectedUUID)
		}
	}
	var stat syscall.Statfs_t
	if err := syscall.Fstatfs(int(dir.Fd()), &stat); err != nil {
		return err
	}
	return checkArtifactCapacity(stat, requiredBytes)
}

func checkArtifactCapacity(stat syscall.Statfs_t, requiredBytes uint64) error {
	blockSize := stat.Frsize
	if blockSize == 0 {
		blockSize = stat.Bsize
	}
	if blockSize <= 0 {
		return errors.New("artifact cache filesystem reports an invalid block size")
	}
	if requiredBytes > 0 && stat.Bavail < 1+(requiredBytes-1)/uint64(blockSize) {
		return fmt.Errorf("artifact cache has insufficient capacity for %d bytes", requiredBytes)
	}
	return nil
}
