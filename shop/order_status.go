package shop

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// OrderStage is one order tracking stage.
type OrderStage uint8

// Order stages.
const (
	Placed OrderStage = iota
	Packed
	Shipped
	Delivered
)

// OrderStatusOptions name the stages; Dates sit under them.
type OrderStatusOptions struct {
	Placed, Packed, Shipped, Delivered string
	Dates                              [4]string
}

var stageGlyphs = [...]string{"receipt", "package", "truck", "house"}

func stageNames(o *OrderStatusOptions) [4]string {
	def(&o.Placed, "Placed")
	def(&o.Packed, "Packed")
	def(&o.Shipped, "Shipped")
	def(&o.Delivered, "Delivered")
	return [4]string{o.Placed, o.Packed, o.Shipped, o.Delivered}
}

func stageInk(t *ui.Theme, s OrderStage) ui.Color {
	return [...]ui.Color{t.TextMuted, t.Warning, t.Accent, t.Success}[s]
}

// OrderStatus tracks an order through icon steps and filling lines.
func OrderStatus(c *ui.Context, current OrderStage, o OrderStatusOptions) ui.Element {
	if current > Delivered {
		panic(fmt.Sprintf("shop: OrderStatus needs a stage up to Delivered, got %d", current))
	}
	names := stageNames(&o)
	t := c.Theme()
	root := ui.Row(c).FillWidth().AlignItems(ui.Start).Role(ui.RoleGroup)
	root.Children(func() {
		for i, name := range names {
			s := OrderStage(i)
			c.Key(s)
			if i > 0 {
				line := ui.Box(c).Grow(1).Height(t.Space(0.5)).Margin(t.Space(4.25), t.Space(1), 0, t.Space(1)).Radius(999).Background(t.Border)
				fill := line.AnimateTo("fill", unit(s <= current), ui.Smooth)
				line.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * fill, H: r.H}, t.Success, 999)
				}))
			}
			done, on := s < current, s == current
			step := ui.Column(c).Shrink(0).Gap(t.Space(1.5)).AlignItems(ui.Center)
			kit.Arrive(step, ui.Pose{Y: t.Space(1), Fade: 1}, time.Duration(i)*60*time.Millisecond).Children(func() {
				dot := ui.Box(c).Size(t.Space(9), t.Space(9)).Center().Radius(999)
				dot.Background(dot.AnimateColor("bg", t.Surface.Mix(t.Success, unit(done)).Mix(t.Accent, unit(on)), ui.Snappy))
				dot.TextColor(dot.AnimateColor("ink", t.TextMuted.Mix(t.AccentText, unit(done || on)), ui.Snappy))
				dot.Border(t.Space(1)*max(0, dot.AnimateTo("ring", unit(on), ui.Bouncy)), t.Accent.Alpha(0.2))
				glyph := stageGlyphs[i]
				if done {
					glyph = "check"
				}
				dot.Children(func() { ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{}) })
				label := ui.Text(c, name).FontSize(t.Rem(0.75)).Shrink(0)
				switch {
				case on:
					label.Bold()
				case !done:
					label.TextColor(t.TextMuted)
				}
				if o.Dates[i] != "" {
					kit.Muted(c, o.Dates[i]).FontSize(t.Rem(0.7)).Shrink(0)
				}
			})
		}
	})
	return root
}
