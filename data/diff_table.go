package data

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Change is how a DiffRow changed.
type Change uint8

// The four changes.
const (
	ChangeSame Change = iota
	ChangeAdded
	ChangeRemoved
	ChangeEdited
)

// DiffRow is cells now, and Before when edited.
type DiffRow struct {
	Change        Change
	Cells, Before []string
}

// DiffTableOptions tune a DiffTable.
type DiffTableOptions struct {
	// Added, Removed and Edited name the marks.
	Added, Removed, Edited string
}

// DiffTable tints rows by change, striking old values.
func DiffTable(c *ui.Context, label string, cols []Column, rows []DiffRow, o DiffTableOptions) ui.Element {
	columns("DiffTable", label, cols)
	for _, r := range rows {
		if r.Change > ChangeEdited || len(r.Cells) != len(cols) || (r.Change == ChangeEdited) != (len(r.Before) == len(cols)) {
			panic(fmt.Sprintf("data: DiffTable %q rows need a known change and a cell per column, before only when edited", label))
		}
	}
	def(&o.Added, "Added")
	def(&o.Removed, "Removed")
	def(&o.Edited, "Changed")
	t := c.Theme()
	look := map[Change]struct {
		ink        ui.Color
		mark, name string
	}{
		ChangeSame:    {t.TextMuted, "", ""},
		ChangeAdded:   {t.Success, "+", o.Added},
		ChangeRemoved: {t.Danger, "−", o.Removed},
		ChangeEdited:  {t.Warning, "~", o.Edited},
	}
	root := frame(c, label)
	root.Children(func() {
		head(c, root, "DiffTable", cols, nil, func() { ui.Box(c).Width(t.Space(4)).Shrink(0) })
		for i, r := range rows {
			c.Key(i)
			l := look[r.Change]
			row := ui.Row(c).Background(l.ink.Alpha(0.1 * unit(r.Change != ChangeSame)))
			kit.Arrive(row, ui.Pose{Y: t.Space(2), Fade: 1}, time.Duration(i)*30*time.Millisecond)
			line(c, row, cols, func() {
				ui.Text(c, l.mark).Bold().TextColor(l.ink).Width(t.Space(4)).Shrink(0).Label(l.name)
			}, func(col int) {
				if r.Change == ChangeEdited && r.Before[col] != r.Cells[col] {
					ui.Row(c).Gap(t.Space(1.5)).Children(func() {
						kit.Muted(c, r.Before[col]).Strikethrough().SingleLine()
						ui.Text(c, r.Cells[col]).Bold().SingleLine()
					})
					return
				}
				cell := ui.Text(c, r.Cells[col]).SingleLine()
				if r.Change == ChangeRemoved {
					cell.Strikethrough().TextColor(t.TextMuted)
				}
			})
		}
	})
	return root
}
