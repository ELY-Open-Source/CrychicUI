package input

import (
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AutoTextAreaOptions tune an AutoTextArea.
type AutoTextAreaOptions struct {
	// MinLines and MaxLines bound height; zero is 3-8.
	MinLines, MaxLines int
	Placeholder        string
}

// AutoTextArea grows with its lines and transitions height.
func AutoTextArea(c *ui.Context, value *string, o AutoTextAreaOptions) ui.Element {
	if value == nil {
		panic("input: AutoTextArea needs a value")
	}
	if o.MinLines == 0 && o.MaxLines == 0 {
		o.MinLines, o.MaxLines = 3, 8
	}
	if o.MinLines <= 0 || o.MaxLines < o.MinLines {
		panic("input: AutoTextArea needs valid line bounds")
	}
	in := ui.TextArea(c, value).FillWidth().Lines(o.MinLines, o.MaxLines).Placeholder(o.Placeholder).
		Padding(c.Theme().Space(3)).Radius(c.Theme().Radius).Background(c.Theme().Background).Border(1, c.Theme().Border).
		Transition(ui.ElementTransition{Size: true})
	if in.Changed() {
		sound.Play(sound.Tick)
	}
	return in
}
