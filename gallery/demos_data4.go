package gallery

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/data"
	"github.com/egoist/mygo/ui"
)

func demoAuditTrail(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() {
		data.AuditTrail(c, "Changes", []data.Audit{{Who: "Grace Hopper", What: "Renamed the report", When: "2 min ago"},
			{Who: "Alan Turing", What: "Shared it with the team", When: "1 hour ago"}, {Who: "Ada Lovelace", What: "Created the report", When: "Yesterday"}})
	})
}

func demoVersionList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	current := intState(root, "current", 2)
	versions := []data.Version{{Name: "v1", When: "May 1", Note: "First draft"}, {Name: "v2", When: "May 3", Note: "Copy edits"}, {Name: "v3", When: "May 6", Note: "Final"}}
	root.Children(func() {
		if data.VersionList(c, "History", versions, current, data.VersionListOptions{}).Changed() {
			log("restored " + versions[*current].Name)
		}
	})
}

func demoCSVPreview(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	header := flag(root, "header")
	root.Children(func() {
		if data.CSVPreview(c, "Import", "name,role,age\nAda,Analyst,36\nAlan,Engineer,41", header, data.CSVPreviewOptions{}).Changed() {
			log(fmt.Sprintf("header %v", *header))
		}
	})
}

func demoColumnMapper(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	mapping := ui.Local(root, "mapping", func() []string { return []string{"Name", "", ""} })
	root.Children(func() {
		if data.ColumnMapper(c, "Columns", []string{"full_name", "job", "years"}, dataFields, mapping, data.ColumnMapperOptions{}).Changed() {
			log(fmt.Sprintf("mapping %q", *mapping))
		}
	})
}

func demoExportMenu(c *ui.Context, log func(string)) {
	root := ui.Row(c).Width(420)
	chosen := intState(root, "chosen", -1)
	open := flag(root, "open")
	exports := []data.Export{{Name: "Copy as CSV", Data: "name,age\nAda,36"}, {Name: "Copy as JSON", Data: `[{"name":"Ada","age":36}]`}}
	root.Children(func() {
		if data.ExportMenu(c, "Export", exports, chosen, open).Changed() {
			log(fmt.Sprintf("open %v, chose %d", *open, *chosen))
		}
	})
}

func demoTableToolbar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(760)
	query := ui.Local(root, "query", func() string { return "" })
	d := ui.Local(root, "density", func() data.Density { return data.DensityDefault })
	root.Children(func() {
		bar, filters, columns := data.TableToolbar(c, "People", query, d, data.TableToolbarOptions{})
		switch {
		case filters.Clicked():
			log("filters")
		case columns.Clicked():
			log("columns")
		case bar.Changed():
			log(fmt.Sprintf("query length %d, density %d", len([]rune(*query)), *d))
		}
	})
}

func demoPinnedTable(c *ui.Context, log func(string)) {
	cols := []data.Column{{Title: "Name", Width: 150}}
	for _, m := range []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug"} {
		cols = append(cols, data.Column{Title: m, Width: 80, End: true})
	}
	ui.Column(c).Width(520).Children(func() {
		data.PinnedTable(c, "Sales", cols, 4, func(r, col int) {
			if col == 0 {
				ui.Text(c, dataPeople[r][0]).SingleLine()
				return
			}
			ui.Text(c, strconv.Itoa((r+1)*col*7)).SingleLine()
		})
	})
}

func demoTotalsRow(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480).Gap(12)
	qty := ui.Local(root, "qty", func() int64 { return 12 })
	cols := []data.Column{{Title: "Item"}, {Title: "Qty", Width: 80, End: true}, {Title: "Cost", Width: 80, End: true}}
	root.Children(func() {
		data.TotalsRow(c, "Totals", cols, []int64{*qty, *qty * 9}, data.TotalsRowOptions{})
		if button.Outline(c, "Add one").Clicked() {
			*qty++
			log(fmt.Sprintf("qty %d", *qty))
		}
	})
}

func demoHeatTable(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520).Gap(12)
	shift := ui.Local(root, "shift", func() int { return 0 })
	root.Children(func() {
		var values [][]float64
		for r := range 3 {
			var row []float64
			for col := range 4 {
				row = append(row, float64((r*4+col+*shift)%12))
			}
			values = append(values, row)
		}
		data.HeatTable(c, "Load", []string{"Mon", "Tue", "Wed"}, []string{"6 AM", "Noon", "6 PM", "Midnight"}, values, data.HeatTableOptions{})
		if button.Outline(c, "Shift").Clicked() {
			*shift++
			log(fmt.Sprintf("shift %d", *shift))
		}
	})
}

func demoMatrix(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	marks := ui.Local(root, "marks", func() ui.Selection[[2]int] { return ui.Selection[[2]int]{} })
	root.Children(func() {
		if data.Matrix(c, "Access", []string{"Ada", "Alan", "Grace"}, []string{"Read", "Write", "Admin"}, marks, data.MatrixOptions{}).Changed() {
			log(fmt.Sprintf("%d marks", marks.Len()))
		}
	})
}

func demoDiffList(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		data.DiffList(c, "Changes", []data.DiffItem{{Change: data.ChangeAdded, Text: "logo.svg"}, {Change: data.ChangeEdited, Text: "index.html"},
			{Change: data.ChangeRemoved, Text: "old-logo.png"}, {Change: data.ChangeSame, Text: "styles.css"}}, data.DiffListOptions{})
	})
}

func demoInventory(c *ui.Context, log func(string)) {
	ui.Column(c).Width(620).Children(func() {
		data.Inventory(c, "Stock", []data.Stock{{SKU: "BL-100", Name: "Bolts", Count: 12, Max: 100}, {SKU: "NT-200", Name: "Nuts", Count: 35, Max: 100},
			{SKU: "WS-300", Name: "Washers", Count: 88, Max: 100}}, data.InventoryOptions{})
	})
}
