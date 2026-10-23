package schedule

import (
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
	def(&o.Label, "Remind me")
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(1)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Text(c, o.Label).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
		for _, option := range options {
			if option == "" {
				panic("schedule: Reminder option must not be empty")
			}
			b := kit.Press(ui.ButtonBase(c).Label(option).Padding(t.Space(1), t.Space(2)).Radius(999).
				Border(1, t.Border).Checked(*choice == option))
			b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*choice == option)), ui.Snappy))
			b.Children(func() { ui.Text(c, option).FontSize(t.Rem(0.75)) })
			if b.Clicked() && *choice != option {
				*choice = option
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Reminder", "selected", "option", *choice)
			}
		}
	})
	return root
}
