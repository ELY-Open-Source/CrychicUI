package game

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// QuestStep is one objective of a quest.
type QuestStep struct {
	Text string
	Done bool
}

// Quest is one quest; Text and Reward are optional.
type Quest struct {
	Title, Text, Reward string
	Steps               []QuestStep
}

// QuestsOptions word the log; Track and Untrack take a title.
type QuestsOptions struct{ Title, Track, Untrack, Done, Steps string }

// Quests is a quest log; tracking plays Select.
func Quests(c *ui.Context, quests []Quest, tracked *string, o QuestsOptions) ui.Element {
	def(&o.Title, "Quests")
	if tracked == nil || len(quests) == 0 {
		panic(fmt.Sprintf("game: Quests %q needs quests and a tracked title", o.Title))
	}
	titles := make([]string, len(quests))
	for i, q := range quests {
		if len(q.Steps) == 0 {
			panic(fmt.Sprintf("game: Quests %q needs steps", q.Title))
		}
		titles[i] = q.Title
	}
	unique("Quests", titles)
	def(&o.Track, "Track %s")
	def(&o.Untrack, "Stop tracking %s")
	def(&o.Done, "Complete")
	def(&o.Steps, "%d of %d")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "scroll", o.Title, func() { kit.Chip(c, fmt.Sprint(len(quests)), t.Accent) })
		list := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3))
		list.Children(func() {
			for i, q := range quests {
				done := 0
				for _, s := range q.Steps {
					if s.Done {
						done++
					}
				}
				whole := done == len(q.Steps)
				on := *tracked == q.Title
				c.Key(q.Title)
				box := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius)
				box.Border(1, box.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
				box.Background(box.AnimateColor("bg", t.Background.Mix(t.Accent.Alpha(0.05), unit(on)), ui.Snappy))
				kit.Stagger(c, list, box, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
						ink := t.Accent
						if whole {
							ink = t.Success
						}
						tile(c, "scroll", ink, t.Space(8))
						ui.Text(c, q.Title).FontSize(t.Rem(0.85)).Bold().Grow(1).MinWidth(0).SingleLine()
						if whole {
							kit.Chip(c, o.Done, t.Success)
						} else {
							kit.Chip(c, fmt.Sprintf(o.Steps, done, len(q.Steps)), t.Accent)
						}
						label := fmt.Sprintf(o.Track, q.Title)
						if on {
							label = fmt.Sprintf(o.Untrack, q.Title)
						}
						b := kit.Press(ui.ButtonBase(c).Label(label).Size(t.Space(8), t.Space(8)).Center().Radius(999).Shrink(0).Checked(on))
						b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent.Alpha(0.15), unit(on)), ui.Snappy))
						b.Children(func() {
							g := glyph(c, "map-pin", 16)
							g.TextColor(g.AnimateColor("ink", t.TextMuted.Mix(t.Accent, unit(on)), ui.Snappy))
						})
						if b.Clicked() {
							if on {
								*tracked = ""
							} else {
								*tracked = q.Title
							}
							root.SetChanged()
							sound.Play(sound.Select)
							kit.Log("Quests", "tracked", "index", i, "on", !on)
						}
					})
					if q.Text != "" {
						kit.Muted(c, q.Text).FontSize(t.Rem(0.75)).FillWidth().Wrap()
					}
					ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
						for _, s := range q.Steps {
							ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
								c.Key(s.Done)
								if s.Done {
									kit.Arrive(glyph(c, "circle-check", 14).TextColor(t.Success), ui.Pose{Scale: 0.5, Fade: 1}, 0)
								} else {
									glyph(c, "circle", 14).TextColor(t.TextMuted)
								}
								text := ui.Text(c, s.Text).FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).Wrap()
								if s.Done {
									text.TextColor(t.TextMuted)
								}
							})
						}
					})
					meter(c, float32(done)/float32(len(q.Steps)), t.Accent.Mix(t.Success, unit(whole))).Height(t.Space(1.25))
					if q.Reward != "" {
						ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							glyph(c, "gift", 14).TextColor(t.Warning)
							ui.Text(c, q.Reward).FontSize(t.Rem(0.72)).Bold().TextColor(t.Warning).MinWidth(0).Wrap()
						})
					}
				})
			}
		})
	})
	return root
}
