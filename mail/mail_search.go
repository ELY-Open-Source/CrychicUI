package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MailQuery is one mail search.
type MailQuery struct {
	Text, From         string
	Attachments, Links bool
}

// MailSearchOptions are the search labels.
type MailSearchOptions struct{ Search, From, Attachments, Links string }

// MailSearch edits text, sender, and filter chips.
func MailSearch(c *ui.Context, query *MailQuery, o MailSearchOptions) ui.Element {
	if query == nil {
		panic("mail: MailSearch needs a query")
	}
	def(&o.Search, "Search mail")
	def(&o.From, "From")
	def(&o.Attachments, "Has attachment")
	def(&o.Links, "Has link")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		search := ui.TextInput(c, &query.Text).Label(o.Search).Placeholder(o.Search).FillWidth()
		from := ui.TextInput(c, &query.From).Label(o.From).Placeholder(o.From).FillWidth()
		if search.Changed() || from.Changed() {
			root.SetChanged()
			sound.Play(sound.Tick)
			kit.Log("MailSearch", "edited", "runes", len([]rune(query.Text))+len([]rune(query.From)))
		}
		ui.Row(c).Wrap().Gap(t.Space(1)).Children(func() {
			filters := []struct {
				label string
				value *bool
			}{{o.Attachments, &query.Attachments}, {o.Links, &query.Links}}
			for i, filter := range filters {
				c.Key(i)
				chip := button(c, filter.label).Radius(999).Checked(*filter.value)
				chip.Children(func() { ui.Text(c, filter.label) })
				if chip.Clicked() {
					*filter.value = !*filter.value
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("MailSearch", "filter", "index", i, "on", *filter.value)
				}
			}
		})
	})
	return root
}
