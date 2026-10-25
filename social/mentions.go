package social

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MentionsOptions word the editor and the empty suggestion list.
type MentionsOptions struct{ Placeholder, Empty string }

// Mentions suggests people after @ while typing; cue Select.
func Mentions(c *ui.Context, text *string, users []User, selected *string, o MentionsOptions) ui.Element {
	if text == nil || len(users) == 0 || selected == nil {
		panic("social: Mentions needs text, users, and selection")
	}
	ids := make([]string, len(users))
	for i, user := range users {
		validUser("Mentions", user)
		if strings.TrimPrefix(user.Handle, "@") == "" {
			panic(fmt.Sprintf("social: Mentions user %q needs a handle", user.Name))
		}
		ids[i] = user.ID
	}
	unique("Mentions", ids)
	def(&o.Placeholder, "Write with @mentions")
	def(&o.Empty, "No one matches")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1.5)).Role(ui.RoleGroup)
	pendingCaret := ui.Local(root, "pendingCaret", func() int { return -1 })
	// The caret outlives focus, which a suggestion click takes.
	lastCaret := ui.Local(root, "lastCaret", func() int { return -1 })
	root.Children(func() {
		field := ui.Column(c).FillWidth().Padding(t.Space(2.5), t.Space(3)).Radius(t.Radius).Background(t.Background)
		field.Border(1, field.AnimateColor("ring", t.Border.Mix(t.Accent, unit(field.FocusWithin())), ui.Snappy))
		var in ui.Element
		field.Children(func() {
			in = ui.TextAreaBase(c, text).Label(o.Placeholder).Placeholder(o.Placeholder).FillWidth().Lines(2, 4).FontSize(t.Rem(0.9))
			if *pendingCaret >= 0 {
				in.Focus().SetTextSelection(*pendingCaret, *pendingCaret)
				*pendingCaret = -1
			}
			if in.Changed() {
				root.SetChanged()
			}
		})
		// The token runs from the last @ to the caret.
		runes := []rune(*text)
		if in.Focused() {
			_, *lastCaret = in.TextSelection()
		}
		caret := len(runes)
		if *lastCaret >= 0 {
			caret = min(*lastCaret, len(runes))
		}
		before, after := string(runes[:caret]), string(runes[caret:])
		at := strings.LastIndex(before, "@")
		if at < 0 {
			return
		}
		query := before[at+1:]
		if strings.ContainsAny(query, " \t\n") {
			return
		}
		c.Key("suggestions")
		panel := ui.Column(c).FillWidth().Padding(t.Space(1)).Radius(t.Radius).Background(t.Background).Border(1, t.Border).
			Shadow(0, 4, 16, 0, t.Text.Alpha(0.12))
		kit.Arrive(panel, ui.Pose{Y: -t.Space(1), Fade: 1}, 0).Children(func() {
			shown := 0
			for _, user := range users {
				handle := strings.TrimPrefix(user.Handle, "@")
				if !strings.HasPrefix(strings.ToLower(handle), strings.ToLower(query)) {
					continue
				}
				c.Key(user.ID)
				b := kit.Press(ui.ButtonBase(c).Label(user.Handle).Row().FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).
					Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
				kit.Stagger(c, panel, b, shown).Children(func() {
					avatar(c, user, t.Space(8))
					ui.Column(c).Grow(1).MinWidth(0).Children(func() {
						ui.Text(c, user.Name).FontSize(t.Rem(0.85)).Bold().FillWidth().MinWidth(0).SingleLine()
						lo, hi := kit.IndexFold(handle, query)
						ui.RichText(c, kit.Marked(t, "@"+handle, lo+1, hi+1)...).FontSize(t.Rem(0.75)).TextColor(t.TextMuted).
							FillWidth().MinWidth(0).SingleLine()
					})
				})
				shown++
				if b.Clicked() {
					head := before[:at] + "@" + handle + " "
					*text = head + strings.TrimLeft(after, " ")
					*pendingCaret = utf8.RuneCountInString(head)
					*lastCaret = *pendingCaret
					*selected = user.ID
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Mentions", "selected")
				}
			}
			if shown == 0 {
				ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).TextColor(t.TextMuted).Children(func() {
					glyph(c, "at-sign", 14)
					ui.Text(c, o.Empty).FontSize(t.Rem(0.8)).Shrink(0)
				})
			}
		})
	})
	return root
}
