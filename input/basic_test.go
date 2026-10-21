package input

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestTextField(t *testing.T) {
	value, changed := "", 0
	errorText := ""
	tt := ui.NewTester(func(c *ui.Context) {
		if TextField(c, "Email", &value, TextFieldOptions{Placeholder: "name@example.com", Help: "We never share it.", Error: errorText}).Changed() {
			changed++
		}
	}, 400, 240)
	find(t, tt, "We never share it.")
	before := baseline()
	r := find(t, tt, "Email")
	tt.ClickAt(r.X+20, r.Y+r.H+20)
	settles(t, tt, "TextField focus")
	tt.Type("bad")
	if changed == 0 || !appended(before, sound.Tick) {
		t.Errorf("changed %d, cues %v", changed, sound.Recent())
	}
	errorText = "Enter a valid address."
	tt.Frame()
	settles(t, tt, "TextField error")
	find(t, tt, errorText)
}

func TestAutoTextArea(t *testing.T) {
	value := ""
	height := float32(0)
	var bounds ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		b := AutoTextArea(c, &value, AutoTextAreaOptions{Placeholder: "Notes"})
		height = b.LayoutBounds().H
		bounds = b.LayoutBounds()
	}, 400, 300)
	tt.Frame()
	initial := height
	before := baseline()
	tt.ClickAt(bounds.X+bounds.W/2, bounds.Y+bounds.H/2)
	tt.Type("one\ntwo\nthree\nfour\nfive")
	tt.Frame()
	settles(t, tt, "AutoTextArea height")
	if height <= initial {
		t.Errorf("height %v, initial %v", height, initial)
	}
	if !appended(before, sound.Tick) {
		t.Error("typing played no cue")
	}
}

func TestPasswordField(t *testing.T) {
	value, changed := "", 0
	tt := ui.NewTester(func(c *ui.Context) {
		if PasswordField(c, "Password", &value, PasswordFieldOptions{Placeholder: "Password"}).Changed() {
			changed++
		}
	}, 400, 240)
	before := baseline()
	r := find(t, tt, "Password")
	tt.ClickAt(r.X+20, r.Y+r.H+10)
	tt.Type("Abcdef1!")
	if value == "" || changed == 0 || !appended(before, sound.Tick) {
		t.Error("password did not change")
	}
	if err := tt.Click("Reveal Password"); err != nil {
		t.Fatal(err)
	}
	if changed < 2 || !appended(before, sound.Toggle) {
		t.Error("reveal did not toggle")
	}
}

func TestSearchInput(t *testing.T) {
	value, changed := "", 0
	tt := ui.NewTester(func(c *ui.Context) {
		root := SearchInput(c, &value, SearchInputOptions{})
		if root.Changed() {
			changed++
		}
	}, 400, 200)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Type("query")
	if value != "query" || changed == 0 {
		t.Fatal("search did not change")
	}
	settles(t, tt, "SearchInput clear")
	if err := tt.Click("Clear search"); err != nil {
		t.Fatal(err)
	}
	if value != "" || !appended(before, sound.Close) {
		t.Error("clear button failed")
	}
	tt.Key(0, ui.KeyTab)
	tt.Type("again")
	if value != "again" {
		t.Fatal("input did not regain focus")
	}
	tt.Key(0, ui.KeyEscape)
	if value != "" {
		t.Error("Escape did not clear")
	}
}

func TestEmailField(t *testing.T) {
	value, changed := "", 0
	tt := ui.NewTester(func(c *ui.Context) {
		if EmailField(c, "Email", &value, EmailFieldOptions{Placeholder: "Email"}).Changed() {
			changed++
		}
	}, 400, 240)
	r := find(t, tt, "Email")
	tt.ClickAt(r.X+20, r.Y+r.H+10)
	tt.Type("bad")
	tt.ClickAt(390, 230)
	settles(t, tt, "EmailField error")
	find(t, tt, "Enter a valid email address.")
	if changed == 0 {
		t.Error("input did not change")
	}
}

func TestURLField(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) {
		URLField(c, "Site", &value, URLFieldOptions{Placeholder: "example.com"})
	}, 400, 240)
	tt.Key(0, ui.KeyTab)
	tt.Type("%")
	tt.ClickAt(390, 230)
	settles(t, tt, "URLField error")
	find(t, tt, "Enter a valid URL host.")
}

func TestPhoneField(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) {
		PhoneField(c, "Phone", &value, PhoneFieldOptions{Placeholder: "555"})
	}, 400, 200)
	r := find(t, tt, "Phone")
	tt.ClickAt(r.X+40, r.Y+r.H+10)
	tt.Type("15551234567")
	if value != "155 512 3456 7" {
		t.Errorf("value %q", value)
	}
}

func TestNumberField(t *testing.T) {
	value := 4.0
	tt := ui.NewTester(func(c *ui.Context) { NumberField(c, "Copies", &value, NumberFieldOptions{}) }, 300, 200)
	before := baseline()
	if err := tt.Click("Increase Copies"); err != nil {
		t.Fatal(err)
	}
	settles(t, tt, "NumberField value")
	if value != 5 || !appended(before, sound.Tick) {
		t.Errorf("value %v, cues %v", value, sound.Recent())
	}
}

func TestCurrencyField(t *testing.T) {
	value := 1.234
	tt := ui.NewTester(func(c *ui.Context) { CurrencyField(c, "Price", &value, CurrencyFieldOptions{}) }, 400, 200)
	if err := tt.Click("Increase"); err != nil {
		t.Fatal(err)
	}
	tt.ClickAt(390, 190)
	if value != 1.24 {
		t.Errorf("value %v", value)
	}
}

func TestPercentField(t *testing.T) {
	value := 99.0
	tt := ui.NewTester(func(c *ui.Context) { PercentField(c, "Progress", &value) }, 300, 200)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyUp)
	if value != 100 {
		t.Errorf("value %v", value)
	}
}

func TestUnitField(t *testing.T) {
	value, changed := 2.0, 0
	unit := "kg"
	tt := ui.NewTester(func(c *ui.Context) {
		if UnitField(c, "Weight", &value, &unit, []string{"g", "kg"}, UnitFieldOptions{}).Changed() {
			changed++
		}
	}, 400, 200)
	before := baseline()
	if err := tt.Click("kg"); err != nil {
		t.Fatal(err)
	}
	if err := tt.ChooseMenuItem("g"); err != nil {
		t.Fatal(err)
	}
	if unit != "g" || changed != 1 || !appended(before, sound.Select) {
		t.Errorf("unit %q, changed %d", unit, changed)
	}
}
