package marketing

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Partner is one tile; Text shows in the top tier.
type Partner struct{ Name, Icon, Text string }

// PartnerTier is a named group of partners, top tier first.
type PartnerTier struct {
	Name     string
	Partners []Partner
}

// PartnersOptions word the section.
type PartnersOptions struct{ Title, Text string }

// Partners shows partner tiers, the top one larger.
func Partners(c *ui.Context, tiers []PartnerTier, o PartnersOptions) ui.Element {
	if len(tiers) == 0 {
		panic(fmt.Sprintf("marketing: Partners %q needs tiers", o.Title))
	}
	names, all := make([]string, len(tiers)), []string{}
	for i, tier := range tiers {
		if len(tier.Partners) == 0 {
			panic(fmt.Sprintf("marketing: Partners tier %q needs partners", tier.Name))
		}
		for _, p := range tier.Partners {
			if strings.TrimSpace(p.Icon) == "" {
				panic(fmt.Sprintf("marketing: Partners partner %q needs an icon", p.Name))
			}
			all = append(all, p.Name)
		}
		names[i] = tier.Name
	}
	unique("Partners", names)
	unique("Partners", all)
	t := c.Theme()
	inks := []ui.Color{t.Warning, t.TextMuted, t.Accent}
	root := ui.Column(c).FillWidth().Gap(t.Space(5)).Role(ui.RoleGroup)
	root.Children(func() {
		heading(c, o.Title, o.Text)
		for r, tier := range tiers {
			top, ink := r == 0, inks[min(r, len(inks)-1)]
			c.Key(strconv.Quote(tier.Name))
			ui.Column(c).FillWidth().Gap(t.Space(2.5)).Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					glyph(c, "medal", 16).TextColor(ink)
					ui.Text(c, tier.Name).FontSize(t.Rem(0.88)).Bold().Grow(1).MinWidth(0).SingleLine()
					kit.Chip(c, strconv.Itoa(len(tier.Partners)), ink)
				})
				grid := ui.Row(c).FillWidth().Wrap().Gap(t.Space(2.5))
				grid.Children(func() {
					for i, p := range tier.Partners {
						basis, side := t.Space(30), t.Space(8)
						if top {
							basis, side = t.Space(48), t.Space(11)
						}
						c.Key(strconv.Quote(p.Name))
						box := ui.Row(c).Grow(1).Basis(basis).MinWidth(0).Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(3)).
							Radius(t.Radius*1.5).Border(1, t.Border).Background(t.Background)
						if top {
							box.Border(1, ink.Alpha(0.4)).Background(t.Background.Mix(ink, 0.05))
						}
						kit.Stagger(c, grid, box, i).Children(func() {
							tile(c, p.Icon, ink, side)
							ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
								ui.Text(c, p.Name).FontSize(t.Rem(0.84)).Bold().FillWidth().Wrap()
								if top && p.Text != "" {
									kit.Muted(c, p.Text).FontSize(t.Rem(0.74)).FillWidth().Wrap()
								}
							})
						})
					}
				})
			})
		}
	})
	return root
}
