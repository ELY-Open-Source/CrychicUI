package chat

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TranslateOptions name the toggle.
type TranslateOptions struct {
	Show, Hide string
}

// Translate toggles original and translation; it cross-fades.
func Translate(c *ui.Context, original, translated string, showTranslation *bool, o TranslateOptions) ui.Element {
	if original == "" || translated == "" || showTranslation == nil {
		panic("chat: Translate needs an original, a translation and a toggle")
	}
	def(&o.Show, "Show translation")
	def(&o.Hide, "Show original")
	t := c.Theme()
	root := shell(c, false)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	fade := ui.ElementTransition{Duration: 160 * time.Millisecond, Enter: &ui.Motion{}, Exit: &ui.Motion{}}
	var toggle ui.Element
	root.Children(func() {
		text, key, label := original, "orig", o.Show
		if *showTranslation {
			text, key, label = translated, "trans", o.Hide
		}
		c.Key(key)
		ui.Text(c, text).Wrap().Transition(fade)
		ui.Row(c).FillWidth().Justify(ui.End).Children(func() {
			toggle = kit.Press(ui.ButtonBase(c).Label(label).Tooltip(label).Row().Gap(t.Space(1)).AlignItems(ui.Center).
				Padding(t.Space(0.5), t.Space(1.5)).Radius(999).TextColor(t.TextMuted))
			toggle.Background(toggle.AnimateColor("bg", t.Text.Alpha(0.08*unit(toggle.Hovered())), ui.Snappy)).Children(func() {
				ui.Icon(c, icon.Must("languages")).Size(16, 16)
				ui.Text(c, label).FontSize(t.Rem(0.75)).SingleLine().ContentTransition(ui.ContentTransition{})
			})
		})
	})
	if toggle.Clicked() {
		*showTranslation = !*showTranslation
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("Translate", "translated", *showTranslation)
	}
	return root
}
