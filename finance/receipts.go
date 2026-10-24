package finance

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Receipt is one scanned receipt; Items counts its lines.
type Receipt struct {
	Merchant, Date string
	Items          int
	Amount         float64
}

// ReceiptsOptions word the card; Item and Items take counts.
type ReceiptsOptions struct{ Title, Total, Item, Items, Empty string }

// Receipts lists scanned receipts and their total.
func Receipts(c *ui.Context, receipts []Receipt, currency string, o ReceiptsOptions) ui.Element {
	if currency == "" {
		panic("finance: Receipts needs a currency")
	}
	total := 0.0
	for _, r := range receipts {
		if r.Merchant == "" || r.Date == "" || r.Items < 0 || r.Amount < 0 || !finite(r.Amount) {
			panic(fmt.Sprintf("finance: Receipts receipt %q needs a date and amounts of zero or more", r.Merchant))
		}
		total += r.Amount
	}
	if !finite(total) {
		panic("finance: Receipts total overflow")
	}
	def(&o.Title, "Receipts")
	def(&o.Total, "Total")
	def(&o.Item, "%d item")
	def(&o.Items, "%d items")
	def(&o.Empty, "No receipts scanned")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "scan-line", o.Title, func() { kit.Chip(c, strconv.Itoa(len(receipts)), t.TextMuted) })
		list := ui.Column(c).FillWidth().Padding(t.Space(1), t.Space(3))
		list.Children(func() {
			if len(receipts) == 0 {
				ui.Column(c).Key("empty").FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(5)).Transition(kit.Enter).
					Children(func() {
						ui.Icon(c, icon.Must("receipt")).Size(22, 22).TextColor(t.TextMuted)
						kit.Muted(c, o.Empty).FontSize(t.Rem(0.85))
					})
			}
			for i, r := range receipts {
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), 0).
					BorderWidth(0, 0, 1, 0).BorderColor(t.Border.Alpha(unit(i < len(receipts)-1)))
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Box(c).Size(t.Space(8), t.Space(10)).Shrink(0).Center().Radius(t.Space(1)).Background(t.Surface).
						Border(1, t.Border).TextColor(palette(t, i)).Children(func() {
						ui.Icon(c, icon.Must("receipt-text")).Size(16, 16)
					})
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, r.Merchant).Bold().FontSize(t.Rem(0.85)).FillWidth().SingleLine()
						meta := r.Date
						if r.Items == 1 {
							meta += " · " + fmt.Sprintf(o.Item, r.Items)
						} else if r.Items > 1 {
							meta += " · " + fmt.Sprintf(o.Items, r.Items)
						}
						kit.Muted(c, meta).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
					})
					money(c, r.Amount, currency).Bold().FontSize(t.Rem(0.85)).Shrink(0)
				})
			}
		})
		ui.Row(c).FillWidth().AlignItems(ui.Center).Padding(t.Space(3)).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			ui.Text(c, o.Total).Bold().Grow(1).MinWidth(0).SingleLine()
			money(c, total, currency).Bold().FontSize(t.Rem(1.1)).Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
		})
	})
	return root
}
