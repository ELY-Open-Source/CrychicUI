package gallery

import (
	"fmt"

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
	return finance.Transaction{Merchant: "Bookstore", Category: "Books", Amount: -24}
}

func demoBalance(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	hidden := boolState(root, "hidden", false)
	root.Children(func() {
		if finance.Balance(c, 12480, "USD", hidden, finance.BalanceOptions{}).Changed() {
			log("hidden " + fmt.Sprint(*hidden))
		}
	})
}

func demoAccountCard(c *ui.Context, log func(string)) {
	finance.AccountCard(c, "My Bank", "1234", finance.AccountCardOptions{})
}

func demoTransactions(c *ui.Context, log func(string)) {
	finance.Transactions(c, []finance.Transaction{demoFinanceTransaction()}, "USD")
}

func demoTransactionRow(c *ui.Context, log func(string)) {
	finance.TransactionRow(c, demoFinanceTransaction(), "USD")
}

func demoSpending(c *ui.Context, log func(string)) {
	rows := []finance.Spend{{Category: "Food", Amount: 30}, {Category: "Books", Amount: 10}}
	finance.Spending(c, rows, "USD")
}

func demoBudget(c *ui.Context, log func(string)) {
	rows := []finance.BudgetItem{{Category: "Food", Used: 120, Limit: 100}}
	finance.Budget(c, rows, "USD", finance.BudgetOptions{})
}

func demoHoldings(c *ui.Context, log func(string)) {
	positions := []finance.Holding{{Symbol: "MYGO", Shares: 2, Price: 50, Change: 1.5}}
	finance.Holdings(c, positions, "USD")
}

func demoTicker(c *ui.Context, log func(string)) {
	finance.Ticker(c, "MYGO", 100, 1.5, "USD")
}

func demoTickerTape(c *ui.Context, log func(string)) {
	items := []finance.TickerItem{{Symbol: "MYGO", Price: 100, Change: 1.5}}
	finance.TickerTape(c, items, "USD")
}

func demoWatchlist(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := intState(root, "selected", -1)
	items := []finance.WatchItem{{Symbol: "MYGO", Price: 100, Change: 1.5, Points: []float64{1, 2, 1.5}}}
	root.Children(func() {
		if finance.Watchlist(c, items, "USD", selected).Changed() {
			log("selected " + fmt.Sprint(*selected))
		}
	})
}

func demoOrderTicket(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	quantity := intState(root, "quantity", 1)
	root.Children(func() {
		root, buy, sell := finance.OrderTicket(c, quantity, 100, "USD", finance.OrderTicketOptions{})
		if root.Changed() {
			log("quantity " + fmt.Sprint(*quantity))
		}
		if buy.Clicked() {
			log("buy " + fmt.Sprint(*quantity))
		}
		if sell.Clicked() {
			log("sell " + fmt.Sprint(*quantity))
		}
	})
}

func demoOrderBook(c *ui.Context, log func(string)) {
	levels := []finance.BookLevel{{Price: 99, Amount: 2}, {Price: 98, Amount: 4}}
	finance.OrderBook(c, levels, levels, "USD", finance.OrderBookOptions{})
}

func demoDepth(c *ui.Context, log func(string)) {
	levels := []finance.BookLevel{{Price: 99, Amount: 2}, {Price: 98, Amount: 4}}
	finance.Depth(c, levels, levels, "USD", finance.DepthOptions{})
}

func demoPriceAlert(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	active := boolState(root, "active", false)
	root.Children(func() {
		if finance.PriceAlert(c, finance.Above, 100, 120, "USD", active, finance.PriceAlertOptions{}).Changed() {
			log("triggered")
		}
	})
}

func demoConverter(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	amount := textState(root, "amount", "10")
	root.Children(func() {
		if finance.Converter(c, amount, 2, "USD", "EUR", finance.ConverterOptions{}).Changed() {
			log("edited")
		}
	})
}

func demoInvoice(c *ui.Context, log func(string)) {
	lines := []finance.Charge{{Label: "Service", Amount: 100}}
	finance.Invoice(c, "INV-1", "Due", lines, 10, "USD", finance.InvoiceOptions{})
}

func demoPaymentRequest(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	paid := boolState(root, "paid", false)
	root.Children(func() {
		if finance.PaymentRequest(c, 100, "Ana", "USD", paid, finance.PaymentRequestOptions{}); *paid {
			log("paid")
			*paid = false
		}
	})
}

func demoBills(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	paid := intState(root, "paid", -1)
	bills := []finance.Bill{{Name: "Electricity", Due: "Friday", Amount: 80}}
	root.Children(func() {
		if finance.Bills(c, bills, paid, "USD", finance.BillsOptions{}).Changed() {
			log("paid " + fmt.Sprint(*paid))
		}
	})
}

func demoSubscriptions(c *ui.Context, log func(string)) {
	rows := []finance.Subscription{{Name: "Music", Next: "Nov 1", Amount: 10}}
	finance.Subscriptions(c, rows, "USD")
}

func demoNetWorth(c *ui.Context, log func(string)) {
	points := []finance.NetPoint{{Label: "Jan", Assets: 100, Debts: 20}, {Label: "Feb", Assets: 150, Debts: 50}}
	finance.NetWorth(c, points, "USD", finance.NetWorthOptions{})
}
