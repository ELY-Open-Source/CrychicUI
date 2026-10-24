// Package shop holds storefront components.
package shop

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// Product is one sellable product.
type Product struct {
	Name, Brand string
	Price       float64
	// OldPrice is the price before a sale, 0 without one.
	OldPrice float64
	Stock    int
	Image    *ui.Bitmap
	// Rating is the mean of Reviews ratings, 0 with none.
	Rating  float64
	Reviews int
}

// CartItem is a product and its quantity.
type CartItem struct {
	Product  Product
	Quantity int
}

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

func validProduct(component string, p Product) {
	if p.Name == "" || p.Price < 0 || p.OldPrice != 0 && p.OldPrice < p.Price || !finite(p.Price) || !finite(p.OldPrice) || p.Stock < 0 ||
		!(p.Rating >= 0 && p.Rating <= 5) || p.Reviews < 0 || p.Rating > 0 && p.Reviews == 0 {
		panic(fmt.Sprintf("shop: %s product %q needs prices, stock and a rating in [0,5] with reviews", component, p.Name))
	}
}

// stars draws five stars, rating of them filled warm.
func stars(c *ui.Context, rating float64, side float32) ui.Element {
	if !(rating >= 0 && rating <= 5) {
		panic(fmt.Sprintf("shop: stars need a rating in [0,5], got %v", rating))
	}
	t := c.Theme()
	gap := side * 0.15
	box := ui.Box(c).Size(5*side+4*gap, side).Shrink(0)
	fill := box.AnimateTo("rating", float32(rating), ui.Smooth)
	return box.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
		row := func(ink ui.Color) {
			for i := range 5 {
				var path ui.Path
				starPath(&path, r.X+float32(i)*(side+gap)+side/2, r.Y+side/2, side/2)
				p.FillPath(&path, ink)
			}
		}
		row(t.Border)
		whole := float32(int(fill))
		w := whole*(side+gap) + (fill-whole)*side
		p.Clip(ui.Rect{X: r.X, Y: r.Y, W: w, H: r.H}, 0, func() { row(t.Warning) })
	}))
}

// starPath traces a five-point star around (x, y).
func starPath(p *ui.Path, x, y, radius float32) {
	for i := range 10 {
		rr := radius
		if i%2 == 1 {
			rr = radius * 0.45
		}
		a := float64(i)*math.Pi/5 - math.Pi/2
		px, py := x+rr*float32(math.Cos(a)), y+rr*float32(math.Sin(a))
		if i == 0 {
			p.MoveTo(px, py)
		} else {
			p.LineTo(px, py)
		}
	}
	p.Close()
}

// stepper is a minus, count, plus pill; callers play cues.
func stepper(c *ui.Context, qty, maxQty int, less, more string) (dec, inc ui.Element) {
	t := c.Theme()
	ui.Row(c).Shrink(0).AlignItems(ui.Center).Radius(999).Border(1, t.Border).Children(func() {
		side := func(name, glyph string, off bool) ui.Element {
			b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(7), t.Space(7)).Center().Radius(999)).Disabled(off)
			b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy))
			return b.Children(func() { ui.Icon(c, icon.Must(glyph)).Size(14, 14) })
		}
		dec = side(less, "minus", qty <= 1)
		kit.Numeric(c, fmt.Sprint(qty)).FontSize(t.Rem(0.85)).Bold().MinWidth(t.Space(5)).TextAlign(ui.Center).
			ContentTransition(ui.ContentTransition{Numeric: true})
		inc = side(more, "plus", qty >= maxQty)
	})
	return dec, inc
}

func money(c *ui.Context, value float64, currency string) ui.Element {
	if currency == "" || value < 0 || !finite(value) {
		panic("shop: money needs a currency and nonnegative value")
	}
	return typography.Money(c, value, currency)
}

func signedMoney(c *ui.Context, value float64, currency string) ui.Element {
	if currency == "" || !finite(value) {
		panic("shop: money needs a currency and finite value")
	}
	return typography.Money(c, value, currency)
}

func moneyText(c *ui.Context, value float64, currency string) string {
	if currency == "" || value < 0 || !finite(value) {
		panic("shop: money needs a currency and nonnegative value")
	}
	o := i18n.NumberOptions{MinFractionDigits: 2, MaxFractionDigits: 2}
	return currency + " " + i18n.FormatNumber(c.Region().FormatLocale(c.Locale()), value, o)
}

func addButton(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	return kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).
		Background(t.Accent).TextColor(t.AccentText))
}
