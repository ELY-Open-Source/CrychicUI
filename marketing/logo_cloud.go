package marketing

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/motion"
	"github.com/egoist/mygo/ui"
)

// Logo is a customer mark: an icon and a wordmark.
type Logo struct{ Name, Icon string }

// LogoCloudOptions word the row.
type LogoCloudOptions struct{ Title string }

// LogoCloud scrolls customer logos in a loop; hover pauses it.
func LogoCloud(c *ui.Context, logos []Logo, o LogoCloudOptions) ui.Element {
	def(&o.Title, "Trusted by teams at")
	if len(logos) < 2 {
		panic(fmt.Sprintf("marketing: LogoCloud %q needs two logos or more", o.Title))
	}
	names := make([]string, len(logos))
	for i, l := range logos {
		if l.Icon == "" {
			panic(fmt.Sprintf("marketing: LogoCloud logo %q needs an icon", l.Name))
		}
		names[i] = l.Name
	}
	unique("LogoCloud", names)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(4), 0).Role(ui.RoleGroup).Label(o.Title)
	root.Children(func() {
		kit.Muted(c, o.Title).FontSize(t.Rem(0.76)).Bold().FillWidth().TextAlign(ui.Center).Wrap()
		motion.Marquee(c, motion.MarqueeOptions{}, func() {
			ui.Row(c).Gap(t.Space(8)).AlignItems(ui.Center).Padding(0, t.Space(4)).TextColor(t.TextMuted).Children(func() {
				for _, l := range logos {
					ui.Row(c).Shrink(0).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
						glyph(c, l.Icon, 20)
						ui.Text(c, l.Name).FontSize(t.Rem(1)).Bold().Shrink(0)
					})
				}
			})
		}).FillWidth()
	})
	return root
}
