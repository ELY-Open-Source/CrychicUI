package gallery

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/code"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

var codeEntries2 = []Entry{
	{ID: 921, Name: "Breadcrumb", Category: "code", Summary: "Folder, file and symbol path; cue Tap.", Demo: demoCodeBreadcrumb},
	{ID: 922, Name: "SymbolSearch", Category: "code", Summary: "Symbols filtered as typed; cue Select.", Demo: demoCodeSymbolSearch},
	{ID: 923, Name: "RegexTester", Category: "code", Summary: "Pattern, text, matches marked and listed.", Demo: demoRegexTester},
	{ID: 924, Name: "JSONFormat", Category: "code", Summary: "Formats or minifies JSON; cue Success, Error.", Demo: demoJSONFormat},
	{ID: 925, Name: "Base64Tool", Category: "code", Summary: "Encode and decode panes; cue Tick.", Demo: demoBase64Tool},
	{ID: 926, Name: "HashTool", Category: "code", Summary: "SHA-256 of text; digest cross-fades.", Demo: demoHashTool},
	{ID: 927, Name: "UnixTime", Category: "code", Summary: "Epoch and date converter; cue Tick.", Demo: demoUnixTime},
	{ID: 928, Name: "Commit", Category: "code", Summary: "Author, message, age and hash; cue Pop on copy.", Demo: demoCodeCommit},
	{ID: 929, Name: "BranchPicker", Category: "code", Summary: "Branches with search and current mark; cue Select.", Demo: demoBranchPicker},
	{ID: 930, Name: "PullRequest", Category: "code", Summary: "Title, status, branches, checks and reviewers.", Demo: demoPullRequest},
	{ID: 931, Name: "Checks", Category: "code", Summary: "CI checks with spinners, results and times.", Demo: demoCodeChecks},
	{ID: 932, Name: "Blame", Category: "code", Summary: "Code beside who changed it and when.", Demo: demoBlame},
	{ID: 933, Name: "CommitGraph", Category: "code", Summary: "Lanes, merges and dots of commits.", Demo: demoCommitGraph},
	{ID: 934, Name: "Language", Category: "code", Summary: "Language bar of shares; segments grow.", Demo: demoCodeLanguage},
	{ID: 935, Name: "KeyTable", Category: "code", Summary: "Keyboard shortcuts in key caps by group.", Demo: demoKeyTable},
	{ID: 936, Name: "Console", Category: "code", Summary: "Log, warn, error rows with filter; cue Toggle.", Demo: demoCodeConsole},
	{ID: 937, Name: "REPL", Category: "code", Summary: "Input evaluating arithmetic with history; cue Send.", Demo: demoREPL},
	{ID: 938, Name: "PackageCard", Category: "code", Summary: "Name, version, stats, install copy; cue Pop.", Demo: demoPackageCard},
	{ID: 939, Name: "Env", Category: "code", Summary: "KEY=value rows with masked values; cue Toggle.", Demo: demoEnv},
	{ID: 940, Name: "Playground", Category: "code", Summary: "Code editor and output pane; run; cue Tap.", Demo: demoPlayground},
}

// billingChecks is one CI run of the billing service.
var billingChecks = []code.Check{
	{Name: "build", State: code.CheckPass, Duration: 72 * time.Second},
	{Name: "unit tests", State: code.CheckPass, Duration: 160 * time.Second},
	{Name: "lint", State: code.CheckPending},
	{Name: "e2e", State: code.CheckFail, Duration: 242 * time.Second},
}

func demoCodeBreadcrumb(c *ui.Context, log func(string)) {
	parts := []string{"billing", "internal", "invoice", "close.go", "Close"}
	root := ui.Column(c).Width(520)
	at := intState(root, "at", 4)
	root.Children(func() {
		if code.Breadcrumb(c, "Path", parts, at).Changed() {
			log("at " + parts[*at])
		}
	})
}

func demoCodeSymbolSearch(c *ui.Context, log func(string)) {
	symbols := []code.Symbol{
		{Kind: code.SymType, Name: "Invoice", Line: 10},
		{Kind: code.SymType, Name: "InvoiceLine", Line: 17},
		{Kind: code.SymFunc, Name: "Close", Line: 24},
		{Kind: code.SymFunc, Name: "sendInvoice", Line: 36},
		{Kind: code.SymFunc, Name: "Price", Line: 52},
		{Kind: code.SymVar, Name: "ledger", Line: 7},
		{Kind: code.SymConst, Name: "defaultTerms", Line: 4},
	}
	root := ui.Column(c).Width(420)
	query, chosen := textState(root, "query", "inv"), intState(root, "chosen", -1)
	root.Children(func() {
		r, field := code.SymbolSearch(c, "Symbols", symbols, query, chosen, code.SymbolSearchOptions{})
		if field.Changed() {
			log(fmt.Sprintf("query %d chars", utf8.RuneCountInString(*query)))
		} else if r.Changed() {
			log("selected " + symbols[*chosen].Name)
		}
	})
}

func demoRegexTester(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	pattern := textState(root, "pattern", `INV-\d{4}`)
	text := textState(root, "text", "Paid INV-0042 and INV-0043; INV-12 is a typo.\nRefund INV-0007 on Friday.")
	root.Children(func() {
		if code.RegexTester(c, "Invoice numbers", pattern, text, code.RegexTesterOptions{}).Changed() {
			log("edited")
		}
	})
}

func demoJSONFormat(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	text := textState(root, "text", `{"invoice":"INV-0042","seats":14,"total":{"amount":12600,"currency":"EUR"},"paid":true}`)
	root.Children(func() {
		r, _, _ := code.JSONFormat(c, "invoice.json", text, code.JSONFormatOptions{})
		if r.Changed() {
			log(fmt.Sprintf("action on %d runes", utf8.RuneCountInString(*text)))
		}
	})
}

func demoBase64Tool(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	text := textState(root, "text", "ws_8f2k:invoice-2026-10")
	encode := boolState(root, "encode", true)
	root.Children(func() {
		r, cp := code.Base64Tool(c, "Base64", text, encode, code.Base64ToolOptions{})
		switch {
		case cp.Clicked():
			log("copied")
		case r.Changed():
			log(fmt.Sprintf("encode %v", *encode))
		}
	})
}

func demoHashTool(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	text := textState(root, "text", "billing-v2.4.1.tar.gz")
	root.Children(func() {
		r, cp := code.HashTool(c, "SHA-256", text, code.HashToolOptions{})
		switch {
		case cp.Clicked():
			log("copied digest")
		case r.Changed():
			log("edited")
		}
	})
}

func demoUnixTime(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	epoch := int64State(root, "epoch", 1761292800)
	root.Children(func() {
		if r, _ := code.UnixTime(c, "Unix time", epoch, code.UnixTimeOptions{}); r.Changed() {
			log(fmt.Sprintf("epoch %d", *epoch))
		}
	})
}

func demoCodeCommit(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	at := ui.Local(root, "at", func() time.Time { return c.Now().Add(-3 * time.Hour) })
	root.Children(func() {
		if _, cp := code.Commit(c, "8fa0fc3e1b2d", "Close invoices on the 1st of each month", "Mei Chen", *at, code.CommitOptions{}); cp.Clicked() {
			log("copied hash")
		}
	})
}

func demoBranchPicker(c *ui.Context, log func(string)) {
	branches := []code.Branch{
		{Name: "main"},
		{Name: "feature/eu-vat", Ahead: 3},
		{Name: "fix/rounding", Ahead: 1, Behind: 4},
		{Name: "release/2.4", Behind: 12},
		{Name: "chore/deps", Ahead: 2},
	}
	root := ui.Column(c).Width(420)
	current, query := intState(root, "current", 0), textState(root, "query", "")
	root.Children(func() {
		r, field := code.BranchPicker(c, "Branches", branches, current, query, code.BranchPickerOptions{})
		if field.Changed() {
			log(fmt.Sprintf("query %d chars", utf8.RuneCountInString(*query)))
		} else if r.Changed() {
			log("selected " + branches[*current].Name)
		}
	})
}

func demoPullRequest(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() {
		code.PullRequest(c, code.PR{
			Number: 1842, Title: "Charge VAT for EU workspaces", Author: "Ana Reyes", Head: "feature/eu-vat", Base: "main",
			State: code.PROpen, Added: 248, Removed: 61, Checks: billingChecks,
			Reviewers: []string{"Kai Mori", "Mei Chen", "Lena Fischer"},
		}, code.PullRequestOptions{})
	})
}

func demoCodeChecks(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() { code.Checks(c, "CI · feature/eu-vat", billingChecks, code.ChecksOptions{}) })
}

func demoBlame(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	now := ui.Local(root, "now", func() time.Time { return c.Now() })
	ago := func(d time.Duration) time.Time { return now.Add(-d) }
	lines := []code.BlameLine{
		{Author: "Mei Chen", Text: "// Close totals the month and sends the invoice.", At: ago(41 * 24 * time.Hour)},
		{Author: "Mei Chen", Text: "func Close(ctx context.Context, ws Workspace) (*Invoice, error) {", At: ago(41 * 24 * time.Hour)},
		{Author: "Kai Mori", Text: "\tseats, err := ws.ActiveSeats(ctx)", At: ago(6 * 24 * time.Hour)},
		{Author: "Kai Mori", Text: "\tif err != nil {", At: ago(6 * 24 * time.Hour)},
		{Author: "Kai Mori", Text: "\t\treturn nil, fmt.Errorf(\"count seats: %w\", err)", At: ago(6 * 24 * time.Hour)},
		{Author: "Kai Mori", Text: "\t}", At: ago(6 * 24 * time.Hour)},
		{Author: "Ana Reyes", Text: "\tinv := &Invoice{Workspace: ws.ID, Seats: seats}", At: ago(3 * time.Hour)},
		{Author: "Ana Reyes", Text: "\tinv.Total = Price(ws.Plan, seats)", At: ago(3 * time.Hour)},
		{Author: "Mei Chen", Text: "\treturn inv, send(ctx, inv, 30*time.Second)", At: ago(41 * 24 * time.Hour)},
		{Author: "Mei Chen", Text: "}", At: ago(41 * 24 * time.Hour)},
	}
	root.Children(func() { code.Blame(c, "invoice.go", lines, code.BlameOptions{Lang: "go"}) })
}

func demoCommitGraph(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() {
		code.CommitGraph(c, "History", []code.GraphCommit{
			{Hash: "4d8da0d9", Message: "Merge feature/eu-vat", Author: "Ana Reyes", Lane: 0, Parents: []string{"35b918c2", "80637aa1"}, Refs: []string{"main"}},
			{Hash: "80637aa1", Message: "Charge VAT on EU invoices", Author: "Kai Mori", Lane: 1, Parents: []string{"5465e35b"}, Refs: []string{"feature/eu-vat"}},
			{Hash: "35b918c2", Message: "Round totals to cents", Author: "Mei Chen", Lane: 0, Parents: []string{"8fa0fc3e"}},
			{Hash: "5465e35b", Message: "Add VAT rates table", Author: "Kai Mori", Lane: 1, Parents: []string{"8fa0fc3e"}},
			{Hash: "8fa0fc3e", Message: "Close invoices on the 1st", Author: "Mei Chen", Lane: 0, Parents: []string{"1849e62a"}},
			{Hash: "1849e62a", Message: "Send invoice emails", Author: "Ana Reyes", Lane: 0, Parents: []string{"0c1d2e3f"}, Refs: []string{"v2.4.0"}},
		})
	})
}

func demoCodeLanguage(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		code.Language(c, "Languages", []code.LangShare{
			{Name: "Go", Share: 0.68}, {Name: "TypeScript", Share: 0.17}, {Name: "Shell", Share: 0.08}, {Name: "CSS", Share: 0.04}, {Name: "Markdown", Share: 0.03},
		})
	})
}

func demoKeyTable(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		code.KeyTable(c, "Shortcuts", []code.KeyGroup{
			{Name: "Editor", Shortcuts: []code.Shortcut{
				{Keys: "Ctrl+S", Action: "Save the file"},
				{Keys: "Ctrl+Shift+P", Action: "Open the command palette"},
				{Keys: "Ctrl+/", Action: "Toggle a line comment"},
			}},
			{Name: "Debug", Shortcuts: []code.Shortcut{
				{Keys: "F5", Action: "Continue"},
				{Keys: "F10", Action: "Step over"},
				{Keys: "Shift+F11", Action: "Step out"},
			}},
		})
	})
}

func demoCodeConsole(c *ui.Context, log func(string)) {
	entries := []code.LogEntry{
		{Level: code.LevelLog, Text: "billing listening on :8080", Source: "main.go:31"},
		{Level: code.LevelLog, Text: "closed 214 invoices for October", Source: "invoice.go:42"},
		{Level: code.LevelWarn, Text: "VAT rate missing for country XK, using 0%", Source: "tax.go:88"},
		{Level: code.LevelLog, Text: "sent 214 invoice emails", Source: "mail.go:17"},
		{Level: code.LevelError, Text: "stripe: card_declined for ws_8f2k", Source: "charge.go:63"},
	}
	root := ui.Column(c).Width(520)
	filter := intState(root, "filter", 7)
	root.Children(func() {
		if code.Console(c, "Console", entries, filter, code.ConsoleOptions{}).Changed() {
			log(fmt.Sprintf("filter %d", *filter))
		}
	})
}

func demoREPL(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	history := ui.Local(root, "history", func() []code.REPLLine {
		return []code.REPLLine{{In: "12600 * 0.19", Out: "2394"}, {In: "(12600 + 2394) / 14", Out: "1071"}}
	})
	draft := textState(root, "draft", "")
	root.Children(func() {
		if r, _ := code.REPL(c, "Calculator", history, draft, code.REPLOptions{}); r.Changed() {
			log(fmt.Sprintf("%d lines", len(*history)))
		}
	})
}

func demoPackageCard(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		_, cp := code.PackageCard(c, code.Package{
			Name: "github.com/egoist/mygo", Version: "v0.5.5", License: "MIT", Downloads: 12400,
			Description: "Native UI for Go: one codebase for desktop, mobile and the web.",
		}, code.PackageCardOptions{})
		if cp.Clicked() {
			log("copied install")
		}
	})
}

func demoEnv(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	shown := boolState(root, "shown", false)
	vars := []code.EnvVar{
		{Key: "DATABASE_URL", Value: "postgres://billing@db.internal:5432/billing"},
		{Key: "STRIPE_SECRET_KEY", Value: "sk_test_sample_9f2c"},
		{Key: "REGION", Value: "eu-central-1"},
		{Key: "MAX_RETRIES", Value: "5"},
		{Key: "LOG_LEVEL", Value: "info"},
	}
	root.Children(func() {
		if code.Env(c, ".env.production", vars, shown, code.EnvOptions{}).Changed() {
			log("values " + map[bool]string{true: "shown", false: "hidden"}[*shown])
		}
	})
}

const playgroundSource = `package main

import "fmt"

func main() {
	seats, rate := 14, 9.00
	fmt.Println("Invoice INV-0042")
	fmt.Printf("Seats %d\n", seats)
	fmt.Printf("Total €%.2f\n", float64(seats)*rate)
}
`

func demoPlayground(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	st := ui.Local(root, "editor", func() *editor.State {
		buf := editor.NewBuffer(playgroundSource)
		buf.SetLanguage(editor.LanguageFor("main.go"))
		return &editor.State{Buffer: buf}
	})
	output := textState(root, "output", "")
	root.Children(func() {
		r, run := code.Playground(c, "main.go", *st, output, func(src string) {
			*output = "Invoice INV-0042\nSeats 14\nTotal €126.00\n\nProgram exited."
		}, code.PlaygroundOptions{Height: 220, OutputHeight: 130})
		switch {
		case run.Clicked():
			log(fmt.Sprintf("ran %d runes", utf8.RuneCountInString((*st).Buffer.String())))
		case r.Changed():
			log(fmt.Sprintf("edited, %d lines", (*st).Buffer.LineCount()))
		}
	})
}

func int64State(root ui.Element, key string, initial int64) *int64 {
	return ui.Local(root, key, func() int64 { return initial })
}
