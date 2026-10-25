package game

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SlidePuzzleOptions word the puzzle; Moves takes a count.
type SlidePuzzleOptions struct{ Moves, Solved, Tile string }

// SlidePuzzle slides tiles into the 0 gap; cues Tick, Success.
func SlidePuzzle(c *ui.Context, tiles *[]int, o SlidePuzzleOptions) ui.Element {
	def(&o.Moves, "%d moves")
	if tiles == nil {
		panic("game: SlidePuzzle needs tiles")
	}
	n := 0
	for n*n < len(*tiles) {
		n++
	}
	if n*n != len(*tiles) || n < 3 || n > 5 {
		panic(fmt.Sprintf("game: SlidePuzzle needs 9, 16 or 25 tiles, got %d", len(*tiles)))
	}
	seen := make([]bool, n*n)
	gap := -1
	for i, v := range *tiles {
		if v < 0 || v >= n*n || seen[v] {
			panic(fmt.Sprintf("game: SlidePuzzle tile %d is %d; tiles must be 0 to %d once each", i, v, n*n-1))
		}
		seen[v] = true
		if v == 0 {
			gap = i
		}
	}
	if !solvable(*tiles, n, gap) {
		panic("game: SlidePuzzle cannot be solved from this order")
	}
	def(&o.Solved, "Solved in %d moves")
	def(&o.Tile, "Tile %d")
	t := c.Theme()
	cell, space := t.Space(15), t.Space(1.5)
	side := float32(n)*cell + float32(n-1)*space
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).AlignItems(ui.Center).Role(ui.RoleGroup)
	moves := ui.Local(root, "moves", func() int { return 0 })
	done := solved(*tiles)
	root.Children(func() {
		c.Key(done)
		if done && *moves > 0 {
			kit.Arrive(kit.Chip(c, fmt.Sprintf(o.Solved, *moves), t.Success), ui.Pose{Scale: 0.7, Fade: 1}, 0)
		} else {
			kit.Numeric(c, fmt.Sprintf(o.Moves, *moves)).FontSize(t.Rem(0.8)).Bold().Shrink(0).
				ContentTransition(ui.ContentTransition{Numeric: true})
		}
		board := ui.Box(c).Size(side, side).Shrink(0).Radius(t.Radius).Background(t.SurfaceHover.Alpha(0.6))
		board.Children(func() {
			for i, v := range *tiles {
				if v == 0 {
					continue
				}
				r, col := i/n, i%n
				near := r == gap/n && (col-gap%n == 1 || gap%n-col == 1) || col == gap%n && (r-gap/n == 1 || gap/n-r == 1)
				c.Key(v)
				b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Tile, v)).Absolute().Size(cell, cell).Center().Radius(t.Radius))
				x := b.AnimateTo("x", float32(col)*(cell+space), ui.Snappy)
				y := b.AnimateTo("y", float32(r)*(cell+space), ui.Snappy)
				b.Left(x).Top(y)
				ink := t.Accent
				if done {
					ink = t.Success
				}
				b.Background(b.AnimateColor("bg", t.Background.Mix(ink.Alpha(0.12), unit(near || done)).Mix(t.SurfaceHover, unit(b.Hovered() && near)), ui.Snappy))
				b.Border(1, b.AnimateColor("line", t.Border.Mix(ink, unit(done)), ui.Snappy))
				b.Children(func() {
					kit.Numeric(c, fmt.Sprint(v)).FontSize(t.Rem(1.1)).Bold().TextColor(t.Text.Mix(ink, unit(done))).Shrink(0)
				})
				if b.Clicked() && near && !done {
					(*tiles)[i], (*tiles)[gap] = 0, v
					*moves++
					root.SetChanged()
					if solved(*tiles) {
						sound.Play(sound.Success)
					} else {
						sound.Play(sound.Tick)
					}
					kit.Log("SlidePuzzle", "slid", "tile", v, "moves", *moves)
				}
			}
		})
	})
	return root
}

// solved says tiles run in order, the gap last.
func solved(tiles []int) bool {
	for i, v := range tiles[:len(tiles)-1] {
		if v != i+1 {
			return false
		}
	}
	return true
}

// solvable checks tile parity against the gap's row.
func solvable(tiles []int, n, gap int) bool {
	inv := 0
	for i, a := range tiles {
		for _, b := range tiles[i+1:] {
			if a != 0 && b != 0 && a > b {
				inv++
			}
		}
	}
	if n%2 == 1 {
		return inv%2 == 0
	}
	fromBottom := n - gap/n
	return (inv+fromBottom)%2 == 1
}
