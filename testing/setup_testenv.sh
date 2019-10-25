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

TEST_DIR_PATH1=$1
TEST_DIR_PATH2=$2

PATH_TO_THIS_SCRIPT="$(dirname $(echo $0))"

# Attempt to create paths specified by Makefile
mkdir -vp "${TEST_DIR_PATH1}"
mkdir -vp "${TEST_DIR_PATH2}"

# Fail if paths were not successfully created
if [[ ! -d "$TEST_DIR_PATH1" ]]; then
    echo "\"$TEST_DIR_PATH1\" not found! Please create and then re-run make command."
    exit 1
elif [[ ! -d "$TEST_DIR_PATH2" ]]; then
    echo "\"$TEST_DIR_PATH2\" not found! Please create and then re-run make command."
    exit 1
fi

while read line
do
    # Get random number between 1-10
    #RAND_NUM="$(shuf -i 1-10 -n 1)"
    # https://stackoverflow.com/questions/2556190/random-number-from-a-range-in-a-bash-script
    RAND_NUM="$(python3 -S -c 'import random; print(random.randrange(1,10))')"


    # Use that random number to create a non-zero byte sized file for testing
    truncate -s ${RAND_NUM}M ${TEST_DIR_PATH1}/${line}
    touch -d $(echo $line | awk -F\- '{print $4}') ${TEST_DIR_PATH1}/${line}

    # Get random number between 1-10
    RAND_NUM="$(python -S -c 'import random; print random.randrange(1,10)')"

    # Use that random number to create a non-zero byte sized file for testing
    truncate -s ${RAND_NUM}K ${TEST_DIR_PATH2}/${line}
    touch -d $(echo $line | awk -F\- '{print $4}') ${TEST_DIR_PATH2}/${line}
done < ${PATH_TO_THIS_SCRIPT}/sample_files_list_dev_web_app_server.txt
