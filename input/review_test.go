package input

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestPercentFieldClampsTypedText(t *testing.T) {
	value := 42.0
	tt := ui.NewTester(func(c *ui.Context) { PercentField(c, "Progress", &value) }, 300, 200)
	tt.Key(0, ui.KeyTab)
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("150")
	if value != 100 {
		t.Errorf("value %v", value)
	}
	before := baseline()
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("NaN")
	if value != 100 || !appended(before) {
		t.Errorf("value %v, cues %v", value, sound.Recent())
	}
}

func TestSuggestStaysDismissed(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { Suggest(c, "Fruit", &value, []string{"Apple", "Pear"}) }, 400, 240)
	tt.Key(0, ui.KeyTab)
	tt.Type("ap")
	if !tt.HasText("Apple") {
		t.Fatal("suggestion missing")
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if tt.HasText("Apple") {
		t.Error("dismissed suggestion reopened")
	}
}

func TestMentionStaysDismissed(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { MentionInput(c, "Comment", &value, []string{"Ada", "Grace"}) }, 400, 240)
	tt.Key(0, ui.KeyTab)
	tt.Type("@a")
	if !tt.HasText("Ada") {
		t.Fatal("mention missing")
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if tt.HasText("Ada") {
		t.Error("dismissed mention reopened")
	}
}

func TestASCIIFormatters(t *testing.T) {
	if got := groupDigits("１２34", []int{3, 3, 4}); got != "34" {
		t.Errorf("phone digits %q", got)
	}
	if got := groupCardDigits("４２42"); got != "42" {
		t.Errorf("card digits %q", got)
	}
	if got := formatExpiry("１２30"); got != "30" {
		t.Errorf("expiry digits %q", got)
	}
	if got := formatIBAN("ＧＢ29"); got != "29" {
		t.Errorf("IBAN digits %q", got)
	}
}

func TestNumberFieldPrecision(t *testing.T) {
	for _, tc := range []struct {
		step float64
		want string
	}{{0.0001, "0.0001"}, {0.25, "0.25"}, {0.125, "0.125"}} {
		t.Run(tc.want, func(t *testing.T) {
			value := 0.0
			tt := ui.NewTester(func(c *ui.Context) {
				NumberField(c, "Precision", &value, NumberFieldOptions{Max: 10, Step: tc.step})
			}, 400, 200)
			if err := tt.Click("Increase Precision"); err != nil {
				t.Fatal(err)
			}
			if value != tc.step || !tt.HasText(tc.want) {
				t.Errorf("value %v, texts %q", value, tt.Texts())
			}
		})
	}
}

func TestNumberFieldNoChangeAtLimit(t *testing.T) {
	value, changed := 0.9999, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if NumberField(c, "Precision", &value, NumberFieldOptions{Max: 1, Step: 0.0001}).Changed() {
			changed++
		}
	}, 400, 200)
	if err := tt.Click("Increase Precision"); err != nil {
		t.Fatal(err)
	}
	if value != 1 || changed != 1 {
		t.Fatalf("value %v, changed %d", value, changed)
	}
	before := baseline()
	if err := tt.Click("Increase Precision"); err != nil {
		t.Fatal(err)
	}
	if changed != 1 || !appended(before) {
		t.Errorf("changed %d, cues %v", changed, sound.Recent())
	}
}

func TestOTPInitialCompletePlaysNoCue(t *testing.T) {
	value := "1234"
	before := baseline()
	ui.NewTester(func(c *ui.Context) { OTPInput(c, &value, 4, OTPInputOptions{}) }, 400, 200)
	if !appended(before) {
		t.Error("initial complete OTP played a cue")
	}
}

func TestOTPProgrammaticFillPlaysNoCue(t *testing.T) {
	value := "12"
	tt := ui.NewTester(func(c *ui.Context) { OTPInput(c, &value, 4, OTPInputOptions{}) }, 400, 200)
	before := baseline()
	value = "1234"
	tt.Frame()
	if !appended(before) {
		t.Error("programmatic OTP completion played a cue")
	}
}

func TestInputValidation(t *testing.T) {
	panics(t, "TextField value", func() { ui.Render(func(c *ui.Context) { TextField(c, "Email", nil, TextFieldOptions{}) }, 10, 10, 1) })
	panics(t, "AutoTextArea value", func() { ui.Render(func(c *ui.Context) { AutoTextArea(c, nil, AutoTextAreaOptions{}) }, 10, 10, 1) })
	panics(t, "PasswordField value", func() {
		ui.Render(func(c *ui.Context) { PasswordField(c, "Password", nil, PasswordFieldOptions{}) }, 10, 10, 1)
	})
	panics(t, "SearchInput value", func() { ui.Render(func(c *ui.Context) { SearchInput(c, nil, SearchInputOptions{}) }, 10, 10, 1) })
	panics(t, "NumberField value", func() {
		ui.Render(func(c *ui.Context) { NumberField(c, "Count", nil, NumberFieldOptions{}) }, 10, 10, 1)
	})
	panics(t, "CurrencyField value", func() {
		ui.Render(func(c *ui.Context) { CurrencyField(c, "Price", nil, CurrencyFieldOptions{}) }, 10, 10, 1)
	})
	panics(t, "PercentField value", func() { ui.Render(func(c *ui.Context) { PercentField(c, "Progress", nil) }, 10, 10, 1) })
	panics(t, "UnitField unit", func() {
		v := 1.0
		unit := "px"
		ui.Render(func(c *ui.Context) { UnitField(c, "Size", &v, &unit, nil, UnitFieldOptions{}) }, 10, 10, 1)
	})
	panics(t, "OTP length", func() { ui.Render(func(c *ui.Context) { OTPInput(c, new(string), 0, OTPInputOptions{}) }, 10, 10, 1) })
	panics(t, "TagInput tags", func() { ui.Render(func(c *ui.Context) { TagInput(c, nil, new(string), TagInputOptions{}) }, 10, 10, 1) })
	panics(t, "MentionInput items", func() { ui.Render(func(c *ui.Context) { MentionInput(c, "Comment", new(string), nil) }, 10, 10, 1) })
	panics(t, "InlineEdit value", func() { ui.Render(func(c *ui.Context) { InlineEdit(c, nil, InlineEditOptions{}) }, 10, 10, 1) })
	panics(t, "Suggest items", func() { ui.Render(func(c *ui.Context) { Suggest(c, "Fruit", new(string), nil) }, 10, 10, 1) })
	panics(t, "ComboField option", func() { ui.Render(func(c *ui.Context) { ComboField(c, "Choice", new(string), nil) }, 10, 10, 1) })
	panics(t, "MultiCombo options", func() { ui.Render(func(c *ui.Context) { MultiCombo(c, "Choose", new([]string), nil) }, 10, 10, 1) })
	panics(t, "AddressField value", func() { ui.Render(func(c *ui.Context) { AddressField(c, nil, AddressFieldOptions{}) }, 10, 10, 1) })
	panics(t, "NameField value", func() {
		ui.Render(func(c *ui.Context) { NameField(c, nil, new(string), NameFieldOptions{}) }, 10, 10, 1)
	})
	panics(t, "CardNumberField value", func() {
		ui.Render(func(c *ui.Context) { CardNumberField(c, nil, CardNumberFieldOptions{}) }, 10, 10, 1)
	})
	panics(t, "ExpiryField value", func() { ui.Render(func(c *ui.Context) { ExpiryField(c, nil, ExpiryFieldOptions{}) }, 10, 10, 1) })
	panics(t, "CVCField value", func() { ui.Render(func(c *ui.Context) { CVCField(c, nil, CVCFieldOptions{}) }, 10, 10, 1) })
	panics(t, "IBANField value", func() { ui.Render(func(c *ui.Context) { IBANField(c, nil, IBANFieldOptions{}) }, 10, 10, 1) })
	panics(t, "SlugField slug", func() {
		ui.Render(func(c *ui.Context) { SlugField(c, "Title", nil, new(bool), SlugFieldOptions{}) }, 10, 10, 1)
	})
}
