package chat

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ThreadOptions tune a Thread.
type ThreadOptions struct {
	// One and Many word the count, "%d" taking it.
	One, Many string
}

// Thread shows a reply count that opens a side thread.
func Thread(c *ui.Context, replies int, open *bool, o ThreadOptions) ui.Element {
	if replies < 0 || open == nil {
		panic("chat: Thread needs a reply count and open")
	}
	def(&o.One, "%d reply")
	def(&o.Many, "%d replies")
	word := o.Many
	if replies == 1 {
		word = o.One
	}
	label := fmt.Sprintf(word, replies)
	t := c.Theme()
	root := kit.Press(ui.ButtonBase(c).Label(label).Row().AlignSelf(ui.Start).Gap(t.Space(1.5)).AlignItems(ui.Center).
		Padding(t.Space(1), t.Space(2)).Radius(t.Radius).Checked(*open))
	root.Background(root.AnimateColor("bg", t.Accent.Alpha(0.12*unit(root.Hovered() || *open)), ui.Snappy)).TextColor(t.Accent)
	root.Children(func() {
		ui.Icon(c, icon.Must("message-circle")).Size(16, 16)
		kit.Numeric(c, label).FontSize(t.Rem(0.9)).Bold().SingleLine()
		ch := ui.Icon(c, icon.Must("chevron-right")).Size(16, 16)
		ch.Rotate(90 * ch.AnimateTo("turn", unit(*open), ui.Snappy))
	})
	if root.Clicked() {
		*open = !*open
		root.SetChanged()
		if *open {
			sound.Play(sound.Open)
		} else {
			sound.Play(sound.Close)
		}
		kit.Log("Thread", "open", *open)
	}
	return root
}
