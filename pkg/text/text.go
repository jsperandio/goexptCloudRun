package text

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func StripAccents(value string) string {
	decomposed := norm.NFD.String(value)

	var strB strings.Builder
	strB.Grow(len(decomposed))

	for _, r := range decomposed {
		if unicode.Is(unicode.Mn, r) {
			continue
		}

		strB.WriteRune(r)
	}

	return strB.String()
}

func EqualIgnoringAccents(left, right string) bool {
	return strings.EqualFold(StripAccents(left), StripAccents(right))
}
