package code

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MinimapOptions tune a Minimap.
type MinimapOptions struct {
	// Height of the strip; zero gives the default.
	Height float32
}

// Minimap is a narrow file view with a viewport box.
func Minimap(c *ui.Context, label string, lines, visible int, at *int, o MinimapOptions) ui.Element {
	if label == "" || lines < 1 || visible < 1 || visible > lines || at == nil || *at < 0 || *at > lines-visible {
		panic(fmt.Sprintf("code: Minimap %q needs lines, a window and a start, got %d of %d", label, *at, lines))
	}
	t := c.Theme()
	size("Minimap", &o.Height, t.Space(40))
	root := ui.ButtonBase(c).Label(label).Width(t.Space(18)).Height(o.Height).Radius(t.Radius).Clip().
		Background(t.Surface).Role(ui.RoleSlider).Range(0, float64(lines), float64(*at))
	height := ui.Local(root, "height", func() float32 { return o.Height })
	pressY := ui.Local(root, "pressY", func() float32 { return -1 })
	kit.Settle(c, root, []float32{*height}, func() []float32 { return []float32{root.LayoutBounds().H} })
	if got := root.LayoutBounds().H; got > 0 {
		*height = got
	}
	frac := float32(*at) / float32(lines)
	span := float32(visible) / float32(lines)
	root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		for i := 0; i < lines; i++ {
			y := r.Y + r.H*float32(i)/float32(lines)
			h := max(r.H/float32(lines), 1)
			w := r.W * (0.3 + 0.6*float32((i*2654435761)%97)/97)
			p.Fill(ui.Rect{X: r.X + t.Space(0.75), Y: y, W: w - t.Space(0.75), H: h}, t.TextMuted.Alpha(0.4), 0)
		}
		vy := r.Y + r.H*frac
		vh := max(r.H*span, t.Space(2))
		p.Fill(ui.Rect{X: r.X, Y: vy, W: r.W, H: vh}, t.Accent.Alpha(0.18), t.Radius)
		p.Stroke(ui.Rect{X: r.X, Y: vy, W: r.W, H: vh}, t.Accent.Alpha(0.5), t.Radius, 1)
	}))
	root.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind == ui.InputPointerDown && ev.Button == 0 {
			*pressY = ev.Y
		}
		return false
	})
	if *pressY >= 0 {
		*at = clampAt(int(*pressY / *height * float32(lines))-visible/2, lines, visible)
		*pressY = -1
		root.SetChanged()
	}
	if _, dy, held := root.Dragged(); held {
		if moved := clampAt(*at+int(dy*float32(lines) / *height), lines, visible); moved != *at {
			*at = moved
			root.SetChanged()
		}
	}
	released(root, "Minimap", float64(*at))
	move := func(to int, event string) {
		if to = clampAt(to, lines, visible); to == *at {
			return
		}
		*at = to
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("Minimap", event, "at", *at)
	}
	switch {
	case root.Shortcut(0, ui.KeyDown), root.Shortcut(0, ui.KeyRight):
		move(*at+1, "down")
	case root.Shortcut(0, ui.KeyUp), root.Shortcut(0, ui.KeyLeft):
		move(*at-1, "up")
	case root.Shortcut(0, ui.KeyPageDown):
		move(*at+visible, "pagedown")
	case root.Shortcut(0, ui.KeyPageUp):
		move(*at-visible, "pageup")
	case root.Shortcut(0, ui.KeyHome):
		move(0, "home")
	case root.Shortcut(0, ui.KeyEnd):
		move(lines-visible, "end")
	}
	return root
}

// clampAt keeps a viewport start inside the file.
func clampAt(at, lines, visible int) int {
	return max(0, min(lines-visible, at))
}
