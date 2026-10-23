package chat

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PinnedOptions tune a Pinned bar.
type PinnedOptions struct {
	// Label names the bar for assistive technology.
	Label string
}

// Pinned is a pinned message bar; a click jumps.
func Pinned(c *ui.Context, sender, text string, o PinnedOptions) ui.Element {
	if sender == "" || text == "" {
		panic(fmt.Sprintf("chat: Pinned %q needs a sender and text", sender))
	}
	def(&o.Label, "Pinned message")
	t := c.Theme()
	root := kit.Press(ui.ButtonBase(c).Label(o.Label).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
		Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Border(1, t.Border))
	root.Background(root.AnimateColor("bg", t.Surface.Mix(t.SurfaceHover, unit(root.Hovered())), ui.Snappy))
	kit.Arrive(root, ui.Pose{Y: -t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Icon(c, icon.Must("pin")).Size(16, 16).TextColor(t.Accent)
		ui.Column(c).Grow(1).Children(func() {
			ui.Text(c, sender).Bold().FontSize(t.Rem(0.8)).TextColor(t.Accent).SingleLine()
			ui.Text(c, text).SingleLine().FontSize(t.Rem(0.9))
		})
	})
	if kit.Clicked(root, sound.Tap) {
		kit.Log("Pinned", "jumped")
	}
	return root
}
