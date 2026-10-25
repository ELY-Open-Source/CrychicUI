package account

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Setting is one on/off preference; Icon is optional.
type Setting struct {
	ID, Title, Detail, Icon string
	On                      bool
}

// PrivacyOptions word the card.
type PrivacyOptions struct{ Title string }

// Privacy toggles privacy settings; cue Toggle.
func Privacy(c *ui.Context, settings *[]Setting, o PrivacyOptions) ui.Element {
	if settings == nil || len(*settings) == 0 {
		panic("account: Privacy needs settings")
	}
	ids := make([]string, len(*settings))
	for i, s := range *settings {
		if s.Title == "" || s.Detail == "" {
			panic(fmt.Sprintf("account: Privacy %q needs a title and detail", s.ID))
		}
		ids[i] = s.ID
	}
	unique("Privacy", ids)
	def(&o.Title, "Privacy")
	t := c.Theme()
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		kit.CardHead(c, "shield", o.Title, nil)
		list := ui.Column(c).FillWidth().Padding(t.Space(1.5))
		list.Children(func() {
			for i, s := range *settings {
				c.Key(s.ID)
				row := kit.Press(ui.ButtonBase(c).Label(s.Title).FillWidth().Padding(t.Space(2.5), t.Space(2)).Radius(t.Radius).
					Role(ui.RoleSwitch).Checked(s.On))
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
						mark := s.Icon
						def(&mark, "eye")
						ink := t.TextMuted
						if s.On {
							ink = t.Accent
						}
						tile(c, mark, ink, t.Space(9))
						ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
							ui.Text(c, s.Title).FontSize(t.Rem(0.85)).Bold().Wrap().FillWidth()
							kit.Muted(c, s.Detail).FontSize(t.Rem(0.75)).Wrap().FillWidth()
						})
						track(c, s.On, row.Hovered())
					})
				})
				if row.Clicked() {
					(*settings)[i].On = !s.On
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("Privacy", "toggled", "index", i, "on", !s.On)
				}
			}
		})
	})
	return root
}
