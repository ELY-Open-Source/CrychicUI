package code

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// LangShare is one language's share of a codebase.
type LangShare struct {
	Name  string
	Share float64
}

// Language is a bar of language shares whose segments grow.
func Language(c *ui.Context, label string, shares []LangShare) ui.Element {
	if label == "" || len(shares) == 0 {
		panic(fmt.Sprintf("code: Language %q needs a label and shares", label))
	}
	total := 0.0
	seen := map[string]bool{}
	for _, s := range shares {
		if s.Name == "" || seen[s.Name] || !(s.Share > 0 && s.Share <= 1) {
			panic(fmt.Sprintf("code: Language %q needs unique names and shares in (0,1]", label))
		}
		seen[s.Name] = true
		total += s.Share
	}
	if total > 1.0001 {
		panic(fmt.Sprintf("code: Language %q shares add to %v, over 1", label, total))
	}
	t := c.Theme()
	spare := []ui.Color{t.Accent, t.Success, t.Warning, t.Danger, t.TextMuted}
	inks := make([]ui.Color, len(shares))
	for i, s := range shares {
		ink, known := langInk(t, s.Name)
		if !known {
			ink = spare[i%len(spare)]
		}
		inks[i] = ink
	}
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		kit.CardHead(c, "code", label, nil)
		ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Children(func() {
			bar := ui.Box(c).FillWidth().Height(t.Space(2.5)).Radius(999).Background(t.Surface)
			g := bar.Keyframes("crychic.grow", []ui.Keyframe{{Value: 0}, {At: 600 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			widths := make([]float32, len(shares))
			for i, s := range shares {
				widths[i] = bar.AnimateTo(fmt.Sprint("seg.", i), float32(s.Share), ui.Smooth) * g
			}
			bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
				x := r.X
				for i, w := range widths {
					if seg := r.W*w - 2; seg > 0 {
						p.Fill(ui.Rect{X: x, Y: r.Y, W: seg, H: r.H}, inks[i], 0)
					}
					x += r.W * w
				}
			}))
			ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(2)).Children(func() {
				for i, s := range shares {
					c.Key(s.Name)
					pct := strconv.Itoa(int(s.Share*100+0.5)) + "%"
					item := ui.Row(c).MinWidth(0).Gap(t.Space(2)).AlignItems(ui.Center).Label(s.Name + " " + pct)
					kit.Stagger(c, root, item, i).Children(func() {
						ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Shrink(0).Radius(999).Background(inks[i])
						ui.Text(c, s.Name).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
						kit.Numeric(c, pct).FontSize(t.Rem(0.8)).Bold().Shrink(0)
					})
				}
			})
		})
	})
	return root
}
