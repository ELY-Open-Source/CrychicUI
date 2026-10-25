package learn

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Lesson is one stop on a path; done lessons lead.
type Lesson struct {
	Title, Icon string
	Minutes     int
	Done        bool
}

// LessonPathOptions word the path; Meta takes number and minutes.
type LessonPathOptions struct{ Meta, Locked string }

// LessonPath winds lessons down a path; cue Select.
func LessonPath(c *ui.Context, lessons []Lesson, picked *int, o LessonPathOptions) ui.Element {
	if picked == nil || len(lessons) == 0 || *picked < -1 || *picked >= len(lessons) {
		panic("learn: LessonPath needs lessons and a pick in range")
	}
	titles := make([]string, len(lessons))
	current := len(lessons)
	for i, l := range lessons {
		if l.Icon == "" || l.Minutes < 0 {
			panic(fmt.Sprintf("learn: LessonPath lesson %q needs an icon and minutes of at least 0", l.Title))
		}
		if l.Done && current < i {
			panic(fmt.Sprintf("learn: LessonPath lesson %q is done after an open one", l.Title))
		}
		if !l.Done && current == len(lessons) {
			current = i
		}
		titles[i] = l.Title
	}
	unique("LessonPath", titles)
	def(&o.Meta, "Lesson %d · %d min")
	def(&o.Locked, "%s, locked")
	t := c.Theme()
	node, rowH := t.Space(12), t.Space(17)
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		canvas := ui.Box(c).FillWidth().Height(float32(len(lessons))*rowH + t.Space(4))
		w := canvas.LayoutBounds().W
		kit.Settle(c, canvas, []float32{w}, func() []float32 { return []float32{canvas.LayoutBounds().W} })
		amp := min(w*0.22, t.Space(18))
		at := func(i int) (float32, float32) {
			return w/2 + amp*float32(math.Sin(float64(i)*math.Pi/2)), t.Space(2) + float32(i)*rowH + rowH/2
		}
		pulse := float32(0)
		if current < len(lessons) {
			pulse = canvas.Loop("pulse", 1400*time.Millisecond, ui.EaseOut)
		}
		trim := canvas.Keyframes("draw", []ui.Keyframe{{Value: 0}, {At: 900 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
		canvas.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			for i := 1; i < len(lessons); i++ {
				x0, y0 := at(i - 1)
				x1, y1 := at(i)
				var seg ui.Path
				mid := (y0 + y1) / 2
				seg.MoveTo(r.X+x0, r.Y+y0).CubeTo(r.X+x0, r.Y+mid, r.X+x1, r.Y+mid, r.X+x1, r.Y+y1)
				ink := t.Border
				if lessons[i].Done || i == current {
					ink = t.Accent.Alpha(0.6)
				}
				span := trim*float32(len(lessons)-1) - float32(i-1)
				if span > 0 {
					p.StrokePath(seg.Trim(0, min(1, span)), t.Space(1), ink)
				}
			}
			if current < len(lessons) {
				x, y := at(current)
				h := node/2 + t.Space(3)*pulse
				p.FillPath(ring(r.X+x, r.Y+y, h), t.Accent.Alpha(0.25*(1-pulse)))
			}
		}))
		canvas.Children(func() {
			for i, l := range lessons {
				x, y := at(i)
				open := i <= current
				ink := t.TextMuted
				switch {
				case l.Done:
					ink = t.Success
				case i == current:
					ink = t.Accent
				}
				label := l.Title
				if !open {
					label = fmt.Sprintf(o.Locked, l.Title)
				}
				c.Key(strconv.Quote(l.Title))
				b := kit.Press(ui.ButtonBase(c).Label(label).Absolute().Left(x-node/2).Top(y-node/2).Size(node, node).Center().
					Radius(999).Checked(*picked == i).Disabled(!open))
				fill := l.Done || i == current
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Mix(ink, unit(fill)), ui.Snappy))
				b.Border(t.Space(0.75), b.AnimateColor("line", t.Background.Mix(t.Text, unit(*picked == i)), ui.Snappy))
				b.TextColor(t.TextMuted.Mix(t.AccentText, unit(fill)))
				kit.Arrive(b, ui.Pose{Scale: 0.5, Fade: 1}, time.Duration(i)*90*time.Millisecond).Children(func() {
					mark := l.Icon
					switch {
					case l.Done:
						mark = "check"
					case !open:
						mark = "lock"
					}
					glyph(c, mark, 22)
				})
				right := x <= w/2
				room := max(t.Space(10), w/2-amp/2-node/2-t.Space(6))
				col := ui.Column(c).Absolute().Top(y - t.Space(4.5)).Width(room).Gap(t.Space(0.5))
				if right {
					col.Left(x + node/2 + t.Space(3))
				} else {
					col.Left(x - node/2 - t.Space(3) - room).AlignItems(ui.End)
				}
				col.Children(func() {
					ui.Text(c, l.Title).FontSize(t.Rem(0.8)).Bold().TextColor(t.Text.Mix(t.TextMuted, unit(!open))).
						MaxWidth(room).SingleLine()
					kit.Muted(c, fmt.Sprintf(o.Meta, i+1, l.Minutes)).FontSize(t.Rem(0.68)).MaxWidth(room).SingleLine()
				})
				if b.Clicked() && open && *picked != i {
					*picked = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("LessonPath", "picked", "lesson", i)
				}
			}
		})
	})
	return root
}
