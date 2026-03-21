package runtime

import (
	"debug/elf"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/dantte-lp/goreveal/schema"
)

// ReadMetadata recovers low-risk runtime layout evidence from ELF binaries.
func ReadMetadata(path string) (schema.RuntimeMetadata, error) {
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
	populateTypelinks(ef, &meta)
	populateItablinks(ef, &meta)
	populateGoModule(ef, &meta)

	return meta, nil
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

func populatePclnTab(ef *elf.File, meta *schema.RuntimeMetadata) {
	section := ef.Section(".gopclntab")
	if section == nil {
		return
	}
	meta.GopclntabAddr = section.Addr
	meta.GopclntabSize = section.Size
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
