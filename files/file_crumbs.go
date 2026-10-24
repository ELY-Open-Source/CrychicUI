package files

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FileCrumbsOptions name the crumb menu action.
type FileCrumbsOptions struct{ Open string }

// FileCrumbs shows path crumbs with folder menus; cue Open.
func FileCrumbs(c *ui.Context, label string, parts []string, at *int, o FileCrumbsOptions) ui.Element {
	if label == "" || len(parts) == 0 || at == nil || *at < 0 || *at >= len(parts) {
		panic(fmt.Sprintf("files: FileCrumbs %q needs a path and a position", label))
	}
	for _, part := range parts {
		if part == "" {
			panic(fmt.Sprintf("files: FileCrumbs %q parts must not be empty", label))
		}
	}
	def(&o.Open, "Open")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(0.5)).AlignItems(ui.Center).FocusGroup(ui.Horizontal).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for i, part := range parts {
			// Keyed by the path prefix: a new path, new crumbs.
			c.Key(fmt.Sprintf("%q", parts[:i+1]))
			ui.Row(c).Gap(t.Space(0.5)).AlignItems(ui.Center).Children(func() {
				if i > 0 {
					ui.Icon(c, icon.Must("chevron-right")).Size(14, 14).TextColor(t.TextMuted.Alpha(0.6))
				}
				on := i == *at
				b := kit.Press(ui.ButtonBase(c).Label(part).Padding(t.Space(1), t.Space(1.5)).Radius(t.Radius).Checked(on))
				kit.Stagger(c, root, b, i)
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(t.Accent.Alpha(0.1), unit(on)), ui.Snappy))
				b.Children(func() {
					ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
						ui.Icon(c, icon.Must("folder")).Size(14, 14).Shrink(0).TextColor(t.Warning)
						txt := ui.Text(c, part).FontSize(t.Rem(0.85)).Shrink(0).MaxWidth(t.Space(36)).SingleLine()
						if on {
							txt.Bold()
						}
						ch := ui.Icon(c, icon.Must("chevron-down")).Size(12, 12).Shrink(0).TextColor(t.TextMuted)
						ch.Opacity(ch.AnimateTo("hover", 0.4+0.6*unit(b.Hovered()), ui.Snappy))
					})
				})
				b.Menu(func(m *ui.Menu) {
					if m.Item(o.Open).Chosen() && *at != i {
						*at = i
						root.SetChanged()
						sound.Play(sound.Open)
						kit.Log("FileCrumbs", "opened", "depth", i)
					}
				})
			})
		}
	})
	return root
}
