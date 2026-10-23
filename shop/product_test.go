package shop

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestProductCard(t *testing.T) {
	added := false
	tt := ui.NewTester(padded(func(c *ui.Context) { ProductCard(c, testProduct(), "USD", &added, ProductCardOptions{}) }), 500, 450)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Add to cart")
	if !added {
		t.Error("product was not added")
	}
	cues(t, before, sound.Pop)
}

func TestProductGallery(t *testing.T) {
	p := testProduct()
	selected := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ProductGallery(c, []*ui.Bitmap{p.Image, p.Image}, &selected, ProductGalleryOptions{})
	}), 500, 450)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Image 2")
	if selected != 1 {
		t.Errorf("selected %d", selected)
	}
	cues(t, before, sound.Select)
}

func TestProductPrice(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { ProductPrice(c, 80, 100, "USD", ProductPriceOptions{}) }), 400, 180)
	settles(t, tt)
	if !tt.HasText("20% off") || !tt.HasText("USD 80.00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestSaleBadge(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { SaleBadge(c, 20, SaleBadgeOptions{}) }), 200, 120)
	settles(t, tt)
	if !tt.HasText("-20%") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestAddToCart(t *testing.T) {
	added := false
	tt := ui.NewTester(padded(func(c *ui.Context) { AddToCart(c, "Add", &added, AddToCartOptions{}) }), 400, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Add")
	if !added {
		t.Error("item was not added")
	}
	cues(t, before, sound.Pop)
}

func TestAddToCartReplay(t *testing.T) {
	added := false
	tt := ui.NewTester(padded(func(c *ui.Context) { AddToCart(c, "Add", &added, AddToCartOptions{}) }), 400, 180)
	settles(t, tt)
	r, ok := tt.Find("Add")
	if !ok {
		t.Fatal("add button missing")
	}
	x, y := int(r.X+r.W-4), int(r.Y+4)
	before := tt.Image().RGBAAt(x, y)
	click(t, tt, "Add")
	tt.Advance(300 * time.Millisecond)
	if got := tt.Image().RGBAAt(x, y); got != before {
		t.Errorf("flight remains after completion: got %v want %v", got, before)
	}
	click(t, tt, "Add")
	if !tt.FramePending() {
		t.Error("second addition has no flight")
	}
}

func TestCartIcon(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { CartIcon(c, 3, CartIconOptions{}) }), 200, 120)
	settles(t, tt)
	if !tt.HasText("3") {
		t.Errorf("texts %q", tt.Texts())
	}
}
