// Package shop holds storefront components.
package shop

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// Product is one sellable product.
type Product struct {
	Name     string
	Price    float64
	OldPrice float64
	Stock    int
	Image    *ui.Bitmap
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
	if p.Name == "" || p.Price < 0 || p.OldPrice < p.Price || !finite(p.Price) || !finite(p.OldPrice) || p.Stock < 0 {
		panic(fmt.Sprintf("shop: %s has an invalid product", component))
	}
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
