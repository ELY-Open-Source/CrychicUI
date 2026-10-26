package marketing

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SiteNavOptions word the bar; Menu labels the narrow menu button.
type SiteNavOptions struct{ Menu string }

// SiteNav folds its links into a menu; cue Open.
func SiteNav(c *ui.Context, brand string, links []string, active *string, act string, o SiteNavOptions) (root, button ui.Element) {
	if strings.TrimSpace(brand) == "" || strings.TrimSpace(act) == "" || len(links) == 0 || active == nil || *active != "" && !slices.Contains(links, *active) {
		panic(fmt.Sprintf("marketing: SiteNav %q needs links, an action and an active link among them", brand))
	}
	unique("SiteNav", links)
	def(&o.Menu, "Menu")
	t := c.Theme()
	root = ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2.5), t.Space(4)).Radius(t.Radius*1.5).
		Border(1, t.Border).Background(t.Background).Role(ui.RoleToolbar)
	kit.Arrive(root, ui.Pose{Y: -t.Space(2), Fade: 1}, 0)
	w := root.LayoutBounds().W
	kit.Settle(c, root, []float32{w}, func() []float32 { return []float32{root.LayoutBounds().W} })
	bold := func(s string, size float32) float32 {
		tw, _ := c.MeasureText(0, ui.Span{Text: s, Size: size, Weight: 700})
		return tw
	}
	need := float32(len(links)-1) * t.Space(1)
	for _, l := range links {
		need += bold(l, t.Rem(0.82)) + 2*t.Space(2)
	}
	brandW := t.Space(10) + min(bold(brand, t.Rem(0.95)), t.Space(40))
	actW := min(bold(act, t.Rem(0.84)), t.Space(60)) + 2*t.Space(4)
	// Room left for links beside brand, action, padding and gaps.
	room := w - 2*t.Space(4) - 2 - brandW - actW - 2*t.Space(3)
	compact := w > 0 && room < need
	// Compact: brand and action share what Menu leaves.
	brandCap, actCap := t.Space(40), t.Space(60)
	if compact {
		free := w - 2*t.Space(4) - 2 - 3*t.Space(3) - t.Space(9) - t.Space(10) - 2*t.Space(4)
		bw, aw := brandW-t.Space(10), actW-2*t.Space(4)
		if bw+aw > free {
			actCap = max(t.Space(6), min(aw, max(free-bw, free/2)))
			brandCap = max(t.Space(6), free-actCap)
		}
	}
	open := ui.Local(root, "open", func() bool { return false })
	if !compact {
		*open = false
	}
	pick := func(i int, l string) {
		*open = false
		if l != *active {
			*active = l
			root.SetChanged()
			sound.Play(sound.Tap)
			kit.Log("SiteNav", "picked", "index", i)
		}
	}
	root.Children(func() {
		ui.Row(c).Shrink(0).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			tile(c, "sparkles", t.Accent, t.Space(8))
			ui.Text(c, brand).FontSize(t.Rem(0.95)).Bold().Shrink(0).MaxWidth(brandCap).SingleLine()
		})
		ui.Row(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Justify(ui.Center).Children(func() {
			if compact {
				return
			}
			for i, l := range links {
				on := l == *active
				c.Key(strconv.Quote(l))
				b := ui.ButtonBase(c).Label(l).Checked(on).Column().Gap(t.Space(0.75)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2)).
					Radius(t.Radius)
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
				b.Children(func() {
					ui.Text(c, l).FontSize(t.Rem(0.82)).Bold().TextColor(t.TextMuted.Mix(t.Text, unit(on || b.Hovered()))).Shrink(0).SingleLine()
					bar := ui.Box(c).Height(t.Space(0.5)).Radius(999).Background(t.Accent)
					bar.Width(t.Space(5) * bar.AnimateTo("lit", unit(on), ui.Snappy))
				})
				if b.Clicked() {
					pick(i, l)
				}
			}
		})
		c.Key("action")
		button = action(c, act, "", true, actCap)
		if button.Clicked() {
			sound.Play(sound.Tap)
			kit.Log("SiteNav", "action")
		}
		if !compact {
			return
		}
		c.Key("menu")
		menu := arrow(c, o.Menu, "menu").Expanded(*open)
		if menu.Clicked() {
			*open = !*open
			if *open {
				sound.Play(sound.Open)
			}
			kit.Log("SiteNav", "menu", "open", *open)
		}
		ui.PopoverBase(c, menu, open, func(panel ui.Element) {
			panel.Margin(t.Space(1), 0, 0, 0).MinWidth(t.Space(48)).Padding(t.Space(1)).Gap(t.Space(0.5)).Radius(t.Radius).
				Background(t.Background).Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).TransformOrigin(1, 0)
			kit.Arrive(panel, ui.Pose{Scale: 0.95, Fade: 1}, 0)
			ui.Scroll(c).FillWidth().MaxHeight(t.Space(72)).Gap(t.Space(0.5)).Children(func() {
				for i, l := range links {
					on := l == *active
					c.Key(strconv.Quote(l))
					item := kit.Press(ui.ButtonBase(c).Label(l).FillWidth().Row().Gap(t.Space(2)).AlignItems(ui.Center).
						Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).Checked(on))
					item.Background(item.AnimateColor("bg", t.SurfaceHover.Alpha(unit(item.Hovered())), ui.Snappy))
					item.Children(func() {
						ui.Text(c, l).FontSize(t.Rem(0.85)).Bold().TextColor(t.Text.Mix(t.Accent, unit(on))).Grow(1).MinWidth(0).SingleLine()
						glyph(c, "check", 14).TextColor(t.Accent).Opacity(unit(on))
					})
					if item.Clicked() {
						pick(i, l)
					}
				}
			})
		})
	})
	return root, button
}
