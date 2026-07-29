# Codex CLI Layout Notes

This reference records the portable repo layout chosen for GoREveal after comparing the current
repository setup with the official Codex docs.

Portable locations:

- root instructions: `AGENTS.md`
- repo-local skills: `.agents/skills/<skill>/SKILL.md`
- repo-local subagents: `.codex/agents/*.toml`
- repo-local config: `.codex/config.toml`

Compatibility note:

- the existing `skills/` tree remains in the repository because current local workflows already
  reference it
- `.agents/skills/` is now the preferred Codex-native layout

Official docs referenced:

- https://developers.openai.com/codex/guides/agents-md
- https://developers.openai.com/codex/skills
- https://developers.openai.com/codex/subagents
- https://developers.openai.com/codex/mcp
- https://developers.openai.com/codex/configuration
