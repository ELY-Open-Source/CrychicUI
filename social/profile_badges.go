package social

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Badge is one earned profile badge.
type Badge struct {
	ID, Name, Icon string
}

// ProfileBadges shows badges earned by a profile.
func ProfileBadges(c *ui.Context, badges []Badge) ui.Element {
	if len(badges) == 0 {
		panic("social: ProfileBadges needs badges")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, badge := range badges {
			if badge.ID == "" || badge.Name == "" || badge.Icon == "" {
				panic("social: invalid badge")
			}
			c.Key(badge.ID)
			chip := ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2)).
				Radius(999).Background(t.Accent.Alpha(0.12))
			kit.Arrive(chip, ui.Pose{Scale: 0.9, Fade: 1}, 0)
			chip.Children(func() {
				ui.Icon(c, icon.Must(badge.Icon)).Size(16, 16).TextColor(t.Accent)
				ui.Text(c, badge.Name)
			})
		}
	})
	return root
}
