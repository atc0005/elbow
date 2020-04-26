#!/bin/bash

# Purpose: Small wrapper script to build and call binary with environment
# variables already configured.


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


BINARY_NAME="$1"
PATH1="$2"
PATH2="$3"

# Use default build options
go build

# See README for complete list of environment variables
export ELBOW_PATHS="${PATH1}"
export ELBOW_FILE_PATTERN="reach-masterdev-"
export ELBOW_EXTENSIONS=".war,.tmp"
export ELBOW_KEEP=1
export ELBOW_FILE_AGE=1
export ELBOW_RECURSE="true"
export ELBOW_KEEP_OLD="true"
export ELBOW_IGNORE_ERRORS="true"
export ELBOW_LOG_LEVEL="lizard"
export ELBOW_USE_SYSLOG="false"
export ELBOW_LOG_FORMAT="json"
export ELBOW_REMOVE="false"
#export ELBOW_LOG_FILE="testing-masterqa-build-removals.txt"
#export ELBOW_CONFIG_FILE="config.example.toml"
#export ELBOW_CONFIG_FILE="config.toml"

echo -e "\n\nCalling ${BINARY_NAME} without flags; rely on env vars\n"
./${BINARY_NAME}

if [[ $? -ne 0 ]]; then
  echo "${BINARY_NAME} execution failed. See earlier output for details."
  sleep 3
fi

echo -e "\n\nCalling ${BINARY_NAME} with 1 path, 2 extensions specified\n"
./${BINARY_NAME} \
  --paths "${PATH1}" \
  --extensions ".war" ".tmp" \
  --pattern "" \
  --keep 1 \
  --recurse \
  --keep-old \
  --ignore-errors \
  --log-level info \
  --use-syslog \
  --log-format text \
  --console-output "stdout"

if [[ $? -ne 0 ]]; then
  echo "${BINARY_NAME} execution failed. See earlier output for details."
  sleep 3
fi


echo -e "\n\nCalling ${BINARY_NAME} with 2 paths and 2 extensions specified\n"
./${BINARY_NAME} \
  --paths "${PATH1}" "${PATH2}" \
  --extensions ".war" ".tmp" \
  --pattern "" \
  --keep 1 \
  --recurse \
  --keep-old \
  --ignore-errors \
  --log-level info \
  --use-syslog \
  --log-format text \
  --console-output "stdout"

if [[ $? -ne 0 ]]; then
  echo "${BINARY_NAME} execution failed. See earlier output for details."
  sleep 3
fi

echo -e "\n\nCalling ${BINARY_NAME} with 2 paths (one good, one bad) and 2 extensions specified\n"
./${BINARY_NAME} \
  --paths "/tmp3" "${PATH1}" \
  --extensions ".war" ".tmp" \
  --pattern "" \
  --keep 1 \
  --recurse \
  --keep-old \
  --ignore-errors \
  --log-level info \
  --use-syslog \
  --log-format text \
  --console-output "stdout"

if [[ $? -ne 0 ]]; then
  echo "${BINARY_NAME} execution failed. See earlier output for details."
  sleep 3
fi

# Confirm that precedence works as expected
echo -e "\n\nCalling ${BINARY_NAME} with flags; override env vars\n"
./${BINARY_NAME} \
  --paths "${PATH1}" \
  --keep 1 \
  --recurse \
  --keep-old \
  --ignore-errors \
  --log-level info \
  --use-syslog \
  --log-format text \
  --console-output "stdout" \
  --remove

if [[ $? -ne 0 ]]; then
  echo "${BINARY_NAME} execution failed. See earlier output for details."
  sleep 3
fi

# Provide invalid option
echo -e "\n\nCalling ${BINARY_NAME} with invalid flag"
echo -e "This will result in several Makefile errors:\n"
echo -e "    Makefile:66: recipe for target 'testrun' failed\n    make: *** [testrun] Error 1\n"
#read -p "Press enter to continue"

./${BINARY_NAME} \
  --paths "${PATH1}" \
  --keep 1 \
  --recurse \
  --keep-old \
  --ignore-errors \
  --log-level info \
  --use-syslog \
  --log-format text \
  --console-output "tacos"

if [[ $? -ne 0 ]]; then
  echo "${BINARY_NAME} execution failed. See earlier output for details."
  sleep 3
fi


# Attempt to use config.toml config file from the root of the repo
echo -e "\n\nCalling ${BINARY_NAME} with config file flag and paths command-line flag"

echo "Unset all ELBOW_* environment variables to prevent shadowing tests with config file"
unset ELBOW_PATHS
unset ELBOW_FILE_PATTERN
unset ELBOW_EXTENSIONS
unset ELBOW_KEEP
unset ELBOW_FILE_AGE
unset ELBOW_RECURSE
unset ELBOW_KEEP_OLD
unset ELBOW_IGNORE_ERRORS
unset ELBOW_LOG_LEVEL
unset ELBOW_USE_SYSLOG
unset ELBOW_LOG_FORMAT
unset ELBOW_REMOVE
unset ELBOW_LOG_FILE
unset ELBOW_CONFIG_FILE

make testenv

if [[ ! -f "config.toml" ]]; then
  echo "config.toml not found."
  echo "Tip: Use config.example.toml as a template."
  echo "e.g., `cp config.example.toml config.toml`"
  exit 1
fi

./${BINARY_NAME} \
  --paths "${PATH1}" \
  --config-file "config.toml"

if [[ $? -ne 0 ]]; then
  echo "${BINARY_NAME} execution failed. See earlier output for details."
  sleep 3
fi

# Attempt to ONLY use config.toml config file from the root of the repo
echo -e "\n\nCalling ${BINARY_NAME} with config file flag ONLY"

make testenv

./${BINARY_NAME} --config-file "config.toml"
