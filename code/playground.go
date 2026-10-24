package code

import (
	"fmt"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

// PlaygroundOptions tune a Playground.
type PlaygroundOptions struct {
	Run, Empty, Output string
	// Heights of the editor and output; zero gives defaults.
	Height, OutputHeight float32
}

// Playground is an editor, Run and output; cue Tap.
func Playground(c *ui.Context, label string, st *editor.State, output *string, onRun func(code string), o PlaygroundOptions) (root, run ui.Element) {
	if label == "" || st == nil || output == nil || onRun == nil {
		panic(fmt.Sprintf("code: Playground %q needs an editor, output and run callback", label))
	}
	def(&o.Run, "Run")
	def(&o.Empty, "Run to see output")
	def(&o.Output, "Output")
	size("Playground", &o.Height, t40(c)*1.6)
	size("Playground", &o.OutputHeight, t40(c))
	t := c.Theme()
	bg, fg := inversePair(t)
	root = ui.Column(c).FillWidth().Radius(t.Radius).Border(1, t.Border).Clip().Role(ui.RoleGroup).Label(label)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	root.Children(func() {
		codeHeader(c, label, "", func() { run = action(c, o.Run, "play", true) })
		c.Key("view")
		editor.View(c, st).Height(o.Height).FillWidth()
		watchEditor(root, st, "Playground")
		c.Key("output")
		ui.Column(c).FillWidth().Height(o.OutputHeight).Background(bg).TextColor(fg).Clip().Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(3)).
				BorderWidth(0, 0, 1, 0).BorderColor(fg.Alpha(0.12)).Children(func() {
				ui.Text(c, o.Output).FontSize(t.Rem(0.7)).Bold().TextColor(fg.Alpha(0.6))
			})
			if *output == "" {
				ui.Text(c, o.Empty).Font("monospace").FontSize(t.Rem(0.8)).TextColor(fg.Alpha(0.5)).Padding(t.Space(2), t.Space(3))
				return
			}
			ui.Scroll(c).FillWidth().Grow(1).Padding(t.Space(2), t.Space(3)).Children(func() {
				// Keyed by output: a new run fades in.
				c.Key(*output)
				out := ui.Text(c, *output).Font("monospace").FontSize(t.Rem(0.8)).TextColor(fg).Wrap()
				kit.Arrive(out, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
			})
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
