module github.com/atc0005/elbow

// Copyright 2020 Adam Chalkley
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

// Based on some reading in https://github.com/golang/go/wiki/Modules, the
// behavior for `go get -u` changed to allow more conservative updates of
// dependencies. The new behavior sounds more natural and is less likely to
// surprise newcomers, so locking the base behavior to Go 1.13 sounds like a
// "Good Thing" to do here.
go 1.13

require (
	github.com/alexflint/go-arg v1.3.0
	github.com/pelletier/go-toml v1.8.1
	github.com/sirupsen/logrus v1.8.0
	github.com/stretchr/testify v1.5.1 // indirect
)
