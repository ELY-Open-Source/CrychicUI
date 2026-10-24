package social

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SocialNotification is one grouped social notification.
type SocialNotification struct {
	ID, Kind, Time string
	Actors         []string
}

// NotificationsFeedOptions are the notification group labels.
type NotificationsFeedOptions struct{ Likes, Follows string }

// NotificationsFeed groups likes and follows.
func NotificationsFeed(c *ui.Context, notifications []SocialNotification, o NotificationsFeedOptions) ui.Element {
	if len(notifications) == 0 {
		panic("social: NotificationsFeed needs notifications")
	}
	def(&o.Likes, "Likes")
	def(&o.Follows, "New followers")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		groups := []struct {
			kind, title string
		}{{"like", o.Likes}, {"follow", o.Follows}}
		for _, group := range groups {
			c.Key(group.kind)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Text(c, group.title).Bold()
				for _, notice := range notifications {
					if notice.ID == "" || notice.Time == "" || len(notice.Actors) == 0 {
						panic("social: invalid notification")
					}
					if notice.Kind != "like" && notice.Kind != "follow" {
						panic("social: unknown notification kind")
					}
					if notice.Kind != group.kind {
						continue
					}
					for _, actor := range notice.Actors {
						if actor == "" {
							panic("social: notification actor must not be empty")
						}
					}
					c.Key(notice.ID)
					ui.Row(c).FillWidth().Children(func() {
						ui.Text(c, strings.Join(notice.Actors, ", ")).Grow(1).Wrap()
						kit.Muted(c, notice.Time)
					})
				}
			})
		}
	})
	return root
}
