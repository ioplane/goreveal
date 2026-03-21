package ingest

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/dantte-lp/goreveal/core/format"
)

// File contains basic metadata gathered during binary ingestion.
type File struct {
	Path   string
	Size   int64
	Format format.Kind
}

// Open reads a binary header and returns the basic ingestion result.
func Open(path string) (File, error) {
	info, err := os.Stat(path)
	if err != nil {
		return File{}, fmt.Errorf("stat %q: %w", path, err)
	}

	f, err := os.Open(path)
	if err != nil {
		return File{}, fmt.Errorf("open %q: %w", path, err)
	}
	defer f.Close()

	header := make([]byte, 16)
	n, err := io.ReadFull(f, header)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		return File{}, fmt.Errorf("read %q header: %w", path, err)
	}

	return File{
		Path:   path,
		Size:   info.Size(),
		Format: format.DetectHeader(header[:n]),
	}, nil
}
