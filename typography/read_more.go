package typography

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ReadMoreOptions tune a ReadMore.
type ReadMoreOptions struct {
	// Label names the expand action.
	Label string
}

// ReadMore fades a clamped tail and expands on click.
func ReadMore(c *ui.Context, text string, lines int, open *bool, o ReadMoreOptions) ui.Element {
	if lines <= 0 || open == nil {
		panic("typography: ReadMore needs lines and open state")
	}
	if o.Label == "" {
		o.Label = "Read more"
	}
	t := c.Theme()
	spring := ui.Snappy
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Transition(ui.ElementTransition{Size: true, Animation: &spring})
	return root.Children(func() {
		body := ui.Box(c).FillWidth()
		body.DrawOver(func(p *ui.Painter, r ui.Rect) {
			if !*open {
				h := min(t.Rem(2), r.H)
				g := ui.LinearGradient{From: t.Background.Alpha(0), To: t.Background, Angle: 180}
				p.FillGradient(ui.Rect{X: r.X, Y: r.Y + r.H - h, W: r.W, H: h}, g, 0)
			}
		})
		body.Children(func() {
			e := ui.Text(c, text).FillWidth()
			if !*open {
				e.MaxLines(lines)
			}
		})
		if !*open {
			button := ui.ButtonBase(c).Key("more").AlignSelf(ui.Start).Padding(t.Space(1), 0).TextColor(t.Accent).
				Transition(ui.ElementTransition{Exit: &ui.Motion{}})
			kit.Press(button).Children(func() { ui.Text(c, o.Label) })
			if kit.Clicked(button, sound.Tap) {
				*open = true
				root.SetChanged()
				kit.Log("ReadMore", "expanded")
			}
		}
	})
}
