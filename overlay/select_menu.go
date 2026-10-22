package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SelectMenuOptions tune a SelectMenu.
type SelectMenuOptions struct {
	// Placeholder shows while nothing is chosen; Empty when nothing matches.
	Placeholder, Search, Empty string
}

// SelectMenu is a button opening a searchable list of choices.
func SelectMenu(c *ui.Context, label string, options []string, chosen *int, open *bool, o SelectMenuOptions) ui.Element {
	if label == "" || len(options) == 0 || chosen == nil || open == nil || *chosen < -1 || *chosen >= len(options) {
		panic(fmt.Sprintf("overlay: SelectMenu %q needs options, a choice or -1, and open", label))
	}
	for _, s := range options {
		if s == "" {
			panic(fmt.Sprintf("overlay: SelectMenu %q options need text", label))
		}
	}
	def(&o.Placeholder, "Select")
	def(&o.Search, "Search")
	def(&o.Empty, "No matches")
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	query := searchText(root, *open)
	field := opening(c, root, *open)
	root.Children(func() {
		trigger := kit.Press(ui.ButtonBase(c).Label(label).MinWidth(t.Space(48)).Gap(t.Space(2)).AlignItems(ui.Center).
			Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).Border(1, t.Border)).Expanded(*open)
		trigger.Background(trigger.AnimateColor("bg", t.SurfaceHover.Alpha(unit(*open || trigger.Hovered())), ui.Snappy)).Children(func() {
			if *chosen < 0 {
				kit.Muted(c, o.Placeholder).Grow(1).SingleLine()
			} else {
				ui.Text(c, options[*chosen]).Grow(1).SingleLine().ContentTransition(ui.ContentTransition{})
			}
			chev := ui.Icon(c, icon.Must("chevrons-up-down")).Size(16, 16).TextColor(t.TextMuted)
			chev.Rotate(180 * chev.AnimateTo("turn", unit(*open), ui.Snappy))
		})
		if trigger.Clicked() {
			*open = !*open
			root.SetChanged()
		}
		if ui.PopoverBase(c, trigger, open, func(panel ui.Element) {
			panel.Margin(t.Space(1), 0, 0, 0).Width(t.Space(64)).Padding(t.Space(1)).Gap(t.Space(1)).Radius(t.Radius).
				Background(t.Background).Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).TransformOrigin(0.5, 0)
			kit.Arrive(panel, ui.Pose{Scale: 0.95, Fade: 1}, 0)
			// A plain input: SearchField's wrapper cannot take focus.
			search := ui.TextInput(c, query).Label(label).Placeholder(o.Search).FillWidth().Bind(field)
			results(c, root, search, *query, len(options), func(i int) string { return options[i] }, func(int) string { return "" }, nil,
				func(i int) {
					*chosen, *open = i, false
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("SelectMenu", "chose", "index", i)
				}, o.Empty)
		}).Dismissed() {
			root.SetChanged()
		}
	})
	moved(root, *open, "SelectMenu")
	return root
}
