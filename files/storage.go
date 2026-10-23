package files

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// StoragePart is one type's storage use.
type StoragePart struct {
	Name  string
	Bytes int64
}

// Storage shows usage bars by file type.
func Storage(c *ui.Context, label string, parts []StoragePart, total int64) ui.Element {
	if label == "" || len(parts) == 0 || total <= 0 {
		panic("files: Storage needs parts and a positive total")
	}
	t := c.Theme()
	inks := []ui.Color{t.Accent, t.Success, t.Warning, t.Danger}
	used := int64(0)
	for _, part := range parts {
		if part.Name == "" || part.Bytes < 0 {
			panic("files: Storage parts need names and nonnegative sizes")
		}
		if part.Bytes > total-used {
			panic("files: Storage usage exceeds total")
		}
		used += part.Bytes
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		kit.Numeric(c, fmt.Sprintf("%s / %s", sizeLabel(used), sizeLabel(total))).Bold()
		for i, part := range parts {
			c.Key(part.Name)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Row(c).FillWidth().Children(func() {
					ui.Text(c, part.Name).FontSize(t.Rem(0.85)).Grow(1)
					sizeText(c, part.Bytes)
				})
				progressBar(c, float64(part.Bytes)/float64(total), inks[i%len(inks)])
			})
		}
	})
	return root
}
