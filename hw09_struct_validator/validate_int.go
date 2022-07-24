package hw09structvalidator

import (
	"strconv"
	"strings"
)

//nolint:gocognit
func validateInt(val int, tags []string) error {
	for _, t := range tags {
		tsplit := strings.SplitN(t, ":", 2)
		if len(tsplit) != 2 {
			return ErrInvalidTag
		}
		tname := tsplit[0]
		tval := tsplit[1]

		switch {
		case tname == "min":
			min, err := strconv.Atoi(tval)
			if err != nil {
				return ErrInvalidTag
			}
			if val < min {
				return ErrInvalidValue
			}

		case tname == "max":
			max, err := strconv.Atoi(tval)
			if err != nil {
				return ErrInvalidTag
			}
			if val > max {
				return ErrInvalidValue
			}

		case tname == "in":
			spl := strings.Split(tval, ",")
			if len(spl) == 0 {
				return ErrInvalidTag
			}
			found := false
			for _, v := range spl {
				n, err := strconv.Atoi(v)
				if err != nil {
					return ErrInvalidTag
				}
				if val == n {
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
