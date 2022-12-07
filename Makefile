
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

# References:
#
# https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies
# https://github.com/mapnik/sphinx-docs/blob/master/Makefile
# https://stackoverflow.com/questions/23843106/how-to-set-child-process-environment-variable-in-makefile
# https://stackoverflow.com/questions/3267145/makefile-execute-another-target
# https://unix.stackexchange.com/questions/124386/using-a-make-rule-to-call-another
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
# https://www.gnu.org/software/make/manual/html_node/Recipe-Syntax.html#Recipe-Syntax
# https://www.gnu.org/software/make/manual/html_node/Special-Variables.html#Special-Variables
# https://github.com/golangci/golangci-lint#install
# https://github.com/golangci/golangci-lint/releases/latest


SHELL = /bin/bash

APPNAME					= elbow

# What package holds the "version" variable used in branding/version output?
# VERSION_VAR_PKG			= $(shell go list -m)
VERSION_VAR_PKG			= main

OUTPUTDIR 				= release_assets

TESTENVBASEDIR			= /tmp/$(APPNAME)
TESTENVDIR1				= $(TESTENVBASEDIR)/path1
TESTENVDIR2				= $(TESTENVBASEDIR)/path2

# https://gist.github.com/TheHippo/7e4d9ec4b7ed4c0d7a39839e6800cc16
# https://github.com/atc0005/elbow/issues/54
VERSION 				= $(shell git describe --always --long --dirty)

# The default `go build` process embeds debugging information. Building
# without that debugging information reduces the binary size by around 28%.
#
# We also include additional flags in an effort to generate static binaries
# that do not have external dependencies. As of Go 1.15 this still appears to
# be a mixed bag, so YMMV.
#
# See https://github.com/golang/go/issues/26492 for more information.
#
# -s
#	Omit the symbol table and debug information.
#
# -w
#	Omit the DWARF symbol table.
#
# -tags 'osusergo,netgo'
#	Use pure Go implementation of user and group id/name resolution.
#	Use pure Go implementation of DNS resolver.
#
# -trimpath
#	https://golang.org/cmd/go/
#   removes all file system paths from the compiled executable, to improve
#   build reproducibility.
#
# CGO_ENABLED=0
#	https://golang.org/cmd/cgo/
#	explicitly disable use of cgo
BUILDCMD				=	CGO_ENABLED=0 go build -mod=vendor -a -trimpath -tags 'osusergo,netgo' -ldflags="-s -w -X $(VERSION_VAR_PKG).version=$(VERSION)"
GOCLEANCMD				=	go clean -mod=vendor ./...
GITCLEANCMD				= 	git clean -xfd
CHECKSUMCMD				=	sha256sum -b

# TODO: Deprecate these?
TESTENVCMD				=   bash testing/setup_testenv.sh
TESTRUNCMD				=   bash testing/run_with_test_settings.sh

.DEFAULT_GOAL := help

  ##########################################################################
  # Targets will not work properly if a file with the same name is ever
  # created in this directory. We explicitly declare our targets to be phony
  # by making them a prerequisite of the special target .PHONY
  ##########################################################################


.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: testenv
## testenv: setup test environment in Windows Subsystem for Linux or other Linux system
testenv:
	@echo "Setting up test environment in \"$(TESTENVDIR1)\" and \"$(TESTENVDIR2)\""
	@$(TESTENVCMD) "$(TESTENVDIR1)" "$(TESTENVDIR2)"
	@echo "Finished creating test files in \"$(TESTENVDIR1)\" and \"$(TESTENVDIR2)\""

.PHONY: testrun
## testrun: use wrapper script to call binary with test settings
testrun:
	@echo "Calling wrapper script: $(TESTRUNCMD)"
	@$(TESTRUNCMD) "$(OUTPUTBASEFILENAME)" "$(TESTENVDIR1)" "$(TESTENVDIR2)"
	@echo "Finished running wrapper script"

.PHONY: lintinstall
## lintinstall: install common linting tools
# https://github.com/golang/go/issues/30515#issuecomment-582044819
lintinstall:
	@echo "Installing linting tools"

	@export PATH="${PATH}:$(go env GOPATH)/bin"

	@echo "Installing latest stable staticcheck version via go install command ..."
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck --version

	@echo Installing latest stable golangci-lint version per official installation script ...
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin
	golangci-lint --version

	@echo "Finished updating linting tools"

.PHONY: linting
## linting: runs common linting checks
linting:
	@echo "Running linting tools ..."

	@echo "Running go vet ..."
	@go vet -mod=vendor $(shell go list -mod=vendor ./... | grep -v /vendor/)

	@echo "Running golangci-lint ..."
	@golangci-lint --version
	@golangci-lint run

	@echo "Running staticcheck ..."
	@staticcheck --version
	@staticcheck $(shell go list -mod=vendor ./... | grep -v /vendor/)

	@echo "Finished running linting checks"

.PHONY: gotests
## gotests: runs go test recursively, verbosely
gotests:
	@echo "Running go tests ..."
	@go test -mod=vendor ./...
	@echo "Finished running go tests"

.PHONY: goclean
## goclean: removes local build artifacts, temporary files, etc
goclean:
	@echo "Removing object files and cached files ..."
	@$(GOCLEANCMD)
	@echo "Removing any existing release assets"
	@mkdir -p "$(OUTPUTDIR)"
	@rm -vf $(wildcard ${OUTPUTDIR}/*/*-linux-*)
	@rm -vf $(wildcard ${OUTPUTDIR}/*/*-windows-*)

.PHONY: clean
## clean: alias for goclean
clean: goclean

.PHONY: gitclean
## gitclean: WARNING - recursively cleans working tree by removing non-versioned files
gitclean:
	@echo "Removing non-versioned files ..."
	@$(GITCLEANCMD)

.PHONY: pristine
## pristine: run goclean and gitclean to remove local changes
pristine: goclean gitclean

.PHONY: all
# https://stackoverflow.com/questions/3267145/makefile-execute-another-target
## all: generates assets for Linux distros and Windows
all: clean windows linux
	@echo "Completed all cross-platform builds ..."

.PHONY: windows
## windows: generates assets for Windows systems
windows:
	@echo "Building release assets for windows ..."

	@mkdir -p $(OUTPUTDIR)/$(APPNAME)

	@echo "Building 386 binaries"
	@env GOOS=windows GOARCH=386 $(BUILDCMD) -o $(OUTPUTDIR)/$(APPNAME)/$(APPNAME)-$(VERSION)-windows-386.exe ${PWD}/cmd/$(APPNAME)

	@echo "Building amd64 binaries"
	@env GOOS=windows GOARCH=amd64 $(BUILDCMD) -o $(OUTPUTDIR)/$(APPNAME)/$(APPNAME)-$(VERSION)-windows-amd64.exe ${PWD}/cmd/$(APPNAME)

	@echo "Generating checksum files"
	@cd $(OUTPUTDIR)/$(APPNAME) && \
	$(CHECKSUMCMD) $(APPNAME)-$(VERSION)-windows-386.exe > $(APPNAME)-$(VERSION)-windows-386.exe.sha256 && \
	$(CHECKSUMCMD) $(APPNAME)-$(VERSION)-windows-amd64.exe > $(APPNAME)-$(VERSION)-windows-amd64.exe.sha256

	@echo "Completed build tasks for windows"

.PHONY: linux
## linux: generates assets for Linux distros
linux:
	@echo "Building release assets for linux ..."

	@mkdir -p $(OUTPUTDIR)/$(APPNAME)

	@echo "Building 386 binaries"
	@env GOOS=linux GOARCH=386 $(BUILDCMD) -o $(OUTPUTDIR)/$(APPNAME)/$(APPNAME)-$(VERSION)-linux-386 ${PWD}/cmd/$(APPNAME)

	@echo "Building amd64 binaries"
	@env GOOS=linux GOARCH=amd64 $(BUILDCMD) -o $(OUTPUTDIR)/$(APPNAME)/$(APPNAME)-$(VERSION)-linux-amd64 ${PWD}/cmd/$(APPNAME)

	@echo "Generating checksum files"
	@cd $(OUTPUTDIR)/$(APPNAME) && \
	$(CHECKSUMCMD) "$(APPNAME)-$(VERSION)-linux-386" > "$(APPNAME)-$(VERSION)-linux-386.sha256" && \
	$(CHECKSUMCMD) "$(APPNAME)-$(VERSION)-linux-amd64" > "$(APPNAME)-$(VERSION)-linux-amd64.sha256"

	@echo "Completed build tasks for linux"
