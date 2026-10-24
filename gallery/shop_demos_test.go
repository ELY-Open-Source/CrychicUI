package gallery

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestShopDemosLog(t *testing.T) {
	press := func(name string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error { return tt.Click(name) }
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"ProductCard":    {press("Add to cart"), "added"},
		"ProductGallery": {press("Image 2"), "selected 1"},
		"AddToCart":      {press("Add Arc desk lamp"), "added"},
		"CartIcon":       {press("Add item"), "count 4"},
		"CartDrawer":     {press("Checkout"), "checkout"},
		"CartRow": {func(tt *ui.Tester) error {
			if err := tt.Click("Remove"); err != nil {
				return err
			}
			tt.Advance(time.Second)
			if tt.HasText("Stoneware mug") {
				return fmt.Errorf("removed row still shows")
			}
			return tt.Click("Undo remove")
		}, "restored"},
		"CheckoutSteps": {press("Payment"), "step 2"},
		"Coupon":        {press("Apply"), "applied true"},
		"Variants":      {press("Brass (2 left)"), "color 0"},
		"ReviewCard":    {press("Helpful (14)"), "helpful"},
		"Wishlist":      {press("Save"), "saved 0 true"},
		"Categories":    {press("Desks"), "selected 3"},
		"ShopFilters":   {press("Fold"), "brands 2"},
		"SortBar": {func(tt *ui.Tester) error {
			if err := tt.Click("Sort by"); err != nil {
				return err
			}
			tt.Advance(time.Second)
			return tt.Click("Top rated")
		}, "sort Top rated"},
	}
	quiet := map[string]bool{"ProductPrice": true, "SaleBadge": true, "OrderSummary": true,
		"Stock": true, "Reviews": true, "Compare": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "shop" || e.ID >= 996 {
			continue
		}
		seen++
		if quiet[e.Name] {
			continue
		}
		act, ok := acts[e.Name]
		if !ok {
			t.Errorf("%s has no logged action", e.Name)
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
		tt.Advance(time.Second)
		if err := act.do(tt); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Frame()
		if !strings.Contains(strings.Join(logs, "\n"), act.want) {
			t.Errorf("%s logs %q, want %q", e.Name, logs, act.want)
		}
	}
	if seen != 20 {
		t.Errorf("checked %d shop demos", seen)
	}
}

func TestCouponDemoLogsEdits(t *testing.T) {
	var logs []string
	tt := ui.NewTester(func(c *ui.Context) { demoCoupon(c, func(s string) { logs = append(logs, s) }) }, 600, 300)
	tt.Advance(time.Second)
	if err := tt.Click("Coupon code"); err != nil {
		t.Fatal(err)
	}
	tt.Type("X")
	tt.Frame()
	if !slices.Contains(logs, "code 9 chars") {
		t.Errorf("logs %q", logs)
	}
}
