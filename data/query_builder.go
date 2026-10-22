package data

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Query joins rules and groups by all, or Any.
type Query struct {
	Any    bool
	Rules  []Rule
	Groups []Query
}

// QueryBuilderOptions tune a QueryBuilder.
type QueryBuilderOptions struct {
	All, Any, AddRule, AddGroup, RemoveGroup string
	// Rule part names; each takes the rule number.
	Field, Op, Value, Remove string
}

// QueryBuilder edits nested AND and OR groups.
func QueryBuilder(c *ui.Context, label string, fields, ops []string, q *Query, o QueryBuilderOptions) ui.Element {
	if label == "" || q == nil {
		panic(fmt.Sprintf("data: QueryBuilder %q needs a query", label))
	}
	var check func(g Query)
	check = func(g Query) {
		rules("QueryBuilder", label, fields, ops, g.Rules)
		for _, s := range g.Groups {
			check(s)
		}
	}
	check(*q)
	def(&o.All, "All")
	def(&o.Any, "Any")
	def(&o.AddRule, "Add rule")
	def(&o.AddGroup, "Add group")
	def(&o.RemoveGroup, "Remove group")
	def(&o.Field, "Field %d")
	def(&o.Op, "Operator %d")
	def(&o.Value, "Value %d")
	def(&o.Remove, "Remove rule %d")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Role(ui.RoleGroup).Label(label)
	say := ruleCopy{o.Field, o.Op, o.Value, o.Remove}
	act := func(cue sound.Cue, event string, depth int) {
		root.SetChanged()
		sound.Play(cue)
		kit.Log("QueryBuilder", event, "depth", depth)
	}
	var group func(g *Query, depth int) bool
	group = func(g *Query, depth int) (gone bool) {
		box := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border).
			Background(t.Surface.Alpha(0.5 * unit(depth > 0)))
		kit.Arrive(box, ui.Pose{Y: t.Space(2), Fade: 1}, 0).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				for i, name := range []string{o.All, o.Any} {
					c.Key(i)
					on := g.Any == (i == 1)
					b := kit.Press(ui.ButtonBase(c).Label(name).Padding(t.Space(1), t.Space(3)).Radius(999).Checked(on))
					b.Background(b.AnimateColor("bg", t.SurfaceHover.Mix(t.Accent, unit(on)), ui.Snappy)).
						TextColor(b.AnimateColor("ink", t.Text.Mix(t.AccentText, unit(on)), ui.Snappy)).Children(func() { ui.Text(c, name).Bold() })
					if b.Clicked() && !on {
						g.Any = i == 1
						act(sound.Toggle, "joined", depth)
					}
				}
				ui.Box(c).Grow(1)
				if pill(c, o.AddRule) {
					g.Rules = append(g.Rules, Rule{Field: fields[0], Op: ops[0]})
					act(sound.Pop, "added rule", depth)
				}
				if pill(c, o.AddGroup) {
					g.Groups = append(g.Groups, Query{Rules: []Rule{{Field: fields[0], Op: ops[0]}}})
					act(sound.Pop, "added group", depth)
				}
				if depth > 0 && pill(c, o.RemoveGroup) {
					gone = true
				}
			})
			drop := -1
			for i := range g.Rules {
				c.Key(fmt.Sprint("rule", i))
				edited, remove := ruleRow(c, fields, ops, &g.Rules[i], i+1, say)
				if edited {
					root.SetChanged()
					kit.Log("QueryBuilder", "edited", "depth", depth, "index", i)
				}
				if remove {
					drop = i
				}
			}
			if drop >= 0 {
				g.Rules = slices.Delete(g.Rules, drop, drop+1)
				act(sound.Delete, "removed rule", depth)
			}
			lost := -1
			for i := range g.Groups {
				c.Key(fmt.Sprint("group", i))
				if group(&g.Groups[i], depth+1) {
					lost = i
				}
			}
			if lost >= 0 {
				g.Groups = slices.Delete(g.Groups, lost, lost+1)
				act(sound.Delete, "removed group", depth+1)
			}
		})
		return gone
	}
	root.Children(func() { group(q, 0) })
	return root
}
