# PR Review Concurrency Probe

This file is a **synthetic load-test marker**. It exists to give the `pr-reviewer-agent`
a real, reviewable diff in `cron` (slot `9` of eleven) so that eleven
concurrent reviews can be measured across both clusters.

The probe backs the measurement in `Prove Eleven of Eleven Slots With an Eleven-PR Run`,
which is the measuring row for goal criterion SC7 — *peak reviewer concurrency reaches
all eleven slots (8 sun + 3 prod) during an eleven-PR run across eleven distinct repos*.

## Why eleven repositories

A `pr-review` task acquires two claim keys: `claims/repo/<owner>-<repo>` (primary) and
`claims/pr/<owner>-<repo>-<number>` (secondary). The repository key is the
mutual-exclusion lock, so two PRs in the same repository serialize rather than run
concurrently. Eleven simultaneous reviewer pods therefore require eleven distinct
repositories, one pull request each.

## What the measurement records

| Field | Source | Why it matters |
|---|---|---|
| Sample timestamp | sampler clock, UTC | anchors the peak to a wall-clock instant |
| Cluster | `sun` / `prod` | the two pools are counted separately |
| Quota hard | `resourcequota/github-pr-review-agent` `.spec.hard.pods` | 8 on sun, 3 on prod |
| Quota used | same object, `.status.used.pods` | corroborates, never proves |
| Active pods | pod phase `Pending` or `Running` | the concurrency figure |
| Active task ids | `agent.benjamin-borbe.de/task-id` label | attributes the count to specific tasks |

The quota line alone is corroboration rather than proof: a namespace-wide
`pods: 8/8` counts every pod that quota admits, so unrelated or leftover review pods
satisfy it with none of the eleven PRs in flight. The per-pod task id is what ties the
count to the set under test.

## The four gates a probe PR must stay inside

| Gate | Value | Failure mode if exceeded |
|---|---|---|
| `MAX_ADDITIONS` | 1200 | parked at `human_review`, never spawns a reviewer |
| `MAX_CHANGED_FILES` | 100 | same |
| `TRIVIAL_AUTO_MERGE` | `"true"` | a mechanically trivial PR self-merges and the review goes moot |
| `REPO_ALLOWLIST` | `github.com/bborbe/*` | outside the allowlist, no task is created |

## Lifecycle of this probe

1. The branch `test/eleven-slot-9` is pushed and a pull request opened against
   `master`.
2. The `do-not-merge` label is applied at creation, so no actor merges it. This
   repository releases automatically, and a merge would cut a tag and ship.
3. The pull request is reviewed like any other — this file is the diff.
4. Once the peak concurrency has been captured, the pull request is **closed**, not
   merged, and the branch is left behind.

## What a pass looks like

A pass is a sample in which eight reviewer pods are active on `sun` and three on `prod`
at the same instant, with the eleven task ids behind those pods all belonging to the
eleven probe pull requests. A run that reaches nine of eleven and explains the shortfall
as scheduling noise is not a pass; a shortfall is recorded as a shortfall, naming the
limit that produced it.

## Notes on repeatability

The sampler reads only `get pods` and `get resourcequota`, so it is safe to run against
production for the length of the measurement. It samples both clusters on every tick, so
the two rows sharing a timestamp describe the same instant rather than two instants a
network round trip apart.

Revision: 20261028T103703Z — re-run marker so the watcher mints a fresh pr-review task.