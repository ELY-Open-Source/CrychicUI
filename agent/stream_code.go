package agent

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// StreamCode is a code block filling as it streams.
func StreamCode(c *ui.Context, lang, code string, done bool) ui.Element {
	if code == "" {
		panic("agent: StreamCode needs code")
	}
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).MaxWidth(t.Space(100)).Gap(t.Space(1)).Padding(t.Space(2), t.Space(3)).
		Radius(t.Radius*2).Background(t.Surface).Border(1, t.Border).Role(ui.RoleGroup).Label(lang)
	shown := code
	streaming := !done && !c.Preferences().ReduceMotion
	if streaming {
		since := ui.Local(root, "since", func() time.Time { return c.Now() })
		const step = 60 * time.Millisecond
		n := int(c.Now().Sub(*since) / step)
		shown = prefixLines(code, n)
		if len(shown) < len(code) {
			c.After(step)
		} else {
			c.After(time.Second)
		}
	}
	return root.Children(func() {
		kit.Muted(c, lang).FontSize(t.Rem(0.75)).Bold().SingleLine()
		ui.Row(c).Gap(t.Space(1)).AlignItems(ui.End).Children(func() {
			ui.Text(c, shown).Font("monospace").FontSize(t.Rem(0.85)).Wrap()
			if streaming {
				caret := ui.Box(c).Width(t.Space(0.5)).Height(t.Rem(1)).Shrink(0).Background(t.Text)
				ph := float64(caret.Loop("blink", time.Second, ui.Linear))
				caret.Opacity(unit(ph < 0.6))
			}
		})
	})
}
