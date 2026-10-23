package code

import (
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SnippetOptions name the copy button.
type SnippetOptions struct {
	Copy string
}

// Snippet is short code with a title and copy.
func Snippet(c *ui.Context, title, lang, code string, o SnippetOptions) (root, copy ui.Element) {
	if title == "" || code == "" {
		panic("code: Snippet needs a title and code")
	}
	def(&o.Copy, "Copy")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup).Label(title)
	root.Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			ui.Text(c, title).Bold().Grow(1).SingleLine()
			copy = button.CopyButton(c, o.Copy, code)
		})
		ui.Column(c).FillWidth().Clip().Radius(t.Radius).Background(t.Surface).Border(1, t.Border).Padding(t.Space(2)).Children(func() {
			codeLines(c, lang, code, nil)
		})
	})
	if copy.Clicked() {
		root.SetChanged()
		kit.Log("Snippet", "copied", "runes", utf8.RuneCountInString(code))
	}
	return root, copy
}
