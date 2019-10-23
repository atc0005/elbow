#!/bin/bash

# This small script is intended to run the exact example commands from the
# README. At some point we can programatically search/replace the example
# output in the file on script run, along with doing the same for the example
# commands.
#
# Run `make testenv` before running this script.


# Text output example
./elbow \
    --paths "/tmp/elbow/path1" "/tmp/elbow/path2" \
    --pattern "reach-master" \
    --keep 1 \
    --recurse \
    --keep-old \
    --ignore-errors \
    --log-level info \
    --use-syslog \
    --log-format text

# JSON output example
./elbow \
    --paths "/tmp/elbow/path1" "/tmp/elbow/path2" \
    --pattern "reach-master" \
    --keep 1 \
    --recurse \
    --keep-old \
    --ignore-errors \
    --log-level info \
    --use-syslog \
    --log-format json

# Help/Usage output example
./elbow --help
