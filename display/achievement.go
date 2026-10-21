package display

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/egoist/mygo/ui"
)

// Achievement shows an icon, name and unlock shine.
func Achievement(c *ui.Context, name, iconName string, unlocked bool) ui.Element {
	if name == "" {
		panic("display: Achievement needs a name")
	}
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).Padding(t.Space(1), t.Space(2)).Radius(999).
		Border(1, t.Border).Clip()
	was := ui.Local(root, "was", func() bool { return unlocked })
	run := ui.Local(root, "run", func() uint64 { return 0 })
	if unlocked != *was {
		*was = unlocked
		if unlocked {
			(*run)++
		}
	}
	bg, fg := t.Surface, t.TextMuted
	if unlocked {
		bg, fg = t.Accent.Alpha(0.14), t.Accent
	}
	root.Background(root.AnimateColor("bg", bg, ui.Snappy)).TextColor(root.AnimateColor("fg", fg, ui.Snappy))
	if *run > 0 {
		shine := root.Keyframes("shine", []ui.Keyframe{{Value: 0}, {At: 500 * time.Millisecond, Value: 1, Ease: ui.EaseOut}},
			ui.KeyframeOptions{Run: *run})
		root.DrawOver(func(p *ui.Painter, r ui.Rect) {
			if shine < 1 {
				w := r.W / 3
				x := r.X - w + (r.W+2*w)*shine
				g := ui.LinearGradient{Angle: 90, Stops: []ui.GradientStop{
					{Color: t.AccentText.Alpha(0), At: 0},
					{Color: t.AccentText.Alpha(0.35), At: 0.5},
					{Color: t.AccentText.Alpha(0), At: 1},
				}}
				p.FillGradient(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, g, 0)
			}
		})
	}
	return root.Children(func() {
		ui.Icon(c, icon.Must(iconName)).Size(16, 16)
		ui.Text(c, name).FontSize(t.Rem(0.85))
	})
}
