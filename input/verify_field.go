package input

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// VerifyFieldOptions tune a VerifyField.
type VerifyFieldOptions struct {
	// Placeholder, Send and Resend are copy; Countdown defaults to 3s.
	Placeholder, Send, Resend string
	Countdown                 time.Duration
}

// VerifyField has a send button with a countdown.
func VerifyField(c *ui.Context, label string, code *string, o VerifyFieldOptions) ui.Element {
	if code == nil {
		panic("input: VerifyField needs a code")
	}
	if o.Send == "" {
		o.Send = "Send code"
	}
	if o.Resend == "" {
		o.Resend = "Resend in"
	}
	if o.Countdown == 0 {
		o.Countdown = 3 * time.Second
	}
	if o.Countdown < 0 {
		panic("input: VerifyField countdown must be positive")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	sent := ui.Local(root, "sent", func() bool { return false })
	due := ui.Local(root, "due", func() time.Time { return time.Time{} })
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			in := ui.TextInput(c, code).Grow(1).Placeholder(o.Placeholder).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
			changed(root, in)
			remaining := 0
			if *sent {
				left := due.Sub(c.Now())
				if left <= 0 {
					*sent = false
				} else {
					remaining = int(math.Ceil(left.Seconds()))
					next := left % time.Second
					if next == 0 {
						next = time.Second
					}
					c.After(next)
				}
			}
			label := o.Send
			if *sent {
				label = fmt.Sprintf("%s %ds", o.Resend, remaining)
			}
			button := ui.ButtonBase(c).Label(label).Disabled(*sent).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Accent).TextColor(t.AccentText).
				Transition(ui.ElementTransition{Size: true})
			kit.Press(button).Children(func() { ui.Text(c, label).SingleLine().ContentTransition(ui.ContentTransition{}) })
			if kit.Clicked(button, sound.Send) {
				*sent = true
				*due = c.Now().Add(o.Countdown)
				c.After(o.Countdown)
				root.SetChanged()
				kit.Log("VerifyField", "sent")
			}
		})
	})
}
