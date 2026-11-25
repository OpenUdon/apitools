//go:build !linux

package catalogupdate

import (
	"errors"
	"os"
)

func checkArtifactFilesystem(_ *os.Root, expectedUUID string, requiredBytes uint64) error {
	if expectedUUID != "" || requiredBytes != 0 {
		return errors.New("artifact cache volume and capacity checks require Linux")
	}
	return nil
}
