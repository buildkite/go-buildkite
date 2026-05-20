# Implementation notes

## 2026-05-20

- Created `lox/optional-update-payloads`.
- Added `Optional[T]` and `Some[T]` as the presence primitive for PATCH request structs.
- Confirmed the repo targets Go 1.25, so `json:",omitzero"` can use `Optional[T].IsZero()` to omit unset fields before `MarshalJSON` runs.
- Chose to make explicit `null` count as present in `UnmarshalJSON`; absent and present-null should not collapse into the same state.
