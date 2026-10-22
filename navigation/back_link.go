package navigation

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BackLink is an arrow and label; hover nudges it.
func BackLink(c *ui.Context, label string) ui.Element {
	if label == "" {
		panic("navigation: BackLink needs a label")
	}
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).AlignSelf(ui.Start).Gap(t.Space(1.5)).Padding(t.Space(1), t.Space(1.5)).
		Radius(t.Radius).TextColor(t.Accent))
	back := float32(-1)
	if c.Direction() == ui.RightToLeft {
		back = 1
	}
	hot := b.Hovered() || b.Focused()
	b.Children(func() {
		arrow := ui.Icon(c, icon.Must("arrow-left")).Size(16, 16)
		arrow.Translate(back*t.Space(1)*arrow.AnimateTo("nudge", unit(hot), ui.Bouncy), 0)
		ui.Text(c, label)
	})
	kit.Clicked(b, sound.Tap)
	return b
}
