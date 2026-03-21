package runtime

import (
	"encoding/binary"
	"testing"
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
	if got.FirstModuleDataFromGoModuleFallback {
		t.Fatalf("ReadMetadata() rich fixture unexpectedly marked go.module fallback = %#v", got)
	}
	if got.GopclntabAddr == 0 || got.GopclntabSize == 0 {
		t.Fatalf("ReadMetadata() gopclntab = %#v", got)
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
