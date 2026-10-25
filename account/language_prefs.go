package account

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// LanguagePrefsOptions word the card; Note takes the region.
type LanguagePrefsOptions struct{ Title, Language, Region, Note string }

// LanguagePrefs picks language and region; cue Open, Select.
func LanguagePrefs(c *ui.Context, languages []string, lang *string, regions []string, region *string, o LanguagePrefsOptions) ui.Element {
	if lang == nil || region == nil || !slices.Contains(languages, *lang) || !slices.Contains(regions, *region) {
		panic(fmt.Sprintf("account: LanguagePrefs %q needs choices holding the current language and region", o.Title))
	}
	unique("LanguagePrefs", languages)
	unique("LanguagePrefs", regions)
	def(&o.Title, "Language and region")
	def(&o.Language, "Language")
	def(&o.Region, "Region")
	def(&o.Note, "Dates, numbers and currency follow %s.")
	t := c.Theme()
	root := card(c).Gap(t.Space(4)).Padding(t.Space(5)).Role(ui.RoleGroup)
	root.Children(func() {
		head(c, "languages", o.Title, "")
		c.Key("language")
		if picker(c, o.Language, "languages", languages, lang) {
			root.SetChanged()
		}
		c.Key("region")
		if picker(c, o.Region, "globe", regions, region) {
			root.SetChanged()
		}
		c.Key("note")
		ui.Column(c).FillWidth().Children(func() {
			c.Key(*region)
			kit.Arrive(kit.Muted(c, fmt.Sprintf(o.Note, *region)).FontSize(t.Rem(0.78)).Wrap().FillWidth(), ui.Pose{Fade: 1}, 0)
		})
	})
	return root
}

// picker is a captioned trigger opening a choice list.
func picker(c *ui.Context, caption, mark string, choices []string, value *string) bool {
	t := c.Theme()
	changed := false
	col := ui.Column(c).FillWidth()
	open := ui.Local(col, "open", func() bool { return false })
	col.Children(func() {
		var trigger ui.Element
		ui.Grid(c).Columns(1).FillWidth().Gap(t.Space(1)).Children(func() {
			trigger = kit.Press(ui.ButtonBase(c).Label(caption).RowStart(2).ColumnStart(1).FillWidth().
				Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Expanded(*open))
			trigger.Background(trigger.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(trigger.Hovered())), ui.Snappy))
			trigger.Border(1, trigger.AnimateColor("line", t.Border.Mix(t.Accent, unit(*open)), ui.Snappy))
			trigger.Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					glyph(c, mark, 16).TextColor(t.TextMuted)
					ui.Text(c, *value).FontSize(t.Rem(0.88)).Grow(1).MinWidth(0).SingleLine().ContentTransition(ui.ContentTransition{})
					chev := glyph(c, "chevron-down", 16).TextColor(t.TextMuted)
					chev.Rotate(180 * chev.AnimateTo("turn", unit(*open), ui.Snappy))
				})
			})
			ui.Text(c, caption).RowStart(1).ColumnStart(1).FontSize(t.Rem(0.78)).Bold().TextColor(t.TextMuted)
		})
		if trigger.Clicked() {
			*open = !*open
			if *open {
				sound.Play(sound.Open)
				kit.Log("LanguagePrefs", "opened")
			}
		}
		ui.PopoverBase(c, trigger, open, func(panel ui.Element) {
			panel.Margin(t.Space(1), 0, 0, 0).MinWidth(t.Space(60)).Padding(t.Space(1)).Gap(t.Space(0.5)).Radius(t.Radius).
				Background(t.Background).Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).TransformOrigin(0.5, 0)
			kit.Arrive(panel, ui.Pose{Scale: 0.95, Fade: 1}, 0)
			ui.Scroll(c).FillWidth().MaxHeight(t.Space(64)).Gap(t.Space(0.5)).Role(ui.RoleMenu).Children(func() {
				for i, ch := range choices {
					c.Key(ch)
					on := ch == *value
					item := kit.Press(ui.ButtonBase(c).Label(ch).FillWidth().Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).
						Role(ui.RoleMenuItem).Checked(on))
					item.Background(item.AnimateColor("bg", t.SurfaceHover.Alpha(unit(item.Hovered())), ui.Snappy))
					item.Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, ch).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
							glyph(c, "check", 14).TextColor(t.Accent).Opacity(unit(on))
						})
					})
					if item.Clicked() {
						*open = false
						if !on {
							*value = ch
							changed = true
							sound.Play(sound.Select)
							kit.Log("LanguagePrefs", "selected", "index", i)
						}
					}
				}
			})
		})
	})
	return changed
}
