package layout

import (
	"slices"

	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Masonry puts each item in the shortest of n columns.
func Masonry(c *ui.Context, columns int, items []func()) ui.Element {
	if columns < 1 {
		panic("layout: Masonry needs a column")
	}
	gap := c.Theme().Space(3)
	root := ui.Box(c).FillWidth()
	w := root.LayoutBounds().W
	colW := max(0, (w-gap*float32(columns-1))/float32(columns))
	used, boxes, sum := make([]float32, len(items)), make([]ui.Element, len(items)), make([]float32, columns)
	root.Children(func() {
		for i, item := range items {
			b := ui.Box(c).Key(i).Absolute().Width(colW).Transition(ui.ElementTransition{Position: true})
			used[i] = b.LayoutBounds().H
			j := slices.Index(sum, slices.Min(sum))
			b.Left(float32(j) * (colW + gap)).Top(sum[j]).Children(item)
			sum[j] += used[i] + gap
			boxes[i] = b
		}
	})
	if len(items) > 0 {
		root.Height(slices.Max(sum) - gap)
	}
	kit.Settle(c, root, append(used, w), func() []float32 {
		h := make([]float32, 0, len(boxes)+1)
		for _, b := range boxes {
			h = append(h, b.LayoutBounds().H)
		}
		return append(h, root.LayoutBounds().W)
	})
	return root
}
