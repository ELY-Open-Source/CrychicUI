package code

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ErrorOverlayOptions tune an ErrorOverlay.
type ErrorOverlayOptions struct {
	Close string
}

// ErrorOverlay slides a compile error panel down over content.
func ErrorOverlay(c *ui.Context, title, detail string, open *bool, o ErrorOverlayOptions, content func()) ui.Element {
	if title == "" || open == nil || content == nil {
		panic("code: ErrorOverlay needs a title, open and content")
	}
	def(&o.Close, "Close")
	t := c.Theme()
	root := ui.Box(c).Clip()
	height := ui.Local(root, "height", func() float32 { return 0 })
	if !*open {
		*height = 0
	}
	root.MinHeight(*height)
	root.Children(func() {
		content()
		if !*open {
			return
		}
		panel := ui.Column(c).Key("err").Absolute().Top(0).Left(0).Right(0).Gap(t.Space(1)).Padding(t.Space(3)).
			Background(t.Background).Border(1, t.Danger).Radius(t.Radius).Shadow(0, 8, 24, 0, kit.Shade(0.16)).
			Role(ui.RoleDialog).Label(title)
		kit.Arrive(panel, ui.Pose{Y: -t.Space(8)}, 0)
		panel.Transition(ui.ElementTransition{Duration: 180 * time.Millisecond, Exit: &ui.Motion{Y: -t.Space(8)}})
		var shut ui.Element
		panel.Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icon.Must("circle-x")).Size(16, 16).TextColor(t.Danger)
				ui.Text(c, title).Bold().Grow(1).SingleLine()
				shut = round(c, o.Close, "x").Size(t.Space(6), t.Space(6))
			})
			if detail != "" {
				ui.Text(c, detail).Font("monospace").FontSize(t.Rem(0.85)).Wrap()
			}
		})
		panel.DrawOver(func(p *ui.Painter, r ui.Rect) {
			if *open && r.H != *height {
				*height = r.H
				c.Invalidate()
			}
		})
		if shut.Clicked() {
			*open = false
			root.SetChanged()
			sound.Play(sound.Close)
			kit.Log("ErrorOverlay", "closed")
		}
	})
	return root
}
