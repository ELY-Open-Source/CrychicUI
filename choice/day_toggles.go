package choice

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DayTogglesOptions tune DayToggles.
type DayTogglesOptions struct {
	// Letters show and Names name the days, Sunday first.
	Letters, Names [7]string
}

// DayToggles is weekday letters, each toggled on or off.
func DayToggles(c *ui.Context, label string, days *[7]bool, o DayTogglesOptions) ui.Element {
	if label == "" || days == nil {
		panic("choice: DayToggles needs a label and days")
	}
	if o.Letters == ([7]string{}) {
		o.Letters = [7]string{"S", "M", "T", "W", "T", "F", "S"}
	}
	if o.Names == ([7]string{}) {
		o.Names = [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	}
	distinct("DayToggles", o.Names[:])
	t := c.Theme()
	root := ui.Row(c).Gap(t.Space(1.5)).AlignSelf(ui.Start).Label(label).FocusGroup(ui.Horizontal)
	return root.Children(func() {
		for i := range days {
			c.Key(i)
			b := ui.ToggleBase(c, &days[i]).Label(o.Names[i]).Tooltip(o.Names[i]).Size(t.Space(9), t.Space(9)).Radius(999)
			bg, fg := t.Surface, t.Text
			if days[i] {
				bg, fg = t.Accent, t.AccentText
			} else if b.Hovered() {
				bg = t.SurfaceHover
			}
			b.Background(b.AnimateColor("bg", bg, quick)).TextColor(b.AnimateColor("fg", fg, quick))
			kit.Press(b).Children(func() { ui.Text(c, o.Letters[i]).Bold() })
			if kit.Changed(b, sound.Toggle) {
				n := 0
				for _, on := range days {
					if on {
						n++
					}
				}
				root.SetChanged()
				kit.Log("DayToggles", "toggled", "count", n)
			}
		}
	})
}
