package runtime

import (
	"encoding/binary"
	"testing"

	"github.com/dantte-lp/goreveal/schema"
)

func TestReadMetadataFixture(t *testing.T) {
	t.Parallel()

	got, err := ReadMetadata("../../corpus/fixtures/go-elf-buildinfo-linux-amd64/fixture.bin")
	if err != nil {
		t.Fatalf("ReadMetadata() error = %v", err)
	}

	if got.FirstModuleDataAddr == 0 {
		t.Fatalf("ReadMetadata() firstmoduledata = %#v", got)
	}
	if got.TrustSummary != "symbol_backed" {
		t.Fatalf("ReadMetadata() trust summary = %q, want %q", got.TrustSummary, "symbol_backed")
	}
	if got.FirstModuleDataFromGoModuleFallback {
		t.Fatalf("ReadMetadata() rich fixture unexpectedly marked go.module fallback = %#v", got)
	}
	if got.GopclntabAddr == 0 || got.GopclntabSize == 0 {
		t.Fatalf("ReadMetadata() gopclntab = %#v", got)
	}
	if got.ELFPclntabHeaderMagic == "" || got.ELFPclntabHeaderMagicKind != "known" || got.ELFPclntabHeaderQuantum != 1 || got.ELFPclntabHeaderPointerSize != 8 {
		t.Fatalf("ReadMetadata() ELF pclntab header = %#v", got)
	}
	if got.ELFPclntabFunctionCountHint == 0 || got.ELFPclntabFileCountHint == 0 {
		t.Fatalf("ReadMetadata() ELF pclntab header hints = %#v", got)
	}
	if got.ELFTextSectionAddr == 0 || got.ELFTextSectionEndInclusive == 0 || got.ELFTextSectionEndInclusive < got.ELFTextSectionAddr {
		t.Fatalf("ReadMetadata() ELF text section range = %#v", got)
	}
	if got.ELFPclntabFuncnametabOffsetHint == 0 || got.ELFPclntabPctabOffsetHint == 0 || got.ELFPclntabFunctabOffsetHint == 0 {
		t.Fatalf("ReadMetadata() ELF pclntab offset hints = %#v", got)
	}
	if got.ELFFunctabLastPCOffsetHint == 0 || !got.ELFFunctabPCOffsetsMonotonic || got.ELFFunctabLastPCOffsetHint < got.ELFFunctabFirstPCOffsetHint {
		t.Fatalf("ReadMetadata() ELF functab pc offset hints = %#v", got)
	}
	if got.ELFFunctabFirstPCAddrHint == 0 || got.ELFFunctabLastPCAddrHint == 0 || !got.ELFFunctabPCAddrHintsWithinText || got.ELFFunctabLastPCAddrHint < got.ELFFunctabFirstPCAddrHint {
		t.Fatalf("ReadMetadata() ELF functab pc addr hints = %#v", got)
	}
	if len(got.ELFFunctabPCAddrSample) == 0 || !got.ELFFunctabPCAddrSampleAllWithinText {
		t.Fatalf("ReadMetadata() ELF functab pc addr sample = %#v", got)
	}
	if got.ModuledataPCHeaderAddr == 0 || !got.ModuledataPCHeaderMatchesGopclntab {
		t.Fatalf("ReadMetadata() moduledata pcheader bridge = %#v", got)
	}
	if got.ModuledataFuncnametabSliceWordIndex != 1 || got.ModuledataFuncnametabAddr == 0 || got.ModuledataFuncnametabLen == 0 || got.ModuledataFuncnametabCap == 0 || !got.ModuledataFuncnametabWithinGopclntab {
		t.Fatalf("ReadMetadata() moduledata funcnametab bridge = %#v", got)
	}
	if got.ModuledataCutabSliceWordIndex != 4 || got.ModuledataCutabAddr == 0 || got.ModuledataCutabLen == 0 || got.ModuledataCutabCap == 0 || !got.ModuledataCutabWithinGopclntab {
		t.Fatalf("ReadMetadata() moduledata cutab bridge = %#v", got)
	}
	if got.ModuledataFiletabSliceWordIndex != 7 || got.ModuledataFiletabAddr == 0 || got.ModuledataFiletabLen == 0 || got.ModuledataFiletabCap == 0 || !got.ModuledataFiletabWithinGopclntab {
		t.Fatalf("ReadMetadata() moduledata filetab bridge = %#v", got)
	}
	if got.ModuledataPctabSliceWordIndex != 10 || got.ModuledataPctabAddr == 0 || got.ModuledataPctabLen == 0 || got.ModuledataPctabCap == 0 || !got.ModuledataPctabWithinGopclntab {
		t.Fatalf("ReadMetadata() moduledata pctab bridge = %#v", got)
	}
	if got.ModuledataPclntableSliceWordIndex != 13 || got.ModuledataPclntableAddr == 0 || got.ModuledataPclntableLen == 0 || got.ModuledataPclntableCap == 0 || !got.ModuledataPclntableWithinGopclntab {
		t.Fatalf("ReadMetadata() moduledata pclntable bridge = %#v", got)
	}
	if got.TypelinkAddr == 0 || got.TypelinkSize == 0 {
		t.Fatalf("ReadMetadata() typelink = %#v", got)
	}
	if got.TypelinkCount == 0 {
		t.Fatalf("ReadMetadata() typelink count = %#v", got)
	}
	if got.ItablinkAddr == 0 || got.ItablinkSize == 0 || got.ItablinkCount == 0 {
		t.Fatalf("ReadMetadata() itablink = %#v", got)
	}
	if len(got.TypelinkSample) == 0 {
		t.Fatalf("ReadMetadata() typelink sample = %#v", got)
	}
	if len(got.TypelinkSample) > 8 {
		t.Fatalf("ReadMetadata() typelink sample too large = %#v", got.TypelinkSample)
	}
	if got.TypelinkMinOffset == 0 || got.TypelinkMaxOffset == 0 || got.TypelinkMaxOffset < got.TypelinkMinOffset {
		t.Fatalf("ReadMetadata() typelink min/max = %#v", got)
	}
	if got.TypelinkNegativeCount+got.TypelinkNonNegativeCount != got.TypelinkCount {
		t.Fatalf("ReadMetadata() typelink sign counts = %#v", got)
	}
	if got.GoModuleAddr == 0 || got.GoModuleSize == 0 {
		t.Fatalf("ReadMetadata() go.module = %#v", got)
	}
	if !got.FirstModuleDataInGoModule {
		t.Fatalf("ReadMetadata() firstmoduledata/go.module cross-check = %#v", got)
	}
	if got.GoModuleWordSize == 0 || len(got.GoModuleWordSample) == 0 {
		t.Fatalf("ReadMetadata() go.module word sample = %#v", got)
	}
	if got.ModuledataTypelinkSliceWordIndex == 0 || got.ModuledataTypelinkLen != got.TypelinkCount || got.ModuledataTypelinkCap != got.TypelinkCount {
		t.Fatalf("ReadMetadata() moduledata typelinks slice = %#v", got)
	}
	if got.ModuledataItablinkSliceWordIndex == 0 || got.ModuledataItablinkLen != got.ItablinkCount || got.ModuledataItablinkCap != got.ItablinkCount {
		t.Fatalf("ReadMetadata() moduledata itablinks slice = %#v", got)
	}
	if got.ModuledataMemoryRangesWordIndex == 0 {
		t.Fatalf("ReadMetadata() moduledata memory ranges word index = %#v", got)
	}
	if got.ModuledataNoptrdataAddr == 0 || got.ModuledataNoptrdataEnd == 0 || got.ModuledataNoptrdataEnd <= got.ModuledataNoptrdataAddr {
		t.Fatalf("ReadMetadata() moduledata noptrdata range = %#v", got)
	}
	if got.ModuledataDataAddr == 0 || got.ModuledataDataEnd == 0 || got.ModuledataDataEnd <= got.ModuledataDataAddr {
		t.Fatalf("ReadMetadata() moduledata data range = %#v", got)
	}
	if got.ModuledataBssAddr == 0 || got.ModuledataBssEnd == 0 || got.ModuledataBssEnd <= got.ModuledataBssAddr {
		t.Fatalf("ReadMetadata() moduledata bss range = %#v", got)
	}
	if got.ModuledataNoptrbssAddr == 0 || got.ModuledataNoptrbssEnd == 0 || got.ModuledataNoptrbssEnd <= got.ModuledataNoptrbssAddr {
		t.Fatalf("ReadMetadata() moduledata noptrbss range = %#v", got)
	}
	if got.ModuledataRodataWordIndex == 0 || got.ModuledataRodataAddr == 0 || got.ModuledataRodataEnd == 0 || got.ModuledataRodataEnd <= got.ModuledataRodataAddr {
		t.Fatalf("ReadMetadata() moduledata rodata range = %#v", got)
	}
	if got.ModuledataTextWordIndex == 0 || got.ModuledataTextAddr == 0 || got.ModuledataTextEndInclusive == 0 || got.ModuledataTextEndInclusive < got.ModuledataTextAddr {
		t.Fatalf("ReadMetadata() moduledata text range = %#v", got)
	}
	if got.TypelinkResolvedBaseAddr == 0 || len(got.TypelinkResolvedSample) == 0 || got.TypelinkResolvedWithinRodataCount == 0 {
		t.Fatalf("ReadMetadata() typelink semantic bridge = %#v", got)
	}
	if !got.TypelinkAllResolvedWithinRodata {
		t.Fatalf("ReadMetadata() typelink all-within-rodata = %#v", got)
	}
	if got.ModuledataTypesAddr == 0 || got.ModuledataETypesAddr == 0 || got.ModuledataETypesAddr <= got.ModuledataTypesAddr {
		t.Fatalf("ReadMetadata() moduledata types range = %#v", got)
	}
	if got.ModuledataTypesRangeWordIndex == 0 {
		t.Fatalf("ReadMetadata() moduledata types range word index = %#v", got)
	}
	if got.TypelinkResolvedWithinTypesCount == 0 || !got.TypelinkAllResolvedWithinTypes {
		t.Fatalf("ReadMetadata() typelink types-range semantics = %#v", got)
	}
	if got.Provenance.Source != "core.runtime.elf" || got.Provenance.Confidence != "medium" {
		t.Fatalf("ReadMetadata() provenance = %#v", got.Provenance)
	}
}

func TestReadMetadataStrippedFixtureFallsBackToGoModule(t *testing.T) {
	t.Parallel()

	got, err := ReadMetadata("../../corpus/fixtures/go-elf-stripped-linux-amd64/fixture.bin")
	if err != nil {
		t.Fatalf("ReadMetadata() error = %v", err)
	}

	if got.GoModuleAddr == 0 || got.GoModuleSize == 0 {
		t.Fatalf("ReadMetadata() stripped go.module = %#v", got)
	}
	if got.ELFPclntabHeaderMagic == "" || got.ELFPclntabHeaderMagicKind != "known" || got.ELFPclntabHeaderQuantum != 1 || got.ELFPclntabHeaderPointerSize != 8 {
		t.Fatalf("ReadMetadata() stripped ELF pclntab header = %#v", got)
	}
	if got.ELFPclntabFunctionCountHint == 0 || got.ELFPclntabFileCountHint == 0 {
		t.Fatalf("ReadMetadata() stripped ELF pclntab header hints = %#v", got)
	}
	if got.ELFTextSectionAddr == 0 || got.ELFTextSectionEndInclusive == 0 || got.ELFTextSectionEndInclusive < got.ELFTextSectionAddr {
		t.Fatalf("ReadMetadata() stripped ELF text section range = %#v", got)
	}
	if got.ELFFunctabLastPCOffsetHint == 0 || !got.ELFFunctabPCOffsetsMonotonic || got.ELFFunctabLastPCOffsetHint < got.ELFFunctabFirstPCOffsetHint {
		t.Fatalf("ReadMetadata() stripped ELF functab pc offset hints = %#v", got)
	}
	if got.ELFFunctabFirstPCAddrHint == 0 || got.ELFFunctabLastPCAddrHint == 0 || !got.ELFFunctabPCAddrHintsWithinText || got.ELFFunctabLastPCAddrHint < got.ELFFunctabFirstPCAddrHint {
		t.Fatalf("ReadMetadata() stripped ELF functab pc addr hints = %#v", got)
	}
	if len(got.ELFFunctabPCAddrSample) == 0 || !got.ELFFunctabPCAddrSampleAllWithinText {
		t.Fatalf("ReadMetadata() stripped ELF functab pc addr sample = %#v", got)
	}
	if got.TrustSummary != "go_module_fallback" {
		t.Fatalf("ReadMetadata() stripped trust summary = %q, want %q", got.TrustSummary, "go_module_fallback")
	}
	if got.FirstModuleDataAddr == 0 {
		t.Fatalf("ReadMetadata() stripped firstmoduledata missing = %#v", got)
	}
	if got.FirstModuleDataAddr != got.GoModuleAddr {
		t.Fatalf("ReadMetadata() stripped fallback addr = %#x, want %#x", got.FirstModuleDataAddr, got.GoModuleAddr)
	}
	if !got.FirstModuleDataFromGoModuleFallback {
		t.Fatalf("ReadMetadata() stripped fallback source bit = %#v", got)
	}
	if !got.FirstModuleDataInGoModule {
		t.Fatalf("ReadMetadata() stripped go.module cross-check = %#v", got)
	}
	if !got.ModuledataPCHeaderMatchesGopclntab || !got.ModuledataFuncnametabWithinGopclntab || !got.ModuledataPclntableWithinGopclntab {
		t.Fatalf("ReadMetadata() stripped pcln bridges = %#v", got)
	}
}

func TestReadMetadataPEFixtureSectionHeuristic(t *testing.T) {
	t.Parallel()

	got, err := ReadMetadata("../../corpus/fixtures/go-pe-buildinfo-windows-amd64/fixture.exe")
	if err != nil {
		t.Fatalf("ReadMetadata() error = %v", err)
	}

	if got.TrustSummary != "section_heuristic" {
		t.Fatalf("ReadMetadata() trust summary = %q, want %q", got.TrustSummary, "section_heuristic")
	}
	if got.PETextSectionAddr != 0x140001000 || got.PETextSectionSize == 0 {
		t.Fatalf("ReadMetadata() PE text section = %#v", got)
	}
	if got.PERdataSectionAddr != 0x1400a5000 || got.PERdataSectionSize == 0 {
		t.Fatalf("ReadMetadata() PE rdata section = %#v", got)
	}
	if got.PEPclntabMagicSection != ".rdata" {
		t.Fatalf("ReadMetadata() PE pclntab magic section = %q", got.PEPclntabMagicSection)
	}
	if got.PEPclntabMagicAddr != 0x1400d6da8 || got.PEPclntabMagicCount != 22 {
		t.Fatalf("ReadMetadata() PE pclntab magic = %#v", got)
	}
	if got.PEPclntabHeaderSection != ".rdata" {
		t.Fatalf("ReadMetadata() PE pclntab header section = %q", got.PEPclntabHeaderSection)
	}
	if got.PEPclntabHeaderAddr != 0x1400dd568 {
		t.Fatalf("ReadMetadata() PE pclntab header addr = %#x", got.PEPclntabHeaderAddr)
	}
	if got.PEPclntabHeaderMagic != "f1ffffff" || got.PEPclntabHeaderQuantum != 1 || got.PEPclntabHeaderPointerSize != 8 {
		t.Fatalf("ReadMetadata() PE pclntab header = %#v", got)
	}
	if got.Provenance.Source != "core.runtime.pe" || got.Provenance.Confidence != "low" {
		t.Fatalf("ReadMetadata() provenance = %#v", got.Provenance)
	}
}

func TestParseELFPclntabHeader(t *testing.T) {
	t.Parallel()

	magic, kind, quantum, pointerSize, ok := parseELFPclntabHeader(
		[]byte{0xf1, 0xff, 0xff, 0xff, 0x00, 0x00, 0x01, 0x08},
	)
	if !ok {
		t.Fatal("parseELFPclntabHeader() unexpectedly failed for known header")
	}
	if magic != "f1ffffff" || kind != "known" || quantum != 1 || pointerSize != 8 {
		t.Fatalf(
			"parseELFPclntabHeader() = (%q, %q, %d, %d), want (%q, %q, %d, %d)",
			magic,
			kind,
			quantum,
			pointerSize,
			"f1ffffff",
			"known",
			1,
			8,
		)
	}

	magic, kind, quantum, pointerSize, ok = parseELFPclntabHeader(
		[]byte{0x12, 0x34, 0x56, 0x78, 0x00, 0x00, 0x01, 0x08},
	)
	if !ok {
		t.Fatal("parseELFPclntabHeader() unexpectedly failed for unknown header")
	}
	if magic != "12345678" || kind != "unknown" || quantum != 1 || pointerSize != 8 {
		t.Fatalf(
			"parseELFPclntabHeader() = (%q, %q, %d, %d), want (%q, %q, %d, %d)",
			magic,
			kind,
			quantum,
			pointerSize,
			"12345678",
			"unknown",
			1,
			8,
		)
	}
}

func TestParseELFPclntabHeaderHints(t *testing.T) {
	t.Parallel()

	header := make([]byte, 8+8*8)
	copy(header[:8], []byte{0x12, 0x34, 0x56, 0x78, 0x00, 0x00, 0x01, 0x08})
	binary.LittleEndian.PutUint64(header[8:], 123)
	binary.LittleEndian.PutUint64(header[16:], 45)
	binary.LittleEndian.PutUint64(header[32:], 0x111)
	binary.LittleEndian.PutUint64(header[40:], 0x222)
	binary.LittleEndian.PutUint64(header[48:], 0x333)
	binary.LittleEndian.PutUint64(header[56:], 0x444)
	binary.LittleEndian.PutUint64(header[64:], 0x555)

	got, ok := parseELFPclntabHeaderHints(binary.LittleEndian, header)
	if !ok {
		t.Fatal("parseELFPclntabHeaderHints() unexpectedly failed")
	}
	if got.functionCount != 123 || got.fileCount != 45 {
		t.Fatalf("parseELFPclntabHeaderHints() counts = %#v", got)
	}
	if got.funcnametabOffset != 0x111 || got.cuOffset != 0x222 || got.filetabOffset != 0x333 || got.pctabOffset != 0x444 || got.functabOffset != 0x555 {
		t.Fatalf("parseELFPclntabHeaderHints() offsets = %#v", got)
	}
}

func TestParseELFFunctabPCOffsetHints(t *testing.T) {
	t.Parallel()

	data := make([]byte, 96)
	const functabOffset = 32
	// 2 function entries plus trailing invalid PC sentinel.
	binary.LittleEndian.PutUint32(data[functabOffset:], 100)
	binary.LittleEndian.PutUint32(data[functabOffset+4:], 1000)
	binary.LittleEndian.PutUint32(data[functabOffset+8:], 200)
	binary.LittleEndian.PutUint32(data[functabOffset+12:], 2000)
	binary.LittleEndian.PutUint32(data[functabOffset+16:], 300)

	first, last, monotonic, ok := parseELFFunctabPCOffsetHints(
		binary.LittleEndian,
		data,
		functabOffset,
		2,
	)
	if !ok {
		t.Fatal("parseELFFunctabPCOffsetHints() unexpectedly failed")
	}
	if first != 100 || last != 300 || !monotonic {
		t.Fatalf(
			"parseELFFunctabPCOffsetHints() = (%d, %d, %t), want (%d, %d, %t)",
			first,
			last,
			monotonic,
			100,
			300,
			true,
		)
	}
}

func TestParseELFFunctabPCOffsetSample(t *testing.T) {
	t.Parallel()

	data := make([]byte, 96)
	const functabOffset = 32
	binary.LittleEndian.PutUint32(data[functabOffset:], 100)
	binary.LittleEndian.PutUint32(data[functabOffset+4:], 1000)
	binary.LittleEndian.PutUint32(data[functabOffset+8:], 200)
	binary.LittleEndian.PutUint32(data[functabOffset+12:], 2000)
	binary.LittleEndian.PutUint32(data[functabOffset+16:], 300)
	binary.LittleEndian.PutUint32(data[functabOffset+20:], 3000)

	got, ok := parseELFFunctabPCOffsetSample(binary.LittleEndian, data, functabOffset, 3, 2)
	if !ok {
		t.Fatal("parseELFFunctabPCOffsetSample() unexpectedly failed")
	}
	if len(got) != 2 || got[0] != 100 || got[1] != 200 {
		t.Fatalf("parseELFFunctabPCOffsetSample() = %#v", got)
	}
}

func TestPopulateELFFunctabAddressHints(t *testing.T) {
	t.Parallel()

	meta := schema.RuntimeMetadata{
		ELFFunctabFirstPCOffsetHint:  0,
		ELFFunctabLastPCOffsetHint:   0x300,
		ELFFunctabPCOffsetsMonotonic: true,
		ELFFunctabPCAddrSample:       []uint64{0, 0x100, 0x200},
		ModuledataTextAddr:           0x401000,
		ModuledataTextEndInclusive:   0x4015ff,
	}

	populateELFFunctabAddressHints(&meta)

	if meta.ELFFunctabFirstPCAddrHint != 0x401000 || meta.ELFFunctabLastPCAddrHint != 0x401300 || !meta.ELFFunctabPCAddrHintsWithinText {
		t.Fatalf("populateELFFunctabAddressHints() = %#v", meta)
	}
	if len(meta.ELFFunctabPCAddrSample) != 3 || meta.ELFFunctabPCAddrSample[1] != 0x401100 || !meta.ELFFunctabPCAddrSampleAllWithinText {
		t.Fatalf("populateELFFunctabAddressHints() sample = %#v", meta)
	}
}

func TestPopulateELFFunctabAddressHintsFallsBackToELFTextSection(t *testing.T) {
	t.Parallel()

	meta := schema.RuntimeMetadata{
		ELFFunctabFirstPCOffsetHint:  0,
		ELFFunctabLastPCOffsetHint:   0x300,
		ELFFunctabPCOffsetsMonotonic: true,
		ELFFunctabPCAddrSample:       []uint64{0, 0x100, 0x200},
		ELFTextSectionAddr:           0x501000,
		ELFTextSectionEndInclusive:   0x5015ff,
	}

	populateELFFunctabAddressHints(&meta)

	if meta.ELFFunctabFirstPCAddrHint != 0x501000 || meta.ELFFunctabLastPCAddrHint != 0x501300 || !meta.ELFFunctabPCAddrHintsWithinText {
		t.Fatalf("populateELFFunctabAddressHints() ELF text fallback = %#v", meta)
	}
	if len(meta.ELFFunctabPCAddrSample) != 3 || meta.ELFFunctabPCAddrSample[1] != 0x501100 || !meta.ELFFunctabPCAddrSampleAllWithinText {
		t.Fatalf("populateELFFunctabAddressHints() ELF text fallback sample = %#v", meta)
	}
}

func TestSummarizeTrust(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		meta struct {
			firstModuleDataAddr                 uint64
			firstModuleDataFromGoModuleFallback bool
			gopclntabAddr                       uint64
			peTextSectionAddr                   uint64
		}
		want string
	}{
		{
			name: "symbol backed",
			meta: struct {
				firstModuleDataAddr                 uint64
				firstModuleDataFromGoModuleFallback bool
				gopclntabAddr                       uint64
				peTextSectionAddr                   uint64
			}{
				firstModuleDataAddr: 0x401000,
			},
			want: "symbol_backed",
		},
		{
			name: "go module fallback",
			meta: struct {
				firstModuleDataAddr                 uint64
				firstModuleDataFromGoModuleFallback bool
				gopclntabAddr                       uint64
				peTextSectionAddr                   uint64
			}{
				firstModuleDataAddr:                 0x401000,
				firstModuleDataFromGoModuleFallback: true,
			},
			want: "go_module_fallback",
		},
		{
			name: "section heuristic",
			meta: struct {
				firstModuleDataAddr                 uint64
				firstModuleDataFromGoModuleFallback bool
				gopclntabAddr                       uint64
				peTextSectionAddr                   uint64
			}{
				gopclntabAddr: 0x500000,
			},
			want: "section_heuristic",
		},
		{
			name: "pe section heuristic",
			meta: struct {
				firstModuleDataAddr                 uint64
				firstModuleDataFromGoModuleFallback bool
				gopclntabAddr                       uint64
				peTextSectionAddr                   uint64
			}{
				peTextSectionAddr: 0x140001000,
			},
			want: "section_heuristic",
		},
		{
			name: "absent",
			want: "absent",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := summarizeTrust(schema.RuntimeMetadata{
				FirstModuleDataAddr:                 tt.meta.firstModuleDataAddr,
				FirstModuleDataFromGoModuleFallback: tt.meta.firstModuleDataFromGoModuleFallback,
				GopclntabAddr:                       tt.meta.gopclntabAddr,
				PETextSectionAddr:                   tt.meta.peTextSectionAddr,
			})

			if string(got) != tt.want {
				t.Fatalf("summarizeTrust() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReadTypelinkSample(t *testing.T) {
	t.Parallel()

	got := readTypelinkSample(binary.LittleEndian, []byte{
		0x04, 0x00, 0x00, 0x00,
		0xf8, 0xff, 0xff, 0xff,
		0x10, 0x00, 0x00, 0x00,
	}, 2)

	if len(got) != 2 || got[0] != 4 || got[1] != -8 {
		t.Fatalf("readTypelinkSample() = %#v", got)
	}
}

func TestSummarizeTypelinks(t *testing.T) {
	t.Parallel()

	minOffset, maxOffset := summarizeTypelinks(binary.LittleEndian, []byte{
		0x04, 0x00, 0x00, 0x00,
		0xf8, 0xff, 0xff, 0xff,
		0x10, 0x00, 0x00, 0x00,
	})

	if minOffset != -8 || maxOffset != 16 {
		t.Fatalf("summarizeTypelinks() = (%d, %d)", minOffset, maxOffset)
	}
}

func TestCountTypelinkSigns(t *testing.T) {
	t.Parallel()

	negativeCount, nonNegativeCount := countTypelinkSigns(binary.LittleEndian, []byte{
		0x04, 0x00, 0x00, 0x00,
		0xf8, 0xff, 0xff, 0xff,
		0x10, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	})

	if negativeCount != 1 || nonNegativeCount != 3 {
		t.Fatalf("countTypelinkSigns() = (%d, %d)", negativeCount, nonNegativeCount)
	}
}

func TestInAddressRange(t *testing.T) {
	t.Parallel()

	ok, offset := inAddressRange(0x110, 0x100, 0x40)
	if !ok || offset != 0x10 {
		t.Fatalf("inAddressRange() = (%v, %#x)", ok, offset)
	}

	ok, offset = inAddressRange(0x200, 0x100, 0x40)
	if ok || offset != 0 {
		t.Fatalf("inAddressRange() out-of-range = (%v, %#x)", ok, offset)
	}
}

func TestReadWordSample(t *testing.T) {
	t.Parallel()

	got64 := readWordSample(binary.LittleEndian, []byte{
		0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}, 8, 4)
	if len(got64) != 2 || got64[0] != 8 || got64[1] != 16 {
		t.Fatalf("readWordSample(64) = %#v", got64)
	}

	got32 := readWordSample(binary.LittleEndian, []byte{
		0x04, 0x00, 0x00, 0x00,
		0x08, 0x00, 0x00, 0x00,
	}, 4, 1)
	if len(got32) != 1 || got32[0] != 4 {
		t.Fatalf("readWordSample(32) = %#v", got32)
	}
}

func TestFindSliceHeader(t *testing.T) {
	t.Parallel()

	index, length, capacity, ok := findSliceHeader([]uint64{
		0x10,
		0x20,
		0x576a80,
		0x202,
		0x202,
	}, 0x576a80, 0x202)

	if !ok || index != 2 || length != 0x202 || capacity != 0x202 {
		t.Fatalf("findSliceHeader() = (%d, %d, %d, %v)", index, length, capacity, ok)
	}
}

func TestRangeWithinSizedRegion(t *testing.T) {
	t.Parallel()

	if !rangeWithinSizedRegion(0x1050, 0x20, 0x1000, 0x100) {
		t.Fatal("rangeWithinSizedRegion() expected true")
	}
	if rangeWithinSizedRegion(0x1100, 0x1, 0x1000, 0x100) {
		t.Fatal("rangeWithinSizedRegion() addr-at-end expected false")
	}
	if rangeWithinSizedRegion(0x10f0, 0x20, 0x1000, 0x100) {
		t.Fatal("rangeWithinSizedRegion() overflow expected false")
	}
}

func TestFindRangeBlock(t *testing.T) {
	t.Parallel()

	index, ok := findRangeBlock(
		[]uint64{
			0x10,
			0x20,
			0x578400,
			0x57d0e2,
			0x57d100,
			0x5827f2,
			0x582800,
			0x5a2cd8,
			0x5a2ce0,
			0x5b82a0,
		},
		[]uint64{
			0x578400,
			0x57d0e2,
			0x57d100,
			0x5827f2,
			0x582800,
			0x5a2cd8,
			0x5a2ce0,
			0x5b82a0,
		},
	)

	if !ok || index != 2 {
		t.Fatalf("findRangeBlock() = (%d, %v)", index, ok)
	}
}

func TestFindRangePair(t *testing.T) {
	t.Parallel()

	index, ok := findRangePair(
		[]uint64{
			0x4cf0a0,
			0x4cf144,
			0x49f000,
			0x4d5382,
		},
		0x49f000,
		0x4d5382,
	)

	if !ok || index != 2 {
		t.Fatalf("findRangePair() = (%d, %v)", index, ok)
	}
}

func TestInclusiveSectionEnd(t *testing.T) {
	t.Parallel()

	if got := inclusiveSectionEnd(0x401000, 0x9d0f1); got != 0x49e0f0 {
		t.Fatalf("inclusiveSectionEnd() = %#x", got)
	}
}

func TestResolveTypelinkSample(t *testing.T) {
	t.Parallel()

	got := resolveTypelinkSample(0x49f000, []int32{0x9d20, 0x8e20})
	if len(got) != 2 || got[0] != 0x4a8d20 || got[1] != 0x4a7e20 {
		t.Fatalf("resolveTypelinkSample() = %#v", got)
	}
}

func TestCountResolvedTypelinksWithinRange(t *testing.T) {
	t.Parallel()

	got := countResolvedTypelinksWithinRange(
		binary.LittleEndian,
		[]byte{
			0x20, 0x9d, 0x00, 0x00,
			0x20, 0x8e, 0x00, 0x00,
			0x00, 0x00, 0x04, 0x00,
		},
		0x49f000,
		0x49f000,
		0x4d5382,
	)

	if got != 2 {
		t.Fatalf("countResolvedTypelinksWithinRange() = %d", got)
	}
}

func TestResolveTypelinksWithinTypesRange(t *testing.T) {
	t.Parallel()

	sample := []int32{0x20, 0x30, 0x200}
	got := countResolvedTypelinksWithinRange(
		binary.LittleEndian,
		[]byte{
			0x20, 0x00, 0x00, 0x00,
			0x30, 0x00, 0x00, 0x00,
			0x00, 0x02, 0x00, 0x00,
		},
		0x490000,
		0x490000,
		0x490100,
	)
	if len(sample) != 3 || got != 2 {
		t.Fatalf("types-range count = %d", got)
	}
}
