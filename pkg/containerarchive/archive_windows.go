package containerarchive

import (
	"io"

	"github.com/docker/docker/pkg/archive"
)

// container is not supported by Windows
func container(path string) error {
	return nil
}

func invokeUnpack(decompressedArchive io.ReadCloser,
	dest string,
	options *archive.TarOptions) error {
	// Windows is different to Linux here because Windows does not support
	// container. Hence there is no point sandboxing a containered process to
	// do the unpack. We call inline instead within the daemon process.
	return archive.Unpack(decompressedArchive, dest, options)
}
