package input

import (
	"strings"
	"unicode"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SlugFieldOptions tune a SlugField.
type SlugFieldOptions struct {
	// Label is caller copy.
	Label string
}

// SlugField derives a slug from a title until locked.
func SlugField(c *ui.Context, title string, slug *string, locked *bool, o SlugFieldOptions) ui.Element {
	if slug == nil || locked == nil {
		panic("input: SlugField needs slug and lock state")
	}
	if o.Label == "" {
		o.Label = "Slug"
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	lastTitle := ui.Local(root, "title", func() string { return "" })
	return root.Children(func() {
		ui.Text(c, o.Label).FontWeight(500)
		if !*locked && title != *lastTitle {
			*slug = slugify(title)
		}
		*lastTitle = title
		row := ui.Row(c).FillWidth().AlignItems(ui.Center).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		row.Children(func() {
			in := ui.TextInput(c, slug).Grow(1).Placeholder(o.Label).Padding(t.Space(2), t.Space(3))
			if changed(root, in) {
				*slug = slugify(*slug)
				*locked = true
			}
			lock := ui.ButtonBase(c).Label("Lock " + o.Label).Padding(t.Space(2)).TextColor(t.TextMuted)
			kit.Press(lock).Children(func() {
				name := "lock-open"
				if *locked {
					name = "lock"
				}
				ui.Icon(c, icon.Must(name)).Size(16, 16).ContentTransition(ui.ContentTransition{})
			})
			if kit.Clicked(lock, sound.Toggle) {
				*locked = !*locked
				root.SetChanged()
				kit.Log("SlugField", "locked", "locked", *locked)
			}
		})
	})
}

// slugify writes lower-case words joined by hyphens.
func slugify(s string) string {
	var out strings.Builder
	lastDash := true
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			out.WriteRune(r)
			lastDash = false
		} else if !lastDash {
			out.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(out.String(), "-")
}
