package game

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Seat is one player in a lobby or ready check.
type Seat struct {
	Name        string
	Ready, Host bool
}

// LobbyOptions word the lobby; Count takes ready and seats.
type LobbyOptions struct{ Title, Count, Waiting, Host, Ready, NotReady string }

// Lobby lists players with ready marks and empty seats.
func Lobby(c *ui.Context, players []Seat, seats int, o LobbyOptions) ui.Element {
	def(&o.Title, "Lobby")
	if len(players) == 0 || seats < len(players) || seats > 16 {
		panic(fmt.Sprintf("game: Lobby %q needs players and at most 16 seats, enough for them", o.Title))
	}
	names := make([]string, len(players))
	ready, hosts := 0, 0
	for i, p := range players {
		names[i] = p.Name
		if p.Ready {
			ready++
		}
		if p.Host {
			hosts++
		}
	}
	unique("Lobby", names)
	if hosts > 1 {
		panic(fmt.Sprintf("game: Lobby %q has %d hosts", o.Title, hosts))
	}
	def(&o.Count, "%d of %d ready")
	def(&o.Waiting, "Waiting for a player")
	def(&o.Host, "Host")
	def(&o.Ready, "Ready")
	def(&o.NotReady, "Not ready")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		tone := t.Warning
		if ready == seats {
			tone = t.Success
		}
		kit.CardHead(c, "users", o.Title, func() { kit.Chip(c, fmt.Sprintf(o.Count, ready, seats), tone) })
		list := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2))
		list.Children(func() {
			for i := range seats {
				if i >= len(players) {
					c.Key(emptySeat(i))
					seat := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).
						Border(1, t.Border).BorderStyle(ui.BorderDashed)
					kit.Stagger(c, list, seat, i).Children(func() {
						d := ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999).Background(t.SurfaceHover)
						d.Opacity(0.5 + 0.5*float32(d.Keyframes("wait", []ui.Keyframe{{Value: 0}, {At: 900 * time.Millisecond, Value: 1, Ease: ui.EaseInOut}},
							ui.KeyframeOptions{Repeat: 1, Alternate: true})))
						d.Children(func() { glyph(c, "user-plus", 14).TextColor(t.TextMuted) })
						kit.Muted(c, o.Waiting).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).SingleLine()
					})
					continue
				}
				p := players[i]
				// Quoted, a name cannot read as a seat key.
				c.Key(strconv.Quote(p.Name))
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())).Mix(t.Success.Alpha(0.06), unit(p.Ready)), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					display.Avatar(c, p.Name, nil).Size(t.Space(8), t.Space(8))
					ui.Text(c, p.Name).FontSize(t.Rem(0.85)).Bold().MinWidth(0).SingleLine()
					if p.Host {
						ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Shrink(0).Children(func() {
							glyph(c, "crown", 14).TextColor(t.Warning)
							ui.Text(c, o.Host).FontSize(t.Rem(0.68)).Bold().TextColor(t.Warning).Shrink(0)
						})
					}
					ui.Box(c).Grow(1)
					readyMark(c, p.Ready, o.Ready, o.NotReady)
				})
			}
		})
	})
	return root
}

// readyMark is a ready or waiting pill that swaps in.
func readyMark(c *ui.Context, ready bool, yes, no string) {
	t := c.Theme()
	c.Key(ready)
	ink, mark, text := t.TextMuted, "hourglass", no
	if ready {
		ink, mark, text = t.Success, "circle-check", yes
	}
	pill := ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(2)).Radius(999).Shrink(0).
		Background(ink.Alpha(0.12)).TextColor(ink)
	kit.Arrive(pill, ui.Pose{Scale: 0.7, Fade: 1}, 0).Children(func() {
		glyph(c, mark, 14)
		ui.Text(c, text).FontSize(t.Rem(0.72)).Bold().Shrink(0)
	})
}

type emptySeat int
