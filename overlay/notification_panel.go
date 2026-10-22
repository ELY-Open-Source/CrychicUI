package overlay

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Notice is one notification.
type Notice struct {
	Title, Text string
	Unread      bool
}

// NoticeGroup is a titled group of notices.
type NoticeGroup struct {
	Title   string
	Notices []Notice
}

// NotificationPanelOptions tune a NotificationPanel.
type NotificationPanelOptions struct {
	// Empty shows when there are no notices.
	Empty string
}

// NotificationPanel is a bell opening grouped notices.
func NotificationPanel(c *ui.Context, label string, groups []NoticeGroup, open *bool, o NotificationPanelOptions) ui.Element {
	if label == "" || open == nil {
		panic(fmt.Sprintf("overlay: NotificationPanel %q needs open", label))
	}
	unread, total := 0, 0
	for _, g := range groups {
		if g.Title == "" {
			panic(fmt.Sprintf("overlay: NotificationPanel %q groups need titles", label))
		}
		for _, n := range g.Notices {
			total++
			if n.Unread {
				unread++
			}
		}
	}
	def(&o.Empty, "You're all caught up")
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	root.Children(func() {
		bell := kit.Press(ui.ButtonBase(c).Label(label).Tooltip(label).Size(t.Space(9), t.Space(9)).Center().Radius(t.Radius)).Expanded(*open)
		bell.Background(bell.AnimateColor("bg", t.SurfaceHover.Alpha(unit(*open || bell.Hovered())), ui.Snappy)).Children(func() {
			ui.Icon(c, icon.Must("bell")).Size(16, 16)
			if unread > 0 {
				c.Key("count")
				badge := ui.Box(c).Absolute().Top(t.Space(0.5)).End(t.Space(0.5)).Padding(0, t.Space(1)).
					Radius(999).Background(t.Danger).TextColor(t.AccentText).FontSize(t.Rem(0.65))
				kit.Arrive(badge, ui.Pose{Scale: 0.4, Fade: 1}, 0).Children(func() {
					kit.Numeric(c, strconv.Itoa(unread)).Bold().ContentTransition(ui.ContentTransition{Numeric: true})
				})
			}
		})
		if bell.Clicked() {
			*open = !*open
			root.SetChanged()
		}
		if ui.PopoverBase(c, bell, open, func(panel ui.Element) {
			panel.Margin(t.Space(1), 0, 0, 0).Width(t.Space(84)).Padding(t.Space(2)).Radius(t.Radius).Background(t.Background).
				Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).TransformOrigin(0.5, 0)
			kit.Arrive(panel, ui.Pose{Scale: 0.95, Fade: 1}, 0)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				if total == 0 {
					kit.Muted(c, o.Empty).Padding(t.Space(3))
					return
				}
				for g, grp := range groups {
					if len(grp.Notices) == 0 {
						continue
					}
					c.Key(g)
					kit.Muted(c, grp.Title).Bold().FontSize(t.Rem(0.75)).Padding(t.Space(1), t.Space(2))
					for i, n := range grp.Notices {
						c.Key(fmt.Sprint(g, "/", i))
						ui.Row(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(1.5), t.Space(2)).AlignItems(ui.Start).Radius(t.Radius).Children(func() {
							dot := ui.Box(c).Size(t.Space(2), t.Space(2)).Margin(t.Space(1.5), 0, 0, 0).Shrink(0).Radius(999)
							dot.Background(dot.AnimateColor("dot", t.Accent.Alpha(unit(n.Unread)), ui.Snappy))
							ui.Column(c).Grow(1).Gap(t.Space(0.5)).Children(func() {
								ui.Text(c, n.Title).Bold()
								kit.Muted(c, n.Text).FontSize(t.Rem(0.85))
							})
						})
					}
				}
			})
		}).Dismissed() {
			root.SetChanged()
		}
	})
	if opened, _ := moved(root, *open, "NotificationPanel"); opened {
		sound.Play(sound.Open)
	}
	return root
}
