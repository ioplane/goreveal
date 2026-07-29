<!--
The PR title must follow Conventional Commits — CI enforces it, and the title
becomes the squash-merge commit and the changelog entry.

  <type>(<scope>): <description>
  feat(core/pclntab): recover funcnametab offsets on stripped ELF

Types: feat fix perf refactor docs test build ci chore revert security
-->

## What and why

<!-- What changes, and what problem it solves. Link the issue: Closes #123 -->

## Evidence

<!--
Every meaningful change carries evidence matched to what it changed. Paste the
relevant output, benchmark diff, or fixture path. "Tests pass" is not evidence.
-->

- [ ] Corpus fixture added or updated
- [ ] Golden snapshot added or updated
- [ ] Differential comparison covered
- [ ] Fuzz target added or extended
- [ ] Benchmark evidence (`benchstat` before/after)
- [ ] Documentation or contract notes updated
- [ ] Not applicable — explain below

<!-- If a golden snapshot changed, explain the diff. Refreshing a snapshot to
     make a test pass is how correctness regressions ship. -->

## Verification

<!-- Paste the commands you ran and their outcome. -->

```text
task lint            # result:
task test            # result:
task lint-scripts    # result:
```

## Claim boundaries

<!--
GoREveal's value is that it does not guess. Confirm this change respects that.
-->

- [ ] Absent evidence still yields `unavailable`, an empty list, or an explicit
      error — never an inferred value
- [ ] Provenance and confidence fields remain accurate for any new output
- [ ] Raw recovered truth is preserved; refinement layers do not overwrite it

## Compliance

- [ ] No code copied, translated, or derived from `gore`, `redress`, `GoReSym`,
      `GoResolver`, `gostringungarbler`, or `AlphaGolang`
- [ ] Module boundaries respected (`core` free of CLI/storage/plugin imports;
      `plugins/` free of recovery logic)
- [ ] No host paths, hostnames, secrets, or internal identifiers introduced
- [ ] Breaking changes to the `schema` contract are called out below

## Breaking changes

<!-- Describe the break and the migration path, or write "None". -->

None
