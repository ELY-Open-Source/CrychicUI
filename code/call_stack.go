package code

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CallStack shows frames with the current one marked.
func CallStack(c *ui.Context, label string, list []StackFrame, current *int) ui.Element {
	if label == "" || len(list) == 0 || current == nil || *current < 0 || *current >= len(list) {
		panic(fmt.Sprintf("code: CallStack %q needs frames and a current one, got %d", label, deref(current)))
	}
	frames("CallStack", list)
	t := c.Theme()
	root := kit.Card(c).FillWidth().Padding(t.Space(1)).Role(ui.RoleList).Label(label)
	return root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2)).Children(func() {
			ui.Icon(c, icon.Must("layers")).Size(14, 14).TextColor(t.TextMuted)
			ui.Text(c, label).FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0).SingleLine()
			kit.Numeric(c, fmt.Sprint(len(list))).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
		})
		ui.Column(c).FillWidth().FocusGroup(ui.Vertical).Children(func() {
			seen := map[string]int{}
			for i, f := range list {
				c.Key(ident(seen, fmt.Sprintf("%q %q %d", f.Func, f.File, f.Line)))
				on := i == *current
				b := kit.Press(ui.ButtonBase(c).Label(f.Func).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(1.25), t.Space(2)).
					Radius(t.Radius).Checked(on))
				kit.Stagger(c, root, b, i)
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(t.Warning.Alpha(0.14), unit(on)), ui.Snappy))
				b.Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						mark := ui.Icon(c, icon.Must("play")).Size(12, 12).Shrink(0).TextColor(t.Warning)
						mark.Opacity(mark.AnimateTo("on", unit(on), ui.Snappy))
						pkg, name := splitFunc(f.Func)
						ui.Row(c).Grow(1).MinWidth(0).Children(func() {
							kit.Muted(c, pkg).Font("monospace").FontSize(t.Rem(0.8)).Shrink(0)
							fn := ui.Text(c, name).Font("monospace").FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).SingleLine()
							if f.App || on {
								fn.Bold()
							} else {
								fn.TextColor(t.TextMuted)
							}
						})
						kit.Muted(c, fmt.Sprintf("%s:%d", f.File, f.Line)).Font("monospace").FontSize(t.Rem(0.75)).Shrink(0)
					})
				})
				if b.Clicked() && !on {
					*current = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("CallStack", "selected", "index", i)
				}
			}
		})
	})
}
