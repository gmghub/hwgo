package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var (
		prevRune rune
		repRune  rune
		repCnt   int
		builder  strings.Builder
	)

	for _, r := range s {
		if unicode.IsDigit(r) {
			if prevRune == '\\' {
				if repRune != 0 {
					builder.WriteRune(repRune)
				}
				repRune = r
				prevRune = r
				continue
			}
			if repRune != 0 {
				repCnt, _ = strconv.Atoi(string(r))
				for i := 0; i < repCnt; i++ {
					builder.WriteRune(repRune)
				}
				repRune = 0
				prevRune = r
				continue
			}
			return "", ErrInvalidString
		}

		if r == '\\' {
			if repRune != 0 {
				builder.WriteRune(repRune)
				repRune = 0
			}
			if prevRune == '\\' {
				prevRune = 0
				repRune = r
			} else {
				prevRune = r
			}
			continue
		}

		if repRune != 0 {
			builder.WriteRune(repRune)
		}
		repRune = r
		prevRune = r
	}

	if repRune != 0 {
		builder.WriteRune(repRune)
	}

	return builder.String(), nil
}
