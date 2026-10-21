package input

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestCountedField(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { CountedField(c, "Bio", &value, CountedFieldOptions{Max: 10}) }, 400, 200)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Type("123456789")
	settles(t, tt, "CountedField warning")
	if value != "123456789" || !tt.HasText("9/10") || !appended(before, sound.Tick) {
		t.Errorf("value %q, texts %q", value, tt.Texts())
	}
}

func TestAdornedField(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) {
		AdornedField(c, "Amount", &value, func() { ui.Text(c, "$") }, func() { ui.Text(c, "USD") }, AdornedFieldOptions{})
	}, 400, 200)
	find(t, tt, "$")
	tt.Key(0, ui.KeyTab)
	settles(t, tt, "AdornedField focus")
	tt.Type("12")
	if value != "12" {
		t.Errorf("value %q", value)
	}
}

func TestClearableField(t *testing.T) {
	value, changed := "", 0
	tt := ui.NewTester(func(c *ui.Context) {
		if ClearableField(c, "Name", &value, ClearableFieldOptions{}).Changed() {
			changed++
		}
	}, 400, 200)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Type("Ada")
	settles(t, tt, "ClearableField clear")
	if err := tt.Click("Clear Name"); err != nil {
		t.Fatal(err)
	}
	if value != "" || changed == 0 || !appended(before, sound.Close) {
		t.Errorf("value %q, changed %d", value, changed)
	}
}

func TestFloatingLabel(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { FloatingLabel(c, "Email", &value, FloatingLabelOptions{}) }, 400, 200)
	before := tt.Image()
	tt.Key(0, ui.KeyTab)
	settles(t, tt, "FloatingLabel focus")
	if !rectChanged(before, tt.Image(), 10, 10, 90, 28) {
		t.Error("label did not float")
	}
	tt.Type("a@b.co")
	if value == "" {
		t.Error("input did not change")
	}
}

func TestUnderlineField(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { UnderlineField(c, "Name", &value, UnderlineFieldOptions{}) }, 400, 200)
	tt.Key(0, ui.KeyTab)
	settles(t, tt, "UnderlineField focus")
	tt.Type("Ada")
	if value != "Ada" {
		t.Errorf("value %q", value)
	}
}

func TestFilledField(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { FilledField(c, "Name", &value, FilledFieldOptions{}) }, 400, 200)
	tt.Key(0, ui.KeyTab)
	settles(t, tt, "FilledField focus")
	tt.Type("Ada")
	if value != "Ada" {
		t.Errorf("value %q", value)
	}
}

func TestJSONField(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { JSONField(c, &value, JSONFieldOptions{}) }, 500, 300)
	tt.Key(0, ui.KeyTab)
	tt.Type("{")
	if !tt.HasText("Enter valid JSON.") {
		t.Fatal("JSON error missing")
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type(`{"ok":true}`)
	if tt.HasText("Enter valid JSON.") {
		t.Error("valid JSON still errors")
	}
}

func TestRegexField(t *testing.T) {
	pattern := "a+"
	tt := ui.NewTester(func(c *ui.Context) { RegexField(c, &pattern, "aa ba", RegexFieldOptions{}) }, 400, 200)
	find(t, tt, "2 matches")
	tt.Key(0, ui.KeyTab)
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("[")
	if !tt.HasText("Enter a valid regular expression.") {
		t.Error("regex error missing")
	}
}

func TestKeyValueField(t *testing.T) {
	rows := []KeyValueRow{}
	tt := ui.NewTester(func(c *ui.Context) { KeyValueField(c, &rows, KeyValueFieldOptions{}) }, 500, 300)
	before := baseline()
	if err := tt.Click("Add row"); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || !appended(before, sound.Pop) {
		t.Fatalf("rows %d", len(rows))
	}
	tt.Key(0, ui.KeyTab)
	tt.Type("host")
	tt.Key(0, ui.KeyTab)
	tt.Type("localhost")
	if rows[0].Key != "host" || rows[0].Value != "localhost" {
		t.Errorf("row %+v", rows[0])
	}
	if err := tt.Click("Remove row"); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 || !appended(before, sound.Delete) {
		t.Errorf("rows %d", len(rows))
	}
}

func TestListField(t *testing.T) {
	items := []string{}
	tt := ui.NewTester(func(c *ui.Context) { ListField(c, "Items", &items, ListFieldOptions{}) }, 500, 300)
	before := baseline()
	if err := tt.Click("Add item"); err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || !appended(before, sound.Pop) {
		t.Fatalf("items %d", len(items))
	}
	tt.Key(0, ui.KeyTab)
	tt.Type("First")
	if items[0] != "First" {
		t.Errorf("item %q", items[0])
	}
	if err := tt.Click("Remove Item"); err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 || !appended(before, sound.Delete) {
		t.Errorf("items %d", len(items))
	}
}
