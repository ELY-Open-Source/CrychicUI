package kit

import (
	"strings"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"
)

// IndexFold finds q in s ignoring case, as byte bounds.
func IndexFold(s, q string) (int, int) {
	n := utf8.RuneCountInString(q)
	for i := range s {
		end := i
		for k := 0; k < n && end < len(s); k++ {
			_, w := utf8.DecodeRuneInString(s[end:])
			end += w
		}
		if strings.EqualFold(s[i:end], q) {
			return i, end
		}
	}
	return -1, -1
}

// Marked spans s with lo..hi bold in the accent.
func Marked(t *ui.Theme, s string, lo, hi int) []ui.Span {
	if lo < 0 || lo == hi {
		return []ui.Span{{Text: s}}
	}
	return []ui.Span{{Text: s[:lo]}, {Text: s[lo:hi], Color: t.Accent, Weight: 700}, {Text: s[hi:]}}
}
