package input

import (
	"unicode"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PasswordFieldOptions tune a PasswordField.
type PasswordFieldOptions struct {
	// Placeholder is caller copy.
	Placeholder string
}

// PasswordField masks input, reveals it and shows strength.
func PasswordField(c *ui.Context, label string, value *string, o PasswordFieldOptions) ui.Element {
	if value == nil {
		panic("input: PasswordField needs a value")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	visible := ui.Local(root, "visible", func() bool { return false })
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		row := ui.Row(c).FillWidth().AlignItems(ui.Center).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		row.Children(func() {
			in := ui.TextInput(c, value).Grow(1).Placeholder(o.Placeholder).Padding(t.Space(2), t.Space(3))
			if !*visible {
				in.Password()
			}
			changed(root, in)
			button := ui.ButtonBase(c).Label("Reveal " + label).Padding(t.Space(2)).TextColor(t.TextMuted)
			kit.Press(button).Children(func() {
				name := "eye"
				if *visible {
					name = "eye-off"
				}
				ui.Icon(c, icon.Must(name)).Size(16, 16).ContentTransition(ui.ContentTransition{})
			})
			if kit.Clicked(button, sound.Toggle) {
				*visible = !*visible
				root.SetChanged()
				kit.Log("PasswordField", "toggled", "visible", *visible)
			}
		})
		score := passwordStrength(*value)
		ui.Row(c).Gap(t.Space(1)).Children(func() {
			for i := range 4 {
				bar := ui.Box(c).Key(i).Grow(1).Basis(0).Height(t.Space(1)).Radius(999)
				color := t.Border
				if i < score {
					color = t.Success
				}
				bar.Background(bar.AnimateColor("score", color, ui.Snappy))
			}
		})
	})
}

// passwordStrength scores length and character classes.
func passwordStrength(s string) int {
	score := 0
	if len([]rune(s)) >= 8 {
		score++
	}
	var lower, upper, digit, symbol bool
	for _, r := range s {
		lower = lower || unicode.IsLower(r)
		upper = upper || unicode.IsUpper(r)
		digit = digit || unicode.IsDigit(r)
		symbol = symbol || !unicode.IsLetter(r) && !unicode.IsDigit(r)
	}
	for _, on := range []bool{lower, upper, digit, symbol} {
		if on {
			score++
		}
	}
	return min(score, 4)
}
