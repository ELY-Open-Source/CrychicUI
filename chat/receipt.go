package chat

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/egoist/mygo/ui"
)

// Delivery is how far a message has gone.
type Delivery uint8

// Delivery states.
const (
	Sent Delivery = iota
	Delivered
	Read
	deliveryCount
)

// ReceiptOptions name the states.
type ReceiptOptions struct {
	Sent, Delivered, Read string
}

// Receipt shows sent, delivered and read checks; checks cross-fade.
func Receipt(c *ui.Context, state Delivery, o ReceiptOptions) ui.Element {
	if state >= deliveryCount {
		panic(fmt.Sprintf("chat: Receipt state %d unknown", state))
	}
	def(&o.Sent, "Sent")
	def(&o.Delivered, "Delivered")
	def(&o.Read, "Read")
	t := c.Theme()
	name, glyph, ink := o.Sent, "check", t.TextMuted
	switch state {
	case Delivered:
		name, glyph = o.Delivered, "check-check"
	case Read:
		name, glyph, ink = o.Read, "check-check", t.Accent
	}
	root := ui.Box(c).AlignSelf(ui.Start).Role(ui.RoleStatus).Label(name).Tooltip(name)
	root.TextColor(root.AnimateColor("ink", ink, ui.Snappy))
	return root.Children(func() {
		ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{})
	})
}
