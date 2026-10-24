package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestFinance2DemosLog(t *testing.T) {
	press := func(name string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error { return tt.Click(name) }
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"Wallet":    {press("Send BTC"), "sent 0"},
		"SplitBill": {press("Increase"), "people 5"},
		"Transfer": {func(tt *ui.Tester) error {
			if err := tt.Click("Amount"); err != nil {
				return err
			}
			tt.Type("2")
			return nil
		}, "7 runes"},
		"CardControls": {press("Freeze card"), "frozen true"},
		"Payees":       {press("Kai Nakamura"), "selected 1"},
		"RiskMeter":    {press("Raise risk"), "level 4"},
		"PnL":          {press("Next day"), "pnl -412.35"},
		"Categorize":   {press("Dining"), "selected 1"},
	}
	quiet := map[string]bool{"Allocation": true, "LoanCalc": true, "Mortgage": true, "TaxSummary": true,
		"SavingsGoal": true, "Statement": true, "CashFlow": true, "Dividends": true,
		"ExchangeRate": true, "Receipts": true, "CreditScore": true, "MarketStatus": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "finance" || e.ID < 1036 {
			continue
		}
		seen++
		if quiet[e.Name] {
			continue
		}
		act, ok := acts[e.Name]
		if !ok {
			t.Errorf("%s has no logged action", e.Name)
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
		tt.Advance(time.Second)
		if err := act.do(tt); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Frame()
		if !strings.Contains(strings.Join(logs, "\n"), act.want) {
			t.Errorf("%s logs %q, want %q", e.Name, logs, act.want)
		}
	}
	if seen != 20 {
		t.Errorf("checked %d second finance demos", seen)
	}
}
