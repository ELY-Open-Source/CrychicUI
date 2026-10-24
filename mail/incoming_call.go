package mail

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// IncomingCallOptions are the call answers.
type IncomingCallOptions struct{ Accept, Decline string }

// IncomingCall accepts or declines a pulsing call; cue Notify.
func IncomingCall(c *ui.Context, caller string, accepted, declined *bool, o IncomingCallOptions) ui.Element {
	if caller == "" || accepted == nil || declined == nil {
		panic("mail: IncomingCall needs caller and state")
	}
	def(&o.Accept, "Accept")
	def(&o.Decline, "Decline")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(4)).Radius(t.Radius).Border(1, t.Border)
	root.Children(func() {
		done := *accepted || *declined
		phase := 0.0
		if !done {
			phase = float64(root.Loop("ring", 1200*time.Millisecond, ui.Linear)) * 2 * math.Pi
		}
		scale := 1 + 0.06*float32(math.Sin(phase))
		ui.Box(c).Size(t.Space(12), t.Space(12)).Radius(999).Background(t.Accent.Alpha(0.2)).Scale(scale, scale).Center().Children(func() {
			ui.Text(c, initial(caller)).TextColor(t.Accent)
		})
		ui.Text(c, caller).Bold()
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			accept := button(c, o.Accept).Background(t.Success).TextColor(t.AccentText).Disabled(done)
			accept.Children(func() { ui.Text(c, o.Accept) })
			if accept.Clicked() && !done {
				*accepted = true
				root.SetChanged()
				sound.Play(sound.Notify)
				kit.Log("IncomingCall", "accepted")
			}
			decline := button(c, o.Decline).Background(t.Danger).TextColor(t.AccentText).Disabled(done)
			decline.Children(func() { ui.Text(c, o.Decline) })
			if decline.Clicked() && !done {
				*declined = true
				root.SetChanged()
				sound.Play(sound.Notify)
				kit.Log("IncomingCall", "declined")
			}
		})
	})
	return root
}
