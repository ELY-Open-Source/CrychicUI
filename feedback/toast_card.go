package feedback

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ToastCardOptions tune a ToastCard.
type ToastCardOptions struct {
	// Action labels an optional action button.
	Action string
	// Close names the close button; "Close" when empty.
	Close string
}

// toastFace styles a toast card's surface.
func toastFace(c *ui.Context, e ui.Element) ui.Element {
	t := c.Theme()
	return e.Row().AlignItems(ui.Start).Gap(t.Space(3)).Width(t.Space(90)).Padding(t.Space(3)).Radius(t.Radius+2).
		Background(t.Background).Border(1, t.Border).Shadow(0, 6, 20, 0, kit.Shade(0.18))
}

// toastBody is a toast's icon, title and description.
func toastBody(c *ui.Context, tone Tone, title, body string) {
	t := c.Theme()
	col, glyph := look(t, tone)
	ui.Icon(c, icon.Must(glyph)).Size(16, 16).TextColor(col).Margin(t.Space(0.5), 0, 0, 0)
	ui.Column(c).Grow(1).Gap(t.Space(1)).Children(func() {
		ui.Text(c, title).Bold()
		if body != "" {
			kit.Muted(c, body)
		}
	})
}

// ToastCard is one toast: icon, title, body, action and close.
func ToastCard(c *ui.Context, tone Tone, title, body string, o ToastCardOptions) (root, act, dismiss ui.Element) {
	if title == "" {
		panic("feedback: ToastCard needs a title")
	}
	if o.Close == "" {
		o.Close = "Close"
	}
	t := c.Theme()
	root = toastFace(c, ui.Row(c)).Role(ui.RoleStatus).Label(title)
	kit.Arrive(root, ui.Pose{X: t.Space(6), Fade: 1}, 0)
	root.Children(func() {
		toastBody(c, tone, title, body)
		ui.Column(c).Gap(t.Space(2)).AlignItems(ui.End).Children(func() {
			dismiss = closer(c, o.Close)
			if kit.Clicked(dismiss, sound.Close) {
				kit.Log("ToastCard", "closed")
			}
			if o.Action != "" {
				act = button(c, o.Action, false)
				if kit.Clicked(act, sound.Tap) {
					kit.Log("ToastCard", "action")
				}
			}
		})
	})
	return root, act, dismiss
}
