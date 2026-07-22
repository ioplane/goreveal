package internalcmd

import (
	"fmt"
	"slices"

	"github.com/dantte-lp/goreveal/schema"
)

type analysisCommand string

const (
	analysisCommandAnalyze          analysisCommand = "analyze"
	analysisCommandInspectFunctions analysisCommand = "inspect_functions"
	analysisCommandInspectPackages  analysisCommand = "inspect_packages"
	analysisCommandInspectTypes     analysisCommand = "inspect_types"
	analysisCommandInspectStrings   analysisCommand = "inspect_strings"
	analysisCommandInspectRuntime   analysisCommand = "inspect_runtime"
	analysisCommandInspectPeeling   analysisCommand = "inspect_peeling"
	analysisCommandPeel             analysisCommand = "peel"
	analysisCommandSourceTree       analysisCommand = "source_tree"
	analysisCommandDeobfuscate      analysisCommand = "deobfuscate"
	analysisCommandExportSQLite     analysisCommand = "export_sqlite"
	analysisCommandExportIDAV1      analysisCommand = "export_ida_v1"
	analysisCommandExportGhidraV1   analysisCommand = "export_ghidra_v1"
)

type stageRequirement struct {
	stage   schema.AnalysisStage
	allowed []schema.StageStatus
}

type commandPolicy struct {
	requirements          []stageRequirement
	rejectFailedAttempted bool
}

var commandPolicies = map[analysisCommand]commandPolicy{
	analysisCommandAnalyze: {},
	analysisCommandInspectFunctions: {
		requirements: availableStages(schema.AnalysisStageFunctions),
	},
	analysisCommandInspectPackages: {
		requirements: availableStages(schema.AnalysisStageFunctions, schema.AnalysisStagePackages),
	},
	analysisCommandInspectTypes: {
		requirements: []stageRequirement{{
			stage:   schema.AnalysisStageTypes,
			allowed: []schema.StageStatus{schema.StageStatusAvailable, schema.StageStatusUnavailable},
		}},
	},
	analysisCommandInspectStrings: {
		requirements: availableStages(schema.AnalysisStageStrings),
	},
	analysisCommandInspectRuntime: {
		requirements: availableStages(schema.AnalysisStageRuntime),
	},
	analysisCommandInspectPeeling: {
		requirements: availableStages(schema.AnalysisStageFunctions, schema.AnalysisStagePeeling),
	},
	analysisCommandPeel: {
		requirements: availableStages(schema.AnalysisStageFunctions, schema.AnalysisStagePeeling),
	},
	analysisCommandSourceTree: {
		requirements: availableStages(schema.AnalysisStageSourceTree),
	},
	analysisCommandDeobfuscate: {
		requirements: availableStages(schema.AnalysisStageRefinement),
	},
	analysisCommandExportSQLite: {
		requirements:          availableStages(schema.AnalysisStageFunctions),
		rejectFailedAttempted: true,
	},
	analysisCommandExportIDAV1: {
		requirements:          availableStages(schema.AnalysisStageFunctions),
		rejectFailedAttempted: true,
	},
	analysisCommandExportGhidraV1: {
		requirements:          availableStages(schema.AnalysisStageFunctions),
		rejectFailedAttempted: true,
	},
}

// analysisPolicyError is stable enough for command tests and machine-facing exit diagnostics.
type analysisPolicyError struct {
	Command analysisCommand
	Stage   schema.AnalysisStage
	Status  schema.StageStatus
	Code    string
	Message string
}

func (e *analysisPolicyError) Error() string {
	if e.Stage == "" {
		return fmt.Sprintf("analysis policy %q: %s", e.Command, e.Message)
	}
	return fmt.Sprintf(
		"analysis policy %q stage %q: status=%q code=%q: %s",
		e.Command,
		e.Stage,
		e.Status,
		e.Code,
		e.Message,
	)
}

func enforceAnalysisPolicy(command analysisCommand, analysis schema.Analysis) error {
	policy, ok := commandPolicies[command]
	if !ok {
		return &analysisPolicyError{
			Command: command,
			Code:    "unknown_analysis_policy",
			Message: "no explicit policy is registered",
		}
	}
	if err := validateDiagnosticVocabulary(command, analysis.Diagnostics); err != nil {
		return err
	}

	byStage := make(map[schema.AnalysisStage]schema.StageDiagnostic, len(analysis.Diagnostics))
	for _, diagnostic := range analysis.Diagnostics {
		if _, duplicate := byStage[diagnostic.Stage]; duplicate {
			return &analysisPolicyError{
				Command: command,
				Stage:   diagnostic.Stage,
				Status:  diagnostic.Status,
				Code:    "duplicate_stage_diagnostic",
				Message: "analysis contains more than one diagnostic for the stage",
			}
		}
		byStage[diagnostic.Stage] = diagnostic

		if policy.rejectFailedAttempted && diagnostic.Status == schema.StageStatusFailed {
			return policyErrorFromDiagnostic(command, diagnostic)
		}
	}

	if err := enforceStageRequirements(command, byStage, policy.requirements); err != nil {
		return err
	}
	if command == analysisCommandDeobfuscate {
		return enforceDeobfuscationFamilyPolicy(command, analysis, byStage)
	}

	return nil
}

func validateDiagnosticVocabulary(command analysisCommand, diagnostics schema.StageDiagnostics) error {
	for _, diagnostic := range diagnostics {
		if !knownAnalysisStage(diagnostic.Stage) {
			return &analysisPolicyError{
				Command: command,
				Stage:   diagnostic.Stage,
				Status:  diagnostic.Status,
				Code:    "unknown_analysis_stage",
				Message: "analysis contains an unknown stage",
			}
		}
		if !knownStageStatus(diagnostic.Status) {
			return &analysisPolicyError{
				Command: command,
				Stage:   diagnostic.Stage,
				Status:  diagnostic.Status,
				Code:    "invalid_stage_status",
				Message: "analysis stage contains an invalid status",
			}
		}
	}

	return nil
}

func knownAnalysisStage(stage schema.AnalysisStage) bool {
	switch stage {
	case schema.AnalysisStageBuildInfo,
		schema.AnalysisStageRuntime,
		schema.AnalysisStageFunctions,
		schema.AnalysisStagePackages,
		schema.AnalysisStageTypes,
		schema.AnalysisStageStrings,
		schema.AnalysisStageSourceTree,
		schema.AnalysisStagePeeling,
		schema.AnalysisStageRefinement:
		return true
	default:
		return false
	}
}

func knownStageStatus(status schema.StageStatus) bool {
	switch status {
	case schema.StageStatusAvailable,
		schema.StageStatusUnavailable,
		schema.StageStatusUnsupported,
		schema.StageStatusFailed:
		return true
	default:
		return false
	}
}

func enforceStageRequirements(
	command analysisCommand,
	byStage map[schema.AnalysisStage]schema.StageDiagnostic,
	requirements []stageRequirement,
) error {
	for _, requirement := range requirements {
		diagnostic, exists := byStage[requirement.stage]
		if !exists {
			return &analysisPolicyError{
				Command: command,
				Stage:   requirement.stage,
				Code:    "stage_not_attempted",
				Message: "required analysis stage was not attempted",
			}
		}
		if !statusAllowed(diagnostic.Status, requirement.allowed) {
			return policyErrorFromDiagnostic(command, diagnostic)
		}
	}

	return nil
}

func enforceDeobfuscationFamilyPolicy(
	command analysisCommand,
	analysis schema.Analysis,
	byStage map[schema.AnalysisStage]schema.StageDiagnostic,
) error {
	if analysis.Refined == nil || !hasRefinedFamilies(*analysis.Refined) {
		return &analysisPolicyError{
			Command: command,
			Stage:   schema.AnalysisStageRefinement,
			Status:  schema.StageStatusAvailable,
			Code:    "refinement_payload_unavailable",
			Message: "refinement stage published no refined families",
		}
	}

	required := make(map[schema.AnalysisStage]bool, 4)
	if len(analysis.Refined.Functions) != 0 {
		required[schema.AnalysisStageFunctions] = true
	}
	if len(analysis.Refined.Packages) != 0 {
		required[schema.AnalysisStageFunctions] = true
		required[schema.AnalysisStagePackages] = true
	}
	if len(analysis.Refined.Types) != 0 {
		required[schema.AnalysisStageTypes] = true
	}
	if len(analysis.Refined.Strings) != 0 {
		required[schema.AnalysisStageStrings] = true
	}

	stages := []schema.AnalysisStage{
		schema.AnalysisStageFunctions,
		schema.AnalysisStagePackages,
		schema.AnalysisStageTypes,
		schema.AnalysisStageStrings,
	}
	requirements := make([]stageRequirement, 0, len(required))
	for _, stage := range stages {
		if required[stage] {
			requirements = append(requirements, availableStages(stage)...)
		}
	}

	return enforceStageRequirements(command, byStage, requirements)
}

func hasRefinedFamilies(refined schema.RefinedAnalysis) bool {
	return len(refined.Functions) != 0 ||
		len(refined.Packages) != 0 ||
		len(refined.Types) != 0 ||
		len(refined.Strings) != 0
}

func policyErrorFromDiagnostic(command analysisCommand, diagnostic schema.StageDiagnostic) error {
	code := diagnostic.Code
	if code == "" {
		code = "stage_status_" + string(diagnostic.Status)
	}
	message := diagnostic.Message
	if message == "" {
		message = "analysis stage does not satisfy command policy"
	}

	return &analysisPolicyError{
		Command: command,
		Stage:   diagnostic.Stage,
		Status:  diagnostic.Status,
		Code:    code,
		Message: message,
	}
}

func availableStages(stages ...schema.AnalysisStage) []stageRequirement {
	requirements := make([]stageRequirement, 0, len(stages))
	for _, stage := range stages {
		requirements = append(requirements, stageRequirement{
			stage:   stage,
			allowed: []schema.StageStatus{schema.StageStatusAvailable},
		})
	}
	return requirements
}

func statusAllowed(status schema.StageStatus, allowed []schema.StageStatus) bool {
	return slices.Contains(allowed, status)
}
