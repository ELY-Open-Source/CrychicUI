package files

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FilePreviewOptions size the preview and word its parts.
type FilePreviewOptions struct {
	Height float32
	// FileType words a type, "%s" taking the extension.
	Empty, FileType, Folder, Modified string
	FormatDate                        func(t time.Time) string
}

// FilePreview shows text, an image or details by file type.
func FilePreview(c *ui.Context, label string, item FileItem, content string, img *ui.Bitmap, o FilePreviewOptions) ui.Element {
	validFile("FilePreview", item)
	if label == "" {
		panic("files: FilePreview needs a label")
	}
	kind := itemGlyph(item)
	if kind == "image" && img == nil {
		panic(fmt.Sprintf("files: FilePreview %q needs image data for %q", label, item.Name))
	}
	def(&o.Empty, "Empty file")
	def(&o.FileType, "%s file")
	def(&o.Folder, "Folder")
	def(&o.Modified, "Modified")
	if o.FormatDate == nil {
		o.FormatDate = shortDate
	}
	t := c.Theme()
	dimension("FilePreview", &o.Height, t.Space(56))
	typeName := o.Folder
	if !item.Dir {
		typeName = fmt.Sprintf(o.FileType, strings.ToUpper(strings.TrimPrefix(filepath.Ext(item.Name), ".")))
	}
	root := kit.Card(c).FillWidth().Clip().Role(ui.RoleGroup).Label(label)
	kit.Arrive(root, ui.Pose{Fade: 1, Scale: 0.98}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2.5), t.Space(3)).
			BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			fileTile(c, item, t.Space(8))
			ui.Column(c).Grow(1).MinWidth(0).Children(func() {
				ui.Text(c, item.Name).Bold().FontSize(t.Rem(0.9)).FillWidth().SingleLine()
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					kit.Muted(c, typeName).FontSize(t.Rem(0.75)).Shrink(0)
					if !item.Dir {
						kit.Muted(c, "·").FontSize(t.Rem(0.75))
						sizeText(c, item.Size).FontSize(t.Rem(0.75))
					}
				})
			})
		})
		body := ui.Column(c).FillWidth().Height(o.Height).Background(t.Surface).Clip()
		body.Children(func() {
			switch {
			case kind == "image":
				ui.Image(c, img).Fill().Fit(ui.Contain)
			case kind == "file-text" || kind == "file-code":
				if content == "" {
					kit.Muted(c, o.Empty).Padding(t.Space(3))
					return
				}
				ui.Scroll(c).Fill().Padding(t.Space(2), 0).Children(func() { previewLines(c, content) })
			default:
				ui.Column(c).Fill().Gap(t.Space(2)).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
					fileTile(c, item, t.Space(16))
					ui.Text(c, item.Name).Bold()
					kit.Muted(c, typeName).FontSize(t.Rem(0.8))
					kit.Muted(c, o.Modified+" "+o.FormatDate(item.Modified)).FontSize(t.Rem(0.8))
				})
			}
		})
	})
	return root
}

// previewLines shows text lines numbered in a quiet gutter.
func previewLines(c *ui.Context, content string) {
	t := c.Theme()
	lines := strings.Split(strings.TrimSuffix(content, "\n"), "\n")
	digits := float32(len(strconv.Itoa(len(lines))))
	for i, line := range lines {
		c.Key(i)
		ui.Row(c).FillWidth().Gap(t.Space(3)).Padding(0, t.Space(3)).Children(func() {
			kit.Numeric(c, strconv.Itoa(i+1)).Font("monospace").FontSize(t.Rem(0.75)).TextColor(t.TextMuted.Alpha(0.7)).
				Width(t.Space(2 * digits)).Shrink(0).TextAlign(ui.End)
			ui.Text(c, strings.ReplaceAll(line, "\t", "    ")).Font("monospace").FontSize(t.Rem(0.8)).LineHeight(1.6).Grow(1).MinWidth(0).Wrap()
		})
	}
}
