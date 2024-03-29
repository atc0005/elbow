#!/bin/bash

# Copyright 2023 Adam Chalkley
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

project_org="atc0005"
project_shortname="elbow"

project_fq_name="${project_org}/${project_shortname}"
project_url_base="https://github.com/${project_org}"
project_repo="${project_url_base}/${project_shortname}"
project_releases="${project_repo}/releases"
project_issues="${project_repo}/issues"
project_discussions="${project_repo}/discussions"

plugin_name=""
plugin_path="/usr/lib64/nagios/plugins"


echo
echo "Thank you for installing packages provided by the ${project_fq_name} project!"
echo
echo "#######################################################################"
echo "NOTE:"
echo
echo "This is a dev build; binaries installed by this package have a _dev"
echo "suffix to allow installation alongside stable versions."
echo
echo "Feedback for all releases is welcome, but especially so for dev builds."
echo "Thank you in advance!"
echo "#######################################################################"
echo
echo "Project resources:"
echo
echo "- Obtain latest release: ${project_releases}"
echo "- View/Ask questions: ${project_discussions}"
echo "- View/Open issues: ${project_issues}"
echo
