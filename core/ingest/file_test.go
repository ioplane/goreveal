package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dantte-lp/goreveal/core/format"
)

func TestOpen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content []byte
		want    format.Kind
	}{
		{name: "elf", content: []byte{0x7f, 'E', 'L', 'F', 0x02}, want: format.ELF},
		{name: "pe", content: []byte{'M', 'Z', 0x90, 0x00}, want: format.PE},
		{name: "macho", content: []byte{0xcf, 0xfa, 0xed, 0xfe}, want: format.MachO},
		{name: "unknown", content: []byte("hello"), want: format.Unknown},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			path := filepath.Join(dir, "sample.bin")
			if err := os.WriteFile(path, tc.content, 0o644); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}

			got, err := Open(path)
			if err != nil {
				t.Fatalf("Open() error = %v", err)
			}

			if got.Format != tc.want {
				t.Fatalf("Open() format = %q, want %q", got.Format, tc.want)
			}
			if got.Path != path {
				t.Fatalf("Open() path = %q, want %q", got.Path, path)
			}
			if got.Size != int64(len(tc.content)) {
				t.Fatalf("Open() size = %d, want %d", got.Size, len(tc.content))
			}
		})
	}
}
