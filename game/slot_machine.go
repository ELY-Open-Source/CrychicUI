package game

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SlotMachineOptions word the machine; Rand returns 0 to n-1.
type SlotMachineOptions struct {
	Spin, Win string
	Rand      func(n int) int
}

// SlotMachine spins three reels; each stop plays Tick.
func SlotMachine(c *ui.Context, symbols []string, reels *[3]int, o SlotMachineOptions) ui.Element {
	def(&o.Spin, "Spin")
	n := len(symbols)
	if reels == nil || n < 3 || n > 10 {
		panic(fmt.Sprintf("game: SlotMachine %q needs 3 to 10 symbols and reels", o.Spin))
	}
	unique("SlotMachine", symbols)
	for i, r := range reels {
		if r < 0 || r >= n {
			panic(fmt.Sprintf("game: SlotMachine %q reel %d shows %d of %d symbols", o.Spin, i, r, n))
		}
	}
	def(&o.Win, "Three of a kind!")
	if o.Rand == nil {
		o.Rand = pick
	}
	t := c.Theme()
	cell := t.Space(16)
	root := card(c).Gap(t.Space(4)).Padding(t.Space(5), t.Space(4)).AlignItems(ui.Center).Role(ui.RoleGroup)
	spins := ui.Local(root, "spins", func() int { return 0 })
	from := ui.Local(root, "from", func() [3]int { return *reels })
	stopped := ui.Local(root, "stopped", func() int { return 3 })
	root.Children(func() {
		frame := ui.Row(c).Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius * 1.5).Background(t.Text.Mix(t.Background, 0.12))
		frame.Children(func() {
			for i := range 3 {
				c.Key([2]int{i, *spins})
				reel(c, symbols, (*from)[i], reels[i], *spins > 0, i, cell, stopped)
			}
		})
		win := *stopped == 3 && reels[0] == reels[1] && reels[1] == reels[2]
		if win {
			c.Key("win")
			kit.Arrive(kit.Chip(c, o.Win, t.Success), ui.Pose{Scale: 0.6, Fade: 1}, 0)
		}
		busy := *stopped < 3
		c.Key("spin")
		b := kit.Press(ui.ButtonBase(c).Label(o.Spin).Padding(t.Space(2), t.Space(5)).Radius(999).Shrink(0).Disabled(busy))
		b.Background(b.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(b.Hovered())), ui.Snappy)).TextColor(t.AccentText)
		b.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				glyph(c, "cherry", 16)
				ui.Text(c, o.Spin).FontSize(t.Rem(0.85)).Bold().Shrink(0)
			})
		})
		if b.Clicked() && !busy {
			*from = *reels
			for i := range reels {
				r := o.Rand(n)
				if r < 0 || r >= n {
					panic(fmt.Sprintf("game: SlotMachine %q Rand gave %d for %d symbols", o.Spin, r, n))
				}
				reels[i] = r
			}
			*spins++
			*stopped = 0
			root.SetChanged()
			kit.Log("SlotMachine", "spun", "reels", *reels)
		}
	})
	return root
}

// reel rolls between symbols; its stop plays Tick.
func reel(c *ui.Context, symbols []string, from, to int, moving bool, i int, cell float32, stopped *int) {
	t := c.Theme()
	n := len(symbols)
	window := ui.Box(c).Size(cell, cell*1.25).Shrink(0).Clip().Radius(t.Radius).Background(t.Background)
	pos := 0.0
	if moving {
		u := float64(window.Keyframes("roll", []ui.Keyframe{{Value: 0}, {At: time.Duration(900+350*i) * time.Millisecond, Value: 1}}))
		travel := float64((2+i)*n + (to-from+n)%n)
		pos = travel * (1 - math.Pow(1-u, 3))
		if u >= 1 && *stopped == i {
			*stopped = i + 1
			sound.Play(sound.Tick)
		}
	}
	whole := int(pos)
	frac := float32(pos - float64(whole))
	top := (from + whole) % n
	if !moving {
		top, frac = to, 0
	}
	window.Children(func() {
		for k := range 2 {
			ui.Box(c).Absolute().Left(0).Top(cell*0.125+(float32(k)-frac)*cell).Size(cell, cell).Center().Children(func() {
				glyph(c, symbols[(top+k)%n], cell*0.55).TextColor(t.Text)
			})
		}
	})
}
