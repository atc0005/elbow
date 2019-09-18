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
