// Code components for the gallery.
package code

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Breadcrumb is a file and symbol path; cue Tap.
func Breadcrumb(c *ui.Context, parts []string, at *int) ui.Element {
	if len(parts) == 0 || at == nil || *at < 0 || *at >= len(parts) {
		panic("code: Breadcrumb needs parts and a position")
	}
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(parts[*at])
	root.Children(func() {
		for i, part := range parts {
			if part == "" {
				panic("code: Breadcrumb needs text in every part")
			}
			c.Key(i)
			if i > 0 {
				ui.Icon(c, iconChevron).Size(14, 14).TextColor(t.TextMuted)
			}
			b := kit.Press(ui.ButtonBase(c).Label(part).Padding(t.Space(0.5), t.Space(1)).Radius(t.Radius))
			b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy))
			b.Children(func() {
				e := ui.Text(c, part).SingleLine()
				if i == *at {
					e.Bold()
				} else {
					e.TextColor(t.TextMuted)
				}
			})
			if b.Clicked() && *at != i {
				*at = i
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("Breadcrumb", "to", "index", i)
			}
		}
	})
	return root
}

// Commit is a commit row; cue Pop on copy.
