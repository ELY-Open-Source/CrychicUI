package fx

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/shader"
)

// corners blends four corner colors across the box.
var corners = shader.MustCompile(`package main

var TL vec4
var TR vec4
var BL vec4
var BR vec4

func Fragment(dst vec4, pos vec2) vec4 {
	u := pos / rectSize()
	u = u * u * (3 - 2*u)
	c := mix(mix(TL, TR, u.x), mix(BL, BR, u.x), u.y)
	return vec4(c.rgb*c.a, c.a)
}
`)

// Mesh shifts a four-corner color mesh behind content.
func Mesh(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	cols := palette(t)
	k := 4 * phase(c, root, 16*time.Second)
	step, f := int(k), k-float32(math.Floor(float64(k)))
	corner := func(i int) ui.Color {
		a, b := cols[(i+step)%4], cols[(i+step+1)%4]
		return a.Mix(b, f).Mix(t.Surface, 0.5)
	}
	o := shader.Options{Radius: t.Radius, Uniforms: map[string]any{"TL": corner(0), "TR": corner(1),
		"BR": corner(2), "BL": corner(3)}}
	return root.Draw(func(p *ui.Painter, r ui.Rect) { corners.Draw(p, r, o) }).Children(content)
}
