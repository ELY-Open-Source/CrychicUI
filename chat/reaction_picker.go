package chat

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ReactionPicker is a row of emoji; it scales in.
func ReactionPicker(c *ui.Context, label string, emoji []string, chosen *string) ui.Element {
	if label == "" || len(emoji) == 0 || chosen == nil {
		panic(fmt.Sprintf("chat: ReactionPicker %q needs a label, emoji and a choice", label))
	}
	for _, e := range emoji {
		if e == "" {
			panic(fmt.Sprintf("chat: ReactionPicker %q emoji must not be empty", label))
		}
	}
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(0.5)).Padding(t.Space(1)).Radius(999).Background(t.Background).
		Border(1, t.Border).Shadow(0, 4, 12, 0, kit.Shade(0.12)).FocusGroup(ui.Horizontal).Role(ui.RoleGroup).Label(label)
	kit.Arrive(root, ui.Pose{Scale: 0.8, Fade: 1}, 0)
	return root.Children(func() {
		for i, e := range emoji {
			c.Key(i)
			b := kit.Press(ui.ButtonBase(c).Label(e).Size(t.Space(8), t.Space(8)).Center().Radius(999).Checked(e == *chosen))
			bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
			if e == *chosen {
				bg = t.Accent.Alpha(0.15)
			}
			b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
				face := ui.Text(c, e).FontSize(t.Rem(1.2))
				k := 1 + 0.25*face.AnimateTo("grow", unit(b.Hovered()), ui.Bouncy)
				face.Scale(k, k)
			})
			if b.Clicked() {
				*chosen = e
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("ReactionPicker", "picked", "index", i)
			}
		}
	})
}
