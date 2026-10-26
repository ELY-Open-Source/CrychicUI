package content

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PrevNextOptions word the two links.
type PrevNextOptions struct{ Prev, Next string }

// PrevNext steps to a neighboring page; cue Tap.
func PrevNext(c *ui.Context, pages []string, current *int, o PrevNextOptions) ui.Element {
	def(&o.Prev, "Previous")
	def(&o.Next, "Next")
	unique("PrevNext", pages)
	if len(pages) == 0 || current == nil || *current < 0 || *current >= len(pages) {
		panic(fmt.Sprintf("content: PrevNext needs pages and a current one, got %d pages", len(pages)))
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).Role(ui.RoleToolbar)
	root.Children(func() {
		link := func(key string, to int, label, mark string, end bool) {
			c.Key(key)
			if to < 0 || to >= len(pages) {
				ui.Box(c).Grow(1).Basis(t.Space(40))
				return
			}
			b := kit.Press(ui.ButtonBase(c).Label(pages[to]).Grow(1).Basis(t.Space(40)).MinWidth(0).Padding(t.Space(3), t.Space(4)).Radius(t.Radius * 1.5))
			hover := unit(b.Hovered())
			b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, hover), ui.Snappy))
			b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, hover), ui.Snappy))
			shift := b.AnimateTo("nudge", hover, ui.Snappy)
			b.Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
					arrow := func() {
						dir := float32(1)
						if !end {
							dir = -1
						}
						glyph(c, mark, 18).TextColor(t.Accent).Translate(dir*t.Space(1)*shift, 0)
					}
					if !end {
						arrow()
					}
					align := ui.Start
					if end {
						align = ui.End
					}
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).AlignItems(align).Children(func() {
						kit.Muted(c, label).FontSize(t.Rem(0.72)).MaxWidth(t.Space(60)).SingleLine()
						ui.Text(c, pages[to]).FontSize(t.Rem(0.9)).Bold().TextColor(t.Accent).FillWidth().MinWidth(0).TextAlign(align).SingleLine()
					})
					if end {
						arrow()
					}
				})
			})
			if b.Clicked() {
				*current = to
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("PrevNext", "page", "index", to)
			}
		}
		link("prev", *current-1, o.Prev, "chevron-left", false)
		link("next", *current+1, o.Next, "chevron-right", true)
	})
	return root
}
