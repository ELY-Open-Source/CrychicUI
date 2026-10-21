package choice

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Listbox is an inline list choosing one; its highlight slides.
func Listbox(c *ui.Context, label string, options []string, selected *string) ui.Element {
	distinct("Listbox", options)
	if label == "" || selected == nil || !slices.Contains(options, *selected) {
		panic("choice: Listbox needs a label and a valid selection")
	}
	t := c.Theme()
	root := ui.Box(c).Padding(t.Space(1)).Radius(t.Radius+2).Background(t.Background).Border(1, t.Border)
	slide := ui.ElementTransition{Animation: &ui.Snappy, Match: thumbMatch{root.ID()}}
	return root.Children(func() {
		changed := false
		ui.RadioGroup(c, func() {
			for _, option := range options {
				c.Key(option)
				row := ui.RadioBase(c, selected, option).Label(option).Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius)
				on := option == *selected
				fg := t.Text
				if on {
					fg = t.Accent
				} else if row.Hovered() {
					row.Background(t.SurfaceHover)
				}
				row.TextColor(row.AnimateColor("fg", fg, quick)).Children(func() {
					if on {
						ui.Box(c).Key("highlight").Absolute().Top(0).Left(0).Right(0).Bottom(0).Radius(t.Radius).
							Background(t.Accent.Alpha(0.12)).Transition(slide)
					}
					ui.Text(c, option).SingleLine()
				})
				changed = kit.Changed(row, sound.Select) || changed
			}
		}).Label(label).Gap(t.Space(0.5))
		if changed {
			root.SetChanged()
			kit.Log("Listbox", "selected", "index", slices.Index(options, *selected))
		}
	})
}
