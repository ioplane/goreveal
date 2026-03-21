# GoREveal Sprint 12 PCLN Checkpoint

> Status: checkpoint note
> Date: 2026-03-20
> Purpose: stop blind `.gopclntab` bridge expansion at the right point and translate the current runtime density into the next roadmap decision.

## Why This Checkpoint Exists

`Sprint 12` has now accumulated a dense bounded `.gopclntab` layout model on the canonical ELF fixture:
- `pcHeader`
- `funcnametab`
- `cutab`
- `filetab`
- `pctab`
- `pclntable`

This is enough to prove that the project is no longer just collecting random runtime evidence.

It is not enough to justify:
- generic pcln-table decoding claims
- cross-version support claims
- rewriting package/type heuristics from runtime truth

At this point, continuing to add more slice-shaped bridges on the same fixture has diminishing return and rising parser-growth risk.

## Checkpoint Conclusion

The next best move is **not** another blind `.gopclntab` bridge by default.

The next best move is one of:
1. validate the current bounded runtime model against a second fixture
2. use the current bounded runtime model to test whether any existing heuristic can be safely reduced
3. deliberately pause `Sprint 12` expansion and return to another lane if neither `1` nor `2` is ready

Current recommendation:
- choose `1` first
- only choose `2` if there is a very small, evidence-backed heuristic adjustment
- do not choose “one more bridge” unless it adds new semantic shape rather than just more of the same layout confidence

## What The Current Runtime Layer Now Gives Us

On the canonical fixture, `GoREveal` now has a bounded runtime model that is strong enough to say:
- the current `.gopclntab` layout is internally self-consistent
- `firstmoduledata` is not just loosely aligned with the ELF sections; it also carries a plausible fixture-local pcln-table shape
- bounded typelink and bounded pcln-table semantics can coexist in the same canonical runtime surface

That is a meaningful milestone.

## What It Still Does Not Give Us

It still does **not** give us:
- general `moduledata` decoding
- typelinks-driven type recovery
- runtime-driven package naming
- runtime-driven type scoping
- support claims outside the current fixture/layout family

## Near-Term Backlog After This Checkpoint

1. Add a second runtime fixture or fixture variant to validate the current pcln-layout chain.
2. Re-check whether the existing `DWARF + source-tree` heuristics can safely consume any of the new runtime truth.
3. Keep `Sprint 7` in maintenance mode and avoid widening claims prematurely.
4. Only return to more `Sprint 12` bridge work if the next slice adds clearly new semantic structure.
5. Re-evaluate whether `Sprint 12` remains the best active lane after the second-fixture question is answered.

## PM/Product Reading

This checkpoint is positive:
- it confirms that `Sprint 12` has produced real product value
- it avoids the common failure mode of overfitting one fixture with an ever-growing parser
- it creates a cleaner handoff into either multi-fixture validation or a more meaningful semantic step

This is the right place to pause blind bridge growth.
