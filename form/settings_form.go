package form

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/egoist/mygo/ui"
)

// Setting is a switch row (On) or select row (Choice).
type Setting struct {
	Title, Description string
	On                 *bool
	Choice             *string
	Choices            []string
}

// SettingGroup is a titled list of settings.
type SettingGroup struct {
	Title    string
	Settings []Setting
}

// SettingsForm lays out grouped switch and select rows.
func SettingsForm(c *ui.Context, groups []SettingGroup) ui.Element {
	if len(groups) == 0 {
		panic("form: SettingsForm needs groups")
	}
	for _, g := range groups {
		for _, s := range g.Settings {
			if s.Title == "" || (s.On == nil) == (s.Choice == nil) {
				panic(fmt.Sprintf("form: setting %q needs a title and exactly one of On or Choice", s.Title))
			}
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(6))
	return root.Children(func() {
		for _, g := range groups {
			layout.Section(c, g.Title, func() {
				for _, s := range g.Settings {
					if s.On != nil {
						if choice.ToggleRow(c, s.Title, s.Description, s.On).Changed() {
							root.SetChanged()
							kit.Log("SettingsForm", "toggled", "setting", s.Title, "on", *s.On)
						}
						continue
					}
					ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
						ui.Column(c).Grow(1).Gap(t.Space(0.5)).Children(func() {
							ui.Text(c, s.Title).FontWeight(500)
							if s.Description != "" {
								kit.Muted(c, s.Description)
							}
						})
						if choice.Dropdown(c, s.Title, s.Choices, s.Choice).Changed() {
							root.SetChanged()
							kit.Log("SettingsForm", "selected", "setting", s.Title)
						}
					})
				}
			})
		}
	})
}
