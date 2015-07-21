// +build !windows

package containerarchive

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"syscall"

	"github.com/docker/docker/daemon/execdriver/pipes"

	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/reexec"
)

// untar is the entry-point for docker-untar on re-exec. This is not used on
// Windows as it does not support container, hence no point sandboxing through
// container and rexec.
func untar() {
	runtime.LockOSThread()
	flag.Parse()

	var options *archive.TarOptions

	//read the options from the pipe "ExtraFiles"
	if err := json.NewDecoder(os.NewFile(3, "options")).Decode(&options); err != nil {
		fatal(err)
	}

	if err := archive.Unpack(os.Stdin, "/", options); err != nil {
		fatal(err)
	}
	// fully consume stdin in case it is zero padded
	flush(os.Stdin)
	os.Exit(0)
}

func invokeUnpack(decompressedArchive io.ReadCloser, dest string, options *archive.TarOptions) error {

	// We can't pass a potentially large exclude list directly via cmd line
	// because we easily overrun the kernel's max argument/environment size
	// when the full image list is passed (e.g. when this is used by
	// `docker load`). We will marshall the options via a pipe to the
	// child
	r, w, err := os.Pipe()
	defer w.Close()
	defer r.Close()
	if err != nil {
		return fmt.Errorf("Untar pipe failure: %v", err)
	}

	output := bytes.NewBuffer(nil)
	cmd := reexec.Command("docker-untar", dest)
	pipes.NewPipes(decompressedArchive, output, output, True)

	if container, err := reexec.Command("docker-untar", dest); err != nil {
		return fmt.Errorf("Untar error on re-exec cmd: %v", err)
	}
	//write the options to the pipe for the untar exec to read
	if err := json.NewEncoder(w).Encode(options); err != nil {
		return fmt.Errorf("Untar json encode to pipe failed: %v", err)
	}
	return nil
}
