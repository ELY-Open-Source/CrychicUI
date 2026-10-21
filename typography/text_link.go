package typography

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// link builds a text link with an animated underline.
func link(c *ui.Context, label, url string, external bool) ui.Element {
	if label == "" || url == "" {
		panic("typography: link needs a label and URL")
	}
	t := c.Theme()
	root := ui.ButtonBase(c).Role(ui.RoleLink).AlignSelf(ui.Start).Gap(t.Space(1)).TextColor(t.Accent)
	target := float32(0)
	if root.Hovered() {
		target = 1
	}
	hover := root.AnimateTo("hover", target, ui.Snappy)
	root.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if hover > 0 {
			p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 1, W: r.W * hover, H: 1}, t.Accent, 0)
		}
	})
	root.Children(func() {
		ui.Text(c, label)
		if external {
			arrow := ui.Icon(c, icon.Must("arrow-up-right")).Size(16, 16)
			arrow.Translate(2*hover, 0)
		}
	})
	if kit.Clicked(root, sound.Tap) {
		c.OpenURL(url)
	}
	return root
}

// TextLink is a link whose underline grows on hover.
func TextLink(c *ui.Context, label, url string) ui.Element {
	return link(c, label, url, false)
}
