/*
Elbow prunes content matching specific patterns, either in a single
directory or recursively through a directory tree.

# Project Home

See our GitHub repo (https://github.com/atc0005/elbow) for the latest code, to
file an issue or submit improvements for review and potential inclusion into
the project.

# Purpose

Prune content matching specific patterns, either in a single directory or
recursively through a directory tree. The primary goal is to use this
application from a cron job to perform routine pruning of generated files that
would otherwise completely clog a filesystem.

# Gotchas

  - File extensions are case-sensitive.
  - File name patterns are case-sensitive.
  - File name patterns, much like shell globs, may match more than intended. Test carefully and do not enable file removal until you have tested and are ready to actually prune the content.

# Features

  - Supports multiple (merged) sources for supplying configuration settings
    Default settings
    TOML format configuration file
    Environment variables
    Command-line flags (with detailed help output)
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
  - Logging in Text or JSON log formats
  - (Optional) Ignore errors encountered when removing files

# Usage

See our main README for supported settings and examples.
*/
package main
