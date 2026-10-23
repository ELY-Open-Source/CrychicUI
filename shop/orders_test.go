package shop

import (
	"fmt"
	"image/color"
	"math"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestOrderStatus(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { OrderStatus(c, Shipped, OrderStatusOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Shipped") || !tt.HasText("Delivered") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestOrderCard(t *testing.T) {
	opened := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		OrderCard(c, "Order 1042", time.Now(), 95, "USD", Shipped, &opened, OrderCardOptions{})
	}), 600, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "View order")
	if !opened {
		t.Error("order did not open")
	}
	cues(t, before, sound.Tap)
}

func TestOrderCardCopy(t *testing.T) {
	opened := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		OrderCard(c, "Order 1042", time.Now(), 95, "USD", Placed, &opened, OrderCardOptions{Placed: "Custom placed"})
	}), 600, 300)
	if !tt.HasText("Custom placed") {
		t.Errorf("custom status %q", tt.Texts())
	}
}

func TestTracking(t *testing.T) {
	events := []TrackEvent{{Label: "Shipped", Location: "Shanghai", Time: time.Now(), Done: true}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Tracking(c, events) }), 500, 200)
	settles(t, tt)
	if !tt.HasText("Shanghai") || !tt.HasText("Shipped") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestDelivery(t *testing.T) {
	selected := 0
	choices := []DeliveryOption{{Name: "Standard", Date: "Friday", Price: 5}, {Name: "Express", Date: "Tomorrow", Price: 12}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Delivery(c, "Delivery", choices, "USD", &selected) }), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Express")
	if selected != 1 {
		t.Errorf("selected %d", selected)
	}
	cues(t, before, sound.Select)
}

func TestFreeShipping(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { FreeShipping(c, 80, 100, "USD", FreeShippingOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("USD 20.00 away from free shipping") {
		t.Errorf("texts %q", tt.Texts())
	}
	free := ui.NewTester(padded(func(c *ui.Context) { FreeShipping(c, 100, 100, "USD", FreeShippingOptions{}) }), 500, 180)
	if !free.HasText("Free shipping unlocked") {
		t.Errorf("texts %q", free.Texts())
	}
}

func TestUpsell(t *testing.T) {
	added := -1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Upsell(c, "Related", []Product{testProduct()}, "USD", &added, UpsellOptions{})
	}), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Add")
	if added != 0 {
		t.Errorf("added %d", added)
	}
	cues(t, before, sound.Pop)
}

func TestRecently(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Recently(c, "Recent", []Product{testProduct()}, "USD") }), 500, 300)
	settles(t, tt)
	if !tt.HasText("Desk lamp") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestFlashSale(t *testing.T) {
	ends := time.Now().Add(90 * time.Minute)
	tt := ui.NewTester(padded(func(c *ui.Context) { FlashSale(c, "Flash sale", ends, FlashSaleOptions{}) }), 500, 180)
	settles(t, tt)
	before := tt.Texts()
	tt.Advance(1500 * time.Millisecond)
	if !tt.HasText("Flash sale") || len(tt.Texts()) == len(before) && tt.Texts()[1] == before[1] {
		t.Errorf("countdown did not update: %q", tt.Texts())
	}
}

func TestFlashSaleDeadline(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		FlashSale(c, "Flash sale", time.Now().Add(500*time.Millisecond), FlashSaleOptions{})
	}), 500, 180)
	tt.Advance(600 * time.Millisecond)
	if !tt.HasText("Sale ended") {
		t.Errorf("deadline texts %q", tt.Texts())
	}
}

func TestGiftCard(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { GiftCard(c, 50, "USD", "For you") }), 400, 220)
	settles(t, tt)
	if !tt.HasText("USD 50.00") || !tt.HasText("For you") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestLoyalty(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Loyalty(c, 200, "Gold", 50, LoyaltyOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("200") || !tt.HasText("50 points to the next tier") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestLoyaltyOverflow(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Loyalty(c, math.MaxInt, "Gold", 1, LoyaltyOptions{}) }), 500, 180)
	settles(t, tt)
	r, ok := tt.Find(fmt.Sprint(math.MaxInt))
	if !ok {
		t.Fatal("points missing")
	}
	accent := ui.LightTheme().Accent
	want := color.RGBA{accent.R, accent.G, accent.B, 255}
	if got := tt.Image().RGBAAt(int(r.X+2), int(r.Y+r.H+10)); got != want {
		t.Errorf("progress got %v want %v", got, want)
	}
}

func TestBundle(t *testing.T) {
	selected := []bool{true}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Bundle(c, []Product{testProduct()}, &selected, "USD", BundleOptions{})
	}), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Desk lamp")
	if selected[0] || !tt.HasText("USD 0.00") {
		t.Errorf("selected %v texts %q", selected, tt.Texts())
	}
	cues(t, before, sound.Toggle)
}

func TestMiniCart(t *testing.T) {
	open := false
	items := []CartItem{{Product: testProduct(), Quantity: 1}}
	tt := ui.NewTester(padded(func(c *ui.Context) { MiniCart(c, items, "USD", &open, MiniCartOptions{}) }), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Cart")
	if !open || !tt.HasText("Desk lamp") {
		t.Errorf("open %v texts %q", open, tt.Texts())
	}
	cues(t, before, sound.Open)
}

func TestMiniCartPopover(t *testing.T) {
	open := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		MiniCart(c, nil, "USD", &open, MiniCartOptions{})
		ui.Text(c, "after")
	}), 500, 300)
	settles(t, tt)
	before, _ := tt.Find("after")
	click(t, tt, "Cart")
	after, _ := tt.Find("after")
	if after.Y != before.Y {
		t.Errorf("popover moved content: before %v after %v", before.Y, after.Y)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open {
		t.Error("Escape left popover open")
	}
}

func TestReceiptSoldOut(t *testing.T) {
	product := testProduct()
	product.Stock = 0
	items := []CartItem{{Product: product, Quantity: 1}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Receipt(c, "R-42", items, "USD", ReceiptOptions{}) }), 500, 250)
	settles(t, tt)
	if !tt.HasText("USD 80.00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestReceipt(t *testing.T) {
	items := []CartItem{{Product: testProduct(), Quantity: 2}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Receipt(c, "R-42", items, "USD", ReceiptOptions{}) }), 500, 250)
	settles(t, tt)
	if !tt.HasText("R-42") || !tt.HasText("USD 160.00") {
		t.Errorf("texts %q", tt.Texts())
	}
}
