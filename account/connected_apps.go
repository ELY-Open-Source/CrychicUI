package account

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// App is one connected third-party app; Icon is optional.
type App struct {
	ID, Name, Access, Since, Icon string
}

// ConnectedAppsOptions word the card; Label takes action and app.
type ConnectedAppsOptions struct{ Title, Disconnect, Empty, Label string }

// ConnectedApps disconnects apps; cue Close.
func ConnectedApps(c *ui.Context, apps *[]App, o ConnectedAppsOptions) ui.Element {
	if apps == nil {
		panic("account: ConnectedApps needs apps")
	}
	ids := make([]string, len(*apps))
	for i, a := range *apps {
		if a.Name == "" || a.Access == "" || a.Since == "" {
			panic(fmt.Sprintf("account: ConnectedApps %q needs a name, access and since", a.ID))
		}
		ids[i] = a.ID
	}
	unique("ConnectedApps", ids)
	def(&o.Title, "Connected apps")
	def(&o.Disconnect, "Disconnect")
	def(&o.Empty, "No apps connected")
	def(&o.Label, "%s %s")
	t := c.Theme()
	inks := []ui.Color{t.Accent, t.Success, t.Warning, t.Danger}
	root := card(c).Role(ui.RoleList)
	cut := -1
	root.Children(func() {
		kit.CardHead(c, "plug", o.Title, func() { kit.Chip(c, fmt.Sprint(len(*apps)), t.Accent) })
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			if len(*apps) == 0 {
				empty(c, "unplug", o.Empty)
			}
			for i, a := range *apps {
				c.Key(a.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).
					Transition(rowExit)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					mark := a.Icon
					def(&mark, "plug")
					tile(c, mark, inks[i%len(inks)], t.Space(9))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, a.Name).FontSize(t.Rem(0.85)).Bold().FillWidth().MinWidth(0).SingleLine()
						ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							kit.Muted(c, a.Access).FontSize(t.Rem(0.75)).MinWidth(0).SingleLine()
							dot(c)
							kit.Muted(c, a.Since).FontSize(t.Rem(0.75)).MinWidth(0).SingleLine()
						})
					})
					b := outline(c, fmt.Sprintf(o.Label, o.Disconnect, a.Name), o.Disconnect, "unplug", t.Danger)
					if b.Clicked() {
						cut = i
					}
				})
			}
		})
	})
	if cut >= 0 {
		*apps = slices.Delete(slices.Clone(*apps), cut, cut+1)
		root.SetChanged()
		sound.Play(sound.Close)
		kit.Log("ConnectedApps", "disconnected", "count", len(*apps))
	}
	return root
}
