package desktop

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// InspectorRow is a labeled value.
type InspectorRow struct{ Label, Value string }

// InspectorSection is a titled group of rows.
type InspectorSection struct {
	Title string
	Rows  []InspectorRow
}

// InspectorPane is folding sections of values; folds report; cue Toggle.
func InspectorPane(c *ui.Context, sections []InspectorSection) ui.Element {
	titles := make([]string, len(sections))
	for i, s := range sections {
		labels := make([]string, len(s.Rows))
		for j, r := range s.Rows {
			labels[j] = r.Label
		}
		if len(s.Rows) == 0 {
			panic(fmt.Sprintf("desktop: InspectorPane section %q needs rows", s.Title))
		}
		unique("InspectorPane", labels)
		titles[i] = s.Title
	}
	unique("InspectorPane", titles)
	if len(sections) == 0 {
		panic("desktop: InspectorPane needs sections")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Background(t.SurfaceHover).BorderWidth(0, 0, 0, 1).BorderColor(t.Border).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{X: t.Space(3), Fade: 1}, 0)
	closed := ui.Local(root, "closed", func() map[string]bool { return map[string]bool{} })
	root.Children(func() {
		for _, s := range sections {
			c.Key(strconv.Quote(s.Title))
			ui.Column(c).FillWidth().BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
				shown := !(*closed)[s.Title]
				h := ui.ButtonBase(c).Label(s.Title).Expanded(shown).FillWidth().Padding(t.Space(2), t.Space(3))
				h.Background(h.AnimateColor("bg", t.Border.Alpha(0.4*unit(h.Hovered())), ui.Snappy))
				turn := h.AnimateTo("turn", unit(shown), ui.Snappy)
				h.Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
						glyph(c, "chevron-right", 12).TextColor(t.TextMuted).Rotate(90 * turn)
						ui.Text(c, s.Title).FontSize(t.Rem(0.74)).Bold().Grow(1).MinWidth(0).SingleLine()
						kit.Muted(c, strconv.Itoa(len(s.Rows))).FontSize(t.Rem(0.68)).Shrink(0)
					})
				})
				if h.Clicked() {
					(*closed)[s.Title] = shown
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("InspectorPane", "section", "open", !shown)
				}
				layout.Collapse(c, shown, func() {
					ui.Column(c).FillWidth().Gap(t.Space(1.5)).Padding(0, t.Space(3), t.Space(3), t.Space(3)).Children(func() {
						for _, r := range s.Rows {
							c.Key(strconv.Quote(r.Label))
							ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
								kit.Muted(c, r.Label).FontSize(t.Rem(0.74)).Width(t.Space(22)).Shrink(0).SingleLine()
								ui.Text(c, r.Value).Font("monospace").FontSize(t.Rem(0.74)).Grow(1).MinWidth(0).Padding(t.Space(0.5), t.Space(1.5)).
									Radius(t.Radius*0.5).Background(t.Background).Border(1, t.Border).SingleLine()
							})
						}
					})
				})
			})
		}
	})
	return root
}
