# elbow

Elbow, Elbow grease.

- [elbow](#elbow)
  - [Project home](#project-home)
  - [Purpose](#purpose)
  - [Gotchas](#gotchas)
  - [Features](#features)
  - [Changelog](#changelog)
  - [Requirements](#requirements)
  - [How to install it](#how-to-install-it)
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
      - [JSON format](#json-format)
    - [Help Output](#help-output)
    - [Prune `.war` files from each branch recursively, keep newest 2](#prune-war-files-from-each-branch-recursively-keep-newest-2)
    - [Keep oldest 1, debug logging, ignore errors, use syslog](#keep-oldest-1-debug-logging-ignore-errors-use-syslog)
    - [Keep log in JSON format, use log file](#keep-log-in-json-format-use-log-file)
    - [Build and run from test area, no options](#build-and-run-from-test-area-no-options)
  - [References](#references)
  - [License](#license)

## Project home

See [our GitHub repo](https://github.com/atc0005/elbow) for the latest code,
to file an issue or submit improvements for review and potential inclusion
into the project.

## Purpose

Prune content matching specific patterns, either in a single directory or
recursively through a directory tree. The primary goal is to use this
application from a cron job to perform routine pruning of generated files that
would otherwise completely clog a filesystem.

## Gotchas

- File extensions are *case-sensitive*
- File name patterns are *case-sensitive*
- File name patterns, much like shell globs, can match more than you might
  wish.
  - Test carefully and do not provide the `--remove` flag until you have
    tested and are ready to actually prune the content.

## Features

- Extensive command-line flags with detailed help output
- (Optional) Use environment variables instead of or in addition to
  command-line arguments
  - Note: See the [Precedence](#precedence) list for how multiple
    configuration sources are processed
- Match on specified file patterns
- Flat (single-level) or recursive search
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
but not yet an official release are also noted in the file under the
`Unreleased` section. A helpful link to the Git commit history since the last
official release is also provided.

## Requirements

- Go 1.13+ (for building)
- Linux (if using Syslog support)
  - macOS and UNIX systems have not been tested
- GCC
  - if building with custom options (as the provided `Makefile` does)
- `make`
  - if using the provided `Makefile`
- UPX
  - if using the provided `Makefile`

Tested using:

- Go 1.13+
- Windows 10 Version 1803+
- Ubuntu Linux 16.04+

## How to install it

1. [Download](https://golang.org/dl/) Go
1. [Install](https://golang.org/doc/install) Go
1. Clone the repo
   1. `cd /tmp`
   1. `git clone https://github.com/atc0005/elbow`
   1. `cd elbow`
1. Install dependencies (optional)
   - for Ubuntu Linux
     - `sudo apt-get install make gcc upx`
   - for CentOS Linux
     1. `sudo yum install make gcc epel-release`
     1. `sudo yum install upx`
1. Build
   - for current operating system with default `go` build options
     - `go build`
   - for all supported platforms
      - `make all`
   - for Windows
      - `make windows`
   - for Linux
     - `make linux`
1. Copy the applicable binary to whatever systems needs to run it
   1. Linux: `/tmp/elbow/elbow`
   1. Windows: `/tmp/elbow/elbow.exe`

## Setup test environment

1. Launch container, VM or WSL instance
1. Clone the repo
   1. `cd /tmp`
   1. `git clone https://github.com/atc0005/elbow`
   1. `cd elbow`
1. Create test files
   - in `/tmp`
     - `make testenv`
     - Note: `/tmp` is the default location
   - in a custom location (e.g., in `$HOME/tmp`)
     - `mkdir -vp $HOME/tmp`
     - `make testenv TESTENVDIR=$HOME/tmp`

See the [Examples](#examples) or the [Configuration
Options](#configuration-options) sections for examples of running `elbow`
against these newly created test files.

## Configuration Options

### Precedence

The priority order is:

1. Command line flags (highest priority)
1. Configuration file
   - **Not supported yet**
1. Environment variables (lowest priority)

### Command-line Arguments

Aside from the built-in `-h`, short flag names are currently not supported.

| Long             | Required | Default        | Repeat | Possible                                                                                                | Description                                                                                              |
| ---------------- | -------- | -------------- | ------ | ------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- |
| `keep`           | Yes      | N/A            | No     | `0+`                                                                                                    | Keep specified number of matching files.                                                                 |
| `path`           | Yes      | N/A            | No     | *valid directory path*                                                                                  | Path to process.                                                                                         |
| `pattern`        | No       | *empty string* | No     | *valid file name characters*                                                                            | Substring pattern to compare filenames against. Wildcards are not supported.                             |
| `extensions`     | No       | *empty list*   | No     | *valid file extensions*                                                                                 | Limit search to specified file extension. Specify as needed to match multiple required extensions.       |
| `recurse`        | No       | `false`        | No     | `true`, `false`                                                                                         | Perform recursive search into subdirectories.                                                            |
| `keep-old`       | No       | `false`        | No     | `true`, `false`                                                                                         | Keep oldest files instead of newer.                                                                      |
| `remove`         | Maybe    | `false`        | No     | `true`, `false`                                                                                         | Remove matched files. The default behavior is to only note what matching files *would* be removed.       |
| `ignore-errors`  | No       | `false`        | No     | `true`, `false`                                                                                         | Ignore errors encountered during file removal.                                                           |
| `log-format`     | No       | `text`         | No     | `text`, `json`                                                                                          | Log formatter used by logging package.                                                                   |
| `log-file`       | No       | *empty string* | No     | *writable directory path*                                                                               | Optional log file used to hold logged messages. If set, log messages are not displayed on the console.   |
| `console-output` | No       | `stdout`       | No     | `stdout`, `stderr`                                                                                      | Specify how log messages are logged to the console.                                                      |
| `log-level`      | No       | `info`         | No     | `emergency`, `alert`, `critical`, `panic`, `fatal`, `error`, `warn`, `info`, `notice`, `debug`, `trace` | Maximum log level at which messages will be logged. Log messages below this threshold will be discarded. |
| `use-syslog`     | No       | `false`        | No     | `true`, `false`                                                                                         | Log messages to syslog in addition to other ouputs. Not supported on Windows.                            |

### Environment Variables

If set, command-line arguments override the equivalent environment variables
listed below. See the [Command-line Arguments](#command-line-arguments) table
for more information.

| Flag Name        | Environment Variable Name | Notes                        | Example                                                     |
| ---------------- | ------------------------- | ---------------------------- | ----------------------------------------------------------- |
| `keep`           | `ELBOW_KEEP`              |                              | `ELBOW_KEEP=1`                                              |
| `path`           | `ELBOW_PATH`              |                              | `ELBOW_PATH="/tmp"`                                         |
| `pattern`        | `ELBOW_FILE_PATTERN`      |                              | `ELBOW_FILE_PATTERN="reach-masterdev-"`                     |
| `extensions`     | `ELBOW_EXTENSIONS`        | *Comma-separated, no spaces* | `ELBOW_EXTENSIONS=".war,.tmp"`                              |
| `recurse`        | `ELBOW_RECURSE`           |                              | `ELBOW_RECURSE="true"`                                      |
| `keep-old`       | `ELBOW_KEEP_OLD`          |                              | `ELBOW_KEEP_OLD="true"`                                     |
| `remove`         | `ELBOW_REMOVE`            |                              | `ELBOW_REMOVE="false"`                                      |
| `ignore-errors`  | `ELBOW_IGNORE_ERRORS`     |                              | `ELBOW_IGNORE_ERRORS="true"`                                |
| `log-format`     | `ELBOW_LOG_FORMAT`        |                              | `ELBOW_LOG_FORMAT="json"`                                   |
| `log-file`       | `ELBOW_LOG_FILE`          |                              | `ELBOW_LOG_FILE="/tmp/testing-masterqa-build-removals.txt"` |
| `console-output` | `ELBOW_CONSOLE_OUTPUT`    |                              | `ELBOW_CONSOLE_OUTPUT="stdout"`                             |
| `log-level`      | `ELBOW_LOG_LEVEL`         |                              | `ELBOW_LOG_LEVEL="debug"`                                   |
| `use-syslog`     | `ELBOW_USE_SYSLOG`        |                              | `ELBOW_USE_SYSLOG="true"`                                   |

### Configuration File

Not yet supported.

## Examples

### Overview

The following steps illustrate a rough, overall idea of what `elbow` is
intended to do. The steps illustrate building and running the application from
within an Ubuntu Linux Subsystem for Windows (WSL) instance. The `t` volume is
present on the Windows host.

The file extension used in the examples below is for `WAR` files that are
generated on a build system that our group maintains. The idea is that `elbow`
could be run as a cron job to help ensure that only X copies (the most recent
in our case) for each of three branches remain on the build box.

There are better aproaches to managing build artifacts (e.g., containers), but
that is the problem that this tool seeks to solve in a simple, "low tech" way.

The particular repo that the build system processes has three branches:

| Branch Name | Type of build |
| ----------- | ------------- |
| `master`    | Production    |
| `masterqa`  | Q/A           |
| `masterdev` | Development   |

We had little control over the name of these branches.

### Log output

#### Text format

```ShellSession
$ ./elbow --path /tmp --pattern "reach-master" --keep 1 --recurse --keep-old --ignore-errors --log-level info --use-syslog --log-format text
```

```ShellSession
ERRO[0000] Failed to enable syslog logging: unable to connect to syslog socket: Unix syslog delivery error
WARN[0000] Proceeding without syslog logging
INFO[0000] Evaluating path: /tmp
INFO[0000] Looking for file pattern: "reach-master"
INFO[0000] Looking for extensions: []
ERRO[0000] error:open /tmp/tmp0dyy3wu9: permission denied  ignore_errors=true
WARN[0000] Error encountered, but continuing as requested.
INFO[0000] 24 files eligible for removal
INFO[0000] 1 files to keep as requested                  keep_oldest=true
INFO[0000] Ignoring file removal errors: true
INFO[0000] File removal not enabled, not removing files
INFO[0000] 0 files successfully removed
INFO[0000] 0 files failed to remove
INFO[0000] Elbow successfully completed.
```

Where supported, the output is colored. Here is a screenshot of the output
from just before the v0.2.0 milestone was completed and the [`v0.2.0`
tag](https://github.com/atc0005/elbow/releases/tag/v0.2.0) created:

![alt text][screenshot]

#### JSON format

```ShellSession
$ ./elbow --path /tmp --pattern "reach-master" --keep 1 --recurse --keep-old --ignore-errors --log-level info --use-syslog --log-format json
```

```json
{"level":"error","msg":"Failed to enable syslog logging: unable to connect to syslog socket: Unix syslog delivery error","time":"2019-09-26T12:38:34-05:00"}
{"level":"warning","msg":"Proceeding without syslog logging","time":"2019-09-26T12:38:34-05:00"}
{"level":"info","msg":"Evaluating path: /tmp","time":"2019-09-26T12:38:34-05:00"}
{"level":"info","msg":"Looking for file pattern: \"reach-master\"","time":"2019-09-26T12:38:34-05:00"}
{"level":"info","msg":"Looking for extensions: []","time":"2019-09-26T12:38:34-05:00"}
{"ignore_errors":true,"level":"error","msg":"error:open /tmp/tmp0dyy3wu9: permission denied","time":"2019-09-26T12:38:34-05:00"}
{"level":"warning","msg":"Error encountered, but continuing as requested.","time":"2019-09-26T12:38:34-05:00"}
{"level":"info","msg":"24 files eligible for removal","time":"2019-09-26T12:38:34-05:00"}
{"keep_oldest":true,"level":"info","msg":"1 files to keep as requested","time":"2019-09-26T12:38:34-05:00"}
{"level":"info","msg":"Ignoring file removal errors: true","time":"2019-09-26T12:38:34-05:00"}
{"level":"info","msg":"File removal not enabled, not removing files","time":"2019-09-26T12:38:34-05:00"}
{"level":"info","msg":"0 files successfully removed","time":"2019-09-26T12:38:34-05:00"}
{"level":"info","msg":"0 files failed to remove","time":"2019-09-26T12:38:34-05:00"}
{"level":"info","msg":"Elbow successfully completed.","time":"2019-09-26T12:38:34-05:00"}
```

### Help Output

```ShellSession
$ ./elbow --help
Elbow prunes content matching specific patterns, either in a single directory or recursively through a directory tree.

ELBOW x.y.z
https://github.com/atc0005/elbow

Usage: elbow [--pattern PATTERN] [--extensions EXTENSIONS] [--path PATH] [--recurse] [--age AGE] [--keep KEEP] [--keep-old] [--remove] [--ignore-errors] [--log-format LOG-FORMAT] [--log-file LOG-FILE] [--console-output CONSOLE-OUTPUT] [--log-level LOG-LEVEL] [--use-syslog]

Options:
  --pattern PATTERN      Substring pattern to compare filenames against. Wildcards are not supported.
  --extensions EXTENSIONS
                         Limit search to specified file extensions. Specify as space separated list to match multiple required extensions.
  --path PATH            Path to process.
  --recurse              Perform recursive search into subdirectories.
  --age AGE              Limit search to files that are the specified number of days old or older.
  --keep KEEP            Keep specified number of matching files.
  --keep-old             Keep oldest files instead of newer.
  --remove               Remove matched files.
  --ignore-errors        Ignore errors encountered during file removal.
  --log-format LOG-FORMAT
                         Log formatter used by logging package. [default: text]
  --log-file LOG-FILE    Optional log file used to hold logged messages. If set, log messages are not displayed on the console.
  --console-output CONSOLE-OUTPUT
                         Specify how log messages are logged to the console. [default: stdout]
  --log-level LOG-LEVEL
                         Maximum log level at which messages will be logged. Log messages below this threshold will be discarded. [default: info]
  --use-syslog           Log messages to syslog in addition to other outputs. Not supported on Windows.
  --help, -h             display this help and exit
  --version              display version and exit
```

### Prune `.war` files from each branch recursively, keep newest 2

Note: Leave off `--remove` to display what *would* be removed.

```ShellSession
cd /mnt/t/github/elbow; go build; cp -vf elbow /tmp/; cd /tmp/
./elbow --path /tmp --extensions ".war" --pattern "reach-master-" --keep 2 --recurse --remove
./elbow --path /tmp --extensions ".war" --pattern "reach-masterqa-" --keep 2 --recurse --remove
./elbow --path /tmp --extensions ".war" --pattern "reach-masterdev-" --keep 2 --recurse --remove
```

### Keep oldest 1, debug logging, ignore errors, use syslog

Note: Leave off `--remove` to display what *would* be removed.

```ShellSession
cd /mnt/t/github/elbow; go build; cp -vf elbow /tmp/; cd /tmp/
./elbow --path /tmp --pattern "reach-master-" --keep 1 --recurse --keep-old --ignore-errors --log-level debug --use-syslog
./elbow --path /tmp --pattern "reach-masterqa-" --keep 1 --recurse --keep-old --ignore-errors --log-level debug --use-syslog
./elbow --path /tmp --pattern "reach-masterdev-" --keep 1 --recurse --keep-old --ignore-errors --log-level debug --use-syslog
```

### Keep log in JSON format, use log file

- These examples attempt to create a log file in the current directory, which is `/tmp` in this case.
- The default logging format is `text` unless overridden; here we specify `json`.
- We attempt to enable syslog logging. This currently fails gracefully on Windows.
- We ignore file removal errors and proceed to the next matching file.

```ShellSession
cd /mnt/t/github/elbow; go build; cp -vf elbow /tmp/; cd /tmp/
./elbow --path /tmp --pattern "reach-master-" --keep 1 --recurse --keep-old --ignore-errors --log-level debug --log-format json --use-syslog --log-file testing-master-build-removals.txt
./elbow --path /tmp --pattern "reach-masterqa-" --keep 1 --recurse --keep-old --ignore-errors --log-level debug --use-syslog --log-format json --log-file testing-masterqa-build-removals.txt
./elbow --path /tmp --pattern "reach-masterdev-" --keep 1 --recurse --keep-old --ignore-errors --log-level debug --use-syslog --log-format json --log-file testing-masterdev-build-removals.txt
```

### Build and run from test area, no options

This results in Help text being displayed. At a minimum, the path to process
has to be provided for the application to proceed.

```ShellSession
cd /mnt/t/github/elbow; go build; cp -vf elbow /tmp/; cd /tmp/; ./elbow
```

## References

See the [docs/references.md](docs/references.md) for details.

## License

Taken directly from the `LICENSE` and `NOTICES.txt` files:

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

[screenshot]: media/elbow_example_text_log_format_2019-09-26.png "Colored text output example screenshot"
