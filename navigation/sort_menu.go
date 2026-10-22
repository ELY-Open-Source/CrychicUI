package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SortMenuOptions tune a SortMenu.
type SortMenuOptions struct {
	Ascending, Descending string
}

// SortMenu picks a sort key and direction; the arrow turns.
func SortMenu(c *ui.Context, label string, keys []string, key *int, desc, open *bool, o SortMenuOptions) ui.Element {
	if label == "" || len(keys) == 0 || key == nil || desc == nil || open == nil || !valid(*key, len(keys)) {
		panic(fmt.Sprintf("navigation: SortMenu %q needs keys, a key, a direction and open", label))
	}
	def(&o.Ascending, "Ascending")
	def(&o.Descending, "Descending")
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	was, shownKey, shownDesc := *open, *key, *desc
	pick := func(event string, value int) {
		*open = false
		root.SetChanged()
		kit.Log("SortMenu", event, "value", value)
	}
	root.Children(func() {
		trigger := kit.Press(ui.ButtonBase(c).Label(label).Gap(t.Space(1.5)).Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius)).Expanded(*open)
		trigger.Background(trigger.AnimateColor("bg", t.SurfaceHover.Alpha(unit(*open || trigger.Hovered())), ui.Snappy)).Children(func() {
			ui.Text(c, keys[*key]).SingleLine().ContentTransition(ui.ContentTransition{})
			arrow := ui.Icon(c, icon.Must("arrow-up")).Size(16, 16)
			arrow.Rotate(180 * arrow.AnimateTo("turn", unit(*desc), ui.Snappy))
		})
		popup(c, trigger, open, func() {
			ui.Column(c).Width(t.Space(48)).Gap(t.Space(0.5)).Children(func() {
				for i, k := range keys {
					c.Key(i)
					if kit.Clicked(checked(c, k, i == shownKey, nil), sound.Select) {
						*key = i
						pick("key", i)
					}
				}
				ui.Box(c).Key("rule").FillWidth().Height(1).Background(t.Border)
				for i, name := range []string{o.Ascending, o.Descending} {
					c.Key(fmt.Sprint("dir", i))
					if kit.Clicked(checked(c, name, shownDesc == (i == 1), nil), sound.Select) {
						*desc = i == 1
						pick("direction", i)
					}
				}
			})
		})
	})
	reopened(root, was, *open, "SortMenu")
	return root
}
