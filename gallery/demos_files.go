package gallery

import (
	"fmt"
	"slices"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/files"
	"github.com/egoist/mygo/ui"
)

var filesEntries = []Entry{
	{ID: 941, Name: "FileList", Category: "files", Summary: "Rows with icon, name, size, date; sort; cue Select.", Demo: demoFileList},
	{ID: 942, Name: "FileGrid", Category: "files", Summary: "Icon tiles with names; cue Select.", Demo: demoFileGrid},
	{ID: 943, Name: "FileIcon", Category: "files", Summary: "Tile chosen by extension; swaps.", Demo: demoFileIcon},
	{ID: 944, Name: "Dropzone", Category: "files", Summary: "Area for dropped files; border animates; cue Drop.", Demo: demoDropzone},
	{ID: 945, Name: "UploadList", Category: "files", Summary: "Files with bars and cancel; cue Close.", Demo: demoUploadList},
	{ID: 946, Name: "FilePreview", Category: "files", Summary: "Text, image or info preview by type.", Demo: demoFilePreview},
	{ID: 947, Name: "FileCard", Category: "files", Summary: "Preview band, name, size and menu; lifts.", Demo: demoFileCard},
	{ID: 948, Name: "FolderCard", Category: "files", Summary: "Folder with item count; opens; cue Open.", Demo: demoFolderCard},
	{ID: 949, Name: "PathBar", Category: "files", Summary: "Folders of a path; cue Tap.", Demo: demoPathBar},
	{ID: 950, Name: "NewFolder", Category: "files", Summary: "Inline row naming a new folder; cue Pop.", Demo: demoNewFolder},
	{ID: 951, Name: "RenameRow", Category: "files", Summary: "Name edited in place; cue Success.", Demo: demoRenameRow},
	{ID: 952, Name: "MoveDialog", Category: "files", Summary: "Folder tree choosing a target; cue Drop.", Demo: demoMoveDialog},
	{ID: 953, Name: "ShareDialog", Category: "files", Summary: "People with roles and a link; cue Pop.", Demo: demoShareDialog},
	{ID: 954, Name: "Permissions", Category: "files", Summary: "Owner, editor, viewer rows; cue Select.", Demo: demoPermissions},
	{ID: 955, Name: "Trash", Category: "files", Summary: "Deleted items with restore; cue Success.", Demo: demoTrash},
	{ID: 956, Name: "Properties", Category: "files", Summary: "Size, type, dates, path.", Demo: demoProperties},
	{ID: 957, Name: "Versions", Category: "files", Summary: "File versions with restore; cue Success.", Demo: demoVersions},
	{ID: 958, Name: "RecentFiles", Category: "files", Summary: "Recently opened rows; cue Tap.", Demo: demoRecentFiles},
}

// fileDay is a fixed demo time before Oct 24, 2026.
func fileDay(daysAgo int, hour int) time.Time {
	return time.Date(2026, 10, 24-daysAgo, 9+hour, 15, 0, 0, time.UTC)
}

func demoFileItems() []files.FileItem {
	return []files.FileItem{
		{Name: "Receipts", Dir: true, Modified: fileDay(3, 0)},
		{Name: "Q3 board report.pdf", Size: 2_480_000, Modified: fileDay(0, 1)},
		{Name: "revenue-by-region.xlsx", Size: 184_000, Modified: fileDay(1, 4)},
		{Name: "launch-notes.md", Size: 6_200, Modified: fileDay(0, 0)},
		{Name: "office-floorplan.png", Size: 1_120_000, Modified: fileDay(6, 2)},
		{Name: "billing.go", Size: 14_800, Modified: fileDay(2, 3)},
	}
}

func demoFileList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	selected, sort := intState(root, "selected", -1), ui.Local(root, "sort", func() files.FileSort { return files.SortName })
	list := demoFileItems()
	root.Children(func() {
		before := *sort
		if files.FileList(c, "Documents", list, selected, sort, files.FileListOptions{}).Changed() {
			if *sort != before {
				log(fmt.Sprintf("sort %d", *sort))
			} else {
				log("selected " + list[*selected].Name)
			}
		}
	})
}

func demoFileGrid(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	selected := intState(root, "selected", -1)
	list := demoFileItems()
	root.Children(func() {
		if files.FileGrid(c, "Documents", list, selected, files.FileGridOptions{}).Changed() {
			log("selected " + list[*selected].Name)
		}
	})
}

func demoFileIcon(c *ui.Context, log func(string)) {
	ui.Row(c).Width(480).Gap(16).Children(func() {
		for _, name := range []string{"report.pdf", "billing.go", "sales.xlsx", "photo.png", "notes.md", "backup.zip"} {
			c.Key(name)
			files.FileIcon(c, name, files.FileIconOptions{})
		}
	})
}

func demoDropzone(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	paths := ui.Local(root, "paths", func() []string { return []string{"/Users/ana/Desktop/Q3 board report.pdf"} })
	root.Children(func() {
		if files.Dropzone(c, "Upload to Documents", paths, files.DropzoneOptions{}).Changed() {
			log(fmt.Sprintf("%d paths", len(*paths)))
		}
	})
}

func demoUploadList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	uploads := ui.Local(root, "uploads", func() []files.UploadItem {
		return []files.UploadItem{
			{Name: "Q3 board report.pdf", Size: 2_480_000, Progress: 1},
			{Name: "office-floorplan.png", Size: 1_120_000, Progress: 0.64},
			{Name: "revenue-by-region.xlsx", Size: 184_000, Progress: 0.22},
		}
	})
	root.Children(func() {
		if files.UploadList(c, "Uploading to Documents", uploads, files.UploadListOptions{}).Changed() {
			log(fmt.Sprintf("%d uploads left", len(*uploads)))
		}
	})
}

const launchNotes = `# Launch notes

- Billing moves to monthly invoices on Nov 1.
- EU workspaces get VAT on every line.
- Support answers in English and German.

Owner: Ana Reyes
`

func demoFilePreview(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		files.FilePreview(c, "Preview", files.FileItem{Name: "launch-notes.md", Size: 6_200, Modified: fileDay(0, 0)}, launchNotes, nil,
			files.FilePreviewOptions{Height: 170})
	})
}

func demoFileCard(c *ui.Context, log func(string)) {
	root := ui.Row(c).Width(480).Gap(12)
	action := ui.Local(root, "action", func() files.FileAction { return files.FileActionNone })
	root.Children(func() {
		for _, f := range demoFileItems()[1:3] {
			c.Key(f.Name)
			ui.Column(c).Grow(1).Basis(0).MinWidth(0).Children(func() {
				if r, _ := files.FileCard(c, f, action, files.FileCardOptions{}); r.Changed() {
					log(fmt.Sprintf("action %d on %s", *action, f.Name))
				}
			})
		}
	})
}

func demoFolderCard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Gap(8)
	opened := boolState(root, "opened", false)
	root.Children(func() {
		if files.FolderCard(c, "Receipts", 48, opened, files.FolderCardOptions{}).Changed() {
			log("opened Receipts")
		}
	})
}

func demoPathBar(c *ui.Context, log func(string)) {
	parts := []string{"Acme", "Finance", "Reports", "2026", "Q3"}
	root := ui.Column(c).Width(480)
	at := intState(root, "at", 4)
	root.Children(func() {
		if files.PathBar(c, "Path", parts, at).Changed() {
			log("at " + parts[*at])
		}
	})
}

func demoNewFolder(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	name := textState(root, "name", "")
	root.Children(func() {
		r, create := files.NewFolder(c, name, files.NewFolderOptions{})
		if create.Changed() {
			log(fmt.Sprintf("created %d runes", utf8.RuneCountInString(*name)))
		} else if r.Changed() {
			log(fmt.Sprintf("typed %d runes", utf8.RuneCountInString(*name)))
		}
	})
}

func demoRenameRow(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	name := textState(root, "name", "draft.md")
	root.Children(func() {
		if files.RenameRow(c, name).Changed() {
			log(fmt.Sprintf("renamed %d runes", utf8.RuneCountInString(*name)))
		}
	})
}

func demoMoveDialog(c *ui.Context, log func(string)) {
	folders := []string{"Acme", "Acme/Finance", "Acme/Finance/Reports", "Acme/Finance/Receipts", "Acme/Marketing", "Acme/Marketing/Launch"}
	root := ui.Column(c).Width(420)
	target := intState(root, "target", -1)
	root.Children(func() {
		r, move := files.MoveDialog(c, "Move Q3 board report.pdf", folders, target, files.MoveDialogOptions{})
		if move.Clicked() {
			log("moved to " + folders[*target])
		} else if r.Changed() {
			log("target " + folders[*target])
		}
	})
}

func demoShareDialog(c *ui.Context, log func(string)) {
	people := []files.SharePerson{
		{Name: "Ana Reyes", Email: "ana@acme.dev", Role: "Owner"},
		{Name: "Kai Mori", Email: "kai@acme.dev", Role: "Editor"},
		{Name: "Lena Fischer", Email: "lena@acme.dev", Role: "Viewer"},
	}
	ui.Column(c).Width(440).Children(func() {
		if _, cp := files.ShareDialog(c, "Share Q3 board report", people, "https://drive.acme.dev/s/q3-board-7f2c", files.ShareDialogOptions{}); cp.Clicked() {
			log("copied link")
		}
	})
}

func demoPermissions(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	role := ui.Local(root, "role", func() files.PermissionRole { return files.Viewer })
	root.Children(func() {
		if files.Permissions(c, "Access for Lena Fischer", role, files.PermissionsOptions{}).Changed() {
			log(fmt.Sprintf("role %d", *role))
		}
	})
}

func demoTrash(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	items := ui.Local(root, "items", func() []files.FileItem {
		return []files.FileItem{
			{Name: "old-pricing.xlsx", Size: 96_000, Modified: fileDay(2, 0), Path: "Finance"},
			{Name: "logo-draft-3.png", Size: 640_000, Modified: fileDay(5, 1), Path: "Marketing/Launch"},
			{Name: "meeting-notes-sept.md", Size: 3_100, Modified: fileDay(12, 0), Path: "Acme"},
		}
	})
	restored := intState(root, "restored", -1)
	root.Children(func() {
		if files.Trash(c, "Trash", *items, restored, files.TrashOptions{}).Changed() {
			name := (*items)[*restored].Name
			*items = slices.Delete(slices.Clone(*items), *restored, *restored+1)
			*restored = -1
			log("restored " + name)
		}
	})
}

func demoProperties(c *ui.Context, log func(string)) {
	ui.Column(c).Width(440).Children(func() {
		files.Properties(c, files.FileItem{Name: "Q3 board report.pdf", Size: 2_480_000, Modified: fileDay(0, 1)},
			"/Acme/Finance/Reports/2026/Q3 board report.pdf", "PDF document", fileDay(9, 0), files.PropertiesOptions{})
	})
}

func demoVersions(c *ui.Context, log func(string)) {
	versions := []files.VersionItem{
		{ID: "v4", Author: "Ana Reyes", Modified: fileDay(0, 1), Size: 2_480_000, Current: true},
		{ID: "v3", Author: "Kai Mori", Modified: fileDay(1, 5), Size: 2_410_000},
		{ID: "v2", Author: "Ana Reyes", Modified: fileDay(3, 2), Size: 1_980_000},
		{ID: "v1", Author: "Lena Fischer", Modified: fileDay(9, 0), Size: 860_000},
	}
	root := ui.Column(c).Width(480)
	restored := intState(root, "restored", -1)
	root.Children(func() {
		if files.Versions(c, "Versions of Q3 board report.pdf", versions, restored, files.VersionsOptions{}).Changed() {
			log("restored " + versions[*restored].ID)
		}
	})
}

func demoRecentFiles(c *ui.Context, log func(string)) {
	list := []files.FileItem{
		{Name: "Q3 board report.pdf", Modified: fileDay(0, 1), Path: "Finance/Reports/2026"},
		{Name: "launch-notes.md", Modified: fileDay(0, 0), Path: "Marketing/Launch"},
		{Name: "revenue-by-region.xlsx", Modified: fileDay(1, 4), Path: "Finance"},
		{Name: "billing.go", Modified: fileDay(2, 3), Path: "Engineering/api"},
	}
	root := ui.Column(c).Width(480)
	opened := intState(root, "opened", -1)
	root.Children(func() {
		if files.RecentFiles(c, "Recent", list, opened, files.RecentFilesOptions{}).Changed() {
			log("opened " + list[*opened].Name)
		}
	})
}
