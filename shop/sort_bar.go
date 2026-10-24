package shop

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SortBarOptions word the bar; Count takes the product count.
type SortBarOptions struct{ Count, Sort string }

// SortBar shows the count and a sort menu; cue Select.
func SortBar(c *ui.Context, count int, sort *string, choices []string, o SortBarOptions) ui.Element {
	if count < 0 || sort == nil || len(choices) == 0 {
		panic("shop: SortBar needs a count and choices")
	}
	for i, ch := range choices {
		if ch == "" || slices.Index(choices, ch) != i {
			panic(fmt.Sprintf("shop: SortBar choice %q must be unique text", ch))
		}
	}
	if !slices.Contains(choices, *sort) {
		panic(fmt.Sprintf("shop: SortBar sort %q is not a choice", *sort))
	}
	def(&o.Count, "%d products")
	def(&o.Sort, "Sort by")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(0, 0, t.Space(2), 0).
		BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	open := ui.Local(root, "open", func() bool { return false })
	root.Children(func() {
		kit.Numeric(c, fmt.Sprintf(o.Count, count)).Bold().Grow(1).MinWidth(0).SingleLine().
			ContentTransition(ui.ContentTransition{Numeric: true})
		trigger := kit.Press(ui.ButtonBase(c).Label(o.Sort).Row().Gap(t.Space(1.5)).AlignItems(ui.Center).
			Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius).Border(1, t.Border).Expanded(*open))
		trigger.Background(trigger.AnimateColor("bg", t.SurfaceHover.Alpha(unit(*open || trigger.Hovered())), ui.Snappy))
		trigger.Children(func() {
			ui.Icon(c, icon.Must("arrow-down-wide-narrow")).Size(16, 16).TextColor(t.TextMuted)
			kit.Muted(c, o.Sort).FontSize(t.Rem(0.85)).Shrink(0)
			ui.Text(c, *sort).FontSize(t.Rem(0.85)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{})
			chev := ui.Icon(c, icon.Must("chevron-down")).Size(16, 16).TextColor(t.TextMuted)
			chev.Rotate(180 * chev.AnimateTo("turn", unit(*open), ui.Snappy))
		})
		if trigger.Clicked() {
			*open = !*open
		}
		ui.PopoverBase(c, trigger, open, func(panel ui.Element) {
			panel.Margin(t.Space(1), 0, 0, 0).MinWidth(t.Space(44)).Padding(t.Space(1)).Gap(t.Space(0.5)).Radius(t.Radius).
				Background(t.Background).Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).TransformOrigin(0.5, 0)
			kit.Arrive(panel, ui.Pose{Scale: 0.95, Fade: 1}, 0)
			ui.Scroll(c).FillWidth().MaxHeight(t.Space(72)).Gap(t.Space(0.5)).Children(func() {
				for i, ch := range choices {
					c.Key(ch)
					on := ch == *sort
					item := kit.Press(ui.ButtonBase(c).Label(ch).FillWidth().Row().Gap(t.Space(2)).AlignItems(ui.Center).
						Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).Checked(on))
					item.Background(item.AnimateColor("bg", t.SurfaceHover.Alpha(unit(item.Hovered())), ui.Snappy))
					item.Children(func() {
						ui.Text(c, ch).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
						ui.Icon(c, icon.Must("check")).Size(14, 14).TextColor(t.Accent).Opacity(unit(on))
					})
					if item.Clicked() {
						*open = false
						if !on {
							*sort = ch
							root.SetChanged()
							sound.Play(sound.Select)
							kit.Log("SortBar", "selected", "index", i)
						}
					}
				}
			})
		})
	})
	return root
}
