package input

import (
	"regexp"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// EmailFieldOptions tune an EmailField.
type EmailFieldOptions struct {
	// Placeholder and Error are caller copy.
	Placeholder, Error string
}

// EmailField validates its address when focus leaves.
func EmailField(c *ui.Context, label string, value *string, o EmailFieldOptions) ui.Element {
	if value == nil {
		panic("input: EmailField needs a value")
	}
	if o.Error == "" {
		o.Error = "Enter a valid email address."
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	wasFocused := ui.Local(root, "focused", func() bool { return false })
	errText := ui.Local(root, "error", func() string { return "" })
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		in := ui.TextInput(c, value).FillWidth().Placeholder(o.Placeholder).Padding(t.Space(2), t.Space(3)).
			Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		changed(root, in)
		if in.Focused() {
			*errText = ""
		} else if *wasFocused && *value != "" && !emailPattern.MatchString(*value) {
			*errText = o.Error
		}
		*wasFocused = in.Focused()
		if *errText != "" {
			kit.Arrive(ui.Text(c, *errText).TextColor(t.Danger).FontSize(t.Rem(0.85)), ui.Pose{Fade: 1}, 0)
		}
	})
}
