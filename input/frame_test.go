package input

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// strokes lists border rows down the column x.
func strokes(tt *ui.Tester, x int, border ui.Color) []int {
	img := tt.Image()
	var ys []int
	on := false
	for y := range img.Bounds().Dy() {
		hit := near(img, x, y, border)
		if hit && !on {
			ys = append(ys, y)
		}
		on = hit
	}
	return ys
}

func TestFieldsDrawOneFrame(t *testing.T) {
	var s string
	var f float64
	var locked bool
	fields := map[string]func(c *ui.Context){
		"SearchInput":     func(c *ui.Context) { SearchInput(c, &s, SearchInputOptions{}) },
		"ClearableField":  func(c *ui.Context) { ClearableField(c, "Name", &s, ClearableFieldOptions{}) },
		"CardNumberField": func(c *ui.Context) { CardNumberField(c, &s, CardNumberFieldOptions{}) },
		"AdornedField":    func(c *ui.Context) { AdornedField(c, "Price", &s, nil, nil, AdornedFieldOptions{}) },
		"PasswordField":   func(c *ui.Context) { PasswordField(c, "Password", &s, PasswordFieldOptions{}) },
		"FileNameField":   func(c *ui.Context) { FileNameField(c, "File", &s, "txt", FileNameFieldOptions{}) },
		"PercentField":    func(c *ui.Context) { PercentField(c, "Rate", &f) },
		"PhoneField":      func(c *ui.Context) { PhoneField(c, "Phone", &s, PhoneFieldOptions{}) },
		"SlugField":       func(c *ui.Context) { SlugField(c, "Title", &s, &locked, SlugFieldOptions{}) },
		"URLField":        func(c *ui.Context) { URLField(c, "Site", &s, URLFieldOptions{}) },
	}
	for name, field := range fields {
		s, f, locked = "", 0, false
		var th *ui.Theme
		tt := ui.NewTester(func(c *ui.Context) {
			th = c.Theme()
			ui.Column(c).Width(400).Padding(10).Children(func() { field(c) })
		}, 420, 160)
		tt.Frame()
		ys := strokes(tt, 240, th.Border)
		if len(ys) < 2 {
			t.Fatalf("%s: frame not found, strokes at %v", name, ys)
		}
		mid := (ys[0] + ys[1]) / 2
		if p := tt.Image().RGBAAt(240, mid); !near(tt.Image(), 240, mid, th.Background) {
			t.Errorf("%s: input fill %v inside the frame, want %v", name, p, th.Background)
		}
	}
}
