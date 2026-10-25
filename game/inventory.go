package game

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Item is one inventory item; blank Name means empty.
type Item struct {
	Name, Icon string
	Count      int
	Rarity     Rarity
}

// InventoryOptions word the grid; Slot takes name and count.
type InventoryOptions struct {
	Title, Empty, Slot, Used string
	Columns                  int
}

// Inventory is a slot grid; dragging swaps slots; cue Drop.
func Inventory(c *ui.Context, slots *[]Item, o InventoryOptions) ui.Element {
	def(&o.Title, "Inventory")
	if slots == nil || len(*slots) == 0 || o.Columns < 0 {
		panic(fmt.Sprintf("game: Inventory %q needs slots", o.Title))
	}
	used := 0
	for i, it := range *slots {
		if it.Name == "" {
			continue
		}
		validRarity("Inventory", it.Name, it.Rarity)
		if it.Icon == "" || it.Count < 1 {
			panic(fmt.Sprintf("game: Inventory %q slot %d needs an icon and a count of at least 1", it.Name, i))
		}
		used++
	}
	def(&o.Empty, "Empty slot %d")
	def(&o.Slot, "%s × %d")
	def(&o.Used, "%d / %d")
	if o.Columns == 0 {
		o.Columns = 5
	}
	t := c.Theme()
	root := card(c).Role(ui.RoleGroup)
	id := ui.Local(root, "id", func() *int { return new(int) })
	from, to := -1, -1
	root.Children(func() {
		kit.CardHead(c, "backpack", o.Title, func() { kit.Chip(c, fmt.Sprintf(o.Used, used, len(*slots)), t.Accent) })
		grid := ui.Grid(c).Columns(o.Columns).FillWidth().Gap(t.Space(2)).Padding(t.Space(3))
		grid.Children(func() {
			for i, it := range *slots {
				c.Key(i)
				slot := ui.Box(c).AspectRatio(1).Center().Radius(t.Radius).Focusable()
				_, over := ui.DragOver[slotDrag](slot)
				ink := rarityInk(t, it.Rarity)
				if it.Name == "" {
					ink = t.Border
					slot.Label(fmt.Sprintf(o.Empty, i+1))
				} else {
					label := it.Name
					if it.Count > 1 {
						label = fmt.Sprintf(o.Slot, it.Name, it.Count)
					}
					slot.Label(label).Cursor(ui.CursorGrab).Drag(slotDrag{grid: *id, from: i})
				}
				slot.Background(slot.AnimateColor("bg", t.SurfaceHover.Alpha(0.5).Mix(ink.Alpha(0.12), unit(it.Name != "")).
					Mix(t.Accent.Alpha(0.2), unit(over)), ui.Snappy))
				slot.Border(1, slot.AnimateColor("line", ink.Alpha(0.6).Mix(t.Accent, unit(over)), ui.Snappy))
				kit.Stagger(c, grid, slot, i).Children(func() {
					if it.Name == "" {
						return
					}
					c.Key(it.Name)
					kit.Arrive(glyph(c, it.Icon, 22).TextColor(ink), ui.Pose{Scale: 0.6, Fade: 1}, 0)
					if it.Count > 1 {
						ui.Box(c).Absolute().Right(t.Space(1)).Bottom(t.Space(0.5)).Children(func() {
							kit.Numeric(c, fmt.Sprint(it.Count)).Font("monospace").FontSize(t.Rem(0.62)).Bold().Shrink(0)
						})
					}
				})
				if d, ok := ui.Drop[slotDrag](slot); ok && d.grid == *id && d.from != i {
					from, to = d.from, i
				}
			}
		})
	})
	if from >= 0 {
		(*slots)[from], (*slots)[to] = (*slots)[to], (*slots)[from]
		root.SetChanged()
		sound.Play(sound.Drop)
		kit.Log("Inventory", "moved", "from", from, "to", to)
	}
	return root
}

type slotDrag struct {
	grid *int
	from int
}
