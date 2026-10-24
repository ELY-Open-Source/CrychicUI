package finance

import (
	"bytes"
	"fmt"
	"log/slog"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestWalletSendKeepsFocus(t *testing.T) {
	coins, sent := []Coin{{Symbol: "BTC", Value: 1}}, -1
	tt := ui.NewTester(padded(func(c *ui.Context) { Wallet(c, coins, "USD", &sent, WalletOptions{}) }), 500, 300)
	settles(t, tt)
	tt.Key(0, ui.KeyTab)
	tt.Frame()
	if !tt.Focused("Send BTC") {
		t.Fatal("Tab did not reach Send BTC")
	}
	before := baseline()
	tt.Key(0, ui.KeySpace)
	tt.Frame()
	if sent != 0 || !tt.Focused("Send BTC") {
		t.Fatalf("sent %d, focus kept %v", sent, tt.Focused("Send BTC"))
	}
	tt.Key(0, ui.KeySpace)
	tt.Frame()
	cues(t, before, sound.Send, sound.Send)
}

func TestPnLRollsAcrossZero(t *testing.T) {
	value := -50.0
	tt := ui.NewTester(padded(func(c *ui.Context) { PnL(c, value, "USD", PnLOptions{}) }), 500, 300)
	settles(t, tt)
	value = 125
	tt.Frame()
	tt.Advance(60 * time.Millisecond)
	mid := tt.Image()
	settles(t, tt)
	end := tt.Image()
	r := rectOf(t, tt, "USD 125.00")
	bg, moved := end.RGBAAt(2, 2), 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if (mid.RGBAAt(x, y) == bg) != (end.RGBAAt(x, y) == bg) {
				moved++
			}
		}
	}
	if moved == 0 {
		t.Error("the number appeared at once instead of rolling")
	}
}

func TestAmortizeRepaysPrincipal(t *testing.T) {
	for _, tc := range []struct {
		principal, rate float64
		years           int
	}{{100000, 1, 50}, {420000, 0.0625, 30}, {1000, 0, 2}, {5000, 1e6, 1}} {
		plan := amortize(loanPayment(tc.principal, tc.rate, tc.years), tc.rate, tc.years)
		paid := 0.0
		for _, y := range plan {
			paid += y.principal
		}
		if math.Abs(paid-tc.principal) > tc.principal*1e-9 || plan[len(plan)-1].principal <= 0 {
			t.Errorf("%v: repaid %v of %v, last year %v", tc, paid, tc.principal, plan[len(plan)-1].principal)
		}
	}
}

func TestInterestStableAtTinyRates(t *testing.T) {
	plan := amortize(loanPayment(1e14, 1e-15, 1), 1e-15, 1)
	if got, want := plan[0].interest, 1e14*(1e-15/12)*13/2; math.Abs(got-want) > want*1e-6 {
		t.Errorf("interest %v, want %v", got, want)
	}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		LoanCalc(c, 1e14, 1e-15, 1, "USD", LoanCalcOptions{})
		Mortgage(c, 1e14, 1e-15, 1, "USD", MortgageOptions{})
	}), 500, 900)
	settles(t, tt)
	if count(tt, "USD 0.05") != 2 || tt.HasText("USD 0.06") || count(tt, "Year 1") != 1 {
		t.Errorf("shown interest %q", tt.Texts())
	}
	panics(t, func() {
		ui.Render(func(c *ui.Context) { LoanCalc(c, 1000, 0.05, 51, "USD", LoanCalcOptions{}) }, 500, 300, 1)
	})
}

func TestMortgageYearOverflow(t *testing.T) {
	panics(t, func() {
		ui.Render(func(c *ui.Context) { Mortgage(c, math.MaxFloat64, 0, 1, "USD", MortgageOptions{}) }, 500, 600, 1)
	})
}

func TestCategorizeLongChipFits(t *testing.T) {
	long := "Professional development and conference travel reimbursements"
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ui.Column(c).Width(320).Children(func() {
			Categorize(c, testTransaction(), []string{"Dining", long}, new(int), "USD", CategorizeOptions{})
		})
	}), 500, 400)
	settles(t, tt)
	chip, card := rectOf(t, tt, long), rectOf(t, tt, "Bookstore")
	if chip.X+chip.W > 10+320+0.5 || chip.W <= 0 {
		t.Errorf("chip %v runs past the 320 column (row %v)", chip, card)
	}
}

func TestMarketStatusFarCountdown(t *testing.T) {
	next := time.Date(2500, 1, 1, 0, 0, 0, 0, time.UTC)
	var want string
	tt := ui.NewTester(padded(func(c *ui.Context) {
		s := next.Unix() - c.Now().Unix()
		if c.Now().Nanosecond() < next.Nanosecond() {
			s++
		}
		want = fmt.Sprintf("%02d:%02d:%02d", s/3600, s/60%60, s%60)
		MarketStatus(c, false, next, MarketStatusOptions{})
	}), 500, 300)
	tt.Frame()
	if !tt.HasText(want) || tt.HasText("2562047:") {
		t.Errorf("want %q in %q", want, tt.Texts())
	}
}

func TestLoanCalcRateOverflow(t *testing.T) {
	panics(t, func() {
		ui.Render(func(c *ui.Context) { LoanCalc(c, 1e-306, 1e307, 1, "USD", LoanCalcOptions{}) }, 500, 300, 1)
	})
}

func TestExchangeRateOneSlash(t *testing.T) {
	panics(t, func() {
		ui.Render(func(c *ui.Context) { ExchangeRate(c, "USD/EUR/JPY", 1, 0, ExchangeRateOptions{}) }, 500, 300, 1)
	})
}

func TestEditsLogCounts(t *testing.T) {
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	t.Cleanup(func() { slog.SetDefault(old) })
	amount, done, query, sel := "", false, "", -1
	tt := ui.NewTester(padded(func(c *ui.Context) { Transfer(c, "Checking", "Savings", &amount, "USD", &done, TransferOptions{}) }), 500, 500)
	settles(t, tt)
	click(t, tt, "Amount")
	tt.Type("42")
	tt.Frame()
	pt := ui.NewTester(padded(func(c *ui.Context) { Payees(c, []Payee{{Name: "Ana"}}, &query, &sel, PayeesOptions{}) }), 500, 500)
	settles(t, pt)
	pt.Key(0, ui.KeyTab)
	pt.Type("an")
	pt.Frame()
	s := buf.String()
	for _, want := range []string{`msg="crychic: Transfer edited" runes=2`, `msg="crychic: Payees query" runes=2`} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %s in %s", want, s)
		}
	}
}
