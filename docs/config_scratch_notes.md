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
