package agent

import (
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"
)

// StreamText shows text token by token with a caret.
func StreamText(c *ui.Context, text string, done bool) ui.Element {
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1)).AlignItems(ui.End)
	shown := text
	streaming := !done && !c.Preferences().ReduceMotion
	if streaming {
		since := ui.Local(root, "since", func() time.Time { return c.Now() })
		const step = 80 * time.Millisecond
		n := int(c.Now().Sub(*since) / step)
		shown = prefixTokens(text, n)
		if len(shown) < len(text) {
			c.After(step)
		} else {
			c.After(time.Second)
		}
	}
	return root.Children(func() {
		ui.Text(c, shown).Wrap()
		if streaming {
			caret := ui.Box(c).Width(t.Space(0.5)).Height(t.Rem(1)).Shrink(0).Background(t.Text)
			ph := float64(caret.Loop("blink", time.Second, ui.Linear))
			caret.Opacity(unit(ph < 0.6))
		}
	})
}

// prefixTokens returns s through the end of its nth token.
func prefixTokens(s string, n int) string {
	if n <= 0 {
		return ""
	}
	tokens, inWord := 0, false
	for i, r := range s {
		switch {
		case unicode.IsSpace(r):
			if inWord {
				tokens++
				inWord = false
				if tokens == n {
					return s[:i]
				}
			}
		case unspaced(r):
			if inWord {
				tokens++
				inWord = false
				if tokens == n {
					return s[:i]
				}
			}
			tokens++
			if tokens == n {
				return s[:i+utf8.RuneLen(r)]
			}
		default:
			inWord = true
		}
	}
	return s
}

// unspaced reports scripts written without spaces, one token per rune.
func unspaced(r rune) bool {
	return unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana) ||
		(r >= 0x3000 && r <= 0x303f) || (r >= 0xff00 && r <= 0xffef)
}
