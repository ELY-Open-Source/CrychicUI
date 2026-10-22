// Package motion holds entrances, text effects, loops and pointer motion.
package motion

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

var arrive = []ui.Keyframe{{Value: 0}, {At: 200 * time.Millisecond, Value: 1, Ease: ui.EaseOut}}

// replay drops first: a replay starts from the shown value.
func replay(d time.Duration, ease ui.Easing) []ui.Keyframe {
	return []ui.Keyframe{{Value: 0}, {At: time.Millisecond, Value: 0}, {At: d, Value: 1, Ease: ease}}
}

// forever repeats a sequence while built.
var forever = ui.KeyframeOptions{Repeat: ui.RepeatForever}

// still reports that the desktop asks for less motion.
func still(c *ui.Context) bool { return c.Preferences().ReduceMotion }

// onChange plays frames whenever run leaves its first value.
func onChange(e ui.Element, run uint64, frames []ui.Keyframe) float32 {
	first := ui.Local(e, "crychic.first", func() uint64 { return run })
	moved := ui.Local(e, "crychic.moved", func() bool { return false })
	if run != *first {
		*moved = true
	}
	if !*moved {
		return frames[0].Value
	}
	return e.Keyframes("crychic.change", frames, ui.KeyframeOptions{Run: run})
}

// finite panics on NaN or infinite v.
func finite(component string, v float64) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		panic("motion: " + component + " needs finite numbers")
	}
}

// unique panics unless items are non-empty, distinct texts.
func unique(component string, items []string) {
	if len(items) == 0 {
		panic("motion: " + component + " needs items")
	}
	if len(set(items)) != len(items) || set(items)[""] {
		panic("motion: " + component + " items need distinct text")
	}
}

// set holds each item once.
func set(items []string) map[string]bool {
	m := make(map[string]bool, len(items))
	for _, s := range items {
		m[s] = true
	}
	return m
}

// unit converts on to an animation target.
func unit(on bool) float32 {
	if on {
		return 1
	}
	return 0
}

func abs(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

// identity is 0 through n-1 in order.
func identity(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i
	}
	return out
}

// positive checks a size knob; zero means fallback.
func positive(component, field string, v, fallback float32) float32 {
	finite(component, float64(v))
	if v < 0 {
		panic("motion: " + component + " needs a positive " + field)
	}
	if v == 0 {
		return fallback
	}
	return v
}

// glyph is one drawn letter and its x offset.
type glyph struct {
	span ui.Span
	x    float32
}

// letters lays out bold headline text letter by letter.
func letters(c *ui.Context, text string) (gs []glyph, w, h float32) {
	t := c.Theme()
	runes := []rune(text)
	for i, r := range runes {
		x, _ := c.MeasureText(0, ui.Span{Text: string(runes[:i]), Size: t.Rem(1.5), Weight: 700})
		gs = append(gs, glyph{ui.Span{Text: string(r), Size: t.Rem(1.5), Weight: 700, Color: t.Text}, x})
	}
	w, h = c.MeasureText(0, ui.Span{Text: text, Size: t.Rem(1.5), Weight: 700})
	return gs, w, h
}

// along is a box's extent on the marquee's axis.
func along(b ui.Rect, vertical bool) float32 {
	if vertical {
		return b.H
	}
	return b.W
}

// marquee loops copies of content enough to fill the view.
func marquee(c *ui.Context, vertical bool, v float32, content func()) ui.Element {
	t := c.Theme()
	gap := t.Space(8)
	root := ui.Box(c)
	if vertical {
		root.ClipY()
	} else {
		root.ClipX()
	}
	offset := ui.Local(root, "offset", func() float32 { return 0 })
	last := ui.Local(root, "last", func() time.Time { return time.Time{} })
	cycle := ui.Local(root, "cycle", func() float32 { return 0 })
	paused := root.Hovered() || still(c)
	return root.Children(func() {
		track := ui.Row(c)
		if vertical {
			track = ui.Column(c)
		}
		used := []float32{*cycle, along(root.LayoutBounds(), vertical)}
		n := 2
		if used[0] > gap && used[1] > 0 {
			n = max(2, int(math.Ceil(float64(used[1]/used[0])))+1)
		}
		var one ui.Element
		track.Gap(gap).AlignSelf(ui.Start).Shrink(0).Children(func() {
			for i := range n {
				b := ui.Box(c.Key(i)).Children(content)
				if i == 0 {
					one = b
				}
			}
		})
		read := func() []float32 {
			return []float32{along(one.LayoutBounds(), vertical) + gap, along(root.LayoutBounds(), vertical)}
		}
		*cycle = read()[0]
		kit.Settle(c, track, used, read)
		now := c.Now()
		if used[0] > gap {
			if !paused && !last.IsZero() {
				*offset += v * float32(now.Sub(*last).Seconds())
			}
			*offset = float32(math.Mod(float64(*offset), float64(used[0])))
		}
		if paused {
			*last = time.Time{}
		} else {
			*last = now
			c.AnimationFrame()
		}
		switch {
		case vertical:
			track.Translate(0, -*offset)
		case c.Direction() == ui.RightToLeft:
			track.Translate(*offset, 0)
		default:
			track.Translate(-*offset, 0)
		}
	})
}
