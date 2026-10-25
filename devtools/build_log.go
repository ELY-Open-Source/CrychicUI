package devtools

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BuildStep is one build stage with its output.
type BuildStep struct {
	Name     string
	State    StepState
	Duration time.Duration
	Lines    []string
}

// BuildLogOptions word the panel.
type BuildLogOptions struct{ Title, Running, Passed, Failed string }

// BuildLog shows timed build steps; folds report changes; cue Toggle.
func BuildLog(c *ui.Context, steps []BuildStep, o BuildLogOptions) ui.Element {
	def(&o.Title, "Build")
	if len(steps) == 0 {
		panic(fmt.Sprintf("devtools: BuildLog %q needs steps", o.Title))
	}
	names := make([]string, len(steps))
	var total time.Duration
	failed, busy := false, false
	for i, s := range steps {
		if s.State < StepPending || s.State > StepFailed || s.Duration < 0 {
			panic(fmt.Sprintf("devtools: BuildLog %q needs a state and a non-negative time", s.Name))
		}
		if s.Duration > math.MaxInt64-total {
			panic(fmt.Sprintf("devtools: BuildLog %q total time overflows", o.Title))
		}
		total += s.Duration
		failed = failed || s.State == StepFailed
		busy = busy || s.State == StepRunning || s.State == StepPending
		names[i] = s.Name
	}
	unique("BuildLog", names)
	def(&o.Running, "Running")
	def(&o.Passed, "Passed")
	def(&o.Failed, "Failed")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	open := ui.Local(root, "open", func() map[string]bool { return map[string]bool{} })
	root.Children(func() {
		kit.CardHead(c, "hammer", o.Title, func() {
			kit.Muted(c, took(total)).Font("monospace").FontSize(t.Rem(0.7)).Shrink(0)
			switch {
			case failed:
				kit.Chip(c, o.Failed, t.Danger)
			case busy:
				kit.Chip(c, o.Running, t.Accent)
			default:
				kit.Chip(c, o.Passed, t.Success)
			}
		})
		list := ui.Column(c).FillWidth().Padding(t.Space(1.5))
		list.Children(func() {
			for i, s := range steps {
				shown, set := (*open)[s.Name]
				if !set {
					shown = s.State == StepFailed
				}
				shown = shown && len(s.Lines) > 0
				c.Key(s.Name)
				step := ui.Column(c).FillWidth()
				kit.Stagger(c, list, step, i).Children(func() {
					head := func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							chev := glyph(c, "chevron-right", 14).TextColor(t.TextMuted).Opacity(unit(len(s.Lines) > 0))
							chev.Rotate(90 * chev.AnimateTo("turn", unit(shown), ui.Snappy))
							stepMark(c, s.State)
							ink := t.Text
							if s.State == StepFailed {
								ink = t.Danger
							}
							ui.Text(c, s.Name).FontSize(t.Rem(0.8)).Bold().TextColor(ink).Grow(1).MinWidth(0).SingleLine()
							if s.State != StepPending {
								kit.Numeric(c, took(s.Duration)).Font("monospace").FontSize(t.Rem(0.7)).TextColor(t.TextMuted).Shrink(0)
							}
						})
					}
					if len(s.Lines) == 0 {
						ui.Column(c).FillWidth().Padding(t.Space(2)).Children(head)
						return
					}
					b := kit.Press(ui.ButtonBase(c).Label(s.Name).FillWidth().Padding(t.Space(2)).Radius(t.Radius).Expanded(shown))
					b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
					b.Children(head)
					if b.Clicked() {
						(*open)[s.Name] = !shown
						root.SetChanged()
						sound.Play(sound.Toggle)
						kit.Log("BuildLog", "toggled", "index", i, "open", !shown)
					}
					layout.Collapse(c, shown, func() { buildLines(c, s.Lines) })
				})
			}
		})
	})
	return root
}

// buildLines is numbered build output.
func buildLines(c *ui.Context, lines []string) {
	t := c.Theme()
	// Padding indents; a margin beside FillWidth overflows.
	ui.Column(c).FillWidth().Padding(0, 0, t.Space(1.5), t.Space(7)).Children(func() {
		ui.Column(c).FillWidth().Padding(t.Space(2), t.Space(2.5)).Radius(t.Radius).Background(t.SurfaceHover.Alpha(0.7)).Children(func() {
			for i, l := range lines {
				ui.Row(c).FillWidth().Gap(t.Space(2.5)).Children(func() {
					kit.Muted(c, strconv.Itoa(i+1)).Font("monospace").FontSize(t.Rem(0.68)).Width(t.Space(6)).Shrink(0).TextAlign(ui.End)
					ui.Text(c, l).Font("monospace").FontSize(t.Rem(0.7)).Grow(1).MinWidth(0).Wrap()
				})
			}
		})
	})
}
