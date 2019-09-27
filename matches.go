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

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
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

func hasMatchingExtension(filename string, config *Config) bool {

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

	log.Debug("hasMatchingExtension: returning false for:", filename)
	log.Debugf("hasMatchingExtension: returning false (%q not in %q)",
		ext, config.FileExtensions)
	return false
}

func hasMatchingFilenamePattern(filename string, config *Config) bool {

	if strings.TrimSpace(config.FilePattern) == "" {
		log.Debug("No FilePattern has been specified!")
		log.Debugf("Considering %s safe for removal\n", filename)
		return true
	}

	// Search for substring
	if strings.Contains(filename, config.FilePattern) {
		log.Debug("hasMatchingFilenamePattern: returning true for:", filename)
		log.Debugf("hasMatchingFilenamePattern: returning true (%q contains %q)",
			filename, config.FilePattern)
		return true
	}

	log.Debug("hasMatchingFilenamePattern: returning false for:", filename)
	log.Debugf("hasMatchingFilenamePattern: returning false (%q does not contain %q)",
		filename, config.FilePattern)
	return false
}

func hasMatchingAge(file os.FileInfo, config *Config) bool {

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
			}).Debug("hasMatchingAge: file mod time is equal to threshold")

		case fileModTime.Before(fileAgeThreshold):
			ageCheckResults = true
			contextLogger.WithFields(logrus.Fields{
				"safe_for_removal": ageCheckResults,
			}).Debug("hasMatchingAge: file mod time is before threshold")

		case fileModTime.After(fileAgeThreshold):
			ageCheckResults = false
			contextLogger.WithFields(logrus.Fields{
				"safe_for_removal": ageCheckResults,
			}).Debug("hasMatchingAge: file mod time is after threshold")

		}

		return ageCheckResults

	}

	contextLogger.WithFields(logrus.Fields{
		"safe_for_removal": ageCheckResults,
	}).Debugf("hasMatchingAge: age flag was not set")

	return true

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
