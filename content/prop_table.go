package content

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Prop is one documented property.
type Prop struct {
	Name, Type, Default, Text string
	Required                  bool
}

// PropTableOptions word the table; Default formats a default.
type PropTableOptions struct{ Title, Required, Default string }

// PropTable lists props with types and defaults.
func PropTable(c *ui.Context, props []Prop, o PropTableOptions) ui.Element {
	def(&o.Title, "Props")
	def(&o.Required, "Required")
	def(&o.Default, "Default: %s")
	names := make([]string, len(props))
	for i, p := range props {
		if strings.TrimSpace(p.Type) == "" {
			panic(fmt.Sprintf("content: PropTable prop %q needs a type", p.Name))
		}
		names[i] = p.Name
	}
	unique("PropTable", names)
	t := c.Theme()
	root := card(c).Role(ui.RoleTable).Label(o.Title)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(4), t.Space(4), t.Space(2), t.Space(4)).Children(func() {
			tile(c, "braces", t.Accent, t.Space(8))
			ui.Text(c, o.Title).FontSize(t.Rem(0.95)).Bold().Grow(1).MinWidth(0).SingleLine()
			kit.Chip(c, strconv.Itoa(len(props)), t.TextMuted)
		})
		list := ui.Column(c).FillWidth()
		list.Children(func() {
			for i, p := range props {
				c.Key(strconv.Quote(p.Name))
				row := ui.Column(c).FillWidth().Gap(t.Space(1.5)).Padding(t.Space(3), t.Space(4)).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Role(ui.RoleRow)
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						ui.Text(c, p.Name).Font("monospace").FontSize(t.Rem(0.84)).Bold().Shrink(0).MaxWidth(t.Space(60)).SingleLine()
						ui.Text(c, p.Type).Font("monospace").FontSize(t.Rem(0.72)).TextColor(t.Accent).Padding(t.Space(0.25), t.Space(1.5)).
							Radius(t.Radius).Background(t.Accent.Alpha(0.1)).Shrink(0).MaxWidth(t.Space(60)).SingleLine()
						if p.Required {
							ui.Text(c, o.Required).FontSize(t.Rem(0.66)).Bold().TextColor(t.Danger).Shrink(0).MaxWidth(t.Space(30)).SingleLine()
						}
					})
					if p.Text != "" {
						kit.Muted(c, p.Text).FontSize(t.Rem(0.8)).FillWidth().Wrap()
					}
					if p.Default != "" {
						ui.Text(c, fmt.Sprintf(o.Default, p.Default)).Font("monospace").FontSize(t.Rem(0.72)).TextColor(t.TextMuted).
							FillWidth().MinWidth(0).SingleLine()
					}
				})
			}
		})
	})
	return root
}
