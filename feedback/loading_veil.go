package feedback

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// LoadingVeilOptions tune a LoadingVeil.
type LoadingVeilOptions struct {
	// Label names the spinner; "Loading" when empty.
	Label string
}

var fade = ui.ElementTransition{Exit: &ui.Motion{}}

// LoadingVeil dims its content under a spinner while loading.
func LoadingVeil(c *ui.Context, loading bool, o LoadingVeilOptions, content func()) ui.Element {
	if content == nil {
		panic("feedback: LoadingVeil needs content")
	}
	if o.Label == "" {
		o.Label = "Loading"
	}
	t := c.Theme()
	root := ui.Box(c).Radius(t.Radius).Clip()
	return root.Children(func() {
		content()
		if !loading {
			return
		}
		c.Key("veil")
		veil := ui.Box(c).Absolute().Top(0).Left(0).Right(0).Bottom(0).Center().Background(t.Background.Alpha(0.7))
		kit.Arrive(veil, ui.Pose{Fade: 1}, 0).Transition(fade).Children(func() {
			Spin(c, o.Label, SpinOptions{})
		})
	})
}
