// Package finance holds money and market components.
package finance

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

func def(s *string, d string) {
	if *s == "" {
		*s = d
	}
}

func unit(v bool) float32 {
	if v {
		return 1
	}
	return 0
}

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func moneyText(c *ui.Context, value float64, currency string) string {
	if currency == "" || !finite(value) {
		panic("finance: invalid money")
	}
	o := i18n.NumberOptions{MinFractionDigits: 2, MaxFractionDigits: 2}
	return currency + " " + i18n.FormatNumber(c.Region().FormatLocale(c.Locale()), value, o)
}

func money(c *ui.Context, value float64, currency string) ui.Element {
	if currency == "" || !finite(value) {
		panic(fmt.Sprintf("finance: invalid money %v", value))
	}
	return typography.Money(c, value, currency)
}

func progress(c *ui.Context, value float64, color ui.Color) ui.Element {
	if value < 0 || value > 1 || !finite(value) {
		panic("finance: progress must be in [0,1]")
	}
	t := c.Theme()
	bar := ui.Box(c).FillWidth().Height(t.Space(1)).Radius(999).Clip().Background(t.SurfaceHover)
	w := bar.AnimateTo("progress", float32(value), ui.Smooth)
	bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * w, H: r.H}, color, 999)
	}))
	return bar
}
