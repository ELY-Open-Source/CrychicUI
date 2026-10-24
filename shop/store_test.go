package shop

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func named(name string, price float64) Product {
	p := testProduct()
	p.Name, p.Price, p.OldPrice = name, price, 0
	return p
}

func TestStorePicksKeepIdentityAcrossReorder(t *testing.T) {
	choices, sel := []DeliveryOption{{Name: "Standard", Date: "Fri"}, {Name: "Express", Date: "Tue", Price: 12}}, 0
	dt := ui.NewTester(padded(func(c *ui.Context) { Delivery(c, "Speed", choices, "USD", &sel, DeliveryOptions{}) }), 500, 400)
	settles(t, dt)
	swapUnderPress(dt, rectOf(t, dt, "Express"), func() { choices[0], choices[1] = choices[1], choices[0] })
	if sel != 0 {
		t.Errorf("delivery chose the option that slid under the pointer: %d", sel)
	}

	items, picked := []Product{named("Lamp", 80), named("Tray", 40)}, []bool{false, false}
	bt := ui.NewTester(padded(func(c *ui.Context) { Bundle(c, items, &picked, "USD", BundleOptions{}) }), 500, 500)
	settles(t, bt)
	swapUnderPress(bt, rectOf(t, bt, "Tray"), func() { items[0], items[1] = items[1], items[0] })
	if picked[0] || picked[1] {
		t.Errorf("bundle toggled the item that slid under the pointer: %v", picked)
	}

	goods, reasons, item, reason := []Product{named("Lamp", 80), named("Tray", 40)}, []string{"Damaged", "Late"}, -1, ""
	rt := ui.NewTester(padded(func(c *ui.Context) { Returns(c, goods, "USD", reasons, &item, &reason, ReturnsOptions{}) }), 500, 500)
	settles(t, rt)
	swapUnderPress(rt, rectOf(t, rt, "Tray"), func() { goods[0], goods[1] = goods[1], goods[0] })
	swapUnderPress(rt, rectOf(t, rt, "Late"), func() { reasons[0], reasons[1] = reasons[1], reasons[0] })
	if item != -1 || reason != "" {
		t.Errorf("returns chose what slid under the pointer: item %d reason %q", item, reason)
	}

	related, added := []Product{named("Mug", 24), named("Vase", 36)}, []bool{false, false}
	ut := ui.NewTester(padded(func(c *ui.Context) { Upsell(c, "More", related, "USD", &added, UpsellOptions{}) }), 500, 500)
	settles(t, ut)
	swapUnderPress(ut, rectOf(t, ut, "Add Vase"), func() { related[0], related[1] = related[1], related[0] })
	if added[0] || added[1] {
		t.Errorf("upsell added the card that slid under the pointer: %v", added)
	}
}

func TestUpsellAddsOnce(t *testing.T) {
	added := []bool{false}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Upsell(c, "More", []Product{named("Mug", 24)}, "USD", &added, UpsellOptions{})
	}), 260, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Add Mug")
	click(t, tt, "Added Mug")
	if !added[0] {
		t.Error("mug was not added")
	}
	cues(t, before, sound.Pop)
}

func TestStoreHoursNextChange(t *testing.T) {
	hours := []DayHours{{Day: "Sun"}, {Day: "Mon", Open: "09:00", Close: "18:00"}, {Day: "Tue"}, {Day: "Wed"},
		{Day: "Thu", Open: "09:00", Close: "17:00"}, {Day: "Fri", Open: "22:00", Close: "02:00"}, {Day: "Sat"}}
	cases := []struct {
		now  time.Time
		want string
	}{
		{time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC), "Opens Thu at 09:00"}, // Wednesday
		{time.Date(2026, 10, 22, 7, 0, 0, 0, time.UTC), "Opens Thu at 09:00"},  // Thursday morning
		{time.Date(2026, 10, 22, 10, 0, 0, 0, time.UTC), "Closes at 17:00"},
		{time.Date(2026, 10, 24, 1, 0, 0, 0, time.UTC), "Closes at 02:00"}, // Friday's late shift
		{time.Date(2026, 10, 24, 3, 0, 0, 0, time.UTC), "Opens Mon at 09:00"},
	}
	for _, tc := range cases {
		tt := ui.NewTester(padded(func(c *ui.Context) { StoreHours(c, tc.now, hours, StoreHoursOptions{}) }), 500, 450)
		tt.Frame()
		if !tt.HasText(tc.want) {
			t.Errorf("%v: want %q in %q", tc.now, tc.want, tt.Texts())
		}
	}
}

func TestBundleSavesOnlyWithTwo(t *testing.T) {
	items, picked := []Product{named("Lamp", 80), named("Tray", 40)}, []bool{true, true}
	tt := ui.NewTester(padded(func(c *ui.Context) { Bundle(c, items, &picked, "USD", BundleOptions{Off: 0.1}) }), 500, 500)
	settles(t, tt)
	if !tt.HasText("Save USD 12.00") || !tt.HasText("USD 108.00") {
		t.Errorf("two items: %q", tt.Texts())
	}
	click(t, tt, "Tray")
	settles(t, tt)
	if tt.HasText("Save") || !tt.HasText("USD 80.00") {
		t.Errorf("one item: %q", tt.Texts())
	}
}

func TestGiftCardShowsLastFourOfCode(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		GiftCard(c, 50, "USD", "Enjoy", GiftCardOptions{From: "From Maya", Code: "NW-GIFT-4821"})
	}), 400, 300)
	settles(t, tt)
	if !tt.HasText("•••• 4821") || strings.Contains(strings.Join(tt.Texts(), " "), "NW-GIFT") || !tt.HasText("From Maya") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestLoyaltyTopTier(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Loyalty(c, 5000, "Platinum", 0, LoyaltyOptions{NextTier: "None"}) }), 500, 250)
	settles(t, tt)
	if !tt.HasText("Top tier reached") || tt.HasText("None") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestMiniCartCountsAndEmpties(t *testing.T) {
	big := testProduct()
	big.Stock = math.MaxInt
	open := true
	items := []CartItem{{Product: big, Quantity: math.MaxInt}, {Product: named("Mug", 24), Quantity: 1}}
	tt := ui.NewTester(padded(func(c *ui.Context) { MiniCart(c, items, "USD", &open, MiniCartOptions{}) }), 500, 500)
	settles(t, tt)
	if !tt.HasText("9223372036854775808") {
		t.Errorf("count %q", tt.Texts())
	}
	items = nil
	tt.Frame()
	settles(t, tt)
	if !tt.HasText("Your cart is empty") {
		t.Errorf("empty %q", tt.Texts())
	}
}

func TestDeliveryReadsFreeAndOrderCardCountsMore(t *testing.T) {
	sel := 0
	dt := ui.NewTester(padded(func(c *ui.Context) {
		Delivery(c, "Speed", []DeliveryOption{{Name: "Standard", Date: "Fri"}}, "USD", &sel, DeliveryOptions{Free: "No charge"})
	}), 500, 200)
	settles(t, dt)
	if !dt.HasText("No charge") || dt.HasText("USD 0.00") {
		t.Errorf("delivery %q", dt.Texts())
	}
	img := testProduct().Image
	opened := false
	ot := ui.NewTester(padded(func(c *ui.Context) {
		OrderCard(c, "Order 1", time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC), 10, "USD", Packed, &opened,
			OrderCardOptions{Images: []*ui.Bitmap{img, img, img, img, img, img}})
	}), 500, 300)
	settles(t, ot)
	if !ot.HasText("+2") || !ot.HasText("Oct 20, 2026") || !ot.HasText("Packed") {
		t.Errorf("order card %q", ot.Texts())
	}
}

func TestTrackingNeedsDoneEventsFirst(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("a done event after a pending one was accepted")
		}
	}()
	events := []TrackEvent{{Label: "Shipped", Location: "Reno"}, {Label: "Packed", Location: "Reno", Done: true}}
	ui.NewTester(func(c *ui.Context) { Tracking(c, events, TrackingOptions{}) }, 300, 300).Frame()
}

func TestFlashSaleTilesTickAndEnd(t *testing.T) {
	var ends time.Time
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if ends.IsZero() {
			ends = c.Now().Add(3*time.Second + 500*time.Millisecond)
		}
		FlashSale(c, "Sale", ends, FlashSaleOptions{})
	}), 500, 200)
	tt.Frame()
	if !tt.HasText("04") {
		t.Fatalf("start %q", tt.Texts())
	}
	tt.Advance(time.Second)
	if !tt.HasText("03") {
		t.Errorf("after 1s %q", tt.Texts())
	}
	tt.Advance(3 * time.Second)
	if !tt.HasText("Sale ended") {
		t.Errorf("after the end %q", tt.Texts())
	}
}

func TestFlashSaleFarDeadlineStaysPositive(t *testing.T) {
	ends := time.Date(2500, 1, 1, 0, 0, 0, 0, time.UTC)
	tt := ui.NewTester(padded(func(c *ui.Context) { FlashSale(c, "Sale", ends, FlashSaleOptions{}) }), 600, 200)
	tt.Frame()
	for _, s := range tt.Texts() {
		if strings.HasPrefix(s, "-") {
			t.Fatalf("negative tile %q in %q", s, tt.Texts())
		}
	}
}

func TestStoreHoursRejectsOverlap(t *testing.T) {
	for _, close := range []string{"12:00", "09:00"} {
		hours := []DayHours{{Day: "Sun"}, {Day: "Mon"}, {Day: "Tue"}, {Day: "Wed"}, {Day: "Thu"},
			{Day: "Fri", Open: "22:00", Close: close}, {Day: "Sat", Open: "09:00", Close: "18:00"}}
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("a shift closing at %s, at or past the next opening, was accepted", close)
				}
			}()
			ui.NewTester(func(c *ui.Context) { StoreHours(c, time.Now(), hours, StoreHoursOptions{}) }, 300, 300).Frame()
		}()
	}
}

// reaches scrolls a popover list until its last row shows.
func reaches(t *testing.T, tt *ui.Tester, first, last string) {
	t.Helper()
	r := rectOf(t, tt, first)
	for range 8 {
		tt.Scroll(r.X+r.W/2, r.Y+r.H/2, 0, 400)
		tt.Advance(500 * time.Millisecond)
	}
	if l, ok := tt.Find(last); !ok || l.H == 0 {
		t.Errorf("%q out of reach: %v", last, l)
	}
}

func TestLongPopoversScroll(t *testing.T) {
	var items []CartItem
	var rows []SizeRow
	var choices []string
	for i := 1; i <= 40; i++ {
		items = append(items, CartItem{Product: named(fmt.Sprintf("Item %02d", i), 10), Quantity: 1})
		rows = append(rows, SizeRow{Size: fmt.Sprintf("S%02d", i), Chest: "90 cm", Waist: "80 cm"})
		choices = append(choices, fmt.Sprintf("Choice %02d", i))
	}
	open := true
	mt := ui.NewTester(padded(func(c *ui.Context) { MiniCart(c, items, "USD", &open, MiniCartOptions{}) }), 500, 600)
	settles(t, mt)
	reaches(t, mt, "Item 01", "Item 40")
	sized := true
	st := ui.NewTester(padded(func(c *ui.Context) { SizeGuide(c, rows, &sized, SizeGuideOptions{}) }), 500, 600)
	settles(t, st)
	reaches(t, st, "S01", "S40")
	sort := "Choice 01"
	bt := ui.NewTester(padded(func(c *ui.Context) { SortBar(c, 9, &sort, choices, SortBarOptions{}) }), 500, 600)
	settles(t, bt)
	click(t, bt, "Sort by")
	settles(t, bt)
	reaches(t, bt, "Choice 02", "Choice 40")
}

func TestReturnReasonFits(t *testing.T) {
	long := "The shade arrived cracked and the arm was bent, so I would like a full refund please"
	item, reason := -1, ""
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ui.Column(c).Width(400).Children(func() {
			Returns(c, []Product{named("Lamp", 80)}, "USD", []string{"Damaged", long}, &item, &reason, ReturnsOptions{})
		})
	}), 600, 500)
	settles(t, tt)
	if r := rectOf(t, tt, long); r.X+r.W > 411 {
		t.Errorf("long reason ends at %v past 410", r.X+r.W)
	}
}

func TestMiniCartBadgeCaps(t *testing.T) {
	big := testProduct()
	big.Stock = math.MaxInt
	open := false
	items := []CartItem{{Product: big, Quantity: 120}}
	tt := ui.NewTester(padded(func(c *ui.Context) { MiniCart(c, items, "USD", &open, MiniCartOptions{}) }), 500, 500)
	settles(t, tt)
	if _, ok := tt.Find("99+"); !ok || tt.HasText("120") {
		t.Errorf("badge %q", tt.Texts())
	}
	click(t, tt, "Cart")
	settles(t, tt)
	if !tt.HasText("120") {
		t.Errorf("full count missing from the panel: %q", tt.Texts())
	}
}

func TestAskAnswersHover(t *testing.T) {
	questions, draft := []Question{}, "Is it dimmable?"
	tt := ui.NewTester(padded(func(c *ui.Context) { QandA(c, &questions, &draft, QandAOptions{}) }), 500, 300)
	settles(t, tt)
	before := bytes.Clone(tt.Image().Pix)
	r := rectOf(t, tt, "Ask")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	settles(t, tt)
	if bytes.Equal(before, tt.Image().Pix) {
		t.Error("hovering Ask changed nothing")
	}
}
