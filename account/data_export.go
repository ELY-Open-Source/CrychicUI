package account

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ExportStatus is where a data export stands.
type ExportStatus int

const (
	ExportIdle ExportStatus = iota
	ExportPreparing
	ExportReady
)

// DataExportOptions word the card; URL is the finished archive.
type DataExportOptions struct{ Title, Note, Request, Preparing, Ready, Download, URL string }

// DataExport requests an export with status; cue Success.
func DataExport(c *ui.Context, status *ExportStatus, progress float64, o DataExportOptions) ui.Element {
	if status == nil || *status < ExportIdle || *status > ExportReady || !finite(progress) || progress < 0 || progress > 1 {
		panic(fmt.Sprintf("account: DataExport %q needs a status and progress in 0..1", o.Title))
	}
	if *status == ExportReady && o.URL == "" {
		panic(fmt.Sprintf("account: DataExport %q is ready without a URL", o.Title))
	}
	def(&o.Title, "Export your data")
	def(&o.Note, "Profile, orders, messages and files as a ZIP archive.")
	def(&o.Request, "Request export")
	def(&o.Preparing, "Preparing your archive…")
	def(&o.Ready, "Your archive is ready")
	def(&o.Download, "Download ZIP")
	t := c.Theme()
	root := card(c).Gap(t.Space(4)).Padding(t.Space(5)).Role(ui.RoleGroup)
	root.Children(func() {
		head(c, "database", o.Title, o.Note)
		c.Key(*status)
		state := ui.Column(c).FillWidth().Gap(t.Space(3))
		kit.Arrive(state, ui.Pose{Y: t.Space(1), Fade: 1}, 0).Children(func() {
			switch *status {
			case ExportIdle:
				b := primary(c, o.Request, "file-archive").FillWidth()
				if b.Clicked() {
					*status = ExportPreparing
					root.SetChanged()
					sound.Play(sound.Success)
					kit.Log("DataExport", "requested")
				}
			case ExportPreparing:
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					spin := glyph(c, "loader-circle", 16).TextColor(t.Accent)
					spin.Rotate(360 * float32(spin.Loop("spin", time.Second, ui.Linear)))
					ui.Text(c, o.Preparing).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).Wrap()
					kit.Numeric(c, fmt.Sprintf("%.0f%%", progress*100)).FontSize(t.Rem(0.78)).Bold().Shrink(0)
				})
				meter(c, float32(progress), t.Accent)
			case ExportReady:
				note(c, "circle-check", o.Ready, t.Success)
				b := primary(c, o.Download, "download").FillWidth()
				if b.Clicked() {
					c.OpenURL(o.URL)
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("DataExport", "downloaded")
				}
			}
		})
	})
	return root
}
