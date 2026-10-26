package geo

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Country is a country and its regions.
type Country struct {
	Name    string
	Regions []string
}

// RegionSelectOptions word the two pickers.
type RegionSelectOptions struct{ Country, Region string }

// RegionSelect picks a country, then a region; cue Select.
func RegionSelect(c *ui.Context, countries []Country, country, region *string, o RegionSelectOptions) ui.Element {
	def(&o.Country, "Country")
	def(&o.Region, "Region")
	names := make([]string, len(countries))
	for i, k := range countries {
		if len(k.Regions) == 0 {
			panic(fmt.Sprintf("geo: RegionSelect country %q needs regions", k.Name))
		}
		unique("RegionSelect", k.Regions)
		names[i] = k.Name
	}
	unique("RegionSelect", names)
	if country == nil || region == nil || !slices.Contains(names, *country) {
		panic(fmt.Sprintf("geo: RegionSelect %q needs a listed country and a region", o.Country))
	}
	regions := countries[slices.Index(names, *country)].Regions
	if *region != "" && !slices.Contains(regions, *region) {
		panic(fmt.Sprintf("geo: RegionSelect region %q is not in %q", *region, *country))
	}
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	open := ui.Local(root, "open", func() bool { return false })
	return root.Children(func() {
		kit.Muted(c, o.Country).FontSize(t.Rem(0.72)).Bold().FillWidth().MinWidth(0).SingleLine()
		c.Key("country")
		ui.Column(c).FillWidth().Children(func() {
			trigger := kit.Press(ui.ButtonBase(c).Label(*country).Role(ui.RoleMenuButton).Expanded(*open).FillWidth().
				Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Border(1, t.Border))
			trigger.Background(trigger.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, max(unit(trigger.Hovered()), unit(*open))), ui.Snappy))
			trigger.Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					glyph(c, "globe", 16).TextColor(t.Accent)
					ui.Text(c, *country).FontSize(t.Rem(0.86)).Bold().Grow(1).MinWidth(0).SingleLine().ContentTransition(ui.ContentTransition{})
					glyph(c, "chevron-down", 14).TextColor(t.TextMuted).Rotate(trigger.AnimateTo("turn", 180*unit(*open), ui.Snappy))
				})
			})
			if trigger.Clicked() {
				*open = !*open
				if *open {
					root.SetChanged()
					sound.Play(sound.Open)
					kit.Log("RegionSelect", "opened")
				}
			}
			ui.PopoverBase(c, trigger, open, func(panel ui.Element) {
				panel.Margin(t.Space(1), 0, 0, 0).MinWidth(t.Space(56)).Padding(t.Space(1)).Radius(t.Radius).
					Background(t.Background).Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).TransformOrigin(0, 0)
				kit.Arrive(panel, ui.Pose{Scale: 0.95, Fade: 1}, 0)
				ui.Scroll(c).FillWidth().MaxHeight(t.Space(64)).Gap(t.Space(0.5)).Role(ui.RoleMenu).Children(func() {
					for i, k := range countries {
						c.Key(strconv.Quote(k.Name))
						on := k.Name == *country
						item := kit.Press(ui.ButtonBase(c).Label(k.Name).Role(ui.RoleMenuItemRadio).Checked(on).FillWidth().Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius))
						item.Background(item.AnimateColor("bg", t.SurfaceHover.Alpha(unit(item.Hovered())), ui.Snappy))
						item.Children(func() {
							ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
								ui.Text(c, k.Name).FontSize(t.Rem(0.84)).Grow(1).MinWidth(0).SingleLine()
								if on {
									glyph(c, "check", 14).TextColor(t.Accent)
								}
							})
						})
						if item.Clicked() {
							*open = false
							if !on {
								*country, *region = k.Name, ""
								root.SetChanged()
								sound.Play(sound.Select)
								kit.Log("RegionSelect", "country", "index", i)
							}
						}
					}
				})
			})
		})
		kit.Muted(c, o.Region).FontSize(t.Rem(0.72)).Bold().FillWidth().MinWidth(0).SingleLine()
		c.Key(strconv.Quote(*country))
		chips := ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).Role(ui.RoleRadioGroup).Label(o.Region)
		chips.Children(func() {
			for i, name := range regions {
				c.Key(strconv.Quote(name))
				on := name == *region
				b := kit.Press(ui.ButtonBase(c).Label(name).Role(ui.RoleRadio).Checked(on).Padding(t.Space(1), t.Space(3)).Radius(999).
					Border(1, t.Border.Mix(t.Accent, unit(on))))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent, unit(on)), ui.Snappy))
				kit.Stagger(c, chips, b, i).Children(func() {
					ui.Text(c, name).FontSize(t.Rem(0.78)).Bold().TextColor(t.Text.Mix(t.AccentText, unit(on))).Shrink(0).MaxWidth(t.Space(48)).SingleLine()
				})
				if b.Clicked() && !on {
					*region = name
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("RegionSelect", "region", "index", i)
				}
			}
		})
	})
}
