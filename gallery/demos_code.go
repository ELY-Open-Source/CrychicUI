package gallery

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/code"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

var codeEntries = []Entry{
	{901, "CodeBlock", "code", "Numbered, highlighted code with copy.", demoCodeBlock},
	{902, "CodeEditor", "code", "MyGo's code editor with a status bar.", demoCodeEditor},
	{903, "DiffView", "code", "MyGo's diff viewer, split or unified.", demoDiffView},
	{904, "Terminal", "code", "Output lines and a blinking prompt.", demoTerminal},
	{905, "CommandLine", "code", "A prompt with history and output rows.", demoCommandLine},
	{906, "CodeTabs", "code", "Files as tabs over numbered code.", demoCodeTabs},
	{907, "FileHeader", "code", "A path, language and actions bar.", demoFileHeader},
	{908, "Snippet", "code", "Short code with a title and copy.", demoSnippet},
	{909, "LineHighlight", "code", "Code with marked lines pulsing once.", demoLineHighlight},
	{910, "Annotated", "code", "Code with numbered notes.", demoAnnotatedCode},
	{911, "StackTrace", "code", "An error over frames; library frames fold.", demoStackTrace},
	{912, "ErrorOverlay", "code", "A compile error panel sliding down.", demoErrorOverlay},
	{913, "Gutter", "code", "Line numbers that toggle breakpoints.", demoGutterCode},
	{914, "DebugBar", "code", "Continue, step over, in, out and stop.", demoDebugBar},
	{915, "Variables", "code", "A tree of names, values and types.", demoVariables},
	{916, "Watches", "code", "Expressions and values; add and remove.", demoWatches},
	{917, "CallStack", "code", "Frames with the current one marked.", demoCallStack},
	{918, "Problems", "code", "Errors and warnings by file.", demoProblems},
	{919, "Outline", "code", "Symbols of a file by kind.", demoOutlineCode},
	{920, "Minimap", "code", "A file drawn small with a viewport box.", demoMinimap},
}

const billingClose = `// Close totals the month and sends the invoice.
func Close(ctx context.Context, ws Workspace) (*Invoice, error) {
	seats, err := ws.ActiveSeats(ctx)
	if err != nil {
		return nil, fmt.Errorf("count seats: %w", err)
	}
	inv := &Invoice{Workspace: ws.ID, Seats: seats}
	inv.Total = Price(ws.Plan, seats)
	return inv, send(ctx, inv, 30*time.Second)
}
`

const billingPrice = `// Price is the monthly cost of a plan.
func Price(p Plan, seats int) Money {
	base := Money(seats) * p.Rate
	if p.Annual && seats >= 5 {
		return base * 10 * 9 / 10
	}
	return base * 12
}
`

const billingTest = `func TestPriceAnnual(t *testing.T) {
	got := Price(Plan{Rate: 900, Annual: true}, 5)
	if got != 40500 {
		t.Errorf("price %d, want 40500", got)
	}
}
`

// billingFrames is a panic inside the billing service.
var billingFrames = []code.StackFrame{
	{Func: "billing.(*Ledger).Add", File: "ledger.go", Line: 58, App: true},
	{Func: "billing.Close", File: "invoice.go", Line: 31, App: true},
	{Func: "net/http.HandlerFunc.ServeHTTP", File: "server.go", Line: 2136},
	{Func: "net/http.serverHandler.ServeHTTP", File: "server.go", Line: 2938},
	{Func: "net/http.(*conn).serve", File: "server.go", Line: 1999},
	{Func: "runtime.goexit", File: "asm_amd64.s", Line: 1700},
}

func goBuffer(src string) *editor.Buffer {
	buf := editor.NewBuffer(src)
	buf.SetLanguage(editor.LanguageFor("main.go"))
	return buf
}

func demoCodeBlock(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() {
		_, cp := code.CodeBlock(c, "Go", billingClose, code.CodeBlockOptions{Title: "billing/invoice.go"})
		if cp.Clicked() {
			log(fmt.Sprintf("copied %d runes", utf8.RuneCountInString(billingClose)))
		}
	})
}

func demoCodeEditor(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	st := ui.Local(root, "editor", func() *editor.State { return &editor.State{Buffer: goBuffer(billingClose)} })
	root.Children(func() {
		if code.CodeEditor(c, *st, code.CodeEditorOptions{Height: 300, Name: "invoice.go"}).Changed() {
			log(fmt.Sprintf("edited, %d lines", (*st).Buffer.LineCount()))
		}
	})
}

func demoDiffView(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	st := ui.Local(root, "diff", func() *editor.DiffState {
		old := strings.Replace(billingPrice, "p.Annual && seats >= 5", "p.Annual", 1)
		return &editor.DiffState{Old: goBuffer(old), New: goBuffer(billingPrice)}
	})
	split := boolState(root, "split", false)
	root.Children(func() {
		if code.DiffView(c, *st, split, code.DiffViewOptions{Height: 280, Name: "billing/price.go"}).Changed() {
			log(map[bool]string{true: "layout split", false: "layout unified"}[*split])
		}
	})
}

func demoTerminal(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480).Gap(12)
	lines := ui.Local(root, "lines", func() []string {
		return []string{
			"$ go build -o bin/api ./cmd/api",
			"$ go test ./billing/...",
			"ok  	github.com/acme/billing	0.412s",
			"ok  	github.com/acme/billing/tax	0.087s",
		}
	})
	runs := intState(root, "runs", 0)
	root.Children(func() {
		code.Terminal(c, "zsh — acme/api", *lines, code.TerminalOptions{Height: 220})
		if button.Outline(c, "Run vet").Clicked() {
			*runs++
			*lines = append(*lines, "$ go vet ./...", fmt.Sprintf("vet: %d packages clean", 11+*runs))
			log(fmt.Sprintf("%d lines", len(*lines)))
		}
	})
}

func demoCommandLine(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	history := ui.Local(root, "history", func() []code.Command {
		return []code.Command{
			{Text: "git status --short", Output: " M billing/price.go\n?? billing/tax/vat.go"},
			{Text: "git log --oneline -2", Output: "8fa0fc3 Close invoices on the 1st\n5465e35 Add VAT for EU workspaces"},
		}
	})
	draft := ui.Local(root, "draft", func() string { return "" })
	root.Children(func() {
		cl, _ := code.CommandLine(c, "zsh — acme/api", *history, draft, func(s string) {
			out := "zsh: command not found: " + strings.Fields(s)[0]
			switch strings.Fields(s)[0] {
			case "pwd":
				out = "/Users/ana/src/acme/api"
			case "ls":
				out = "billing  cmd  go.mod  go.sum  README.md"
			case "date":
				out = "Fri Oct 24 09:41:00 CEST 2026"
			}
			*history = append(*history, code.Command{Text: s, Output: out})
			log(fmt.Sprintf("ran %d runes", utf8.RuneCountInString(s)))
		}, code.CommandLineOptions{Height: 200})
		if cl.Changed() {
			log(fmt.Sprintf("draft %d runes", utf8.RuneCountInString(*draft)))
		}
	})
}

func demoCodeTabs(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	current := intState(root, "current", 0)
	files := []code.CodeFile{
		{Name: "invoice.go", Lang: "go", Code: billingClose},
		{Name: "price.go", Lang: "go", Code: billingPrice},
		{Name: "price_test.go", Lang: "go", Code: billingTest},
	}
	root.Children(func() {
		if code.CodeTabs(c, files, current).Changed() {
			log("selected " + files[*current].Name)
		}
	})
}

func demoFileHeader(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		_, cp := code.FileHeader(c, "billing/internal/invoice/close.go", "Go", code.FileHeaderOptions{Meta: "84 lines · 2.6 KB"})
		if cp.Clicked() {
			log("copied the path")
		}
	})
}

func demoSnippet(c *ui.Context, log func(string)) {
	src := "curl -fsSL https://get.acme.dev/cli | sh"
	ui.Column(c).Width(480).Children(func() {
		_, cp := code.Snippet(c, "Install the CLI", "sh", src, code.SnippetOptions{Description: "Puts acme in /usr/local/bin and checks its signature."})
		if cp.Clicked() {
			log(fmt.Sprintf("copied %d runes", utf8.RuneCountInString(src)))
		}
	})
}

func demoLineHighlight(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() { code.LineHighlight(c, "go", billingClose, []int{4, 5, 6}) })
}

func demoAnnotatedCode(c *ui.Context, log func(string)) {
	notes := []code.Note{
		{Line: 3, Text: "Seats count per workspace, not per user."},
		{Line: 4, Text: "Annual plans of five seats or more get ten percent off."},
		{Line: 7, Text: "Monthly plans pay a full year up front."},
	}
	ui.Column(c).Width(520).Children(func() { code.Annotated(c, "go", billingPrice, notes) })
}

func demoStackTrace(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	open := boolState(root, "open", false)
	root.Children(func() {
		if code.StackTrace(c, "panic: assignment to entry in nil map", billingFrames, open, code.StackTraceOptions{}).Changed() {
			log(fmt.Sprintf("open %v", *open))
		}
	})
}

func demoErrorOverlay(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520).Gap(12)
	open := boolState(root, "open", true)
	root.Children(func() {
		box := code.ErrorOverlay(c, "Build failed", "billing/invoice.go:9:14: undefined: send\n\treturn inv, send(ctx, inv, 30*time.Second)",
			open, code.ErrorOverlayOptions{}, func() {
				if _, cp := code.CodeBlock(c, "Go", billingClose, code.CodeBlockOptions{Title: "billing/invoice.go"}); cp.Clicked() {
					log(fmt.Sprintf("copied %d runes", utf8.RuneCountInString(billingClose)))
				}
			})
		if box.Changed() {
			log(fmt.Sprintf("open %v", *open))
		}
		if button.Outline(c, "Rebuild").Clicked() {
			*open = true
			log("open true")
		}
	})
}

func demoGutterCode(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	breakpoints := ui.Local(root, "breakpoints", func() []int { return []int{4} })
	root.Children(func() {
		if code.Gutter(c, "go", billingClose, breakpoints, code.GutterOptions{}).Changed() {
			log(fmt.Sprintf("%d breakpoints", len(*breakpoints)))
		}
	})
}

func demoDebugBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	running := boolState(root, "running", false)
	names := []string{"continue", "over", "into", "out", "stop"}
	root.Children(func() {
		_, cont, over, into, out, stop := code.DebugBar(c, *running, code.DebugBarOptions{})
		for i, b := range []ui.Element{cont, over, into, out, stop} {
			if b.Clicked() {
				*running = i == 0
				log(names[i])
			}
		}
		where := "Paused on invoice.go:8 in billing.Close"
		if *running {
			where = "Running billing.Close…"
		}
		ui.Text(c, where).FontSize(c.Theme().Rem(0.85)).TextColor(c.Theme().TextMuted)
	})
}

func demoVariables(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	vars := []code.Variable{
		{Name: "ws", Type: "Workspace", Value: "{…}", Children: []code.Variable{
			{Name: "ID", Type: "string", Value: `"ws_8f2k"`},
			{Name: "Plan", Type: "Plan", Value: "{…}", Children: []code.Variable{
				{Name: "Rate", Type: "Money", Value: "900"},
				{Name: "Annual", Type: "bool", Value: "true"},
			}},
		}},
		{Name: "seats", Type: "int", Value: "14"},
		{Name: "inv", Type: "*Invoice", Value: "nil"},
		{Name: "err", Type: "error", Value: "nil"},
	}
	root.Children(func() {
		if code.Variables(c, "Locals", vars).Changed() {
			log("toggled")
		}
	})
}

func demoWatches(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	watches := ui.Local(root, "watches", func() []code.Watch {
		return []code.Watch{{Expr: "seats * p.Rate", Value: "12600"}, {Expr: "ws.Plan.Annual", Value: "true"}, {Expr: "inv.Total", Value: ""}}
	})
	root.Children(func() {
		if code.Watches(c, "Watch", watches, code.WatchesOptions{}).Changed() {
			log(fmt.Sprintf("%d watches", len(*watches)))
		}
	})
}

func demoCallStack(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	current := intState(root, "current", 0)
	root.Children(func() {
		if code.CallStack(c, "Call stack", billingFrames, current).Changed() {
			log(fmt.Sprintf("selected %d", *current))
		}
	})
}

func demoProblems(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	selected := intState(root, "selected", -1)
	problems := []code.Problem{
		{Severity: code.ProbError, File: "invoice.go", Line: 9, Message: "undefined: send"},
		{Severity: code.ProbWarning, File: "invoice.go", Line: 7, Message: "inv escapes to the heap"},
		{Severity: code.ProbError, File: "price.go", Line: 3, Message: "cannot use seats (variable of type int64) as int value"},
		{Severity: code.ProbWarning, File: "price.go", Line: 5, Message: "integer division truncates the discount"},
	}
	root.Children(func() {
		if code.Problems(c, "Problems", problems, selected, code.ProblemsOptions{}).Changed() {
			log("selected " + problems[*selected].Message)
		}
	})
}

func demoOutlineCode(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	selected := intState(root, "selected", -1)
	symbols := []code.Symbol{
		{Kind: code.SymConst, Name: "defaultTerms", Line: 4},
		{Kind: code.SymVar, Name: "ledger", Line: 7},
		{Kind: code.SymType, Name: "Invoice", Line: 10},
		{Kind: code.SymType, Name: "Line", Line: 17},
		{Kind: code.SymFunc, Name: "Close", Line: 24},
		{Kind: code.SymFunc, Name: "send", Line: 36},
	}
	root.Children(func() {
		if code.Outline(c, "invoice.go", symbols, selected, code.OutlineOptions{}).Changed() {
			log("selected " + symbols[*selected].Name)
		}
	})
}

// ledgerSource is a long file for the minimap.
func ledgerSource() string {
	var b strings.Builder
	b.WriteString("package billing\n\nimport \"fmt\"\n\n")
	for _, name := range []string{"Add", "Remove", "Total", "Close", "Reopen", "Export", "Import", "Audit"} {
		fmt.Fprintf(&b, "// %s edits one ledger.\nfunc (l *Ledger) %s(e Entry) error {\n", name, name)
		b.WriteString("\tl.mu.Lock()\n\tdefer l.mu.Unlock()\n\tif e.Amount == 0 {\n\t\treturn nil\n\t}\n")
		fmt.Fprintf(&b, "\tif err := l.check(e); err != nil {\n\t\treturn fmt.Errorf(\"%s: %%w\", err)\n\t}\n", strings.ToLower(name))
		b.WriteString("\tl.entries = append(l.entries, e)\n\tl.total += e.Amount\n\treturn l.flush()\n}\n\n")
	}
	return b.String()
}

func demoMinimap(c *ui.Context, log func(string)) {
	src := ledgerSource()
	lines := strings.Split(src, "\n")
	const visible = 16
	root := ui.Row(c).Width(520).Gap(12)
	at := intState(root, "at", 20)
	root.Children(func() {
		ui.Column(c).Grow(1).MinWidth(0).Children(func() {
			title := fmt.Sprintf("ledger.go · lines %d–%d", *at+1, *at+visible)
			shown := strings.Join(lines[*at:*at+visible], "\n")
			if _, cp := code.Snippet(c, title, "go", shown, code.SnippetOptions{}); cp.Clicked() {
				log(fmt.Sprintf("copied lines %d–%d", *at+1, *at+visible))
			}
		})
		if code.Minimap(c, "Editor map", "go", src, visible, at, code.MinimapOptions{Height: 360}).Changed() {
			log(fmt.Sprintf("at %d", *at))
		}
	})
}
