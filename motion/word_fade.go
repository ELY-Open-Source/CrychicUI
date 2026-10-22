package motion

import (
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

// WordFade fades the words of text in one after another.
func WordFade(c *ui.Context, text string) ui.Element {
	words := strings.Fields(text)
	if len(words) == 0 {
		panic("motion: WordFade needs words")
	}
	t := c.Theme()
	root := ui.Box(c)
	return root.Children(func() {
		ui.Row(c.Key(text)).Wrap().Gap(t.Rem(0.3)).Children(func() {
			for i, w := range words {
				c.Key(i)
				e := ui.Text(c, w)
				k := e.Keyframes("crychic.in", arrive, ui.KeyframeOptions{Delay: time.Duration(i) * 90 * time.Millisecond})
				e.Opacity(k).Translate(0, (1-k)*t.Space(1.5))
			}
		})
	})
}
