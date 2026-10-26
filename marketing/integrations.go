package marketing

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Integration is one app tile.
type Integration struct {
	Name, Text, Icon string
	Connected        bool
}

// IntegrationsOptions word the section; Connected tags linked apps.
type IntegrationsOptions struct{ Title, Text, Connected string }

// Integrations is a grid of app tiles; cue Tap.
func Integrations(c *ui.Context, apps []Integration, opened *string, o IntegrationsOptions) ui.Element {
	def(&o.Connected, "Connected")
	names := make([]string, len(apps))
	for i, a := range apps {
		if strings.TrimSpace(a.Icon) == "" {
			panic(fmt.Sprintf("marketing: Integrations %q app %q needs an icon", o.Title, a.Name))
		}
		names[i] = a.Name
	}
	if len(apps) == 0 || opened == nil || *opened != "" && !slices.Contains(names, *opened) {
		panic(fmt.Sprintf("marketing: Integrations %q needs apps and an opened app among them", o.Title))
	}
	unique("Integrations", names)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		heading(c, o.Title, o.Text)
		c.Key("grid")
		grid := ui.Row(c).FillWidth().Wrap().Gap(t.Space(3))
		grid.Children(func() {
			for i, a := range apps {
				on := a.Name == *opened
				c.Key(strconv.Quote(a.Name))
				b := kit.Press(ui.ButtonBase(c).Label(a.Name).Checked(on).Column().AlignItems(ui.Stretch).Grow(1).Basis(t.Space(52)).
					MinWidth(0).Gap(t.Space(2)).Padding(t.Space(4)).Radius(t.Radius * 1.5))
				lit := unit(on || b.Hovered())
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.Accent, 0.04*lit), ui.Snappy))
				b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, lit), ui.Snappy))
				kit.Stagger(c, grid, b, i).Children(func() {
					ui.Row(c).FillWidth().AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
						tile(c, a.Icon, t.Accent, t.Space(10))
						if a.Connected {
							kit.Chip(c, o.Connected, t.Success)
						}
					})
					ui.Text(c, a.Name).FontSize(t.Rem(0.9)).Bold().FillWidth().MinWidth(0).SingleLine()
					if a.Text != "" {
						kit.Muted(c, a.Text).FontSize(t.Rem(0.78)).FillWidth().Wrap()
					}
				})
				if b.Clicked() {
					sound.Play(sound.Tap)
					kit.Log("Integrations", "opened", "index", i)
					if !on {
						*opened = a.Name
						root.SetChanged()
					}
				}
			}
		})
	})
	return root
}
