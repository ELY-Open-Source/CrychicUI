package data

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/egoist/mygo/ui"
)

// Feature is a name and, per plan, Has.
type Feature struct {
	Name string
	Has  []bool
}

// CompareTableOptions tune a CompareTable.
type CompareTableOptions struct {
	// Yes and No name the marks.
	Yes, No string
}

// CompareTable shows features across plans; checks draw in.
func CompareTable(c *ui.Context, label string, plans []string, features []Feature, o CompareTableOptions) ui.Element {
	if label == "" || len(plans) == 0 || len(features) == 0 {
		panic(fmt.Sprintf("data: CompareTable %q needs plans and features", label))
	}
	unique("CompareTable", label, plans)
	for _, f := range features {
		if f.Name == "" || len(f.Has) != len(plans) {
			panic(fmt.Sprintf("data: CompareTable %q features need names and a value per plan", label))
		}
	}
	def(&o.Yes, "Included")
	def(&o.No, "Not included")
	t := c.Theme()
	cols := []Column{{Title: label}}
	for _, p := range plans {
		cols = append(cols, Column{Title: p, Width: t.Space(24)})
	}
	root := frame(c, label)
	root.Children(func() {
		head(c, root, "CompareTable", cols, nil, nil)
		for r, f := range features {
			c.Key(r)
			line(c, hovered(c, false), cols, nil, func(col int) {
				if col == 0 {
					ui.Text(c, f.Name)
					return
				}
				if !f.Has[col-1] {
					ui.Icon(c, icon.Must("minus")).Size(16, 16).TextColor(t.TextMuted).Label(o.No)
					return
				}
				mark := ui.Box(c).Size(16, 16).Role(ui.RoleImage).Label(o.Yes)
				k := mark.Keyframes("draw", []ui.Keyframe{{Value: 0}, {At: 400 * time.Millisecond, Value: 1, Ease: ui.EaseOut}},
					ui.KeyframeOptions{Delay: time.Duration(r*len(plans)+col) * 40 * time.Millisecond})
				mark.Draw(func(p *ui.Painter, b ui.Rect) {
					var tick ui.Path
					tick.MoveTo(b.X+b.W*0.2, b.Y+b.H*0.55).LineTo(b.X+b.W*0.42, b.Y+b.H*0.75).LineTo(b.X+b.W*0.8, b.Y+b.H*0.3)
					p.StrokePath(tick.Trim(0, k), t.Space(0.5), t.Success)
				})
			})
		}
	})
	return root
}
