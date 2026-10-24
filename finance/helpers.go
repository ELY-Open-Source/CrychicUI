// Package finance holds money and market components.
package finance

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/icon"
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

// tile is a rounded square holding a tinted icon.
func tile(c *ui.Context, glyph string, ink ui.Color) ui.Element {
	t := c.Theme()
	return ui.Box(c).Size(t.Space(10), t.Space(10)).Shrink(0).Center().Radius(t.Radius).Background(ink.Alpha(0.12)).
		TextColor(ink).Children(func() { ui.Icon(c, icon.Must(glyph)).Size(18, 18) })
}

// symbolTile shows a symbol's first letters in its palette tint.
func symbolTile(c *ui.Context, symbol string) ui.Element {
	t := c.Theme()
	ink := palette(t, len(symbol)+int(symbol[0]))
	r := []rune(symbol)
	return ui.Box(c).Size(t.Space(10), t.Space(10)).Shrink(0).Center().Radius(t.Radius).Background(ink.Alpha(0.14)).
		Children(func() {
			ui.Text(c, string(r[:min(2, len(r))])).FontSize(t.Rem(0.75)).Bold().TextColor(ink).Shrink(0)
		})
}

// changeInk is green for gains, red for losses, else muted.
func changeInk(t *ui.Theme, change float64) ui.Color {
	switch {
	case change > 0:
		return t.Success
	case change < 0:
		return t.Danger
	}
	return t.TextMuted
}

// changePill is a signed percent with a trend arrow.
func changePill(c *ui.Context, change float64) ui.Element {
	t := c.Theme()
	ink := changeInk(t, change)
	glyph := "minus"
	switch {
	case change > 0:
		glyph = "trending-up"
	case change < 0:
		glyph = "trending-down"
	}
	pill := ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(2)).Radius(999)
	pill.Background(pill.AnimateColor("bg", ink.Alpha(0.12), ui.Snappy)).TextColor(pill.AnimateColor("ink", ink, ui.Snappy))
	return pill.Children(func() {
		ui.Icon(c, icon.Must(glyph)).Size(12, 12).ContentTransition(ui.ContentTransition{})
		kit.Numeric(c, fmt.Sprintf("%+.2f%%", change)).FontSize(t.Rem(0.75)).Bold().Shrink(0).
			ContentTransition(ui.ContentTransition{Numeric: true})
	})
}

// palette is the i-th category color from the theme.
func palette(t *ui.Theme, i int) ui.Color {
	inks := [...]ui.Color{t.Accent, t.Success, t.Warning, t.Danger, t.Accent.Mix(t.Success, 0.5), t.Warning.Mix(t.Danger, 0.5)}
	return inks[(i%len(inks)+len(inks))%len(inks)]
}

// inversePair is the inverse fill and text, each defaulted alone.
func inversePair(t *ui.Theme) (ui.Color, ui.Color) {
	fill, text := t.Inverse, t.InverseText
	if fill == ui.Transparent {
		fill = t.Text
	}
	if text == ui.Transparent {
		text = t.Background
	}
	return fill, text
}
