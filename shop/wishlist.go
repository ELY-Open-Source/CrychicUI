package shop

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// WishlistOptions are the wishlist labels.
type WishlistOptions struct{ Save, Saved string }

// Wishlist toggles a saved product; cue Toggle.
func Wishlist(c *ui.Context, p Product, saved *bool, o WishlistOptions) ui.Element {
	validProduct("Wishlist", p)
	if saved == nil {
		panic("shop: Wishlist needs a bound state")
	}
	def(&o.Save, "Save")
	def(&o.Saved, "Saved")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border)
	root.Children(func() {
		ui.Text(c, p.Name).SingleLine().Grow(1)
		name, ink := o.Save, t.TextMuted
		if *saved {
			name, ink = o.Saved, t.Danger
		}
		b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(8), t.Space(8)).Center().Radius(999).TextColor(ink).
			Checked(*saved))
		b.Children(func() { ui.Icon(c, icon.Must("heart")).Size(16, 16).ContentTransition(ui.ContentTransition{}) })
		if b.Clicked() {
			*saved = !*saved
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("Wishlist", "saved", "value", *saved)
		}
	})
	return root
}
