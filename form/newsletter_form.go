package form

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// NewsletterFormOptions tune a NewsletterForm.
type NewsletterFormOptions struct {
	// Label, Placeholder, Subscribe and Done are caller copy.
	Label, Placeholder, Subscribe, Done string
}

// NewsletterForm takes an email; on subscribe a thank-you swaps in.
func NewsletterForm(c *ui.Context, email *string, o NewsletterFormOptions) (root, subscribe ui.Element) {
	if email == nil {
		panic("form: NewsletterForm needs an email")
	}
	def(&o.Label, "Email")
	def(&o.Placeholder, "you@example.com")
	def(&o.Subscribe, "Subscribe")
	def(&o.Done, "You're subscribed.")
	t := c.Theme()
	root = ui.Box(c).FillWidth()
	done := ui.Local(root, "done", func() bool { return false })
	swap := ui.ElementTransition{Exit: &ui.Motion{Opacity: 0, Y: -t.Space(2)}}
	root.Children(func() {
		if *done {
			row := ui.Row(c).Key("done").Gap(t.Space(2)).AlignItems(ui.Center).TextColor(t.Success).Transition(swap)
			kit.Arrive(row, ui.Pose{Y: t.Space(2), Fade: 1}, 0).Children(func() {
				ui.Icon(c, icon.Must("circle-check")).Size(16, 16)
				ui.Text(c, o.Done).FontWeight(500)
			})
			return
		}
		ui.Row(c).Key("form").FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Transition(swap).Children(func() {
			edited(root, ui.TextInput(c, email).Label(o.Label).Placeholder(o.Placeholder).Grow(1))
			subscribe = primary(c, o.Subscribe).Disabled(!emailPattern.MatchString(*email))
		})
	})
	if kit.Clicked(subscribe, sound.Success) {
		*done = true
		c.AnimationFrame()
		kit.Log("NewsletterForm", "subscribed")
	}
	return root, subscribe
}
