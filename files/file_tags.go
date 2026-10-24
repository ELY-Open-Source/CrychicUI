package files

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FileTags toggles colored file tags, checked when on; cue Toggle.
func FileTags(c *ui.Context, label string, tags []string, selected *[]bool) ui.Element {
	if label == "" || len(tags) == 0 || selected == nil || len(*selected) != len(tags) {
		panic(fmt.Sprintf("files: FileTags %q needs tags and one flag each", label))
	}
	seen := map[string]bool{}
	for _, tag := range tags {
		if tag == "" || seen[tag] {
			panic(fmt.Sprintf("files: FileTags %q needs unique tags", label))
		}
		seen[tag] = true
	}
	t := c.Theme()
	inks := []ui.Color{t.Accent, t.Success, t.Warning, t.Danger, t.TextMuted}
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).FocusGroup(ui.Horizontal).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for i, tag := range tags {
			c.Key(tag)
			ink := inks[i%len(inks)]
			on := (*selected)[i]
			b := kit.Press(ui.ButtonBase(c).Label(tag).Tooltip(tag).Padding(t.Space(1), t.Space(2.5), t.Space(1), t.Space(2)).Radius(999).Checked(on))
			kit.Stagger(c, root, b, i)
			b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(ink.Alpha(0.14), unit(on)), ui.Snappy))
			b.Border(1, b.AnimateColor("line", t.Border.Mix(ink, unit(on)), ui.Snappy))
			k := b.AnimateTo("on", unit(on), ui.Bouncy)
			b.Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					ui.Grid(c).Columns(1).Size(t.Space(3.5), t.Space(3.5)).Center().Children(func() {
						dot := ui.Box(c).ColumnStart(1).RowStart(1).Size(t.Space(2.5), t.Space(2.5)).Radius(999).Background(ink)
						dot.Scale(1-k, 1-k)
						ui.Icon(c, icon.Must("check")).ColumnStart(1).RowStart(1).Size(14, 14).TextColor(ink).Scale(k, k).Opacity(min(1, k))
					})
					ui.Text(c, tag).FontSize(t.Rem(0.8)).Bold().Shrink(0).MaxWidth(t.Space(40)).SingleLine()
				})
			})
			if b.Clicked() {
				(*selected)[i] = !on
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("FileTags", "toggled", "index", i, "selected", !on)
			}
		}
	})
	return root
}
