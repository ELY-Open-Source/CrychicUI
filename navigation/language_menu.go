package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// LanguageMenu lists languages with a check on the chosen one.
func LanguageMenu(c *ui.Context, label string, langs []string, chosen *int, open *bool) ui.Element {
	if label == "" || len(langs) == 0 || chosen == nil || open == nil || !valid(*chosen, len(langs)) {
		panic(fmt.Sprintf("navigation: LanguageMenu %q needs languages, a choice and open", label))
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	was, shown := *open, *chosen
	root.Children(func() {
		trigger := kit.Press(ui.ButtonBase(c).Label(label).Gap(t.Space(1.5)).Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius)).Expanded(*open)
		trigger.Background(trigger.AnimateColor("bg", t.SurfaceHover.Alpha(unit(*open || trigger.Hovered())), ui.Snappy)).Children(func() {
			ui.Icon(c, icon.Must("globe")).Size(16, 16)
			ui.Text(c, langs[*chosen]).SingleLine().ContentTransition(ui.ContentTransition{})
		})
		popup(c, trigger, open, func() {
			ui.Column(c).Width(t.Space(48)).Gap(t.Space(0.5)).Children(func() {
				for i, lang := range langs {
					c.Key(i)
					if kit.Clicked(checked(c, lang, i == shown, nil), sound.Select) {
						*chosen, *open = i, false
						root.SetChanged()
						kit.Log("LanguageMenu", "chose", "index", i)
					}
				}
			})
		})
	})
	reopened(root, was, *open, "LanguageMenu")
	return root
}
