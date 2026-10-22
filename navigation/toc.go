package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TOCOptions tune a TOC.
type TOCOptions struct {
	// Title heads the list of headings.
	Title string
}

// TOC lists headings beside their sections; the marker follows scrolling.
func TOC(c *ui.Context, label string, headings []string, section func(i int), o TOCOptions) ui.Element {
	if label == "" || len(headings) == 0 || section == nil {
		panic(fmt.Sprintf("navigation: TOC %q needs headings and a section builder", label))
	}
	def(&o.Title, "On this page")
	t := c.Theme()
	root := ui.Row(c).AlignItems(ui.Stretch).Gap(t.Space(4)).Label(label)
	st := ui.Local(root, "scroll", func() ui.ScrollState { return ui.ScrollState{} })
	var tops func() []float32
	active := 0
	root.Children(func() {
		tops = sectioned(c, st, headings, section)
		active = reading(tops(), st)
		nav := ui.Column(c).Width(t.Space(44)).Shrink(0).Gap(t.Space(0.5)).FocusGroup(ui.Vertical)
		id := nav.ID()
		nav.Children(func() {
			kit.Muted(c, o.Title).FontSize(t.Rem(0.75)).Bold().Padding(t.Space(1), t.Space(2.5))
			for i, h := range headings {
				c.Key(i)
				b := kit.Press(ui.ButtonBase(c).Label(h).FillWidth().Checked(i == active))
				fg := t.TextMuted
				if i == active || b.Hovered() {
					fg = t.Text
				}
				b.TextColor(b.AnimateColor("fg", fg, ui.Snappy)).Children(func() {
					ui.Row(c).FillWidth().AlignItems(ui.Stretch).Children(func() {
						bar := ui.Box(c).Width(2).Radius(999)
						if i == active {
							bar.Background(t.Accent).Transition(glide(id))
						}
						ui.Text(c, h).SingleLine().Padding(t.Space(1), t.Space(2.5))
					})
				})
				if kit.Clicked(b, sound.Tap) {
					st.ScrollTo(0, min(tops()[i], st.MaxY), ui.Smooth)
					root.SetChanged()
					kit.Log("TOC", "jumped", "index", i)
				}
			}
		})
	})
	kit.Settle(c, root, []float32{float32(active)}, func() []float32 { return []float32{float32(reading(tops(), st))} })
	return root
}

// reading is the section the scroll has reached.
func reading(tops []float32, st *ui.ScrollState) int {
	if len(tops) > 1 && tops[len(tops)-1] == 0 {
		return 0
	}
	if st.MaxY > 0 && st.Y >= st.MaxY-1 {
		return len(tops) - 1
	}
	at := 0
	for i, top := range tops {
		if top <= st.Y+8 {
			at = i
		}
	}
	return at
}
