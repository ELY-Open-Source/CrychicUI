package gallery

import (
	"fmt"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/finance"
	"github.com/egoist/mygo/ui"
)

var financeEntries = []Entry{
	{ID: 1016, Name: "Balance", Category: "finance", Summary: "Balance with a hide toggle; digits roll; cue Toggle.", Demo: demoBalance},
	{ID: 1017, Name: "AccountCard", Category: "finance", Summary: "Bank card face with masked number.", Demo: demoAccountCard},
	{ID: 1018, Name: "Transactions", Category: "finance", Summary: "Rows with merchant, category, amount sign.", Demo: demoTransactions},
	{ID: 1019, Name: "TransactionRow", Category: "finance", Summary: "One row with icon and amount; fades in.", Demo: demoTransactionRow},
	{ID: 1020, Name: "Spending", Category: "finance", Summary: "Categories with bars of spend.", Demo: demoSpending},
	{ID: 1021, Name: "Budget", Category: "finance", Summary: "Category budgets with used bars.", Demo: demoBudget},
	{ID: 1022, Name: "Holdings", Category: "finance", Summary: "Positions with value and change.", Demo: demoHoldings},
	{ID: 1023, Name: "Ticker", Category: "finance", Summary: "Symbol, price, change; flashes on change.", Demo: demoTicker},
	{ID: 1024, Name: "TickerTape", Category: "finance", Summary: "Scrolling row of tickers.", Demo: demoTickerTape},
	{ID: 1025, Name: "Watchlist", Category: "finance", Summary: "Symbols with sparkline and change; cue Select.", Demo: demoWatchlist},
	{ID: 1026, Name: "OrderTicket", Category: "finance", Summary: "Buy and sell with quantity and price; cue Success.", Demo: demoOrderTicket},
	{ID: 1027, Name: "OrderBook", Category: "finance", Summary: "Bids and asks with depth bars; rows flash.", Demo: demoOrderBook},
	{ID: 1028, Name: "Depth", Category: "finance", Summary: "Cumulative bid and ask areas.", Demo: demoDepth},
	{ID: 1029, Name: "PriceAlert", Category: "finance", Summary: "Above or below a price; cue Notify.", Demo: demoPriceAlert},
	{ID: 1030, Name: "Converter", Category: "finance", Summary: "Amount between two currencies at a rate; cue Tick.", Demo: demoConverter},
	{ID: 1031, Name: "Invoice", Category: "finance", Summary: "Lines, tax, total and status.", Demo: demoInvoice},
	{ID: 1032, Name: "PaymentRequest", Category: "finance", Summary: "Amount, payer and pay; cue Success.", Demo: demoPaymentRequest},
	{ID: 1033, Name: "Bills", Category: "finance", Summary: "Due bills with pay; cue Success.", Demo: demoBills},
	{ID: 1034, Name: "Subscriptions", Category: "finance", Summary: "Recurring charges with next date.", Demo: demoSubscriptions},
	{ID: 1035, Name: "NetWorth", Category: "finance", Summary: "Assets minus debts over time.", Demo: demoNetWorth},
}

func demoFinanceTransaction() finance.Transaction {
	return finance.Transaction{Merchant: "Whole Foods Market", Category: "Groceries", Amount: -84.32, When: "Today", Icon: "shopping-cart"}
}

func demoBalance(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	hidden := boolState(root, "hidden", false)
	root.Children(func() {
		o := finance.BalanceOptions{Change: "+USD 1,240.00 since last month"}
		if finance.Balance(c, 12480.55, "USD", hidden, o).Changed() {
			log("hidden " + fmt.Sprint(*hidden))
		}
	})
}

func demoAccountCard(c *ui.Context, log func(string)) {
	ui.Column(c).Width(380).Children(func() {
		finance.AccountCard(c, "Northwind Bank", "4821",
			finance.AccountCardOptions{Mask: "•••• •••• •••• %s", Holder: "MAYA CHEN", Expires: "09/29", Network: "VISA"})
	})
}

func demoTransactionList() []finance.Transaction {
	return []finance.Transaction{
		{Merchant: "Blue Bottle Coffee", Category: "Food & drink", Amount: -6.5, When: "Today", Icon: "coffee"},
		demoFinanceTransaction(),
		{Merchant: "Acme Robotics payroll", Category: "Salary", Amount: 4200, When: "Yesterday", Icon: "briefcase"},
		{Merchant: "Lyft", Category: "Transport", Amount: -18.4, When: "Yesterday", Icon: "car"},
		{Merchant: "Netflix", Category: "Entertainment", Amount: -15.49, When: "Oct 22", Icon: "tv"},
	}
}

func demoTransactions(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() {
		finance.Transactions(c, demoTransactionList(), "USD", finance.TransactionsOptions{Title: "Recent activity"})
	})
}

func demoTransactionRow(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() { finance.TransactionRow(c, demoFinanceTransaction(), "USD") })
}

func demoSpending(c *ui.Context, log func(string)) {
	rows := []finance.Spend{{Category: "Housing", Amount: 1850}, {Category: "Groceries", Amount: 520},
		{Category: "Dining out", Amount: 340}, {Category: "Transport", Amount: 210}, {Category: "Entertainment", Amount: 120}}
	ui.Column(c).Width(400).Children(func() { finance.Spending(c, rows, "USD", finance.SpendingOptions{Title: "October spending"}) })
}

func demoBudget(c *ui.Context, log func(string)) {
	rows := []finance.BudgetItem{{Category: "Groceries", Used: 420, Limit: 500}, {Category: "Dining out", Used: 380, Limit: 300},
		{Category: "Shopping", Used: 230, Limit: 250}, {Category: "Transport", Used: 150, Limit: 250}}
	ui.Column(c).Width(400).Children(func() { finance.Budget(c, rows, "USD", finance.BudgetOptions{}) })
}

func demoHoldings(c *ui.Context, log func(string)) {
	positions := []finance.Holding{
		{Symbol: "AAPL", Name: "Apple", Shares: 12, Price: 228.52, Change: 1.24},
		{Symbol: "MSFT", Name: "Microsoft", Shares: 8, Price: 431.1, Change: -0.62},
		{Symbol: "NVDA", Name: "NVIDIA", Shares: 20, Price: 138.07, Change: 3.41},
		{Symbol: "VTI", Name: "Vanguard Total Market", Shares: 30, Price: 284.9, Change: 0.35},
	}
	ui.Column(c).Width(440).Children(func() { finance.Holdings(c, positions, "USD", finance.HoldingsOptions{}) })
}

func demoTicker(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380).Gap(12)
	price := numberState(root, "price", 228.52)
	ticks := intState(root, "ticks", 0)
	root.Children(func() {
		finance.Ticker(c, "AAPL", *price, 100*(*price-225.74)/225.74, "USD", finance.TickerOptions{Name: "Apple Inc."})
		if button.Outline(c, "Next tick").Clicked() {
			*ticks++
			*price += []float64{0.84, -1.22, 0.37}[*ticks%3]
			log(fmt.Sprintf("price %.2f", *price))
		}
	})
}

func demoTickerTape(c *ui.Context, log func(string)) {
	items := []finance.TickerItem{{Symbol: "AAPL", Price: 228.52, Change: 1.24}, {Symbol: "MSFT", Price: 431.1, Change: -0.62},
		{Symbol: "NVDA", Price: 138.07, Change: 3.41}, {Symbol: "AMZN", Price: 186.38, Change: 0.92},
		{Symbol: "TSLA", Price: 218.85, Change: -2.17}, {Symbol: "GOOGL", Price: 165.27, Change: 0.48}}
	ui.Column(c).Width(480).Children(func() { finance.TickerTape(c, items, "USD") })
}

func demoWatchlist(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	selected := intState(root, "selected", -1)
	items := []finance.WatchItem{
		{Symbol: "AAPL", Name: "Apple", Price: 228.52, Change: 1.24, Points: []float64{221, 224, 223, 226, 225, 228.5}},
		{Symbol: "TSLA", Name: "Tesla", Price: 218.85, Change: -2.17, Points: []float64{232, 228, 230, 224, 221, 218.9}},
		{Symbol: "AMZN", Name: "Amazon", Price: 186.38, Change: 0.92, Points: []float64{182, 181, 184, 183, 185, 186.4}},
		{Symbol: "NVDA", Name: "NVIDIA", Price: 138.07, Change: 3.41, Points: []float64{126, 129, 131, 130, 135, 138.1}},
	}
	root.Children(func() {
		if finance.Watchlist(c, items, "USD", selected, finance.WatchlistOptions{}).Changed() {
			log("selected " + fmt.Sprint(*selected))
		}
	})
}

func demoOrderTicket(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	quantity := intState(root, "quantity", 1)
	root.Children(func() {
		root, buy, sell := finance.OrderTicket(c, quantity, 228.52, "USD", finance.OrderTicketOptions{Title: "Trade AAPL"})
		switch {
		case buy.Clicked():
			log("buy " + fmt.Sprint(*quantity))
		case sell.Clicked():
			log("sell " + fmt.Sprint(*quantity))
		case root.Changed():
			log("quantity " + fmt.Sprint(*quantity))
		}
	})
}

func demoBook() (bids, asks []finance.BookLevel) {
	bids = []finance.BookLevel{{Price: 228.5, Amount: 120}, {Price: 228.48, Amount: 340}, {Price: 228.45, Amount: 90},
		{Price: 228.4, Amount: 510}, {Price: 228.35, Amount: 260}}
	asks = []finance.BookLevel{{Price: 228.54, Amount: 80}, {Price: 228.57, Amount: 260}, {Price: 228.6, Amount: 150},
		{Price: 228.65, Amount: 420}, {Price: 228.7, Amount: 310}}
	return bids, asks
}

func demoOrderBook(c *ui.Context, log func(string)) {
	bids, asks := demoBook()
	ui.Column(c).Width(360).Children(func() { finance.OrderBook(c, bids, asks, "USD", finance.OrderBookOptions{Title: "AAPL order book"}) })
}

func demoDepth(c *ui.Context, log func(string)) {
	bids, asks := demoBook()
	ui.Column(c).Width(440).Children(func() { finance.Depth(c, bids, asks, "USD", finance.DepthOptions{Title: "AAPL depth"}) })
}

func demoPriceAlert(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400).Gap(12)
	price := numberState(root, "price", 224.1)
	active := boolState(root, "active", false)
	root.Children(func() {
		if finance.PriceAlert(c, finance.Above, 230, *price, "USD", active, finance.PriceAlertOptions{}).Changed() {
			log("triggered")
		}
		if button.Outline(c, "Simulate a rally").Clicked() {
			*price += 8
			log(fmt.Sprintf("price %.2f", *price))
		}
	})
}

func demoConverter(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	amount := textState(root, "amount", "250")
	root.Children(func() {
		if finance.Converter(c, amount, 0.9214, "USD", "EUR", finance.ConverterOptions{}).Changed() {
			log(fmt.Sprintf("amount %d chars", utf8.RuneCountInString(*amount)))
		}
	})
}

func demoInvoice(c *ui.Context, log func(string)) {
	lines := []finance.Charge{{Label: "Brand identity design", Amount: 2400}, {Label: "Website build, 40 hours", Amount: 4800},
		{Label: "Hosting setup", Amount: 250}}
	o := finance.InvoiceOptions{Tax: "Tax (8%)", To: "Acme Robotics", Date: "Issued Oct 15, 2026"}
	ui.Column(c).Width(420).Children(func() { finance.Invoice(c, "INV-2026-0142", "Due Nov 14", lines, 596, "USD", o) })
}

func demoPaymentRequest(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	paid := boolState(root, "paid", false)
	root.Children(func() {
		o := finance.PaymentRequestOptions{Note: "Dinner at Nopa, split three ways"}
		if root, _ := finance.PaymentRequest(c, 42.5, "Daniel Okafor", "USD", paid, o); root.Changed() {
			log("paid")
		}
	})
}

func demoBills(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	paid := ui.Local(root, "paid", func() []bool { return make([]bool, 4) })
	bills := []finance.Bill{
		{Name: "Internet", Due: "Due Oct 20, overdue", Amount: 65, Icon: "wifi", Overdue: true},
		{Name: "Electricity", Due: "Due Oct 28", Amount: 86.4, Icon: "zap"},
		{Name: "Rent", Due: "Due Nov 1", Amount: 2100, Icon: "house"},
		{Name: "Phone", Due: "Due Nov 4", Amount: 45, Icon: "smartphone"},
	}
	root.Children(func() {
		if finance.Bills(c, bills, paid, "USD", finance.BillsOptions{}).Changed() {
			n := 0
			for _, on := range *paid {
				if on {
					n++
				}
			}
			log(fmt.Sprintf("%d paid", n))
		}
	})
}

func demoSubscriptions(c *ui.Context, log func(string)) {
	rows := []finance.Subscription{
		{Name: "Spotify", Plan: "Premium Family", Next: "Nov 3", Amount: 16.99, Icon: "music"},
		{Name: "Netflix", Plan: "Standard", Next: "Nov 12", Amount: 15.49, Icon: "tv"},
		{Name: "iCloud+", Plan: "200 GB", Next: "Nov 18", Amount: 2.99, Icon: "cloud"},
		{Name: "Equinox", Plan: "Monthly", Next: "Nov 1", Amount: 49, Icon: "dumbbell"},
		{Name: "The New York Times", Plan: "Digital", Next: "Nov 21", Amount: 4, Icon: "newspaper"},
	}
	ui.Column(c).Width(420).Children(func() { finance.Subscriptions(c, rows, "USD", finance.SubscriptionsOptions{}) })
}

func demoNetWorth(c *ui.Context, log func(string)) {
	points := []finance.NetPoint{{Label: "May", Assets: 182000, Debts: 64000}, {Label: "Jun", Assets: 186500, Debts: 62800},
		{Label: "Jul", Assets: 184200, Debts: 61500}, {Label: "Aug", Assets: 191800, Debts: 60300},
		{Label: "Sep", Assets: 196400, Debts: 59100}, {Label: "Oct", Assets: 203900, Debts: 57800}}
	ui.Column(c).Width(440).Children(func() { finance.NetWorth(c, points, "USD", finance.NetWorthOptions{}) })
}
