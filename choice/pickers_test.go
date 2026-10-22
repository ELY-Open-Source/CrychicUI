package choice

import (
	"image"
	"image/color"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestYesNo(t *testing.T) {
	v := -1
	tt := ui.NewTester(func(c *ui.Context) { YesNo(c, "Attending", &v, YesNoOptions{No: "Nope"}) }, 240, 80)
	before := baseline()
	click(t, tt, "Nope")
	if v != 0 {
		t.Fatalf("value %d", v)
	}
	click(t, tt, "Yes")
	if v != 1 {
		t.Fatalf("value %d", v)
	}
	cues(t, before, sound.Select, sound.Select)
	settles(t, tt, "YesNo tint")
}

func TestImageAndColorChoice(t *testing.T) {
	img := ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 4, 4)))
	pic, col := 0, 0
	var th ui.Theme
	tt := ui.NewTester(func(c *ui.Context) {
		th = *c.Theme()
		ui.Column(c).Gap(8).Children(func() {
			ImageChoice(c, "Cover", []ImageItem{{"Dunes", img}, {"Sea", img}}, &pic)
			ColorChoice(c, "Color", []Swatch{{"Paper", th.Background}, {"Red", th.Danger}}, &col)
		})
	}, 500, 300)
	before := baseline()
	click(t, tt, "Sea")
	click(t, tt, "Red")
	if pic != 1 || col != 1 {
		t.Fatalf("pic %d col %d", pic, col)
	}
	cues(t, before, sound.Select, sound.Select)
	settles(t, tt, "rings")
	red := find(t, tt, "Red")
	shot := tt.Image()
	if p := shot.RGBAAt(int(red.X+red.W/2), int(red.Y+red.H/2)); !nearColor(p, th.Danger) {
		t.Errorf("swatch center %v, want %v", p, th.Danger)
	}
	if p := shot.RGBAAt(int(red.X+red.W/2), int(red.Y)+1); !nearColor(p, th.Accent) {
		t.Errorf("ring %v, want accent %v", p, th.Accent)
	}
	panics(t, "nil image", func() {
		ui.Render(func(c *ui.Context) { ImageChoice(c, "x", []ImageItem{{"A", nil}}, &pic) }, 50, 50, 1)
	})
}

func TestSizeChoiceStrikesCornerToCorner(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		SizeChoice(c, "Size", []string{"S", "L"}, []string{"L"}, new(string))
	}, 400, 80)
	r := find(t, tt, "L")
	x := int(r.X + 0.15*r.W)
	if a, b := tt.Image().RGBAAt(x, int(r.Y+0.85*r.H)), tt.Image().RGBAAt(x, int(r.Y+0.5*r.H)); a == b {
		t.Errorf("no diagonal strike at %d: %v matches the box %v", x, a, b)
	}
}

func TestSizeChoice(t *testing.T) {
	size := ""
	tt := ui.NewTester(func(c *ui.Context) {
		SizeChoice(c, "Size", []string{"S", "M", "L", "XL"}, []string{"L"}, &size)
	}, 400, 80)
	before := baseline()
	click(t, tt, "L")
	if size != "" {
		t.Fatalf("unavailable size chosen: %q", size)
	}
	click(t, tt, "XL")
	if size != "XL" {
		t.Fatalf("size %q", size)
	}
	cues(t, before, sound.Select)
	settles(t, tt, "SizeChoice box")
	panics(t, "unavailable selection", func() {
		bad := "L"
		ui.Render(func(c *ui.Context) { SizeChoice(c, "x", []string{"L"}, []string{"L"}, &bad) }, 50, 50, 1)
	})
}

func TestPlanChoice(t *testing.T) {
	plan := "Free"
	plans := []Plan{{"Free", "$0", "/mo", "For one."}, {"Pro", "$12", "/mo", "For teams."}}
	tt := ui.NewTester(func(c *ui.Context) { PlanChoice(c, "Plan", plans, &plan) }, 600, 200)
	find(t, tt, "$12")
	before := baseline()
	click(t, tt, "Pro")
	if plan != "Pro" {
		t.Fatalf("plan %q", plan)
	}
	cues(t, before, sound.Select)
	settles(t, tt, "PlanChoice border")
}

func TestChoiceMatrix(t *testing.T) {
	sel := []int{-1, -1}
	var cells []ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		ChoiceMatrix(c, []string{"Speed", "Price"}, []string{"Bad", "Good"}, &sel)
	}, 500, 200)
	good := find(t, tt, "Good")
	speed, price := find(t, tt, "Speed"), find(t, tt, "Price")
	cells = append(cells, ui.Rect{X: good.X + good.W/2, Y: speed.Y + speed.H/2}, ui.Rect{X: good.X + good.W/2, Y: price.Y + price.H/2})
	before := baseline()
	tt.ClickAt(cells[1].X, cells[1].Y)
	if !slices.Equal(sel, []int{-1, 1}) {
		t.Fatalf("selection %v", sel)
	}
	tt.Key(0, ui.KeyLeft)
	if !slices.Equal(sel, []int{-1, 0}) {
		t.Errorf("arrow gave %v", sel)
	}
	cues(t, before, sound.Select, sound.Select)
	settles(t, tt, "ChoiceMatrix dot")
	panics(t, "short answers", func() {
		short := []int{0}
		ui.Render(func(c *ui.Context) { ChoiceMatrix(c, []string{"A", "B"}, []string{"X"}, &short) }, 50, 50, 1)
	})
}

func TestRanking(t *testing.T) {
	items := []string{"Speed", "Price", "Support"}
	tt := ui.NewTester(func(c *ui.Context) { Ranking(c, "Priorities", &items) }, 300, 200)
	before := baseline()
	a, b := find(t, tt, "Speed"), find(t, tt, "Support")
	tt.Press(a.X+a.W/2, a.Y+a.H/2)
	for i := 1; i <= 8; i++ {
		tt.Move(a.X+a.W/2, a.Y+a.H/2+(b.Y-a.Y)*float32(i)/8)
	}
	tt.Release(b.X+b.W/2, b.Y+b.H/2)
	if !slices.Equal(items, []string{"Price", "Support", "Speed"}) {
		t.Fatalf("drag gave %v", items)
	}
	settles(t, tt, "Ranking slide")
	click(t, tt, "Support")
	tt.Key(ui.Alt, ui.KeyUp)
	if !slices.Equal(items, []string{"Support", "Price", "Speed"}) {
		t.Fatalf("Alt+Up gave %v", items)
	}
	cues(t, before, sound.Drop, sound.Drop)
}

func TestSlotPicker(t *testing.T) {
	slot := ""
	tt := ui.NewTester(func(c *ui.Context) {
		SlotPicker(c, "Time", []string{"9:00", "9:30", "10:00"}, []string{"9:30"}, &slot)
	}, 400, 80)
	before := baseline()
	click(t, tt, "9:30")
	click(t, tt, "10:00")
	if slot != "10:00" {
		t.Fatalf("slot %q", slot)
	}
	cues(t, before, sound.Select)
	settles(t, tt, "SlotPicker chip")
}

func TestOptionTiles(t *testing.T) {
	var sel []string
	tt := ui.NewTester(func(c *ui.Context) {
		OptionTiles(c, "Goals", []Tile{{"Fitness", "dumbbell"}, {"Sleep", "moon"}}, &sel)
	}, 400, 160)
	before := baseline()
	click(t, tt, "Sleep")
	click(t, tt, "Fitness")
	click(t, tt, "Sleep")
	if !slices.Equal(sel, []string{"Fitness"}) {
		t.Fatalf("selection %v", sel)
	}
	cues(t, before, sound.Toggle, sound.Toggle, sound.Toggle)
	settles(t, tt, "OptionTiles badge")
}

func TestPickN(t *testing.T) {
	sel := []string{"Go"}
	refused := 0
	o := PickNOptions{Refused: func() { refused++ }}
	tt := ui.NewTester(func(c *ui.Context) { PickN(c, "Languages", []string{"Go", "Rust", "Zig"}, 2, &sel, o) }, 400, 80)
	before := baseline()
	click(t, tt, "Rust")
	click(t, tt, "Zig")
	if !slices.Equal(sel, []string{"Go", "Rust"}) || refused != 1 {
		t.Fatalf("selection %v refused %d", sel, refused)
	}
	z := find(t, tt, "Zig")
	tt.Advance(100 * time.Millisecond)
	if moved := find(t, tt, "Zig"); moved.X == z.X {
		t.Error("refused chip did not shake")
	}
	cues(t, before, sound.Toggle, sound.Error)
	settles(t, tt, "PickN shake")
	panics(t, "n too large", func() {
		ui.Render(func(c *ui.Context) { PickN(c, "x", []string{"A"}, 2, &sel, PickNOptions{}) }, 50, 50, 1)
	})
}

func nearColor(p color.RGBA, want ui.Color) bool {
	d := func(a, b uint8) bool { return a-b < 8 || b-a < 8 }
	return d(p.R, want.R) && d(p.G, want.G) && d(p.B, want.B)
}
