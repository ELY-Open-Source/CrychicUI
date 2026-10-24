package gallery

import (
	"fmt"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/shop"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

var shopEntries = []Entry{
	{ID: 976, Name: "ProductCard", Category: "shop", Summary: "Image, name, price, add; image zooms; cue Pop.", Demo: demoProductCard},
	{ID: 977, Name: "ProductGallery", Category: "shop", Summary: "Main image and thumbs; cue Select.", Demo: demoProductGallery},
	{ID: 978, Name: "ProductPrice", Category: "shop", Summary: "Price, old price and saving badge.", Demo: demoProductPrice},
	{ID: 979, Name: "SaleBadge", Category: "shop", Summary: "-20% badge; wobbles in.", Demo: demoSaleBadge},
	{ID: 980, Name: "AddToCart", Category: "shop", Summary: "Adds with a fly-to-cart dot; cue Pop.", Demo: demoAddToCart},
	{ID: 981, Name: "CartIcon", Category: "shop", Summary: "Bag with a count that rolls.", Demo: demoCartIcon},
	{ID: 982, Name: "CartDrawer", Category: "shop", Summary: "Items, quantities, subtotal, checkout; cue Open.", Demo: demoCartDrawer},
	{ID: 983, Name: "CartRow", Category: "shop", Summary: "Image, name, quantity, remove; collapses out; cue Delete.", Demo: demoCartRow},
	{ID: 984, Name: "CheckoutSteps", Category: "shop", Summary: "Cart, shipping, payment, done; cue Tick.", Demo: demoCheckoutSteps},
	{ID: 985, Name: "OrderSummary", Category: "shop", Summary: "Lines, discounts, tax, total; total rolls.", Demo: demoOrderSummary},
	{ID: 986, Name: "Coupon", Category: "shop", Summary: "Code field that applies or errors; cue Success, Error.", Demo: demoCoupon},
	{ID: 987, Name: "Variants", Category: "shop", Summary: "Color and size selectors with stock; cue Select.", Demo: demoVariants},
	{ID: 988, Name: "Stock", Category: "shop", Summary: "In stock, low, out with dot.", Demo: demoStock},
	{ID: 989, Name: "Reviews", Category: "shop", Summary: "Average, bars per star and list.", Demo: demoReviews},
	{ID: 990, Name: "ReviewCard", Category: "shop", Summary: "Stars, author, text, helpful; cue Pop.", Demo: demoReviewCard},
	{ID: 991, Name: "Wishlist", Category: "shop", Summary: "Heart toggles saved; cue Toggle.", Demo: demoWishlist},
	{ID: 992, Name: "Compare", Category: "shop", Summary: "Products side by side with specs.", Demo: demoShopCompare},
	{ID: 993, Name: "Categories", Category: "shop", Summary: "Image tiles of categories; cue Tap.", Demo: demoCategories},
	{ID: 994, Name: "ShopFilters", Category: "shop", Summary: "Price range, brands, rating; cue Toggle.", Demo: demoShopFilters},
	{ID: 995, Name: "SortBar", Category: "shop", Summary: "Count and sort select; cue Select.", Demo: demoSortBar},
}

// demoProducts is one home goods store's catalog.
func demoProducts() []shop.Product {
	s := productShots()
	return []shop.Product{
		{Name: "Arc desk lamp", Brand: "Lumen & Co", Price: 89, OldPrice: 119, Stock: 6, Image: s[0], Rating: 4.6, Reviews: 128},
		{Name: "Stoneware mug", Brand: "Kiln Studio", Price: 24, Stock: 3, Image: s[1], Rating: 4.8, Reviews: 412},
		{Name: "Globe pendant", Brand: "Lumen & Co", Price: 149, Stock: 12, Image: s[2], Rating: 4.4, Reviews: 86},
		{Name: "Oak desk tray", Brand: "Grain", Price: 42, OldPrice: 52, Stock: 20, Image: s[3], Rating: 4.7, Reviews: 203},
		{Name: "Linen floor cushion", Brand: "Fold", Price: 58, Image: s[4], Rating: 4.2, Reviews: 37},
		{Name: "Glazed bud vase", Brand: "Kiln Studio", Price: 36, Stock: 9, Image: s[5], Rating: 4.5, Reviews: 64},
	}
}

func demoShopProduct() shop.Product { return demoProducts()[0] }

func demoProductCard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320)
	added := boolState(root, "added", false)
	root.Children(func() {
		if root, _ := shop.ProductCard(c, demoShopProduct(), "USD", added, shop.ProductCardOptions{}); root.Changed() {
			log("added")
		}
	})
}

func demoProductGallery(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	selected := intState(root, "selected", 0)
	s := productShots()
	root.Children(func() {
		if shop.ProductGallery(c, []*ui.Bitmap{s[0], s[6], s[7]}, selected, shop.ProductGalleryOptions{}).Changed() {
			log("selected " + fmt.Sprint(*selected))
		}
	})
}

func demoProductPrice(c *ui.Context, log func(string)) {
	ui.Column(c).Width(320).Children(func() { shop.ProductPrice(c, 89, 119, "USD", shop.ProductPriceOptions{}) })
}

func demoSaleBadge(c *ui.Context, log func(string)) {
	ui.Row(c).Width(320).Gap(12).Children(func() {
		for _, off := range []int{10, 25, 40} {
			c.Key(off)
			shop.SaleBadge(c, off, shop.SaleBadgeOptions{})
		}
	})
}

func demoAddToCart(c *ui.Context, log func(string)) {
	root := ui.Row(c).Width(420).Gap(12).AlignItems(ui.Center)
	added := boolState(root, "added", false)
	p := demoShopProduct()
	root.Children(func() {
		ui.Image(c, p.Image).Size(56, 56).Fit(ui.Cover).Radius(8)
		ui.Column(c).Grow(1).MinWidth(0).Gap(2).Children(func() {
			ui.Text(c, p.Name).Bold()
			typography.Money(c, p.Price, "USD")
		})
		if shop.AddToCart(c, "Add Arc desk lamp", added, shop.AddToCartOptions{}).Changed() {
			log("added")
		}
	})
}

func demoCartIcon(c *ui.Context, log func(string)) {
	root := ui.Row(c).Width(320).Gap(16).AlignItems(ui.Center)
	count := intState(root, "count", 3)
	root.Children(func() {
		shop.CartIcon(c, *count, shop.CartIconOptions{})
		if button.Outline(c, "Add item").Clicked() {
			*count++
			log(fmt.Sprintf("count %d", *count))
		}
	})
}

func demoCartDrawer(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	checkout := boolState(root, "checkout", false)
	p := demoProducts()
	items := []shop.CartItem{{Product: p[0], Quantity: 1}, {Product: p[1], Quantity: 2}, {Product: p[5], Quantity: 1}}
	root.Children(func() {
		shop.CartDrawer(c, items, "USD", checkout, shop.CartDrawerOptions{})
		if *checkout {
			log("checkout")
			*checkout = false
		}
	})
}

func demoCartRow(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	item := ui.Local(root, "item", func() shop.CartItem { return shop.CartItem{Product: demoProducts()[1], Quantity: 1} })
	removed := boolState(root, "removed", false)
	root.Children(func() {
		if !*removed {
			c.Key("cart-row")
			if shop.CartRow(c, item, "USD", removed, shop.CartRowOptions{}).Changed() {
				log(fmt.Sprintf("quantity %d removed %v", item.Quantity, *removed))
			}
			return
		}
		if button.Outline(c, "Undo remove").Clicked() {
			*removed = false
			log("restored")
		}
	})
}

func demoCheckoutSteps(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	current := ui.Local(root, "current", func() shop.CheckoutStep { return shop.Shipping })
	root.Children(func() {
		if shop.CheckoutSteps(c, current, shop.CheckoutStepsOptions{}).Changed() {
			log("step " + fmt.Sprint(*current))
		}
	})
}

func demoOrderSummary(c *ui.Context, log func(string)) {
	lines := []shop.SummaryLine{{Label: "Subtotal", Amount: 207}, {Label: "Spring sale", Amount: -31.05},
		{Label: "Shipping", Amount: 0}, {Label: "Tax", Amount: 14.08}}
	ui.Column(c).Width(320).Children(func() { shop.OrderSummary(c, lines, "USD", shop.OrderSummaryOptions{}) })
}

func demoCoupon(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	code, applied := textState(root, "code", "SPRING15"), textState(root, "applied", "")
	root.Children(func() {
		edit, apply := shop.Coupon(c, code, []string{"SPRING15", "WELCOME10"}, applied, shop.CouponOptions{})
		switch {
		case apply.Changed():
			log("applied " + fmt.Sprint(*applied != ""))
		case edit.Changed():
			log(fmt.Sprintf("code %d chars", utf8.RuneCountInString(*code)))
		}
	})
}

func demoVariants(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	color, size := intState(root, "color", -1), intState(root, "size", 1)
	colors := []shop.VariantStock{{Name: "Brass", Stock: 2, Swatch: ui.RGB(196, 150, 80)},
		{Name: "Matte black", Stock: 8, Swatch: ui.RGB(48, 48, 52)}, {Name: "Opal white", Stock: 0, Swatch: ui.RGB(236, 234, 226)}}
	sizes := []shop.VariantStock{{Name: "Small", Stock: 4}, {Name: "Medium", Stock: 12}, {Name: "Large", Stock: 0}}
	root.Children(func() {
		if shop.Variants(c, colors, sizes, color, size, shop.VariantsOptions{}).Changed() {
			log(fmt.Sprintf("color %d size %d", *color, *size))
		}
	})
}

func demoStock(c *ui.Context, log func(string)) {
	ui.Column(c).Width(320).Gap(10).Children(func() {
		for _, n := range []int{24, 3, 0} {
			c.Key(n)
			shop.Stock(c, n, shop.StockOptions{})
		}
	})
}

func demoReviews(c *ui.Context, log func(string)) {
	reviews := []shop.Review{
		{Rating: 5, Author: "Maya Chen", Title: "Perfect reading light", Text: "The arm holds any angle and the warm setting is easy on the eyes.", Helpful: 14},
		{Rating: 5, Author: "Daniel Okafor", Text: "Heavier than it looks, in a good way. The brass matches my shelf."},
		{Rating: 4, Author: "Priya Raman", Text: "Lovely light. It sits a little high for a small desk."},
		{Rating: 2, Author: "Sofia Rossi", Title: "Switch feels cheap", Text: "The switch wobbled after a month."},
	}
	ui.Column(c).Width(480).Children(func() { shop.Reviews(c, reviews, shop.ReviewsOptions{}) })
}

func demoReviewCard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	helpful := boolState(root, "helpful", false)
	review := shop.Review{Rating: 5, Author: "Maya Chen", Title: "Perfect reading light",
		Text: "The arm holds any angle and the warm setting is easy on the eyes late at night. Worth every cent.", Helpful: 14}
	root.Children(func() {
		if shop.ReviewCard(c, review, helpful, shop.ReviewCardOptions{}).Changed() {
			log("helpful")
		}
	})
}

func demoWishlist(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380).Gap(8)
	saved := ui.Local(root, "saved", func() []bool { return []bool{false, true, false} })
	p := demoProducts()
	root.Children(func() {
		for i, product := range []shop.Product{p[0], p[5], p[2]} {
			c.Key(product.Name)
			if shop.Wishlist(c, product, "USD", &(*saved)[i], shop.WishlistOptions{}).Changed() {
				log(fmt.Sprintf("saved %d %v", i, (*saved)[i]))
			}
		}
	})
}

func demoShopCompare(c *ui.Context, log func(string)) {
	s := productShots()
	products := []shop.Product{demoShopProduct(),
		{Name: "Matte task lamp", Brand: "Northlight", Price: 69, Stock: 14, Image: s[6], Rating: 4.3, Reviews: 54},
		{Name: "Opal desk lamp", Brand: "Hale", Price: 99, Stock: 7, Image: s[7], Rating: 4.5, Reviews: 210}}
	specs := []shop.CompareSpec{
		{Name: "Brightness", Values: []string{"450 lm", "380 lm", "520 lm"}},
		{Name: "Color temp", Values: []string{"2700–4000 K", "3000 K", "2700–6500 K"}},
		{Name: "Dimmable", Values: []string{"Yes", "No", "Yes"}},
		{Name: "Material", Values: []string{"Brass", "Steel", "Aluminum"}},
		{Name: "Warranty", Values: []string{"2 years", "2 years", "2 years"}},
	}
	ui.Column(c).Width(520).Children(func() { shop.Compare(c, products, specs, "USD") })
}

func demoCategories(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	selected := intState(root, "selected", 0)
	s := productShots()
	categories := []shop.Category{{Name: "Lighting", Image: s[0], Count: 124}, {Name: "Tableware", Image: s[1], Count: 86},
		{Name: "Ceiling", Image: s[2], Count: 42}, {Name: "Desks", Image: s[3], Count: 57},
		{Name: "Textiles", Image: s[4], Count: 38}, {Name: "Decor", Image: s[5], Count: 73}}
	root.Children(func() {
		if shop.Categories(c, categories, selected, shop.CategoriesOptions{}).Changed() {
			log("selected " + fmt.Sprint(*selected))
		}
	})
}

func demoShopFilters(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320)
	low, high := numberState(root, "low", 40), numberState(root, "high", 320)
	rating := intState(root, "rating", 4)
	chosen := ui.Local(root, "chosen", func() []bool { return []bool{true, false, false, false} })
	brands := []string{"Lumen & Co", "Kiln Studio", "Fold", "Grain"}
	root.Children(func() {
		o := shop.ShopFiltersOptions{Max: 500, Step: 5}
		if shop.ShopFilters(c, low, high, brands, chosen, rating, o).Changed() {
			n := 0
			for _, on := range *chosen {
				if on {
					n++
				}
			}
			log(fmt.Sprintf("rating %d brands %d", *rating, n))
		}
	})
}

func demoSortBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	sort := textState(root, "sort", "Featured")
	choices := []string{"Featured", "Price: low to high", "Price: high to low", "Top rated", "Newest"}
	root.Children(func() {
		if shop.SortBar(c, 128, sort, choices, shop.SortBarOptions{}).Changed() {
			log("sort " + *sort)
		}
	})
}
