package schedule

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MeetingCardOptions are the meeting labels.
type MeetingCardOptions struct{ Join, Joined, DateFormat, People string }

// MeetingCard shows a meeting with join and attendees; cue Tap.
func MeetingCard(c *ui.Context, title string, when time.Time, attendees []string, joined *bool, o MeetingCardOptions) ui.Element {
	if title == "" || when.IsZero() || len(attendees) == 0 || joined == nil {
		panic(fmt.Sprintf("schedule: MeetingCard %q has invalid data", title))
	}
	for _, name := range attendees {
		if name == "" {
			panic(fmt.Sprintf("schedule: MeetingCard %q has an empty attendee", title))
		}
	}
	def(&o.Join, "Join")
	def(&o.Joined, "Joined")
	def(&o.DateFormat, "Mon, Jan 2 · 15:04")
	def(&o.People, "%d people")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(3)).Children(func() {
			tile(c, "video", t.Accent)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, title).Bold().SingleLine().MinWidth(0)
				meta(c, "clock", when.Format(o.DateFormat))
			})
			c.Key("join")
			b := kit.Press(ui.ButtonBase(c).Label(o.Join).Shrink(0).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).Disabled(*joined))
			ink := t.Accent.Mix(t.AccentHover, unit(b.Hovered())).Mix(t.Success, unit(*joined))
			b.Background(b.AnimateColor("bg", ink, ui.Snappy)).TextColor(t.AccentText)
			b.Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					name, text := "video", o.Join
					if *joined {
						name, text = "check", o.Joined
					}
					c.Key(name)
					kit.Arrive(glyph(c, name, 14), ui.Pose{Scale: 0.5, Fade: 1}, 0)
					ui.Text(c, text).Bold().Shrink(0).ContentTransition(ui.ContentTransition{})
				})
			})
			if b.Clicked() && !*joined {
				*joined = true
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("MeetingCard", "joined", "people", len(attendees))
			}
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).
			BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			display.AvatarGroup(c, attendees, display.AvatarGroupOptions{Max: 4})
			kit.Muted(c, fmt.Sprintf(o.People, len(attendees))).FontSize(t.Rem(0.75)).Shrink(0)
		})
	})
	return root
}
