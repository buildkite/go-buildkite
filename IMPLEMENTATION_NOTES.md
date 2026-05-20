# Implementation notes

## 2026-05-20

- Created `lox/optional-update-payloads`.
- Added `Optional[T]` and `Some[T]` as the presence primitive for PATCH request structs.
- Confirmed the repo targets Go 1.25, so `json:",omitzero"` can use `Optional[T].IsZero()` to omit unset fields before `MarshalJSON` runs.
- Chose to make explicit `null` count as present in `UnmarshalJSON`; absent and present-null should not collapse into the same state.
- Converted dedicated update structs to `Optional[T]`: clusters, cluster queues, cluster secrets, package registries, pipelines, pipeline schedules, team pipelines, and team suites.
- Added raw request-body assertions for update tests so omitted keys and explicit empty values are visible to tests.
- Added regression coverage for clearing `UpdatePipelineSchedule.Env`, clearing `UpdatePipeline.Tags`, and sending `skip_queued_branch_builds: false` deliberately.
- Split shared create/update structs for cluster tokens, pipeline templates, teams, and test suites. Create structs keep plain values; update structs use `Optional[T]`.
- Fixed `TestSuitesService.Update` to return the request construction error instead of `nil`.
- Added raw request-body coverage for shared-type updates, including explicit false for pipeline template availability and team permission booleans.
