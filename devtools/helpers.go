// Package devtools holds developer and operations components.
package devtools

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

func def(s *string, d string) {
	if *s == "" {
		*s = d
	}
}

func unit(v bool) float32 {
	if v {
		return 1
	}
	return 0
}

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

// glyph is a Lucide icon size DIP square.
func glyph(c *ui.Context, name string, size float32) ui.Element {
	return ui.Icon(c, icon.Must(name)).Size(size, size).Shrink(0)
}

// card is the shared arriving card surface.
func card(c *ui.Context) ui.Element {
	t := c.Theme()
	root := ui.Column(c).FillWidth().Radius(t.Radius*1.5).Border(1, t.Border).Background(t.Background).
		Shadow(0, 1, 3, 0, t.Text.Alpha(0.06))
	return kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
}

// tile is a tinted icon square.
func tile(c *ui.Context, mark string, ink ui.Color, size float32) ui.Element {
	t := c.Theme()
	box := ui.Box(c).Size(size, size).Shrink(0).Center().Radius(t.Radius).Background(ink.Alpha(0.12)).TextColor(ink)
	return box.Children(func() { glyph(c, mark, size*0.5) })
}

// empty is a centered empty-list state.
func empty(c *ui.Context, mark, text string) {
	t := c.Theme()
	col := ui.Column(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(6), t.Space(3))
	kit.Arrive(col, ui.Pose{Scale: 0.9, Fade: 1}, 0).Children(func() {
		tile(c, mark, t.TextMuted, t.Space(10))
		kit.Muted(c, text).FontSize(t.Rem(0.85))
	})
}

// meter is a bar filled to share in ink.
func meter(c *ui.Context, share float32, ink ui.Color) ui.Element {
	t := c.Theme()
	bar := ui.Box(c).FillWidth().Height(t.Space(1.5)).Radius(999).Background(t.SurfaceHover)
	grow := bar.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 600 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	w := grow * bar.AnimateTo("share", max(0, min(1, share)), ui.Smooth)
	col := bar.AnimateColor("ink", ink, ui.Snappy)
	return bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * w, H: r.H}, col, 999)
	}))
}

// track is a switch drawing with a sliding knob.
func track(c *ui.Context, on, hover bool) {
	t := c.Theme()
	b := ui.Box(c).Size(t.Space(9), t.Space(5)).Shrink(0).Radius(999)
	b.Background(b.AnimateColor("track", t.SurfaceHover.Mix(t.Border, 0.6*unit(hover)).Mix(t.Accent, unit(on)), ui.Snappy))
	x := b.AnimateTo("knob", unit(on), ui.Snappy)
	pad := t.Space(0.5)
	b.Draw(func(p *ui.Painter, r ui.Rect) {
		d := r.H - 2*pad
		p.Fill(ui.Rect{X: r.X + pad + x*(r.W-r.H), Y: r.Y + pad, W: d, H: d}, t.Background, d/2)
	})
}

// outline is a bordered action with an icon.
func outline(c *ui.Context, label, text, mark string, ink ui.Color) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius).Shrink(0))
	b.Background(b.AnimateColor("bg", t.Background.Mix(ink.Alpha(0.1), unit(b.Hovered())), ui.Snappy))
	b.Border(1, b.AnimateColor("line", t.Border.Mix(ink, unit(b.Hovered())*0.6), ui.Snappy)).TextColor(ink)
	return b.Children(func() {
		ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
			glyph(c, mark, 14)
			ui.Text(c, text).FontSize(t.Rem(0.78)).Bold().Shrink(0)
		})
	})
}

// Stage is a deploy environment.
type Stage int

const (
	StageDevelopment Stage = iota
	StagePreview
	StageStaging
	StageProduction
)

func stageInk(t *ui.Theme, s Stage) ui.Color {
	return [...]ui.Color{t.Success, t.Accent, t.Warning, t.Danger}[s]
}

func validStage(component string, s Stage) {
	if s < StageDevelopment || s > StageProduction {
		panic(fmt.Sprintf("devtools: %s stage %d is not Development to Production", component, s))
	}
}

// StepState is where a pipeline step stands.
type StepState int

const (
	StepPending StepState = iota
	StepRunning
	StepDone
	StepFailed
)

// stepMark is a step's icon; running spins.
func stepMark(c *ui.Context, s StepState) {
	t := c.Theme()
	c.Key(s)
	switch s {
	case StepRunning:
		g := glyph(c, "loader-circle", 16).TextColor(t.Accent)
		g.Rotate(360 * float32(g.Loop("spin", time.Second, ui.Linear)))
	case StepDone:
		kit.Arrive(glyph(c, "circle-check", 16).TextColor(t.Success), ui.Pose{Scale: 0.4, Fade: 1}, 0)
	case StepFailed:
		kit.Arrive(glyph(c, "circle-x", 16).TextColor(t.Danger), ui.Pose{Scale: 0.4, Fade: 1}, 0)
	default:
		glyph(c, "circle-dashed", 16).TextColor(t.TextMuted)
	}
}

// dot joins two meta texts with a drawn dot.
func dot(c *ui.Context) {
	t := c.Theme()
	ui.Box(c).Size(t.Space(0.75), t.Space(0.75)).Shrink(0).Radius(999).Background(t.TextMuted.Alpha(0.6))
}

// unique panics when items repeat or are blank.
func unique(component string, items []string) {
	seen := map[string]bool{}
	for _, s := range items {
		if strings.TrimSpace(s) == "" || seen[s] {
			panic(fmt.Sprintf("devtools: %s item %q is blank or repeated", component, s))
		}
		seen[s] = true
	}
}

// Level is a log level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// levelInk colors a level.
func levelInk(t *ui.Theme, l Level) ui.Color {
	return [...]ui.Color{t.TextMuted, t.Accent, t.Warning, t.Danger}[l]
}

func validLevel(component string, l Level) {
	if l < LevelDebug || l > LevelError {
		panic(fmt.Sprintf("devtools: %s level %d is not Debug to Error", component, l))
	}
}

// statusInk colors an HTTP status.
func statusInk(t *ui.Theme, code int) ui.Color {
	switch {
	case code >= 500:
		return t.Danger
	case code >= 400:
		return t.Warning
	case code >= 300:
		return t.Accent
	}
	return t.Success
}

// methodInk colors an HTTP method.
func methodInk(t *ui.Theme, method string) ui.Color {
	switch method {
	case "GET":
		return t.Accent
	case "POST":
		return t.Success
	case "DELETE":
		return t.Danger
	}
	return t.Warning
}

// took formats a duration as ms or s.
func took(d time.Duration) string {
	if d < time.Second {
		return strconv.FormatFloat(math.Round(float64(d)/float64(time.Millisecond)*10)/10, 'f', -1, 64) + " ms"
	}
	return strconv.FormatFloat(math.Round(d.Seconds()*100)/100, 'f', -1, 64) + " s"
}

// spark is a line through values scaled into r.
func spark(r ui.Rect, values []float64, lo, hi float64) *ui.Path {
	var p ui.Path
	span := hi - lo
	if span <= 0 {
		span = 1
	}
	for i, v := range values {
		x := r.X
		if len(values) > 1 {
			x += r.W * float32(i) / float32(len(values)-1)
		}
		y := r.Y + r.H - r.H*float32((v-lo)/span)
		if i == 0 {
			p.MoveTo(x, y)
		} else {
			p.LineTo(x, y)
		}
	}
	return &p
}

// bounds is the least and greatest of values.
func bounds(component string, values []float64) (lo, hi float64) {
	if len(values) == 0 {
		panic(fmt.Sprintf("devtools: %s needs values", component))
	}
	lo, hi = values[0], values[0]
	for _, v := range values {
		if !finite(v) {
			panic(fmt.Sprintf("devtools: %s value %v is not finite", component, v))
		}
		lo, hi = min(lo, v), max(hi, v)
	}
	if !finite(hi - lo) {
		panic(fmt.Sprintf("devtools: %s range overflows", component))
	}
	return lo, hi
}
