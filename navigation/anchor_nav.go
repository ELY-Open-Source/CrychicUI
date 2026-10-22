package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AnchorNav is section links that scroll to their sections below.
func AnchorNav(c *ui.Context, label string, sections []string, section func(i int)) ui.Element {
	if label == "" || len(sections) == 0 || section == nil {
		panic(fmt.Sprintf("navigation: AnchorNav %q needs sections and a section builder", label))
	}
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(2)).Label(label)
	st := ui.Local(root, "scroll", func() ui.ScrollState { return ui.ScrollState{} })
	jump := -1
	root.Children(func() {
		ui.Row(c).Gap(t.Space(1)).FocusGroup(ui.Horizontal).Children(func() {
			for i, s := range sections {
				c.Key(i)
				b := kit.Press(ui.ButtonBase(c).Label(s).Padding(t.Space(1), t.Space(2)).Radius(999))
				bg, fg := t.Surface.Alpha(0), t.TextMuted
				if b.Hovered() {
					bg, fg = t.Surface, t.Accent
				}
				b.Background(b.AnimateColor("bg", bg, ui.Snappy)).TextColor(b.AnimateColor("fg", fg, ui.Snappy)).Children(func() {
					ui.Text(c, s).SingleLine()
				})
				if kit.Clicked(b, sound.Tap) {
					jump = i
				}
			}
		})
		tops := sectioned(c, st, sections, section)
		if jump >= 0 {
			st.ScrollTo(0, min(tops()[jump], st.MaxY), ui.Smooth)
			root.SetChanged()
			kit.Log("AnchorNav", "jumped", "index", jump)
		}
	})
	return root
}
