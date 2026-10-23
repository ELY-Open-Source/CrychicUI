package gallery

import (
	"fmt"
	"image"

	"github.com/ELY-Open-Source/CrychicUI/shop"
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

func demoShopProduct() shop.Product {
	return shop.Product{Name: "Desk lamp", Price: 80, OldPrice: 100, Stock: 4,
		Image: ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 2, 2)))}
}

func demoProductCard(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	added := boolState(root, "added", false)
	root.Children(func() {
		if root, _ := shop.ProductCard(c, demoShopProduct(), "USD", added, shop.ProductCardOptions{}); root.Changed() {
			log("added")
		}
	})
}

func demoProductGallery(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := intState(root, "selected", 0)
	p := demoShopProduct()
	root.Children(func() {
		if shop.ProductGallery(c, []*ui.Bitmap{p.Image, p.Image}, selected, shop.ProductGalleryOptions{}).Changed() {
			log("selected " + fmt.Sprint(*selected))
		}
	})
}

func demoProductPrice(c *ui.Context, log func(string)) {
	shop.ProductPrice(c, 80, 100, "USD", shop.ProductPriceOptions{})
}

func demoSaleBadge(c *ui.Context, log func(string)) {
	shop.SaleBadge(c, 20, shop.SaleBadgeOptions{})
}

func demoAddToCart(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	added := boolState(root, "added", false)
	root.Children(func() {
		if shop.AddToCart(c, "Add", added, shop.AddToCartOptions{}).Changed() {
			log("added")
		}
	})
}

func demoCartIcon(c *ui.Context, log func(string)) {
	shop.CartIcon(c, 3, shop.CartIconOptions{})
}

func demoCartDrawer(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	checkout := boolState(root, "checkout", false)
	items := []shop.CartItem{{Product: demoShopProduct(), Quantity: 2}}
	root.Children(func() {
		shop.CartDrawer(c, items, "USD", checkout, shop.CartDrawerOptions{})
		if *checkout {
			log("checkout")
			*checkout = false
		}
	})
}

func demoCartRow(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	item := ui.Local(root, "item", func() shop.CartItem { return shop.CartItem{Product: demoShopProduct(), Quantity: 1} })
	removed := boolState(root, "removed", false)
	root.Children(func() {
		c.Key("cart-row")
		if shop.CartRow(c, item, "USD", removed, shop.CartRowOptions{}).Changed() {
			log(fmt.Sprintf("quantity %d removed %v", item.Quantity, *removed))
			*removed = false
		}
	})
}

func demoCheckoutSteps(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	current := ui.Local(root, "current", func() shop.CheckoutStep { return shop.Cart })
	root.Children(func() {
		if shop.CheckoutSteps(c, current, shop.CheckoutStepsOptions{}).Changed() {
			log("step " + fmt.Sprint(*current))
		}
	})
}

func demoOrderSummary(c *ui.Context, log func(string)) {
	lines := []shop.SummaryLine{{Label: "Subtotal", Amount: 100}, {Label: "Discount", Amount: -10}, {Label: "Tax", Amount: 5}}
	shop.OrderSummary(c, lines, "USD", shop.OrderSummaryOptions{})
}

func demoCoupon(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	code, applied := textState(root, "code", "SAVE"), textState(root, "applied", "")
	root.Children(func() {
		_, apply := shop.Coupon(c, code, []string{"SAVE"}, applied, shop.CouponOptions{})
		if apply.Changed() {
			log("applied " + fmt.Sprint(*applied != ""))
		}
	})
}

func demoVariants(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	color, size := intState(root, "color", -1), intState(root, "size", -1)
	colors := []shop.VariantStock{{Name: "Red", Stock: 2}, {Name: "Blue", Stock: 1}}
	sizes := []shop.VariantStock{{Name: "M", Stock: 3}, {Name: "L", Stock: 0}}
	root.Children(func() {
		if shop.Variants(c, colors, sizes, color, size, shop.VariantsOptions{}).Changed() {
			log(fmt.Sprintf("color %d size %d", *color, *size))
		}
	})
}

func demoStock(c *ui.Context, log func(string)) {
	shop.Stock(c, 2, shop.StockOptions{})
}

func demoReviews(c *ui.Context, log func(string)) {
	reviews := []shop.Review{{Rating: 5, Author: "Ana", Text: "Great lamp", Helpful: 2}, {Rating: 4, Author: "Kai", Text: "Warm light"}}
	shop.Reviews(c, reviews, shop.ReviewsOptions{})
}

func demoReviewCard(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	helpful := boolState(root, "helpful", false)
	review := shop.Review{Rating: 5, Author: "Ana", Text: "Great lamp", Helpful: 2}
	root.Children(func() {
		if shop.ReviewCard(c, review, helpful, shop.ReviewCardOptions{}).Changed() {
			log("helpful")
		}
	})
}

func demoWishlist(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	saved := boolState(root, "saved", false)
	root.Children(func() {
		if shop.Wishlist(c, demoShopProduct(), saved, shop.WishlistOptions{}).Changed() {
			log("saved " + fmt.Sprint(*saved))
		}
	})
}

func demoShopCompare(c *ui.Context, log func(string)) {
	products := []shop.Product{demoShopProduct(), demoShopProduct()}
	products[1].Name = "Floor lamp"
	specs := []shop.CompareSpec{{Name: "Material", Values: []string{"Wood", "Metal"}}}
	shop.Compare(c, products, specs, "USD")
}

func demoCategories(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := intState(root, "selected", -1)
	p := demoShopProduct()
	categories := []shop.Category{{Name: "Lighting", Image: p.Image}, {Name: "Desks", Image: p.Image}}
	root.Children(func() {
		if shop.Categories(c, categories, selected).Changed() {
			log("selected " + fmt.Sprint(*selected))
		}
	})
}

func demoShopFilters(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	low, high := numberState(root, "low", 0), numberState(root, "high", 1000)
	rating := intState(root, "rating", 0)
	chosen := ui.Local(root, "chosen", func() []bool { return []bool{false, false} })
	root.Children(func() {
		if shop.ShopFilters(c, low, high, []string{"A", "B"}, chosen, rating, shop.ShopFiltersOptions{}).Changed() {
			log(fmt.Sprintf("rating %d chosen %v", *rating, *chosen))
		}
	})
}

func demoSortBar(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	sort := textState(root, "sort", "Featured")
	root.Children(func() {
		if shop.SortBar(c, 24, sort, []string{"Featured", "Price"}, shop.SortBarOptions{}).Changed() {
			log("sort " + *sort)
		}
	})
}
