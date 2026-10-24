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
	Cards []string
}

// BoardOptions are the board labels.
type BoardOptions struct{ Add, New string }

// Board adds draggable Kanban cards; cue Drop.
func Board(c *ui.Context, lanes *[]Lane, o BoardOptions) ui.Element {
	if lanes == nil || len(*lanes) == 0 {
		panic("project: Board needs lanes")
	}
	def(&o.Add, "Add card")
	def(&o.New, "New card")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	gridID := ui.Local(root, "gridID", func() *int { return new(int) })
	addTitle, addLane := "", -1
	moveTitle, moveTo := "", [2]int{-1, -1}
	moveFrom, moveIndex := -1, -1
	root.Children(func() {
		for li, lane := range *lanes {
			if lane.Title == "" {
				panic("project: Board lane title must not be empty")
			}
			c.Key(li)
			zone := ui.Column(c).Grow(1).Basis(0).Gap(t.Space(1)).Padding(t.Space(1)).Radius(t.Radius).Background(t.Surface)
			zone.Children(func() {
				ui.Row(c).FillWidth().Children(func() {
					ui.Text(c, lane.Title).Bold().Grow(1)
					kit.Numeric(c, fmt.Sprint(len(lane.Cards))).TextColor(t.TextMuted)
				})
				for ci, card := range lane.Cards {
					if card == "" {
						panic("project: Board card must not be empty")
					}
					c.Key(fmt.Sprintf("%s.%d", card, ci))
					row := ui.Box(c).FillWidth().Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border).
						Focusable().Label(card).Cursor(ui.CursorGrab).Drag(boardDrag{grid: *gridID, lane: li, index: ci, card: card})
					row.Children(func() { ui.Text(c, card) })
					if d, ok := ui.Drop[boardDrag](row); ok && d.grid == *gridID {
						moveTitle, moveTo, moveFrom, moveIndex = d.card, [2]int{li, ci}, d.lane, d.index
					}
				}
				add := kit.Press(ui.ButtonBase(c).Label(o.Add).Padding(t.Space(1)).Radius(t.Radius))
				add.Children(func() { ui.Text(c, o.Add).FontSize(t.Rem(0.75)).TextColor(t.TextMuted) })
				if add.Clicked() {
					addTitle, addLane = o.New, li
				}
			})
			if d, ok := ui.Drop[boardDrag](zone); ok && d.grid == *gridID {
				moveTitle, moveTo, moveFrom, moveIndex = d.card, [2]int{li, len(lane.Cards)}, d.lane, d.index
			}
		}
	})
	if addTitle != "" {
		(*lanes)[addLane].Cards = append((*lanes)[addLane].Cards, addTitle)
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("Board", "added", "lane", addLane)
	}
	if moveTitle != "" {
		if moveFrom >= 0 && moveIndex >= 0 {
			(*lanes)[moveFrom].Cards = slices.Delete((*lanes)[moveFrom].Cards, moveIndex, moveIndex+1)
			to := min(moveTo[1], len((*lanes)[moveTo[0]].Cards))
			(*lanes)[moveTo[0]].Cards = slices.Insert((*lanes)[moveTo[0]].Cards, to, moveTitle)
			root.SetChanged()
			sound.Play(sound.Drop)
			kit.Log("Board", "moved", "from", moveFrom, "to", moveTo)
		}
	}
	return root
}

type boardDrag struct {
	grid  *int
	lane  int
	index int
	card  string
}
