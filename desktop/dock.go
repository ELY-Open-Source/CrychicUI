package desktop

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DockApp is a dock app; a positive Badge shows.
type DockApp struct {
	Name, Icon string
	Running    bool
	Badge      int
}

// Dock magnifies under the pointer; clicks report; cue Tap.
func Dock(c *ui.Context, apps []DockApp, opened *string) ui.Element {
	names := make([]string, len(apps))
	for i, a := range apps {
		if a.Icon == "" || a.Badge < 0 {
			panic(fmt.Sprintf("desktop: Dock app %q needs an icon and a badge of 0 or more", a.Name))
		}
		names[i] = a.Name
	}
	unique("Dock", names)
	if len(apps) == 0 || opened == nil {
		panic("desktop: Dock needs apps and an opened name")
	}
	t := c.Theme()
	root := ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.End).Padding(t.Space(1.5)).Radius(t.Space(4)).
		Background(t.SurfaceHover).Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.12)).Role(ui.RoleToolbar)
	kit.Arrive(root, ui.Pose{Y: t.Space(4), Fade: 1}, 0)
	hover := ui.Local(root, "hover", func() int { return -1 })
	root.Children(func() {
		now := -1
		for i, a := range apps {
			c.Key(strconv.Quote(a.Name))
			b := ui.ButtonBase(c).Label(a.Name).Column().Size(t.Space(11), t.Space(13)).AlignItems(ui.Center).Justify(ui.Start).TransformOrigin(0.5, 1)
			if b.Hovered() {
				now = i
			}
			target := float32(1)
			if *hover >= 0 {
				switch d := max(i-*hover, *hover-i); d {
				case 0:
					target = 1.35
				case 1:
					target = 1.15
				}
			}
			scale := b.AnimateTo("scale", target, ui.Snappy)
			bounced := ui.Local(b, "bounced", func() time.Time { return time.Time{} })
			up := c.Now().Sub(*bounced) < 180*time.Millisecond
			if up {
				c.After(180*time.Millisecond - c.Now().Sub(*bounced))
			}
			lift := b.AnimateTo("lift", -t.Space(3)*unit(up), ui.Bouncy)
			b.Scale(scale, scale).Translate(0, lift)
			b.Children(func() {
				ui.Box(c).Size(t.Space(11), t.Space(11)).Children(func() {
					tile(c, a.Icon, t.Accent, t.Space(11)).Background(t.Background).Border(1, t.Border).Radius(t.Radius * 2)
					if a.Badge > 0 {
						ui.Box(c).Absolute().Right(-t.Space(1.5)).Top(-t.Space(1.5)).Size(t.Space(5.5), t.Space(4)).Center().Radius(999).Background(t.Danger).
							Children(func() {
								ui.Text(c, strconv.Itoa(min(a.Badge, 99))).FontSize(t.Rem(0.6)).Bold().TextColor(t.AccentText).Shrink(0).SingleLine().ContentTransition(ui.ContentTransition{})
							})
					}
				})
				ui.Box(c).Size(t.Space(1), t.Space(1)).Margin(t.Space(0.75), 0, 0, 0).Radius(999).Background(t.Text.Alpha(0.6 * unit(a.Running)))
			})
			if b.Clicked() {
				*opened = a.Name
				*bounced = c.Now()
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("Dock", "opened", "index", i)
			}
		}
		if now != *hover {
			*hover = now
			c.After(time.Millisecond)
		}
	})
	return root
}
