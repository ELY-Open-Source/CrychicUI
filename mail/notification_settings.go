package mail

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// NotificationSetting is one event's channel settings.
type NotificationSetting struct {
	Event       string
	Email, Push bool
}

// NotificationSettingsOptions word the card; Toggle joins event and channel.
type NotificationSettingsOptions struct{ Title, Email, Push, Toggle string }

// NotificationSettings toggles channels by event; cue Toggle.
func NotificationSettings(c *ui.Context, settings *[]NotificationSetting, o NotificationSettingsOptions) ui.Element {
	if settings == nil || len(*settings) == 0 {
		panic("mail: NotificationSettings needs settings")
	}
	def(&o.Title, "Notifications")
	def(&o.Email, "Email")
	def(&o.Push, "Push")
	def(&o.Toggle, "%s %s")
	events := make([]string, len(*settings))
	for i, s := range *settings {
		events[i] = s.Event
	}
	unique("NotificationSettings", events)
	t := c.Theme()
	col := t.Space(14)
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "bell-ring", o.Title, nil)
		ui.Row(c).FillWidth().AlignItems(ui.Center).Padding(t.Space(2), t.Space(3), t.Space(1), t.Space(3)).Children(func() {
			ui.Box(c).Grow(1)
			for _, ch := range []struct{ name, mark string }{{o.Email, "mail"}, {o.Push, "bell"}} {
				ui.Row(c).Width(col).Shrink(0).Gap(t.Space(1)).Justify(ui.Center).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
					glyph(c, ch.mark, 12)
					ui.Text(c, ch.name).FontSize(t.Rem(0.7)).Bold().Shrink(0)
				})
			}
		})
		list := ui.Column(c).FillWidth().Padding(0, 0, t.Space(1.5), 0)
		list.Children(func() {
			for i, setting := range *settings {
				c.Key(setting.Event)
				line := ui.Row(c).FillWidth().AlignItems(ui.Center).Padding(t.Space(2), t.Space(3))
				kit.Stagger(c, list, line, i).Children(func() {
					ui.Text(c, setting.Event).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).Wrap()
					channels := []struct {
						label string
						value *bool
					}{{o.Email, &(*settings)[i].Email}, {o.Push, &(*settings)[i].Push}}
					for j, channel := range channels {
						c.Key(j)
						ui.Row(c).Width(col).Shrink(0).Justify(ui.Center).Children(func() {
							b := toggle(c, fmt.Sprintf(o.Toggle, setting.Event, channel.label), *channel.value)
							if b.Clicked() {
								*channel.value = !*channel.value
								root.SetChanged()
								sound.Play(sound.Toggle)
								kit.Log("NotificationSettings", "toggled", "event", i, "channel", j, "on", *channel.value)
							}
						})
					}
				})
			}
		})
	})
	return root
}
