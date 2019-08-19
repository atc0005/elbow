package main

import (
	"os"
	"sort"
)

// FileMatch represents a superset of statistics (including os.FileInfo) for a
// file matched by provided search criteria.
type FileMatch struct {
	os.FileInfo
	Path string
}

// FileMatches is a slice of FileMatch objects
type FileMatches []FileMatch

func (fm FileMatches) sortByModTimeAsc() {

	// https://stackoverflow.com/questions/46746862/list-files-in-a-directory-sorted-by-creation-time
	sort.Slice(fm, func(i, j int) bool {
		return fm[i].ModTime().Before(fm[j].ModTime())
	})

}

func (fm FileMatches) sortByModTimeDesc() {

	// https://stackoverflow.com/questions/46746862/list-files-in-a-directory-sorted-by-creation-time
	sort.Slice(fm, func(i, j int) bool {
		return fm[i].ModTime().After(fm[j].ModTime())
	})

}
