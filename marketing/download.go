package marketing

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Platform is one download target.
type Platform struct{ Name, Icon, Note, URL string }

// DownloadOptions word it; Get takes a platform, Opened each URL.
type DownloadOptions struct {
	Title, Text, Get string
	Opened           *string
}

// Download offers the app per platform; cue Tap.
func Download(c *ui.Context, platforms []Platform, o DownloadOptions) ui.Element {
	def(&o.Get, "Download for %s")
	if len(platforms) == 0 {
		panic(fmt.Sprintf("marketing: Download %q needs platforms", o.Title))
	}
	names := make([]string, len(platforms))
	for i, p := range platforms {
		if strings.TrimSpace(p.Icon) == "" || strings.TrimSpace(p.URL) == "" {
			panic(fmt.Sprintf("marketing: Download %q platform %q needs an icon and a URL", o.Title, p.Name))
		}
		names[i] = p.Name
	}
	unique("Download", names)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		heading(c, o.Title, o.Text)
		c.Key("grid")
		grid := ui.Row(c).FillWidth().Wrap().Gap(t.Space(3))
		grid.Children(func() {
			for i, p := range platforms {
				c.Key(strconv.Quote(p.Name))
				b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Get, p.Name)).Column().AlignItems(ui.Stretch).Grow(1).Basis(t.Space(52)).
					MinWidth(0).Gap(t.Space(2)).Padding(t.Space(4)).Radius(t.Radius * 1.5))
				hover := b.AnimateTo("hover", unit(b.Hovered()), ui.Snappy)
				b.Background(t.Background.Mix(t.Accent, 0.04*hover)).Border(1, t.Border.Mix(t.Accent, hover))
				kit.Stagger(c, grid, b, i).Children(func() {
					ui.Row(c).FillWidth().AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
						tile(c, p.Icon, t.Accent, t.Space(10))
						glyph(c, "download", 18).TextColor(t.TextMuted.Mix(t.Accent, hover)).Translate(0, t.Space(0.75)*hover)
					})
					ui.Text(c, p.Name).FontSize(t.Rem(0.92)).Bold().FillWidth().MinWidth(0).SingleLine()
					if p.Note != "" {
						kit.Muted(c, p.Note).FontSize(t.Rem(0.74)).FillWidth().MinWidth(0).SingleLine()
					}
				})
				if b.Clicked() {
					c.OpenURL(p.URL)
					if o.Opened != nil {
						*o.Opened = p.URL
					}
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("Download", "opened", "index", i)
				}
			}
		})
	})
	return root
}
