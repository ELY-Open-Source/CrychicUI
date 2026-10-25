package social

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Story is one user's story ring.
type Story struct {
	ID   string
	User User
	Seen bool
}

// StoryRing opens unseen story avatars; cue Open.
func StoryRing(c *ui.Context, stories []Story, opened *string) ui.Element {
	if len(stories) == 0 || opened == nil {
		panic("social: StoryRing needs stories and state")
	}
	ids := make([]string, len(stories))
	for i, story := range stories {
		validUser("StoryRing", story.User)
		ids[i] = story.ID
	}
	unique("StoryRing", ids)
	t := c.Theme()
	root := ui.ScrollHorizontal(c).FillWidth().FocusGroup(ui.Horizontal).Role(ui.RoleGroup)
	root.Children(func() {
		strip := ui.Row(c).Gap(t.Space(3)).Padding(t.Space(1))
		strip.Children(func() {
			for i, story := range stories {
				c.Key(story.ID)
				b := kit.Press(ui.ButtonBase(c).Label(story.User.Name).Column().Gap(t.Space(1.5)).AlignItems(ui.Center).
					Padding(t.Space(1)).Radius(t.Radius).Shrink(0).Checked(*opened == story.ID))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
				kit.Stagger(c, strip, b, i).Children(func() {
					frame := ui.Box(c).Size(t.Space(17), t.Space(17)).Shrink(0).Center()
					seen := frame.AnimateTo("seen", unit(story.Seen), ui.Smooth)
					frame.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
						cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, r.W/2-t.Space(0.5)
						p.StrokePath(ring(cx, cy, rad), 1, t.Border)
						if seen < 1 {
							p.StrokePathConicGradient(ring(cx, cy, rad), t.Space(0.75), ui.ConicGradient{Stops: []ui.GradientStop{
								{Color: t.Accent.Alpha(1 - seen), At: 0}, {Color: t.Danger.Alpha(1 - seen), At: 0.5}, {Color: t.Warning.Alpha(1 - seen), At: 0.75},
								{Color: t.Accent.Alpha(1 - seen), At: 1}}})
						}
					}))
					frame.Children(func() { avatar(c, story.User, t.Space(14)) })
					ui.Text(c, story.User.Name).FontSize(t.Rem(0.72)).MaxWidth(t.Space(17)).SingleLine().
						TextColor(t.Text.Mix(t.TextMuted, unit(story.Seen)))
				})
				if b.Clicked() {
					*opened = story.ID
					root.SetChanged()
					sound.Play(sound.Open)
					kit.Log("StoryRing", "opened", "index", i)
				}
			}
		})
	})
	return root
}
