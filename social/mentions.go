package social

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MentionsOptions are the mention editor labels.
type MentionsOptions struct{ Placeholder string }

// Mentions suggests people after @ while typing; cue Select.
func Mentions(c *ui.Context, text *string, users []User, selected *string, o MentionsOptions) ui.Element {
	if text == nil || len(users) == 0 || selected == nil {
		panic("social: Mentions needs text, users, and selection")
	}
	def(&o.Placeholder, "Write with @mentions")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	pendingCaret := ui.Local(root, "pendingCaret", func() int { return -1 })
	root.Children(func() {
		in := ui.TextArea(c, text).Label(o.Placeholder).Placeholder(o.Placeholder).FillWidth().Lines(2, 4)
		if *pendingCaret >= 0 {
			in.Focus().SetTextSelection(*pendingCaret, *pendingCaret)
			*pendingCaret = -1
		}
		if in.Changed() {
			root.SetChanged()
		}
		at := strings.LastIndex(*text, "@")
		if at < 0 {
			return
		}
		query := (*text)[at+1:]
		if strings.ContainsAny(query, " \t\n") {
			return
		}
		for _, user := range users {
			validUser(user)
			handle := strings.TrimPrefix(user.Handle, "@")
			if handle == "" {
				panic("social: mention user needs handle")
			}
			if !strings.HasPrefix(strings.ToLower(handle), strings.ToLower(query)) {
				continue
			}
			c.Key(user.ID)
			b := button(c, user.Handle).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
			b.Children(func() {
				avatar(c, user, t.Space(7))
				ui.Text(c, user.Name).Grow(1)
				kit.Muted(c, user.Handle)
			})
			if b.Clicked() {
				*text = (*text)[:at] + "@" + handle + " "
				end := len([]rune(*text))
				*pendingCaret = end
				*selected = user.ID
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Mentions", "selected")
			}
		}
	})
	return root
}
