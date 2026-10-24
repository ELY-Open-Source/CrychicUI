package dashboard

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ReportExportOptions word the card; Done takes the format.
type ReportExportOptions struct{ Title, Export, Done string }

// ReportExport picks a format and exports; cue Success.
func ReportExport(c *ui.Context, format *string, formats []string, exported *bool, o ReportExportOptions) (root, export ui.Element) {
	if format == nil || len(formats) == 0 || exported == nil {
		panic("dashboard: ReportExport needs a format, formats and an exported flag")
	}
	unique("ReportExport", formats)
	if !slices.Contains(formats, *format) {
		panic(fmt.Sprintf("dashboard: ReportExport format %q is not among %q", *format, formats))
	}
	def(&o.Title, "Export report")
	def(&o.Export, "Export")
	def(&o.Done, "%s report ready")
	t := c.Theme()
	root = card(c)
	root.Children(func() {
		kit.CardHead(c, "download", o.Title, nil)
		ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Children(func() {
			grid := ui.Grid(c).Columns(min(len(formats), 4)).FillWidth().Gap(t.Space(2))
			grid.Children(func() {
				for i, choice := range formats {
					c.Key(choice)
					on := choice == *format
					b := kit.Press(ui.ButtonBase(c).Label(choice).Column().Gap(t.Space(1.5)).Padding(t.Space(2.5), t.Space(1)).
						MinWidth(0).Radius(t.Radius).Border(1, t.Border).Checked(on))
					b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(t.Accent.Alpha(0.1), unit(on)), ui.Snappy)).
						BorderColor(b.AnimateColor("edge", t.Border.Mix(t.Accent, unit(on)), ui.Snappy)).
						TextColor(b.AnimateColor("ink", t.Text.Mix(t.Accent, unit(on)), ui.Snappy))
					kit.Stagger(c, grid, b, i).Children(func() {
						ui.Icon(c, icon.Must(formatGlyph(choice))).Size(20, 20)
						ui.Text(c, choice).FontSize(t.Rem(0.8)).Bold().MaxWidth(t.Space(30)).SingleLine()
					})
					if b.Clicked() && !on {
						*format = choice
						*exported = false
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("ReportExport", "format", "index", i)
					}
				}
			})
			if *exported {
				c.Key("done")
				banner := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).
					Background(t.Success.Alpha(0.1)).TextColor(t.Success)
				kit.Arrive(banner, ui.Pose{Scale: 0.95, Fade: 1}, 0).Children(func() {
					ui.Icon(c, icon.Must("circle-check")).Size(16, 16).Shrink(0)
					ui.Text(c, fmt.Sprintf(o.Done, *format)).Bold().FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).Wrap()
				})
			}
			c.Key("export")
			export = kit.Press(ui.ButtonBase(c).Label(o.Export).FillWidth().Center().Padding(t.Space(2.5)).Radius(t.Radius).TextColor(t.AccentText))
			export.Background(export.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(export.Hovered())), ui.Snappy))
			export.Children(func() {
				ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					ui.Icon(c, icon.Must("download")).Size(16, 16)
					ui.Text(c, o.Export).Bold().Shrink(0)
				})
			})
		})
	})
	if export.Clicked() {
		*exported = true
		root.SetChanged()
		sound.Play(sound.Success)
		kit.Log("ReportExport", "exported", "index", slices.Index(formats, *format))
	}
	return root, export
}

// formatGlyph is a file icon for a format name.
func formatGlyph(format string) string {
	switch strings.ToLower(format) {
	case "csv", "xlsx", "xls":
		return "file-spreadsheet"
	case "json":
		return "file-json"
	case "png", "jpg", "svg":
		return "file-image"
	case "pdf":
		return "file-text"
	}
	return "file"
}
