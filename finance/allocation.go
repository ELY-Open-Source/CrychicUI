package finance

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// AllocationPart is one asset class share.
type AllocationPart struct {
	Name  string
	Share float64
}

// Allocation shows asset shares as a donut.
func Allocation(c *ui.Context, label string, parts []AllocationPart) ui.Element {
	if label == "" || len(parts) == 0 {
		panic("finance: Allocation needs a label and parts")
	}
	t := c.Theme()
	inks := []ui.Color{t.Accent, t.Success, t.Warning, t.Danger}
	total := 0.0
	for _, part := range parts {
		if part.Name == "" || part.Share <= 0 || part.Share > 1 || !finite(part.Share) {
			panic("finance: Allocation has an invalid share")
		}
		total += part.Share
	}
	if total > 1.0001 {
		panic("finance: Allocation shares exceed one")
	}
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		donut := ui.Box(c).Size(t.Space(24), t.Space(24))
		donut.Draw(func(p *ui.Painter, r ui.Rect) {
			start := 0.0
			for i, part := range parts {
				fraction := part.Share / total
				var path ui.Path
				path.Circle(r.X+r.W/2, r.Y+r.H/2, r.W/2-t.Space(1)).Trim(float32(start), float32(start+fraction))
				p.StrokePath(&path, t.Space(2), inks[i%len(inks)])
				start += fraction
			}
		})
		ui.Column(c).Gap(t.Space(1)).Children(func() {
			for i, part := range parts {
				c.Key(i)
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(inks[i%len(inks)])
					ui.Text(c, part.Name).FontSize(t.Rem(0.8))
					kit.Numeric(c, fmt.Sprintf("%.0f%%", part.Share/total*100)).TextColor(t.TextMuted)
				})
			}
		})
	})
	return root
}
