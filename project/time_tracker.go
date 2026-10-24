package project

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TimeTrackerOptions are the tracker labels; Task names the work.
type TimeTrackerOptions struct{ Start, Stop, Task, Running, Paused string }

// TimeTracker starts and stops with total; cue Toggle.
func TimeTracker(c *ui.Context, running *bool, total *time.Duration, o TimeTrackerOptions) ui.Element {
	if running == nil || total == nil || *total < 0 {
		panic("project: TimeTracker needs state")
	}
	def(&o.Start, "Start")
	def(&o.Stop, "Stop")
	def(&o.Task, "Time tracked")
	def(&o.Running, "Running")
	def(&o.Paused, "Paused")
	t := c.Theme()
	root := card(c).Row().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(3))
	started := ui.Local(root, "started", func() time.Time { return time.Time{} })
	wasRunning := ui.Local(root, "wasRunning", func() bool { return *running })
	if *running != *wasRunning {
		if *running {
			*started = c.Now()
		}
		*wasRunning = *running
	}
	if *running && (*started).IsZero() {
		*started = c.Now()
	}
	// Seconds add without overflow; Durations near the limit would not.
	secs := int64(*total / time.Second)
	if *running {
		elapsed := c.Now().Sub(*started)
		frac := *total%time.Second + elapsed%time.Second
		secs += int64(elapsed/time.Second) + int64(frac/time.Second)
		c.After(time.Second - frac%time.Second)
	}
	ink := t.TextMuted.Mix(t.Success, unit(*running))
	root.Children(func() {
		ui.Box(c).Size(t.Space(10), t.Space(10)).Shrink(0).Center().Radius(999).
			Background(root.AnimateColor("tile", ink.Alpha(0.14), ui.Snappy)).TextColor(root.AnimateColor("tint", ink, ui.Snappy)).
			Children(func() { glyph(c, "timer", 18) })
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, o.Task).FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0).SingleLine()
				state := o.Paused
				if *running {
					state = o.Running
				}
				c.Key(*running)
				kit.Arrive(tag(c, state, ink), ui.Pose{Scale: 0.8, Fade: 1}, 0)
			})
			ui.Text(c, fmt.Sprintf("%d:%02d:%02d", secs/3600, secs/60%60, secs%60)).Font("monospace").FontSize(t.Rem(1.4)).Bold()
		})
		name, icon, fill := o.Start, "play", t.Accent
		if *running {
			name, icon, fill = o.Stop, "square", t.Danger
		}
		c.Key("toggle")
		b := kit.Press(ui.ButtonBase(c).Label(name).Size(t.Space(10), t.Space(10)).Shrink(0).Center().Radius(999).TextColor(t.AccentText))
		b.Background(b.AnimateColor("bg", fill.Mix(t.Text, 0.12*unit(b.Hovered())), ui.Snappy))
		b.Children(func() {
			c.Key(icon)
			kit.Arrive(glyph(c, icon, 16), ui.Pose{Scale: 0.5, Rotate: -45, Fade: 1}, 0)
		})
		if b.Clicked() {
			if *running {
				elapsed := c.Now().Sub(*started)
				if elapsed > time.Duration(math.MaxInt64)-*total {
					panic("project: TimeTracker total overflow")
				}
				*total += elapsed
				*running = false
			} else {
				*started = c.Now()
				*running = true
			}
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("TimeTracker", "toggled", "running", *running)
		}
	})
	return root
}
