package input

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SignaturePadOptions tune a SignaturePad.
type SignaturePadOptions struct {
	// Width and Height size the pad; Clear labels its action.
	Width, Height float32
	Clear         string
}

type signatureStroke struct {
	path ui.Path
	at   time.Time
}

// SignaturePad draws pointer strokes and clears them.
func SignaturePad(c *ui.Context, o SignaturePadOptions) ui.Element {
	if o.Width == 0 {
		o.Width = c.Theme().Rem(20)
	}
	if o.Height == 0 {
		o.Height = c.Theme().Rem(12)
	}
	if o.Clear == "" {
		o.Clear = "Clear signature"
	}
	if !finite(float64(o.Width)) || !finite(float64(o.Height)) || o.Width <= 0 || o.Height <= 0 {
		panic("input: SignaturePad needs a positive size")
	}
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(2))
	strokes := ui.Local(root, "strokes", func() []signatureStroke { return nil })
	current := ui.Local(root, "current", func() *ui.Path { return nil })
	drawing := ui.Local(root, "drawing", func() bool { return false })
	completed := ui.Local(root, "completed", func() bool { return false })
	return root.Children(func() {
		pad := ui.Box(c).Size(o.Width, o.Height).Radius(t.Radius).Background(t.Surface).Border(1, t.Border).Clip()
		pad.HandleInput(func(ev ui.InputEvent) bool {
			switch ev.Kind {
			case ui.InputPointerDown:
				if ev.Button != 0 {
					return false
				}
				path := (&ui.Path{}).MoveTo(ev.X, ev.Y)
				*current = path
				*drawing = true
				return true
			case ui.InputPointerMove:
				if *drawing && *current != nil {
					(*current).LineTo(ev.X, ev.Y)
				}
				return *drawing
			case ui.InputPointerUp:
				if *drawing && *current != nil {
					*strokes = append(*strokes, signatureStroke{path: **current, at: ev.Time})
					*current = nil
					*drawing = false
					*completed = true
				}
				return true
			default:
				return false
			}
		})
		pad.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Transform(ui.Translation(r.X, r.Y), 0, 0, func() {
				for _, stroke := range *strokes {
					age := c.Now().Sub(stroke.at)
					alpha := float32(age) / float32(200*time.Millisecond)
					if age < 200*time.Millisecond {
						c.AnimationFrame()
					}
					p.StrokePath(&stroke.path, t.Space(0.5), t.Text.Alpha(min(1, alpha)))
				}
				if *current != nil {
					p.StrokePath(*current, t.Space(0.5), t.Text)
				}
			})
		})
		if *completed {
			*completed = false
			root.SetChanged()
			sound.Play(sound.Tick)
			kit.Log("SignaturePad", "stroke", "count", len(*strokes))
		}
		clear := ui.ButtonBase(c).Label(o.Clear).Gap(t.Space(1)).Padding(t.Space(1), t.Space(2)).Radius(t.Radius).Background(t.Surface).Border(1, t.Border)
		kit.Press(clear).Children(func() {
			ui.Icon(c, icon.Must("x")).Size(16, 16)
			ui.Text(c, o.Clear)
		})
		if kit.Clicked(clear, sound.Close) {
			*strokes = nil
			root.SetChanged()
			kit.Log("SignaturePad", "cleared")
		}
	})
}
