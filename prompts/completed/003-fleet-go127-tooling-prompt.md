---
status: completed
summary: 'Ensured Go 1.27 tooling compatibility: moved gofmt -w to run last in the format target with the canonical comment, bumped golangci-lint to v2.13.1 and errcheck to v1.20.0 in tools.env, ran gofmt normalization, and added a CHANGELOG Unreleased entry'
execution_id: cron-go127-exec-003-fleet-go127-tooling-prompt
dark-factory-version: dev
created: "2026-08-21T12:02:45Z"
queued: "2026-08-21T12:02:45Z"
started: "2026-08-21T12:02:47Z"
completed: "2026-08-21T12:16:48Z"
---

# Ensure Go 1.27 tooling compatibility in format target and tools.env

<summary>
- Ensure `gofmt -w` runs last in the `format` target so golines' wrapping is normalized before the gofmt lint check
- Bump golangci-lint to v2.13.1 (fixes staticcheck `buildir` panic on Go 1.27 AST)
- Bump errcheck to v1.20.0 (fixes `package "context" without types` on Go 1.27)
- Reformat any go files golines over-indented so the gofmt check passes
- No runtime behavior change — build-tooling only, matching the canonical fix in bborbe/go-skeleton
</summary>

<objective>
Make this repo's `make precommit` pass on the Go 1.27.0 toolchain: gofmt runs last in `format`, and tools.env pins golangci-lint v2.13.1 + errcheck v1.20.0 — the same canonical fix merged in bborbe/go-skeleton (commits 0a0800b + d42ee30, merge e7818706).
</objective>

<context>
Read the repo's `format:` target (lives in `Makefile`, or in `Makefile.precommit` when `Makefile` does `include tools.env` / `include Makefile.precommit`) and `tools.env`.
The canonical fix is bborbe/go-skeleton (commits 0a0800b + d42ee30): `gofmt -w` moved to run last with the comment `# golines last, then gofmt last so its wrapping is normalized and the gofmt lint check passes`, `ERRCHECK_VERSION` bumped to v1.20.0, `GOLANGCI_LINT_VERSION` bumped to v2.13.1.
</context>

<requirements>
1. In the `format:` target, ensure the `gofmt -w` find line runs LAST, after any golines line. If it already runs last, leave it; otherwise move it there. If the `format:` target has no `gofmt -w` line (or no `format:` target exists), skip the reorder and note it. Add the comment `# golines last, then gofmt last so its wrapping is normalized and the gofmt lint check passes` immediately above the `gofmt -w` line if it is not already present.
2. In `tools.env`, ensure `GOLANGCI_LINT_VERSION` is at least v2.13.1 and `ERRCHECK_VERSION` is at least v1.20.0 (bump if lower; leave unchanged if already at or above). If `tools.env` is absent or does not declare these variables, skip the bump and note it.
3. Run `gofmt -w` on every `.go` file excluding `vendor/` (e.g. `find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +`) to normalize any over-indentation golines introduced.
4. Change nothing else — no other target, no other tool version, no code logic.
</requirements>

<constraints>
- Do NOT commit — dark-factory handles git
- Do not change any other tool versions, Makefile targets, or source logic
- Keep the change minimal and mechanical
- Repo-relative paths only
</constraints>

<verification>
Run `make precommit` -- must pass.
Confirm the `gofmt -w` line is the LAST command line in the `format:` target, after the golines line (grep the Makefile/Makefile.precommit), and that `tools.env` shows GOLANGCI_LINT_VERSION v2.13.1+ and ERRCHECK_VERSION v1.20.0+. If precommit still fails on Go 1.27 after the three changes, investigate whether another pinned tool version needs the same bump; do not change unrelated config.
</verification>
