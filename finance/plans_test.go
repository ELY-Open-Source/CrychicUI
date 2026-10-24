package finance

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func count(tt *ui.Tester, text string) int {
	n := 0
	for _, s := range tt.Texts() {
		if s == text {
			n++
		}
	}
	return n
}

func hover(t *testing.T, tt *ui.Tester, name string) {
	t.Helper()
	r := rectOf(t, tt, name)
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Advance(time.Second)
}

func TestPickersKeepIdentityAcrossReorder(t *testing.T) {
	coins := []Coin{{Symbol: "BTC", Value: 1}, {Symbol: "ETH", Value: 1}}
	sent := -1
	wt := ui.NewTester(padded(func(c *ui.Context) { Wallet(c, coins, "USD", &sent, WalletOptions{}) }), 500, 400)
	settles(t, wt)
	swapUnderPress(wt, rectOf(t, wt, "Send ETH"), func() { coins[0], coins[1] = coins[1], coins[0] })
	if sent != -1 {
		t.Errorf("wallet sent the coin that slid under the pointer: %d", sent)
	}
	payees, sel := []Payee{{Name: "Ana"}, {Name: "Kai"}}, -1
	query := ""
	pt := ui.NewTester(padded(func(c *ui.Context) { Payees(c, payees, &query, &sel, PayeesOptions{}) }), 500, 400)
	settles(t, pt)
	swapUnderPress(pt, rectOf(t, pt, "Kai"), func() { payees[0], payees[1] = payees[1], payees[0] })
	if sel != -1 {
		t.Errorf("payees chose the row that slid under the pointer: %d", sel)
	}
	cats, pick := []string{"Dining", "Travel"}, -1
	ct := ui.NewTester(padded(func(c *ui.Context) { Categorize(c, testTransaction(), cats, &pick, "USD", CategorizeOptions{}) }), 500, 300)
	settles(t, ct)
	swapUnderPress(ct, rectOf(t, ct, "Travel"), func() { cats[0], cats[1] = cats[1], cats[0] })
	if pick != -1 {
		t.Errorf("categorize chose the chip that slid under the pointer: %d", pick)
	}
}

func TestPickersRejectDuplicates(t *testing.T) {
	sent, sel, q := -1, -1, ""
	panics(t, func() {
		ui.Render(func(c *ui.Context) {
			Wallet(c, []Coin{{Symbol: "BTC"}, {Symbol: "BTC"}}, "USD", &sent, WalletOptions{})
		}, 500, 300, 1)
	})
	panics(t, func() {
		ui.Render(func(c *ui.Context) { Payees(c, []Payee{{Name: "Ana"}, {Name: "Ana"}}, &q, &sel, PayeesOptions{}) }, 500, 300, 1)
	})
	panics(t, func() {
		ui.Render(func(c *ui.Context) {
			Categorize(c, testTransaction(), []string{"A1", "A1"}, &sel, "USD", CategorizeOptions{})
		}, 500, 300, 1)
	})
}

func TestPayeesFilterAndEmpty(t *testing.T) {
	query, sel := "", -1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Payees(c, []Payee{{Name: "Ana Martínez"}, {Name: "Kai Nakamura"}}, &query, &sel, PayeesOptions{})
	}), 500, 400)
	settles(t, tt)
	query = "kai"
	tt.Frame()
	settles(t, tt)
	if tt.HasText("Ana Martínez") || !tt.HasText("Kai Nakamura") {
		t.Errorf("filter %q", tt.Texts())
	}
	query = "zz"
	tt.Frame()
	settles(t, tt)
	if !tt.HasText("No payees match") {
		t.Errorf("empty %q", tt.Texts())
	}
}

func TestCategorizeShowsPick(t *testing.T) {
	pick := -1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Categorize(c, testTransaction(), []string{"Dining", "Travel"}, &pick, "USD", CategorizeOptions{})
	}), 500, 300)
	settles(t, tt)
	if count(tt, "Books") != 1 {
		t.Fatalf("header %q", tt.Texts())
	}
	click(t, tt, "Travel")
	if pick != 1 || count(tt, "Books") != 0 {
		t.Errorf("pick %d texts %q", pick, tt.Texts())
	}
}

func TestAllocationHoverShowsClass(t *testing.T) {
	parts := []AllocationPart{{Name: "Stocks", Share: 0.7}, {Name: "Bonds", Share: 0.3}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Allocation(c, "Mix", parts) }), 500, 400)
	settles(t, tt)
	if count(tt, "Bonds") != 1 || count(tt, "Stocks") != 2 {
		t.Fatalf("center starts on the largest class: %q", tt.Texts())
	}
	hover(t, tt, "Bonds")
	if count(tt, "Bonds") != 2 || count(tt, "30%") != 2 {
		t.Errorf("hover did not reach the center: %q", tt.Texts())
	}
}

func TestCashFlowHoverShowsMonth(t *testing.T) {
	months := []FlowMonth{{Label: "Jun", In: 100, Out: 40}, {Label: "Jul", In: 300, Out: 90}}
	tt := ui.NewTester(padded(func(c *ui.Context) { CashFlow(c, months, "USD", CashFlowOptions{}) }), 500, 400)
	settles(t, tt)
	if !tt.HasText("USD 400.00") || !tt.HasText("Net +USD 270.00") {
		t.Fatalf("totals %q", tt.Texts())
	}
	hover(t, tt, "Jul")
	if !tt.HasText("USD 300.00") || !tt.HasText("USD 90.00") {
		t.Errorf("hovered month %q", tt.Texts())
	}
}

func TestTransferStates(t *testing.T) {
	amount, done := "", false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Transfer(c, "Checking", "Savings", &amount, "USD", &done, TransferOptions{})
	}), 500, 400)
	settles(t, tt)
	if tt.HasText("Enter a positive amount") {
		t.Error("an empty field already shows the error")
	}
	amount = "abc"
	tt.Frame()
	settles(t, tt)
	click(t, tt, "Confirm")
	if !tt.HasText("Enter a positive amount") || done {
		t.Errorf("invalid amount %q", tt.Texts())
	}
	amount = "250"
	tt.Frame()
	settles(t, tt)
	before := baseline()
	click(t, tt, "Confirm")
	settles(t, tt)
	if !done || !tt.HasText("USD 250.00 sent to Savings") || !tt.HasText("Sent") {
		t.Errorf("confirmed %v texts %q", done, tt.Texts())
	}
	click(t, tt, "Confirm")
	cues(t, before, sound.Send)
}

func TestCardControlsSwitchMoves(t *testing.T) {
	frozen, online, limit := false, true, 2500.0
	tt := ui.NewTester(padded(func(c *ui.Context) { CardControls(c, &frozen, &online, &limit, "USD", CardControlsOptions{}) }), 500, 400)
	settles(t, tt)
	if !tt.HasText("USD 2,500.00") || tt.HasText("Frozen") {
		t.Fatalf("texts %q", tt.Texts())
	}
	click(t, tt, "Freeze card")
	if !tt.FramePending() {
		t.Error("the switch does not slide")
	}
	settles(t, tt)
	if !frozen || !tt.HasText("Frozen") {
		t.Errorf("frozen %v texts %q", frozen, tt.Texts())
	}
}

func TestSavingsGoalStates(t *testing.T) {
	saved := 400.0
	tt := ui.NewTester(padded(func(c *ui.Context) { SavingsGoal(c, "Trip", saved, 1000, "Apr 1", "USD", SavingsGoalOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("USD 600.00 to go") || !tt.HasText("40%") || tt.HasText("Goal reached") {
		t.Fatalf("texts %q", tt.Texts())
	}
	saved = 1000
	tt.Frame()
	settles(t, tt)
	if !tt.HasText("Goal reached") || !tt.HasText("100%") || tt.HasText("to go") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestReceiptsEmptyAndItems(t *testing.T) {
	receipts := []Receipt{{Merchant: "Shell", Date: "Oct 19", Items: 1, Amount: 48.3}, {Merchant: "Uber", Date: "Oct 12", Items: 3, Amount: 23.9}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Receipts(c, receipts, "USD", ReceiptsOptions{}) }), 500, 400)
	settles(t, tt)
	if !tt.HasText("Oct 19 · 1 item") || !tt.HasText("Oct 12 · 3 items") || !tt.HasText("USD 72.20") {
		t.Fatalf("texts %q", tt.Texts())
	}
	receipts = nil
	tt.Frame()
	settles(t, tt)
	if !tt.HasText("No receipts scanned") || !tt.HasText("USD 0.00") {
		t.Errorf("empty %q", tt.Texts())
	}
}

func TestTaxSummaryZeroIncome(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { TaxSummary(c, 0, 0, 0.2, "USD", TaxSummaryOptions{}) }), 500, 400)
	settles(t, tt)
	if tt.HasText("effective") || tt.HasText("Take-home") {
		t.Errorf("zero income shows a rate: %q", tt.Texts())
	}
	tt = ui.NewTester(padded(func(c *ui.Context) { TaxSummary(c, 1000, 200, 0.25, "USD", TaxSummaryOptions{}) }), 500, 400)
	settles(t, tt)
	if !tt.HasText("USD 800.00") || !tt.HasText("USD 200.00") || !tt.HasText("20.0% effective") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestLoanAndMortgageCopy(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { LoanCalc(c, 1000, 0.12, 1, "USD", LoanCalcOptions{}) }), 500, 400)
	settles(t, tt)
	if !tt.HasText("Term") || count(tt, "1 year") != 1 || tt.HasText("1 years") || !tt.HasText("12.00%") || !tt.HasText("USD 66.19") {
		t.Errorf("loan %q", tt.Texts())
	}
	panics(t, func() {
		ui.Render(func(c *ui.Context) { Mortgage(c, 1000, 0.05, 51, "USD", MortgageOptions{}) }, 500, 400, 1)
	})
}

func TestPnLSignAndChange(t *testing.T) {
	value := 125.0
	tt := ui.NewTester(padded(func(c *ui.Context) { PnL(c, value, "USD", PnLOptions{Basis: 5000}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("+") || !tt.HasText("+2.50%") {
		t.Fatalf("gain %q", tt.Texts())
	}
	value = -50
	tt.Frame()
	settles(t, tt)
	if tt.HasText("+") || !tt.HasText("-1.00%") || !tt.HasText("USD -50.00") {
		t.Errorf("loss %q", tt.Texts())
	}
	panics(t, func() {
		ui.Render(func(c *ui.Context) { PnL(c, math.MaxFloat64, "USD", PnLOptions{Basis: 1e-300}) }, 500, 300, 1)
	})
}

func TestExchangeRateInverseAndPair(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { ExchangeRate(c, "USD/EUR", 0.8, 0, ExchangeRateOptions{Points: []float64{1, 2}}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("1 EUR = 1.2500 USD") || !tt.HasText("+0.00%") {
		t.Errorf("texts %q", tt.Texts())
	}
	for _, pair := range []string{"USDEUR", "USD/", "/EUR"} {
		panics(t, func() {
			ui.Render(func(c *ui.Context) { ExchangeRate(c, pair, 1, 0, ExchangeRateOptions{}) }, 500, 300, 1)
		})
	}
	panics(t, func() {
		ui.Render(func(c *ui.Context) { ExchangeRate(c, "USD/EUR", 1, 0, ExchangeRateOptions{Points: []float64{1}}) }, 500, 300, 1)
	})
	panics(t, func() {
		ui.Render(func(c *ui.Context) { ExchangeRate(c, "USD/EUR", 1e-310, 0, ExchangeRateOptions{}) }, 500, 300, 1)
	})
}

func TestCreditScoreBands(t *testing.T) {
	for score, band := range map[int]string{300: "Poor", 549: "Poor", 550: "Fair", 649: "Fair", 650: "Good", 749: "Good", 750: "Excellent", 850: "Excellent"} {
		tt := ui.NewTester(padded(func(c *ui.Context) { CreditScore(c, score, CreditScoreOptions{}) }), 500, 400)
		tt.Frame()
		if !tt.HasText(band) {
			t.Errorf("%d: %q", score, tt.Texts())
		}
	}
}

func TestMarketStatusCountdown(t *testing.T) {
	var next time.Time
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if next.IsZero() {
			next = c.Now().Add(2*time.Hour + 14*time.Minute)
		}
		MarketStatus(c, true, next, MarketStatusOptions{})
	}), 500, 300)
	tt.Frame()
	if !tt.HasText("02:14:00") || !tt.HasText("Closes in") {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt.Advance(61500 * time.Millisecond)
	if !tt.HasText("02:12:59") {
		t.Errorf("after a minute %q", tt.Texts())
	}
}

func TestRiskMeterNames(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { RiskMeter(c, 5, RiskMeterOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("Very high") || !tt.HasText("5 of 5") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestSplitBillOverflowFaces(t *testing.T) {
	people := 6
	tt := ui.NewTester(padded(func(c *ui.Context) { SplitBill(c, 70, &people, "USD", SplitBillOptions{}) }), 500, 300)
	settles(t, tt)
	if tt.HasText("+1") {
		t.Fatalf("six people need no overflow: %q", tt.Texts())
	}
	click(t, tt, "Increase")
	settles(t, tt)
	if !tt.HasText("+1") || !tt.HasText("USD 10.00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestStatementTotals(t *testing.T) {
	rows := []Transaction{{Merchant: "Payroll", Category: "Salary", Amount: 100}, {Merchant: "Rent", Category: "Home", Amount: -40},
		{Merchant: "Cafe", Category: "Dining", Amount: -5}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Statement(c, "May", rows, "USD", StatementOptions{}) }), 500, 600)
	settles(t, tt)
	for _, want := range []string{"USD 100.00", "USD 45.00", "USD 55.00", "Money in", "Money out"} {
		if !tt.HasText(want) {
			t.Errorf("missing %q in %q", want, strings.Join(tt.Texts(), " | "))
		}
	}
}
