package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// WishlistOptions name the heart's two states.
type WishlistOptions struct{ Save, Saved string }

// Wishlist is a product row with a heart; cue Toggle.
func Wishlist(c *ui.Context, p Product, currency string, saved *bool, o WishlistOptions) ui.Element {
	validProduct("Wishlist", p)
	if saved == nil || currency == "" {
		panic(fmt.Sprintf("shop: Wishlist %q needs a currency and bound state", p.Name))
	}
	def(&o.Save, "Save")
	def(&o.Saved, "Saved")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(2)).Role(ui.RoleGroup).Label(p.Name)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	var b ui.Element
	root.Row().AlignItems(ui.Center).Children(func() {
		if p.Image != nil {
			ui.Image(c, p.Image).Size(t.Space(12), t.Space(12)).Shrink(0).Fit(ui.Cover).Radius(t.Radius)
		}
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
			ui.Text(c, p.Name).Bold().FillWidth().SingleLine()
			money(c, p.Price, currency).FontSize(t.Rem(0.8)).TextColor(t.TextMuted)
		})
		name := o.Save
		if *saved {
			name = o.Saved
		}
		// Keyed: an optional image before must not move its focus.
		c.Key("heart")
		b = kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(10), t.Space(10)).Shrink(0).Center().Radius(999).Checked(*saved))
		b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(t.Danger.Alpha(0.12), unit(*saved)), ui.Snappy))
		b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Danger, unit(*saved)), ui.Snappy))
		pop := b.AnimateTo("pop", unit(*saved), ui.Bouncy)
		b.Children(func() {
			heart := ui.Icon(c, icon.Must("heart")).Size(18, 18)
			heart.Scale(1+0.3*pop*(1-pop)*4, 1+0.3*pop*(1-pop)*4)
		})
	})
	if b.Clicked() {
		*saved = !*saved
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("Wishlist", "saved", "value", *saved)
	}
	return root
}
