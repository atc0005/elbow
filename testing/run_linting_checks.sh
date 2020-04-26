#!/bin/bash

# Copyright 2019 Adam Chalkley
#
# https://github.com/atc0005/elbow
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# Purpose: Run common linters locally to confirm code quality


# Go ahead and append $GOPATH/bin to $PATH in an effort to locate
# the go linters referenced in this script.
export PATH=${PATH}:$(go env GOPATH)/bin

# Assume all is well starting out
final_exit_code=0
failed_app=""

###########################################################
# Run linters
###########################################################


# https://stackoverflow.com/a/42510278/903870
diff -u <(echo -n) <(gofmt -l -e -d .)

status=$?
if [[ $status -ne 0 ]]; then
    final_exit_code=$status
    failed_app="gofmt"
    echo "Non-zero exit code from $failed_app: $status"
fi

go vet ./...

status=$?
if [[ $status -ne 0 ]]; then
    final_exit_code=$status
    failed_app="go vet"
    echo "Non-zero exit code from $failed_app: $status"
fi

if ! which golint > /dev/null; then
cat <<\EOF
Error: Unable to locate "golint"

Install golint with the following command:

make lintinstall

EOF
    exit 1
else
    golint -set_exit_status ./...
fi

# TODO: This might not be needed based on use of "-set_exit_status"
status=$?
if [[ $status -ne 0 ]]; then
    final_exit_code=$status
    failed_app="staticcheck"
    echo "Non-zero exit code from $failed_app: $status"
fi

if ! which golangci-lint > /dev/null; then
cat <<\EOF
Error: Unable to locate "golangci-lint"

Install golangci-lint with the following command:

make lintinstall

EOF
    exit 1
else
    golangci-lint run \
        -E goimports \
        -E gosec \
        -E stylecheck \
        -E goconst \
        -E depguard \
        -E prealloc \
        -E misspell \
        -E maligned \
        -E dupl \
        -E unconvert \
        -E golint \
        -E gocritic \
        -E scopelint
fi

status=$?
if [[ $status -ne 0 ]]; then
    final_exit_code=$status
    failed_app="golangci-lint"
    echo "Non-zero exit code from $failed_app: $status"
fi

if ! which staticcheck > /dev/null; then
cat <<\EOF
Error: Unable to locate "staticcheck"

Install staticcheck with the following command:

make lintinstall

EOF
    exit 1
else
    staticcheck ./...
fi

status=$?
if [[ $status -ne 0 ]]; then
    final_exit_code=$status
    failed_app="staticcheck"
    echo "Non-zero exit code from $failed_app: $status"
fi

# Give feedback on linting failure cause
if [[ $final_exit_code -ne 0 ]]; then
    echo "Linting failed, most recent failure: $failed_app"
fi

exit $final_exit_code
