package project

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// QuickAddOptions word the field, due chip and note.
type QuickAddOptions struct{ Placeholder, Add, Due, Added, NoDue, DateFormat string }

// QuickAdd enters a task with parsed due; cue Pop.
func QuickAdd(c *ui.Context, title *string, due *time.Time, created *bool, o QuickAddOptions) (root, add ui.Element) {
	if title == nil || due == nil || created == nil {
		panic("project: QuickAdd needs title, due, and state")
	}
	def(&o.Placeholder, "Task and due")
	def(&o.Add, "Add")
	def(&o.Due, "Due %s")
	def(&o.Added, "Added · %s")
	def(&o.NoDue, "no due date")
	def(&o.DateFormat, "Mon, Jan 2")
	t := c.Theme()
	root = card(c).Gap(t.Space(2)).Padding(t.Space(2.5))
	pending := ui.Local(root, "pending", func() bool { return false })
	note := ui.Local(root, "note", func() string { return "" })
	parsed := parseQuickDue(c.Now(), *title)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			field := ui.Row(c).Grow(1).MinWidth(0).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2.5)).
				Radius(t.Radius).Border(1, t.Border)
			field.Children(func() {
				glyph(c, "circle-plus", 16).TextColor(t.TextMuted)
				in := ui.TextInputBase(c, title).Label(o.Placeholder).Placeholder(o.Placeholder).Grow(1).MinWidth(0)
				in.OnTextKey(func(k ui.TextKey) bool {
					if k.Key != ui.KeyEnter || strings.TrimSpace(k.Text()) == "" {
						return false
					}
					*pending = true
					return true
				})
				field.BorderColor(field.AnimateColor("line", t.Border.Mix(t.Accent, unit(in.Focused())), ui.Snappy))
				if in.Changed() {
					*note = ""
					root.SetChanged()
				}
				c.Key(parsed.IsZero())
				if !parsed.IsZero() {
					chip := ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(2)).Radius(999).
						Background(t.Accent.Alpha(0.12)).TextColor(t.Accent)
					kit.Arrive(chip, ui.Pose{Scale: 0.7, Fade: 1}, 0).Children(func() {
						glyph(c, "calendar", 12)
						ui.Text(c, fmt.Sprintf(o.Due, parsed.Format(o.DateFormat))).FontSize(t.Rem(0.7)).Bold().Shrink(0).
							ContentTransition(ui.ContentTransition{})
					})
				} else {
					glyph(c, "corner-down-left", 14).TextColor(t.TextMuted)
				}
			})
			add = kit.Press(ui.ButtonBase(c).Label(o.Add).Padding(t.Space(2), t.Space(3.5)).Radius(t.Radius).Shrink(0).
				TextColor(t.AccentText).Disabled(strings.TrimSpace(*title) == ""))
			add.Background(add.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(add.Hovered())), ui.Snappy))
			add.Children(func() { ui.Text(c, o.Add).FontSize(t.Rem(0.85)).Bold().Shrink(0) })
		})
		if *note != "" {
			row := ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(0, t.Space(1)).TextColor(t.Success)
			kit.Arrive(row, ui.Pose{Y: t.Space(1), Fade: 1}, 0).Children(func() {
				glyph(c, "circle-check", 14)
				ui.Text(c, *note).FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).SingleLine()
			})
		}
	})
	if *pending || add.Clicked() {
		*pending = false
		*due = parseQuickDue(c.Now(), *title)
		*created = true
		when := o.NoDue
		if !due.IsZero() {
			when = fmt.Sprintf(o.Due, due.Format(o.DateFormat))
		}
		*note = fmt.Sprintf(o.Added, when)
		add.SetChanged()
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("QuickAdd", "added", "hasDue", !due.IsZero())
	}
	return root, add
}

func parseQuickDue(now time.Time, title string) time.Time {
	text := strings.ToLower(strings.TrimSpace(title))
	for _, field := range strings.Fields(text) {
		switch field {
		case "today":
			return now
		case "tomorrow":
			return now.AddDate(0, 0, 1)
		default:
			if parsed, err := time.ParseInLocation("2006-01-02", field, now.Location()); err == nil {
				return parsed
			}
		}
	}
	return time.Time{}
}
