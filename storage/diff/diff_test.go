package diff

import (
	"testing"

	"github.com/dantte-lp/goreveal/schema"
)

func TestCompare(t *testing.T) {
	t.Parallel()

	left := schema.Analysis{
		Input: schema.Input{
			Path:   "/tmp/sample.bin",
			Format: "elf",
		},
		BuildInfo: &schema.BuildInfo{
			GoVersion: "go1.26.1",
			Path:      "example.com/sample",
		},
		Functions: []schema.Function{
			{Name: "main.main", ImportPath: "example.com/sample", SourceFile: "main.go", SourceLine: 10, Entry: 0x1000, End: 0x1100},
			{Name: "main.worker", ImportPath: "example.com/sample", SourceFile: "worker.go", SourceLine: 22, Entry: 0x1200, End: 0x1300},
			{Name: "main.service", ImportPath: "example.com/sample", Package: "main", SourceFile: "service.go", SourceLine: 40, Entry: 0x1400, End: 0x1500},
			{Name: "example.com/sample/internal/app.Handler", ImportPath: "example.com/sample/internal/app", Package: "internal/app", Entry: 0x1600, End: 0x1700},
		},
		Packages: []schema.Package{
			{Name: "main", FunctionCount: 1},
		},
		Types: []schema.Type{
			{Name: "main.counter", Kind: "struct"},
		},
		Strings: []schema.StringCandidate{
			{Value: "hello", Region: ".rodata", Offset: 8},
		},
		SourceTree: &schema.SourceTree{
			Root:  "example.com/sample",
			Files: []string{"main.go"},
		},
		Refined: &schema.RefinedAnalysis{
			Passes: []string{"synthetic-function-names"},
		},
		Peeling: &schema.PeelingAnalysis{
			Functions: []schema.PeelingFunction{
				{Name: "main.main", ImportPath: "example.com/sample", SourceFile: "main.go", SourceLine: 10, Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceModuleLocal},
				{Name: "main.worker", ImportPath: "example.com/sample", SourceFile: "worker.go", SourceLine: 22, Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceModuleLocal},
				{Name: "main.service", Package: "main", ImportPath: "example.com/sample", SourceFile: "service.go", SourceLine: 40, Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceBuildInfoPath},
				{Name: "example.com/sample/internal/app.Handler", Package: "internal/app", ImportPath: "example.com/sample/internal/app", Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceModuleLocal},
			},
		},
	}

	right := schema.Analysis{
		Input: schema.Input{
			Path:   "/tmp/sample.bin",
			Format: "elf",
		},
		BuildInfo: &schema.BuildInfo{
			GoVersion: "go1.26.2",
			Path:      "example.com/sample/v2",
		},
		Functions: []schema.Function{
			{Name: "main.main", ImportPath: "example.com/sample/v2", SourceFile: "main.go", SourceLine: 10, Entry: 0x1000, End: 0x1100},
			{Name: "main.renamedWorker", ImportPath: "example.com/sample/v2", SourceFile: "worker.go", SourceLine: 22, Entry: 0x1200, End: 0x1300},
			{Name: "main.serviceV2", ImportPath: "example.com/sample/v2", Package: "main", SourceFile: "service.go", SourceLine: 48, Entry: 0x1400, End: 0x1500},
			{Name: "example.com/sample/v2/internal/app.Handler", ImportPath: "example.com/sample/v2/internal/app", Package: "internal/app", Entry: 0x1600, End: 0x1700},
			{Name: "main.helper", Entry: 0x1200, End: 0x1220},
		},
		Packages: []schema.Package{
			{Name: "main", FunctionCount: 2},
			{Name: "internal/app", FunctionCount: 1},
		},
		Types: []schema.Type{
			{Name: "main.counter", Kind: "struct"},
			{Name: "main.extra", Kind: "string"},
		},
		Strings: []schema.StringCandidate{
			{Value: "world", Region: ".rodata", Offset: 24},
		},
		SourceTree: &schema.SourceTree{
			Root:  "example.com/sample/v2",
			Files: []string{"main.go", "helper.go"},
		},
		Refined: &schema.RefinedAnalysis{
			Passes: []string{"synthetic-function-names", "string-segments"},
		},
		Peeling: &schema.PeelingAnalysis{
			Functions: []schema.PeelingFunction{
				{Name: "main.main", ImportPath: "example.com/sample/v2", SourceFile: "main.go", SourceLine: 10, Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceModuleLocal},
				{Name: "main.renamedWorker", ImportPath: "example.com/sample/v2", SourceFile: "worker.go", SourceLine: 22, Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceModuleLocal},
				{Name: "main.serviceV2", Package: "main", ImportPath: "example.com/sample/v2", SourceFile: "service.go", SourceLine: 48, Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceBuildInfoPath},
				{Name: "example.com/sample/v2/internal/app.Handler", Package: "internal/app", ImportPath: "example.com/sample/v2/internal/app", Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceModuleLocal},
				{Name: "main.helper", Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceModuleLocal},
			},
		},
	}

	got := Compare(left, right)

	if !got.BuildInfoChanged {
		t.Fatal("BuildInfoChanged = false, want true")
	}
	if got.LeftCounts.Functions != 4 || got.RightCounts.Functions != 5 {
		t.Fatalf("function counts = %#v %#v", got.LeftCounts, got.RightCounts)
	}
	if len(got.AddedFunctions) != 1 || got.AddedFunctions[0] != "main.helper" {
		t.Fatalf("added functions = %#v", got.AddedFunctions)
	}
	if len(got.RemovedFunctions) != 0 {
		t.Fatalf("removed functions = %#v", got.RemovedFunctions)
	}
	if len(got.MatchedFunctions) != 4 {
		t.Fatalf("matched functions = %#v", got.MatchedFunctions)
	}
	if len(got.TransferCandidates) != 4 {
		t.Fatalf("transfer candidates = %#v", got.TransferCandidates)
	}
	if len(got.AcceptedTransfers) != 3 {
		t.Fatalf("accepted transfers = %#v", got.AcceptedTransfers)
	}
	if got.TransferReview == nil {
		t.Fatalf("transfer review = %#v", got.TransferReview)
	}
	if len(got.TransferReviewPackages) != 1 {
		t.Fatalf("transfer review packages = %#v", got.TransferReviewPackages)
	}
	if len(got.TransferReviewPlan) != 1 {
		t.Fatalf("transfer review plan = %#v", got.TransferReviewPlan)
	}
	if got.TransferReviewFocus == nil {
		t.Fatalf("transfer review focus = %#v", got.TransferReviewFocus)
	}
	if len(got.TransferPackages) != 2 {
		t.Fatalf("transfer packages = %#v", got.TransferPackages)
	}
	exactNameMatch, ok := findFunctionMatch(got.MatchedFunctions, "main.main", "main.main")
	if !ok {
		t.Fatalf("missing exact-name match in %#v", got.MatchedFunctions)
	}
	if exactNameMatch.Score != 100 || exactNameMatch.Reason != FunctionMatchReasonExactName {
		t.Fatalf("exact-name match = %#v", exactNameMatch)
	}
	sourceLocationMatch, ok := findFunctionMatch(got.MatchedFunctions, "main.worker", "main.renamedWorker")
	if !ok {
		t.Fatalf("missing source-location match in %#v", got.MatchedFunctions)
	}
	if sourceLocationMatch.Score != 95 || sourceLocationMatch.Reason != FunctionMatchReasonSourceLocation {
		t.Fatalf("source-location match = %#v", sourceLocationMatch)
	}
	if sourceLocationMatch.LeftClassification != schema.PeelingClassUser || sourceLocationMatch.RightClassification != schema.PeelingClassUser {
		t.Fatalf("matched classifications = %#v", sourceLocationMatch)
	}
	sourceFileMatch, ok := findFunctionMatch(got.MatchedFunctions, "main.service", "main.serviceV2")
	if !ok {
		t.Fatalf("missing source-file match in %#v", got.MatchedFunctions)
	}
	if sourceFileMatch.Score != 90 || sourceFileMatch.Reason != FunctionMatchReasonSourceFile {
		t.Fatalf("source-file match = %#v", sourceFileMatch)
	}
	moduleLocalMatch, ok := findFunctionMatch(got.MatchedFunctions, "example.com/sample/internal/app.Handler", "example.com/sample/v2/internal/app.Handler")
	if !ok {
		t.Fatalf("missing module-local normalized-name match in %#v", got.MatchedFunctions)
	}
	if moduleLocalMatch.Score != 92 || moduleLocalMatch.Reason != FunctionMatchReasonModuleLocalNormalizedName {
		t.Fatalf("module-local normalized-name match = %#v", moduleLocalMatch)
	}
	sourceFileTransfer, ok := findTransferCandidate(got.TransferCandidates, "main.service", "main.serviceV2")
	if !ok {
		t.Fatalf("missing source-file transfer candidate in %#v", got.TransferCandidates)
	}
	if sourceFileTransfer.Disposition != TransferDispositionReview || sourceFileTransfer.ProjectedClassificationReason != schema.PeelingEvidenceBuildInfoPath {
		t.Fatalf("source-file transfer candidate = %#v", sourceFileTransfer)
	}
	if sourceFileTransfer.ProjectedPackage != "main" {
		t.Fatalf("source-file transfer package = %#v", sourceFileTransfer)
	}
	moduleLocalTransfer, ok := findTransferCandidate(got.TransferCandidates, "example.com/sample/internal/app.Handler", "example.com/sample/v2/internal/app.Handler")
	if !ok {
		t.Fatalf("missing module-local transfer candidate in %#v", got.TransferCandidates)
	}
	if moduleLocalTransfer.Disposition != TransferDispositionReady || moduleLocalTransfer.ProjectedName != "example.com/sample/internal/app.Handler" {
		t.Fatalf("module-local transfer candidate = %#v", moduleLocalTransfer)
	}
	if moduleLocalTransfer.ProjectedPackage != "internal/app" {
		t.Fatalf("module-local transfer package = %#v", moduleLocalTransfer)
	}
	exactAccepted, ok := findAcceptedTransfer(got.AcceptedTransfers, "main.main", "main.main")
	if !ok {
		t.Fatalf("missing exact accepted transfer in %#v", got.AcceptedTransfers)
	}
	if exactAccepted.AcceptedBy != FunctionMatchReasonExactName || exactAccepted.ProjectedClassification != schema.PeelingClassUser {
		t.Fatalf("exact accepted transfer = %#v", exactAccepted)
	}
	if exactAccepted.ProjectedPackage != "main" {
		t.Fatalf("exact accepted transfer package = %#v", exactAccepted)
	}
	if _, ok := findAcceptedTransfer(got.AcceptedTransfers, "main.service", "main.serviceV2"); ok {
		t.Fatalf("source-file review candidate unexpectedly accepted in %#v", got.AcceptedTransfers)
	}
	if got.TransferReview.ReviewCount != 1 || got.TransferReview.AutoAcceptedCount != 3 || got.TransferReview.ReviewPackageCount != 1 {
		t.Fatalf("transfer review = %#v", got.TransferReview)
	}
	if len(got.TransferReview.Items) != 1 {
		t.Fatalf("transfer review items = %#v", got.TransferReview)
	}
	reviewItem := got.TransferReview.Items[0]
	if reviewItem.Package != "main" ||
		reviewItem.ImportPath != "example.com/sample" ||
		reviewItem.LeftName != "main.service" ||
		reviewItem.RightName != "main.serviceV2" ||
		reviewItem.MatchReason != FunctionMatchReasonSourceFile ||
		reviewItem.ProjectedClassificationCause != schema.PeelingEvidenceBuildInfoPath {
		t.Fatalf("transfer review item = %#v", reviewItem)
	}
	reviewPackage := got.TransferReviewPackages[0]
	if reviewPackage.Name != "main" ||
		reviewPackage.ImportPath != "example.com/sample" ||
		reviewPackage.Classification != schema.PeelingClassUser ||
		reviewPackage.ReviewCount != 1 ||
		reviewPackage.AutoAcceptedCount != 2 ||
		reviewPackage.HighestMatchScore != 90 ||
		reviewPackage.HighestMatchReason != FunctionMatchReasonSourceFile ||
		reviewPackage.SampleLeftName != "main.service" ||
		reviewPackage.SampleRightName != "main.serviceV2" {
		t.Fatalf("transfer review package = %#v", reviewPackage)
	}
	reviewFocus := got.TransferReviewFocus
	reviewPlan := got.TransferReviewPlan[0]
	if reviewPlan.Action != "review_package" ||
		reviewPlan.Package != "main" ||
		reviewPlan.ImportPath != "example.com/sample" ||
		reviewPlan.HighestMatchReason != FunctionMatchReasonSourceFile ||
		reviewPlan.ItemCount != 1 ||
		len(reviewPlan.Items) != 1 {
		t.Fatalf("transfer review plan = %#v", got.TransferReviewPlan)
	}
	if reviewPlan.Items[0].LeftName != "main.service" ||
		reviewPlan.Items[0].RightName != "main.serviceV2" ||
		reviewPlan.Items[0].MatchReason != FunctionMatchReasonSourceFile {
		t.Fatalf("transfer review plan items = %#v", reviewPlan.Items)
	}
	if reviewFocus.Action != "review_package" ||
		reviewFocus.Package != "main" ||
		reviewFocus.ImportPath != "example.com/sample" ||
		reviewFocus.Classification != schema.PeelingClassUser ||
		reviewFocus.ReviewCount != 1 ||
		reviewFocus.AutoAcceptedCount != 2 ||
		reviewFocus.HighestMatchScore != 90 ||
		reviewFocus.HighestMatchReason != FunctionMatchReasonSourceFile ||
		reviewFocus.SampleLeftName != "main.service" ||
		reviewFocus.SampleRightName != "main.serviceV2" ||
		reviewFocus.ItemCount != 1 ||
		len(reviewFocus.Items) != 1 {
		t.Fatalf("transfer review focus = %#v", reviewFocus)
	}
	if reviewFocus.Items[0].LeftName != "main.service" ||
		reviewFocus.Items[0].RightName != "main.serviceV2" ||
		reviewFocus.Items[0].MatchReason != FunctionMatchReasonSourceFile {
		t.Fatalf("transfer review focus items = %#v", reviewFocus.Items)
	}
	mainPackage, ok := findTransferPackage(got.TransferPackages, "main", "example.com/sample")
	if !ok {
		t.Fatalf("missing main transfer package in %#v", got.TransferPackages)
	}
	if mainPackage.CandidateCount != 3 || mainPackage.ReadyCount != 2 || mainPackage.ReviewCount != 1 || mainPackage.AcceptedCount != 2 {
		t.Fatalf("main transfer package = %#v", mainPackage)
	}
	if mainPackage.HighestAcceptedBy != FunctionMatchReasonExactName || mainPackage.HighestCandidateReason != FunctionMatchReasonExactName {
		t.Fatalf("main transfer package reasons = %#v", mainPackage)
	}
	internalPackage, ok := findTransferPackage(got.TransferPackages, "internal/app", "example.com/sample/internal/app")
	if !ok {
		t.Fatalf("missing internal/app transfer package in %#v", got.TransferPackages)
	}
	if internalPackage.CandidateCount != 1 || internalPackage.ReadyCount != 1 || internalPackage.AcceptedCount != 1 {
		t.Fatalf("internal/app transfer package = %#v", internalPackage)
	}
	if internalPackage.ModuleLocal || internalPackage.Classification != schema.PeelingClassUser {
		t.Fatalf("internal/app transfer package posture = %#v", internalPackage)
	}
	if len(got.AddedPackages) != 1 || got.AddedPackages[0] != "internal/app" {
		t.Fatalf("added packages = %#v", got.AddedPackages)
	}
	if len(got.AddedTypes) != 1 || got.AddedTypes[0] != "main.extra" {
		t.Fatalf("added types = %#v", got.AddedTypes)
	}
	if len(got.RemovedStrings) != 1 || got.RemovedStrings[0] != "hello" {
		t.Fatalf("removed strings = %#v", got.RemovedStrings)
	}
	if len(got.AddedSourceFiles) != 1 || got.AddedSourceFiles[0] != "helper.go" {
		t.Fatalf("added source files = %#v", got.AddedSourceFiles)
	}
	if len(got.AddedRefinedPasses) != 1 || got.AddedRefinedPasses[0] != "string-segments" {
		t.Fatalf("added refined passes = %#v", got.AddedRefinedPasses)
	}
}

func TestCompareTransferReviewPlanOrdersByPriority(t *testing.T) {
	t.Parallel()

	left := schema.Analysis{
		BuildInfo: &schema.BuildInfo{Path: "example.com/sample"},
		Functions: []schema.Function{
			{Name: "main.service", Package: "main", ImportPath: "example.com/sample", SourceFile: "service.go", SourceLine: 40},
			{Name: "main.alpha", Package: "main", ImportPath: "example.com/sample", SourceFile: "alpha.go", SourceLine: 10},
			{Name: "example.com/sample/internal/app.Handler", Package: "internal/app", ImportPath: "example.com/sample/internal/app", SourceFile: "handler.go", SourceLine: 7},
		},
		Peeling: &schema.PeelingAnalysis{
			Functions: []schema.PeelingFunction{
				{Name: "main.service", Package: "main", ImportPath: "example.com/sample", SourceFile: "service.go", SourceLine: 40, Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceBuildInfoPath},
				{Name: "main.alpha", Package: "main", ImportPath: "example.com/sample", SourceFile: "alpha.go", SourceLine: 10, Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceBuildInfoPath},
				{Name: "example.com/sample/internal/app.Handler", Package: "internal/app", ImportPath: "example.com/sample/internal/app", SourceFile: "handler.go", SourceLine: 7, Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceModuleLocal},
			},
		},
	}
	right := schema.Analysis{
		BuildInfo: &schema.BuildInfo{Path: "example.com/sample/v2"},
		Functions: []schema.Function{
			{Name: "main.serviceV2", Package: "main", ImportPath: "example.com/sample/v2", SourceFile: "service.go", SourceLine: 48},
			{Name: "main.alphaV2", Package: "main", ImportPath: "example.com/sample/v2", SourceFile: "alpha.go", SourceLine: 15},
			{Name: "example.com/sample/v2/internal/app.HandlerV2", Package: "internal/app", ImportPath: "example.com/sample/v2/internal/app", SourceFile: "handler.go", SourceLine: 21},
		},
		Peeling: &schema.PeelingAnalysis{
			Functions: []schema.PeelingFunction{
				{Name: "main.serviceV2", Package: "main", ImportPath: "example.com/sample/v2", SourceFile: "service.go", SourceLine: 48, Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceBuildInfoPath},
				{Name: "main.alphaV2", Package: "main", ImportPath: "example.com/sample/v2", SourceFile: "alpha.go", SourceLine: 15, Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceBuildInfoPath},
				{Name: "example.com/sample/v2/internal/app.HandlerV2", Package: "internal/app", ImportPath: "example.com/sample/v2/internal/app", SourceFile: "handler.go", SourceLine: 21, Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceModuleLocal},
			},
		},
	}

	got := Compare(left, right)

	if len(got.TransferReviewPlan) != 2 {
		t.Fatalf("transfer review plan = %#v", got.TransferReviewPlan)
	}
	if got.TransferReviewPlan[0].Package != "main" || got.TransferReviewPlan[0].ReviewCount != 2 || got.TransferReviewPlan[0].ItemCount != 2 || len(got.TransferReviewPlan[0].Items) != 2 {
		t.Fatalf("transfer review plan[0] = %#v", got.TransferReviewPlan[0])
	}
	if got.TransferReviewPlan[1].Package != "internal/app" || got.TransferReviewPlan[1].ReviewCount != 1 || got.TransferReviewPlan[1].ItemCount != 1 || len(got.TransferReviewPlan[1].Items) != 1 {
		t.Fatalf("transfer review plan[1] = %#v", got.TransferReviewPlan[1])
	}
	if got.TransferReviewFocus == nil || got.TransferReviewFocus.Package != "main" {
		t.Fatalf("transfer review focus = %#v", got.TransferReviewFocus)
	}
}

func findFunctionMatch(matches []FunctionMatch, leftName, rightName string) (FunctionMatch, bool) {
	for _, match := range matches {
		if match.LeftName == leftName && match.RightName == rightName {
			return match, true
		}
	}

	return FunctionMatch{}, false
}

func findTransferCandidate(candidates []TransferCandidate, leftName, rightName string) (TransferCandidate, bool) {
	for _, candidate := range candidates {
		if candidate.LeftName == leftName && candidate.RightName == rightName {
			return candidate, true
		}
	}

	return TransferCandidate{}, false
}

func findAcceptedTransfer(candidates []AcceptedTransfer, leftName, rightName string) (AcceptedTransfer, bool) {
	for _, candidate := range candidates {
		if candidate.LeftName == leftName && candidate.RightName == rightName {
			return candidate, true
		}
	}

	return AcceptedTransfer{}, false
}

func findTransferPackage(summaries []TransferPackageSummary, name, importPath string) (TransferPackageSummary, bool) {
	for _, summary := range summaries {
		if summary.Name == name && summary.ImportPath == importPath {
			return summary, true
		}
	}

	return TransferPackageSummary{}, false
}
