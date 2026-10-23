package files

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DownloadRowOptions are the download labels.
type DownloadRowOptions struct{ Open string }

// DownloadRow shows progress, speed, and open; cue Success.
func DownloadRow(c *ui.Context, item FileItem, progress float64, speed int64, opened *bool, o DownloadRowOptions) ui.Element {
	validFile("DownloadRow", item)
	if speed < 0 || opened == nil {
		panic("files: DownloadRow needs speed and bound state")
	}
	def(&o.Open, "Open")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border)
	root.Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			ui.Text(c, item.Name).SingleLine().Grow(1)
			kit.Numeric(c, fmt.Sprintf("%s/s", sizeLabel(speed))).FontSize(t.Rem(0.8)).TextColor(t.TextMuted)
			if progress >= 1 {
				b := kit.Press(ui.ButtonBase(c).Label(o.Open).Padding(t.Space(0.5), t.Space(2)).Radius(t.Radius))
				b.Children(func() { ui.Text(c, o.Open).FontSize(t.Rem(0.8)) })
				if b.Clicked() && !*opened {
					*opened = true
					root.SetChanged()
					sound.Play(sound.Success)
					kit.Log("DownloadRow", "opened")
				}
			}
		})
		progressBar(c, progress, t.Accent)
	})
	return root
}
