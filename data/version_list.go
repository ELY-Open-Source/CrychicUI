package data

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Version is a named, dated save with a note.
type Version struct{ Name, When, Note string }

// VersionListOptions tune a VersionList.
type VersionListOptions struct {
	// Restore names the button; Current marks the live one.
	Restore, Current string
}

// VersionList lists versions; restoring makes one current.
func VersionList(c *ui.Context, label string, versions []Version, current *int, o VersionListOptions) ui.Element {
	if label == "" || len(versions) == 0 || current == nil || !valid(*current, len(versions)) {
		panic(fmt.Sprintf("data: VersionList %q needs versions and a current one", label))
	}
	var names []string
	for _, v := range versions {
		names = append(names, v.Name)
	}
	unique("VersionList", label, names)
	def(&o.Restore, "Restore")
	def(&o.Current, "Current")
	t := c.Theme()
	root := frame(c, label).Dividers(1, t.Border)
	shown := *current
	return root.Children(func() {
		for i, v := range versions {
			c.Key(v.Name)
			on := i == shown
			row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2.5), t.Space(3))
			row.Background(row.AnimateColor("bg", t.Accent.Alpha(0.08*unit(on)), ui.Snappy))
			kit.Arrive(row, ui.Pose{Y: t.Space(2), Fade: 1}, time.Duration(i)*40*time.Millisecond).Children(func() {
				ui.Column(c).Grow(1).Gap(t.Space(0.5)).Children(func() {
					ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						ui.Text(c, v.Name).Bold()
						kit.Muted(c, v.When)
					})
					if v.Note != "" {
						kit.Muted(c, v.Note)
					}
				})
				if on {
					ui.Text(c, o.Current).FontSize(t.Rem(0.8)).Bold().TextColor(t.Success).Padding(t.Space(0.5), t.Space(2)).
						Radius(999).Background(t.Success.Alpha(0.12)).ContentTransition(ui.ContentTransition{})
					return
				}
				name := o.Restore + " " + v.Name
				b := kit.Press(ui.ButtonBase(c).Label(name).Padding(t.Space(1), t.Space(2.5)).Radius(t.Radius).Border(1, t.Border))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy)).Children(func() {
					ui.Text(c, o.Restore).FontSize(t.Rem(0.85))
				})
				if b.Clicked() {
					*current = i
					root.SetChanged()
					sound.Play(sound.Success)
					kit.Log("VersionList", "restored", "index", i)
				}
			})
		}
	})
}
