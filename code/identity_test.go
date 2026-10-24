package code

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

func TestCommandLineKeepsFocusAsHistoryAppears(t *testing.T) {
	var history []Command
	draft := ""
	tt := ui.NewTester(padded(func(c *ui.Context) {
		CommandLine(c, "Shell", history, &draft, func(s string) { history = append(history, Command{Text: s}) },
			CommandLineOptions{Run: "Execute"})
	}), 500, 400)
	settles(t, tt, "prompt")
	if _, ok := tt.Find("Execute"); !ok {
		t.Fatalf("run copy not used: %q", tt.Texts())
	}
	click(t, tt, "Run a command")
	tt.Type("pwd")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	tt.Type("ls")
	tt.Frame()
	if len(history) != 1 || draft != "ls" {
		t.Errorf("history %v draft %q: the input lost focus", history, draft)
	}
}

func TestCodeEditorKeepsFocusAsBarAppears(t *testing.T) {
	st := &editor.State{Buffer: editor.NewBuffer("")}
	name := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { CodeEditor(c, st, CodeEditorOptions{Name: name}) }), 500, 400)
	tt.Frame()
	tt.Key(0, ui.KeyTab)
	tt.Type("x")
	tt.Frame()
	name = "main.go"
	tt.Frame()
	tt.Type("y")
	tt.Frame()
	if got := st.Buffer.String(); got != "xy" {
		t.Errorf("buffer %q: the editor lost focus", got)
	}
}

// swapUnderPress presses at, reorders, then releases there.
func swapUnderPress(tt *ui.Tester, r ui.Rect, swap func()) {
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Frame()
	swap()
	tt.Frame()
	tt.Release(x, y)
	tt.Frame()
}

func TestRowsKeepIdentityAcrossReorder(t *testing.T) {
	// Fields that join alike must still key apart.
	frames := []StackFrame{{Func: "Widget::f()", File: "const a.cpp", Line: 10}, {Func: "Widget::f() const", File: "a.cpp", Line: 10}}
	current := 0
	cs := ui.NewTester(padded(func(c *ui.Context) { CallStack(c, "Stack", frames, &current) }), 500, 300)
	settles(t, cs, "stack")
	swapUnderPress(cs, find(t, cs, "Widget::f() const"), func() { frames[0], frames[1] = frames[1], frames[0] })
	if current != 0 {
		t.Errorf("call stack chose the row that slid under the pointer: %d", current)
	}

	problems := []Problem{{File: "a.go", Line: 3, Message: "first"}, {File: "a.go", Line: 9, Message: "second"}}
	picked := -1
	ps := ui.NewTester(padded(func(c *ui.Context) { Problems(c, "Problems", problems, &picked, ProblemsOptions{}) }), 500, 300)
	settles(t, ps, "problems")
	swapUnderPress(ps, find(t, ps, "second"), func() { problems[0], problems[1] = problems[1], problems[0] })
	if picked != -1 {
		t.Errorf("problems chose the row that slid under the pointer: %d", picked)
	}

	symbols := []Symbol{{Kind: SymFunc, Name: "Open", Line: 3}, {Kind: SymFunc, Name: "Shut", Line: 9}}
	chosen := -1
	os := ui.NewTester(padded(func(c *ui.Context) { Outline(c, "Outline", symbols, &chosen, OutlineOptions{}) }), 500, 300)
	settles(t, os, "outline")
	swapUnderPress(os, find(t, os, "Shut"), func() { symbols[0], symbols[1] = symbols[1], symbols[0] })
	if chosen != -1 {
		t.Errorf("outline chose the row that slid under the pointer: %d", chosen)
	}
}

func TestErrorOverlayContentKeysAreItsOwn(t *testing.T) {
	open := true
	render(func(c *ui.Context) {
		ErrorOverlay(c, "Build failed", "x", &open, ErrorOverlayOptions{}, func() {
			c.Key("panel")
			ui.Text(c, "a")
			c.Key("scrim")
			ui.Text(c, "b")
		})
	})
}

func TestCodeTabsScrollToEveryTab(t *testing.T) {
	var files []CodeFile
	for _, n := range []string{"invoice_reminders.go", "invoice_reminders_test.go", "payment_retries.go", "ledger_snapshots.go"} {
		files = append(files, CodeFile{Name: n, Lang: "go", Code: "package billing // " + strings.TrimSuffix(n, ".go")})
	}
	current := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ui.Column(c).Width(320).Children(func() { CodeTabs(c, files, &current) })
	}), 600, 300)
	settles(t, tt, "tabs")
	first, last := find(t, tt, "invoice_reminders.go"), find(t, tt, "ledger_snapshots.go")
	if last.X+last.W <= 330 {
		t.Fatalf("tabs fit; the test needs overflow: %v", last)
	}
	tt.Scroll(first.X+first.W/2, first.Y+first.H/2, 600, 0)
	settles(t, tt, "tabs")
	last = find(t, tt, "ledger_snapshots.go")
	if last.X+last.W > 331 {
		t.Fatalf("last tab still cut off at %v", last.X+last.W)
	}
	click(t, tt, "ledger_snapshots.go")
	if current != 3 {
		t.Errorf("current %d", current)
	}
}
