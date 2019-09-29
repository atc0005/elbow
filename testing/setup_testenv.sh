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

TEST_DIR_PATH=$1

PATH_TO_THIS_SCRIPT="$(dirname $(echo $0))"

if [[ ! -d "$TEST_DIR_PATH" ]]; then
    echo "\"$TEST_DIR_PATH\" not found! Please create and then re-run make command."
    exit 1
fi

while read line
do
    touch -d $(echo $line | awk -F\- '{print $4}') ${TEST_DIR_PATH}/${line}
done < ${PATH_TO_THIS_SCRIPT}/sample_files_list_dev_web_app_server.txt
