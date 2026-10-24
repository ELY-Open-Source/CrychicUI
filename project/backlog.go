package project

import (
	"fmt"
	"math"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BacklogOptions word the list; Total takes issues and points.
type BacklogOptions struct{ Title, Total string }

// Backlog orders issues by drag; cue Drop.
func Backlog(c *ui.Context, issues *[]Issue, o BacklogOptions) ui.Element {
	if issues == nil || len(*issues) == 0 {
		panic("project: Backlog needs issues")
	}
	def(&o.Title, "Backlog")
	keys := make([]string, len(*issues))
	points := 0
	for i, issue := range *issues {
		validIssue("Backlog", issue)
		keys[i] = issue.Key
		if issue.Points > math.MaxInt-points {
			panic(fmt.Sprintf("project: Backlog %q points overflow", o.Title))
		}
		points += issue.Points
	}
	unique("Backlog", keys)
	def(&o.Total, "%d issues · %d pts")
	t := c.Theme()
	root := card(c)
	gridID := ui.Local(root, "gridID", func() *int { return new(int) })
	moveKey, moveTo := "", -1
	root.Children(func() {
		kit.CardHead(c, "list-ordered", o.Title, func() {
			kit.Chip(c, fmt.Sprintf(o.Total, len(*issues), points), t.TextMuted)
		})
		list := ui.Column(c).FillWidth().Padding(t.Space(1))
		list.Children(func() {
			for i, issue := range *issues {
				kind, kindInk := issueLook(t, issue.Type)
				level, levelInk := priorityLook(t, issue.Priority)
				c.Key(issue.Key)
				row := ui.Row(c).FillWidth().MinHeight(t.Space(9)).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).
					Focusable().Label(issue.Title).Cursor(ui.CursorGrab).Drag(backlogDrag{grid: *gridID, key: issue.Key}).
					Transition(ui.ElementTransition{Animation: &ui.Snappy, Position: true})
				d, over := ui.DragOver[backlogDrag](row)
				over = over && d.key != issue.Key
				row.Background(row.AnimateColor("bg", t.SurfaceHover.Alpha(0.7*unit(row.Hovered())).Mix(t.Accent.Alpha(0.12), unit(over)), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					glyph(c, "grip-vertical", 14).TextColor(t.TextMuted)
					kit.Numeric(c, fmt.Sprint(i+1)).FontSize(t.Rem(0.7)).TextColor(t.TextMuted).Width(t.Space(4)).Shrink(0)
					glyph(c, kind, 14).TextColor(kindInk)
					ui.Text(c, issue.Key).Font("monospace").FontSize(t.Rem(0.7)).Bold().TextColor(t.TextMuted).Shrink(0)
					ui.Text(c, issue.Title).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).SingleLine()
					glyph(c, level, 14).TextColor(levelInk)
					if issue.Points > 0 {
						kit.Chip(c, fmt.Sprint(issue.Points), t.TextMuted)
					}
					if issue.Assignee != "" {
						display.Avatar(c, issue.Assignee, nil).Size(t.Space(6), t.Space(6)).FontSize(t.Rem(0.6)).Tooltip(issue.Assignee)
					}
				})
				if d, ok := ui.Drop[backlogDrag](row); ok && d.grid == *gridID {
					moveKey, moveTo = d.key, i
				}
			}
		})
	})
	if moveKey != "" {
		from := slices.IndexFunc(*issues, func(i Issue) bool { return i.Key == moveKey })
		if from >= 0 && from != moveTo {
			issue := (*issues)[from]
			*issues = slices.Insert(slices.Delete(*issues, from, from+1), moveTo, issue)
			root.SetChanged()
			sound.Play(sound.Drop)
			kit.Log("Backlog", "moved", "from", from, "to", moveTo)
		}
	}
	return root
}

type backlogDrag struct {
	grid *int
	key  string
}
