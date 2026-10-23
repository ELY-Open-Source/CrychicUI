package chat

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ContactMessageOptions name the add button.
type ContactMessageOptions struct {
	Add, Added string
}

// ContactMessage is a contact card in a bubble with add.
func ContactMessage(c *ui.Context, name string, image *ui.Bitmap, added *bool, o ContactMessageOptions) ui.Element {
	if name == "" || added == nil {
		panic(fmt.Sprintf("chat: ContactMessage %q needs a name and added", name))
	}
	def(&o.Add, "Add")
	def(&o.Added, "Added")
	t := c.Theme()
	root := shell(c, false).Row().Gap(t.Space(3)).AlignItems(ui.Center)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	var btn ui.Element
	root.Children(func() {
		display.Avatar(c, name, image)
		ui.Text(c, name).Bold().Grow(1).Wrap()
		label, glyph := o.Add, "user-plus"
		if *added {
			label, glyph = o.Added, "check"
		}
		btn = kit.Press(ui.ButtonBase(c).Label(label).Row().Gap(t.Space(1.5)).AlignItems(ui.Center).
			Padding(t.Space(1.5), t.Space(3)).Radius(999).Checked(*added))
		bg, fg := t.Accent.Alpha(0.12*unit(btn.Hovered())), t.Accent
		if *added {
			bg = t.Accent
			fg = t.AccentText
		}
		btn.Background(btn.AnimateColor("bg", bg, ui.Snappy)).TextColor(btn.AnimateColor("fg", fg, ui.Snappy)).
			Border(1, t.Accent).Children(func() {
			ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{})
			ui.Text(c, label).FontSize(t.Rem(0.9)).Bold().SingleLine().ContentTransition(ui.ContentTransition{})
		})
	})
	if btn.Clicked() {
		*added = !*added
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("ContactMessage", "added", *added)
	}
	return root
}
