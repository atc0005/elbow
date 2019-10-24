/*

Elbow prunes content matching specific patterns, either in a single
directory or recursively through a directory tree.

PROJECT HOME

See our GitHub repo (https://github.com/atc0005/elbow) for the latest code, to
file an issue or submit improvements for review and potential inclusion into
the project.

PURPOSE

Prune content matching specific patterns, either in a single directory or
recursively through a directory tree. The primary goal is to use this
application from a cron job to perform routine pruning of generated files that
would otherwise completely clog a filesystem.

GOTCHAS

• File extensions are case-sensitive.

• File name patterns are case-sensitive.

• File name patterns, much like shell globs, may match more than intended. Test carefully and do not enable file removal until you have tested and are ready to actually prune the content.

FEATURES

• Extensive command-line flags with detailed help output

• (Optional) Use environment variables instead of or in addition to command-line arguments

• Match on specified file patterns

• Flat (single-level) or recursive search

• Process one or many paths

• Age-based threshold for matches (e.g., match files X days old or older)

• Keep a specified number of older or newer matches

• Limit search to specified list of file extensions

• Toggle file removal (read-only by default)

• Extensive, leveled-logging

• (Optional) Syslog logging (not supported on Windows)

• (Optional) Logging to a file (if enabled, mutes console output)

• Logging in Text or JSON log formats

• (Optional) Ignore errors encountered when removing files

USAGE

Help output is below. See the README for examples.

   $ ./elbow --help
   Elbow prunes content matching specific patterns, either in a single directory or recursively through a directory tree.

   ELBOW x.y.z
   https://github.com/atc0005/elbow

   Usage: elbow [--pattern PATTERN] [--extensions EXTENSIONS] [--age AGE] [--keep KEEP] [--keep-old] [--remove] [--ignore-errors] [--log-level LOG-LEVEL] [--log-format LOG-FORMAT] [--log-file LOG-FILE] [--console-output CONSOLE-OUTPUT] [--use-syslog] [--paths PATHS] [--recurse]

   Options:
   --pattern PATTERN      Substring pattern to compare filenames against. Wildcards are not supported.
   --extensions EXTENSIONS
                           Limit search to specified file extensions. Specify as space separated list to match multiple required extensions.
   --age AGE              Limit search to files that are the specified number of days old or older.
   --keep KEEP            Keep specified number of matching files per provided path.
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
   --help, -h             display this help and exit
   --version              display version and exit

*/
package main
