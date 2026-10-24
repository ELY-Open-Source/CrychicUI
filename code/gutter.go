package code

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// GutterOptions word each line's toggle, "%d" taking the line.
type GutterOptions struct {
	Toggle string
}

// Gutter shows code whose line numbers toggle breakpoints.
func Gutter(c *ui.Context, lang, code string, breakpoints *[]int, o GutterOptions) ui.Element {
	if code == "" || breakpoints == nil {
		panic("code: Gutter needs code and breakpoints")
	}
	lines := linesOf(code)
	for _, b := range *breakpoints {
		if b < 1 || b > len(lines) {
			panic(fmt.Sprintf("code: Gutter breakpoints must be lines from 1 to %d", len(lines)))
		}
	}
	def(&o.Toggle, "Breakpoint on line %d")
	t := c.Theme()
	digits := float32(len(strconv.Itoa(len(lines))))
	root := ui.Column(c).FillWidth().Clip().Radius(t.Radius).Background(t.Surface).Border(1, t.Border).Padding(t.Space(2), 0).
		Role(ui.RoleGroup).Label(lang)
	kit.Arrive(root, ui.Pose{Fade: 1}, 0)
	h := newHighlighter(c, lang)
	state := uint64(0)
	root.Children(func() {
		for i, line := range lines {
			n := i + 1
			c.Key(n)
			var spans []ui.Span
			spans, state = h.line(state, line)
			on := slices.Contains(*breakpoints, n)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(0, t.Space(3), 0, t.Space(1))
			row.Background(row.AnimateColor("bg", t.Danger.Alpha(0.1*unit(on)), ui.Snappy))
			row.Children(func() {
				b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Toggle, n)).Checked(on).Width(t.Space(2*digits+6)).Shrink(0).
					Padding(0, t.Space(1)).Radius(t.Radius))
				b.Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
						dot := ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Shrink(0).Radius(999)
						k := dot.AnimateTo("on", unit(on), ui.Bouncy)
						dot.Background(dot.AnimateColor("ink", t.Danger.Alpha(max(unit(on), 0.35*unit(b.Hovered()))), ui.Snappy)).Scale(0.6+0.4*k, 0.6+0.4*k)
						num := kit.Numeric(c, strconv.Itoa(n)).Font("monospace").FontSize(t.Rem(0.8)).Grow(1).TextAlign(ui.End)
						num.TextColor(num.AnimateColor("ink", t.TextMuted.Alpha(0.7).Mix(t.Danger, unit(on)), ui.Snappy))
					})
				})
				ui.RichText(c, spans...).Font("monospace").FontSize(t.Rem(0.8)).LineHeight(1.6).Grow(1).MinWidth(0).Wrap()
				if b.Clicked() {
					if on {
						*breakpoints = slices.DeleteFunc(slices.Clone(*breakpoints), func(v int) bool { return v == n })
					} else {
						*breakpoints = append(slices.Clone(*breakpoints), n)
					}
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("Gutter", "breakpoints", "count", len(*breakpoints))
				}
			})
		}
	})
	return root
}
