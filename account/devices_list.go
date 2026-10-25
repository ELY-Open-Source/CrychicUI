package account

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Device is one remembered device; Icon is optional.
type Device struct {
	ID, Name, Last, Icon string
}

// DevicesListOptions word the card; Seen and Label are formats.
type DevicesListOptions struct{ Title, Remove, Empty, Seen, Label string }

// DevicesList removes devices; cue Delete.
func DevicesList(c *ui.Context, devices *[]Device, o DevicesListOptions) ui.Element {
	if devices == nil {
		panic("account: DevicesList needs devices")
	}
	ids := make([]string, len(*devices))
	for i, d := range *devices {
		if d.Name == "" || d.Last == "" {
			panic(fmt.Sprintf("account: DevicesList %q needs a name and last seen", d.ID))
		}
		ids[i] = d.ID
	}
	unique("DevicesList", ids)
	def(&o.Title, "Devices")
	def(&o.Remove, "Remove")
	def(&o.Empty, "No remembered devices")
	def(&o.Seen, "Last seen %s")
	def(&o.Label, "%s %s")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	remove := -1
	root.Children(func() {
		kit.CardHead(c, "laptop", o.Title, func() { kit.Chip(c, fmt.Sprint(len(*devices)), t.Accent) })
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			if len(*devices) == 0 {
				empty(c, "laptop", o.Empty)
			}
			for i, d := range *devices {
				c.Key(d.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).
					Transition(rowExit)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					mark := d.Icon
					def(&mark, "laptop")
					tile(c, mark, t.Accent, t.Space(9))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, d.Name).FontSize(t.Rem(0.85)).Bold().FillWidth().MinWidth(0).SingleLine()
						kit.Muted(c, fmt.Sprintf(o.Seen, d.Last)).FontSize(t.Rem(0.75)).FillWidth().MinWidth(0).SingleLine()
					})
					b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Label, o.Remove, d.Name)).Size(t.Space(8), t.Space(8)).
						Shrink(0).Center().Radius(999))
					b.Background(b.AnimateColor("bg", t.Danger.Alpha(0).Mix(t.Danger.Alpha(0.12), unit(b.Hovered())), ui.Snappy))
					b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Danger, unit(b.Hovered())), ui.Snappy))
					b.Children(func() { glyph(c, "trash-2", 16) })
					if b.Clicked() {
						remove = i
					}
				})
			}
		})
	})
	if remove >= 0 {
		*devices = slices.Delete(slices.Clone(*devices), remove, remove+1)
		root.SetChanged()
		sound.Play(sound.Delete)
		kit.Log("DevicesList", "removed", "count", len(*devices))
	}
	return root
}
