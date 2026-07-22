package schema

// RuntimeTrustSummary compresses the current runtime evidence posture into one bounded field.
type RuntimeTrustSummary string

const (
	RuntimeTrustSummaryAbsent           RuntimeTrustSummary = "absent"
	RuntimeTrustSummarySectionHeuristic RuntimeTrustSummary = "section_heuristic"
	RuntimeTrustSummaryGoModuleFallback RuntimeTrustSummary = "go_module_fallback"
	RuntimeTrustSummarySymbolBacked     RuntimeTrustSummary = "symbol_backed"
)

// SourceEvidenceKind classifies how source/file-related truth entered the current projection.
type SourceEvidenceKind string

const (
	SourceEvidenceKindDWARFPaths      SourceEvidenceKind = "dwarf_paths"
	SourceEvidenceKindLineTableFiles  SourceEvidenceKind = "line_table_files"
	SourceEvidenceKindPackageFallback SourceEvidenceKind = "package_fallback"
)

// SourceEvidenceSummary compresses the source/file evidence landscape into a bounded cue.
type SourceEvidenceSummary struct {
	TreeKind                    SourceEvidenceKind `json:"tree_kind,omitempty"`
	DWARFPathPackageCount       int                `json:"dwarf_path_package_count,omitempty"`
	DWARFPathFileCount          int                `json:"dwarf_path_file_count,omitempty"`
	LineTablePackageCount       int                `json:"line_table_package_count,omitempty"`
	LineTableFileCount          int                `json:"line_table_file_count,omitempty"`
	PackageFallbackPackageCount int                `json:"package_fallback_package_count,omitempty"`
	PackageFallbackFileCount    int                `json:"package_fallback_file_count,omitempty"`
	MixedPackageEvidenceKinds   bool               `json:"mixed_package_evidence_kinds,omitempty"`
}

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
	TrustSummary                         RuntimeTrustSummary `json:"trust_summary,omitempty"`
	ELFPclntabHeaderMagic                string              `json:"elf_pclntab_header_magic,omitempty"`
	ELFPclntabHeaderMagicKind            string              `json:"elf_pclntab_header_magic_kind,omitempty"`
	ELFPclntabHeaderQuantum              uint64              `json:"elf_pclntab_header_quantum,omitempty"`
	ELFPclntabHeaderPointerSize          uint64              `json:"elf_pclntab_header_pointer_size,omitempty"`
	ELFPclntabFunctionCountHint          uint64              `json:"elf_pclntab_function_count_hint,omitempty"`
	ELFPclntabFileCountHint              uint64              `json:"elf_pclntab_file_count_hint,omitempty"`
	ELFPclntabFuncnametabOffsetHint      uint64              `json:"elf_pclntab_funcnametab_offset_hint,omitempty"`
	ELFPclntabCuOffsetHint               uint64              `json:"elf_pclntab_cu_offset_hint,omitempty"`
	ELFPclntabFiletabOffsetHint          uint64              `json:"elf_pclntab_filetab_offset_hint,omitempty"`
	ELFPclntabPctabOffsetHint            uint64              `json:"elf_pclntab_pctab_offset_hint,omitempty"`
	ELFPclntabFunctabOffsetHint          uint64              `json:"elf_pclntab_functab_offset_hint,omitempty"`
	ELFFunctabFirstPCOffsetHint          uint64              `json:"elf_functab_first_pc_offset_hint,omitempty"`
	ELFFunctabLastPCOffsetHint           uint64              `json:"elf_functab_last_pc_offset_hint,omitempty"`
	ELFFunctabPCOffsetsMonotonic         bool                `json:"elf_functab_pc_offsets_monotonic,omitempty"`
	ELFTextSectionAddr                   uint64              `json:"elf_text_section_addr,omitempty"`
	ELFTextSectionEndInclusive           uint64              `json:"elf_text_section_end_inclusive,omitempty"`
	ELFFunctabFirstPCAddrHint            uint64              `json:"elf_functab_first_pc_addr_hint,omitempty"`
	ELFFunctabLastPCAddrHint             uint64              `json:"elf_functab_last_pc_addr_hint,omitempty"`
	ELFFunctabPCAddrHintsWithinText      bool                `json:"elf_functab_pc_addr_hints_within_text,omitempty"`
	ELFFunctabPCAddrSample               []uint64            `json:"elf_functab_pc_addr_sample,omitempty"`
	ELFFunctabPCAddrSampleAllWithinText  bool                `json:"elf_functab_pc_addr_sample_all_within_text,omitempty"`
	ELFFunctionFoothold                  string              `json:"elf_function_foothold,omitempty"`
	ELFFunctionFootholdCountHint         uint64              `json:"elf_function_foothold_count_hint,omitempty"`
	ELFFunctionFootholdTextSource        string              `json:"elf_function_foothold_text_source,omitempty"`
	ELFFunctionFootholdStartAddr         uint64              `json:"elf_function_foothold_start_addr,omitempty"`
	ELFFunctionFootholdEndAddr           uint64              `json:"elf_function_foothold_end_addr,omitempty"`
	ELFFunctionRecoveryBlocker           string              `json:"elf_function_recovery_blocker,omitempty"`
	PETextSectionAddr                    uint64              `json:"pe_text_section_addr,omitempty"`
	PETextSectionSize                    uint64              `json:"pe_text_section_size,omitempty"`
	PERdataSectionAddr                   uint64              `json:"pe_rdata_section_addr,omitempty"`
	PERdataSectionSize                   uint64              `json:"pe_rdata_section_size,omitempty"`
	PEPclntabMagicSection                string              `json:"pe_pclntab_magic_section,omitempty"`
	PEPclntabMagicAddr                   uint64              `json:"pe_pclntab_magic_addr,omitempty"`
	PEPclntabMagicCount                  uint64              `json:"pe_pclntab_magic_count,omitempty"`
	PEPclntabHeaderSection               string              `json:"pe_pclntab_header_section,omitempty"`
	PEPclntabHeaderAddr                  uint64              `json:"pe_pclntab_header_addr,omitempty"`
	PEPclntabHeaderMagic                 string              `json:"pe_pclntab_header_magic,omitempty"`
	PEPclntabHeaderQuantum               uint64              `json:"pe_pclntab_header_quantum,omitempty"`
	PEPclntabHeaderPointerSize           uint64              `json:"pe_pclntab_header_pointer_size,omitempty"`
	FirstModuleDataAddr                  uint64              `json:"firstmoduledata_addr,omitempty"`
	FirstModuleDataFromGoModuleFallback  bool                `json:"firstmoduledata_from_go_module_fallback,omitempty"`
	GopclntabAddr                        uint64              `json:"gopclntab_addr,omitempty"`
	GopclntabSize                        uint64              `json:"gopclntab_size,omitempty"`
	TypelinkAddr                         uint64              `json:"typelink_addr,omitempty"`
	TypelinkSize                         uint64              `json:"typelink_size,omitempty"`
	TypelinkCount                        uint64              `json:"typelink_count,omitempty"`
	ItablinkAddr                         uint64              `json:"itablink_addr,omitempty"`
	ItablinkSize                         uint64              `json:"itablink_size,omitempty"`
	ItablinkCount                        uint64              `json:"itablink_count,omitempty"`
	TypelinkSample                       []int32             `json:"typelink_sample,omitempty"`
	TypelinkResolvedBaseAddr             uint64              `json:"typelink_resolved_base_addr,omitempty"`
	TypelinkResolvedSample               []uint64            `json:"typelink_resolved_sample,omitempty"`
	TypelinkResolvedWithinRodataCount    uint64              `json:"typelink_resolved_within_rodata_count,omitempty"`
	TypelinkAllResolvedWithinRodata      bool                `json:"typelink_all_resolved_within_rodata,omitempty"`
	TypelinkMinOffset                    int32               `json:"typelink_min_offset,omitempty"`
	TypelinkMaxOffset                    int32               `json:"typelink_max_offset,omitempty"`
	TypelinkNegativeCount                uint64              `json:"typelink_negative_count,omitempty"`
	TypelinkNonNegativeCount             uint64              `json:"typelink_non_negative_count,omitempty"`
	GoModuleAddr                         uint64              `json:"go_module_addr,omitempty"`
	GoModuleSize                         uint64              `json:"go_module_size,omitempty"`
	FirstModuleDataInGoModule            bool                `json:"firstmoduledata_in_go_module,omitempty"`
	FirstModuleDataGoModuleOffset        uint64              `json:"firstmoduledata_go_module_offset,omitempty"`
	GoModuleWordSize                     uint64              `json:"go_module_word_size,omitempty"`
	GoModuleWordSample                   []uint64            `json:"go_module_word_sample,omitempty"`
	ModuledataPCHeaderAddr               uint64              `json:"moduledata_pcheader_addr,omitempty"`
	ModuledataPCHeaderMatchesGopclntab   bool                `json:"moduledata_pcheader_matches_gopclntab,omitempty"`
	ModuledataFuncnametabSliceWordIndex  uint64              `json:"moduledata_funcnametab_slice_word_index,omitempty"`
	ModuledataFuncnametabAddr            uint64              `json:"moduledata_funcnametab_addr,omitempty"`
	ModuledataFuncnametabLen             uint64              `json:"moduledata_funcnametab_len,omitempty"`
	ModuledataFuncnametabCap             uint64              `json:"moduledata_funcnametab_cap,omitempty"`
	ModuledataFuncnametabWithinGopclntab bool                `json:"moduledata_funcnametab_within_gopclntab,omitempty"`
	ModuledataCutabSliceWordIndex        uint64              `json:"moduledata_cutab_slice_word_index,omitempty"`
	ModuledataCutabAddr                  uint64              `json:"moduledata_cutab_addr,omitempty"`
	ModuledataCutabLen                   uint64              `json:"moduledata_cutab_len,omitempty"`
	ModuledataCutabCap                   uint64              `json:"moduledata_cutab_cap,omitempty"`
	ModuledataCutabWithinGopclntab       bool                `json:"moduledata_cutab_within_gopclntab,omitempty"`
	ModuledataFiletabSliceWordIndex      uint64              `json:"moduledata_filetab_slice_word_index,omitempty"`
	ModuledataFiletabAddr                uint64              `json:"moduledata_filetab_addr,omitempty"`
	ModuledataFiletabLen                 uint64              `json:"moduledata_filetab_len,omitempty"`
	ModuledataFiletabCap                 uint64              `json:"moduledata_filetab_cap,omitempty"`
	ModuledataFiletabWithinGopclntab     bool                `json:"moduledata_filetab_within_gopclntab,omitempty"`
	ModuledataPctabSliceWordIndex        uint64              `json:"moduledata_pctab_slice_word_index,omitempty"`
	ModuledataPctabAddr                  uint64              `json:"moduledata_pctab_addr,omitempty"`
	ModuledataPctabLen                   uint64              `json:"moduledata_pctab_len,omitempty"`
	ModuledataPctabCap                   uint64              `json:"moduledata_pctab_cap,omitempty"`
	ModuledataPctabWithinGopclntab       bool                `json:"moduledata_pctab_within_gopclntab,omitempty"`
	ModuledataPclntableSliceWordIndex    uint64              `json:"moduledata_pclntable_slice_word_index,omitempty"`
	ModuledataPclntableAddr              uint64              `json:"moduledata_pclntable_addr,omitempty"`
	ModuledataPclntableLen               uint64              `json:"moduledata_pclntable_len,omitempty"`
	ModuledataPclntableCap               uint64              `json:"moduledata_pclntable_cap,omitempty"`
	ModuledataPclntableWithinGopclntab   bool                `json:"moduledata_pclntable_within_gopclntab,omitempty"`
	ModuledataTypelinkSliceWordIndex     uint64              `json:"moduledata_typelink_slice_word_index,omitempty"`
	ModuledataTypelinkLen                uint64              `json:"moduledata_typelink_len,omitempty"`
	ModuledataTypelinkCap                uint64              `json:"moduledata_typelink_cap,omitempty"`
	ModuledataItablinkSliceWordIndex     uint64              `json:"moduledata_itablink_slice_word_index,omitempty"`
	ModuledataItablinkLen                uint64              `json:"moduledata_itablink_len,omitempty"`
	ModuledataItablinkCap                uint64              `json:"moduledata_itablink_cap,omitempty"`
	ModuledataMemoryRangesWordIndex      uint64              `json:"moduledata_memory_ranges_word_index,omitempty"`
	ModuledataNoptrdataAddr              uint64              `json:"moduledata_noptrdata_addr,omitempty"`
	ModuledataNoptrdataEnd               uint64              `json:"moduledata_noptrdata_end,omitempty"`
	ModuledataDataAddr                   uint64              `json:"moduledata_data_addr,omitempty"`
	ModuledataDataEnd                    uint64              `json:"moduledata_data_end,omitempty"`
	ModuledataBssAddr                    uint64              `json:"moduledata_bss_addr,omitempty"`
	ModuledataBssEnd                     uint64              `json:"moduledata_bss_end,omitempty"`
	ModuledataNoptrbssAddr               uint64              `json:"moduledata_noptrbss_addr,omitempty"`
	ModuledataNoptrbssEnd                uint64              `json:"moduledata_noptrbss_end,omitempty"`
	ModuledataRodataWordIndex            uint64              `json:"moduledata_rodata_word_index,omitempty"`
	ModuledataRodataAddr                 uint64              `json:"moduledata_rodata_addr,omitempty"`
	ModuledataRodataEnd                  uint64              `json:"moduledata_rodata_end,omitempty"`
	ModuledataTextWordIndex              uint64              `json:"moduledata_text_word_index,omitempty"`
	ModuledataTextAddr                   uint64              `json:"moduledata_text_addr,omitempty"`
	ModuledataTextEndInclusive           uint64              `json:"moduledata_text_end_inclusive,omitempty"`
	ModuledataTypesRangeWordIndex        uint64              `json:"moduledata_types_range_word_index,omitempty"`
	ModuledataTypesAddr                  uint64              `json:"moduledata_types_addr,omitempty"`
	ModuledataETypesAddr                 uint64              `json:"moduledata_etypes_addr,omitempty"`
	TypelinkResolvedWithinTypesCount     uint64              `json:"typelink_resolved_within_types_count,omitempty"`
	TypelinkAllResolvedWithinTypes       bool                `json:"typelink_all_resolved_within_types,omitempty"`
	Provenance                           Provenance          `json:"provenance"`
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
	Name               string             `json:"name"`
	ImportPath         string             `json:"import_path,omitempty"`
	SourceFileCount    int                `json:"source_file_count,omitempty"`
	FunctionCount      int                `json:"function_count"`
	HasSourceEvidence  bool               `json:"has_source_evidence"`
	SourceEvidenceKind SourceEvidenceKind `json:"source_evidence_kind,omitempty"`
	ModuleLocal        bool               `json:"module_local,omitempty"`
	Provenance         Provenance         `json:"provenance"`
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
	Name               string             `json:"name"`
	ImportPath         string             `json:"import_path"`
	FunctionCount      int                `json:"function_count,omitempty"`
	HasFileEvidence    bool               `json:"has_file_evidence"`
	SourceEvidenceKind SourceEvidenceKind `json:"source_evidence_kind,omitempty"`
	Files              []string           `json:"files"`
}

// SourceTree is a minimal source-like projection recovered from file paths.
type SourceTree struct {
	Root                  string                `json:"root"`
	SourceEvidenceKind    SourceEvidenceKind    `json:"source_evidence_kind,omitempty"`
	SourceEvidenceSummary SourceEvidenceSummary `json:"source_evidence_summary,omitempty"`
	PathlessFileEvidence  bool                  `json:"pathless_file_evidence,omitempty"`
	Files                 []string              `json:"files"`
	Packages              []SourcePackage       `json:"packages"`
	ExternalPackages      []SourcePackage       `json:"external_packages,omitempty"`
}

// Analysis is the minimal canonical analysis result for Sprint 1.
type Analysis struct {
	Input         Input             `json:"input"`
	Provenance    Provenance        `json:"provenance"`
	Diagnostics   []StageDiagnostic `json:"diagnostics"`
	BuildInfo     *BuildInfo        `json:"build_info,omitempty"`
	Runtime       *RuntimeMetadata  `json:"runtime,omitempty"`
	Functions     []Function        `json:"functions,omitempty"`
	Packages      []Package         `json:"packages,omitempty"`
	Types         []Type            `json:"types"`
	StringRegions []StringRegion    `json:"string_regions,omitempty"`
	Strings       []StringCandidate `json:"strings,omitempty"`
	SourceTree    *SourceTree       `json:"source_tree,omitempty"`
	Peeling       *PeelingAnalysis  `json:"peeling,omitempty"`
	Refined       *RefinedAnalysis  `json:"refined,omitempty"`
}
