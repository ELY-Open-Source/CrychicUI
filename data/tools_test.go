package data

import (
	"slices"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestCSVPreview(t *testing.T) {
	text, header := "name,age\nAda,36\nGrace,45", true
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if CSVPreview(c, "Import", text, &header, CSVPreviewOptions{}).Changed() {
			changed++
		}
	}, 600, 400)
	settles(t, tt, "preview")
	if !tt.HasText("Grace") || tt.HasText("Column 1") {
		t.Fatalf("with header %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "First row is a header")
	if header || changed != 1 || !tt.HasText("Column 1") || !tt.FramePending() {
		t.Fatalf("header %v, changed %d, %q", header, changed, tt.Texts())
	}
	settles(t, tt, "fade")
	cues(t, before, sound.Toggle)
	text = "a,b\nc"
	tt.Frame()
	if !tt.HasText("wrong number of fields") {
		t.Errorf("ragged rows show %q", tt.Texts())
	}
	text = ""
	tt.Frame()
	if !tt.HasText("No rows") {
		t.Errorf("empty shows %q", tt.Texts())
	}
}

func TestColumnMapper(t *testing.T) {
	mapping := []string{"", "Email"}
	tt := ui.NewTester(func(c *ui.Context) {
		ColumnMapper(c, "Columns", []string{"full name", "mail"}, []string{"Name", "Email"}, &mapping, ColumnMapperOptions{})
	}, 600, 400)
	click(t, tt, "Map full name")
	settles(t, tt, "menu")
	before := baseline()
	click(t, tt, "Name")
	if !slices.Equal(mapping, []string{"Name", "Email"}) {
		t.Fatalf("mapping %q", mapping)
	}
	cues(t, before, sound.Select)
	panics(t, "unknown field", func() {
		render(func(c *ui.Context) {
			m := []string{"Phone"}
			ColumnMapper(c, "x", []string{"a"}, []string{"Name"}, &m, ColumnMapperOptions{})
		})
	})
	panics(t, "short mapping", func() {
		render(func(c *ui.Context) {
			var m []string
			ColumnMapper(c, "x", []string{"a"}, []string{"Name"}, &m, ColumnMapperOptions{})
		})
	})
}

func TestExportMenu(t *testing.T) {
	chosen, open := -1, false
	tt := ui.NewTester(func(c *ui.Context) {
		ExportMenu(c, "Export", []Export{{"CSV", "a,b"}, {"JSON", `["a","b"]`}}, &chosen, &open)
	}, 500, 400)
	before := baseline()
	click(t, tt, "Export")
	settles(t, tt, "menu")
	click(t, tt, "JSON")
	if chosen != 1 || open || tt.Clipboard() != `["a","b"]` {
		t.Fatalf("chosen %d, open %v, clipboard %q", chosen, open, tt.Clipboard())
	}
	cues(t, before, sound.Open, sound.Pop)
}

func TestTableToolbarAndDensity(t *testing.T) {
	query, d := "", DensityDefault
	changed := 0
	var hits int
	tt := ui.NewTester(func(c *ui.Context) {
		root, filters, _ := TableToolbar(c, "People", &query, &d, TableToolbarOptions{})
		if root.Changed() {
			changed++
		}
		if filters.Clicked() {
			hits++
		}
	}, 900, 200)
	before := baseline()
	click(t, tt, "Filters")
	click(t, tt, "Roomy")
	if hits != 1 || d != DensityRoomy || changed != 1 {
		t.Fatalf("filters %d, density %v, changed %d", hits, d, changed)
	}
	click(t, tt, "Search")
	tt.Type("ada")
	tt.Frame()
	if query != "ada" || changed != 2 {
		t.Errorf("query %q, changed %d", query, changed)
	}
	cues(t, before, sound.Tap, sound.Select)
	th := ui.LightTheme()
	if DensityCompact.Pad(th) >= DensityDefault.Pad(th) || DensityRoomy.Pad(th) <= DensityDefault.Pad(th) {
		t.Error("densities out of order")
	}
	own := ui.NewTester(func(c *ui.Context) {
		TableToolbar(c, "x", &query, &d, TableToolbarOptions{Densities: [3]string{"Dicht", "Normal", "Weit"}})
	}, 900, 200)
	if !own.HasText("Dicht") || own.HasText("Compact") {
		t.Errorf("density copy not passed: %q", own.Texts())
	}
	panics(t, "density", func() {
		render(func(c *ui.Context) { d := Density(9); DensitySwitch(c, "x", &d, DensitySwitchOptions{}) })
	})
}
