package finance

import (
	"fmt"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// AllocationPart is one asset class share.
type AllocationPart struct {
	Name  string
	Share float64
}

// Allocation shows shares as a donut; hover lifts a class.
func Allocation(c *ui.Context, label string, parts []AllocationPart) ui.Element {
	if label == "" || len(parts) == 0 {
		panic("finance: Allocation needs a label and parts")
	}
	total, top := 0.0, 0
	for i, part := range parts {
		if part.Name == "" || part.Share <= 0 || part.Share > 1 || !finite(part.Share) ||
			slices.IndexFunc(parts, func(p AllocationPart) bool { return p.Name == part.Name }) != i {
			panic(fmt.Sprintf("finance: Allocation %q part %q needs a unique name and a share in (0,1]", label, part.Name))
		}
		total += part.Share
		if part.Share > parts[top].Share {
			top = i
		}
	}
	if total > 1.0001 {
		panic(fmt.Sprintf("finance: Allocation %q shares exceed one", label))
	}
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(label)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "chart-pie", label, nil)
		ui.Row(c).FillWidth().Gap(t.Space(4)).AlignItems(ui.Center).Padding(t.Space(3)).Children(func() {
			hover := -1
			list := ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5))
			list.Children(func() {
				for i, part := range parts {
					c.Key(part.Name)
					row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius)
					if row.Hovered() {
						hover = i
					}
					row.Background(row.AnimateColor("bg", t.SurfaceHover.Alpha(unit(row.Hovered())), ui.Snappy))
					kit.Stagger(c, list, row, i).Children(func() {
						ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Shrink(0).Radius(999).Background(palette(t, i))
						ui.Text(c, part.Name).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
						kit.Numeric(c, fmt.Sprintf("%.0f%%", 100*(part.Share/total))).FontSize(t.Rem(0.85)).Bold().Shrink(0)
					})
				}
			})
			shown := top
			if hover >= 0 {
				shown = hover
			}
			ring := ui.Box(c).Size(t.Space(32), t.Space(32)).Shrink(0).Center()
			grow := ring.Keyframes("crychic.grow", []ui.Keyframe{{Value: 0}, {At: 800 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			widths := make([]float32, len(parts))
			for i := range parts {
				widths[i] = ring.AnimateTo(fmt.Sprint("w", i), t.Space(3)+t.Space(1.5)*unit(i == hover), ui.Snappy)
			}
			ring.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				gap := float32(0.006) * unit(len(parts) > 1)
				start := float32(0)
				for i, part := range parts {
					share := float32(part.Share/total) * grow
					var path ui.Path
					path.Circle(r.X+r.W/2, r.Y+r.H/2, r.W/2-t.Space(2.5))
					p.StrokePath(path.Trim(start+gap, max(start+gap, start+share-gap)), widths[i], palette(t, i))
					start += share
				}
			}))
			ring.Children(func() {
				ui.Column(c).AlignItems(ui.Center).Gap(t.Space(0.25)).Children(func() {
					kit.Numeric(c, fmt.Sprintf("%.0f%%", 100*(parts[shown].Share/total))).Bold().FontSize(t.Rem(1.3)).Shrink(0).
						ContentTransition(ui.ContentTransition{Numeric: true})
					kit.Muted(c, parts[shown].Name).FontSize(t.Rem(0.7)).MaxWidth(t.Space(20)).SingleLine().
						ContentTransition(ui.ContentTransition{})
				})
			})
		})
	})
	return root
}
