package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SortBarOptions are the sort bar labels.
type SortBarOptions struct{ Count string }

// SortBar shows product count and sort choices; cue Select.
func SortBar(c *ui.Context, count int, sort *string, choices []string, o SortBarOptions) ui.Element {
	if count < 0 || sort == nil || len(choices) == 0 {
		panic("shop: SortBar needs count and choices")
	}
	def(&o.Count, "%d products")
	valid := false
	for _, choice := range choices {
		if choice == "" {
			panic("shop: SortBar choice must not be empty")
		}
		valid = valid || *sort == choice
	}
	if !valid {
		panic("shop: SortBar needs a current choice")
	}
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		kit.Numeric(c, fmt.Sprintf(o.Count, count)).Bold().Grow(1)
		for _, choice := range choices {
			b := kit.Press(ui.ButtonBase(c).Label(choice).Padding(t.Space(1), t.Space(2)).Radius(999).
				Border(1, t.Border).Checked(*sort == choice))
			b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*sort == choice)), ui.Snappy))
			b.Children(func() { ui.Text(c, choice).FontSize(t.Rem(0.8)).SingleLine() })
			if b.Clicked() && *sort != choice {
				*sort = choice
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("SortBar", "selected", "sort", *sort)
			}
		}
	})
	return root
}
