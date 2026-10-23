package media

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CassetteOptions tune a Cassette.
type CassetteOptions struct {
	// Width of the tape; zero gives the default.
	Width float32
}

// Cassette turns its reels while playing.
func Cassette(c *ui.Context, label, title string, playing bool, o CassetteOptions) ui.Element {
	if label == "" {
		panic(fmt.Sprintf("media: Cassette %q needs a label", label))
	}
	t := c.Theme()
	size("Cassette", label, &o.Width, t.Space(64))
	w, h := o.Width, o.Width*0.62
	root := ui.Box(c).Size(w, h).Radius(t.Radius).Background(t.Surface).Border(1, t.Border).Role(ui.RoleImage).Label(label)
	turn := float64(spin(c, root, playing, 0.8))
	root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X + r.W*0.08, Y: r.Y + r.H*0.1, W: r.W * 0.84, H: r.H * 0.52}, t.Background, t.Radius)
		win := ui.Rect{X: r.X + r.W*0.24, Y: r.Y + r.H*0.32, W: r.W * 0.52, H: r.H * 0.24}
		p.Fill(win, kit.Shade(0.75), win.H/2)
		rad := win.H * 0.42
		for _, cx := range []float32{win.X + win.H/2, win.X + win.W - win.H/2} {
			cy := win.Y + win.H/2
			p.Fill(ui.Rect{X: cx - rad, Y: cy - rad, W: 2 * rad, H: 2 * rad}, t.Background, rad)
			for j := range 6 {
				a := 2 * math.Pi * (turn + float64(j)/6)
				dx, dy := float32(math.Cos(a)), float32(math.Sin(a))
				var spoke ui.Path
				spoke.MoveTo(cx+dx*rad*0.3, cy+dy*rad*0.3).LineTo(cx+dx*rad*0.82, cy+dy*rad*0.82)
				p.StrokePath(&spoke, t.Space(0.5), t.TextMuted)
			}
		}
		var base ui.Path
		base.MoveTo(r.X+r.W*0.2, r.Y+r.H).LineTo(r.X+r.W*0.27, r.Y+r.H*0.8).LineTo(r.X+r.W*0.73, r.Y+r.H*0.8).
			LineTo(r.X+r.W*0.8, r.Y+r.H).Close()
		p.FillPath(&base, t.SurfaceHover)
		p.StrokePath(&base, 1, t.Border)
	}))
	return root.Children(func() {
		ui.Row(c).Absolute().Left(0).Right(0).Top(h * 0.14).Justify(ui.Center).Children(func() {
			ui.Text(c, title).SingleLine().Bold().FontSize(t.Rem(0.8))
		})
	})
}
