package choice

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SegmentItem is a segment's label and optional Lucide icon.
type SegmentItem struct{ Label, Icon string }

type thumbMatch struct{ id uint64 }

// segmentControl is a segmented control whose thumb slides.
func segmentControl(c *ui.Context, component, label string, items []SegmentItem, selected *int, iconOnly bool) ui.Element {
	if label == "" || len(items) == 0 || selected == nil || *selected < 0 || *selected >= len(items) {
		panic("choice: " + component + " needs a label, items and a valid selection")
	}
	t := c.Theme()
	parts := ui.SegmentedBase(c, selected, len(items))
	root := parts.Track.Label(label).AlignSelf(ui.Start).Gap(t.Space(0.5)).Padding(t.Space(0.5)).Radius(t.Radius).
		Background(t.Surface).Border(1, t.Border)
	id, r := root.ID(), t.Radius-t.Space(0.5)
	root.Children(func() {
		for i, item := range items {
			c.Key(i)
			seg := parts.Segment(i).Label(item.Label).Radius(r)
			if iconOnly {
				seg.Tooltip(item.Label)
			}
			// One cell stacks the thumb under the label.
			seg.Children(func() {
				ui.Grid(c).Columns(1).Children(func() {
					if i == *selected {
						ui.Box(c).Key("thumb").ColumnStart(1).RowStart(1).Radius(r).Background(t.Background).
							Shadow(0, 1, 2, 0, kit.Shade(0.12)).Transition(ui.ElementTransition{Animation: &ui.Snappy, Match: thumbMatch{id}})
					}
					face := ui.Row(c).Key("face").ColumnStart(1).RowStart(1).Gap(t.Space(1.5)).Padding(t.Space(1), t.Space(3)).
						AlignItems(ui.Center)
					if iconOnly {
						face.Padding(t.Space(1.5))
					}
					face.Children(func() {
						if item.Icon != "" {
							ui.Icon(c, icon.Must(item.Icon)).Size(16, 16)
						}
						if !iconOnly {
							ui.Text(c, item.Label).SingleLine()
						}
					})
				})
			})
		}
	})
	if kit.Changed(root, sound.Select) {
		kit.Log(component, "selected", "index", *selected)
	}
	return root
}

// Segments is a segmented control whose thumb slides.
func Segments(c *ui.Context, label string, labels []string, selected *int) ui.Element {
	distinct("Segments", labels)
	items := make([]SegmentItem, len(labels))
	for i, l := range labels {
		items[i] = SegmentItem{Label: l}
	}
	return segmentControl(c, "Segments", label, items, selected, false)
}
