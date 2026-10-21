package input

import (
	"github.com/egoist/mygo/ui"
)

// AutosizeInputOptions tune an AutosizeInput.
type AutosizeInputOptions struct {
	// Placeholder is caller copy; MinWidth and MaxWidth bound the field.
	Placeholder        string
	MinWidth, MaxWidth float32
}

// AutosizeInput follows its text width with a transition.
func AutosizeInput(c *ui.Context, value *string, o AutosizeInputOptions) ui.Element {
	if value == nil {
		panic("input: AutosizeInput needs a value")
	}
	if o.MinWidth == 0 {
		o.MinWidth = c.Theme().Rem(8)
	}
	if o.MaxWidth == 0 {
		o.MaxWidth = c.Theme().Rem(40)
	}
	if !finite(float64(o.MinWidth)) || !finite(float64(o.MaxWidth)) || o.MinWidth <= 0 || o.MaxWidth < o.MinWidth {
		panic("input: AutosizeInput needs valid width bounds")
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	return root.Children(func() {
		text := *value
		if text == "" {
			text = o.Placeholder
		}
		textWidth, _ := c.MeasureText(0, ui.Span{Text: text})
		target := min(o.MaxWidth, max(o.MinWidth, textWidth+t.Space(6)))
		width := root.AnimateTo("width", target, ui.Snappy)
		in := ui.TextInputBase(c, value).Width(width).Placeholder(o.Placeholder).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		changed(root, in)
	})
}
