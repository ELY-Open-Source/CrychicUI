package chat

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MembersOptions name the two groups.
type MembersOptions struct {
	Online, Offline string
}

// Members lists online and offline groups; a tap picks one.
func Members(c *ui.Context, label string, online, offline []string, chosen *int, o MembersOptions) ui.Element {
	total := len(online) + len(offline)
	if label == "" || total == 0 || chosen == nil || *chosen < -1 || *chosen >= total {
		panic(fmt.Sprintf("chat: Members %q needs a label, members and a choice or -1", label))
	}
	for _, n := range online {
		if n == "" {
			panic(fmt.Sprintf("chat: Members %q names must not be empty", label))
		}
	}
	for _, n := range offline {
		if n == "" {
			panic(fmt.Sprintf("chat: Members %q names must not be empty", label))
		}
	}
	def(&o.Online, "Online")
	def(&o.Offline, "Offline")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(0.5)).FocusGroup(ui.Vertical).Role(ui.RoleGroup).Label(label)
	head := func(title string, n int) {
		kit.Muted(c, title+" · "+strconv.Itoa(n)).Bold().FontSize(t.Rem(0.75)).Padding(t.Space(1), t.Space(2))
	}
	row := func(i int, name string, live bool) {
		on := i == *chosen
		b := kit.Press(ui.ButtonBase(c).Label(name).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
			Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).Checked(on))
		bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
		if on {
			bg = t.Accent.Alpha(0.1)
		}
		b.Background(b.AnimateColor(fmt.Sprint("bg.", i), bg, ui.Snappy)).Children(func() {
			dot := t.TextMuted
			if live {
				dot = t.Success
			}
			ui.Box(c).Size(t.Space(2), t.Space(2)).Shrink(0).Radius(999).Background(dot)
			ui.Text(c, name).Grow(1).SingleLine()
		})
		if !live {
			b.Opacity(0.55)
		}
		if b.Clicked() && !on {
			*chosen = i
			root.SetChanged()
			sound.Play(sound.Tap)
			kit.Log("Members", "picked", "index", i)
		}
	}
	return root.Children(func() {
		i := 0
		if len(online) > 0 {
			head(o.Online, len(online))
			for _, n := range online {
				c.Key("on." + n)
				row(i, n, true)
				i++
			}
		}
		if len(offline) > 0 {
			head(o.Offline, len(offline))
			for _, n := range offline {
				c.Key("off." + n)
				row(i, n, false)
				i++
			}
		}
	})
}
