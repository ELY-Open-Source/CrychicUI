package gallery

import (
	"fmt"
	"image"
	"image/color"

	"github.com/ELY-Open-Source/CrychicUI/files"
	"github.com/egoist/mygo/ui"
)

var filesEntries2 = []Entry{
	{ID: 959, Name: "Storage", Category: "files", Summary: "Usage by type with bars.", Demo: demoFileStorage},
	{ID: 960, Name: "SyncRow", Category: "files", Summary: "File with sync state icon; spins while syncing.", Demo: demoFileSyncRow},
	{ID: 961, Name: "DownloadRow", Category: "files", Summary: "Name, bar, speed, open; cue Success.", Demo: demoDownloadRow},
	{ID: 962, Name: "ZipList", Category: "files", Summary: "Archive entries with sizes; cue Toggle.", Demo: demoZipList},
	{ID: 963, Name: "ImageFiles", Category: "files", Summary: "Thumbnails with names; cue Select.", Demo: demoImageFiles},
	{ID: 964, Name: "Compress", Category: "files", Summary: "Before and after size with a bar.", Demo: demoCompress},
	{ID: 965, Name: "FileSearch", Category: "files", Summary: "Name search with type filters; cue Select.", Demo: demoFileSearch},
	{ID: 966, Name: "ColumnBrowser", Category: "files", Summary: "Finder-style columns; columns slide; cue Select.", Demo: demoColumnBrowser},
	{ID: 967, Name: "FileCrumbs", Category: "files", Summary: "Path crumbs with folder menus; cue Open.", Demo: demoFileCrumbs},
	{ID: 968, Name: "QuotaWarning", Category: "files", Summary: "Near-full banner with upgrade; cue Warning.", Demo: demoQuotaWarning},
	{ID: 969, Name: "Starred", Category: "files", Summary: "Starred files; star toggles; cue Toggle.", Demo: demoStarred},
	{ID: 970, Name: "Conflict", Category: "files", Summary: "Keep mine, keep theirs, keep both; cue Select.", Demo: demoConflict},
	{ID: 971, Name: "SelectBar", Category: "files", Summary: "3 selected with actions; slides; cue Tap.", Demo: demoSelectBar},
	{ID: 972, Name: "FileTags", Category: "files", Summary: "Colored tags on files; cue Toggle.", Demo: demoFileTags},
	{ID: 973, Name: "UploadButtonRow", Category: "files", Summary: "Upload and new menu; cue Open.", Demo: demoUploadButtonRow},
	{ID: 974, Name: "FileTypeFilter", Category: "files", Summary: "Type chips with counts; cue Toggle.", Demo: demoFileTypeFilter},
	{ID: 975, Name: "Thumbs", Category: "files", Summary: "Large icon view with zoom slider; cue Tick.", Demo: demoFileThumbs},
}

func demoFileStorage(c *ui.Context, log func(string)) {
	parts := []files.StoragePart{{Name: "Images", Bytes: 48000}, {Name: "Code", Bytes: 16000}, {Name: "Documents", Bytes: 8000}}
	files.Storage(c, "Storage", parts, 100000)
}

func demoFileSyncRow(c *ui.Context, log func(string)) {
	files.SyncRow(c, files.FileItem{Name: "notes.md", Size: 2400}, files.Syncing, files.SyncRowOptions{})
}

func demoDownloadRow(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	opened := boolState(root, "opened", false)
	root.Children(func() {
		item := files.FileItem{Name: "gallery.wasm", Size: 33545277}
		if files.DownloadRow(c, item, 1, 1200000, opened, files.DownloadRowOptions{}).Changed() {
			log("opened download")
		}
	})
}

func demoZipList(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	opened := intState(root, "opened", -1)
	entries := []files.ZipEntry{{Name: "main.go", Size: 8400, Compressed: 2600}, {Name: "README.md", Size: 3200, Compressed: 1200}}
	root.Children(func() {
		if files.ZipList(c, "Archive", entries, opened, files.ZipListOptions{}).Changed() {
			log("open " + fmt.Sprint(*opened))
		}
	})
}

func demoImageFiles(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := intState(root, "selected", -1)
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 90, G: 120, B: 220, A: 255})
	bitmap := ui.NewBitmap(img)
	images := []files.ImageFile{
		{Item: files.FileItem{Name: "one.png", Size: 2000}, Image: bitmap},
		{Item: files.FileItem{Name: "two.png", Size: 3000}, Image: bitmap},
	}
	root.Children(func() {
		if files.ImageFiles(c, "Images", images, selected).Changed() {
			log("selected " + fmt.Sprint(*selected))
		}
	})
}

func demoCompress(c *ui.Context, log func(string)) {
	files.Compress(c, "Compress", 100000, 42000, files.CompressOptions{})
}

func demoFileSearch(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	query, kind := textState(root, "query", ""), textState(root, "kind", "")
	selected := intState(root, "selected", -1)
	root.Children(func() {
		if files.FileSearch(c, "Search", demoFileItems(), query, kind, selected, files.FileSearchOptions{}).Changed() {
			log(fmt.Sprintf("selected %d kind %s", *selected, *kind))
		}
	})
}

func demoColumnBrowser(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := ui.Local(root, "selected", func() []int { return []int{0, 0} })
	columns := [][]string{{"home", "work"}, {"docs", "images"}}
	root.Children(func() {
		if files.ColumnBrowser(c, "Browse", columns, selected).Changed() {
			log(fmt.Sprint(*selected))
		}
	})
}

func demoFileCrumbs(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	at := intState(root, "at", 0)
	root.Children(func() {
		if files.FileCrumbs(c, []string{"home", "work", "docs"}, at, files.FileCrumbsOptions{}).Changed() {
			log("opened " + fmt.Sprint(*at))
		}
	})
}

func demoQuotaWarning(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	upgraded := boolState(root, "upgraded", false)
	root.Children(func() {
		if files.QuotaWarning(c, 92000, 100000, upgraded, files.QuotaWarningOptions{}).Changed() {
			log("upgraded")
		}
	})
}

func demoStarred(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	starred := ui.Local(root, "starred", func() []bool { return []bool{true, false, false} })
	root.Children(func() {
		if files.Starred(c, "Starred", demoFileItems(), starred, files.StarredOptions{}).Changed() {
			log(fmt.Sprint(*starred))
		}
	})
}

func demoConflict(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	choice := ui.Local(root, "choice", func() files.ConflictChoice { return files.ConflictNone })
	root.Children(func() {
		if files.Conflict(c, "Conflict", choice, files.ConflictOptions{}).Changed() {
			log("choice " + fmt.Sprint(*choice))
		}
	})
}

func demoSelectBar(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	action := ui.Local(root, "action", func() files.SelectAction { return files.SelectNone })
	root.Children(func() {
		if files.SelectBar(c, 3, action, files.SelectBarOptions{}).Changed() {
			log("action " + fmt.Sprint(*action))
			*action = files.SelectNone
		}
	})
}

func demoFileTags(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := ui.Local(root, "selected", func() []bool { return []bool{true, false, false} })
	root.Children(func() {
		if files.FileTags(c, "Tags", []string{"work", "urgent", "design"}, selected).Changed() {
			log(fmt.Sprint(*selected))
		}
	})
}

func demoUploadButtonRow(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	action := ui.Local(root, "action", func() files.UploadAction { return files.UploadNone })
	root.Children(func() {
		files.UploadButtonRow(c, action, files.UploadButtonRowOptions{})
		if *action != files.UploadNone {
			log("action " + fmt.Sprint(*action))
			*action = files.UploadNone
		}
	})
}

func demoFileTypeFilter(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := intState(root, "selected", -1)
	types := []files.TypeCount{{Kind: "Documents", Count: 12}, {Kind: "Code", Count: 8}, {Kind: "Images", Count: 5}}
	root.Children(func() {
		if files.FileTypeFilter(c, "Types", types, selected).Changed() {
			log("selected " + fmt.Sprint(*selected))
		}
	})
}

func demoFileThumbs(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	zoom := numberState(root, "zoom", 1)
	selected := intState(root, "selected", -1)
	root.Children(func() {
		if files.Thumbs(c, "Thumbs", demoFileItems(), zoom, selected, files.ThumbsOptions{}).Changed() {
			log(fmt.Sprintf("selected %d zoom %.2f", *selected, *zoom))
		}
	})
}
