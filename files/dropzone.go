package files

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DropzoneOptions are the drop area labels.
type DropzoneOptions struct {
	Prompt, Detail string
	Height         float32
}

// Dropzone accepts dropped file paths; cue Drop.
func Dropzone(c *ui.Context, label string, paths *[]string, o DropzoneOptions) ui.Element {
	if label == "" || paths == nil {
		panic("files: Dropzone needs a label and paths")
	}
	def(&o.Prompt, "Drop files here")
	def(&o.Detail, "Files keep their full paths")
	t := c.Theme()
	dimension("Dropzone", &o.Height, t.Space(32))
	root := ui.Column(c).FillWidth().Height(o.Height).Gap(t.Space(2)).AlignItems(ui.Center).Justify(ui.Center).
		Radius(t.Radius).Background(t.Surface).Role(ui.RoleGroup).Label(label)
	root.Border(1, root.AnimateColor("border", t.Border.Mix(t.Accent, unit(root.FileDragOver())), ui.Snappy))
	root.Children(func() {
		ui.Icon(c, icon.Must("upload")).Size(t.Space(7), t.Space(7)).TextColor(t.Accent)
		ui.Text(c, o.Prompt).Bold()
		kit.Muted(c, o.Detail).FontSize(t.Rem(0.8))
	})
	appendDrop(root, paths, root.DroppedFiles())
	return root
}

func appendDrop(root ui.Element, paths *[]string, dropped []string) {
	if len(dropped) == 0 {
		return
	}
	*paths = append(*paths, dropped...)
	root.SetChanged()
	sound.Play(sound.Drop)
	kit.Log("Dropzone", "dropped", "count", len(dropped))
}
