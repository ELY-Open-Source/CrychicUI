package schedule

import (
	"fmt"
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
	for i, option := range options {
		if option.IsZero() {
			panic("schedule: ScheduleSend option must not be zero")
		}
		for _, prior := range options[:i] {
			if prior.Equal(option) {
				panic(fmt.Sprintf("schedule: ScheduleSend option %v repeats", option))
			}
		}
	}
	def(&o.Label, "Schedule send")
	def(&o.DateFormat, "Jan 02 15:04")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "send", o.Label, nil)
		list := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2))
		list.Children(func() {
			for i, option := range options {
				text, on := option.Format(o.DateFormat), choice.Equal(option)
				name, ink := "moon", t.Accent
				switch h := option.Hour(); {
				case h < 12:
					name, ink = "sunrise", t.Warning
				case h < 18:
					name, ink = "sun", t.Warning
				}
				c.Key(instant(option))
				b := kit.Press(ui.ButtonBase(c).Label(text).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(1.5), t.Space(2)).
					Radius(t.Radius).Checked(on))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(0.7*unit(b.Hovered())).Mix(t.Accent.Alpha(0.12), unit(on)), ui.Snappy))
				kit.Stagger(c, list, b, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
						tile(c, name, ink)
						ui.Text(c, text).FontSize(t.Rem(0.85)).Bold().Grow(1).MinWidth(0).SingleLine()
						if on {
							c.Key("check")
							kit.Arrive(glyph(c, "circle-check", 16).TextColor(t.Accent), ui.Pose{Scale: 0.4, Fade: 1}, 0)
						}
					})
				})
				if b.Clicked() && !on {
					*choice = option
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("ScheduleSend", "selected", "index", i)
				}
			}
		})
	})
	return root
}
