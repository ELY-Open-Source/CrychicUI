package input

import (
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PortFieldOptions tune a PortField.
type PortFieldOptions struct {
	// Label and Error are caller copy.
	Label, Error string
}

// PortField edits a port from 1 through 65535.
func PortField(c *ui.Context, label string, value *int, o PortFieldOptions) ui.Element {
	if value == nil || *value < 0 {
		panic("input: PortField needs a value")
	}
	if o.Label == "" {
		o.Label = label
	}
	if o.Error == "" {
		o.Error = "Use a port from 1 through 65535."
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	text := ui.Local(root, "text", func() string { return strconv.Itoa(*value) })
	return root.Children(func() {
		ui.Text(c, o.Label).FontWeight(500)
		in := ui.TextInput(c, text).FillWidth().Label(o.Label).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		if in.Changed() {
			parsed, err := strconv.Atoi(*text)
			if err == nil && parsed >= 1 && parsed <= 65535 {
				*value = parsed
				root.SetChanged()
				sound.Play(sound.Tick)
			}
		}
		if !in.Focused() {
			*text = strconv.Itoa(*value)
		}
		if *text != "" {
			parsed, err := strconv.Atoi(*text)
			if err != nil || parsed < 1 || parsed > 65535 {
				kit.Arrive(ui.Text(c, o.Error).TextColor(t.Danger).FontSize(t.Rem(0.85)), ui.Pose{Fade: 1}, 0)
			}
		}
	})
}
