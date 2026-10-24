package code

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// REPLLine is one evaluated input and its result.
type REPLLine struct {
	In, Out string
	Err     bool
}

// REPLOptions are the copy knobs of REPL.
type REPLOptions struct {
	Placeholder, Run string
	FormatError      func(err error) string
	// Height caps the history; zero gives the default.
	Height float32
}

// REPL evaluates arithmetic in a terminal with history; cue Send.
func REPL(c *ui.Context, label string, history *[]REPLLine, draft *string, o REPLOptions) (root, send ui.Element) {
	if label == "" || history == nil || draft == nil {
		panic(fmt.Sprintf("code: REPL %q needs a label and bound values", label))
	}
	def(&o.Placeholder, "1 + 2 * 3")
	def(&o.Run, "Run")
	if o.FormatError == nil {
		o.FormatError = func(err error) string { return "error: " + err.Error() }
	}
	size("REPL", &o.Height, t40(c))
	t := c.Theme()
	bg, fg := inversePair(t)
	root = ui.Column(c).FillWidth().Radius(t.Radius).Clip().Background(bg).TextColor(fg).
		Shadow(0, 1, 3, 0, kit.Shade(0.06)).Role(ui.RoleGroup).Label(label)
	pending := ui.Local(root, "pending", func() string { return "" })
	root.Children(func() {
		termChrome(c, label, fg)
		c.Key("history")
		ui.Scroll(c).FillWidth().MaxHeight(o.Height).Padding(t.Space(2), 0).Children(func() {
			for i, l := range *history {
				c.Key(i)
				entry := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(0.5), t.Space(3))
				stagger(c, root, entry, i).Children(func() {
					ui.Text(c, "> "+l.In).Font("monospace").FontSize(t.Rem(0.8)).TextColor(fg.Alpha(0.7)).Wrap()
					out := ui.Text(c, "= "+l.Out).Font("monospace").FontSize(t.Rem(0.8)).Bold().Wrap()
					if l.Err {
						out.TextColor(t.Danger)
					} else {
						out.TextColor(t.Success)
					}
				})
			}
		})
		c.Key("prompt")
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).
			BorderWidth(1, 0, 0, 0).BorderColor(fg.Alpha(0.12)).Children(func() {
			ui.Text(c, "›").Font("monospace").FontSize(t.Rem(0.95)).Bold().TextColor(t.Accent)
			in := ui.TextInput(c, draft).Label(o.Placeholder).Placeholder(o.Placeholder).Font("monospace").FontSize(t.Rem(0.85)).Grow(1)
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
			send = kit.Press(ui.ButtonBase(c).Label(o.Run).Tooltip(o.Run).Size(t.Space(7), t.Space(7)).Shrink(0).Center().Radius(999)).
				Disabled(strings.TrimSpace(*draft) == "")
			send.Background(send.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(send.Hovered())), ui.Snappy)).TextColor(t.AccentText).
				Children(func() { ui.Icon(c, icon.Must("corner-down-left")).Size(14, 14) })
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
		kit.Log("REPL", "ran", "runes", utf8.RuneCountInString(src), "error", line.Err)
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
