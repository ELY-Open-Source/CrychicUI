package files

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RecentFilesOptions word when a file was opened.
type RecentFilesOptions struct {
	FormatDate func(t time.Time) string
}

// RecentFiles lists files by when they were opened; cue Tap.
func RecentFiles(c *ui.Context, label string, list []FileItem, opened *int, o RecentFilesOptions) ui.Element {
	validFiles("RecentFiles", label, list)
	if opened == nil || *opened < -1 || *opened >= len(list) {
		panic(fmt.Sprintf("files: RecentFiles %q needs an open index or -1", label))
	}
	if o.FormatDate == nil {
		o.FormatDate = func(t time.Time) string { return t.Format("Jan 2, 15:04") }
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(0.5)).FocusGroup(ui.Vertical).Role(ui.RoleList).Label(label)
	root.Children(func() {
		seen := map[string]int{}
		for i, f := range list {
			key := fmt.Sprintf("%q %q", f.Path, f.Name)
			seen[key]++
			c.Key(fmt.Sprint(key, seen[key]))
			on := i == *opened
			row := kit.Press(ui.ButtonBase(c).Label(f.Name).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(1.5), t.Space(2)).
				Radius(t.Radius).Checked(on))
			kit.Stagger(c, root, row, i)
			hover := row.AnimateTo("hover", unit(row.Hovered()), ui.Snappy)
			row.Background(row.AnimateColor("bg", t.SurfaceHover.Alpha(unit(row.Hovered())).Mix(t.Accent.Alpha(0.12), unit(on)), ui.Snappy))
			row.Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
					fileTile(c, f, t.Space(8))
					ui.Column(c).Grow(1).MinWidth(0).Children(func() {
						ui.Text(c, f.Name).FontSize(t.Rem(0.85)).FillWidth().SingleLine()
						if f.Path != "" {
							kit.Muted(c, f.Path).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
						}
					})
					kit.Muted(c, o.FormatDate(f.Modified)).FontSize(t.Rem(0.75)).Shrink(0)
					arrow := ui.Icon(c, icon.Must("arrow-up-right")).Size(14, 14).Shrink(0).TextColor(t.TextMuted)
					arrow.Opacity(hover).Translate(t.Space(1)*(1-hover), 0)
				})
			})
			if row.Clicked() {
				*opened = i
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("RecentFiles", "opened", "index", i)
			}
		}
	})
	return root
}
