package finance

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestCashFlow(t *testing.T) {
	months := []FlowMonth{{Label: "Jan", In: 100, Out: 40}}
	tt := ui.NewTester(padded(func(c *ui.Context) { CashFlow(c, months, "USD") }), 500, 220)
	settles(t, tt)
	if !tt.HasText("Jan") || !tt.HasText("USD 100.00") || !tt.HasText("USD 40.00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestDividends(t *testing.T) {
	payments := []Dividend{{Symbol: "MYGO", Date: "Nov 1", Amount: 5}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Dividends(c, payments, "USD") }), 500, 180)
	settles(t, tt)
	if !tt.HasText("MYGO") || !tt.HasText("Nov 1") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestRiskMeter(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { RiskMeter(c, 4, RiskMeterOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Low risk") || !tt.HasText("High risk") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestRiskMeterMotion(t *testing.T) {
	level := 1
	tt := ui.NewTester(padded(func(c *ui.Context) { RiskMeter(c, level, RiskMeterOptions{}) }), 500, 180)
	settles(t, tt)
	level = 4
	tt.Frame()
	if !tt.FramePending() {
		t.Error("risk colors do not animate")
	}
}

func TestExchangeRate(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { ExchangeRate(c, "USD/EUR", 0.85, 1.2) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("USD/EUR") || !tt.HasText("0.8500") || !tt.HasText("+1.20%") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestPayees(t *testing.T) {
	query, selected := "", -1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Payees(c, []string{"Ana", "Kai"}, &query, &selected, PayeesOptions{})
	}), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Kai")
	if selected != 1 {
		t.Errorf("selected %d", selected)
	}
	cues(t, before, sound.Select)
}

func TestCategorize(t *testing.T) {
	selected := -1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Categorize(c, testTransaction(), []string{"Books", "Food"}, &selected)
	}), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Books")
	if selected != 0 {
		t.Errorf("selected %d", selected)
	}
	cues(t, before, sound.Select)
}

func TestReceipts(t *testing.T) {
	receipts := []Receipt{{Merchant: "Bookstore", Date: "Oct 1", Amount: 24}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Receipts(c, receipts, "USD", ReceiptsOptions{}) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("Bookstore") || !tt.HasText("USD 24.00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestCreditScore(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { CreditScore(c, 800, CreditScoreOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("800") || !tt.HasText("Excellent") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestMarketStatus(t *testing.T) {
	next := time.Now().Add(time.Hour)
	tt := ui.NewTester(padded(func(c *ui.Context) { MarketStatus(c, true, next, MarketStatusOptions{}) }), 500, 180)
	tt.Frame()
	if !tt.HasText("Market open") || !tt.HasText("Next") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestMarketStatusCopy(t *testing.T) {
	next := time.Date(2026, 12, 1, 9, 30, 0, 0, time.UTC)
	tt := ui.NewTester(padded(func(c *ui.Context) {
		MarketStatus(c, false, next, MarketStatusOptions{DateFormat: "2006-01-02"})
	}), 500, 180)
	if !tt.HasText("Next 2026-12-01") {
		t.Errorf("custom date format %q", tt.Texts())
	}
}

func TestPnL(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { PnL(c, -25, "USD", PnLOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("USD -25.00") || !tt.HasText("Profit and loss") {
		t.Errorf("texts %q", tt.Texts())
	}
}
