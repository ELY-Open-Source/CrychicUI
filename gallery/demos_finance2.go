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
	parts := []finance.AllocationPart{{Name: "Stocks", Share: 0.6}, {Name: "Cash", Share: 0.4}}
	finance.Allocation(c, "Allocation", parts)
}

func demoWallet(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	sent := intState(root, "sent", -1)
	coins := []finance.Coin{{Symbol: "BTC", Amount: 0.5, Value: 20000}}
	root.Children(func() {
		if finance.Wallet(c, coins, "USD", sent, finance.WalletOptions{}).Changed() {
			log("sent " + fmt.Sprint(*sent))
		}
	})
}

func demoLoanCalc(c *ui.Context, log func(string)) {
	finance.LoanCalc(c, 1000, 0.12, 1, "USD", finance.LoanCalcOptions{})
}

func demoMortgage(c *ui.Context, log func(string)) {
	finance.Mortgage(c, 1000, 0.12, 1, "USD", finance.MortgageOptions{})
}

func demoTaxSummary(c *ui.Context, log func(string)) {
	finance.TaxSummary(c, 1000, 100, 0.1, "USD", finance.TaxSummaryOptions{})
}

func demoSplitBill(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	people := intState(root, "people", 2)
	root.Children(func() {
		if finance.SplitBill(c, 100, people, "USD", finance.SplitBillOptions{}).Changed() {
			log("people " + fmt.Sprint(*people))
		}
	})
}

func demoSavingsGoal(c *ui.Context, log func(string)) {
	finance.SavingsGoal(c, "Home", 500, 1000, "2027-01-01", "USD", finance.SavingsGoalOptions{})
}

func demoFinanceTransfer(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	amount := textState(root, "amount", "100")
	confirmed := boolState(root, "confirmed", false)
	root.Children(func() {
		root, _ := finance.Transfer(c, "Checking", "Savings", amount, "USD", confirmed, finance.TransferOptions{})
		if root.Changed() {
			log("amount " + fmt.Sprint(utf8.RuneCountInString(*amount)) + " runes")
		}
		if *confirmed {
			log("confirmed")
			*confirmed = false
		}
	})
}

func demoCardControls(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	frozen, online := boolState(root, "frozen", false), boolState(root, "online", true)
	limit := numberState(root, "limit", 1000)
	root.Children(func() {
		if finance.CardControls(c, frozen, online, limit, finance.CardControlsOptions{}).Changed() {
			log(fmt.Sprintf("frozen %v online %v", *frozen, *online))
		}
	})
}

func demoStatement(c *ui.Context, log func(string)) {
	finance.Statement(c, "October", []finance.Transaction{demoFinanceTransaction()}, "USD", finance.StatementOptions{})
}

func demoCashFlow(c *ui.Context, log func(string)) {
	months := []finance.FlowMonth{{Label: "Jan", In: 100, Out: 40}}
	finance.CashFlow(c, months, "USD")
}

func demoDividends(c *ui.Context, log func(string)) {
	payments := []finance.Dividend{{Symbol: "MYGO", Date: "Nov 1", Amount: 5}}
	finance.Dividends(c, payments, "USD")
}

func demoRiskMeter(c *ui.Context, log func(string)) {
	finance.RiskMeter(c, 4, finance.RiskMeterOptions{})
}

func demoExchangeRate(c *ui.Context, log func(string)) {
	finance.ExchangeRate(c, "USD/EUR", 0.85, 1.2)
}

func demoPayees(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	query, selected := textState(root, "query", ""), intState(root, "selected", -1)
	root.Children(func() {
		if finance.Payees(c, []string{"Ana", "Kai"}, query, selected, finance.PayeesOptions{}).Changed() {
			log("selected " + fmt.Sprint(*selected))
		}
	})
}

func demoCategorize(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := intState(root, "selected", -1)
	root.Children(func() {
		if finance.Categorize(c, demoFinanceTransaction(), []string{"Books", "Food"}, selected).Changed() {
			log("selected " + fmt.Sprint(*selected))
		}
	})
}

func demoFinanceReceipts(c *ui.Context, log func(string)) {
	receipts := []finance.Receipt{{Merchant: "Bookstore", Date: "Oct 1", Amount: 24}}
	finance.Receipts(c, receipts, "USD", finance.ReceiptsOptions{})
}

func demoCreditScore(c *ui.Context, log func(string)) {
	finance.CreditScore(c, 800, finance.CreditScoreOptions{})
}

func demoMarketStatus(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	next := ui.Local(root, "next", func() time.Time { return c.Now().Add(time.Hour) })
	root.Children(func() {
		finance.MarketStatus(c, true, *next, finance.MarketStatusOptions{})
	})
}

func demoPnL(c *ui.Context, log func(string)) {
	finance.PnL(c, -25, "USD", finance.PnLOptions{})
}
