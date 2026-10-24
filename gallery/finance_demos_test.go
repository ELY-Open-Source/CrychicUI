package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestFinanceDemosLog(t *testing.T) {
	press := func(name string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error { return tt.Click(name) }
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"Balance":     {press("Hide balance"), "hidden true"},
		"Watchlist":   {press("TSLA"), "selected 1"},
		"Ticker":      {press("Next tick"), "price 227.30"},
		"PriceAlert":  {press("Simulate a rally"), "triggered"},
		"OrderTicket": {press("Increase"), "quantity 2"},
		"Converter": {func(tt *ui.Tester) error {
			if err := tt.Click("USD"); err != nil {
				return err
			}
			tt.Type("2")
			return nil
		}, "amount "},
		"PaymentRequest": {press("Pay"), "paid"},
		"Bills":          {press("Pay Rent"), "1 paid"},
	}
	quiet := map[string]bool{"AccountCard": true, "Transactions": true, "TransactionRow": true, "Spending": true,
		"Budget": true, "Holdings": true, "TickerTape": true, "OrderBook": true, "Depth": true,
		"Invoice": true, "Subscriptions": true, "NetWorth": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "finance" || e.ID >= 1036 {
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
		t.Errorf("checked %d finance demos", seen)
	}
}
