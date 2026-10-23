package gallery

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/files"
	"github.com/egoist/mygo/ui"
)

var filesEntries = []Entry{
	{ID: 941, Name: "FileList", Category: "files", Summary: "Rows with icon, name, size, date; sort; cue Select.", Demo: demoFileList},
	{ID: 942, Name: "FileGrid", Category: "files", Summary: "Icon tiles with names; cue Select.", Demo: demoFileGrid},
	{ID: 943, Name: "FileIcon", Category: "files", Summary: "Icon chosen by extension; swaps.", Demo: demoFileIcon},
	{ID: 944, Name: "Dropzone", Category: "files", Summary: "Area for dropped files; border animates; cue Drop.", Demo: demoDropzone},
	{ID: 945, Name: "UploadList", Category: "files", Summary: "Files with bars and cancel; cue Close.", Demo: demoUploadList},
	{ID: 946, Name: "FilePreview", Category: "files", Summary: "Text, image or info preview by type.", Demo: demoFilePreview},
	{ID: 947, Name: "FileCard", Category: "files", Summary: "Icon, name, size and menu; lifts.", Demo: demoFileCard},
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

func demoFileItems() []files.FileItem {
	now := time.Now()
	return []files.FileItem{
		{Name: "notes.md", Size: 2400, Modified: now.Add(-time.Hour)},
		{Name: "main.go", Size: 8400, Modified: now.Add(-2 * time.Hour)},
		{Name: "photo.png", Size: 42000, Modified: now.Add(-3 * time.Hour)},
	}
}

func demoFileList(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected, sort := intState(root, "selected", -1), ui.Local(root, "sort", func() files.FileSort { return files.SortName })
	root.Children(func() {
		if files.FileList(c, "Files", demoFileItems(), selected, sort, files.FileListOptions{}).Changed() {
			log(fmt.Sprintf("selected %d sort %d", *selected, *sort))
		}
	})
}

func demoFileGrid(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := intState(root, "selected", -1)
	root.Children(func() {
		if files.FileGrid(c, "Files", demoFileItems(), selected).Changed() {
			log("selected " + fmt.Sprint(*selected))
		}
	})
}

func demoFileIcon(c *ui.Context, log func(string)) {
	files.FileIcon(c, "main.go")
}

func demoDropzone(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	paths := listState(root, "paths")
	root.Children(func() {
		if files.Dropzone(c, "Upload", paths, files.DropzoneOptions{}).Changed() {
			log(fmt.Sprint(len(*paths)) + " paths")
		}
	})
}

func demoUploadList(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	uploads := ui.Local(root, "uploads", func() []files.UploadItem {
		return []files.UploadItem{{Name: "report.pdf", Progress: 0.72}, {Name: "photo.png", Progress: 0.38}}
	})
	root.Children(func() {
		if files.UploadList(c, "Uploads", uploads, files.UploadListOptions{}).Changed() {
			log(fmt.Sprint(len(*uploads)) + " uploads left")
		}
	})
}

func demoFilePreview(c *ui.Context, log func(string)) {
	item := files.FileItem{Name: "notes.md", Size: 48, Modified: time.Now()}
	files.FilePreview(c, "Preview", item, "# Notes\n\nA text preview.", nil, files.FilePreviewOptions{})
}

func demoFileCard(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	action := ui.Local(root, "action", func() files.FileAction { return files.FileActionNone })
	root.Children(func() {
		view, _ := files.FileCard(c, files.FileItem{Name: "notes.md", Size: 2400, Modified: time.Now()}, action, files.FileCardOptions{})
		if view.Changed() {
			log("action " + fmt.Sprint(*action))
			*action = files.FileActionNone
		}
	})
}

func demoFolderCard(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	opened := boolState(root, "opened", false)
	root.Children(func() {
		if files.FolderCard(c, "Documents", 12, opened, files.FolderCardOptions{}).Changed() {
			log("opened Documents")
		}
	})
}

func demoPathBar(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	at := intState(root, "at", 2)
	parts := []string{"home", "projects", "crychic"}
	root.Children(func() {
		if files.PathBar(c, parts, at).Changed() {
			log("at " + parts[*at])
		}
	})
}

func demoNewFolder(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	name := textState(root, "name", "")
	root.Children(func() {
		_, create := files.NewFolder(c, name, files.NewFolderOptions{})
		if create.Changed() {
			log("created " + fmt.Sprint(utf8.RuneCountInString(*name)) + " runes")
		}
	})
}

func demoRenameRow(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	name := textState(root, "name", "draft.md")
	root.Children(func() {
		if files.RenameRow(c, name).Changed() {
			log("renamed " + fmt.Sprint(utf8.RuneCountInString(*name)) + " runes")
		}
	})
}

func demoMoveDialog(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	target := intState(root, "target", -1)
	root.Children(func() {
		_, move := files.MoveDialog(c, "Move file", []string{"home", "work/docs", "work/images"}, target, files.MoveDialogOptions{})
		if move.Clicked() {
			log("moved to " + fmt.Sprint(*target))
		}
	})
}

func demoShareDialog(c *ui.Context, log func(string)) {
	people := []files.SharePerson{{Name: "Ana", Role: "Owner"}, {Name: "Kai", Role: "Editor"}}
	_, copy := files.ShareDialog(c, "Share", people, "https://example.test/f/crychic", files.ShareDialogOptions{})
	if copy.Clicked() {
		log("copied link")
	}
}

func demoPermissions(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	role := ui.Local(root, "role", func() files.PermissionRole { return files.Viewer })
	root.Children(func() {
		if files.Permissions(c, "Access", role, files.PermissionsOptions{}).Changed() {
			log("role " + fmt.Sprint(*role))
		}
	})
}

func demoTrash(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	items := listState(root, "items", "old.md", "unused.go")
	restored := intState(root, "restored", -1)
	root.Children(func() {
		if files.Trash(c, "Trash", *items, restored, files.TrashOptions{}).Changed() {
			log("restored " + fmt.Sprint(*restored))
			*items = append((*items)[:*restored], (*items)[*restored+1:]...)
			*restored = -1
		}
	})
}

func demoProperties(c *ui.Context, log func(string)) {
	item := files.FileItem{Name: "notes.md", Size: 2400, Modified: time.Now()}
	files.Properties(c, item, "/home/docs/notes.md", "Markdown", time.Now().Add(-24*time.Hour), files.PropertiesOptions{})
}

func demoVersions(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	restored := intState(root, "restored", -1)
	versions := []files.VersionItem{
		{ID: "v2", Size: 2400, Modified: time.Now(), Current: true},
		{ID: "v1", Size: 1800, Modified: time.Now().Add(-time.Hour)},
	}
	root.Children(func() {
		if files.Versions(c, "Versions", versions, restored, files.VersionsOptions{}).Changed() {
			log("restored " + versions[*restored].ID)
		}
	})
}

func demoRecentFiles(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	opened := intState(root, "opened", -1)
	list := demoFileItems()
	root.Children(func() {
		if files.RecentFiles(c, "Recent", list, opened).Changed() {
			log("opened " + list[*opened].Name)
		}
	})
}
