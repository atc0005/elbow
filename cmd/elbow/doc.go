// Copyright 2019 Adam Chalkley
//
// https://github.com/atc0005/elbow
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Prune content matching specific patterns, either in a single directory or
// recursively through a directory tree. The primary goal is to use this
// application from a cron job to perform routine pruning of generated files
// that would otherwise completely clog a filesystem.
//
// See our [GitHub repo]:
//
//   - to review documentation (including examples)
//   - for the latest code
//   - to file an issue or submit improvements for review and potential
//     inclusion into the project
//
// [GitHub repo]: https://github.com/atc0005/elbow
package main
