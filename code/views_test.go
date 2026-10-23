package code

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

func TestCodeBlock(t *testing.T) {
	code := "func main() {\n\tprintln(\"hi\")\n}"
	copies := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, cp := CodeBlock(c, "go", code, CodeBlockOptions{})
		if cp.Clicked() {
			copies++
		}
	}), 500, 300)
	settles(t, tt, "block")
	if !tt.HasText("func main() {") || !tt.HasText("go") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Copy")
	if copies != 1 || tt.Clipboard() != code {
		t.Errorf("copies %d clipboard %q", copies, tt.Clipboard())
	}
	cues(t, before, sound.Pop)
	panics(t, "empty code", func() { render(func(c *ui.Context) { CodeBlock(c, "go", "", CodeBlockOptions{}) }) })
}

func TestCodeEditor(t *testing.T) {
	buf := editor.NewBuffer("package main\n\n")
	buf.SetLanguage(editor.LanguageFor("main.go"))
	st := &editor.State{Buffer: buf}
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if CodeEditor(c, st, CodeEditorOptions{}).Changed() {
			changed++
		}
	}), 500, 400)
	settles(t, tt, "editor")
	if !tt.HasText("Go") || !tt.HasText("Ln 1, Col 1") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.TypeKey(0, ui.KeyX, "x")
	tt.Frame()
	if changed == 0 || !tt.HasText("Ln 1, Col 2") {
		t.Errorf("changed %d texts %q", changed, tt.Texts())
	}
	cues(t, before, sound.Tick)
	changed = 0
	recent := baseline()
	buf.Edit(editor.Edit{Start: 0, End: 0, Text: "// external\n"})
	tt.Frame()
	if changed != 1 {
		t.Errorf("a programmatic edit reported %d changes", changed)
	}
	cues(t, recent, sound.Tick)

	changed = 0
	recent = baseline()
	tt.KeyDown(ui.KeyRight)
	buf.Edit(editor.Edit{Start: 0, End: 0, Text: "z"})
	tt.Frame()
	tt.KeyUp(ui.KeyRight)
	if changed != 1 {
		t.Errorf("navigation then programmatic reported %d changes", changed)
	}
	cues(t, recent, sound.Tick)

	calls := 0
	per := ui.NewTester(padded(func(c *ui.Context) {
		st.OnInput = func(ev ui.InputEvent) bool {
			calls++
			return false
		}
		CodeEditor(c, st, CodeEditorOptions{})
	}), 500, 400)
	per.Frame()
	per.Key(0, ui.KeyTab)
	per.TypeKey(0, ui.KeyX, "x")
	per.Frame()
	if calls == 0 {
		t.Fatalf("caller hook lost: calls %d", calls)
	}
	cur := st
	swap := ui.NewTester(padded(func(c *ui.Context) {
		CodeEditor(c, cur, CodeEditorOptions{})
	}), 500, 400)
	swap.Frame()
	cur = &editor.State{Buffer: editor.NewBuffer("new")}
	swap.Key(0, ui.KeyTab)
	swap.TypeKey(0, ui.KeyX, "q")
	swap.Frame()
	if !swap.HasText("Ln 1, Col 2") {
		t.Fatalf("replacement state texts %q", swap.Texts())
	}
	panics(t, "nil state", func() { render(func(c *ui.Context) { CodeEditor(c, nil, CodeEditorOptions{}) }) })
}

func TestCodeEditorCompositionAndMenu(t *testing.T) {
	st := &editor.State{Buffer: editor.NewBuffer("")}
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if CodeEditor(c, st, CodeEditorOptions{}).Changed() {
			changed++
		}
	}), 500, 400)
	tt.Frame()
	tt.Key(0, ui.KeyTab)
	before := baseline()
	tt.Compose("に", 1)
	tt.Frame()
	tt.Type("に")
	tt.Frame()
	if changed == 0 || st.Buffer.String() != "に" {
		t.Fatalf("composition changed %d text %q", changed, st.Buffer.String())
	}
	if recent := sound.Recent(); recent[len(recent)-1] != sound.Tick {
		t.Errorf("composition cues %v", recent)
	}

	changed = 0
	before = baseline()
	tt.SetClipboard("menu")
	if err := tt.RightClick("Code editor"); err != nil {
		t.Fatal(err)
	}
	if err := tt.ChooseMenuItem("Paste"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if changed == 0 || st.Buffer.String() != "にmenu" {
		t.Errorf("menu paste changed %d text %q", changed, st.Buffer.String())
	}
	cues(t, before, sound.Tick)
}

func TestCodeEditorStateCallbacks(t *testing.T) {
	a := &editor.State{Buffer: editor.NewBuffer("a")}
	b := &editor.State{Buffer: editor.NewBuffer("b")}
	aCalls, bCalls := 0, 0
	a.OnInput = func(ui.InputEvent) bool { aCalls++; return false }
	b.OnInput = func(ui.InputEvent) bool { bCalls++; return false }
	cur := a
	tt := ui.NewTester(padded(func(c *ui.Context) { CodeEditor(c, cur, CodeEditorOptions{}) }), 500, 400)
	tt.Frame()
	cur = b
	tt.Frame()
	cur = a
	tt.Frame()
	tt.Key(0, ui.KeyTab)
	tt.TypeKey(0, ui.KeyX, "x")
	tt.Frame()
	if aCalls == 0 || bCalls != 0 {
		t.Errorf("callbacks A %d B %d", aCalls, bCalls)
	}
}

func TestDiffView(t *testing.T) {
	oldBuf, newBuf := editor.NewBuffer("a\nb\nc\n"), editor.NewBuffer("a\nB\nc\n")
	oldBuf.SetLanguage(editor.LanguageFor("main.go"))
	newBuf.SetLanguage(editor.LanguageFor("main.go"))
	st := &editor.DiffState{Old: oldBuf, New: newBuf}
	split, changed := true, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if DiffView(c, st, &split, DiffViewOptions{}).Changed() {
			changed++
		}
	}), 600, 400)
	settles(t, tt, "diff")
	before := baseline()
	click(t, tt, "Unified")
	if split || st.Layout != editor.DiffUnified || changed != 1 {
		t.Fatalf("split %v layout %v changed %d", split, st.Layout, changed)
	}
	click(t, tt, "Split")
	if !split || st.Layout != editor.DiffSplit || changed != 2 {
		t.Errorf("split %v layout %v changed %d", split, st.Layout, changed)
	}
	cues(t, before, sound.Select, sound.Select)
	panics(t, "nil flag", func() { render(func(c *ui.Context) { DiffView(c, st, nil, DiffViewOptions{}) }) })
}

func TestTerminal(t *testing.T) {
	lines := []string{"building…", "done in 1.2s"}
	tt := ui.NewTester(padded(func(c *ui.Context) { Terminal(c, "Shell", lines, TerminalOptions{}) }), 500, 300)
	tt.Frame()
	if !moves(tt, 650*time.Millisecond) {
		t.Error("caret does not blink")
	}
	if !tt.HasText("done in 1.2s") || !tt.HasText("$") {
		t.Errorf("texts %q", tt.Texts())
	}
	lines = append(lines, "test ok")
	tt.Frame()
	if !tt.HasText("test ok") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "no label", func() { render(func(c *ui.Context) { Terminal(c, "", nil, TerminalOptions{}) }) })
}

func TestCommandLine(t *testing.T) {
	history := []Command{{Text: "ls", Output: "deck.md"}}
	draft := ""
	var ran []string
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, _ := CommandLine(c, "Shell", history, &draft, func(s string) { ran = append(ran, s) }, CommandLineOptions{})
		if root.Changed() {
			changed++
		}
	}), 500, 400)
	settles(t, tt, "prompt")
	if !tt.HasText("deck.md") || !tt.HasText("$ ls") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Run a command")
	tt.Type("pwd")
	tt.Frame()
	if changed != 1 {
		t.Fatalf("typing reported %d changes", changed)
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if len(ran) != 1 || ran[0] != "pwd" || draft != "" || changed != 2 {
		t.Fatalf("ran %q draft %q changed %d", ran, draft, changed)
	}
	cues(t, before, sound.Send)
	tt.Key(0, ui.KeyUp)
	tt.Frame()
	if draft != "ls" || changed != 3 {
		t.Errorf("draft %q changed %d after Up", draft, changed)
	}
	tt.Key(0, ui.KeyDown)
	tt.Frame()
	if draft != "" {
		t.Errorf("draft %q after Down", draft)
	}
	history = nil
	tt.Frame()
	click(t, tt, "Run a command")
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyDown)
	tt.Frame()
	cues(t, before, sound.Send, sound.Tick, sound.Tick)
	panics(t, "nil run", func() { render(func(c *ui.Context) { CommandLine(c, "x", nil, &draft, nil, CommandLineOptions{}) }) })
}

func TestCodeTabs(t *testing.T) {
	files := []CodeFile{
		{Name: "main.go", Lang: "go", Code: "package main"},
		{Name: "util.go", Lang: "go", Code: "package util"},
	}
	current, changed := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if CodeTabs(c, files, &current).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "tabs")
	if !tt.HasText("package main") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "util.go")
	if current != 1 || changed != 1 || !tt.HasText("package util") {
		t.Fatalf("current %d changed %d texts %q", current, changed, tt.Texts())
	}
	cues(t, before, sound.Select)
	panics(t, "bad tab", func() {
		n := 2
		render(func(c *ui.Context) { CodeTabs(c, files, &n) })
	})
}

func TestFileHeader(t *testing.T) {
	copies := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, cp := FileHeader(c, "cmd/gallery/main.go", "Go", FileHeaderOptions{})
		if cp.Clicked() {
			copies++
		}
	}), 500, 200)
	settles(t, tt, "header")
	before := baseline()
	click(t, tt, "Copy path")
	if copies != 1 || tt.Clipboard() != "cmd/gallery/main.go" {
		t.Errorf("copies %d clipboard %q", copies, tt.Clipboard())
	}
	cues(t, before, sound.Pop)
	panics(t, "empty path", func() { render(func(c *ui.Context) { FileHeader(c, "", "Go", FileHeaderOptions{}) }) })
}

func TestSnippet(t *testing.T) {
	code := "http.ListenAndServe(\":8080\", nil)"
	copies := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, cp := Snippet(c, "Serve the gallery", "go", code, SnippetOptions{})
		if cp.Clicked() {
			copies++
		}
	}), 500, 300)
	settles(t, tt, "snippet")
	before := baseline()
	click(t, tt, "Copy")
	if copies != 1 || tt.Clipboard() != code {
		t.Errorf("copies %d clipboard %q", copies, tt.Clipboard())
	}
	cues(t, before, sound.Pop)
	panics(t, "empty title", func() { render(func(c *ui.Context) { Snippet(c, "", "go", "x", SnippetOptions{}) }) })
}
