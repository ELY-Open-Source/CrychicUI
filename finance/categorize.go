package finance

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CategorizeOptions word the card; Icons holds one per category.
type CategorizeOptions struct {
	Title string
	Icons []string
}

// Categorize picks a transaction's category from chips; cue Select.
func Categorize(c *ui.Context, transaction Transaction, categories []string, selected *int, currency string, o CategorizeOptions) ui.Element {
	if transaction.Merchant == "" || len(categories) == 0 || selected == nil || *selected < -1 || *selected >= len(categories) || currency == "" {
		panic("finance: Categorize needs a merchant, categories, a selection and a currency")
	}
	if len(o.Icons) != 0 && len(o.Icons) != len(categories) {
		panic(fmt.Sprintf("finance: Categorize %q needs one icon per category", transaction.Merchant))
	}
	for i, name := range categories {
		if name == "" || slices.Index(categories, name) != i {
			panic(fmt.Sprintf("finance: Categorize %q category %q must be unique and not empty", transaction.Merchant, name))
		}
	}
	def(&o.Title, "Category")
	if *selected >= 0 {
		transaction.Category = categories[*selected]
	}
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Column(c).FillWidth().Padding(t.Space(1), t.Space(3)).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			txRow(c, transaction, currency)
		})
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Children(func() {
			kit.Muted(c, o.Title).FontSize(t.Rem(0.75)).Bold()
			chips := ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5))
			chips.Children(func() {
				for i, name := range categories {
					c.Key(name)
					on := i == *selected
					b := kit.Press(ui.ButtonBase(c).Label(name).Radius(999).Border(1, t.Border).Checked(on))
					b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(t.Accent, unit(on)), ui.Snappy)).
						TextColor(b.AnimateColor("ink", t.Text.Mix(t.AccentText, unit(on)), ui.Snappy))
					kit.Stagger(c, chips, b, i).Children(func() {
						ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(3)).Children(func() {
							glyph := "tag"
							if len(o.Icons) > 0 {
								glyph = o.Icons[i]
							}
							if on {
								glyph = "check"
							}
							ui.Icon(c, icon.Must(glyph)).Size(14, 14).ContentTransition(ui.ContentTransition{})
							ui.Text(c, name).FontSize(t.Rem(0.8)).Bold().Shrink(0).MaxWidth(t.Space(44)).Wrap()
						})
					})
					if b.Clicked() && !on {
						*selected = i
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("Categorize", "selected", "index", i)
					}
				}
			})
		})
	})
	return root
}
