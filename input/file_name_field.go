package input

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FileNameFieldOptions tune a FileNameField.
type FileNameFieldOptions struct {
	// Placeholder is caller copy.
	Placeholder string
}

// FileNameField edits a name while its extension stays fixed.
func FileNameField(c *ui.Context, label string, name *string, extension string, o FileNameFieldOptions) ui.Element {
	if name == nil || strings.TrimSpace(extension) == "" {
		panic("input: FileNameField needs a name and extension")
	}
	extension = strings.TrimPrefix(extension, ".")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		row := ui.Row(c).FillWidth().AlignItems(ui.Center).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		row.Children(func() {
			in := ui.TextInput(c, name).Grow(1).Placeholder(o.Placeholder).Padding(t.Space(2), t.Space(3))
			changed(root, in)
			kit.Muted(c, "."+extension).Padding(t.Space(2), t.Space(3), t.Space(2), 0)
		})
	})
}
