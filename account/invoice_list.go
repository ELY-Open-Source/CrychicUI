package account

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Invoice is one downloadable billing invoice.
type Invoice struct {
	ID, Number, Date, Amount, URL string
}

// InvoiceListOptions are the invoice labels.
type InvoiceListOptions struct{ Download string }

// InvoiceList downloads invoices; cue Tap.
func InvoiceList(c *ui.Context, invoices []Invoice, downloaded *string, o InvoiceListOptions) ui.Element {
	if len(invoices) == 0 || downloaded == nil {
		panic("account: InvoiceList needs invoices and state")
	}
	def(&o.Download, "Download")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleList)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for _, invoice := range invoices {
			if invoice.ID == "" || invoice.Number == "" || invoice.Date == "" || invoice.Amount == "" || invoice.URL == "" {
				panic("account: invalid invoice")
			}
			c.Key(invoice.ID)
			ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Column(c).Grow(1).Children(func() {
					ui.Text(c, invoice.Number)
					kit.Muted(c, invoice.Date+" · "+invoice.Amount)
				})
				b := button(c, invoice.ID+" "+o.Download).Role(ui.RoleLink)
				b.Children(func() { ui.Text(c, o.Download) })
				if b.Clicked() {
					c.OpenURL(invoice.URL)
					*downloaded = invoice.ID
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("InvoiceList", "downloaded")
				}
			})
		}
	})
	return root
}
