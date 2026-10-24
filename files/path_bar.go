package files

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PathBar shows the folders of a path; cue Tap.
func PathBar(c *ui.Context, label string, parts []string, at *int) ui.Element {
	if label == "" || len(parts) == 0 || at == nil || *at < 0 || *at >= len(parts) {
		panic(fmt.Sprintf("files: PathBar %q needs a path and a position", label))
	}
	for _, part := range parts {
		if part == "" {
			panic(fmt.Sprintf("files: PathBar %q parts must not be empty", label))
		}
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(0.5)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(1.5)).
		Radius(t.Radius).Background(t.Surface).Border(1, t.Border).FocusGroup(ui.Horizontal).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for i, part := range parts {
			// Keyed by the path prefix: a new path, new folders.
			c.Key(fmt.Sprintf("%q", parts[:i+1]))
			ui.Row(c).Gap(t.Space(0.5)).AlignItems(ui.Center).Children(func() {
				if i > 0 {
					ui.Icon(c, icon.Must("chevron-right")).Size(14, 14).TextColor(t.TextMuted.Alpha(0.6))
				}
				on := i == *at
				b := kit.Press(ui.ButtonBase(c).Label(part).Tooltip(part).Padding(t.Space(0.75), t.Space(1.5)).Radius(t.Radius).Checked(on))
				kit.Stagger(c, root, b, i)
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(t.Background, unit(on)), ui.Snappy))
				b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Text, unit(on || b.Hovered())), ui.Snappy))
				b.Children(func() {
					ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
						glyph := "folder"
						if i == 0 {
							glyph = "house"
						}
						ui.Icon(c, icon.Must(glyph)).Size(14, 14).Shrink(0)
						txt := ui.Text(c, part).FontSize(t.Rem(0.85)).Shrink(0).MaxWidth(t.Space(40)).SingleLine()
						if on {
							txt.Bold()
						}
					})
				})
				if b.Clicked() && !on {
					*at = i
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("PathBar", "opened", "depth", i)
				}
			})
		}
	})
	return root
}
