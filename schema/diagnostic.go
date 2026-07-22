package schema

// AnalysisStage identifies one recovery, derivation, or refinement stage.
type AnalysisStage string

const (
	AnalysisStageBuildInfo  AnalysisStage = "build_info"
	AnalysisStageRuntime    AnalysisStage = "runtime"
	AnalysisStageFunctions  AnalysisStage = "functions"
	AnalysisStagePackages   AnalysisStage = "packages"
	AnalysisStageTypes      AnalysisStage = "types"
	AnalysisStageStrings    AnalysisStage = "strings"
	AnalysisStageSourceTree AnalysisStage = "source_tree"
	AnalysisStagePeeling    AnalysisStage = "peeling"
	AnalysisStageRefinement AnalysisStage = "refinement"
)

// StageStatus records the outcome of one attempted analysis stage.
type StageStatus string

const (
	StageStatusAvailable   StageStatus = "available"
	StageStatusUnavailable StageStatus = "unavailable"
	StageStatusUnsupported StageStatus = "unsupported"
	StageStatusFailed      StageStatus = "failed"
)

// StageDiagnostic makes partial analysis outcomes explicit and machine-readable.
type StageDiagnostic struct {
	Stage   AnalysisStage `json:"stage"`
	Status  StageStatus   `json:"status"`
	Code    string        `json:"code,omitempty"`
	Message string        `json:"message,omitempty"`
}
