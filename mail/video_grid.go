package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// VideoGrid shows participant tiles and speaker rings.
func VideoGrid(c *ui.Context, participants []Participant) ui.Element {
	if len(participants) == 0 {
		panic("mail: VideoGrid needs participants")
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, person := range participants {
			if person.ID == "" || person.Name == "" {
				panic("mail: invalid participant")
			}
			c.Key(person.ID)
			tile := ui.Box(c).Height(t.Space(32)).Radius(t.Radius).Background(t.Surface).Center()
			ring := tile.AnimateColor("speaker", t.Border.Mix(t.Accent, unit(person.Speaking)), ui.Snappy)
			tile.Border(1, ring)
			kit.Arrive(tile, ui.Pose{Scale: 0.96, Fade: 1}, 0)
			tile.Children(func() {
				ui.Box(c).Size(t.Space(10), t.Space(10)).Radius(999).Background(t.Accent.Alpha(0.15)).Center().Children(func() {
					ui.Text(c, initial(person.Name)).TextColor(t.Accent)
				})
				ui.Row(c).Absolute().Left(t.Space(2)).Bottom(t.Space(2)).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
					ui.Text(c, person.Name)
					if person.Muted {
						ui.Icon(c, icon.Must("mic-off")).Size(16, 16).TextColor(t.Danger)
					}
				})
			})
		}
	})
	return root
}
