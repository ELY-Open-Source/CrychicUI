package choice

import (
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FilterChips is chips choosing several, with a rolling count.
func FilterChips(c *ui.Context, label string, options []string, selected *[]string) ui.Element {
	distinct(who("FilterChips", label), options)
	subset(who("FilterChips", label), options, selected)
	if label == "" {
		panic("choice: " + who("FilterChips", label) + " needs a label")
	}
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(2))
	return root.Children(func() {
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			ui.Text(c, label).Bold()
			kit.Numeric(c, strconv.Itoa(len(*selected))).TextColor(t.TextMuted)
		})
		ui.Row(c).Wrap().Gap(t.Space(2)).Label(label).FocusGroup(ui.Horizontal).Children(func() {
			for _, option := range options {
				on := slices.Contains(*selected, option)
				c.Key(option)
				chip := ui.ToggleBase(c, &on).Label(option).Gap(t.Space(1)).Padding(t.Space(1), t.Space(3)).Radius(999)
				bg, fg, edge := t.Background, t.Text, t.Border
				if on {
					bg, fg, edge = t.Accent.Alpha(0.14), t.Accent, t.Accent.Alpha(0.4)
				} else if chip.Hovered() {
					bg = t.SurfaceHover
				}
				chip.Background(chip.AnimateColor("bg", bg, quick)).Border(1, chip.AnimateColor("edge", edge, quick)).
					TextColor(chip.AnimateColor("fg", fg, quick))
				shown := on
				kit.Press(chip).Children(func() {
					check := ui.Icon(c, icon.Must("check")).Size(16, 16)
					w := check.AnimateTo("w", 16*unit(shown), ui.Snappy)
					check.Width(max(w, 0)).Opacity(min(max(w/16, 0), 1))
					ui.Text(c, option).SingleLine()
				})
				if kit.Changed(chip, sound.Toggle) {
					*selected = flip(*selected, option)
					root.SetChanged()
					kit.Log("FilterChips", "toggled", "count", len(*selected))
				}
			}
		})
	})
}
