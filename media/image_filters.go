package media

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Filter is a look ImageFilters applies.
type Filter uint8

const (
	Original Filter = iota
	Mono
	Soft
	Warm
	Cool
	Faded
	filterCount
)

// ImageFiltersOptions tune an ImageFilters.
type ImageFiltersOptions struct {
	// Names label the filters in order; English when empty.
	Names []string
	// Height of the preview; zero gives the default.
	Height float32
}

// ImageFilters previews filter presets on a thumbnail strip.
func ImageFilters(c *ui.Context, label string, image *ui.Bitmap, filter *Filter, o ImageFiltersOptions) ui.Element {
	pictures("ImageFilters", label, []*ui.Bitmap{image})
	if filter == nil || *filter >= filterCount || len(o.Names) != 0 && len(o.Names) != int(filterCount) {
		panic(fmt.Sprintf("media: ImageFilters %q needs a known filter and a name for each or none", label))
	}
	if len(o.Names) == 0 {
		o.Names = []string{"Original", "Mono", "Soft", "Warm", "Cool", "Faded"}
	}
	t := c.Theme()
	size("ImageFilters", label, &o.Height, t.Space(56))
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).Label(label)
	shown := *filter
	pick := func(f Filter) {
		if f >= filterCount || f == shown {
			return
		}
		*filter, shown = f, f
		root.SetChanged()
		sound.Play(sound.Select)
		kit.Log("ImageFilters", "picked", "filter", int(f))
	}
	_, _, prev, next := ends(c)
	fade := ui.ElementTransition{Enter: &ui.Motion{}, Exit: &ui.Motion{}}
	return root.Children(func() {
		main := ui.Grid(c).Columns(1).FillWidth().Height(o.Height).Radius(t.Radius).Clip()
		thumbs := make([]ui.Element, filterCount)
		strip := ui.Row(c).FillWidth().Gap(t.Space(2)).FocusGroup(ui.Horizontal).Children(func() {
			for f := range filterCount {
				c.Key(int(f))
				b := kit.Press(ui.ButtonBase(c).Label(o.Names[f]).Column().Gap(t.Space(1)).AlignItems(ui.Center).Shrink(0))
				thumbs[f] = b.Children(func() {
					ui.Box(c).Size(t.Space(16), t.Space(16)).Radius(t.Radius).Clip().Children(func() { look(c, image, f) })
					ui.Text(c, o.Names[f]).SingleLine().FontSize(t.Rem(0.75))
				})
			}
		})
		for f, b := range thumbs {
			if b.Clicked() {
				pick(Filter(f))
			}
		}
		switch {
		case strip.Shortcut(0, prev) && shown > 0:
			pick(shown - 1)
		case strip.Shortcut(0, next):
			pick(shown + 1)
		}
		for f, b := range thumbs {
			fg := t.TextMuted
			if Filter(f) == shown {
				fg = t.Accent
			}
			b.Checked(Filter(f) == shown).TextColor(b.AnimateColor("fg", fg, ui.Snappy))
		}
		main.Children(func() {
			c.Key(int(shown))
			look(c, image, shown).ColumnStart(1).RowStart(1).Transition(fade)
		})
	})
}
