package code

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

// CodeEditorOptions tune a CodeEditor.
type CodeEditorOptions struct {
	// Height of the editor; zero gives the default.
	Height float32
	// Name is the file shown above; empty hides the bar.
	Name string
	// Format words the caret; Lines counts; Edited names the dot.
	Format, Lines, Edited string
}

// CodeEditor is MyGo's code editor between file and status bars.
func CodeEditor(c *ui.Context, st *editor.State, o CodeEditorOptions) ui.Element {
	if st == nil {
		panic("code: CodeEditor needs a state")
	}
	t := c.Theme()
	size("CodeEditor", &o.Height, t.Space(80))
	def(&o.Format, "Ln %d, Col %d")
	def(&o.Lines, "%d lines")
	def(&o.Edited, "Edited")
	root := ui.Column(c).FillWidth().Height(o.Height).Radius(t.Radius).Border(1, t.Border).Clip().Role(ui.RoleGroup).Label(o.Name)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	root.Children(func() {
		if o.Name != "" {
			c.Key("bar")
			codeHeader(c, o.Name, "", func() {
				edited := st.Buffer != nil && st.Buffer.Modified()
				dot := ui.Box(c).Size(t.Space(2), t.Space(2)).Shrink(0).Margin(0, t.Space(2)).Radius(999).Background(t.Accent).Label(o.Edited)
				k := dot.AnimateTo("on", unit(edited), ui.Bouncy)
				dot.Scale(k, k).Opacity(k)
			})
		}
		// Keyed: the bar appearing must not move the editor.
		c.Key("view")
		editor.View(c, st).Grow(1).FillWidth()
		watchEditor(root, st, "CodeEditor")
		c.Key("status")
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(3)).Background(t.Surface).
			BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			lang, lines := "Plain text", 0
			if st.Buffer != nil {
				lines = st.Buffer.LineCount()
				if l := st.Buffer.Language(); l != nil && l.Name != "" {
					lang = l.Name
				}
			}
			langChip(c, lang)
			kit.Numeric(c, fmt.Sprintf(o.Lines, lines)).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
			ui.Box(c).Grow(1)
			pos := st.CursorPos()
			kit.Numeric(c, fmt.Sprintf(o.Format, pos.Line+1, pos.Col+1)).FontSize(t.Rem(0.75)).TextColor(t.TextMuted).
				ContentTransition(ui.ContentTransition{Numeric: true})
		})
	})
	return root
}
