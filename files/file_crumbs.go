package files

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FileCrumbsOptions are the crumb menu labels.
type FileCrumbsOptions struct{ Open string }

// FileCrumbs shows path crumbs with folder menus; cue Open.
func FileCrumbs(c *ui.Context, parts []string, at *int, o FileCrumbsOptions) ui.Element {
	if len(parts) == 0 || at == nil || *at < 0 || *at >= len(parts) {
		panic("files: FileCrumbs needs a path and valid index")
	}
	def(&o.Open, "Open")
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(1)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		for i, part := range parts {
			if part == "" {
				panic("files: FileCrumbs part must not be empty")
			}
			if i > 0 {
				ui.Icon(c, icon.Must("chevron-right")).Size(16, 16).TextColor(t.TextMuted)
			}
			c.Key(i)
			b := kit.Press(ui.ButtonBase(c).Label(part).Row().Gap(t.Space(1)).AlignItems(ui.Center).
				Padding(t.Space(1), t.Space(1.5)).Radius(t.Radius))
			b.Children(func() {
				ui.Icon(c, icon.Must("folder")).Size(16, 16).TextColor(t.TextMuted)
				ui.Text(c, part).FontSize(t.Rem(0.85))
			})
			b.Menu(func(m *ui.Menu) {
				if m.Item(o.Open).Chosen() {
					*at = i
					root.SetChanged()
					sound.Play(sound.Open)
					kit.Log("FileCrumbs", "opened", "depth", i)
				}
			})
		}
	})
	return root
}
