# Phase 6 — Verification + cleanup

## Goal

Final green-check across the monorepo, docs updated, and anything that remains manual clearly flagged for the user.

## Tasks

- 6.1 `make test` green in every sub-project (automatable).
- 6.2 `make dev.start` smoke — manual; hand to user. Verify POST /v1/tags round-trip via curl or UI.
- 6.3 `atlas migrate status` shows composed migrations in order: 0000–00NN (core), 0100–01NN (users), 0200–0201 (tags) — manual.
- 6.4 grep sanity: `grep -rn "tags" core-module/model` returns no references (tags live in tags-module only); `grep -rn "tags" users-module/model/migrations` returns nothing (tags migrations live only in tags-module/model; users-module only composes).
- 6.5 Audit log entries for tag create/update/delete — manual.
- 6.6 Update root README + users-module summary to mention tags-module dependency.

## Deliverables

- `user-components/README.md` includes tags-module in module layout (done in Phase 1 Task 6; re-check).
- `users-module/plan/summary.md` gains a "Depends on tags-module" note (new, brief).
- All TODO.md items in this plan checked off or explicitly flagged as manual-pending.

## Reports

Write progress notes into `report.<N>.<topic>.md` in this directory as work proceeds. At minimum, write `report.6.smoke.md` summarizing which verify items passed automated checks and which require user smoke.
