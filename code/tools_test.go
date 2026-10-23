package code

import (
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestRegexTester(t *testing.T) {
	pattern, text := "a+", "banana"
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if RegexTester(c, "Rx", &pattern, &text, RegexTesterOptions{}).Changed() {
			changed++
		}
	}), 500, 400)
	settles(t, tt, "regex")
	if !tt.HasText("3 matches") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Pattern")
	tt.Key(0, ui.KeyHome)
	tt.TypeKey(0, ui.KeyX, "x")
	tt.Frame()
	if pattern != "xa+" || changed != 1 {
		t.Fatalf("pattern %q changed %d", pattern, changed)
	}
	cues(t, before, sound.Tick)
	click(t, tt, "Pattern")
	tt.Key(ui.Ctrl, ui.KeyA)
	tt.Type("([")
	tt.Frame()
	if !tt.HasText("bad pattern: error parsing regexp: missing closing ]: `[`") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "no fields", func() {
		render(func(c *ui.Context) { RegexTester(c, "x", nil, nil, RegexTesterOptions{}) })
	})

	pattern, text = "foo|bar", "foo\nbar"
	tt.Frame()
	settles(t, tt, "regex")
	foo, bar := find(t, tt, "foo"), find(t, tt, "bar")
	if foo.Y == bar.Y {
		t.Errorf("matches share a line: %v %v", foo, bar)
	}

	custom := ui.NewTester(padded(func(c *ui.Context) {
		RegexTester(c, "Rx", &pattern, &text, RegexTesterOptions{Pattern: "Motif", Matches: "%d Treffer"})
	}), 500, 400)
	settles(t, custom, "custom")
	if !custom.HasText("Motif") || !custom.HasText("2 Treffer") {
		t.Errorf("custom copy %q", custom.Texts())
	}

	pattern, text = "(?s)foo.*bar", "foo\nbar"
	tt.Frame()
	settles(t, tt, "regex")
	foos := 0
	for _, s := range tt.Texts() {
		if strings.Contains(s, "foo") {
			foos++
		}
	}
	if foos != 1 {
		t.Errorf("foo appears %d times: %q", foos, tt.Texts())
	}

	pattern, text = "foo|bar", "foo\n\nbar"
	tt.Frame()
	settles(t, tt, "regex")
	foo2, bar2 := find(t, tt, "foo"), find(t, tt, "bar")
	if bar2.Y-foo2.Y < 2*foo2.H {
		t.Errorf("blank line collapsed: foo %v bar %v", foo2, bar2)
	}

	pattern, text = "bar", "foo bar"
	tt.Frame()
	if !tt.HasText("foo bar") {
		t.Errorf("space beside match lost: %q", tt.Texts())
	}

	errCustom := ui.NewTester(padded(func(c *ui.Context) {
		bad := "["
		RegexTester(c, "Rx", &bad, &text, RegexTesterOptions{FormatError: func(error) string { return "Motif invalide" }})
	}), 500, 300)
	settles(t, errCustom, "err custom")
	if !errCustom.HasText("Motif invalide") {
		t.Errorf("custom error %q", errCustom.Texts())
	}
}

func TestJSONFormat(t *testing.T) {
	text := `{"a":1}`
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, _, _ := JSONFormat(c, "JSON", &text, JSONFormatOptions{})
		if root.Changed() {
			changed++
		}
	}), 500, 400)
	settles(t, tt, "json")
	before := baseline()
	click(t, tt, "Format")
	if !strings.Contains(text, "\n  \"a\": 1\n") || changed != 1 {
		t.Fatalf("text %q changed %d", text, changed)
	}
	cues(t, before, sound.Success)
	click(t, tt, "Minify")
	if text != `{"a":1}` {
		t.Fatalf("text %q", text)
	}
	text = "{bad"
	before = baseline()
	click(t, tt, "Format")
	cues(t, before, sound.Error)
	if !tt.HasText("invalid character 'b' looking for beginning of object key string") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "no text", func() {
		render(func(c *ui.Context) { JSONFormat(c, "x", nil, JSONFormatOptions{}) })
	})

	text = `{"id":9007199254740993}`
	click(t, tt, "Minify")
	if text != `{"id":9007199254740993}` {
		t.Fatalf("numbers not preserved: %q", text)
	}

	text = `{"a":1,"a":2}`
	click(t, tt, "Minify")
	if text != `{"a":1,"a":2}` {
		t.Fatalf("duplicate members not preserved: %q", text)
	}

	text = `{"a":1} {"b":2}`
	before = baseline()
	click(t, tt, "Format")
	if text != `{"a":1} {"b":2}` {
		t.Fatalf("trailing content kept: %q", text)
	}
	cues(t, before, sound.Error)
	if !tt.HasText("invalid JSON: more than one JSON value") {
		t.Errorf("texts %q", tt.Texts())
	}

	custom := ui.NewTester(padded(func(c *ui.Context) {
		bad := "{"
		JSONFormat(c, "JSON", &bad, JSONFormatOptions{FormatError: func(error) string { return "JSON invalide" }})
	}), 500, 300)
	settles(t, custom, "custom json")
	click(t, custom, "Format")
	if !custom.HasText("JSON invalide") {
		t.Errorf("custom error %q", custom.Texts())
	}
}

func TestBase64Tool(t *testing.T) {
	text, encode := "hi", true
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Base64Tool(c, "Base64", &text, &encode, Base64ToolOptions{}).Changed() {
			changed++
		}
	}), 500, 400)
	settles(t, tt, "base64")
	if !tt.HasText("aGk=") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Decode")
	if encode || changed != 1 {
		t.Fatalf("encode %v changed %d", encode, changed)
	}
	cues(t, before, sound.Tick)
	if !tt.HasText("invalid base64") {
		t.Errorf("texts %q", tt.Texts())
	}
	text = "aGk="
	tt.Frame()
	if !tt.HasText("hi") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "no encode", func() {
		render(func(c *ui.Context) { Base64Tool(c, "x", &text, nil, Base64ToolOptions{}) })
	})
}

func TestHashTool(t *testing.T) {
	text := "abc"
	tt := ui.NewTester(padded(func(c *ui.Context) {
		HashTool(c, "Hash", &text)
	}), 500, 300)
	settles(t, tt, "hash")
	if !tt.HasText("ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad") {
		t.Errorf("texts %q", tt.Texts())
	}
	text = "abcd"
	tt.Frame()
	if !moves(tt, 300*time.Millisecond) {
		t.Error("digest does not cross-fade")
	}
	settles(t, tt, "hash")
	panics(t, "no text", func() {
		render(func(c *ui.Context) { HashTool(c, "x", nil) })
	})
}

func TestUnixTime(t *testing.T) {
	epoch := int64(1700000000)
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if UnixTime(c, "Unix", &epoch, UnixTimeOptions{}).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "unix")

	before := baseline()
	click(t, tt, "Unix seconds")
	tt.Key(ui.Ctrl, ui.KeyA)
	tt.Type("1704067200")
	tt.Frame()
	if epoch != 1704067200 || changed != 1 {
		t.Fatalf("epoch %d changed %d", epoch, changed)
	}
	cues(t, before, sound.Tick)
	epoch = 1735689600
	tt.Frame()
	click(t, tt, "Unix seconds")
	tt.Key(ui.Ctrl, ui.KeyA)
	tt.Type("1700000000")
	tt.Frame()
	if epoch != 1700000000 {
		t.Fatalf("epoch after external update %d", epoch)
	}
	click(t, tt, "Unix seconds")
	tt.Key(ui.Ctrl, ui.KeyA)
	tt.Type("x")
	tt.Frame()
	if !tt.HasText("not a whole number") {
		t.Fatalf("no error %q", tt.Texts())
	}
	epoch = 1704067200
	tt.Frame()
	tt.Frame()
	if tt.HasText("not a whole number") {
		t.Fatalf("stale error %q", tt.Texts())
	}
	panics(t, "no epoch", func() {
		render(func(c *ui.Context) { UnixTime(c, "x", nil, UnixTimeOptions{}) })
	})
}

func TestPackageCard(t *testing.T) {
	changed, copied := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, copy := PackageCard(c, "github.com/egoist/mygo", "v0.5.5", 12400, PackageCardOptions{})
		if root.Changed() {
			changed++
		}
		if copy.Clicked() {
			copied++
		}
	}), 500, 300)
	settles(t, tt, "package")
	if !tt.HasText("12.4k downloads") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Copy install")
	if tt.Clipboard() != "go get github.com/egoist/mygo@v0.5.5" || copied != 1 {
		t.Fatalf("clipboard %q copied %d", tt.Clipboard(), copied)
	}
	cues(t, before, sound.Pop)
	panics(t, "empty version", func() {
		render(func(c *ui.Context) { PackageCard(c, "x", "", 1, PackageCardOptions{}) })
	})

	custom := ui.NewTester(padded(func(c *ui.Context) {
		PackageCard(c, "pkg", "1.0", 5, PackageCardOptions{Install: "install %s@%s"})
	}), 500, 300)
	settles(t, custom, "custom card")
	if !custom.HasText("install pkg@1.0") {
		t.Errorf("custom install %q", custom.Texts())
	}
}

func TestREPLStaleDraft(t *testing.T) {
	history := []REPLLine{}
	draft := ""
	tt := ui.NewTester(padded(func(c *ui.Context) {
		REPL(c, "REPL", &history, &draft, REPLOptions{})
	}), 500, 400)
	settles(t, tt, "stale repl")
	click(t, tt, "1 + 2 * 3")
	tt.Type("2+2")
	draft = ""
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if len(history) != 1 || history[0].In != "2+2" || history[0].Out != "4" {
		t.Fatalf("stale draft discarded input: %v", history)
	}
}

func TestREPL(t *testing.T) {
	history := []REPLLine{{In: "1+1", Out: "2"}}
	draft := ""
	changed, runs := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, send := REPL(c, "REPL", &history, &draft, REPLOptions{})
		if root.Changed() {
			changed++
		}
		if send.Clicked() {
			runs++
		}
	}), 500, 400)
	settles(t, tt, "repl")
	if !tt.HasText("> 1+1") || !tt.HasText("= 2") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "1 + 2 * 3")
	tt.Type("2*(3+4)")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if len(history) != 2 || history[1].Out != "14" || changed != 2 {
		t.Fatalf("history %v changed %d", history, changed)
	}
	cues(t, before, sound.Send)
	click(t, tt, "1 + 2 * 3")
	tt.Type("1/0")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if len(history) != 3 || !history[2].Err || history[2].Out != "error: division by zero" {
		t.Fatalf("history %v", history)
	}
	panics(t, "no history", func() {
		render(func(c *ui.Context) { REPL(c, "x", nil, &draft, REPLOptions{}) })
	})

	h := []REPLLine{}
	d := ""
	errCustom := ui.NewTester(padded(func(c *ui.Context) {
		REPL(c, "REPL", &h, &d, REPLOptions{FormatError: func(error) string { return "Erreur personnalisee" }})
	}), 500, 400)
	settles(t, errCustom, "err repl")
	click(t, errCustom, "1 + 2 * 3")
	tt2 := errCustom
	tt2.Type("1/0")
	tt2.Key(0, ui.KeyEnter)
	tt2.Frame()
	if !errCustom.HasText("Erreur personnalisee") {
		t.Errorf("custom error %q", errCustom.Texts())
	}

	click(t, tt, "1 + 2 * 3")
	tt.Type("1.2.3+1")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if len(history) != 4 || !history[3].Err || history[3].Out != `error: bad number "1.2.3"` {
		t.Fatalf("history %v", history)
	}
	click(t, tt, "1 + 2 * 3")
	tt.Type("0.0000001*2")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if len(history) != 5 || history[4].Out != "2e-07" || history[4].Err {
		t.Fatalf("history %v", history)
	}
}
