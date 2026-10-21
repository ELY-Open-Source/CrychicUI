package input

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestOTPInput(t *testing.T) {
	value, changed := "", 0
	tt := ui.NewTester(func(c *ui.Context) {
		if OTPInput(c, &value, 4, OTPInputOptions{}).Changed() {
			changed++
		}
	}, 400, 200)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Type("1234")
	if value != "1234" || changed == 0 || !appended(before, sound.Tick, sound.Success) {
		t.Errorf("value %q, changed %d, cues %v", value, changed, sound.Recent())
	}
}

func TestPinInput(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { PinInput(c, &value, 4, PinInputOptions{}) }, 400, 200)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Type("12")
	settles(t, tt, "PinInput dots")
	if value != "12" || textCount(tt, "•") != 2 || !appended(before, sound.Tick) {
		t.Errorf("value %q, texts %q", value, tt.Texts())
	}
}

func textCount(tt *ui.Tester, want string) int {
	n := 0
	for _, s := range tt.Texts() {
		if s == want {
			n++
		}
	}
	return n
}

func TestTagInput(t *testing.T) {
	tags := []string{}
	draft := ""
	tt := ui.NewTester(func(c *ui.Context) { TagInput(c, &tags, &draft, TagInputOptions{}) }, 400, 240)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Type("design,")
	if len(tags) != 1 || tags[0] != "design" || !appended(before, sound.Pop) {
		t.Fatalf("tags %q, cues %v", tags, sound.Recent())
	}
	tt.Key(0, ui.KeyBackspace)
	if len(tags) != 0 || !appended(before, sound.Delete) {
		t.Errorf("tags %q, cues %v", tags, sound.Recent())
	}
}

func TestMentionInput(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { MentionInput(c, "Comment", &value, []string{"Ada", "Grace"}) }, 400, 240)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Type("Hello @a")
	if err := tt.Click("Ada"); err != nil {
		t.Fatal(err)
	}
	if value != "Hello @Ada " || !appended(before, sound.Select) {
		t.Errorf("value %q, cues %v", value, sound.Recent())
	}
}

func TestCommandInput(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { CommandInput(c, "Command", &value, []string{"delete", "duplicate"}) }, 400, 240)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Type("/del")
	if err := tt.Click("delete"); err != nil {
		t.Fatal(err)
	}
	if value != "/delete " || !appended(before, sound.Select) {
		t.Errorf("value %q", value)
	}
}

func TestInlineEdit(t *testing.T) {
	value, changed := "Draft", 0
	tt := ui.NewTester(func(c *ui.Context) {
		if InlineEdit(c, &value, InlineEditOptions{}).Changed() {
			changed++
		}
	}, 400, 200)
	before := baseline()
	if err := tt.Click("Edit text"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("New title")
	tt.Key(0, ui.KeyEnter)
	if value != "New title" || changed != 1 || !appended(before, sound.Success) {
		t.Errorf("value %q, changed %d", value, changed)
	}
}

func TestEditableTitle(t *testing.T) {
	value := "Draft"
	tt := ui.NewTester(func(c *ui.Context) { EditableTitle(c, &value, EditableTitleOptions{}) }, 400, 200)
	if err := tt.Click("Edit text"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("New")
	tt.Key(0, ui.KeyEnter)
	if value != "New" {
		t.Errorf("value %q", value)
	}
}

func TestSuggest(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { Suggest(c, "Fruit", &value, []string{"Apple", "Pear"}) }, 400, 240)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Type("ap")
	if err := tt.Click("Apple"); err != nil {
		t.Fatal(err)
	}
	if value != "Apple" || !appended(before, sound.Select) {
		t.Errorf("value %q", value)
	}
	tt.Frame()
	if tt.HasText("Apple") {
		t.Error("selected suggestion reopened")
	}
}

func TestComboField(t *testing.T) {
	selected, changed := "One", 0
	tt := ui.NewTester(func(c *ui.Context) {
		if ComboField(c, "Choice", &selected, []string{"One", "Two"}).Changed() {
			changed++
		}
	}, 400, 240)
	before := baseline()
	if err := tt.Click("Choice"); err != nil {
		t.Fatal(err)
	}
	settles(t, tt, "ComboField popup")
	if changed != 0 || selected != "One" {
		t.Fatalf("opening chose %q", selected)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("Two")
	tt.Key(0, ui.KeyEnter)
	if selected != "Two" || changed != 1 || !appended(before, sound.Select) {
		t.Errorf("selected %q, changed %d", selected, changed)
	}
}

func TestMultiCombo(t *testing.T) {
	chosen := []string{}
	tt := ui.NewTester(func(c *ui.Context) { MultiCombo(c, "Choose", &chosen, []string{"One", "Two"}) }, 400, 240)
	before := baseline()
	if err := tt.Click("Choose"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Two")
	tt.Key(0, ui.KeyEnter)
	if len(chosen) != 1 || chosen[0] != "Two" || !appended(before, sound.Pop) {
		t.Fatalf("chosen %q", chosen)
	}
	if err := tt.Click("Remove Two"); err != nil {
		t.Fatal(err)
	}
	if len(chosen) != 0 || !appended(before, sound.Delete) {
		t.Errorf("chosen %q", chosen)
	}
}

func TestAddressField(t *testing.T) {
	value := Address{}
	tt := ui.NewTester(func(c *ui.Context) { AddressField(c, &value, AddressFieldOptions{}) }, 500, 300)
	tt.Key(0, ui.KeyTab)
	tt.Type("1 Main")
	tt.Key(0, ui.KeyTab)
	tt.Type("Springfield")
	if value.Street != "1 Main" || value.City != "Springfield" {
		t.Errorf("address %+v", value)
	}
}

func TestNameField(t *testing.T) {
	given, family := "", ""
	tt := ui.NewTester(func(c *ui.Context) { NameField(c, &given, &family, NameFieldOptions{}) }, 500, 200)
	tt.Key(0, ui.KeyTab)
	tt.Type("Ada")
	tt.Key(0, ui.KeyTab)
	tt.Type("Lovelace")
	if given != "Ada" || family != "Lovelace" {
		t.Errorf("%q %q", given, family)
	}
}

func TestCardNumberField(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { CardNumberField(c, &value, CardNumberFieldOptions{}) }, 400, 200)
	tt.Key(0, ui.KeyTab)
	tt.Type("4242424242424242")
	if value != "4242 4242 4242 4242" {
		t.Errorf("value %q", value)
	}
}

func TestExpiryField(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { ExpiryField(c, &value, ExpiryFieldOptions{}) }, 400, 200)
	tt.Key(0, ui.KeyTab)
	tt.Type("1330")
	if value != "13/30" || !tt.HasText("Use a month from 01 to 12.") {
		t.Errorf("value %q, texts %q", value, tt.Texts())
	}
}

func TestCVCField(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { CVCField(c, &value, CVCFieldOptions{}) }, 400, 200)
	tt.Key(0, ui.KeyTab)
	tt.Type("12")
	if value != "12" || !tt.HasText("Use three or four digits.") {
		t.Errorf("value %q, texts %q", value, tt.Texts())
	}
}

func TestIBANField(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { IBANField(c, &value, IBANFieldOptions{}) }, 500, 240)
	tt.Key(0, ui.KeyTab)
	tt.Type("GB29NWBK60161331926819")
	if value != "GB29 NWBK 6016 1331 9268 19" || tt.HasText("Check the IBAN checksum.") {
		t.Errorf("value %q, texts %q", value, tt.Texts())
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("GB29NWBK60161331926818")
	if !tt.HasText("Check the IBAN checksum.") {
		t.Error("invalid checksum missing")
	}
}

func TestSlugField(t *testing.T) {
	title := "My First Post"
	slug, locked, changed := "", false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if SlugField(c, title, &slug, &locked, SlugFieldOptions{}).Changed() {
			changed++
		}
	}, 400, 200)
	if slug != "my-first-post" {
		t.Fatalf("slug %q", slug)
	}
	before := baseline()
	if err := tt.Click("Lock Slug"); err != nil {
		t.Fatal(err)
	}
	if !locked || changed == 0 || !appended(before, sound.Toggle) {
		t.Errorf("locked %v, changed %d", locked, changed)
	}
	title = "New Title"
	tt.Frame()
	if slug != "my-first-post" {
		t.Error("locked slug followed the title")
	}
}
