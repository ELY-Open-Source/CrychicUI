package schedule

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ScheduleSendOptions are the send-later labels.
type ScheduleSendOptions struct{ Label, DateFormat string }

// ScheduleSend selects send-later choices; cue Select.
func ScheduleSend(c *ui.Context, choice *time.Time, options []time.Time, o ScheduleSendOptions) ui.Element {
	if choice == nil || len(options) == 0 {
		panic("schedule: ScheduleSend needs choices")
	}
	def(&o.Label, "Send later")
	def(&o.DateFormat, "Jan 02 15:04")
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(1)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Text(c, o.Label).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
		for _, option := range options {
			if option.IsZero() {
				panic("schedule: ScheduleSend option must not be zero")
			}
			text := option.Format(o.DateFormat)
			b := kit.Press(ui.ButtonBase(c).Label(text).Padding(t.Space(1), t.Space(2)).Radius(999).
				Border(1, t.Border).Checked(choice.Equal(option)))
			b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(choice.Equal(option))), ui.Snappy))
			b.Children(func() { ui.Text(c, text).FontSize(t.Rem(0.75)) })
			if b.Clicked() && !choice.Equal(option) {
				*choice = option
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("ScheduleSend", "selected")
			}
		}
	})
	return root
}
