package files

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FilePreviewOptions are the preview labels and size.
type FilePreviewOptions struct {
	Height                  float32
	Empty, Type, Size, Date string
}

// FilePreview shows text, image, or metadata by type.
func FilePreview(c *ui.Context, label string, item FileItem, content string, img *ui.Bitmap, o FilePreviewOptions) ui.Element {
	validFile("FilePreview", item)
	if label == "" {
		panic("files: FilePreview needs a label")
	}
	def(&o.Empty, "Empty file")
	def(&o.Type, "Type")
	def(&o.Size, "Size")
	def(&o.Date, "Modified")
	t := c.Theme()
	dimension("FilePreview", &o.Height, t.Space(48))
	root := ui.Column(c).FillWidth().Height(o.Height).Radius(t.Radius).Border(1, t.Border).Clip().
		Background(t.Surface).Role(ui.RoleGroup).Label(label)
	kit.Arrive(root, ui.Pose{Fade: 1, Scale: 0.98}, 0)
	root.Children(func() {
		kind := glyph(item.Name)
		if item.Dir {
			kind = "folder"
		}
		switch kind {
		case "image":
			if img == nil {
				panic("files: FilePreview image needs image data")
			}
			ui.Image(c, img).Fill().Fit(ui.Contain)
		case "file-text", "file-code":
			if content == "" {
				kit.Muted(c, o.Empty).Padding(t.Space(3))
				return
			}
			ui.Scroll(c).Fill().Children(func() {
				ui.Text(c, content).Font("monospace").FontSize(t.Rem(0.85)).Padding(t.Space(3)).Wrap()
			})
		default:
			ui.Column(c).Fill().Gap(t.Space(2)).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
				ui.Icon(c, icon.Must(itemGlyph(item))).Size(16, 16)
				ui.Text(c, item.Name).Bold()
				kit.Muted(c, o.Type+": "+kind).FontSize(t.Rem(0.8))
				sizeText(c, item.Size)
				kit.Muted(c, o.Date+": "+date(item.Modified)).FontSize(t.Rem(0.8))
			})
		}
	})
	return root
}
