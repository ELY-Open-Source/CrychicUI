// Code components for the gallery.
package code

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
	"unicode/utf8"
)

type PlaygroundOptions struct {
	Run    string
	Empty  string
	Height float32
	Output float32
}

// Playground is an editor with an output pane; cue Tap.

func Playground(c *ui.Context, st *editor.State, output *string, onRun func(code string), o PlaygroundOptions) (root, run ui.Element) {
	if st == nil || output == nil || onRun == nil {
		panic("code: Playground needs an editor, output and run callback")
	}
	def(&o.Run, "Run")
	def(&o.Empty, "Run to see output")
	size("Playground", &o.Height, t40(c)*1.6)
	size("Playground", &o.Output, t40(c))
	t := c.Theme()
	root = ui.Column(c).AlignSelf(ui.Start).Width(t.Space(100)).Gap(t.Space(2)).Role(ui.RoleGroup).Label("Playground")
	root.Children(func() {
		editor.View(c, st).Height(o.Height).FillWidth().Border(1, t.Border).Radius(t.Radius).Clip()
		watchEditor(root, st, "Playground")
		ui.Row(c).FillWidth().Children(func() {
			ui.Box(c).Grow(1)
			run = kit.Press(ui.ButtonBase(c).Label(o.Run).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).
				Background(t.Accent).TextColor(t.AccentText))
			run.Children(func() {
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
					ui.Icon(c, icon.Must("play")).Size(16, 16)
					ui.Text(c, o.Run).SingleLine()
				})
			})
		})
		inv, invText := inversePair(t)
		pane := ui.Column(c).FillWidth().Height(o.Output).Padding(t.Space(2)).Background(inv).Radius(t.Radius).Clip()
		pane.Children(func() {
			if *output == "" {
				ui.Text(c, o.Empty).Font("monospace").FontSize(t.Rem(0.85)).TextColor(invText.Alpha(0.5))
			} else {
				ui.Scroll(c).FillWidth().Children(func() {
					ui.Text(c, *output).Font("monospace").FontSize(t.Rem(0.85)).TextColor(invText).Wrap()
				})
			}
		})
	})
	if run.Clicked() {
		src := st.Buffer.String()
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("Playground", "run", "runes", utf8.RuneCountInString(src))
		onRun(src)
	}
	return root, run
}

// inversePair is the terminal palette with fallbacks.
