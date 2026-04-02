package schema

// PeelingClass is the bounded analyst-facing function classification used by code peeling.
type PeelingClass string

const (
	PeelingClassUser       PeelingClass = "user"
	PeelingClassStdlib     PeelingClass = "stdlib"
	PeelingClassRuntime    PeelingClass = "runtime"
	PeelingClassThirdParty PeelingClass = "third_party"
)

// PeelingEvidence records the bounded reason behind a code-peeling classification.
type PeelingEvidence string

const (
	PeelingEvidenceModuleLocal              PeelingEvidence = "module_local"
	PeelingEvidenceBuildInfoPath            PeelingEvidence = "build_info_path"
	PeelingEvidencePackageMain              PeelingEvidence = "package_main"
	PeelingEvidenceRuntimeImportPath        PeelingEvidence = "runtime_import_path"
	PeelingEvidenceRuntimeNameFingerprint   PeelingEvidence = "runtime_name_fingerprint"
	PeelingEvidenceRuntimeSourceFingerprint PeelingEvidence = "runtime_source_fingerprint"
	PeelingEvidenceThirdPartyImportPath     PeelingEvidence = "third_party_import_path"
	PeelingEvidenceStdlibImportPath         PeelingEvidence = "stdlib_import_path"
	PeelingEvidenceStdlibNameFingerprint    PeelingEvidence = "stdlib_name_fingerprint"
	PeelingEvidenceStdlibSourceFingerprint  PeelingEvidence = "stdlib_source_fingerprint"
	PeelingEvidenceDefaultStdlib            PeelingEvidence = "default_stdlib"
)

// PeelingFunction is a function-level classification derived from already-known raw truth.
type PeelingFunction struct {
	Name                   string          `json:"name"`
	Package                string          `json:"package,omitempty"`
	ImportPath             string          `json:"import_path,omitempty"`
	SourceFile             string          `json:"source_file,omitempty"`
	SourceLine             int             `json:"source_line,omitempty"`
	Entry                  uint64          `json:"entry"`
	End                    uint64          `json:"end"`
	ModuleLocal            bool            `json:"module_local,omitempty"`
	Classification         PeelingClass    `json:"classification"`
	ClassificationEvidence PeelingEvidence `json:"classification_evidence,omitempty"`
}

// PeelingPackage summarizes function-level code peeling for one package/import path.
type PeelingPackage struct {
	Name                    string       `json:"name"`
	ImportPath              string       `json:"import_path,omitempty"`
	ModuleLocal             bool         `json:"module_local,omitempty"`
	FunctionCount           int          `json:"function_count"`
	UserFunctionCount       int          `json:"user_function_count,omitempty"`
	StdlibFunctionCount     int          `json:"stdlib_function_count,omitempty"`
	RuntimeFunctionCount    int          `json:"runtime_function_count,omitempty"`
	ThirdPartyFunctionCount int          `json:"third_party_function_count,omitempty"`
	PrimaryClassification   PeelingClass `json:"primary_classification"`
}

// PeelingAnalysis is the bounded analyst-facing code-peeling layer derived from canonical truth.
type PeelingAnalysis struct {
	Functions  []PeelingFunction `json:"functions,omitempty"`
	Packages   []PeelingPackage  `json:"packages,omitempty"`
	Provenance Provenance        `json:"provenance"`
}
