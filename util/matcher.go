package util

import "regexp"

type RegexesMatcher struct {
	regexes []*regexp.Regexp
}

func (m *RegexesMatcher) Match(input string) bool {
	for _, regex := range m.regexes {
		if regex == nil {
			continue
		}
		if regex.MatchString(input) {
			return true
		}
	}
	return false
}
