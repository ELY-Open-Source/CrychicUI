package schedule

import (
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MeetingCardOptions are the meeting labels.
type MeetingCardOptions struct{ Join, DateFormat string }

// MeetingCard shows a meeting with join and attendees; cue Tap.
func MeetingCard(c *ui.Context, title string, when time.Time, attendees []string, joined *bool, o MeetingCardOptions) ui.Element {
	if title == "" || when.IsZero() || len(attendees) == 0 || joined == nil {
		panic("schedule: MeetingCard has invalid data")
	}
	def(&o.Join, "Join")
	def(&o.DateFormat, "Jan 02 15:04")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			ui.Column(c).Grow(1).Children(func() {
				ui.Text(c, title).Bold()
				kit.Muted(c, when.Format(o.DateFormat)).FontSize(t.Rem(0.75))
			})
			b := kit.Press(ui.ButtonBase(c).Label(o.Join).Padding(t.Space(1), t.Space(2)).Radius(t.Radius).
				Background(t.Accent).TextColor(t.AccentText).Disabled(*joined))
			b.Children(func() { ui.Text(c, o.Join).SingleLine() })
			if b.Clicked() && !*joined {
				*joined = true
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("MeetingCard", "joined")
			}
		})
		ui.Row(c).Gap(t.Space(1)).Children(func() {
			for _, name := range attendees {
				if name == "" {
					panic("schedule: MeetingCard attendee must not be empty")
				}
				ui.Box(c).Size(t.Space(7), t.Space(7)).Radius(999).Background(t.SurfaceHover).Center().Children(func() {
					first, _ := utf8.DecodeRuneInString(name)
					ui.Text(c, string(first)).FontSize(t.Rem(0.65)).TextColor(t.TextMuted)
				})
			}
		})
	})
	return root
}
