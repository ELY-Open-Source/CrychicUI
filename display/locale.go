package display

import "github.com/egoist/mygo/ui"

// Locale is a language pill cross-fading on change.
func Locale(c *ui.Context, label string) ui.Element {
	if label == "" {
		panic("display: Locale needs a label")
	}
	t := c.Theme()
	return ui.Text(c, label).Padding(t.Space(0.5), t.Space(2)).Radius(999).Background(t.Accent.Alpha(0.12)).
		TextColor(t.Accent).ContentTransition(ui.ContentTransition{})
}
