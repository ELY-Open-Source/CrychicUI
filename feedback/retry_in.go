package feedback

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RetryInOptions tune a RetryIn.
type RetryInOptions struct {
	// Wait writes seconds left; "Retry in %d s".
	Wait string
	// Retry labels the enabled button; "Retry" when empty.
	Retry string
}

// RetryIn counts down whole seconds, then enables Retry.
func RetryIn(c *ui.Context, at time.Time, o RetryInOptions) ui.Element {
	if o.Wait == "" {
		o.Wait = "Retry in %d s"
	}
	if o.Retry == "" {
		o.Retry = "Retry"
	}
	left := at.Sub(c.Now())
	text, waiting := o.Retry, left > 0
	if waiting {
		secs := left / time.Second
		if left%time.Second != 0 {
			secs++
		}
		text = fmt.Sprintf(o.Wait, secs)
		c.After(left - (secs-1)*time.Second)
	}
	t := c.Theme()
	b := ui.ButtonBase(c).Label(text).Disabled(waiting).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	bg := t.Surface
	if b.Hovered() {
		bg = t.SurfaceHover
	}
	b.Background(b.AnimateColor("bg", bg, quick))
	kit.Press(b).Children(func() { kit.Numeric(c, text).ContentTransition(ui.ContentTransition{Numeric: true, Down: true}) })
	if kit.Clicked(b, sound.Tap) {
		kit.Log("RetryIn", "retry")
	}
	return b
}
