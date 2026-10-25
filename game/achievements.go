package game

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Achievement is one goal; Goal 0 means no progress bar.
type Achievement struct {
	Name, Text, Icon, When string
	Goal, Progress         int
	Unlocked               bool
}

// AchievementsOptions word the list; Count takes unlocked and total.
type AchievementsOptions struct{ Title, Count, Progress string }

// Achievements lists locked and unlocked achievements.
func Achievements(c *ui.Context, list []Achievement, o AchievementsOptions) ui.Element {
	def(&o.Title, "Achievements")
	if len(list) == 0 {
		panic(fmt.Sprintf("game: Achievements %q needs achievements", o.Title))
	}
	names := make([]string, len(list))
	got := 0
	for i, a := range list {
		if a.Icon == "" || a.Goal < 0 || a.Progress < 0 || a.Progress > a.Goal {
			panic(fmt.Sprintf("game: Achievements %q needs an icon and 0 <= progress <= goal", a.Name))
		}
		names[i] = a.Name
		if a.Unlocked {
			got++
		}
	}
	unique("Achievements", names)
	def(&o.Count, "%d of %d unlocked")
	def(&o.Progress, "%d / %d")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "medal", o.Title, func() { kit.Chip(c, fmt.Sprintf(o.Count, got, len(list)), t.Warning) })
		rows := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		rows.Children(func() {
			for i, a := range list {
				c.Key(a.Name)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, rows, row, i).Children(func() {
					c.Key(a.Unlocked)
					if a.Unlocked {
						kit.Arrive(tile(c, a.Icon, t.Warning, t.Space(10)), ui.Pose{Scale: 0.6, Rotate: -15, Fade: 1}, 0)
					} else {
						tile(c, "lock", t.TextMuted, t.Space(10))
					}
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.75)).Children(func() {
						name := ui.Text(c, a.Name).FontSize(t.Rem(0.85)).Bold().FillWidth().MinWidth(0).SingleLine()
						if !a.Unlocked {
							name.TextColor(t.TextMuted)
						}
						if a.Text != "" {
							kit.Muted(c, a.Text).FontSize(t.Rem(0.72)).FillWidth().Wrap()
						}
						if !a.Unlocked && a.Goal > 0 {
							ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
								ui.Box(c).Grow(1).MinWidth(0).Children(func() {
									meter(c, float32(float64(a.Progress)/float64(a.Goal)), t.Accent).Height(t.Space(1.25))
								})
								kit.Muted(c, fmt.Sprintf(o.Progress, a.Progress, a.Goal)).Font("monospace").FontSize(t.Rem(0.65)).Shrink(0)
							})
						}
					})
					if a.Unlocked {
						ui.Column(c).Shrink(0).Gap(t.Space(0.5)).AlignItems(ui.End).Children(func() {
							glyph(c, "circle-check", 16).TextColor(t.Success)
							if a.When != "" {
								kit.Muted(c, a.When).FontSize(t.Rem(0.65)).Shrink(0)
							}
						})
					}
				})
			}
		})
	})
	return root
}
