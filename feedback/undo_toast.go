package feedback

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// UndoToastOptions tune an UndoToast.
type UndoToastOptions struct {
	// Undo labels the button; "Undo" when empty.
	Undo string
	// Close names the close button; "Close" when empty.
	Close string
	// For is how long it shows; 5 s when zero.
	For time.Duration
}

// UndoToast is a toast with Undo and a countdown ring.
func UndoToast(c *ui.Context, message string, open *bool, o UndoToastOptions) (root, undo ui.Element) {
	if message == "" || open == nil || o.For < 0 {
		panic("feedback: UndoToast needs a message, open state and a positive For")
	}
	if o.Undo == "" {
		o.Undo = "Undo"
	}
	if o.Close == "" {
		o.Close = "Close"
	}
	if o.For == 0 {
		o.For = 5 * time.Second
	}
	t := c.Theme()
	root = ui.Column(c).AlignItems(ui.Start)
	root.Children(func() {
		if !*open {
			return
		}
		c.Key("toast")
		card := toastFace(c, ui.Row(c)).AlignItems(ui.Center).Role(ui.RoleStatus).Label(message)
		ends := ui.Local(card, "ends", func() time.Time { return c.Now().Add(o.For) })
		if !c.Now().Before(*ends) {
			*open = false
			root.SetChanged()
			kit.Log("UndoToast", "expired")
			return
		}
		c.After(ends.Sub(c.Now()))
		kit.Arrive(card, ui.Pose{Y: t.Space(4), Fade: 1}, 0).Transition(ui.ElementTransition{Exit: &ui.Motion{X: t.Space(10)}})
		card.Children(func() {
			d, w := t.Space(5), t.Space(0.5)
			ui.Box(c).Size(d, d).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
				left := float32(max(ends.Sub(p.Now()), 0)) / float32(o.For)
				if left > 0 {
					p.AnimationFrame()
				}
				cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, r.W/2-w
				var ring, arc ui.Path
				p.StrokePath(ring.Circle(cx, cy, rad), w, t.Border)
				p.Transform(ui.Rotation(-90), cx, cy, func() { p.StrokePath(arc.Circle(cx, cy, rad).Trim(0, left), w, t.Accent) })
			})
			ui.Text(c, message).Grow(1)
			undo = button(c, o.Undo, false)
			if kit.Clicked(undo, sound.Tap) {
				*open = false
				root.SetChanged()
				kit.Log("UndoToast", "undone")
			}
			if kit.Clicked(closer(c, o.Close), sound.Close) {
				*open = false
				root.SetChanged()
				kit.Log("UndoToast", "closed")
			}
		})
	})
	return root, undo
}
