package navigation

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DrilldownOptions tune a Drilldown.
type DrilldownOptions struct {
	// Back names the row leading up.
	Back string
}

// Drilldown slides into submenus; a back row leads out.
func Drilldown(c *ui.Context, label string, list []Node, chosen *string, trail *[]string, o DrilldownOptions) ui.Element {
	if label == "" || len(list) == 0 || chosen == nil || trail == nil {
		panic(fmt.Sprintf("navigation: Drilldown %q needs nodes, a choice and a trail", label))
	}
	nodes("Drilldown", list)
	def(&o.Back, "Back")
	t := c.Theme()
	root := ui.Column(c).Width(t.Space(64)).Padding(t.Space(1)).Radius(t.Radius).Border(1, t.Border).Clip().Label(label)
	level := list
	for _, name := range *trail {
		at := slices.IndexFunc(level, func(n Node) bool { return n.Label == name && n.Children != nil })
		if at < 0 {
			panic(fmt.Sprintf("navigation: Drilldown %q trail %v leads to no submenu", label, *trail))
		}
		level = level[at].Children
	}
	tap := func(event string) {
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("Drilldown", event, "depth", len(*trail))
	}
	return root.Children(func() {
		kit.Slide(c, len(*trail), func() {
			ui.Column(c).FillWidth().Gap(t.Space(0.5)).FocusGroup(ui.Vertical).Children(func() {
				if len(*trail) > 0 {
					c.Key("back")
					up := kit.Press(ui.ButtonBase(c).Label(o.Back).FillWidth().Radius(t.Radius))
					up.Background(up.AnimateColor("bg", t.SurfaceHover.Alpha(unit(up.Hovered())), ui.Snappy)).Children(func() {
						ui.Row(c).Gap(t.Space(2)).Padding(t.Space(1.5), t.Space(2)).AlignItems(ui.Center).TextColor(t.Accent).Children(func() {
							ui.Icon(c, icon.Must("chevron-left")).Size(16, 16)
							ui.Text(c, (*trail)[len(*trail)-1]).Bold().SingleLine()
						})
					})
					if up.Clicked() {
						*trail = (*trail)[:len(*trail)-1]
						tap("up")
					}
				}
				for _, n := range level {
					c.Key(n.Label)
					r := rowBase(c, n.Label)
					r.Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, n.Label).Grow(1).SingleLine()
							if n.Children != nil {
								ui.Icon(c, icon.Must("chevron-right")).Size(16, 16).TextColor(t.TextMuted)
							}
						})
					})
					if r.Clicked() {
						if n.Children != nil {
							*trail = append(*trail, n.Label)
							tap("down")
						} else {
							*chosen = strings.Join(append(append([]string{}, *trail...), n.Label), "/")
							tap("chose")
						}
					}
				}
			})
		})
	})
}
