package feedback

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BannerOptions tune a Banner.
type BannerOptions struct {
	// Tone colors the bar; ToneInfo when zero.
	Tone Tone
	// Close names the close button; "Dismiss" when empty.
	Close string
}

// Banner is a full-width bar that slides down.
func Banner(c *ui.Context, message, action string, open *bool, o BannerOptions) (root, act ui.Element) {
	if message == "" || open == nil {
		panic("feedback: Banner needs a message and open state")
	}
	if o.Close == "" {
		o.Close = "Dismiss"
	}
	t := c.Theme()
	col, glyph := look(t, o.Tone)
	root = ui.Column(c).FillWidth().ClipY()
	root.Children(func() {
		if !*open {
			return
		}
		c.Key("bar")
		bar := ui.Row(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(2), t.Space(4)).Background(col.Alpha(0.12)).
			BorderWidth(0, 0, 1, 0).BorderColor(col.Alpha(0.3)).Role(ui.RoleStatus).Label(message)
		kit.Arrive(bar, ui.Pose{Y: -t.Space(12)}, 0).Transition(collapse).Children(func() {
			ui.Icon(c, icon.Must(glyph)).Size(16, 16).TextColor(col)
			ui.Text(c, message).Grow(1)
			if action != "" {
				act = button(c, action, true)
				if kit.Clicked(act, sound.Tap) {
					kit.Log("Banner", "action")
				}
			}
			if kit.Clicked(closer(c, o.Close), sound.Close) {
				*open = false
				root.SetChanged()
				kit.Log("Banner", "closed")
			}
		})
	})
	return root, act
}
