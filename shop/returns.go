package shop

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ReturnsOptions word the reason heading.
type ReturnsOptions struct{ Reason string }

// Returns picks an item and a reason chip; cue Select.
func Returns(c *ui.Context, items []Product, currency string, reasons []string, selected *int, reason *string, o ReturnsOptions) ui.Element {
	if len(items) == 0 || currency == "" || len(reasons) == 0 || selected == nil || reason == nil || *selected < -1 || *selected >= len(items) {
		panic("shop: Returns needs items, a currency, reasons, an item or -1 and a reason")
	}
	for i, p := range items {
		validProduct("Returns", p)
		if slices.IndexFunc(items, func(q Product) bool { return q.Name == p.Name }) != i {
			panic(fmt.Sprintf("shop: Returns item %q must be unique", p.Name))
		}
	}
	for i, r := range reasons {
		if r == "" || slices.Index(reasons, r) != i {
			panic(fmt.Sprintf("shop: Returns reason %q must be unique text", r))
		}
	}
	if *reason != "" && !slices.Contains(reasons, *reason) {
		panic(fmt.Sprintf("shop: Returns reason %q is not offered", *reason))
	}
	def(&o.Reason, "Reason for return")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).Role(ui.RoleGroup)
	root.Children(func() {
		list := ui.Column(c).FillWidth().Gap(t.Space(2)).FocusGroup(ui.Vertical)
		list.Children(func() {
			for i, p := range items {
				c.Key(p.Name)
				on := i == *selected
				b := kit.Press(ui.ButtonBase(c).Label(p.Name).Row().FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).
					Padding(t.Space(2)).Radius(t.Radius).Checked(on))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent.Alpha(0.08), unit(on)), ui.Snappy))
				b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
				kit.Stagger(c, list, b, i).Children(func() {
					radio(c, on)
					thumb(c, p.Image, t.Space(11))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, p.Name).Bold().FontSize(t.Rem(0.85)).FillWidth().SingleLine()
						money(c, p.Price, currency).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
					})
				})
				if b.Clicked() && !on {
					*selected = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Returns", "selected", "index", i)
				}
			}
		})
		ui.Text(c, o.Reason).Bold().FontSize(t.Rem(0.85)).Shrink(0)
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).Children(func() {
			for i, r := range reasons {
				c.Key(r)
				on := *reason == r
				b := kit.Press(ui.ButtonBase(c).Label(r).Radius(999).Checked(on))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent.Alpha(0.12), unit(on)), ui.Snappy))
				b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
				b.TextColor(b.AnimateColor("ink", t.Text.Mix(t.Accent, unit(on)), ui.Snappy))
				b.Children(func() {
					ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2.5)).Children(func() {
						if on {
							kit.Arrive(ui.Icon(c, icon.Must("check")).Size(12, 12), ui.Pose{Scale: 0.4, Fade: 1}, 0)
						}
						ui.Text(c, r).FontSize(t.Rem(0.8)).Shrink(0).MaxWidth(t.Space(60)).Wrap()
					})
				})
				if b.Clicked() && !on {
					*reason = r
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Returns", "reason", "index", i)
				}
			}
		})
	})
	return root
}
