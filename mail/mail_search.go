package mail

import (
	"unicode/utf8"

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
	root := card(c).Gap(t.Space(2)).Padding(t.Space(2.5))
	root.Children(func() {
		field := func(mark, label string, text *string) {
			box := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2.5)).Radius(999).Border(1, t.Border)
			box.Children(func() {
				glyph(c, mark, 14).TextColor(t.TextMuted)
				in := ui.TextInputBase(c, text).Label(label).Placeholder(label).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0)
				box.BorderColor(box.AnimateColor("line", t.Border.Mix(t.Accent, unit(in.Focused())), ui.Snappy))
				if in.Changed() {
					root.SetChanged()
					sound.Play(sound.Tick)
					kit.Log("MailSearch", "edited", "runes", utf8.RuneCountInString(query.Text)+utf8.RuneCountInString(query.From))
				}
			})
		}
		field("search", o.Search, &query.Text)
		field("user", o.From, &query.From)
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).Children(func() {
			filters := []struct {
				label, mark string
				value       *bool
			}{{o.Attachments, "paperclip", &query.Attachments}, {o.Links, "link", &query.Links}}
			for i, filter := range filters {
				c.Key(i)
				on := *filter.value
				chip := kit.Press(ui.ButtonBase(c).Label(filter.label).Padding(t.Space(1), t.Space(2.5)).Radius(999).Shrink(0).Checked(on))
				chip.Background(chip.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(chip.Hovered())).Mix(t.Accent.Alpha(0.12), unit(on)), ui.Snappy))
				chip.Border(1, chip.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
				chip.TextColor(chip.AnimateColor("ink", t.Text.Mix(t.Accent, unit(on)), ui.Snappy))
				chip.Children(func() {
					ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
						c.Key(on)
						mark := filter.mark
						if on {
							mark = "check"
						}
						kit.Arrive(glyph(c, mark, 12), ui.Pose{Scale: 0.5, Fade: 1}, 0)
						ui.Text(c, filter.label).FontSize(t.Rem(0.78)).Shrink(0)
					})
				})
				if chip.Clicked() {
					*filter.value = !on
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("MailSearch", "filter", "index", i, "on", *filter.value)
				}
			}
		})
	})
	return root
}
