package input

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/richtext"
)

func TestRichEditor(t *testing.T) {
	doc := richtext.NewDocument(richtext.NewText("Hello"))
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if RichEditor(c, doc).Changed() {
			changed++
		}
	}, 500, 300)
	before := baseline()
	if err := tt.Click("Bold"); err != nil {
		t.Fatal(err)
	}
	if doc.Formatting().Bold != richtext.FormatOn || changed != 1 || !appended(before, sound.Toggle) {
		t.Errorf("bold %v, changed %d", doc.Formatting().Bold, changed)
	}
}

func TestMarkdownEditor(t *testing.T) {
	source, changed := "", 0
	tt := ui.NewTester(func(c *ui.Context) {
		if MarkdownEditor(c, &source, MarkdownEditorOptions{}).Changed() {
			changed++
		}
	}, 700, 300)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Type("# Hello")
	settles(t, tt, "MarkdownEditor preview")
	if source != "# Hello" || changed == 0 || !tt.HasText("Hello") || !appended(before, sound.Tick) {
		t.Errorf("source %q, changed %d, texts %q", source, changed, tt.Texts())
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("# Goodbye")
	if !tt.FramePending() {
		t.Error("second edit did not cross-fade")
	}
	settles(t, tt, "MarkdownEditor second preview")
	if !tt.HasText("Goodbye") {
		t.Error("second preview missing")
	}
}

func TestNoteInput(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { NoteInput(c, &value, "Write a note") }, 500, 300)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	settles(t, tt, "NoteInput placeholder")
	tt.Type("A note")
	if value != "A note" || !appended(before, sound.Tick) {
		t.Errorf("value %q", value)
	}
}

func TestQuantityInput(t *testing.T) {
	value, changed := 1, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if QuantityInput(c, "Items", &value, QuantityInputOptions{}).Changed() {
			changed++
		}
	}, 400, 200)
	before := baseline()
	if err := tt.Click("Increase Items"); err != nil {
		t.Fatal(err)
	}
	settles(t, tt, "QuantityInput value")
	if value != 2 || changed != 1 || !appended(before, sound.Tick) {
		t.Errorf("value %d, changed %d", value, changed)
	}
}

func TestIPField(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { IPField(c, &value, IPFieldOptions{}) }, 400, 200)
	tt.Key(0, ui.KeyTab)
	tt.Type("192168001001")
	if value != "192.168.001.001" || tt.HasText("Each octet") {
		t.Errorf("value %q, texts %q", value, tt.Texts())
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("9999")
	if !tt.HasText("Each octet must be 0 through 255.") {
		t.Error("IP range error missing")
	}
}

func TestPortField(t *testing.T) {
	value := 80
	tt := ui.NewTester(func(c *ui.Context) { PortField(c, "Port", &value, PortFieldOptions{}) }, 400, 200)
	tt.Key(0, ui.KeyTab)
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("70000")
	if value != 80 || !tt.HasText("Use a port from 1 through 65535.") {
		t.Errorf("value %d, texts %q", value, tt.Texts())
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("8080")
	if value != 8080 {
		t.Errorf("value %d", value)
	}
}

func TestMaskedField(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) { MaskedField(c, "Serial", &value, "###-##", MaskedFieldOptions{}) }, 400, 200)
	tt.Key(0, ui.KeyTab)
	tt.Type("123456")
	if value != "123-45" {
		t.Errorf("value %q", value)
	}
}

func TestTemplateField(t *testing.T) {
	value := "Hello "
	tt := ui.NewTester(func(c *ui.Context) {
		TemplateField(c, "Template", &value, []string{"name", "date"}, TemplateFieldOptions{})
	}, 500, 300)
	before := baseline()
	if err := tt.Click("Insert variable"); err != nil {
		t.Fatal(err)
	}
	if err := tt.ChooseMenuItem("name"); err != nil {
		t.Fatal(err)
	}
	if value != "Hello {name}" || !appended(before, sound.Select) {
		t.Errorf("value %q", value)
	}
}

func TestAutosizeInput(t *testing.T) {
	value := ""
	width := float32(0)
	tt := ui.NewTester(func(c *ui.Context) {
		width = AutosizeInput(c, &value, AutosizeInputOptions{Placeholder: "Name"}).LayoutBounds().W
	}, 500, 200)
	tt.Frame()
	initial := width
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Type("A much longer name")
	tt.Frame()
	settles(t, tt, "AutosizeInput width")
	if width <= initial || !appended(before, sound.Tick) {
		t.Errorf("width %v, initial %v", width, initial)
	}
}

func TestInputGroup(t *testing.T) {
	host, port, clicked := "", "", 0
	tt := ui.NewTester(func(c *ui.Context) {
		_, action := InputGroup(c, []string{"Host", "Port"}, []*string{&host, &port}, "Connect")
		if action.Clicked() {
			clicked++
		}
	}, 500, 200)
	tt.Key(0, ui.KeyTab)
	tt.Type("localhost")
	tt.Key(0, ui.KeyTab)
	tt.Type("8080")
	if host != "localhost" || port != "8080" {
		t.Errorf("%q %q", host, port)
	}
	before := baseline()
	if err := tt.Click("Connect"); err != nil {
		t.Fatal(err)
	}
	if clicked != 1 || !appended(before, sound.Tap) {
		t.Errorf("clicked %d", clicked)
	}
}

func TestFileNameField(t *testing.T) {
	name := ""
	tt := ui.NewTester(func(c *ui.Context) { FileNameField(c, "File", &name, "pdf", FileNameFieldOptions{}) }, 400, 200)
	tt.Key(0, ui.KeyTab)
	tt.Type("report")
	if name != "report" || !tt.HasText(".pdf") {
		t.Errorf("name %q, texts %q", name, tt.Texts())
	}
}

func TestConfirmPassword(t *testing.T) {
	password, confirm := "", ""
	tt := ui.NewTester(func(c *ui.Context) { ConfirmPassword(c, &password, &confirm, ConfirmPasswordOptions{}) }, 400, 240)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Type("secret")
	tt.Key(0, ui.KeyTab)
	tt.Type("secret")
	if !tt.HasText("Confirm password") || !appended(before, sound.Success) {
		t.Errorf("texts %q, cues %v", tt.Texts(), sound.Recent())
	}
}
