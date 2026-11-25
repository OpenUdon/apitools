package catalogupdate

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
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

func checkStateFilesystem(file *os.File) error {
	var stat syscall.Statfs_t
	if err := syscall.Fstatfs(int(file.Fd()), &stat); err != nil {
		return err
	}
	return rejectStateFilesystem(uint32(stat.Type))
}

func rejectStateFilesystem(kind uint32) error {
	switch kind {
	case 0xef53, 0x58465342, 0x9123683e, 0xf2f52010, 0x2fc12fc1, // ext2/3/4, XFS, Btrfs, F2FS, ZFS
		0x3153464a, 0x52654973, 0x24051905, 0xca451a4e: // JFS, ReiserFS, UBIFS, Bcachefs
		return nil
	default:
		return fmt.Errorf("catalog state needs a supported native local filesystem (found %#x)", kind)
	}
}

func openObjectRead(root *os.Root, name string) (*os.File, error) {
	return openPinnedFile(root, name, syscall.O_RDONLY)
}

func openStatePath(root *os.Root, name string) (*os.File, error) {
	// Closing an ordinary file descriptor can release this process's
	// SQLite advisory locks. A path-only descriptor supports stat checks
	// without taking part in file locking.
	return openPinnedFile(root, name, unix.O_PATH)
}

func openPinnedFile(root *os.Root, name string, flags int) (*os.File, error) {
	dir, err := root.Open(".")
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	fd, err := syscall.Openat(int(dir.Fd()), name, flags|syscall.O_NOFOLLOW|syscall.O_NONBLOCK|syscall.O_CLOEXEC, 0o600)
	if err != nil {
		return nil, err
	}
	return os.NewFile(uintptr(fd), name), nil
}

func syncArtifactFilesystem(root *os.Root) error {
	dir, err := root.Open(".")
	if err != nil {
		return err
	}
	return errors.Join(unix.Syncfs(int(dir.Fd())), dir.Close())
}

func lockCatalogCache(ctx context.Context, root *os.Root) (*os.File, error) {
	return lockCatalogFile(ctx, root, ".catalog-update.lock", true)
}

func checkStateWritable(path string) error {
	return unix.Access(path, unix.W_OK)
}

func lockCatalogFile(ctx context.Context, root *os.Root, lockName string, wait bool) (*os.File, error) {
	file, err := openPinnedFile(root, lockName, syscall.O_CREAT|syscall.O_RDWR)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() {
		_ = file.Close()
		return nil, errors.New("invalid catalog cache lock file")
	}
	for {
		err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			if current, err := root.Lstat(lockName); err != nil || !os.SameFile(info, current) {
				_ = file.Close()
				return nil, errors.New("catalog cache lock changed during acquisition")
			}
			return file, nil
		}
		if err != syscall.EWOULDBLOCK && err != syscall.EINTR {
			_ = file.Close()
			return nil, err
		}
		if !wait && err == syscall.EWOULDBLOCK {
			_ = file.Close()
			return nil, errors.New("catalog maintenance is already running for this state database")
		}
		select {
		case <-ctx.Done():
			_ = file.Close()
			return nil, ctx.Err()
		case <-time.After(20 * time.Millisecond):
		}
	}
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
