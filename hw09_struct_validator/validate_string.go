package hw09structvalidator

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// nolint:gocognit
func validateString(val string, tags []string) error {
	for _, t := range tags {
		tsplit := strings.SplitN(t, ":", 2)
		if len(tsplit) != 2 {
			return ErrInvalidTag
		}
		tname := tsplit[0]
		tval := tsplit[1]

		switch {
		case tname == "len":
			l, err := strconv.Atoi(tval)
			if err != nil || l < 0 {
				return ErrInvalidTag
			}
			if utf8.RuneCountInString(val) != l {
				return ErrInvalidLen
			}

		case tname == "regexp":
			rex, err := regexp.Compile(tval)
			if err != nil {
				return ErrInvalidTag
			}
			if !rex.Match([]byte(val)) {
				return ErrInvalidRegexpNotMatch
			}

		case tname == "in":
			spl := strings.Split(tval, ",")
			if len(spl) == 0 {
				return ErrInvalidTag
			}
			found := false
			for _, v := range spl {
				if val == v {
					found = true
					break
				}
			}
			if !found {
				return ErrInvalidExcludeIn
			}

		default:
			return ErrInvalidTag
		}
	}

	return nil
}
