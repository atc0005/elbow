package main

import (
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
// TODO: Do I really need to abstract the fact that FileMatches is a slice of
// FileMatch objects? It seems that by hiding this it makes it harder to see
// that we're working with a slice?
type FileMatches []FileMatch

func hasValidExtension(filename string, config *Config) bool {

	// NOTE: We do NOT compare extensions insensitively. We can add that
	// functionality in the future if needed.
	ext := filepath.Ext(filename)

	if len(config.FileExtensions) == 0 {
		log.Debug("No extension limits have been set!")
		log.Debugf("Considering %s safe for removal\n", filename)
		return true
	}

	if inList(ext, config.FileExtensions) {
		log.Debugf("%s has a valid extension for removal\n", filename)
		return true
	}

	log.Debug("hasValidExtension: returning false for:", filename)
	log.Debugf("hasValidExtension: returning false (%q not in %q)",
		ext, config.FileExtensions)
	return false
}

func hasValidFilenamePattern(filename string, config *Config) bool {

	if strings.TrimSpace(config.FilePattern) == "" {
		log.Debug("No FilePattern has been specified!")
		log.Debugf("Considering %s safe for removal\n", filename)
		return true
	}

	// Search for substring
	if strings.Contains(filename, config.FilePattern) {
		log.Debug("hasValidFilenamePattern: returning true for:", filename)
		log.Debugf("hasValidFilenamePattern: returning true (%q contains %q)",
			filename, config.FilePattern)
		return true
	}

	log.Debug("hasValidFilenamePattern: returning false for:", filename)
	log.Debugf("hasValidFilenamePattern: returning false (%q does not contain %q)",
		filename, config.FilePattern)
	return false
}

// inList is a helper function to emulate Python's `if "x"
// in list:` functionality
func inList(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

// TODO: Two methods, or one method with a boolean flag determining behavior?
func (fm FileMatches) sortByModTimeAsc() {
	sort.Slice(fm, func(i, j int) bool {
		return fm[i].ModTime().Before(fm[j].ModTime())
	})
}

func (fm FileMatches) sortByModTimeDesc() {
	sort.Slice(fm, func(i, j int) bool {
		return fm[i].ModTime().After(fm[j].ModTime())
	})
}
