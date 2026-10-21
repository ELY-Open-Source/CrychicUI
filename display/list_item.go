package display

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ListItemOptions tune a ListItem.
type ListItemOptions struct {
	// Icon leads the row; Clickable gives it a tap cue.
	Icon      string
	Clickable bool
}

// ListItem is an icon, text and trailing slot row.
func ListItem(c *ui.Context, title, subtitle string, o ListItemOptions, trailing func()) ui.Element {
	t := c.Theme()
	var root ui.Element
	if o.Clickable {
		root = ui.ButtonBase(c).Label(title).FillWidth().Justify(ui.Start)
	} else {
		root = ui.Row(c).FillWidth()
	}
	root.Gap(t.Space(3)).Padding(t.Space(2), t.Space(3)).Radius(t.Radius)
	bg := t.Background.Alpha(0)
	if root.Hovered() {
		bg = t.SurfaceHover
	}
	root.Background(root.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
		if o.Icon != "" {
			ui.Icon(c, icon.Must(o.Icon)).Size(16, 16).TextColor(t.TextMuted)
		}
		ui.Column(c).Grow(1).Gap(t.Space(0.5)).Children(func() {
			ui.Text(c, title)
			if subtitle != "" {
				kit.Muted(c, subtitle).FontSize(t.Rem(0.85))
			}
		})
		if trailing != nil {
			ui.Row(c).Children(trailing)
		}
	})
	if o.Clickable {
		kit.Clicked(root, sound.Tap)
	}
	return root
}
