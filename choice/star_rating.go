package choice

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// StarRatingOptions tune a StarRating.
type StarRatingOptions struct {
	// Max is the number of stars; 5 when zero.
	Max int
}

var bounce = []ui.Keyframe{{Value: 1}, {At: 120 * time.Millisecond, Value: 1.3, Ease: ui.EaseOut}, {At: 360 * time.Millisecond, Value: 1, Ease: ui.EaseInOut}}

// StarRating is previewing stars; the chosen star bounces.
func StarRating(c *ui.Context, label string, value *int, o StarRatingOptions) ui.Element {
	if o.Max == 0 {
		o.Max = 5
	}
	if label == "" || value == nil || o.Max < 1 || *value < 0 || *value > o.Max {
		panic("choice: StarRating needs a label and a value from 0 to Max")
	}
	t := c.Theme()
	root := ui.Row(c).Gap(t.Space(0.5)).AlignSelf(ui.Start).Focusable().Label(label).Role(ui.RoleSlider)
	runs := ui.Local(root, "runs", func() uint64 { return 0 })
	set := func(v int) {
		v = min(o.Max, max(0, v))
		if v != *value {
			*value = v
			*runs++
			root.SetChanged()
			sound.Play(sound.Select)
			kit.Log("StarRating", "rated", "value", v)
		}
	}
	forth, back := ui.KeyRight, ui.KeyLeft
	if c.Direction() == ui.RightToLeft {
		forth, back = back, forth
	}
	switch {
	case root.Shortcut(0, forth), root.Shortcut(0, ui.KeyUp):
		set(*value + 1)
	case root.Shortcut(0, back), root.Shortcut(0, ui.KeyDown):
		set(*value - 1)
	case root.Shortcut(0, ui.KeyHome):
		set(0)
	case root.Shortcut(0, ui.KeyEnd):
		set(o.Max)
	}
	root.Range(0, float64(o.Max), float64(*value)).Step(1)
	hover := -1
	return root.Children(func() {
		for i := range o.Max {
			c.Key(i)
			star := ui.Box(c).Size(t.Rem(1.5), t.Rem(1.5)).Cursor(ui.CursorPointer)
			if star.Clicked() {
				if *value == i+1 {
					set(0)
				} else {
					set(i + 1)
				}
				root.Focus()
			}
			if star.Hovered() {
				hover = i
			}
			if i == *value-1 && *runs > 0 {
				k := star.Keyframes("bounce", bounce, ui.KeyframeOptions{Run: *runs})
				star.Scale(k, k)
			}
			star.Draw(func(p *ui.Painter, r ui.Rect) {
				filled, col := i < *value, t.Warning
				if hover >= 0 {
					filled, col = i <= hover, t.Warning.Alpha(0.6)
				}
				if filled {
					p.FillPath(starPath(r), col)
				} else {
					p.StrokePath(starPath(r), 1.2, t.TextMuted)
				}
			})
		}
	})
}

// starPath is a five-pointed star in r.
func starPath(r ui.Rect) *ui.Path {
	var path ui.Path
	cx, cy := r.X+r.W/2, r.Y+r.H/2+r.H*0.04
	for i := range 10 {
		rad := r.W * 0.48
		if i%2 == 1 {
			rad = r.W * 0.2
		}
		a := float64(i)*math.Pi/5 - math.Pi/2
		x, y := cx+rad*float32(math.Cos(a)), cy+rad*float32(math.Sin(a))
		if i == 0 {
			path.MoveTo(x, y)
		} else {
			path.LineTo(x, y)
		}
	}
	path.Close()
	return &path
}
