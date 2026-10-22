package datetime

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PomodoroOptions tune a Pomodoro.
type PomodoroOptions struct {
	// Work and Break default to 25 and 5 minutes.
	Work, Break time.Duration
	// Phase names and button copy.
	Focus, Rest, Start, Pause string
}

// Pomodoro alternates work and break; each switch cues Notify.
func Pomodoro(c *ui.Context, o PomodoroOptions) ui.Element {
	if o.Work == 0 {
		o.Work = 25 * time.Minute
	}
	if o.Break == 0 {
		o.Break = 5 * time.Minute
	}
	if o.Work < 0 || o.Break < 0 {
		panic(fmt.Sprintf("datetime: Pomodoro phases %v and %v must be positive", o.Work, o.Break))
	}
	def(&o.Focus, "Focus")
	def(&o.Rest, "Break")
	def(&o.Start, "Start")
	def(&o.Pause, "Pause")
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).AlignItems(ui.Center).Gap(t.Space(3))
	resting := ui.Local(root, "resting", func() bool { return false })
	running := ui.Local(root, "running", func() bool { return false })
	ends := ui.Local(root, "ends", func() time.Time { return time.Time{} })
	left := ui.Local(root, "left", func() time.Duration { return o.Work })
	now := c.Now()
	if *running {
		*left = max(ends.Sub(now), 0)
		if *left == 0 {
			*resting = !*resting
			*left = o.Work
			if *resting {
				*left = o.Break
			}
			*ends = now.Add(*left)
			sound.Play(sound.Notify)
			root.SetChanged()
			kit.Log("Pomodoro", "switched", "resting", *resting)
		}
		c.After(min(*left, *left%time.Second+time.Millisecond))
	}
	phase, span, tint := o.Focus, o.Work, t.Accent
	if *resting {
		phase, span, tint = o.Rest, o.Break, t.Success
	}
	return root.Children(func() {
		ring(c, float32(*left)/float32(span), tint, func() {
			ui.Column(c).AlignItems(ui.Center).Children(func() {
				kit.Numeric(c, clock(*left)).FontSize(t.Rem(1.75)).FontWeight(600)
				kit.Muted(c, phase).ContentTransition(ui.ContentTransition{})
			})
		})
		if *running {
			if button.Outline(c, o.Pause).Clicked() {
				*running = false
				root.SetChanged()
				kit.Log("Pomodoro", "paused")
			}
		} else if button.Solid(c, o.Start).Clicked() {
			*ends, *running = now.Add(*left), true
			root.SetChanged()
			kit.Log("Pomodoro", "started")
		}
	})
}
