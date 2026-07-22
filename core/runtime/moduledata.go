package runtime

import (
	"bytes"
	"debug/elf"
	stdpe "debug/pe"
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	binaryformat "github.com/dantte-lp/goreveal/core/format"
	"github.com/dantte-lp/goreveal/core/recoveryerr"
	"github.com/dantte-lp/goreveal/schema"
)

// ReadMetadata recovers bounded runtime layout evidence from supported binaries.
func ReadMetadata(path string) (schema.RuntimeMetadata, error) {
	kind, err := binaryformat.DetectFile(path)
	if err != nil {
		return schema.RuntimeMetadata{}, fmt.Errorf("detect runtime container: %w", err)
	}

	var meta schema.RuntimeMetadata
	switch kind {
	case binaryformat.ELF:
		meta, err = readELFMetadata(path)
	case binaryformat.PE:
		meta, err = readPEMetadata(path)
	case binaryformat.MachO:
		return schema.RuntimeMetadata{}, recoveryerr.NewUnsupported(
			recoveryerr.CodeRuntimeUnsupportedContainer,
			"runtime metadata recovery does not support Mach-O",
			nil,
		)
	case binaryformat.Unknown:
		return schema.RuntimeMetadata{}, errors.New("detect runtime container: unknown binary format")
	default:
		return schema.RuntimeMetadata{}, fmt.Errorf("detect runtime container: unhandled format %q", kind)
	}
	if err != nil {
		return schema.RuntimeMetadata{}, err
	}
	if meta.TrustSummary == schema.RuntimeTrustSummaryAbsent {
		return schema.RuntimeMetadata{}, recoveryerr.NewUnavailable(
			recoveryerr.CodeRuntimeMetadataNotFound,
			"runtime metadata evidence is absent",
			nil,
		)
	}

	return meta, nil
}

func readELFMetadata(path string) (schema.RuntimeMetadata, error) {
	fh, err := os.Open(path)
	if err != nil {
		return schema.RuntimeMetadata{}, fmt.Errorf("open file: %w", err)
	}
	defer fh.Close()

	ef, err := elf.NewFile(fh)
	if err != nil {
		return schema.RuntimeMetadata{}, fmt.Errorf("open ELF: %w", err)
	}

	meta := schema.RuntimeMetadata{
		Provenance: schema.Provenance{
			Source:     "core.runtime.elf",
			Confidence: "medium",
		},
	}

	populateFirstModuleData(ef, &meta)
	populatePclnTab(ef, &meta)
	populateELFTextSection(ef, &meta)
	populateTypelinks(ef, &meta)
	populateItablinks(ef, &meta)
	populateGoModule(ef, &meta)
	populateELFFunctabAddressHints(&meta)
	meta.TrustSummary = summarizeTrust(meta)

	return meta, nil
}

func readPEMetadata(path string) (schema.RuntimeMetadata, error) {
	fh, err := os.Open(path)
	if err != nil {
		return schema.RuntimeMetadata{}, fmt.Errorf("open file: %w", err)
	}
	defer fh.Close()

	pf, err := stdpe.NewFile(fh)
	if err != nil {
		return schema.RuntimeMetadata{}, fmt.Errorf("open PE: %w", err)
	}

	meta := schema.RuntimeMetadata{
		Provenance: schema.Provenance{
			Source:     "core.runtime.pe",
			Confidence: "low",
		},
	}

	populatePESectionRange(pf, ".text", &meta.PETextSectionAddr, &meta.PETextSectionSize)
	populatePESectionRange(pf, ".rdata", &meta.PERdataSectionAddr, &meta.PERdataSectionSize)
	populatePEPclntabMagic(pf, &meta)
	populatePEPclntabHeader(pf, &meta)
	meta.TrustSummary = summarizeTrust(meta)

	return meta, nil
}

func summarizeTrust(meta schema.RuntimeMetadata) schema.RuntimeTrustSummary {
	if meta.FirstModuleDataAddr != 0 {
		if meta.FirstModuleDataFromGoModuleFallback {
			return schema.RuntimeTrustSummaryGoModuleFallback
		}

		return schema.RuntimeTrustSummarySymbolBacked
	}

	if meta.GopclntabAddr != 0 ||
		meta.TypelinkAddr != 0 ||
		meta.ItablinkAddr != 0 ||
		meta.GoModuleAddr != 0 ||
		meta.PETextSectionAddr != 0 ||
		meta.PERdataSectionAddr != 0 ||
		meta.PEPclntabMagicAddr != 0 {
		return schema.RuntimeTrustSummarySectionHeuristic
	}

	return schema.RuntimeTrustSummaryAbsent
}

func populateFirstModuleData(ef *elf.File, meta *schema.RuntimeMetadata) {
	syms, err := ef.Symbols()
	if err == nil {
		for _, sym := range syms {
			if sym.Name == "runtime.firstmoduledata" {
				meta.FirstModuleDataAddr = sym.Value
				return
			}
		}
	}

	// On stripped fixtures the symbol may be absent while `.go.module` still starts
	// at the same runtime.firstmoduledata address for the current ELF layout family.
	section := ef.Section(".go.module")
	if section != nil {
		meta.FirstModuleDataAddr = section.Addr
		meta.FirstModuleDataFromGoModuleFallback = true
	}
}

func populatePESectionRange(pf *stdpe.File, name string, outAddr, outSize *uint64) {
	section := peSection(pf, name)
	if section == nil {
		return
	}

	*outAddr = peImageBase(pf) + uint64(section.VirtualAddress)
	*outSize = peSectionSize(section)
}

func populatePEPclntabMagic(pf *stdpe.File, meta *schema.RuntimeMetadata) {
	for _, sectionName := range []string{".rdata", ".text"} {
		section := peSection(pf, sectionName)
		if section == nil {
			continue
		}

		data, err := section.Data()
		if err != nil {
			continue
		}

		for _, magic := range pclntabMagics() {
			firstOffset, count, ok := findMagic(data, magic)
			if !ok {
				continue
			}

			meta.PEPclntabMagicSection = sectionName
			meta.PEPclntabMagicAddr = peImageBase(pf) + uint64(section.VirtualAddress) + uint64(firstOffset)
			meta.PEPclntabMagicCount = uint64(count)
			return
		}
	}
}

func populatePEPclntabHeader(pf *stdpe.File, meta *schema.RuntimeMetadata) {
	for _, sectionName := range []string{".rdata", ".text"} {
		section := peSection(pf, sectionName)
		if section == nil {
			continue
		}

		data, err := section.Data()
		if err != nil {
			continue
		}

		offset, magic, quantum, pointerSize, ok := findPEPclntabHeaderCandidate(data)
		if !ok {
			continue
		}

		meta.PEPclntabHeaderSection = sectionName
		meta.PEPclntabHeaderAddr = peImageBase(pf) + uint64(section.VirtualAddress) + uint64(offset)
		meta.PEPclntabHeaderMagic = magic
		meta.PEPclntabHeaderQuantum = uint64(quantum)
		meta.PEPclntabHeaderPointerSize = uint64(pointerSize)
		return
	}
}

func populatePclnTab(ef *elf.File, meta *schema.RuntimeMetadata) {
	section := ef.Section(".gopclntab")
	if section == nil {
		return
	}
	meta.GopclntabAddr = section.Addr
	meta.GopclntabSize = section.Size

	data, err := section.Data()
	if err != nil {
		return
	}

	magic, kind, quantum, pointerSize, ok := parseELFPclntabHeader(data)
	if !ok {
		return
	}

	meta.ELFPclntabHeaderMagic = magic
	meta.ELFPclntabHeaderMagicKind = kind
	meta.ELFPclntabHeaderQuantum = uint64(quantum)
	meta.ELFPclntabHeaderPointerSize = uint64(pointerSize)

	headerHints, ok := parseELFPclntabHeaderHints(ef.ByteOrder, data)
	if !ok {
		return
	}

	meta.ELFPclntabFunctionCountHint = headerHints.functionCount
	meta.ELFPclntabFileCountHint = headerHints.fileCount
	meta.ELFPclntabFuncnametabOffsetHint = headerHints.funcnametabOffset
	meta.ELFPclntabCuOffsetHint = headerHints.cuOffset
	meta.ELFPclntabFiletabOffsetHint = headerHints.filetabOffset
	meta.ELFPclntabPctabOffsetHint = headerHints.pctabOffset
	meta.ELFPclntabFunctabOffsetHint = headerHints.functabOffset

	firstPC, lastPC, monotonic, ok := parseELFFunctabPCOffsetHints(
		ef.ByteOrder,
		data,
		headerHints.functabOffset,
		headerHints.functionCount,
	)
	if !ok {
		return
	}

	meta.ELFFunctabFirstPCOffsetHint = firstPC
	meta.ELFFunctabLastPCOffsetHint = lastPC
	meta.ELFFunctabPCOffsetsMonotonic = monotonic

	pcOffsetSample, ok := parseELFFunctabPCOffsetSample(
		ef.ByteOrder,
		data,
		headerHints.functabOffset,
		headerHints.functionCount,
		8,
	)
	if !ok {
		return
	}

	meta.ELFFunctabPCAddrSample = pcOffsetSample
}

func populateELFTextSection(ef *elf.File, meta *schema.RuntimeMetadata) {
	section := ef.Section(".text")
	if section == nil {
		return
	}

	meta.ELFTextSectionAddr = section.Addr
	meta.ELFTextSectionEndInclusive = inclusiveSectionEnd(section.Addr, section.Size)
}

func populateTypelinks(ef *elf.File, meta *schema.RuntimeMetadata) {
	section := ef.Section(".typelink")
	if section == nil {
		return
	}

	meta.TypelinkAddr = section.Addr
	meta.TypelinkSize = section.Size
	meta.TypelinkCount = section.Size / 4

	data, err := section.Data()
	if err != nil {
		return
	}

	meta.TypelinkSample = readTypelinkSample(ef.ByteOrder, data, 8)
	meta.TypelinkMinOffset, meta.TypelinkMaxOffset = summarizeTypelinks(ef.ByteOrder, data)
	meta.TypelinkNegativeCount, meta.TypelinkNonNegativeCount = countTypelinkSigns(ef.ByteOrder, data)

	rodata := ef.Section(".rodata")
	if rodata == nil {
		return
	}

	meta.TypelinkResolvedBaseAddr = rodata.Addr
	meta.TypelinkResolvedSample = resolveTypelinkSample(rodata.Addr, meta.TypelinkSample)
	meta.TypelinkResolvedWithinRodataCount = countResolvedTypelinksWithinRange(
		ef.ByteOrder,
		data,
		rodata.Addr,
		rodata.Addr,
		rodata.Addr+rodata.Size,
	)
	meta.TypelinkAllResolvedWithinRodata = meta.TypelinkResolvedWithinRodataCount == meta.TypelinkCount
}

func populateItablinks(ef *elf.File, meta *schema.RuntimeMetadata) {
	section := ef.Section(".itablink")
	if section == nil {
		return
	}

	meta.ItablinkAddr = section.Addr
	meta.ItablinkSize = section.Size
	wordSize := uint64(elfWordSize(ef.Class))
	if wordSize > 0 {
		meta.ItablinkCount = section.Size / wordSize
	}
}

func populateGoModule(ef *elf.File, meta *schema.RuntimeMetadata) {
	section := ef.Section(".go.module")
	if section == nil {
		return
	}

	meta.GoModuleAddr = section.Addr
	meta.GoModuleSize = section.Size
	meta.FirstModuleDataInGoModule, meta.FirstModuleDataGoModuleOffset = inAddressRange(meta.FirstModuleDataAddr, section.Addr, section.Size)

	data, err := section.Data()
	if err != nil {
		return
	}

	wordSize := elfWordSize(ef.Class)
	meta.GoModuleWordSize = uint64(wordSize)
	meta.GoModuleWordSample = readWordSample(ef.ByteOrder, data, wordSize, 4)
	if !meta.FirstModuleDataInGoModule {
		return
	}

	words := readWords(ef.ByteOrder, data[meta.FirstModuleDataGoModuleOffset:], wordSize)
	populateGoModuleSliceHeaders(words, meta)
	populateGoModuleMemoryRanges(ef, words, meta)
	populateGoModuleRodata(ef, words, meta)
	populateGoModuleText(ef, words, meta)
}

func populateGoModuleSliceHeaders(words []uint64, meta *schema.RuntimeMetadata) {
	if len(words) >= 1 {
		meta.ModuledataPCHeaderAddr = words[0]
		meta.ModuledataPCHeaderMatchesGopclntab = words[0] == meta.GopclntabAddr
	}
	assignPclnBridge(words, 1, meta.GopclntabAddr, meta.GopclntabSize, &meta.ModuledataFuncnametabSliceWordIndex, &meta.ModuledataFuncnametabAddr, &meta.ModuledataFuncnametabLen, &meta.ModuledataFuncnametabCap, &meta.ModuledataFuncnametabWithinGopclntab)
	assignPclnBridge(words, 4, meta.GopclntabAddr, meta.GopclntabSize, &meta.ModuledataCutabSliceWordIndex, &meta.ModuledataCutabAddr, &meta.ModuledataCutabLen, &meta.ModuledataCutabCap, &meta.ModuledataCutabWithinGopclntab)
	assignPclnBridge(words, 7, meta.GopclntabAddr, meta.GopclntabSize, &meta.ModuledataFiletabSliceWordIndex, &meta.ModuledataFiletabAddr, &meta.ModuledataFiletabLen, &meta.ModuledataFiletabCap, &meta.ModuledataFiletabWithinGopclntab)
	assignPclnBridge(words, 10, meta.GopclntabAddr, meta.GopclntabSize, &meta.ModuledataPctabSliceWordIndex, &meta.ModuledataPctabAddr, &meta.ModuledataPctabLen, &meta.ModuledataPctabCap, &meta.ModuledataPctabWithinGopclntab)
	assignPclnBridge(words, 13, meta.GopclntabAddr, meta.GopclntabSize, &meta.ModuledataPclntableSliceWordIndex, &meta.ModuledataPclntableAddr, &meta.ModuledataPclntableLen, &meta.ModuledataPclntableCap, &meta.ModuledataPclntableWithinGopclntab)
	if index, length, capacity, ok := findSliceHeader(words, meta.TypelinkAddr, meta.TypelinkCount); ok {
		meta.ModuledataTypelinkSliceWordIndex = index
		meta.ModuledataTypelinkLen = length
		meta.ModuledataTypelinkCap = capacity
	}
	if index, length, capacity, ok := findSliceHeader(words, meta.ItablinkAddr, meta.ItablinkCount); ok {
		meta.ModuledataItablinkSliceWordIndex = index
		meta.ModuledataItablinkLen = length
		meta.ModuledataItablinkCap = capacity
	}
}

func assignPclnBridge(words []uint64, index int, base, size uint64, outIndex, outAddr, outLen, outCap *uint64, outWithin *bool) {
	if len(words) < index+3 {
		return
	}

	*outIndex = uint64(index)
	*outAddr = words[index]
	*outLen = words[index+1]
	*outCap = words[index+2]
	*outWithin = rangeWithinSizedRegion(*outAddr, *outLen, base, size)
}

func populateGoModuleMemoryRanges(ef *elf.File, words []uint64, meta *schema.RuntimeMetadata) {
	noptrdata, dataRange, bss, noptrbss, index, ok := findMemoryRangeBlock(ef, words)
	if !ok {
		return
	}

	meta.ModuledataMemoryRangesWordIndex = index
	meta.ModuledataNoptrdataAddr = noptrdata[0]
	meta.ModuledataNoptrdataEnd = noptrdata[1]
	meta.ModuledataDataAddr = dataRange[0]
	meta.ModuledataDataEnd = dataRange[1]
	meta.ModuledataBssAddr = bss[0]
	meta.ModuledataBssEnd = bss[1]
	meta.ModuledataNoptrbssAddr = noptrbss[0]
	meta.ModuledataNoptrbssEnd = noptrbss[1]
}

func populateGoModuleRodata(ef *elf.File, words []uint64, meta *schema.RuntimeMetadata) {
	section := ef.Section(".rodata")
	if section == nil {
		return
	}

	if index, ok := findRangePair(words, section.Addr, section.Addr+section.Size); ok {
		meta.ModuledataRodataWordIndex = index
		meta.ModuledataRodataAddr = section.Addr
		meta.ModuledataRodataEnd = section.Addr + section.Size
	}

	index, ok := findRangePair(words, section.Addr, section.Addr+section.Size)
	if !ok {
		return
	}

	meta.ModuledataTypesRangeWordIndex = index
	meta.ModuledataTypesAddr = section.Addr
	meta.ModuledataETypesAddr = section.Addr + section.Size

	typelink := ef.Section(".typelink")
	if typelink == nil {
		return
	}
	tlData, err := typelink.Data()
	if err != nil {
		return
	}

	meta.TypelinkResolvedWithinTypesCount = countResolvedTypelinksWithinRange(
		ef.ByteOrder,
		tlData,
		meta.ModuledataTypesAddr,
		meta.ModuledataTypesAddr,
		meta.ModuledataETypesAddr,
	)
	meta.TypelinkAllResolvedWithinTypes = meta.TypelinkResolvedWithinTypesCount == meta.TypelinkCount
}

func populateGoModuleText(ef *elf.File, words []uint64, meta *schema.RuntimeMetadata) {
	section := ef.Section(".text")
	if section == nil {
		return
	}
	if index, ok := findRangePair(words, section.Addr, inclusiveSectionEnd(section.Addr, section.Size)); ok {
		meta.ModuledataTextWordIndex = index
		meta.ModuledataTextAddr = section.Addr
		meta.ModuledataTextEndInclusive = inclusiveSectionEnd(section.Addr, section.Size)
	}
}

func populateELFFunctabAddressHints(meta *schema.RuntimeMetadata) {
	if meta == nil {
		return
	}
	if !meta.ELFFunctabPCOffsetsMonotonic {
		return
	}

	textAddr, textEndInclusive := elfTextRange(meta)
	if textAddr == 0 || textEndInclusive == 0 {
		return
	}

	meta.ELFFunctabFirstPCAddrHint = textAddr + meta.ELFFunctabFirstPCOffsetHint
	meta.ELFFunctabLastPCAddrHint = textAddr + meta.ELFFunctabLastPCOffsetHint

	textEndExclusive, ok := inclusiveEndToExclusive(textEndInclusive)
	if !ok {
		return
	}
	if meta.ELFFunctabFirstPCAddrHint < textAddr {
		return
	}
	if meta.ELFFunctabLastPCAddrHint < meta.ELFFunctabFirstPCAddrHint {
		return
	}

	meta.ELFFunctabPCAddrHintsWithinText = meta.ELFFunctabLastPCAddrHint <= textEndExclusive
	if len(meta.ELFFunctabPCAddrSample) == 0 {
		return
	}

	sample := make([]uint64, 0, len(meta.ELFFunctabPCAddrSample))
	allWithinText := true
	for _, offset := range meta.ELFFunctabPCAddrSample {
		addr := textAddr + offset
		sample = append(sample, addr)
		if addr < textAddr || addr > textEndExclusive {
			allWithinText = false
		}
	}
	meta.ELFFunctabPCAddrSample = sample
	meta.ELFFunctabPCAddrSampleAllWithinText = allWithinText
}

func elfTextRange(meta *schema.RuntimeMetadata) (uint64, uint64) {
	if meta == nil {
		return 0, 0
	}
	if meta.ModuledataTextAddr != 0 && meta.ModuledataTextEndInclusive != 0 {
		return meta.ModuledataTextAddr, meta.ModuledataTextEndInclusive
	}
	return meta.ELFTextSectionAddr, meta.ELFTextSectionEndInclusive
}

func elfTextSource(meta *schema.RuntimeMetadata) string {
	if meta == nil {
		return ""
	}
	if meta.ModuledataTextAddr != 0 && meta.ModuledataTextEndInclusive != 0 {
		return "moduledata_text"
	}
	if meta.ELFTextSectionAddr != 0 && meta.ELFTextSectionEndInclusive != 0 {
		return "elf_text_section"
	}
	return ""
}

// ELFTextSourceForProjection exposes the currently truthful ELF text-range source
// for analyst-facing projection layers without exporting internal helpers wholesale.
func ELFTextSourceForProjection(meta *schema.RuntimeMetadata) string {
	return elfTextSource(meta)
}

func inclusiveEndToExclusive(endInclusive uint64) (uint64, bool) {
	if endInclusive == ^uint64(0) {
		return 0, false
	}
	return endInclusive + 1, true
}

func readTypelinkSample(order binary.ByteOrder, data []byte, limit int) []int32 {
	if len(data) < 4 || limit <= 0 {
		return nil
	}

	count := len(data) / 4
	count = min(count, limit)

	sample := make([]int32, 0, count)
	for i := 0; i < count; i++ {
		offset := int32(order.Uint32(data[i*4 : i*4+4]))
		sample = append(sample, offset)
	}

	return sample
}

func summarizeTypelinks(order binary.ByteOrder, data []byte) (int32, int32) {
	if len(data) < 4 {
		return 0, 0
	}

	first := int32(order.Uint32(data[:4]))
	minOffset := first
	maxOffset := first

	for i := 1; i < len(data)/4; i++ {
		offset := int32(order.Uint32(data[i*4 : i*4+4]))
		if offset < minOffset {
			minOffset = offset
		}
		if offset > maxOffset {
			maxOffset = offset
		}
	}

	return minOffset, maxOffset
}

func countTypelinkSigns(order binary.ByteOrder, data []byte) (uint64, uint64) {
	var negativeCount uint64
	var nonNegativeCount uint64
	for i := range len(data) / 4 {
		offset := int32(order.Uint32(data[i*4 : i*4+4]))
		if offset < 0 {
			negativeCount++
			continue
		}
		nonNegativeCount++
	}

	return negativeCount, nonNegativeCount
}

func resolveTypelinkSample(base uint64, sample []int32) []uint64 {
	if base == 0 || len(sample) == 0 {
		return nil
	}

	resolved := make([]uint64, 0, len(sample))
	for _, offset := range sample {
		if offset < 0 {
			continue
		}
		resolved = append(resolved, base+uint64(offset))
	}

	return resolved
}

func countResolvedTypelinksWithinRange(order binary.ByteOrder, data []byte, base, start, end uint64) uint64 {
	var count uint64
	for i := range len(data) / 4 {
		offset := int32(order.Uint32(data[i*4 : i*4+4]))
		if offset < 0 {
			continue
		}
		va := base + uint64(offset)
		if va >= start && va < end {
			count++
		}
	}
	return count
}

func inAddressRange(addr, base, size uint64) (bool, uint64) {
	if addr < base || addr >= base+size {
		return false, 0
	}
	return true, addr - base
}

func rangeWithinSizedRegion(addr, length, base, size uint64) bool {
	if size == 0 {
		return false
	}
	if addr < base {
		return false
	}
	end := base + size
	if addr >= end {
		return false
	}
	return length <= end-addr
}

func elfWordSize(class elf.Class) int {
	if class == elf.ELFCLASS32 {
		return 4
	}
	return 8
}

func readWordSample(order binary.ByteOrder, data []byte, wordSize, limit int) []uint64 {
	if wordSize <= 0 || len(data) < wordSize || limit <= 0 {
		return nil
	}

	count := len(data) / wordSize
	count = min(count, limit)

	sample := make([]uint64, 0, count)
	for i := 0; i < count; i++ {
		start := i * wordSize
		end := start + wordSize
		switch wordSize {
		case 4:
			sample = append(sample, uint64(order.Uint32(data[start:end])))
		case 8:
			sample = append(sample, order.Uint64(data[start:end]))
		default:
			return sample
		}
	}

	return sample
}

func readWords(order binary.ByteOrder, data []byte, wordSize int) []uint64 {
	if wordSize <= 0 || len(data) < wordSize {
		return nil
	}

	count := len(data) / wordSize
	words := make([]uint64, 0, count)
	for i := range count {
		start := i * wordSize
		end := start + wordSize
		switch wordSize {
		case 4:
			words = append(words, uint64(order.Uint32(data[start:end])))
		case 8:
			words = append(words, order.Uint64(data[start:end]))
		default:
			return words
		}
	}

	return words
}

func findSliceHeader(words []uint64, dataAddr, dataCount uint64) (uint64, uint64, uint64, bool) {
	for i := 0; i+2 < len(words); i++ {
		if words[i] != dataAddr {
			continue
		}
		if words[i+1] != dataCount || words[i+2] != dataCount {
			continue
		}
		return uint64(i), words[i+1], words[i+2], true
	}

	return 0, 0, 0, false
}

func findRangeBlock(words []uint64, block []uint64) (uint64, bool) {
	if len(block) == 0 || len(words) < len(block) {
		return 0, false
	}

outer:
	for i := 0; i+len(block) <= len(words); i++ {
		for j := range len(block) {
			if words[i+j] != block[j] {
				continue outer
			}
		}
		return uint64(i), true
	}

	return 0, false
}

func findRangePair(words []uint64, start, end uint64) (uint64, bool) {
	return findRangeBlock(words, []uint64{start, end})
}

func inclusiveSectionEnd(addr, size uint64) uint64 {
	if size == 0 {
		return addr
	}
	return addr + size - 1
}

func findMemoryRangeBlock(ef *elf.File, words []uint64) ([2]uint64, [2]uint64, [2]uint64, [2]uint64, uint64, bool) {
	names := []string{".noptrdata", ".data", ".bss", ".noptrbss"}
	block := make([]uint64, 0, len(names)*2)
	ranges := make([][2]uint64, 0, len(names))
	for _, name := range names {
		section := ef.Section(name)
		if section == nil {
			return [2]uint64{}, [2]uint64{}, [2]uint64{}, [2]uint64{}, 0, false
		}
		r := [2]uint64{section.Addr, section.Addr + section.Size}
		ranges = append(ranges, r)
		block = append(block, r[0], r[1])
	}

	index, ok := findRangeBlock(words, block)
	if !ok {
		return [2]uint64{}, [2]uint64{}, [2]uint64{}, [2]uint64{}, 0, false
	}

	return ranges[0], ranges[1], ranges[2], ranges[3], index, true
}

func peSection(pf *stdpe.File, name string) *stdpe.Section {
	for _, section := range pf.Sections {
		if section.Name == name {
			return section
		}
	}

	return nil
}

func peImageBase(pf *stdpe.File) uint64 {
	switch header := pf.OptionalHeader.(type) {
	case *stdpe.OptionalHeader32:
		return uint64(header.ImageBase)
	case *stdpe.OptionalHeader64:
		return header.ImageBase
	default:
		return 0
	}
}

func peSectionSize(section *stdpe.Section) uint64 {
	if section.VirtualSize != 0 {
		return uint64(section.VirtualSize)
	}

	return uint64(section.Size)
}

func pclntabMagics() [][]byte {
	return [][]byte{
		{0xf0, 0xff, 0xff, 0xff},
		{0xf1, 0xff, 0xff, 0xff},
		{0xfa, 0xff, 0xff, 0xff},
		{0xfb, 0xff, 0xff, 0xff},
	}
}

func parseELFPclntabHeader(data []byte) (string, string, byte, byte, bool) {
	if len(data) < 8 {
		return "", "", 0, 0, false
	}

	magic := data[:4]
	quantum := data[6]
	pointerSize := data[7]
	if quantum != 1 && quantum != 2 && quantum != 4 {
		return "", "", 0, 0, false
	}
	if pointerSize != 4 && pointerSize != 8 {
		return "", "", 0, 0, false
	}

	kind := "unknown"
	for _, known := range pclntabMagics() {
		if bytes.Equal(magic, known) {
			kind = "known"
			break
		}
	}

	return fmt.Sprintf("%x", magic), kind, quantum, pointerSize, true
}

type elfPclntabHeaderHints struct {
	functionCount     uint64
	fileCount         uint64
	funcnametabOffset uint64
	cuOffset          uint64
	filetabOffset     uint64
	pctabOffset       uint64
	functabOffset     uint64
}

func parseELFPclntabHeaderHints(order binary.ByteOrder, data []byte) (elfPclntabHeaderHints, bool) {
	if order == nil {
		return elfPclntabHeaderHints{}, false
	}

	_, _, _, pointerSize, ok := parseELFPclntabHeader(data)
	if !ok {
		return elfPclntabHeaderHints{}, false
	}

	wordSize := int(pointerSize)
	requiredWords := 8
	headerSize := 8 + requiredWords*wordSize
	if len(data) < headerSize {
		return elfPclntabHeaderHints{}, false
	}

	word := func(index int) uint64 {
		start := 8 + index*wordSize
		end := start + wordSize
		switch wordSize {
		case 4:
			return uint64(order.Uint32(data[start:end]))
		case 8:
			return order.Uint64(data[start:end])
		default:
			return 0
		}
	}

	return elfPclntabHeaderHints{
		functionCount:     word(0),
		fileCount:         word(1),
		funcnametabOffset: word(3),
		cuOffset:          word(4),
		filetabOffset:     word(5),
		pctabOffset:       word(6),
		functabOffset:     word(7),
	}, true
}

func parseELFFunctabPCOffsetHints(
	order binary.ByteOrder,
	data []byte,
	functabOffset uint64,
	functionCount uint64,
) (uint64, uint64, bool, bool) {
	if order == nil || functionCount == 0 {
		return 0, 0, false, false
	}
	if functabOffset >= uint64(len(data)) {
		return 0, 0, false, false
	}

	const entryFieldSize = 4
	entrySize := entryFieldSize * 2
	required := functabOffset + uint64((functionCount*2+1)*entryFieldSize)
	if required > uint64(len(data)) {
		return 0, 0, false, false
	}

	readPCOffset := func(index uint64) uint64 {
		start := int(functabOffset + index*uint64(entrySize))
		return uint64(order.Uint32(data[start : start+entryFieldSize]))
	}

	first := readPCOffset(0)
	last := readPCOffset(functionCount)
	monotonic := true
	prev := first
	for i := uint64(1); i <= functionCount; i++ {
		current := readPCOffset(i)
		if current < prev {
			monotonic = false
			break
		}
		prev = current
	}

	return first, last, monotonic, true
}

func parseELFFunctabPCOffsetSample(
	order binary.ByteOrder,
	data []byte,
	functabOffset uint64,
	functionCount uint64,
	limit int,
) ([]uint64, bool) {
	if order == nil || functionCount == 0 || limit <= 0 {
		return nil, false
	}
	if functabOffset >= uint64(len(data)) {
		return nil, false
	}

	const entryFieldSize = 4
	entrySize := entryFieldSize * 2
	maxEntries := functionCount
	if uint64(limit) < maxEntries {
		maxEntries = uint64(limit)
	}
	required := functabOffset + uint64((maxEntries-1)*uint64(entrySize)) + entryFieldSize
	if required > uint64(len(data)) {
		return nil, false
	}

	sample := make([]uint64, 0, maxEntries)
	for i := uint64(0); i < maxEntries; i++ {
		start := int(functabOffset + i*uint64(entrySize))
		sample = append(sample, uint64(order.Uint32(data[start:start+entryFieldSize])))
	}

	return sample, true
}

func findMagic(data, magic []byte) (int, int, bool) {
	if len(data) < len(magic) || len(magic) == 0 {
		return 0, 0, false
	}

	first := -1
	count := 0
	start := 0
	for {
		index := bytes.Index(data[start:], magic)
		if index < 0 {
			break
		}

		absolute := start + index
		if first < 0 {
			first = absolute
		}
		count++
		start = absolute + 1
	}

	return first, count, first >= 0
}

func findPEPclntabHeaderCandidate(data []byte) (int, string, byte, byte, bool) {
	for _, magic := range pclntabMagics() {
		start := 0
		for {
			index := bytes.Index(data[start:], magic)
			if index < 0 {
				break
			}

			absolute := start + index
			quantum, pointerSize, ok := parsePEPclntabHeaderCandidate(data, absolute)
			if ok {
				return absolute, fmt.Sprintf("%x", magic), quantum, pointerSize, true
			}

			start = absolute + 1
		}
	}

	return 0, "", 0, 0, false
}

func parsePEPclntabHeaderCandidate(data []byte, offset int) (byte, byte, bool) {
	if len(data) < offset+8 {
		return 0, 0, false
	}
	if data[offset+4] != 0 || data[offset+5] != 0 {
		return 0, 0, false
	}

	quantum := data[offset+6]
	if quantum != 1 && quantum != 2 && quantum != 4 {
		return 0, 0, false
	}

	pointerSize := data[offset+7]
	if pointerSize != 4 && pointerSize != 8 {
		return 0, 0, false
	}

	return quantum, pointerSize, true
}
