package input

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ClearableFieldOptions tune a ClearableField.
type ClearableFieldOptions struct {
	// Placeholder and Clear label the field and action.
	Placeholder, Clear string
}

// ClearableField shows a clear button while text is present.
func ClearableField(c *ui.Context, label string, value *string, o ClearableFieldOptions) ui.Element {
	if value == nil {
		panic("input: ClearableField needs a value")
	}
	if o.Clear == "" {
		o.Clear = "Clear " + label
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		row := ui.Row(c).FillWidth().AlignItems(ui.Center).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		row.Children(func() {
			in := ui.TextInput(c, value).Grow(1).Placeholder(o.Placeholder).Padding(t.Space(2), t.Space(3))
			changed(root, in)
			if *value != "" {
				clear := ui.ButtonBase(c).Key("clear").Label(o.Clear).Padding(t.Space(2)).TextColor(t.TextMuted)
				kit.Arrive(kit.Press(clear), ui.Pose{Scale: 0.5, Fade: 1}, 0).Children(func() {
					ui.Icon(c, icon.Must("x")).Size(16, 16)
				})
				if kit.Clicked(clear, sound.Close) {
					*value = ""
					root.SetChanged()
					kit.Log("ClearableField", "cleared")
				}
			}
		})
	})
}
