package feedback

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ToasterOptions tune a Toaster.
type ToasterOptions struct {
	// Close names each close button; "Close" when empty.
	Close string
}

var toastMotion = ui.ElementTransition{Animation: &ui.Snappy, Exit: &ui.Motion{X: 40}}

// toastTone reads a toast's Type as a tone.
func toastTone(kind string) Tone {
	switch kind {
	case "", "info":
		return ToneInfo
	case "success":
		return ToneSuccess
	case "warning":
		return ToneWarning
	case "error":
		return ToneDanger
	}
	panic("feedback: unknown toast type " + kind)
}

// Toaster styles window toasts; Type picks info, success, warning, error.
func Toaster(c *ui.Context, o ToasterOptions) ui.Element {
	if o.Close == "" {
		o.Close = "Close"
	}
	t := c.Theme()
	return ui.ToastViewportBase(c, func(viewport ui.Element, toasts []ui.Toast) {
		viewport.Padding(t.Space(4)).Gap(t.Space(2)).AlignItems(ui.End)
		seen := ui.Local(viewport, "seen", func() map[uint64]bool { return map[uint64]bool{} })
		for _, toast := range toasts {
			tone := toastTone(toast.Type)
			c.Key(toast.ID)
			parts := ui.ToastBase(c, toast)
			if id := parts.Root.ID(); !(*seen)[id] {
				(*seen)[id] = true
				sound.Play(sound.Notify)
				kit.Log("Toaster", "shown", "count", len(toasts))
			}
			kit.Arrive(toastFace(c, parts.Root), ui.Pose{Y: t.Space(4), Fade: 1}, 0).Transition(toastMotion).Children(func() {
				toastBody(c, tone, toast.Title, toast.Description)
				keys := ui.Column(c).Gap(t.Space(2)).AlignItems(ui.End)
				if keys.Shortcut(0, ui.KeyEscape) {
					c.CloseToast(toast.ID)
					sound.Play(sound.Close)
					kit.Log("Toaster", "escaped")
				}
				keys.Children(func() {
					if kit.Clicked(cross(c, parts.CloseButton().Label(o.Close)), sound.Close) {
						kit.Log("Toaster", "closed")
					}
					if toast.Action != "" {
						a := parts.ActionButton().Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).Border(1, t.Border)
						bg := t.Surface
						if a.Hovered() {
							bg = t.SurfaceHover
						}
						a.Background(a.AnimateColor("bg", bg, quick))
						kit.Press(a).Children(func() { ui.Text(c, toast.Action).SingleLine() })
						if kit.Clicked(a, sound.Tap) {
							kit.Log("Toaster", "action")
						}
					}
				})
			})
		}
	})
}
