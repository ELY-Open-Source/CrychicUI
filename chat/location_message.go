package chat

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// LocationMessage is a map tile with pin and address.
func LocationMessage(c *ui.Context, address string) ui.Element {
	if address == "" {
		panic("chat: LocationMessage needs an address")
	}
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Width(t.Space(70)).Clip().Radius(t.Radius*2).Border(1, t.Border).
		Background(t.Surface).Role(ui.RoleGroup).Label(address)
	return root.Children(func() {
		tile := ui.Box(c).FillWidth().AspectRatio(2).Center()
		tile.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, t.SurfaceHover.Alpha(0.5), 0)
			for i := 1; i < 5; i++ {
				x := r.X + r.W*float32(i)/5
				p.Line(x, r.Y, x, r.Y+r.H, 1, t.Border.Alpha(0.6))
			}
			for i := 1; i < 3; i++ {
				y := r.Y + r.H*float32(i)/3
				p.Line(r.X, y, r.X+r.W, y, 1, t.Border.Alpha(0.6))
			}
		}))
		tile.Children(func() {
			pin := ui.Icon(c, icon.Must("map-pin")).Size(16, 16).TextColor(t.Danger)
			kit.Arrive(pin, ui.Pose{Y: -t.Space(3), Fade: 1}, 0)
		})
		ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Children(func() {
			ui.Icon(c, icon.Must("map-pin")).Size(16, 16).TextColor(t.TextMuted)
			ui.Text(c, address).FontSize(t.Rem(0.9)).Wrap()
		})
	})
}
