package text

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func StripAccents(value string) string {
	decomposed := norm.NFD.String(value)

	var stripped strings.Builder
	stripped.Grow(len(decomposed))

	for _, r := range decomposed {
		if unicode.Is(unicode.Mn, r) {
			continue
		}

		stripped.WriteRune(r)
	}

	return stripped.String()
}

func EqualIgnoringAccents(left, right string) bool {
	return strings.EqualFold(StripAccents(left), StripAccents(right))
}
