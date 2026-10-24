package finance

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestBillsBannerKeyFreeOfNames(t *testing.T) {
	bills, paid := []Bill{{Name: "all", Due: "Nov 1", Amount: 900}, {Name: "Phone", Due: "Nov 4", Amount: 40}}, []bool{true, true}
	tt := ui.NewTester(padded(func(c *ui.Context) { Bills(c, bills, &paid, "USD", BillsOptions{}) }), 500, 400)
	settles(t, tt)
	if !tt.HasText("All bills paid") {
		t.Errorf("no banner: %q", tt.Texts())
	}
}

func TestOrderTicketChangedOnlyOnQuantity(t *testing.T) {
	quantity, changed := 1, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if root, _, _ := OrderTicket(c, &quantity, 100, "USD", OrderTicketOptions{}); root.Changed() {
			changed++
		}
	}), 500, 400)
	settles(t, tt)
	click(t, tt, "Buy")
	click(t, tt, "Sell")
	if changed != 0 {
		t.Errorf("Buy and Sell reported %d quantity changes", changed)
	}
	click(t, tt, "Increase")
	if changed != 1 || quantity != 2 {
		t.Errorf("Increase: changed %d, quantity %d", changed, quantity)
	}
}

func TestBookMidpointNearMaxFloat(t *testing.T) {
	bids, asks := []BookLevel{{Price: 1e308, Amount: 1}}, []BookLevel{{Price: 1.1e308, Amount: 1}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		OrderBook(c, bids, asks, "USD", OrderBookOptions{})
		Depth(c, bids, asks, "USD", DepthOptions{})
	}), 500, 700)
	settles(t, tt)
}

func TestNetWorthChangeNearMaxFloat(t *testing.T) {
	points := []NetPoint{{Label: "Jan", Assets: 1e308}, {Label: "Feb", Assets: 1.1e308}}
	tt := ui.NewTester(padded(func(c *ui.Context) { NetWorth(c, points, "USD", NetWorthOptions{}) }), 500, 500)
	settles(t, tt)
	if !tt.HasText("+10.00%") {
		t.Errorf("change: %q", tt.Texts())
	}
	tiny := []NetPoint{{Label: "Jan", Assets: 1e-306}, {Label: "Feb", Assets: 10}}
	panics(t, func() { ui.Render(func(c *ui.Context) { NetWorth(c, tiny, "USD", NetWorthOptions{}) }, 500, 500, 1) })
}

func TestDepthCenterLabelMatchesAxis(t *testing.T) {
	bids, asks := []BookLevel{{Price: 90, Amount: 1}, {Price: 99, Amount: 1}}, []BookLevel{{Price: 101, Amount: 1}, {Price: 102, Amount: 1}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Depth(c, bids, asks, "USD", DepthOptions{}) }), 500, 400)
	settles(t, tt)
	if !tt.HasText("USD 96.00") || tt.HasText("USD 100.00") {
		t.Errorf("axis labels: %q", tt.Texts())
	}
}

func TestOrderBookRowsStagger(t *testing.T) {
	var bids, asks []BookLevel
	for i := range 4 {
		bids = append(bids, BookLevel{Price: float64(90 + i), Amount: 1})
		asks = append(asks, BookLevel{Price: float64(110 + i), Amount: 1})
	}
	tt := ui.NewTester(padded(func(c *ui.Context) { OrderBook(c, bids, asks, "USD", OrderBookOptions{}) }), 500, 600)
	tt.Frame()
	tt.Advance(250 * time.Millisecond)
	early := tt.Image()
	settles(t, tt)
	late := tt.Image()
	diff := func(name string) int {
		r := rectOf(t, tt, name)
		n := 0
		for y := int(r.Y); y < int(r.Y+r.H); y++ {
			for x := int(r.X); x < int(r.X+r.W); x++ {
				if early.RGBAAt(x, y) != late.RGBAAt(x, y) {
					n++
				}
			}
		}
		return n
	}
	if n := diff("USD 113.00"); n != 0 {
		t.Errorf("top ask still arriving at 250ms: %d pixels", n)
	}
	if n := diff("USD 90.00"); n == 0 {
		t.Error("bottom bid already in place at 250ms")
	}
}
