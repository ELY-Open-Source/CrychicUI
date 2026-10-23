package dashboard

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ReportExportOptions are the export labels.
type ReportExportOptions struct{ Export string }

// ReportExport chooses format and exports; cue Success.
func ReportExport(c *ui.Context, format *string, formats []string, exported *bool, o ReportExportOptions) (root, export ui.Element) {
	if format == nil || len(formats) == 0 || exported == nil {
		panic("dashboard: ReportExport needs format and state")
	}
	def(&o.Export, "Export")
	valid := false
	for _, choice := range formats {
		valid = valid || *format == choice
	}
	if !valid {
		panic("dashboard: ReportExport needs a current format")
	}
	t := c.Theme()
	root = ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		for _, choice := range formats {
			if choice == "" {
				panic("dashboard: ReportExport format must not be empty")
			}
			b := kit.Press(ui.ButtonBase(c).Label(choice).Padding(t.Space(1), t.Space(2)).Radius(999).
				Border(1, t.Border).Checked(*format == choice))
			b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*format == choice)), ui.Snappy))
			b.Children(func() { ui.Text(c, choice).FontSize(t.Rem(0.8)) })
			if b.Clicked() && *format != choice {
				*format = choice
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("ReportExport", "format", "format", *format)
			}
		}
		export = kit.Press(ui.ButtonBase(c).Label(o.Export).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).
			Background(t.Accent).TextColor(t.AccentText))
		export.Children(func() { ui.Text(c, o.Export).SingleLine() })
	})
	if export.Clicked() {
		*exported = true
		root.SetChanged()
		sound.Play(sound.Success)
		kit.Log("ReportExport", "exported", "format", *format)
	}
	return root, export
}
