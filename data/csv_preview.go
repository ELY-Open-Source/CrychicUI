package data

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CSVPreviewOptions tune a CSVPreview.
type CSVPreviewOptions struct {
	// Header names the switch; Column takes the number.
	Header, Column, Empty string
	// Height scrolls the rows within it; zero shows them all.
	Height float32
}

// CSVPreview tables CSV text; a switch sets the header.
func CSVPreview(c *ui.Context, label, text string, header *bool, o CSVPreviewOptions) ui.Element {
	if label == "" || header == nil {
		panic("data: CSVPreview needs a label and a header switch")
	}
	def(&o.Header, "First row is a header")
	def(&o.Column, "Column %d")
	def(&o.Empty, "No rows")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).Role(ui.RoleGroup).Label(label)
	recs, err := csv.NewReader(strings.NewReader(text)).ReadAll()
	last := ui.Local(root, "text", func() string { return "" })
	if text != *last {
		*last = text
		kit.Log("CSVPreview", "parsed", "rows", len(recs), "failed", err != nil)
	}
	fade := ui.ElementTransition{Animation: &ui.Snappy, Enter: &ui.Motion{}, Exit: &ui.Motion{}}
	return root.Children(func() {
		if choice.Toggle(c, o.Header, header).Changed() {
			root.SetChanged()
			kit.Log("CSVPreview", "header", "on", *header)
		}
		switch {
		case err != nil:
			ui.Text(c, err.Error()).TextColor(t.Danger)
			return
		case len(recs) == 0:
			kit.Muted(c, o.Empty)
			return
		}
		var cols []Column
		for i, name := range recs[0] {
			if !*header || name == "" {
				name = fmt.Sprintf(o.Column, i+1)
			}
			cols = append(cols, Column{Title: name})
		}
		rows := recs
		if *header {
			rows = recs[1:]
		}
		c.Key(*header)
		table := frame(c, label).Transition(fade)
		list := ui.Local(table, "list", func() ui.ListState { return ui.ListState{} })
		table.Children(func() {
			head(c, table, "CSVPreview", cols, nil, nil)
			body(c, table, list, len(rows), o.Height, func(i int) {
				line(c, hovered(c, false), cols, nil, func(col int) { ui.Text(c, rows[i][col]).SingleLine() })
			})
		})
	})
}
