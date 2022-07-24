package hw09structvalidator

import (
	"errors"
	"testing"
)

func TestValidateString(t *testing.T) {
	//nolint:lll
	cases := []struct {
		name     string
		in       string
		tags     []string
		expected error
	}{
		{name: "no tag", in: "", tags: []string{}, expected: nil},
		{name: "invalid tag", in: "", tags: []string{"len"}, expected: ErrInvalidTag},

		{name: "invalid (len < 5)", in: "abcd", tags: []string{"len:5"}, expected: ErrInvalidLen},
		{name: "invalid (len > 5)", in: "abcdef", tags: []string{"len:5"}, expected: ErrInvalidLen},
		{name: "valid (len = 5)", in: "abcde", tags: []string{"len:5"}, expected: nil},

		{name: "invalid tag (regexp invalid)", in: "abc", tags: []string{"regexp:[0-9]++"}, expected: ErrInvalidTag},
		{name: "invalid (regexp not match)", in: "abc", tags: []string{"regexp:^\\w+@\\w+\\.\\w+$"}, expected: ErrInvalidRegexpNotMatch},
		{name: "valid (regexp match)", in: "abc@email.com", tags: []string{"regexp:^\\w+@\\w+\\.\\w+$"}, expected: nil},
		{name: "valid (regexp with :)", in: "abc:cba", tags: []string{"regexp:^\\w+:\\w+$"}, expected: nil},

		{name: "invalid (tag in error)", in: "abc", tags: []string{"in"}, expected: ErrInvalidTag},
		{name: "invalid (excludes in)", in: "abc", tags: []string{"in:aaa,bbb,ccc"}, expected: ErrInvalidExcludeIn},
		{name: "valid (includes in)", in: "abc", tags: []string{"in:aaa,abc,bbb"}, expected: nil},

		{name: "invalid (all rules)", in: "abc", tags: []string{"len:3", "regexp:^\\w+$", "in:aaa,bbb,ccc"}, expected: ErrInvalidExcludeIn},
		{name: "valid (all rules)", in: "abc", tags: []string{"len:3", "regexp:^\\w+$", "in:aaa,bbb,ccc,abc"}, expected: nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validateString(c.in, c.tags)
			if !errors.Is(err, c.expected) {
				t.Logf("invalid return value: in='%v' tags='%v' expected='%v' get='%v'", c.in, c.tags, c.expected, err)
				t.Fail()
			}
		})
	}
}
