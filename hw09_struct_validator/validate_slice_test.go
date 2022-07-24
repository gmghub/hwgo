package hw09structvalidator

import (
	"errors"
	"testing"
)

func TestValidateSlice(t *testing.T) {
	// nolint: lll
	cases := []struct {
		name     string
		in       interface{}
		tags     []string
		expected error
	}{
		{name: "nil in", in: nil, tags: []string{}, expected: nil},

		// // []int tests
		{name: "empty in", in: []int{}, tags: []string{"min:10"}, expected: nil},

		{name: "invalid (in < min)", in: []int{11, 10, 9}, tags: []string{"min:10"}, expected: ErrInvalidValue},
		{name: "valid (in >= min)", in: []int{10, 11, 12}, tags: []string{"min:10"}, expected: nil},

		{name: "invalid (in > max)", in: []int{0, 50, 51}, tags: []string{"max:50"}, expected: ErrInvalidValue},
		{name: "valid (in <= max)", in: []int{0, 1, 2, 49, 50}, tags: []string{"max:50"}, expected: nil},

		{name: "invalid (excludes in)", in: []int{10, 20, 31}, tags: []string{"in:10,20,30"}, expected: ErrInvalidExcludeIn},
		{name: "valid (includes in)", in: []int{10, 20, 30}, tags: []string{"in:10,20,30"}, expected: nil},

		{name: "invalid (all rules)", in: []int{10, 20, 31}, tags: []string{"min:10", "max:50", "in:10,20,30"}, expected: ErrInvalidExcludeIn},
		{name: "valid (all rules)", in: []int{10, 20, 30}, tags: []string{"min:10", "max:50", "in:10,20,30"}, expected: nil},

		// []string tests
		{name: "empty in", in: []string{}, tags: []string{"len:3"}, expected: nil},

		{name: "invalid (len < 3)", in: []string{"aaa", "bbb", "cc"}, tags: []string{"len:3"}, expected: ErrInvalidLen},
		{name: "valid (len = 5)", in: []string{"aaa", "bbb", "ccc"}, tags: []string{"len:3"}, expected: nil},

		{name: "invalid tag (regexp invalid)", in: []string{"aaa", "bbb", "abc"}, tags: []string{"regexp:[0-9]++"}, expected: ErrInvalidTag},
		{name: "invalid (regexp not match)", in: []string{"aaa@b.com", "bbb@b.com", "abc@"}, tags: []string{"regexp:^\\w+@\\w+\\.\\w+$"}, expected: ErrInvalidRegexpNotMatch},
		{name: "valid (regexp match)", in: []string{"aaa@b.com", "bbb@b.com", "abc@email.com"}, tags: []string{"regexp:^\\w+@\\w+\\.\\w+$"}, expected: nil},
		{name: "valid (regexp with :)", in: []string{"abc:cba", "bcd:dcb", "ecb:bce"}, tags: []string{"regexp:^\\w+:\\w+$"}, expected: nil},

		{name: "invalid (excludes in)", in: []string{"aaa", "bbb", "ccc"}, tags: []string{"in:aaa,bbb"}, expected: ErrInvalidExcludeIn},
		{name: "valid (includes in)", in: []string{"aaa", "bbb", "ccc"}, tags: []string{"in:aaa,bbb,ccc"}, expected: nil},

		{name: "invalid (all rules)", in: []string{"101", "202", "30a"}, tags: []string{"len:3", "regexp:\\w{3}", "in:101,202,303"}, expected: ErrInvalidExcludeIn},
		{name: "valid (all rules)", in: []string{"101", "202", "303"}, tags: []string{"len:3", "regexp:\\w{3}", "in:101,202,303"}, expected: nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validateSlice(c.in, c.tags)
			if !errors.Is(err, c.expected) {
				t.Logf("invalid return value: in='%v' tags='%v' expected='%v' get='%v'", c.in, c.tags, c.expected, err)
				t.Fail()
			}
		})
	}
}
