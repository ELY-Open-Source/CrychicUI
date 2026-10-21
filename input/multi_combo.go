package input

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MultiCombo adds chosen items as removable chips.
func MultiCombo(c *ui.Context, label string, chosen *[]string, options []string) ui.Element {
	if chosen == nil || len(options) == 0 {
		panic("input: MultiCombo needs chosen items and options")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	draft := ui.Local(root, "draft", func() string { return "" })
	return root.Children(func() {
		if len(*chosen) > 0 {
			ui.Row(c).Wrap().Gap(t.Space(2)).Children(func() {
				for i, item := range *chosen {
					chip := ui.Row(c).Key(i).Gap(t.Space(1)).Padding(t.Space(1), t.Space(2)).Radius(999).
						Background(t.Accent.Alpha(0.12)).TextColor(t.Accent)
					chip.Children(func() {
						ui.Text(c, item)
						remove := ui.ButtonBase(c).Label("Remove " + item).TextColor(t.Accent)
						kit.Press(remove).Children(func() { ui.Text(c, "×") })
						if kit.Clicked(remove, sound.Delete) {
							*chosen = slices.Delete(*chosen, i, i+1)
							root.SetChanged()
							kit.Log("MultiCombo", "removed", "count", len(*chosen))
						}
					})
				}
			})
		}
		combo := ui.Combobox(c, draft, options).Label(label).FillWidth()
		if combo.Changed() && slices.Contains(options, *draft) {
			*chosen = append(*chosen, *draft)
			*draft = ""
			root.SetChanged()
			sound.Play(sound.Pop)
			kit.Log("MultiCombo", "added", "count", len(*chosen))
		}
	})
}
