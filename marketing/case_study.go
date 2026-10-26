package marketing

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Metric is one headline result.
type Metric struct{ Value, Label string }

// Case is a customer story; nil Cover hides the band.
type Case struct {
	Company, Icon, Quote, Name, Role string
	Photo, Cover                     *ui.Bitmap
	Metrics                          []Metric
}

// CaseStudyOptions word the card; Read labels the action.
type CaseStudyOptions struct{ Read string }

// CaseStudy is a customer story card; cue Tap.
func CaseStudy(c *ui.Context, cs Case, o CaseStudyOptions) (root, read ui.Element) {
	def(&o.Read, "Read the story")
	if strings.TrimSpace(cs.Company) == "" || strings.TrimSpace(cs.Icon) == "" || strings.TrimSpace(cs.Quote) == "" ||
		strings.TrimSpace(cs.Name) == "" || len(cs.Metrics) == 0 || len(cs.Metrics) > 4 {
		panic(fmt.Sprintf("marketing: CaseStudy %q needs an icon, a quote, a name and 1 to 4 metrics", cs.Company))
	}
	for _, m := range cs.Metrics {
		if strings.TrimSpace(m.Value) == "" || strings.TrimSpace(m.Label) == "" {
			panic(fmt.Sprintf("marketing: CaseStudy %q metric needs a value and a label", cs.Company))
		}
	}
	t := c.Theme()
	root = card(c).Clip().Role(ui.RoleGroup)
	root.Children(func() {
		if cs.Cover != nil {
			band := ui.Box(c).FillWidth().AspectRatio(2.4).Clip().Background(t.SurfaceHover)
			band.Children(func() {
				img := ui.Image(c, cs.Cover).Fill().Fit(ui.Cover)
				z := img.Keyframes("zoom", []ui.Keyframe{{Value: 1.08}, {At: 900 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
				img.Scale(z, z)
			})
		}
		c.Key("body")
		ui.Column(c).FillWidth().Gap(t.Space(4)).Padding(t.Space(5)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
				tile(c, cs.Icon, t.Accent, t.Space(9))
				ui.Text(c, cs.Company).FontSize(t.Rem(1)).Bold().Grow(1).MinWidth(0).SingleLine()
			})
			stats := ui.Row(c).FillWidth().Gap(t.Space(3))
			stats.Children(func() {
				for i, m := range cs.Metrics {
					cell := ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(0.5)).Padding(0, 0, 0, t.Space(3)).
						BorderWidth(0, 0, 0, t.Space(0.5)).BorderColor(t.Accent.Alpha(0.5))
					kit.Stagger(c, stats, cell, i).Children(func() {
						ui.Text(c, m.Value).FontSize(t.Rem(1.4)).Bold().TextColor(t.Accent).FillWidth().MinWidth(0).SingleLine()
						kit.Muted(c, m.Label).FontSize(t.Rem(0.74)).FillWidth().Wrap()
					})
				}
			})
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				glyph(c, "quote", 18).TextColor(t.Accent.Alpha(0.6))
				ui.Text(c, cs.Quote).FontSize(t.Rem(0.9)).Grow(1).MinWidth(0).Wrap()
			})
			ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
				ui.Row(c).Grow(1).MinWidth(0).Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
					avatar(c, cs.Name, cs.Photo, 9)
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
						ui.Text(c, cs.Name).FontSize(t.Rem(0.84)).Bold().FillWidth().MinWidth(0).SingleLine()
						if cs.Role != "" {
							kit.Muted(c, cs.Role).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
						}
					})
				})
				c.Key("read")
				read = action(c, o.Read, "arrow-right", false, t.Space(60))
			})
		})
	})
	if read.Clicked() {
		sound.Play(sound.Tap)
		kit.Log("CaseStudy", "read")
	}
	return root, read
}
