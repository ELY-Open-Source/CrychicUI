package shop

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestReviews(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Reviews(c, testReviews(), ReviewsOptions{}) }), 500, 450)
	settles(t, tt)
	if !tt.HasText("4.0") || !tt.HasText("2 reviews") || !tt.HasText("Great") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestReviewCard(t *testing.T) {
	helpful := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ReviewCard(c, testReviews()[0], &helpful, ReviewCardOptions{})
	}), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Helpful (2)")
	if !helpful || !tt.HasText("Helpful (3)") {
		t.Errorf("helpful %v texts %q", helpful, tt.Texts())
	}
	cues(t, before, sound.Pop)
}

func TestReviewCardOverflow(t *testing.T) {
	helpful := false
	review := testReviews()[0]
	review.Helpful = math.MaxInt
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ReviewCard(c, review, &helpful, ReviewCardOptions{})
	}), 500, 250)
	settles(t, tt)
	click(t, tt, fmt.Sprintf("Helpful (%d)", math.MaxInt))
	if helpful || !tt.HasText(fmt.Sprintf("Helpful (%d)", math.MaxInt)) {
		t.Errorf("helpful overflow %v texts %q", helpful, tt.Texts())
	}
}

func TestWishlist(t *testing.T) {
	saved := false
	tt := ui.NewTester(padded(func(c *ui.Context) { Wishlist(c, testProduct(), &saved, WishlistOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Save")
	if !saved {
		t.Error("product was not saved")
	}
	cues(t, before, sound.Toggle)
}

func TestCompare(t *testing.T) {
	products := []Product{testProduct(), testProduct()}
	products[1].Name = "Floor lamp"
	specs := []CompareSpec{{Name: "Material", Values: []string{"Wood", "Metal"}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Compare(c, products, specs, "USD") }), 500, 300)
	settles(t, tt)
	if !tt.HasText("Floor lamp") || !tt.HasText("Metal") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestCategories(t *testing.T) {
	selected := -1
	p := testProduct()
	categories := []Category{{Name: "Lighting", Image: p.Image}, {Name: "Desks", Image: p.Image}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Categories(c, categories, &selected) }), 500, 350)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Desks")
	if selected != 1 {
		t.Errorf("selected %d", selected)
	}
	cues(t, before, sound.Tap)
}

func TestShopFilters(t *testing.T) {
	low, high, rating := 0.0, 1000.0, 0
	chosen := []bool{false, false}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ShopFilters(c, &low, &high, []string{"A", "B"}, &chosen, &rating, ShopFiltersOptions{})
	}), 500, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "B")
	click(t, tt, "4 stars")
	if !chosen[1] || rating != 4 {
		t.Errorf("chosen %v rating %d", chosen, rating)
	}
	cues(t, before, sound.Toggle, sound.Toggle)
}

func TestShopFiltersRatingVisible(t *testing.T) {
	low, high, rating := 0.0, 1000.0, 1
	chosen := []bool{false}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ShopFilters(c, &low, &high, []string{"A"}, &chosen, &rating, ShopFiltersOptions{})
	}), 500, 400)
	settles(t, tt)
	before := bytes.Clone(tt.Image().Pix)
	rating = 4
	tt.Frame()
	settles(t, tt)
	if bytes.Equal(before, tt.Image().Pix) {
		t.Error("rating selection is not visible")
	}
}

func TestShopFiltersCopy(t *testing.T) {
	low, high, rating := 0.0, 1000.0, 0
	chosen := []bool{false}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ShopFilters(c, &low, &high, []string{"A"}, &chosen, &rating, ShopFiltersOptions{Stars: "%d stars minimum"})
	}), 500, 400)
	if !tt.HasText("4 stars minimum") {
		t.Errorf("custom rating label %q", tt.Texts())
	}
}

func TestSortBar(t *testing.T) {
	sort := "Featured"
	tt := ui.NewTester(padded(func(c *ui.Context) {
		SortBar(c, 2, &sort, []string{"Featured", "Price"}, SortBarOptions{})
	}), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Price")
	if sort != "Price" || !tt.HasText("2 products") {
		t.Errorf("sort %q texts %q", sort, tt.Texts())
	}
	cues(t, before, sound.Select)
}
