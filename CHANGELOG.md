# Changelog

## Overview

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Please [open an issue](https://github.com/atc0005/elbow/issues) for any
deviations that you spot; I'm still learning!.

## Types of changes

The following types of changes will be recorded in this file:

- `Added` for new features.
- `Changed` for changes in existing functionality.
- `Deprecated` for soon-to-be removed features.
- `Removed` for now removed features.
- `Fixed` for any bug fixes.
- `Security` in case of vulnerabilities.

## [Unreleased]

- placeholder

## [v0.8.1] - 2023-07-17

### Overview

- Bug fixes
- Dependency updates
- GitHub Actions workflow updates
- built using Go 1.19.11
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.19.8` to `1.19.11`
  - `atc0005/go-ci`
    - `go-ci-oldstable-build-v0.10.4` to `go-ci-oldstable-build-v0.11.4`
  - `sirupsen/logrus`
    - `v1.9.0` to `v1.9.3`
  - `pelletier/go-toml`
    - `v2.0.7` to `v2.0.9`
  - `golang.org/x/sys`
    - `v0.7.0` to `v0.10.0`
- (GH-459) Update vuln analysis GHAW to remove on.push hook

### Fixed

- (GH-456) Disable depguard linter
- (GH-461) Restore local CodeQL workflow

## [v0.8.0] - 2023-04-14

### Overview

- Add support for generating DEB, RPM packages
- Build improvements
- Bug fixes
- Generated binary changes
  - filename patterns
  - compression (~ 66% smaller)
  - executable metadata
- built using Go 1.19.8
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Added

- (GH-438) Generate RPM/DEB packages using nFPM
- (GH-441) Add version details to Windows executables

### Changed

- (GH-443) Switch to semantic versioning (semver) compatible versioning
  pattern
- (GH-442) Makefile: Compress binaries & use fixed filenames
- (GH-439) Makefile: Refresh recipes to add "standard" set, new
  package-related options
- (GH-440) Build dev/stable releases using go-ci Docker image
- (GH-444) Move internal packages to internal path

### Fixed

- (GH-436) Fix various revive linting errors
- (GH-437) Fix errcheck linting errors

## [v0.7.24] - 2023-04-13

### Overview

- Bug fixes
- Dependency updates
- GitHub Actions workflow updates
- built using Go 1.19.8
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Added

- (GH-422) Add Go Module Validation, Dependency Updates jobs

### Changed

- Dependencies
  - `Go`
    - `1.19.4` to `1.19.8`
  - `pelletier/go-toml`
    - `v2.0.6` to `v2.0.7`
  - `golang.org/x/sys`
    - `v0.3.0` to `v0.7.0`
- CI
  - (GH-425) Drop `Push Validation` workflow
  - (GH-426) Rework workflow scheduling
  - (GH-428) Remove `Push Validation` workflow status badge

### Fixed

- (GH-432) Update vuln analysis GHAW to use on.push hook

## [v0.7.23] - 2022-12-12

### Overview

- Bug fixes
- Dependency updates
- GitHub Actions Workflows updates
- built using Go 1.19.4
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.19.1` to `1.19.4`
  - `pelletier/go-toml`
    - `v2.0.5` to `v2.0.6`
  - `github.com/alexflint/go-scalar`
    - `v1.1.0` to `v1.2.0`
  - `golang.org/x/sys`
    - `v0.0.0-20220715151400-c0bba94af5f8` to `v0.3.0`
- (GH-412) Refactor GitHub Actions workflows to import logic

### Fixed

- (GH-417) Fix Makefile Go module base path detection
- (GH-403) Add missing cmd doc file

## [v0.7.22] - 2022-09-22

### Overview

- Bug fixes
- Dependency updates
- GitHub Actions Workflows updates
- built using Go 1.19.1
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.17.12` to `1.19.1`
  - `pelletier/go-toml`
    - `v2.0.2` to `v2.0.5`
  - `github/codeql-action`
    - `v2.1.22` to `v2.1.25`
- (GH-404) Update project to Go 1.19
- (GH-405) Update Makefile and GitHub Actions Workflows

### Fixed

- (GH-402) Apply linting fixes for Go 1.19 release
- (GH-403) Add missing cmd doc file

## [v0.7.21] - 2022-07-22

### Overview

- Bug fixes
- Dependency updates
- built using Go 1.17.12
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.17.10` to `1.17.12`
  - `sirupsen/logrus`
    - `v1.8.1` to `v1.9.0`
  - `pelletier/go-toml`
    - `v2.0.1` to `v2.0.2`

### Fixed

- (GH-394) Update lintinstall Makefile recipe

## [v0.7.20] - 2022-05-11

### Overview

- Dependency updates
- built using Go 1.17.10
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.17.9` to `1.17.10`
  - `pelletier/go-toml`
    - `v2.0.1-0.20220509164502-c5ca2c682b57` to `v2.0.1`

## [v0.7.19] - 2022-05-10

### Overview

- Dependency updates
- built using Go 1.17.9
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.17.7` to `1.17.9`
  - `pelletier/go-toml`
    - `v1.9.4` to `v2.0.1-0.20220509164502-c5ca2c682b57`

## [v0.7.18] - 2022-03-02

### Overview

- Dependency updates
- CI / linting improvements
- built using Go 1.17.7
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.12` to `1.17.7`
    - (GH-371) Update go.mod file, canary Dockerfile to reflect current
      dependencies
  - `alexflint/go-arg`
    - `v1.4.2` to `v1.4.3`
  - `actions/checkout`
    - `v2.4.0` to `v3`
  - `actions/setup-node`
    - `v2.5.1` to `v3`

- (GH-374) Expand linting GitHub Actions Workflow to include `oldstable`,
  `unstable` container images
- (GH-375) Switch Docker image source from Docker Hub to GitHub Container
  Registry (GHCR)

## [v0.7.17] - 2021-12-29

### Overview

- Dependency updates
- built using Go 1.16.12
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.10` to `1.16.12`
  - `actions/setup-node`
    - `v2.4.1` to `v2.5.1`

## [v0.7.16] - 2021-11-10

### Overview

- Dependency updates
- built using Go 1.16.10
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.7` to `1.16.10`
  - `pelletier/go-toml`
    - `v1.9.3` to `v1.9.4`
  - `actions/checkout`
    - `v2.3.4` to `v2.4.0`
  - `actions/setup-node`
    - `v2.4.0` to `v2.4.1`

### Fixed

- (GH-362) False positive `G307: Deferring unsafe method "Close" on type
  "*os.File" (gosec)` linting error
- (GH-350) Build tag format changed between Go 1.16 and 1.17
- (GH-348) Lock Go version to the latest "oldstable" series

## [v0.7.15] - 2021-08-09

### Overview

- Dependency updates
- built using Go 1.16.7
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.6` to `1.16.7`
  - `actions/setup-node`
    - updated from `v2.2.0` to `v2.4.0`

## [v0.7.14] - 2021-07-20

### Overview

- Dependency updates
- Bug fixes
- built using Go 1.16.6
  - Statically linked
  - Linux (x86, x64)

### Added

- Add "canary" Dockerfile to track stable Go releases, serve as a reminder to
  generate fresh binaries

### Changed

- Show the date for the specific number of days when using the `--age` flag

- Dependencies
  - `Go`
    - `1.16.2` to `1.16.6`
  - `alexflint/go-arg`
    - `v1.3.0` to `v1.4.2`
  - `pelletier/go-toml`
    - `v1.8.1` to `v1.9.3`
  - `actions/setup-node`
    - `v2.1.5` to `v2.2.0`
    - update `node-version` value to always use latest LTS version instead of
      hard-coded version

### Fixed

- cmd/elbow/main.go:97:17: ST1023: should omit type int from declaration; it
  will be inferred from the right-hand side (stylecheck)
- Fix doc comment field references

## [v0.7.13] - 2021-03-15

### Overview

- Dependency updates
- Bug fixes
- Built using Go 1.16.2

### Changed

- Compare file extensions case-insensitively

- Dependencies
  - Built using Go 1.16.2
    - **Statically linked**
    - Windows (x86, x64)
    - Linux (x86, x64)
  - `sirupsen/logrus`
    - `v1.8.0` to `v1.8.1`
  - `actions/setup-node`
    - `v2.1.4` to `v2.1.5`

### Fixed

- Prune unneeded `\n` escape character from log messages

- Compare file extensions without leading dot

## [v0.7.12] - 2021-02-21

### Overview

- Dependency updates
- Built using Go 1.15.8

### Changed

- Swap out GoDoc badge for pkg.go.dev badge

- Dependencies
  - Built using Go 1.15.8
    - **Statically linked**
    - Windows (x86, x64)
    - Linux (x86, x64)
  - `sirupsen/logrus`
    - `v1.7.0` to `v1.8.0`
  - `actions/setup-node`
    - `v2.1.2` to `v2.1.4`

## [v0.7.11] - 2020-11-16

### Changed

- Binary release
  - Built using Go 1.15.5
  - **Statically linked**
  - Windows
    - x86
    - x64
  - Linux
    - x86
    - x64

### Fixed

- Logic error in PathExists function
- LogBuffer.Flush method attempt to empty slice
- Minor doc comment typo

## [v0.7.10] - 2020-10-11

### Added

- Binary release
  - Built using Go 1.15.2
  - **Statically linked**
  - Windows
    - x86
    - x64
  - Linux
    - x86
    - x64

### Changed

- Dependencies
  - `actions/setup-node`
    - `v2.1.1` to `v2.1.2`

- Add `-trimpath` build flag

### Fixed

- Makefile build options do not generate static binaries
- Makefile generates checksums with qualified path

## [v0.7.9] - 2020-09-29

### Added

- First (limited) binary release
  - Built using Go 1.15.2
  - Windows
    - x86
    - x64
  - Linux
    - x86
    - x64

### Changed

- Dependencies
  - upgrade `actions/checkout`
    - `v2.3.1` to `v2.3.3`
  - upgrade `sirupsen/logrus`
    - `v1.6.0` to `v1.7.0`
  - upgrade `pelletier/go-toml`
    - `v1.8.0` to `v1.8.1`

### Fixed

- Miscellaneous linting issues
  - comment formatting
  - exitAfterDefer: os.Exit clutters defer
- Makefile generates checksums with qualified path

## [v0.7.8] - 2020-08-03

### Added

- Docker-based GitHub Actions Workflows
  - Replace native GitHub Actions with containers created and managed through
    the `atc0005/go-ci` project.

  - New, primary workflow
    - with parallel linting, testing and building tasks
    - with three Go environments
      - "old stable" - currently `Go 1.13.14`
      - "stable" - currently `Go 1.14.6`
      - "unstable" - currently `Go 1.15rc1`
    - Makefile is *not* used in this workflow
    - staticcheck linting using latest stable version provided by the
      `atc0005/go-ci` containers

  - Separate Makefile-based linting and building workflow
    - intended to help ensure that local Makefile-based builds that are
      referenced in project README files continue to work as advertised until
      a better local tool can be discovered/explored further
    - use `golang:latest` container to allow for Makefile-based linting
      tooling installation testing since the `atc0005/go-ci` project provides
      containers with those tools already pre-installed
      - linting tasks use container-provided `golangci-lint` config file
        *except* for the Makefile-driven linting task which continues to use
        the repo-provided copy of the `golangci-lint` configuration file

  - Add Quick Validation workflow
    - run on every push, everything else on pull request updates
    - linting via `golangci-lint` only
    - testing
    - no builds

### Changed

- README
  - Link badges to applicable GitHub Actions workflows results

- Linting
  - local
    - `golangci-lint`
      - disable default exclusions
    - `Makefile`
      - install latest stable `golangci-lint` binary instead of using a fixed
        version
  - CI
    - remove repo-provided copy of `golangci-lint` config file at start of
      linting task in order to force use of Docker container-provided config
      file

- Dependencies
  - upgrade `actions/setup-node`
    - `v2.1.0` to `v2.1.1`
  - upgrade `actions/setup-go`
    - `v2.1.0` to `v2.1.1`
    - note: since replaced with a Docker container

### Fixed

- Miscellaneous linting issues
  - `errcheck`
  - `gosec`
    - log file permissions
    - file inclusion via variable
  - `stylecheck`

## [v0.7.7] - 2020-07-19

### Added

- Dependabot
  - enable version updates
  - enable GitHub Actions updates

### Changed

- Dependencies
  - upgrade `pelletier/go-toml`
    - `v1.7.0` to `v1.8.0`
  - upgrade `actions/setup-go`
    - `v1` to `v2.1.0`
  - upgrade `actions/checkout`
    - `v1` to `v2.3.1`
  - upgrade `actions/setup-node`
    - `v1` to `v2.1.0`

### Fixed

- Remove duplicate defaultAppName const
- Fix CHANGELOG section order

## [v0.7.6] - 2020-05-03

### Changed

- `sirupsen/logrus` updated from `v1.5.0` to `v1.6.0`

### Fixed

- Version string/output was not shown when `-h` or `--version` flags were used

## [v0.7.5] - 2020-04-29

### Changed

- Update golangci-lint to v1.25.1
- Enable `gofmt` linter in golangci-lint config file

### Fixed

- Linting errors
  - Unused field in table test struct
  - Remove explicit struct type for each struct entry in table test slice

- Remove gofmt and golint as separate checks, enable these linters in
  golangci-lint config

- Update README to list accurate build/deploy steps based on recent
  restructuring work

## [v0.7.4] - 2020-04-26

### Changed

- Use `cmd/elbow` directory to match best practices

- Vendor dependencies

- README
  - update one-off build examples to include new cmd subdirectory

- Update GitHub Actions Workflows
  - specify new cmd subdir path for builds
  - Disable running `go get` after checking out code
  - Exclude `vendor` folder from ...
  - Markdown linting checks
  - tests
  - basic build

- Update `.gitignore`
  - add ignored paths for binaries
  - add `release_assets`

- Update Makefile
  - replace two external shell scripts with equivalent embedded commands
  - borrow heavily from existing `Makefile` for `atc0005/send2teams` project
  - generate binaries within `release_assets` subdirectory structure
  - dynamically determine go module path for version tag use
    - disabled for now as I have not moved this into a subpkg (e.g., `config`)
      yet
  - include `-mod=vendor` build flag for applicable `go` commands to reflect
    Go 1.13 vendoring
    - this includes specifying `-mod=vendor` even for `go list` commands,
      which unless specified results in dependencies being downloaded, even
      when they're already provided in a local, top-level `vendor` directory
  - dynamic help/menu output based on recipe "doc comment"

- Linting
  - Enabled `scopelint` linter
  - Moved `golangci-lint` config settings to external file

- Copyright date bump

### Fixed

- Linting
  - dogsled linting warnings regarding use of `runtime.Caller(1)`
    - applied `// nolint:dogsled` hotfix (GH-237)

## [v0.7.3] - 2020-04-26

### Changed

- GitHub Actions Workflow
  - Go v1.12.x dropped
  - Go v1.14.x added
  - Echo Go version used by workflow

- Dependencies updated
  - `pelletier/go-toml`
    - `v1.6.0` to `v1.7.0`
  - `sirupsen/logrus`
    - `v1.4.2` to `v1.5.0`
  - `alexflint/go-arg`
    - `v1.2.0` to `v1.3.0`

- Install `golangci-lint` binary via `make lintinstall` instead of building
  from source

### Fixed

- Correct filename reference

## [v0.7.2] - 2020-02-12

### Fixed

- Fix unhandled error condition by enforcing early exit as intended

## [v0.7.1] - 2020-01-16

### Fixed

- Fix release links in this CHANGELOG file

## [v0.7.0] - 2020-01-16

### Added

- Add support for importing settings from a TOML configuration file
- Add support for merging multiple config sources
  - the defined precedence decides what config source wins for a specific setting in case of a
  conflict.
  - non-conflicting settings are merged into one comprehensive configuration
- Extend GitHub Actions Workflow to include Markdown linting
- Add logger object for use in pre-config initialization
  - this allows delaying the filtering of leveled log messages until the user
    has indicated what logging level they prefer
- Add tests for bulk of core config source handling which includes validation
  of options, precedence/merge rules and other functionality
  - TODO: additional coverage is needed
  - TODO: many of the tests are inefficient and need further work
- Expand list of golangci-lint linters
  - `prealloc`
  - `misspell`
  - `maligned`
  - `dupl`
  - `unconvert`
  - `golint`
  - `gocritic`

### Changed

- Increase use of constants for common string comparison cases
- Configuration source precedence/priority has changed
  1. Command line flags (highest priority)
  1. Environment variables
  1. Environment variables loaded from `.env` files
      - **Not supported yet**
  1. Configuration file
  1. Default settings (lowest priority)

### Fixed

- "Successfully completed" message shown when failures occur during file removal
- Non-recursive directory processing sets wrong value for `FileMatch.Path`
- Anonymous function for `filepath.Walk()` doesn't check `Config.IgnoreErrors`
  before returning error
- Non-recursive directory processing attempts to open a directory as a file
- "NotSet successfully completed" message shown at end of test run
- Miscellaneous doc comments for updated function outdated
- Various linting errors for lack of constants, misspellings, inefficient logic
- `linting` Makefile recipe returns zero exit code even when `golangci-lint`
  reports problems
- Minor typos in Makefile output

## [v0.6.2] - 2019-11-04

### Fixed

- Add missing ALT text to CI badge

## [v0.6.1] - 2019-10-31

### Fixed

- Remove stray phrase from CHANGELOG
- Update build status badge on README to reflect recent workflow name change

## [v0.6.0] - 2019-10-31

### Added

- GitHub Actions Workflow
  - Run linting checks
  - Run build checks

- Documentation
  - Add CI badge to README to note current codebase state

- Makefile
  - new recipe: `make linting`
    - runs suite of linting checks, suggests `make lintinstall` if tools are
      not found
  - new recipe: `make lintinstall`
    - this recipe is used by the GitHub Actions Workflow as part of pre-test setup
    - this same recipe can be used locally on-demand as part of a
      pre-commit/pre-push workflow

- Report total size of files (eligible for removal, failed and success) for
  each path that is processed

### Changed

- Makefile
  - `make all` now builds x86 binaries for Linux and Windows in addition to
    the existing x64 binaries

### Fixed

- Additional godoc coverage for subpackages
- Fix various linting errors exposed by CI-related work

## [v0.5.2] - 2019-10-24

### Fixed

- nil pointer dereference due to not initializing the logger handle before use

## [v0.5.1] - 2019-10-24

### Changed

- Logging level for file removal "intent" messages changed from INFO to DEBUG
  level in order to de-duplicate the coverage for removing files (e.g., noting
  before and after)
- Refactored `main` package to create multiple sub-packages in the hope that
  this will make it easier to create unit tests later
- Apply default settings via `default` struct tag
- `alexflint/go-arg` package updated to v1.2.0
- README, godoc coverage
  - updates to reflect updated Help output
- the `--keep` command-line flag defaults to `0` instead of being a required
  flag

### Removed

- `required` constraint for the `--keep` flag (now defaults to `0`)

### Fixed

- README
  - `--age` command-line flag notes had description details in the wrong column
- syntax of godoc documentation so that it properly renders on godoc.org or
  local godoc instance
- golangci-lint linting errors
  - unintentional duplicate `arg` struct tags
  - unused function (refactored to separate package)
- `Makefile`
  - minor tweaks to output to adjust spacing after removal of UPX calls

## [v0.5.0] - 2019-10-23

### Added

- Add brief first draft of `godoc` compatible usage information
- Support for processing  one or more paths via command-line or environment
  variable
- Logging output
  - total paths provided
  - current iteration in loop across provided paths
  - ending result noted as successful completion or completion with warnings
  - misc logging tweaks to record additional field values that could be useful
    when troubleshooting

### Changed

- Updated README to cover new support for processing multiple paths
  - examples now reference `/tmp/elbow` as the base path with
    `/tmp/elbow/path1` and `/tmp/elbow/path2` as example multi-path arguments
- single path command-line flag `--path` replaced by multi-path `--paths`
  command-line flag
- `ELBOW_PATH` environment variable replaced by `ELBOW_PATHS` which now
  supports comma-separated list of paths
- `Makefile`
  - `Makefile` and test scripts updated to work with multiple paths
  - `make testenv` now prepares `/tmp/elbow/path1` and `/tmp/elbow/path2` by
    default
  - override variable exposed as `TESTENVBASEDIR` (covered in README)
  - UPX is no longer used to compress generated executables

### Removed

- Use of UPX for compressing executables
  - using UPX to compress executables disables use of `go version` and `go
    version -m -v` to determine the version of Go and associated modules used
    to build them
- `--path` command-line flag (see earlier notes)
- `ELBOW_PATH` environment variable (see earlier notes)

### Fixed

- Early exit logic
  - hard-coded `os.Exit(1)` calls (from before support for ignoring errors
    were added) were adjusted to respect the `IgnoreErrors` option
- README
  - Add missing `--age` command-line flag option
  - Add missing `ELBOW_FILE_AGE` environment variable

## [v0.4.0] - 2019-10-17

### Added

- Latest Release badge on README
- Support for environment variables via `alexflint/go-arg` package
- `Makefile`
  - command: `testrun`
  - Set `version` global variable in `main` package based on `git describe`

### Changed

- `--extension` (multi-use) flag is now `--extensions` (single call, multiple values supported
  - See [README](README.md) for usage
- Replaced `jessevdk/go-flags` package with `alexflint/go-arg`
- Improve configuration validation to accommodate lack of native `go-arg`
  support for enforcing specific flag values
- `Makefile`
  - TODO: Add more info here
- `go.mod` updated to use Go 1.13 as the base version
  - Based on some reading in <https://github.com/golang/go/wiki/Modules>, the
    behavior for `go get -u` changed to allow more conservative updates of
    dependencies. The new behavior sounds more natural and is less likely to
    surprise newcomers, so locking the base behavior to Go 1.13 sounds like a
    "Good Thing" to do here.
- README updated to note Go 1.13 as the base version

### Removed

- `jessevdk/go-flags` package replaced with `alexflint/go-arg`

### Fixed

- Typo in license text
- Replace lightweight Git tags with annotated tags

## [v0.3.2] - 2019-10-16

### Fixed

- README: Add package dependencies and install steps
- Fix miscellaneous spelling errors
  - credit: `Code Spell Checker` VSCode extension
- Update `Config.String()` method to include new fields
- Log config field values after setting logging level so that they're visible
  when choosing to log at `debug` level
- Remove placeholder text from README file that has since been superseded by
  real content
- Fix file removal bug by using fully-qualified path to file instead of
  shortname
  - the bug was due to an unintentional assumption that the file to be removed
    was within the current working directory

## [v0.3.1] - 2019-09-29

### Changed

- Update `Makefile` (and the called `testing/setup_testenv.sh` script) to allow
  for providing a custom location for generation of test files)

### Fixed

- Minor tweak to logging output to reduce duplication between main log message
  and the structured field

## [v0.3.0] - 2019-09-27

### Added

- Add Makefile
  - handle cleaning the build environment, the local Git repo and other
    temporary content
  - handle building binaries for Windows and Linux
- Add test environment setup script
  - used by `Makefile`; usable separately if desired
- (optional) Match files based on age threshold (in number of days)

### Changed

- Update .gitignore file to exclude temporary files generated by UPX
- Update README to provide build and, test environment setup directions
- Refactor threshold helper functions in an effort to more clearly reflect
  their purpose
- Update logging to include more structured fields

### Fixed

- Fix link to section in README

## [v0.2.0] - 2019-09-26

### Added

Documentation:

- `LICENSE` file
- `CHANGELOG.md` file
- `NOTICE.txt` file
- License header to all source files

Logging:

- Apply leveled logging to better filter desired logging levels
- Add (optional on Linux, unavailable on Windows) syslog logging support
- Add optional log file support

Configuration options:

- (optional) Ignore errors when removing files
- (optional) Log format (text or json, defaults to text)
- (optional) Log level (large list, mapped where possible to syslog logging
  levels)
- (optional) Console output toggle (stdout or stderr)
- (optional) Log file path (logging to a log file mutes console output)

### Changed

- Extensive updates to main `README.md` file
- Additional polish for "feedback" log statements; work towards having all
  required information set to INFO log level (which is the default)
- Use `jessevdk/go-flags` for command-line flag support
- Short flag names dropped
  - There are some issues with `go-flags` misdetecting leading dashes in file
    patterns as short flags, so instead of dealing with that right now I've
    opted to only support long flag names
  - `go-flags` only supports single letter short flags, and with the number of
    flags that we're using I decided to keep things simple for now and only
    use long flag names
- The number of files to keep from match results is now a required flag

### Removed

- Feature: Short flag names
- Package: `integrii/flaggy`
- Package: `r3labs/diff`

## [v0.1.0] - 2019-09-17

### Added

This initial prototype supports:

- Command-line flags support via `integrii/flaggy` package
- Matching on specified file patterns
- Flat (single-level) or recursive search
- Keeping a specified number of older or newer matches
- Limiting search to specified list of extensions
- Toggling file removal (read-only by default)
- Go modules (vs classic GOPATH setup)
- Brief overview, examples for testing purposes

[Unreleased]: https://github.com/atc0005/elbow/compare/v0.8.1...HEAD
[v0.8.1]: https://github.com/atc0005/elbow/releases/tag/v0.8.1
[v0.8.0]: https://github.com/atc0005/elbow/releases/tag/v0.8.0
[v0.7.24]: https://github.com/atc0005/elbow/releases/tag/v0.7.24
[v0.7.23]: https://github.com/atc0005/elbow/releases/tag/v0.7.23
[v0.7.22]: https://github.com/atc0005/elbow/releases/tag/v0.7.22
[v0.7.21]: https://github.com/atc0005/elbow/releases/tag/v0.7.21
[v0.7.20]: https://github.com/atc0005/elbow/releases/tag/v0.7.20
[v0.7.19]: https://github.com/atc0005/elbow/releases/tag/v0.7.19
[v0.7.18]: https://github.com/atc0005/elbow/releases/tag/v0.7.18
[v0.7.17]: https://github.com/atc0005/elbow/releases/tag/v0.7.17
[v0.7.16]: https://github.com/atc0005/elbow/releases/tag/v0.7.16
[v0.7.15]: https://github.com/atc0005/elbow/releases/tag/v0.7.15
[v0.7.14]: https://github.com/atc0005/elbow/releases/tag/v0.7.14
[v0.7.13]: https://github.com/atc0005/elbow/releases/tag/v0.7.13
[v0.7.12]: https://github.com/atc0005/elbow/releases/tag/v0.7.12
[v0.7.11]: https://github.com/atc0005/elbow/releases/tag/v0.7.11
[v0.7.10]: https://github.com/atc0005/elbow/releases/tag/v0.7.10
[v0.7.9]: https://github.com/atc0005/elbow/releases/tag/v0.7.9
[v0.7.8]: https://github.com/atc0005/elbow/releases/tag/v0.7.8
[v0.7.7]: https://github.com/atc0005/elbow/releases/tag/v0.7.7
[v0.7.6]: https://github.com/atc0005/elbow/releases/tag/v0.7.6
[v0.7.5]: https://github.com/atc0005/elbow/releases/tag/v0.7.5
[v0.7.4]: https://github.com/atc0005/elbow/releases/tag/v0.7.4
[v0.7.3]: https://github.com/atc0005/elbow/releases/tag/v0.7.3
[v0.7.2]: https://github.com/atc0005/elbow/releases/tag/v0.7.2
[v0.7.1]: https://github.com/atc0005/elbow/releases/tag/v0.7.1
[v0.7.0]: https://github.com/atc0005/elbow/releases/tag/v0.7.0
[v0.6.2]: https://github.com/atc0005/elbow/releases/tag/v0.6.2
[v0.6.1]: https://github.com/atc0005/elbow/releases/tag/v0.6.1
[v0.6.0]: https://github.com/atc0005/elbow/releases/tag/v0.6.0
[v0.5.2]: https://github.com/atc0005/elbow/releases/tag/v0.5.2
[v0.5.1]: https://github.com/atc0005/elbow/releases/tag/v0.5.1
[v0.5.0]: https://github.com/atc0005/elbow/releases/tag/v0.5.0
[v0.4.0]: https://github.com/atc0005/elbow/releases/tag/v0.4.0
[v0.3.2]: https://github.com/atc0005/elbow/releases/tag/v0.3.2
[v0.3.1]: https://github.com/atc0005/elbow/releases/tag/v0.3.1
[v0.3.0]: https://github.com/atc0005/elbow/releases/tag/v0.3.0
[v0.2.0]: https://github.com/atc0005/elbow/releases/tag/v0.2.0
[v0.1.0]: https://github.com/atc0005/elbow/releases/tag/v0.1.0
