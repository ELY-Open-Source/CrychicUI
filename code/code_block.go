package code

import (
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CodeBlockOptions name the copy button.
type CodeBlockOptions struct {
	Copy string
}

// CodeBlock is highlighted code with a language and copy.
func CodeBlock(c *ui.Context, lang, code string, o CodeBlockOptions) (root, copy ui.Element) {
	if code == "" {
		panic("code: CodeBlock needs code")
	}
	def(&o.Copy, "Copy")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Clip().Radius(t.Radius).Background(t.Surface).Border(1, t.Border).
		Role(ui.RoleGroup).Label(lang)
	root.Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Padding(t.Space(1), t.Space(2)).Children(func() {
			kit.Muted(c, lang).FontSize(t.Rem(0.75)).Bold().Grow(1).SingleLine()
			copy = button.CopyButton(c, o.Copy, code)
		})
		ui.Box(c).FillWidth().Height(1).Background(t.Border)
		ui.Column(c).FillWidth().Padding(t.Space(2)).Children(func() {
			codeLines(c, lang, code, nil)
		})
	})
	if copy.Clicked() {
		root.SetChanged()
		kit.Log("CodeBlock", "copied", "runes", utf8.RuneCountInString(code))
	}
	return root, copy
}
