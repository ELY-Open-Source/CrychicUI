package gallery

import (
	"fmt"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/files"
	"github.com/egoist/mygo/ui"
)

var filesEntries2 = []Entry{
	{ID: 959, Name: "Storage", Category: "files", Summary: "Usage by type with a growing bar.", Demo: demoFileStorage},
	{ID: 960, Name: "SyncRow", Category: "files", Summary: "File with sync state icon; spins while syncing.", Demo: demoFileSyncRow},
	{ID: 961, Name: "DownloadRow", Category: "files", Summary: "Name, bar, speed, time left, open; cue Success.", Demo: demoDownloadRow},
	{ID: 962, Name: "ZipList", Category: "files", Summary: "Archive entries with sizes; cue Toggle.", Demo: demoZipList},
	{ID: 963, Name: "ImageFiles", Category: "files", Summary: "Thumbnails with names; cue Select.", Demo: demoImageFiles},
	{ID: 964, Name: "Compress", Category: "files", Summary: "Before and after size as bars.", Demo: demoCompress},
	{ID: 965, Name: "FileSearch", Category: "files", Summary: "Name search with type filters; cue Select.", Demo: demoFileSearch},
	{ID: 966, Name: "ColumnBrowser", Category: "files", Summary: "Finder-style columns; columns slide; cue Select.", Demo: demoColumnBrowser},
	{ID: 967, Name: "FileCrumbs", Category: "files", Summary: "Path crumbs with folder menus; cue Open.", Demo: demoFileCrumbs},
	{ID: 968, Name: "QuotaWarning", Category: "files", Summary: "Near-full banner with upgrade; cue Warning.", Demo: demoQuotaWarning},
	{ID: 969, Name: "Starred", Category: "files", Summary: "Starred files; star toggles; cue Toggle.", Demo: demoStarred},
	{ID: 970, Name: "Conflict", Category: "files", Summary: "Keep mine, keep theirs, keep both; cue Select.", Demo: demoConflict},
	{ID: 971, Name: "SelectBar", Category: "files", Summary: "Selected count with actions; slides; cue Tap.", Demo: demoSelectBar},
	{ID: 972, Name: "FileTags", Category: "files", Summary: "Colored tags on files; cue Toggle.", Demo: demoFileTags},
	{ID: 973, Name: "UploadButtonRow", Category: "files", Summary: "Upload and new menu; cue Open.", Demo: demoUploadButtonRow},
	{ID: 974, Name: "FileTypeFilter", Category: "files", Summary: "Type chips with counts; cue Toggle.", Demo: demoFileTypeFilter},
	{ID: 975, Name: "Thumbs", Category: "files", Summary: "Large icon view with zoom slider; cue Tick.", Demo: demoFileThumbs},
}

func demoFileStorage(c *ui.Context, log func(string)) {
	ui.Column(c).Width(440).Children(func() {
		files.Storage(c, "Acme workspace", []files.StoragePart{
			{Name: "Documents", Bytes: 21_400_000_000},
			{Name: "Images", Bytes: 14_800_000_000},
			{Name: "Video", Bytes: 9_600_000_000},
			{Name: "Archives", Bytes: 3_100_000_000},
		}, 100_000_000_000, files.StorageOptions{})
	})
}

func demoFileSyncRow(c *ui.Context, log func(string)) {
	ui.Column(c).Width(440).Gap(8).Children(func() {
		for i, row := range []struct {
			item  files.FileItem
			state files.SyncState
		}{
			{files.FileItem{Name: "Q3 board report.pdf", Size: 2_480_000, Path: "Finance/Reports"}, files.Synced},
			{files.FileItem{Name: "office-floorplan.png", Size: 1_120_000, Path: "Facilities"}, files.Syncing},
			{files.FileItem{Name: "payroll-october.xlsx", Size: 310_000, Path: "HR"}, files.SyncFailed},
		} {
			c.Key(i)
			files.SyncRow(c, row.item, row.state, files.SyncRowOptions{})
		}
	})
}

func demoDownloadRow(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440).Gap(8)
	opened := boolState(root, "opened", false)
	root.Children(func() {
		files.DownloadRow(c, files.FileItem{Name: "customer-export.csv", Size: 48_000_000}, 0.42, 3_200_000, opened, files.DownloadRowOptions{})
		c.Key("done")
		if _, open := files.DownloadRow(c, files.FileItem{Name: "brand-assets.zip", Size: 126_000_000}, 1, 0, opened, files.DownloadRowOptions{}); open.Clicked() {
			log("opened download")
		}
	})
}

func demoZipList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	opened := intState(root, "opened", -1)
	root.Children(func() {
		if files.ZipList(c, "brand-assets.zip", []files.ZipEntry{
			{Name: "logo.svg", Size: 18_400, Compressed: 6_100},
			{Name: "guidelines.pdf", Size: 4_200_000, Compressed: 3_800_000},
			{Name: "palette.json", Size: 2_300, Compressed: 640},
			{Name: "launch-hero.png", Size: 2_900_000, Compressed: 2_860_000},
		}, opened, files.ZipListOptions{}).Changed() {
			log(fmt.Sprintf("open %d", *opened))
		}
	})
}

func demoImageFiles(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	selected := intState(root, "selected", 0)
	p := photos()
	images := []files.ImageFile{
		{Item: files.FileItem{Name: "harbor-dawn.jpg", Size: 3_400_000}, Image: p[0]},
		{Item: files.FileItem{Name: "lake-house.jpg", Size: 2_900_000}, Image: p[1]},
		{Item: files.FileItem{Name: "meadow.jpg", Size: 4_100_000}, Image: p[2]},
	}
	root.Children(func() {
		if files.ImageFiles(c, "Photos", images, selected).Changed() {
			log("selected " + images[*selected].Item.Name)
		}
	})
}

func demoCompress(c *ui.Context, log func(string)) {
	ui.Column(c).Width(440).Children(func() {
		files.Compress(c, "brand-assets.zip", 7_120_000, 2_650_000, files.CompressOptions{})
	})
}

func demoFileSearch(c *ui.Context, log func(string)) {
	list := []files.FileItem{
		{Name: "Q3 board report.pdf", Size: 2_480_000, Path: "Finance/Reports"},
		{Name: "Q3 revenue.xlsx", Size: 184_000, Path: "Finance"},
		{Name: "q3-launch-notes.md", Size: 6_200, Path: "Marketing/Launch"},
		{Name: "billing.go", Size: 14_800, Path: "Engineering/api"},
		{Name: "Q2 board report.pdf", Size: 2_310_000, Path: "Finance/Reports"},
	}
	root := ui.Column(c).Width(480)
	query, kind, selected := textState(root, "query", "q3"), textState(root, "kind", ""), intState(root, "selected", -1)
	root.Children(func() {
		kindBefore, queryBefore, selBefore := *kind, *query, *selected
		if files.FileSearch(c, "Search files", list, query, kind, selected, files.FileSearchOptions{}).Changed() {
			switch {
			case *kind != kindBefore:
				log("kind " + *kind)
			case *query != queryBefore:
				log(fmt.Sprintf("query %d chars", utf8.RuneCountInString(*query)))
			case *selected != selBefore:
				log("selected " + list[*selected].Name)
			}
		}
	})
}

func demoColumnBrowser(c *ui.Context, log func(string)) {
	columns := [][]string{
		{"Finance", "Marketing", "Engineering"},
		{"Reports", "Receipts", "Payroll"},
		{"Q3 board report.pdf", "Q2 board report.pdf", "Annual 2025.pdf"},
	}
	root := ui.Column(c).Width(520)
	selected := ui.Local(root, "selected", func() []int { return []int{0, 0, -1} })
	root.Children(func() {
		if files.ColumnBrowser(c, "Browse Acme", columns, selected).Changed() {
			log(fmt.Sprint(*selected))
		}
	})
}

func demoFileCrumbs(c *ui.Context, log func(string)) {
	parts := []string{"Acme", "Finance", "Reports", "2026"}
	root := ui.Column(c).Width(480)
	at := intState(root, "at", 3)
	root.Children(func() {
		if files.FileCrumbs(c, "Path", parts, at, files.FileCrumbsOptions{}).Changed() {
			log(fmt.Sprintf("opened %d", *at))
		}
	})
}

func demoQuotaWarning(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	upgraded := boolState(root, "upgraded", false)
	root.Children(func() {
		if r, _ := files.QuotaWarning(c, 91_400_000_000, 100_000_000_000, upgraded, files.QuotaWarningOptions{}); r.Changed() {
			log("upgraded")
		}
	})
}

func demoStarred(c *ui.Context, log func(string)) {
	list := []files.FileItem{
		{Name: "Q3 board report.pdf", Path: "Finance/Reports"},
		{Name: "brand-guidelines.pdf", Path: "Marketing"},
		{Name: "roadmap-2027.md", Path: "Product"},
		{Name: "office-floorplan.png", Path: "Facilities"},
	}
	root := ui.Column(c).Width(440)
	starred := ui.Local(root, "starred", func() []bool { return []bool{true, true, false, false} })
	root.Children(func() {
		if files.Starred(c, "Starred", list, starred, files.StarredOptions{}).Changed() {
			log(fmt.Sprint(*starred))
		}
	})
}

func demoConflict(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	choice := ui.Local(root, "choice", func() files.ConflictChoice { return files.ConflictNone })
	mine := files.FileItem{Name: "budget-2027.xlsx", Size: 212_000, Modified: fileDay(0, 1)}
	theirs := files.FileItem{Name: "budget-2027.xlsx", Size: 198_000, Modified: fileDay(0, 0)}
	root.Children(func() {
		if files.Conflict(c, "budget-2027.xlsx changed in two places", mine, theirs, choice, files.ConflictOptions{}).Changed() {
			log(fmt.Sprintf("choice %d", *choice))
		}
	})
}

func demoSelectBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	action := ui.Local(root, "action", func() files.SelectAction { return files.SelectNone })
	root.Children(func() {
		if files.SelectBar(c, 3, action, files.SelectBarOptions{}).Changed() {
			log(fmt.Sprintf("action %d", *action))
		}
	})
}

func demoFileTags(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	selected := ui.Local(root, "selected", func() []bool { return []bool{true, false, false, true, false} })
	root.Children(func() {
		if files.FileTags(c, "Tags", []string{"Invoices", "Urgent", "Contracts", "Q3", "Archive"}, selected).Changed() {
			log(fmt.Sprint(*selected))
		}
	})
}

func demoUploadButtonRow(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	action := ui.Local(root, "action", func() files.UploadAction { return files.UploadNone })
	root.Children(func() {
		if r, _, _ := files.UploadButtonRow(c, action, files.UploadButtonRowOptions{}); r.Changed() {
			log(fmt.Sprintf("action %d", *action))
		}
	})
}

func demoFileTypeFilter(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	selected := intState(root, "selected", -1)
	types := []files.TypeCount{
		{Kind: "Documents", Sample: "a.pdf", Count: 128},
		{Kind: "Spreadsheets", Sample: "a.xlsx", Count: 42},
		{Kind: "Images", Sample: "a.png", Count: 311},
		{Kind: "Code", Sample: "a.go", Count: 57},
		{Kind: "Archives", Sample: "a.zip", Count: 9},
	}
	root.Children(func() {
		if files.FileTypeFilter(c, "File types", types, selected).Changed() {
			log(fmt.Sprintf("selected %d", *selected))
		}
	})
}

func demoFileThumbs(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	zoom := ui.Local(root, "zoom", func() float64 { return 1 })
	selected := intState(root, "selected", -1)
	list := demoFileItems()
	root.Children(func() {
		if files.Thumbs(c, "Documents", list, zoom, selected, files.ThumbsOptions{}).Changed() {
			log(fmt.Sprintf("selected %d zoom %.2f", *selected, *zoom))
		}
	})
}
