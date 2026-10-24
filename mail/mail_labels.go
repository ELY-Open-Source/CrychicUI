package mail

import (
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
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, label := range labels {
			if label.ID == "" || label.Name == "" {
				panic("mail: invalid label")
			}
			c.Key(label.ID)
			on := (*selected)[label.ID]
			chip := button(c, label.Name).Row().Gap(t.Space(1)).AlignItems(ui.Center).Radius(999).Checked(on)
			chip.Background(chip.AnimateColor("selected", t.Surface.Mix(t.Accent.Alpha(0.12), unit(on)), ui.Snappy))
			chip.Children(func() {
				ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(label.Color)
				ui.Text(c, label.Name)
			})
			if chip.Clicked() {
				(*selected)[label.ID] = !on
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("MailLabels", "toggled", "on", !on)
			}
		}
	})
	return root
}
