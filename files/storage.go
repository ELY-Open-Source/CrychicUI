package files

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// StoragePart is one type's storage use.
type StoragePart struct {
	Name  string
	Bytes int64
}

// StorageOptions word the totals and the free part.
type StorageOptions struct {
	Of, Free string
}

// Storage shows usage by type as a growing stacked bar.
func Storage(c *ui.Context, label string, parts []StoragePart, total int64, o StorageOptions) ui.Element {
	if label == "" || len(parts) == 0 || total <= 0 {
		panic(fmt.Sprintf("files: Storage %q needs parts and a positive total", label))
	}
	used := int64(0)
	seen := map[string]bool{}
	for _, part := range parts {
		if part.Name == "" || seen[part.Name] || part.Bytes < 0 || part.Bytes > total-used {
			panic(fmt.Sprintf("files: Storage %q needs unique named parts within the total", label))
		}
		seen[part.Name] = true
		used += part.Bytes
	}
	def(&o.Of, "of")
	def(&o.Free, "Free")
	t := c.Theme()
	inks := []ui.Color{t.Accent, t.Success, t.Warning, t.Danger, t.TextMuted}
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		kit.CardHead(c, "hard-drive", label, func() {
			kit.Numeric(c, strconv.Itoa(int(100*float64(used)/float64(total)+0.5))+"%").FontSize(t.Rem(0.75)).Bold().TextColor(t.Accent)
		})
		ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.End).Children(func() {
				sizeText(c, used).FontSize(t.Rem(1.4)).Bold().TextColor(t.Text)
				kit.Muted(c, o.Of).FontSize(t.Rem(0.85))
				sizeText(c, total).FontSize(t.Rem(0.85))
			})
			bar := ui.Box(c).FillWidth().Height(t.Space(3)).Radius(999).Background(t.Surface)
			g := bar.Keyframes("crychic.grow", []ui.Keyframe{{Value: 0}, {At: 700 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			shares := make([]float32, len(parts))
			for i, part := range parts {
				shares[i] = bar.AnimateTo(fmt.Sprint("part.", i), float32(float64(part.Bytes)/float64(total)), ui.Smooth) * g
			}
			bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
				x := r.X
				for i, s := range shares {
					if w := r.W*s - 2; w > 0 {
						p.Fill(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, inks[i%len(inks)], 0)
					}
					x += r.W * s
				}
			}))
			ui.Column(c).FillWidth().Gap(t.Space(1.5)).Children(func() {
				row := func(i int, key, name string, bytes int64, ink ui.Color) {
					c.Key(key)
					item := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
					kit.Stagger(c, root, item, i).Children(func() {
						ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Shrink(0).Radius(999).Background(ink)
						ui.Text(c, name).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
						sizeText(c, bytes)
						kit.Numeric(c, strconv.Itoa(int(100*float64(bytes)/float64(total)+0.5))+"%").FontSize(t.Rem(0.75)).
							TextColor(t.TextMuted).Width(t.Space(10)).Shrink(0).TextAlign(ui.End)
					})
				}
				for i, part := range parts {
					row(i, fmt.Sprintf("part %q", part.Name), part.Name, part.Bytes, inks[i%len(inks)])
				}
				row(len(parts), "free", o.Free, total-used, t.Border)
			})
		})
	})
	return root
}
