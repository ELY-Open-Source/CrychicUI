package code

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// GraphCommit is one commit in a graph, newest first.
type GraphCommit struct {
	Hash, Message, Author string
	Lane                  int
	// Parents are hashes; Refs are its branch and tag names.
	Parents, Refs []string
}

// CommitGraph draws lanes, merges and dots beside commits.
func CommitGraph(c *ui.Context, label string, commits []GraphCommit) ui.Element {
	if label == "" || len(commits) == 0 {
		panic(fmt.Sprintf("code: CommitGraph %q needs a label and commits", label))
	}
	row := map[string]int{}
	lanes := 0
	for i, k := range commits {
		if k.Hash == "" || k.Message == "" || k.Lane < 0 || k.Lane > 7 || row[k.Hash] > 0 {
			panic(fmt.Sprintf("code: CommitGraph %q needs unique hashes, messages and lanes 0..7", label))
		}
		row[k.Hash] = i + 1
		lanes = max(lanes, k.Lane+1)
	}
	for i, k := range commits {
		for _, p := range k.Parents {
			if j := row[p]; j > 0 && j-1 <= i {
				panic(fmt.Sprintf("code: CommitGraph %q lists %s before its parent", label, k.Hash))
			}
		}
	}
	t := c.Theme()
	inks := []ui.Color{t.Accent, t.Success, t.Warning, t.Danger}
	lane := t.Space(4)
	h := t.Space(9)
	root := ui.Column(c).FillWidth().Padding(t.Space(1), 0).Role(ui.RoleList).Label(label)
	g := root.Keyframes("crychic.draw", []ui.Keyframe{{Value: 0}, {At: 700 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	root.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
		x := func(l int) float32 { return r.X + t.Space(3) + float32(l)*lane }
		y := func(i int) float32 { return r.Y + t.Space(1) + float32(i)*h + h/2 }
		bottom := r.Y + r.H
		for i, k := range commits {
			for _, par := range k.Parents {
				j, ok := row[par]
				to, pl := bottom, k.Lane
				if ok {
					to, pl = y(j-1), commits[j-1].Lane
				}
				var path ui.Path
				path.MoveTo(x(k.Lane), y(i))
				if pl != k.Lane {
					path.CubeTo(x(k.Lane), y(i)+h*0.6, x(pl), y(i)+h*0.4, x(pl), y(i)+h)
				}
				path.LineTo(x(pl), to)
				p.StrokePath(path.Trim(0, g), 2, inks[pl%len(inks)].Alpha(0.7))
			}
		}
		for i, k := range commits {
			ink := inks[k.Lane%len(inks)]
			rad := t.Space(1.25)
			if len(k.Parents) > 1 {
				rad = t.Space(1.75)
			}
			p.Fill(ui.Rect{X: x(k.Lane) - rad - 2, Y: y(i) - rad - 2, W: 2*rad + 4, H: 2*rad + 4}, t.Background, 999)
			p.Fill(ui.Rect{X: x(k.Lane) - rad, Y: y(i) - rad, W: 2 * rad, H: 2 * rad}, ink, 999)
		}
	}))
	return root.Children(func() {
		for i, k := range commits {
			c.Key(k.Hash)
			item := ui.Row(c).FillWidth().Height(h).Gap(t.Space(2)).AlignItems(ui.Center).
				Padding(0, t.Space(3), 0, t.Space(3)+float32(lanes)*lane).Role(ui.RoleListItem).Label(k.Message)
			stagger(c, root, item, i).Children(func() {
				ui.Text(c, short(k.Hash)).Font("monospace").FontSize(t.Rem(0.75)).TextColor(t.TextMuted).Shrink(0)
				for _, ref := range k.Refs {
					c.Key("ref." + ref)
					chip(c, ref, inks[k.Lane%len(inks)])
				}
				ui.Text(c, k.Message).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
				if k.Author != "" {
					kit.Muted(c, k.Author).FontSize(t.Rem(0.75)).Shrink(0)
				}
			})
		}
	})
}
