package media

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ThumbStripOptions tune a ThumbStrip.
type ThumbStripOptions struct {
	// Height of the main image; zero gives the default.
	Height float32
	// Item names a thumbnail by number; "Photo %d" by default.
	Item string
}

// ThumbStrip shows the image at picked by its thumbnail.
func ThumbStrip(c *ui.Context, label string, images []*ui.Bitmap, at *int, o ThumbStripOptions) ui.Element {
	pictures("ThumbStrip", label, images)
	chosen("ThumbStrip", label, at, 0, len(images))
	t := c.Theme()
	size("ThumbStrip", label, &o.Height, t.Space(56))
	def(&o.Item, "Photo %d")
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Label(label)
	shown := *at
	pick := func(i int) {
		if i < 0 || i >= len(images) || i == shown {
			return
		}
		*at, shown = i, i
		root.SetChanged()
		sound.Play(sound.Select)
		kit.Log("ThumbStrip", "picked", "index", i)
	}
	_, _, prev, next := ends(c)
	fade := ui.ElementTransition{Enter: &ui.Motion{}, Exit: &ui.Motion{}}
	return root.Children(func() {
		main := ui.Grid(c).Columns(1).FillWidth().Height(o.Height).Radius(t.Radius).Clip().Background(t.Surface)
		strip := ui.Row(c).FillWidth().Gap(t.Space(2)).FocusGroup(ui.Horizontal)
		thumbs := make([]ui.Element, len(images))
		strip.Children(func() {
			for i, img := range images {
				c.Key(i)
				name := fmt.Sprintf(o.Item, i+1)
				thumbs[i] = kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(16), t.Space(12)).Shrink(0).
					Radius(t.Radius).Clip()).Children(func() { ui.Image(c, img).Fill().Fit(ui.Cover) })
			}
		})
		for i, b := range thumbs {
			if b.Clicked() {
				pick(i)
			}
		}
		switch {
		case strip.Shortcut(0, prev):
			pick(shown - 1)
		case strip.Shortcut(0, next):
			pick(shown + 1)
		}
		for i, b := range thumbs {
			edge := ui.Color{}
			if i == shown {
				edge = t.Accent
			}
			b.Checked(i == shown).Border(2, b.AnimateColor("edge", edge, ui.Snappy)).
				Opacity(b.AnimateTo("lit", min(1, 0.55+0.3*unit(b.Hovered())+0.45*unit(i == shown)), ui.Snappy))
		}
		main.Children(func() {
			c.Key(shown)
			ui.Image(c, images[shown]).ColumnStart(1).RowStart(1).Fill().Label(fmt.Sprintf(o.Item, shown+1)).
				Transition(fade)
		})
	})
}
