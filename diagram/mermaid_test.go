package diagram

import (
	"bytes"
	"context"
	"image"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	pele "github.com/Broderick-Westrope/pele-go"
	"github.com/egoist/mygo/ui"
)

const flow = "flowchart LR\n  accTitle: Order flow\n  Cart[Cart] --> Pay{Paid?}\n  Pay -- yes --> Ship([Ship order])\n" +
	"  Pay -- no --> Cart\n  style Ship fill:#0a7,stroke:#064,color:#fff"

// show opens a window holding one Mermaid of *src.
func show(src *string, o MermaidOptions, w int, theme *ui.Theme) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		if theme != nil {
			*theme = *c.Theme()
		}
		Mermaid(c, *src, o)
	}, w, 400)
	tt.Advance(time.Second)
	return tt
}

// logged counts log lines holding s.
func logged(t *testing.T) func(s string) int {
	t.Helper()
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	t.Cleanup(func() { slog.SetDefault(old) })
	return func(s string) int { return strings.Count(buf.String(), s) }
}

func picture(t *testing.T, tt *ui.Tester, name string) ui.Rect {
	t.Helper()
	info, ok := tt.AccessByRole(ui.RoleImage, name)
	if !ok {
		t.Fatalf("no image named %q", name)
	}
	return info.Bounds
}

func at(img *image.RGBA, x, y float32) ui.Color {
	c := img.RGBAAt(int(x), int(y))
	return ui.Color{R: c.R, G: c.G, B: c.B, A: c.A}
}

// spots finds where Pele's SVG puts each node and arrowhead.
func spots(t *testing.T, theme *ui.Theme) (nodes map[string][2]float32, heads [][2]float32, w float32) {
	t.Helper()
	off := false
	r, err := pele.Render(context.Background(), flow, pele.Options{Measurer: measurer{theme.Font},
		FontSize: pele.Ptr(float64(theme.Rem(1))), Links: &off, Images: &off})
	if err != nil {
		t.Fatal(err)
	}
	num := func(s string) float32 {
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			t.Fatal(err)
		}
		return float32(f)
	}
	nodes = map[string][2]float32{}
	for _, m := range regexp.MustCompile(`data-id="(\w+)" transform="translate\(([\d.]+),([\d.]+)\)"`).FindAllStringSubmatch(r.SVG, -1) {
		nodes[m[1]] = [2]float32{num(m[2]), num(m[3])}
	}
	for _, m := range regexp.MustCompile(`class="pele-marker" d="M([\d.]+),([\d.]+)L([\d.]+),`).FindAllStringSubmatch(r.SVG, -1) {
		heads = append(heads, [2]float32{(num(m[1]) + 2*num(m[3])) / 3, num(m[2])})
	}
	if len(nodes) != 3 || len(heads) != 3 {
		t.Fatalf("%d nodes and %d arrowheads in %s", len(nodes), len(heads), r.SVG)
	}
	return nodes, heads, float32(r.Width)
}

func TestMermaidDrawsWherePeleDraws(t *testing.T) {
	src, theme := flow, ui.Theme{}
	tt := show(&src, MermaidOptions{}, 700, &theme)
	nodes, heads, w := spots(t, &theme)
	if box := picture(t, tt, "Order flow"); box.X != 0 || box.Y != 0 || box.W != w {
		t.Fatalf("picture %+v, want Pele's width %v at the corner", box, w)
	}
	img := tt.Image()
	for id, fill := range map[string]ui.Color{"Cart": theme.Surface, "Pay": theme.Surface, "Ship": ui.Hex("#00aa77")} {
		spot := nodes[id]
		if got := at(img, spot[0], spot[1]-14); got != fill {
			t.Errorf("node %s draws %+v, want %+v", id, got, fill)
		}
		if ink(img, spot[0]-12, spot[1]-6, 24, 12) == 0 {
			t.Errorf("node %s has no label", id)
		}
	}
	for _, h := range heads {
		if got := at(img, h[0], h[1]); got != theme.Text {
			t.Errorf("arrowhead at %v draws %+v, want %+v", h, got, theme.Text)
		}
	}
	cart, pay := nodes["Cart"], nodes["Pay"]
	if got := at(img, (cart[0]+pay[0])/2-20, heads[0][1]); got == theme.Background {
		t.Error("no edge between Cart and Pay")
	}
	if got := at(img, w+20, 20); got != theme.Background {
		t.Errorf("beside the diagram %+v, want the background", got)
	}
}

// ink counts pixels unlike the box's corner.
func ink(img *image.RGBA, x, y, w, h float32) int {
	n, ground := 0, at(img, x, y)
	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			if at(img, px, py) != ground {
				n++
			}
		}
	}
	return n
}

func TestMermaidFitsItsWidth(t *testing.T) {
	src, theme := flow, ui.Theme{}
	tt := show(&src, MermaidOptions{}, 260, &theme)
	nodes, _, w := spots(t, &theme)
	sc, err := render(flow, themed(&theme), theme.Rem(1), theme.Font)
	if err != nil {
		t.Fatal(err)
	}
	box := picture(t, tt, "Order flow")
	if want := 260 * sc.h / sc.w; w <= 260 || sc.w != w || box.W != 260 || box.H < want-0.5 || box.H > want+0.5 {
		t.Fatalf("picture %+v, want 260 x %v of a scene %v wide", box, want, w)
	}
	k := 260 / w
	if got := at(tt.Image(), nodes["Ship"][0]*k, (nodes["Ship"][1]-14)*k); got != ui.Hex("#00aa77") {
		t.Errorf("shrunk node draws %+v", got)
	}
}

func TestMermaidFollowsThemeAndTextScale(t *testing.T) {
	count := logged(t)
	src, theme := flow, ui.Theme{}
	tt := show(&src, MermaidOptions{}, 900, &theme)
	for range 5 {
		tt.Frame()
	}
	light, size := theme.Surface, picture(t, tt, "Order flow")
	if n := count("Mermaid drew"); n != 1 {
		t.Fatalf("drew %d times over six frames", n)
	}
	tt.SetDark(true)
	tt.Advance(time.Second)
	if theme.Surface == light || count("Mermaid drew") != 2 {
		t.Fatalf("dark theme: surface %+v, drew %d", theme.Surface, count("Mermaid drew"))
	}
	nodes, _, _ := spots(t, &theme)
	if got := at(tt.Image(), nodes["Cart"][0], nodes["Cart"][1]-14); got != theme.Surface {
		t.Errorf("Cart draws %+v in the dark, want %+v", got, theme.Surface)
	}
	tt.SetPreferences(ui.Preferences{TextScale: 1.5})
	tt.Advance(time.Second)
	if grown := picture(t, tt, "Order flow"); grown.W <= size.W || count("Mermaid drew") != 3 {
		t.Errorf("text scale 1.5: width %v from %v, drew %d", grown.W, size.W, count("Mermaid drew"))
	}
	if tt.FramePending() {
		t.Error("still moving")
	}
}

func TestMermaidShowsPelesError(t *testing.T) {
	count := logged(t)
	src := "flowchart LR\n  Cart -->"
	tt := show(&src, MermaidOptions{}, 500, nil)
	for range 3 {
		tt.Frame()
	}
	if !tt.HasText("Diagram failed") || count("Mermaid failed") != 1 || count("line=3 column=1") != 1 {
		t.Fatalf("texts %q, failed logged %d", tt.Texts(), count("Mermaid failed"))
	}
	if count("Cart") != 0 {
		t.Error("the log holds the source")
	}
	told := false
	for _, s := range tt.Texts() {
		told = told || strings.HasPrefix(s, "3:1 Parse error on line 3:") && strings.Contains(s, "Expecting 'NODE_STRING', got 'EOF'")
	}
	if !told {
		t.Errorf("texts %q lack Pele's message", tt.Texts())
	}
	if _, ok := tt.AccessByRole(ui.RoleImage, "Diagram"); ok {
		t.Error("a failed diagram still draws")
	}
	src = "flowchart LR\n  Cart --> Pay"
	tt.Frame()
	tt.Advance(time.Second)
	picture(t, tt, "Diagram")
	if tt.HasText("Diagram failed") || count("Mermaid drew") != 1 {
		t.Errorf("after the fix: texts %q, drew %d", tt.Texts(), count("Mermaid drew"))
	}
	src = "no such diagram"
	tt.Frame()
	if !tt.HasText("No diagram type detected for the given text.") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestMermaidLabels(t *testing.T) {
	src := flow
	picture(t, show(&src, MermaidOptions{Label: "Checkout"}, 500, nil), "Checkout")
	src = "pie\n  \"Tea\" : 3\n  \"Coffee\" : 5"
	picture(t, show(&src, MermaidOptions{}, 500, nil), "Diagram")
	src = "not a diagram"
	if tt := show(&src, MermaidOptions{Failed: "图表出错"}, 500, nil); !tt.HasText("图表出错") {
		t.Errorf("texts %q", tt.Texts())
	}
}
