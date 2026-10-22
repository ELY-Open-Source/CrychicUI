package choice

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Tile is a labeled Lucide icon tile.
type Tile struct{ Label, Icon string }

// OptionTiles is icon tiles choosing several; checks pop in.
func OptionTiles(c *ui.Context, label string, tiles []Tile, selected *[]string) ui.Element {
	labels := make([]string, len(tiles))
	for i, tl := range tiles {
		labels[i] = tl.Label
		icon.Must(tl.Icon)
	}
	distinct("OptionTiles", labels)
	subset("OptionTiles", labels, selected)
	if label == "" {
		panic("choice: OptionTiles needs a label")
	}
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(3)).Label(label).FocusGroup(ui.Horizontal)
	return root.Children(func() {
		for _, tl := range tiles {
			on := slices.Contains(*selected, tl.Label)
			shown := on
			c.Key(tl.Label)
			b := ui.ToggleBase(c, &on).Label(tl.Label).Column().Gap(t.Space(2)).Size(t.Space(24), t.Space(20)).Radius(t.Radius + 2)
			edge, fill, fg := t.Border, t.Background, t.Text
			if shown {
				edge, fill, fg = t.Accent, t.Accent.Alpha(0.06), t.Accent
			} else if b.Hovered() {
				edge = t.Border.Mix(t.Text, 0.25)
			}
			b.Background(b.AnimateColor("fill", fill, ui.Snappy)).Border(1, b.AnimateColor("edge", edge, ui.Snappy))
			kit.Press(b).Children(func() {
				ui.Icon(c, icon.Must(tl.Icon)).Size(16, 16).TextColor(b.AnimateColor("fg", fg, quick))
				ui.Text(c, tl.Label).FontSize(t.Rem(0.85))
				badge := ui.Box(c).Attach(ui.AnchorTopRight, ui.AnchorTopRight).Top(t.Space(1.5)).Right(t.Space(1.5)).
					Size(t.Space(5), t.Space(5)).Center().Radius(999).Background(t.Accent).TextColor(t.AccentText)
				k := badge.AnimateTo("pop", unit(shown), ui.Bouncy)
				badge.Scale(max(k, 0), max(k, 0)).Opacity(min(max(k, 0), 1)).Children(func() {
					ui.Icon(c, icon.Must("check")).Size(16, 16)
				})
			})
			if kit.Changed(b, sound.Toggle) {
				*selected = flip(*selected, tl.Label)
				root.SetChanged()
				kit.Log("OptionTiles", "toggled", "count", len(*selected))
			}
		}
	})
}
