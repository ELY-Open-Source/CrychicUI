package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/motion"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TabBarBadges are tabs with count badges whose digits roll.
func TabBarBadges(c *ui.Context, label string, names []string, counts []int, active *int) ui.Element {
	tabs("TabBarBadges", label, names, active)
	if len(counts) != len(names) {
		panic(fmt.Sprintf("navigation: TabBarBadges %q needs a count per tab", label))
	}
	for _, n := range counts {
		if n < 0 {
			panic(fmt.Sprintf("navigation: TabBarBadges %q counts must not be negative", label))
		}
	}
	t := c.Theme()
	shown := *active
	parts := ui.TabsBase(c, active, len(names))
	root := parts.List.Label(label).Gap(t.Space(1))
	id := root.ID()
	root.Children(func() {
		for i, name := range names {
			c.Key(i)
			underlined(c, parts.Tab(i), name, i == shown, id, func() {
				badge := ui.Row(c).Padding(0, t.Space(1.5)).Radius(999).FontSize(t.Rem(0.75))
				bg, fg := t.Surface, t.TextMuted
				if i == shown {
					bg, fg = t.Accent, t.AccentText
				}
				badge.Background(badge.AnimateColor("bg", bg, ui.Snappy)).TextColor(badge.AnimateColor("fg", fg, ui.Snappy)).
					Children(func() { motion.RollingDigits(c, int64(counts[i])).FontSize(t.Rem(0.75)) })
			})
		}
	})
	if kit.Changed(root, sound.Select) {
		kit.Log("TabBarBadges", "picked", "index", *active)
	}
	return root
}
