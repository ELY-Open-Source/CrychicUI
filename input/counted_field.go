package input

import (
	"strconv"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CountedFieldOptions tune a CountedField.
type CountedFieldOptions struct {
	// Placeholder is caller copy; Max is the character limit.
	Placeholder string
	Max         int
}

// CountedField counts characters and warns near the limit.
func CountedField(c *ui.Context, label string, value *string, o CountedFieldOptions) ui.Element {
	if o.Max == 0 {
		o.Max = 280
	}
	if value == nil || o.Max < 0 {
		panic("input: CountedField needs a value and positive limit")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		in := ui.TextInput(c, value).FillWidth().Placeholder(o.Placeholder).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).
			Background(t.Background).Border(1, t.Border)
		changed(root, in)
		n := utf8.RuneCountInString(*value)
		color := t.TextMuted
		if n >= int(float64(o.Max)*0.9) {
			color = t.Warning
		}
		counter := kit.Muted(c, strconv.Itoa(n)+"/"+strconv.Itoa(o.Max)).AlignSelf(ui.End).FontSize(t.Rem(0.85))
		counter.TextColor(counter.AnimateColor("count", color, ui.Snappy))
	})
}
