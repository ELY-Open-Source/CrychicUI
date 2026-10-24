package shop

import (
	"bytes"
	"fmt"
	"image"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func rectOf(t *testing.T, tt *ui.Tester, name string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(name)
	if !ok {
		t.Fatalf("no %q in %q", name, tt.Texts())
	}
	return r
}

// swapUnderPress presses at r, reorders, then releases there.
func swapUnderPress(tt *ui.Tester, r ui.Rect, swap func()) {
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Frame()
	swap()
	tt.Frame()
	tt.Release(x, y)
	tt.Frame()
}

func bitmap() *ui.Bitmap { return ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 2, 2))) }

// solid is an opaque dark bitmap, visible wherever drawn.
func solid() *ui.Bitmap {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	for i := range img.Pix {
		img.Pix[i] = 40 + 215*uint8(i%4/3)
	}
	return ui.NewBitmap(img)
}

func TestPicksKeepIdentityAcrossReorder(t *testing.T) {
	cats := []Category{{Name: "Lighting", Image: bitmap()}, {Name: "Desks", Image: bitmap()}}
	sel := -1
	ct := ui.NewTester(padded(func(c *ui.Context) { Categories(c, cats, &sel, CategoriesOptions{}) }), 500, 350)
	settles(t, ct)
	swapUnderPress(ct, rectOf(t, ct, "Desks"), func() { cats[0], cats[1] = cats[1], cats[0] })
	if sel != -1 {
		t.Errorf("categories chose the tile that slid under the pointer: %d", sel)
	}

	brands, chosen, rating := []string{"Lumen & Co", "Fold"}, []bool{false, false}, 0
	low, high := 0.0, 1000.0
	ft := ui.NewTester(padded(func(c *ui.Context) {
		ShopFilters(c, &low, &high, brands, &chosen, &rating, ShopFiltersOptions{})
	}), 500, 600)
	settles(t, ft)
	swapUnderPress(ft, rectOf(t, ft, "Fold"), func() { brands[0], brands[1] = brands[1], brands[0] })
	if chosen[0] || chosen[1] {
		t.Errorf("filters chose the brand that slid under the pointer: %v", chosen)
	}

	choices, sort := []string{"Featured", "Newest", "Top rated"}, "Featured"
	st := ui.NewTester(padded(func(c *ui.Context) { SortBar(c, 9, &sort, choices, SortBarOptions{}) }), 500, 400)
	settles(t, st)
	click(t, st, "Sort by")
	settles(t, st)
	swapUnderPress(st, rectOf(t, st, "Top rated"), func() { choices[1], choices[2] = choices[2], choices[1] })
	if sort != "Featured" {
		t.Errorf("sort chose the choice that slid under the pointer: %q", sort)
	}

	images, selected := []*ui.Bitmap{bitmap(), bitmap()}, 0
	gt := ui.NewTester(padded(func(c *ui.Context) { ProductGallery(c, images, &selected, ProductGalleryOptions{}) }), 500, 500)
	settles(t, gt)
	swapUnderPress(gt, rectOf(t, gt, "Image 2"), func() { images[0], images[1] = images[1], images[0] })
	if selected != 0 {
		t.Errorf("gallery chose the thumbnail that slid under the pointer: %d", selected)
	}
}

func TestActionsKeepFocusAsOptionalRowsAppear(t *testing.T) {
	added, helpful, saved := false, false, false
	cases := []struct {
		name, action string
		view         func(c *ui.Context, on bool)
	}{
		{"ProductCard", "Add to cart", func(c *ui.Context, on bool) {
			p := testProduct()
			if on {
				p.Brand = "Lumen & Co"
			}
			ProductCard(c, p, "USD", &added, ProductCardOptions{})
		}},
		{"ReviewCard", "Helpful (2)", func(c *ui.Context, on bool) {
			r := testReviews()[0]
			if on {
				r.Title = "Great lamp"
			}
			ReviewCard(c, r, &helpful, ReviewCardOptions{})
		}},
		{"Wishlist", "Save", func(c *ui.Context, on bool) {
			p := testProduct()
			if !on {
				p.Image = nil
			}
			Wishlist(c, p, "USD", &saved, WishlistOptions{})
		}},
	}
	for _, tc := range cases {
		on := false
		tt := ui.NewTester(padded(func(c *ui.Context) { tc.view(c, on) }), 300, 600)
		settles(t, tt)
		tt.Key(0, ui.KeyTab)
		tt.Frame()
		if !tt.Focused(tc.action) {
			t.Fatalf("%s: Tab did not reach %q", tc.name, tc.action)
		}
		on = true
		tt.Frame()
		if !tt.Focused(tc.action) {
			t.Errorf("%s: %q lost focus when the optional row appeared", tc.name, tc.action)
		}
	}
}

func TestSortBarChangesOnlyOnNewSort(t *testing.T) {
	sort, changed := "Featured", 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if SortBar(c, 9, &sort, []string{"Featured", "Newest"}, SortBarOptions{}).Changed() {
			changed++
		}
	}), 500, 400)
	settles(t, tt)
	click(t, tt, "Sort by")
	settles(t, tt)
	click(t, tt, "Featured")
	settles(t, tt)
	if changed != 0 {
		t.Errorf("opening or keeping the sort reported %d changes", changed)
	}
	click(t, tt, "Sort by")
	settles(t, tt)
	click(t, tt, "Newest")
	if sort != "Newest" || changed != 1 {
		t.Errorf("sort %q changed %d", sort, changed)
	}
}

func TestProductGalleryScrollsToEveryThumbnail(t *testing.T) {
	images := make([]*ui.Bitmap, 8)
	for i := range images {
		images[i] = bitmap()
	}
	selected := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ui.Column(c).Width(340).Children(func() { ProductGallery(c, images, &selected, ProductGalleryOptions{}) })
	}), 600, 600)
	settles(t, tt)
	first := rectOf(t, tt, "Image 1")
	tt.Scroll(first.X+first.W/2, first.Y+first.H/2, 800, 0)
	settles(t, tt)
	if last := rectOf(t, tt, "Image 8"); last.W == 0 || last.X+last.W > 351 {
		t.Fatalf("last thumbnail out of reach: %v", last)
	}
	click(t, tt, "Image 8")
	if selected != 7 {
		t.Errorf("selected %d", selected)
	}
}

func TestPricesWithoutSaleOrTinySaving(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { ProductPrice(c, 24, 0, "USD", ProductPriceOptions{}) }), 400, 180)
	settles(t, tt)
	if !tt.HasText("USD 24.00") || tt.HasText("off") || tt.HasText("save") {
		t.Errorf("no-sale price %q", tt.Texts())
	}
	tt = ui.NewTester(padded(func(c *ui.Context) { ProductPrice(c, 99.9, 100, "USD", ProductPriceOptions{}) }), 400, 180)
	settles(t, tt)
	if tt.HasText("% off") || !tt.HasText("You save USD 0.10") {
		t.Errorf("a 0.1%% saving read %q", tt.Texts())
	}
	if got := percentOff(80, 100); got != 20 {
		t.Errorf("percentOff(80, 100) = %d", got)
	}
	if got := percentOff(8e307, 1e308); got != 20 {
		t.Errorf("percentOff(8e307, 1e308) = %d", got)
	}
}

func TestCartDrawerCountsPastInt(t *testing.T) {
	big, small := testProduct(), testProduct()
	big.Stock, small.Name = math.MaxInt, "Mug"
	items := []CartItem{{Product: big, Quantity: math.MaxInt}, {Product: small, Quantity: 1}}
	checkout := false
	tt := ui.NewTester(padded(func(c *ui.Context) { CartDrawer(c, items, "USD", &checkout, CartDrawerOptions{}) }), 500, 600)
	settles(t, tt)
	if want := fmt.Sprint(uint64(math.MaxInt) + 1); !tt.HasText(want) || strings.Contains(strings.Join(tt.Texts(), " "), "-9") {
		t.Errorf("count %q, want %s", tt.Texts(), want)
	}
}

func TestCheckoutStepsAnswerHover(t *testing.T) {
	step := Cart
	tt := ui.NewTester(padded(func(c *ui.Context) { CheckoutSteps(c, &step, CheckoutStepsOptions{}) }), 500, 200)
	settles(t, tt)
	before := bytes.Clone(tt.Image().Pix)
	r := rectOf(t, tt, "Payment")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("hovering a step started no animation")
	}
	settles(t, tt)
	if bytes.Equal(before, tt.Image().Pix) {
		t.Error("hovering a step changed nothing")
	}
}

func TestCategoryTilesArriveStaggered(t *testing.T) {
	var cats []Category
	for _, n := range []string{"Lighting", "Tableware", "Ceiling", "Desks", "Textiles", "Decor"} {
		cats = append(cats, Category{Name: n, Image: bitmap()})
	}
	sel := -1
	view := padded(func(c *ui.Context) { Categories(c, cats, &sel, CategoriesOptions{}) })
	settled := ui.NewTester(view, 500, 400)
	settles(t, settled)
	tt := ui.NewTester(view, 500, 400)
	tt.Frame()
	tt.Advance(210 * time.Millisecond)
	now, end := tt.Image(), settled.Image()
	same := func(name string) bool {
		r := rectOf(t, settled, name)
		for y := int(r.Y); y < int(r.Y+r.H); y++ {
			for x := int(r.X); x < int(r.X+r.W); x++ {
				if now.RGBAAt(x, y) != end.RGBAAt(x, y) {
					return false
				}
			}
		}
		return true
	}
	if !same("Lighting") || same("Decor") {
		t.Errorf("at 210ms first tile settled %v, last tile settled %v; want true, false", same("Lighting"), same("Decor"))
	}
}

func TestGalleryThumbsWaitForTheirArrival(t *testing.T) {
	images, selected := []*ui.Bitmap{solid(), solid(), solid()}, 0
	view := padded(func(c *ui.Context) { ProductGallery(c, images, &selected, ProductGalleryOptions{}) })
	settled := ui.NewTester(view, 500, 500)
	settles(t, settled)
	r := rectOf(t, settled, "Image 3")
	x, y := int(r.X+r.W/2), int(r.Y+r.H/2)
	blank := ui.NewTester(padded(func(c *ui.Context) {}), 500, 500).Image().RGBAAt(x, y)
	tt := ui.NewTester(view, 500, 500)
	tt.Frame()
	tt.Advance(20 * time.Millisecond)
	if got := tt.Image().RGBAAt(x, y); got != blank {
		t.Errorf("third thumbnail shows before its delay: %v, blank %v", got, blank)
	}
}
