package datetime

import (
	"fmt"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// TimeSlots picks one of a day's slots; taken ones dim.
func TimeSlots(c *ui.Context, label string, slots, taken []time.Time, selected *time.Time) ui.Element {
	if label == "" || len(slots) == 0 || selected == nil {
		panic("datetime: TimeSlots needs a label, slots and a selection")
	}
	busy := make([]bool, len(slots))
	for i, s := range slots {
		if slices.ContainsFunc(slots[:i], s.Equal) {
			panic(fmt.Sprintf("datetime: TimeSlots slot %d repeats an earlier one", i))
		}
		busy[i] = slices.ContainsFunc(taken, s.Equal)
	}
	for i, s := range taken {
		if !slices.ContainsFunc(slots, s.Equal) {
			panic(fmt.Sprintf("datetime: TimeSlots taken %d is not a slot", i))
		}
	}
	pick := -1
	if !selected.IsZero() {
		if pick = slices.IndexFunc(slots, selected.Equal); pick < 0 || busy[pick] {
			panic("datetime: TimeSlots selection must be a free slot")
		}
	}
	t := c.Theme()
	loc := locale(c)
	shown, changed := pick, false
	root := ui.RadioGroup(c, func() {
		for i, s := range slots {
			name := i18n.FormatDate(loc, s, i18n.DateOptions{Skeleton: "jm"})
			c.Key(i)
			b := ui.RadioBase(c, &pick, i).Label(name).Disabled(busy[i]).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius)
			bg, edge, fg := t.Background, t.Border, t.Text
			switch {
			case i == shown:
				bg, edge, fg = t.Accent, t.Accent, t.AccentText
			case !busy[i] && b.Hovered():
				bg = t.SurfaceHover
			}
			b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Border(1, b.AnimateColor("edge", edge, ui.Snappy)).
				TextColor(b.AnimateColor("fg", fg, ui.Snappy))
			kit.Press(b).Children(func() {
				text := kit.Numeric(c, name)
				if busy[i] {
					text.Strikethrough()
				}
			})
			changed = kit.Changed(b, sound.Select) || changed
		}
	}).Label(label).Row().Wrap().Gap(t.Space(2))
	if changed {
		*selected = slots[pick]
		root.SetChanged()
		kit.Log("TimeSlots", "picked", "index", pick)
	}
	return root
}
