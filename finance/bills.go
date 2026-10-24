package finance

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Bill is one due bill; Icon is a Lucide name.
type Bill struct {
	Name, Due string
	Amount    float64
	Icon      string
	Overdue   bool
}

// BillsOptions word the card; PayFor names each pay button.
type BillsOptions struct{ Title, Pay, PayFor, Paid, AllPaid string }

// Bills lists due bills under their unpaid total; cue Success.
func Bills(c *ui.Context, bills []Bill, paid *[]bool, currency string, o BillsOptions) ui.Element {
	if len(bills) == 0 || paid == nil || len(*paid) != len(bills) || currency == "" {
		panic("finance: Bills needs bills, one paid flag each and a currency")
	}
	due := 0.0
	for i, b := range bills {
		if b.Name == "" || b.Due == "" || b.Amount <= 0 || !finite(b.Amount) ||
			slices.IndexFunc(bills, func(x Bill) bool { return x.Name == b.Name }) != i {
			panic(fmt.Sprintf("finance: Bills bill %q needs a unique name, a due date and a positive amount", b.Name))
		}
		if !(*paid)[i] {
			due += b.Amount
		}
	}
	if !finite(due) {
		panic("finance: Bills total overflows")
	}
	def(&o.Title, "Upcoming bills")
	def(&o.Pay, "Pay")
	def(&o.PayFor, "Pay %s")
	def(&o.Paid, "Paid")
	def(&o.AllPaid, "All bills paid")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "calendar", o.Title, func() {
			money(c, due, currency).Bold().Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
		})
		if due == 0 {
			c.Key("all")
			banner := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).
				Margin(t.Space(2), t.Space(3), 0, t.Space(3)).Radius(t.Radius).Background(t.Success.Alpha(0.1)).TextColor(t.Success)
			kit.Arrive(banner, ui.Pose{Scale: 0.95, Fade: 1}, 0).Children(func() {
				ui.Icon(c, icon.Must("circle-check")).Size(16, 16)
				ui.Text(c, o.AllPaid).Bold().FontSize(t.Rem(0.85)).Shrink(0)
			})
		}
		c.Key("rows")
		list := ui.Column(c).FillWidth().Padding(t.Space(1), t.Space(3), t.Space(2), t.Space(3))
		list.Children(func() {
			for i, b := range bills {
				c.Key(b.Name)
				done := (*paid)[i]
				late := b.Overdue && !done
				glyph := b.Icon
				if glyph == "" {
					glyph = "receipt"
				}
				ink := palette(t, i)
				if late {
					ink = t.Danger
				}
				var pay ui.Element
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), 0).
					BorderWidth(0, 0, 1, 0).BorderColor(t.Border.Alpha(unit(i < len(bills)-1)))
				kit.Stagger(c, list, row, i).Children(func() {
					tile(c, glyph, ink)
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, b.Name).Bold().FontSize(t.Rem(0.9)).FillWidth().SingleLine()
						when := kit.Muted(c, b.Due).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
						if late {
							when.TextColor(t.Danger)
						}
					})
					amount := money(c, b.Amount, currency).Bold().FontSize(t.Rem(0.85)).Shrink(0)
					if done {
						amount.TextColor(t.TextMuted).Strikethrough()
						kit.Arrive(kit.Chip(c, o.Paid, t.Success), ui.Pose{Scale: 0.8, Fade: 1}, 0)
						return
					}
					name := fmt.Sprintf(o.PayFor, b.Name)
					pay = kit.Press(ui.ButtonBase(c).Label(name).Radius(999).Shrink(0).TextColor(t.AccentText))
					pay.Background(pay.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(pay.Hovered())), ui.Snappy))
					pay.Children(func() {
						ui.Text(c, o.Pay).FontSize(t.Rem(0.8)).Bold().Padding(t.Space(1), t.Space(3)).Shrink(0)
					})
				})
				if !done && pay.Clicked() {
					(*paid)[i] = true
					root.SetChanged()
					sound.Play(sound.Success)
					kit.Log("Bills", "paid", "index", i)
				}
			}
		})
	})
	return root
}
