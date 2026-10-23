package shop

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestCartDrawer(t *testing.T) {
	checkout := false
	items := []CartItem{{Product: testProduct(), Quantity: 2}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		CartDrawer(c, items, "USD", &checkout, CartDrawerOptions{})
	}), 500, 300)
	settles(t, tt)
	if !tt.HasText("USD 160.00") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Checkout")
	if !checkout {
		t.Error("checkout did not open")
	}
	cues(t, before, sound.Open)
}

func TestCartRow(t *testing.T) {
	item := CartItem{Product: testProduct(), Quantity: 1}
	removed := false
	tt := ui.NewTester(padded(func(c *ui.Context) { CartRow(c, &item, "USD", &removed, CartRowOptions{}) }), 600, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Increase")
	if item.Quantity != 2 {
		t.Errorf("quantity %d", item.Quantity)
	}
	click(t, tt, "Remove")
	if !removed {
		t.Error("row was not removed")
	}
	cues(t, before, sound.Tick, sound.Delete)
}

func TestCartRowCollapse(t *testing.T) {
	item := CartItem{Product: testProduct(), Quantity: 1}
	removed := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if !removed {
			c.Key("row")
			CartRow(c, &item, "USD", &removed, CartRowOptions{})
		}
		ui.Text(c, "after")
	}), 600, 250)
	settles(t, tt)
	before, ok := tt.Find("after")
	if !ok {
		t.Fatal("following text missing")
	}
	click(t, tt, "Remove")
	tt.Advance(50 * time.Millisecond)
	after, ok := tt.Find("after")
	t.Logf("collapse before %v after %v", before.Y, after.Y)
	if !tt.FramePending() {
		t.Error("removed row has no collapse transition")
	}
	settles(t, tt)
}

func TestCartRowCallerKeys(t *testing.T) {
	first := CartItem{Product: testProduct(), Quantity: 1}
	second := CartItem{Product: testProduct(), Quantity: 2}
	removed := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		c.Key("first")
		CartRow(c, &first, "USD", &removed, CartRowOptions{})
		c.Key("second")
		CartRow(c, &second, "USD", &removed, CartRowOptions{})
	}), 600, 350)
	settles(t, tt)
	if !tt.HasText("Desk lamp") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestCheckoutSteps(t *testing.T) {
	current := Cart
	tt := ui.NewTester(padded(func(c *ui.Context) { CheckoutSteps(c, &current, CheckoutStepsOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Payment")
	if current != Payment {
		t.Errorf("step %v", current)
	}
	cues(t, before, sound.Tick)
}

func TestOrderSummary(t *testing.T) {
	lines := []SummaryLine{{Label: "Subtotal", Amount: 100}, {Label: "Discount", Amount: -10}, {Label: "Tax", Amount: 5}}
	tt := ui.NewTester(padded(func(c *ui.Context) { OrderSummary(c, lines, "USD", OrderSummaryOptions{}) }), 500, 250)
	settles(t, tt)
	if !tt.HasText("USD 95.00") || !tt.HasText("USD -10.00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestCoupon(t *testing.T) {
	code, applied := "SAVE", ""
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Coupon(c, &code, []string{"SAVE"}, &applied, CouponOptions{})
	}), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Apply")
	if applied != "SAVE" {
		t.Errorf("applied %q", applied)
	}
	code = "BAD"
	click(t, tt, "Apply")
	if !tt.HasText("Invalid coupon") {
		t.Errorf("texts %q", tt.Texts())
	}
	cues(t, before, sound.Success, sound.Error)
}

func TestCouponEdit(t *testing.T) {
	code, applied, changed := "", "", 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, _ := Coupon(c, &code, []string{"SAVE"}, &applied, CouponOptions{})
		if root.Changed() {
			changed++
		}
	}), 500, 220)
	settles(t, tt)
	click(t, tt, "Coupon code")
	tt.Type("S")
	tt.Frame()
	if changed != 1 {
		t.Errorf("edits reported %d changes", changed)
	}
}

func TestVariants(t *testing.T) {
	colors := []VariantStock{{Name: "Red", Stock: 2}, {Name: "Blue", Stock: 0}}
	sizes := []VariantStock{{Name: "M", Stock: 3}}
	color, size := -1, -1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Variants(c, colors, sizes, &color, &size, VariantsOptions{})
	}), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Red (2 left)")
	click(t, tt, "M (3 left)")
	if color != 0 || size != 0 {
		t.Errorf("color %d size %d", color, size)
	}
	cues(t, before, sound.Select, sound.Select)
}

func TestStock(t *testing.T) {
	for count, want := range map[int]string{0: "Out of stock", 2: "Low stock", 8: "In stock"} {
		tt := ui.NewTester(padded(func(c *ui.Context) { Stock(c, count, StockOptions{}) }), 250, 120)
		settles(t, tt)
		if !tt.HasText(want) {
			t.Errorf("count %d texts %q", count, tt.Texts())
		}
	}
}
