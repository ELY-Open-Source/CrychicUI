package health

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// IntervalsOptions word the timer; Round takes round and rounds.
type IntervalsOptions struct{ Work, Rest, Round, Start, Pause, Reset, Done string }

type intervalState struct {
	round         int
	resting, done bool
	running       bool
	ends          time.Time
	left          time.Duration
}

// Intervals alternates work and rest for rounds; cue Notify.
func Intervals(c *ui.Context, work, rest time.Duration, rounds int, o IntervalsOptions) (root, run, reset ui.Element) {
	if work < time.Second || rest < time.Second || work > time.Hour || rest > time.Hour || rounds < 1 || rounds > 50 {
		panic(fmt.Sprintf("health: Intervals needs work and rest of 1 s to 1 h and 1 to 50 rounds, got %v, %v, %d", work, rest, rounds))
	}
	def(&o.Work, "Work")
	def(&o.Rest, "Rest")
	def(&o.Round, "Round %d of %d")
	def(&o.Start, "Start")
	def(&o.Pause, "Pause")
	def(&o.Reset, "Reset")
	def(&o.Done, "Done")
	t := c.Theme()
	root = card(c).Gap(t.Space(4)).Padding(t.Space(5), t.Space(4)).AlignItems(ui.Center).Role(ui.RoleGroup)
	fresh := func() intervalState { return intervalState{round: 1, left: work} }
	st := ui.Local(root, "state", fresh)
	now := c.Now()
	for st.running && !now.Before(st.ends) {
		at := st.ends
		switch {
		case st.resting:
			st.round, st.resting, st.ends, st.left = st.round+1, false, at.Add(work), work
		case st.round == rounds:
			st.running, st.done, st.left = false, true, 0
		default:
			st.resting, st.ends, st.left = true, at.Add(rest), rest
		}
		sound.Play(sound.Notify)
		kit.Log("Intervals", "phase", "round", st.round, "resting", st.resting, "done", st.done)
	}
	span := work
	if st.resting {
		span = rest
	}
	var secs int64
	if st.running {
		secs = countdown(c, st.ends)
	} else {
		secs = int64((st.left + time.Second - 1) / time.Second)
	}
	ink, phase := t.Danger, o.Work
	switch {
	case st.done:
		ink, phase = t.Success, o.Done
	case st.resting:
		ink, phase = t.Success, o.Rest
	}
	root.Children(func() {
		dial(c, t.Space(40), float32(float64(secs)/span.Seconds()), ink, func() {
			ui.Column(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
				c.Key(phase)
				kit.Arrive(ui.Text(c, phase).FontSize(t.Rem(0.8)).Bold().TextColor(ink).MaxWidth(t.Space(28)).SingleLine(), ui.Pose{Scale: 0.7, Fade: 1}, 0)
				kit.Numeric(c, clock(secs)).Font("monospace").FontSize(t.Rem(1.7)).Bold().Shrink(0).
					ContentTransition(ui.ContentTransition{Numeric: true})
			})
		})
		finished := st.round - 1
		if st.done {
			finished = rounds
		}
		ui.Column(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
			kit.Muted(c, fmt.Sprintf(o.Round, st.round, rounds)).FontSize(t.Rem(0.78)).Shrink(0)
			meter(c, float32(finished)/float32(rounds), t.Success)
		})
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			word, mark := o.Start, "play"
			if st.running {
				word, mark = o.Pause, "pause"
			}
			c.Key("run")
			run = kit.Press(ui.ButtonBase(c).Label(word).Padding(t.Space(2), t.Space(5)).Radius(999).Shrink(0).Disabled(st.done))
			run.Background(run.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(run.Hovered())), ui.Snappy)).TextColor(t.AccentText)
			run.Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					glyph(c, mark, 16)
					ui.Text(c, word).FontSize(t.Rem(0.85)).Bold().Shrink(0)
				})
			})
			if run.Clicked() && !st.done {
				if st.running {
					st.left, st.running = st.ends.Sub(now), false
				} else {
					st.ends, st.running = now.Add(st.left), true
				}
				sound.Play(sound.Tap)
				kit.Log("Intervals", "toggled", "running", st.running)
			}
			c.Key("reset")
			reset = kit.Press(ui.ButtonBase(c).Label(o.Reset).Padding(t.Space(2), t.Space(4)).Radius(999).Border(1, t.Border).Shrink(0))
			reset.Background(reset.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(reset.Hovered())), ui.Snappy))
			reset.Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					glyph(c, "rotate-ccw", 14)
					ui.Text(c, o.Reset).FontSize(t.Rem(0.8)).Bold().Shrink(0)
				})
			})
			if reset.Clicked() {
				*st = fresh()
				sound.Play(sound.Tick)
				kit.Log("Intervals", "reset")
			}
		})
	})
	return root, run, reset
}
