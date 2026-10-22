package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Spreadsheet is a grid of cells; Enter edits.
func Spreadsheet(c *ui.Context, label string, cells [][]string, at *[2]int) ui.Element {
	if label == "" || len(cells) == 0 || len(cells[0]) == 0 || at == nil || !valid(at[0], len(cells)) || !valid(at[1], len(cells[0])) {
		panic(fmt.Sprintf("data: Spreadsheet %q needs cells and a chosen cell in them", label))
	}
	for _, r := range cells {
		if len(r) != len(cells[0]) {
			panic(fmt.Sprintf("data: Spreadsheet %q rows need the same length", label))
		}
	}
	t := c.Theme()
	w, h := t.Space(24), t.Space(8)
	root := ui.Column(c).AlignSelf(ui.Start).Focusable().Role(ui.RoleGroup).Label(label).Radius(t.Radius).Border(1, t.Border).Clip()
	sheet := ui.Local(root, "sheet", func() ui.Handle { return ui.Handle{} })
	field := ui.Local(root, "field", func() ui.Handle { return ui.Handle{} })
	editing := ui.Local(root, "editing", func() bool { return false })
	draft := ui.Local(root, "draft", func() string { return "" })
	root.Bind(sheet)
	pick := func(r, col int) {
		if [2]int{r, col} == *at {
			return
		}
		*at, *editing = [2]int{r, col}, false
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("Spreadsheet", "chose", "row", r, "column", col)
	}
	edit := func() {
		*editing, *draft = true, cells[at[0]][at[1]]
		field.Focus(c)
		root.SetChanged()
		kit.Log("Spreadsheet", "editing")
	}
	if !*editing {
		r, col := at[0], at[1]
		switch {
		case root.Shortcut(0, ui.KeyUp) && r > 0:
			pick(r-1, col)
		case root.Shortcut(0, ui.KeyDown) && r < len(cells)-1:
			pick(r+1, col)
		case root.Shortcut(0, ui.KeyLeft) && col > 0:
			pick(r, col-1)
		case root.Shortcut(0, ui.KeyRight) && col < len(cells[0])-1:
			pick(r, col+1)
		case root.Shortcut(0, ui.KeyEnter):
			edit()
		}
	}
	gridLine := func(content func()) { ui.Row(c).Children(content) }
	headCell := func(s string, width float32) {
		ui.Row(c).Size(width, h).Center().Background(t.Surface).Border(1, t.Border).Children(func() {
			kit.Muted(c, s).FontSize(t.Rem(0.75))
		})
	}
	root.Children(func() {
		gridLine(func() {
			headCell("", t.Space(10))
			for col := range cells[0] {
				c.Key(col)
				headCell(letters(col), w)
			}
		})
		for r, row := range cells {
			c.Key(r)
			gridLine(func() {
				headCell(fmt.Sprint(r+1), t.Space(10))
				for col := range row {
					c.Key(col)
					on := *at == [2]int{r, col}
					box := ui.Row(c).Size(w, h).AlignItems(ui.Center).Padding(0, t.Space(2)).Border(1, t.Border)
					box.Background(box.AnimateColor("bg", t.Background.Mix(t.Accent, 0.14*unit(on)), ui.Snappy))
					box.Children(func() {
						if on && *editing {
							in := ui.TextInput(c, draft).Label(fmt.Sprintf("%s%d", letters(col), r+1)).FillWidth().Bind(field)
							switch {
							case in.Submitted():
								cells[r][col], *editing = *draft, false
								sheet.Focus(c)
								root.SetChanged()
								sound.Play(sound.Tick)
								kit.Log("Spreadsheet", "saved", "row", r, "column", col)
							case in.Shortcut(0, ui.KeyEscape):
								*editing = false
								sheet.Focus(c)
								root.SetChanged()
							}
							return
						}
						ui.Text(c, row[col]).SingleLine()
					})
					if box.Clicked() {
						sheet.Focus(c)
						if on {
							edit()
						} else {
							pick(r, col)
						}
					}
				}
			})
		}
	})
	return root
}
