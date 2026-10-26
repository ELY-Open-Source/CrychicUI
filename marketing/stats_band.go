package marketing

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Stat is one big number; Decimals sets its places.
type Stat struct {
	Value                 float64
	Decimals              int
	Prefix, Suffix, Label string
}

// StatsBandOptions word the band.
type StatsBandOptions struct{ Title string }

// StatsBand shows big numbers that count up on arrival.
func StatsBand(c *ui.Context, stats []Stat, o StatsBandOptions) ui.Element {
	if len(stats) == 0 || len(stats) > 6 {
		panic(fmt.Sprintf("marketing: StatsBand %q needs 1 to 6 stats", o.Title))
	}
	labels := make([]string, len(stats))
	for i, s := range stats {
		if !(math.Abs(s.Value) <= 1e15) || s.Decimals < 0 || s.Decimals > 3 {
			panic(fmt.Sprintf("marketing: StatsBand stat %q needs a value within 1e15 and 0 to 3 places", s.Label))
		}
		labels[i] = s.Label
	}
	unique("StatsBand", labels)
	t := c.Theme()
	root := card(c).Gap(t.Space(4)).Padding(t.Space(6), t.Space(4)).Background(t.Background.Mix(t.Accent, 0.05)).Role(ui.RoleGroup)
	root.Children(func() {
		if o.Title != "" {
			ui.Text(c, o.Title).FontSize(t.Rem(1.1)).Bold().FillWidth().TextAlign(ui.Center).Wrap()
		}
		band := ui.Row(c).FillWidth().Wrap().Gap(t.Space(4)).Justify(ui.Center)
		band.Children(func() {
			for i, s := range stats {
				c.Key(strconv.Quote(s.Label))
				cell := ui.Column(c).Grow(1).Basis(t.Space(28)).MinWidth(0).Gap(t.Space(1)).AlignItems(ui.Center)
				delay := time.Duration(i) * 120 * time.Millisecond
				frames := []ui.Keyframe{{Value: 0}, {At: delay + 1200*time.Millisecond, Value: 1, Ease: ui.EaseOut}}
				if delay > 0 {
					frames = []ui.Keyframe{{Value: 0}, {At: delay, Value: 0}, {At: delay + 1200*time.Millisecond, Value: 1, Ease: ui.EaseOut}}
				}
				k := cell.Keyframes("count", frames)
				shown := s.Value
				if k < 1 {
					shown = s.Value * float64(k)
				}
				kit.Stagger(c, band, cell, i).Children(func() {
					kit.Numeric(c, s.Prefix+decimal(c, shown, s.Decimals)+s.Suffix).FontSize(t.Rem(1.9)).Bold().TextColor(t.Accent).
						FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine()
					kit.Muted(c, s.Label).FontSize(t.Rem(0.8)).FillWidth().TextAlign(ui.Center).Wrap()
				})
			}
		})
	})
	return root
}
