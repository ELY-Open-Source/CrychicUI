package datetime

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// StopwatchOptions tune a Stopwatch.
type StopwatchOptions struct {
	// Button copy; Lap formats the lap number in.
	Start, Stop, Lap, Reset string
}

// Stopwatch counts up with start, lap and reset.
func Stopwatch(c *ui.Context, o StopwatchOptions) ui.Element {
	def(&o.Start, "Start")
	def(&o.Stop, "Stop")
	def(&o.Lap, "Lap %d")
	def(&o.Reset, "Reset")
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(3)).MinWidth(t.Space(64))
	running := ui.Local(root, "running", func() bool { return false })
	from := ui.Local(root, "from", func() time.Time { return time.Time{} })
	banked := ui.Local(root, "banked", func() time.Duration { return 0 })
	laps := ui.Local(root, "laps", func() []time.Duration { return nil })
	now := c.Now()
	elapsed := *banked
	if *running {
		elapsed += now.Sub(*from)
		c.AnimationFrame()
	}
	return root.Children(func() {
		kit.Numeric(c, stopwatch(elapsed)).FontSize(t.Rem(2.25)).FontWeight(600)
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			toggle := o.Start
			if *running {
				toggle = o.Stop
			}
			if act(c, toggle, true).Clicked() {
				if *running {
					*banked, *running = elapsed, false
				} else {
					*from, *running = now, true
				}
				sound.Play(sound.Tick)
				root.SetChanged()
				kit.Log("Stopwatch", "running", "on", *running)
			}
			if act(c, fmt.Sprintf(o.Lap, len(*laps)+1), false).Disabled(!*running).Clicked() {
				*laps = append(*laps, elapsed)
				sound.Play(sound.Tick)
				root.SetChanged()
				kit.Log("Stopwatch", "lap", "count", len(*laps))
			}
			if act(c, o.Reset, false).Disabled(*running || elapsed == 0).Clicked() {
				*banked, *laps = 0, nil
				sound.Play(sound.Tick)
				root.SetChanged()
				kit.Log("Stopwatch", "reset")
			}
		})
		for i := len(*laps) - 1; i >= 0; i-- {
			split := (*laps)[i]
			if i > 0 {
				split -= (*laps)[i-1]
			}
			ui.Row(c).Key(i).FillWidth().Transition(ui.ElementTransition{Enter: &ui.Motion{Y: -t.Space(2)}}).Children(func() {
				kit.Muted(c, fmt.Sprintf(o.Lap, i+1)).Grow(1)
				ui.Text(c, stopwatch(split))
			})
		}
	})
}

// stopwatch writes d as minutes, seconds and hundredths.
func stopwatch(d time.Duration) string {
	cs := d.Milliseconds() / 10
	return fmt.Sprintf("%02d:%02d.%02d", cs/6000, cs/100%60, cs%100)
}
