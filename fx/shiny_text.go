package fx

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/shader"
)

// shine lifts Base to Glint in a band at At.
var shine = shader.MustCompile(`package main

var Base vec4
var Glint vec4
var At float

func Fragment(dst vec4, pos vec2) vec4 {
	x := pos.x / rectSize().x
	c := mix(Base, Glint, smoothstep(0.15, 0, abs(x-At)))
	return vec4(c.rgb*c.a, c.a)
}
`)

// ShinyText passes a sheen over muted text in a loop.
func ShinyText(c *ui.Context, text string) ui.Element {
	s := headline(c, text)
	t := c.Theme()
	e := filled(c, s)
	at := -0.2 + 1.4*phase(c, e, 2500*time.Millisecond)
	o := shader.Options{Uniforms: map[string]any{"Base": t.TextMuted, "Glint": t.Text, "At": at}}
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		shine.Fill(p, r, o, func() { p.RichText(r.X, r.Y, 0, s) })
	})
}
