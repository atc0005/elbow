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


###########################################################
# Run linters
###########################################################


# https://stackoverflow.com/a/42510278/903870
diff -u <(echo -n) <(gofmt -l -e -d .)


go vet ./...


if ! which golint > /dev/null; then
cat <<\EOF
Error: Unable to locate "golint"

Change your current working directory (e.g., "cd $HOME") and install golint
with the following command:

go get -u golang.org/x/lint/golint
EOF
    exit 1
else
    golint -set_exit_status ./...
fi

if ! which golangci-lint > /dev/null; then
cat <<\EOF
Error: Unable to locate "golangci-lint"

Change your current working directory (e.g., "cd $HOME") and install golangci-lint
with the following command:

go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
EOF
    exit 1
else
    golangci-lint run \
        -E goimports \
        -E gosec \
        -E stylecheck \
        -E goconst \
        -E depguard \
        -E prealloc
fi


if ! which staticcheck > /dev/null; then
cat <<\EOF
Error: Unable to locate "staticcheck"

Change your current working directory (e.g., "cd $HOME") and install staticcheck
with the following command:

go get -u honnef.co/go/tools/cmd/staticcheck
EOF
    exit 1
else
    staticcheck ./...
fi
