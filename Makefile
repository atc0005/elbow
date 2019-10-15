
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


OUTPUTBASEFILENAME		= elbow
TESTENVDIR				= /tmp

# The default `go build` process embeds debugging information. Building
# without that debugging information reduces the binary size by around 28%.
BUILDCMD				=	go build -a -ldflags="-s -w"
BINARYPACKCMD			=	upx -f --brute
GOCLEANCMD				=	go clean
GITCLEANCMD				= 	git clean -xfd
TESTENVCMD				=   bash testing/setup_testenv.sh
TESTRUNCMD				=   bash testing/run_with_test_settings.sh

.DEFAULT_GOAL := help

# Targets will not work properly if a file with the same name is ever created
# in this directory. We explicitly declare our targets to be phony by
# making them a prerequisite of the special target .PHONY
.PHONY: help clean goclean gitclean pristine all windows linux build testenv

# WARNING: Make expects you to use tabs to introduce recipe lines
help:
	@echo "Please use \`make <target>' where <target> is one of"
	@echo "  clean          go clean to remove local build artifacts, temporary files, etc"
	@echo "  pristine       go clean and git clean local changes"
	@echo "  all            cross-compile for multiple operating systems"
	@echo "  windows        to generate a binary file for Windows"
	@echo "  linux          to generate a binary file for Linux distros"
	@echo "  testenv        setup test environment in Windows Subsystem for Linux or other Linux system"
	@echo "  testrun        use wrapper script to call binary with test settings"

testenv:
	@echo "Setting up test environment in \"$(TESTENVDIR)\""
	@$(TESTENVCMD) "$(TESTENVDIR)"
	@echo "Finished creating test files in \"$(TESTENVDIR)\""

testrun:
	@echo "Calling wrapper script: $(TESTRUNCMD)"
	@$(TESTRUNCMD) "$(OUTPUTBASEFILENAME)" "$(TESTENVDIR)"
	@echo "Finished running wrapper script"

goclean:
	@echo "Removing object files and cached files ..."
	@$(GOCLEANCMD)

# Setup alias for user reference
clean: goclean

gitclean:
	@echo "Recursively cleaning working tree by removing non-versioned files ..."
	@$(GITCLEANCMD)

pristine: goclean gitclean


# https://stackoverflow.com/questions/3267145/makefile-execute-another-target
all: clean windows linux
	@echo "Completed all cross-platform builds ..."

windows: OS=windows
windows: OUTPUT_FILENAME=$(OUTPUTBASEFILENAME).exe
# FIXME: Figure out how to have one `build` recipe and have it called
# for `windows` and for `linux`.
windows:
	@echo "Building $(OUTPUTBASEFILENAME) for $(OS) ..."
	@env GOOS=$(OS) $(BUILDCMD) -o $(OUTPUT_FILENAME)
	@echo "Running executable packer to compress binary ..."
	@$(BINARYPACKCMD) $(OUTPUT_FILENAME)
	@echo
	@echo "Completed build for $(OS)"

linux: OS=linux
linux: OUTPUT_FILENAME=$(OUTPUTBASEFILENAME)
# FIXME: Figure out how to have one `build` recipe and have it called
# for `windows` and for `linux`.
linux:
	@echo "Building $(OUTPUTBASEFILENAME) for $(OS) ..."
	@env GOOS=$(OS) $(BUILDCMD) -o $(OUTPUT_FILENAME)
	@echo "Running executable packer to compress binary ..."
	@$(BINARYPACKCMD) $(OUTPUT_FILENAME)
	@echo
	@echo "Completed build for $(OS)"
