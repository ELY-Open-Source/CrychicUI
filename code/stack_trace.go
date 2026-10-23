package code

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// StackTraceOptions word the toggle.
type StackTraceOptions struct {
	// More and Less word the toggle, "%d" taking the count.
	More, Less string
}

// StackTrace shows frames; app frames are bold.
func StackTrace(c *ui.Context, list []StackFrame, open *bool, o StackTraceOptions) ui.Element {
	if len(list) == 0 || open == nil {
		panic("code: StackTrace needs frames and open")
	}
	frames("StackTrace", list)
	def(&o.More, "%d more frames")
	def(&o.Less, "Show less")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Clip().Radius(t.Radius).Background(t.Surface).Border(1, t.Border).
		Role(ui.RoleGroup).Label(list[0].Func)
	var toggle ui.Element
	expand := len(list) > 1
	root.Children(func() {
		shown := 1
		if *open {
			shown = len(list)
		}
		for i := 0; i < shown; i++ {
			c.Key(i)
			f := list[i]
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(1), t.Space(3))
			row.Children(func() {
				name := ui.Text(c, f.Func).Font("monospace").FontSize(t.Rem(0.85)).Grow(1).SingleLine()
				if f.App {
					name.Bold()
				} else {
					name.TextColor(t.TextMuted)
				}
				kit.Muted(c, fmt.Sprintf("%s:%d", f.File, f.Line)).Font("monospace").FontSize(t.Rem(0.8)).SingleLine()
			})
		}
		if expand {
			label := fmt.Sprintf(o.More, len(list)-1)
			if *open {
				label = o.Less
			}
			toggle = kit.Press(ui.ButtonBase(c).Label(label).Row().FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).
				Padding(t.Space(1), t.Space(3)))
			toggle.Background(toggle.AnimateColor("bg", t.SurfaceHover.Alpha(unit(toggle.Hovered())), ui.Snappy)).Children(func() {
				ch := ui.Icon(c, icon.Must("chevron-right")).Size(16, 16).TextColor(t.TextMuted)
				ch.Rotate(90 * ch.AnimateTo("turn", unit(*open), ui.Snappy))
				kit.Muted(c, label).FontSize(t.Rem(0.8)).SingleLine()
			})
		}
	})
	if expand && toggle.Clicked() {
		*open = !*open
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("StackTrace", "open", *open)
	}
	return root
}
