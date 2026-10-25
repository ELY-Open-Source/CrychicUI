package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// GroupData is one social group; Icon and Faces are optional.
type GroupData struct {
	ID, Name, Privacy, Icon string
	Private                 bool
	Members                 int
	Faces                   []string
}

// GroupCardOptions format the grouped member count.
type GroupCardOptions struct{ Members string }

// GroupCard shows group members and privacy.
func GroupCard(c *ui.Context, group GroupData, o GroupCardOptions) ui.Element {
	if group.ID == "" || group.Name == "" || group.Privacy == "" || group.Members < 0 {
		panic(fmt.Sprintf("social: GroupCard %q needs id, privacy and members", group.Name))
	}
	def(&o.Members, "%s members")
	def(&group.Icon, "users-round")
	t := c.Theme()
	lock, ink := "globe", t.Success
	if group.Private {
		lock, ink = "lock", t.Warning
	}
	root := card(c).Row().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(3))
	root.Background(root.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(root.Hovered())), ui.Snappy))
	root.Children(func() {
		ui.Box(c).Size(t.Space(12), t.Space(12)).Shrink(0).Center().Radius(t.Radius * 1.5).Background(t.Accent.Alpha(0.14)).
			TextColor(t.Accent).Children(func() { glyph(c, group.Icon, 22) })
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
			ui.Text(c, group.Name).FontSize(t.Rem(0.92)).Bold().FillWidth().MinWidth(0).SingleLine()
			ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.25), t.Space(2)).Radius(999).
					Background(ink.Alpha(0.12)).TextColor(ink).Children(func() {
					glyph(c, lock, 12)
					ui.Text(c, group.Privacy).FontSize(t.Rem(0.7)).Bold().Shrink(0)
				})
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
					glyph(c, "users", 12)
					kit.Numeric(c, fmt.Sprintf(o.Members, number(c, group.Members))).FontSize(t.Rem(0.75)).Shrink(0)
				})
			})
		})
		if len(group.Faces) > 0 {
			display.AvatarGroup(c, group.Faces, display.AvatarGroupOptions{}).AlignSelf(ui.Center)
		}
	})
	return root
}
