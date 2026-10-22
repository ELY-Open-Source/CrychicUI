package overlay

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// fadeOut fades a veil away as it goes.
var fadeOut = ui.ElementTransition{Duration: 200 * time.Millisecond, Exit: &ui.Motion{}}

// cover veils content in one cell, disabling it.
func cover(c *ui.Context, root ui.Element, on bool, blur float32, content, veil func()) {
	root.Children(func() {
		ui.Grid(c).Columns(1).FillWidth().Children(func() {
			body := ui.Box(c).Key("content").ColumnStart(1).RowStart(1).Disabled(on)
			body.Blur(body.AnimateTo("blur", blur*unit(on), ui.Smooth)).Children(content)
			if on {
				v := ui.Column(c).Key("veil").ColumnStart(1).RowStart(1).FillWidth().FillHeight().Center().Transition(fadeOut)
				kit.Arrive(v, ui.Pose{Fade: 1}, 0).Children(veil)
			}
		})
	})
}

// card is a centered panel: icon, title, text, then actions.
func card(c *ui.Context, title, text, glyph string, actions func()) {
	t := c.Theme()
	ui.Column(c).Width(t.Space(72)).Padding(t.Space(5)).Gap(t.Space(2)).AlignItems(ui.Center).Radius(t.Radius).
		Background(t.Background).Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).Role(ui.RoleGroup).Label(title).Children(func() {
		ui.Box(c).Size(t.Space(10), t.Space(10)).Center().Radius(999).Background(t.Accent.Alpha(0.12)).Children(func() {
			ui.Icon(c, icon.Must(glyph)).Size(16, 16).TextColor(t.Accent)
		})
		ui.Text(c, title).Bold()
		if text != "" {
			kit.Muted(c, text).TextAlign(ui.Center)
		}
		actions()
	})
}
