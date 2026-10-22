package motion

import (
	"strings"
	"time"
	"unicode"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TypewriterOptions tune a Typewriter.
type TypewriterOptions struct {
	// Every is the time per letter; 50 ms.
	Every time.Duration
}

var blink = []ui.Keyframe{{Value: 1}, {At: 500 * time.Millisecond, Value: 1}, {At: 520 * time.Millisecond, Value: 0},
	{At: 1000 * time.Millisecond, Value: 0}, {At: 1020 * time.Millisecond, Value: 1}}

// Typewriter types text letter by letter, ticking per word.
func Typewriter(c *ui.Context, text string, o TypewriterOptions) ui.Element {
	if text == "" || o.Every < 0 {
		panic("motion: Typewriter needs text and a positive Every")
	}
	if o.Every == 0 {
		o.Every = 50 * time.Millisecond
	}
	t := c.Theme()
	runes := []rune(text)
	root := ui.Box(c)
	return root.Children(func() {
		line := ui.Row(c.Key(text)).AlignItems(ui.Center)
		n := int(line.Keyframes("crychic.type", []ui.Keyframe{{Value: 0},
			{At: time.Duration(len(runes)) * o.Every, Value: float32(len(runes)), Ease: ui.Linear}}))
		shown := string(runes[:n])
		done := len(strings.Fields(shown))
		if n < len(runes) && n > 0 && !unicode.IsSpace(runes[n-1]) && !unicode.IsSpace(runes[n]) {
			done--
		}
		typed := ui.Local(line, "words", func() int { return 0 })
		for ; *typed < done; *typed++ {
			sound.Play(sound.Tick)
		}
		line.Children(func() {
			ui.Text(c, shown).Label(text)
			caret := ui.Box(c).Width(t.Space(0.5)).Height(t.Rem(1.2)).Radius(t.Space(0.25)).Background(t.Accent)
			caret.Opacity(caret.Keyframes("blink", blink, forever))
		})
	})
}
