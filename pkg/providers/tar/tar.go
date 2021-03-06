package tar

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

// Dependencies .
type Dependencies struct {
	Logger zerolog.Logger
}

// Tarer is a tarer.
type Tarer struct {
	Logger zerolog.Logger
}

// NewTarer creates a provider which can tar / untar archives.
func NewTarer(logger zerolog.Logger) *Tarer {
	return &Tarer{
		Logger: logger,
	}
}

// Untar takes a destination path and a reader; a tar reader loops over the tar file
// creating the file structure at 'dst' along the way, and writing any files
func (t *Tarer) Untar(dst string, r io.Reader) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer func(gzr *gzip.Reader) {
		_ = gzr.Close()
	}(gzr)

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there is
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		if header.Typeflag != tar.TypeReg {
			continue
		}
		f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
		if err != nil {
			return err
		}
		// copy over contents
		if _, err := io.Copy(f, tr); err != nil {
			return err
		}
		// manually close here after each file operation; deferring would cause each file close
		// to wait until all operations have completed.
		if err := f.Close(); err != nil {
			return err
		}
	}
}
