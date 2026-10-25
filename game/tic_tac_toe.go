package game

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Mark is a tic-tac-toe cell.
type Mark int

const (
	Blank Mark = iota
	MarkX
	MarkO
)

// TicTacToeOptions word the game; Taken takes a cell and mark.
type TicTacToeOptions struct {
	Turn, Win, Draw, Reset, Cell, Taken string
	Marks                               [2]string
}

var lines = [8][3]int{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, {0, 3, 6}, {1, 4, 7}, {2, 5, 8}, {0, 4, 8}, {2, 4, 6}}

// TicTacToe draws its win line; cues Pop, Success.
func TicTacToe(c *ui.Context, board *[9]Mark, o TicTacToeOptions) ui.Element {
	if board == nil {
		panic("game: TicTacToe needs a board")
	}
	xs, os := 0, 0
	for i, m := range board {
		switch m {
		case MarkX:
			xs++
		case MarkO:
			os++
		case Blank:
		default:
			panic(fmt.Sprintf("game: TicTacToe cell %d holds %d", i, m))
		}
	}
	if xs != os && xs != os+1 {
		panic(fmt.Sprintf("game: TicTacToe has %d X and %d O; X moves first", xs, os))
	}
	def(&o.Turn, "%s to move")
	def(&o.Win, "%s wins")
	def(&o.Draw, "Draw")
	def(&o.Reset, "New game")
	def(&o.Cell, "Cell %d")
	def(&o.Taken, "Cell %d, %s")
	if o.Marks == [2]string{} {
		o.Marks = [2]string{"X", "O"}
	}
	t := c.Theme()
	won, line := winner(board)
	full := xs+os == 9
	next := MarkX
	if xs > os {
		next = MarkO
	}
	inks := [3]ui.Color{t.TextMuted, t.Accent, t.Danger}
	cell, gap := t.Space(18), t.Space(2)
	side := 3*cell + 2*gap
	root := card(c).Gap(t.Space(4)).Padding(t.Space(5), t.Space(4)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		text, tone := fmt.Sprintf(o.Turn, o.Marks[next-1]), inks[next]
		switch {
		case won != Blank:
			text, tone = fmt.Sprintf(o.Win, o.Marks[won-1]), inks[won]
		case full:
			text, tone = o.Draw, t.TextMuted
		}
		c.Key(text)
		kit.Arrive(ui.Text(c, text).FontSize(t.Rem(1)).Bold().TextColor(tone).Shrink(0), ui.Pose{Y: t.Space(1), Fade: 1}, 0)
		ui.Box(c).Size(side, side).Shrink(0).Children(func() {
			for i, m := range board {
				r, col := i/3, i%3
				x, y := float32(col)*(cell+gap), float32(r)*(cell+gap)
				// Taken cells are not buttons, so stay undimmed.
				if m != Blank {
					c.Key(takenCell(i))
					ui.Box(c).Absolute().Left(x).Top(y).Size(cell, cell).Radius(t.Radius).Background(t.SurfaceHover.Alpha(0.6)).
						Label(fmt.Sprintf(o.Taken, i+1, o.Marks[m-1])).Children(func() { markBox(c, m, inks[m], cell) })
					continue
				}
				c.Key(i)
				b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Cell, i+1)).Absolute().Left(x).Top(y).Size(cell, cell).
					Radius(t.Radius).Disabled(won != Blank))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(0.6).Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
				if b.Clicked() && won == Blank {
					board[i] = next
					root.SetChanged()
					if w, _ := winner(board); w != Blank {
						sound.Play(sound.Success)
					} else {
						sound.Play(sound.Pop)
					}
					kit.Log("TicTacToe", "moved", "cell", i, "mark", next)
				}
			}
			if won != Blank {
				// Cells take int keys; the line takes its own type.
				c.Key(winLine(line))
				over := ui.Box(c).Absolute().Left(0).Top(0).Size(side, side)
				draw := over.Keyframes("win", []ui.Keyframe{{Value: 0}, {At: 400 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
				a, z := lines[line][0], lines[line][2]
				center := func(i int) (float32, float32) {
					return float32(i%3)*(cell+gap) + cell/2, float32(i/3)*(cell+gap) + cell/2
				}
				over.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
					x0, y0 := center(a)
					x1, y1 := center(z)
					var path ui.Path
					path.MoveTo(r.X+x0, r.Y+y0).LineTo(r.X+x1, r.Y+y1)
					p.StrokePath(path.Trim(0, draw), t.Space(1.25), inks[won].Alpha(0.85))
				}))
			}
		})
		b := kit.Press(ui.ButtonBase(c).Label(o.Reset).Padding(t.Space(1.5), t.Space(4)).Radius(999).Border(1, t.Border).Shrink(0))
		b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
		b.Children(func() { ui.Text(c, o.Reset).FontSize(t.Rem(0.8)).Bold().Shrink(0) })
		if b.Clicked() {
			*board = [9]Mark{}
			root.SetChanged()
			sound.Play(sound.Pop)
			kit.Log("TicTacToe", "reset")
		}
	})
	return root
}

// markBox draws an X or O that draws in.
func markBox(c *ui.Context, m Mark, ink ui.Color, size float32) {
	t := c.Theme()
	box := ui.Box(c).Size(size, size)
	draw := box.Keyframes("draw", []ui.Keyframe{{Value: 0}, {At: 260 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		in := r.W * 0.28
		w := t.Space(1.25)
		if m == MarkO {
			p.StrokePath(ring(r.X+r.W/2, r.Y+r.H/2, r.W/2-in).Trim(0, draw), w, ink)
			return
		}
		var a, b ui.Path
		a.MoveTo(r.X+in, r.Y+in).LineTo(r.X+r.W-in, r.Y+r.H-in)
		b.MoveTo(r.X+r.W-in, r.Y+in).LineTo(r.X+in, r.Y+r.H-in)
		p.StrokePath(a.Trim(0, min(1, draw*2)), w, ink)
		p.StrokePath(b.Trim(0, max(0, draw*2-1)), w, ink)
	})
}

// winner is the mark holding a line and that line.
func winner(board *[9]Mark) (Mark, int) {
	for i, l := range lines {
		if m := board[l[0]]; m != Blank && board[l[1]] == m && board[l[2]] == m {
			return m, i
		}
	}
	return Blank, -1
}

type winLine int

type takenCell int
