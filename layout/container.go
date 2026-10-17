package layout

import "github.com/egoist/mygo/ui"

// ContainerOptions tune a Container.
type ContainerOptions struct {
	// MaxWidth is 64 rem when zero.
	MaxWidth float32
}

// Container is a centered, padded page column; its width animates.
func Container(c *ui.Context, o ContainerOptions, content func()) ui.Element {
	t := c.Theme()
	if o.MaxWidth == 0 {
		o.MaxWidth = t.Rem(64)
	}
	if !finite(o.MaxWidth) || o.MaxWidth < 0 {
		panic("layout: Container MaxWidth must be finite and non-negative")
	}
	return ui.Column(c).FillWidth().MaxWidth(o.MaxWidth).Margin(0, ui.Auto).Padding(t.Space(6)).
		Transition(ui.ElementTransition{Size: true}).Children(content)
}
