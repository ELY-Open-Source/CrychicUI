package desktop

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TitleTabsOptions word the new-tab button.
type TitleTabsOptions struct{ New string }

// TitleTabs are tabs in a title bar; cue Select.
func TitleTabs(c *ui.Context, tabs []string, picked *int, o TitleTabsOptions) (root, add ui.Element) {
	def(&o.New, "New tab")
	unique("TitleTabs", tabs)
	if len(tabs) == 0 || picked == nil || *picked < 0 || *picked >= len(tabs) {
		panic(fmt.Sprintf("desktop: TitleTabs needs tabs and a picked one, got %d tabs", len(tabs)))
	}
	t := c.Theme()
	root = ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.End).Padding(t.Space(1.5), t.Space(2), 0, t.Space(2)).
		Background(t.SurfaceHover).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Role(ui.RoleTabList)
	kit.Arrive(root, ui.Pose{Y: -t.Space(1), Fade: 1}, 0)
	root.Children(func() {
		lights(c, true)
		strip := ui.Row(c).Grow(1).MinWidth(0).Gap(t.Space(1)).AlignItems(ui.End)
		strip.Children(func() {
			for i, name := range tabs {
				c.Key(strconv.Quote(name))
				on := i == *picked
				b := ui.ButtonBase(c).Label(name).Role(ui.RoleTab).Checked(on).Grow(1).Basis(0).MinWidth(0).MaxWidth(t.Space(48)).
					Padding(t.Space(2), t.Space(3)).Radius(t.Radius, t.Radius, 0, 0)
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Mix(t.Border, 0.4*unit(b.Hovered())).Mix(t.Background, unit(on)), ui.Snappy))
				b.Children(func() {
					ink := t.TextMuted.Mix(t.Text, unit(on))
					ui.Text(c, name).FontSize(t.Rem(0.78)).TextColor(b.AnimateColor("ink", ink, ui.Snappy)).FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine()
				})
				if b.Clicked() && !on {
					*picked = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("TitleTabs", "picked", "index", i)
				}
			}
		})
		c.Key("add")
		add = kit.Press(ui.ButtonBase(c).Label(o.New).Shrink(0).Padding(t.Space(1.5)).Margin(0, 0, t.Space(1), 0).Radius(999))
		add.Background(add.AnimateColor("bg", t.Border.Alpha(0.6*unit(add.Hovered())), ui.Snappy))
		add.Children(func() { glyph(c, "plus", 14).TextColor(t.TextMuted) })
		if add.Clicked() {
			sound.Play(sound.Tap)
			kit.Log("TitleTabs", "new tab")
		}
	})
	return root, add
}
