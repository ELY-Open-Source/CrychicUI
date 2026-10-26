package marketing

import (
	"fmt"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// WaitlistOptions word the signup; Count is how many wait already.
type WaitlistOptions struct {
	Title, Text, Label, Join, Waiting, Place, Joined string
	Count                                            int
	Recent                                           []string
}

// Waitlist joins an email, showing its place; cue Success.
func Waitlist(c *ui.Context, email *string, place *int, o WaitlistOptions) (root, join ui.Element) {
	def(&o.Title, "Join the waitlist")
	if email == nil || place == nil || o.Count < 0 || o.Count >= 1e12 || *place < 0 || *place > o.Count+1 {
		panic(fmt.Sprintf("marketing: Waitlist %q needs state, a count of 0 to 1e12 and a place of 0 to count+1", o.Title))
	}
	unique("Waitlist", o.Recent)
	def(&o.Label, "Email address")
	def(&o.Join, "Join")
	def(&o.Waiting, "%s people waiting")
	def(&o.Place, "#%s")
	def(&o.Joined, "You're in line")
	t := c.Theme()
	root = card(c).Gap(t.Space(4)).Padding(t.Space(6), t.Space(5)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		tile(c, "hourglass", t.Accent, t.Space(11))
		ui.Text(c, o.Title).FontSize(t.Rem(1.2)).Bold().FillWidth().TextAlign(ui.Center).Wrap()
		if o.Text != "" {
			kit.Muted(c, o.Text).FontSize(t.Rem(0.84)).FillWidth().MaxWidth(t.Space(100)).TextAlign(ui.Center).Wrap()
		}
		c.Key(*place > 0)
		ui.Column(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			if *place > 0 {
				box := ui.Column(c).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(3), t.Space(6)).Radius(t.Radius * 1.5).
					Background(t.Success.Alpha(0.1))
				kit.Arrive(box, ui.Pose{Scale: 0.9, Fade: 1}, 0).Children(func() {
					ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(t.Success).Children(func() {
						glyph(c, "party-popper", 16)
						ui.Text(c, o.Joined).FontSize(t.Rem(0.82)).Bold().Shrink(0)
					})
					kit.Numeric(c, fmt.Sprintf(o.Place, number(c, *place))).FontSize(t.Rem(2)).Bold().Shrink(0).
						ContentTransition(ui.ContentTransition{Numeric: true})
				})
				return
			}
			ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				in := emailField(c, o.Label, email)
				if in.Changed() {
					root.SetChanged()
					kit.Log("Waitlist", "edited", "runes", utf8.RuneCountInString(*email))
				}
				ok := validEmail(*email)
				c.Key("join")
				join = action(c, o.Join, "arrow-right", true, t.Space(60)).Disabled(!ok)
				if (join.Clicked() || in.Submitted()) && ok {
					*place = o.Count + 1
					root.SetChanged()
					sound.Play(sound.Success)
					kit.Log("Waitlist", "joined", "place", *place)
				}
			})
		})
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			if len(o.Recent) > 0 {
				display.AvatarGroup(c, o.Recent, display.AvatarGroupOptions{})
			}
			waiting := o.Count
			if *place > 0 {
				waiting++
			}
			kit.Muted(c, fmt.Sprintf(o.Waiting, number(c, waiting))).FontSize(t.Rem(0.78)).Shrink(0).
				ContentTransition(ui.ContentTransition{Numeric: true})
		})
	})
	return root, join
}
