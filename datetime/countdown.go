package datetime

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CountdownOptions tune a Countdown.
type CountdownOptions struct {
	// Start, Pause and Reset are button copy.
	Start, Pause, Reset string
}

// Countdown runs a ring down to zero, then cues Warning.
func Countdown(c *ui.Context, total time.Duration, o CountdownOptions) ui.Element {
	if total <= 0 {
		panic(fmt.Sprintf("datetime: Countdown total %v must be positive", total))
	}
	def(&o.Start, "Start")
	def(&o.Pause, "Pause")
	def(&o.Reset, "Reset")
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).AlignItems(ui.Center).Gap(t.Space(3))
	running := ui.Local(root, "running", func() bool { return false })
	ends := ui.Local(root, "ends", func() time.Time { return time.Time{} })
	left := ui.Local(root, "left", func() time.Duration { return total })
	now := c.Now()
	if *running {
		*left = max(ends.Sub(now), 0)
		if *left == 0 {
			*running = false
			sound.Play(sound.Warning)
			root.SetChanged()
			kit.Log("Countdown", "done")
		} else {
			c.After(min(*left, *left%time.Second+time.Millisecond))
		}
	}
	return root.Children(func() {
		ring(c, float32(*left)/float32(total), t.Accent, func() {
			kit.Numeric(c, clock(*left)).FontSize(t.Rem(1.75)).FontWeight(600)
		})
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			if *running {
				if button.Outline(c, o.Pause).Clicked() {
					*running = false
					root.SetChanged()
					kit.Log("Countdown", "paused")
				}
			} else if button.Solid(c, o.Start).Disabled(*left == 0).Clicked() {
				*ends, *running = now.Add(*left), true
				root.SetChanged()
				kit.Log("Countdown", "started")
			}
			if button.Ghost(c, o.Reset).Disabled(*left == total).Clicked() {
				*left, *running = total, false
				root.SetChanged()
				kit.Log("Countdown", "reset")
			}
		})
	})
}

// clock writes d as minutes and seconds, rounding up.
func clock(d time.Duration) string {
	s := int64((d + time.Second - 1) / time.Second)
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}
