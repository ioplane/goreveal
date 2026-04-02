# GoREveal Progress Assessment

> Status: product/program checkpoint
> Date: 2026-03-20
> Purpose: quantify current progress by feature block, baseline-tool comparison, and sprint lane so planning can use a stable percentage view instead of only narrative notes.

## Method

These percentages are deliberately heuristic.

They mean:
- `0%`: effectively absent
- `25%`: early foothold exists, but not reliable for regular use
- `50%`: useful and real, but clearly incomplete
- `75%`: strong enough for regular use within current scope
- `100%`: mature for the currently declared product scope

They are not:
- velocity
- story points
- release readiness
- parity claims against external tools

## Overall Progress

- Platform as a product: `99%`
- Accuracy/recovery engine: `74%`
- Overall roadmap completion: `99%`

## Process Health

| Area | Reading |
| --- | --- |
| Delivery discipline | `high` — bounded vertical slices, Podman-first verification, and doc-sync discipline are now routine rather than occasional |
| Strategic focus | `high` — the roadmap is much less reactive than before, and the latest protected-binary work now ends at a real decision boundary instead of another open-ended parser detour; the main remaining risk is now failing to keep the active PM lane on the explicitly ranked `Sprint 16` question instead of drifting back into parser/deobfuscation work without a named workflow-shaped pain point |
| Evidence quality | `high` — fixtures, snapshots, baseline comparison, external-binary validation, the widened protected-binary matrix, and the first real `garble` rows on both `linux/amd64` and `linux/arm64` are all now measured, and the canonical protected-matrix task now prefers a local upstream `garble` checkout when available instead of being pinned to the stale `v0.15.0` release path; the protected gap is now explained by measured runtime/blocker evidence, bounded absolute `PC` address hints, a first sampled absolute `PC` foothold, and an analyst-facing `address_only` foothold that now survives on both Linux architectures rather than only one, with the same compact runtime surface now locked into thin export-contract tests and plugin consumer tests; the latest intermediate rerun has also been refreshed into `.tmp/fresh-eval/go-baselines.json`, re-confirming real file visibility on measured `ELF`, `PE`, and `Mach-O` samples instead of file-blind operator output, while the latest bounded timing sample on `rclone-linux-amd64` still does not show an operational efficiency regression that would outrank workflow or source-confidence work |
| Narrative accuracy | `high` — planning and README now mostly describe landed capability instead of optimistic posture, and the dev-image/tooling contract is now closer to the actual operator workflow |

## Feature Block Progress

| Block | Progress | Reading |
| --- | ---: | --- |
| Core Analysis | `90%` | Real and useful across rich ELF, stripped ELF, the bounded PE checkpoint, and the first bounded Mach-O function foothold, and now includes a thin line-table-backed source-visibility increment plus raw ELF `pclntab` header evidence, nonzero header-level function/file-count hints on garbled rows, monotonic `functab` PC-offset hints under garble, bounded absolute `PC` address hints within `.text`, a first sampled absolute `PC` address foothold, a first explicit analyst-facing `address_only` function foothold that now survives on both `linux/amd64` and `linux/arm64`, a compact operator-facing explanation of whether that foothold is backed by `moduledata_text` or the ELF `.text` section, and a first bounded source-confidence contract over already-known file/package truth through `source_evidence_kind` plus compact `source_evidence_summary` breakdowns; fresh external measurements now also confirm real file visibility on representative `ELF`, `PE`, and `Mach-O` binaries |
| CLI Surface | `99%` | The main operator-facing commands are coherent, `inspect peeling` plus `peel` now make the first differentiator visible on real binaries outside the fixture corpus, and the review workflow now has dedicated `diff review sqlite`, `diff handoff sqlite`, and `diff next sqlite` operator paths instead of only generic diff output; the handoff bridge now also carries structured `target_profiles` plus a self-describing artifact bundle for workstation consumers with explicit contract, transport, workspace-phase, host-action, binding-entrypoint, required-artifact, and expected-outcome hints, while `diff next sqlite` now also exposes a compact `review_checklist`, a compact `review_snapshot`, explicit `review_progress`, a compact `up_next` package snapshot, and an `upcoming_packages` horizon with sample pair and strongest-match context, so the CLI now reads more like a usable local workflow than a collection of raw subcommands |
| Storage | `100%` | SQLite persistence and diffing are already strong product assets, and now include explainable function matches, transfer-candidate preview, deterministic accepted transfers, a first package-level transfer queue summary, a compact analyst-facing `transfer_review` surface for pending human-review items, a package-first `transfer_review_packages` triage surface, a bounded `transfer_review_focus` first-pass bundle over the same workflow state, a compact `transfer_review_plan` action queue, a dedicated `diff review sqlite` CLI/operator projection over that review path, a compact machine-readable `handoff` block with left/right input context and recommended workstation targets, a self-describing artifact bundle, a structured per-target `target_profiles` handoff contract with explicit export contracts, preferred transport, artifact roles, workspace phases, host actions, binding-entrypoint hints, required-artifact hints, and expected host-outcome hints, plus dedicated `diff handoff sqlite` and `diff next sqlite` operator-facing projections over the same bridge; the `diff next sqlite` path now also carries self-contained `recommended_actions`, explicit `review_progress`, a compact `up_next` package snapshot, and an `upcoming_packages` horizon with sample pair and strongest-match context, so the immediate review pass no longer requires manual command reconstruction or ad hoc queue math |
| Validation / Evidence | `94%` | Good discipline, snapshots, differential reports, external-binary smoke checks, the first open-source cross-platform matrix, the widened purpose-built protected-binary matrix across `amd64` and `arm64`, the first Mach-O corpus fixture, and the first real `garble` rows on both `linux/amd64` and `linux/arm64` all exist; the current intermediate repo-local rerun also revalidated the rich `ELF` overlap contract plus bounded `PE`/`Mach-O` footholds and the measured `rclone-linux-amd64` comparison, while a bounded timing sample confirmed that operational efficiency is not the next blocking concern |
| Integrations | `82%` | Thin IDA/Ghidra adapters are real and validated, and they now mirror more canonical truth without plugin-side recomputation; the newest protected runtime surface is now asserted both at export-contract level and at plugin consumer level, while the handoff bridge is now close to a target-bindable workstation contract rather than only an operator hint blob |
| Capability Transfer | `98%` | Bounded runtime trust, the PE checkpoint, engine-owned code peeling, bounded function matching, transfer-candidate preview, accepted-transfer surface, package-level transfer summaries, the compact `transfer_review` queue, the package-first `transfer_review_packages` triage layer, the `transfer_review_focus` first-pass bundle, the new `transfer_review_plan` action queue with package-ordered attached review items, the `diff review sqlite` CLI/operator projection, the machine-readable `handoff` block, the `diff handoff sqlite` operator-facing bridge, and the `diff next sqlite` actionability projection now together form a more complete operator bridge because the next-step artifact also carries self-contained `recommended_actions` for review, handoff, and target export, a compact `review_checklist`, a compact `review_snapshot`, explicit `review_progress`, a compact `up_next` package snapshot, and an `upcoming_packages` horizon with sample pair and strongest-match context, while the self-describing artifact bundle, explicit target binding hints, and expected host-outcome semantics now make the local review/handoff story feel substantially product-shaped rather than preview-shaped; the fresh intermediate rerun did not reopen parser-first priorities and instead reinforced that the product edge is currently workflow/value-shaped |
| Documentation / Planning | `100%` | Architecture, sprint, and product docs are already a strength of the repo and now track code-peeling, PE checkpoints, the Mach-O foothold, transfer workflow progress, cross-platform external-binary validation, the current garble-blocker triage posture, the measured RE-lab/workstation context from `rehelp`, the new review-handoff operator path, and the dev-image/operator-tool contract directly |
| Performance / SIMD | `5%` | Still intentionally deferred behind accuracy work |
| Service / API | `0%` | Not meaningfully started |
| Deobfuscation Depth | `18%` | Scaffolding and first passes exist, but this is still an early lane |

## Module Readiness

| Module | Readiness | Reading |
| --- | --- | --- |
| `schema` | `high` | Canonical contracts are dense, provenance-aware, and now include bounded runtime trust plus `peeling` surfaces |
| `engine` | `high` | Orchestration is stable, source-tree is real, and engine-owned code peeling now sits in the correct boundary |
| `cmd/goreveal` | `high` | `analyze`, `inspect *`, `peel`, exports, and diffs already form a coherent operator surface |
| `docs/ops` | `high` | Plans, architecture docs, AGENTS/skills, Podman-first flow, script linting, and the dev-image operator toolbox are all mature for the current scope |
| `core` | `high-medium` | Real recovery exists for current ELF fixtures, external stripped ELF binaries, the bounded PE checkpoint, and the first bounded Mach-O function foothold, and line-table-backed file evidence plus raw protected-binary runtime signals are now projected without new parser claims; fresh external runs now confirm that this bounded file-evidence path is real on representative `ELF`, `PE`, and `Mach-O` binaries, not only on fixtures. The remaining protected gap is now a measured custom-`pclntab` named-recovery problem with preserved header-level hints, monotonic `functab` PC-offset hints, bounded absolute `PC` address hints, a first sampled absolute `PC` foothold, a first explicit `address_only` function foothold now portable across Linux architectures through a separate ELF `.text` section range, a compact explanation surface for whether that foothold is `moduledata_text`-backed or `elf_text_section`-backed, and preserved `moduledata` bridges rather than a vague environment issue, but deeper runtime/type semantics and broader format coverage are still missing |
| `storage` | `high` | Persistence and stored-run diffing now support explainable matching, projected transfer truth, deterministic accepted transfers, package-oriented queue summaries, a compact analyst-facing review queue, package-first review triage, and a first-pass focus bundle, making storage a real workflow layer rather than only a persistence layer |
| `plugins` | `high-medium` | Thin IDA/Ghidra adapters are still intentionally narrow, but the contract is now healthier: canonical runtime/export truth is checked at export and consumer boundaries rather than depending on informal adapter behavior, and the remaining gap is analyst convenience rather than contract drift |
| `corpus` | `high` | Rich ELF, stripped ELF, PE, and now a bounded Mach-O fixture exist with snapshots and differential discipline, and external plus protected matrix validation now spans `amd64` and `arm64`, though breadth across families and edge cases is still limited |
| `deobfuscation` | `low` | Architecturally correct, but still a thin capability area rather than a strong product lane |
| `api` | `low` | Still effectively unstarted by design |

## Capability-Transfer Progress By Baseline

| Tool | Progress | Reading |
| --- | ---: | --- |
| `gore` | `62%` | Product shape already stronger; still behind on less-heuristic type/package truth |
| `redress` | `60%` | Fresh external comparison confirms `redress` is still the strongest source/file-oriented practical reference, but `GoREveal` now exceeds it on function/package/transfer workflow shape |
| `GoReSym` | `38%` | Bounded runtime truth is now substantial; true semantic runtime/type depth is still missing, but fresh external comparison no longer supports treating GoReSym as the dominant practical product baseline outside runtime semantics |
| `GoResolver` | `10%` | Architecture is ready, but capability transfer is mostly future work |
| `gostringungarbler` | `15%` | Early refined string work exists; targeted deobfuscation depth does not |
| `AlphaGolang` | `48%` | Integration boundary is cleaner; analyst convenience inside IDA is still thinner |

## Sprint Progress

| Sprint / Chunk | Progress | Reading |
| --- | ---: | --- |
| Chunk 0 | `100%` | Foundations, docs, skills, repo contract |
| Sprint 1 | `100%` | Minimal end-to-end analysis |
| Sprint 2 | `100%` | Build info, functions, packages |
| Sprint 3 | `100%` | Types, strings |
| Sprint 4 | `100%` | Source-tree v1 |
| Sprint 5 | `100%` | Refined layer and deobfuscation scaffold |
| Sprint 6 | `100%` | SQLite persistence and diffing |
| Sprint 7 | `68%` | Good proof lane, but still not a broad evidence lane |
| Sprint 8 | `100%` | Export contracts and thin adapters landed |
| Sprint 9 | `5%` | Deliberately deferred |
| Sprint 10 | `0%` | Not started |
| Sprint 11 | `100%` | Package/type/source usability checkpoint complete |
| Sprint 12 | `99%` | Runtime trust, stripped/rich stability, the PE runtime checkpoint, the first PE and Mach-O function footholds, code-peeling MVP, bounded function matching, transfer-candidate preview, accepted-transfer surface, package-level transfer summaries, thin source-visibility response, external cross-platform smoke checks, the protected-binary triage lane through a portable `address_only` foothold, the first focused review bundle, the `diff review sqlite` operator projection, its compact workstation-facing `handoff` block, and the dedicated `diff handoff sqlite` bridge are all landed |
| Sprint 13 | `78%` | Re-baselined as workstation handoff contract work; the dedicated `diff handoff sqlite` path, structured `target_profiles`, a self-describing artifact bundle, contract/transport hints, artifact roles, workspace phases, host actions, binding-entrypoint hints, required-artifact hints, and expected host-outcome hints are already landed, and the first `Sprint 14` slice now exists beside that contract, so the remaining work is now only a small lock-or-stop decision rather than shape invention |
| Sprint 14 | `52%` | Review workflow actionability is now materially underway through the compact `transfer_review_plan` queue and the dedicated `goreveal diff next sqlite <database> <left-id> <right-id>` operator projection, and that plan now carries package-ordered attached review items plus self-contained `recommended_actions`, a compact `review_checklist`, a compact `review_snapshot`, explicit `review_progress`, a compact `up_next` package snapshot, and an `upcoming_packages` horizon with sample pair and strongest-match context, so the operator can act on the next bundle without falling back to the larger raw queue payload, reconstructing the obvious follow-on commands by hand, manually counting what remains after the current package, guessing what “done for this bundle” means, or parsing the full queue just to see what comes next; the lane now also has an explicit stop-condition, and the current weighted reading is that the declared local operator loop is likely already complete enough for this scope unless a new measured gap appears |
| Sprint 15 | `42%` | The lane now has a bounded-success shape: `source_tree` and enriched package surfaces project a bounded `source_evidence_kind = dwarf_paths | line_table_files | package_fallback` vocabulary over already-known truth, `source_tree` carries a compact `source_evidence_summary`, and that summary also exposes per-evidence-class file counts, so operators no longer need to manually infer either the package landscape or the file-density shape from the full package list; the current weighted reading is to freeze this lane by default unless one new named inference step is explicitly demonstrated |
| Sprint 16 | `30%` | The PM-side target ranking and stop-condition are now landed, and the first bounded DEV rerun is landed too: `garble`-class workflows rank first because they have the highest current analyst ambiguity and commercial relevance among already-measured protected classes, the first workflow-shaped pain point is explicit as “no review-ready anchor set on current garbled rows for peel/review/handoff/next despite the existing `address_only` foothold”, and the first neighboring-build rerun now confirms that both measured Linux neighbor pairs currently yield zero matched functions, zero transfer/review packages, no handoff, and no next-step review projection |
| Sprint 17 | `0%` | Server control-plane foundations remain architecturally chosen but intentionally deferred until the local workflow and handoff story is mature enough to justify remote orchestration |
| Sprint 18 | `0%` | Metadata and remote interop platform work remains downstream of both local workflow maturity and server control-plane foundations |
| Sprint 19 | `0%` | Public release readiness and licensing are now explicit later-horizon work rather than vague release polish |
| Sprint 20 | `0%` | Evidence expansion and comparative automation are now explicit follow-on work rather than an implicit maintenance bucket |
| Sprint 21 | `0%` | Build correlation and version tracking remain the strongest workflow epic after the current local product and release/evidence horizons settle |
| Sprint 22 | `0%` | Metadata knowledge network remains the strongest long-term moat, but only after version-tracking and remote metadata foundations are real |
| Sprint 23 | `0%` | Analyst workspace automation and replay should come only after remote metadata and build-correlation value are proven, not before |
| Sprint 24 | `0%` | Comparative knowledge packs and decision support should be the late-stage synthesis lane after evidence automation, metadata, and workspace automation are all real |

## Comparative Reading

### Where GoREveal is already stronger than the baseline tools as a product

- canonical schema
- container-first reproducibility
- SQLite persistence and diffing
- export-driven integrations
- planning discipline and explicit sprint strategy

### Where GoREveal is competitive but not yet dominant

- function/package UX
- source-tree usefulness on the canonical ELF fixture
- thin IDA/Ghidra integration flow
- evidence-backed operator workflows
- first code-peeling usefulness for analyst triage
- stripped-ELF usefulness on real external Go binaries
- broad Linux server-binary usefulness
- build-to-build transfer workflow value beyond preview and queue stages

### Where GoREveal is still clearly behind

- runtime-semantic `moduledata` decoding
- typelinks-driven type recovery
- cross-version runtime confidence
- deobfuscation depth
- rich RE-tool analyst workflows
- fixture breadth across PE families and other future formats
- Mach-O runtime and semantic depth beyond the first function/package/peeling foothold

## Planning Implication

The percentages confirm the current strategy:

1. `Sprint 12` stays the primary execution lane.
2. `Sprint 7` stays in maintenance mode.
3. `Sprint 11` remains a completed checkpoint unless runtime truth later forces reclassification.
4. The real comparison pass is now landed and the first bounded `PE` function foothold is landed too.
5. The post-`PE` comparison rerun is now completed for the current highest-signal slices.
6. The protected-binary matrix is now widened through `arm64`, the canonical operator path is now truthful for that wider lane, and the first real `garble` rows are now measured through a local upstream `garble` checkout on both `linux/amd64` and `linux/arm64`; the latest bounded triage slices converted the old arm64 split into a portable section-backed `address_only` foothold and then locked that compact runtime surface through export-contract and plugin-consumer tests.
7. That changes the weighted next move: the protected lane now looks more stabilized than blocked, so the default next step should be workflow/value work unless a new measured protected-specific gap appears.
8. The first concrete workflow/value response is now landed too through `storage/diff transfer_review`, `transfer_review_packages`, and `transfer_review_focus`, so the repo is no longer only planning that return.
9. The bounded CLI/operator projection for that focused review pass is now landed too through `goreveal diff review sqlite <database> <left-id> <right-id>`.
10. That same operator path now also carries a compact machine-readable `handoff` block with left/right input context and recommended workstation targets, so explicit host-platform handoff no longer starts from raw diff state alone.
11. A dedicated `goreveal diff handoff sqlite <database> <left-id> <right-id>` projection is now landed too, turning that review bridge into its own operator-facing artifact instead of leaving it embedded only inside `diff review sqlite`.
12. The measured `rehelp` / RE-lab inventory now makes workstation interop more concrete too: `ida-pro-mcp`, `Diaphora`, `BinExport`, `rizin`, and multiple dynamic/symbolic sidecars are real adjacent tools in the operator environment, not generic backlog names.
13. `Sprint 9`, `Sprint 10`, and richer post-`Sprint 13` expansion work remain intentionally behind accuracy depth.
14. the post-`Sprint 12` sprint sequence is now explicit: `Sprint 13` workstation handoff contract, `Sprint 14` review workflow actionability, `Sprint 15` thin semantic/source confidence, `Sprint 16` protected workflow/orchestration
15. `Sprint 13` is no longer only a planning hypothesis: the remaining decision is whether one final thin contract lock is worth it before stopping that lane, not whether that lane should exist at all
16. `Sprint 14` is now started in bounded form through `transfer_review_plan` and `goreveal diff next sqlite <database> <left-id> <right-id>`, and that plan now carries package-ordered attached review items plus self-contained `recommended_actions`, explicit `review_progress`, a compact `up_next` snapshot, and an `upcoming_packages` horizon with sample pair context, so review actionability is no longer only a sequencing note.
17. the later horizon is now explicit too: `Sprint 17` server control-plane foundations and `Sprint 18` metadata/remote interop platform stay ordered but deferred behind local workflow maturity

## Near-Term Priority

The next highest-value low-regret work is now:
- keep the new bounded `Mach-O` foothold stable: function/package/peeling only, no broad runtime claims
- keep the new bounded `PE` foothold stable: function/package/peeling plus existing runtime posture, no broad PE semantic claims
- treat the completed post-`PE` comparison rerun plus the verified protected matrix as the evidence baseline, because the main open question is now protected obfuscation behavior rather than a missing format foothold
- use the newly stabilized protected lane as a floor, not as the default next implementation lane
- keep the re-opened workflow/value lane stable now that `transfer_review`, `transfer_review_packages`, `transfer_review_focus`, `transfer_review_plan`, `diff review sqlite`, `diff handoff sqlite`, and `diff next sqlite` are all landed, with the next-step path now carrying self-contained `recommended_actions`, and use the fresh external comparison as evidence that this lane is backed by real cross-format utility rather than fixture-local success
- after that, harden explicit host-platform MCP and workstation handoff planning on top of the now-landed `diff handoff sqlite` bridge, using the measured `rehelp` / RE-lab environment as the real interop baseline
- the first bounded hardening slices inside that lane are now landed too through structured per-target `target_profiles` for `ida` and `ghidra`, plus explicit export-contract IDs, preferred transport hints, artifact roles, workspace phases, and host actions; the next handoff step should build on that shape instead of replacing it
- the current low-regret next step is therefore to treat both `Sprint 14` and `Sprint 15` as frozen by default for the current scope, move active PM ranking into `Sprint 16`, and keep protected work corpus/comparison-first rather than reopening parser or deobfuscation lanes by drift
- after the current local sequence settles, the next later product horizon should move into server/control-plane and remote interop only in that order: first bounded control-plane foundations, then remote metadata/MCP platform work
- only reopen the protected lane immediately if one new protected-specific analyst surface clearly promises materially better signal than the current `address_only` + text-source + span projection
- avoid reopening broad parser or heuristic-rewrite work

Current deferred-resume handoff:
- `docs/plans/2026-03-20-goreveal-deferred-continuation.md`

Current product-direction brainstorming:
- `docs/plans/2026-03-20-goreveal-market-killer-features-brainstorm.md`

Current cross-platform empirical validation:
- `docs/plans/2026-03-31-goreveal-external-binary-matrix-evaluation.md`

Current real baseline comparison:
- `docs/plans/2026-04-01-goreveal-initial-baseline-comparison-results.md`

Current universal-workbench comparison:
- `docs/plans/2026-04-01-goreveal-universal-re-workbench-comparison.md`

Current execution order:
- `docs/plans/2026-04-01-goreveal-next-execution-plan.md`

Current sprint baseline:
- `docs/plans/2026-04-01-goreveal-post-sprint12-sprint-plan.md`
