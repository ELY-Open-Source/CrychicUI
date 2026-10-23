package finance

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CardControlsOptions are the card control labels.
type CardControlsOptions struct{ Freeze, Online, Limit string }

// CardControls toggles freeze, online, and limits; cue Toggle.
func CardControls(c *ui.Context, frozen, online *bool, limit *float64, o CardControlsOptions) ui.Element {
	if frozen == nil || online == nil || limit == nil || *limit < 0 || *limit > 10000 || !finite(*limit) {
		panic("finance: CardControls has invalid data")
	}
	def(&o.Freeze, "Freeze card")
	def(&o.Online, "Online payments")
	def(&o.Limit, "Spending limit")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, item := range []struct {
			name  string
			value *bool
		}{{o.Freeze, frozen}, {o.Online, online}} {
			row := kit.Press(ui.ButtonBase(c).Label(item.name).Row().FillWidth().Padding(t.Space(2)).Radius(t.Radius).
				Border(1, t.Border).Checked(*item.value))
			row.Background(row.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*item.value)), ui.Snappy))
			row.Children(func() { ui.Text(c, item.name).SingleLine() })
			if row.Clicked() {
				*item.value = !*item.value
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("CardControls", "toggled", "control", item.name, "value", *item.value)
			}
		}
		ui.Text(c, o.Limit).FontSize(t.Rem(0.8))
		slider := ui.Slider(c, limit, 0, 10000).Label(o.Limit).FillWidth()
		if kit.Changed(slider, sound.Tick) {
			root.SetChanged()
			kit.Log("CardControls", "limit", "value", *limit)
		}
	})
	return root
}
