package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CommunityData is one community card.
type CommunityData struct {
	ID, Name, Description string
	Members               int
}

// CommunityOptions are the community labels.
type CommunityOptions struct{ Join, Joined, Members string }

// Community joins a community card; cue Toggle.
func Community(c *ui.Context, community CommunityData, joined *bool, o CommunityOptions) ui.Element {
	if community.ID == "" || community.Name == "" || community.Description == "" || community.Members < 0 || joined == nil {
		panic("social: invalid community")
	}
	def(&o.Join, "Join")
	def(&o.Joined, "Joined")
	def(&o.Members, "%d members")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Text(c, community.Name).Bold()
		ui.Text(c, community.Description).Wrap()
		ui.Row(c).AlignItems(ui.Center).Children(func() {
			kit.Numeric(c, fmt.Sprintf(o.Members, community.Members)).TextColor(t.TextMuted).Grow(1)
			label := o.Join
			if *joined {
				label = o.Joined
			}
			b := button(c, label).Checked(*joined)
			b.Children(func() { ui.Text(c, label) })
			if b.Clicked() {
				*joined = !*joined
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("Community", "joined", *joined)
			}
		})
	})
	return root
}
