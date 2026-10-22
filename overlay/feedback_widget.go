package overlay

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FeedbackWidgetOptions tune a FeedbackWidget.
type FeedbackWidgetOptions struct {
	Label, Title, Placeholder, Send string
}

// FeedbackWidget is a round button opening a form.
func FeedbackWidget(c *ui.Context, text *string, open *bool, o FeedbackWidgetOptions) (root, send ui.Element) {
	if text == nil || open == nil {
		panic("overlay: FeedbackWidget needs text and open")
	}
	def(&o.Label, "Feedback")
	def(&o.Title, "Send feedback")
	def(&o.Placeholder, "What can we improve?")
	def(&o.Send, "Send")
	t := c.Theme()
	root = ui.Box(c).AlignSelf(ui.End)
	field := opening(c, root, *open)
	root.Children(func() {
		b := kit.Press(ui.ButtonBase(c).Label(o.Label).Tooltip(o.Label).Size(t.Space(12), t.Space(12)).Center().Radius(999).
			Shadow(0, 4, 12, 0, kit.Shade(0.2))).Expanded(*open)
		b.Background(b.AnimateColor("bg", t.Accent.Mix(t.AccentText, 0.12*unit(b.Hovered())), ui.Snappy)).TextColor(t.AccentText).
			Children(func() {
				glyph := "message-square"
				if *open {
					glyph = "x"
				}
				ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{})
			})
		if b.Clicked() {
			*open = !*open
			root.SetChanged()
		}
		if ui.PopoverBase(c, b, open, func(panel ui.Element) {
			panel.Margin(0, 0, t.Space(2), 0).Width(t.Space(80)).Padding(t.Space(4)).Gap(t.Space(3)).Radius(t.Radius*1.5).
				Background(t.Background).Border(1, t.Border).Shadow(0, 12, 32, 0, kit.Shade(0.24)).TransformOrigin(1, 1)
			kit.Arrive(panel, ui.Pose{Scale: 0.92, Fade: 1}, 0)
			ui.Text(c, o.Title).Bold()
			area := ui.TextArea(c, text).Label(o.Title).Placeholder(o.Placeholder).FillWidth().Lines(3, 8).Bind(field)
			if area.Changed() {
				root.SetChanged()
			}
			ui.Row(c).FillWidth().Justify(ui.End).Children(func() {
				send = primary(c, o.Send, t.Accent, t.AccentText).Disabled(strings.TrimSpace(*text) == "")
				if send.Clicked() {
					*open = false
					root.SetChanged()
					sound.Play(sound.Send)
					kit.Log("FeedbackWidget", "sent", "length", len(*text))
				}
			})
		}).Dismissed() {
			root.SetChanged()
		}
	})
	if opened, _ := moved(root, *open, "FeedbackWidget"); opened {
		sound.Play(sound.Open)
	}
	return root, send
}
