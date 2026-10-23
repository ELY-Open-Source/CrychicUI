// Code components for the gallery.
package code

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	"time"
)

type GraphCommit struct {
	Hash, Message string
	Lane          int
}

// CommitGraph is lanes and dots of commits.

func CommitGraph(c *ui.Context, label string, commits []GraphCommit) ui.Element {
	if label == "" || len(commits) == 0 {
		panic("code: CommitGraph needs a label and commits")
	}
	lanes := 0
	for _, k := range commits {
		if k.Hash == "" || k.Message == "" || k.Lane < 0 || k.Lane > 7 {
			panic("code: CommitGraph needs a hash, message and lane 0..7 in every commit")
		}
		if k.Lane+1 > lanes {
			lanes = k.Lane + 1
		}
	}
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Role(ui.RoleGroup).Label(label)
	laneInks := []ui.Color{t.Accent, t.Success, t.Warning, t.Danger}
	root.Children(func() {
		for i, k := range commits {
			c.Key(k.Hash)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Height(t.Space(6))
			kit.Arrive(row, ui.Pose{X: -t.Space(2), Fade: 1}, time.Duration(i*30)*time.Millisecond)
			row.Children(func() {
				strip := ui.Box(c).Width(t.Space(float32(lanes) * 3)).Height(t.Space(6)).Shrink(0)
				strip.Draw(func(p *ui.Painter, r ui.Rect) {
					for l := 0; l < lanes; l++ {
						x := r.X + t.Space(float32(l)*3+1)
						p.Line(x, r.Y, x, r.Y+r.H, 1, laneInks[l%len(laneInks)].Alpha(0.35))
						if l == k.Lane {
							p.Fill(ui.Rect{X: x - t.Space(1), Y: r.Y + r.H/2 - t.Space(1), W: t.Space(2), H: t.Space(2)}, laneInks[l%len(laneInks)], t.Space(1))
						}
					}
				})
				ui.Text(c, short(k.Hash)).Font("monospace").FontSize(t.Rem(0.8)).TextColor(t.TextMuted).Shrink(0)
				ui.Text(c, k.Message).SingleLine()
			})
		}
	})
	return root
}

// BlameLine is one line of blame.
