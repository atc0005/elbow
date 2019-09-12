package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileMatch represents a superset of statistics (including os.FileInfo) for a
// file matched by provided search criteria. This allows us to record the
// original full path while also
type FileMatch struct {
	os.FileInfo
	Path string
}

// FileMatches is a slice of FileMatch objects
// TODO: Do I really need to abstract the fact that FileMatches is a slide of
// FileMatch objects? It seems that by hiding this it makes it harder to see
// that we're working with a slice?
type FileMatches []FileMatch

func hasValidExtension(filename string, config *Config) bool {

	// NOTE: We do NOT compare extensions insensitively. We can add that
	// functionality in the future if needed.
	ext := filepath.Ext(filename)

	if len(config.FileExtensions) == 0 {
		// DEBUG
		log.Println("No extension limits have been set!")
		log.Printf("Considering %s safe for removal\n", filename)
		return true
	}

	if inFileExtensionsPatterns(ext, config.FileExtensions) {
		// DEBUG
		log.Printf("%s has a valid extension for removal\n", filename)
		return true
	}

	// DEBUG
	log.Println("hasValidExtension: returning false for:", filename)

	log.Printf("hasValidExtension: returning false (%q not in %q)",
		ext, config.FileExtensions)
	return false
}

func hasValidFilenamePattern(filename string, config *Config) bool {

	if strings.TrimSpace(config.FilePattern) == "" {
		// DEBUG
		log.Println("No FilePattern has been specified!")
		log.Printf("Considering %s safe for removal\n", filename)
		return true
	}

	// Search for substring
	if strings.Contains(filename, config.FilePattern) {
		return true
	}

	// DEBUG
	log.Printf("hasValidFilenamePattern: returning false (%q does not contain %q)",
		filename, config.FilePattern)
	return false
}

// inFileExtensionsPatterns is a helper function to emulate Python's `if "x"
// in list:` functionality
func inFileExtensionsPatterns(ext string, exts []string) bool {
	for _, pattern := range exts {
		if ext == pattern {
			return true
		}
	}
	return false
}

// TODO: Two methods, or one method with a boolean flag determining behavior?
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
