package game

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MemoryOptions word the game; Turns, if set, counts cards turned.
type MemoryOptions struct {
	Moves, Done, Hidden string
	Columns             int
	Turns               *int
}

// Memory turns cards to find pairs; cues Swipe, Success.
func Memory(c *ui.Context, faces []string, matched *[]bool, o MemoryOptions) ui.Element {
	def(&o.Moves, "%d moves")
	if matched == nil || len(faces) < 4 || len(faces)%2 == 1 || len(*matched) != len(faces) || o.Columns < 0 ||
		o.Turns != nil && *o.Turns < 0 {
		panic(fmt.Sprintf("game: Memory needs an even count of at least 4 faces and a matched flag each, got %d", len(faces)))
	}
	count := map[string]int{}
	for _, f := range faces {
		count[f]++
	}
	// Card order keeps the named face stable.
	for _, f := range faces {
		if f == "" || count[f] != 2 {
			panic(fmt.Sprintf("game: Memory face %q appears %d times, not twice", f, count[f]))
		}
	}
	for i := range faces {
		if (*matched)[i] && !(*matched)[pair(faces, i)] {
			panic(fmt.Sprintf("game: Memory card %d is matched without its pair", i))
		}
	}
	def(&o.Done, "All pairs in %d moves")
	def(&o.Hidden, "Card %d")
	if o.Columns == 0 {
		o.Columns = 4
	}
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	open := ui.Local(root, "open", func() []int { return nil })
	moves := ui.Local(root, "moves", func() int { return 0 })
	shut := ui.Local(root, "shut", func() time.Time { return time.Time{} })
	now := c.Now()
	if len(*open) == 2 && !now.Before(*shut) {
		*open = nil
	}
	if len(*open) == 2 {
		c.After(shut.Sub(now))
	}
	done := true
	for _, m := range *matched {
		done = done && m
	}
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			glyph(c, "brain", 16).TextColor(t.Accent)
			c.Key(done)
			if done {
				kit.Arrive(kit.Chip(c, fmt.Sprintf(o.Done, *moves), t.Success), ui.Pose{Scale: 0.7, Fade: 1}, 0)
			} else {
				kit.Numeric(c, fmt.Sprintf(o.Moves, *moves)).FontSize(t.Rem(0.8)).Bold().Shrink(0).
					ContentTransition(ui.ContentTransition{Numeric: true})
			}
		})
		grid := ui.Grid(c).Columns(o.Columns).FillWidth().Gap(t.Space(2))
		grid.Children(func() {
			for i, f := range faces {
				up := (*matched)[i] || slices.Contains(*open, i)
				label := fmt.Sprintf(o.Hidden, i+1)
				if up {
					label = f
				}
				c.Key(i)
				b := kit.Press(ui.ButtonBase(c).Label(label).AspectRatio(0.8).Radius(t.Radius).Disabled((*matched)[i]))
				ink := t.Accent
				if (*matched)[i] {
					ink = t.Success
				}
				b.Background(b.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(b.Hovered())).Mix(t.Background, unit(up)), ui.Snappy))
				b.Border(1, b.AnimateColor("line", t.Accent.Mix(ink, unit(up)), ui.Snappy))
				kit.Stagger(c, grid, b, i).Children(func() {
					c.Key(up)
					content := ui.Box(c).FillWidth().FillHeight().Center()
					content.Children(func() {
						if up {
							glyph(c, f, 24).TextColor(ink)
						} else {
							glyph(c, "sparkles", 16).TextColor(t.AccentText.Alpha(0.6))
						}
					})
					// The face swaps edge-on, like a turned card.
					k := content.Keyframes("turn", []ui.Keyframe{{Value: 0}, {At: 220 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
					content.Scale(max(0.05, k), 1)
				})
				if b.Clicked() && !up && len(*open) < 2 && !(*matched)[i] {
					*open = append(*open, i)
					cue := sound.Swipe
					if len(*open) == 2 {
						*moves++
						a, z := (*open)[0], (*open)[1]
						if faces[a] == faces[z] {
							(*matched)[a], (*matched)[z] = true, true
							*open = nil
							cue = sound.Success
							root.SetChanged()
						} else {
							*shut = now.Add(800 * time.Millisecond)
						}
					}
					if o.Turns != nil && *o.Turns < math.MaxInt {
						*o.Turns++
						root.SetChanged()
					}
					sound.Play(cue)
					kit.Log("Memory", "turned", "card", i, "moves", *moves)
				}
			}
		})
	})
	return root
}

// pair is the other card showing face i.
func pair(faces []string, i int) int {
	for j, f := range faces {
		if j != i && f == faces[i] {
			return j
		}
	}
	return -1
}
