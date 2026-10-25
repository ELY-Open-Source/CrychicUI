package mail

import (
	"cmp"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ContactsOptions word the empty list.
type ContactsOptions struct{ Empty string }

// Contacts lists people under letters and selects one; cue Select.
func Contacts(c *ui.Context, contacts []Contact, selected *string, o ContactsOptions) ui.Element {
	if selected == nil {
		panic("mail: Contacts needs selection")
	}
	def(&o.Empty, "No contacts")
	ids := make([]string, len(contacts))
	for i, p := range contacts {
		validContact("Contacts", p)
		ids[i] = p.ID
	}
	unique("Contacts", ids)
	sorted := slices.Clone(contacts)
	slices.SortStableFunc(sorted, func(a, b Contact) int {
		return cmp.Or(strings.Compare(initial(a.Name), initial(b.Name)), strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)))
	})
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleList)
	root.Children(func() {
		if len(sorted) == 0 {
			ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).TextColor(t.TextMuted).Children(func() {
				glyph(c, "users", 16)
				ui.Text(c, o.Empty).FontSize(t.Rem(0.8)).Shrink(0)
			})
		}
		n := 0
		for start := 0; start < len(sorted); {
			letter := initial(sorted[start].Name)
			end := start
			for end < len(sorted) && initial(sorted[end].Name) == letter {
				end++
			}
			group := sorted[start:end]
			c.Key(letter)
			ui.Column(c).FillWidth().Gap(t.Space(0.5)).Children(func() {
				kit.Muted(c, letter).FontSize(t.Rem(0.7)).Bold().Padding(0, t.Space(2.5))
				for _, p := range group {
					c.Key(p.ID)
					on := *selected == p.ID
					r := row(c, p.Name, on)
					kit.Stagger(c, root, r, n).Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
							avatar(c, p.Name, 8)
							ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
								ui.Text(c, p.Name).FontSize(t.Rem(0.85)).Bold().FillWidth().MinWidth(0).SingleLine()
								meta := p.Email
								if p.Role != "" {
									meta = p.Role
								}
								kit.Muted(c, meta).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
							})
						})
					})
					n++
					if r.Clicked() && !on {
						*selected = p.ID
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("Contacts", "selected")
					}
				}
			})
			start = end
		}
	})
	return root
}

// initial is a name's first letter, upper case.
func initial(name string) string {
	r, _ := utf8.DecodeRuneInString(name)
	return string(unicode.ToUpper(r))
}
