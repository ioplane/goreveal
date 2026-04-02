package diff

import (
	"slices"
	"strings"

	"github.com/dantte-lp/goreveal/schema"
)

// FunctionMatchReason explains the bounded reason for a version-tracking-adjacent function match.
type FunctionMatchReason string

const (
	FunctionMatchReasonExactName                 FunctionMatchReason = "exact_name"
	FunctionMatchReasonSourceLocation            FunctionMatchReason = "source_location"
	FunctionMatchReasonSourceFile                FunctionMatchReason = "source_file"
	FunctionMatchReasonModuleLocalNormalizedName FunctionMatchReason = "module_local_normalized_name"
)

// FunctionMatch is a bounded similarity record between two functions across stored analyses.
type FunctionMatch struct {
	LeftName            string              `json:"left_name"`
	RightName           string              `json:"right_name"`
	Score               int                 `json:"score"`
	Reason              FunctionMatchReason `json:"reason"`
	LeftImportPath      string              `json:"left_import_path,omitempty"`
	RightImportPath     string              `json:"right_import_path,omitempty"`
	LeftSourceFile      string              `json:"left_source_file,omitempty"`
	RightSourceFile     string              `json:"right_source_file,omitempty"`
	LeftSourceLine      int                 `json:"left_source_line,omitempty"`
	RightSourceLine     int                 `json:"right_source_line,omitempty"`
	LeftClassification  schema.PeelingClass `json:"left_classification,omitempty"`
	RightClassification schema.PeelingClass `json:"right_classification,omitempty"`
}

// TransferDisposition captures how directly a bounded match can be treated as a transfer preview.
type TransferDisposition string

const (
	TransferDispositionReady  TransferDisposition = "ready"
	TransferDispositionReview TransferDisposition = "review"
)

// TransferCandidate is a bounded left-to-right truth projection preview built over existing function matches.
type TransferCandidate struct {
	LeftName                      string                 `json:"left_name"`
	RightName                     string                 `json:"right_name"`
	MatchScore                    int                    `json:"match_score"`
	MatchReason                   FunctionMatchReason    `json:"match_reason"`
	Disposition                   TransferDisposition    `json:"disposition"`
	ProjectedName                 string                 `json:"projected_name"`
	ProjectedPackage              string                 `json:"projected_package,omitempty"`
	ProjectedImportPath           string                 `json:"projected_import_path,omitempty"`
	ProjectedSourceFile           string                 `json:"projected_source_file,omitempty"`
	ProjectedSourceLine           int                    `json:"projected_source_line,omitempty"`
	ProjectedModuleLocal          bool                   `json:"projected_module_local,omitempty"`
	ProjectedClassification       schema.PeelingClass    `json:"projected_classification,omitempty"`
	ProjectedClassificationReason schema.PeelingEvidence `json:"projected_classification_evidence,omitempty"`
}

// AcceptedTransfer is a bounded transfer candidate that the current rules accept without manual review.
type AcceptedTransfer struct {
	LeftName                      string                 `json:"left_name"`
	RightName                     string                 `json:"right_name"`
	AcceptedBy                    FunctionMatchReason    `json:"accepted_by"`
	ProjectedName                 string                 `json:"projected_name"`
	ProjectedPackage              string                 `json:"projected_package,omitempty"`
	ProjectedImportPath           string                 `json:"projected_import_path,omitempty"`
	ProjectedSourceFile           string                 `json:"projected_source_file,omitempty"`
	ProjectedSourceLine           int                    `json:"projected_source_line,omitempty"`
	ProjectedModuleLocal          bool                   `json:"projected_module_local,omitempty"`
	ProjectedClassification       schema.PeelingClass    `json:"projected_classification,omitempty"`
	ProjectedClassificationReason schema.PeelingEvidence `json:"projected_classification_evidence,omitempty"`
}

// TransferReviewItem is one pending human-review item derived from existing transfer candidates.
type TransferReviewItem struct {
	Package                      string                 `json:"package,omitempty"`
	ImportPath                   string                 `json:"import_path,omitempty"`
	LeftName                     string                 `json:"left_name"`
	RightName                    string                 `json:"right_name"`
	MatchScore                   int                    `json:"match_score"`
	MatchReason                  FunctionMatchReason    `json:"match_reason"`
	ProjectedName                string                 `json:"projected_name"`
	ProjectedSourceFile          string                 `json:"projected_source_file,omitempty"`
	ProjectedSourceLine          int                    `json:"projected_source_line,omitempty"`
	ProjectedClassification      schema.PeelingClass    `json:"projected_classification,omitempty"`
	ProjectedClassificationCause schema.PeelingEvidence `json:"projected_classification_evidence,omitempty"`
}

// TransferReviewQueue is a compact analyst-facing review surface built over existing transfer state.
type TransferReviewQueue struct {
	ReviewCount        int                  `json:"review_count"`
	AutoAcceptedCount  int                  `json:"auto_accepted_count"`
	ReviewPackageCount int                  `json:"review_package_count"`
	Items              []TransferReviewItem `json:"items,omitempty"`
}

// TransferReviewPackage is a compact package-first triage summary over pending review items.
type TransferReviewPackage struct {
	Name               string              `json:"name"`
	ImportPath         string              `json:"import_path,omitempty"`
	Classification     schema.PeelingClass `json:"classification,omitempty"`
	ReviewCount        int                 `json:"review_count"`
	AutoAcceptedCount  int                 `json:"auto_accepted_count,omitempty"`
	HighestMatchScore  int                 `json:"highest_match_score,omitempty"`
	HighestMatchReason FunctionMatchReason `json:"highest_match_reason,omitempty"`
	SampleLeftName     string              `json:"sample_left_name,omitempty"`
	SampleRightName    string              `json:"sample_right_name,omitempty"`
}

// TransferReviewAction is a compact recommended next operator step over the current review queue.
type TransferReviewAction struct {
	Action             string               `json:"action"`
	Package            string               `json:"package,omitempty"`
	ImportPath         string               `json:"import_path,omitempty"`
	Classification     schema.PeelingClass  `json:"classification,omitempty"`
	ReviewCount        int                  `json:"review_count"`
	AutoAcceptedCount  int                  `json:"auto_accepted_count,omitempty"`
	HighestMatchScore  int                  `json:"highest_match_score,omitempty"`
	HighestMatchReason FunctionMatchReason  `json:"highest_match_reason,omitempty"`
	SampleLeftName     string               `json:"sample_left_name,omitempty"`
	SampleRightName    string               `json:"sample_right_name,omitempty"`
	ItemCount          int                  `json:"item_count,omitempty"`
	Items              []TransferReviewItem `json:"items,omitempty"`
}

// TransferPackageSummary aggregates current transfer workflow state at the package boundary.
type TransferPackageSummary struct {
	Name                   string              `json:"name"`
	ImportPath             string              `json:"import_path,omitempty"`
	ModuleLocal            bool                `json:"module_local,omitempty"`
	Classification         schema.PeelingClass `json:"classification,omitempty"`
	CandidateCount         int                 `json:"candidate_count"`
	ReadyCount             int                 `json:"ready_count,omitempty"`
	ReviewCount            int                 `json:"review_count,omitempty"`
	AcceptedCount          int                 `json:"accepted_count,omitempty"`
	HighestAcceptedBy      FunctionMatchReason `json:"highest_accepted_by,omitempty"`
	HighestCandidateScore  int                 `json:"highest_candidate_score,omitempty"`
	HighestCandidateReason FunctionMatchReason `json:"highest_candidate_reason,omitempty"`
}

// Counts summarizes cardinalities at the schema boundary.
type Counts struct {
	Functions int `json:"functions"`
	Packages  int `json:"packages"`
	Types     int `json:"types"`
	Strings   int `json:"strings"`
	Files     int `json:"files"`
}

// Summary is the stable v1 diff payload between two stored analyses.
type Summary struct {
	BuildInfoChanged       bool                     `json:"build_info_changed"`
	LeftCounts             Counts                   `json:"left_counts"`
	RightCounts            Counts                   `json:"right_counts"`
	AddedFunctions         []string                 `json:"added_functions,omitempty"`
	RemovedFunctions       []string                 `json:"removed_functions,omitempty"`
	MatchedFunctions       []FunctionMatch          `json:"matched_functions,omitempty"`
	TransferCandidates     []TransferCandidate      `json:"transfer_candidates,omitempty"`
	AcceptedTransfers      []AcceptedTransfer       `json:"accepted_transfers,omitempty"`
	TransferReview         *TransferReviewQueue     `json:"transfer_review,omitempty"`
	TransferReviewPackages []TransferReviewPackage  `json:"transfer_review_packages,omitempty"`
	TransferReviewPlan     []TransferReviewAction   `json:"transfer_review_plan,omitempty"`
	TransferReviewFocus    *TransferReviewAction    `json:"transfer_review_focus,omitempty"`
	TransferPackages       []TransferPackageSummary `json:"transfer_packages,omitempty"`
	AddedPackages          []string                 `json:"added_packages,omitempty"`
	RemovedPackages        []string                 `json:"removed_packages,omitempty"`
	AddedTypes             []string                 `json:"added_types,omitempty"`
	RemovedTypes           []string                 `json:"removed_types,omitempty"`
	AddedStrings           []string                 `json:"added_strings,omitempty"`
	RemovedStrings         []string                 `json:"removed_strings,omitempty"`
	AddedSourceFiles       []string                 `json:"added_source_files,omitempty"`
	RemovedSourceFiles     []string                 `json:"removed_source_files,omitempty"`
	AddedRefinedPasses     []string                 `json:"added_refined_passes,omitempty"`
	RemovedRefinedPasses   []string                 `json:"removed_refined_passes,omitempty"`
}

// Compare computes a stable summary diff between two canonical analyses.
func Compare(left, right schema.Analysis) Summary {
	matches := matchFunctions(left, right)
	transferCandidates := buildTransferCandidates(left, right, matches)
	acceptedTransfers := buildAcceptedTransfers(transferCandidates)

	reviewQueue := buildTransferReviewQueue(transferCandidates, acceptedTransfers)
	reviewPackages := buildTransferReviewPackages(transferCandidates, acceptedTransfers)
	reviewPlan := buildTransferReviewPlan(reviewPackages, reviewQueue)
	out := Summary{
		BuildInfoChanged:       buildInfoChanged(left.BuildInfo, right.BuildInfo),
		LeftCounts:             counts(left),
		RightCounts:            counts(right),
		AddedFunctions:         difference(functionNames(right.Functions), matchedRightNames(matches, left.Functions)),
		RemovedFunctions:       difference(functionNames(left.Functions), matchedLeftNames(matches, right.Functions)),
		MatchedFunctions:       matches,
		TransferCandidates:     transferCandidates,
		AcceptedTransfers:      acceptedTransfers,
		TransferReview:         reviewQueue,
		TransferReviewPackages: reviewPackages,
		TransferReviewPlan:     reviewPlan,
		TransferReviewFocus:    buildTransferReviewFocus(reviewPackages, reviewQueue),
		TransferPackages:       buildTransferPackageSummaries(left, transferCandidates, acceptedTransfers),
		AddedPackages:          difference(packageNames(right.Packages), packageNames(left.Packages)),
		RemovedPackages:        difference(packageNames(left.Packages), packageNames(right.Packages)),
		AddedTypes:             difference(typeNames(right.Types), typeNames(left.Types)),
		RemovedTypes:           difference(typeNames(left.Types), typeNames(right.Types)),
		AddedStrings:           difference(stringValues(right.Strings), stringValues(left.Strings)),
		RemovedStrings:         difference(stringValues(left.Strings), stringValues(right.Strings)),
		AddedSourceFiles:       difference(sourceFiles(right.SourceTree), sourceFiles(left.SourceTree)),
		RemovedSourceFiles:     difference(sourceFiles(left.SourceTree), sourceFiles(right.SourceTree)),
		AddedRefinedPasses:     difference(refinedPasses(right.Refined), refinedPasses(left.Refined)),
		RemovedRefinedPasses:   difference(refinedPasses(left.Refined), refinedPasses(right.Refined)),
	}
	return out
}

type functionSourceKey struct {
	packageName string
	sourceFile  string
	sourceLine  int
}

type functionSourceFileKey struct {
	packageName    string
	sourceFile     string
	classification schema.PeelingClass
}

type normalizedNameKey struct {
	name           string
	importPath     string
	classification schema.PeelingClass
}

func matchFunctions(left, right schema.Analysis) []FunctionMatch {
	if len(left.Functions) == 0 || len(right.Functions) == 0 {
		return nil
	}

	leftByName := make(map[string]schema.Function)
	rightByName := make(map[string]schema.Function)
	for _, fn := range left.Functions {
		leftByName[fn.Name] = fn
	}
	for _, fn := range right.Functions {
		rightByName[fn.Name] = fn
	}

	leftPeeling := peelingByFunctionName(left.Peeling)
	rightPeeling := peelingByFunctionName(right.Peeling)

	var matches []FunctionMatch
	matchedLeft := make(map[string]bool)
	matchedRight := make(map[string]bool)

	commonNames := intersection(functionNames(left.Functions), functionNames(right.Functions))
	for _, name := range commonNames {
		leftFn := leftByName[name]
		rightFn := rightByName[name]
		matches = append(matches, newFunctionMatch(leftFn, rightFn, leftPeeling[name], rightPeeling[name], 100, FunctionMatchReasonExactName))
		matchedLeft[name] = true
		matchedRight[name] = true
	}

	rightBySource := make(map[functionSourceKey]schema.Function)
	for _, fn := range right.Functions {
		if matchedRight[fn.Name] {
			continue
		}
		key, ok := sourceKey(fn)
		if !ok {
			continue
		}
		rightBySource[key] = fn
	}

	for _, fn := range left.Functions {
		if matchedLeft[fn.Name] {
			continue
		}
		key, ok := sourceKey(fn)
		if !ok {
			continue
		}
		rightFn, ok := rightBySource[key]
		if !ok || matchedRight[rightFn.Name] {
			continue
		}
		matches = append(matches, newFunctionMatch(fn, rightFn, leftPeeling[fn.Name], rightPeeling[rightFn.Name], 95, FunctionMatchReasonSourceLocation))
		matchedLeft[fn.Name] = true
		matchedRight[rightFn.Name] = true
	}

	rightByNormalizedName := make(map[normalizedNameKey]schema.Function)
	rightNormalizedCollisions := make(map[normalizedNameKey]bool)
	for _, fn := range right.Functions {
		if matchedRight[fn.Name] {
			continue
		}
		rightPeel, ok := rightPeeling[fn.Name]
		if !ok {
			continue
		}
		key, ok := normalizedNameMatchKey(fn, right.BuildInfo, rightPeel)
		if !ok {
			continue
		}
		if _, exists := rightByNormalizedName[key]; exists {
			rightNormalizedCollisions[key] = true
			continue
		}
		rightByNormalizedName[key] = fn
	}

	for _, fn := range left.Functions {
		if matchedLeft[fn.Name] {
			continue
		}
		leftPeel, ok := leftPeeling[fn.Name]
		if !ok {
			continue
		}
		key, ok := normalizedNameMatchKey(fn, left.BuildInfo, leftPeel)
		if !ok || rightNormalizedCollisions[key] {
			continue
		}
		rightFn, ok := rightByNormalizedName[key]
		if !ok || matchedRight[rightFn.Name] {
			continue
		}
		rightPeel := rightPeeling[rightFn.Name]
		matches = append(matches, newFunctionMatch(fn, rightFn, leftPeel, rightPeel, 92, FunctionMatchReasonModuleLocalNormalizedName))
		matchedLeft[fn.Name] = true
		matchedRight[rightFn.Name] = true
	}

	rightBySourceFile := make(map[functionSourceFileKey]schema.Function)
	rightSourceFileCollisions := make(map[functionSourceFileKey]bool)
	for _, fn := range right.Functions {
		if matchedRight[fn.Name] {
			continue
		}
		rightPeel, ok := rightPeeling[fn.Name]
		if !ok || rightPeel.Classification == "" {
			continue
		}
		key, ok := sourceFileKey(fn, rightPeel.Classification)
		if !ok {
			continue
		}
		if _, exists := rightBySourceFile[key]; exists {
			rightSourceFileCollisions[key] = true
			continue
		}
		rightBySourceFile[key] = fn
	}

	for _, fn := range left.Functions {
		if matchedLeft[fn.Name] {
			continue
		}
		leftPeel, ok := leftPeeling[fn.Name]
		if !ok || leftPeel.Classification == "" {
			continue
		}
		key, ok := sourceFileKey(fn, leftPeel.Classification)
		if !ok || rightSourceFileCollisions[key] {
			continue
		}
		rightFn, ok := rightBySourceFile[key]
		if !ok || matchedRight[rightFn.Name] {
			continue
		}
		rightPeel := rightPeeling[rightFn.Name]
		matches = append(matches, newFunctionMatch(fn, rightFn, leftPeel, rightPeel, 90, FunctionMatchReasonSourceFile))
		matchedLeft[fn.Name] = true
		matchedRight[rightFn.Name] = true
	}

	slices.SortFunc(matches, func(a, b FunctionMatch) int {
		if n := compareStrings(a.LeftName, b.LeftName); n != 0 {
			return n
		}
		return compareStrings(a.RightName, b.RightName)
	})

	return matches
}

func sourceKey(fn schema.Function) (functionSourceKey, bool) {
	if fn.SourceFile == "" || fn.SourceLine == 0 {
		return functionSourceKey{}, false
	}

	return functionSourceKey{
		packageName: fn.Package,
		sourceFile:  fn.SourceFile,
		sourceLine:  fn.SourceLine,
	}, true
}

func sourceFileKey(fn schema.Function, classification schema.PeelingClass) (functionSourceFileKey, bool) {
	if fn.SourceFile == "" || fn.Package == "" || classification == "" {
		return functionSourceFileKey{}, false
	}

	return functionSourceFileKey{
		packageName:    fn.Package,
		sourceFile:     fn.SourceFile,
		classification: classification,
	}, true
}

func normalizedNameMatchKey(fn schema.Function, info *schema.BuildInfo, peeling schema.PeelingFunction) (normalizedNameKey, bool) {
	if peeling.Classification != schema.PeelingClassUser || info == nil || info.Path == "" {
		return normalizedNameKey{}, false
	}

	normalizedName := normalizeModuleLocalValue(fn.Name, info.Path)
	normalizedImportPath := normalizeModuleLocalValue(fn.ImportPath, info.Path)
	if normalizedName == fn.Name && normalizedImportPath == fn.ImportPath {
		return normalizedNameKey{}, false
	}
	if normalizedName == "" {
		return normalizedNameKey{}, false
	}

	return normalizedNameKey{
		name:           normalizedName,
		importPath:     normalizedImportPath,
		classification: peeling.Classification,
	}, true
}

func normalizeModuleLocalValue(value, modulePath string) string {
	if value == "" || modulePath == "" {
		return value
	}
	if value == modulePath {
		return ""
	}
	if trimmed, ok := strings.CutPrefix(value, modulePath+"/"); ok {
		return trimmed
	}

	return value
}

func matchedLeftNames(matches []FunctionMatch, right []schema.Function) []string {
	out := functionNames(right)
	for _, match := range matches {
		if !slices.Contains(out, match.LeftName) {
			out = append(out, match.LeftName)
		}
	}
	slices.Sort(out)
	return slices.Compact(out)
}

func matchedRightNames(matches []FunctionMatch, left []schema.Function) []string {
	out := functionNames(left)
	for _, match := range matches {
		if !slices.Contains(out, match.RightName) {
			out = append(out, match.RightName)
		}
	}
	slices.Sort(out)
	return slices.Compact(out)
}

func newFunctionMatch(left, right schema.Function, leftPeel, rightPeel schema.PeelingFunction, score int, reason FunctionMatchReason) FunctionMatch {
	return FunctionMatch{
		LeftName:            left.Name,
		RightName:           right.Name,
		Score:               score,
		Reason:              reason,
		LeftImportPath:      left.ImportPath,
		RightImportPath:     right.ImportPath,
		LeftSourceFile:      left.SourceFile,
		RightSourceFile:     right.SourceFile,
		LeftSourceLine:      left.SourceLine,
		RightSourceLine:     right.SourceLine,
		LeftClassification:  leftPeel.Classification,
		RightClassification: rightPeel.Classification,
	}
}

func buildTransferCandidates(left, right schema.Analysis, matches []FunctionMatch) []TransferCandidate {
	if len(matches) == 0 {
		return nil
	}

	leftByName := make(map[string]schema.Function, len(left.Functions))
	for _, fn := range left.Functions {
		leftByName[fn.Name] = fn
	}
	rightByName := make(map[string]schema.Function, len(right.Functions))
	for _, fn := range right.Functions {
		rightByName[fn.Name] = fn
	}

	leftPeeling := peelingByFunctionName(left.Peeling)
	rightPeeling := peelingByFunctionName(right.Peeling)

	var out []TransferCandidate
	for _, match := range matches {
		leftFn, ok := leftByName[match.LeftName]
		if !ok {
			continue
		}
		if _, ok := rightByName[match.RightName]; !ok {
			continue
		}
		leftPeel, leftOK := leftPeeling[match.LeftName]
		rightPeel, rightOK := rightPeeling[match.RightName]
		if !leftOK || !rightOK {
			continue
		}
		if leftPeel.Classification != schema.PeelingClassUser || rightPeel.Classification != schema.PeelingClassUser {
			continue
		}

		out = append(out, TransferCandidate{
			LeftName:                      match.LeftName,
			RightName:                     match.RightName,
			MatchScore:                    match.Score,
			MatchReason:                   match.Reason,
			Disposition:                   transferDisposition(match.Reason),
			ProjectedName:                 leftFn.Name,
			ProjectedPackage:              transferPackageName(leftFn, left.BuildInfo, TransferCandidate{ProjectedImportPath: leftFn.ImportPath}),
			ProjectedImportPath:           leftFn.ImportPath,
			ProjectedSourceFile:           leftFn.SourceFile,
			ProjectedSourceLine:           leftFn.SourceLine,
			ProjectedModuleLocal:          leftFn.ModuleLocal,
			ProjectedClassification:       leftPeel.Classification,
			ProjectedClassificationReason: leftPeel.ClassificationEvidence,
		})
	}

	slices.SortFunc(out, func(a, b TransferCandidate) int {
		if n := compareStrings(a.LeftName, b.LeftName); n != 0 {
			return n
		}
		return compareStrings(a.RightName, b.RightName)
	})

	return out
}

func transferDisposition(reason FunctionMatchReason) TransferDisposition {
	switch reason {
	case FunctionMatchReasonExactName, FunctionMatchReasonSourceLocation, FunctionMatchReasonModuleLocalNormalizedName:
		return TransferDispositionReady
	case FunctionMatchReasonSourceFile:
		return TransferDispositionReview
	default:
		return TransferDispositionReview
	}
}

func buildAcceptedTransfers(candidates []TransferCandidate) []AcceptedTransfer {
	if len(candidates) == 0 {
		return nil
	}

	var out []AcceptedTransfer
	for _, candidate := range candidates {
		if candidate.Disposition != TransferDispositionReady {
			continue
		}
		out = append(out, AcceptedTransfer{
			LeftName:                      candidate.LeftName,
			RightName:                     candidate.RightName,
			AcceptedBy:                    candidate.MatchReason,
			ProjectedName:                 candidate.ProjectedName,
			ProjectedPackage:              candidate.ProjectedPackage,
			ProjectedImportPath:           candidate.ProjectedImportPath,
			ProjectedSourceFile:           candidate.ProjectedSourceFile,
			ProjectedSourceLine:           candidate.ProjectedSourceLine,
			ProjectedModuleLocal:          candidate.ProjectedModuleLocal,
			ProjectedClassification:       candidate.ProjectedClassification,
			ProjectedClassificationReason: candidate.ProjectedClassificationReason,
		})
	}

	slices.SortFunc(out, func(a, b AcceptedTransfer) int {
		if n := compareStrings(a.LeftName, b.LeftName); n != 0 {
			return n
		}
		return compareStrings(a.RightName, b.RightName)
	})

	return out
}

func buildTransferReviewQueue(candidates []TransferCandidate, accepted []AcceptedTransfer) *TransferReviewQueue {
	if len(candidates) == 0 {
		return nil
	}

	queue := &TransferReviewQueue{
		AutoAcceptedCount: len(accepted),
	}
	reviewPackages := make(map[string]bool)
	for _, candidate := range candidates {
		if candidate.Disposition != TransferDispositionReview {
			continue
		}
		queue.Items = append(queue.Items, TransferReviewItem{
			Package:                      candidate.ProjectedPackage,
			ImportPath:                   candidate.ProjectedImportPath,
			LeftName:                     candidate.LeftName,
			RightName:                    candidate.RightName,
			MatchScore:                   candidate.MatchScore,
			MatchReason:                  candidate.MatchReason,
			ProjectedName:                candidate.ProjectedName,
			ProjectedSourceFile:          candidate.ProjectedSourceFile,
			ProjectedSourceLine:          candidate.ProjectedSourceLine,
			ProjectedClassification:      candidate.ProjectedClassification,
			ProjectedClassificationCause: candidate.ProjectedClassificationReason,
		})
		queue.ReviewCount++
		reviewPackages[candidate.ProjectedPackage+"\x00"+candidate.ProjectedImportPath] = true
	}
	if queue.ReviewCount == 0 {
		return nil
	}
	queue.ReviewPackageCount = len(reviewPackages)
	slices.SortFunc(queue.Items, func(a, b TransferReviewItem) int {
		if n := compareStrings(a.Package, b.Package); n != 0 {
			return n
		}
		if n := compareStrings(a.LeftName, b.LeftName); n != 0 {
			return n
		}
		return compareStrings(a.RightName, b.RightName)
	})
	return queue
}

type transferReviewPackageKey struct {
	name           string
	importPath     string
	classification schema.PeelingClass
}

func buildTransferReviewPackages(candidates []TransferCandidate, accepted []AcceptedTransfer) []TransferReviewPackage {
	if len(candidates) == 0 {
		return nil
	}

	autoAcceptedByPackage := make(map[string]int, len(accepted))
	for _, acceptedTransfer := range accepted {
		key := acceptedTransfer.ProjectedPackage + "\x00" + acceptedTransfer.ProjectedImportPath
		autoAcceptedByPackage[key]++
	}

	summaries := make(map[transferReviewPackageKey]*TransferReviewPackage)
	for _, candidate := range candidates {
		if candidate.Disposition != TransferDispositionReview {
			continue
		}
		key := transferReviewPackageKey{
			name:           candidate.ProjectedPackage,
			importPath:     candidate.ProjectedImportPath,
			classification: candidate.ProjectedClassification,
		}
		summary, ok := summaries[key]
		if !ok {
			summary = &TransferReviewPackage{
				Name:              key.name,
				ImportPath:        key.importPath,
				Classification:    key.classification,
				AutoAcceptedCount: autoAcceptedByPackage[key.name+"\x00"+key.importPath],
			}
			summaries[key] = summary
		}
		summary.ReviewCount++
		if candidate.MatchScore > summary.HighestMatchScore {
			summary.HighestMatchScore = candidate.MatchScore
			summary.HighestMatchReason = candidate.MatchReason
			summary.SampleLeftName = candidate.LeftName
			summary.SampleRightName = candidate.RightName
		}
	}

	if len(summaries) == 0 {
		return nil
	}

	out := make([]TransferReviewPackage, 0, len(summaries))
	for _, summary := range summaries {
		out = append(out, *summary)
	}
	slices.SortFunc(out, func(a, b TransferReviewPackage) int {
		if n := compareStrings(a.ImportPath, b.ImportPath); n != 0 {
			return n
		}
		return compareStrings(a.Name, b.Name)
	})
	return out
}

func buildTransferReviewFocus(reviewPackages []TransferReviewPackage, queue *TransferReviewQueue) *TransferReviewAction {
	if len(reviewPackages) == 0 {
		return nil
	}

	focus := reviewPackages[0]
	for _, candidate := range reviewPackages[1:] {
		if transferReviewPackageLess(focus, candidate) {
			focus = candidate
		}
	}

	action := &TransferReviewAction{
		Action:             "review_package",
		Package:            focus.Name,
		ImportPath:         focus.ImportPath,
		Classification:     focus.Classification,
		ReviewCount:        focus.ReviewCount,
		AutoAcceptedCount:  focus.AutoAcceptedCount,
		HighestMatchScore:  focus.HighestMatchScore,
		HighestMatchReason: focus.HighestMatchReason,
		SampleLeftName:     focus.SampleLeftName,
		SampleRightName:    focus.SampleRightName,
	}
	if queue == nil {
		return action
	}
	for _, item := range queue.Items {
		if item.Package != focus.Name || item.ImportPath != focus.ImportPath {
			continue
		}
		action.Items = append(action.Items, item)
	}
	action.ItemCount = len(action.Items)
	return action
}

func buildTransferReviewPlan(reviewPackages []TransferReviewPackage, queue *TransferReviewQueue) []TransferReviewAction {
	if len(reviewPackages) == 0 {
		return nil
	}

	ordered := append([]TransferReviewPackage(nil), reviewPackages...)
	slices.SortFunc(ordered, func(a, b TransferReviewPackage) int {
		switch {
		case transferReviewPackageLess(a, b):
			return 1
		case transferReviewPackageLess(b, a):
			return -1
		default:
			return 0
		}
	})

	out := make([]TransferReviewAction, 0, len(ordered))
	for _, item := range ordered {
		action := TransferReviewAction{
			Action:             "review_package",
			Package:            item.Name,
			ImportPath:         item.ImportPath,
			Classification:     item.Classification,
			ReviewCount:        item.ReviewCount,
			AutoAcceptedCount:  item.AutoAcceptedCount,
			HighestMatchScore:  item.HighestMatchScore,
			HighestMatchReason: item.HighestMatchReason,
			SampleLeftName:     item.SampleLeftName,
			SampleRightName:    item.SampleRightName,
		}
		if queue != nil {
			for _, reviewItem := range queue.Items {
				if reviewItem.Package != item.Name || reviewItem.ImportPath != item.ImportPath {
					continue
				}
				action.Items = append(action.Items, reviewItem)
			}
		}
		action.ItemCount = len(action.Items)
		out = append(out, action)
	}
	return out
}

func transferReviewPackageLess(left, right TransferReviewPackage) bool {
	if left.ReviewCount != right.ReviewCount {
		return left.ReviewCount < right.ReviewCount
	}
	if left.HighestMatchScore != right.HighestMatchScore {
		return left.HighestMatchScore < right.HighestMatchScore
	}
	if left.AutoAcceptedCount != right.AutoAcceptedCount {
		return left.AutoAcceptedCount < right.AutoAcceptedCount
	}
	if n := compareStrings(left.ImportPath, right.ImportPath); n != 0 {
		return n > 0
	}
	return compareStrings(left.Name, right.Name) > 0
}

type transferPackageKey struct {
	name           string
	importPath     string
	moduleLocal    bool
	classification schema.PeelingClass
}

func buildTransferPackageSummaries(left schema.Analysis, candidates []TransferCandidate, accepted []AcceptedTransfer) []TransferPackageSummary {
	if len(candidates) == 0 {
		return nil
	}

	leftByName := make(map[string]schema.Function, len(left.Functions))
	for _, fn := range left.Functions {
		leftByName[fn.Name] = fn
	}

	summaries := make(map[transferPackageKey]*TransferPackageSummary)
	candidatePackageKeys := make(map[string]transferPackageKey, len(candidates))
	for _, candidate := range candidates {
		leftFn, ok := leftByName[candidate.LeftName]
		if !ok {
			continue
		}
		key := transferPackageSummaryKey(leftFn, left.BuildInfo, candidate)
		candidatePackageKeys[candidate.LeftName+"\x00"+candidate.RightName] = key
		summary := ensureTransferPackageSummary(summaries, key)
		summary.CandidateCount++
		switch candidate.Disposition {
		case TransferDispositionReady:
			summary.ReadyCount++
		case TransferDispositionReview:
			summary.ReviewCount++
		}
		if candidate.MatchScore > summary.HighestCandidateScore {
			summary.HighestCandidateScore = candidate.MatchScore
			summary.HighestCandidateReason = candidate.MatchReason
		}
	}

	for _, acceptedTransfer := range accepted {
		key, ok := candidatePackageKeys[acceptedTransfer.LeftName+"\x00"+acceptedTransfer.RightName]
		if !ok {
			continue
		}
		summary := ensureTransferPackageSummary(summaries, key)
		summary.AcceptedCount++
		if acceptedReasonRank(acceptedTransfer.AcceptedBy) > acceptedReasonRank(summary.HighestAcceptedBy) {
			summary.HighestAcceptedBy = acceptedTransfer.AcceptedBy
		}
	}

	out := make([]TransferPackageSummary, 0, len(summaries))
	for _, summary := range summaries {
		out = append(out, *summary)
	}

	slices.SortFunc(out, func(a, b TransferPackageSummary) int {
		if n := compareStrings(a.ImportPath, b.ImportPath); n != 0 {
			return n
		}
		return compareStrings(a.Name, b.Name)
	})

	return out
}

func transferPackageSummaryKey(leftFn schema.Function, info *schema.BuildInfo, candidate TransferCandidate) transferPackageKey {
	return transferPackageKey{
		name:           transferPackageName(leftFn, info, candidate),
		importPath:     candidate.ProjectedImportPath,
		moduleLocal:    candidate.ProjectedModuleLocal,
		classification: candidate.ProjectedClassification,
	}
}

func transferPackageName(fn schema.Function, info *schema.BuildInfo, candidate TransferCandidate) string {
	if fn.Package != "" {
		return fn.Package
	}
	if info != nil && info.Path != "" && candidate.ProjectedImportPath == info.Path {
		return "main"
	}
	if strings.HasPrefix(fn.Name, "main.") {
		return "main"
	}
	return pathBase(fn.ImportPath)
}

func pathBase(path string) string {
	if path == "" {
		return ""
	}
	if idx := strings.LastIndexByte(path, '/'); idx >= 0 {
		return path[idx+1:]
	}
	return path
}

func ensureTransferPackageSummary(summaries map[transferPackageKey]*TransferPackageSummary, key transferPackageKey) *TransferPackageSummary {
	if summary, ok := summaries[key]; ok {
		return summary
	}
	summary := &TransferPackageSummary{
		Name:           key.name,
		ImportPath:     key.importPath,
		ModuleLocal:    key.moduleLocal,
		Classification: key.classification,
	}
	summaries[key] = summary
	return summary
}

func acceptedReasonRank(reason FunctionMatchReason) int {
	switch reason {
	case FunctionMatchReasonExactName:
		return 4
	case FunctionMatchReasonSourceLocation:
		return 3
	case FunctionMatchReasonModuleLocalNormalizedName:
		return 2
	case FunctionMatchReasonSourceFile:
		return 1
	default:
		return 0
	}
}

func peelingByFunctionName(peeling *schema.PeelingAnalysis) map[string]schema.PeelingFunction {
	if peeling == nil {
		return nil
	}

	out := make(map[string]schema.PeelingFunction, len(peeling.Functions))
	for _, fn := range peeling.Functions {
		out[fn.Name] = fn
	}
	return out
}

func intersection(left, right []string) []string {
	var out []string
	for _, item := range left {
		if slices.Contains(right, item) {
			out = append(out, item)
		}
	}
	return out
}

func compareStrings(left, right string) int {
	switch {
	case left < right:
		return -1
	case left > right:
		return 1
	default:
		return 0
	}
}

func buildInfoChanged(left, right *schema.BuildInfo) bool {
	switch {
	case left == nil && right == nil:
		return false
	case left == nil || right == nil:
		return true
	default:
		return left.GoVersion != right.GoVersion || left.Path != right.Path
	}
}

func counts(analysis schema.Analysis) Counts {
	return Counts{
		Functions: len(analysis.Functions),
		Packages:  len(analysis.Packages),
		Types:     len(analysis.Types),
		Strings:   len(analysis.Strings),
		Files:     len(sourceFiles(analysis.SourceTree)),
	}
}

func functionNames(items []schema.Function) []string {
	names := make([]string, 0, len(items))
	for _, item := range items {
		names = append(names, item.Name)
	}
	slices.Sort(names)
	return slices.Compact(names)
}

func packageNames(items []schema.Package) []string {
	names := make([]string, 0, len(items))
	for _, item := range items {
		names = append(names, item.Name)
	}
	slices.Sort(names)
	return slices.Compact(names)
}

func typeNames(items []schema.Type) []string {
	names := make([]string, 0, len(items))
	for _, item := range items {
		names = append(names, item.Name)
	}
	slices.Sort(names)
	return slices.Compact(names)
}

func stringValues(items []schema.StringCandidate) []string {
	values := make([]string, 0, len(items))
	for _, item := range items {
		values = append(values, item.Value)
	}
	slices.Sort(values)
	return slices.Compact(values)
}

func sourceFiles(tree *schema.SourceTree) []string {
	if tree == nil {
		return nil
	}
	files := append([]string(nil), tree.Files...)
	slices.Sort(files)
	return slices.Compact(files)
}

func refinedPasses(refined *schema.RefinedAnalysis) []string {
	if refined == nil {
		return nil
	}
	passes := append([]string(nil), refined.Passes...)
	slices.Sort(passes)
	return slices.Compact(passes)
}

func difference(left, right []string) []string {
	var out []string
	for _, item := range left {
		if !slices.Contains(right, item) {
			out = append(out, item)
		}
	}
	return out
}
