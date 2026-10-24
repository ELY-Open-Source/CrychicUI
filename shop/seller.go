package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SellerOptions word the card; Meta sits under the rating.
type SellerOptions struct {
	Follow, Following, Meta string
	Verified                bool
	Image                   *ui.Bitmap
}

// Seller shows avatar, rating and a follow pill; cue Toggle.
func Seller(c *ui.Context, name string, rating float64, followed *bool, o SellerOptions) ui.Element {
	if name == "" || !(rating >= 0 && rating <= 5) || followed == nil {
		panic(fmt.Sprintf("shop: Seller %q needs a rating in [0,5] and bound state", name))
	}
	def(&o.Follow, "Follow")
	def(&o.Following, "Following")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	var b ui.Element
	root.Row().AlignItems(ui.Center).Children(func() {
		display.Avatar(c, name, o.Image).Size(t.Space(12), t.Space(12)).Shrink(0)
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.75)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, name).Bold().MinWidth(0).SingleLine()
				if o.Verified {
					ui.Icon(c, icon.Must("badge-check")).Size(16, 16).TextColor(t.Accent)
				}
			})
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				stars(c, rating, t.Space(3))
				kit.Numeric(c, fmt.Sprintf("%.1f", rating)).FontSize(t.Rem(0.75)).Bold().Shrink(0)
			})
			if o.Meta != "" {
				kit.Muted(c, o.Meta).FontSize(t.Rem(0.75)).FillWidth().Wrap()
			}
		})
		name, glyph := o.Follow, "user-plus"
		if *followed {
			name, glyph = o.Following, "user-check"
		}
		b = kit.Press(ui.ButtonBase(c).Label(name).Radius(999).Shrink(0).Checked(*followed))
		b.Background(b.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(b.Hovered())).Mix(t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), unit(*followed)), ui.Snappy))
		b.Border(1, b.AnimateColor("line", t.Accent.Mix(t.Border, unit(*followed)), ui.Snappy))
		b.TextColor(b.AnimateColor("ink", t.AccentText.Mix(t.Text, unit(*followed)), ui.Snappy))
		b.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(3)).Children(func() {
				ui.Icon(c, icon.Must(glyph)).Size(14, 14).ContentTransition(ui.ContentTransition{})
				ui.Text(c, name).FontSize(t.Rem(0.85)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{})
			})
		})
	})
	if b.Clicked() {
		*followed = !*followed
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("Seller", "followed", "value", *followed)
	}
	return root
}
