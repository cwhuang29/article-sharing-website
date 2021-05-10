package utils

import (
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/russross/blackfriday"
)

// Each mandarin symbol takes 3 - 4 bytes.
// The following `limit` is a value to measured how many "characters" can be displayed,
// It measures not only number of characters, but also the width of each character.
// `ratio` is the width ratio of Chinese words versus English characters
func DecodeRuneStringForFrontend(s string, limit float64, ratio float64) string {
	idx := 0

	for cnt := 0.; cnt < limit; {
		_, width := utf8.DecodeRuneInString(s[idx:])
		// fmt.Printf("%#U starts at byte position %d\n", runeValue, charWidth)

		idx += width
		if width == 1 {
			cnt += 1 // e.g. English alphabets
		} else {
			cnt += ratio
		}
	}
	return s[:idx]
}

func ParseMarkdownToHTML(s string) string {
	/*
		It is such a bad idea to self-implement markdown parser
		links := regexp.MustCompile(`\[([^\s]+)\]\(([^\s]+)\)`)
		code := regexp.MustCompile("`([^\r|\n]*)`")
		s = links.ReplaceAllString(s, `<a href="$2">$1</a>`)
		s = bold.ReplaceAllString(s, `<strong>$1</strong>`)
	*/
	byteS := blackfriday.MarkdownCommon([]byte(s))
	return string(byteS)
}

func GetUUID() string {
	return uuid.NewString()
}
