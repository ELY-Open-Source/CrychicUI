package account

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Invoice is one downloadable billing invoice.
type Invoice struct {
	ID, Number, Date, Amount, URL string
	Paid                          bool
}

// InvoiceListOptions word the card; Label takes action and number.
type InvoiceListOptions struct{ Title, Download, Paid, Due, Label string }

// InvoiceList downloads invoices; cue Tap.
func InvoiceList(c *ui.Context, invoices []Invoice, downloaded *string, o InvoiceListOptions) ui.Element {
	if len(invoices) == 0 || downloaded == nil {
		panic("account: InvoiceList needs invoices and state")
	}
	ids := make([]string, len(invoices))
	for i, inv := range invoices {
		if inv.Number == "" || inv.Date == "" || inv.Amount == "" || inv.URL == "" {
			panic(fmt.Sprintf("account: InvoiceList %q needs a number, date, amount and URL", inv.ID))
		}
		ids[i] = inv.ID
	}
	unique("InvoiceList", ids)
	def(&o.Title, "Invoices")
	def(&o.Download, "Download")
	def(&o.Paid, "Paid")
	def(&o.Due, "Due")
	def(&o.Label, "%s %s")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "receipt", o.Title, func() { kit.Chip(c, fmt.Sprint(len(invoices)), t.Accent) })
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			for i, inv := range invoices {
				c.Key(inv.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					tile(c, "file-text", t.Accent, t.Space(9))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, inv.Number).FontSize(t.Rem(0.85)).Bold().FillWidth().MinWidth(0).SingleLine()
						kit.Muted(c, inv.Date).FontSize(t.Rem(0.75)).FillWidth().MinWidth(0).SingleLine()
					})
					ui.Column(c).Gap(t.Space(0.5)).AlignItems(ui.End).Shrink(0).Children(func() {
						kit.Numeric(c, inv.Amount).FontSize(t.Rem(0.85)).Bold().Shrink(0)
						if inv.Paid {
							kit.Chip(c, o.Paid, t.Success)
						} else {
							kit.Chip(c, o.Due, t.Warning)
						}
					})
					done := *downloaded == inv.ID
					b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Label, o.Download, inv.Number)).Role(ui.RoleLink).
						Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999))
					b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
					b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Accent, unit(b.Hovered())).Mix(t.Success, unit(done)), ui.Snappy))
					b.Children(func() {
						c.Key(done)
						mark := "download"
						if done {
							mark = "check"
						}
						kit.Arrive(glyph(c, mark, 16), ui.Pose{Y: -t.Space(1), Fade: 1}, 0)
					})
					if b.Clicked() {
						c.OpenURL(inv.URL)
						*downloaded = inv.ID
						root.SetChanged()
						sound.Play(sound.Tap)
						kit.Log("InvoiceList", "downloaded", "index", i)
					}
				})
			}
		})
	})
	return root
}
