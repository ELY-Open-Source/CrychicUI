package desktop

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TrafficLightsOptions name the three buttons.
type TrafficLightsOptions struct{ Close, Minimize, Zoom string }

// TrafficLights close, minimize and zoom; cue Tap.
func TrafficLights(c *ui.Context, o TrafficLightsOptions) (root, close, minimize, zoom ui.Element) {
	def(&o.Close, "Close")
	def(&o.Minimize, "Minimize")
	def(&o.Zoom, "Zoom")
	t := c.Theme()
	root = ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleToolbar)
	kit.Arrive(root, ui.Pose{Scale: 0.8, Fade: 1}, 0)
	root.Children(func() {
		light := func(label, mark string, ink ui.Color) ui.Element {
			b := kit.Press(ui.ButtonBase(c).Label(label).Size(t.Space(3.5), t.Space(3.5)).Center().Radius(999))
			hover := b.AnimateTo("hover", unit(b.Hovered()), ui.Snappy)
			b.Background(ink.Mix(t.Text, 0.12*hover)).Border(1, ink.Mix(t.Text, 0.2))
			b.Children(func() {
				glyph(c, mark, t.Space(2.5)).TextColor(t.Text.Alpha(0.7)).Opacity(hover)
			})
			if b.Clicked() {
				sound.Play(sound.Tap)
				kit.Log("TrafficLights", "pressed", "button", label)
			}
			return b
		}
		close = light(o.Close, "x", t.Danger)
		minimize = light(o.Minimize, "minus", t.Warning)
		zoom = light(o.Zoom, "maximize-2", t.Success)
	})
	return root, close, minimize, zoom
}
