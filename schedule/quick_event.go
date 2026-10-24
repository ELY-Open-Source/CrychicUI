package schedule

import (
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// QuickEventOptions are the quick event labels.
type QuickEventOptions struct{ Placeholder, Create, Hint, Added string }

// QuickEvent enters a one-line event; cue Pop.
func QuickEvent(c *ui.Context, title *string, created *bool, o QuickEventOptions) (root, create ui.Element) {
	if title == nil || created == nil {
		panic("schedule: QuickEvent needs title and state")
	}
	def(&o.Placeholder, "Event title")
	def(&o.Create, "Create")
	def(&o.Hint, "Press Enter to add it to your calendar")
	def(&o.Added, "Added to calendar")
	t := c.Theme()
	empty := strings.TrimSpace(*title) == ""
	root = card(c).Gap(t.Space(2)).Padding(t.Space(3))
	pending := ui.Local(root, "pending", func() bool { return false })
	root.Children(func() {
		field := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(1), t.Space(1), t.Space(3)).
			Radius(t.Radius).Border(1, t.Border)
		field.Children(func() {
			glyph(c, "calendar-plus", 16).TextColor(t.TextMuted)
			in := ui.TextInput(c, title).Label(o.Placeholder).Placeholder(o.Placeholder).Grow(1).MinWidth(0)
			in.OnTextKey(func(k ui.TextKey) bool {
				if k.Key != ui.KeyEnter || strings.TrimSpace(k.Text()) == "" {
					return false
				}
				*pending = true
				return true
			})
			field.BorderColor(field.AnimateColor("line", t.Border.Mix(t.Accent, unit(in.Focused())), ui.Snappy))
			if in.Changed() {
				*created = false
				root.SetChanged()
			}
			c.Key("create")
			create = kit.Press(ui.ButtonBase(c).Label(o.Create).Shrink(0).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).
				TextColor(t.AccentText).Disabled(empty))
			create.Background(create.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(create.Hovered())), ui.Snappy))
			create.Children(func() {
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
					glyph(c, "plus", 14)
					ui.Text(c, o.Create).Bold().Shrink(0)
				})
			})
		})
		c.Key(*created)
		line := ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(t.TextMuted)
		if *created {
			line.TextColor(t.Success)
		}
		kit.Arrive(line, ui.Pose{Y: t.Space(1), Fade: 1}, 0).Children(func() {
			name, text := "corner-down-left", o.Hint
			if *created {
				name, text = "circle-check", o.Added
			}
			glyph(c, name, 12)
			ui.Text(c, text).FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).SingleLine()
		})
	})
	if *pending || create.Clicked() {
		*pending = false
		*created = true
		create.SetChanged()
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("QuickEvent", "created", "runes", utf8.RuneCountInString(*title))
	}
	return root, create
}
