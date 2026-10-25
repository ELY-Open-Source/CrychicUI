package health

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// HabitsOptions word the grid; Cell takes habit and day.
type HabitsOptions struct{ Title, Cell, Streak string }

// Habits toggles habits by day with streaks; cue Toggle.
func Habits(c *ui.Context, habits, days []string, done *[][]bool, o HabitsOptions) ui.Element {
	def(&o.Title, "Habits")
	if done == nil || len(habits) == 0 || len(days) == 0 || len(days) > 14 || len(*done) != len(habits) {
		panic(fmt.Sprintf("health: Habits %q needs habits, 1 to 14 days and a row of marks per habit", o.Title))
	}
	for i, row := range *done {
		if len(row) != len(days) {
			panic(fmt.Sprintf("health: Habits %q row %q has %d marks for %d days", o.Title, habits[i], len(row), len(days)))
		}
	}
	unique("Habits", habits)
	unique("Habits", days)
	def(&o.Cell, "%s on %s")
	def(&o.Streak, "%d")
	t := c.Theme()
	cell := t.Space(8)
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		kit.CardHead(c, "calendar-check", o.Title, nil)
		view := ui.ScrollHorizontal(c).FillWidth()
		view.Children(func() {
			grid := ui.Column(c).Gap(t.Space(1.5)).Padding(t.Space(3))
			grid.Children(func() {
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Width(t.Space(26)).Shrink(0)
					for _, d := range days {
						kit.Muted(c, d).FontSize(t.Rem(0.62)).Width(cell).Shrink(0).TextAlign(ui.Center).SingleLine()
					}
					ui.Box(c).Width(t.Space(10)).Shrink(0)
				})
				for i, h := range habits {
					marks := (*done)[i]
					streak := 0
					for j := len(marks) - 1; j >= 0 && marks[j]; j-- {
						streak++
					}
					c.Key(strconv.Quote(h))
					row := ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center)
					kit.Stagger(c, grid, row, i).Children(func() {
						ui.Text(c, h).FontSize(t.Rem(0.8)).Bold().Width(t.Space(26)).Shrink(0).SingleLine()
						for j, on := range marks {
							c.Key(strconv.Quote(days[j]))
							b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Cell, h, days[j])).Checked(on).Size(cell, cell).Shrink(0).Center().
								Radius(t.Space(1.5)))
							b.Background(b.AnimateColor("bg", t.SurfaceHover.Mix(t.Border, 0.5*unit(b.Hovered())).Mix(t.Success, unit(on)), ui.Snappy))
							b.TextColor(t.AccentText)
							b.Children(func() {
								if on {
									c.Key("check")
									kit.Arrive(glyph(c, "check", 14), ui.Pose{Scale: 0.3, Fade: 1}, 0)
								}
							})
							if b.Clicked() {
								(*done)[i][j] = !on
								root.SetChanged()
								sound.Play(sound.Toggle)
								kit.Log("Habits", "toggled", "habit", i, "day", j, "on", !on)
							}
						}
						ui.Row(c).Width(t.Space(10)).Shrink(0).Gap(t.Space(0.5)).AlignItems(ui.Center).Justify(ui.End).
							TextColor(t.Warning.Mix(t.TextMuted, unit(streak == 0))).Children(func() {
							glyph(c, "flame", 14)
							kit.Numeric(c, fmt.Sprintf(o.Streak, streak)).FontSize(t.Rem(0.76)).Bold().Shrink(0).
								ContentTransition(ui.ContentTransition{Numeric: true})
						})
					})
				}
			})
		})
	})
	return root
}
