package code

import (
	"image/color"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestLineHighlight(t *testing.T) {
	var surface color.RGBA
	tt := ui.NewTester(padded(func(c *ui.Context) {
		surface = rgba(c.Theme().Surface)
		LineHighlight(c, "go", "one\ntwo\nthree", []int{2})
	}), 500, 300)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("marked line does not pulse")
	}
	r := find(t, tt, "two")
	other := find(t, tt, "one")
	px := func(x, y float32) color.RGBA { return tt.Image().RGBAAt(int(x), int(y)) }
	if px(r.X+1, r.Y+r.H/2) == px(other.X+1, other.Y+other.H/2) {
		t.Error("marked and unmarked lines share a tint during the pulse")
	}
	settles(t, tt, "lines")
	if !at(tt, r.X+1, r.Y+r.H/2, surface) {
		t.Error("marked line keeps its tint after the pulse")
	}
	panics(t, "mark out of range", func() {
		render(func(c *ui.Context) { LineHighlight(c, "go", "one", []int{2}) })
	})
	panics(t, "duplicate mark", func() {
		render(func(c *ui.Context) { LineHighlight(c, "go", "one\ntwo", []int{1, 1}) })
	})
}

func TestAnnotated(t *testing.T) {
	notes := []Note{{Line: 2, Text: "This runs first."}, {Line: 3, Text: "Then this."}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Annotated(c, "go", "one\ntwo\nthree", notes) }), 500, 300)
	settles(t, tt, "annotated")
	if !tt.HasText("This runs first.") || !tt.HasText("2") {
		t.Errorf("texts %q", tt.Texts())
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
		{Func: "main.main", File: "main.go", Line: 12, App: true},
		{Func: "http.serve", File: "server.go", Line: 99},
		{Func: "runtime.goexit", File: "asm.s", Line: 1},
	}
	open, changed := false, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if StackTrace(c, list, &open, StackTraceOptions{}).Changed() {
			changed++
		}
	}), 500, 400)
	settles(t, tt, "trace")
	if !tt.HasText("2 more frames") || tt.HasText("http.serve") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "2 more frames")
	if !open || changed != 1 || !tt.HasText("http.serve") || !tt.HasText("Show less") {
		t.Fatalf("open %v changed %d texts %q", open, changed, tt.Texts())
	}
	cues(t, before, sound.Toggle)
	click(t, tt, "Show less")
	if open || tt.HasText("http.serve") {
		t.Errorf("open %v texts %q", open, tt.Texts())
	}
	panics(t, "bad frame", func() {
		render(func(c *ui.Context) { StackTrace(c, []StackFrame{{Func: "f"}}, &open, StackTraceOptions{}) })
	})
}

func TestErrorOverlay(t *testing.T) {
	open, changed := true, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root := ErrorOverlay(c, "build failed", "main.go:12: undefined: deck", &open, ErrorOverlayOptions{}, func() {
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
	if !tt.HasText("build failed") || !tt.HasText("editor content") {
		t.Fatalf("texts %q", tt.Texts())
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
		if Gutter(c, 5, &breakpoints).Changed() {
			changed++
		}
	}), 300, 300)
	settles(t, tt, "gutter")
	before := baseline()
	click(t, tt, "2")
	if len(breakpoints) != 1 || breakpoints[0] != 2 || changed != 1 {
		t.Fatalf("breakpoints %v changed %d", breakpoints, changed)
	}
	click(t, tt, "2")
	if len(breakpoints) != 0 || changed != 2 {
		t.Errorf("breakpoints %v changed %d", breakpoints, changed)
	}
	cues(t, before, sound.Toggle, sound.Toggle)
	panics(t, "zero lines", func() { render(func(c *ui.Context) { Gutter(c, 0, &breakpoints) }) })
	panics(t, "out of range", func() {
		b := []int{9}
		render(func(c *ui.Context) { Gutter(c, 5, &b) })
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
	before := baseline()
	click(t, tt, "Continue")
	if cont != 1 || changed != 1 {
		t.Fatalf("cont %d changed %d", cont, changed)
	}
	click(t, tt, "Stop")
	if stop != 0 || changed != 1 {
		t.Errorf("disabled stop fired: stop %d changed %d", stop, changed)
	}
	cues(t, before, sound.Tap)
	running = true
	tt.Frame()
	again := baseline()
	click(t, tt, "Continue")
	if cont != 1 {
		t.Error("disabled continue fired")
	}
	click(t, tt, "Stop")
	if stop != 1 || changed != 2 {
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
		if Variables(c, "Variables", vars).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "variables")
	if !tt.HasText("port") || !tt.HasText("8080") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "cfg")
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
}

func TestWatches(t *testing.T) {
	watches := []Watch{{Expr: "len(deck)", Value: "12"}}
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Watches(c, "Watches", &watches, WatchesOptions{}).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "watches")
	before := baseline()
	click(t, tt, "Watch an expression")
	tt.Type("x + 1")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if len(watches) != 2 || watches[1].Expr != "x + 1" || changed != 1 {
		t.Fatalf("watches %v changed %d", watches, changed)
	}
	if !tt.HasText("x + 1") || !tt.HasText("—") {
		t.Errorf("texts %q", tt.Texts())
	}
	cues(t, before, sound.Pop)
	click(t, tt, "Remove")
	if len(watches) != 1 || changed != 2 {
		t.Errorf("watches %v changed %d", watches, changed)
	}
	cues(t, before, sound.Pop, sound.Pop)
	panics(t, "empty expr", func() {
		w := []Watch{{Expr: ""}}
		render(func(c *ui.Context) { Watches(c, "x", &w, WatchesOptions{}) })
	})
}

func TestCallStack(t *testing.T) {
	list := []StackFrame{
		{Func: "main.main", File: "main.go", Line: 12, App: true},
		{Func: "http.serve", File: "server.go", Line: 99},
	}
	current, changed := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if CallStack(c, list, &current).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "stack")
	before := baseline()
	click(t, tt, "http.serve")
	if current != 1 || changed != 1 {
		t.Fatalf("current %d changed %d", current, changed)
	}
	cues(t, before, sound.Select)
	panics(t, "bad current", func() {
		n := 5
		render(func(c *ui.Context) { CallStack(c, list, &n) })
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
		if Problems(c, "Problems", problems, &selected).Changed() {
			changed++
		}
	}), 500, 400)
	settles(t, tt, "problems")
	if !tt.HasText("util.go") || !tt.HasText("undefined: deck") {
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
			Problems(c, "x", []Problem{{Severity: ProbError, File: "f", Message: "m"}}, &selected)
		})
	})
}

func TestOutline(t *testing.T) {
	symbols := []Symbol{
		{Kind: SymType, Name: "Config", Line: 3},
		{Kind: SymFunc, Name: "Serve", Line: 20},
		{Kind: SymVar, Name: "cache", Line: 41},
	}
	selected, changed := -1, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Outline(c, "Outline", symbols, &selected).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "outline")
	if !tt.HasText("f") || !tt.HasText("T") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Serve")
	if selected != 1 || changed != 1 {
		t.Fatalf("selected %d changed %d", selected, changed)
	}
	cues(t, before, sound.Select)
	panics(t, "unknown kind", func() {
		render(func(c *ui.Context) { Outline(c, "x", []Symbol{{Kind: 9, Name: "a", Line: 1}}, &selected) })
	})
}

func TestMinimap(t *testing.T) {
	at, changed := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Minimap(c, "Editor map", 200, 20, &at, MinimapOptions{}).Changed() {
			changed++
		}
	}), 300, 300)
	settles(t, tt, "minimap")
	r := find(t, tt, "Editor map")
	before := baseline()
	tt.Press(r.X+r.W/2, r.Y+r.H-6)
	tt.Frame()
	if at < 150 || changed == 0 {
		t.Fatalf("at %d changed %d after a press", at, changed)
	}
	tt.Move(r.X+r.W/2, r.Y+r.H-40)
	tt.Frame()
	tt.Release(r.X+r.W/2, r.Y+r.H-40)
	tt.Frame()
	if at >= 190 || at == 0 {
		t.Errorf("at %d after a drag", at)
	}
	if len(sound.Recent()) == 0 {
		t.Error("no cue after the drag")
	}
	cues(t, before, sound.Tick)
	keys := ui.NewTester(padded(func(c *ui.Context) {
		if Minimap(c, "Keys", 200, 20, &at, MinimapOptions{}).Changed() {
			changed++
		}
	}), 300, 300)
	keys.Frame()
	keys.Key(0, ui.KeyTab)
	if !keys.Focused("Keys") {
		t.Fatal("minimap does not take the focus")
	}
	at = 0
	keys.Key(0, ui.KeyEnd)
	keys.Frame()
	if at != 180 {
		t.Errorf("at %d after End, want 180", at)
	}
	panics(t, "window over lines", func() {
		n := 0
		render(func(c *ui.Context) { Minimap(c, "x", 10, 20, &n, MinimapOptions{}) })
	})
	panics(t, "start out of range", func() {
		n := 200
		render(func(c *ui.Context) { Minimap(c, "x", 200, 20, &n, MinimapOptions{}) })
	})
}

func TestMinimapSlowDrag(t *testing.T) {
	at := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Minimap(c, "Slow map", 10, 2, &at, MinimapOptions{Height: 160})
	}), 300, 300)
	settles(t, tt, "slow map")
	r := find(t, tt, "Slow map")
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Frame()
	for i := 1; i <= 64; i++ {
		tt.Move(x, y+float32(i))
		tt.Frame()
	}
	tt.Release(x, y+64)
	tt.Frame()
	if at != 8 {
		t.Errorf("at %d after slow drag, want 8", at)
	}
}

func TestMinimapUnscrollable(t *testing.T) {
	at, changed := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Minimap(c, "Whole file", 2, 2, &at, MinimapOptions{}).Changed() {
			changed++
		}
	}), 300, 300)
	settles(t, tt, "whole file")
	before := baseline()
	click(t, tt, "Whole file")
	if at != 0 || changed != 0 {
		t.Errorf("at %d changed %d without scroll range", at, changed)
	}
	cues(t, before)
}

func TestMinimapPressMotion(t *testing.T) {
	at := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Minimap(c, "Motion map", 20, 2, &at, MinimapOptions{})
	}), 300, 300)
	settles(t, tt, "motion map")
	r := find(t, tt, "Motion map")
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	if !tt.FramePending() {
		t.Error("press has no motion")
	}
	tt.Release(r.X+r.W/2, r.Y+r.H/2)
	settles(t, tt, "motion map")
}
