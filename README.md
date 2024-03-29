# elbow

Elbow, Elbow grease.

[![Latest Release](https://img.shields.io/github/release/atc0005/elbow.svg?style=flat-square)](https://github.com/atc0005/elbow/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/atc0005/elbow.svg)](https://pkg.go.dev/github.com/atc0005/elbow)
[![go.mod Go version](https://img.shields.io/github/go-mod/go-version/atc0005/elbow)](https://github.com/atc0005/elbow)
[![Lint and Build](https://github.com/atc0005/elbow/actions/workflows/lint-and-build.yml/badge.svg)](https://github.com/atc0005/elbow/actions/workflows/lint-and-build.yml)
[![Project Analysis](https://github.com/atc0005/elbow/actions/workflows/project-analysis.yml/badge.svg)](https://github.com/atc0005/elbow/actions/workflows/project-analysis.yml)

- [elbow](#elbow)
  - [Project home](#project-home)
  - [Purpose](#purpose)
  - [Gotchas](#gotchas)
  - [Features](#features)
  - [Changelog](#changelog)
  - [Requirements](#requirements)
    - [Building source code](#building-source-code)
    - [Running](#running)
  - [Installation](#installation)
    - [From source](#from-source)
    - [Using release binaries](#using-release-binaries)
  - [Setup test environment](#setup-test-environment)
  - [Configuration Options](#configuration-options)
    - [Precedence](#precedence)
    - [Command-line Arguments](#command-line-arguments)
    - [Environment Variables](#environment-variables)
    - [Configuration File](#configuration-file)
  - [Examples](#examples)
    - [Overview](#overview)
    - [Log output](#log-output)
      - [Text format](#text-format)
        - [Screenshots](#screenshots)
          - [Original implementation](#original-implementation)
          - [Multiple paths](#multiple-paths)
      - [JSON format](#json-format)
    - [Help Output](#help-output)
    - [Prune `.war` files from each branch recursively, keep newest 2](#prune-war-files-from-each-branch-recursively-keep-newest-2)
    - [Keep oldest 1, debug logging, ignore errors, use syslog](#keep-oldest-1-debug-logging-ignore-errors-use-syslog)
    - [Log to a file in JSON format](#log-to-a-file-in-json-format)
  - [License](#license)
  - [References](#references)
    - [Flag packages](#flag-packages)
    - [Configuration object](#configuration-object)
    - [Sorting files](#sorting-files)
    - [Path/File Existence](#pathfile-existence)
    - [Slice management](#slice-management)

## Project home

See [our GitHub repo][repo-url] for the latest code, to file an issue or
submit improvements for review and potential inclusion into the project.

## Purpose

Prune content matching specific patterns, either in a single directory or
recursively through a directory tree. The primary goal is to use this
application from a cron job to perform routine pruning of generated files that
would otherwise completely clog a filesystem.

## Gotchas

- File extensions are *case-sensitive*
- File name patterns are *case-sensitive*
- File name patterns, much like shell globs, may match more than intended.
  - Test carefully and do not provide the `--remove` flag until you have
    tested and are ready to actually prune the content.

## Features

- Supports multiple (merged) sources for supplying configuration settings
  - Environment variables
  - TOML format configuration file
  - Command-line flags (with detailed help output)
  - Note: See the [Precedence](#precedence) list for how multiple
    configuration sources are processed
- Match on specified file patterns
- Flat (single-level) or recursive search
- Process one or many paths
- Age-based threshold for matches (e.g., match files X days old or older)
- Keep a specified number of older or newer matches
- Limit search to specified list of file extensions
- Toggle file removal (read-only by default)
- Extensive, leveled-logging
  - (Optional) Syslog logging (not supported on Windows)
  - (Optional) Logging to a file (if enabled, mutes console output)
  - Text or JSON log formats
- (Optional) Ignore errors encountered when removing files

Worth noting: This project uses Go modules (vs classic `GOPATH` setup)

## Changelog

See the [`CHANGELOG.md`](CHANGELOG.md) file for the changes associated with
each release of this application. Changes that have been merged to `master`,
but not yet an official release may also be noted in the file under the
`Unreleased` section. A helpful link to the Git commit history since the last
official release is also provided for further review.

## Requirements

The following is a loose guideline. Other combinations of Go and operating
systems for building and running tools from this repo may work, but have not
been tested.

### Building source code

- Go
  - see this project's `go.mod` file for *preferred* version
  - this project tests against [officially supported Go
    releases][go-supported-releases]
    - the most recent stable release (aka, "stable")
    - the prior, but still supported release (aka, "oldstable")
- GCC
  - if building with custom options (as the provided `Makefile` does)
- `make`
  - if using the provided `Makefile`

### Running

- Windows 10
- Ubuntu Linux 18.04+
- Red Hat Enterprise Linux 7+

## Installation

### From source

1. [Download][go-docs-download] Go
1. [Install][go-docs-install] Go
   - NOTE: Pay special attention to the remarks about `$HOME/.profile`
1. Clone the repo
   1. `cd /tmp`
   1. `git clone https://github.com/atc0005/elbow`
   1. `cd elbow`
1. Install dependencies (optional)
   - for Ubuntu Linux
     - `sudo apt-get install make gcc`
   - for CentOS Linux
     1. `sudo yum install make gcc`
1. Build
   - for current operating system
     - `go build -mod=vendor ./cmd/elbow/`
       - *forces build to use bundled dependencies in top-level `vendor`
         folder*
   - for all supported platforms (where `make` is installed)
      - `make all`
   - for Windows
      - `make windows`
   - for Linux
     - `make linux`
1. Copy the applicable binary to whatever systems needs to run it
   - if using `Makefile`: look in `/tmp/release_assets/elbow/`
   - if using `go build`: look in `/tmp/elbow/`

**NOTE**: Depending on which `Makefile` recipe you use the generated binary
may be compressed and have an `xz` extension. If so, you should decompress the
binary first before deploying it (e.g., `xz -d elbow-linux-amd64.xz`).

### Using release binaries

1. Download the [latest
   release](https://github.com/atc0005/elbow/releases/latest) binaries
1. Decompress binaries
   - e.g., `xz -d elbow-linux-amd64.xz`
1. Deploy
   - Place `elbow` in a location of your choice
     - e.g., `/usr/local/bin/`

**NOTE**:

DEB and RPM packages are provided as an alternative to manually deploying
binaries.

## Setup test environment

1. Launch container, VM or WSL instance
1. Clone the repo
   1. `cd /tmp`
   1. `git clone https://github.com/atc0005/elbow`
   1. `cd elbow`
1. Create test files
   - in subdirectories of `/tmp/elbow`
     - `make testenv`
   - in a custom location (e.g., in `$HOME`)
     - `mkdir -vp $HOME/tmp/elbow`
     - `make testenv TESTENVBASEDIR="$HOME/tmp/elbow"`

See the [Examples](#examples) or the [Configuration
Options](#configuration-options) sections for examples of running `elbow`
against these newly created test files.

## Configuration Options

### Precedence

**Note: This behavior is subject to change based on feedback.**

The priority order is (mostly):

1. Command line flags (highest priority)
1. Environment variables
1. Environment variables loaded from `.env` files
   - **Not supported yet**
1. Configuration file
1. Default settings (lowest priority)

Configuration sources lower in the list are loaded first, with configuration
sources above loaded sequentially (if enabled) after. Settings are *merged*,
with settings specifically defined in sources with higher precedence
overriding values set by configuration sources with lower precedence.

For example, if the configuration file defines `/tmp/elbow/path1` as the path
to process, an environment variable defines `/tmp/elbow/path2` and the
command-line flag for that setting specifies `/tmp/elbow/path3`, the
command-line flag will win and `/tmp/elbow/path3` will be used.

The intent of this behavior is to provide a *feathered* layering of
configuration settings; if a configuration file provides all settings that you
want other than one, you can use the configuration file for the other settings
and specify the settings that you wish to override via environment variable or
command-line flag.

**Note: This behavior is subject to change based on feedback.**

### Command-line Arguments

Aside from the built-in `-h`, short flag names are currently not supported.

| Long             | Required | Default        | Repeat | Possible                                                                                                | Description                                                                                                                                                    |
| ---------------- | -------- | -------------- | ------ | ------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `keep`           | No       | `0`            | No     | `0+`                                                                                                    | Keep specified number of matching files.                                                                                                                       |
| `paths`          | Yes      | N/A            | No     | *one or more valid directory paths*                                                                     | List of comma or space-separated paths to process.                                                                                                             |
| `pattern`        | No       | *empty string* | No     | *valid file name characters*                                                                            | Substring pattern to compare filenames against. Wildcards are not supported.                                                                                   |
| `extensions`     | No       | *empty list*   | No     | *valid file extensions*                                                                                 | Limit search to specified file extension. Specify as space separated list to match multiple required extensions. Comparisons are performed case-insensitively. |
| `recurse`        | No       | `false`        | No     | `true`, `false`                                                                                         | Perform recursive search into subdirectories.                                                                                                                  |
| `keep-old`       | No       | `false`        | No     | `true`, `false`                                                                                         | Keep oldest files instead of newer.                                                                                                                            |
| `age`            | No       | `0`            | No     | `0+`                                                                                                    | Limit search to files that are the specified number of days old or older.                                                                                      |
| `remove`         | Maybe    | `false`        | No     | `true`, `false`                                                                                         | Remove matched files. The default behavior is to only note what matching files *would* be removed.                                                             |
| `ignore-errors`  | No       | `false`        | No     | `true`, `false`                                                                                         | Ignore errors encountered during file removal.                                                                                                                 |
| `log-format`     | No       | `text`         | No     | `text`, `json`                                                                                          | Log formatter used by logging package.                                                                                                                         |
| `log-file`       | No       | *empty string* | No     | *writable directory path*                                                                               | Optional log file used to hold logged messages. If set, log messages are not displayed on the console.                                                         |
| `console-output` | No       | `stdout`       | No     | `stdout`, `stderr`                                                                                      | Specify how log messages are logged to the console.                                                                                                            |
| `log-level`      | No       | `info`         | No     | `emergency`, `alert`, `critical`, `panic`, `fatal`, `error`, `warn`, `info`, `notice`, `debug`, `trace` | Maximum log level at which messages will be logged. Log messages below this threshold will be discarded.                                                       |
| `use-syslog`     | No       | `false`        | No     | `true`, `false`                                                                                         | Log messages to syslog in addition to other ouputs. Not supported on Windows.                                                                                  |
| `config-file`    | No       | *empty string* | No     | *valid path to config file*                                                                             | Full path to optional TOML-formatted configuration file. See `config.example.toml` for a starter template.                                                     |

### Environment Variables

If set, environment variables override settings provided by a configuration
file. If used, command-line arguments override the equivalent environment
variables listed below. See the [Command-line
Arguments](#command-line-arguments) table for more information.

| Flag Name        | Environment Variable Name | Notes                        | Example                                                                             |
| ---------------- | ------------------------- | ---------------------------- | ----------------------------------------------------------------------------------- |
| `keep`           | `ELBOW_KEEP`              |                              | `ELBOW_KEEP=1`                                                                      |
| `paths`          | `ELBOW_PATHS`             |                              | `ELBOW_PATHS="/tmp/elbow/path1"`, `ELBOW_PATHS="/tmp/elbow/path1,/tmp/elbow/path2"` |
| `pattern`        | `ELBOW_FILE_PATTERN`      |                              | `ELBOW_FILE_PATTERN="reach-masterdev-"`                                             |
| `extensions`     | `ELBOW_EXTENSIONS`        | *Comma-separated, no spaces* | `ELBOW_EXTENSIONS=".war,.tmp"`                                                      |
| `recurse`        | `ELBOW_RECURSE`           |                              | `ELBOW_RECURSE="true"`                                                              |
| `keep-old`       | `ELBOW_KEEP_OLD`          |                              | `ELBOW_KEEP_OLD="true"`                                                             |
| `age`            | `ELBOW_FILE_AGE`          |                              | `ELBOW_FILE_AGE=120`                                                                |
| `remove`         | `ELBOW_REMOVE`            |                              | `ELBOW_REMOVE="false"`                                                              |
| `ignore-errors`  | `ELBOW_IGNORE_ERRORS`     |                              | `ELBOW_IGNORE_ERRORS="true"`                                                        |
| `log-format`     | `ELBOW_LOG_FORMAT`        |                              | `ELBOW_LOG_FORMAT="json"`                                                           |
| `log-file`       | `ELBOW_LOG_FILE`          |                              | `ELBOW_LOG_FILE="/tmp/testing-masterqa-build-removals.txt"`                         |
| `console-output` | `ELBOW_CONSOLE_OUTPUT`    |                              | `ELBOW_CONSOLE_OUTPUT="stdout"`                                                     |
| `log-level`      | `ELBOW_LOG_LEVEL`         |                              | `ELBOW_LOG_LEVEL="debug"`                                                           |
| `use-syslog`     | `ELBOW_USE_SYSLOG`        |                              | `ELBOW_USE_SYSLOG="true"`                                                           |
| `config-file`    | `ELBOW_CONFIG_FILE`       |                              | `ELBOW_CONFIG_FILE="/usr/local/elbow/config.toml"`                                  |

### Configuration File

Configuration file settings have the lowest priority and are overridden by
settings specified in other configuration sources, except for default values.
See the [Command-line Arguments](#command-line-arguments) table for more
information, including the available values for the listed configuration
settings.

| Flag Name        | Config file Setting Name | Section Name   | Notes                                                                    |
| ---------------- | ------------------------ | -------------- | ------------------------------------------------------------------------ |
| `pattern`        | `pattern`                | `filehandling` |                                                                          |
| `extensions`     | `file_extensions`        | `filehandling` |                                                                          |
| `age`            | `file_age`               | `filehandling` |                                                                          |
| `keep`           | `files_to_keep`          | `filehandling` |                                                                          |
| `keep-old`       | `keep_oldest`            | `filehandling` |                                                                          |
| `remove`         | `remove`                 | `filehandling` |                                                                          |
| `ignore-errors`  | `ignore_errors`          | `filehandling` |                                                                          |
| `paths`          | `paths`                  | `search`       | [Multi-line array](https://github.com/toml-lang/toml#user-content-array) |
| `recurse`        | `recursive_search`       | `search`       |                                                                          |
| `log-level`      | `log_level`              | `logging`      |                                                                          |
| `log-format`     | `log_format`             | `logging`      |                                                                          |
| `log-file`       | `log_file_path`          | `logging`      |                                                                          |
| `console-output` | `console_output`         | `logging`      |                                                                          |
| `use-syslog`     | `use_syslog`             | `logging`      |                                                                          |

See the [`config.example.toml`](config.example.toml) file for an example of
how to use these settings.

## Examples

### Overview

The following steps illustrate a rough, overall idea of what `elbow` is
intended to do. The steps illustrate building and running the application from
within an Ubuntu Linux Subsystem for Windows (WSL) instance. The `t` volume is
present on the Windows host.

The file extension used in the examples below is for `WAR` files that are
generated on a build system that our group used to maintain. The idea is that
`elbow` could be run as a cron job to help ensure that only X copies (the most
recent in our case) for each of three branches remain on the build box.

There are better approaches to managing build artifacts (e.g., containers);
this tool seeks to solve in a simple, "low tech" way.

The particular repo that the build system processed has three branches:

| Branch Name | Type of build |
| ----------- | ------------- |
| `master`    | Production    |
| `masterqa`  | Q/A           |
| `masterdev` | Development   |

We had little control over the name of these branches.

### Log output

#### Text format

```ShellSession
./elbow --paths "/tmp/elbow/path1" "/tmp/elbow/path2" --pattern "reach-master" --keep 1 --recurse --keep-old --ignore-errors --log-level info --use-syslog --log-format text
```

```ShellSession
INFO[0000] Syslog logging requested, attempting to enable it  use_syslog=true
ERRO[0000] Failed to enable syslog logging: unable to connect to syslog socket: Unix syslog delivery error  use_syslog=true
WARN[0000] Proceeding without syslog logging             use_syslog=true
INFO[0000] Starting evaluation of paths list             extensions="[]" file_age=0 file_pattern=reach-master paths="[/tmp/elbow/path1 /tmp/elbow/path2]"
INFO[0000] Beginning processing of path "/tmp/elbow/path1" (1 of 2)  ignore_errors=true iteration=1 total_paths=2
INFO[0000] 183 files eligible for removal (910.0 MiB)    extensions="[]" file_age=0 file_pattern=reach-master iteration=1 path=/tmp/elbow/path1 total_file_size=954204160
INFO[0000] 1 files to keep as requested                  iteration=1 keep_oldest=true
INFO[0000] Ignoring file removal errors: true
INFO[0000] File removal not enabled, not removing files
INFO[0000] 0 files successfully removed (0 B)
INFO[0000] 0 files failed to remove (0 B)
INFO[0000] Ending processing of path "/tmp/elbow/path1" (1 of 2)  ignore_errors=true iteration=1 total_paths=2
INFO[0000] Beginning processing of path "/tmp/elbow/path2" (2 of 2)  ignore_errors=true iteration=2 total_paths=2
INFO[0000] 183 files eligible for removal (954.0 KiB)    extensions="[]" file_age=0 file_pattern=reach-master iteration=2 path=/tmp/elbow/path2 total_file_size=976896
INFO[0000] 1 files to keep as requested                  iteration=2 keep_oldest=true
INFO[0000] Ignoring file removal errors: true
INFO[0000] File removal not enabled, not removing files
INFO[0000] 0 files successfully removed (0 B)
INFO[0000] 0 files failed to remove (0 B)
INFO[0000] Ending processing of path "/tmp/elbow/path2" (2 of 2)  ignore_errors=true iteration=2 total_paths=2
INFO[0000] Elbow successfully completed.                 eligible_remove=366 eligible_size="910.9 MiB" failed_removed=0 failed_size="0 B" success_removed=0 success_size="0 B"
```

Where supported, the output is colored.

##### Screenshots

###### Original implementation

From just before the v0.2.0 milestone was completed and the [`v0.2.0`
tag](https://github.com/atc0005/elbow/releases/tag/v0.2.0) created:

![Colored text output example screenshot][screenshot-v0.2.0]

###### Multiple paths

While working on support for multiple paths per [issue
32](https://github.com/atc0005/elbow/issues/32):

![Colored text output example screenshot][screenshot-issue32]

#### JSON format

```ShellSession
./elbow --paths "/tmp/elbow/path1" "/tmp/elbow/path2" --pattern "reach-master" --keep 1 --recurse --keep-old --ignore-errors --log-level info --use-syslog --log-format json
```

```json
{"level":"info","msg":"Syslog logging requested, attempting to enable it","time":"2019-11-12T06:51:50-06:00","use_syslog":true}
{"level":"error","msg":"Failed to enable syslog logging: unable to connect to syslog socket: Unix syslog delivery error","time":"2019-11-12T06:51:50-06:00","use_syslog":true}
{"level":"warning","msg":"Proceeding without syslog logging","time":"2019-11-12T06:51:50-06:00","use_syslog":true}
{"extensions":null,"file_age":0,"file_pattern":"reach-master","level":"info","msg":"Starting evaluation of paths list","paths":["/tmp/elbow/path1","/tmp/elbow/path2"],"time":"2019-11-12T06:51:50-06:00"}
{"ignore_errors":true,"iteration":1,"level":"info","msg":"Beginning processing of path \"/tmp/elbow/path1\" (1 of 2)","time":"2019-11-12T06:51:50-06:00","total_paths":2}
{"extensions":null,"file_age":0,"file_pattern":"reach-master","iteration":1,"level":"info","msg":"183 files eligible for removal (910.0 MiB)","path":"/tmp/elbow/path1","time":"2019-11-12T06:51:50-06:00","total_file_size":954204160}
{"iteration":1,"keep_oldest":true,"level":"info","msg":"1 files to keep as requested","time":"2019-11-12T06:51:50-06:00"}
{"level":"info","msg":"Ignoring file removal errors: true","time":"2019-11-12T06:51:50-06:00"}
{"level":"info","msg":"File removal not enabled, not removing files","time":"2019-11-12T06:51:50-06:00"}
{"level":"info","msg":"0 files successfully removed (0 B)","time":"2019-11-12T06:51:50-06:00"}
{"level":"info","msg":"0 files failed to remove (0 B)","time":"2019-11-12T06:51:50-06:00"}
{"ignore_errors":true,"iteration":1,"level":"info","msg":"Ending processing of path \"/tmp/elbow/path1\" (1 of 2)","time":"2019-11-12T06:51:50-06:00","total_paths":2}
{"ignore_errors":true,"iteration":2,"level":"info","msg":"Beginning processing of path \"/tmp/elbow/path2\" (2 of 2)","time":"2019-11-12T06:51:50-06:00","total_paths":2}
{"extensions":null,"file_age":0,"file_pattern":"reach-master","iteration":2,"level":"info","msg":"183 files eligible for removal (954.0 KiB)","path":"/tmp/elbow/path2","time":"2019-11-12T06:51:50-06:00","total_file_size":976896}
{"iteration":2,"keep_oldest":true,"level":"info","msg":"1 files to keep as requested","time":"2019-11-12T06:51:50-06:00"}
{"level":"info","msg":"Ignoring file removal errors: true","time":"2019-11-12T06:51:50-06:00"}
{"level":"info","msg":"File removal not enabled, not removing files","time":"2019-11-12T06:51:50-06:00"}
{"level":"info","msg":"0 files successfully removed (0 B)","time":"2019-11-12T06:51:50-06:00"}
{"level":"info","msg":"0 files failed to remove (0 B)","time":"2019-11-12T06:51:50-06:00"}
{"ignore_errors":true,"iteration":2,"level":"info","msg":"Ending processing of path \"/tmp/elbow/path2\" (2 of 2)","time":"2019-11-12T06:51:50-06:00","total_paths":2}
{"eligible_remove":366,"eligible_size":"910.9 MiB","failed_removed":0,"failed_size":"0 B","level":"info","msg":"Elbow successfully completed.","success_removed":0,"success_size":"0 B","time":"2019-11-12T06:51:50-06:00"}
```

### Help Output

```ShellSession
$ ./elbow --help
Elbow prunes content matching specific patterns, either in a single directory or recursively through a directory tree.

ELBOW x.y.z
https://github.com/atc0005/elbow

Usage: elbow [--pattern PATTERN] [--extensions EXTENSIONS] [--age AGE] [--keep KEEP] [--keep-old] [--remove] [--ignore-errors] [--log-level LOG-LEVEL] [--log-format LOG-FORMAT] [--log-file LOG-FILE] [--console-output CONSOLE-OUTPUT] [--use-syslog] [--paths PATHS] [--recurse] [--config-file CONFIG-FILE]

Options:
  --pattern PATTERN      Substring pattern to compare filenames against. Wildcards are not supported.
  --extensions EXTENSIONS
                         Limit search to specified file extensions. Specify as space separated list to match multiple required extensions. Comparisons are performed case-insensitively.
  --age AGE              Limit search to files that are the specified number of days old or older.
  --keep KEEP            Keep specified number of matching files per provided path. [default: -1]
  --keep-old             Keep oldest files instead of newer per provided path.
  --remove               Remove matched files per provided path.
  --ignore-errors        Ignore errors encountered during file removal.
  --log-level LOG-LEVEL
                         Maximum log level at which messages will be logged. Log messages below this threshold will be discarded. [default: info]
  --log-format LOG-FORMAT
                         Log formatter used by logging package. [default: text]
  --log-file LOG-FILE    Optional log file used to hold logged messages. If set, log messages are not displayed on the console.
  --console-output CONSOLE-OUTPUT
                         Specify how log messages are logged to the console. [default: stdout]
  --use-syslog           Log messages to syslog in addition to other outputs. Not supported on Windows.
  --paths PATHS          List of comma or space-separated paths to process.
  --recurse              Perform recursive search into subdirectories per provided path.
  --config-file CONFIG-FILE
                         Full path to optional TOML-formatted configuration file. See config.example.toml for a starter template.
  --help, -h             display this help and exit
  --version              display version and exit
```

### Prune `.war` files from each branch recursively, keep newest 2

Note: Leave off `--remove` to display what *would* be removed.

```ShellSession
./elbow --paths "/tmp/elbow/path1" "/tmp/elbow/path2" --extensions ".war" --pattern "reach-master-" --keep 2 --recurse --remove
./elbow --paths "/tmp/elbow/path1" "/tmp/elbow/path2" --extensions ".war" --pattern "reach-masterqa-" --keep 2 --recurse --remove
./elbow --paths "/tmp/elbow/path1" "/tmp/elbow/path2" --extensions ".war" --pattern "reach-masterdev-" --keep 2 --recurse --remove
```

### Keep oldest 1, debug logging, ignore errors, use syslog

Note: Leave off `--remove` to display what *would* be removed.

```ShellSession
./elbow --paths "/tmp/elbow/path1" "/tmp/elbow/path2" --pattern "reach-master-" --keep 1 --recurse --keep-old --ignore-errors --log-level debug --use-syslog
./elbow --paths "/tmp/elbow/path1" "/tmp/elbow/path2" --pattern "reach-masterqa-" --keep 1 --recurse --keep-old --ignore-errors --log-level debug --use-syslog
./elbow --paths "/tmp/elbow/path1" "/tmp/elbow/path2" --pattern "reach-masterdev-" --keep 1 --recurse --keep-old --ignore-errors --log-level debug --use-syslog
```

### Log to a file in JSON format

- These examples attempt to create a log file in the current directory.
- The default logging format is `text` unless overridden; here we specify `json`.
- We attempt to enable syslog logging. This currently fails gracefully on Windows.
- We ignore file removal errors and proceed to the next matching file.

```ShellSession
./elbow --paths "/tmp/elbow/path1" "/tmp/elbow/path2" --pattern "reach-master-" --keep 1 --recurse --keep-old --ignore-errors --log-level debug --log-format json --use-syslog --log-file testing-master-build-removals.txt
./elbow --paths "/tmp/elbow/path1" "/tmp/elbow/path2" --pattern "reach-masterqa-" --keep 1 --recurse --keep-old --ignore-errors --log-level debug --use-syslog --log-format json --log-file testing-masterqa-build-removals.txt
./elbow --paths "/tmp/elbow/path1" "/tmp/elbow/path2" --pattern "reach-masterdev-" --keep 1 --recurse --keep-old --ignore-errors --log-level debug --use-syslog --log-format json --log-file testing-masterdev-build-removals.txt
```

## License

Taken directly from the `LICENSE` and `NOTICE.txt` files:

```License
Copyright 2019-Present Adam Chalkley

https://github.com/atc0005/elbow/blob/master/LICENSE

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
```

## References

### Flag packages

The following list of packages were either used or seriously considered at
some point during the development of this project.

- <https://github.com/jessevdk/go-flags#example>
- <https://github.com/sirupsen/logrus>
- <https://github.com/integrii/flaggy>

### Configuration object

- <https://github.com/go-sql-driver/mysql/blob/877a9775f06853f611fb2d4e817d92479242d1cd/dsn.go#L67>
- <https://github.com/aws/aws-sdk-go/blob/10878ad0389c5b3069815112ce888b191c8cd325/aws/config.go#L251>
- <https://github.com/aws/aws-sdk-go/blob/master/aws/config.go>
- <https://github.com/aws/aws-sdk-go/blob/10878ad0389c5b3069815112ce888b191c8cd325/awstesting/integration/performance/s3GetObject/config.go#L25>
- <https://github.com/aws/aws-sdk-go/blob/10878ad0389c5b3069815112ce888b191c8cd325/awstesting/integration/performance/s3GetObject/main.go#L25>

### Sorting files

- <https://stackoverflow.com/questions/46746862/list-files-in-a-directory-sorted-by-creation-time>

### Path/File Existence

- <https://gist.github.com/mattes/d13e273314c3b3ade33f>

### Slice management

- <https://yourbasic.org/golang/delete-element-slice/>
- <https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang>
- <https://github.com/golang/go/wiki/SliceTricks>

<!-- Footnotes here  -->

[repo-url]: <https://github.com/atc0005/check-elbow>  "This project's GitHub repo"

[go-docs-download]: <https://golang.org/dl>  "Download Go"

[go-docs-install]: <https://golang.org/doc/install>  "Install Go"

[go-supported-releases]: <https://go.dev/doc/devel/release#policy> "Go Release Policy"

[screenshot-v0.2.0]: media/elbow_example_text_log_format_2019-09-26.png "Colored text output example screenshot"
[screenshot-issue32]: media/elbow_example_text_log_format_2019-10-22.png "Colored text output example screenshot"

<!-- []: PLACEHOLDER "DESCRIPTION_HERE" -->
