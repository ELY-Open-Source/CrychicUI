package agent

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// StreamCodeOptions name the copy button.
type StreamCodeOptions struct {
	Copy string
}

// StreamCode is a code block filling as it streams.
func StreamCode(c *ui.Context, lang, code string, done bool, o StreamCodeOptions) (root, copy ui.Element) {
	if lang == "" || code == "" {
		panic(fmt.Sprintf("agent: StreamCode %q needs a language and code", lang))
	}
	def(&o.Copy, "Copy")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Radius(t.Radius).Background(t.Surface).Border(1, t.Border).Clip().Role(ui.RoleGroup).Label(lang)
	shown := code
	streaming := !done && !c.Preferences().ReduceMotion
	if streaming {
		since := ui.Local(root, "since", func() time.Time { return c.Now() })
		const step = 60 * time.Millisecond
		n := int(c.Now().Sub(*since) / step)
		shown = prefixLines(code, n)
		if len(shown) < len(code) {
			c.After(step)
		} else {
			c.After(time.Second)
		}
	}
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(1), t.Space(1), t.Space(3)).
			Background(t.SurfaceHover).Children(func() {
			ui.Icon(c, icon.Must("file-code")).Size(16, 16).TextColor(t.TextMuted)
			kit.Muted(c, lang).FontSize(t.Rem(0.75)).Bold().Grow(1).SingleLine()
			copy = button.CopyButton(c, o.Copy, code).Disabled(streaming && len(shown) < len(code))
		})
		lines := strings.Split(strings.TrimSuffix(shown, "\n"), "\n")
		ui.Scroll(c).FillWidth().MaxHeight(t.Space(80)).Padding(t.Space(2), t.Space(3)).Children(func() {
			for i, l := range lines {
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(3))
				kit.Arrive(row, ui.Pose{X: -t.Space(1), Fade: 1}, 0).Children(func() {
					kit.Numeric(c, strconv.Itoa(i+1)).Font("monospace").FontSize(t.Rem(0.75)).TextColor(t.TextMuted).Width(t.Space(5)).Shrink(0)
					ui.Text(c, strings.ReplaceAll(l, "	", "    ")).Font("monospace").FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).Wrap()
					if streaming && i == len(lines)-1 {
						caret := ui.Box(c).Width(t.Space(0.5)).Height(t.Rem(0.9)).Shrink(0).Background(t.Accent)
						ph := float64(caret.Loop("blink", time.Second, ui.Linear))
						caret.Opacity(unit(ph < 0.6))
					}
				})
			}
		})
	})
	if copy.Clicked() {
		kit.Log("StreamCode", "copied", "lines", strings.Count(strings.TrimSuffix(code, "\n"), "\n")+1)
	}
	return root, copy
}
