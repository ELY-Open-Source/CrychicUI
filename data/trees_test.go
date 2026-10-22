package data

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestCheckTree(t *testing.T) {
	nodes := []Node{{"Fruit", []Node{{"Apple", nil}, {"Pear", nil}}}, {"Kale", nil}}
	var checked ui.Selection[string]
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if CheckTree(c, "Food", nodes, &checked).Changed() {
			changed++
		}
	}, 400, 400)
	before := baseline()
	click(t, tt, "Fruit")
	if !checked.Has("Fruit/Apple") || !checked.Has("Fruit/Pear") || checked.Len() != 2 {
		t.Fatalf("parent check: %d checked", checked.Len())
	}
	click(t, tt, "Apple")
	if checked.Has("Fruit/Apple") || !checked.Has("Fruit/Pear") {
		t.Fatalf("leaf uncheck failed")
	}
	click(t, tt, "Fruit")
	if checked.Len() != 2 || changed != 3 {
		t.Errorf("mixed parent click: %d checked, changed %d", checked.Len(), changed)
	}
	cues(t, before, sound.Toggle, sound.Toggle, sound.Toggle)
	panics(t, "slash", func() {
		render(func(c *ui.Context) { var s ui.Selection[string]; CheckTree(c, "x", []Node{{"a/b", nil}}, &s) })
	})
}

func TestDirTree(t *testing.T) {
	dirs := []Dir{{"src", 12, []Dir{{"lib", 4, nil}}}, {"docs", 3, nil}}
	var open ui.Selection[string]
	chosen := ""
	tt := ui.NewTester(func(c *ui.Context) { DirTree(c, "Files", dirs, &open, &chosen, DirTreeOptions{}) }, 400, 400)
	if tt.HasText("lib") || !tt.HasText("12") {
		t.Fatalf("start %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Expand src")
	if !open.Has("src") || !tt.HasText("lib") || !tt.FramePending() {
		t.Fatalf("open: %q", tt.Texts())
	}
	settles(t, tt, "opened")
	click(t, tt, "lib")
	if chosen != "src/lib" {
		t.Errorf("chose %q", chosen)
	}
	cues(t, before, sound.Toggle, sound.Select)
	panics(t, "negative count", func() {
		render(func(c *ui.Context) {
			var o ui.Selection[string]
			s := ""
			DirTree(c, "x", []Dir{{"a", -1, nil}}, &o, &s, DirTreeOptions{})
		})
	})
}

func TestJSONTree(t *testing.T) {
	var v any
	if err := json.Unmarshal([]byte(`{"name":"Ada","age":36,"tags":["x"],"ok":true,"none":null}`), &v); err != nil {
		t.Fatal(err)
	}
	th := ui.LightTheme()
	var open ui.Selection[string]
	tt := ui.NewTester(func(c *ui.Context) { JSONTree(c, "Payload", v, &open, JSONTreeOptions{}) }, 400, 400)
	if !tt.HasText("{5}") || tt.HasText(`"Ada"`) {
		t.Fatalf("closed %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Expand $")
	settles(t, tt, "opened")
	for _, want := range []string{`"Ada"`, "36", "true", "null", "[1]"} {
		if !tt.HasText(want) {
			t.Errorf("no %s in %q", want, tt.Texts())
		}
	}
	r := find(t, tt, `"Ada"`)
	if !tinted(tt, r, th.Success) {
		t.Error("strings are not in the success color")
	}
	cues(t, before, sound.Toggle)
	panics(t, "int", func() {
		render(func(c *ui.Context) { var o ui.Selection[string]; JSONTree(c, "x", 3, &o, JSONTreeOptions{}) })
	})
}

// tinted reports whether r holds ink near color k.
func tinted(tt *ui.Tester, r ui.Rect, k ui.Color) bool {
	img := tt.Image()
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if p := img.RGBAAt(x, y); near(p.R, k.R) && near(p.G, k.G) && near(p.B, k.B) {
				return true
			}
		}
	}
	return false
}

func TestKVTable(t *testing.T) {
	pairs := []Pair{{"Host", "a.dev"}, {"Port", "443"}}
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if KVTable(c, "Server", pairs, KVTableOptions{}).Changed() {
			changed++
		}
	}, 500, 300)
	before := baseline()
	click(t, tt, "Copy Port")
	if tt.Clipboard() != "443" || changed != 1 || !tt.FramePending() {
		t.Fatalf("copy: %q, changed %d", tt.Clipboard(), changed)
	}
	th := ui.LightTheme()
	tt.Advance(500 * time.Millisecond)
	if !tinted(tt, find(t, tt, "Copy Port"), th.Success) {
		t.Error("no check after copying")
	}
	// The check holds 1.5 s, then crossfades back.
	tt.Advance(2 * time.Second)
	if tt.FramePending() || tinted(tt, find(t, tt, "Copy Port"), th.Success) {
		t.Error("check did not turn back")
	}
	cues(t, before, sound.Pop)
	panics(t, "twin keys", func() {
		render(func(c *ui.Context) { KVTable(c, "x", []Pair{{"a", "1"}, {"a", "2"}}, KVTableOptions{}) })
	})
}

func TestPropertyGrid(t *testing.T) {
	name, visible, size := "Box", false, 10.0
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if PropertyGrid(c, "Layer", []Property{{Name: "Name", Value: &name}, {Name: "Visible", Value: &visible},
			{Name: "Size", Value: &size, Min: 0, Max: 100}}).Changed() {
			changed++
		}
	}, 500, 300)
	settles(t, tt, "grid")
	before := baseline()
	// Name texts match the editors' labels: go by Tab.
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	tt.Frame()
	if !visible {
		t.Fatal("switch did not flip")
	}
	tt.Key(ui.Shift, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	tt.Type("2")
	tt.Frame()
	if name != "Box2" {
		t.Fatalf("name %q", name)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyUp)
	tt.Frame()
	if size != 11 || changed != 3 {
		t.Errorf("size %v, changed %d", size, changed)
	}
	cues(t, before, sound.Tick, sound.Tick, sound.Tick)
	panics(t, "bad bounds", func() {
		render(func(c *ui.Context) {
			v := 5.0
			PropertyGrid(c, "x", []Property{{Name: "n", Value: &v, Min: 10, Max: 1}})
		})
	})
	panics(t, "int", func() {
		render(func(c *ui.Context) { v := 5; PropertyGrid(c, "x", []Property{{Name: "n", Value: &v}}) })
	})
}
