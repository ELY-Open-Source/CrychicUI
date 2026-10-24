package files

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// VersionItem describes one stored file version.
type VersionItem struct {
	ID, Author string
	Modified   time.Time
	Size       int64
	Current    bool
}

// VersionsOptions word the versions.
type VersionsOptions struct {
	Restore, Current string
	FormatDate       func(t time.Time) string
}

// Versions shows file versions on a rail; cue Success.
func Versions(c *ui.Context, label string, versions []VersionItem, restored *int, o VersionsOptions) ui.Element {
	if label == "" || len(versions) == 0 || restored == nil || *restored < -1 || *restored >= len(versions) {
		panic(fmt.Sprintf("files: Versions %q needs versions and a restore index or -1", label))
	}
	seen := map[string]bool{}
	for _, v := range versions {
		if v.ID == "" || seen[v.ID] || v.Size < 0 {
			panic(fmt.Sprintf("files: Versions %q needs unique IDs and sizes not below zero", label))
		}
		seen[v.ID] = true
	}
	def(&o.Restore, "Restore")
	def(&o.Current, "Current")
	if o.FormatDate == nil {
		o.FormatDate = func(t time.Time) string { return t.Format("Jan 2, 15:04") }
	}
	t := c.Theme()
	root := kit.Card(c).FillWidth().Padding(t.Space(1), 0).Role(ui.RoleList).Label(label)
	root.Children(func() {
		for i, v := range versions {
			c.Key(v.ID)
			row := ui.Row(c).FillWidth().Gap(t.Space(3)).Padding(0, t.Space(3))
			kit.Stagger(c, root, row, i).Children(func() {
				ui.Column(c).Width(t.Space(4)).Shrink(0).AlignItems(ui.Center).FillHeight().Children(func() {
					top := t.Border.Alpha(unit(i > 0))
					ui.Box(c).Width(1).Height(t.Space(3)).Background(top)
					ink := t.TextMuted
					if v.Current {
						ink = t.Accent
					}
					ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Shrink(0).Radius(999).Background(ink).Border(t.Space(0.5), t.Background)
					ui.Box(c).Width(1).Grow(1).Background(t.Border.Alpha(unit(i < len(versions)-1)))
				})
				ui.Row(c).Grow(1).MinWidth(0).Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2), 0).Children(func() {
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, v.ID).Font("monospace").FontSize(t.Rem(0.8)).Bold().Shrink(0)
							if v.Current {
								kit.Chip(c, o.Current, t.Accent)
							}
						})
						ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							if v.Author != "" {
								display.Avatar(c, v.Author, nil).Size(t.Space(4.5), t.Space(4.5)).FontSize(t.Rem(0.5)).Shrink(0)
								kit.Muted(c, v.Author).FontSize(t.Rem(0.75)).Shrink(0)
								kit.Muted(c, "·").FontSize(t.Rem(0.75))
							}
							kit.Muted(c, o.FormatDate(v.Modified)).FontSize(t.Rem(0.75)).Shrink(0)
							kit.Muted(c, "·").FontSize(t.Rem(0.75))
							sizeText(c, v.Size).FontSize(t.Rem(0.75))
						})
					})
					if v.Current {
						return
					}
					name := o.Restore + " " + v.ID
					b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(o.Restore).Padding(t.Space(1), t.Space(2.5)).Radius(t.Radius).Border(1, t.Border))
					b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy)).Children(func() {
						ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							ui.Icon(c, icon.Must("history")).Size(12, 12)
							ui.Text(c, o.Restore).FontSize(t.Rem(0.8)).Shrink(0)
						})
					})
					if b.Clicked() {
						*restored = i
						root.SetChanged()
						sound.Play(sound.Success)
						kit.Log("Versions", "restored", "index", i)
					}
				})
			})
		}
	})
	return root
}
