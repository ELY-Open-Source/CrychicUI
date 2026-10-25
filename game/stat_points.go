package game

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// StatPointsOptions word the panel; Raise and Lower take names.
type StatPointsOptions struct{ Title, Points, Raise, Lower string }

// StatPoints spends points, never below start values; cue Tick.
func StatPoints(c *ui.Context, stats *[]Stat, points *int, o StatPointsOptions) ui.Element {
	def(&o.Title, "Attributes")
	if stats == nil || points == nil || len(*stats) == 0 || *points < 0 {
		panic(fmt.Sprintf("game: StatPoints %q needs stats and points of at least 0", o.Title))
	}
	names := make([]string, len(*stats))
	top := 1
	for i, s := range *stats {
		if s.Value < 0 {
			panic(fmt.Sprintf("game: StatPoints %q value %d is negative", s.Name, s.Value))
		}
		names[i] = s.Name
		top = max(top, s.Value)
	}
	unique("StatPoints", names)
	def(&o.Points, "%d points")
	def(&o.Raise, "Raise %s")
	def(&o.Lower, "Lower %s")
	t := c.Theme()
	root := card(c).Role(ui.RoleGroup)
	base := ui.Local(root, "base", func() map[string]int {
		m := map[string]int{}
		for _, s := range *stats {
			m[s.Name] = s.Value
		}
		return m
	})
	root.Children(func() {
		kit.CardHead(c, "swords", o.Title, func() {
			tone := t.TextMuted
			if *points > 0 {
				tone = t.Warning
			}
			kit.Chip(c, fmt.Sprintf(o.Points, *points), tone)
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2), t.Space(3))
		list.Children(func() {
			for i, s := range *stats {
				floor, seen := (*base)[s.Name]
				if !seen {
					(*base)[s.Name] = s.Value
					floor = s.Value
				}
				c.Key(s.Name)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(1), 0)
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Text(c, s.Name).FontSize(t.Rem(0.82)).Bold().Width(t.Space(22)).Shrink(0).SingleLine()
					ui.Box(c).Grow(1).MinWidth(0).Children(func() {
						meter(c, float32(s.Value)/float32(top), t.Accent.Mix(t.Warning, unit(s.Value > floor))).Height(t.Space(1.5))
					})
					step := func(label, mark string, on bool) ui.Element {
						b := kit.Press(ui.ButtonBase(c).Label(label).Size(t.Space(7), t.Space(7)).Center().Radius(999).Shrink(0).
							Border(1, t.Border).Disabled(!on))
						b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
						return b.Children(func() { glyph(c, mark, 14) })
					}
					refund := s.Value > floor && *points < math.MaxInt
					down := step(fmt.Sprintf(o.Lower, s.Name), "minus", refund)
					n := kit.Numeric(c, fmt.Sprint(s.Value)).FontSize(t.Rem(0.95)).Bold().Width(t.Space(8)).Shrink(0).TextAlign(ui.Center).
						ContentTransition(ui.ContentTransition{Numeric: true})
					n.TextColor(n.AnimateColor("ink", t.Text.Mix(t.Warning, unit(s.Value > floor)), ui.Snappy))
					up := step(fmt.Sprintf(o.Raise, s.Name), "plus", *points > 0 && s.Value < math.MaxInt)
					delta := 0
					if down.Clicked() && refund {
						delta = -1
					}
					if up.Clicked() && *points > 0 && s.Value < math.MaxInt {
						delta = 1
					}
					if delta != 0 {
						(*stats)[i].Value += delta
						*points -= delta
						root.SetChanged()
						sound.Play(sound.Tick)
						kit.Log("StatPoints", "spent", "stat", i, "delta", delta, "left", *points)
					}
				})
			}
		})
	})
	return root
}
