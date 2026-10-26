package desktop

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DesktopNote is one notification.
type DesktopNote struct {
	ID, App, Icon, Title, Text string
	Time                       time.Time
}

// NotificationCenterOptions word the panel; Time is a layout.
type NotificationCenterOptions struct{ Title, Clear, Empty, Time string }

// NotificationCenter lists notifications by app; clear cues Delete.
func NotificationCenter(c *ui.Context, notes []DesktopNote, o NotificationCenterOptions) (root, clear ui.Element) {
	def(&o.Title, "Notifications")
	def(&o.Clear, "Clear all")
	def(&o.Empty, "No notifications")
	def(&o.Time, "15:04")
	ids := make([]string, len(notes))
	var apps []string
	for i, n := range notes {
		if strings.TrimSpace(n.App) == "" || n.Icon == "" || strings.TrimSpace(n.Title) == "" {
			panic(fmt.Sprintf("desktop: NotificationCenter note %q needs an app, an icon and a title", n.ID))
		}
		ids[i] = n.ID
		if !slices.Contains(apps, n.App) {
			apps = append(apps, n.App)
		}
	}
	unique("NotificationCenter", ids)
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Radius(t.Radius*2).Background(t.SurfaceHover).Border(1, t.Border).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{X: t.Space(4), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Text(c, o.Title).FontSize(t.Rem(0.95)).Bold().Grow(1).MinWidth(0).SingleLine()
			clear = kit.Press(ui.ButtonBase(c).Label(o.Clear).Padding(t.Space(1), t.Space(2.5)).Radius(999).Disabled(len(notes) == 0))
			clear.Background(clear.AnimateColor("bg", t.Background.Mix(t.Border, unit(clear.Hovered())), ui.Snappy))
			clear.Children(func() {
				ui.Text(c, o.Clear).FontSize(t.Rem(0.72)).Bold().Shrink(0).MaxWidth(t.Space(30)).SingleLine()
			})
			if clear.Clicked() {
				sound.Play(sound.Delete)
				kit.Log("NotificationCenter", "cleared", "count", len(notes))
			}
		})
		if len(notes) == 0 {
			ui.Column(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(8), 0).Children(func() {
				kit.Arrive(glyph(c, "bell-off", 24).TextColor(t.TextMuted), ui.Pose{Scale: 0.7, Fade: 1}, 0)
				kit.Muted(c, o.Empty).FontSize(t.Rem(0.84)).FillWidth().TextAlign(ui.Center).Wrap()
			})
			return
		}
		groups := ui.Column(c).FillWidth().Gap(t.Space(3))
		groups.Children(func() {
			for gi, app := range apps {
				c.Key(strconv.Quote(app))
				var mine []DesktopNote
				for _, n := range notes {
					if n.App == app {
						mine = append(mine, n)
					}
				}
				g := ui.Column(c).FillWidth().Gap(t.Space(1.5))
				kit.Stagger(c, groups, g, gi).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(0, t.Space(1)).Children(func() {
						glyph(c, mine[0].Icon, 14).TextColor(t.TextMuted)
						kit.Muted(c, app).FontSize(t.Rem(0.72)).Bold().Grow(1).MinWidth(0).SingleLine()
						kit.Chip(c, strconv.Itoa(len(mine)), t.TextMuted)
					})
					for _, n := range mine {
						c.Key(strconv.Quote(n.ID))
						ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Start).Padding(t.Space(3)).Radius(t.Radius*1.5).Background(t.Background).
							Shadow(0, 1, 3, 0, t.Text.Alpha(0.06)).Children(func() {
							tile(c, n.Icon, t.Accent, t.Space(8))
							ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
								ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
									ui.Text(c, n.Title).FontSize(t.Rem(0.82)).Bold().Grow(1).MinWidth(0).SingleLine()
									if !n.Time.IsZero() {
										kit.Muted(c, n.Time.Format(o.Time)).FontSize(t.Rem(0.68)).Shrink(0).MaxWidth(t.Space(20)).SingleLine()
									}
								})
								if n.Text != "" {
									kit.Muted(c, n.Text).FontSize(t.Rem(0.76)).FillWidth().MaxLines(2).Wrap()
								}
							})
						})
					}
				})
			}
		})
	})
	return root, clear
}
