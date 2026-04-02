package internalcmd

import (
	"context"
	"fmt"
	"io"

	"github.com/dantte-lp/goreveal/schema"
	storagediff "github.com/dantte-lp/goreveal/storage/diff"
	storesqlite "github.com/dantte-lp/goreveal/storage/sqlite"
)

type storedRunDiff struct {
	LeftID  int64               `json:"left_id"`
	RightID int64               `json:"right_id"`
	Summary storagediff.Summary `json:"summary"`
}

type storedRunDiffReview struct {
	LeftID         int64                               `json:"left_id"`
	RightID        int64                               `json:"right_id"`
	LeftInput      schema.Input                        `json:"left_input"`
	RightInput     schema.Input                        `json:"right_input"`
	TransferReview *storagediff.TransferReviewQueue    `json:"transfer_review,omitempty"`
	Packages       []storagediff.TransferReviewPackage `json:"transfer_review_packages,omitempty"`
	Focus          *storagediff.TransferReviewAction   `json:"transfer_review_focus,omitempty"`
	Handoff        *diffReviewHandoff                  `json:"handoff,omitempty"`
}

type storedRunDiffHandoff struct {
	LeftID     int64                             `json:"left_id"`
	RightID    int64                             `json:"right_id"`
	LeftInput  schema.Input                      `json:"left_input"`
	RightInput schema.Input                      `json:"right_input"`
	Focus      *storagediff.TransferReviewAction `json:"transfer_review_focus,omitempty"`
	Handoff    *diffReviewHandoff                `json:"handoff,omitempty"`
}

type storedRunDiffNext struct {
	LeftID         int64                              `json:"left_id"`
	RightID        int64                              `json:"right_id"`
	LeftInput      schema.Input                       `json:"left_input"`
	RightInput     schema.Input                       `json:"right_input"`
	ReviewPlan     []storagediff.TransferReviewAction `json:"transfer_review_plan,omitempty"`
	Focus          *storagediff.TransferReviewAction  `json:"transfer_review_focus,omitempty"`
	UpNext         *storagediff.TransferReviewAction  `json:"up_next,omitempty"`
	Upcoming       []storedRunDiffNextPackage         `json:"upcoming_packages,omitempty"`
	Progress       *storedRunDiffNextProgress         `json:"review_progress,omitempty"`
	Recommended    []storedRunDiffNextAction          `json:"recommended_actions,omitempty"`
	Checklist      []string                           `json:"review_checklist,omitempty"`
	Snapshot       *storedRunDiffNextSnapshot         `json:"review_snapshot,omitempty"`
	ReviewCommand  string                             `json:"review_command,omitempty"`
	HandoffCommand string                             `json:"handoff_command,omitempty"`
}

type storedRunDiffNextPackage struct {
	Package            string                          `json:"package,omitempty"`
	ImportPath         string                          `json:"import_path,omitempty"`
	ReviewCount        int                             `json:"review_count,omitempty"`
	ItemCount          int                             `json:"item_count,omitempty"`
	HighestMatchScore  int                             `json:"highest_match_score,omitempty"`
	HighestMatchReason storagediff.FunctionMatchReason `json:"highest_match_reason,omitempty"`
	SampleLeftName     string                          `json:"sample_left_name,omitempty"`
	SampleRightName    string                          `json:"sample_right_name,omitempty"`
}

type storedRunDiffNextProgress struct {
	CurrentStep              int    `json:"current_step"`
	TotalSteps               int    `json:"total_steps"`
	CurrentPackage           string `json:"current_package,omitempty"`
	CurrentImportPath        string `json:"current_import_path,omitempty"`
	CurrentItemCount         int    `json:"current_item_count,omitempty"`
	RemainingPackageCount    int    `json:"remaining_package_count,omitempty"`
	RemainingReviewItemCount int    `json:"remaining_review_item_count,omitempty"`
	NextPackage              string `json:"next_package,omitempty"`
	NextImportPath           string `json:"next_import_path,omitempty"`
	NextItemCount            int    `json:"next_item_count,omitempty"`
}

type storedRunDiffNextAction struct {
	Kind        string `json:"kind"`
	Target      string `json:"target,omitempty"`
	Package     string `json:"package,omitempty"`
	ImportPath  string `json:"import_path,omitempty"`
	Command     string `json:"command"`
	Description string `json:"description,omitempty"`
}

type storedRunDiffNextSnapshot struct {
	CurrentPackage            string                          `json:"current_package,omitempty"`
	CurrentImportPath         string                          `json:"current_import_path,omitempty"`
	CurrentItemCount          int                             `json:"current_item_count,omitempty"`
	CurrentHighestMatchScore  int                             `json:"current_highest_match_score,omitempty"`
	CurrentHighestMatchReason storagediff.FunctionMatchReason `json:"current_highest_match_reason,omitempty"`
	NextPackage               string                          `json:"next_package,omitempty"`
	NextImportPath            string                          `json:"next_import_path,omitempty"`
	RemainingReviewItemCount  int                             `json:"remaining_review_item_count,omitempty"`
	RecommendedActionCount    int                             `json:"recommended_action_count,omitempty"`
}

type diffReviewHandoff struct {
	HandoffContract       string                      `json:"handoff_contract"`
	ArtifactRole          string                      `json:"artifact_role"`
	Mode                  string                      `json:"mode"`
	RecommendedPath       string                      `json:"recommended_path,omitempty"`
	Package               string                      `json:"package,omitempty"`
	ImportPath            string                      `json:"import_path,omitempty"`
	RecommendedTargets    []string                    `json:"recommended_targets,omitempty"`
	RecommendedMCPServers []string                    `json:"recommended_mcp_servers,omitempty"`
	RecommendedExports    []string                    `json:"recommended_exports,omitempty"`
	Artifacts             []diffReviewHandoffArtifact `json:"artifacts,omitempty"`
	ReviewCommand         string                      `json:"review_command,omitempty"`
	ExportCommands        []string                    `json:"export_commands,omitempty"`
	OperatorSteps         []string                    `json:"operator_steps,omitempty"`
	TargetProfiles        []diffReviewHandoffTarget   `json:"target_profiles,omitempty"`
}

type diffReviewHandoffArtifact struct {
	ID              string `json:"id"`
	ArtifactRole    string `json:"artifact_role"`
	Contract        string `json:"contract"`
	Format          string `json:"format"`
	ProducerCommand string `json:"producer_command,omitempty"`
}

type diffReviewHandoffTarget struct {
	Target             string   `json:"target"`
	RecommendedMCP     string   `json:"recommended_mcp_server,omitempty"`
	ExportFormat       string   `json:"export_format"`
	ExportContract     string   `json:"export_contract"`
	ExportCommand      string   `json:"export_command"`
	ArtifactRole       string   `json:"artifact_role"`
	BindingMode        string   `json:"binding_mode,omitempty"`
	HostEntrypoint     string   `json:"host_entrypoint,omitempty"`
	ImportMode         string   `json:"import_mode"`
	PreferredTransport string   `json:"preferred_transport,omitempty"`
	WorkspacePhase     string   `json:"workspace_phase,omitempty"`
	WorkspaceAction    string   `json:"workspace_action"`
	ExpectedHostResult string   `json:"expected_host_result,omitempty"`
	CompletionSignal   string   `json:"completion_signal,omitempty"`
	HostActions        []string `json:"host_actions,omitempty"`
	RequiredArtifacts  []string `json:"required_artifacts,omitempty"`
	OperatorSteps      []string `json:"operator_steps,omitempty"`
}

// RunDiffSQLite compares two stored analyses from the same SQLite database.
func RunDiffSQLite(ctx context.Context, stdout io.Writer, dbPath string, leftID, rightID int64) (err error) {
	_, _, summary, err := loadStoredRunDiffSummary(ctx, dbPath, leftID, rightID)
	if err != nil {
		return err
	}

	return writeJSON(stdout, storedRunDiff{
		LeftID:  leftID,
		RightID: rightID,
		Summary: summary,
	})
}

// RunDiffReviewSQLite writes the compact review-oriented projection over the current diff state.
func RunDiffReviewSQLite(ctx context.Context, stdout io.Writer, dbPath string, leftID, rightID int64) error {
	left, right, summary, err := loadStoredRunDiffSummary(ctx, dbPath, leftID, rightID)
	if err != nil {
		return err
	}

	return writeJSON(stdout, storedRunDiffReview{
		LeftID:         leftID,
		RightID:        rightID,
		LeftInput:      left.Input,
		RightInput:     right.Input,
		TransferReview: summary.TransferReview,
		Packages:       summary.TransferReviewPackages,
		Focus:          summary.TransferReviewFocus,
		Handoff:        buildDiffReviewHandoff(dbPath, leftID, rightID, left.Input, right.Input, summary.TransferReviewFocus),
	})
}

// RunDiffHandoffSQLite writes the compact workstation/MCP handoff projection over the current review state.
func RunDiffHandoffSQLite(ctx context.Context, stdout io.Writer, dbPath string, leftID, rightID int64) error {
	left, right, summary, err := loadStoredRunDiffSummary(ctx, dbPath, leftID, rightID)
	if err != nil {
		return err
	}

	return writeJSON(stdout, storedRunDiffHandoff{
		LeftID:     leftID,
		RightID:    rightID,
		LeftInput:  left.Input,
		RightInput: right.Input,
		Focus:      summary.TransferReviewFocus,
		Handoff:    buildDiffReviewHandoff(dbPath, leftID, rightID, left.Input, right.Input, summary.TransferReviewFocus),
	})
}

// RunDiffNextSQLite writes the compact next-step review projection over the current review state.
func RunDiffNextSQLite(ctx context.Context, stdout io.Writer, dbPath string, leftID, rightID int64) error {
	left, right, summary, err := loadStoredRunDiffSummary(ctx, dbPath, leftID, rightID)
	if err != nil {
		return err
	}

	upNext := buildDiffNextUpNext(summary.TransferReviewPlan, summary.TransferReviewFocus)
	upcoming := buildDiffNextUpcoming(summary.TransferReviewPlan, summary.TransferReviewFocus, 3)
	progress := buildDiffNextProgress(summary.TransferReviewPlan, summary.TransferReviewFocus)
	recommended := buildDiffNextActions(dbPath, leftID, rightID, right.Input, summary.TransferReviewFocus)

	return writeJSON(stdout, storedRunDiffNext{
		LeftID:         leftID,
		RightID:        rightID,
		LeftInput:      left.Input,
		RightInput:     right.Input,
		ReviewPlan:     summary.TransferReviewPlan,
		Focus:          summary.TransferReviewFocus,
		UpNext:         upNext,
		Upcoming:       upcoming,
		Progress:       progress,
		Recommended:    recommended,
		Checklist:      buildDiffNextChecklist(summary.TransferReviewFocus, upNext),
		Snapshot:       buildDiffNextSnapshot(summary.TransferReviewFocus, upNext, progress, recommended),
		ReviewCommand:  fmt.Sprintf("goreveal diff review sqlite %s %d %d", dbPath, leftID, rightID),
		HandoffCommand: fmt.Sprintf("goreveal diff handoff sqlite %s %d %d", dbPath, leftID, rightID),
	})
}

func loadStoredRunDiffSummary(ctx context.Context, dbPath string, leftID, rightID int64) (left schema.Analysis, right schema.Analysis, summary storagediff.Summary, err error) {
	store, err := storesqlite.Open(ctx, dbPath)
	if err != nil {
		return schema.Analysis{}, schema.Analysis{}, storagediff.Summary{}, fmt.Errorf("open sqlite store %q: %w", dbPath, err)
	}
	defer func() {
		if closeErr := store.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close sqlite store %q: %w", dbPath, closeErr)
		}
	}()

	left, err = store.LoadAnalysis(ctx, leftID)
	if err != nil {
		return schema.Analysis{}, schema.Analysis{}, storagediff.Summary{}, fmt.Errorf("load left analysis %d: %w", leftID, err)
	}
	right, err = store.LoadAnalysis(ctx, rightID)
	if err != nil {
		return schema.Analysis{}, schema.Analysis{}, storagediff.Summary{}, fmt.Errorf("load right analysis %d: %w", rightID, err)
	}

	return left, right, storagediff.Compare(left, right), nil
}

func buildDiffReviewHandoff(dbPath string, leftID, rightID int64, leftInput, rightInput schema.Input, focus *storagediff.TransferReviewAction) *diffReviewHandoff {
	if focus == nil {
		return nil
	}

	reviewCommand := fmt.Sprintf("goreveal diff review sqlite %s %d %d", dbPath, leftID, rightID)

	return &diffReviewHandoff{
		HandoffContract:       "goreveal.review_handoff/v1",
		ArtifactRole:          "review_handoff",
		Mode:                  "host_platform_review",
		RecommendedPath:       "export_then_import",
		Package:               focus.Package,
		ImportPath:            focus.ImportPath,
		RecommendedTargets:    []string{"ida", "ghidra"},
		RecommendedMCPServers: []string{"ida-pro-mcp"},
		RecommendedExports:    []string{"ida", "ghidra"},
		Artifacts: []diffReviewHandoffArtifact{
			{
				ID:           "review_handoff",
				ArtifactRole: "review_handoff",
				Contract:     "goreveal.review_handoff/v1",
				Format:       "json",
			},
			{
				ID:              "ida_export",
				ArtifactRole:    "go_metadata_export",
				Contract:        "goreveal.export.ida/v1",
				Format:          "ida",
				ProducerCommand: fmt.Sprintf("goreveal export ida %s", rightInput.Path),
			},
			{
				ID:              "ghidra_export",
				ArtifactRole:    "go_metadata_export",
				Contract:        "goreveal.export.ghidra/v1",
				Format:          "ghidra",
				ProducerCommand: fmt.Sprintf("goreveal export ghidra %s", rightInput.Path),
			},
		},
		ReviewCommand: reviewCommand,
		ExportCommands: []string{
			fmt.Sprintf("goreveal export ida %s", rightInput.Path),
			fmt.Sprintf("goreveal export ghidra %s", rightInput.Path),
		},
		OperatorSteps: []string{
			reviewCommand,
			fmt.Sprintf("goreveal export ida %s", rightInput.Path),
			fmt.Sprintf("goreveal export ghidra %s", rightInput.Path),
			fmt.Sprintf("handoff runtime/package review for %s from %s to host platform MCP or workspace import", focus.Package, leftInput.Path),
		},
		TargetProfiles: []diffReviewHandoffTarget{
			{
				Target:             "ida",
				RecommendedMCP:     "ida-pro-mcp",
				ExportFormat:       "ida",
				ExportContract:     "goreveal.export.ida/v1",
				ExportCommand:      fmt.Sprintf("goreveal export ida %s", rightInput.Path),
				ArtifactRole:       "go_metadata_export",
				BindingMode:        "mcp_server",
				HostEntrypoint:     "ida-pro-mcp.import_export_payload",
				ImportMode:         "mcp_or_workspace_import",
				PreferredTransport: "mcp",
				WorkspacePhase:     "import_then_annotate",
				WorkspaceAction:    "apply_go_specific_annotations",
				ExpectedHostResult: "annotated_workspace_review_ready",
				CompletionSignal:   "names_comments_and_runtime_context_applied",
				HostActions: []string{
					"import_export_payload",
					"apply_names_and_comments",
					"review_runtime_and_package_context",
				},
				RequiredArtifacts: []string{"review_handoff", "ida_export"},
				OperatorSteps: []string{
					reviewCommand,
					fmt.Sprintf("goreveal export ida %s", rightInput.Path),
					"pass the export payload to ida-pro-mcp or import it into the IDA workspace",
				},
			},
			{
				Target:             "ghidra",
				ExportFormat:       "ghidra",
				ExportContract:     "goreveal.export.ghidra/v1",
				ExportCommand:      fmt.Sprintf("goreveal export ghidra %s", rightInput.Path),
				ArtifactRole:       "go_metadata_export",
				BindingMode:        "workspace_loader",
				HostEntrypoint:     "ghidra.workspace_import",
				ImportMode:         "workspace_import",
				PreferredTransport: "file_or_workspace_import",
				WorkspacePhase:     "import_then_annotate",
				WorkspaceAction:    "apply_go_specific_annotations",
				ExpectedHostResult: "annotated_workspace_review_ready",
				CompletionSignal:   "names_comments_and_runtime_context_applied",
				HostActions: []string{
					"import_export_payload",
					"apply_names_and_comments",
					"review_runtime_and_package_context",
				},
				RequiredArtifacts: []string{"review_handoff", "ghidra_export"},
				OperatorSteps: []string{
					reviewCommand,
					fmt.Sprintf("goreveal export ghidra %s", rightInput.Path),
					"import the export payload into the Ghidra workspace or future host-platform bridge",
				},
			},
		},
	}
}

func buildDiffNextActions(dbPath string, leftID, rightID int64, rightInput schema.Input, focus *storagediff.TransferReviewAction) []storedRunDiffNextAction {
	if focus == nil {
		return nil
	}

	reviewCommand := fmt.Sprintf("goreveal diff review sqlite %s %d %d", dbPath, leftID, rightID)
	handoffCommand := fmt.Sprintf("goreveal diff handoff sqlite %s %d %d", dbPath, leftID, rightID)
	idaExport := fmt.Sprintf("goreveal export ida %s", rightInput.Path)
	ghidraExport := fmt.Sprintf("goreveal export ghidra %s", rightInput.Path)

	return []storedRunDiffNextAction{
		{
			Kind:        "review_bundle",
			Package:     focus.Package,
			ImportPath:  focus.ImportPath,
			Command:     reviewCommand,
			Description: "review the current package bundle against the focused transfer items",
		},
		{
			Kind:        "handoff_bundle",
			Package:     focus.Package,
			ImportPath:  focus.ImportPath,
			Command:     handoffCommand,
			Description: "emit the workstation handoff artifact for the current package bundle",
		},
		{
			Kind:        "export_target",
			Target:      "ida",
			Package:     focus.Package,
			ImportPath:  focus.ImportPath,
			Command:     idaExport,
			Description: "export the right-hand analysis for IDA-oriented review",
		},
		{
			Kind:        "export_target",
			Target:      "ghidra",
			Package:     focus.Package,
			ImportPath:  focus.ImportPath,
			Command:     ghidraExport,
			Description: "export the right-hand analysis for Ghidra-oriented review",
		},
	}
}

func buildDiffNextProgress(plan []storagediff.TransferReviewAction, focus *storagediff.TransferReviewAction) *storedRunDiffNextProgress {
	if len(plan) == 0 || focus == nil {
		return nil
	}

	currentIndex := -1
	for i, item := range plan {
		if item.Package == focus.Package && item.ImportPath == focus.ImportPath {
			currentIndex = i
			break
		}
	}
	if currentIndex < 0 {
		currentIndex = 0
	}

	progress := &storedRunDiffNextProgress{
		CurrentStep:       currentIndex + 1,
		TotalSteps:        len(plan),
		CurrentPackage:    focus.Package,
		CurrentImportPath: focus.ImportPath,
		CurrentItemCount:  focus.ItemCount,
	}

	if currentIndex+1 < len(plan) {
		progress.RemainingPackageCount = len(plan) - (currentIndex + 1)
		for _, item := range plan[currentIndex+1:] {
			progress.RemainingReviewItemCount += item.ItemCount
		}
		next := plan[currentIndex+1]
		progress.NextPackage = next.Package
		progress.NextImportPath = next.ImportPath
		progress.NextItemCount = next.ItemCount
	}

	return progress
}

func buildDiffNextUpNext(plan []storagediff.TransferReviewAction, focus *storagediff.TransferReviewAction) *storagediff.TransferReviewAction {
	if len(plan) == 0 || focus == nil {
		return nil
	}

	for i, item := range plan {
		if item.Package != focus.Package || item.ImportPath != focus.ImportPath {
			continue
		}
		if i+1 >= len(plan) {
			return nil
		}
		next := plan[i+1]
		nextCopy := next
		return &nextCopy
	}

	if len(plan) > 1 {
		next := plan[1]
		nextCopy := next
		return &nextCopy
	}

	return nil
}

func buildDiffNextUpcoming(plan []storagediff.TransferReviewAction, focus *storagediff.TransferReviewAction, limit int) []storedRunDiffNextPackage {
	if len(plan) == 0 || focus == nil || limit <= 0 {
		return nil
	}

	start := -1
	for i, item := range plan {
		if item.Package == focus.Package && item.ImportPath == focus.ImportPath {
			start = i + 1
			break
		}
	}
	if start < 0 {
		start = 1
	}
	if start >= len(plan) {
		return nil
	}

	end := start + limit
	if end > len(plan) {
		end = len(plan)
	}

	out := make([]storedRunDiffNextPackage, 0, end-start)
	for _, item := range plan[start:end] {
		out = append(out, storedRunDiffNextPackage{
			Package:            item.Package,
			ImportPath:         item.ImportPath,
			ReviewCount:        item.ReviewCount,
			ItemCount:          item.ItemCount,
			HighestMatchScore:  item.HighestMatchScore,
			HighestMatchReason: item.HighestMatchReason,
			SampleLeftName:     item.SampleLeftName,
			SampleRightName:    item.SampleRightName,
		})
	}
	return out
}

func buildDiffNextChecklist(focus, upNext *storagediff.TransferReviewAction) []string {
	if focus == nil {
		return nil
	}

	checklist := []string{
		fmt.Sprintf("review all %d pending transfer items for %s", focus.ItemCount, focus.Package),
		fmt.Sprintf("confirm the strongest pending match reason for %s remains %s", focus.Package, focus.HighestMatchReason),
		fmt.Sprintf("emit or reuse the handoff bundle for %s before host-platform review", focus.Package),
	}
	if upNext != nil {
		checklist = append(checklist, fmt.Sprintf("after %s, continue with %s", focus.Package, upNext.Package))
	}

	return checklist
}

func buildDiffNextSnapshot(
	focus, upNext *storagediff.TransferReviewAction,
	progress *storedRunDiffNextProgress,
	recommended []storedRunDiffNextAction,
) *storedRunDiffNextSnapshot {
	if focus == nil {
		return nil
	}

	snapshot := &storedRunDiffNextSnapshot{
		CurrentPackage:            focus.Package,
		CurrentImportPath:         focus.ImportPath,
		CurrentItemCount:          focus.ItemCount,
		CurrentHighestMatchScore:  focus.HighestMatchScore,
		CurrentHighestMatchReason: focus.HighestMatchReason,
		RecommendedActionCount:    len(recommended),
	}
	if upNext != nil {
		snapshot.NextPackage = upNext.Package
		snapshot.NextImportPath = upNext.ImportPath
	}
	if progress != nil {
		snapshot.RemainingReviewItemCount = progress.RemainingReviewItemCount
	}

	return snapshot
}
