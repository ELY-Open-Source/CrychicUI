package game

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Level is one level tile; Stars is 0 to 3.
type Level struct {
	Name   string
	Stars  int
	Locked bool
}

// LevelPickOptions word the grid; Label and Locked take the name.
type LevelPickOptions struct {
	Title, Label, Locked, Total string
	Columns                     int
}

// LevelPick is a grid of starred levels; cue Select.
func LevelPick(c *ui.Context, levels []Level, picked *int, o LevelPickOptions) ui.Element {
	def(&o.Title, "Levels")
	if picked == nil || len(levels) == 0 || *picked < -1 || *picked >= len(levels) || o.Columns < 0 {
		panic(fmt.Sprintf("game: LevelPick %q needs levels and a pick in range", o.Title))
	}
	names := make([]string, len(levels))
	total := 0
	for i, l := range levels {
		if l.Stars < 0 || l.Stars > 3 || l.Locked && l.Stars > 0 {
			panic(fmt.Sprintf("game: LevelPick %q needs 0 to 3 stars, and none while locked", l.Name))
		}
		names[i] = l.Name
		total += l.Stars
	}
	unique("LevelPick", names)
	if *picked >= 0 && levels[*picked].Locked {
		panic(fmt.Sprintf("game: LevelPick %q pick %q is locked", o.Title, levels[*picked].Name))
	}
	def(&o.Label, "Level %s")
	def(&o.Locked, "Level %s, locked")
	def(&o.Total, "%d / %d")
	if o.Columns == 0 {
		o.Columns = 4
	}
	t := c.Theme()
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		kit.CardHead(c, "map", o.Title, func() {
			ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
				stars(c, 1, t.Space(3), t.Warning)
				kit.Muted(c, fmt.Sprintf(o.Total, total, 3*len(levels))).FontSize(t.Rem(0.72)).Bold().Shrink(0)
			})
		})
		grid := ui.Grid(c).Columns(o.Columns).FillWidth().Gap(t.Space(2)).Padding(t.Space(3))
		grid.Children(func() {
			for i, l := range levels {
				on := *picked == i
				label := fmt.Sprintf(o.Label, l.Name)
				if l.Locked {
					label = fmt.Sprintf(o.Locked, l.Name)
				}
				c.Key(l.Name)
				b := kit.Press(ui.ButtonBase(c).Label(label).AspectRatio(1).Column().Center().Gap(t.Space(1)).Radius(t.Radius * 1.2).
					Disabled(l.Locked).Checked(on))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(0.6).Mix(t.SurfaceHover, unit(b.Hovered())).
					Mix(t.Accent.Alpha(0.12), unit(on)), ui.Snappy))
				b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
				kit.Stagger(c, grid, b, i).Children(func() {
					if l.Locked {
						glyph(c, "lock", 18).TextColor(t.TextMuted)
						return
					}
					ui.Text(c, l.Name).FontSize(t.Rem(1.05)).Bold().TextColor(t.Text.Mix(t.Accent, unit(on))).MaxWidth(t.Space(16)).SingleLine()
					stars(c, l.Stars, t.Space(2.5), t.Warning)
				})
				if b.Clicked() && !on && !l.Locked {
					*picked = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("LevelPick", "picked", "index", i)
				}
			}
		})
	})
	return root
}
