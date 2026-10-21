package input

import (
	"encoding/json"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// JSONFieldOptions tune a JSONField.
type JSONFieldOptions struct {
	// Label, Placeholder and Error are caller copy.
	Label, Placeholder, Error string
}

// JSONField edits JSON and fades in validation errors.
func JSONField(c *ui.Context, value *string, o JSONFieldOptions) ui.Element {
	if value == nil {
		panic("input: JSONField needs a value")
	}
	if o.Label == "" {
		o.Label = "JSON"
	}
	if o.Placeholder == "" {
		o.Placeholder = `{"name":"Crychic"}`
	}
	if o.Error == "" {
		o.Error = "Enter valid JSON."
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		ui.Text(c, o.Label).FontWeight(500)
		in := ui.TextArea(c, value).FillWidth().Lines(4, 10).Placeholder(o.Placeholder).Padding(t.Space(3)).Radius(t.Radius).
			Background(t.Background).Border(1, t.Border)
		changed(root, in)
		if strings.TrimSpace(*value) != "" && !json.Valid([]byte(*value)) {
			kit.Arrive(ui.Text(c, o.Error).TextColor(t.Danger).FontSize(t.Rem(0.85)), ui.Pose{Fade: 1}, 0)
		}
	})
}
