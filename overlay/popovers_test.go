package overlay

import (
	"strings"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestCalloutPopover(t *testing.T) {
	th := ui.LightTheme()
	open, changed := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		a := ui.ButtonBase(c).Label("Tip").Size(80, 32)
		if a.Clicked() {
			open = !open
		}
		if CalloutPopover(c, a, &open, feedback.ToneDanger, "Careful", "This cannot be undone.").Changed() {
			changed++
		}
	}, 500, 400)
	before := baseline()
	click(t, tt, "Tip")
	settles(t, tt, "callout")
	r := find(t, tt, "Careful")
	edge := 0
	for y := int(r.Y) + 4; y < int(r.Y+r.H)-4; y++ {
		if p := tt.Image().RGBAAt(int(r.X)+2, y); p.R == th.Danger.R && p.G == th.Danger.G && p.B == th.Danger.B {
			edge++
		}
	}
	if edge < int(r.H)/2 {
		t.Errorf("danger edge %d px of %v", edge, r.H)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open || changed != 1 {
		t.Errorf("Escape: open %v, changed %d", open, changed)
	}
	cues(t, before, sound.Open)
	panics(t, "tone 9", func() {
		render(func(c *ui.Context) { o := true; CalloutPopover(c, ui.Box(c), &o, 9, "x", "") })
	})
}

func TestInfoPopover(t *testing.T) {
	open := false
	tt := ui.NewTester(func(c *ui.Context) { InfoPopover(c, "About tax", "Prices include VAT.", &open) }, 500, 400)
	before := baseline()
	click(t, tt, "About tax")
	if !open || !tt.HasText("Prices include VAT.") {
		t.Fatalf("open %v, %q", open, tt.Texts())
	}
	settles(t, tt, "info")
	click(t, tt, "About tax")
	if open {
		t.Error("second click did not close")
	}
	cues(t, before, sound.Open)
}

func TestConfirmPopover(t *testing.T) {
	open, yes := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		a := ui.ButtonBase(c).Label("Delete file").Size(120, 32)
		if a.Clicked() {
			open = !open
		}
		if _, y := ConfirmPopover(c, a, &open, ConfirmPopoverOptions{Question: "Delete it?"}); y.Clicked() {
			yes++
		}
	}, 500, 400)
	before := baseline()
	click(t, tt, "Delete file")
	settles(t, tt, "confirm")
	click(t, tt, "No")
	if open || yes != 0 {
		t.Fatalf("No: open %v, yes %d", open, yes)
	}
	click(t, tt, "Delete file")
	if !tt.HasText("Delete it?") {
		t.Fatalf("question %q", tt.Texts())
	}
	settles(t, tt, "confirm again")
	click(t, tt, "Yes")
	if open || yes != 1 {
		t.Errorf("Yes: open %v, yes %d", open, yes)
	}
	cues(t, before, sound.Tap, sound.Tap)
}

func TestSelectMenu(t *testing.T) {
	chosen, open := -1, false
	tt := ui.NewTester(func(c *ui.Context) {
		SelectMenu(c, "Color", []string{"Red", "Green", "Blue"}, &chosen, &open, SelectMenuOptions{})
	}, 500, 400)
	if !tt.HasText("Select") {
		t.Fatalf("placeholder %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Color")
	settles(t, tt, "menu")
	tt.Type("bl")
	tt.Frame()
	if _, ok := tt.Find("Red"); ok {
		t.Fatalf("filter bl kept Red: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if chosen != 2 || open || !tt.HasText("Blue") {
		t.Errorf("Enter: chosen %d, open %v, %q", chosen, open, tt.Texts())
	}
	cues(t, before, sound.Select)
}

func TestMentionMenu(t *testing.T) {
	text, chosen, open := "", -1, false
	tt := ui.NewTester(func(c *ui.Context) {
		field := ui.TextInput(c, &text).Label("Message").Width(300)
		if field.Changed() && strings.HasSuffix(text, "@") {
			open = true
		}
		q := text[strings.LastIndex(text, "@")+1:]
		MentionMenu(c, "People", field, q, []string{"Ada", "Adam", "Bob"}, &chosen, &open, MentionMenuOptions{})
	}, 500, 400)
	click(t, tt, "Message")
	tt.Type("hi @")
	tt.Frame()
	if !open || !tt.HasText("Bob") {
		t.Fatalf("@ did not open: %q", tt.Texts())
	}
	tt.Type("ad")
	tt.Frame()
	if _, ok := tt.Find("Bob"); ok {
		t.Fatalf("filter kept Bob: %q", tt.Texts())
	}
	before := baseline()
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if chosen != 1 || open {
		t.Fatalf("pick: chosen %d, open %v", chosen, open)
	}
	cues(t, before, sound.Select)
	tt.Type(" @")
	tt.Frame()
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open || !strings.HasSuffix(text, " @") {
		t.Errorf("Escape: open %v, text %q", open, text)
	}
}

func TestSlashMenu(t *testing.T) {
	text, chosen, open := "", -1, false
	commands := []MenuItem{{"Heading", "heading", ""}, {"Image", "image", ""}, {"Table", "table", ""}}
	tt := ui.NewTester(func(c *ui.Context) {
		field := ui.TextInput(c, &text).Label("Note").Width(300)
		if field.Changed() && strings.HasPrefix(text, "/") {
			open = true
		}
		SlashMenu(c, "Blocks", field, strings.TrimPrefix(text, "/"), commands, &chosen, &open, SlashMenuOptions{})
	}, 500, 400)
	click(t, tt, "Note")
	tt.Type("/ta")
	tt.Frame()
	r := find(t, tt, "Table")
	// A press sliding off a row keeps typing here.
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	tt.Move(r.X+r.W/2, r.Y-40)
	tt.Release(r.X+r.W/2, r.Y-40)
	tt.Frame()
	tt.Type("b")
	tt.Frame()
	if text != "/tab" || !open {
		t.Fatalf("slid-off press: text %q, open %v", text, open)
	}
	tt.Key(0, ui.KeyBackspace)
	tt.Frame()
	before := baseline()
	r = find(t, tt, "Table")
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	if chosen != 2 || open {
		t.Fatalf("click: chosen %d, open %v", chosen, open)
	}
	cues(t, before, sound.Select)
	tt.Type("x")
	tt.Frame()
	if text != "/tax" {
		t.Errorf("the field lost focus to the menu: %q", text)
	}
}
