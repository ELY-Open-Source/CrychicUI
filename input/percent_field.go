package input

import (
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PercentField is a percent suffix clamped from zero to 100.
func PercentField(c *ui.Context, label string, value *float64) ui.Element {
	if value == nil || !finite(*value) {
		panic("input: PercentField needs a finite value")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	text := ui.Local(root, "text", func() string { return strconv.FormatFloat(*value, 'f', -1, 64) })
	step := ui.Local(root, "step", func() int { return 0 })
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		row := ui.Row(c).FillWidth().AlignItems(ui.Center).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		row.Children(func() {
			in := ui.TextInput(c, text).Grow(1).Label(label).Padding(t.Space(2), 0, t.Space(2), t.Space(3))
			in.OnTextKey(func(k ui.TextKey) bool {
				if k.Key == ui.KeyUp {
					*step = 1
					return true
				}
				if k.Key == ui.KeyDown {
					*step = -1
					return true
				}
				return false
			})
			if in.Changed() {
				parsed, err := strconv.ParseFloat(strings.TrimSpace(*text), 64)
				if err == nil && finite(parsed) {
					next := min(100, max(0, parsed))
					if next != *value {
						*value = next
						root.SetChanged()
						sound.Play(sound.Tick)
					}
				}
			}
			if *step != 0 {
				next := min(100, max(0, *value+float64(*step)))
				*step = 0
				if next != *value {
					*value = next
					*text = strconv.FormatFloat(*value, 'f', -1, 64)
					root.SetChanged()
					sound.Play(sound.Tick)
				}
			}
			if !in.Focused() {
				formatted := strconv.FormatFloat(*value, 'f', -1, 64)
				if *text != formatted {
					*text = formatted
				}
			}
			kit.Muted(c, "%").Padding(t.Space(2), t.Space(3), t.Space(2), t.Space(1))
		})
	})
}
