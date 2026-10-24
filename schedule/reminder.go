package schedule

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ReminderOptions are the reminder labels.
type ReminderOptions struct{ Label string }

// Reminder selects remind-me options; cue Select.
func Reminder(c *ui.Context, choice *string, options []string, o ReminderOptions) ui.Element {
	if choice == nil || len(options) == 0 {
		panic("schedule: Reminder needs options")
	}
	for i, option := range options {
		if option == "" || slices.Contains(options[:i], option) {
			panic(fmt.Sprintf("schedule: Reminder option %q is empty or repeated", option))
		}
	}
	def(&o.Label, "Remind me")
	t := c.Theme()
	set := slices.Contains(options, *choice)
	root := card(c).Row().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(3))
	root.Children(func() {
		ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(t.Radius).
			Background(root.AnimateColor("tile", t.TextMuted.Alpha(0.1).Mix(t.Warning.Alpha(0.14), unit(set)), ui.Snappy)).
			TextColor(t.TextMuted.Mix(t.Warning, unit(set))).Children(func() {
			name := "bell"
			if set {
				name = "bell-ring"
			}
			c.Key(name)
			kit.Arrive(glyph(c, name, 16), ui.Pose{Rotate: -30, Fade: 1}, 0)
		})
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1.5)).Children(func() {
			ui.Text(c, o.Label).FontSize(t.Rem(0.8)).Bold()
			ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).Children(func() {
				for _, option := range options {
					c.Key(option)
					b := pill(c, option, *choice == option, t.Warning, func() { ui.Text(c, option).FontSize(t.Rem(0.75)).Bold().Shrink(0) })
					if b.Clicked() && *choice != option {
						*choice = option
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("Reminder", "selected", "option", option)
					}
				}
			})
		})
	})
	return root
}
