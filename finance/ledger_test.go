package finance

import (
	"math"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
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

func TestRowsKeepIdentityAcrossReorder(t *testing.T) {
	items := []WatchItem{{Symbol: "AAPL", Price: 1, Points: []float64{1, 2}}, {Symbol: "TSLA", Price: 1, Points: []float64{2, 1}}}
	sel := -1
	wt := ui.NewTester(padded(func(c *ui.Context) { Watchlist(c, items, "USD", &sel, WatchlistOptions{}) }), 500, 400)
	settles(t, wt)
	swapUnderPress(wt, rectOf(t, wt, "TSLA"), func() { items[0], items[1] = items[1], items[0] })
	if sel != -1 {
		t.Errorf("watchlist chose the row that slid under the pointer: %d", sel)
	}
	bills, paid := []Bill{{Name: "Rent", Due: "Nov 1", Amount: 900}, {Name: "Phone", Due: "Nov 4", Amount: 40}}, []bool{false, false}
	bt := ui.NewTester(padded(func(c *ui.Context) { Bills(c, bills, &paid, "USD", BillsOptions{}) }), 500, 400)
	settles(t, bt)
	swapUnderPress(bt, rectOf(t, bt, "Pay Phone"), func() { bills[0], bills[1] = bills[1], bills[0] })
	if paid[0] || paid[1] {
		t.Errorf("bills paid the row that slid under the pointer: %v", paid)
	}
}

func TestBillsPayDownToAllPaid(t *testing.T) {
	bills, paid := []Bill{{Name: "Rent", Due: "Nov 1", Amount: 900}, {Name: "Phone", Due: "Nov 4", Amount: 40}}, []bool{false, false}
	tt := ui.NewTester(padded(func(c *ui.Context) { Bills(c, bills, &paid, "USD", BillsOptions{}) }), 500, 400)
	settles(t, tt)
	if !tt.HasText("USD 940.00") {
		t.Fatalf("due total %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Pay Rent")
	click(t, tt, "Pay Phone")
	settles(t, tt)
	if !paid[0] || !paid[1] || !tt.HasText("All bills paid") || !tt.HasText("USD 0.00") {
		t.Errorf("paid %v texts %q", paid, tt.Texts())
	}
	if _, ok := tt.Find("Pay Rent"); ok {
		t.Error("a paid bill still offers Pay")
	}
	cues(t, before, sound.Success, sound.Success)
}

func TestCrossedBooksPanic(t *testing.T) {
	bids, asks := []BookLevel{{Price: 101, Amount: 1}}, testAsks()
	panics(t, func() {
		ui.Render(func(c *ui.Context) { OrderBook(c, bids, asks, "USD", OrderBookOptions{}) }, 500, 400, 1)
	})
	panics(t, func() { ui.Render(func(c *ui.Context) { Depth(c, bids, asks, "USD", DepthOptions{}) }, 500, 400, 1) })
}

func TestOrderBookShowsSpread(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { OrderBook(c, testLevels(), testAsks(), "USD", OrderBookOptions{}) }), 500, 500)
	settles(t, tt)
	if !tt.HasText("Spread USD 2.00") || !tt.HasText("USD 100.00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestOrderTicketStopsBeforeTotalOverflows(t *testing.T) {
	quantity := 2
	price := math.MaxFloat64 / 2.5
	tt := ui.NewTester(padded(func(c *ui.Context) { OrderTicket(c, &quantity, price, "USD", OrderTicketOptions{}) }), 500, 400)
	settles(t, tt)
	click(t, tt, "Increase")
	if quantity != 2 {
		t.Errorf("quantity %d: the total would overflow", quantity)
	}
}

func TestPaymentRequestPaysOnce(t *testing.T) {
	paid := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		PaymentRequest(c, 42.5, "Daniel Okafor", "USD", &paid, PaymentRequestOptions{})
	}), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Pay")
	settles(t, tt)
	click(t, tt, "Paid")
	if !paid || !tt.HasText("Daniel Okafor requests") {
		t.Errorf("paid %v texts %q", paid, tt.Texts())
	}
	cues(t, before, sound.Success)
}

func TestTickerTapeItemsDoNotTouch(t *testing.T) {
	items := []TickerItem{{Symbol: "AAPL", Price: 228.52, Change: 1.24}, {Symbol: "MSFT", Price: 431.1, Change: -0.62}}
	tt := ui.NewTester(padded(func(c *ui.Context) { TickerTape(c, items, "USD") }), 300, 120)
	tt.Frame()
	tt.Frame()
	change, next := rectOf(t, tt, "+1.24%"), rectOf(t, tt, "MSFT")
	if change.X+change.W > next.X {
		t.Errorf("change ends at %v past the next symbol at %v", change.X+change.W, next.X)
	}
}

func TestTransactionsGroupAndEmpty(t *testing.T) {
	rows := []Transaction{{Merchant: "Lyft", Category: "Transport", Amount: -18.4, When: "Today"},
		{Merchant: "Payroll", Category: "Salary", Amount: 4200, When: "Today"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Transactions(c, rows, "USD", TransactionsOptions{}) }), 500, 400)
	settles(t, tt)
	n := 0
	for _, s := range tt.Texts() {
		if s == "TODAY" {
			n++
		}
	}
	if n != 1 {
		t.Errorf("day heading shown %d times in %q", n, tt.Texts())
	}
	rows = nil
	tt.Frame()
	if !tt.HasText("No transactions yet") {
		t.Errorf("empty %q", tt.Texts())
	}
}

func TestPriceAlertBelowAtZero(t *testing.T) {
	active := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		PriceAlert(c, Below, 10, 0, "USD", &active, PriceAlertOptions{})
	}), 500, 200)
	settles(t, tt)
	if !active || !tt.HasText("Price below USD 10.00") {
		t.Errorf("active %v texts %q", active, tt.Texts())
	}
}
