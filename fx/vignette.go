package fx

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/shader"
)

// shade darkens toward the corners by Strength.
var shade = shader.MustCompile(`package main

var Edge vec4
var Strength float

func Fragment(dst vec4, pos vec2) vec4 {
	d := length((pos/rectSize() - 0.5) * 2) / 1.4142
	a := Edge.a * Strength * smoothstep(0.3, 1, d)
	return vec4(Edge.rgb*a, a)
}
`)

// Vignette darkens content's edges toward the corners, breathing.
func Vignette(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	k := phase(c, root, 6*time.Second)
	strength := 0.7 + 0.3*float32(math.Cos(2*math.Pi*float64(k)))
	o := shader.Options{Radius: t.Radius, Uniforms: map[string]any{"Edge": kit.Shade(0.4), "Strength": strength}}
	return root.Children(content).DrawOver(func(p *ui.Painter, r ui.Rect) { shade.Draw(p, r, o) })
}
