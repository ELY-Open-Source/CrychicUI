package code

import (
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SnippetOptions add a description and name the copy button.
type SnippetOptions struct {
	Description string
	Copy        string
}

// Snippet is short code with a title and copy.
func Snippet(c *ui.Context, title, lang, code string, o SnippetOptions) (root, copy ui.Element) {
	if title == "" || code == "" {
		panic("code: Snippet needs a title and code")
	}
	def(&o.Copy, "Copy")
	t := c.Theme()
	root = kit.Card(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Role(ui.RoleGroup).Label(title)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icon.Must("square-terminal")).Size(16, 16).Shrink(0).TextColor(t.Accent)
			ui.Text(c, title).Bold().Grow(1).MinWidth(0).SingleLine()
			if lang != "" {
				langChip(c, lang)
			}
		})
		if o.Description != "" {
			kit.Muted(c, o.Description).FontSize(t.Rem(0.85)).Wrap()
		}
		ui.Row(c).FillWidth().AlignItems(ui.Center).Radius(t.Radius).Background(t.Surface).Border(1, t.Border).
			Padding(0, t.Space(1), 0, 0).Children(func() {
			ui.Column(c).Grow(1).MinWidth(0).Children(func() { codeBody(c, lang, code, false, nil, nil) })
			copy = button.CopyButton(c, o.Copy, code)
		})
	})
	if copy.Clicked() {
		root.SetChanged()
		kit.Log("Snippet", "copied", "runes", utf8.RuneCountInString(code))
	}
	return root, copy
}
