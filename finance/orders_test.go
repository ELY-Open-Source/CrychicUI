package finance

import (
	"math"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestOrderTicket(t *testing.T) {
	quantity := 1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		OrderTicket(c, &quantity, 100, "USD", OrderTicketOptions{})
	}), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Increase")
	click(t, tt, "Buy")
	if quantity != 2 {
		t.Errorf("quantity %d", quantity)
	}
	cues(t, before, sound.Tick, sound.Success)
}

func TestOrderTicketOverflow(t *testing.T) {
	quantity := math.MaxInt
	tt := ui.NewTester(padded(func(c *ui.Context) {
		OrderTicket(c, &quantity, 100, "USD", OrderTicketOptions{})
	}), 500, 250)
	settles(t, tt)
	click(t, tt, "Increase")
	if quantity != math.MaxInt {
		t.Errorf("quantity overflowed to %d", quantity)
	}
}

func TestOrderBook(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		OrderBook(c, testLevels(), testAsks(), "USD", OrderBookOptions{})
	}), 500, 300)
	settles(t, tt)
	if !tt.HasText("Bids") || !tt.HasText("USD 99.00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestOrderBookFlashLifetime(t *testing.T) {
	rows := testLevels()
	tt := ui.NewTester(padded(func(c *ui.Context) {
		OrderBook(c, rows, testAsks(), "USD", OrderBookOptions{})
	}), 500, 300)
	settles(t, tt)
	settles(t, tt)
	r, ok := tt.Find("USD 99.00")
	if !ok {
		t.Fatal("price missing")
	}
	x, y := int(r.X+r.W-2), int(r.Y+r.H/2)
	before := tt.Image().RGBAAt(x, y)
	rows[0].Amount = 8
	rows[0].Price = 100
	tt.Frame()
	tt.Advance(100 * time.Millisecond)
	if got := tt.Image().RGBAAt(x, y); got == before {
		t.Errorf("row flash not visible at 100ms: %v", got)
	}
}

func TestDepth(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Depth(c, testLevels(), testAsks(), "USD", DepthOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("Bids") || !tt.HasText("Asks") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestPriceAlert(t *testing.T) {
	active := false
	before := baseline()
	tt := ui.NewTester(padded(func(c *ui.Context) {
		PriceAlert(c, Above, 100, 120, "USD", &active, PriceAlertOptions{})
	}), 500, 180)
	settles(t, tt)
	if !active || !tt.HasText("Price alert triggered") {
		t.Errorf("active %v texts %q", active, tt.Texts())
	}
	cues(t, before, sound.Notify)
}

func TestConverter(t *testing.T) {
	amount := "10"
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Converter(c, &amount, 2, "USD", "EUR", ConverterOptions{})
	}), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "USD")
	tt.Key(ui.Ctrl, ui.KeyA)
	tt.Type("12")
	tt.Frame()
	if !tt.HasText("EUR 24.00") {
		t.Errorf("texts %q", tt.Texts())
	}
	cues(t, before, sound.Tick)
}

func TestConverterOverflow(t *testing.T) {
	amount := "1e308"
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Converter(c, &amount, 2, "USD", "EUR", ConverterOptions{})
	}), 500, 220)
	settles(t, tt)
	if !tt.HasText("Enter a number") {
		t.Errorf("overflow texts %q", tt.Texts())
	}
}

func TestInvoice(t *testing.T) {
	lines := []Charge{{Label: "Service", Amount: 100}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Invoice(c, "INV-1", "Due", lines, 10, "USD", InvoiceOptions{})
	}), 500, 300)
	settles(t, tt)
	if !tt.HasText("INV-1") || !tt.HasText("USD 110.00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestPaymentRequest(t *testing.T) {
	paid := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		PaymentRequest(c, 100, "Ana", "USD", &paid, PaymentRequestOptions{})
	}), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Pay")
	if !paid {
		t.Error("request was not paid")
	}
	cues(t, before, sound.Success)
}

func TestBills(t *testing.T) {
	paid := []bool{false}
	bills := []Bill{{Name: "Electricity", Due: "Friday", Amount: 80}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Bills(c, bills, &paid, "USD", BillsOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Pay")
	if !paid[0] || !tt.HasText("Paid") {
		t.Errorf("paid %v texts %q", paid, tt.Texts())
	}
	cues(t, before, sound.Success)
}

func TestSubscriptions(t *testing.T) {
	rows := []Subscription{{Name: "Music", Next: "Nov 1", Amount: 10}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Subscriptions(c, rows, "USD", SubscriptionsOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Music") || !tt.HasText("Nov 1") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestNetWorth(t *testing.T) {
	points := []NetPoint{{Label: "Jan", Assets: 100, Debts: 20}, {Label: "Feb", Assets: 150, Debts: 50}}
	tt := ui.NewTester(padded(func(c *ui.Context) { NetWorth(c, points, "USD", NetWorthOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("USD 100.00") || !tt.HasText("Net worth") {
		t.Errorf("texts %q", tt.Texts())
	}
}
