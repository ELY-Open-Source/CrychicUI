package dashboard

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DeviceShare is one device share; Icon overrides the default.
type DeviceShare struct {
	Kind, Icon string
	Share      float64
}

// DevicesOptions word the card.
type DevicesOptions struct{ Title string }

// Devices shows device shares as tiles over one split bar.
func Devices(c *ui.Context, shares []DeviceShare, o DevicesOptions) ui.Element {
	if len(shares) == 0 || len(shares) > 4 {
		panic("dashboard: Devices needs 1–4 shares")
	}
	names, total := make([]string, len(shares)), 0.0
	for i, share := range shares {
		if share.Share < 0 || share.Share > 1 || !finite(share.Share) {
			panic(fmt.Sprintf("dashboard: Devices %q needs a share in [0,1]", share.Kind))
		}
		names[i] = share.Kind
		total += share.Share
	}
	unique("Devices", names)
	if total > 1.0001 {
		panic("dashboard: Devices shares exceed one")
	}
	def(&o.Title, "Devices")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "monitor-smartphone", o.Title, nil)
		ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Children(func() {
			row := ui.Row(c).FillWidth().Gap(t.Space(2))
			row.Children(func() {
				for i, share := range shares {
					c.Key(share.Kind)
					glyph := share.Icon
					if glyph == "" {
						glyph = map[string]string{"desktop": "monitor", "mobile": "smartphone", "tablet": "tablet"}[strings.ToLower(share.Kind)]
					}
					if glyph == "" {
						glyph = "monitor-smartphone"
					}
					box := ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(3), t.Space(1)).
						Radius(t.Radius)
					box.Background(box.AnimateColor("bg", t.Surface.Mix(palette(t, i).Alpha(0.1), unit(box.Hovered())), ui.Snappy))
					kit.Stagger(c, row, box, i).Children(func() {
						tile(c, glyph, palette(t, i))
						kit.Numeric(c, fmt.Sprintf("%.0f%%", 100*share.Share)).Bold().FontSize(t.Rem(1.2)).Shrink(0).
							ContentTransition(ui.ContentTransition{Numeric: true})
						kit.Muted(c, share.Kind).FontSize(t.Rem(0.75)).MaxWidth(t.Space(30)).SingleLine()
					})
				}
			})
			bar := ui.Box(c).FillWidth().Height(t.Space(2.5))
			ws := make([]float32, len(shares))
			for i, share := range shares {
				ws[i] = bar.AnimateTo(fmt.Sprint("s", i), float32(share.Share), ui.Smooth)
			}
			bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
				p.Fill(r, t.SurfaceHover, 999)
				x := r.X
				for i, s := range ws {
					if w := r.W * s; w > 0 {
						p.Fill(ui.Rect{X: x, Y: r.Y, W: max(0, w-t.Space(0.5)), H: r.H}, palette(t, i), 999)
						x += w
					}
				}
			}))
		})
	})
	return root
}
