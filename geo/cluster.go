package geo

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ClusterOptions word the cluster; Label formats the count.
type ClusterOptions struct{ Label, More string }

// Cluster is grouped pins with their count.
func Cluster(c *ui.Context, count int, o ClusterOptions) ui.Element {
	def(&o.Label, "%d places")
	def(&o.More, "999+")
	if count < 2 {
		panic(fmt.Sprintf("geo: Cluster %q needs at least 2 places, got %d", o.Label, count))
	}
	t := c.Theme()
	size := float32(10)
	if count >= 100 {
		size = 14
	} else if count >= 10 {
		size = 12
	}
	text := strconv.Itoa(count)
	if count > 999 {
		text = o.More
	}
	root := ui.Box(c).Size(t.Space(size+6), t.Space(size+6)).Shrink(0).Center().Role(ui.RoleGroup).Label(fmt.Sprintf(o.Label, count))
	kit.Arrive(root, ui.Pose{Scale: 0.4, Fade: 1}, 0)
	root.Draw(func(p *ui.Painter, r ui.Rect) {
		var ring ui.Path
		p.FillPath(ring.Circle(r.X+r.W/2, r.Y+r.H/2, r.W/2), t.Accent.Alpha(0.18))
	})
	return root.Children(func() {
		ui.Box(c).Size(t.Space(size), t.Space(size)).Center().Radius(999).Background(t.Accent).Border(2, t.Background).Children(func() {
			kit.Numeric(c, text).FontSize(t.Rem(0.62 + size*0.012)).Bold().TextColor(t.AccentText).Shrink(0).MaxWidth(t.Space(size - 2)).SingleLine()
		})
	})
}
