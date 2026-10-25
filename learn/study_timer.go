package learn

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// StudyTimerOptions set the phases, the goal and the copy.
type StudyTimerOptions struct {
	Focus, Rest                                 time.Duration
	Goal                                        int
	Start, Pause, Skip, Working, Resting, Count string
}

type timerState struct {
	rest, running bool
	ends          time.Time
	left          time.Duration
}

// StudyTimer runs focus and break, counting sessions; cue Notify.
func StudyTimer(c *ui.Context, sessions *int, o StudyTimerOptions) (root, run, skip ui.Element) {
	if o.Focus == 0 {
		o.Focus = 25 * time.Minute
	}
	if o.Rest == 0 {
		o.Rest = 5 * time.Minute
	}
	if o.Goal == 0 {
		o.Goal = 4
	}
	if sessions == nil || *sessions < 0 || o.Focus < time.Second || o.Rest < time.Second || o.Focus > 24*time.Hour || o.Rest > 24*time.Hour || o.Goal < 1 || o.Goal > 12 {
		panic("learn: StudyTimer needs sessions of at least 0, phases of 1 s to 24 h and a goal of 1 to 12")
	}
	def(&o.Start, "Start")
	def(&o.Pause, "Pause")
	def(&o.Skip, "Skip")
	def(&o.Working, "Focus")
	def(&o.Resting, "Break")
	def(&o.Count, "%d of %d sessions")
	t := c.Theme()
	root = card(c).Gap(t.Space(4)).Padding(t.Space(5), t.Space(4)).AlignItems(ui.Center).Role(ui.RoleGroup)
	st := ui.Local(root, "state", func() timerState { return timerState{left: o.Focus} })
	now := c.Now()
	// next ends the phase; a focus counts only when finished.
	next := func(finished bool, at time.Time) {
		if st.rest {
			*st = timerState{left: o.Focus}
			return
		}
		if finished && *sessions < math.MaxInt {
			*sessions++
			root.SetChanged()
		}
		*st = timerState{rest: true, running: st.running, ends: at.Add(o.Rest), left: o.Rest}
	}
	for st.running && !now.Before(st.ends) {
		next(true, st.ends)
		sound.Play(sound.Notify)
		kit.Log("StudyTimer", "phase", "rest", st.rest, "sessions", *sessions)
	}
	span := o.Focus
	if st.rest {
		span = o.Rest
	}
	var secs int64
	if st.running {
		secs = countdown(c, st.ends)
	} else {
		secs = int64((st.left + time.Second - 1) / time.Second)
	}
	share := float32(float64(secs) / span.Seconds())
	ink, phase, mark := t.Accent, o.Working, "brain"
	if st.rest {
		ink, phase, mark = t.Success, o.Resting, "coffee"
	}
	root.Children(func() {
		dial(c, t.Space(40), share, ink, func() {
			ui.Column(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
				c.Key(st.rest)
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).TextColor(ink).Children(func() {
					kit.Arrive(glyph(c, mark, 14), ui.Pose{Scale: 0.5, Fade: 1}, 0)
					ui.Text(c, phase).FontSize(t.Rem(0.74)).Bold().MaxWidth(t.Space(26)).SingleLine()
				})
				kit.Numeric(c, clock(secs)).Font("monospace").FontSize(t.Rem(1.7)).Bold().Shrink(0).
					ContentTransition(ui.ContentTransition{Numeric: true})
			})
		})
		ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
			for i := range o.Goal {
				on := i < *sessions
				c.Key(i)
				d := ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Shrink(0).Radius(999)
				d.Background(d.AnimateColor("bg", t.SurfaceHover.Mix(t.Accent, unit(on)), ui.Snappy))
				k := d.AnimateTo("pop", 1+0.4*unit(on), ui.Bouncy)
				d.Scale(k, k)
			}
		})
		kit.Muted(c, fmt.Sprintf(o.Count, min(*sessions, o.Goal), o.Goal)).FontSize(t.Rem(0.76)).Shrink(0)
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			word, glyphName := o.Start, "play"
			if st.running {
				word, glyphName = o.Pause, "pause"
			}
			c.Key("run")
			run = kit.Press(ui.ButtonBase(c).Label(word).Padding(t.Space(2), t.Space(5)).Radius(999).Shrink(0))
			run.Background(run.AnimateColor("bg", ink.Mix(t.AccentHover, 0.3*unit(run.Hovered())), ui.Snappy)).TextColor(t.AccentText)
			run.Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					glyph(c, glyphName, 16)
					ui.Text(c, word).FontSize(t.Rem(0.85)).Bold().Shrink(0)
				})
			})
			if run.Clicked() {
				if st.running {
					st.left, st.running = st.ends.Sub(now), false
				} else {
					st.ends, st.running = now.Add(st.left), true
				}
				sound.Play(sound.Tap)
				kit.Log("StudyTimer", "toggled", "running", st.running)
			}
			c.Key("skip")
			skip = kit.Press(ui.ButtonBase(c).Label(o.Skip).Padding(t.Space(2), t.Space(4)).Radius(999).Border(1, t.Border).Shrink(0))
			skip.Background(skip.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(skip.Hovered())), ui.Snappy))
			skip.Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					glyph(c, "skip-forward", 14)
					ui.Text(c, o.Skip).FontSize(t.Rem(0.8)).Bold().Shrink(0)
				})
			})
			if skip.Clicked() {
				next(false, now)
				sound.Play(sound.Tick)
				kit.Log("StudyTimer", "skipped", "rest", st.rest)
			}
		})
	})
	return root, run, skip
}
