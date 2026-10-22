package fx

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/shader"
)

// field colors Cell-sized cells A, B, C by three sines.
var field = shader.MustCompile(`package main

var A vec4
var B vec4
var C vec4
var T float
var Cell float

func Fragment(dst vec4, pos vec2) vec4 {
	u := (floor(pos/Cell) + 0.5) * Cell / rectSize()
	v := sin(6.2832*(1.5*u.x+T)) + sin(6.2832*(2*u.y-T)) + sin(6.2832*(u.x+u.y+2*T))
	f := v/6 + 0.5
	c := mix(mix(A, B, smoothstep(0, 0.5, f)), C, smoothstep(0.5, 1, f))
	return vec4(c.rgb*c.a, c.a)
}
`)

// Plasma colors cells by moving sine fields behind content.
func Plasma(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	mute := func(col ui.Color) ui.Color { return col.Mix(t.Surface, 0.55) }
	o := shader.Options{Radius: t.Radius, Uniforms: map[string]any{"A": mute(t.Accent), "B": mute(t.Success),
		"C": mute(t.Warning), "T": phase(c, root, 10*time.Second), "Cell": t.Space(2)}}
	return root.Draw(func(p *ui.Painter, r ui.Rect) { field.Draw(p, r, o) }).Children(content)
}
