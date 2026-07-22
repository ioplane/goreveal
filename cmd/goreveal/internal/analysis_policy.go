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
		requirements: availableStages(
			schema.AnalysisStageFunctions,
			schema.AnalysisStagePackages,
			schema.AnalysisStageTypes,
			schema.AnalysisStageStrings,
			schema.AnalysisStageRefinement,
		),
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

	for _, requirement := range policy.requirements {
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
