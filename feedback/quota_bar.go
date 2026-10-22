package feedback

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// QuotaBarOptions tune a QuotaBar.
type QuotaBarOptions struct {
	// Format writes used and limit; "%.0f of %.0f" when empty.
	Format string
}

// QuotaBar is used of a limit; Warning past 80%.
func QuotaBar(c *ui.Context, label string, used, limit float64, o QuotaBarOptions) ui.Element {
	if label == "" || !(limit > 0) || math.IsInf(limit, 0) || !(used >= 0) || math.IsInf(used, 0) {
		panic("feedback: QuotaBar needs a label, used >= 0 and a finite limit > 0")
	}
	if o.Format == "" {
		o.Format = "%.0f of %.0f"
	}
	t := c.Theme()
	k := min(used/limit, 1)
	tint := t.Accent
	switch {
	case used >= limit:
		tint = t.Danger
	case k > 0.8:
		tint = t.Warning
	}
	rad := t.Space(0.75)
	root := ui.Column(c).FillWidth().Gap(t.Space(1.5))
	return root.Children(func() {
		ui.Row(c).Children(func() {
			ui.Text(c, label).Grow(1)
			kit.Numeric(c, fmt.Sprintf(o.Format, used, limit)).TextColor(t.TextMuted).FontSize(t.Rem(0.85))
		})
		bar := ui.Box(c).FillWidth().Height(t.Space(1.5)).Radius(rad).Background(t.Border).Clip().
			Role(ui.RoleMeter).Label(label).Range(0, limit, used)
		bar.Children(func() {
			fill := ui.Box(c).Absolute().Start(0).Top(0).Bottom(0).Radius(rad)
			fill.Background(fill.AnimateColor("tint", tint, quick)).WidthPercent(fill.AnimateTo("w", float32(k)*100, ui.Smooth))
		})
	})
}
