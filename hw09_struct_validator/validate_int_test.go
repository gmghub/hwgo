package hw09structvalidator

import (
	"errors"
	"testing"
)

func TestValidateInt(t *testing.T) {
	// nolint: lll
	cases := []struct {
		name     string
		in       int
		tags     []string
		expected error
	}{
		{name: "no tag", in: 0, tags: []string{}, expected: nil},

		{name: "invalid (zero in)", in: 0, tags: []string{"min:10"}, expected: ErrInvalidValue},
		{name: "invalid (in < min)", in: 9, tags: []string{"min:10"}, expected: ErrInvalidValue},
		{name: "valid (in > min)", in: 11, tags: []string{"min:10"}, expected: nil},
		{name: "valid (in == min)", in: 10, tags: []string{"min:10", "max:50"}, expected: nil},
		{name: "invalid (negative in)", in: -15, tags: []string{"min:-10"}, expected: ErrInvalidValue},
		{name: "valid (negative in)", in: -5, tags: []string{"min:-10"}, expected: nil},

		{name: "invalid (in > max)", in: 51, tags: []string{"max:50"}, expected: ErrInvalidValue},
		{name: "valid (in < max)", in: 0, tags: []string{"max:50"}, expected: nil},
		{name: "valid (in == max)", in: 50, tags: []string{"max:50"}, expected: nil},

		{name: "invalid (tag in error)", in: 0, tags: []string{"in"}, expected: ErrInvalidTag},
		{name: "invalid (excludes in)", in: 0, tags: []string{"in:10"}, expected: ErrInvalidExcludeIn},
		{name: "valid (includes in)", in: 10, tags: []string{"in:10"}, expected: nil},

		{name: "invalid (all rules)", in: 11, tags: []string{"min:10", "max:50", "in:10,20,30"}, expected: ErrInvalidExcludeIn},
		{name: "valid (all rules)", in: 10, tags: []string{"min:10", "max:50", "in:10,20,30"}, expected: nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validateInt(c.in, c.tags)
			if !errors.Is(err, c.expected) {
				t.Logf("invalid return value: in='%v' tags='%v' expected='%v' get='%v'", c.in, c.tags, c.expected, err)
				t.Fail()
			}
		})
	}
}
