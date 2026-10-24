package shop

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// OrderCardOptions word the card; Images show up to four items.
type OrderCardOptions struct {
	Open                               string
	Placed, Packed, Shipped, Delivered string
	// Date lays out the order date; More counts hidden images.
	Date, More string
	Images     []*ui.Bitmap
}

// OrderCard shows number, date, total and status; cue Tap.
func OrderCard(c *ui.Context, number string, when time.Time, total float64, currency string, status OrderStage, opened *bool, o OrderCardOptions) ui.Element {
	if number == "" || total < 0 || !finite(total) || currency == "" || status > Delivered || opened == nil {
		panic(fmt.Sprintf("shop: OrderCard %q needs a total, currency, known status and bound state", number))
	}
	for _, img := range o.Images {
		if img == nil {
			panic(fmt.Sprintf("shop: OrderCard %q images must not be nil", number))
		}
	}
	def(&o.Open, "View order")
	def(&o.Date, "Jan 2, 2006")
	def(&o.More, "+%d")
	names := stageNames(&OrderStatusOptions{Placed: o.Placed, Packed: o.Packed, Shipped: o.Shipped, Delivered: o.Delivered})
	t := c.Theme()
	ink := stageInk(t, status)
	root := kit.Lift(kit.Press(ui.ButtonBase(c).Label(o.Open).Column().FillWidth().Gap(t.Space(3)).AlignItems(ui.Stretch).
		Padding(t.Space(3)).Radius(t.Radius).Shadow(0, 1, 3, 0, kit.Shade(0.06))))
	root.Background(root.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, 0.6*unit(root.Hovered())), ui.Snappy))
	root.Border(1, root.AnimateColor("line", t.Border.Mix(t.Accent, 0.5*unit(root.Hovered())), ui.Snappy))
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(t.Space(10), t.Space(10)).Shrink(0).Center().Radius(t.Radius).Background(ink.Alpha(0.12)).TextColor(ink).
				Children(func() { ui.Icon(c, icon.Must(stageGlyphs[status])).Size(18, 18) })
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, number).Bold().FillWidth().SingleLine()
				kit.Muted(c, when.Format(o.Date)).FontSize(t.Rem(0.75))
			})
			ui.Column(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.End).Children(func() {
				money(c, total, currency).Bold()
				kit.Chip(c, names[status], ink)
			})
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Row(c).Grow(1).MinWidth(0).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				for i, img := range o.Images[:min(4, len(o.Images))] {
					c.Key(i)
					thumb(c, img, t.Space(10))
				}
				if n := len(o.Images) - 4; n > 0 {
					ui.Box(c).Size(t.Space(10), t.Space(10)).Shrink(0).Center().Radius(t.Radius).Background(t.Surface).Children(func() {
						kit.Numeric(c, fmt.Sprintf(o.More, n)).FontSize(t.Rem(0.75)).Bold().TextColor(t.TextMuted)
					})
				}
				if len(o.Images) == 0 {
					track(c, status)
				}
			})
			chev := ui.Icon(c, icon.Must("chevron-right")).Size(16, 16).TextColor(t.TextMuted)
			chev.Translate(t.Space(1)*chev.AnimateTo("nudge", unit(root.Hovered()), ui.Snappy), 0)
		})
		if len(o.Images) > 0 {
			track(c, status)
		}
	})
	if root.Clicked() && !*opened {
		*opened = true
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("OrderCard", "opened")
	}
	return root
}

// track is four segments, filled through the order's stage.
func track(c *ui.Context, status OrderStage) {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(1)).Children(func() {
		for i := range 4 {
			c.Key(i)
			seg := ui.Box(c).Grow(1).Basis(0).Height(t.Space(1)).Radius(999)
			seg.Background(seg.AnimateColor("fill", t.Border.Mix(stageInk(t, status), unit(OrderStage(i) <= status)), ui.Smooth))
		}
	})
}
