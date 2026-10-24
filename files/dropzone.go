package files

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DropzoneOptions word the drop area.
type DropzoneOptions struct {
	Prompt, Detail string
	// One and Many word the dropped count, "%d" taking it.
	One, Many string
	Height    float32
}

// Dropzone takes dropped file paths; its border wakes; cue Drop.
func Dropzone(c *ui.Context, label string, paths *[]string, o DropzoneOptions) ui.Element {
	if label == "" || paths == nil {
		panic(fmt.Sprintf("files: Dropzone %q needs a label and paths", label))
	}
	def(&o.Prompt, "Drop files here")
	def(&o.Detail, "Files keep their full paths")
	def(&o.One, "%d file added")
	def(&o.Many, "%d files added")
	t := c.Theme()
	dimension("Dropzone", &o.Height, t.Space(36))
	root := ui.Column(c).FillWidth().Height(o.Height).Gap(t.Space(2)).AlignItems(ui.Center).Justify(ui.Center).
		Radius(t.Radius).Role(ui.RoleGroup).Label(label)
	over := root.AnimateTo("over", unit(root.FileDragOver()), ui.Snappy)
	root.Background(t.Surface.Mix(t.Accent.Alpha(0.08), over))
	line := t.Border.Mix(t.Accent, over)
	root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		p.StrokeDashed(ui.Rect{X: r.X + 0.5, Y: r.Y + 0.5, W: r.W - 1, H: r.H - 1}, line, t.Radius, 1+over)
	}))
	root.Children(func() {
		tile := ui.Box(c).Size(t.Space(11), t.Space(11)).Center().Radius(999).Background(t.Accent.Alpha(0.12)).TextColor(t.Accent)
		lift := float32(0)
		if over > 0 {
			lift = over * (1 - 2*abs(tile.Loop("bob", 900*time.Millisecond, ui.EaseInOut)-0.5))
		}
		tile.Translate(0, -t.Space(1)*lift).Scale(1+0.08*over, 1+0.08*over)
		tile.Children(func() { ui.Icon(c, icon.Must("cloud-upload")).Size(20, 20) })
		ui.Text(c, o.Prompt).Bold()
		kit.Muted(c, o.Detail).FontSize(t.Rem(0.8))
		if n := len(*paths); n > 0 {
			format := o.Many
			if n == 1 {
				format = o.One
			}
			c.Key("count")
			kit.Chip(c, fmt.Sprintf(format, n), t.Success)
		}
	})
	appendDrop(root, paths, root.DroppedFiles())
	return root
}

func abs(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
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
