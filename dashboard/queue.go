package dashboard

import (
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// QueueOptions word the card.
type QueueOptions struct{ Title, Waiting, Running string }

// Queue shows jobs waiting and running over one split bar.
func Queue(c *ui.Context, waiting, running int, o QueueOptions) ui.Element {
	if waiting < 0 || running < 0 || waiting > math.MaxInt-running {
		panic("dashboard: Queue needs counts of zero or more that sum without overflow")
	}
	def(&o.Title, "Job queue")
	def(&o.Waiting, "Waiting")
	def(&o.Running, "Running")
	total := waiting + running
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "list-ordered", o.Title, nil)
		ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Children(func() {
			row := ui.Row(c).FillWidth().Gap(t.Space(2))
			row.Children(func() {
				for i, stat := range []struct {
					name, glyph string
					value       int
					ink         ui.Color
				}{{o.Waiting, "hourglass", waiting, t.Warning}, {o.Running, "play", running, t.Accent}} {
					c.Key(i)
					box := ui.Row(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2.5)).
						Radius(t.Radius).Background(t.Surface)
					kit.Stagger(c, row, box, i).Children(func() {
						tile(c, stat.glyph, stat.ink)
						ui.Column(c).Grow(1).MinWidth(0).Children(func() {
							kit.Muted(c, stat.name).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
							kit.Numeric(c, number(c, stat.value)).Bold().FontSize(t.Rem(1.3)).FillWidth().SingleLine().
								ContentTransition(ui.ContentTransition{Numeric: true})
						})
					})
				}
			})
			bar := ui.Box(c).FillWidth().Height(t.Space(2))
			share := float32(0)
			if total > 0 {
				share = float32(float64(running) / float64(total))
			}
			w := bar.AnimateTo("running", share, ui.Smooth)
			rest := t.SurfaceHover
			if total > 0 {
				rest = t.Warning.Alpha(0.35)
			}
			bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
				p.Fill(r, rest, 999)
				if w > 0 {
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: max(r.H, r.W*w), H: r.H}, t.Accent, 999)
				}
			}))
		})
	})
	return root
}
