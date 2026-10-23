package files

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// VersionItem describes one stored file version.
type VersionItem struct {
	ID       string
	Modified time.Time
	Size     int64
	Current  bool
}

// VersionsOptions are the version labels.
type VersionsOptions struct{ Restore string }

// Versions shows file versions with restore; cue Success.
func Versions(c *ui.Context, label string, versions []VersionItem, restored *int, o VersionsOptions) ui.Element {
	if label == "" || len(versions) == 0 || restored == nil || *restored < -1 || *restored >= len(versions) {
		panic("files: Versions needs versions and a valid restore index")
	}
	def(&o.Restore, "Restore")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for i, version := range versions {
			if version.ID == "" || version.Size < 0 {
				panic("files: Versions need IDs and nonnegative sizes")
			}
			c.Key(version.ID)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).
				Radius(t.Radius).Border(1, t.Border)
			kit.Arrive(row, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
			row.Children(func() {
				ui.Icon(c, icon.Must("clock")).Size(16, 16).TextColor(t.TextMuted)
				ui.Column(c).Grow(1).Children(func() {
					ui.Text(c, version.ID).Font("monospace").FontSize(t.Rem(0.85))
					kit.Muted(c, date(version.Modified)).FontSize(t.Rem(0.75))
				})
				sizeText(c, version.Size)
				if !version.Current {
					b := kit.Press(ui.ButtonBase(c).Label(o.Restore).Padding(t.Space(1), t.Space(2)).Radius(t.Radius))
					b.Children(func() { ui.Text(c, o.Restore).FontSize(t.Rem(0.8)).SingleLine() })
					if b.Clicked() {
						*restored = i
						root.SetChanged()
						sound.Play(sound.Success)
						kit.Log("Versions", "restored", "index", i)
					}
				}
			})
		}
	})
	return root
}
