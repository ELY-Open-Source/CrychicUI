package finance

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestAllocation(t *testing.T) {
	parts := []AllocationPart{{Name: "Stocks", Share: 0.6}, {Name: "Cash", Share: 0.4}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Allocation(c, "Allocation", parts) }), 500, 250)
	settles(t, tt)
	if !tt.HasText("Stocks") || !tt.HasText("60%") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestWallet(t *testing.T) {
	sent := -1
	coins := []Coin{{Symbol: "BTC", Amount: 0.5, Value: 20000}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Wallet(c, coins, "USD", &sent, WalletOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Send")
	if sent != 0 {
		t.Errorf("sent %d", sent)
	}
	cues(t, before, sound.Send)
}

func TestLoanCalc(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { LoanCalc(c, 1000, 0.12, 1, "USD", LoanCalcOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("USD 88.85") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestLoanCalcTinyRate(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { LoanCalc(c, 1000, 1e-16, 1, "USD", LoanCalcOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Monthly payment") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestMortgage(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Mortgage(c, 1000, 0.12, 2, "USD", MortgageOptions{}) }), 500, 250)
	settles(t, tt)
	if !tt.HasText("USD 1,000.00") || !tt.HasText("Year 2") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestTaxSummary(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { TaxSummary(c, 1000, 100, 0.1, "USD", TaxSummaryOptions{}) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("USD 90.00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestSplitBill(t *testing.T) {
	people := 2
	tt := ui.NewTester(padded(func(c *ui.Context) { SplitBill(c, 100, &people, "USD", SplitBillOptions{}) }), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Increase")
	if people != 3 || !tt.HasText("USD 33.33") {
		t.Errorf("people %d texts %q", people, tt.Texts())
	}
	cues(t, before, sound.Tick)
}

func TestSplitBillBounds(t *testing.T) {
	people := 101
	panics(t, func() {
		ui.Render(func(c *ui.Context) { SplitBill(c, 100, &people, "USD", SplitBillOptions{}) }, 500, 200, 1)
	})
}

func TestSplitBillCopy(t *testing.T) {
	people := 2
	tt := ui.NewTester(padded(func(c *ui.Context) {
		SplitBill(c, 100, &people, "USD", SplitBillOptions{Decrease: "Fewer", Increase: "More"})
	}), 500, 220)
	if !tt.HasText("Fewer") || !tt.HasText("More") {
		t.Errorf("custom labels %q", tt.Texts())
	}
}

func TestSavingsGoal(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		SavingsGoal(c, "Home", 500, 1000, "2027-01-01", "USD", SavingsGoalOptions{})
	}), 500, 220)
	settles(t, tt)
	if !tt.HasText("Target 2027-01-01") || !tt.HasText("USD 500.00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestTransfer(t *testing.T) {
	amount, confirmed := "100", false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Transfer(c, "Checking", "Savings", &amount, "USD", &confirmed, TransferOptions{})
	}), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Confirm")
	if !confirmed {
		t.Error("transfer was not confirmed")
	}
	cues(t, before, sound.Send)
}

func TestCardControls(t *testing.T) {
	frozen, online, limit := false, true, 1000.0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		CardControls(c, &frozen, &online, &limit, CardControlsOptions{})
	}), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Freeze card")
	click(t, tt, "Online payments")
	if !frozen || online {
		t.Errorf("frozen %v online %v", frozen, online)
	}
	cues(t, before, sound.Toggle, sound.Toggle)
}

func TestStatement(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Statement(c, "October", []Transaction{testTransaction()}, "USD", StatementOptions{})
	}), 500, 220)
	settles(t, tt)
	if !tt.HasText("October") || !tt.HasText("USD -24.00") {
		t.Errorf("texts %q", tt.Texts())
	}
}
