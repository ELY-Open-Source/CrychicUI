package gallery

import (
	"encoding/json"
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/data"
	"github.com/egoist/mygo/ui"
)

func demoCheckTree(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	checked := ui.Local(root, "checked", func() ui.Selection[string] { return ui.Selection[string]{} })
	nodes := []data.Node{{Label: "Notifications", Children: []data.Node{{Label: "Email"}, {Label: "Push"}, {Label: "SMS"}}},
		{Label: "Privacy", Children: []data.Node{{Label: "Show profile"}, {Label: "Share usage"}}}}
	root.Children(func() {
		if data.CheckTree(c, "Settings", nodes, checked).Changed() {
			log(fmt.Sprint(checked.Len(), " checked"))
		}
	})
}

func demoDirTree(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	open := ui.Local(root, "open", func() ui.Selection[string] { return ui.Selection[string]{} })
	chosen := ui.Local(root, "chosen", func() string { return "" })
	dirs := []data.Dir{{Name: "Documents", Count: 128, Children: []data.Dir{{Name: "Invoices", Count: 42}, {Name: "Letters", Count: 9}}},
		{Name: "Pictures", Count: 2048, Children: []data.Dir{{Name: "2025", Count: 860}}}, {Name: "Music", Count: 311}}
	root.Children(func() {
		if data.DirTree(c, "Folders", dirs, open, chosen, data.DirTreeOptions{}).Changed() {
			log(fmt.Sprintf("chose %q, %d open", *chosen, open.Len()))
		}
	})
}

var demoJSON = func() any {
	var v any
	if err := json.Unmarshal([]byte(`{"user":{"name":"Ada","age":36,"admin":true},"tags":["math","engines"],"manager":null}`), &v); err != nil {
		panic(err)
	}
	return v
}()

func demoJSONTree(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	open := ui.Local(root, "open", func() ui.Selection[string] { s := ui.Selection[string]{}; s.Add("$"); return s })
	root.Children(func() {
		if data.JSONTree(c, "Response", demoJSON, open, data.JSONTreeOptions{}).Changed() {
			log(fmt.Sprint(open.Len(), " open"))
		}
	})
}

func demoKVTable(c *ui.Context, log func(string)) {
	pairs := []data.Pair{{Key: "Host", Value: "api.crychic.dev"}, {Key: "Region", Value: "eu-west-1"}, {Key: "Token", Value: "pk_live_42"}}
	ui.Column(c).Width(460).Children(func() {
		if data.KVTable(c, "Connection", pairs, data.KVTableOptions{}).Changed() {
			log("copied")
		}
	})
}

func demoPropertyGrid(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	name := ui.Local(root, "name", func() string { return "Header" })
	visible := ui.Local(root, "visible", func() bool { return true })
	opacity := ui.Local(root, "opacity", func() float64 { return 80 })
	root.Children(func() {
		if data.PropertyGrid(c, "Layer", []data.Property{{Name: "Name", Value: name}, {Name: "Visible", Value: visible},
			{Name: "Opacity", Value: opacity, Min: 0, Max: 100, Step: 5}}).Changed() {
			log(fmt.Sprintf("visible %v, opacity %.0f, name %d characters", *visible, *opacity, len(*name)))
		}
	})
}

func demoPivot(c *ui.Context, log func(string)) {
	records := []data.Record{{Row: "North", Col: "Q1", Value: 120}, {Row: "North", Col: "Q2", Value: 140}, {Row: "South", Col: "Q1", Value: 90},
		{Row: "South", Col: "Q2", Value: 110}, {Row: "West", Col: "Q1", Value: 75}, {Row: "West", Col: "Q2", Value: 95}}
	ui.Column(c).Width(460).Children(func() { data.Pivot(c, "Region", records, data.PivotOptions{}) })
}

func demoCompareTable(c *ui.Context, log func(string)) {
	features := []data.Feature{{Name: "Projects", Has: []bool{true, true, true}}, {Name: "Shared boards", Has: []bool{false, true, true}},
		{Name: "SSO", Has: []bool{false, false, true}}, {Name: "Audit log", Has: []bool{false, false, true}}}
	ui.Column(c).Width(500).Children(func() {
		data.CompareTable(c, "Feature", []string{"Free", "Team", "Business"}, features, data.CompareTableOptions{})
	})
}

func demoDiffTable(c *ui.Context, log func(string)) {
	cols := []data.Column{{Title: "Setting"}, {Title: "Value"}}
	rows := []data.DiffRow{{Cells: []string{"region", "eu-west-1"}}, {Change: data.ChangeEdited, Cells: []string{"replicas", "5"}, Before: []string{"replicas", "3"}},
		{Change: data.ChangeAdded, Cells: []string{"cache", "on"}}, {Change: data.ChangeRemoved, Cells: []string{"debug", "true"}}}
	ui.Column(c).Width(460).Children(func() { data.DiffTable(c, "Deploy", cols, rows, data.DiffTableOptions{}) })
}

func demoSpreadsheet(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	cells := ui.Local(root, "cells", func() [][]string {
		return [][]string{{"Item", "Qty", "Price"}, {"Pens", "12", "1.50"}, {"Paper", "3", "4.20"}, {"Ink", "2", "9.90"}}
	})
	at := ui.Local(root, "at", func() [2]int { return [2]int{1, 1} })
	root.Children(func() {
		if data.Spreadsheet(c, "Order", *cells, at).Changed() {
			log(fmt.Sprint("cell ", *at))
		}
	})
}

func demoColumnPicker(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	order := ui.Local(root, "order", func() []string { return []string{"Name", "Role", "Age", "Team", "Started"} })
	hidden := ui.Local(root, "hidden", func() ui.Selection[string] { return ui.Selection[string]{} })
	root.Children(func() {
		if data.ColumnPicker(c, "Columns", order, hidden).Changed() {
			log(fmt.Sprintf("%d hidden, first %s", hidden.Len(), (*order)[0]))
		}
	})
}

func demoFilterBuilder(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(560)
	rules := ui.Local(root, "rules", func() []data.Rule { return []data.Rule{{Field: "Role", Op: "is", Value: "Engineer"}} })
	root.Children(func() {
		if data.FilterBuilder(c, "Filters", dataFields, dataOps, rules, data.FilterBuilderOptions{}).Changed() {
			log(fmt.Sprint(len(*rules), " rules"))
		}
	})
}

func demoQueryBuilder(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(540)
	q := ui.Local(root, "query", func() data.Query {
		return data.Query{Rules: []data.Rule{{Field: "Age", Op: "is not", Value: "30"}},
			Groups: []data.Query{{Any: true, Rules: []data.Rule{{Field: "Role", Op: "is", Value: "Admiral"}}}}}
	})
	root.Children(func() {
		if data.QueryBuilder(c, "Query", dataFields, dataOps, q, data.QueryBuilderOptions{}).Changed() {
			log(fmt.Sprintf("any %v, %d rules, %d groups", q.Any, len(q.Rules), len(q.Groups)))
		}
	})
}

func demoFacets(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(260)
	chosen := ui.Local(root, "chosen", func() ui.Selection[string] { return ui.Selection[string]{} })
	facets := []data.Facet{{Title: "Brand", Options: []data.FacetOption{{Label: "Acme", Count: 24}, {Label: "Globex", Count: 9}}},
		{Title: "Size", Options: []data.FacetOption{{Label: "S", Count: 12}, {Label: "M", Count: 31}, {Label: "L", Count: 7}}}}
	root.Children(func() {
		if data.Facets(c, "Filters", facets, chosen).Changed() {
			log(fmt.Sprint(chosen.Len(), " chosen"))
		}
	})
}

func demoResults(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	chosen := ui.Local(root, "chosen", func() int { return -1 })
	results := []data.Result{{Title: "Go tour", Snippet: "A tour of the Go language, step by step."},
		{Title: "Effective Go", Snippet: "Tips for writing clear, idiomatic Go code."}, {Title: "Learn Rust", Snippet: "Ownership and borrowing."}}
	root.Children(func() {
		if data.Results(c, "Search", "go", results, chosen, data.ResultsOptions{}).Changed() {
			log("opened " + results[*chosen].Title)
		}
	})
}

func demoTableFooter(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	page := ui.Local(root, "page", func() int { return 0 })
	size := ui.Local(root, "size", func() int { return 10 })
	root.Children(func() {
		if data.TableFooter(c, 128, page, size, data.TableFooterOptions{}).Changed() {
			log(fmt.Sprintf("page %d, %d per page", *page+1, *size))
		}
	})
}

func demoBulkBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520).Gap(12)
	chosen := ui.Local(root, "chosen", func() ui.Selection[int] { return ui.Selection[int]{} })
	pick := ui.Local(root, "pick", func() int { return -1 })
	actions := []data.Action{{Label: "Archive", Icon: "archive"}, {Label: "Delete", Icon: "trash-2"}}
	root.Children(func() {
		if data.SelectTable(c, "People", dataCols, 3, chosen, personCell(c, dataPeople), data.SelectTableOptions{}).Changed() {
			log(fmt.Sprintf("selected %d", chosen.Len()))
		}
		r, clear := data.BulkBar(c, chosen.Len(), actions, pick, data.BulkBarOptions{})
		switch {
		case clear.Clicked():
			chosen.Clear()
			log("cleared")
		case r.Changed():
			log(fmt.Sprintf("%s %d", actions[*pick].Label, chosen.Len()))
		}
	})
}

func demoRowActions(c *ui.Context, log func(string)) {
	root := ui.Row(c).Width(420).Gap(8).AlignItems(ui.Center)
	open := flag(root, "open")
	chosen := ui.Local(root, "chosen", func() int { return -1 })
	actions := []data.Action{{Label: "Rename", Icon: "pencil"}, {Label: "Duplicate", Icon: "copy"}, {Label: "Delete", Icon: "trash-2"}}
	root.Children(func() {
		ui.Text(c, "Quarterly report.pdf").Grow(1)
		if data.RowActions(c, "More actions", actions, chosen, open).Changed() {
			log(fmt.Sprintf("open %v, chose %d", *open, *chosen))
		}
	})
}
