package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Link is one dependency: an issue key, title and status.
type Link struct {
	Key, Title string
	Status     TaskStatus
}

// DependenciesOptions word the card; Open counts unfinished blockers.
type DependenciesOptions struct{ Title, BlockedBy, Blocks, None, Open, Clear string }

// Dependencies shows blocked-by and blocks links; cue Open.
func Dependencies(c *ui.Context, blockedBy, blocks []Link, opened *string, o DependenciesOptions) ui.Element {
	if opened == nil {
		panic("project: Dependencies needs opened state")
	}
	def(&o.Title, "Dependencies")
	def(&o.BlockedBy, "Blocked by")
	def(&o.Blocks, "Blocks")
	def(&o.None, "None")
	def(&o.Open, "%d blocking")
	def(&o.Clear, "Unblocked")
	open := 0
	for _, group := range [][]Link{blockedBy, blocks} {
		keys := make([]string, len(group))
		for i, link := range group {
			if link.Title == "" || link.Status > Done {
				panic(fmt.Sprintf("project: Dependencies %q has an invalid link %q", o.Title, link.Key))
			}
			keys[i] = link.Key
		}
		unique("Dependencies", keys)
	}
	for _, link := range blockedBy {
		if link.Status != Done {
			open++
		}
	}
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "git-branch", o.Title, func() {
			if open > 0 {
				kit.Chip(c, fmt.Sprintf(o.Open, open), t.Danger)
			} else {
				kit.Chip(c, o.Clear, t.Success)
			}
		})
		body := ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(2.5), t.Space(2), t.Space(3), t.Space(2))
		body.Children(func() {
			groups := []struct {
				name, glyph string
				ink         ui.Color
				links       []Link
			}{{o.BlockedBy, "octagon-x", t.Danger, blockedBy}, {o.Blocks, "arrow-right", t.Accent, blocks}}
			for g, group := range groups {
				c.Key(g)
				ui.Column(c).FillWidth().Gap(t.Space(0.5)).Children(func() {
					ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(0, t.Space(1)).TextColor(group.ink).Children(func() {
						glyph(c, group.glyph, 14)
						ui.Text(c, group.name).FontSize(t.Rem(0.75)).Bold().Shrink(0)
						kit.Muted(c, fmt.Sprint(len(group.links))).FontSize(t.Rem(0.7))
					})
					if len(group.links) == 0 {
						kit.Muted(c, o.None).FontSize(t.Rem(0.8)).Padding(t.Space(1), t.Space(1))
					}
					list := ui.Column(c).FillWidth()
					list.Children(func() {
						for i, link := range group.links {
							c.Key(link.Key)
							b := kit.Press(ui.ButtonBase(c).Label(link.Key).Role(ui.RoleLink).FillWidth().AlignItems(ui.Stretch).
								Padding(t.Space(1.5), t.Space(1)).Radius(t.Radius))
							b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
							nudge := b.AnimateTo("nudge", t.Space(0.75)*unit(b.Hovered()), ui.Snappy)
							kit.Stagger(c, list, b, i).Children(func() {
								ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
									name, ink := statusLook(t, link.Status)
									glyph(c, name, 14).TextColor(ink)
									kit.Muted(c, link.Key).FontSize(t.Rem(0.7)).Bold().Shrink(0)
									text := ui.Text(c, link.Title).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).SingleLine()
									if link.Status == Done {
										text.Strikethrough().TextColor(t.TextMuted)
									}
									glyph(c, "chevron-right", 14).TextColor(t.TextMuted).Translate(nudge, 0)
								})
							})
							if b.Clicked() {
								*opened = link.Key
								root.SetChanged()
								sound.Play(sound.Open)
								kit.Log("Dependencies", "opened")
							}
						}
					})
				})
			}
		})
	})
	return root
}
