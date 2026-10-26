package desktop

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Network is a Wi-Fi network; Bars is 0 to 3.
type Network struct {
	Name   string
	Bars   int
	Secure bool
}

// WifiItemOptions word the item, its menu and the empty state.
type WifiItemOptions struct{ Label, Title, Empty string }

// WifiItem is a Wi-Fi menu; opening reports; cue Select.
func WifiItem(c *ui.Context, networks []Network, current *string, o WifiItemOptions) ui.Element {
	def(&o.Label, "Wi-Fi")
	def(&o.Title, "Networks")
	def(&o.Empty, "No networks found")
	names := make([]string, len(networks))
	for i, n := range networks {
		if n.Bars < 0 || n.Bars > 3 {
			panic(fmt.Sprintf("desktop: WifiItem network %q needs 0 to 3 bars", n.Name))
		}
		names[i] = n.Name
	}
	unique("WifiItem", names)
	if current == nil || (*current != "" && !slices.Contains(names, *current)) {
		panic(fmt.Sprintf("desktop: WifiItem %q needs a current network, or none", o.Label))
	}
	t := c.Theme()
	signal := func(bars int) string { return [4]string{"wifi-zero", "wifi-low", "wifi-high", "wifi"}[bars] }
	at := slices.Index(names, *current)
	root := ui.Column(c).AlignItems(ui.Start).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: -t.Space(1), Fade: 1}, 0)
	open := ui.Local(root, "open", func() bool { return false })
	root.Children(func() {
		trigger := kit.Press(ui.ButtonBase(c).Label(o.Label).Role(ui.RoleMenuButton).Expanded(*open).Padding(t.Space(1), t.Space(1.5)).Radius(t.Radius))
		trigger.Background(trigger.AnimateColor("bg", t.Border.Alpha(0.6*max(unit(trigger.Hovered()), unit(*open))), ui.Snappy))
		trigger.Children(func() {
			mark := "wifi-off"
			if at >= 0 {
				mark = signal(networks[at].Bars)
			}
			glyph(c, mark, 16).ContentTransition(ui.ContentTransition{})
		})
		if trigger.Clicked() {
			*open = !*open
			if *open {
				root.SetChanged()
				sound.Play(sound.Open)
				kit.Log("WifiItem", "opened")
			}
		}
		ui.PopoverBase(c, trigger, open, func(panel ui.Element) {
			panel.Margin(t.Space(1), 0, 0, 0).MinWidth(t.Space(60)).Padding(t.Space(1)).Gap(t.Space(0.5)).Radius(t.Radius).
				Background(t.Background).Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).TransformOrigin(0, 0)
			kit.Arrive(panel, ui.Pose{Scale: 0.95, Fade: 1}, 0)
			kit.Muted(c, o.Title).FontSize(t.Rem(0.72)).Bold().Padding(t.Space(1), t.Space(2))
			if len(networks) == 0 {
				c.Key("empty")
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Children(func() {
					glyph(c, "wifi-off", 14).TextColor(t.TextMuted)
					kit.Muted(c, o.Empty).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).Wrap()
				})
			}
			ui.Scroll(c).FillWidth().MaxHeight(t.Space(72)).Gap(t.Space(0.5)).Role(ui.RoleMenu).Children(func() {
				for i, n := range networks {
					c.Key(strconv.Quote(n.Name))
					on := i == at
					item := kit.Press(ui.ButtonBase(c).Label(n.Name).Role(ui.RoleMenuItem).Checked(on).FillWidth().Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius))
					item.Background(item.AnimateColor("bg", t.SurfaceHover.Alpha(unit(item.Hovered())), ui.Snappy))
					item.Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							dot := ui.Box(c).Size(t.Space(6), t.Space(6)).Shrink(0).Center().Radius(999)
							dot.Background(t.SurfaceHover.Mix(t.Accent, unit(on)))
							dot.Children(func() { glyph(c, signal(n.Bars), 14).TextColor(t.Text.Mix(t.AccentText, unit(on))) })
							ui.Text(c, n.Name).FontSize(t.Rem(0.82)).Grow(1).MinWidth(0).SingleLine()
							if n.Secure {
								glyph(c, "lock", 12).TextColor(t.TextMuted)
							}
						})
					})
					if item.Clicked() {
						*open = false
						if !on {
							*current = n.Name
							root.SetChanged()
							sound.Play(sound.Select)
							kit.Log("WifiItem", "joined", "index", i)
						}
					}
				}
			})
		})
	})
	return root
}
