package input

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestKeyValueFieldRemoveFirst(t *testing.T) {
	rows := []KeyValueRow{{Key: "a", Value: "1"}, {Key: "b", Value: "2"}}
	tt := ui.NewTester(func(c *ui.Context) { KeyValueField(c, &rows, KeyValueFieldOptions{}) }, 500, 300)
	if err := tt.Click("Remove row"); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Key != "b" {
		t.Errorf("rows %+v", rows)
	}
}

func TestListFieldRemoveFirst(t *testing.T) {
	items := []string{"a", "b"}
	tt := ui.NewTester(func(c *ui.Context) { ListField(c, "Items", &items, ListFieldOptions{}) }, 500, 300)
	if err := tt.Click("Remove Item"); err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0] != "b" {
		t.Errorf("items %q", items)
	}
}

func TestMarkdownPreviewStaysVisible(t *testing.T) {
	source := ""
	var surface ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		surface = c.Theme().Surface
		MarkdownEditor(c, &source, MarkdownEditorOptions{})
	}, 700, 300)
	tt.Key(0, ui.KeyTab)
	tt.Type("# Hello")
	settles(t, tt, "MarkdownEditor preview")
	r := find(t, tt, "Hello")
	inked := false
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			inked = inked || !near(tt.Image(), x, y, surface)
		}
	}
	if !inked {
		t.Error("preview text is transparent")
	}
}

func TestIPFieldPreservesDots(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { IPField(c, &value, IPFieldOptions{}) }, 400, 200)
	tt.Key(0, ui.KeyTab)
	tt.Type("192.168.1.1")
	if value != "192.168.1.1" || tt.HasText("Each octet") {
		t.Errorf("value %q, texts %q", value, tt.Texts())
	}
}

func TestIPFieldIncrementalEntry(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { IPField(c, &value, IPFieldOptions{}) }, 400, 200)
	tt.Key(0, ui.KeyTab)
	for _, step := range []struct{ text, want string }{
		{"192", "192"}, {"168", "192.168"}, {"0", "192.168.0"}, {"1", "192.168.01"}, {"8", "192.168.018"}, {"1", "192.168.018.1"},
	} {
		tt.Type(step.text)
		tt.Frame()
		if value != step.want {
			t.Fatalf("after %q: %q, want %q", step.text, value, step.want)
		}
	}
	tt.Key(ui.Cmd, ui.KeyA)
	for _, text := range []string{"192.", "168.", "0.", "1"} {
		tt.Type(text)
		tt.Frame()
	}
	if value != "192.168.0.1" {
		t.Errorf("dotted value %q", value)
	}
}

func TestIPFieldEditsFirstOctet(t *testing.T) {
	value := "192.168.0.1"
	tt := ui.NewTester(func(c *ui.Context) { IPField(c, &value, IPFieldOptions{}) }, 400, 200)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyDelete)
	tt.Type("1")
	if value != "192.168.0.1" {
		t.Errorf("value %q", value)
	}
}

func TestIPFieldRejectsNegativeOctets(t *testing.T) {
	value := "-1.2.3.4"
	tt := ui.NewTester(func(c *ui.Context) { IPField(c, &value, IPFieldOptions{}) }, 400, 200)
	if !tt.HasText("Each octet must be 0 through 255.") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestIPFieldPlaceholderOverride(t *testing.T) {
	value := ""
	standard := ui.NewTester(func(c *ui.Context) { IPField(c, &value, IPFieldOptions{}) }, 400, 200)
	custom := ""
	customField := ui.NewTester(func(c *ui.Context) {
		IPField(c, &custom, IPFieldOptions{Placeholder: "IPv4 address"})
	}, 400, 200)
	if !rectChanged(standard.Image(), customField.Image(), 0, 20, 300, 80) {
		t.Error("IP placeholder override did not render")
	}
}

func TestSignaturePadOffsetAndInk(t *testing.T) {
	changed := 0
	padAt := float32(50)
	var background ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		background = c.Theme().Background
		ui.Column(c).Padding(padAt).Children(func() {
			if SignaturePad(c, SignaturePadOptions{}).Changed() {
				changed++
			}
		})
	}, 500, 400)
	before := tt.Image()
	tt.Press(60, 60)
	tt.Move(100, 60)
	tt.Release(100, 60)
	tt.Advance(100 * time.Millisecond)
	if !rectChanged(before, tt.Image(), 55, 55, 110, 70) {
		t.Fatal("stroke did not draw at the pad's offset")
	}
	tt.Advance(200 * time.Millisecond)
	padAt = 100
	tt.Frame()
	if !near(tt.Image(), 60, 60, background) {
		t.Error("old stroke remained after the pad moved")
	}
	if near(tt.Image(), 110, 110, background) {
		t.Error("stroke did not move with the pad")
	}
	if changed != 1 || tt.FramePending() {
		t.Errorf("changed %d, pending %v", changed, tt.FramePending())
	}
}

func TestVerifyFieldCountsDown(t *testing.T) {
	code, changed := "", 0
	tt := ui.NewTester(func(c *ui.Context) {
		if VerifyField(c, "Code", &code, VerifyFieldOptions{}).Changed() {
			changed++
		}
	}, 400, 240)
	before := baseline()
	if err := tt.Click("Send code"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Resend in 3s") || !appended(before, sound.Send) {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt.Advance(time.Second)
	if !tt.HasText("Resend in 2s") {
		t.Errorf("countdown texts %q", tt.Texts())
	}
	tt.Advance(2 * time.Second)
	if !tt.HasText("Send code") {
		t.Error("resend did not re-enable")
	}
	if err := tt.Click("Send code"); err != nil {
		t.Fatal(err)
	}
	if changed != 2 || !appended(before, sound.Send, sound.Send) {
		t.Errorf("changed %d, cues %v", changed, sound.Recent())
	}
}

func TestFloatingLabelClickFocuses(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { FloatingLabel(c, "Email", &value, FloatingLabelOptions{}) }, 400, 200)
	tt.ClickAt(20, 32)
	tt.Type("a@b.co")
	if value != "a@b.co" {
		t.Errorf("value %q", value)
	}
}

func TestNoteInputClickFocuses(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { NoteInput(c, &value, "Write a note") }, 500, 300)
	tt.ClickAt(20, 24)
	tt.Type("A note")
	if value != "A note" {
		t.Errorf("value %q", value)
	}
}

func TestConfirmPasswordProgrammaticMatchPlaysNoCue(t *testing.T) {
	password, confirm := "secret", "secret"
	before := baseline()
	ui.NewTester(func(c *ui.Context) { ConfirmPassword(c, &password, &confirm, ConfirmPasswordOptions{}) }, 400, 240)
	if !appended(before) {
		t.Error("initial matching passwords played Success")
	}
}

func TestCountedFieldZeroOptions(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { CountedField(c, "Bio", &value, CountedFieldOptions{}) }, 400, 200)
	if !tt.HasText("0/280") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestQuantityInputAtMaximum(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	value, changed := maxInt, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if QuantityInput(c, "Items", &value, QuantityInputOptions{Max: maxInt}).Changed() {
			changed++
		}
	}, 400, 200)
	before := baseline()
	if err := tt.Click("Increase Items"); err != nil {
		t.Fatal(err)
	}
	if value != maxInt || changed != 0 || !appended(before) {
		t.Errorf("value %d, changed %d, cues %v", value, changed, sound.Recent())
	}
}

func TestRegexAndJSONOverrides(t *testing.T) {
	pattern := "a+"
	ttRegex := ui.NewTester(func(c *ui.Context) {
		RegexField(c, &pattern, "aa", RegexFieldOptions{Matches: "hits"})
	}, 400, 200)
	if !ttRegex.HasText("1 hits") {
		t.Errorf("regex texts %q", ttRegex.Texts())
	}
	value := ""
	standard := ui.NewTester(func(c *ui.Context) {
		JSONField(c, &value, JSONFieldOptions{})
	}, 500, 300)
	custom := ""
	customField := ui.NewTester(func(c *ui.Context) {
		JSONField(c, &custom, JSONFieldOptions{Placeholder: "Paste JSON"})
	}, 500, 300)
	if !rectChanged(standard.Image(), customField.Image(), 0, 20, 400, 120) {
		t.Error("placeholder override did not change the field")
	}
}

func TestPortFieldLabelOverride(t *testing.T) {
	value := 80
	tt := ui.NewTester(func(c *ui.Context) { PortField(c, "Port", &value, PortFieldOptions{Label: "Public port"}) }, 400, 200)
	if !tt.HasText("Public port") {
		t.Errorf("texts %q", tt.Texts())
	}
}
