package desktop

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SwitchWindow is an open window: its app, title and icon.
type SwitchWindow struct{ App, Title, Icon string }

// WindowSwitcherOptions word the next button.
type WindowSwitcherOptions struct{ Next string }

// WindowSwitcher is an alt-tab strip; next cycles; cue Tick.
func WindowSwitcher(c *ui.Context, windows []SwitchWindow, current *int, o WindowSwitcherOptions) (root, next ui.Element) {
	def(&o.Next, "Next window")
	titles := make([]string, len(windows))
	for i, w := range windows {
		if w.Icon == "" || w.App == "" {
			panic(fmt.Sprintf("desktop: WindowSwitcher window %q needs an app and an icon", w.Title))
		}
		titles[i] = w.Title
	}
	unique("WindowSwitcher", titles)
	if len(windows) == 0 || current == nil || *current < 0 || *current >= len(windows) {
		panic(fmt.Sprintf("desktop: WindowSwitcher needs windows and a current one, got %d", len(windows)))
	}
	t := c.Theme()
	pick := func(i int) {
		*current = i
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("WindowSwitcher", "current", "index", i)
	}
	root = ui.Column(c).Gap(t.Space(3)).Padding(t.Space(3)).Radius(t.Radius*2).Background(t.SurfaceHover.Alpha(0.96)).
		Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).AlignItems(ui.Center).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Scale: 0.95, Fade: 1}, 0)
	root.Children(func() {
		ui.ScrollHorizontal(c).MaxWidth(t.Space(120)).Children(func() {
			ui.Row(c).Gap(t.Space(2)).Children(func() {
				for i, w := range windows {
					c.Key(strconv.Quote(w.Title))
					on := i == *current
					b := ui.ButtonBase(c).Label(w.Title).Checked(on).Padding(t.Space(2)).Radius(t.Radius*1.5).Border(2, t.Accent.Alpha(unit(on)))
					b.Background(b.AnimateColor("bg", t.Border.Alpha(0.5*unit(b.Hovered())).Mix(t.Accent.Alpha(0.12), unit(on)), ui.Snappy))
					b.Children(func() {
						ui.Column(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							thumb := ui.Box(c).Size(t.Space(22), t.Space(15)).Radius(t.Radius).Background(t.Background).Border(1, t.Border).Center()
							thumb.Draw(func(p *ui.Painter, r ui.Rect) {
								p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W, H: t.Space(2.5)}, t.SurfaceHover, 0)
								for k := range 3 {
									p.Fill(ui.Rect{X: r.X + t.Space(2), Y: r.Y + t.Space(4) + float32(k)*t.Space(2), W: r.W * (0.7 - 0.15*float32(k)), H: t.Space(0.75)}, t.Border, 999)
								}
							})
							thumb.Children(func() { tile(c, w.Icon, t.Accent, t.Space(8)).Background(t.Background).Border(1, t.Border) })
							ui.Text(c, w.App).FontSize(t.Rem(0.72)).Bold().MaxWidth(t.Space(22)).SingleLine()
						})
					})
					if b.Clicked() && !on {
						pick(i)
					}
				}
			})
		})
		c.Key(strconv.Quote(windows[*current].Title))
		kit.Arrive(kit.Muted(c, windows[*current].Title).FontSize(t.Rem(0.78)).MaxWidth(t.Space(100)).SingleLine(), ui.Pose{Y: t.Space(1), Fade: 1}, 0)
		c.Key("next")
		next = kit.Press(ui.ButtonBase(c).Label(o.Next).Padding(t.Space(1), t.Space(3)).Radius(999).Disabled(len(windows) < 2))
		next.Background(next.AnimateColor("bg", t.Background.Mix(t.Border, unit(next.Hovered())), ui.Snappy)).Border(1, t.Border)
		next.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, o.Next).FontSize(t.Rem(0.74)).Bold().Shrink(0).MaxWidth(t.Space(40)).SingleLine()
				glyph(c, "arrow-right", 12)
			})
		})
		if next.Clicked() && len(windows) > 1 {
			pick((*current + 1) % len(windows))
		}
	})
	return root, next
}
