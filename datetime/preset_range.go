package datetime

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Preset is the last Days days, today included.
type Preset struct {
	Label string
	Days  int
}

// PresetRangeOptions tune a PresetRange.
type PresetRangeOptions struct {
	// Presets default to today, 7 days and 30 days.
	Presets []Preset
	// Custom, From and To label the custom choice.
	Custom, From, To string
}

// PresetRange picks a preset range, or custom dates.
func PresetRange(c *ui.Context, label string, start, end *time.Time, o PresetRangeOptions) ui.Element {
	if label == "" || start == nil || end == nil {
		panic("datetime: PresetRange needs a label and dates")
	}
	if o.Presets == nil {
		o.Presets = []Preset{{"Today", 1}, {"7 days", 7}, {"30 days", 30}}
	}
	labels := make([]string, 0, len(o.Presets)+1)
	for _, p := range o.Presets {
		if p.Label == "" || p.Days < 1 {
			panic(fmt.Sprintf("datetime: PresetRange preset %q needs a label and a day or more", p.Label))
		}
		labels = append(labels, p.Label)
	}
	def(&o.Custom, "Custom")
	def(&o.From, "From")
	def(&o.To, "To")
	labels = append(labels, o.Custom)
	t := c.Theme()
	today := c.Now()
	root := ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(3))
	custom := ui.Local(root, "custom", func() bool { return false })
	pick := len(o.Presets)
	for i, p := range o.Presets {
		if !*custom && sameDay(*end, today) && sameDay(*start, today.AddDate(0, 0, 1-p.Days)) {
			pick = i
		}
	}
	return root.Children(func() {
		if choice.Segments(c, label, labels, &pick).Changed() {
			if *custom = pick == len(o.Presets); !*custom {
				*start, *end = onDay(today.AddDate(0, 0, 1-o.Presets[pick].Days), *start), onDay(today, *end)
			}
			root.SetChanged()
			kit.Log("PresetRange", "picked", "index", pick)
		}
		if pick == len(o.Presets) {
			ui.Row(c).Key("custom").Gap(t.Space(2)).Transition(ui.ElementTransition{Enter: &ui.Motion{Y: -t.Space(2)},
				Exit: &ui.Motion{Collapse: true}}).Children(func() {
				from := ui.DateInput(c, start).Label(o.From)
				to := ui.DateInput(c, end).Label(o.To)
				if kit.Changed(from, sound.Select) || kit.Changed(to, sound.Select) {
					if dayBefore(*end, *start) {
						*start, *end = onDay(*end, *start), onDay(*start, *end)
					}
					root.SetChanged()
					kit.Log("PresetRange", "custom")
				}
			})
		}
	})
}
