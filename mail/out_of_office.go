package mail

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AutoReply is one automatic reply configuration.
type AutoReply struct {
	Enabled    bool
	Start, End time.Time
	Message    string
}

// OutOfOfficeOptions word the switch, fields and status line.
type OutOfOfficeOptions struct{ Enabled, Start, End, Message, Range, Off, DateFormat string }

// OutOfOffice edits automatic reply dates and message; cue Toggle.
func OutOfOffice(c *ui.Context, value *AutoReply, o OutOfOfficeOptions) ui.Element {
	if value == nil || value.Start.IsZero() || value.End.IsZero() || value.End.Before(value.Start) {
		panic("mail: OutOfOffice needs a valid range")
	}
	def(&o.Enabled, "Out of office")
	def(&o.Start, "Start")
	def(&o.End, "End")
	def(&o.Message, "Message")
	def(&o.Range, "Replying %s – %s")
	def(&o.Off, "Automatic replies are off")
	def(&o.DateFormat, "Jan 2")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		head := kit.Press(ui.ButtonBase(c).Label(o.Enabled).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(3)).Checked(value.Enabled))
		head.Background(head.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, 0.6*unit(head.Hovered())), ui.Snappy))
		head.Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
				ink := t.TextMuted
				if value.Enabled {
					ink = t.Warning
				}
				tile(c, "plane-takeoff", ink)
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
					ui.Text(c, o.Enabled).FontSize(t.Rem(0.9)).Bold().SingleLine()
					status := o.Off
					if value.Enabled {
						status = fmt.Sprintf(o.Range, value.Start.Format(o.DateFormat), value.End.Format(o.DateFormat))
					}
					kit.Muted(c, status).FontSize(t.Rem(0.75)).FillWidth().MinWidth(0).SingleLine().ContentTransition(ui.ContentTransition{})
				})
				track(c, ui.Box(c).Size(t.Space(9), t.Space(5)).Shrink(0).Radius(999), value.Enabled, head.Hovered())
			})
		})
		if head.Clicked() {
			value.Enabled = !value.Enabled
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("OutOfOffice", "toggled", "enabled", value.Enabled)
		}
		layout.Collapse(c, value.Enabled, func() {
			ui.Column(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(3)).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					start := ui.DateInput(c, &value.Start).Label(o.Start).Grow(1).Basis(0).MinWidth(0)
					glyph(c, "arrow-right", 14).TextColor(t.TextMuted)
					end := ui.DateInput(c, &value.End).Label(o.End).Grow(1).Basis(0).MinWidth(0)
					if start.Changed() || end.Changed() {
						if value.End.Before(value.Start) {
							if start.Changed() {
								value.End = value.Start
							} else {
								value.Start = value.End
							}
						}
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("OutOfOffice", "dates")
					}
				})
				in := ui.TextArea(c, &value.Message).Label(o.Message).Placeholder(o.Message).FontSize(t.Rem(0.85)).FillWidth().Lines(3, 6)
				if in.Changed() {
					root.SetChanged()
					sound.Play(sound.Tick)
					kit.Log("OutOfOffice", "edited", "runes", utf8.RuneCountInString(value.Message))
				}
			})
		})
	})
	return root
}
