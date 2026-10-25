package social

import (
	"fmt"

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
	ids := make([]string, len(links))
	for i, link := range links {
		if link.Name == "" || link.Icon == "" || link.URL == "" {
			panic(fmt.Sprintf("social: SocialLinks link %q needs a name, icon and URL", link.ID))
		}
		ids[i] = link.ID
	}
	unique("SocialLinks", ids)
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, link := range links {
			c.Key(link.ID)
			b := kit.Press(ui.ButtonBase(c).Label(link.Name).Role(ui.RoleLink).Size(t.Space(10), t.Space(10)).Shrink(0).Radius(999).
				Tooltip(link.Name).Checked(*opened == link.URL))
			hover := unit(b.Hovered() || *opened == link.URL)
			b.Background(b.AnimateColor("bg", t.Background.Mix(t.Accent.Alpha(0.12), hover), ui.Snappy))
			b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent.Alpha(0.6), hover), ui.Snappy))
			b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Accent, hover), ui.Snappy))
			kit.Stagger(c, root, b, i).Children(func() { glyph(c, link.Icon, 18) })
			if b.Clicked() {
				c.OpenURL(link.URL)
				*opened = link.URL
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("SocialLinks", "opened", "index", i)
			}
		}
	})
	return root
}
