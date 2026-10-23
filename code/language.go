// Code components for the gallery.
package code

import (
	"fmt"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	"strconv"
	"time"
)

type LangShare struct {
	Name  string
	Share float64
}

// Language is a language bar of shares; segments grow.

func Language(c *ui.Context, label string, shares []LangShare) ui.Element {
	if label == "" || len(shares) == 0 {
		panic("code: Language needs a label and shares")
	}
	total := 0.0
	for _, s := range shares {
		if s.Name == "" || s.Share <= 0 || s.Share > 1 {
			panic("code: Language needs a name and a share in (0,1] in every entry")
		}
		total += s.Share
	}
	if total > 1.0001 {
		panic(fmt.Sprintf("code: Language shares add to %v, over 1", total))
	}
	t := c.Theme()
	inks := []ui.Color{t.Accent, t.Success, t.Warning, t.Danger}
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		bar := ui.Box(c).FillWidth().Height(t.Space(2)).Radius(999).Clip()
		g := bar.Keyframes("crychic.grow", []ui.Keyframe{{Value: 0}, {At: 500 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
		bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
			x := r.X
			for i, s := range shares {
				grow := bar.AnimateTo(fmt.Sprint("seg.", i), float32(s.Share), ui.Smooth)
				w := r.W * grow * g
				p.Fill(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, inks[i%len(inks)], 0)
				x += w
			}
		}))
		row := ui.Row(c).FillWidth().Gap(t.Space(4))
		row.Children(func() {
			for i, s := range shares {
				c.Key(s.Name)
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(inks[i%len(inks)])
					kit.Muted(c, s.Name+" "+strconv.Itoa(int(s.Share*100+0.5))+"%").FontSize(t.Rem(0.8))
				})
			}
		})
	})
	return root
}

// Shortcut is one keyboard shortcut.
