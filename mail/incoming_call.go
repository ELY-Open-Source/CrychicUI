package mail

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// IncomingCallOptions word the answers and status lines.
type IncomingCallOptions struct{ Accept, Decline, Calling, Accepted, Declined string }

// IncomingCall accepts or declines a pulsing call; cue Notify.
func IncomingCall(c *ui.Context, caller string, accepted, declined *bool, o IncomingCallOptions) ui.Element {
	if caller == "" || accepted == nil || declined == nil {
		panic(fmt.Sprintf("mail: IncomingCall %q needs caller and state", caller))
	}
	def(&o.Accept, "Accept")
	def(&o.Decline, "Decline")
	def(&o.Calling, "Incoming call")
	def(&o.Accepted, "Call accepted")
	def(&o.Declined, "Call declined")
	t := c.Theme()
	root := card(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(5), t.Space(4))
	root.Children(func() {
		done := *accepted || *declined
		rings := ui.Box(c).Size(t.Space(30), t.Space(30)).Center()
		phase := float32(0)
		if !done {
			phase = rings.Loop("ring", 1500*time.Millisecond, ui.Linear)
		}
		rings.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			if done {
				return
			}
			cx, cy := r.X+r.W/2, r.Y+r.H/2
			for _, lag := range []float32{0, 0.5} {
				k := phase + lag
				if k > 1 {
					k--
				}
				rad := t.Space(8) + (r.W/2-t.Space(8))*k
				var ring ui.Path
				ring.Circle(cx, cy, rad)
				p.StrokePath(&ring, 1.5, t.Accent.Alpha(0.6*(1-k)))
			}
		}))
		rings.Children(func() { avatar(c, caller, 16) })
		ui.Text(c, caller).FontSize(t.Rem(1.05)).Bold()
		status := o.Calling
		switch {
		case *accepted:
			status = o.Accepted
		case *declined:
			status = o.Declined
		}
		kit.Muted(c, status).FontSize(t.Rem(0.8)).ContentTransition(ui.ContentTransition{})
		ui.Row(c).Gap(t.Space(12)).Padding(t.Space(3), 0, 0, 0).Children(func() {
			decline := roundAction(c, o.Decline, "phone-off", t.Danger, t.AccentText).Disabled(done)
			accept := roundAction(c, o.Accept, "phone", t.Success, t.AccentText).Disabled(done)
			if accept.Clicked() && !done {
				*accepted = true
				root.SetChanged()
				sound.Play(sound.Notify)
				kit.Log("IncomingCall", "accepted")
			}
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
