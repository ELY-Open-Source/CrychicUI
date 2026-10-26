package content

import (
	"bytes"
	"image"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestDocTOCMarkerStaysInside(t *testing.T) {
	items, active := outline(), 4
	tt := ui.NewTester(padded(func(c *ui.Context) { DocTOC(c, items, &active, DocTOCOptions{}) }), 400, 400)
	settles(t, tt)
	items, active = items[:2], 1
	tt.Frame()
	tt.Advance(16 * time.Millisecond)
	r := rectOf(t, tt, "Install")
	img, n := tt.Image(), 0
	for y := int(r.Y+r.H) + 2; y < int(r.Y+r.H)+150 && y < img.Bounds().Dy(); y++ {
		for x := 8; x < 16; x++ {
			p := img.RGBAAt(x, y)
			if int(p.B) > int(p.R)+60 && int(p.B) > int(p.G)+30 {
				n++
			}
		}
	}
	if n > 0 {
		t.Errorf("marker paints %d pixels below the outline", n)
	}
}

func TestMarkdownBodyKeepsFocus(t *testing.T) {
	title := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Markdown(c, "```go\nvar a = 1\n```", MarkdownOptions{Title: title}) }), 480, 400)
	settles(t, tt)
	click(t, tt, "Copy")
	title = "README.md"
	tt.Frame()
	if !tt.Focused("Copy code") {
		t.Error("the code copy button lost focus when the bar appeared")
	}
	panics(t, `"README.md" needs text`, func(c *ui.Context) { Markdown(c, " ", MarkdownOptions{Title: "README.md"}) })
}

func TestDocCodeTabsCopyFollowsSample(t *testing.T) {
	samples := []CodeSample{{Name: "Go", Lang: "go", Code: "var a = 1"}, {Name: "JS", Lang: "javascript", Code: "let a = 1"}}
	picked := 0
	tt := ui.NewTester(padded(func(c *ui.Context) { DocCodeTabs(c, samples, &picked, DocCodeTabsOptions{}) }), 520, 300)
	settles(t, tt)
	click(t, tt, "Copy code")
	click(t, tt, "JS")
	tt.Advance(300 * time.Millisecond)
	b := rectOf(t, tt, "Copy code")
	crop := func(img *image.RGBA) []byte {
		var out []byte
		for y := int(b.Y); y < int(b.Y+b.H); y++ {
			out = append(out, img.Pix[y*img.Stride+int(b.X)*4:y*img.Stride+int(b.X+b.W)*4]...)
		}
		return out
	}
	now := crop(tt.Image())
	tt.Advance(3 * time.Second)
	if !bytes.Equal(now, crop(tt.Image())) {
		t.Error("the button still shows Go's copy check over JS")
	}
}

func TestCompoundCopyFits(t *testing.T) {
	tt := ui.NewTester(within(320, func(c *ui.Context) {
		ReadingTime(c, 1611, ReadingTimeOptions{Progress: 0.4, Read: "Approximate reading time %d minutes", Left: "About %d minutes still to go"})
	}), 480, 200)
	settles(t, tt)
	// A 480 window, so overflow past 320 is not clipped.
	inside(t, "read", rectOf(t, tt, "Approximate reading time 8 minutes"), 0, 320)
	inside(t, "left", rectOf(t, tt, "About 5 minutes still to go"), 0, 320)
	vote := Unrated
	o := WasHelpfulOptions{Yes: "Yes, it answered my question", No: "No, something is missing here"}
	tt = ui.NewTester(within(320, func(c *ui.Context) { WasHelpful(c, &vote, o) }), 480, 400)
	settles(t, tt)
	// Card padding is 16 DIPs.
	inside(t, "yes", rectOf(t, tt, o.Yes), 0, 320-16)
	inside(t, "no", rectOf(t, tt, o.No), 0, 320-16)
}
