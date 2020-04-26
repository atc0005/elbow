#!/bin/bash

# Copyright 2020 Adam Chalkley
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


# Purpose: Helper script for installing linting tools used by this project

export PATH=${PATH}:$(go env GOPATH)/bin

# https://github.com/golangci/golangci-lint#install
# https://github.com/golangci/golangci-lint/releases/latest
GOLANGCI_LINT_VERSION="v1.25.0"

# Temporarily disable module-aware mode so that we can install linting tools
# without modifying this project's go.mod and go.sum files
export GO111MODULE="off"
go get -u golang.org/x/lint/golint
go get -u honnef.co/go/tools/cmd/staticcheck

curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION}
golangci-lint --version


# Reset GO111MODULE back to the default
export GO111MODULE=""
