package typography

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestInlineCodeHugsText(t *testing.T) {
	width := float32(0)
	tt := ui.NewTester(func(c *ui.Context) {
		width = InlineCode(c, "x").LayoutBounds().W
	}, 300, 200)
	tt.Frame()
	if width >= 300 {
		t.Errorf("code pill stretched to %v", width)
	}
}

func TestHighlightHugsText(t *testing.T) {
	width := float32(0)
	tt := ui.NewTester(func(c *ui.Context) {
		width = Highlight(c, "mm").LayoutBounds().W
	}, 300, 200)
	tt.Frame()
	if width >= 300 {
		t.Errorf("highlight stretched to %v", width)
	}
}

func TestAbbrHugsText(t *testing.T) {
	width := float32(0)
	tt := ui.NewTester(func(c *ui.Context) {
		width = Abbr(c, "API", "Application Programming Interface").LayoutBounds().W
	}, 300, 200)
	tt.Frame()
	if width >= 300 {
		t.Errorf("abbreviation stretched to %v", width)
	}
}
