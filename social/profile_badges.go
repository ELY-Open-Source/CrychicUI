package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Badge is one earned profile badge; Detail is optional.
type Badge struct {
	ID, Name, Icon, Detail string
}

// ProfileBadges shows badges earned by a profile.
func ProfileBadges(c *ui.Context, badges []Badge) ui.Element {
	if len(badges) == 0 {
		panic("social: ProfileBadges needs badges")
	}
	ids := make([]string, len(badges))
	for i, badge := range badges {
		if badge.Name == "" || badge.Icon == "" {
			panic(fmt.Sprintf("social: ProfileBadges badge %q needs a name and icon", badge.ID))
		}
		ids[i] = badge.ID
	}
	unique("ProfileBadges", ids)
	t := c.Theme()
	inks := []ui.Color{t.Accent, t.Warning, t.Success, t.Danger}
	root := ui.Grid(c).Columns(3).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, badge := range badges {
			ink := inks[i%len(inks)]
			c.Key(badge.ID)
			tile := ui.Column(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(3), t.Space(2)).Radius(t.Radius).
				Border(1, t.Border).Background(t.Background).MinWidth(0)
			kit.Stagger(c, root, tile, i).Children(func() {
				medal := ui.Box(c).Size(t.Space(11), t.Space(11)).Shrink(0).Center().Radius(999).Background(ink.Alpha(0.14)).TextColor(ink).
					Border(1, ink.Alpha(0.35))
				kit.Arrive(medal, ui.Pose{Scale: 0.5, Rotate: -20, Fade: 1}, 0).Children(func() { glyph(c, badge.Icon, 20) })
				ui.Text(c, badge.Name).FontSize(t.Rem(0.78)).Bold().FillWidth().MinWidth(0).SingleLine().TextAlign(ui.Center)
				if badge.Detail != "" {
					kit.Muted(c, badge.Detail).FontSize(t.Rem(0.68)).FillWidth().MinWidth(0).SingleLine().TextAlign(ui.Center)
				}
			})
		}
	})
	return root
}
