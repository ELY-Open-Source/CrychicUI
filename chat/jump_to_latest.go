package chat

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// JumpToLatestOptions word the count, "%d" taking it.
type JumpToLatestOptions struct {
	One, Many string
}

// JumpToLatest is a pill with a count after scrolling up.
func JumpToLatest(c *ui.Context, count int, o JumpToLatestOptions) ui.Element {
	if count < 0 {
		panic("chat: JumpToLatest needs a count from zero")
	}
	def(&o.One, "%d new message")
	def(&o.Many, "%d new messages")
	if count == 0 {
		return ui.Box(c)
	}
	word := o.Many
	if count == 1 {
		word = o.One
	}
	label := fmt.Sprintf(word, count)
	t := c.Theme()
	root := kit.Press(ui.ButtonBase(c).Label(label).Row().AlignSelf(ui.Start).Gap(t.Space(1.5)).AlignItems(ui.Center).
		Padding(t.Space(1.5), t.Space(3)).Radius(999).Background(t.Accent).TextColor(t.AccentText))
	kit.Arrive(root, ui.Pose{Scale: 0.8, Fade: 1}, 0)
	root.Children(func() {
		ui.Icon(c, icon.Must("arrow-down")).Size(16, 16)
		kit.Numeric(c, label).FontSize(t.Rem(0.9)).Bold().SingleLine()
	})
	if kit.Clicked(root, sound.Tap) {
		root.SetChanged()
		kit.Log("JumpToLatest", "jumped", "count", count)
	}
	return root
}
