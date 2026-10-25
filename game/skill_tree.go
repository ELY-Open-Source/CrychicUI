package game

import (
	"fmt"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Skill is one tree node; it needs all of Requires.
type Skill struct {
	ID, Name, Icon string
	Requires       []string
	Cost           int
}

// SkillTreeOptions word the tree; Learn takes a skill name.
type SkillTreeOptions struct{ Title, Points, Learn, Learned, Cost string }

// SkillTree unlocks skills along their links; cue Success.
func SkillTree(c *ui.Context, skills []Skill, learned *[]string, points *int, o SkillTreeOptions) ui.Element {
	def(&o.Title, "Skills")
	if learned == nil || points == nil || len(skills) == 0 || *points < 0 {
		panic(fmt.Sprintf("game: SkillTree %q needs skills, learned IDs and points of at least 0", o.Title))
	}
	ids := make([]string, len(skills))
	index := map[string]int{}
	for i, s := range skills {
		if s.Name == "" || s.Icon == "" || s.Cost < 1 {
			panic(fmt.Sprintf("game: SkillTree %q needs a name, icon and a cost of at least 1", s.ID))
		}
		ids[i] = s.ID
		index[s.ID] = i
	}
	unique("SkillTree", ids)
	for _, s := range skills {
		for _, r := range s.Requires {
			if _, ok := index[r]; !ok || r == s.ID {
				panic(fmt.Sprintf("game: SkillTree %q requires unknown skill %q", s.ID, r))
			}
		}
	}
	for _, id := range *learned {
		if _, ok := index[id]; !ok {
			panic(fmt.Sprintf("game: SkillTree %q learned unknown skill %q", o.Title, id))
		}
	}
	def(&o.Points, "%d points")
	def(&o.Learn, "Learn %s")
	def(&o.Learned, "%s, learned")
	def(&o.Cost, "%d pt")
	depth := skillDepth(o.Title, skills, index)
	rows, slot := map[int]int{}, make([]int, len(skills))
	deepest := 0
	for i := range skills {
		slot[i] = rows[depth[i]]
		rows[depth[i]]++
		deepest = max(deepest, depth[i])
	}
	has := func(id string) bool { return slices.Contains(*learned, id) }
	open := func(s Skill) bool {
		for _, r := range s.Requires {
			if !has(r) {
				return false
			}
		}
		return !has(s.ID) && *points >= s.Cost
	}
	t := c.Theme()
	node, rowH := t.Space(13), t.Space(26)
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		kit.CardHead(c, "git-fork", o.Title, func() {
			tone := t.TextMuted
			if *points > 0 {
				tone = t.Warning
			}
			kit.Chip(c, fmt.Sprintf(o.Points, *points), tone)
		})
		view := ui.ScrollHorizontal(c).FillWidth()
		seen := view.LayoutBounds().W
		kit.Settle(c, view, []float32{seen}, func() []float32 { return []float32{view.LayoutBounds().W} })
		widest := 0
		for _, n := range rows {
			widest = max(widest, n)
		}
		// Each node keeps room for its label; wide rows scroll.
		w := max(seen, float32(widest)*t.Space(30))
		at := func(i int) (float32, float32) {
			n := rows[depth[i]]
			return w * (float32(slot[i]) + 0.5) / float32(n), t.Space(4) + float32(depth[i])*rowH + node/2
		}
		view.Children(func() {
			canvas := ui.Box(c).Width(w).Height(float32(deepest+1)*rowH + t.Space(4))
			canvas.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				for i, s := range skills {
					for _, req := range s.Requires {
						x0, y0 := at(index[req])
						x1, y1 := at(i)
						ink := t.Border
						if has(req) && has(s.ID) {
							ink = t.Accent
						} else if has(req) {
							ink = t.Accent.Alpha(0.4)
						}
						var link ui.Path
						// Links leave below the name and cost labels.
						y0, y1 = r.Y+y0+node/2+t.Space(8), r.Y+y1-node/2
						mid := (y0 + y1) / 2
						link.MoveTo(r.X+x0, y0).CubeTo(r.X+x0, mid, r.X+x1, mid, r.X+x1, y1)
						p.StrokePath(&link, t.Space(0.5), ink)
					}
				}
			}))
			canvas.Children(func() {
				for i, s := range skills {
					x, y := at(i)
					learnedNow, can := has(s.ID), open(s)
					label := fmt.Sprintf(o.Learn, s.Name)
					if learnedNow {
						label = fmt.Sprintf(o.Learned, s.Name)
					}
					c.Key(s.ID)
					b := kit.Press(ui.ButtonBase(c).Label(label).Absolute().Left(x-node/2).Top(y-node/2).Size(node, node).Center().
						Radius(999).Disabled(!can && !learnedNow))
					ink := t.TextMuted
					switch {
					case learnedNow:
						ink = t.Accent
					case can:
						ink = t.Warning
					}
					b.Background(b.AnimateColor("bg", t.SurfaceHover.Mix(t.Accent, unit(learnedNow)).Mix(t.Warning.Alpha(0.15), unit(can)), ui.Snappy))
					b.Border(t.Space(0.5), b.AnimateColor("line", t.Border.Mix(ink, unit(learnedNow || can)), ui.Snappy))
					kit.Arrive(b, ui.Pose{Scale: 0.6, Fade: 1}, time.Duration(depth[i])*80*time.Millisecond).Children(func() {
						g := glyph(c, s.Icon, 20)
						g.TextColor(g.AnimateColor("ink", t.TextMuted.Mix(t.AccentText, unit(learnedNow)).Mix(t.Warning, unit(can)), ui.Snappy))
					})
					ui.Column(c).Absolute().Left(x - t.Space(14)).Top(y + node/2 + t.Space(1)).Width(t.Space(28)).AlignItems(ui.Center).Children(func() {
						ui.Text(c, s.Name).FontSize(t.Rem(0.7)).Bold().TextColor(t.Text.Mix(t.TextMuted, unit(!learnedNow && !can))).
							FillWidth().TextAlign(ui.Center).SingleLine()
						if !learnedNow {
							kit.Muted(c, fmt.Sprintf(o.Cost, s.Cost)).FontSize(t.Rem(0.6)).Shrink(0)
						}
					})
					if b.Clicked() && can {
						*learned = append(*learned, s.ID)
						*points -= s.Cost
						root.SetChanged()
						sound.Play(sound.Success)
						kit.Log("SkillTree", "learned", "index", i, "left", *points)
					}
				}
			})
		})
	})
	return root
}

// skillDepth is each skill's row along requirements; loops panic.
func skillDepth(title string, skills []Skill, index map[string]int) []int {
	depth := make([]int, len(skills))
	state := make([]int, len(skills))
	var visit func(i int) int
	visit = func(i int) int {
		switch state[i] {
		case 1:
			panic(fmt.Sprintf("game: SkillTree %q requirements form a loop at %q", title, skills[i].ID))
		case 2:
			return depth[i]
		}
		state[i] = 1
		d := 0
		for _, r := range skills[i].Requires {
			d = max(d, visit(index[r])+1)
		}
		depth[i], state[i] = d, 2
		return d
	}
	for i := range skills {
		visit(i)
	}
	return depth
}
