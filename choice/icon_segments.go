package choice

import "github.com/egoist/mygo/ui"

// IconSegments is icon-only segments named by their labels.
func IconSegments(c *ui.Context, label string, items []SegmentItem, selected *int) ui.Element {
	for _, item := range items {
		if item.Icon == "" || item.Label == "" {
			panic("choice: IconSegments items need an icon and a label")
		}
	}
	return segmentControl(c, "IconSegments", label, items, selected, true)
}
