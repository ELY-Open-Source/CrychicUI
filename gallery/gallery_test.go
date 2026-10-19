package gallery

import (
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func pingDemo(c *ui.Context, log func(string)) {
	if ui.Button(c, "Ping").Clicked() {
		log("pinged")
	}
}

var pingEntry = Entry{ID: 1, Name: "PingBox", Category: "layout", Summary: "Answers a ping.", Demo: pingDemo}

// withEntry registers pingEntry for one test.
func withEntry(t *testing.T) {
	t.Helper()
	old, oldCode := entries, codes["pingDemo"]
	entries, codes["pingDemo"] = []Entry{pingEntry}, "func pingDemo() {}\n"
	t.Cleanup(func() { entries, codes["pingDemo"] = old, oldCode })
}

func open(path string) (*App, *ui.Tester) {
	a := New(path)
	return a, ui.NewTester(a.View, 1280, 820)
}

func TestHome(t *testing.T) {
	withEntry(t)
	a, tt := open("/")
	for _, s := range []string{"Components", "Categories", "Sound cues", "Layout", "Geo", "1 components"} {
		if !tt.HasText(s) {
			t.Errorf("home lacks %q", s)
		}
	}
	if s := a.Snapshot(); !s.Built || s.Page != "/" {
		t.Errorf("snapshot %+v", s)
	}
	if err := tt.Click("Layout"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if a.Snapshot().Page != "/c/layout" || !tt.HasText("PingBox") {
		t.Errorf("page %q, texts %q", a.Snapshot().Page, tt.Texts())
	}
}

func TestComponentPage(t *testing.T) {
	withEntry(t)
	a, tt := open(pingEntry.Path())
	for _, s := range []string{"PingBox", "Answers a ping.", "0001 · layout", "Code"} {
		if !tt.HasText(s) {
			t.Errorf("page lacks %q", s)
		}
	}
	if got := a.codes[pingEntry.Path()].Buffer.String(); got != codes["pingDemo"] {
		t.Errorf("code = %q", got)
	}
	if err := tt.Click("Ping"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("PingBox: pinged") {
		t.Errorf("log strip lacks the event: %q", tt.Texts())
	}
}

func TestNotFound(t *testing.T) {
	withEntry(t)
	for _, p := range []string{"/c/nope", "/c/layout/nope", "/x"} {
		if _, tt := open(p); !tt.HasText("Not found") {
			t.Errorf("%s: no Not found", p)
		}
	}
}

func TestSearch(t *testing.T) {
	withEntry(t)
	a, tt := open("/")
	a.query = "ping"
	tt.Frame()
	if !tt.HasText("1 results") {
		t.Fatalf("texts %q", tt.Texts())
	}
	if err := tt.Click("PingBox"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if a.query != "" || a.Snapshot().Page != pingEntry.Path() {
		t.Errorf("query %q, page %q", a.query, a.Snapshot().Page)
	}
}

func TestAppearance(t *testing.T) {
	_, tt := open("/")
	px := func() ui.Color { p := tt.Image().RGBAAt(1200, 400); return ui.Color{R: p.R, G: p.G, B: p.B, A: 255} }
	if err := tt.Click("Dark"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if got, want := px(), ui.DarkTheme().Background; got.R != want.R || got.G != want.G || got.B != want.B {
		t.Errorf("dark background %v, want %v", got, want)
	}
	if err := tt.Click("Light"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if got, want := px(), ui.LightTheme().Background; got.R != want.R || got.G != want.G || got.B != want.B {
		t.Errorf("light background %v, want %v", got, want)
	}
}

func TestSoundSwitch(t *testing.T) {
	t.Cleanup(func() { sound.SetEnabled(false) })
	_, tt := open("/")
	if err := tt.Click("Sound"); err != nil {
		t.Fatal(err)
	}
	if !sound.Enabled() {
		t.Fatal("switch did not enable sound")
	}
	if err := tt.Click("Sound"); err != nil {
		t.Fatal(err)
	}
	if sound.Enabled() {
		t.Error("switch did not mute sound")
	}
}

func TestCueStrip(t *testing.T) {
	sound.Play(sound.Success)
	if _, tt := open("/"); !strings.Contains(strings.Join(tt.Texts(), " "), "success") {
		t.Errorf("cue strip lacks success: %q", tt.Texts())
	}
}

func TestSources(t *testing.T) {
	raw := "func demoR() string {\n\treturn `a\n}\nb`\n}"
	got := sources(fstest.MapFS{"demos_x.go": {Data: []byte(
		"package gallery\r\n\r\n// A doc.\r\nfunc demoA(c int) {\r\n\tx()\r\n}\r\n\r\nfunc demoB() { x() }\r\n\r\n" +
			"func (e Entry) M() {\r\n}\r\n\r\n" + strings.ReplaceAll(raw, "\n", "\r\n") + "\r\n")}})
	want := map[string]string{"demoA": "func demoA(c int) {\n\tx()\n}", "demoB": "func demoB() { x() }", "demoR": raw}
	if len(got) != len(want) {
		t.Errorf("sources = %q", got)
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s = %q, want %q", k, got[k], v)
		}
	}
	defer func() {
		if recover() == nil {
			t.Error("broken file did not panic")
		}
	}()
	sources(fstest.MapFS{"demos_y.go": {Data: []byte("package gallery\nfunc demoB() {\n")}})
}

func TestSlug(t *testing.T) {
	for in, want := range map[string]string{"Stack": "stack", "KPIGrid": "kpi-grid", "Base64Tool": "base64-tool"} {
		if got := (Entry{Name: in}).Slug(); got != want {
			t.Errorf("slug(%s) = %s, want %s", in, got, want)
		}
	}
}

func TestCategories(t *testing.T) {
	seen := map[string]bool{}
	for _, k := range Categories {
		icon.Must(k.Icon)
		if seen[k.Key] {
			t.Errorf("category %s twice", k.Key)
		}
		seen[k.Key] = true
	}
	if len(Categories) != 38 {
		t.Errorf("categories = %d, want 38", len(Categories))
	}
}

func TestEveryEntry(t *testing.T) {
	for i, e := range Entries() {
		if i > 0 && e.ID <= Entries()[i-1].ID {
			t.Errorf("%s: ID %d out of order", e.Name, e.ID)
		}
		if _, ok := category(e.Category); !ok {
			t.Errorf("%s: unknown category %q", e.Name, e.Category)
		}
		if !strings.Contains(e.Code(), "func "+funcName(e.Demo)+"(") {
			t.Errorf("%s: code panel lacks its demo", e.Name)
		}
		for _, dark := range []bool{false, true} {
			a, tt := open(e.Path())
			tt.SetDark(dark)
			tt.Advance(3 * time.Second)
			if a.Snapshot().Page != e.Path() || !tt.HasText(e.Name) || tt.FramePending() {
				t.Errorf("%s dark=%v: page %q, pending %v", e.Name, dark, a.Snapshot().Page, tt.FramePending())
			}
		}
	}
}
