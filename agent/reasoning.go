package agent

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ReasoningOptions tune a Reasoning panel.
type ReasoningOptions struct {
	// Took is the thinking time; zero shows Title alone.
	Took time.Duration
	// Title is the plain header; Format words it with Took.
	Title, Format string
}

// Reasoning is a collapsible reasoning panel; height springs.
func Reasoning(c *ui.Context, body string, open *bool, o ReasoningOptions) ui.Element {
	if body == "" || open == nil || o.Took < 0 {
		panic("agent: Reasoning needs a body, open and a time not negative")
	}
	def(&o.Title, "Reasoning")
	def(&o.Format, "Thought for %ds")
	title := o.Title
	if o.Took > 0 {
		title = fmt.Sprintf(o.Format, int64((o.Took+time.Second/2)/time.Second))
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	var head ui.Element
	root.Children(func() {
		head = kit.Press(ui.ButtonBase(c).Label(title).Row().AlignSelf(ui.Start).Gap(t.Space(1.5)).AlignItems(ui.Center).
			Padding(t.Space(1), t.Space(2)).Radius(t.Radius).Checked(*open))
		head.Background(head.AnimateColor("bg", t.SurfaceHover.Alpha(unit(head.Hovered())), ui.Snappy))
		head.TextColor(head.AnimateColor("ink", t.TextMuted.Mix(t.Text, unit(head.Hovered())), ui.Snappy))
		head.Children(func() {
			ui.Icon(c, icon.Must("brain")).Size(16, 16)
			ui.Text(c, title).FontSize(t.Rem(0.9)).Shrink(0)
			ch := ui.Icon(c, icon.Must("chevron-down")).Size(16, 16)
			ch.Rotate(-90 + 90*ch.AnimateTo("turn", unit(*open), ui.Snappy))
		})
		layout.Collapse(c, *open, func() {
			ui.Row(c).FillWidth().Gap(t.Space(3)).PaddingStart(t.Space(3)).Children(func() {
				ui.Box(c).Width(2).Radius(1).Background(t.Border)
				kit.Muted(c, body).FontSize(t.Rem(0.9)).Grow(1).Wrap()
			})
		})
	})
	if head.Clicked() {
		*open = !*open
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("Reasoning", "open", *open)
	}
	return root
}
