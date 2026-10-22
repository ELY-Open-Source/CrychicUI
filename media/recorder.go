package media

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Take is a Recorder's state.
type Take uint8

const (
	Idle Take = iota
	Recording
	Paused
	takeCount
)

// RecorderOptions tune a Recorder.
type RecorderOptions struct {
	Record, Pause, Resume, Stop string
}

// Recorder records, pauses and stops with level and time.
func Recorder(c *ui.Context, label string, state *Take, level float64, o RecorderOptions) ui.Element {
	if label == "" || state == nil || *state >= takeCount || !(level >= 0 && level <= 1) {
		panic(fmt.Sprintf("media: Recorder %q needs a label, a known state and a level from 0 to 1", label))
	}
	def(&o.Record, "Record")
	def(&o.Pause, "Pause")
	def(&o.Resume, "Resume")
	def(&o.Stop, "Stop")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(4)).Radius(t.Radius).Border(1, t.Border).Label(label)
	now := c.Now()
	banked := ui.Local(root, "banked", func() time.Duration { return 0 })
	since := ui.Local(root, "since", func() time.Time { return now })
	seen := ui.Local(root, "seen", func() Take { return Idle })
	var dot, meter, main, stop ui.Element
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			dot = ui.Box(c).Size(t.Space(3), t.Space(3)).Radius(999).Background(t.Danger)
			meter = ui.Row(c).Grow(1).AlignItems(ui.Center).Gap(t.Space(3))
		})
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			main = kit.Press(ui.ButtonBase(c).Size(t.Space(11), t.Space(11)).Center().Radius(999))
			stop = kit.Press(ui.ButtonBase(c).Label(o.Stop).Tooltip(o.Stop).Size(t.Space(11), t.Space(11)).Center().Radius(999).
				Border(1, t.Border))
		})
	})
	act := func(to Take, event string) {
		*state = to
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("Recorder", event)
	}
	switch {
	case main.Clicked() && *state == Recording:
		act(Paused, "paused")
	case main.Clicked():
		act(Recording, "recording")
	case stop.Clicked() && *state != Idle:
		act(Idle, "stopped")
	}
	if *state != *seen {
		if *seen == Recording {
			*banked += now.Sub(*since)
		}
		if *state == Recording {
			*since = now
		}
		if *state == Idle {
			kit.Log("Recorder", "length", "seconds", banked.Seconds())
			*banked = 0
		}
		*seen = *state
	}
	took := *banked
	on := *state == Recording
	if on {
		took += now.Sub(*since)
		c.After(time.Second - took%time.Second)
	}
	glow := float32(0.3)
	if on {
		glow = 0.4 + 0.6*float32(math.Abs(math.Sin(math.Pi*float64(dot.Loop("pulse", 1200*time.Millisecond, ui.Linear)))))
	}
	dot.Opacity(dot.AnimateTo("glow", glow, ui.Snappy))
	meter.Children(func() {
		kit.Numeric(c, clock(took)).FontSize(t.Rem(1.4)).Bold()
		feedback.LevelMeter(c, label, level*float64(unit(on))).Grow(1)
	})
	name, glyph, fg, bg := o.Record, "circle", t.Danger, t.Danger.Alpha(0.12)
	switch *state {
	case Recording:
		name, glyph, fg, bg = o.Pause, "pause", t.Text, t.SurfaceHover
	case Paused:
		name = o.Resume
	}
	main.Label(name).Tooltip(name).TextColor(main.AnimateColor("fg", fg, ui.Snappy)).
		Background(main.AnimateColor("bg", bg, ui.Snappy)).
		Children(func() { ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{}) })
	stop.Disabled(*state == Idle).Children(func() { ui.Icon(c, icon.Must("square")).Size(16, 16) })
	return root
}
