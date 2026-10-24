package code

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// StackTraceOptions word the toggle, "%d" taking the count.
type StackTraceOptions struct {
	More, Less string
}

var frameFold = ui.ElementTransition{Animation: &ui.Snappy, Enter: &ui.Motion{Collapse: true}, Exit: &ui.Motion{Collapse: true}}

// StackTrace shows an error over frames; open shows library frames.
func StackTrace(c *ui.Context, message string, list []StackFrame, open *bool, o StackTraceOptions) ui.Element {
	if message == "" || len(list) == 0 || open == nil {
		panic(fmt.Sprintf("code: StackTrace %q needs a message, frames and open", message))
	}
	frames("StackTrace", list)
	def(&o.More, "Show %d library frames")
	def(&o.Less, "Hide library frames")
	// Closed, app frames show; without any, the first does.
	app := slices.ContainsFunc(list, func(f StackFrame) bool { return f.App })
	closed := func(i int) bool { return list[i].App || !app && i == 0 }
	hidden := 0
	for i := range list {
		if !closed(i) {
			hidden++
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Clip().Radius(t.Radius).Background(t.Surface).Border(1, t.Border).
		Role(ui.RoleGroup).Label(message)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	var toggle ui.Element
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(2.5), t.Space(3)).Background(t.Danger.Alpha(0.08)).
			BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			ui.Icon(c, icon.Must("circle-x")).Size(16, 16).Shrink(0).TextColor(t.Danger)
			ui.Text(c, message).Font("monospace").FontSize(t.Rem(0.8)).Bold().TextColor(t.Danger).Grow(1).MinWidth(0).Wrap()
		})
		ui.Column(c).FillWidth().Padding(t.Space(1), 0).Children(func() {
			for i, f := range list {
				if !*open && !closed(i) {
					continue
				}
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(3)).Transition(frameFold)
				kit.Stagger(c, root, row, i).Children(func() { frameRow(c, f) })
			}
		})
		if hidden > 0 {
			label := fmt.Sprintf(o.More, hidden)
			if *open {
				label = o.Less
			}
			toggle = kit.Press(ui.ButtonBase(c).Label(label).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(1.5), t.Space(3)).
				BorderWidth(1, 0, 0, 0).BorderColor(t.Border))
			toggle.Background(toggle.AnimateColor("bg", t.SurfaceHover.Alpha(unit(toggle.Hovered())), ui.Snappy)).Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					ch := ui.Icon(c, icon.Must("chevron-right")).Size(14, 14).TextColor(t.TextMuted)
					ch.Rotate(90 * ch.AnimateTo("turn", unit(*open), ui.Snappy))
					kit.Muted(c, label).FontSize(t.Rem(0.8)).Shrink(0)
				})
			})
		}
	})
	if hidden > 0 && toggle.Clicked() {
		*open = !*open
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("StackTrace", "open", *open, "frames", len(list))
	}
	return root
}

// frameRow shows a frame: package, name, file and line.
func frameRow(c *ui.Context, f StackFrame) {
	t := c.Theme()
	dot := t.Border
	if f.App {
		dot = t.Accent
	}
	ui.Box(c).Size(t.Space(1.5), t.Space(1.5)).Shrink(0).Radius(999).Background(dot)
	pkg, name := splitFunc(f.Func)
	ui.Row(c).Grow(1).MinWidth(0).Children(func() {
		kit.Muted(c, pkg).Font("monospace").FontSize(t.Rem(0.8)).Shrink(0)
		fn := ui.Text(c, name).Font("monospace").FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).SingleLine()
		if f.App {
			fn.Bold()
		} else {
			fn.TextColor(t.TextMuted)
		}
	})
	kit.Muted(c, fmt.Sprintf("%s:%d", f.File, f.Line)).Font("monospace").FontSize(t.Rem(0.75)).Shrink(0)
}

// splitFunc cuts "net/http.(*conn).serve" after its package.
func splitFunc(fn string) (pkg, name string) {
	slash := strings.LastIndex(fn, "/") + 1
	if dot := strings.Index(fn[slash:], "."); dot >= 0 {
		return fn[:slash+dot+1], fn[slash+dot+1:]
	}
	return "", fn
}
