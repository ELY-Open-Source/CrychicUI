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

// Command is one ran command and its output.
type Command struct {
	Text, Output string
}

// CommandLineOptions tune a CommandLine.
type CommandLineOptions struct {
	// Height of the output; zero gives the default.
	Height      float32
	Prompt      string
	Placeholder string
	Run         string
}

// CommandLine is a prompt input with history and output rows.
func CommandLine(c *ui.Context, label string, history []Command, draft *string, onRun func(cmd string), o CommandLineOptions) (root, send ui.Element) {
	if label == "" || draft == nil || onRun == nil {
		panic(fmt.Sprintf("code: CommandLine %q needs a label, a draft and a run func", label))
	}
	for _, k := range history {
		if k.Text == "" {
			panic(fmt.Sprintf("code: CommandLine %q history needs commands", label))
		}
	}
	def(&o.Prompt, "$")
	def(&o.Placeholder, "Run a command")
	def(&o.Run, "Run")
	t := c.Theme()
	size("CommandLine", &o.Height, t.Space(48))
	bg, fg := inversePair(t)
	root = ui.Column(c).FillWidth().Radius(t.Radius).Clip().Background(bg).TextColor(fg).
		Shadow(0, 1, 3, 0, kit.Shade(0.06)).Role(ui.RoleGroup).Label(label)
	pending := ui.Local(root, "pending", func() string { return "" })
	recall := ui.Local(root, "recall", func() int { return -1 })
	if *recall >= len(history) {
		*recall = len(history) - 1
	}
	moved := ui.Local(root, "moved", func() bool { return false })
	blank := strings.TrimSpace(*draft) == ""
	root.Children(func() {
		termChrome(c, label, fg)
		if len(history) > 0 {
			c.Key("history")
			ui.Scroll(c).FillWidth().MaxHeight(o.Height).Padding(t.Space(2), 0).Children(func() {
				for i, k := range history {
					c.Key(i)
					entry := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(0.5), t.Space(3))
					kit.Stagger(c, root, entry, i).Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(1.5)).Children(func() { termLine(c, o.Prompt, o.Prompt+" "+k.Text, fg) })
						if k.Output != "" {
							for j, out := range linesOf(k.Output) {
								c.Key(j)
								ui.Row(c).FillWidth().Children(func() { termLine(c, o.Prompt, out, fg) })
							}
						}
					})
				}
			})
		}
		// Keyed: the history appearing must not move the input.
		c.Key("prompt")
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).
			BorderWidth(1, 0, 0, 0).BorderColor(fg.Alpha(0.12)).Children(func() {
			ui.Text(c, o.Prompt).Font("monospace").FontSize(t.Rem(0.85)).Bold().TextColor(t.Success)
			in := ui.TextInput(c, draft).Grow(1).Placeholder(o.Placeholder).Label(o.Placeholder).Font("monospace").FontSize(t.Rem(0.85))
			if in.Changed() {
				root.SetChanged()
			}
			in.OnTextKey(func(k ui.TextKey) bool {
				switch {
				case k.Key == ui.KeyEnter && k.Mods == 0:
					if text := strings.TrimSpace(k.Text()); text != "" {
						*pending = text
						k.SetText("")
					}
					return true
				case k.Key == ui.KeyUp && k.Mods == 0 && len(history) > 0:
					if *recall < len(history)-1 {
						*recall++
						*moved = true
						k.SetText(history[len(history)-1-*recall].Text)
					}
					return true
				case k.Key == ui.KeyDown && k.Mods == 0:
					if *recall > 0 {
						*recall--
						*moved = true
						k.SetText(history[len(history)-1-*recall].Text)
					} else if *recall == 0 {
						*recall = -1
						*moved = true
						k.SetText("")
					}
					return true
				}
				return false
			})
			send = kit.Press(ui.ButtonBase(c).Label(o.Run).Tooltip(o.Run).Size(t.Space(7), t.Space(7)).Shrink(0).Center().Radius(999)).Disabled(blank)
			send.Background(send.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(send.Hovered())), ui.Snappy)).
				TextColor(t.AccentText).Children(func() { ui.Icon(c, icon.Must("play")).Size(16, 16) })
		})
	})
	if *moved {
		*moved = false
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("CommandLine", "recalled")
	}
	text := *pending
	*pending = ""
	if text == "" && !blank && send.Clicked() {
		text = strings.TrimSpace(*draft)
		*draft = ""
	}
	if text != "" {
		*recall = -1
		onRun(text)
		root.SetChanged()
		sound.Play(sound.Send)
		kit.Log("CommandLine", "ran", "runes", utf8.RuneCountInString(text))
	}
	return root, send
}
