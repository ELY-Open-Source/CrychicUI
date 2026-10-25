package learn

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Pair is a left item and its matching right item.
type Pair struct{ Left, Right string }

// MatchPairsOptions word the board; Picks, if set, counts clicks.
type MatchPairsOptions struct {
	Title, Count, Done string
	Picks              *int
}

// MatchPairs links left to right; lines draw; cue Success.
func MatchPairs(c *ui.Context, pairs []Pair, order []int, matched *[]bool, o MatchPairsOptions) ui.Element {
	def(&o.Title, "Match the pairs")
	n := len(pairs)
	if matched == nil || n < 2 || n > 8 || len(order) != n || len(*matched) != n || o.Picks != nil && *o.Picks < 0 {
		panic(fmt.Sprintf("learn: MatchPairs %q needs 2 to 8 pairs, an order and a matched flag each", o.Title))
	}
	lefts, rights := make([]string, n), make([]string, n)
	for i, p := range pairs {
		lefts[i], rights[i] = p.Left, p.Right
	}
	unique("MatchPairs", lefts)
	unique("MatchPairs", rights)
	sorted := slices.Sorted(slices.Values(order))
	for i, v := range sorted {
		if v != i {
			panic(fmt.Sprintf("learn: MatchPairs %q order %v is not a permutation", o.Title, order))
		}
	}
	def(&o.Count, "%d of %d matched")
	def(&o.Done, "All matched")
	t := c.Theme()
	rowH, gap := t.Space(11), t.Space(2)
	place := make([]int, n)
	for j, i := range order {
		place[i] = j
	}
	done := 0
	for _, m := range *matched {
		if m {
			done++
		}
	}
	root := card(c).Role(ui.RoleGroup)
	// Picks are held by text, so a reorder keeps them.
	held := ui.Local(root, "sel", func() [2]string { return [2]string{} })
	missed := ui.Local(root, "miss", func() [2]string { return [2]string{} })
	find := func(k int, text string) int {
		for i, p := range pairs {
			if text != "" && [2]string{p.Left, p.Right}[k] == text {
				return i
			}
		}
		return -1
	}
	sel := [2]int{find(0, held[0]), find(1, held[1])}
	miss := [2]int{find(0, missed[0]), find(1, missed[1])}
	missAt := ui.Local(root, "at", func() time.Time { return time.Time{} })
	now := c.Now()
	flash := now.Sub(*missAt) < 450*time.Millisecond
	if flash {
		c.After(450*time.Millisecond - now.Sub(*missAt))
	}
	for k := range sel {
		if sel[k] < 0 || (*matched)[sel[k]] {
			sel[k], held[k] = -1, ""
		}
	}
	// pick toggles side k to pair i, then judges both.
	pick := func(k, i int) {
		if sel[k] == i {
			sel[k], held[k] = -1, ""
		} else {
			sel[k], held[k] = i, [2]string{pairs[i].Left, pairs[i].Right}[k]
		}
		if o.Picks != nil && *o.Picks < math.MaxInt {
			*o.Picks++
			root.SetChanged()
		}
		if sel[0] < 0 || sel[1] < 0 {
			sound.Play(sound.Tap)
			kit.Log("MatchPairs", "picked", "side", k, "pair", sel[k])
			return
		}
		right := sel[0] == sel[1]
		if right {
			(*matched)[sel[0]] = true
			root.SetChanged()
			sound.Play(sound.Success)
		} else {
			*missed, *missAt = *held, now
			sound.Play(sound.Error)
		}
		kit.Log("MatchPairs", "linked", "left", sel[0], "right", place[sel[1]], "match", right)
		*held = [2]string{}
	}
	root.Children(func() {
		kit.CardHead(c, "link", o.Title, func() {
			c.Key(done == n)
			if done == n {
				kit.Arrive(kit.Chip(c, o.Done, t.Success), ui.Pose{Scale: 0.7, Fade: 1}, 0)
			} else {
				kit.Chip(c, fmt.Sprintf(o.Count, done, n), t.Accent)
			}
		})
		ui.Row(c).FillWidth().Padding(t.Space(3)).Children(func() {
			item := func(label string, on, good, bad bool, k int) ui.Element {
				ink := t.Accent
				switch {
				case good:
					ink = t.Success
				case bad:
					ink = t.Danger
				}
				lit := on || good || bad
				b := kit.Press(ui.ButtonBase(c).Label(label).FillWidth().Height(rowH).Padding(0, t.Space(3)).Radius(t.Radius).Disabled(good))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(ink.Alpha(0.12), unit(lit)), ui.Snappy))
				b.Border(1, b.AnimateColor("line", t.Border.Mix(ink, unit(lit)), ui.Snappy))
				return kit.Arrive(b, ui.Pose{Y: t.Space(1.5), Fade: 1}, time.Duration(k)*50*time.Millisecond).Children(func() {
					ui.Text(c, label).FontSize(t.Rem(0.82)).Bold().FillWidth().MinWidth(0).SingleLine().TextColor(t.Text.Mix(ink, 0.5*unit(lit)))
				})
			}
			ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(gap).Children(func() {
				for i, p := range pairs {
					c.Key(strconv.Quote(p.Left))
					b := item(p.Left, sel[0] == i, (*matched)[i], flash && miss[0] == i, i)
					if b.Clicked() && !(*matched)[i] {
						pick(0, i)
					}
				}
			})
			canvas := ui.Box(c).Width(t.Space(14)).Shrink(0).Height(float32(n)*rowH + float32(n-1)*gap)
			lines := make([]float32, n)
			for i, m := range *matched {
				if m {
					lines[i] = canvas.Keyframes(pairs[i], []ui.Keyframe{{Value: 0}, {At: 400 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
				}
			}
			canvas.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				for i, k := range lines {
					if k <= 0 {
						continue
					}
					y0 := r.Y + float32(i)*(rowH+gap) + rowH/2
					y1 := r.Y + float32(place[i])*(rowH+gap) + rowH/2
					var link ui.Path
					link.MoveTo(r.X, y0).CubeTo(r.X+r.W/2, y0, r.X+r.W/2, y1, r.X+r.W, y1)
					p.StrokePath(link.Trim(0, k), t.Space(0.5), t.Success)
				}
			}))
			ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(gap).Children(func() {
				for j, i := range order {
					c.Key(strconv.Quote(pairs[i].Right))
					b := item(pairs[i].Right, sel[1] == i, (*matched)[i], flash && miss[1] == i, j)
					if b.Clicked() && !(*matched)[i] {
						pick(1, i)
					}
				}
			})
		})
	})
	return root
}
