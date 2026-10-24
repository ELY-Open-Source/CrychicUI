package finance

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SplitBillOptions word the card.
type SplitBillOptions struct{ Title, People, Each, Decrease, Increase string }

// SplitBill splits a total among people; cue Tick.
func SplitBill(c *ui.Context, total float64, people *int, currency string, o SplitBillOptions) ui.Element {
	if total < 0 || !finite(total) || people == nil || *people < 1 || *people > 100 || currency == "" {
		panic("finance: SplitBill has invalid data")
	}
	def(&o.Title, "Split the bill")
	def(&o.People, "People")
	def(&o.Each, "Per person")
	def(&o.Decrease, "Decrease")
	def(&o.Increase, "Increase")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "receipt", o.Title, func() { money(c, total, currency).Bold().Shrink(0) })
		ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1.5)).Children(func() {
					kit.Muted(c, o.People).FontSize(t.Rem(0.8))
					ui.Row(c).Children(func() {
						shown := min(*people, 6)
						for i := range shown {
							c.Key(i)
							face := ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999).
								Background(palette(t, i).Mix(t.Background, 0.8)).TextColor(palette(t, i)).Border(2, t.Background)
							if i > 0 {
								face.MarginStart(-t.Space(2))
							}
							face.Transition(ui.ElementTransition{Animation: &ui.Bouncy, Enter: &ui.Motion{Scale: 0.4},
								Exit: &ui.Motion{Scale: 0.4}})
							face.Children(func() { ui.Icon(c, icon.Must("user")).Size(14, 14) })
						}
						if *people > shown {
							c.Key("more")
							ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999).MarginStart(-t.Space(2)).
								Background(t.SurfaceHover).Border(2, t.Background).Children(func() {
								kit.Numeric(c, "+"+strconv.Itoa(*people-shown)).FontSize(t.Rem(0.7)).Bold().Shrink(0).
									ContentTransition(ui.ContentTransition{Numeric: true})
							})
						}
					})
				})
				ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.5)).Radius(999).Border(1, t.Border).
					Children(func() {
						for _, step := range []struct {
							name, glyph string
							delta       int
						}{{o.Decrease, "minus", -1}, {o.Increase, "plus", 1}} {
							can := *people+step.delta >= 1 && *people+step.delta <= 100
							b := kit.Press(ui.ButtonBase(c).Label(step.name).Tooltip(step.name).Size(t.Space(8), t.Space(8)).
								Center().Radius(999)).Disabled(!can)
							b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy))
							b.Children(func() { ui.Icon(c, icon.Must(step.glyph)).Size(14, 14) })
							if step.delta < 0 {
								kit.Numeric(c, fmt.Sprint(*people)).Bold().MinWidth(t.Space(6)).TextAlign(ui.Center).
									ContentTransition(ui.ContentTransition{Numeric: true})
							}
							if b.Clicked() && can {
								*people += step.delta
								root.SetChanged()
								sound.Play(sound.Tick)
								kit.Log("SplitBill", "people", "value", *people)
							}
						}
					})
			})
			ui.Row(c).FillWidth().AlignItems(ui.Center).Padding(t.Space(3)).Radius(t.Radius).Background(t.Accent.Alpha(0.08)).
				Children(func() {
					ui.Text(c, o.Each).Bold().Grow(1).MinWidth(0).SingleLine()
					money(c, total/float64(*people), currency).Bold().FontSize(t.Rem(1.4)).TextColor(t.Accent).Shrink(0).
						ContentTransition(ui.ContentTransition{Numeric: true})
				})
		})
	})
	return root
}
