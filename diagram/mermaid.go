package diagram

import (
	"errors"
	"fmt"

	pele "github.com/Broderick-Westrope/pele-go"
	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// MermaidOptions tune a Mermaid.
type MermaidOptions struct {
	// Label names the diagram; empty takes the source's title.
	Label string
	// Failed heads the error box; empty is English.
	Failed string
}

// inputs are what a scene is made from.
type inputs struct {
	source, font string
	pal          palette
	size         float32
}

// drawn is a scene, or its error, and their inputs.
type drawn struct {
	in  inputs
	sc  *scene
	err error
}

// Mermaid draws Mermaid source in theme colors, shrunk to fit.
func Mermaid(c *ui.Context, source string, o MermaidOptions) ui.Element {
	t := c.Theme()
	in := inputs{source: source, font: t.Font, pal: themed(t), size: t.Rem(1)}
	root := ui.Column(c).FillWidth()
	d := ui.Local(root, "drawn", func() drawn { return drawn{} })
	if d.in != in {
		*d = drawn{in: in}
		d.sc, d.err = render(source, in.pal, in.size, in.font)
		report(d, len(source))
	}
	return kit.Arrive(root.Children(func() {
		if d.err != nil {
			failed(c, d.err, o.Failed)
			return
		}
		sc := d.sc
		label := o.Label
		if label == "" {
			label = sc.title
		}
		if label == "" {
			label = "Diagram"
		}
		ui.Box(c).FillWidth().MaxWidth(sc.w).AspectRatio(sc.w / sc.h).Role(ui.RoleImage).Label(label).
			Draw(kit.Inside(0, sc.paint))
	}), ui.Pose{Y: 6, Fade: 1}, 0)
}

// themed maps Pele's tokens to the theme.
func themed(t *ui.Theme) palette {
	return palette{bg: t.Background, fg: t.Text, muted: t.TextMuted, line: t.Text, surface: t.Surface,
		alt: t.SurfaceHover, border: t.Border, accent: t.Accent, critical: t.Danger, radius: t.Radius,
		series: [8]ui.Color{t.Accent, t.Success, t.Warning, t.Danger, t.Accent.Mix(t.Success, 0.5),
			t.Warning.Mix(t.Danger, 0.5), t.Accent.Mix(t.Danger, 0.5), t.TextMuted}}
}

// report logs one render, never the source.
func report(d *drawn, bytes int) {
	var pe *pele.Error
	switch {
	case d.err == nil:
		kit.Log("Mermaid", "drew", "bytes", bytes, "runs", len(d.sc.runs), "width", d.sc.w, "height", d.sc.h)
	case errors.As(d.err, &pe):
		kit.Log("Mermaid", "failed", "bytes", bytes, "code", pe.Code, "line", pe.Line, "column", pe.Column)
	default:
		kit.Log("Mermaid", "failed", "bytes", bytes, "code", "walk")
	}
}

// failed shows the error; Pele's gets its line and column.
func failed(c *ui.Context, err error, title string) {
	if title == "" {
		title = "Diagram failed"
	}
	body := err.Error()
	var pe *pele.Error
	if errors.As(err, &pe) && pe.Line > 0 {
		body = fmt.Sprintf("%d:%d %s", pe.Line, pe.Column, body)
	}
	feedback.Alert(c, feedback.ToneDanger, title, body)
}

// halo is where an outlined text's stroke is stamped.
var halo = [8][2]float32{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-0.7, -0.7}, {0.7, -0.7}, {-0.7, 0.7}, {0.7, 0.7}}

// paint draws the scene scaled to r's width.
func (sc *scene) paint(p *ui.Painter, r ui.Rect) {
	k := r.W / sc.w
	view := ui.Scaling(k, k).Then(ui.Translation(r.X, r.Y))
	for _, run := range sc.runs {
		p.Transform(run.t.Then(view), 0, 0, func() {
			for _, m := range run.marks {
				switch {
				case m.text != nil && m.halo > 0:
					for _, at := range halo {
						p.TextLayout(m.text, m.x+at[0]*m.halo, m.y+at[1]*m.halo, m.fill)
					}
				case m.text != nil:
					p.TextLayout(m.text, m.x, m.y, m.fill)
				default:
					if m.fill.A > 0 {
						p.FillPath(m.path, m.fill)
					}
					if m.stroke.A > 0 {
						p.StrokePath(m.path, m.width, m.stroke)
					}
				}
			}
		})
	}
}
