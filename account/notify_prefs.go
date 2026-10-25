package account

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Topic is one notification topic with its channels.
type Topic struct {
	ID, Name, Detail string
	Email, Push, SMS bool
}

// NotifyPrefsOptions word the card; Label takes channel and topic.
type NotifyPrefsOptions struct{ Title, Email, Push, SMS, Label string }

// NotifyPrefs toggles email, push, and SMS by topic; cue Toggle.
func NotifyPrefs(c *ui.Context, topics *[]Topic, o NotifyPrefsOptions) ui.Element {
	if topics == nil || len(*topics) == 0 {
		panic("account: NotifyPrefs needs topics")
	}
	ids := make([]string, len(*topics))
	for i, tp := range *topics {
		if tp.Name == "" {
			panic(fmt.Sprintf("account: NotifyPrefs %q needs a name", tp.ID))
		}
		ids[i] = tp.ID
	}
	unique("NotifyPrefs", ids)
	def(&o.Title, "Notifications")
	def(&o.Email, "Email")
	def(&o.Push, "Push")
	def(&o.SMS, "SMS")
	def(&o.Label, "%s %s")
	t := c.Theme()
	channels := []struct {
		name, mark string
	}{{o.Email, "mail"}, {o.Push, "bell"}, {o.SMS, "message-square"}}
	cell := t.Space(14)
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		kit.CardHead(c, "bell-ring", o.Title, nil)
		ui.Row(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(2), t.Space(3.5)).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).
			Children(func() {
				ui.Box(c).Grow(1)
				for _, ch := range channels {
					ui.Column(c).Width(cell).Shrink(0).Gap(t.Space(0.5)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
						glyph(c, ch.mark, 14)
						ui.Text(c, ch.name).FontSize(t.Rem(0.7)).Bold().Shrink(0)
					})
				}
			})
		list := ui.Column(c).FillWidth().Padding(t.Space(1.5))
		list.Children(func() {
			for i, tp := range *topics {
				c.Key(tp.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, tp.Name).FontSize(t.Rem(0.85)).Bold().Wrap().FillWidth()
						if tp.Detail != "" {
							kit.Muted(c, tp.Detail).FontSize(t.Rem(0.75)).Wrap().FillWidth()
						}
					})
					flags := []*bool{&(*topics)[i].Email, &(*topics)[i].Push, &(*topics)[i].SMS}
					for k, ch := range channels {
						on := *flags[k]
						b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Label, ch.name, tp.Name)).Width(cell).Shrink(0).
							Padding(t.Space(1), 0).Radius(t.Radius).Checked(on))
						b.Children(func() { tick(c, on, b.Hovered()) })
						if b.Clicked() {
							*flags[k] = !on
							root.SetChanged()
							sound.Play(sound.Toggle)
							kit.Log("NotifyPrefs", "toggled", "topic", i, "channel", k, "on", !on)
						}
					}
				})
			}
		})
	})
	return root
}
