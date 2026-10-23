package gallery

import (
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
		"AddToCart":      {press("Add"), "added"},
		"CartDrawer":     {press("Checkout"), "checkout"},
		"CartRow":        {press("Increase"), "quantity 2"},
		"CheckoutSteps":  {press("Payment"), "step 2"},
		"Coupon":         {press("Apply"), "applied true"},
		"Variants":       {press("Red (2 left)"), "color 0"},
		"ReviewCard":     {press("Helpful (2)"), "helpful"},
		"Wishlist":       {press("Save"), "saved true"},
		"Categories":     {press("Desks"), "selected 1"},
		"ShopFilters":    {press("B"), "false true"},
		"SortBar":        {press("Price"), "sort Price"},
	}
	quiet := map[string]bool{"ProductPrice": true, "SaleBadge": true, "CartIcon": true, "OrderSummary": true,
		"Stock": true, "Reviews": true, "Compare": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "shop" {
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
