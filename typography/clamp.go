package typography

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ClampOptions tune a Clamp.
type ClampOptions struct {
	// More and Less label the toggle.
	More, Less string
}

// Clamp limits text lines and springs when toggled.
func Clamp(c *ui.Context, text string, lines int, open *bool, o ClampOptions) ui.Element {
	if lines <= 0 || open == nil {
		panic("typography: Clamp needs lines and open")
	}
	if o.More == "" {
		o.More = "More"
	}
	if o.Less == "" {
		o.Less = "Less"
	}
	t := c.Theme()
	spring := ui.Snappy
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Transition(ui.ElementTransition{Size: true, Animation: &spring})
	return root.Children(func() {
		body := ui.Text(c, text).FillWidth()
		if !*open {
			body.MaxLines(lines)
		}
		label := o.More
		if *open {
			label = o.Less
		}
		button := ui.ButtonBase(c).AlignSelf(ui.Start).Padding(t.Space(1), 0).TextColor(t.Accent)
		kit.Press(button).Children(func() { ui.Text(c, label).ContentTransition(ui.ContentTransition{}) })
		if kit.Clicked(button, sound.Toggle) {
			*open = !*open
			root.SetChanged()
			kit.Log("Clamp", "toggled", "open", *open)
		}
	})
}
