package project

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Lane is one Kanban column.
type Lane struct {
	Title string
	Cards []Task
}

// BoardOptions are the board labels; New titles added cards.
type BoardOptions struct{ Add, New, Empty string }

// Board adds draggable Kanban cards; cue Drop.
func Board(c *ui.Context, lanes *[]Lane, o BoardOptions) ui.Element {
	if lanes == nil || len(*lanes) == 0 {
		panic("project: Board needs lanes")
	}
	var titles, ids []string
	for _, lane := range *lanes {
		titles = append(titles, lane.Title)
		for _, task := range lane.Cards {
			validTask("Board", task)
			ids = append(ids, task.ID)
		}
	}
	unique("Board", titles)
	unique("Board", ids)
	def(&o.Add, "Add card")
	def(&o.New, "New card")
	def(&o.Empty, "Drop cards here")
	t := c.Theme()
	heads := []ui.Color{t.TextMuted, t.Accent, t.Warning, t.Success}
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Start).Role(ui.RoleGroup)
	gridID := ui.Local(root, "gridID", func() *int { return new(int) })
	made := ui.Local(root, "made", func() int { return 0 })
	addLane, moveID, toLane, before := -1, "", -1, ""
	root.Children(func() {
		for li, lane := range *lanes {
			c.Key(lane.Title)
			zone := ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(1.5)).Padding(t.Space(2)).Radius(t.Radius)
			_, over := ui.DragOver[boardDrag](zone)
			zone.Background(zone.AnimateColor("bg", t.SurfaceHover.Alpha(0.5).Mix(t.Accent.Alpha(0.1), unit(over)), ui.Snappy))
			zone.Border(1, zone.AnimateColor("line", t.Border.Alpha(0).Mix(t.Accent.Alpha(0.5), unit(over)), ui.Snappy))
			zone.Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(1)).Children(func() {
					ui.Box(c).Size(t.Space(2), t.Space(2)).Shrink(0).Radius(999).Background(heads[li%len(heads)])
					ui.Text(c, lane.Title).FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0).SingleLine()
					kit.Chip(c, fmt.Sprint(len(lane.Cards)), t.TextMuted)
				})
				// Cards share a parent only with cards.
				cards := ui.Column(c).FillWidth().Gap(t.Space(1.5))
				cards.Children(func() {
					for ci, task := range lane.Cards {
						c.Key(task.ID)
						row := kit.Lift(kit.Card(c).FillWidth().Gap(t.Space(1.5)).Padding(t.Space(2.5)).Focusable().Label(task.Title).
							Cursor(ui.CursorGrab).Drag(boardDrag{grid: *gridID, id: task.ID}).
							Transition(ui.ElementTransition{Animation: &ui.Snappy, Position: true}))
						kit.Stagger(c, cards, row, ci).Children(func() { taskFace(c, task) })
						if d, ok := ui.Drop[boardDrag](row); ok && d.grid == *gridID {
							moveID, toLane, before = d.id, li, task.ID
						}
					}
					if len(lane.Cards) == 0 {
						ui.Box(c).FillWidth().Height(t.Space(14)).Center().Radius(t.Radius).Border(1, t.Border).
							BorderStyle(ui.BorderDashed).Children(func() { kit.Muted(c, o.Empty).FontSize(t.Rem(0.7)) })
					}
				})
				add := kit.Press(ui.ButtonBase(c).Label(o.Add).FillWidth().Padding(t.Space(1.5)).Radius(t.Radius))
				add.Background(add.AnimateColor("bg", t.SurfaceHover.Alpha(unit(add.Hovered())), ui.Snappy))
				add.TextColor(add.AnimateColor("ink", t.TextMuted.Mix(t.Text, unit(add.Hovered())), ui.Snappy))
				add.Children(func() {
					ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
						glyph(c, "plus", 14)
						ui.Text(c, o.Add).FontSize(t.Rem(0.75)).Shrink(0)
					})
				})
				if add.Clicked() {
					addLane = li
				}
			})
			if d, ok := ui.Drop[boardDrag](zone); ok && d.grid == *gridID && moveID == "" {
				moveID, toLane, before = d.id, li, ""
			}
		}
	})
	if addLane >= 0 {
		id := ""
		for id == "" || slices.Contains(ids, id) {
			*made++
			id = fmt.Sprint("new-", *made)
		}
		(*lanes)[addLane].Cards = append((*lanes)[addLane].Cards, Task{ID: id, Title: o.New})
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("Board", "added", "lane", addLane)
	}
	if moveID != "" && moveID != before && moveCard(*lanes, moveID, toLane, before) {
		root.SetChanged()
		sound.Play(sound.Drop)
		kit.Log("Board", "moved", "lane", toLane)
	}
	return root
}

// moveCard moves card id before card before; false if unmoved.
func moveCard(lanes []Lane, id string, to int, before string) bool {
	var task Task
	from, was := -1, -1
	for li := range lanes {
		if i := slices.IndexFunc(lanes[li].Cards, func(t Task) bool { return t.ID == id }); i >= 0 {
			task, from, was = lanes[li].Cards[i], li, i
			lanes[li].Cards = slices.Delete(lanes[li].Cards, i, i+1)
		}
	}
	at := slices.IndexFunc(lanes[to].Cards, func(t Task) bool { return t.ID == before })
	if at < 0 {
		at = len(lanes[to].Cards)
	}
	lanes[to].Cards = slices.Insert(lanes[to].Cards, at, task)
	return from != to || was != at
}

type boardDrag struct {
	grid *int
	id   string
}
