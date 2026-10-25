package game

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Team is one side of a scoreboard.
type Team struct {
	Name  string
	Score int
}

// ScoreboardOptions word the board; Add takes a team name.
type ScoreboardOptions struct{ Period, Add, Point string }

// Scoreboard shows two scores that roll and bump; cue Pop.
func Scoreboard(c *ui.Context, teams *[2]Team, o ScoreboardOptions) ui.Element {
	if teams == nil {
		panic("game: Scoreboard needs two teams")
	}
	for _, tm := range teams {
		if strings.TrimSpace(tm.Name) == "" || tm.Score < 0 {
			panic(fmt.Sprintf("game: Scoreboard team %q needs a name and a score of at least 0", tm.Name))
		}
	}
	if teams[0].Name == teams[1].Name {
		panic(fmt.Sprintf("game: Scoreboard has two teams named %q", teams[0].Name))
	}
	def(&o.Add, "Add a point to %s")
	def(&o.Point, "+1")
	t := c.Theme()
	root := card(c).Row().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(5), t.Space(4)).Role(ui.RoleGroup)
	bumps := ui.Local(root, "bumps", func() [2]time.Time { return [2]time.Time{} })
	now := c.Now()
	root.Children(func() {
		side := func(i int) {
			tm := teams[i]
			lead := tm.Score > teams[1-i].Score
			c.Key(tm.Name)
			ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				display.Avatar(c, tm.Name, nil).Size(t.Space(10), t.Space(10))
				ui.Text(c, tm.Name).FontSize(t.Rem(0.85)).Bold().MaxWidth(t.Space(36)).SingleLine()
				fresh := now.Sub(bumps[i]) < 120*time.Millisecond
				if fresh {
					c.After(120*time.Millisecond - now.Sub(bumps[i]))
				}
				n := kit.Numeric(c, fmt.Sprint(tm.Score)).FontSize(t.Rem(2.6)).Bold().FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine().
					ContentTransition(ui.ContentTransition{Numeric: true})
				s := n.AnimateTo("bump", 1+0.25*unit(fresh), ui.Bouncy)
				n.Scale(s, s).TextColor(n.AnimateColor("ink", t.Text.Mix(t.Accent, unit(lead)), ui.Snappy))
				b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Add, tm.Name)).Padding(t.Space(1), t.Space(3.5)).
					Radius(999).Border(1, t.Border).Shrink(0).Disabled(tm.Score == math.MaxInt))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
				b.Children(func() { ui.Text(c, o.Point).FontSize(t.Rem(0.8)).Bold().Shrink(0) })
				if b.Clicked() && tm.Score < math.MaxInt {
					teams[i].Score++
					bumps[i] = now
					root.SetChanged()
					sound.Play(sound.Pop)
					kit.Log("Scoreboard", "scored", "side", i, "score", teams[i].Score)
				}
			})
		}
		side(0)
		ui.Column(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(t.Space(4), t.Space(0.75)).Radius(999).Background(t.TextMuted)
			if o.Period != "" {
				kit.Chip(c, o.Period, t.Accent)
			}
		})
		side(1)
	})
	return root
}
