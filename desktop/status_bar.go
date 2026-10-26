package desktop

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// StatusItem is a fact; Tone: muted, good, warning, bad.
type StatusItem struct {
	Icon, Text string
	Tone       int
}

// StatusBar is a window's status bar, items left and right.
func StatusBar(c *ui.Context, left, right []StatusItem) ui.Element {
	for _, it := range append(append([]StatusItem{}, left...), right...) {
		if strings.TrimSpace(it.Text) == "" || it.Tone < 0 || it.Tone > 3 {
			panic(fmt.Sprintf("desktop: StatusBar item %q needs text and tone 0 to 3", it.Text))
		}
	}
	t := c.Theme()
	inks := [4]ui.Color{t.TextMuted, t.Success, t.Warning, t.Danger}
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(3)).
		Background(t.SurfaceHover).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Role(ui.RoleStatus)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	// Facts never shrink; a full line wraps.
	fact := func(side, i int, it StatusItem) {
		c.Key([2]int{side, i})
		ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Shrink(0).Children(func() {
			if it.Icon != "" {
				glyph(c, it.Icon, 12).TextColor(inks[it.Tone])
			} else if it.Tone > 0 {
				ui.Box(c).Size(t.Space(1.5), t.Space(1.5)).Shrink(0).Radius(999).Background(inks[it.Tone])
			}
			ui.Text(c, it.Text).FontSize(t.Rem(0.7)).TextColor(t.TextMuted).Shrink(0).MaxWidth(t.Space(40)).SingleLine().
				ContentTransition(ui.ContentTransition{})
		})
	}
	return root.Children(func() {
		for i, it := range left {
			fact(0, i, it)
		}
		c.Key("spacer")
		ui.Box(c).Grow(1)
		for i, it := range right {
			fact(1, i, it)
		}
	})
}
