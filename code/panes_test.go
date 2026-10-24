package code

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

func TestSymbolSearch(t *testing.T) {
	symbols := []Symbol{
		{SymFunc, "runServer", 12},
		{SymType, "Server", 30},
		{SymVar, "port", 5},
	}
	query, chosen := "", -1
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, _ := SymbolSearch(c, "Symbols", symbols, &query, &chosen, SymbolSearchOptions{})
		if root.Changed() {
			changed++
		}
	}), 500, 400)
	settles(t, tt, "symbols")
	tt.Key(0, ui.KeyTab)
	tt.Type("serv")
	tt.Frame()
	if tt.HasText("port") || !tt.HasText("runServer") || !tt.HasText("Server") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "runServer")
	if chosen != 0 {
		t.Fatalf("chosen %d", chosen)
	}
	cues(t, before, sound.Select)
	query = "zzz"
	tt.Frame()
	if !tt.HasText("No symbols match") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "no query", func() {
		render(func(c *ui.Context) { SymbolSearch(c, "x", symbols, nil, &chosen, SymbolSearchOptions{}) })
	})

	many := make([]Symbol, 13)
	for i := range many {
		many[i] = Symbol{Kind: SymFunc, Name: "foo", Line: i + 1}
	}
	many[2].Name = "foo1"
	collide := ui.NewTester(padded(func(c *ui.Context) {
		q, ch := "", -1
		SymbolSearch(c, "Symbols", many, &q, &ch, SymbolSearchOptions{})
	}), 500, 500)
	settles(t, collide, "collide")
	if !collide.HasText("foo1") {
		t.Errorf("colliding names %q", collide.Texts()[:4])
	}
}

func TestBranchPicker(t *testing.T) {
	branches := []Branch{{Name: "main"}, {Name: "dev", Ahead: 2}, {Name: "feature/review", Behind: 5}}
	current, query := 0, ""
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, _ := BranchPicker(c, "Branches", branches, &current, &query, BranchPickerOptions{})
		if root.Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "branches")
	tt.Key(0, ui.KeyTab)
	tt.Type("feat")
	tt.Frame()
	if tt.HasText("main") || !tt.HasText("feature/review") {
		t.Errorf("texts %q", tt.Texts())
	}
	tt.Key(ui.Ctrl, ui.KeyA)
	tt.Key(0, ui.KeyBackspace)
	tt.Frame()
	was := changed
	before := baseline()
	click(t, tt, "dev")
	if current != 1 || changed != was+1 {
		t.Fatalf("current %d changed %d", current, changed)
	}
	cues(t, before, sound.Select)
	panics(t, "bad current", func() {
		render(func(c *ui.Context) { BranchPicker(c, "x", branches, &[]int{9}[0], &query, BranchPickerOptions{}) })
	})
}

func TestKeyTable(t *testing.T) {
	groups := []KeyGroup{{
		Name: "Editor",
		Shortcuts: []Shortcut{
			{"Ctrl+S", "Save the file"},
			{"Ctrl+Shift+P", "Open the palette"},
		},
	}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		KeyTable(c, "Shortcuts", groups)
	}), 500, 300)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("groups do not rise")
	}
	settles(t, tt, "keys")
	if !tt.HasText("Ctrl") || !tt.HasText("Shift") || !tt.HasText("Save the file") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "no keys", func() {
		render(func(c *ui.Context) { KeyTable(c, "x", []KeyGroup{{"g", []Shortcut{{"", "a"}}}}) })
	})
}

func TestConsole(t *testing.T) {
	entries := []LogEntry{
		{Level: LevelLog, Text: "server started"},
		{Level: LevelWarn, Text: "slow query", Source: "db.go:88"},
		{Level: LevelError, Text: "write failed"},
	}
	filter := 7
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Console(c, "Console", entries, &filter, ConsoleOptions{}).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "console")
	if !tt.HasText("server started") || !tt.HasText("Error 1") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Log 1")
	if filter != 6 || changed != 1 {
		t.Fatalf("filter %d changed %d", filter, changed)
	}
	cues(t, before, sound.Toggle)
	settles(t, tt, "console")
	if tt.HasText("server started") || !tt.HasText("db.go:88") {
		t.Errorf("log row still shown %q", tt.Texts())
	}
	filter = 0
	tt.Frame()
	settles(t, tt, "console")
	if !tt.HasText("No messages") {
		t.Errorf("no empty state %q", tt.Texts())
	}
	panics(t, "bad level", func() {
		render(func(c *ui.Context) {
			Console(c, "x", []LogEntry{{Level: logLevelCount, Text: "y"}}, &filter, ConsoleOptions{})
		})
	})
}

func TestEnv(t *testing.T) {
	vars := []EnvVar{{"API_KEY", "secret-one"}, {"REGION", "us-east"}}
	shown := false
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Env(c, "Env", vars, &shown, EnvOptions{}).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "env")
	if tt.HasText("secret-one") || !tt.HasText("API_KEY") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Show values")
	if !shown || changed != 1 {
		t.Fatalf("shown %v changed %d", shown, changed)
	}
	cues(t, before, sound.Toggle)
	if !tt.HasText("secret-one") {
		t.Errorf("value still masked %q", tt.Texts())
	}
	panics(t, "empty key", func() {
		render(func(c *ui.Context) { Env(c, "x", []EnvVar{{"", "v"}}, &shown, EnvOptions{}) })
	})
}

func TestPlayground(t *testing.T) {
	buf := editor.NewBuffer("println(\"hi\")\n")
	buf.SetLanguage(editor.LanguageFor("main.go"))
	st := &editor.State{Buffer: buf}
	output := ""
	var ran []string
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, _ := Playground(c, "Playground", st, &output, func(code string) {
			ran = append(ran, code)
			output = "hi"
		}, PlaygroundOptions{})
		if root.Changed() {
			changed++
		}
	}), 500, 500)
	settles(t, tt, "playground")
	was := changed
	before := baseline()
	click(t, tt, "Code editor")
	tt.TypeKey(0, ui.KeyX, "x")
	tt.Frame()
	if changed != was+1 {
		t.Fatalf("edits reported %d changes", changed)
	}
	cues(t, before, sound.Tick)
	was = changed
	before = baseline()
	buf.Edit(editor.Edit{Start: 0, End: 0, Text: "// external"})
	tt.Frame()
	if changed != was+1 {
		t.Fatalf("programmatic edit reported %d changes", changed)
	}
	cues(t, before, sound.Tick)
	before = baseline()
	src := buf.String()
	was = changed
	click(t, tt, "Run")
	if len(ran) != 1 || ran[0] != src || changed != was+1 {
		t.Fatalf("ran %q changed %d", ran, changed)
	}
	cues(t, before, sound.Tap)
	if !tt.HasText("hi") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "no run", func() {
		render(func(c *ui.Context) { Playground(c, "Playground", st, &output, nil, PlaygroundOptions{}) })
	})

	calls := 0
	buf2 := editor.NewBuffer("abc")
	st2 := &editor.State{Buffer: buf2}
	st2.OnInput = func(ev ui.InputEvent) bool {
		calls++
		return false
	}
	native := ui.NewTester(padded(func(c *ui.Context) {
		CodeEditor(c, st2, CodeEditorOptions{})
	}), 500, 400)
	settles(t, native, "native")
	native.Key(0, ui.KeyTab)
	native.Key(0, ui.KeyRight)
	native.Frame()
	if calls == 0 || !native.HasText("Ln 1, Col 2") {
		t.Fatalf("native input broken: calls %d texts %q", calls, native.Texts())
	}

	was = changed
	before = baseline()
	r := find(t, tt, "Playground")
	tt.Move(r.X+r.W/2, r.Y+20)
	tt.Frame()
	buf.Edit(editor.Edit{Start: 0, End: 0, Text: "y"})
	tt.Frame()
	if changed != was+1 {
		t.Fatalf("pointer then programmatic reported %d changes", changed)
	}
	cues(t, before, sound.Tick)

	output = "1\n2\n3\n4\n5\n6\n7\n8\n9\n10\n11\n12"
	small := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(10).Children(func() {
			Playground(c, "Playground", st, &output, func(string) {}, PlaygroundOptions{Height: 100, OutputHeight: 60})
			ui.Text(c, "after")
		})
	}, 500, 300)
	settles(t, small, "small")
	after := find(t, small, "after")
	pane := find(t, small, output)
	if pane.Y+pane.H > after.Y {
		t.Errorf("output overflows: pane at %v, after at %v", pane, after)
	}
}

func TestEnvMaskKeepsLengthSecret(t *testing.T) {
	vars := []EnvVar{{"A", "x"}, {"B", "a much longer secret value"}}
	shown := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Env(c, "Env", vars, &shown, EnvOptions{})
	}), 500, 300)
	settles(t, tt, "env")
	if tt.HasText("a much longer secret value") || tt.HasText("x") {
		t.Errorf("values leak: %q", tt.Texts())
	}
	if !tt.HasText("••••••••") {
		t.Errorf("no mask: %q", tt.Texts())
	}
}
