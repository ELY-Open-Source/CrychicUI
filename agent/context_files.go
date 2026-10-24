package agent

import (
	"fmt"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

// ContextFile is one file attached to a prompt.
type ContextFile struct {
	Name string
	Size int64
}

// ContextFilesOptions tune the file chips.
type ContextFilesOptions struct {
	Remove, Empty string
}

// fileGlyph picks a Lucide icon by extension.
func fileGlyph(name string) string {
	switch strings.ToLower(path.Ext(name)) {
	case ".go", ".js", ".ts", ".py", ".rs", ".json", ".html", ".css":
		return "file-code"
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg":
		return "file-image"
	case ".csv", ".xlsx", ".tsv":
		return "file-spreadsheet"
	}
	return "file-text"
}

// ContextFiles shows attached files as chips with remove.
func ContextFiles(c *ui.Context, label string, files *[]ContextFile, o ContextFilesOptions) ui.Element {
	if label == "" || files == nil {
		panic(fmt.Sprintf("agent: ContextFiles %q needs a label and files", label))
	}
	seen := map[string]bool{}
	for _, f := range *files {
		if f.Name == "" || f.Size < 0 || seen[f.Name] {
			panic(fmt.Sprintf("agent: ContextFiles %q files need unique names and sizes", label))
		}
		seen[f.Name] = true
	}
	def(&o.Remove, "Remove")
	def(&o.Empty, "No files attached")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).FocusGroup(ui.Horizontal).Role(ui.RoleGroup).Label(label)
	delAt := -1
	root.Children(func() {
		if len(*files) == 0 {
			empty := ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).
				Border(1, t.Border)
			kit.Arrive(empty, ui.Pose{Fade: 1}, 0).Children(func() {
				ui.Icon(c, icon.Must("paperclip")).Size(16, 16).TextColor(t.TextMuted)
				kit.Muted(c, o.Empty).FontSize(t.Rem(0.85))
			})
		}
		for i, f := range *files {
			c.Key(f.Name)
			chip := ui.Row(c).MaxWidthPercent(100).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(1), t.Space(1), t.Space(1.5)).
				Radius(t.Radius).Background(t.Background).Border(1, t.Border).Transition(kit.Pop)
			kit.Arrive(chip, ui.Pose{Scale: 0.8, Fade: 1}, time.Duration(i)*40*time.Millisecond)
			var del ui.Element
			chip.Children(func() {
				tile(c, fileGlyph(f.Name), t.Accent, t.Space(7))
				ui.Column(c).Shrink(0).Children(func() {
					ui.Text(c, f.Name).FontSize(t.Rem(0.85)).Shrink(0).MaxWidth(t.Space(48)).SingleLine()
					typography.Bytes(c, f.Size).FontSize(t.Rem(0.7)).TextColor(t.TextMuted)
				})
				del = kit.Press(ui.ButtonBase(c).Label(o.Remove+" "+f.Name).Tooltip(o.Remove).Size(t.Space(6), t.Space(6)).Shrink(0).Center().Radius(999))
				del.Background(del.AnimateColor("bg", t.Danger.Alpha(0.15*unit(del.Hovered())), ui.Snappy)).
					TextColor(del.AnimateColor("ink", t.TextMuted.Mix(t.Danger, unit(del.Hovered())), ui.Snappy)).
					Children(func() { ui.Icon(c, icon.Must("x")).Size(16, 16) })
			})
			if del.Clicked() {
				delAt = i
			}
		}
	})
	if delAt >= 0 {
		*files = slices.Delete(slices.Clone(*files), delAt, delAt+1)
		root.SetChanged()
		sound.Play(sound.Delete)
		kit.Log("ContextFiles", "removed", "left", len(*files))
	}
	return root
}
