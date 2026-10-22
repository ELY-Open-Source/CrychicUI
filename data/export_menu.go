package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Export is a format name and the text it copies.
type Export struct{ Name, Data string }

// ExportMenu opens formats; choosing one copies its data.
func ExportMenu(c *ui.Context, label string, exports []Export, chosen *int, open *bool) ui.Element {
	if label == "" || len(exports) == 0 || chosen == nil || open == nil || *chosen < -1 || *chosen >= len(exports) {
		panic(fmt.Sprintf("data: ExportMenu %q needs exports, a choice or -1, and open", label))
	}
	var actions []Action
	var names []string
	for _, e := range exports {
		actions = append(actions, Action{Label: e.Name, Icon: "copy"})
		names = append(names, e.Name)
	}
	unique("ExportMenu", label, names)
	t := c.Theme()
	root := ui.Box(c).Shrink(0)
	root.Children(func() {
		b := kit.Press(ui.ButtonBase(c).Label(label).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(3)).
			Radius(t.Radius).Border(1, t.Border)).Expanded(*open)
		b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(*open || b.Hovered())), ui.Snappy)).Children(func() {
			ui.Icon(c, icon.Must("download")).Size(16, 16)
			ui.Text(c, label)
		})
		menu(c, root, b, "ExportMenu", label, actions, open, func(i int) {
			*chosen = i
			c.WriteClipboard(exports[i].Data)
			sound.Play(sound.Pop)
			kit.Log("ExportMenu", "copied", "index", i, "bytes", len(exports[i].Data))
		})
	})
	return root
}
