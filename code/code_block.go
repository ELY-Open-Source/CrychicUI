package code

import (
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CodeBlockOptions name the file and the copy button.
type CodeBlockOptions struct {
	// Title is a file name shown before the language.
	Title string
	Copy  string
}

// CodeBlock is numbered, highlighted code with a language and copy.
func CodeBlock(c *ui.Context, lang, code string, o CodeBlockOptions) (root, copy ui.Element) {
	if lang == "" || code == "" {
		panic("code: CodeBlock needs a language and code")
	}
	def(&o.Copy, "Copy")
	def(&o.Title, lang)
	t := c.Theme()
	root = ui.Column(c).FillWidth().Clip().Radius(t.Radius).Background(t.Surface).Border(1, t.Border).
		Role(ui.RoleGroup).Label(o.Title)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	root.Children(func() {
		codeHeader(c, o.Title, lang, func() { copy = button.CopyButton(c, o.Copy, code) })
		codeBody(c, lang, code, true, nil, nil)
	})
	if copy.Clicked() {
		root.SetChanged()
		kit.Log("CodeBlock", "copied", "runes", utf8.RuneCountInString(code))
	}
	return root, copy
}
