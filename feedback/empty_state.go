package feedback

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// state is a centered icon, title, text and optional action.
func state(c *ui.Context, svg *ui.SVG, tint ui.Color, title, text, action string, primary bool, glyph func(ui.Element)) (root, act ui.Element) {
	if title == "" {
		panic("feedback: a state needs a title")
	}
	t := c.Theme()
	root = ui.Column(c).AlignSelf(ui.Center).MaxWidth(t.Space(100)).AlignItems(ui.Center).Gap(t.Space(2)).
		Padding(t.Space(6)).Role(ui.RoleStatus).Label(title)
	kit.Arrive(root, ui.Pose{Y: t.Space(3), Fade: 1}, 0)
	root.Children(func() {
		disc := ui.Box(c).Size(t.Space(12), t.Space(12)).Center().Radius(999).Background(tint.Alpha(0.12)).
			TextColor(tint).Margin(0, 0, t.Space(1), 0)
		disc.Children(func() {
			i := ui.Icon(c, svg).Size(t.Space(6), t.Space(6))
			if glyph != nil {
				glyph(i)
			}
		})
		ui.Text(c, title).Bold().FontSize(t.Rem(1.15)).TextAlign(ui.Center)
		if text != "" {
			kit.Muted(c, text).TextAlign(ui.Center)
		}
		if action != "" {
			act = button(c, action, primary).Margin(t.Space(2), 0, 0, 0)
		}
	})
	return root, act
}

// EmptyState is icon, title, text and action rising in.
func EmptyState(c *ui.Context, glyph, title, text, action string) (root, act ui.Element) {
	root, act = state(c, icon.Must(glyph), c.Theme().TextMuted, title, text, action, true, nil)
	if act.Valid() && kit.Clicked(act, sound.Tap) {
		kit.Log("EmptyState", "action")
	}
	return root, act
}
