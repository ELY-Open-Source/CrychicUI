package devtools

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DeployRecord is one past deploy.
type DeployRecord struct {
	ID, Version, Commit, Author, When string
	Current, Failed                   bool
}

// DeploysOptions word the panel; Label takes action and version.
type DeploysOptions struct{ Title, Current, Failed, Rollback, Label string }

// Deploys lists deploy history with rollback; cue Warning.
func Deploys(c *ui.Context, deploys []DeployRecord, rollback *string, o DeploysOptions) ui.Element {
	if len(deploys) == 0 || rollback == nil {
		panic("devtools: Deploys needs deploys and a rollback target")
	}
	ids := make([]string, len(deploys))
	current := 0
	for i, d := range deploys {
		if d.Version == "" || d.Commit == "" || d.Author == "" || d.When == "" || d.Current && d.Failed {
			panic(fmt.Sprintf("devtools: Deploys %q needs a version, commit, author and time, and cannot be current and failed", d.ID))
		}
		ids[i] = d.ID
		if d.Current {
			current++
		}
	}
	unique("Deploys", ids)
	if current != 1 {
		panic(fmt.Sprintf("devtools: Deploys needs one current deploy, has %d", current))
	}
	def(&o.Title, "Deploys")
	def(&o.Current, "Current")
	def(&o.Failed, "Failed")
	def(&o.Rollback, "Roll back")
	def(&o.Label, "%s to %s")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "history", o.Title, nil)
		list := ui.Column(c).FillWidth().Padding(t.Space(3), t.Space(4))
		list.Children(func() {
			for i, d := range deploys {
				ink := t.TextMuted
				switch {
				case d.Current:
					ink = t.Success
				case d.Failed:
					ink = t.Danger
				}
				c.Key(d.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Stretch)
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Column(c).Width(t.Space(3)).Shrink(0).AlignItems(ui.Center).Padding(t.Space(1.5), 0, 0, 0).Children(func() {
						dot := ui.Box(c).Size(t.Space(3), t.Space(3)).Shrink(0).Radius(999).Background(ink)
						if d.Current {
							dot.Border(t.Space(0.75), ink.Alpha(0.3))
						}
						if i < len(deploys)-1 {
							ui.Box(c).Width(1).Grow(1).MinHeight(t.Space(3)).Background(t.Border)
						}
					})
					ui.Row(c).Grow(1).MinWidth(0).Gap(t.Space(2)).AlignItems(ui.Start).Padding(0, 0, t.Space(3), 0).Children(func() {
						ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
							ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
								ui.Text(c, d.Version).FontSize(t.Rem(0.85)).Bold().Shrink(0)
								ui.Text(c, d.Commit).Font("monospace").FontSize(t.Rem(0.7)).TextColor(t.TextMuted).MinWidth(0).SingleLine()
							})
							ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
								kit.Muted(c, d.Author).FontSize(t.Rem(0.75)).MinWidth(0).SingleLine()
								dot(c)
								kit.Muted(c, d.When).FontSize(t.Rem(0.75)).MinWidth(0).SingleLine()
							})
						})
						switch {
						case d.Current:
							kit.Chip(c, o.Current, t.Success)
						case d.Failed:
							kit.Chip(c, o.Failed, t.Danger)
						default:
							b := outline(c, fmt.Sprintf(o.Label, o.Rollback, d.Version), o.Rollback, "undo-2", t.Warning)
							if b.Clicked() {
								*rollback = d.ID
								root.SetChanged()
								sound.Play(sound.Warning)
								kit.Log("Deploys", "rollback", "index", i)
							}
						}
					})
				})
			}
		})
	})
	return root
}
