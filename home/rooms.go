package home

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Room is one tab; On counts devices switched on.
type Room struct {
	Name, Icon string
	On         int
}

// RoomsOptions word the tabs; On takes the count.
type RoomsOptions struct{ Label, On string }

// Rooms picks a room from a tab strip; cue Select.
func Rooms(c *ui.Context, rooms []Room, picked *string, o RoomsOptions) ui.Element {
	def(&o.Label, "Rooms")
	def(&o.On, "%d on")
	names := make([]string, len(rooms))
	for i, r := range rooms {
		if r.Icon == "" || r.On < 0 {
			panic(fmt.Sprintf("home: Rooms room %q needs an icon and an on count of 0 or more", r.Name))
		}
		names[i] = r.Name
	}
	if len(rooms) == 0 || picked == nil || !slices.Contains(names, *picked) {
		panic(fmt.Sprintf("home: Rooms %q needs rooms and a picked name among them", o.Label))
	}
	unique("Rooms", names)
	t := c.Theme()
	root := ui.ScrollHorizontal(c).FillWidth().Role(ui.RoleGroup).Label(o.Label)
	root.Children(func() {
		strip := ui.Row(c).Gap(t.Space(2)).Padding(t.Space(1), 0)
		strip.Children(func() {
			for i, r := range rooms {
				on := r.Name == *picked
				c.Key(strconv.Quote(r.Name))
				b := kit.Press(ui.ButtonBase(c).Label(r.Name).Checked(on).Shrink(0).Padding(t.Space(2), t.Space(3.5)).Radius(999))
				b.Border(1, b.AnimateColor("edge", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent, unit(on)), ui.Snappy))
				b.TextColor(b.AnimateColor("ink", t.Text.Mix(t.AccentText, unit(on)), ui.Snappy))
				kit.Stagger(c, strip, b, i).Children(func() {
					ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
						glyph(c, r.Icon, 16)
						ui.Text(c, r.Name).FontSize(t.Rem(0.8)).Bold().Shrink(0)
						if r.On > 0 {
							c.Key("count")
							kit.Numeric(c, fmt.Sprintf(o.On, r.On)).FontSize(t.Rem(0.66)).Shrink(0).Padding(0, t.Space(1.5)).Radius(999).
								Background(t.Background.Alpha(0.2*unit(on)).Mix(t.SurfaceHover, 1-unit(on)))
						}
					})
				})
				if b.Clicked() && !on {
					*picked = r.Name
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Rooms", "picked", "index", i)
				}
			}
		})
	})
	return root
}
