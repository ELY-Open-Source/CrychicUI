package datetime

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// Relative is a choice some days from today.
type Relative struct {
	Label string
	Days  int
}

// RelativeDateOptions tune a RelativeDate.
type RelativeDateOptions struct {
	// Choices default to today through next week.
	Choices []Relative
}

// RelativeDate picks "in N days" chips and shows the date.
func RelativeDate(c *ui.Context, label string, date *time.Time, o RelativeDateOptions) ui.Element {
	if label == "" || date == nil {
		panic("datetime: RelativeDate needs a label and a date")
	}
	if o.Choices == nil {
		o.Choices = []Relative{{"Today", 0}, {"Tomorrow", 1}, {"In 3 days", 3}, {"Next week", 7}}
	}
	for _, r := range o.Choices {
		if r.Label == "" {
			panic(fmt.Sprintf("datetime: RelativeDate choice of %d days needs a label", r.Days))
		}
	}
	t := c.Theme()
	today := c.Now()
	root := ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(2)).Role(ui.RoleRadioGroup).Label(label)
	return root.Children(func() {
		ui.Row(c).Gap(t.Space(1.5)).Children(func() {
			for _, r := range o.Choices {
				day := today.AddDate(0, 0, r.Days)
				on := sameDay(day, *date)
				chip := kit.Press(ui.ButtonBase(c).Label(r.Label).Padding(t.Space(1), t.Space(3)).Radius(999)).Checked(on)
				bg, fg := t.Surface, t.Text
				if on {
					bg, fg = t.Accent, t.AccentText
				}
				chip.Background(chip.AnimateColor("bg", bg, ui.Snappy)).TextColor(chip.AnimateColor("fg", fg, ui.Snappy)).
					Children(func() { ui.Text(c, r.Label) })
				if chip.Clicked() {
					*date = onDay(day, *date)
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("RelativeDate", "picked", "days", r.Days)
				}
			}
		})
		kit.Muted(c, i18n.FormatDate(locale(c), *date, i18n.DateOptions{Skeleton: "yMMMEd"})).
			ContentTransition(ui.ContentTransition{})
	})
}
