package feedback

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// StoragePart is a named share in the caller's color.
type StoragePart struct {
	Name  string
	Size  float64
	Color ui.Color
}

// StorageBarOptions tune a StorageBar.
type StorageBarOptions struct {
	// Format writes each part's size; "%.0f GB" when empty.
	Format string
}

// StorageBar is a usage bar of growing parts, with legend.
func StorageBar(c *ui.Context, label string, parts []StoragePart, total float64, o StorageBarOptions) ui.Element {
	sum := 0.0
	names := make([]string, len(parts))
	for i, p := range parts {
		if p.Name == "" || slices.Contains(names[:i], p.Name) || !(p.Size >= 0) || math.IsInf(p.Size, 0) {
			panic("feedback: StorageBar parts need distinct names and finite sizes")
		}
		names[i] = p.Name
		sum += p.Size
	}
	if label == "" || len(parts) == 0 || math.IsInf(total, 0) || !(sum <= total) {
		panic("feedback: StorageBar needs a label, parts and a finite total over their sum")
	}
	if o.Format == "" {
		o.Format = "%.0f GB"
	}
	t := c.Theme()
	rad := t.Space(1)
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	return root.Children(func() {
		bar := ui.Row(c).FillWidth().Height(t.Space(2.5)).Radius(rad).Background(t.Border).Clip().Gap(1).
			Role(ui.RoleMeter).Label(label).Range(0, total, sum)
		bar.Children(func() {
			for i, p := range parts {
				c.Key(p.Name)
				a := ui.Smooth
				a.Delay = time.Duration(i) * 80 * time.Millisecond
				seg := ui.Box(c).FillHeight().Background(p.Color)
				g := seg.Keyframes("grow", grow, ui.KeyframeOptions{Delay: a.Delay})
				seg.WidthPercent(seg.AnimateTo("w", float32(p.Size/total)*100, a) * g)
			}
		})
		ui.Row(c).Wrap().Gap(t.Space(4)).Children(func() {
			for _, p := range parts {
				c.Key(p.Name)
				ui.Row(c).Gap(t.Space(1.5)).Children(func() {
					ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Radius(t.Space(0.75)).Background(p.Color)
					ui.Text(c, p.Name).FontSize(t.Rem(0.85))
					kit.Numeric(c, fmt.Sprintf(o.Format, p.Size)).FontSize(t.Rem(0.85)).TextColor(t.TextMuted)
				})
			}
		})
	})
}
