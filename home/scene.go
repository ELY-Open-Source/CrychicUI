package home

import (
	"fmt"
	"math"
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SceneItem is one scene button.
type SceneItem struct{ Name, Icon string }

// SceneOptions word the grid; Runs counts every press.
type SceneOptions struct {
	Title string
	Runs  *int
}

// Scene runs one of the home's scenes; cue Tap.
func Scene(c *ui.Context, scenes []SceneItem, active *string, o SceneOptions) ui.Element {
	def(&o.Title, "Scenes")
	names := make([]string, len(scenes))
	for i, s := range scenes {
		if s.Icon == "" {
			panic(fmt.Sprintf("home: Scene %q needs an icon", s.Name))
		}
		names[i] = s.Name
	}
	if len(scenes) == 0 || active == nil || *active != "" && !slices.Contains(names, *active) {
		panic(fmt.Sprintf("home: Scene %q needs scenes and an active name among them", o.Title))
	}
	unique("Scene", names)
	if o.Runs != nil && *o.Runs < 0 {
		panic(fmt.Sprintf("home: Scene %q needs a run count of 0 or more", o.Title))
	}
	// A full run count shuts the buttons; it never wraps.
	full := o.Runs != nil && *o.Runs == math.MaxInt
	t := c.Theme()
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		kit.CardHead(c, "sparkles", o.Title, nil)
		grid := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(2)).Padding(t.Space(3))
		grid.Children(func() {
			for i, s := range scenes {
				on := s.Name == *active
				c.Key(strconv.Quote(s.Name))
				b := kit.Press(ui.ButtonBase(c).Label(s.Name).Checked(on).MinWidth(0).Padding(t.Space(2.5)).Radius(t.Radius).
					AlignItems(ui.Stretch).Disabled(full))
				b.Border(1, b.AnimateColor("edge", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent, 0.1*unit(on)), ui.Snappy))
				kit.Stagger(c, grid, b, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						tile(c, s.Icon, t.Accent, t.Space(8), on)
						ui.Text(c, s.Name).FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0).SingleLine()
					})
				})
				if b.Clicked() && !full {
					if o.Runs != nil {
						*o.Runs++
					}
					sound.Play(sound.Tap)
					kit.Log("Scene", "ran", "index", i)
					if !on {
						*active = s.Name
						root.SetChanged()
					}
				}
			}
		})
	})
	return root
}
