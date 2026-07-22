package schema

// The private v1 wire graph deliberately duplicates every serialized field.
// Canonical schema structs evolve; these DTOs may change only after an explicit
// v1 compatibility decision and fixture update.

type idaExportV1Wire struct {
	Contract   string                 `json:"contract"`
	Input      inputV1Wire            `json:"input"`
	BuildInfo  *buildInfoV1Wire       `json:"build_info,omitempty"`
	Runtime    *runtimeMetadataV1Wire `json:"runtime,omitempty"`
	Peeling    *peelingAnalysisV1Wire `json:"peeling,omitempty"`
	Functions  []idaFunctionV1Wire    `json:"functions,omitempty"`
	Packages   []packageV1Wire        `json:"packages,omitempty"`
	Types      []idaTypeV1Wire        `json:"types,omitempty"`
	Strings    []idaStringV1Wire      `json:"strings,omitempty"`
	SourceTree *sourceTreeV1Wire      `json:"source_tree,omitempty"`
	Refined    *refinedSummaryV1Wire  `json:"refined,omitempty"`
}

type ghidraExportV1Wire struct {
	Contract   string                 `json:"contract"`
	Program    ghidraProgramV1Wire    `json:"program"`
	Runtime    *runtimeMetadataV1Wire `json:"runtime,omitempty"`
	Peeling    *peelingAnalysisV1Wire `json:"peeling,omitempty"`
	Symbols    []ghidraSymbolV1Wire   `json:"symbols,omitempty"`
	Packages   []packageV1Wire        `json:"packages,omitempty"`
	Types      []ghidraTypeV1Wire     `json:"types,omitempty"`
	Strings    []ghidraStringV1Wire   `json:"strings,omitempty"`
	SourceTree *sourceTreeV1Wire      `json:"source_tree,omitempty"`
	Refined    *refinedSummaryV1Wire  `json:"refined,omitempty"`
}

type inputV1Wire struct {
	Path   string `json:"path"`
	Size   int64  `json:"size"`
	Format string `json:"format"`
}

type buildInfoV1Wire struct {
	GoVersion  string           `json:"go_version"`
	Path       string           `json:"path"`
	Provenance provenanceV1Wire `json:"provenance"`
}

type runtimeMetadataV1Wire struct {
	TrustSummary                         string           `json:"trust_summary,omitempty"`
	ELFPclntabHeaderMagic                string           `json:"elf_pclntab_header_magic,omitempty"`
	ELFPclntabHeaderMagicKind            string           `json:"elf_pclntab_header_magic_kind,omitempty"`
	ELFPclntabHeaderQuantum              uint64           `json:"elf_pclntab_header_quantum,omitempty"`
	ELFPclntabHeaderPointerSize          uint64           `json:"elf_pclntab_header_pointer_size,omitempty"`
	ELFPclntabFunctionCountHint          uint64           `json:"elf_pclntab_function_count_hint,omitempty"`
	ELFPclntabFileCountHint              uint64           `json:"elf_pclntab_file_count_hint,omitempty"`
	ELFPclntabFuncnametabOffsetHint      uint64           `json:"elf_pclntab_funcnametab_offset_hint,omitempty"`
	ELFPclntabCuOffsetHint               uint64           `json:"elf_pclntab_cu_offset_hint,omitempty"`
	ELFPclntabFiletabOffsetHint          uint64           `json:"elf_pclntab_filetab_offset_hint,omitempty"`
	ELFPclntabPctabOffsetHint            uint64           `json:"elf_pclntab_pctab_offset_hint,omitempty"`
	ELFPclntabFunctabOffsetHint          uint64           `json:"elf_pclntab_functab_offset_hint,omitempty"`
	ELFFunctabFirstPCOffsetHint          uint64           `json:"elf_functab_first_pc_offset_hint,omitempty"`
	ELFFunctabLastPCOffsetHint           uint64           `json:"elf_functab_last_pc_offset_hint,omitempty"`
	ELFFunctabPCOffsetsMonotonic         bool             `json:"elf_functab_pc_offsets_monotonic,omitempty"`
	ELFTextSectionAddr                   uint64           `json:"elf_text_section_addr,omitempty"`
	ELFTextSectionEndInclusive           uint64           `json:"elf_text_section_end_inclusive,omitempty"`
	ELFFunctabFirstPCAddrHint            uint64           `json:"elf_functab_first_pc_addr_hint,omitempty"`
	ELFFunctabLastPCAddrHint             uint64           `json:"elf_functab_last_pc_addr_hint,omitempty"`
	ELFFunctabPCAddrHintsWithinText      bool             `json:"elf_functab_pc_addr_hints_within_text,omitempty"`
	ELFFunctabPCAddrSample               []uint64         `json:"elf_functab_pc_addr_sample,omitempty"`
	ELFFunctabPCAddrSampleAllWithinText  bool             `json:"elf_functab_pc_addr_sample_all_within_text,omitempty"`
	ELFFunctionFoothold                  string           `json:"elf_function_foothold,omitempty"`
	ELFFunctionFootholdCountHint         uint64           `json:"elf_function_foothold_count_hint,omitempty"`
	ELFFunctionFootholdTextSource        string           `json:"elf_function_foothold_text_source,omitempty"`
	ELFFunctionFootholdStartAddr         uint64           `json:"elf_function_foothold_start_addr,omitempty"`
	ELFFunctionFootholdEndAddr           uint64           `json:"elf_function_foothold_end_addr,omitempty"`
	ELFFunctionRecoveryBlocker           string           `json:"elf_function_recovery_blocker,omitempty"`
	PETextSectionAddr                    uint64           `json:"pe_text_section_addr,omitempty"`
	PETextSectionSize                    uint64           `json:"pe_text_section_size,omitempty"`
	PERdataSectionAddr                   uint64           `json:"pe_rdata_section_addr,omitempty"`
	PERdataSectionSize                   uint64           `json:"pe_rdata_section_size,omitempty"`
	PEPclntabMagicSection                string           `json:"pe_pclntab_magic_section,omitempty"`
	PEPclntabMagicAddr                   uint64           `json:"pe_pclntab_magic_addr,omitempty"`
	PEPclntabMagicCount                  uint64           `json:"pe_pclntab_magic_count,omitempty"`
	PEPclntabHeaderSection               string           `json:"pe_pclntab_header_section,omitempty"`
	PEPclntabHeaderAddr                  uint64           `json:"pe_pclntab_header_addr,omitempty"`
	PEPclntabHeaderMagic                 string           `json:"pe_pclntab_header_magic,omitempty"`
	PEPclntabHeaderQuantum               uint64           `json:"pe_pclntab_header_quantum,omitempty"`
	PEPclntabHeaderPointerSize           uint64           `json:"pe_pclntab_header_pointer_size,omitempty"`
	FirstModuleDataAddr                  uint64           `json:"firstmoduledata_addr,omitempty"`
	FirstModuleDataFromGoModuleFallback  bool             `json:"firstmoduledata_from_go_module_fallback,omitempty"`
	GopclntabAddr                        uint64           `json:"gopclntab_addr,omitempty"`
	GopclntabSize                        uint64           `json:"gopclntab_size,omitempty"`
	TypelinkAddr                         uint64           `json:"typelink_addr,omitempty"`
	TypelinkSize                         uint64           `json:"typelink_size,omitempty"`
	TypelinkCount                        uint64           `json:"typelink_count,omitempty"`
	ItablinkAddr                         uint64           `json:"itablink_addr,omitempty"`
	ItablinkSize                         uint64           `json:"itablink_size,omitempty"`
	ItablinkCount                        uint64           `json:"itablink_count,omitempty"`
	TypelinkSample                       []int32          `json:"typelink_sample,omitempty"`
	TypelinkResolvedBaseAddr             uint64           `json:"typelink_resolved_base_addr,omitempty"`
	TypelinkResolvedSample               []uint64         `json:"typelink_resolved_sample,omitempty"`
	TypelinkResolvedWithinRodataCount    uint64           `json:"typelink_resolved_within_rodata_count,omitempty"`
	TypelinkAllResolvedWithinRodata      bool             `json:"typelink_all_resolved_within_rodata,omitempty"`
	TypelinkMinOffset                    int32            `json:"typelink_min_offset,omitempty"`
	TypelinkMaxOffset                    int32            `json:"typelink_max_offset,omitempty"`
	TypelinkNegativeCount                uint64           `json:"typelink_negative_count,omitempty"`
	TypelinkNonNegativeCount             uint64           `json:"typelink_non_negative_count,omitempty"`
	GoModuleAddr                         uint64           `json:"go_module_addr,omitempty"`
	GoModuleSize                         uint64           `json:"go_module_size,omitempty"`
	FirstModuleDataInGoModule            bool             `json:"firstmoduledata_in_go_module,omitempty"`
	FirstModuleDataGoModuleOffset        uint64           `json:"firstmoduledata_go_module_offset,omitempty"`
	GoModuleWordSize                     uint64           `json:"go_module_word_size,omitempty"`
	GoModuleWordSample                   []uint64         `json:"go_module_word_sample,omitempty"`
	ModuledataPCHeaderAddr               uint64           `json:"moduledata_pcheader_addr,omitempty"`
	ModuledataPCHeaderMatchesGopclntab   bool             `json:"moduledata_pcheader_matches_gopclntab,omitempty"`
	ModuledataFuncnametabSliceWordIndex  uint64           `json:"moduledata_funcnametab_slice_word_index,omitempty"`
	ModuledataFuncnametabAddr            uint64           `json:"moduledata_funcnametab_addr,omitempty"`
	ModuledataFuncnametabLen             uint64           `json:"moduledata_funcnametab_len,omitempty"`
	ModuledataFuncnametabCap             uint64           `json:"moduledata_funcnametab_cap,omitempty"`
	ModuledataFuncnametabWithinGopclntab bool             `json:"moduledata_funcnametab_within_gopclntab,omitempty"`
	ModuledataCutabSliceWordIndex        uint64           `json:"moduledata_cutab_slice_word_index,omitempty"`
	ModuledataCutabAddr                  uint64           `json:"moduledata_cutab_addr,omitempty"`
	ModuledataCutabLen                   uint64           `json:"moduledata_cutab_len,omitempty"`
	ModuledataCutabCap                   uint64           `json:"moduledata_cutab_cap,omitempty"`
	ModuledataCutabWithinGopclntab       bool             `json:"moduledata_cutab_within_gopclntab,omitempty"`
	ModuledataFiletabSliceWordIndex      uint64           `json:"moduledata_filetab_slice_word_index,omitempty"`
	ModuledataFiletabAddr                uint64           `json:"moduledata_filetab_addr,omitempty"`
	ModuledataFiletabLen                 uint64           `json:"moduledata_filetab_len,omitempty"`
	ModuledataFiletabCap                 uint64           `json:"moduledata_filetab_cap,omitempty"`
	ModuledataFiletabWithinGopclntab     bool             `json:"moduledata_filetab_within_gopclntab,omitempty"`
	ModuledataPctabSliceWordIndex        uint64           `json:"moduledata_pctab_slice_word_index,omitempty"`
	ModuledataPctabAddr                  uint64           `json:"moduledata_pctab_addr,omitempty"`
	ModuledataPctabLen                   uint64           `json:"moduledata_pctab_len,omitempty"`
	ModuledataPctabCap                   uint64           `json:"moduledata_pctab_cap,omitempty"`
	ModuledataPctabWithinGopclntab       bool             `json:"moduledata_pctab_within_gopclntab,omitempty"`
	ModuledataPclntableSliceWordIndex    uint64           `json:"moduledata_pclntable_slice_word_index,omitempty"`
	ModuledataPclntableAddr              uint64           `json:"moduledata_pclntable_addr,omitempty"`
	ModuledataPclntableLen               uint64           `json:"moduledata_pclntable_len,omitempty"`
	ModuledataPclntableCap               uint64           `json:"moduledata_pclntable_cap,omitempty"`
	ModuledataPclntableWithinGopclntab   bool             `json:"moduledata_pclntable_within_gopclntab,omitempty"`
	ModuledataTypelinkSliceWordIndex     uint64           `json:"moduledata_typelink_slice_word_index,omitempty"`
	ModuledataTypelinkLen                uint64           `json:"moduledata_typelink_len,omitempty"`
	ModuledataTypelinkCap                uint64           `json:"moduledata_typelink_cap,omitempty"`
	ModuledataItablinkSliceWordIndex     uint64           `json:"moduledata_itablink_slice_word_index,omitempty"`
	ModuledataItablinkLen                uint64           `json:"moduledata_itablink_len,omitempty"`
	ModuledataItablinkCap                uint64           `json:"moduledata_itablink_cap,omitempty"`
	ModuledataMemoryRangesWordIndex      uint64           `json:"moduledata_memory_ranges_word_index,omitempty"`
	ModuledataNoptrdataAddr              uint64           `json:"moduledata_noptrdata_addr,omitempty"`
	ModuledataNoptrdataEnd               uint64           `json:"moduledata_noptrdata_end,omitempty"`
	ModuledataDataAddr                   uint64           `json:"moduledata_data_addr,omitempty"`
	ModuledataDataEnd                    uint64           `json:"moduledata_data_end,omitempty"`
	ModuledataBssAddr                    uint64           `json:"moduledata_bss_addr,omitempty"`
	ModuledataBssEnd                     uint64           `json:"moduledata_bss_end,omitempty"`
	ModuledataNoptrbssAddr               uint64           `json:"moduledata_noptrbss_addr,omitempty"`
	ModuledataNoptrbssEnd                uint64           `json:"moduledata_noptrbss_end,omitempty"`
	ModuledataRodataWordIndex            uint64           `json:"moduledata_rodata_word_index,omitempty"`
	ModuledataRodataAddr                 uint64           `json:"moduledata_rodata_addr,omitempty"`
	ModuledataRodataEnd                  uint64           `json:"moduledata_rodata_end,omitempty"`
	ModuledataTextWordIndex              uint64           `json:"moduledata_text_word_index,omitempty"`
	ModuledataTextAddr                   uint64           `json:"moduledata_text_addr,omitempty"`
	ModuledataTextEndInclusive           uint64           `json:"moduledata_text_end_inclusive,omitempty"`
	ModuledataTypesRangeWordIndex        uint64           `json:"moduledata_types_range_word_index,omitempty"`
	ModuledataTypesAddr                  uint64           `json:"moduledata_types_addr,omitempty"`
	ModuledataETypesAddr                 uint64           `json:"moduledata_etypes_addr,omitempty"`
	TypelinkResolvedWithinTypesCount     uint64           `json:"typelink_resolved_within_types_count,omitempty"`
	TypelinkAllResolvedWithinTypes       bool             `json:"typelink_all_resolved_within_types,omitempty"`
	Provenance                           provenanceV1Wire `json:"provenance"`
}

type peelingAnalysisV1Wire struct {
	Functions  []peelingFunctionV1Wire `json:"functions,omitempty"`
	Packages   []peelingPackageV1Wire  `json:"packages,omitempty"`
	Provenance provenanceV1Wire        `json:"provenance"`
}

type peelingFunctionV1Wire struct {
	Name                   string `json:"name"`
	Package                string `json:"package,omitempty"`
	ImportPath             string `json:"import_path,omitempty"`
	SourceFile             string `json:"source_file,omitempty"`
	SourceLine             int    `json:"source_line,omitempty"`
	Entry                  uint64 `json:"entry"`
	End                    uint64 `json:"end"`
	ModuleLocal            bool   `json:"module_local,omitempty"`
	Classification         string `json:"classification"`
	ClassificationEvidence string `json:"classification_evidence,omitempty"`
}

type peelingPackageV1Wire struct {
	Name                    string `json:"name"`
	ImportPath              string `json:"import_path,omitempty"`
	ModuleLocal             bool   `json:"module_local,omitempty"`
	FunctionCount           int    `json:"function_count"`
	UserFunctionCount       int    `json:"user_function_count,omitempty"`
	StdlibFunctionCount     int    `json:"stdlib_function_count,omitempty"`
	RuntimeFunctionCount    int    `json:"runtime_function_count,omitempty"`
	ThirdPartyFunctionCount int    `json:"third_party_function_count,omitempty"`
	PrimaryClassification   string `json:"primary_classification"`
}

type packageV1Wire struct {
	Name               string           `json:"name"`
	ImportPath         string           `json:"import_path,omitempty"`
	SourceFileCount    int              `json:"source_file_count,omitempty"`
	FunctionCount      int              `json:"function_count"`
	HasSourceEvidence  bool             `json:"has_source_evidence"`
	SourceEvidenceKind string           `json:"source_evidence_kind,omitempty"`
	ModuleLocal        bool             `json:"module_local,omitempty"`
	Provenance         provenanceV1Wire `json:"provenance"`
}

type sourceEvidenceSummaryV1Wire struct {
	TreeKind                    string `json:"tree_kind,omitempty"`
	DWARFPathPackageCount       int    `json:"dwarf_path_package_count,omitempty"`
	DWARFPathFileCount          int    `json:"dwarf_path_file_count,omitempty"`
	LineTablePackageCount       int    `json:"line_table_package_count,omitempty"`
	LineTableFileCount          int    `json:"line_table_file_count,omitempty"`
	PackageFallbackPackageCount int    `json:"package_fallback_package_count,omitempty"`
	PackageFallbackFileCount    int    `json:"package_fallback_file_count,omitempty"`
	MixedPackageEvidenceKinds   bool   `json:"mixed_package_evidence_kinds,omitempty"`
}

type sourcePackageV1Wire struct {
	Name               string   `json:"name"`
	ImportPath         string   `json:"import_path"`
	FunctionCount      int      `json:"function_count,omitempty"`
	HasFileEvidence    bool     `json:"has_file_evidence"`
	SourceEvidenceKind string   `json:"source_evidence_kind,omitempty"`
	Files              []string `json:"files"`
}

type sourceTreeV1Wire struct {
	Root                  string                      `json:"root"`
	SourceEvidenceKind    string                      `json:"source_evidence_kind,omitempty"`
	SourceEvidenceSummary sourceEvidenceSummaryV1Wire `json:"source_evidence_summary,omitempty"`
	PathlessFileEvidence  bool                        `json:"pathless_file_evidence,omitempty"`
	Files                 []string                    `json:"files"`
	Packages              []sourcePackageV1Wire       `json:"packages"`
	ExternalPackages      []sourcePackageV1Wire       `json:"external_packages,omitempty"`
}

type idaFunctionV1Wire struct {
	Name          string           `json:"name"`
	RefinedName   string           `json:"refined_name,omitempty"`
	Package       string           `json:"package,omitempty"`
	ImportPath    string           `json:"import_path,omitempty"`
	SourceFile    string           `json:"source_file,omitempty"`
	SourceLine    int              `json:"source_line,omitempty"`
	Autogenerated bool             `json:"autogenerated,omitempty"`
	Entry         uint64           `json:"entry"`
	End           uint64           `json:"end"`
	ModuleLocal   bool             `json:"module_local,omitempty"`
	Provenance    provenanceV1Wire `json:"provenance"`
}

type idaTypeV1Wire struct {
	Name            string           `json:"name"`
	RefinedName     string           `json:"refined_name,omitempty"`
	Package         string           `json:"package,omitempty"`
	ImportPath      string           `json:"import_path,omitempty"`
	Kind            string           `json:"kind"`
	SourceFileCount int              `json:"source_file_count,omitempty"`
	ModuleLocal     bool             `json:"module_local,omitempty"`
	UserMeaningful  bool             `json:"user_meaningful,omitempty"`
	Provenance      provenanceV1Wire `json:"provenance"`
}

type idaStringV1Wire struct {
	Value        string           `json:"value"`
	RefinedValue string           `json:"refined_value,omitempty"`
	Region       string           `json:"region"`
	Address      uint64           `json:"address,omitempty"`
	Offset       uint64           `json:"offset"`
	Provenance   provenanceV1Wire `json:"provenance"`
}

type ghidraProgramV1Wire struct {
	Path       string           `json:"path"`
	Format     string           `json:"format"`
	ModulePath string           `json:"module_path,omitempty"`
	GoVersion  string           `json:"go_version,omitempty"`
	Provenance provenanceV1Wire `json:"provenance"`
}

type ghidraSymbolV1Wire struct {
	Name          string           `json:"name"`
	RefinedName   string           `json:"refined_name,omitempty"`
	Package       string           `json:"package,omitempty"`
	ImportPath    string           `json:"import_path,omitempty"`
	SourceFile    string           `json:"source_file,omitempty"`
	SourceLine    int              `json:"source_line,omitempty"`
	Autogenerated bool             `json:"autogenerated,omitempty"`
	Address       uint64           `json:"address"`
	End           uint64           `json:"end"`
	ModuleLocal   bool             `json:"module_local,omitempty"`
	Provenance    provenanceV1Wire `json:"provenance"`
}

type ghidraTypeV1Wire struct {
	Name            string           `json:"name"`
	RefinedName     string           `json:"refined_name,omitempty"`
	Package         string           `json:"package,omitempty"`
	ImportPath      string           `json:"import_path,omitempty"`
	Kind            string           `json:"kind"`
	SourceFileCount int              `json:"source_file_count,omitempty"`
	ModuleLocal     bool             `json:"module_local,omitempty"`
	UserMeaningful  bool             `json:"user_meaningful,omitempty"`
	Provenance      provenanceV1Wire `json:"provenance"`
}

type ghidraStringV1Wire struct {
	Value        string           `json:"value"`
	RefinedValue string           `json:"refined_value,omitempty"`
	Region       string           `json:"region"`
	Address      uint64           `json:"address,omitempty"`
	Offset       uint64           `json:"offset"`
	Provenance   provenanceV1Wire `json:"provenance"`
}

type refinedSummaryV1Wire struct {
	Passes []string `json:"passes,omitempty"`
}

type provenanceV1Wire struct {
	Source     string `json:"source"`
	Confidence string `json:"confidence"`
}

func newIDAExportV1Wire(source IDAExport) idaExportV1Wire {
	return idaExportV1Wire{
		Contract:   source.Contract,
		Input:      newInputV1Wire(source.Input),
		BuildInfo:  newBuildInfoV1Wire(source.BuildInfo),
		Runtime:    newRuntimeMetadataV1Wire(source.Runtime),
		Peeling:    newPeelingAnalysisV1Wire(source.Peeling),
		Functions:  mapV1Slice(source.Functions, newIDAFunctionV1Wire),
		Packages:   mapV1Slice(source.Packages, newPackageV1Wire),
		Types:      mapV1Slice(source.Types, newIDATypeV1Wire),
		Strings:    mapV1Slice(source.Strings, newIDAStringV1Wire),
		SourceTree: newSourceTreeV1Wire(source.SourceTree),
		Refined:    newRefinedSummaryV1Wire(source.Refined),
	}
}

func newGhidraExportV1Wire(source GhidraExport) ghidraExportV1Wire {
	return ghidraExportV1Wire{
		Contract:   source.Contract,
		Program:    newGhidraProgramV1Wire(source.Program),
		Runtime:    newRuntimeMetadataV1Wire(source.Runtime),
		Peeling:    newPeelingAnalysisV1Wire(source.Peeling),
		Symbols:    mapV1Slice(source.Symbols, newGhidraSymbolV1Wire),
		Packages:   mapV1Slice(source.Packages, newPackageV1Wire),
		Types:      mapV1Slice(source.Types, newGhidraTypeV1Wire),
		Strings:    mapV1Slice(source.Strings, newGhidraStringV1Wire),
		SourceTree: newSourceTreeV1Wire(source.SourceTree),
		Refined:    newRefinedSummaryV1Wire(source.Refined),
	}
}

func newInputV1Wire(source Input) inputV1Wire {
	return inputV1Wire{
		Path:   source.Path,
		Size:   source.Size,
		Format: source.Format,
	}
}

func newBuildInfoV1Wire(source *BuildInfo) *buildInfoV1Wire {
	if source == nil {
		return nil
	}
	return &buildInfoV1Wire{
		GoVersion:  source.GoVersion,
		Path:       source.Path,
		Provenance: newProvenanceV1Wire(source.Provenance),
	}
}

func newRuntimeMetadataV1Wire(source *RuntimeMetadata) *runtimeMetadataV1Wire {
	if source == nil {
		return nil
	}
	return &runtimeMetadataV1Wire{
		TrustSummary:                         string(source.TrustSummary),
		ELFPclntabHeaderMagic:                source.ELFPclntabHeaderMagic,
		ELFPclntabHeaderMagicKind:            source.ELFPclntabHeaderMagicKind,
		ELFPclntabHeaderQuantum:              source.ELFPclntabHeaderQuantum,
		ELFPclntabHeaderPointerSize:          source.ELFPclntabHeaderPointerSize,
		ELFPclntabFunctionCountHint:          source.ELFPclntabFunctionCountHint,
		ELFPclntabFileCountHint:              source.ELFPclntabFileCountHint,
		ELFPclntabFuncnametabOffsetHint:      source.ELFPclntabFuncnametabOffsetHint,
		ELFPclntabCuOffsetHint:               source.ELFPclntabCuOffsetHint,
		ELFPclntabFiletabOffsetHint:          source.ELFPclntabFiletabOffsetHint,
		ELFPclntabPctabOffsetHint:            source.ELFPclntabPctabOffsetHint,
		ELFPclntabFunctabOffsetHint:          source.ELFPclntabFunctabOffsetHint,
		ELFFunctabFirstPCOffsetHint:          source.ELFFunctabFirstPCOffsetHint,
		ELFFunctabLastPCOffsetHint:           source.ELFFunctabLastPCOffsetHint,
		ELFFunctabPCOffsetsMonotonic:         source.ELFFunctabPCOffsetsMonotonic,
		ELFTextSectionAddr:                   source.ELFTextSectionAddr,
		ELFTextSectionEndInclusive:           source.ELFTextSectionEndInclusive,
		ELFFunctabFirstPCAddrHint:            source.ELFFunctabFirstPCAddrHint,
		ELFFunctabLastPCAddrHint:             source.ELFFunctabLastPCAddrHint,
		ELFFunctabPCAddrHintsWithinText:      source.ELFFunctabPCAddrHintsWithinText,
		ELFFunctabPCAddrSample:               cloneV1Slice(source.ELFFunctabPCAddrSample),
		ELFFunctabPCAddrSampleAllWithinText:  source.ELFFunctabPCAddrSampleAllWithinText,
		ELFFunctionFoothold:                  source.ELFFunctionFoothold,
		ELFFunctionFootholdCountHint:         source.ELFFunctionFootholdCountHint,
		ELFFunctionFootholdTextSource:        source.ELFFunctionFootholdTextSource,
		ELFFunctionFootholdStartAddr:         source.ELFFunctionFootholdStartAddr,
		ELFFunctionFootholdEndAddr:           source.ELFFunctionFootholdEndAddr,
		ELFFunctionRecoveryBlocker:           source.ELFFunctionRecoveryBlocker,
		PETextSectionAddr:                    source.PETextSectionAddr,
		PETextSectionSize:                    source.PETextSectionSize,
		PERdataSectionAddr:                   source.PERdataSectionAddr,
		PERdataSectionSize:                   source.PERdataSectionSize,
		PEPclntabMagicSection:                source.PEPclntabMagicSection,
		PEPclntabMagicAddr:                   source.PEPclntabMagicAddr,
		PEPclntabMagicCount:                  source.PEPclntabMagicCount,
		PEPclntabHeaderSection:               source.PEPclntabHeaderSection,
		PEPclntabHeaderAddr:                  source.PEPclntabHeaderAddr,
		PEPclntabHeaderMagic:                 source.PEPclntabHeaderMagic,
		PEPclntabHeaderQuantum:               source.PEPclntabHeaderQuantum,
		PEPclntabHeaderPointerSize:           source.PEPclntabHeaderPointerSize,
		FirstModuleDataAddr:                  source.FirstModuleDataAddr,
		FirstModuleDataFromGoModuleFallback:  source.FirstModuleDataFromGoModuleFallback,
		GopclntabAddr:                        source.GopclntabAddr,
		GopclntabSize:                        source.GopclntabSize,
		TypelinkAddr:                         source.TypelinkAddr,
		TypelinkSize:                         source.TypelinkSize,
		TypelinkCount:                        source.TypelinkCount,
		ItablinkAddr:                         source.ItablinkAddr,
		ItablinkSize:                         source.ItablinkSize,
		ItablinkCount:                        source.ItablinkCount,
		TypelinkSample:                       cloneV1Slice(source.TypelinkSample),
		TypelinkResolvedBaseAddr:             source.TypelinkResolvedBaseAddr,
		TypelinkResolvedSample:               cloneV1Slice(source.TypelinkResolvedSample),
		TypelinkResolvedWithinRodataCount:    source.TypelinkResolvedWithinRodataCount,
		TypelinkAllResolvedWithinRodata:      source.TypelinkAllResolvedWithinRodata,
		TypelinkMinOffset:                    source.TypelinkMinOffset,
		TypelinkMaxOffset:                    source.TypelinkMaxOffset,
		TypelinkNegativeCount:                source.TypelinkNegativeCount,
		TypelinkNonNegativeCount:             source.TypelinkNonNegativeCount,
		GoModuleAddr:                         source.GoModuleAddr,
		GoModuleSize:                         source.GoModuleSize,
		FirstModuleDataInGoModule:            source.FirstModuleDataInGoModule,
		FirstModuleDataGoModuleOffset:        source.FirstModuleDataGoModuleOffset,
		GoModuleWordSize:                     source.GoModuleWordSize,
		GoModuleWordSample:                   cloneV1Slice(source.GoModuleWordSample),
		ModuledataPCHeaderAddr:               source.ModuledataPCHeaderAddr,
		ModuledataPCHeaderMatchesGopclntab:   source.ModuledataPCHeaderMatchesGopclntab,
		ModuledataFuncnametabSliceWordIndex:  source.ModuledataFuncnametabSliceWordIndex,
		ModuledataFuncnametabAddr:            source.ModuledataFuncnametabAddr,
		ModuledataFuncnametabLen:             source.ModuledataFuncnametabLen,
		ModuledataFuncnametabCap:             source.ModuledataFuncnametabCap,
		ModuledataFuncnametabWithinGopclntab: source.ModuledataFuncnametabWithinGopclntab,
		ModuledataCutabSliceWordIndex:        source.ModuledataCutabSliceWordIndex,
		ModuledataCutabAddr:                  source.ModuledataCutabAddr,
		ModuledataCutabLen:                   source.ModuledataCutabLen,
		ModuledataCutabCap:                   source.ModuledataCutabCap,
		ModuledataCutabWithinGopclntab:       source.ModuledataCutabWithinGopclntab,
		ModuledataFiletabSliceWordIndex:      source.ModuledataFiletabSliceWordIndex,
		ModuledataFiletabAddr:                source.ModuledataFiletabAddr,
		ModuledataFiletabLen:                 source.ModuledataFiletabLen,
		ModuledataFiletabCap:                 source.ModuledataFiletabCap,
		ModuledataFiletabWithinGopclntab:     source.ModuledataFiletabWithinGopclntab,
		ModuledataPctabSliceWordIndex:        source.ModuledataPctabSliceWordIndex,
		ModuledataPctabAddr:                  source.ModuledataPctabAddr,
		ModuledataPctabLen:                   source.ModuledataPctabLen,
		ModuledataPctabCap:                   source.ModuledataPctabCap,
		ModuledataPctabWithinGopclntab:       source.ModuledataPctabWithinGopclntab,
		ModuledataPclntableSliceWordIndex:    source.ModuledataPclntableSliceWordIndex,
		ModuledataPclntableAddr:              source.ModuledataPclntableAddr,
		ModuledataPclntableLen:               source.ModuledataPclntableLen,
		ModuledataPclntableCap:               source.ModuledataPclntableCap,
		ModuledataPclntableWithinGopclntab:   source.ModuledataPclntableWithinGopclntab,
		ModuledataTypelinkSliceWordIndex:     source.ModuledataTypelinkSliceWordIndex,
		ModuledataTypelinkLen:                source.ModuledataTypelinkLen,
		ModuledataTypelinkCap:                source.ModuledataTypelinkCap,
		ModuledataItablinkSliceWordIndex:     source.ModuledataItablinkSliceWordIndex,
		ModuledataItablinkLen:                source.ModuledataItablinkLen,
		ModuledataItablinkCap:                source.ModuledataItablinkCap,
		ModuledataMemoryRangesWordIndex:      source.ModuledataMemoryRangesWordIndex,
		ModuledataNoptrdataAddr:              source.ModuledataNoptrdataAddr,
		ModuledataNoptrdataEnd:               source.ModuledataNoptrdataEnd,
		ModuledataDataAddr:                   source.ModuledataDataAddr,
		ModuledataDataEnd:                    source.ModuledataDataEnd,
		ModuledataBssAddr:                    source.ModuledataBssAddr,
		ModuledataBssEnd:                     source.ModuledataBssEnd,
		ModuledataNoptrbssAddr:               source.ModuledataNoptrbssAddr,
		ModuledataNoptrbssEnd:                source.ModuledataNoptrbssEnd,
		ModuledataRodataWordIndex:            source.ModuledataRodataWordIndex,
		ModuledataRodataAddr:                 source.ModuledataRodataAddr,
		ModuledataRodataEnd:                  source.ModuledataRodataEnd,
		ModuledataTextWordIndex:              source.ModuledataTextWordIndex,
		ModuledataTextAddr:                   source.ModuledataTextAddr,
		ModuledataTextEndInclusive:           source.ModuledataTextEndInclusive,
		ModuledataTypesRangeWordIndex:        source.ModuledataTypesRangeWordIndex,
		ModuledataTypesAddr:                  source.ModuledataTypesAddr,
		ModuledataETypesAddr:                 source.ModuledataETypesAddr,
		TypelinkResolvedWithinTypesCount:     source.TypelinkResolvedWithinTypesCount,
		TypelinkAllResolvedWithinTypes:       source.TypelinkAllResolvedWithinTypes,
		Provenance:                           newProvenanceV1Wire(source.Provenance),
	}
}

func newPeelingAnalysisV1Wire(source *PeelingAnalysis) *peelingAnalysisV1Wire {
	if source == nil {
		return nil
	}
	return &peelingAnalysisV1Wire{
		Functions:  mapV1Slice(source.Functions, newPeelingFunctionV1Wire),
		Packages:   mapV1Slice(source.Packages, newPeelingPackageV1Wire),
		Provenance: newProvenanceV1Wire(source.Provenance),
	}
}

func newPeelingFunctionV1Wire(source PeelingFunction) peelingFunctionV1Wire {
	return peelingFunctionV1Wire{
		Name:                   source.Name,
		Package:                source.Package,
		ImportPath:             source.ImportPath,
		SourceFile:             source.SourceFile,
		SourceLine:             source.SourceLine,
		Entry:                  source.Entry,
		End:                    source.End,
		ModuleLocal:            source.ModuleLocal,
		Classification:         string(source.Classification),
		ClassificationEvidence: string(source.ClassificationEvidence),
	}
}

func newPeelingPackageV1Wire(source PeelingPackage) peelingPackageV1Wire {
	return peelingPackageV1Wire{
		Name:                    source.Name,
		ImportPath:              source.ImportPath,
		ModuleLocal:             source.ModuleLocal,
		FunctionCount:           source.FunctionCount,
		UserFunctionCount:       source.UserFunctionCount,
		StdlibFunctionCount:     source.StdlibFunctionCount,
		RuntimeFunctionCount:    source.RuntimeFunctionCount,
		ThirdPartyFunctionCount: source.ThirdPartyFunctionCount,
		PrimaryClassification:   string(source.PrimaryClassification),
	}
}

func newPackageV1Wire(source Package) packageV1Wire {
	return packageV1Wire{
		Name:               source.Name,
		ImportPath:         source.ImportPath,
		SourceFileCount:    source.SourceFileCount,
		FunctionCount:      source.FunctionCount,
		HasSourceEvidence:  source.HasSourceEvidence,
		SourceEvidenceKind: string(source.SourceEvidenceKind),
		ModuleLocal:        source.ModuleLocal,
		Provenance:         newProvenanceV1Wire(source.Provenance),
	}
}

func newSourceTreeV1Wire(source *SourceTree) *sourceTreeV1Wire {
	if source == nil {
		return nil
	}
	return &sourceTreeV1Wire{
		Root:                  source.Root,
		SourceEvidenceKind:    string(source.SourceEvidenceKind),
		SourceEvidenceSummary: newSourceEvidenceSummaryV1Wire(source.SourceEvidenceSummary),
		PathlessFileEvidence:  source.PathlessFileEvidence,
		Files:                 cloneV1Slice(source.Files),
		Packages:              mapV1Slice(source.Packages, newSourcePackageV1Wire),
		ExternalPackages:      mapV1Slice(source.ExternalPackages, newSourcePackageV1Wire),
	}
}

func newSourceEvidenceSummaryV1Wire(source SourceEvidenceSummary) sourceEvidenceSummaryV1Wire {
	return sourceEvidenceSummaryV1Wire{
		TreeKind:                    string(source.TreeKind),
		DWARFPathPackageCount:       source.DWARFPathPackageCount,
		DWARFPathFileCount:          source.DWARFPathFileCount,
		LineTablePackageCount:       source.LineTablePackageCount,
		LineTableFileCount:          source.LineTableFileCount,
		PackageFallbackPackageCount: source.PackageFallbackPackageCount,
		PackageFallbackFileCount:    source.PackageFallbackFileCount,
		MixedPackageEvidenceKinds:   source.MixedPackageEvidenceKinds,
	}
}

func newSourcePackageV1Wire(source SourcePackage) sourcePackageV1Wire {
	return sourcePackageV1Wire{
		Name:               source.Name,
		ImportPath:         source.ImportPath,
		FunctionCount:      source.FunctionCount,
		HasFileEvidence:    source.HasFileEvidence,
		SourceEvidenceKind: string(source.SourceEvidenceKind),
		Files:              cloneV1Slice(source.Files),
	}
}

func newIDAFunctionV1Wire(source IDAFunction) idaFunctionV1Wire {
	return idaFunctionV1Wire{
		Name:          source.Name,
		RefinedName:   source.RefinedName,
		Package:       source.Package,
		ImportPath:    source.ImportPath,
		SourceFile:    source.SourceFile,
		SourceLine:    source.SourceLine,
		Autogenerated: source.Autogenerated,
		Entry:         source.Entry,
		End:           source.End,
		ModuleLocal:   source.ModuleLocal,
		Provenance:    newProvenanceV1Wire(source.Provenance),
	}
}

func newIDATypeV1Wire(source IDAType) idaTypeV1Wire {
	return idaTypeV1Wire{
		Name:            source.Name,
		RefinedName:     source.RefinedName,
		Package:         source.Package,
		ImportPath:      source.ImportPath,
		Kind:            source.Kind,
		SourceFileCount: source.SourceFileCount,
		ModuleLocal:     source.ModuleLocal,
		UserMeaningful:  source.UserMeaningful,
		Provenance:      newProvenanceV1Wire(source.Provenance),
	}
}

func newIDAStringV1Wire(source IDAString) idaStringV1Wire {
	return idaStringV1Wire{
		Value:        source.Value,
		RefinedValue: source.RefinedValue,
		Region:       source.Region,
		Address:      source.Address,
		Offset:       source.Offset,
		Provenance:   newProvenanceV1Wire(source.Provenance),
	}
}

func newGhidraProgramV1Wire(source GhidraProgram) ghidraProgramV1Wire {
	return ghidraProgramV1Wire{
		Path:       source.Path,
		Format:     source.Format,
		ModulePath: source.ModulePath,
		GoVersion:  source.GoVersion,
		Provenance: newProvenanceV1Wire(source.Provenance),
	}
}

func newGhidraSymbolV1Wire(source GhidraSymbol) ghidraSymbolV1Wire {
	return ghidraSymbolV1Wire{
		Name:          source.Name,
		RefinedName:   source.RefinedName,
		Package:       source.Package,
		ImportPath:    source.ImportPath,
		SourceFile:    source.SourceFile,
		SourceLine:    source.SourceLine,
		Autogenerated: source.Autogenerated,
		Address:       source.Address,
		End:           source.End,
		ModuleLocal:   source.ModuleLocal,
		Provenance:    newProvenanceV1Wire(source.Provenance),
	}
}

func newGhidraTypeV1Wire(source GhidraType) ghidraTypeV1Wire {
	return ghidraTypeV1Wire{
		Name:            source.Name,
		RefinedName:     source.RefinedName,
		Package:         source.Package,
		ImportPath:      source.ImportPath,
		Kind:            source.Kind,
		SourceFileCount: source.SourceFileCount,
		ModuleLocal:     source.ModuleLocal,
		UserMeaningful:  source.UserMeaningful,
		Provenance:      newProvenanceV1Wire(source.Provenance),
	}
}

func newGhidraStringV1Wire(source GhidraString) ghidraStringV1Wire {
	return ghidraStringV1Wire{
		Value:        source.Value,
		RefinedValue: source.RefinedValue,
		Region:       source.Region,
		Address:      source.Address,
		Offset:       source.Offset,
		Provenance:   newProvenanceV1Wire(source.Provenance),
	}
}

func newRefinedSummaryV1Wire(source *RefinedSummary) *refinedSummaryV1Wire {
	if source == nil {
		return nil
	}
	return &refinedSummaryV1Wire{Passes: cloneV1Slice(source.Passes)}
}

func newProvenanceV1Wire(source Provenance) provenanceV1Wire {
	return provenanceV1Wire{
		Source:     source.Source,
		Confidence: source.Confidence,
	}
}

func mapV1Slice[Source, Destination any](source []Source, project func(Source) Destination) []Destination {
	if source == nil {
		return nil
	}
	destination := make([]Destination, len(source))
	for index, item := range source {
		destination[index] = project(item)
	}
	return destination
}

func cloneV1Slice[Item any](source []Item) []Item {
	if source == nil {
		return nil
	}
	destination := make([]Item, len(source))
	copy(destination, source)
	return destination
}
