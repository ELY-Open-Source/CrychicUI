package shop

import (
	"fmt"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SellerOptions are the seller labels.
type SellerOptions struct{ Follow, Following string }

// Seller shows seller rating and follow; cue Toggle.
func Seller(c *ui.Context, name string, rating float64, followed *bool, o SellerOptions) ui.Element {
	if name == "" || rating < 0 || rating > 5 || !finite(rating) || followed == nil {
		panic("shop: Seller has invalid data")
	}
	def(&o.Follow, "Follow")
	def(&o.Following, "Following")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border)
	root.Children(func() {
		ui.Box(c).Size(t.Space(9), t.Space(9)).Radius(999).Background(t.Accent.Alpha(0.15)).Center().Children(func() {
			first, _ := utf8.DecodeRuneInString(name)
			ui.Text(c, string(first)).Bold().TextColor(t.Accent)
		})
		ui.Column(c).Grow(1).Children(func() {
			ui.Text(c, name).Bold().SingleLine()
			kit.Numeric(c, fmt.Sprintf("%.1f", rating)).FontSize(t.Rem(0.8)).TextColor(t.Warning)
		})
		name := o.Follow
		if *followed {
			name = o.Following
		}
		b := kit.Press(ui.ButtonBase(c).Label(name).Padding(t.Space(1), t.Space(2)).Radius(999).Checked(*followed))
		b.Children(func() { ui.Text(c, name).SingleLine() })
		if b.Clicked() {
			*followed = !*followed
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("Seller", "followed", "value", *followed)
		}
	})
	return root
}
