package chat

import (
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CodeMessageOptions name the copy button.
type CodeMessageOptions struct {
	Copy string
}

// CodeMessage is a code block bubble with copy.
func CodeMessage(c *ui.Context, lang, code string, o CodeMessageOptions) (root, copy ui.Element) {
	if code == "" {
		panic("chat: CodeMessage needs code")
	}
	def(&o.Copy, "Copy")
	t := c.Theme()
	root = ui.Column(c).AlignSelf(ui.Start).MaxWidth(t.Space(100)).Gap(t.Space(1)).Padding(t.Space(2), t.Space(3)).
		Radius(t.Radius*2).Background(t.Surface).Border(1, t.Border).Role(ui.RoleGroup).Label(lang)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			kit.Muted(c, lang).FontSize(t.Rem(0.75)).Bold().Grow(1).SingleLine()
			copy = button.CopyButton(c, o.Copy, code)
		})
		ui.Text(c, code).Font("monospace").FontSize(t.Rem(0.85)).Wrap()
	})
	if copy.Clicked() {
		kit.Log("CodeMessage", "copied", "runes", utf8.RuneCountInString(code))
	}
	return root, copy
}
