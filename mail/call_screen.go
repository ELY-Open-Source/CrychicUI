package mail

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CallScreenOptions word the end button and status lines.
type CallScreenOptions struct{ End, Live, Ended string }

// CallScreen shows caller and timer; end cues Close.
func CallScreen(c *ui.Context, caller string, started time.Time, ended *bool, o CallScreenOptions) ui.Element {
	if caller == "" || started.IsZero() || started.After(c.Now()) || ended == nil {
		panic(fmt.Sprintf("mail: CallScreen %q needs a past start and state", caller))
	}
	def(&o.End, "End")
	def(&o.Live, "Connected")
	def(&o.Ended, "Call ended")
	t := c.Theme()
	paper, ink := t.Text, t.Background
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(6), t.Space(4)).Radius(t.Radius).
		Background(paper).TextColor(ink).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	frozen := ui.Local(root, "elapsed", func() time.Duration { return c.Now().Sub(started).Truncate(time.Second) })
	wasEnded := ui.Local(root, "wasEnded", func() bool { return *ended })
	elapsed := c.Now().Sub(started).Truncate(time.Second)
	if *ended && !*wasEnded {
		*frozen = elapsed
	}
	*wasEnded = *ended
	if *ended {
		elapsed = *frozen
	}
	if !*ended {
		c.After(time.Second - c.Now().Sub(started)%time.Second)
	}
	root.Children(func() {
		halo := ui.Box(c).Size(t.Space(28), t.Space(28)).Center()
		pulse := float32(0)
		if !*ended {
			pulse = halo.Loop("halo", 1600*time.Millisecond, ui.EaseOut)
		}
		halo.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			if *ended {
				return
			}
			cx, cy := r.X+r.W/2, r.Y+r.H/2
			rad := t.Space(10) + (r.W/2-t.Space(10))*pulse
			p.Fill(ui.Rect{X: cx - rad, Y: cy - rad, W: 2 * rad, H: 2 * rad}, t.Success.Alpha(0.35*(1-pulse)), rad)
		}))
		halo.Children(func() { avatar(c, caller, 20) })
		ui.Text(c, caller).FontSize(t.Rem(1.1)).Bold()
		status, mark := o.Live, t.Success
		if *ended {
			status, mark = o.Ended, t.Danger
		}
		ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(t.Space(1.5), t.Space(1.5)).Radius(999).Background(mark)
			ui.Text(c, status).FontSize(t.Rem(0.78)).TextColor(ink.Alpha(0.75)).ContentTransition(ui.ContentTransition{})
		})
		kit.Numeric(c, clock(elapsed)).FontSize(t.Rem(1.4)).Bold().Padding(t.Space(1), 0, t.Space(3), 0)
		end := roundAction(c, o.End, "phone-off", t.Danger, t.AccentText).Disabled(*ended)
		if end.Clicked() && !*ended {
			*ended = true
			root.SetChanged()
			sound.Play(sound.Close)
			kit.Log("CallScreen", "ended", "seconds", int64(elapsed.Seconds()))
		}
	})
	return root
}

// clock formats a length as mm:ss or h:mm:ss.
func clock(d time.Duration) string {
	s := int64(d / time.Second)
	if s >= 3600 {
		return fmt.Sprintf("%d:%02d:%02d", s/3600, s/60%60, s%60)
	}
	return fmt.Sprintf("%02d:%02d", s/60, s%60)
}
