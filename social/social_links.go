package social

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SocialLink is one external profile link.
type SocialLink struct {
	ID, Name, Icon, URL string
}

// SocialLinks opens icon links; cue Tap.
func SocialLinks(c *ui.Context, links []SocialLink, opened *string) ui.Element {
	if len(links) == 0 || opened == nil {
		panic("social: SocialLinks needs links and state")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, link := range links {
			if link.ID == "" || link.Name == "" || link.Icon == "" || link.URL == "" {
				panic("social: invalid social link")
			}
			c.Key(link.ID)
			b := button(c, link.Name).Role(ui.RoleLink).Radius(999).Tooltip(link.Name)
			b.Children(func() { ui.Icon(c, icon.Must(link.Icon)).Size(16, 16).TextColor(t.Accent) })
			if b.Clicked() {
				c.OpenURL(link.URL)
				*opened = link.URL
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("SocialLinks", "opened")
			}
		}
	})
	return root
}
