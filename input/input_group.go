package input

import (
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// InputGroup joins inputs and one action with shared borders.
func InputGroup(c *ui.Context, labels []string, values []*string, actionLabel string) (root, action ui.Element) {
	if len(labels) == 0 || len(labels) != len(values) || actionLabel == "" {
		panic("input: InputGroup needs one value per label and an action")
	}
	for _, value := range values {
		if value == nil {
			panic("input: InputGroup needs every value")
		}
	}
	t := c.Theme()
	root = ui.Row(c).FillWidth()
	root.Children(func() {
		for i := range labels {
			c.Key(i)
			field := ui.TextInputBase(c, values[i]).Grow(1).Basis(0).Placeholder(labels[i]).Padding(t.Space(2), t.Space(3)).Background(t.Background)
			if i == 0 {
				field.BorderWidth(1, 0, 1, 1).BorderColor(t.Border).Radius(t.Radius, 0, 0, t.Radius)
			} else {
				field.BorderWidth(1, 0, 1, 0).BorderColor(t.Border).Radius(0)
			}
			changed(root, field)
		}
		action = ui.ButtonBase(c).Label(actionLabel).Padding(t.Space(2), t.Space(3)).Background(t.Accent).TextColor(t.AccentText).
			BorderWidth(1, 1, 1, 0).BorderColor(t.Accent).Radius(0, t.Radius, t.Radius, 0)
		action.Children(func() { ui.Text(c, actionLabel).SingleLine() })
		if action.Clicked() {
			sound.Play(sound.Tap)
		}
	})
	return root, action
}
