package gallery

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/finance"
	"github.com/egoist/mygo/ui"
)

var financeEntries2 = []Entry{
	{ID: 1036, Name: "Allocation", Category: "finance", Summary: "Shares by asset class as a donut.", Demo: demoAllocation},
	{ID: 1037, Name: "Wallet", Category: "finance", Summary: "Coins with balances and send; cue Send.", Demo: demoWallet},
	{ID: 1038, Name: "LoanCalc", Category: "finance", Summary: "Payment from amount, rate, term; rolls.", Demo: demoLoanCalc},
	{ID: 1039, Name: "Mortgage", Category: "finance", Summary: "Principal and interest over years.", Demo: demoMortgage},
	{ID: 1040, Name: "TaxSummary", Category: "finance", Summary: "Income, deductions, tax due.", Demo: demoTaxSummary},
	{ID: 1041, Name: "SplitBill", Category: "finance", Summary: "Total split among people; cue Tick.", Demo: demoSplitBill},
	{ID: 1042, Name: "SavingsGoal", Category: "finance", Summary: "Goal with saved bar and target date.", Demo: demoSavingsGoal},
	{ID: 1043, Name: "Transfer", Category: "finance", Summary: "From, to, amount, confirm; cue Send.", Demo: demoFinanceTransfer},
	{ID: 1044, Name: "CardControls", Category: "finance", Summary: "Freeze, limits, online toggles; cue Toggle.", Demo: demoCardControls},
	{ID: 1045, Name: "Statement", Category: "finance", Summary: "Month statement rows and totals.", Demo: demoStatement},
	{ID: 1046, Name: "CashFlow", Category: "finance", Summary: "In and out bars by month.", Demo: demoCashFlow},
	{ID: 1047, Name: "Dividends", Category: "finance", Summary: "Payments by date.", Demo: demoDividends},
	{ID: 1048, Name: "RiskMeter", Category: "finance", Summary: "Risk 1–5 on a scale.", Demo: demoRiskMeter},
	{ID: 1049, Name: "ExchangeRate", Category: "finance", Summary: "Pair with rate and change.", Demo: demoExchangeRate},
	{ID: 1050, Name: "Payees", Category: "finance", Summary: "Saved payees with search; cue Select.", Demo: demoPayees},
	{ID: 1051, Name: "Categorize", Category: "finance", Summary: "Category chips for a transaction; cue Select.", Demo: demoCategorize},
	{ID: 1052, Name: "Receipts", Category: "finance", Summary: "Scanned receipt list with totals.", Demo: demoFinanceReceipts},
	{ID: 1053, Name: "CreditScore", Category: "finance", Summary: "Score on an arc with band.", Demo: demoCreditScore},
	{ID: 1054, Name: "MarketStatus", Category: "finance", Summary: "Open or closed with time to next.", Demo: demoMarketStatus},
	{ID: 1055, Name: "PnL", Category: "finance", Summary: "Profit and loss with sign color; rolls.", Demo: demoPnL},
}

func demoAllocation(c *ui.Context, log func(string)) {
	parts := []finance.AllocationPart{{Name: "US stocks", Share: 0.45}, {Name: "International", Share: 0.2},
		{Name: "Bonds", Share: 0.2}, {Name: "Real estate", Share: 0.1}, {Name: "Cash", Share: 0.05}}
	ui.Column(c).Width(440).Children(func() { finance.Allocation(c, "Portfolio mix", parts) })
}

func demoWallet(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	sent := intState(root, "sent", -1)
	coins := []finance.Coin{{Symbol: "BTC", Name: "Bitcoin", Amount: 0.1842, Value: 12480.55},
		{Symbol: "ETH", Name: "Ethereum", Amount: 2.35, Value: 6122.40}, {Symbol: "SOL", Name: "Solana", Amount: 18.4, Value: 2719.52},
		{Symbol: "USDC", Name: "USD Coin", Amount: 850, Value: 850}}
	root.Children(func() {
		if finance.Wallet(c, coins, "USD", sent, finance.WalletOptions{}).Changed() {
			log("sent " + fmt.Sprint(*sent))
		}
	})
}

func demoLoanCalc(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	amount := numberState(root, "amount", 24000)
	root.Children(func() {
		finance.LoanCalc(c, *amount, 0.069, 5, "USD", finance.LoanCalcOptions{Title: "Car loan"})
		if ui.Slider(c, amount, 5000, 60000).Label("Loan amount").FillWidth().Changed() {
			log(fmt.Sprintf("amount %.0f", *amount))
		}
	})
}

func demoMortgage(c *ui.Context, log func(string)) {
	ui.Column(c).Width(460).Children(func() {
		finance.Mortgage(c, 420000, 0.0625, 30, "USD", finance.MortgageOptions{})
	})
}

func demoTaxSummary(c *ui.Context, log func(string)) {
	ui.Column(c).Width(400).Children(func() {
		finance.TaxSummary(c, 96500, 14600, 0.22, "USD", finance.TaxSummaryOptions{Title: "2026 federal estimate"})
	})
}

func demoSplitBill(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	people := intState(root, "people", 4)
	root.Children(func() {
		if finance.SplitBill(c, 186.40, people, "USD", finance.SplitBillOptions{Title: "Dinner at Nopa"}).Changed() {
			log("people " + fmt.Sprint(*people))
		}
	})
}

func demoSavingsGoal(c *ui.Context, log func(string)) {
	ui.Column(c).Width(380).Gap(12).Children(func() {
		finance.SavingsGoal(c, "Kyoto trip", 2840, 4500, "Apr 1, 2027", "USD", finance.SavingsGoalOptions{Icon: "plane"})
		finance.SavingsGoal(c, "Emergency fund", 12000, 12000, "Dec 31, 2026", "USD", finance.SavingsGoalOptions{Icon: "shield"})
	})
}

func demoFinanceTransfer(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	amount := textState(root, "amount", "250.00")
	confirmed := boolState(root, "confirmed", false)
	root.Children(func() {
		root, _ := finance.Transfer(c, "Everyday checking", "High-yield savings", amount, "USD", confirmed, finance.TransferOptions{})
		if root.Changed() {
			if *confirmed {
				log("confirmed")
			} else {
				log("amount " + fmt.Sprint(utf8.RuneCountInString(*amount)) + " runes")
			}
		}
		if *confirmed && ui.Button(c, "New transfer").Clicked() {
			*confirmed = false
			log("reset")
		}
	})
}

func demoCardControls(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	frozen, online := boolState(root, "frozen", false), boolState(root, "online", true)
	limit := numberState(root, "limit", 2500)
	root.Children(func() {
		if finance.CardControls(c, frozen, online, limit, "USD", finance.CardControlsOptions{Title: "Visa ending 4242"}).Changed() {
			log(fmt.Sprintf("frozen %v online %v limit %.0f", *frozen, *online, *limit))
		}
	})
}

func demoStatement(c *ui.Context, log func(string)) {
	rows := []finance.Transaction{
		{Merchant: "Acme Corp payroll", Category: "Salary", Amount: 4200, When: "Oct 1"},
		{Merchant: "Bayview Apartments", Category: "Rent", Amount: -2100, When: "Oct 1", Icon: "house"},
		demoFinanceTransaction(),
		{Merchant: "PG&E", Category: "Utilities", Amount: -86.40, When: "Oct 14", Icon: "zap"},
	}
	ui.Column(c).Width(440).Children(func() {
		finance.Statement(c, "October 2026", rows, "USD", finance.StatementOptions{})
	})
}

func demoCashFlow(c *ui.Context, log func(string)) {
	months := []finance.FlowMonth{{Label: "May", In: 5200, Out: 3900}, {Label: "Jun", In: 5200, Out: 4650},
		{Label: "Jul", In: 6100, Out: 4200}, {Label: "Aug", In: 5200, Out: 5480}, {Label: "Sep", In: 5400, Out: 3720},
		{Label: "Oct", In: 5900, Out: 4110}}
	ui.Column(c).Width(460).Children(func() { finance.CashFlow(c, months, "USD", finance.CashFlowOptions{}) })
}

func demoDividends(c *ui.Context, log func(string)) {
	payments := []finance.Dividend{{Symbol: "AAPL", Name: "Apple", Date: "Nov 14", Amount: 31.25},
		{Symbol: "MSFT", Name: "Microsoft", Date: "Dec 11", Amount: 46.48}, {Symbol: "KO", Name: "Coca-Cola", Date: "Dec 15", Amount: 61.20},
		{Symbol: "VTI", Name: "Vanguard Total Market", Date: "Dec 26", Amount: 118.90}}
	ui.Column(c).Width(420).Children(func() {
		finance.Dividends(c, payments, "USD", finance.DividendsOptions{Title: "Upcoming dividends"})
	})
}

func demoRiskMeter(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400).Gap(12)
	level := intState(root, "level", 3)
	root.Children(func() {
		finance.RiskMeter(c, *level, finance.RiskMeterOptions{Title: "Growth portfolio"})
		if ui.Button(c, "Raise risk").Clicked() {
			*level = *level%5 + 1
			log(fmt.Sprintf("level %d", *level))
		}
	})
}

func demoExchangeRate(c *ui.Context, log func(string)) {
	points := []float64{0.9312, 0.9288, 0.9301, 0.9254, 0.9230, 0.9246, 0.9198, 0.9214}
	ui.Column(c).Width(400).Gap(12).Children(func() {
		finance.ExchangeRate(c, "USD/EUR", 0.9214, -1.05, finance.ExchangeRateOptions{Points: points})
		finance.ExchangeRate(c, "GBP/USD", 1.2768, 0.42, finance.ExchangeRateOptions{Points: []float64{1.262, 1.265, 1.268, 1.266, 1.2768}})
	})
}

func demoPayees(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	query, selected := textState(root, "query", ""), intState(root, "selected", -1)
	payees := []finance.Payee{{Name: "Ana Martínez", Detail: "Chase •••• 4821"}, {Name: "Kai Nakamura", Detail: "Wells Fargo •••• 0937"},
		{Name: "Bayview Apartments", Detail: "Rent · Bank of America •••• 5512"}, {Name: "Priya Shah", Detail: "Venmo @priya-s"},
		{Name: "Golden Gate Dental", Detail: "Billing •••• 2290"}}
	root.Children(func() {
		if finance.Payees(c, payees, query, selected, finance.PayeesOptions{}).Changed() {
			log(fmt.Sprintf("selected %d query %d", *selected, utf8.RuneCountInString(*query)))
		}
	})
}

func demoCategorize(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	selected := intState(root, "selected", 0)
	categories := []string{"Groceries", "Dining", "Household", "Health", "Gifts"}
	icons := []string{"shopping-cart", "utensils", "house", "heart-pulse", "gift"}
	root.Children(func() {
		if finance.Categorize(c, demoFinanceTransaction(), categories, selected, "USD", finance.CategorizeOptions{Icons: icons}).Changed() {
			log("selected " + fmt.Sprint(*selected))
		}
	})
}

func demoFinanceReceipts(c *ui.Context, log func(string)) {
	receipts := []finance.Receipt{{Merchant: "Trader Joe's", Date: "Oct 21", Items: 14, Amount: 67.18},
		{Merchant: "Shell", Date: "Oct 19", Items: 1, Amount: 48.30}, {Merchant: "Office Depot", Date: "Oct 16", Items: 6, Amount: 112.45},
		{Merchant: "Uber", Date: "Oct 12", Items: 1, Amount: 23.90}}
	ui.Column(c).Width(420).Children(func() { finance.Receipts(c, receipts, "USD", finance.ReceiptsOptions{}) })
}

func demoCreditScore(c *ui.Context, log func(string)) {
	ui.Column(c).Width(400).Children(func() {
		finance.CreditScore(c, 742, finance.CreditScoreOptions{Note: "Updated Oct 20 · Experian"})
	})
}

func demoMarketStatus(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400).Gap(12)
	start := ui.Local(root, "start", func() time.Time { return c.Now() })
	root.Children(func() {
		finance.MarketStatus(c, true, start.Add(2*time.Hour+14*time.Minute), finance.MarketStatusOptions{Open: "NYSE open"})
		finance.MarketStatus(c, false, start.Add(17*time.Hour+44*time.Minute), finance.MarketStatusOptions{Closed: "LSE closed"})
	})
}

func demoPnL(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380).Gap(12)
	day := intState(root, "day", 0)
	values := []float64{1284.60, -412.35, 96.10}
	root.Children(func() {
		finance.PnL(c, values[*day], "USD", finance.PnLOptions{Period: "Today", Basis: 48250})
		if ui.Button(c, "Next day").Clicked() {
			*day = (*day + 1) % len(values)
			log(fmt.Sprintf("pnl %.2f", values[*day]))
		}
	})
}
