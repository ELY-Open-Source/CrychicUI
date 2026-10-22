package gallery

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/color"
	"github.com/egoist/mygo/ui"
)

var colorEntries = []Entry{
	{536, "Picker", "color", "A color field opening a picker.", demoPicker},
	{537, "Well", "color", "A swatch button opening a picker.", demoColorWell},
	{538, "HexField", "color", "Hex input with a swatch.", demoHexField},
	{539, "RGBFields", "color", "Red, green and blue fields.", demoRGBFields},
	{540, "HSLFields", "color", "Hue, saturation and lightness fields.", demoHSLFields},
	{541, "OklchFields", "color", "Oklch fields with a gamut note.", demoOklchFields},
	{542, "HueRing", "color", "A ring dragged for hue.", demoHueRing},
	{543, "ShadePad", "color", "Saturation and value on a pad.", demoShadePad},
	{544, "AlphaSlider", "color", "Opacity over a checkered track.", demoAlphaSlider},
	{545, "SwatchGrid", "color", "A grid of swatches; the chosen rings.", demoSwatchGrid},
	{546, "Palette", "color", "Shades of one hue, light to dark.", demoPalette},
	{547, "Contrast", "color", "Contrast ratio with AA and AAA marks.", demoContrast},
	{548, "GradientEditor", "color", "Stops dragged and added on a bar.", demoGradientEditor},
	{549, "RecentColors", "color", "Last colors chosen, newest first.", demoRecentColors},
	{550, "ColorName", "color", "The nearest named color.", demoColorName},
	{551, "Harmony", "color", "Complement, triad and analogous.", demoHarmony},
	{552, "BeforeAfter", "color", "Two swatches split by an arrow.", demoBeforeAfter},
	{553, "ThemePreview", "color", "A small window drawn in a theme.", demoThemePreview},
	{554, "Eyedropper", "color", "Pick a color from a strip.", demoEyedropper},
	{555, "ColorBlind", "color", "Colors as color-blind eyes see them.", demoColorBlind},
}

// swatches is a demo palette.
var swatches = []ui.Color{ui.Hex("#EF4444"), ui.Hex("#F97316"), ui.Hex("#EAB308"), ui.Hex("#22C55E"),
	ui.Hex("#06B6D4"), ui.Hex("#3B82F6"), ui.Hex("#8B5CF6"), ui.Hex("#EC4899")}

// tint is a demo color kept across builds.
func tint(root ui.Element) *ui.Color {
	return ui.Local(root, "tint", func() ui.Color { return ui.Hex("#3B82F6") })
}

// hex writes a demo log line for a color.
func hex(col ui.Color) string { return fmt.Sprintf("#%02X%02X%02X%02X", col.R, col.G, col.B, col.A) }

func demoPicker(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	v := tint(root)
	root.Children(func() {
		if color.Picker(c, "Brand", v).Changed() {
			log(hex(*v))
		}
	})
}

func demoColorWell(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	v := tint(root)
	root.Children(func() {
		if color.Well(c, "Accent", v).Changed() {
			log(hex(*v))
		}
	})
}

func demoHexField(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	v := tint(root)
	root.Children(func() {
		if color.HexField(c, "Hex color", v, color.HexFieldOptions{}).Changed() {
			log(hex(*v))
		}
	})
}

func demoRGBFields(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	v := tint(root)
	root.Children(func() {
		if color.RGBFields(c, "Color", v, color.RGBFieldsOptions{}).Changed() {
			log(hex(*v))
		}
	})
}

func demoHSLFields(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	v := tint(root)
	root.Children(func() {
		if color.HSLFields(c, "Color", v, color.HSLFieldsOptions{}).Changed() {
			log(hex(*v))
		}
	})
}

func demoOklchFields(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	v := tint(root)
	root.Children(func() {
		if color.OklchFields(c, "Color", v, color.OklchFieldsOptions{}).Changed() {
			log(hex(*v))
		}
	})
}

func demoHueRing(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	h := ui.Local(root, "hue", func() float64 { return 210 })
	root.Children(func() {
		if color.HueRing(c, "Hue", h).Changed() {
			log(fmt.Sprintf("hue %.0f", *h))
		}
	})
}

func demoShadePad(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	v := tint(root)
	root.Children(func() {
		if color.ShadePad(c, "Shade", v).Changed() {
			log(hex(*v))
		}
	})
}

func demoAlphaSlider(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	v := tint(root)
	root.Children(func() {
		if color.AlphaSlider(c, "Opacity", v).Changed() {
			log(fmt.Sprintf("alpha %d", v.A))
		}
	})
}

func demoSwatchGrid(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	v := tint(root)
	root.Children(func() {
		if color.SwatchGrid(c, "Swatches", swatches, v, color.SwatchGridOptions{Columns: 4}).Changed() {
			log(hex(*v))
		}
	})
}

func demoPalette(c *ui.Context, log func(string)) {
	color.Palette(c, ui.Hex("#3B82F6"), color.PaletteOptions{})
}

func demoContrast(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	fg := ui.Local(root, "fg", func() ui.Color { return ui.Hex("#767676") })
	root.Children(func() {
		color.Contrast(c, *fg, ui.Hex("#FFFFFF"), color.ContrastOptions{})
		if button.Outline(c, "Darker text").Clicked() {
			*fg = fg.Mix(ui.Hex("#000000"), 0.3)
			log(hex(*fg))
		}
	})
}

func demoGradientEditor(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	stops := ui.Local(root, "stops", func() []ui.GradientStop {
		return []ui.GradientStop{{Color: ui.Hex("#3B82F6"), At: 0}, {Color: ui.Hex("#EC4899"), At: 1}}
	})
	root.Children(func() {
		if color.GradientEditor(c, "Fill", stops, color.GradientEditorOptions{}).Changed() {
			log(fmt.Sprintf("%d stops", len(*stops)))
		}
	})
}

func demoRecentColors(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	v := tint(root)
	recent := ui.Local(root, "recent", func() []ui.Color { return swatches[5:] })
	next := ui.Local(root, "next", func() int { return 0 })
	root.Children(func() {
		if color.RecentColors(c, "Recent", *recent, v, color.RecentColorsOptions{}).Changed() {
			log(hex(*v))
		}
		if button.Outline(c, "Choose another").Clicked() {
			*recent = append([]ui.Color{swatches[*next%5]}, *recent...)
			*next++
			log(fmt.Sprintf("%d recent", len(*recent)))
		}
	})
}

func demoColorName(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	i := ui.Local(root, "i", func() int { return 0 })
	names := []color.NamedColor{{Name: "Tomato", Color: ui.Hex("#FF6347")}, {Name: "Gold", Color: ui.Hex("#FFD700")},
		{Name: "Sea green", Color: ui.Hex("#2E8B57")}, {Name: "Royal blue", Color: ui.Hex("#4169E1")},
		{Name: "Orchid", Color: ui.Hex("#DA70D6")}, {Name: "Slate gray", Color: ui.Hex("#708090")}}
	root.Children(func() {
		color.ColorName(c, swatches[*i%len(swatches)], names)
		if button.Outline(c, "Next color").Clicked() {
			*i++
			log(hex(swatches[*i%len(swatches)]))
		}
	})
}

func demoHarmony(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	v := tint(root)
	root.Children(func() {
		if color.Harmony(c, "Harmony", ui.Hex("#3B82F6"), v, color.HarmonyOptions{}).Changed() {
			log(hex(*v))
		}
	})
}

func demoBeforeAfter(c *ui.Context, log func(string)) {
	color.BeforeAfter(c, ui.Hex("#94A3B8"), ui.Hex("#3B82F6"), color.BeforeAfterOptions{})
}

func demoThemePreview(c *ui.Context, log func(string)) {
	root := ui.Row(c).Gap(12).AlignItems(ui.Center)
	dark := ui.Local(root, "dark", func() bool { return false })
	root.Children(func() {
		th := ui.LightTheme()
		if *dark {
			th = ui.DarkTheme()
		}
		color.ThemePreview(c, "Theme", th)
		if button.Outline(c, "Switch theme").Clicked() {
			*dark = !*dark
			log(fmt.Sprintf("dark %v", *dark))
		}
	})
}

func demoEyedropper(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	v := tint(root)
	root.Children(func() {
		if color.Eyedropper(c, "Strip", swatches, v).Changed() {
			log(hex(*v))
		}
	})
}

func demoColorBlind(c *ui.Context, log func(string)) {
	color.ColorBlind(c, swatches, color.ColorBlindOptions{})
}
