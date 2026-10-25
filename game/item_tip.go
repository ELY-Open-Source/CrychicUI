package game

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Stat is one stat; Delta compares with worn gear.
type Stat struct {
	Name         string
	Value, Delta int
}

// ItemDetail is an item with stats; Text, Level optional.
type ItemDetail struct {
	Name, Kind, Icon, Text string
	Rarity                 Rarity
	Level                  int
	Stats                  []Stat
}

// ItemTipOptions word the card; Rarities names Common to Legendary.
type ItemTipOptions struct {
	Level    string
	Rarities [5]string
}

// ItemTip is an item card with its stats.
func ItemTip(c *ui.Context, it ItemDetail, o ItemTipOptions) ui.Element {
	if it.Name == "" || it.Icon == "" || it.Kind == "" || it.Level < 0 {
		panic(fmt.Sprintf("game: ItemTip %q needs a name, icon, kind and a level of at least 0", it.Name))
	}
	validRarity("ItemTip", it.Name, it.Rarity)
	names := make([]string, len(it.Stats))
	for i, s := range it.Stats {
		names[i] = s.Name
	}
	if len(names) > 0 {
		unique("ItemTip", names)
	}
	def(&o.Level, "Requires level %d")
	for i, name := range [...]string{"Common", "Uncommon", "Rare", "Epic", "Legendary"} {
		def(&o.Rarities[i], name)
	}
	t := c.Theme()
	ink := rarityInk(t, it.Rarity)
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Border(1, ink.Alpha(0.55)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			kit.Arrive(tile(c, it.Icon, ink, t.Space(12)), ui.Pose{Scale: 0.7, Rotate: -10, Fade: 1}, 0)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, it.Name).FontSize(t.Rem(0.95)).Bold().TextColor(ink).FillWidth().Wrap()
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					ui.Text(c, o.Rarities[it.Rarity]).FontSize(t.Rem(0.72)).Bold().TextColor(ink).Shrink(0)
					ui.Box(c).Size(t.Space(0.75), t.Space(0.75)).Shrink(0).Radius(999).Background(t.TextMuted.Alpha(0.6))
					kit.Muted(c, it.Kind).FontSize(t.Rem(0.72)).MinWidth(0).SingleLine()
				})
			})
		})
		if len(it.Stats) > 0 {
			list := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2.5), 0).BorderWidth(1, 0, 1, 0).BorderColor(t.Border)
			list.Children(func() {
				for i, s := range it.Stats {
					c.Key(s.Name)
					row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
					kit.Stagger(c, list, row, i).Children(func() {
						kit.Muted(c, s.Name).FontSize(t.Rem(0.78)).Grow(1).MinWidth(0).SingleLine()
						kit.Numeric(c, fmt.Sprint(s.Value)).FontSize(t.Rem(0.82)).Bold().Shrink(0)
						if s.Delta != 0 {
							tone := t.Success
							if s.Delta < 0 {
								tone = t.Danger
							}
							kit.Chip(c, fmt.Sprintf("%+d", s.Delta), tone)
						}
					})
				}
			})
		}
		if it.Text != "" {
			kit.Muted(c, it.Text).FontSize(t.Rem(0.78)).FillWidth().Wrap()
		}
		if it.Level > 0 {
			kit.Muted(c, fmt.Sprintf(o.Level, it.Level)).FontSize(t.Rem(0.7)).Bold()
		}
	})
	return root
}
