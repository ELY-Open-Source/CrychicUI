package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AccordionList opens one row's detail at a time.
func AccordionList(c *ui.Context, label string, titles []string, open *int, detail func(i int)) ui.Element {
	if label == "" || len(titles) == 0 || open == nil || detail == nil || *open < -1 || *open >= len(titles) {
		panic(fmt.Sprintf("data: AccordionList %q needs titles, an open row or -1, and a detail", label))
	}
	unique("AccordionList", label, titles)
	t := c.Theme()
	root := frame(c, label).Dividers(1, t.Border)
	shown := *open
	return root.Children(func() {
		for i, title := range titles {
			c.Key(title)
			on := i == shown
			ui.Column(c).FillWidth().Children(func() {
				b := kit.Press(ui.ButtonBase(c).Label(title).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
					Padding(t.Space(2.5), t.Space(3))).Expanded(on)
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy)).Children(func() {
					ui.Text(c, title).Bold().Grow(1)
					chevron(c, on)
				})
				if b.Clicked() {
					*open = i
					if on {
						*open = -1
					}
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("AccordionList", "toggled", "index", i, "open", !on)
				}
				kids(c, on, func() {
					ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(0, t.Space(3), t.Space(3), t.Space(3)).Children(func() { detail(i) })
				})
			})
		}
	})
}
