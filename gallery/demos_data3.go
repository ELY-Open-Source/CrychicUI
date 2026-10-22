package gallery

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/data"
	"github.com/egoist/mygo/ui"
)

var dataEntries2 = []Entry{
	{684, "EmptyTable", "data", "An empty table with an action.", demoEmptyTable},
	{685, "LoadingTable", "data", "Sheen rows while loading.", demoLoadingTable},
	{686, "StatList", "data", "Label and value rows; values roll.", demoStatList},
	{687, "RankList", "data", "Ranked rows with growing bars.", demoRankList},
	{688, "TimelineList", "data", "Dated rows on a line.", demoTimelineList},
	{689, "LogList", "data", "Leveled lines that follow the end.", demoLogList},
	{690, "Feed", "data", "Cards arriving at the top.", demoFeed},
	{691, "MasterDetail", "data", "A list beside its detail.", demoMasterDetail},
	{692, "Kanban", "data", "Cards dragged between lanes.", demoKanban},
	{693, "AccordionList", "data", "One detail open at a time.", demoAccordionList},
	{694, "TagCloud", "data", "Words sized by weight.", demoTagCloud},
	{695, "DataCards", "data", "Records as cards in a grid.", demoDataCards},
	{696, "RecordView", "data", "One record in sections.", demoRecordView},
	{697, "AuditTrail", "data", "Who changed what, and when.", demoAuditTrail},
	{698, "VersionList", "data", "Versions to restore.", demoVersionList},
	{699, "CSVPreview", "data", "Parsed CSV with a header switch.", demoCSVPreview},
	{700, "ColumnMapper", "data", "Source columns mapped to fields.", demoColumnMapper},
	{701, "ExportMenu", "data", "Formats that copy the data.", demoExportMenu},
	{702, "TableToolbar", "data", "Search, filters, columns, density.", demoTableToolbar},
	{703, "DensitySwitch", "data", "Compact, default or roomy rows.", demoDensitySwitch},
	{704, "PinnedTable", "data", "A first column fixed while scrolling.", demoPinnedTable},
	{705, "TotalsRow", "data", "A footer of rolling sums.", demoTotalsRow},
	{706, "HeatTable", "data", "Cells tinted by value.", demoHeatTable},
	{707, "Matrix", "data", "Checks by row and column.", demoMatrix},
	{708, "Paginated", "data", "A table paged with cross-fades.", demoPaginated},
	{709, "DiffList", "data", "Added, removed and changed items.", demoDiffList},
	{710, "Inventory", "data", "Stock bars colored by level.", demoInventory},
}

func demoEmptyTable(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	root.Children(func() {
		if _, act := data.EmptyTable(c, "People", dataCols, "No people yet.", "Add person", data.EmptyTableOptions{Icon: "users"}); act.Clicked() {
			log("add person")
		}
	})
}

func demoLoadingTable(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() { data.LoadingTable(c, "People", dataCols, 4) })
}

func demoStatList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	visits := ui.Local(root, "visits", func() int64 { return 1280 })
	root.Children(func() {
		data.StatList(c, "Today", []data.Stat{{Label: "Visits", Value: *visits}, {Label: "Orders", Value: 42, Unit: "paid"}})
		if button.Outline(c, "Add a visit").Clicked() {
			*visits++
			log(fmt.Sprintf("visits %d", *visits))
		}
	})
}

func demoRankList(c *ui.Context, log func(string)) {
	ui.Column(c).Width(560).Children(func() {
		data.RankList(c, "Top pages", []data.Ranked{{Label: "/pricing", Value: 4210}, {Label: "/docs", Value: 3120}, {Label: "/blog", Value: 1800}, {Label: "/about", Value: 640}}, data.RankListOptions{})
	})
}

func demoTimelineList(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() {
		data.TimelineList(c, "History", []data.Moment{{Date: "Mar 2", Title: "Founded", Text: "Two people in a garage."},
			{Date: "Jun 9", Title: "First release", Text: "Version 1.0 ships."}, {Date: "Oct 21", Title: "Ten thousand users"}})
	})
}

func demoLogList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(560).Gap(12)
	lines := ui.Local(root, "lines", func() []data.LogLine {
		return []data.LogLine{{Level: data.LevelInfo, Time: "09:00:01", Text: "server started"}, {Level: data.LevelDebug, Time: "09:00:02", Text: "cache warm"}}
	})
	root.Children(func() {
		data.LogList(c, "Server log", *lines, data.LogListOptions{Height: 220})
		if button.Outline(c, "Add a line").Clicked() {
			n := len(*lines)
			*lines = append(*lines, data.LogLine{Level: data.Level(n % 4), Time: fmt.Sprintf("09:00:%02d", n+1), Text: fmt.Sprintf("event %d", n)})
			log(fmt.Sprintf("%d lines", len(*lines)))
		}
	})
}

func demoFeed(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480).Gap(12)
	posts := ui.Local(root, "posts", func() []data.Post { return []data.Post{{ID: "1", Title: "Welcome", Text: "Your feed starts here."}} })
	root.Children(func() {
		if button.Outline(c, "New post").Clicked() {
			n := len(*posts) + 1
			*posts = append([]data.Post{{ID: fmt.Sprint(n), Title: fmt.Sprintf("Post %d", n), Text: "Fresh from the server."}}, *posts...)
			log(fmt.Sprintf("%d posts", n))
		}
		data.Feed(c, "News", *posts)
	})
}

func demoMasterDetail(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(620)
	chosen := intState(root, "chosen", -1)
	root.Children(func() {
		names := []string{dataPeople[0][0], dataPeople[1][0], dataPeople[2][0]}
		if data.MasterDetail(c, "People", names, chosen, func(i int) {
			ui.Text(c, dataPeople[i][0]).Bold()
			ui.Text(c, dataPeople[i][1])
		}, data.MasterDetailOptions{}).Changed() {
			log(fmt.Sprintf("chose %d", *chosen))
		}
	})
}

func demoKanban(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(760)
	lanes := ui.Local(root, "lanes", func() []data.Lane {
		return []data.Lane{{Title: "To do", Cards: []string{"Write brief", "Sketch logo"}}, {Title: "Doing", Cards: []string{"Build site"}}, {Title: "Done", Cards: []string{"Kickoff"}}}
	})
	root.Children(func() {
		if data.Kanban(c, "Launch", lanes).Changed() {
			var counts []int
			for _, l := range *lanes {
				counts = append(counts, len(l.Cards))
			}
			log(fmt.Sprintf("lanes %v", counts))
		}
	})
}

func demoAccordionList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	open := intState(root, "open", -1)
	answers := []string{"Orders ship in two days.", "Returns are free for thirty days.", "Write to help@example.com."}
	root.Children(func() {
		if data.AccordionList(c, "FAQ", []string{"Shipping", "Returns", "Support"}, open, func(i int) { ui.Text(c, answers[i]) }).Changed() {
			log(fmt.Sprintf("open %d", *open))
		}
	})
}

func demoTagCloud(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	chosen := intState(root, "chosen", -1)
	tags := []data.Tag{{Text: "go", Weight: 40}, {Text: "rust", Weight: 22}, {Text: "wasm", Weight: 30}, {Text: "ui", Weight: 35}, {Text: "audio", Weight: 9}, {Text: "motion", Weight: 18}, {Text: "tests", Weight: 12}}
	root.Children(func() {
		if data.TagCloud(c, "Topics", tags, chosen).Changed() {
			log("chose " + tags[*chosen].Text)
		}
	})
}

func demoDataCards(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(720)
	chosen := intState(root, "chosen", -1)
	var entries []data.Entry
	for _, p := range dataPeople[:3] {
		entries = append(entries, data.Entry{Title: p[0], Fields: []data.Pair{{Key: "Role", Value: p[1]}, {Key: "Age", Value: p[2]}}})
	}
	root.Children(func() {
		if data.DataCards(c, "People", entries, chosen).Changed() {
			log(fmt.Sprintf("chose %d", *chosen))
		}
	})
}

func demoRecordView(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		data.RecordView(c, "Order 1042", []data.Section{
			{Title: "Buyer", Fields: []data.Pair{{Key: "Name", Value: "Ada Lovelace"}, {Key: "Email", Value: "ada@example.com"}}},
			{Title: "Payment", Fields: []data.Pair{{Key: "Card", Value: "•••• 4242"}, {Key: "Total", Value: "$84.00"}}},
		})
	})
}

func demoDensitySwitch(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	d := ui.Local(root, "density", func() data.Density { return data.DensityDefault })
	t := c.Theme()
	root.Children(func() {
		if data.DensitySwitch(c, "Rows", d, data.DensitySwitchOptions{}).Changed() {
			log(fmt.Sprintf("density %d", *d))
		}
		ui.Column(c).FillWidth().Radius(t.Radius).Border(1, t.Border).Dividers(1, t.Border).Children(func() {
			for i, p := range dataPeople[:3] {
				c.Key(i)
				row := ui.Row(c).FillWidth()
				row.Padding(row.AnimateTo("pad", d.Pad(t), ui.Snappy), 12).Children(func() { ui.Text(c, p[0]) })
			}
		})
	})
}

func demoPaginated(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	page := intState(root, "page", 0)
	rows := slices.Repeat(dataPeople, 5)
	root.Children(func() {
		if data.Paginated(c, "People", dataCols, len(rows), personCell(c, rows), page, data.PaginatedOptions{Size: 5}).Changed() {
			log(fmt.Sprintf("page %d", *page+1))
		}
	})
}
