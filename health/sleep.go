package health

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Stage is a sleep stage.
type Stage int

const (
	Awake Stage = iota
	REM
	Light
	Deep
)

// SleepSpan is a stretch of the night in one stage.
type SleepSpan struct {
	Stage   Stage
	Minutes int
}

// SleepOptions word the chart; Start, if set, labels the axis.
type SleepOptions struct {
	Title, Asleep, Span, Short, Clock string
	Names                             [4]string
	Start                             time.Time
}

// Sleep draws the night's stages as stepped bars.
func Sleep(c *ui.Context, spans []SleepSpan, o SleepOptions) ui.Element {
	def(&o.Title, "Sleep")
	if len(spans) == 0 {
		panic(fmt.Sprintf("health: Sleep %q needs spans", o.Title))
	}
	var each [4]int
	total := 0
	for i, s := range spans {
		if s.Stage < Awake || s.Stage > Deep || s.Minutes < 1 || s.Minutes > 24*60-total {
			panic(fmt.Sprintf("health: Sleep %q span %d needs a stage and 1 or more minutes within a day", o.Title, i))
		}
		each[s.Stage] += s.Minutes
		total += s.Minutes
	}
	def(&o.Asleep, "%s asleep")
	def(&o.Span, "%dh %02dm")
	def(&o.Short, "%dm")
	def(&o.Clock, "15:04")
	for i, n := range [...]string{"Awake", "REM", "Light", "Deep"} {
		def(&o.Names[i], n)
	}
	t := c.Theme()
	inks := [4]ui.Color{t.Warning, t.Accent.Mix(t.Danger, 0.35), t.Accent.Alpha(0.55), t.Accent.Mix(t.Text, 0.35)}
	span := func(m int) string {
		if m < 60 {
			return fmt.Sprintf(o.Short, m)
		}
		return fmt.Sprintf(o.Span, m/60, m%60)
	}
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			kit.Arrive(tile(c, "moon", t.Accent, t.Space(11)), ui.Pose{Scale: 0.6, Rotate: -20, Fade: 1}, 0)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				kit.Muted(c, o.Title).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
				kit.Numeric(c, fmt.Sprintf(o.Asleep, span(total-each[Awake]))).FontSize(t.Rem(1.2)).Bold().FillWidth().MinWidth(0).SingleLine()
			})
		})
		chart := ui.Box(c).FillWidth().Height(t.Space(28))
		k := chart.Keyframes("draw", []ui.Keyframe{{Value: 0}, {At: 900 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
		chart.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			lane := r.H / 4
			for i := range 4 {
				p.Line(r.X, r.Y+lane*float32(i+1), r.X+r.W, r.Y+lane*float32(i+1), 1, t.Border.Alpha(0.6))
			}
			p.Clip(ui.Rect{X: r.X, Y: r.Y, W: r.W * k, H: r.H}, 0, func() {
				at := 0
				for _, s := range spans {
					x := r.X + r.W*float32(at)/float32(total)
					w := float32(math.Max(1, float64(r.W*float32(s.Minutes)/float32(total))))
					p.Fill(ui.Rect{X: x, Y: r.Y + lane*float32(s.Stage) + lane*0.15, W: w, H: lane * 0.7}, inks[s.Stage], t.Space(0.75))
					at += s.Minutes
				}
			})
		}))
		if !o.Start.IsZero() {
			ui.Row(c).FillWidth().Justify(ui.SpaceBetween).Children(func() {
				kit.Muted(c, o.Start.Format(o.Clock)).FontSize(t.Rem(0.66)).Shrink(0)
				kit.Muted(c, o.Start.Add(time.Duration(total)*time.Minute).Format(o.Clock)).FontSize(t.Rem(0.66)).Shrink(0)
			})
		}
		legend := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(2))
		legend.Children(func() {
			for i, name := range o.Names {
				cell := ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).MinWidth(0)
				kit.Stagger(c, legend, cell, i).Children(func() {
					ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Shrink(0).Radius(t.Space(0.75)).Background(inks[i])
					kit.Muted(c, name).FontSize(t.Rem(0.72)).Grow(1).MinWidth(0).SingleLine()
					ui.Text(c, span(each[i])).FontSize(t.Rem(0.76)).Bold().Shrink(0)
				})
			}
		})
	})
	return root
}
