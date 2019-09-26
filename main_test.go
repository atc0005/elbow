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

package main

import "testing"

func TestMain(t *testing.T) {

	defaultConfig := NewConfig()

	var emptySlice = []string{}
	var nilSlice []string

	t.Logf("%v\n", emptySlice)
	t.Log(len(emptySlice))
	t.Log("emptySlice is nil:", emptySlice == nil)
	t.Log("-------------------------")

	t.Logf("%v\n", nilSlice)
	t.Log(len(nilSlice))
	t.Log("nilSlice is nil:", nilSlice == nil)
	t.Log("-------------------------")

	t.Logf("%v\n", defaultConfig.FileExtensions)
	t.Log(len(defaultConfig.FileExtensions))
	t.Log("defaultConfig.FileExtensions is nil:", defaultConfig.FileExtensions == nil)

}
