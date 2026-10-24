package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Share is one named fraction; Icon is a Lucide name.
type Share struct {
	Name, Icon string
	Share      float64
}

// TrafficOptions word the card.
type TrafficOptions struct{ Title string }

// Traffic shows sources as one split bar and ranked rows.
func Traffic(c *ui.Context, sources []Share, o TrafficOptions) ui.Element {
	if len(sources) == 0 {
		panic("dashboard: Traffic needs sources")
	}
	names, total := make([]string, len(sources)), 0.0
	for i, source := range sources {
		if source.Share < 0 || source.Share > 1 || !finite(source.Share) {
			panic(fmt.Sprintf("dashboard: Traffic source %q needs a share in [0,1]", source.Name))
		}
		names[i] = source.Name
		total += source.Share
	}
	unique("Traffic", names)
	if total > 1.0001 {
		panic("dashboard: Traffic shares exceed one")
	}
	def(&o.Title, "Traffic sources")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "globe", o.Title, nil)
		ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Children(func() {
			bar := ui.Box(c).FillWidth().Height(t.Space(2.5))
			shares := make([]float32, len(sources))
			for i, source := range sources {
				shares[i] = bar.AnimateTo(fmt.Sprint("s", i), float32(source.Share), ui.Smooth)
			}
			bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
				p.Fill(r, t.SurfaceHover, 999)
				x := r.X
				for i, s := range shares {
					if w := r.W * s; w > 0 {
						p.Fill(ui.Rect{X: x, Y: r.Y, W: max(0, w-t.Space(0.5)), H: r.H}, palette(t, i), 999)
						x += w
					}
				}
			}))
			list := ui.Column(c).FillWidth().Gap(t.Space(1))
			list.Children(func() {
				for i, source := range sources {
					c.Key(source.Name)
					row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius)
					row.Background(row.AnimateColor("bg", t.SurfaceHover.Alpha(unit(row.Hovered())), ui.Snappy))
					kit.Stagger(c, list, row, i).Children(func() {
						if source.Icon != "" {
							ui.Icon(c, icon.Must(source.Icon)).Size(16, 16).Shrink(0).TextColor(palette(t, i))
						} else {
							ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Shrink(0).Radius(999).Background(palette(t, i))
						}
						ui.Text(c, source.Name).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
						kit.Numeric(c, fmt.Sprintf("%.0f%%", 100*source.Share)).Bold().FontSize(t.Rem(0.85)).Shrink(0)
					})
				}
			})
		})
	})
	return root
}
