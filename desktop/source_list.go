package desktop

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SourceItem is a sidebar entry; a positive Count shows.
type SourceItem struct {
	Name, Icon string
	Count      int
}

// SourceGroup is a titled run of entries.
type SourceGroup struct {
	Title string
	Items []SourceItem
}

// SourceList is a grouped sidebar list; cue Select.
func SourceList(c *ui.Context, groups []SourceGroup, current *string) ui.Element {
	titles, names := make([]string, len(groups)), []string{}
	for i, g := range groups {
		for _, it := range g.Items {
			if it.Icon == "" || it.Count < 0 {
				panic(fmt.Sprintf("desktop: SourceList item %q needs an icon and a count of 0 or more", it.Name))
			}
			names = append(names, it.Name)
		}
		titles[i] = g.Title
	}
	unique("SourceList", titles)
	unique("SourceList", names)
	if current == nil || !slices.Contains(names, *current) {
		panic("desktop: SourceList needs a current item among its items")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(2)).Background(t.SurfaceHover).Role(ui.RoleList)
	kit.Arrive(root, ui.Pose{X: -t.Space(3), Fade: 1}, 0)
	root.Children(func() {
		for _, g := range groups {
			c.Key(strconv.Quote(g.Title))
			list := ui.Column(c).FillWidth().Gap(t.Space(0.5))
			list.Children(func() {
				kit.Muted(c, g.Title).FontSize(t.Rem(0.68)).Bold().FillWidth().MinWidth(0).SingleLine().Padding(0, t.Space(2), t.Space(0.5), t.Space(2))
				for i, it := range g.Items {
					c.Key(strconv.Quote(it.Name))
					on := it.Name == *current
					b := ui.ButtonBase(c).Label(it.Name).Role(ui.RoleListItem).Checked(on).FillWidth().Padding(t.Space(1.25), t.Space(2)).Radius(t.Radius)
					b.Background(b.AnimateColor("bg", t.Border.Alpha(0.5*unit(b.Hovered())).Mix(t.Accent.Alpha(0.16), unit(on)), ui.Snappy))
					kit.Stagger(c, list, b, i).Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							glyph(c, it.Icon, 16).TextColor(b.AnimateColor("mark", t.TextMuted.Mix(t.Accent, unit(on)), ui.Snappy))
							ui.Text(c, it.Name).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).SingleLine()
							if it.Count > 0 {
								kit.Muted(c, strconv.Itoa(it.Count)).FontSize(t.Rem(0.7)).Shrink(0).ContentTransition(ui.ContentTransition{})
							}
						})
					})
					if b.Clicked() && !on {
						*current = it.Name
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("SourceList", "picked")
					}
				}
			})
		}
	})
	return root
}
