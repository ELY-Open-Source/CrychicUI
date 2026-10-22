package feedback

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SaveStatusOptions tune a SaveStatus; empty copy reads English.
type SaveStatusOptions struct {
	Saving, Saved string
}

// SaveStatus cross-fades a saving spinner to a saved check.
func SaveStatus(c *ui.Context, saving bool, o SaveStatusOptions) ui.Element {
	if o.Saving == "" {
		o.Saving = "Saving…"
	}
	if o.Saved == "" {
		o.Saved = "Saved"
	}
	t := c.Theme()
	text := o.Saved
	if saving {
		text = o.Saving
	}
	root := ui.Row(c).Gap(t.Space(1.5)).AlignSelf(ui.Start).Role(ui.RoleStatus).Label(text).TextColor(t.TextMuted)
	return root.Children(func() {
		ui.Box(c).Size(16, 16).Center().Children(func() {
			if saving {
				c.Key("spin")
				Spin(c, o.Saving, SpinOptions{Size: SpinSmall}).Transition(fade)
				return
			}
			c.Key("saved")
			check := ui.Icon(c, icon.Must("check")).Size(16, 16).TextColor(t.Success)
			kit.Arrive(check, ui.Pose{Scale: 0.4, Fade: 1}, 0).Transition(fade)
		})
		ui.Text(c, text).FontSize(t.Rem(0.85)).ContentTransition(ui.ContentTransition{})
	})
}
