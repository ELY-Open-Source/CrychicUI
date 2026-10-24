package code

import (
	"image/color"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestLineHighlight(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		LineHighlight(c, "go", "one\ntwo\nthree", []int{2})
	}), 500, 300)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("marked line does not pulse")
	}
	r := find(t, tt, "two")
	other := find(t, tt, "one")
	px := func(r ui.Rect) color.RGBA { return tt.Image().RGBAAt(int(r.X+r.W+8), int(r.Y+r.H/2)) }
	pulsing := px(r)
	settles(t, tt, "lines")
	if px(r) == pulsing {
		t.Error("the pulse never fades")
	}
	if px(r) == px(other) {
		t.Error("the marked line loses its tint after the pulse")
	}
	panics(t, "mark out of range", func() {
		render(func(c *ui.Context) { LineHighlight(c, "go", "one", []int{2}) })
	})
	panics(t, "duplicate mark", func() {
		render(func(c *ui.Context) { LineHighlight(c, "go", "one\ntwo", []int{1, 1}) })
	})
}

func TestCodeKeepsIndentation(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		LineHighlight(c, "go", "func f() {\n\treturn\n}", nil)
	}), 500, 300)
	settles(t, tt, "indent")
	if !tt.HasText("    return") {
		t.Errorf("tab not shown as spaces: %q", tt.Texts())
	}
}

func TestAnnotated(t *testing.T) {
	notes := []Note{{Line: 2, Text: "This runs first."}, {Line: 3, Text: "Then this."}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Annotated(c, "go", "one\ntwo\nthree", notes) }), 500, 400)
	settles(t, tt, "annotated")
	code, note := find(t, tt, "two"), find(t, tt, "This runs first.")
	if note.X <= code.X+code.W || note.Y > code.Y+code.H || note.Y+note.H < code.Y {
		t.Errorf("note not beside its line: note %v line %v", note, code)
	}
	panics(t, "unsorted notes", func() {
		render(func(c *ui.Context) {
			Annotated(c, "go", "one\ntwo", []Note{{Line: 2, Text: "a"}, {Line: 1, Text: "b"}})
		})
	})
	panics(t, "empty note", func() {
		render(func(c *ui.Context) { Annotated(c, "go", "one", []Note{{Line: 1}}) })
	})
}

func TestStackTrace(t *testing.T) {
	list := []StackFrame{
		{Func: "main.(*Server).handle", File: "server.go", Line: 42, App: true},
		{Func: "net/http.HandlerFunc.ServeHTTP", File: "server.go", Line: 2136},
		{Func: "net/http.(*conn).serve", File: "server.go", Line: 2009},
	}
	open, changed := false, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if StackTrace(c, "panic: nil map", list, &open, StackTraceOptions{}).Changed() {
			changed++
		}
	}), 500, 400)
	settles(t, tt, "trace")
	if !tt.HasText("panic: nil map") || !tt.HasText("(*Server).handle") || !tt.HasText("Show 2 library frames") || tt.HasText("(*conn).serve") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Show 2 library frames")
	if !open || changed != 1 || !tt.HasText("(*conn).serve") || !tt.HasText("Hide library frames") {
		t.Fatalf("open %v changed %d texts %q", open, changed, tt.Texts())
	}
	cues(t, before, sound.Toggle)
	click(t, tt, "Hide library frames")
	settles(t, tt, "trace")
	if open || tt.HasText("(*conn).serve") {
		t.Errorf("open %v texts %q", open, tt.Texts())
	}
	lib := []StackFrame{{Func: "runtime.goexit", File: "asm.s", Line: 1}, {Func: "runtime.main", File: "proc.go", Line: 9}}
	closed := false
	only := ui.NewTester(padded(func(c *ui.Context) { StackTrace(c, "fatal", lib, &closed, StackTraceOptions{}) }), 500, 300)
	settles(t, only, "library trace")
	if !only.HasText("goexit") || only.HasText("proc.go:9") || !only.HasText("Show 1 library frames") {
		t.Errorf("library-only texts %q", only.Texts())
	}
	panics(t, "bad frame", func() {
		render(func(c *ui.Context) { StackTrace(c, "x", []StackFrame{{Func: "f"}}, &open, StackTraceOptions{}) })
	})
}

func TestSplitFunc(t *testing.T) {
	for in, want := range map[string][2]string{
		"net/http.(*conn).serve": {"net/http.", "(*conn).serve"},
		"main.main":              {"main.", "main"},
		"goexit":                 {"", "goexit"},
	} {
		if pkg, name := splitFunc(in); pkg != want[0] || name != want[1] {
			t.Errorf("%q split %q %q", in, pkg, name)
		}
	}
}

func TestErrorOverlay(t *testing.T) {
	open, changed := true, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root := ErrorOverlay(c, "Build failed", "main.go:12:5: undefined: deck\n\tdeck.Ship()", &open, ErrorOverlayOptions{}, func() {
			ui.Text(c, "editor content")
		})
		if root.Changed() {
			changed++
		}
	}), 500, 300)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("panel does not slide down")
	}
	settles(t, tt, "overlay")
	if !tt.HasText("Build failed") || !tt.HasText("editor content") || !tt.HasText("    deck.Ship()") {
		t.Fatalf("texts %q", tt.Texts())
	}
	panel, detail := find(t, tt, "Build failed"), find(t, tt, "    deck.Ship()")
	if detail.Y+detail.H > 300 || detail.Y < panel.Y {
		t.Errorf("detail outside the box: %v", detail)
	}
	before := baseline()
	click(t, tt, "Close")
	if open || changed != 1 {
		t.Fatalf("open %v changed %d", open, changed)
	}
	cues(t, before, sound.Close)
	panics(t, "nil content", func() {
		render(func(c *ui.Context) { ErrorOverlay(c, "x", "", &open, ErrorOverlayOptions{}, nil) })
	})
}

func TestGutter(t *testing.T) {
	breakpoints := []int{}
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Gutter(c, "go", "a := 1\nb := 2\nc := a + b\nprint(c)\nreturn", &breakpoints, GutterOptions{}).Changed() {
			changed++
		}
	}), 400, 300)
	settles(t, tt, "gutter")
	if !tt.HasText("c := a + b") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Breakpoint on line 2")
	if len(breakpoints) != 1 || breakpoints[0] != 2 || changed != 1 {
		t.Fatalf("breakpoints %v changed %d", breakpoints, changed)
	}
	click(t, tt, "Breakpoint on line 2")
	if len(breakpoints) != 0 || changed != 2 {
		t.Errorf("breakpoints %v changed %d", breakpoints, changed)
	}
	cues(t, before, sound.Toggle, sound.Toggle)
	panics(t, "no code", func() { render(func(c *ui.Context) { Gutter(c, "go", "", &breakpoints, GutterOptions{}) }) })
	panics(t, "out of range", func() {
		b := []int{9}
		render(func(c *ui.Context) { Gutter(c, "go", "a\nb", &b, GutterOptions{}) })
	})
}

func TestDebugBar(t *testing.T) {
	running := false
	cont, stop, changed := 0, 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, cn, _, _, _, st := DebugBar(c, running, DebugBarOptions{})
		if root.Changed() {
			changed++
		}
		if cn.Clicked() {
			cont++
		}
		if st.Clicked() {
			stop++
		}
	}), 500, 200)
	settles(t, tt, "debug bar")
	if !tt.HasText("Paused") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Continue")
	click(t, tt, "Stop")
	if cont != 1 || stop != 1 || changed != 2 {
		t.Fatalf("cont %d stop %d changed %d", cont, stop, changed)
	}
	cues(t, before, sound.Tap, sound.Tap)
	running = true
	tt.Frame()
	if !tt.HasText("Running") || !moves(tt, 600*time.Millisecond) {
		t.Errorf("running pill still: %q", tt.Texts())
	}
	again := baseline()
	click(t, tt, "Continue")
	if cont != 1 {
		t.Error("disabled continue fired")
	}
	click(t, tt, "Stop")
	if stop != 2 || changed != 3 {
		t.Errorf("stop %d changed %d", stop, changed)
	}
	cues(t, again, sound.Tap)
}

func TestVariables(t *testing.T) {
	vars := []Variable{
		{Name: "cfg", Type: "Config", Value: "{…}", Children: []Variable{{Name: "port", Type: "int", Value: "8080"}}},
		{Name: "ready", Type: "bool", Value: "true"},
	}
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Variables(c, "Locals", vars).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "variables")
	if !tt.HasText("Locals") || !tt.HasText("port") || !tt.HasText("8080") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "cfg")
	settles(t, tt, "variables")
	if changed != 1 || tt.HasText("port") {
		t.Fatalf("changed %d texts %q", changed, tt.Texts())
	}
	cues(t, before, sound.Toggle)
	click(t, tt, "cfg")
	if !tt.HasText("port") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "empty name", func() {
		render(func(c *ui.Context) { Variables(c, "x", []Variable{{Name: ""}}) })
	})
	panics(t, "twin names", func() {
		render(func(c *ui.Context) { Variables(c, "x", []Variable{{Name: "a"}, {Name: "a"}}) })
	})
}

func TestValueInk(t *testing.T) {
	th := ui.LightTheme()
	for v, want := range map[string]ui.Color{`"hi"`: th.Success, "true": th.Warning, "nil": th.Warning, "-3.5": th.Accent, "Config{}": th.Text} {
		if got := valueInk(th, v); got != want {
			t.Errorf("%q ink %v, want %v", v, got, want)
		}
	}
}

func TestWatches(t *testing.T) {
	watches := []Watch{{Expr: "len(deck)", Value: "12"}}
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Watches(c, "Watch", &watches, WatchesOptions{}).Changed() {
			changed++
		}
	}), 500, 400)
	settles(t, tt, "watches")
	before := baseline()
	click(t, tt, "Watch an expression")
	tt.Type("x + 1")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if len(watches) != 2 || watches[1].Expr != "x + 1" || changed != 1 {
		t.Fatalf("watches %v changed %d", watches, changed)
	}
	if !tt.HasText("x + 1") || !tt.HasText("not available") {
		t.Errorf("texts %q", tt.Texts())
	}
	cues(t, before, sound.Pop)
	click(t, tt, "Watch an expression")
	tt.Type("x + 1")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if len(watches) != 2 {
		t.Errorf("duplicate added: %v", watches)
	}
	cues(t, before, sound.Pop)
	click(t, tt, "Remove len(deck)")
	settles(t, tt, "watches")
	click(t, tt, "Remove x + 1")
	settles(t, tt, "watches")
	if len(watches) != 0 || !tt.HasText("No expressions watched") {
		t.Errorf("watches %v texts %q", watches, tt.Texts())
	}
	cues(t, before, sound.Pop, sound.Pop, sound.Pop)
	panics(t, "empty expr", func() {
		w := []Watch{{Expr: ""}}
		render(func(c *ui.Context) { Watches(c, "x", &w, WatchesOptions{}) })
	})
}

func TestCallStack(t *testing.T) {
	list := []StackFrame{
		{Func: "main.(*Server).handle", File: "server.go", Line: 42, App: true},
		{Func: "net/http.(*conn).serve", File: "server.go", Line: 2009},
	}
	current, changed := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if CallStack(c, "Call stack", list, &current).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "stack")
	if !tt.HasText("Call stack") || !tt.HasText("net/http.") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "net/http.(*conn).serve")
	if current != 1 || changed != 1 {
		t.Fatalf("current %d changed %d", current, changed)
	}
	cues(t, before, sound.Select)
	panics(t, "bad current", func() {
		n := 5
		render(func(c *ui.Context) { CallStack(c, "x", list, &n) })
	})
}

func TestProblems(t *testing.T) {
	problems := []Problem{
		{Severity: ProbError, File: "main.go", Line: 12, Message: "undefined: deck"},
		{Severity: ProbWarning, File: "main.go", Line: 30, Message: "unused variable"},
		{Severity: ProbWarning, File: "util.go", Line: 4, Message: "missing return"},
	}
	selected, changed := -1, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Problems(c, "Problems", problems, &selected, ProblemsOptions{}).Changed() {
			changed++
		}
	}), 500, 400)
	settles(t, tt, "problems")
	if !tt.HasText("util.go") || !tt.HasText("undefined: deck") || !tt.HasText("Ln 30") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "missing return")
	if selected != 2 || changed != 1 {
		t.Fatalf("selected %d changed %d", selected, changed)
	}
	cues(t, before, sound.Select)
	panics(t, "zero line", func() {
		render(func(c *ui.Context) {
			Problems(c, "x", []Problem{{Severity: ProbError, File: "f", Message: "m"}}, &selected, ProblemsOptions{})
		})
	})
}

func TestOutline(t *testing.T) {
	symbols := []Symbol{
		{Kind: SymFunc, Name: "Serve", Line: 20},
		{Kind: SymType, Name: "Config", Line: 3},
		{Kind: SymVar, Name: "cache", Line: 41},
	}
	selected, changed := -1, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Outline(c, "Outline", symbols, &selected, OutlineOptions{}).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "outline")
	if !tt.HasText("Types") || !tt.HasText("Functions") || tt.HasText("Constants") {
		t.Fatalf("texts %q", tt.Texts())
	}
	if find(t, tt, "Config").Y > find(t, tt, "Serve").Y {
		t.Error("types do not lead")
	}
	before := baseline()
	click(t, tt, "Serve")
	if selected != 0 || changed != 1 {
		t.Fatalf("selected %d changed %d", selected, changed)
	}
	cues(t, before, sound.Select)
	panics(t, "unknown kind", func() {
		render(func(c *ui.Context) {
			Outline(c, "x", []Symbol{{Kind: 9, Name: "a", Line: 1}}, &selected, OutlineOptions{})
		})
	})
}
