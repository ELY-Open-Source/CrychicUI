package marketing

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Showcase is one feature told beside its media.
type Showcase struct {
	Eyebrow, Title, Text string
	Points               []string
	Media                *ui.Bitmap
}

// FeatureRowOptions word the section.
type FeatureRowOptions struct{ Title, Text string }

// FeatureRow sets text beside media, sides alternating.
func FeatureRow(c *ui.Context, rows []Showcase, o FeatureRowOptions) ui.Element {
	if len(rows) == 0 {
		panic(fmt.Sprintf("marketing: FeatureRow %q needs rows", o.Title))
	}
	titles := make([]string, len(rows))
	for i, r := range rows {
		if r.Media == nil {
			panic(fmt.Sprintf("marketing: FeatureRow row %q needs media", r.Title))
		}
		unique("FeatureRow", r.Points)
		titles[i] = r.Title
	}
	unique("FeatureRow", titles)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(6)).Role(ui.RoleGroup)
	root.Children(func() {
		heading(c, o.Title, o.Text)
		for i, r := range rows {
			c.Key(strconv.Quote(r.Title))
			line := ui.Row(c).FillWidth().Wrap().Gap(t.Space(5)).AlignItems(ui.Center)
			kit.Stagger(c, root, line, i).Children(func() {
				media := func() {
					frame := ui.Box(c).Grow(1).Basis(t.Space(50)).MinWidth(0).AspectRatio(16.0 / 10).Radius(t.Radius * 1.5).Clip().
						Background(t.SurfaceHover)
					kit.Lift(frame).Children(func() { ui.Image(c, r.Media).Fill().Fit(ui.Cover) })
				}
				if i%2 == 1 {
					media()
				}
				ui.Column(c).Grow(1).Basis(t.Space(50)).MinWidth(0).Gap(t.Space(2)).Children(func() {
					if r.Eyebrow != "" {
						ui.Text(c, r.Eyebrow).FontSize(t.Rem(0.72)).Bold().TextColor(t.Accent).FillWidth().MinWidth(0).SingleLine()
					}
					ui.Text(c, r.Title).FontSize(t.Rem(1.1)).Bold().FillWidth().Wrap()
					kit.Muted(c, r.Text).FontSize(t.Rem(0.84)).FillWidth().Wrap()
					for _, p := range r.Points {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							glyph(c, "check", 14).TextColor(t.Success)
							ui.Text(c, p).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).Wrap()
						})
					}
				})
				if i%2 == 0 {
					media()
				}
			})
		}
	})
	return root
}
