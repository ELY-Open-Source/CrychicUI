// Code components for the gallery.
package code

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"strings"
	"unicode/utf8"
)

type REPLLine struct {
	In, Out string
	Err     bool
}

// REPLOptions are the copy knobs of REPL.

type REPLOptions struct {
	Placeholder string
	Run         string
	FormatError func(err error) string
	Height      float32
}

// REPL evaluates arithmetic with history; cue Send.

func REPL(c *ui.Context, label string, history *[]REPLLine, draft *string, o REPLOptions) (root, send ui.Element) {
	if label == "" || history == nil || draft == nil {
		panic("code: REPL needs a label and bound values")
	}
	def(&o.Placeholder, "1 + 2 * 3")
	def(&o.Run, "Run")
	if o.FormatError == nil {
		o.FormatError = func(err error) string { return "error: " + err.Error() }
	}
	size("REPL", &o.Height, t40(c))
	t := c.Theme()
	root = ui.Column(c).AlignSelf(ui.Start).Width(t.Space(90)).Gap(t.Space(1)).Role(ui.RoleGroup).Label(label)
	pending := ui.Local(root, "pending", func() string { return "" })
	root.Children(func() {
		inv, invText := inversePair(t)
		scroll := ui.Column(c).FillWidth().MaxHeight(o.Height).Padding(t.Space(2)).Background(inv).Radius(t.Radius)
		scroll.Children(func() {
			ui.Scroll(c).FillWidth().Children(func() {
				col := ui.Column(c).FillWidth().Gap(t.Space(0.5))
				col.Children(func() {
					for i, l := range *history {
						c.Key(i)
						ui.Text(c, "> "+l.In).Font("monospace").FontSize(t.Rem(0.85)).TextColor(invText.Alpha(0.7)).Wrap()
						out := ui.Text(c, "= "+l.Out).Font("monospace").FontSize(t.Rem(0.85)).Wrap()
						if l.Err {
							out.TextColor(t.Danger)
						} else {
							out.TextColor(invText)
						}
					}
				})
			})
		})
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			in := ui.TextInput(c, draft).Label(o.Placeholder).Font("monospace").Grow(1)
			in.OnTextKey(func(k ui.TextKey) bool {
				if k.Key == ui.KeyEnter {
					*pending = strings.TrimSpace(k.Text())
					k.SetText("")
					return true
				}
				return false
			})
			if in.Changed() {
				root.SetChanged()
			}
			send = kit.Press(ui.ButtonBase(c).Label(o.Run).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999).
				Background(t.Accent).TextColor(t.AccentText).Disabled(strings.TrimSpace(*draft) == ""))
			send.Children(func() {
				ui.Icon(c, icon.Must("play")).Size(16, 16)
			})
		})
	})
	run := func(src string) {
		line := REPLLine{In: src}
		v, err := eval(src)
		if err != nil {
			line.Out = o.FormatError(err)
			line.Err = true
		} else {
			line.Out = evalText(v)
		}
		*history = append(*history, line)
		root.SetChanged()
		sound.Play(sound.Send)
		kit.Log("REPL", "ran", "runes", utf8.RuneCountInString(src))
	}
	if *pending != "" {
		run(*pending)
		*pending = ""
	}
	if send.Clicked() && strings.TrimSpace(*draft) != "" {
		run(strings.TrimSpace(*draft))
		*draft = ""
	}
	return root, send
}

// EnvVar is one environment variable.
