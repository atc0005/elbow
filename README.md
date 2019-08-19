# elbow

Elbow, Elbow grease.

## Purpose

Prune content matching specific patterns, either in a single directory or
recursively through a directory tree. The primary goal is to use this
application from a cron job to perform routine pruning of generated files that
would otherwise completely clog a filesystem.

## Setup test environment

1. cd /path/to/create/test/files
1. `touch $(cat /path/to/this/repo/sample_files_list_dev_web_app_server.txt)`
1. Build app
1. Pass in path to `/path/to/create/test/files`
