package schema

// Input describes the analyzed binary at the schema boundary.
type Input struct {
	Path   string `json:"path"`
	Size   int64  `json:"size"`
	Format string `json:"format"`
}

// BuildInfo captures Go build metadata recovered from the binary.
type BuildInfo struct {
	GoVersion  string     `json:"go_version"`
	Path       string     `json:"path"`
	Provenance Provenance `json:"provenance"`
}

// RuntimeMetadata captures low-risk runtime layout evidence recovered from the binary.
type RuntimeMetadata struct {
	FirstModuleDataAddr                  uint64     `json:"firstmoduledata_addr,omitempty"`
	FirstModuleDataFromGoModuleFallback  bool       `json:"firstmoduledata_from_go_module_fallback,omitempty"`
	GopclntabAddr                        uint64     `json:"gopclntab_addr,omitempty"`
	GopclntabSize                        uint64     `json:"gopclntab_size,omitempty"`
	TypelinkAddr                         uint64     `json:"typelink_addr,omitempty"`
	TypelinkSize                         uint64     `json:"typelink_size,omitempty"`
	TypelinkCount                        uint64     `json:"typelink_count,omitempty"`
	ItablinkAddr                         uint64     `json:"itablink_addr,omitempty"`
	ItablinkSize                         uint64     `json:"itablink_size,omitempty"`
	ItablinkCount                        uint64     `json:"itablink_count,omitempty"`
	TypelinkSample                       []int32    `json:"typelink_sample,omitempty"`
	TypelinkResolvedBaseAddr             uint64     `json:"typelink_resolved_base_addr,omitempty"`
	TypelinkResolvedSample               []uint64   `json:"typelink_resolved_sample,omitempty"`
	TypelinkResolvedWithinRodataCount    uint64     `json:"typelink_resolved_within_rodata_count,omitempty"`
	TypelinkAllResolvedWithinRodata      bool       `json:"typelink_all_resolved_within_rodata,omitempty"`
	TypelinkMinOffset                    int32      `json:"typelink_min_offset,omitempty"`
	TypelinkMaxOffset                    int32      `json:"typelink_max_offset,omitempty"`
	TypelinkNegativeCount                uint64     `json:"typelink_negative_count,omitempty"`
	TypelinkNonNegativeCount             uint64     `json:"typelink_non_negative_count,omitempty"`
	GoModuleAddr                         uint64     `json:"go_module_addr,omitempty"`
	GoModuleSize                         uint64     `json:"go_module_size,omitempty"`
	FirstModuleDataInGoModule            bool       `json:"firstmoduledata_in_go_module,omitempty"`
	FirstModuleDataGoModuleOffset        uint64     `json:"firstmoduledata_go_module_offset,omitempty"`
	GoModuleWordSize                     uint64     `json:"go_module_word_size,omitempty"`
	GoModuleWordSample                   []uint64   `json:"go_module_word_sample,omitempty"`
	ModuledataPCHeaderAddr               uint64     `json:"moduledata_pcheader_addr,omitempty"`
	ModuledataPCHeaderMatchesGopclntab   bool       `json:"moduledata_pcheader_matches_gopclntab,omitempty"`
	ModuledataFuncnametabSliceWordIndex  uint64     `json:"moduledata_funcnametab_slice_word_index,omitempty"`
	ModuledataFuncnametabAddr            uint64     `json:"moduledata_funcnametab_addr,omitempty"`
	ModuledataFuncnametabLen             uint64     `json:"moduledata_funcnametab_len,omitempty"`
	ModuledataFuncnametabCap             uint64     `json:"moduledata_funcnametab_cap,omitempty"`
	ModuledataFuncnametabWithinGopclntab bool       `json:"moduledata_funcnametab_within_gopclntab,omitempty"`
	ModuledataCutabSliceWordIndex        uint64     `json:"moduledata_cutab_slice_word_index,omitempty"`
	ModuledataCutabAddr                  uint64     `json:"moduledata_cutab_addr,omitempty"`
	ModuledataCutabLen                   uint64     `json:"moduledata_cutab_len,omitempty"`
	ModuledataCutabCap                   uint64     `json:"moduledata_cutab_cap,omitempty"`
	ModuledataCutabWithinGopclntab       bool       `json:"moduledata_cutab_within_gopclntab,omitempty"`
	ModuledataFiletabSliceWordIndex      uint64     `json:"moduledata_filetab_slice_word_index,omitempty"`
	ModuledataFiletabAddr                uint64     `json:"moduledata_filetab_addr,omitempty"`
	ModuledataFiletabLen                 uint64     `json:"moduledata_filetab_len,omitempty"`
	ModuledataFiletabCap                 uint64     `json:"moduledata_filetab_cap,omitempty"`
	ModuledataFiletabWithinGopclntab     bool       `json:"moduledata_filetab_within_gopclntab,omitempty"`
	ModuledataPctabSliceWordIndex        uint64     `json:"moduledata_pctab_slice_word_index,omitempty"`
	ModuledataPctabAddr                  uint64     `json:"moduledata_pctab_addr,omitempty"`
	ModuledataPctabLen                   uint64     `json:"moduledata_pctab_len,omitempty"`
	ModuledataPctabCap                   uint64     `json:"moduledata_pctab_cap,omitempty"`
	ModuledataPctabWithinGopclntab       bool       `json:"moduledata_pctab_within_gopclntab,omitempty"`
	ModuledataPclntableSliceWordIndex    uint64     `json:"moduledata_pclntable_slice_word_index,omitempty"`
	ModuledataPclntableAddr              uint64     `json:"moduledata_pclntable_addr,omitempty"`
	ModuledataPclntableLen               uint64     `json:"moduledata_pclntable_len,omitempty"`
	ModuledataPclntableCap               uint64     `json:"moduledata_pclntable_cap,omitempty"`
	ModuledataPclntableWithinGopclntab   bool       `json:"moduledata_pclntable_within_gopclntab,omitempty"`
	ModuledataTypelinkSliceWordIndex     uint64     `json:"moduledata_typelink_slice_word_index,omitempty"`
	ModuledataTypelinkLen                uint64     `json:"moduledata_typelink_len,omitempty"`
	ModuledataTypelinkCap                uint64     `json:"moduledata_typelink_cap,omitempty"`
	ModuledataItablinkSliceWordIndex     uint64     `json:"moduledata_itablink_slice_word_index,omitempty"`
	ModuledataItablinkLen                uint64     `json:"moduledata_itablink_len,omitempty"`
	ModuledataItablinkCap                uint64     `json:"moduledata_itablink_cap,omitempty"`
	ModuledataMemoryRangesWordIndex      uint64     `json:"moduledata_memory_ranges_word_index,omitempty"`
	ModuledataNoptrdataAddr              uint64     `json:"moduledata_noptrdata_addr,omitempty"`
	ModuledataNoptrdataEnd               uint64     `json:"moduledata_noptrdata_end,omitempty"`
	ModuledataDataAddr                   uint64     `json:"moduledata_data_addr,omitempty"`
	ModuledataDataEnd                    uint64     `json:"moduledata_data_end,omitempty"`
	ModuledataBssAddr                    uint64     `json:"moduledata_bss_addr,omitempty"`
	ModuledataBssEnd                     uint64     `json:"moduledata_bss_end,omitempty"`
	ModuledataNoptrbssAddr               uint64     `json:"moduledata_noptrbss_addr,omitempty"`
	ModuledataNoptrbssEnd                uint64     `json:"moduledata_noptrbss_end,omitempty"`
	ModuledataRodataWordIndex            uint64     `json:"moduledata_rodata_word_index,omitempty"`
	ModuledataRodataAddr                 uint64     `json:"moduledata_rodata_addr,omitempty"`
	ModuledataRodataEnd                  uint64     `json:"moduledata_rodata_end,omitempty"`
	ModuledataTextWordIndex              uint64     `json:"moduledata_text_word_index,omitempty"`
	ModuledataTextAddr                   uint64     `json:"moduledata_text_addr,omitempty"`
	ModuledataTextEndInclusive           uint64     `json:"moduledata_text_end_inclusive,omitempty"`
	ModuledataTypesRangeWordIndex        uint64     `json:"moduledata_types_range_word_index,omitempty"`
	ModuledataTypesAddr                  uint64     `json:"moduledata_types_addr,omitempty"`
	ModuledataETypesAddr                 uint64     `json:"moduledata_etypes_addr,omitempty"`
	TypelinkResolvedWithinTypesCount     uint64     `json:"typelink_resolved_within_types_count,omitempty"`
	TypelinkAllResolvedWithinTypes       bool       `json:"typelink_all_resolved_within_types,omitempty"`
	Provenance                           Provenance `json:"provenance"`
}

// Function captures a recovered function symbol and its address range.
type Function struct {
	Name          string     `json:"name"`
	Package       string     `json:"package,omitempty"`
	ImportPath    string     `json:"import_path,omitempty"`
	SourceFile    string     `json:"source_file,omitempty"`
	SourceLine    int        `json:"source_line,omitempty"`
	Autogenerated bool       `json:"autogenerated,omitempty"`
	Entry         uint64     `json:"entry"`
	End           uint64     `json:"end"`
	ModuleLocal   bool       `json:"module_local,omitempty"`
	Provenance    Provenance `json:"provenance"`
}

// Package captures a recovered package classification derived from functions.
type Package struct {
	Name              string     `json:"name"`
	ImportPath        string     `json:"import_path,omitempty"`
	SourceFileCount   int        `json:"source_file_count,omitempty"`
	FunctionCount     int        `json:"function_count"`
	HasSourceEvidence bool       `json:"has_source_evidence"`
	ModuleLocal       bool       `json:"module_local,omitempty"`
	Provenance        Provenance `json:"provenance"`
}

// Type captures a recovered type and the metadata source used to derive it.
type Type struct {
	Name            string     `json:"name"`
	Package         string     `json:"package,omitempty"`
	ImportPath      string     `json:"import_path,omitempty"`
	Kind            string     `json:"kind"`
	SourceFileCount int        `json:"source_file_count,omitempty"`
	ModuleLocal     bool       `json:"module_local,omitempty"`
	UserMeaningful  bool       `json:"user_meaningful,omitempty"`
	Provenance      Provenance `json:"provenance"`
}

// StringRegion describes a binary region scanned for string candidates.
type StringRegion struct {
	Name   string `json:"name"`
	Addr   uint64 `json:"addr"`
	Size   uint64 `json:"size"`
	Source string `json:"source"`
}

// StringCandidate captures a printable string recovered from a scan region.
type StringCandidate struct {
	Value      string     `json:"value"`
	Region     string     `json:"region"`
	Addr       uint64     `json:"addr,omitempty"`
	Offset     uint64     `json:"offset"`
	Provenance Provenance `json:"provenance"`
}

// SourcePackage captures a projected package node in a source-like tree.
type SourcePackage struct {
	Name            string   `json:"name"`
	ImportPath      string   `json:"import_path"`
	FunctionCount   int      `json:"function_count,omitempty"`
	HasFileEvidence bool     `json:"has_file_evidence"`
	Files           []string `json:"files"`
}

// SourceTree is a minimal source-like projection recovered from file paths.
type SourceTree struct {
	Root             string          `json:"root"`
	Files            []string        `json:"files"`
	Packages         []SourcePackage `json:"packages"`
	ExternalPackages []SourcePackage `json:"external_packages,omitempty"`
}

// Analysis is the minimal canonical analysis result for Sprint 1.
type Analysis struct {
	Input         Input             `json:"input"`
	Provenance    Provenance        `json:"provenance"`
	BuildInfo     *BuildInfo        `json:"build_info,omitempty"`
	Runtime       *RuntimeMetadata  `json:"runtime,omitempty"`
	Functions     []Function        `json:"functions,omitempty"`
	Packages      []Package         `json:"packages,omitempty"`
	Types         []Type            `json:"types"`
	StringRegions []StringRegion    `json:"string_regions,omitempty"`
	Strings       []StringCandidate `json:"strings,omitempty"`
	SourceTree    *SourceTree       `json:"source_tree,omitempty"`
	Refined       *RefinedAnalysis  `json:"refined,omitempty"`
}
