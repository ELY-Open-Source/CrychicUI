package health

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BreatheOptions set the phase lengths and their words.
type BreatheOptions struct {
	Inhale, Hold, Exhale              time.Duration
	In, Held, Out, Ready, Start, Stop string
}

// Breathe grows and shrinks a circle with breath; cue Tap.
func Breathe(c *ui.Context, running *bool, o BreatheOptions) ui.Element {
	if o.Inhale == 0 {
		o.Inhale = 4 * time.Second
	}
	if o.Hold == 0 {
		o.Hold = 4 * time.Second
	}
	if o.Exhale == 0 {
		o.Exhale = 6 * time.Second
	}
	for _, d := range []time.Duration{o.Inhale, o.Hold, o.Exhale} {
		if d < time.Second || d > time.Minute {
			panic(fmt.Sprintf("health: Breathe phases need 1 to 60 s, got %v", d))
		}
	}
	if running == nil {
		panic("health: Breathe needs a running state")
	}
	def(&o.In, "Breathe in")
	def(&o.Held, "Hold")
	def(&o.Out, "Breathe out")
	def(&o.Ready, "Ready when you are")
	def(&o.Start, "Start")
	def(&o.Stop, "Stop")
	t := c.Theme()
	total := o.Inhale + o.Hold + o.Exhale
	root := card(c).Gap(t.Space(4)).Padding(t.Space(5), t.Space(4)).AlignItems(ui.Center).Role(ui.RoleGroup)
	since := ui.Local(root, "since", func() time.Time { return c.Now() })
	starts := ui.Local(root, "starts", func() int { return 0 })
	was := ui.Local(root, "was", func() bool { return false })
	now := c.Now()
	if *running && !*was {
		*since = now
		*starts++
	}
	*was = *running
	word, left := o.Ready, time.Duration(0)
	if *running {
		in := now.Sub(*since) % total
		switch {
		case in < o.Inhale:
			word, left = o.In, o.Inhale-in
		case in < o.Inhale+o.Hold:
			word, left = o.Held, o.Inhale+o.Hold-in
		default:
			word, left = o.Out, total-in
		}
		frac := left % time.Second
		if frac == 0 {
			frac = time.Second
		}
		c.After(frac + time.Millisecond)
	}
	root.Children(func() {
		size := t.Space(44)
		orb := ui.Box(c).Size(size, size).Shrink(0).Center()
		// One sequence key: a stop eases home.
		frames := []ui.Keyframe{{Value: 0.55}, {At: 600 * time.Millisecond, Value: 0.55, Ease: ui.EaseOut}}
		run := ui.KeyframeOptions{Run: uint64(2**starts + 1)}
		if *running {
			frames = []ui.Keyframe{{Value: 0.55}, {At: o.Inhale, Value: 1, Ease: ui.EaseInOut},
				{At: o.Inhale + o.Hold, Value: 1}, {At: total, Value: 0.55, Ease: ui.EaseInOut}}
			run = ui.KeyframeOptions{Run: uint64(2 * *starts), Repeat: ui.RepeatForever}
		}
		k := orb.Keyframes("breath", frames, run)
		orb.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			cx, cy, rr := r.X+r.W/2, r.Y+r.H/2, r.W/2
			p.FillPath(ring(cx, cy, rr*k), t.Accent.Alpha(0.1))
			p.FillPath(ring(cx, cy, rr*k*0.78), t.Accent.Alpha(0.16))
			p.FillPath(ring(cx, cy, rr*k*0.56), t.Accent.Alpha(0.85))
		}))
		orb.Children(func() {
			if *running {
				secs := int64((left + time.Second - 1) / time.Second)
				kit.Numeric(c, fmt.Sprint(secs)).FontSize(t.Rem(1.3)).Bold().TextColor(t.AccentText).Shrink(0).
					ContentTransition(ui.ContentTransition{Numeric: true})
			}
		})
		// Own scope: caller words never meet the button key.
		ui.Column(c).FillWidth().Children(func() {
			c.Key(word)
			kit.Arrive(ui.Text(c, word).FontSize(t.Rem(0.95)).Bold().FillWidth().TextAlign(ui.Center).Wrap(), ui.Pose{Y: t.Space(1), Fade: 1}, 0)
		})
		label, mark := o.Start, "play"
		if *running {
			label, mark = o.Stop, "square"
		}
		c.Key("toggle")
		b := kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(2), t.Space(5)).Radius(999).Border(1, t.Border).Shrink(0))
		b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
		b.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				glyph(c, mark, 14)
				ui.Text(c, label).FontSize(t.Rem(0.82)).Bold().Shrink(0)
			})
		})
		if b.Clicked() {
			*running = !*running
			root.SetChanged()
			sound.Play(sound.Tap)
			kit.Log("Breathe", "running", "on", *running)
		}
	})
	return root
}
