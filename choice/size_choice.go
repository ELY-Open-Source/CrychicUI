package choice

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SizeChoice is size boxes; unavailable ones are struck through.
func SizeChoice(c *ui.Context, label string, sizes, unavailable []string, selected *string) ui.Element {
	distinct("SizeChoice", sizes)
	for _, u := range unavailable {
		if !slices.Contains(sizes, u) {
			panic("choice: SizeChoice unavailable size is unknown")
		}
	}
	if label == "" || selected == nil || *selected != "" && (!slices.Contains(sizes, *selected) || slices.Contains(unavailable, *selected)) {
		panic("choice: SizeChoice needs a label and an empty or available selection")
	}
	t := c.Theme()
	changed := false
	root := ui.RadioGroup(c, func() {
		for _, s := range sizes {
			out := slices.Contains(unavailable, s)
			c.Key(s)
			b := ui.RadioBase(c, selected, s).Label(s).Disabled(out).MinWidth(t.Space(11)).Height(t.Space(9)).
				Padding(0, t.Space(2)).Center().Radius(t.Radius)
			bg, edge, fg := t.Background, t.Border, t.Text
			switch {
			case *selected == s:
				bg, edge, fg = t.Accent.Alpha(0.1), t.Accent, t.Accent
			case !out && b.Hovered():
				edge = t.Border.Mix(t.Text, 0.3)
			}
			b.Background(b.AnimateColor("bg", bg, quick)).Border(1, b.AnimateColor("edge", edge, quick)).
				TextColor(b.AnimateColor("fg", fg, quick))
			if out {
				b.Clip().DrawOver(func(p *ui.Painter, r ui.Rect) { p.Line(r.X, r.Y+r.H, r.X+r.W, r.Y, 1, t.Border) })
			}
			kit.Press(b).Children(func() { ui.Text(c, s).Bold() })
			changed = kit.Changed(b, sound.Select) || changed
		}
	}).Label(label).Row().Wrap().Gap(t.Space(2))
	if changed {
		root.SetChanged()
		kit.Log("SizeChoice", "selected", "index", slices.Index(sizes, *selected))
	}
	return root
}
