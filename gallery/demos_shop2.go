package gallery

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/shop"
	"github.com/egoist/mygo/ui"
)

var shopEntries2 = []Entry{
	{ID: 996, Name: "OrderStatus", Category: "shop", Summary: "Placed, packed, shipped, delivered track.", Demo: demoOrderStatus},
	{ID: 997, Name: "OrderCard", Category: "shop", Summary: "Number, date, total, status; cue Tap.", Demo: demoOrderCard},
	{ID: 998, Name: "Tracking", Category: "shop", Summary: "Carrier events on a line.", Demo: demoTracking},
	{ID: 999, Name: "Delivery", Category: "shop", Summary: "Speed options with dates and prices; cue Select.", Demo: demoDelivery},
	{ID: 1000, Name: "FreeShipping", Category: "shop", Summary: "Bar toward the free-shipping amount.", Demo: demoFreeShipping},
	{ID: 1001, Name: "Upsell", Category: "shop", Summary: "You may also like row; cue Pop.", Demo: demoUpsell},
	{ID: 1002, Name: "Recently", Category: "shop", Summary: "Recently viewed strip.", Demo: demoRecently},
	{ID: 1003, Name: "FlashSale", Category: "shop", Summary: "Sale header with a countdown.", Demo: demoFlashSale},
	{ID: 1004, Name: "GiftCard", Category: "shop", Summary: "Card face with amount and message.", Demo: demoGiftCard},
	{ID: 1005, Name: "Loyalty", Category: "shop", Summary: "Points, tier and bar to the next.", Demo: demoLoyalty},
	{ID: 1006, Name: "Bundle", Category: "shop", Summary: "Pick items for a bundle price; cue Toggle.", Demo: demoBundle},
	{ID: 1007, Name: "MiniCart", Category: "shop", Summary: "Popover of the cart; cue Open.", Demo: demoMiniCart},
	{ID: 1008, Name: "Receipt", Category: "shop", Summary: "Itemized receipt with a total.", Demo: demoShopReceipt},
	{ID: 1009, Name: "Returns", Category: "shop", Summary: "Items to return with reasons; cue Select.", Demo: demoReturns},
	{ID: 1010, Name: "Seller", Category: "shop", Summary: "Seller avatar, rating, follow; cue Toggle.", Demo: demoSeller},
	{ID: 1011, Name: "QandA", Category: "shop", Summary: "Questions with answers and ask; cue Send.", Demo: demoQandA},
	{ID: 1012, Name: "SizeGuide", Category: "shop", Summary: "Table of sizes in a popover; cue Open.", Demo: demoSizeGuide},
	{ID: 1013, Name: "Specs", Category: "shop", Summary: "Specification rows by group.", Demo: demoSpecs},
	{ID: 1014, Name: "StoreHours", Category: "shop", Summary: "Open now badge and week hours.", Demo: demoStoreHours},
	{ID: 1015, Name: "ThankYou", Category: "shop", Summary: "Order confirmed with number; check draws; cue Success.", Demo: demoThankYou},
}

// demoOrderDay is when the demo order was placed.
var demoOrderDay = time.Date(2026, 10, 20, 9, 12, 0, 0, time.Local)

func demoCartItems() []shop.CartItem {
	p := demoProducts()
	return []shop.CartItem{{Product: p[0], Quantity: 1}, {Product: p[1], Quantity: 2}, {Product: p[5], Quantity: 1}}
}

func demoOrderStatus(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		shop.OrderStatus(c, shop.Shipped, shop.OrderStatusOptions{Dates: [4]string{"Oct 20", "Oct 21", "Oct 22", "Est. Oct 26"}})
	})
}

func demoOrderCard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	opened := boolState(root, "opened", false)
	p := demoProducts()
	o := shop.OrderCardOptions{Images: []*ui.Bitmap{p[0].Image, p[1].Image, p[5].Image}}
	root.Children(func() {
		if shop.OrderCard(c, "Order #NW-20417", demoOrderDay, 173, "USD", shop.Shipped, opened, o).Changed() {
			log("opened order")
		}
	})
}

func demoTracking(c *ui.Context, log func(string)) {
	at := func(day, hour, minute int) time.Time { return time.Date(2026, 10, day, hour, minute, 0, 0, time.Local) }
	events := []shop.TrackEvent{
		{Label: "Picked up by carrier", Location: "Portland, OR", Time: at(21, 9, 12), Done: true},
		{Label: "Departed sorting facility", Location: "Portland, OR", Time: at(21, 18, 40), Done: true},
		{Label: "Arrived at regional hub", Location: "Reno, NV", Time: at(22, 7, 5), Done: true},
		{Label: "Out for delivery", Location: "San Francisco, CA", Time: at(26, 8, 0)},
		{Label: "Delivered", Location: "San Francisco, CA", Time: at(26, 17, 0)},
	}
	ui.Column(c).Width(420).Children(func() { shop.Tracking(c, events, shop.TrackingOptions{}) })
}

func demoDelivery(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	selected := intState(root, "selected", 0)
	choices := []shop.DeliveryOption{{Name: "Standard", Date: "Arrives Thu, Oct 29"},
		{Name: "Express", Date: "Arrives Tue, Oct 27", Price: 12}, {Name: "Next day", Date: "Arrives Mon, Oct 26", Price: 24}}
	root.Children(func() {
		if shop.Delivery(c, "Delivery speed", choices, "USD", selected, shop.DeliveryOptions{}).Changed() {
			log("selected " + fmt.Sprint(*selected))
		}
	})
}

func demoFreeShipping(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	subtotal := numberState(root, "subtotal", 61)
	root.Children(func() {
		shop.FreeShipping(c, *subtotal, 75, "USD", shop.FreeShippingOptions{})
		if button.Outline(c, "Add a mug").Clicked() {
			*subtotal += 24
			log(fmt.Sprintf("subtotal %.0f", *subtotal))
		}
	})
}

func demoUpsell(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	added := ui.Local(root, "added", func() []bool { return make([]bool, 3) })
	p := demoProducts()
	root.Children(func() {
		if shop.Upsell(c, "You may also like", []shop.Product{p[1], p[5], p[3]}, "USD", added, shop.UpsellOptions{}).Changed() {
			n := 0
			for _, on := range *added {
				if on {
					n++
				}
			}
			log(fmt.Sprintf("%d added", n))
		}
	})
}

func demoRecently(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() { shop.Recently(c, "Recently viewed", demoProducts(), "USD") })
}

func demoFlashSale(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	ends := ui.Local(root, "ends", func() time.Time { return c.Now().Add(2*time.Hour + 14*time.Minute) })
	root.Children(func() {
		shop.FlashSale(c, "Autumn lighting sale", *ends, shop.FlashSaleOptions{Subtitle: "Up to 40% off lamps and pendants"})
	})
}

func demoGiftCard(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Children(func() {
		shop.GiftCard(c, 50, "USD", "Happy housewarming! Pick something cozy for the new place.",
			shop.GiftCardOptions{From: "From Maya & Daniel", Code: "NW-GIFT-4821"})
	})
}

func demoLoyalty(c *ui.Context, log func(string)) {
	ui.Column(c).Width(400).Children(func() {
		shop.Loyalty(c, 1240, "Gold member", 260, shop.LoyaltyOptions{NextTier: "Platinum"})
	})
}

func demoBundle(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	selected := ui.Local(root, "selected", func() []bool { return []bool{true, true, true} })
	p := demoProducts()
	root.Children(func() {
		if shop.Bundle(c, []shop.Product{p[0], p[3], p[1]}, selected, "USD", shop.BundleOptions{Off: 0.1}).Changed() {
			n := 0
			for _, on := range *selected {
				if on {
					n++
				}
			}
			log(fmt.Sprintf("%d selected", n))
		}
	})
}

func demoMiniCart(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	open := boolState(root, "open", false)
	root.Children(func() {
		if shop.MiniCart(c, demoCartItems(), "USD", open, shop.MiniCartOptions{}).Changed() {
			log("open " + fmt.Sprint(*open))
		}
	})
}

func demoShopReceipt(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Children(func() {
		shop.Receipt(c, "NW-20417", demoCartItems(), "USD", shop.ReceiptOptions{Meta: "Oct 20, 2026 · Visa •••• 4242"})
	})
}

func demoReturns(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	selected := intState(root, "selected", 0)
	reason := textState(root, "reason", "")
	p := demoProducts()
	reasons := []string{"Arrived damaged", "Wrong item", "No longer needed", "Better price elsewhere"}
	root.Children(func() {
		if shop.Returns(c, []shop.Product{p[0], p[1], p[5]}, "USD", reasons, selected, reason, shop.ReturnsOptions{}).Changed() {
			log(fmt.Sprintf("item %d reason %s", *selected, *reason))
		}
	})
}

func demoSeller(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	followed := boolState(root, "followed", false)
	o := shop.SellerOptions{Meta: "12.4k sales · Ships from Portland, OR", Verified: true}
	root.Children(func() {
		if shop.Seller(c, "Lumen Lighting", 4.8, followed, o).Changed() {
			log("followed " + fmt.Sprint(*followed))
		}
	})
}

func demoQandA(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	questions := ui.Local(root, "questions", func() []shop.Question {
		return []shop.Question{
			{Text: "Does the arm lock in place?", Answer: "Yes. The joints have friction locks and hold any angle."},
			{Text: "Can I use a smart bulb with it?"},
			{Text: "How long is the cord?", Answer: "1.8 m, with an inline switch."},
		}
	})
	draft := textState(root, "draft", "")
	root.Children(func() {
		edit, ask := shop.QandA(c, questions, draft, shop.QandAOptions{})
		switch {
		case ask.Changed():
			log(fmt.Sprint(len(*questions)) + " questions")
		case edit.Changed():
			log(fmt.Sprintf("draft %d chars", utf8.RuneCountInString(*draft)))
		}
	})
}

func demoSizeGuide(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	open := boolState(root, "open", false)
	rows := []shop.SizeRow{{Size: "XS", Chest: "81–86 cm", Waist: "66–71 cm"}, {Size: "S", Chest: "86–91 cm", Waist: "71–76 cm"},
		{Size: "M", Chest: "96–101 cm", Waist: "81–86 cm"}, {Size: "L", Chest: "106–111 cm", Waist: "91–96 cm"},
		{Size: "XL", Chest: "116–121 cm", Waist: "101–106 cm"}}
	root.Children(func() {
		if shop.SizeGuide(c, rows, open, shop.SizeGuideOptions{}).Changed() {
			log("open " + fmt.Sprint(*open))
		}
	})
}

func demoSpecs(c *ui.Context, log func(string)) {
	groups := []shop.SpecGroup{
		{Name: "Light", Icon: "lightbulb", Rows: []shop.SpecRow{{Name: "Brightness", Value: "450 lm"},
			{Name: "Color temperature", Value: "2700–4000 K, three steps"}, {Name: "Bulb", Value: "Integrated LED, 6 W"}}},
		{Name: "Dimensions", Icon: "ruler", Rows: []shop.SpecRow{{Name: "Height", Value: "48 cm"},
			{Name: "Arm reach", Value: "42 cm"}, {Name: "Weight", Value: "2.1 kg"}}},
		{Name: "Materials", Icon: "layers", Rows: []shop.SpecRow{{Name: "Body", Value: "Brass-plated steel"},
			{Name: "Shade", Value: "Spun aluminum"}}},
	}
	ui.Column(c).Width(440).Children(func() { shop.Specs(c, groups) })
}

func demoStoreHours(c *ui.Context, log func(string)) {
	hours := []shop.DayHours{
		{Day: "Sunday", Open: "10:00", Close: "16:00"}, {Day: "Monday", Open: "09:00", Close: "18:00"},
		{Day: "Tuesday", Open: "09:00", Close: "18:00"}, {Day: "Wednesday"},
		{Day: "Thursday", Open: "09:00", Close: "20:00"}, {Day: "Friday", Open: "09:00", Close: "20:00"},
		{Day: "Saturday", Open: "10:00", Close: "18:00"},
	}
	now := time.Date(2026, 10, 23, 14, 30, 0, 0, time.Local)
	ui.Column(c).Width(380).Children(func() { shop.StoreHours(c, now, hours, shop.StoreHoursOptions{}) })
}

func demoThankYou(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() {
		shop.ThankYou(c, "NW-20417", shop.ThankYouOptions{
			Note: "A receipt is on its way to maya.chen@example.com. We'll email tracking once it ships."})
	})
}
