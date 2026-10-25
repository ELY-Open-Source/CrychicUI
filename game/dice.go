package game

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DiceOptions word the roller; Rand returns 0 to n-1.
type DiceOptions struct {
	Roll, Total string
	Rand        func(n int) int
}

// pips place each face's dots on a 3×3 grid.
var pips = [7][][2]int{{}, {{1, 1}}, {{0, 0}, {2, 2}}, {{0, 0}, {1, 1}, {2, 2}}, {{0, 0}, {2, 0}, {0, 2}, {2, 2}},
	{{0, 0}, {2, 0}, {1, 1}, {0, 2}, {2, 2}}, {{0, 0}, {2, 0}, {0, 1}, {2, 1}, {0, 2}, {2, 2}}}

// Dice rolls one to six dice that tumble; cue Drop.
func Dice(c *ui.Context, faces *[]int, o DiceOptions) ui.Element {
	def(&o.Roll, "Roll")
	if faces == nil || len(*faces) == 0 || len(*faces) > 6 {
		panic(fmt.Sprintf("game: Dice %q needs one to six dice", o.Roll))
	}
	total := 0
	for i, f := range *faces {
		if f < 1 || f > 6 {
			panic(fmt.Sprintf("game: Dice %q die %d shows %d, not 1 to 6", o.Roll, i, f))
		}
		total += f
	}
	def(&o.Total, "Total %d")
	if o.Rand == nil {
		o.Rand = pick
	}
	t := c.Theme()
	root := card(c).Gap(t.Space(4)).Padding(t.Space(5), t.Space(4)).AlignItems(ui.Center).Role(ui.RoleGroup)
	rolls := ui.Local(root, "rolls", func() int { return 0 })
	root.Children(func() {
		ui.Row(c).Wrap().Gap(t.Space(3)).Justify(ui.Center).Children(func() {
			for i, f := range *faces {
				// A new key per roll replays the tumble.
				c.Key([2]int{i, *rolls})
				die(c, f, time.Duration(i)*80*time.Millisecond)
			}
		})
		kit.Numeric(c, fmt.Sprintf(o.Total, total)).FontSize(t.Rem(0.95)).Bold().Shrink(0).
			ContentTransition(ui.ContentTransition{Numeric: true})
		b := kit.Press(ui.ButtonBase(c).Label(o.Roll).Padding(t.Space(2), t.Space(5)).Radius(999).Shrink(0))
		b.Background(b.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(b.Hovered())), ui.Snappy)).TextColor(t.AccentText)
		b.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				glyph(c, "dices", 16)
				ui.Text(c, o.Roll).FontSize(t.Rem(0.85)).Bold().Shrink(0)
			})
		})
		if b.Clicked() {
			for i := range *faces {
				n := o.Rand(6)
				if n < 0 || n > 5 {
					panic(fmt.Sprintf("game: Dice %q Rand gave %d for 6 faces", o.Roll, n))
				}
				(*faces)[i] = n + 1
			}
			*rolls++
			root.SetChanged()
			sound.Play(sound.Drop)
			kit.Log("Dice", "rolled", "dice", len(*faces))
		}
	})
	return root
}

// die is one face that tumbles in after delay.
func die(c *ui.Context, f int, delay time.Duration) {
	t := c.Theme()
	size := t.Space(15)
	d := ui.Box(c).Size(size, size).Shrink(0).Radius(t.Radius*1.4).Background(t.Background).Border(1, t.Border).
		Shadow(0, 1, 3, 0, t.Text.Alpha(0.08)).Label(fmt.Sprint(f))
	k := d.Keyframes("tumble", []ui.Keyframe{{Value: 1}, {At: 650 * time.Millisecond, Value: 0, Ease: ui.EaseOut}},
		ui.KeyframeOptions{Delay: delay})
	d.Translate(0, -t.Space(6)*k).Rotate(-540*k).Scale(1+0.15*k, 1+0.15*k)
	d.Draw(func(p *ui.Painter, r ui.Rect) {
		cell := r.W / 4
		for _, at := range pips[f] {
			x, y := r.X+cell*float32(at[0]+1), r.Y+cell*float32(at[1]+1)
			p.FillPath(ring(x, y, cell*0.32), t.Text)
		}
	})
}
