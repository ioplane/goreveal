package format

import "testing"

func TestDetectHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		header []byte
		want   Kind
	}{
		{name: "elf", header: []byte{0x7f, 'E', 'L', 'F'}, want: ELF},
		{name: "pe", header: []byte{'M', 'Z', 0x90, 0x00}, want: PE},
		{name: "macho 64", header: []byte{0xcf, 0xfa, 0xed, 0xfe}, want: MachO},
		{name: "fat macho", header: []byte{0xca, 0xfe, 0xba, 0xbe}, want: MachO},
		{name: "unknown", header: []byte{0x00, 0x01, 0x02, 0x03}, want: Unknown},
		{name: "short", header: []byte{0x7f, 'E'}, want: Unknown},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := DetectHeader(tc.header)
			if got != tc.want {
				t.Fatalf("DetectHeader() = %q, want %q", got, tc.want)
			}
		})
	}
}
