package mail

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MailLabel is one colored mail label.
type MailLabel struct {
	ID, Name string
	Color    ui.Color
}

// MailLabels toggles colored labels; cue Toggle.
func MailLabels(c *ui.Context, labels []MailLabel, selected *map[string]bool) ui.Element {
	if len(labels) == 0 || selected == nil || *selected == nil {
		panic("mail: MailLabels needs labels and state")
	}
	ids := make([]string, len(labels))
	for i, label := range labels {
		if label.Name == "" {
			panic(fmt.Sprintf("mail: MailLabels label %q needs a name", label.ID))
		}
		ids[i] = label.ID
	}
	unique("MailLabels", ids)
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).Role(ui.RoleGroup)
	room := root.LayoutBounds().W
	kit.Settle(c, root, []float32{room}, func() []float32 { return []float32{root.LayoutBounds().W} })
	// Padding, dot, check and gaps around the name.
	text := max(0, room-t.Space(10)-14)
	root.Children(func() {
		for i, label := range labels {
			c.Key(label.ID)
			on := (*selected)[label.ID]
			chip := kit.Press(ui.ButtonBase(c).Label(label.Name).Padding(t.Space(1), t.Space(2.5)).Radius(999).Shrink(0).Checked(on))
			chip.Background(chip.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(chip.Hovered())).Mix(label.Color.Alpha(0.14), unit(on)), ui.Snappy))
			chip.Border(1, chip.AnimateColor("line", t.Border.Mix(label.Color, unit(on)), ui.Snappy))
			kit.Stagger(c, root, chip, i).Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Size(t.Space(2), t.Space(2)).Shrink(0).Radius(999).Background(label.Color)
					ui.Text(c, label.Name).FontSize(t.Rem(0.8)).Shrink(0).MaxWidth(text).SingleLine()
					c.Key(on)
					if on {
						kit.Arrive(glyph(c, "check", 12).TextColor(label.Color), ui.Pose{Scale: 0.4, Fade: 1}, 0)
					}
				})
			})
			if chip.Clicked() {
				(*selected)[label.ID] = !on
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("MailLabels", "toggled", "index", i, "on", !on)
			}
		}
	})
	return root
}
