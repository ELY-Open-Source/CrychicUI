package input

import (
	"net/url"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// URLFieldOptions tune a URLField.
type URLFieldOptions struct {
	// Scheme, Placeholder and Error are caller copy.
	Scheme, Placeholder, Error string
}

// URLField keeps a scheme and validates hosts on blur.
func URLField(c *ui.Context, label string, value *string, o URLFieldOptions) ui.Element {
	if value == nil {
		panic("input: URLField needs a value")
	}
	if o.Scheme == "" {
		o.Scheme = "https://"
	}
	if o.Error == "" {
		o.Error = "Enter a valid URL host."
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	wasFocused := ui.Local(root, "focused", func() bool { return false })
	errText := ui.Local(root, "error", func() string { return "" })
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		row := ui.Row(c).FillWidth().AlignItems(ui.Center).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		row.Children(func() {
			kit.Muted(c, o.Scheme).Padding(t.Space(2), 0, t.Space(2), t.Space(3))
			in := ui.TextInput(c, value).Grow(1).Placeholder(o.Placeholder).Padding(t.Space(2), t.Space(3), t.Space(2), t.Space(1))
			changed(root, in)
			if in.Focused() {
				*errText = ""
			} else if *wasFocused && *value != "" {
				u, err := url.Parse(o.Scheme + *value)
				if err != nil || u.Hostname() == "" {
					*errText = o.Error
				}
			}
			*wasFocused = in.Focused()
		})
		if *errText != "" {
			kit.Arrive(ui.Text(c, *errText).TextColor(t.Danger).FontSize(t.Rem(0.85)), ui.Pose{Fade: 1}, 0)
		}
	})
}
