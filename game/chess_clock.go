package game

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ChessClockOptions word the clock; Press takes a side name.
type ChessClockOptions struct {
	Names              [2]string
	Press, Pause, Flag string
}

// ChessClock runs one side's clock; turn -1 pauses; cue Tick.
func ChessClock(c *ui.Context, left *[2]time.Duration, turn *int, o ChessClockOptions) ui.Element {
	if left == nil || turn == nil || *turn < -1 || *turn > 1 || left[0] < 0 || left[1] < 0 {
		panic("game: ChessClock needs non-negative times and a turn of -1, 0 or 1")
	}
	if o.Names == [2]string{} {
		o.Names = [2]string{"White", "Black"}
	}
	def(&o.Press, "End %s's move")
	def(&o.Pause, "Pause")
	def(&o.Flag, "Time")
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(3)).Role(ui.RoleGroup)
	since := ui.Local(root, "since", func() time.Time { return c.Now() })
	last := ui.Local(root, "turn", func() int { return *turn })
	now := c.Now()
	// A caller's turn change keeps the last side's used time.
	if *turn != *last {
		if *last >= 0 {
			left[*last] = max(0, left[*last]-now.Sub(*since))
		}
		*since, *last = now, *turn
	}
	shown := *left
	if *turn >= 0 {
		shown[*turn] = max(0, left[*turn]-now.Sub(*since))
	}
	flagged := *turn >= 0 && shown[*turn] == 0
	if *turn >= 0 && !flagged {
		run := shown[*turn]
		next := run % time.Second
		if run < 10*time.Second {
			next = run % (100 * time.Millisecond)
		}
		if next == 0 {
			next = min(run, 100*time.Millisecond)
		}
		c.After(next + time.Millisecond)
	}
	commit := func() {
		if *turn >= 0 {
			left[*turn] = shown[*turn]
		}
	}
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
			for i := range 2 {
				on := *turn == i
				low := shown[i] < 10*time.Second
				out := shown[i] == 0 && flagged && on
				ink := t.Accent
				switch {
				case out:
					ink = t.Danger
				case low:
					ink = t.Warning
				}
				c.Key(i)
				b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Press, o.Names[i])).Grow(1).Basis(0).MinWidth(0).Height(t.Space(30)).
					Column().Center().Gap(t.Space(1)).Radius(t.Radius * 1.5).Disabled(flagged || *turn == 1-i))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Mix(ink, unit(on)), ui.Snappy))
				b.TextColor(b.AnimateColor("ink", t.Text.Mix(t.AccentText, unit(on)), ui.Snappy))
				b.Children(func() {
					ui.Text(c, o.Names[i]).FontSize(t.Rem(0.8)).Bold().MaxWidth(t.Space(30)).SingleLine()
					text := tenths(shown[i])
					if out {
						text = o.Flag
					}
					kit.Numeric(c, text).Font("monospace").FontSize(t.Rem(1.9)).Bold().Shrink(0).
						ContentTransition(ui.ContentTransition{Numeric: true})
				})
				if b.Clicked() && !flagged && *turn != 1-i {
					commit()
					*turn, *since, *last = 1-i, now, 1-i
					root.SetChanged()
					sound.Play(sound.Tick)
					kit.Log("ChessClock", "pressed", "side", i)
				}
			}
		})
		p := kit.Press(ui.ButtonBase(c).Label(o.Pause).Padding(t.Space(1.5), t.Space(4)).Radius(999).Border(1, t.Border).Shrink(0).
			Disabled(*turn < 0 || flagged))
		p.Background(p.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(p.Hovered())), ui.Snappy))
		p.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				glyph(c, "pause", 14)
				ui.Text(c, o.Pause).FontSize(t.Rem(0.8)).Bold().Shrink(0)
			})
		})
		if p.Clicked() && *turn >= 0 && !flagged {
			commit()
			*turn, *last = -1, -1
			root.SetChanged()
			sound.Play(sound.Tick)
			kit.Log("ChessClock", "paused")
		}
	})
	return root
}

// tenths formats m:ss, with tenths under ten seconds.
func tenths(d time.Duration) string {
	if d < 10*time.Second {
		return fmt.Sprintf("0:%02d.%d", d/time.Second, d%time.Second/(100*time.Millisecond))
	}
	secs := int64(d / time.Second)
	if d%time.Second != 0 {
		secs++
	}
	return clock(secs)
}
