package catalog

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

// Run FIFO probes in a subprocess so a blocking open is a bounded test failure.
func TestSnapshotFIFORefusal(t *testing.T) {
	if os.Getenv("APITOOLS_TEST_FIFO") == "1" {
		src := t.TempDir()
		if err := syscall.Mkfifo(filepath.Join(src, "pipe"), 0600); err != nil {
			t.Fatal(err)
		}
		paths, err := ResolveRoot(RootOptions{Directory: src})
		if err != nil {
			t.Fatal(err)
		}
		if err := SnapshotRoot(paths, filepath.Join(t.TempDir(), "snapshot"), SnapshotOptions{SelectedSourcePath: "pipe"}); err == nil {
			t.Fatal("accepted FIFO")
		}
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestSnapshotFIFORefusal$")
	cmd.Env = append(os.Environ(), "APITOOLS_TEST_FIFO=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("FIFO probe: %v (%v) %s", err, ctx.Err(), out)
	}
}
