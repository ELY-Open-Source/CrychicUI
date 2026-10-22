package gallery

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/data"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

var dataEntries = []Entry{
	{656, "DataTable", "data", "Header clicks sort; the arrow turns.", demoDataTable},
	{657, "SelectTable", "data", "Row checks and a header check.", demoSelectTable},
	{658, "ExpandTable", "data", "Rows open a springing detail.", demoExpandTable},
	{659, "EditTable", "data", "Cells edited in place.", demoEditTable},
	{660, "TreeTable", "data", "An outline of nested rows.", demoTreeTable},
	{661, "VirtualList", "data", "Ten thousand rows, built in view.", demoVirtualList},
	{662, "GroupedList", "data", "Group titles pinned while scrolling.", demoGroupedList},
	{663, "ReorderList", "data", "Rows dragged to reorder.", demoReorderList},
	{664, "SwipeList", "data", "Swipe to archive or delete.", demoSwipeList},
	{665, "CardList", "data", "Cards that lift on hover.", demoCardList},
	{666, "GridList", "data", "A grid of tiles to choose.", demoGridList},
	{667, "CheckTree", "data", "Checks that cascade and go mixed.", demoCheckTree},
	{668, "DirTree", "data", "Folders with counts.", demoDirTree},
	{669, "JSONTree", "data", "JSON colored by type.", demoJSONTree},
	{670, "KVTable", "data", "Keys and values to copy.", demoKVTable},
	{671, "PropertyGrid", "data", "Typed property editors.", demoPropertyGrid},
	{672, "Pivot", "data", "Sums by row and column.", demoPivot},
	{673, "CompareTable", "data", "Features across plans.", demoCompareTable},
	{674, "DiffTable", "data", "Rows tinted by change.", demoDiffTable},
	{675, "Spreadsheet", "data", "A grid of editable cells.", demoSpreadsheet},
	{676, "ColumnPicker", "data", "Show, hide and order columns.", demoColumnPicker},
	{677, "FilterBuilder", "data", "Rules of field, operator, value.", demoFilterBuilder},
	{678, "QueryBuilder", "data", "Nested AND and OR groups.", demoQueryBuilder},
	{679, "Facets", "data", "Facet checks with counts.", demoFacets},
	{680, "Results", "data", "Results with matches marked.", demoResults},
	{681, "TableFooter", "data", "Rows per page and range.", demoTableFooter},
	{682, "BulkBar", "data", "Actions for the selection.", demoBulkBar},
	{683, "RowActions", "data", "A row's … menu.", demoRowActions},
}

var (
	dataPeople = [][]string{{"Ada Lovelace", "Analyst", "36"}, {"Alan Turing", "Engineer", "41"}, {"Grace Hopper", "Admiral", "45"},
		{"Linus Torvalds", "Maintainer", "28"}, {"Margaret Hamilton", "Director", "33"}}
	dataCols   = []data.Column{{Title: "Name"}, {Title: "Role"}, {Title: "Age", Width: 60, End: true}}
	dataFields = []string{"Name", "Role", "Age"}
	dataOps    = []string{"is", "contains", "is not"}
)

// personCell shows dataPeople's cell.
func personCell(c *ui.Context, rows [][]string) func(r, col int) {
	return func(r, col int) { ui.Text(c, rows[r][col]).SingleLine() }
}

func demoDataTable(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	sort := ui.Local(root, "sort", func() ui.SortOrder { return ui.SortOrder{} })
	root.Children(func() {
		rows := slices.Clone(dataPeople)
		col := slices.IndexFunc(dataCols, func(k data.Column) bool { return k.Title == sort.Column })
		if col >= 0 {
			slices.SortStableFunc(rows, func(a, b []string) int {
				if sort.Descending {
					a, b = b, a
				}
				return strings.Compare(a[col], b[col])
			})
		}
		if data.DataTable(c, "People", dataCols, len(rows), sort, personCell(c, rows), data.DataTableOptions{}).Changed() {
			log(fmt.Sprintf("sorted by %s, descending %v", sort.Column, sort.Descending))
		}
	})
}

func demoSelectTable(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	chosen := ui.Local(root, "chosen", func() ui.Selection[int] { return ui.Selection[int]{} })
	root.Children(func() {
		if data.SelectTable(c, "People", dataCols, len(dataPeople), chosen, personCell(c, dataPeople), data.SelectTableOptions{}).Changed() {
			log(fmt.Sprint(chosen.Len(), " chosen"))
		}
	})
}

func demoExpandTable(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	open := ui.Local(root, "open", func() ui.Selection[int] { return ui.Selection[int]{} })
	root.Children(func() {
		if data.ExpandTable(c, "People", dataCols, len(dataPeople), open, personCell(c, dataPeople), func(r int) {
			kit.Muted(c, fmt.Sprintf("%s joined as %s.", dataPeople[r][0], strings.ToLower(dataPeople[r][1])))
		}, data.ExpandTableOptions{}).Changed() {
			log(fmt.Sprint(open.Len(), " open"))
		}
	})
}

func demoEditTable(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	rows := ui.Local(root, "rows", func() [][]string {
		var out [][]string
		for _, r := range dataPeople[:3] {
			out = append(out, slices.Clone(r))
		}
		return out
	})
	at := ui.Local(root, "at", func() [2]int { return [2]int{-1, -1} })
	root.Children(func() {
		if data.EditTable(c, "People", dataCols, *rows, at, data.EditTableOptions{}).Changed() {
			log(fmt.Sprint("editing ", *at))
		}
	})
}

func demoTreeTable(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	open := ui.Local(root, "open", func() ui.Selection[string] { return ui.Selection[string]{} })
	rows := []data.TreeRow{
		{Cells: []string{"src", "3"}, Children: []data.TreeRow{{Cells: []string{"main.go", "1"}}, {Cells: []string{"ui", "2"},
			Children: []data.TreeRow{{Cells: []string{"view.go", "1"}}, {Cells: []string{"theme.go", "1"}}}}}},
		{Cells: []string{"go.mod", "1"}},
	}
	root.Children(func() {
		if data.TreeTable(c, "Files", []data.Column{{Title: "Name"}, {Title: "Files", Width: 60, End: true}}, rows, open, data.TreeTableOptions{}).Changed() {
			log(fmt.Sprint(open.Len(), " open"))
		}
	})
}

func demoVirtualList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	chosen := ui.Local(root, "chosen", func() int { return -1 })
	root.Children(func() {
		if data.VirtualList(c, "Invoices", 10000, chosen, func(i int) {
			ui.Text(c, fmt.Sprintf("Invoice %05d", i+1)).Grow(1)
			kit.Muted(c, fmt.Sprintf("$%d", (i*37)%900+100))
		}, data.VirtualListOptions{Height: 260}).Changed() {
			log(fmt.Sprint("chose ", *chosen))
		}
	})
}

func demoGroupedList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	chosen := ui.Local(root, "chosen", func() [2]int { return [2]int{-1, -1} })
	groups := []data.Group{{Title: "A", Items: []string{"Ada", "Alan", "Anita"}}, {Title: "G", Items: []string{"Grace", "Guido"}},
		{Title: "L", Items: []string{"Linus", "Lisa", "Leslie"}}, {Title: "M", Items: []string{"Margaret", "Mary"}}}
	root.Children(func() {
		if data.GroupedList(c, "Contacts", groups, chosen, data.GroupedListOptions{Height: 260}).Changed() {
			log(fmt.Sprint("chose ", *chosen))
		}
	})
}

func demoReorderList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	items := ui.Local(root, "items", func() []string { return []string{"Plan", "Design", "Build", "Test", "Ship"} })
	root.Children(func() {
		if data.ReorderList(c, "Steps", items, data.ReorderListOptions{Height: 240}).Changed() {
			log(strings.Join(*items, ", "))
		}
	})
}

func demoSwipeList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Gap(8)
	items := ui.Local(root, "items", func() []string { return []string{"Invoice from Acme", "Team lunch", "Weekly report", "Build failed"} })
	archived := ui.Local(root, "archived", func() []string { return nil })
	root.Children(func() {
		kit.Muted(c, "Swipe right to archive, left to delete.")
		if data.SwipeList(c, "Inbox", items, archived).Changed() {
			log(fmt.Sprintf("%d left, %d archived", len(*items), len(*archived)))
		}
	})
}

func demoCardList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	chosen := ui.Local(root, "chosen", func() int { return -1 })
	cards := []data.Card{{Title: "Starter", Text: "One seat, basic reports.", Icon: "user"},
		{Title: "Team", Text: "Five seats and shared boards.", Icon: "users"}, {Title: "Business", Text: "SSO and audit logs.", Icon: "building-2"}}
	root.Children(func() {
		if data.CardList(c, "Plans", cards, chosen).Changed() {
			log("chose " + cards[*chosen].Title)
		}
	})
}

func demoGridList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	chosen := ui.Local(root, "chosen", func() int { return -1 })
	tiles := []data.Tile{{Label: "Photos", Icon: "image"}, {Label: "Music", Icon: "music"}, {Label: "Docs", Icon: "file-text"},
		{Label: "Mail", Icon: "mail"}, {Label: "Maps", Icon: "map"}, {Label: "Notes", Icon: "notebook-pen"}}
	root.Children(func() {
		if data.GridList(c, "Apps", tiles, chosen, data.GridListOptions{Height: 220}).Changed() {
			log("chose " + tiles[*chosen].Label)
		}
	})
}
