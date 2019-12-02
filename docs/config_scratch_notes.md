# Scratch notes for config merging

## Default config

These values represent the default set of values used to create the
`baseConfig` struct. They're intended to reflect a sane set of defaults that
in most cases would not be overridden by the user, but can be as the situation
requires. This set of defaults will likely be adjusted once sufficient "field
testing" has been performed.

| Field Name        | Default Value                                                                                                      |
| ----------------- | ------------------------------------------------------------------------------------------------------------------ |
| `AppName`         | `Elbow`                                                                                                            |
| `AppDescription`  | `prunes content matching specific patterns, either in a single directory or recursively through a directory tree.` |
| `AppURL`          | `https://github.com/atc0005/elbow`                                                                                 |
| `AppVersion`      | `dev`                                                                                                              |
| `FilePattern`     | `""`                                                                                                               |
| `FileAge`         | `0`                                                                                                                |
| `NumFilesToKeep`  | `0`                                                                                                                |
| `KeepOldest`      | `false`                                                                                                            |
| `Remove`          | `false`                                                                                                            |
| `IgnoreErrors`    | `false`                                                                                                            |
| `RecursiveSearch` | `false`                                                                                                            |
| `LogLevel`        | `info`                                                                                                             |
| `LogFormat`       | `text`                                                                                                             |
| `LogFilePath`     | `""`                                                                                                               |
| `ConsoleOutput`   | `stdout`                                                                                                           |
| `UseSyslog`       | `false`                                                                                                            |
| `ConfigFile`      | `""`                                                                                                               |
| `Paths`           | `nil`                                                                                                              |
| `FileExtensions`  | `nil`                                                                                                              |

## Fields that are not merged

- flagParser
- logger
- logFileHandle

Those fields are populated on the `baseConfig` struct, not the others; the
`fileConfig` and `argsConfig` structs are temporarily created and then merged
into the `baseConfig` struct.

## Test cases

The following sets of values are intentionally chosen to conflict with one
another. The goal is to merge them in order and after each merge confirm the
results. Any deviation would help identify a potential logic problem, either
in my thinking or with the code.

### Test: Default

Mostly as noted before, but for `Paths` and `FileExtensions`. All other fields
are default values.

| Field Name        | Default Value                                                                                                      |
| ----------------- | ------------------------------------------------------------------------------------------------------------------ |
| `AppName`         | `Elbow`                                                                                                            |
| `AppDescription`  | `prunes content matching specific patterns, either in a single directory or recursively through a directory tree.` |
| `AppURL`          | `https://github.com/atc0005/elbow`                                                                                 |
| `AppVersion`      | `dev`                                                                                                              |
| `FilePattern`     | `""`                                                                                                               |
| `FileAge`         | `0`                                                                                                                |
| `NumFilesToKeep`  | `0`                                                                                                                |
| `KeepOldest`      | `false`                                                                                                            |
| `Remove`          | `false`                                                                                                            |
| `IgnoreErrors`    | `false`                                                                                                            |
| `RecursiveSearch` | `false`                                                                                                            |
| `LogLevel`        | `info`                                                                                                             |
| `LogFormat`       | `text`                                                                                                             |
| `LogFilePath`     | `""`                                                                                                               |
| `ConsoleOutput`   | `stdout`                                                                                                           |
| `UseSyslog`       | `false`                                                                                                            |
| `ConfigFile`      | `""`                                                                                                               |
| `Paths`           | `/tmp/elbow/path1`                                                                                                 |
| `FileExtensions`  | `.yaml, .json`                                                                                                     |

### Test: File Config

**NOTE**: This is an in-memory "file" configuration, not the
`config.example.toml` config file kept at the root of this repo. That config
file has many values preset that are too similar to the default config
settings in order to properly contrast the changes between merging a set of
new changes into the default, `baseConfig` struct.

| Field Name        | Default Value                          |
| ----------------- | -------------------------------------- |
| `AppName`         | `toml_app_name`                        |
| `AppDescription`  | `toml_app_description`                 |
| `AppURL`          | `toml_app_url`                         |
| `AppVersion`      | `toml_app_version`                     |
| `FilePattern`     | `reach-masterdev-`                     |
| `FileAge`         | `1`                                    |
| `NumFilesToKeep`  | `2`                                    |
| `KeepOldest`      | `true`                                 |
| `Remove`          | `true`                                 |
| `IgnoreErrors`    | `true`                                 |
| `RecursiveSearch` | `true`                                 |
| `LogLevel`        | `debug`                                |
| `LogFormat`       | `json`                                 |
| `LogFilePath`     | `/var/log/elbow.log`                   |
| `ConsoleOutput`   | `stderr`                               |
| `UseSyslog`       | `true`                                 |
| `ConfigFile`      | `/usr/local/etc/elbow/config.toml`     |
| `Paths`           | `/tmp/elbow/path1`, `/tmp/elbow/path2` |
| `FileExtensions`  | `.tmp`, `.war`                         |

## Test: Environment variables

Our current docs note that these settings are below the config file in
precedence, but the reality is otherwise. The `alexflint/go-arg` package has
no awareness of other configurations and assumes that it will be the sole
handler of config sources.

We'll need to update our docs to note that environment variables take
precedence over configuration files.

| Field Name        | Default Value            |
| ----------------- | ------------------------ |
| `AppName`         | `ElbowEnvVar`            |
| `AppDescription`  | `something nifty here`   |
| `AppURL`          | `https://example.com`    |
| `AppVersion`      | `x.y.z`                  |
| `FilePattern`     | `reach-masterqa-`        |
| `FileAge`         | `3`                      |
| `NumFilesToKeep`  | `4`                      |
| `KeepOldest`      | `false`                  |
| `Remove`          | `false`                  |
| `IgnoreErrors`    | `false`                  |
| `RecursiveSearch` | `false`                  |
| `LogLevel`        | `warning`                |
| `LogFormat`       | `text`                   |
| `LogFilePath`     | `/var/log/elbow/env.log` |
| `ConsoleOutput`   | `stdout`                 |
| `UseSyslog`       | `false`                  |
| `ConfigFile`      | `/tmp/config.toml`       |
| `Paths`           | `/tmp/elbow/path3`       |
| `FileExtensions`  | `.docx, .pptx`           |

## Test: Command-line flags

This set of values should override all others.

**NOTE**: We will have to fake the `App*` settings since we do not offer
command-line flags to override those values. It's likely that the ability to
override those values will be removed once this branch is rebased.

| Field Name        | Default Value              |
| ----------------- | -------------------------- |
| `AppName`         | `ElbowFlagVar`             |
| `AppDescription`  | `nothing fancy`            |
| `AppURL`          | `https://example.org`      |
| `AppVersion`      | `0.1.2`                    |
| `FilePattern`     | `reach-master-`            |
| `FileAge`         | `5`                        |
| `NumFilesToKeep`  | `6`                        |
| `KeepOldest`      | `true`                     |
| `Remove`          | `true`                     |
| `IgnoreErrors`    | `true`                     |
| `RecursiveSearch` | `true`                     |
| `LogLevel`        | `info`                     |
| `LogFormat`       | `json`                     |
| `LogFilePath`     | `/var/log/elbow/flags.log` |
| `ConsoleOutput`   | `stderr`                   |
| `UseSyslog`       | `true`                     |
| `ConfigFile`      | `/tmp/configfile.toml`     |
| `Paths`           | `/tmp/elbow/path4`         |
| `FileExtensions`  | `.java, .class`            |
