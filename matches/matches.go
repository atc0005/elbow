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

// Package matches provides types and functions intended to help with
// collecting and validating file search results against required criteria.
package matches

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/atc0005/elbow/config"
	"github.com/atc0005/elbow/units"
	"github.com/sirupsen/logrus"
)

// FileMatch represents a superset of statistics (including os.FileInfo) for a
// file matched by provided search criteria. This allows us to record the
// original full path while also recording file metadata used in later
// calculations.
type FileMatch struct {
	os.FileInfo
	Path string
}

// FileMatches is a slice of FileMatch objects that represents the search
// results based on user-specified criteria.
type FileMatches []FileMatch

// TotalFileSize returns the cumulative size of all files in the slice in bytes
func (fm FileMatches) TotalFileSize() int64 {

	var totalSize int64

	for _, file := range fm {

		totalSize += file.Size()
	}

	return totalSize

}

// TotalFileSizeHR returns a human-readable string of the cumulative size of
// all files in the slice of bytes
func (fm FileMatches) TotalFileSizeHR() string {
	return units.ByteCountIEC(fm.TotalFileSize())
}

// SizeHR returns a human-readable string of the size of a FileMatch object.
func (fm FileMatch) SizeHR() string {
	return units.ByteCountIEC(fm.Size())
}

// HasMatchingExtension validates whether a file has the desired extension
func HasMatchingExtension(filename string, config *config.Config) bool {

	log := config.Logger

	// NOTE: We do NOT compare extensions insensitively. We can add that
	// functionality in the future if needed.
	ext := filepath.Ext(filename)

	if len(config.FileExtensions) == 0 {
		log.Debug("No extension limits have been set!")
		log.Debugf("Considering %s safe for removal\n", filename)
		return true
	}

	if InList(ext, config.FileExtensions) {
		log.Debugf("%s has a valid extension for removal\n", filename)
		return true
	}

	log.Debug("HasMatchingExtension: returning false for:", filename)
	log.Debugf("HasMatchingExtension: returning false (%q not in %q)",
		ext, config.FileExtensions)
	return false
}

// HasMatchingFilenamePattern validates whether a filename matches the desired
// pattern
func HasMatchingFilenamePattern(filename string, config *config.Config) bool {

	log := config.Logger

	if strings.TrimSpace(config.FilePattern) == "" {
		log.Debug("No FilePattern has been specified!")
		log.Debugf("Considering %s safe for removal\n", filename)
		return true
	}

	// Search for substring
	if strings.Contains(filename, config.FilePattern) {
		log.Debug("HasMatchingFilenamePattern: returning true for:", filename)
		log.Debugf("HasMatchingFilenamePattern: returning true (%q contains %q)",
			filename, config.FilePattern)
		return true
	}

	log.Debug("HasMatchingFilenamePattern: returning false for:", filename)
	log.Debugf("HasMatchingFilenamePattern: returning false (%q does not contain %q)",
		filename, config.FilePattern)
	return false
}

// HasMatchingAge validates whether a file matches the desired age threshold
func HasMatchingAge(file os.FileInfo, config *config.Config) bool {

	log := config.Logger

	// used by this function's context logger and for return code
	var ageCheckResults bool

	now := time.Now()
	fileModTime := file.ModTime()

	// common fields that we can apply to all messages in this function
	contextLogger := log.WithFields(logrus.Fields{
		"file_mod_time": fileModTime.Format(time.RFC3339),
		"current_time":  now.Format(time.RFC3339),
		"file_age_flag": config.FileAge,
		"filename":      file.Name(),
	})

	// The default for this flag is 0, so only a positive, non-zero number
	// is considered for use with age matching.
	if config.FileAge > 0 {

		// Flip user specified number of days negative so that we can wind
		// back that many days from the file modification time. This gives
		// us our threshold to compare file modification times against.
		daysBack := -(config.FileAge)
		fileAgeThreshold := now.AddDate(0, 0, daysBack)

		// Bundle more fields now that we have access to the data
		contextLogger = contextLogger.WithFields(logrus.Fields{
			"file_age_threshold": fileAgeThreshold.Format(time.RFC3339),
			"days_back":          daysBack,
		})

		contextLogger.Debug("Before age check")

		switch {
		case fileModTime.Equal(fileAgeThreshold):
			ageCheckResults = true
			contextLogger.WithFields(logrus.Fields{
				"safe_for_removal": ageCheckResults,
			}).Debug("HasMatchingAge: file mod time is equal to threshold")

		case fileModTime.Before(fileAgeThreshold):
			ageCheckResults = true
			contextLogger.WithFields(logrus.Fields{
				"safe_for_removal": ageCheckResults,
			}).Debug("HasMatchingAge: file mod time is before threshold")

		case fileModTime.After(fileAgeThreshold):
			ageCheckResults = false
			contextLogger.WithFields(logrus.Fields{
				"safe_for_removal": ageCheckResults,
			}).Debug("HasMatchingAge: file mod time is after threshold")

		}

		return ageCheckResults

	}

	contextLogger.WithFields(logrus.Fields{
		"safe_for_removal": ageCheckResults,
	}).Debugf("HasMatchingAge: age flag was not set")

	return true

}

// InList is a helper function to emulate Python's `if "x"
// in list:` functionality
func InList(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

// SortByModTimeAsc sorts slice of FileMatches in ascending order
func (fm FileMatches) SortByModTimeAsc() {
	sort.Slice(fm, func(i, j int) bool {
		return fm[i].ModTime().Before(fm[j].ModTime())
	})
}

// SortByModTimeDesc sorts slice of FileMatches in descending order
func (fm FileMatches) SortByModTimeDesc() {
	sort.Slice(fm, func(i, j int) bool {
		return fm[i].ModTime().After(fm[j].ModTime())
	})
}