package typography

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestTextTransitions(t *testing.T) {
	cases := []struct {
		name  string
		want  string
		build func(*ui.Context, string) ui.Element
	}{
		{"Title", "second", Title},
		{"Subtitle", "second", Subtitle},
		{"Paragraph", "second", Paragraph},
		{"Caption", "second", Caption},
		{"Overline", "SECOND", Overline},
		{"InlineCode", "second", InlineCode},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			text := "first"
			tt := ui.NewTester(func(c *ui.Context) { tc.build(c, text) }, 300, 200)
			text = "second"
			tt.Frame()
			settles(t, tt, tc.name)
			if !tt.HasText(tc.want) || tt.HasText("first") || tt.HasText("FIRST") {
				t.Errorf("texts %q", tt.Texts())
			}
		})
	}
}

func TestHeading(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Children(func() {
			Heading(c, "Highest", 1)
			Heading(c, "Lowest", 6)
		})
	}, 300, 200)
	settles(t, tt, "Heading")
	if find(t, tt, "Highest").H <= find(t, tt, "Lowest").H {
		t.Error("heading levels did not change size")
	}
	for _, level := range []int{0, 7} {
		panics(t, "Heading level", func() { ui.Render(func(c *ui.Context) { Heading(c, "Bad", level) }, 10, 10, 1) })
	}
}

func TestTextArrivals(t *testing.T) {
	cases := []struct {
		name  string
		build func(*ui.Context, string) ui.Element
	}{
		{"Display", Display},
		{"Lead", Lead},
		{"Eyebrow", Eyebrow},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) { tc.build(c, "Words") }, 300, 200)
			settles(t, tt, tc.name)
			find(t, tt, "Words")
		})
	}
}

func TestLabel(t *testing.T) {
	text := "Email"
	tt := ui.NewTester(func(c *ui.Context) {
		Label(c, text, LabelOptions{Required: true})
	}, 300, 200)
	settles(t, tt, "Label mark")
	if !tt.HasText("*") {
		t.Error("required mark missing")
	}
	text = "Account"
	tt.Frame()
	settles(t, tt, "Label text")
	find(t, tt, "Account")
}
