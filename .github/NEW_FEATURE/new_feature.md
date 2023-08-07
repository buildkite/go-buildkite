<type: fix | feat | build | chore | ci | docs | style | refactor | perf | test >(scope)[!]: [description]

## PR checklist:
- [ ] tests added
- [ ] examples of each function added to the relevant `examples/` folder (create if new)
- [ ] `CHANGELOG.md` updated with pending release information

### Example
feat(pipeline service): Add extra parameter to CreatePipeline struct

## PR checklist:
- [ ] tests added
- [x] examples of each function added to the relevant `examples/` folder (create if new)
- [x] `CHANGELOG.md` updated with pending release information

#### To ! or not to !
`!` denotes a breaking change, in this example the test does **not** cause a breaking change and so `!` is not required.

A `BREAKING CHANGE` footer may also be used:

```
feat: strip non-ASCII chars from pipeline name
BREAKING CHANGE: no longer replaces non-ASCII with "blank symbols"
```
