package gallery

import (
	"fmt"
	"time"

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

func demoOrderStatus(c *ui.Context, log func(string)) {
	shop.OrderStatus(c, shop.Shipped, shop.OrderStatusOptions{})
}

func demoOrderCard(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	opened := boolState(root, "opened", false)
	root.Children(func() {
		if shop.OrderCard(c, "Order 1042", time.Now(), 95, "USD", shop.Shipped, opened, shop.OrderCardOptions{}).Changed() {
			log("opened order")
		}
	})
}

func demoTracking(c *ui.Context, log func(string)) {
	events := []shop.TrackEvent{{Label: "Shipped", Location: "Shanghai", Time: time.Now(), Done: true}}
	shop.Tracking(c, events)
}

func demoDelivery(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := intState(root, "selected", 0)
	choices := []shop.DeliveryOption{{Name: "Standard", Date: "Friday", Price: 5}, {Name: "Express", Date: "Tomorrow", Price: 12}}
	root.Children(func() {
		if shop.Delivery(c, "Delivery", choices, "USD", selected).Changed() {
			log("selected " + fmt.Sprint(*selected))
		}
	})
}

func demoFreeShipping(c *ui.Context, log func(string)) {
	shop.FreeShipping(c, 80, 100, "USD", shop.FreeShippingOptions{})
}

func demoUpsell(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	added := intState(root, "added", -1)
	root.Children(func() {
		if shop.Upsell(c, "Related", []shop.Product{demoShopProduct()}, "USD", added, shop.UpsellOptions{}).Changed() {
			log("added " + fmt.Sprint(*added))
		}
	})
}

func demoRecently(c *ui.Context, log func(string)) {
	shop.Recently(c, "Recently viewed", []shop.Product{demoShopProduct()}, "USD")
}

func demoFlashSale(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	ends := ui.Local(root, "ends", func() time.Time { return c.Now().Add(90 * time.Minute) })
	root.Children(func() {
		shop.FlashSale(c, "Flash sale", *ends, shop.FlashSaleOptions{})
	})
}

func demoGiftCard(c *ui.Context, log func(string)) {
	shop.GiftCard(c, 50, "USD", "For you")
}

func demoLoyalty(c *ui.Context, log func(string)) {
	shop.Loyalty(c, 200, "Gold", 50, shop.LoyaltyOptions{})
}

func demoBundle(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := ui.Local(root, "selected", func() []bool { return []bool{true} })
	root.Children(func() {
		if shop.Bundle(c, []shop.Product{demoShopProduct()}, selected, "USD", shop.BundleOptions{}).Changed() {
			log(fmt.Sprint(*selected))
		}
	})
}

func demoMiniCart(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := boolState(root, "open", false)
	items := []shop.CartItem{{Product: demoShopProduct(), Quantity: 1}}
	root.Children(func() {
		if shop.MiniCart(c, items, "USD", open, shop.MiniCartOptions{}).Changed() {
			log("open " + fmt.Sprint(*open))
		}
	})
}

func demoShopReceipt(c *ui.Context, log func(string)) {
	items := []shop.CartItem{{Product: demoShopProduct(), Quantity: 2}}
	shop.Receipt(c, "R-42", items, "USD", shop.ReceiptOptions{})
}

func demoReturns(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := intState(root, "selected", -1)
	reason := textState(root, "reason", "")
	root.Children(func() {
		if shop.Returns(c, []shop.Product{demoShopProduct()}, []string{"Damaged", "Wrong size"}, selected, reason, shop.ReturnsOptions{}).Changed() {
			log(fmt.Sprintf("selected %d reason %s", *selected, *reason))
		}
	})
}

func demoSeller(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	followed := boolState(root, "followed", false)
	root.Children(func() {
		if shop.Seller(c, "Mei Shop", 4.8, followed, shop.SellerOptions{}).Changed() {
			log("followed " + fmt.Sprint(*followed))
		}
	})
}

func demoQandA(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	questions := ui.Local(root, "questions", func() []shop.Question { return nil })
	draft := textState(root, "draft", "")
	root.Children(func() {
		_, ask := shop.QandA(c, questions, draft, shop.QandAOptions{})
		if ask.Changed() {
			log(fmt.Sprint(len(*questions)) + " questions")
		}
	})
}

func demoSizeGuide(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := boolState(root, "open", false)
	rows := []shop.SizeRow{{Size: "M", Chest: "96 cm", Waist: "80 cm"}}
	root.Children(func() {
		if shop.SizeGuide(c, rows, open, shop.SizeGuideOptions{}).Changed() {
			log("open " + fmt.Sprint(*open))
		}
	})
}

func demoSpecs(c *ui.Context, log func(string)) {
	groups := []shop.SpecGroup{{Name: "Dimensions", Rows: []shop.SpecRow{{Name: "Height", Value: "40 cm"}}}}
	shop.Specs(c, groups)
}

func demoStoreHours(c *ui.Context, log func(string)) {
	hours := []shop.DayHours{
		{Day: "Sun", Open: "00:00", Close: "23:59"}, {Day: "Mon", Open: "00:00", Close: "23:59"},
		{Day: "Tue", Open: "00:00", Close: "23:59"}, {Day: "Wed", Open: "00:00", Close: "23:59"},
		{Day: "Thu", Open: "00:00", Close: "23:59"}, {Day: "Fri", Open: "00:00", Close: "23:59"},
		{Day: "Sat", Open: "00:00", Close: "23:59"},
	}
	shop.StoreHours(c, time.Now(), hours, shop.StoreHoursOptions{})
}

func demoThankYou(c *ui.Context, log func(string)) {
	shop.ThankYou(c, "1042", shop.ThankYouOptions{})
}
