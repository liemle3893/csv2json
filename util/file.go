package util

import (
	"io/ioutil"
	"os"
	"regexp"
)

// ListFiles list all files in `directory` that match `includePatterns` and unmatch `excludePatterns`
func ListFiles(directory string,
	includePatterns []string,
	excludePatterns []string) ([]os.FileInfo, error) {

	var result []os.FileInfo = []os.FileInfo{}
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return []os.FileInfo{}, err
	}

	includeMatcher := createPatterns(includePatterns)
	excludeMatcher := createPatterns(excludePatterns)
	for _, file := range files {
		fileName := file.Name()
		include := includeMatcher.Match(fileName)
		exclude := excludeMatcher.Match(fileName)
		if include && !exclude {
			result = append(result, file)
		}
	}
	return result, nil
}

func createPatterns(patterns []string) *RegexesMatcher {
	regexes := make([]*regexp.Regexp, len(patterns))
	for _, pattern := range patterns {
		r, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}
		regexes = append(regexes, r)
	}
	return &RegexesMatcher{regexes: regexes}
}
