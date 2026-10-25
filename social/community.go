package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CommunityData is one community; Icon, Cover, Faces are optional.
type CommunityData struct {
	ID, Name, Description, Icon string
	Members                     int
	Cover                       ui.ImageSource
	Faces                       []string
}

// CommunityOptions word the join pill; Members formats the grouped count.
type CommunityOptions struct{ Join, Joined, Members string }

// Community joins a community card; cue Toggle.
func Community(c *ui.Context, community CommunityData, joined *bool, o CommunityOptions) ui.Element {
	if community.ID == "" || community.Name == "" || community.Description == "" || community.Members < 0 || joined == nil {
		panic(fmt.Sprintf("social: Community %q needs id, description, members and state", community.Name))
	}
	def(&o.Join, "Join")
	def(&o.Joined, "Joined")
	def(&o.Members, "%s members")
	def(&community.Icon, "users")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		band := ui.Box(c).FillWidth().Height(t.Space(16)).Clip().Background(t.Accent.Alpha(0.18))
		band.Children(func() {
			if community.Cover != nil {
				ui.Image(c, community.Cover).Fill().Fit(ui.Cover)
			}
		})
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(0, t.Space(3.5), t.Space(3.5), t.Space(3.5)).Children(func() {
			tile := ui.Box(c).Size(t.Space(12), t.Space(12)).Margin(-t.Space(6), 0, 0, 0).Center().Radius(t.Radius*1.5).
				Background(t.Background).Border(t.Space(0.75), t.Background).Shadow(0, 1, 3, 0, t.Text.Alpha(0.12))
			kit.Arrive(tile, ui.Pose{Scale: 0.6, Fade: 1}, 0).Children(func() {
				ui.Box(c).Fill().Center().Radius(t.Radius).Background(t.Accent.Alpha(0.14)).TextColor(t.Accent).Children(func() {
					glyph(c, community.Icon, 22)
				})
			})
			ui.Text(c, community.Name).FontSize(t.Rem(1.05)).Bold().Wrap().FillWidth()
			kit.Muted(c, community.Description).FontSize(t.Rem(0.82)).Wrap().FillWidth()
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), 0, 0, 0).Children(func() {
				if len(community.Faces) > 0 {
					display.AvatarGroup(c, community.Faces, display.AvatarGroupOptions{})
				}
				kit.Numeric(c, fmt.Sprintf(o.Members, number(c, community.Members))).FontSize(t.Rem(0.78)).TextColor(t.TextMuted).
					Grow(1).MinWidth(0).SingleLine()
				text := o.Join
				if *joined {
					text = o.Joined
				}
				c.Key("join")
				b := follow(c, text, text, "plus", *joined)
				if b.Clicked() {
					*joined = !*joined
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("Community", "joined", *joined)
				}
			})
		})
	})
	return root
}
