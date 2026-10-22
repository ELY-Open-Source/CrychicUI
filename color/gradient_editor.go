package color

import (
	"cmp"
	"fmt"
	"math"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// GradientEditorOptions tune a GradientEditor.
type GradientEditorOptions struct {
	// Stop names each handle; it formats the stop's percent.
	Stop string
}

// GradientEditor drags stops along a bar; a click adds one.
func GradientEditor(c *ui.Context, label string, stops *[]ui.GradientStop, o GradientEditorOptions) ui.Element {
	if label == "" || stops == nil || len(*stops) < 2 {
		panic("color: GradientEditor needs a label and two stops")
	}
	for i, s := range *stops {
		if !(s.At >= 0 && s.At <= 1) {
			panic(fmt.Sprintf("color: GradientEditor stop %d is at %v, outside 0 to 1", i, s.At))
		}
	}
	def(&o.Stop, "Stop at %.0f%%")
	t := c.Theme()
	w, h, hs := t.Space(72), t.Space(8), t.Space(5)
	root := ui.Column(c).AlignSelf(ui.Start).AlignItems(ui.Center).Gap(t.Space(1)).Label(label)
	return root.Children(func() {
		bar := kit.Press(ui.ButtonBase(c).Label(label).Size(w, h).Radius(t.Radius).Border(1, t.Border))
		if bar.Clicked() {
			at := float32(0.5)
			if x, _, over := bar.PointerPosition(); over {
				at = max(0, min(x/w, 1))
			}
			*stops = append(*stops, ui.GradientStop{Color: sample(*stops, at), At: at})
			root.SetChanged()
			sound.Play(sound.Pop)
			kit.Log("GradientEditor", "added", "stops", len(*stops))
		}
		g := ui.LinearGradient{Angle: 90, Stops: sorted(*stops)}
		bar.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
			checker(p, r, t)
			p.FillGradient(r, g, 0)
		}))
		ui.Box(c).Size(w+hs, hs).Children(func() {
			for i := range *stops {
				c.Key(i)
				handle(c, root, stops, i, w, hs, o.Stop)
			}
		})
	})
}

// handle is one draggable stop; it cues Drop.
func handle(c *ui.Context, root ui.Element, stops *[]ui.GradientStop, i int, w, hs float32, name string) {
	t := c.Theme()
	s := &(*stops)[i]
	b := ui.Box(c).Absolute().Top(0).Left(s.At*w).Size(hs, hs).Radius(999).Border(2, t.Text).
		Focusable().Role(ui.RoleSlider).Label(fmt.Sprintf(name, s.At*100)).Range(0, 1, float64(s.At)).Cursor(ui.CursorGrab)
	held := ui.Local(b, "held", func() bool { return false })
	moved := false
	if dx, _, ok := b.Dragged(); ok {
		*held = true
		if dx != 0 {
			s.At, moved = max(0, min(s.At+dx/w, 1)), true
		}
	} else if *held {
		*held = false
		sound.Play(sound.Drop)
		kit.Log("GradientEditor", "dropped", "stop", i)
	}
	for _, k := range []struct {
		key  ui.Key
		step float32
	}{{ui.KeyLeft, -0.01}, {ui.KeyRight, 0.01}} {
		if b.Shortcut(0, k.key) {
			s.At, moved = max(0, min(float32(math.Round(float64((s.At+k.step)*100)))/100, 1)), true
			sound.Play(sound.Drop)
		}
	}
	if moved {
		root.SetChanged()
	}
	k := 1 + 0.2*b.AnimateTo("grip", unit(*held), ui.Snappy)
	b.Background(b.AnimateColor("col", s.Color, ui.Snappy)).Scale(k, k)
}

// sorted copies stops in order along the bar.
func sorted(stops []ui.GradientStop) []ui.GradientStop {
	out := slices.Clone(stops)
	slices.SortStableFunc(out, func(a, b ui.GradientStop) int { return cmp.Compare(a.At, b.At) })
	return out
}

// sample is the gradient's color at at.
func sample(stops []ui.GradientStop, at float32) ui.Color {
	s := sorted(stops)
	if at <= s[0].At {
		return s[0].Color
	}
	for i := 1; i < len(s); i++ {
		if at <= s[i].At {
			return s[i-1].Color.Mix(s[i].Color, (at-s[i-1].At)/max(s[i].At-s[i-1].At, 1e-6))
		}
	}
	return s[len(s)-1].Color
}
