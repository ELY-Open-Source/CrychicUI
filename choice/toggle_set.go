package choice

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ToggleSet is joined toggles choosing several.
func ToggleSet(c *ui.Context, label string, items []string, selected *[]string) ui.Element {
	distinct("ToggleSet", items)
	subset("ToggleSet", items, selected)
	if label == "" {
		panic("choice: ToggleSet needs a label")
	}
	t := c.Theme()
	raised := t.Background
	if t.Dark {
		raised = t.SurfacePressed
	}
	changed := false
	root := ui.ToggleGroup(c, func() {
		for _, item := range items {
			on := slices.Contains(*selected, item)
			c.Key(item)
			b := ui.ToggleBase(c, &on).Label(item).Padding(t.Space(1), t.Space(3)).Radius(t.Radius - t.Space(0.5))
			face := raised.Alpha(0)
			if on {
				face = raised
				b.Shadow(0, 1, 2, 0, kit.Shade(0.12))
			}
			b.Background(b.AnimateColor("face", face, quick))
			kit.Press(b).Children(func() { ui.Text(c, item).SingleLine() })
			if kit.Changed(b, sound.Toggle) {
				*selected = flip(*selected, item)
				changed = true
			}
		}
	}).Label(label).AlignSelf(ui.Start)
	if changed {
		root.SetChanged()
		kit.Log("ToggleSet", "toggled", "count", len(*selected))
	}
	return root
}
