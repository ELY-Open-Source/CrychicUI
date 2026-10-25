package mail

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ParticipantsOptions word the card and the speaking mark.
type ParticipantsOptions struct{ Title, Speaking, Muted string }

// Participants lists call participants with mute marks.
func Participants(c *ui.Context, participants []Participant, o ParticipantsOptions) ui.Element {
	def(&o.Title, "In this call")
	def(&o.Speaking, "Speaking")
	def(&o.Muted, "Muted")
	if len(participants) == 0 {
		panic(fmt.Sprintf("mail: Participants %q needs people", o.Title))
	}
	ids := make([]string, len(participants))
	for i, p := range participants {
		if p.Name == "" {
			panic(fmt.Sprintf("mail: Participants person %q needs a name", p.ID))
		}
		ids[i] = p.ID
	}
	unique("Participants", ids)
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "users", o.Title, func() { kit.Chip(c, fmt.Sprint(len(participants)), t.TextMuted) })
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5)).Role(ui.RoleList)
		list.Children(func() {
			for i, p := range participants {
				c.Key(p.ID)
				line := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius)
				line.Background(line.AnimateColor("bg", t.Background.Mix(t.Accent.Alpha(0.08), unit(p.Speaking)), ui.Snappy))
				kit.Stagger(c, list, line, i).Children(func() {
					ring := ui.Box(c).Size(t.Space(9), t.Space(9)).Shrink(0).Center().Radius(999)
					ring.Border(t.Space(0.5), ring.AnimateColor("ring", t.Accent.Alpha(unit(p.Speaking)), ui.Snappy))
					ring.Children(func() { avatar(c, p.Name, 7) })
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
						ui.Text(c, p.Name).FontSize(t.Rem(0.85)).FillWidth().MinWidth(0).SingleLine()
						switch {
						case p.Speaking:
							ui.Text(c, o.Speaking).FontSize(t.Rem(0.7)).TextColor(t.Accent)
						case p.Muted:
							kit.Muted(c, o.Muted).FontSize(t.Rem(0.7))
						}
					})
					switch {
					case p.Muted:
						glyph(c, "mic-off", 16).TextColor(t.Danger)
					case p.Speaking:
						glyph(c, "audio-lines", 16).TextColor(t.Accent)
					default:
						glyph(c, "mic", 16).TextColor(t.TextMuted)
					}
				})
			}
		})
	})
	return root
}
