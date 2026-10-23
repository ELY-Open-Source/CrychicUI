package files

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PathBar shows folder breadcrumbs; cue Tap.
func PathBar(c *ui.Context, parts []string, at *int) ui.Element {
	if len(parts) == 0 || at == nil || *at < 0 || *at >= len(parts) {
		panic("files: PathBar needs a path and valid index")
	}
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(1)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		for i, part := range parts {
			if part == "" {
				panic("files: PathBar parts must not be empty")
			}
			if i > 0 {
				ui.Icon(c, icon.Must("chevron-right")).Size(16, 16).TextColor(t.TextMuted)
			}
			c.Key(i)
			b := kit.Press(ui.ButtonBase(c).Label(part).Padding(t.Space(1), t.Space(1.5)).Radius(t.Radius).Checked(i == *at))
			b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy))
			b.Children(func() { ui.Text(c, part).FontSize(t.Rem(0.85)).SingleLine() })
			if b.Clicked() && *at != i {
				*at = i
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("PathBar", "opened", "depth", i)
			}
		}
	})
	return root
}
