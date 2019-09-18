# elbow

Elbow, Elbow grease.

- [elbow](#elbow)
  - [Purpose](#purpose)
  - [Gotchas](#gotchas)
  - [Setup test environment](#setup-test-environment)
  - [Examples](#examples)
    - [Overview](#overview)
    - [Prune `.war` files from each branch recursively, keep newest 2](#prune-war-files-from-each-branch-recursively-keep-newest-2)
    - [... keep oldest 1, debug logging, ignore errors, use syslog](#keep-oldest-1-debug-logging-ignore-errors-use-syslog)
    - [... log in JSON format, use log file](#log-in-json-format-use-log-file)
    - [Build and run from test area, no options](#build-and-run-from-test-area-no-options)
  - [References](#references)
    - [Configuration object](#configuration-object)
    - [Sorting files](#sorting-files)
    - [Path/File Existence](#pathfile-existence)
    - [Slice management](#slice-management)

## Purpose

Prune content matching specific patterns, either in a single directory or
recursively through a directory tree. The primary goal is to use this
application from a cron job to perform routine pruning of generated files that
would otherwise completely clog a filesystem.

## Gotchas

- File extensions are *case-sensitive*
- File name patterns are *case-sensitive*
- File name patterns, much like shell globs, can match more than you might
  wish. Test carefully and do not provide the `--remove` flag until you are
  ready to actually prune the content.

## Setup test environment

1. Launch container, VM or WSL instance
1. `cd /path/to/create/test/files`
1. `touch $(cat /path/to/this/repo/testing/sample_files_list_dev_web_app_server.txt)`
1. `cd /path/to/this/repo`
1. `go build`

See next section for examples of running the app against the test files.

## Examples

### Overview

The following steps illustrate a rough, overall idea of what this application
is intended to do. The steps illustrate building and running the application
from within an Ubuntu Linux Subsystem for Windows (WSL) instance. The `/t`
volume is present on the Windows host.

The file extension used in the examples is for a `WAR` file that is generated
on a build system that our team maintains. The idea is that this application
could be run as a cron job to help ensure that only X copies (the most recent)
for each of three branches remain on the build box.

There are better aproaches to managing those build artifacts, but that is the
problem that this tool seeks to solve in a simple way.

### Prune `.war` files from each branch recursively, keep newest 2

Note: Leave off `--remove` to display what *would* be removed.

```ShellSession
cd /mnt/t/github/elbow; go build; cp -vf elbow /tmp/; cd /tmp/; ./elbow --path /tmp --extension ".war" --pattern "reach-master-" --keep 2 --recurse --remove
cd /mnt/t/github/elbow; go build; cp -vf elbow /tmp/; cd /tmp/; ./elbow --path /tmp --extension ".war" --pattern "reach-masterqa-" --keep 2 --recurse --remove
cd /mnt/t/github/elbow; go build; cp -vf elbow /tmp/; cd /tmp/; ./elbow --path /tmp --extension ".war" --pattern "reach-masterdev-" --keep 2 --recurse --remove
```

### ... keep oldest 1, debug logging, ignore errors, use syslog

Note: Leave off `--remove` to display what *would* be removed.

```ShellSession
cd /mnt/t/github/elbow; go build; cp -vf elbow /tmp/; cd /tmp/; ./elbow --path /tmp --pattern "reach-master-" --keep 1 --recurse --keep-old --ignore-errors --log-level debug --use-syslog
cd /mnt/t/github/elbow; go build; cp -vf elbow /tmp/; cd /tmp/; ./elbow --path /tmp --pattern "reach-masterqa-" --keep 1 --recurse --keep-old --ignore-errors --log-level debug --use-syslog
cd /mnt/t/github/elbow; go build; cp -vf elbow /tmp/; cd /tmp/; ./elbow --path /tmp --pattern "reach-masterdev-" --keep 1 --recurse --keep-old --ignore-errors --log-level debug --use-syslog
```

### ... log in JSON format, use log file

- These examples attempt to create a log file in the current directory, which is `/tmp` in this case.
- The default logging format is `text` unless overridden; here we specify `json`.
- We attempt to enable syslog logging. This currently fails gracefully on Windows.
- We ignore file removal errors and proceed to the next matching file.

```ShellSession
cd /mnt/t/github/elbow; go build; cp -vf elbow /tmp/; cd /tmp/; ./elbow --path /tmp --pattern "reach-master-" --keep 1 --recurse --keep-old --ignore-errors --log-level debug --log-format json --use-syslog --log-file testing-master-build-removals.txt
cd /mnt/t/github/elbow; go build; cp -vf elbow /tmp/; cd /tmp/; ./elbow --path /tmp --pattern "reach-masterqa-" --keep 1 --recurse --keep-old --ignore-errors --log-level debug --use-syslog --log-format json --log-file testing-masterqa-build-removals.txt
cd /mnt/t/github/elbow; go build; cp -vf elbow /tmp/; cd /tmp/; ./elbow --path /tmp --pattern "reach-masterdev-" --keep 1 --recurse --keep-old --ignore-errors --log-level debug --use-syslog --log-format json --log-file testing-masterdev-build-removals.txt
```

### Build and run from test area, no options

This results in Help text being displayed. At a minimum, the path to process
has to be provided for the application to proceed.

```ShellSession
cd /mnt/t/github/elbow; go build; cp -vf elbow /tmp/; cd /tmp/; ./elbow
```

## References

The following unordered list of sites/examples provided guidance while
developing this application. Depending on when consulted, the original code
written based on that guidance may no longer be present in the active version
of this application.

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
