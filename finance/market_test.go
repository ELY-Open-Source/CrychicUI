package finance

import (
	"fmt"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestBalance(t *testing.T) {
	hidden := false
	tt := ui.NewTester(padded(func(c *ui.Context) { Balance(c, 100, "USD", &hidden, BalanceOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Hide balance")
	if !hidden || !tt.HasText("••••") {
		t.Errorf("hidden %v texts %q", hidden, tt.Texts())
	}
	cues(t, before, sound.Toggle)
}

func TestBalanceOverdraft(t *testing.T) {
	hidden := false
	tt := ui.NewTester(padded(func(c *ui.Context) { Balance(c, -25, "USD", &hidden, BalanceOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("USD -25.00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestAccountCard(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { AccountCard(c, "My Bank", "1234", AccountCardOptions{}) }), 400, 220)
	settles(t, tt)
	if !tt.HasText("•••• 1234") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestTransactions(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Transactions(c, []Transaction{testTransaction()}, "USD", TransactionsOptions{})
	}), 500, 180)
	settles(t, tt)
	if !tt.HasText("Bookstore") || !tt.HasText("USD -24.00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestTransactionRow(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { TransactionRow(c, testTransaction(), "USD") }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Books") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestSpending(t *testing.T) {
	rows := []Spend{{Category: "Food", Amount: 30}, {Category: "Books", Amount: 10}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Spending(c, rows, "USD", SpendingOptions{}) }), 500, 250)
	settles(t, tt)
	if !tt.HasText("Food") || !tt.HasText("USD 30.00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestSpendingZero(t *testing.T) {
	rows := []Spend{{Category: "Food", Amount: 0}, {Category: "Books", Amount: 0}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Spending(c, rows, "USD", SpendingOptions{}) }), 500, 250)
	settles(t, tt)
	if !tt.HasText("USD 0.00") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestDerivedOverflow(t *testing.T) {
	panics(t, func() {
		ui.Render(func(c *ui.Context) {
			Spending(c, []Spend{{Category: "A", Amount: 1e308}, {Category: "B", Amount: 1e308}}, "USD", SpendingOptions{})
		}, 500, 300, 1)
	})
	panics(t, func() {
		ui.Render(func(c *ui.Context) {
			Depth(c, []BookLevel{{Price: 1, Amount: 1e308}, {Price: 2, Amount: 1e308}}, testLevels(), "USD", DepthOptions{})
		}, 500, 300, 1)
	})
	panics(t, func() {
		ui.Render(func(c *ui.Context) {
			NetWorth(c, []NetPoint{{Label: "A", Assets: 0, Debts: 1e308}, {Label: "B", Assets: 1e308}}, "USD", NetWorthOptions{})
		}, 500, 300, 1)
	})
}

func TestBudget(t *testing.T) {
	rows := []BudgetItem{{Category: "Food", Used: 120, Limit: 100}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Budget(c, rows, "USD", BudgetOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("USD 20.00 over") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestHoldings(t *testing.T) {
	positions := []Holding{{Symbol: "MYGO", Shares: 2, Price: 50, Change: 1.5}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Holdings(c, positions, "USD", HoldingsOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("MYGO") || !tt.HasText("USD 100.00") || !tt.HasText("+1.50%") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestTicker(t *testing.T) {
	price := 100.0
	tt := ui.NewTester(padded(func(c *ui.Context) { Ticker(c, "MYGO", price, 1.5, "USD", TickerOptions{}) }), 500, 180)
	settles(t, tt)
	price = 110
	tt.Frame()
	if !tt.FramePending() {
		t.Error("price change did not flash")
	}
}

func TestTickerTape(t *testing.T) {
	items := []TickerItem{{Symbol: "MYGO", Price: 100, Change: 1.5}}
	tt := ui.NewTester(padded(func(c *ui.Context) { TickerTape(c, items, "USD") }), 500, 180)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("ticker tape stopped")
	}
}

func TestTickerTapeLongList(t *testing.T) {
	items := make([]TickerItem, 10)
	for i := range items {
		items[i] = TickerItem{Symbol: fmt.Sprintf("SYM%d", i), Price: 100, Change: 1}
	}
	tt := ui.NewTester(padded(func(c *ui.Context) { TickerTape(c, items, "USD") }), 500, 180)
	tt.Advance(11400 * time.Millisecond)
	r, ok := tt.Find("SYM9")
	if !ok || r.X < 0 || r.X > 500 {
		t.Errorf("last symbol never arrived: %v", r)
	}
}

func TestTickerTapeWrap(t *testing.T) {
	items := make([]TickerItem, 10)
	for i := range items {
		items[i] = TickerItem{Symbol: fmt.Sprintf("SYM%d", i), Price: 100, Change: 1}
	}
	tt := ui.NewTester(padded(func(c *ui.Context) { TickerTape(c, items, "USD") }), 500, 180)
	tt.Advance(11990 * time.Millisecond)
	before, _ := tt.Find("SYM1")
	tt.Advance(20 * time.Millisecond)
	after, _ := tt.Find("SYM1")
	t.Logf("wrap before %v after %v", before.X, after.X)
	if delta := after.X - before.X; delta < -1 || delta > 1 {
		t.Errorf("cycle jumps by %v", delta)
	}
}

func TestWatchlist(t *testing.T) {
	selected := -1
	items := []WatchItem{{Symbol: "MYGO", Price: 100, Change: 1.5, Points: []float64{1, 2, 1.5}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Watchlist(c, items, "USD", &selected, WatchlistOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "MYGO")
	if selected != 0 {
		t.Errorf("selected %d", selected)
	}
	cues(t, before, sound.Select)
}
