## [v2.4.0](https://github.com/buildkite/go-buildkite/tree/v2.4.0) (2019-02-15)
[Full Changelog](https://github.com/buildkite/go-buildkite/compare/v2.1.0...v2.4.0)

### Changed
- Added support for rebuilding builds [#27](https://github.com/buildkite/go-buildkite/pull/27) (@kevsmith)
- Add typed provider settings [#23](https://github.com/buildkite/go-buildkite/pull/23) (@haines)
- Add blocked to build struct [#30](https://github.com/buildkite/go-buildkite/pull/30) (@jsm)
- Add JobsService with unblockJob method [#31](https://github.com/buildkite/go-buildkite/pull/31) (@dschofie)
- Support Job.runnable_at [#26](https://github.com/buildkite/go-buildkite/pull/26) (@lox)
- Parse AgentQueryRules and TimeoutInMinutes for Steps [#25](https://github.com/buildkite/go-buildkite/pull/25) (@lox)
- Parse the timestamp format webhooks use [#19](https://github.com/buildkite/go-buildkite/pull/19) (@euank)

## [v2.3.1](https://github.com/buildkite/go-buildkite/tree/v2.3.1) (2020-02-19)
[Full Changelog](https://github.com/buildkite/go-buildkite/compare/v2.3.0...v2.3.1)

### Changed
- Go Modules: /v2 suffix on module path [#63](https://github.com/buildkite/go-buildkite/pull/63) (@pda)

## [v2.3.0](https://github.com/buildkite/go-buildkite/tree/v2.3.0) (2020-02-18)
[Full Changelog](https://github.com/buildkite/go-buildkite/compare/v2.2.0...v2.3.0)

### Changed
- Add Annotations API support [#62](https://github.com/buildkite/go-buildkite/pull/62) (@toolmantim)
- Add the `step_key` field to Jobs [#61](https://github.com/buildkite/go-buildkite/pull/61) (@toolmantim)
- Add header times to job log [#59](https://github.com/buildkite/go-buildkite/pull/59) (@kushmansingh)
- run tests on go 1.13 [#58](https://github.com/buildkite/go-buildkite/pull/58) (@yob)
- Fix tests: compile error, segfault. [#57](https://github.com/buildkite/go-buildkite/pull/57) (@pda)
- Go modules; delete vendor/ [#56](https://github.com/buildkite/go-buildkite/pull/56) (@pda)
- License/Copyright tidy [#55](https://github.com/buildkite/go-buildkite/pull/55) (@pda)
- Add functionality for teams endpoint [#54](https://github.com/buildkite/go-buildkite/pull/54) (@pyasi)
- Add client func for the Artifacts ListByJob endpoint [#53](https://github.com/buildkite/go-buildkite/pull/53) (@aashishkapur)
- Add support for the Job ArtifactsURL field [#52](https://github.com/buildkite/go-buildkite/pull/52) (@aashishkapur)
- Added include_retried_jobs for BuildsListOptions [#51](https://github.com/buildkite/go-buildkite/pull/51) (@NorseGaud)
- Add support for getting a job's log output [#48](https://github.com/buildkite/go-buildkite/pull/48) (@y-yagi)
- Fixed spelling errors in comment [#46](https://github.com/buildkite/go-buildkite/pull/46) (@eightseventhreethree)
- Add YAML support to types [#44](https://github.com/buildkite/go-buildkite/pull/44) (@gdhagger)
- Include `source` field for the builds API response [#49](https://github.com/buildkite/go-buildkite/pull/49) (@angulito)
- Add missing AgentListOptions. [#50](https://github.com/buildkite/go-buildkite/pull/50) (@philwo)
-  Add filtering option for listing builds by commit [#47](https://github.com/buildkite/go-buildkite/pull/47) (@srmocher)
- Add cancel build function [#45](https://github.com/buildkite/go-buildkite/pull/45) (@dschofie)
- Add pull request build information to CreateBuild [#43](https://github.com/buildkite/go-buildkite/pull/43) (@jradtilbrook)
- Allow updating pipeline provider settings [#41](https://github.com/buildkite/go-buildkite/pull/41) (@gdhagger)
- Explicitly set JSON content-type on API requests [#40](https://github.com/buildkite/go-buildkite/pull/40) (@matthewd)

## [v2.2.0](https://github.com/buildkite/go-buildkite/tree/v2.2.0) (2019-02-20)
[Full Changelog](https://github.com/buildkite/go-buildkite/compare/v2.1.0...v2.2.0)

### Changed
- Update comments to have func name, run gofmt and golint [#36](https://github.com/buildkite/go-buildkite/pull/36) (@charlottestjohn)
- Added support for rebuilding builds [#27](https://github.com/buildkite/go-buildkite/pull/27) (@kevsmith)
- Add typed provider settings [#23](https://github.com/buildkite/go-buildkite/pull/23) (@haines)
- Add blocked to build struct [#30](https://github.com/buildkite/go-buildkite/pull/30) (@jsm)
- Add JobsService with unblockJob method [#31](https://github.com/buildkite/go-buildkite/pull/31) (@dschofie)
- Support Job.runnable_at [#26](https://github.com/buildkite/go-buildkite/pull/26) (@lox)
- Parse AgentQueryRules and TimeoutInMinutes for Steps [#25](https://github.com/buildkite/go-buildkite/pull/25) (@lox)
-  Parse the timestamp format webhooks use [#19](https://github.com/buildkite/go-buildkite/pull/19) (@euank)
- Pipeline CRUD [#22](https://github.com/buildkite/go-buildkite/pull/22) (@mubeta06)
- Actually bump version to 2.1.1 [#17](https://github.com/buildkite/go-buildkite/pull/17) (@lox)

## [v2.1.0](https://github.com/buildkite/go-buildkite/tree/v2.1.0) (2017-11-23)
[Full Changelog](https://github.com/buildkite/go-buildkite/compare/v2.0.0...v2.1.0)

### Changed
- Add retry with exp backoff for GET/429 responses [#16](https://github.com/buildkite/go-buildkite/pull/16) (@lox)
- Add Create() to BuildsService [#11](https://github.com/buildkite/go-buildkite/pull/11) (@bshi)
- Update README.md [#15](https://github.com/buildkite/go-buildkite/pull/15) (@lox)
- Fix the build and bump to 1.7.3 [#13](https://github.com/buildkite/go-buildkite/pull/13) (@toolmantim)
- Add Artifacts Service [#10](https://github.com/buildkite/go-buildkite/pull/10) (@pquerna)
- Removed the deprecated/dead featured_build pipeline property [#7](https://github.com/buildkite/go-buildkite/pull/7) (@tobyjoe)
- Added BadgeURL property to Pipeline [#6](https://github.com/buildkite/go-buildkite/pull/6) (@tobyjoe)
- Changed the wolfeidau refs to buildkite [#9](https://github.com/buildkite/go-buildkite/pull/9) (@tobyjoe)

## [v2.0.0](https://github.com/buildkite/go-buildkite/tree/v2.0.0) (2016-03-31)
[Full Changelog](https://github.com/buildkite/go-buildkite/compare/v1.0.0...v2.0.0)

### Changed
- Changes for v2 [#5](https://github.com/buildkite/go-buildkite/pull/5) (@lox)

## [v1.0.0](https://github.com/buildkite/go-buildkite/tree/v1.0.0) (2016-01-10)
[Full Changelog](https://github.com/buildkite/go-buildkite/compare/44957ec...v1.0.0)

- Initial implementation
