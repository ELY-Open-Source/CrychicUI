package devtools

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// QueryOptions word the editor; Line and Lines take a count.
type QueryOptions struct{ Title, Database, Placeholder, Run, Line, Lines string }

// Query edits SQL; blank SQL never runs; cue Tap.
func Query(c *ui.Context, sql *string, o QueryOptions) (root, run ui.Element) {
	def(&o.Title, "Query")
	if sql == nil {
		panic(fmt.Sprintf("devtools: Query %q needs SQL state", o.Title))
	}
	def(&o.Placeholder, "SELECT * FROM orders LIMIT 10;")
	def(&o.Run, "Run")
	def(&o.Line, "%d line")
	def(&o.Lines, "%d lines")
	t := c.Theme()
	root = card(c).Role(ui.RoleGroup)
	root.Children(func() {
		kit.CardHead(c, "square-terminal", o.Title, func() {
			if o.Database != "" {
				kit.Chip(c, o.Database, t.Accent)
			}
		})
		frame := ui.Column(c).FillWidth().Padding(t.Space(2.5), t.Space(3)).Background(t.SurfaceHover.Alpha(0.5))
		frame.BorderWidth(0, 0, 1, 0).BorderColor(frame.AnimateColor("ring", t.Border.Mix(t.Accent, unit(frame.FocusWithin())), ui.Snappy))
		frame.Children(func() {
			in := ui.TextAreaBase(c, sql).Label(o.Title).Placeholder(o.Placeholder).Font("monospace").FontSize(t.Rem(0.8)).
				FillWidth().Lines(4, 10)
			if in.Changed() {
				root.SetChanged()
				kit.Log("Query", "edited", "runes", utf8.RuneCountInString(*sql))
			}
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Children(func() {
			n := strings.Count(*sql, "\n") + 1
			format := o.Lines
			if n == 1 {
				format = o.Line
			}
			kit.Muted(c, fmt.Sprintf(format, n)).FontSize(t.Rem(0.72)).Grow(1).MinWidth(0).SingleLine()
			blank := strings.TrimSpace(*sql) == ""
			run = kit.Press(ui.ButtonBase(c).Label(o.Run).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).Shrink(0).Disabled(blank))
			run.Background(run.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(run.Hovered())), ui.Snappy)).TextColor(t.AccentText)
			run.Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					glyph(c, "play", 14)
					ui.Text(c, o.Run).FontSize(t.Rem(0.8)).Bold().Shrink(0)
				})
			})
			if run.Clicked() && !blank {
				sound.Play(sound.Tap)
				kit.Log("Query", "run", "runes", utf8.RuneCountInString(*sql))
			}
		})
	})
	return root, run
}
