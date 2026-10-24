package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// NotificationSetting is one event's channel settings.
type NotificationSetting struct {
	Event       string
	Email, Push bool
}

// NotificationSettingsOptions are the channel labels.
type NotificationSettingsOptions struct{ Email, Push string }

// NotificationSettings toggles channels by event; cue Toggle.
func NotificationSettings(c *ui.Context, settings *[]NotificationSetting, o NotificationSettingsOptions) ui.Element {
	if settings == nil || len(*settings) == 0 {
		panic("mail: NotificationSettings needs settings")
	}
	def(&o.Email, "Email")
	def(&o.Push, "Push")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, setting := range *settings {
			if setting.Event == "" {
				panic("mail: notification event must not be empty")
			}
			c.Key(setting.Event)
			ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, setting.Event).Grow(1)
				channels := []struct {
					label string
					value *bool
				}{{o.Email, &(*settings)[i].Email}, {o.Push, &(*settings)[i].Push}}
				for j, channel := range channels {
					c.Key(j)
					b := button(c, setting.Event+" "+channel.label).Checked(*channel.value)
					b.Children(func() { ui.Text(c, channel.label) })
					if b.Clicked() {
						*channel.value = !*channel.value
						root.SetChanged()
						sound.Play(sound.Toggle)
						kit.Log("NotificationSettings", "toggled", "event", i, "channel", j, "on", *channel.value)
					}
				}
			})
		}
	})
	return root
}
