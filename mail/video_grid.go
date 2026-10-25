package mail

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// VideoGrid shows participant tiles and speaker rings.
func VideoGrid(c *ui.Context, participants []Participant) ui.Element {
	if len(participants) == 0 {
		panic("mail: VideoGrid needs participants")
	}
	ids := make([]string, len(participants))
	for i, p := range participants {
		if p.Name == "" {
			panic(fmt.Sprintf("mail: VideoGrid participant %q needs a name", p.ID))
		}
		ids[i] = p.ID
	}
	unique("VideoGrid", ids)
	t := c.Theme()
	cols := 3
	switch {
	case len(participants) == 1:
		cols = 1
	case len(participants) <= 4:
		cols = 2
	}
	root := ui.Grid(c).Columns(cols).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, p := range participants {
			c.Key(p.ID)
			tile := ui.Box(c).Height(t.Space(30)).Radius(t.Radius).Background(t.Text.Mix(t.Accent, 0.12)).Center()
			tile.Border(t.Space(0.5), tile.AnimateColor("speaker", t.Accent.Alpha(unit(p.Speaking)), ui.Snappy))
			room := tile.LayoutBounds().W
			kit.Settle(c, tile, []float32{room}, func() []float32 { return []float32{tile.LayoutBounds().W} })
			// Insets, pill padding, icon and gap around the name.
			text := min(t.Space(30), max(0, room-t.Space(10)-12))
			kit.Stagger(c, root, tile, i).Children(func() {
				avatar(c, p.Name, 14)
				pill := ui.Row(c).Absolute().Left(t.Space(2)).Bottom(t.Space(2)).Gap(t.Space(1)).AlignItems(ui.Center).
					Padding(t.Space(0.5), t.Space(2)).Radius(999).Background(t.Text.Alpha(0.55)).TextColor(t.Background)
				pill.Children(func() {
					switch {
					case p.Muted:
						glyph(c, "mic-off", 12).TextColor(t.Danger.Mix(t.Background, 0.3))
					case p.Speaking:
						glyph(c, "audio-lines", 12)
					}
					ui.Text(c, p.Name).FontSize(t.Rem(0.72)).Bold().Shrink(0).MaxWidth(text).SingleLine()
				})
			})
		}
	})
	return root
}
