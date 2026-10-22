package fx

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/shader"
)

// frost tints the blurred backdrop by Tint.
var frost = shader.MustCompile(`package main

var Tint vec4

func Fragment(dst vec4, pos vec2) vec4 {
	b := backdropSample(pos)
	return mix(b, vec4(Tint.rgb, 1), Tint.a)
}
`)

// Frosted is a panel that blurs what lies behind it.
func Frosted(c *ui.Context, content func()) ui.Element {
	if content == nil {
		panic("fx: Frosted needs content")
	}
	t := c.Theme()
	o := shader.Options{BackdropBlur: t.Space(3), Radius: t.Radius, Uniforms: map[string]any{"Tint": t.Surface.Alpha(0.6)}}
	root := ui.Box(c).AlignSelf(ui.Start).Padding(t.Space(4)).Radius(t.Radius).Border(1, t.Border).
		Draw(func(p *ui.Painter, r ui.Rect) { frost.Draw(p, r, o) })
	return kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0).Children(content)
}
