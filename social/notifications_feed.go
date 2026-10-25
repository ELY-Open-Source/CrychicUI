package social

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SocialNotification is one grouped social notification.
type SocialNotification struct {
	ID, Kind, Time string
	Actors         []string
}

// NotificationsFeedOptions word the groups; Like and Follow take the actors.
type NotificationsFeedOptions struct{ Likes, Follows, Like, Follow, Two, Many string }

// NotificationsFeed groups likes and follows.
func NotificationsFeed(c *ui.Context, notifications []SocialNotification, o NotificationsFeedOptions) ui.Element {
	if len(notifications) == 0 {
		panic("social: NotificationsFeed needs notifications")
	}
	ids := make([]string, len(notifications))
	for i, notice := range notifications {
		if notice.Time == "" || len(notice.Actors) == 0 || slices.Contains(notice.Actors, "") {
			panic(fmt.Sprintf("social: NotificationsFeed notice %q needs a time and actors", notice.ID))
		}
		if notice.Kind != "like" && notice.Kind != "follow" {
			panic(fmt.Sprintf("social: NotificationsFeed notice %q has unknown kind %q", notice.ID, notice.Kind))
		}
		ids[i] = notice.ID
	}
	unique("NotificationsFeed", ids)
	def(&o.Likes, "Likes")
	def(&o.Follows, "New followers")
	def(&o.Like, "%s liked your post")
	def(&o.Follow, "%s followed you")
	def(&o.Two, "%s and %s")
	def(&o.Many, "%s and %d others")
	t := c.Theme()
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		groups := []struct {
			kind, title, format, mark string
			ink                       ui.Color
		}{{"like", o.Likes, o.Like, "heart", t.Danger}, {"follow", o.Follows, o.Follow, "user-plus", t.Accent}}
		n := 0
		for g, group := range groups {
			c.Key(group.kind)
			section := ui.Column(c).FillWidth().Padding(t.Space(1.5))
			if g > 0 {
				section.BorderWidth(1, 0, 0, 0).BorderColor(t.Border)
			}
			section.Children(func() {
				header := false
				for _, notice := range notifications {
					if notice.Kind != group.kind {
						continue
					}
					if !header {
						kit.Muted(c, group.title).FontSize(t.Rem(0.72)).Bold().Padding(t.Space(1), t.Space(2))
						header = true
					}
					c.Key(notice.ID)
					row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
					row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
					kit.Stagger(c, root, row, n).Children(func() {
						ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999).Background(group.ink.Alpha(0.14)).
							TextColor(group.ink).Children(func() { glyph(c, group.mark, 16) })
						ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
							display.AvatarGroup(c, notice.Actors, display.AvatarGroupOptions{})
							ui.RichText(c, phrase(group.format, actors(notice.Actors, o))...).FontSize(t.Rem(0.82)).Wrap().FillWidth()
						})
						kit.Muted(c, notice.Time).FontSize(t.Rem(0.72)).Shrink(0).AlignSelf(ui.Start)
					})
					n++
				}
			})
		}
	})
	return root
}

// actors names one or two, else the first and others.
func actors(names []string, o NotificationsFeedOptions) string {
	switch len(names) {
	case 1:
		return names[0]
	case 2:
		return fmt.Sprintf(o.Two, names[0], names[1])
	}
	return fmt.Sprintf(o.Many, names[0], len(names)-1)
}

// phrase fills format's %s with who, set in bold.
func phrase(format, who string) []ui.Span {
	if strings.Count(format, "%s") != 1 {
		panic(fmt.Sprintf("social: NotificationsFeed format %q needs one %%s", format))
	}
	before, after, _ := strings.Cut(format, "%s")
	return []ui.Span{{Text: before}, {Text: who, Weight: 700}, {Text: after}}
}
