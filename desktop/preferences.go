package desktop

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PrefTab is a preferences pane's name and icon.
type PrefTab struct{ Name, Icon string }

// Preferences is a window of icon tabs; cue Select.
func Preferences(c *ui.Context, tabs []PrefTab, picked *int, body func()) ui.Element {
	names := make([]string, len(tabs))
	for i, tb := range tabs {
		if tb.Icon == "" {
			panic(fmt.Sprintf("desktop: Preferences tab %q needs an icon", tb.Name))
		}
		names[i] = tb.Name
	}
	unique("Preferences", names)
	if len(tabs) == 0 || picked == nil || *picked < 0 || *picked >= len(tabs) || body == nil {
		panic(fmt.Sprintf("desktop: Preferences needs tabs, a picked one and a body, got %d tabs", len(tabs)))
	}
	t := c.Theme()
	root := window(c).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Column(c).FillWidth().Background(t.SurfaceHover).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			titleBar(c, tabs[*picked].Name, true, nil)
			ui.ScrollHorizontal(c).FillWidth().Role(ui.RoleTabList).Children(func() {
				ui.Row(c).Gap(t.Space(1)).Padding(0, t.Space(3), t.Space(2), t.Space(3)).Children(func() {
					for i, tb := range tabs {
						c.Key(strconv.Quote(tb.Name))
						on := i == *picked
						b := kit.Press(ui.ButtonBase(c).Label(tb.Name).Role(ui.RoleTab).Checked(on).Shrink(0).Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius))
						b.Background(b.AnimateColor("bg", t.Border.Alpha(0.5*unit(b.Hovered())).Mix(t.Border, unit(on)), ui.Snappy))
						ink := t.TextMuted.Mix(t.Accent, unit(on))
						b.Children(func() {
							ui.Column(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
								glyph(c, tb.Icon, 20).TextColor(b.AnimateColor("ink", ink, ui.Snappy))
								ui.Text(c, tb.Name).FontSize(t.Rem(0.7)).TextColor(t.Text.Mix(t.Accent, unit(on))).Shrink(0).MaxWidth(t.Space(24)).SingleLine()
							})
						})
						if b.Clicked() && !on {
							*picked = i
							root.SetChanged()
							sound.Play(sound.Select)
							kit.Log("Preferences", "pane", "index", i)
						}
					}
				})
			})
		})
		c.Key(strconv.Quote(tabs[*picked].Name))
		kit.Arrive(ui.Column(c).FillWidth().Padding(t.Space(5)), ui.Pose{Y: t.Space(1), Fade: 1}, 0).Children(body)
	})
	return root
}
