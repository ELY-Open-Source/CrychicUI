package mail

import (
	"fmt"
	"net/mail"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RecipientsOptions word the groups; Remove formats each chip's label.
type RecipientsOptions struct{ To, Cc, Bcc, Remove string }

// Recipients edits To, Cc, and Bcc chips; cue Pop.
func Recipients(c *ui.Context, to, cc, bcc *[]string, o RecipientsOptions) ui.Element {
	if to == nil || cc == nil || bcc == nil {
		panic("mail: Recipients needs all groups")
	}
	def(&o.To, "To")
	def(&o.Cc, "Cc")
	def(&o.Bcc, "Bcc")
	def(&o.Remove, "Remove %s")
	groups := []struct {
		label  string
		people *[]string
	}{{o.To, to}, {o.Cc, cc}, {o.Bcc, bcc}}
	for _, group := range groups {
		for _, address := range *group.people {
			if _, err := mail.ParseAddress(address); err != nil {
				panic(fmt.Sprintf("mail: Recipients %s has an invalid address %q", group.label, address))
			}
		}
		unique("Recipients", *group.people)
	}
	t := c.Theme()
	root := card(c)
	removeGroup, removeAt := -1, -1
	root.Children(func() {
		for i, group := range groups {
			c.Key(i)
			line := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3))
			if i < len(groups)-1 {
				line.BorderWidth(0, 0, 1, 0).BorderColor(t.Border)
			}
			line.Children(func() {
				kit.Muted(c, group.label).FontSize(t.Rem(0.8)).Width(t.Space(10)).Shrink(0)
				chips := ui.Row(c).Grow(1).MinWidth(0).Wrap().Gap(t.Space(1.5))
				room := chips.LayoutBounds().W
				kit.Settle(c, chips, []float32{room}, func() []float32 { return []float32{chips.LayoutBounds().W} })
				// Avatar, cross, padding and gaps around the address.
				text := min(t.Space(52), max(0, room-t.Space(10)-12))
				chips.Children(func() {
					for j, address := range *group.people {
						c.Key(address)
						chip := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Remove, address)).Padding(t.Space(0.5), t.Space(1.5), t.Space(0.5), t.Space(0.5)).
							Radius(999).Shrink(0))
						chip.Background(chip.AnimateColor("bg", t.SurfaceHover.Mix(t.Danger.Alpha(0.14), unit(chip.Hovered())), ui.Snappy))
						chip.TextColor(chip.AnimateColor("ink", t.Text.Mix(t.Danger, unit(chip.Hovered())), ui.Snappy))
						kit.Arrive(chip, ui.Pose{Scale: 0.8, Fade: 1}, 0).Children(func() {
							ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
								avatar(c, address, 5)
								ui.Text(c, address).FontSize(t.Rem(0.78)).Shrink(0).MaxWidth(text).SingleLine()
								glyph(c, "x", 12)
							})
						})
						if chip.Clicked() {
							removeGroup, removeAt = i, j
						}
					}
				})
			})
		}
	})
	if removeAt >= 0 {
		people := groups[removeGroup].people
		*people = slices.Delete(*people, removeAt, removeAt+1)
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("Recipients", "removed", "group", removeGroup, "count", len(*people))
	}
	return root
}
