package code

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

// CodeEditorOptions tune a CodeEditor.
type CodeEditorOptions struct {
	// Height of the editor; zero gives the default.
	Height float32
	Format string
}

// CodeEditor is MyGo's code editor with a status bar.
func CodeEditor(c *ui.Context, st *editor.State, o CodeEditorOptions) ui.Element {
	if st == nil {
		panic("code: CodeEditor needs a state")
	}
	t := c.Theme()
	size("CodeEditor", &o.Height, t.Space(80))
	def(&o.Format, "Ln %d, Col %d")
	root := ui.Column(c).FillWidth().Height(o.Height).Radius(t.Radius).Border(1, t.Border).Clip().Role(ui.RoleGroup)
	root.Children(func() {
		view := editor.View(c, st).Grow(1).FillWidth()
		kit.Arrive(view, ui.Pose{Fade: 1}, 0)
		version := ui.Local(view, "version", func() uint64 { return st.Buffer.Version() })
		typed := ui.Local(view, "typed", func() int { return 0 })
		view.HandleInput(func(ev ui.InputEvent) bool {
			*typed = 2
			return false
		})
		if st.Buffer.Version() != *version {
			*version = st.Buffer.Version()
			if *typed > 0 {
				*typed = 0
				root.SetChanged()
				sound.Play(sound.Tick)
				kit.Log("CodeEditor", "edited")
			}
		} else if *typed > 0 {
			*typed--
		}
		ui.Box(c).FillWidth().Height(1).Background(t.Border)
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2)).Background(t.Surface).Children(func() {
			lang := "Plain text"
			if st.Buffer != nil {
				if l := st.Buffer.Language(); l != nil && l.Name != "" {
					lang = l.Name
				}
			}
			kit.Muted(c, lang).FontSize(t.Rem(0.8)).SingleLine()
			ui.Box(c).Grow(1)
			pos := st.CursorPos()
			kit.Numeric(c, fmt.Sprintf(o.Format, pos.Line+1, pos.Col+1)).FontSize(t.Rem(0.8)).TextColor(t.TextMuted)
		})
	})
	return root
}
