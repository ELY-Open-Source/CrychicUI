package datetime

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// Unit is a step of calendar time.
type Unit int

// Units of calendar time.
const (
	UnitDay Unit = iota
	UnitWeek
	UnitMonth
	UnitYear
)

// Rule repeats every Every units, on Days for weeks.
type Rule struct {
	Every int
	Unit  Unit
	// Days are weekdays, Sunday first.
	Days [7]bool
}

// RecurrenceOptions tune a Recurrence.
type RecurrenceOptions struct {
	// Units name the units in the select.
	Units [4]string
	// One summarizes a single unit; Many formats a count.
	One, Many [4]string
	// On joins a summary and its weekdays.
	On string
	// Interval, Unit and Days name the three controls.
	Interval, Unit, Days string
}

// Recurrence edits a repeat rule and summarizes it.
func Recurrence(c *ui.Context, label string, rule *Rule, o RecurrenceOptions) ui.Element {
	if label == "" || rule == nil || rule.Every < 1 || rule.Every > 99 || rule.Unit < UnitDay || rule.Unit > UnitYear {
		panic("datetime: Recurrence needs a label and a rule of 1 to 99 known units")
	}
	if o.Units == ([4]string{}) {
		o.Units = [4]string{"Days", "Weeks", "Months", "Years"}
	}
	if o.One == ([4]string{}) {
		o.One = [4]string{"Every day", "Every week", "Every month", "Every year"}
	}
	if o.Many == ([4]string{}) {
		o.Many = [4]string{"Every %d days", "Every %d weeks", "Every %d months", "Every %d years"}
	}
	def(&o.On, "%s on %s")
	def(&o.Interval, "Interval")
	def(&o.Unit, "Unit")
	def(&o.Days, "Days")
	t := c.Theme()
	loc := locale(c)
	root := ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(3)).Label(label)
	return kit.Arrive(root, ui.Pose{Y: 4, Fade: 1}, 0).Children(func() {
		every, unit, changed := float64(rule.Every), o.Units[rule.Unit], false
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			if choice.NumberStepper(c, o.Interval, &every, 1, 99, 1, choice.NumberStepperOptions{}).Changed() {
				rule.Every, changed = int(every), true
			}
			if choice.Dropdown(c, o.Unit, o.Units[:], &unit).Changed() {
				rule.Unit, changed = Unit(slices.Index(o.Units[:], unit)), true
			}
		})
		if rule.Unit == UnitWeek {
			c.Key("days")
			ui.Box(c).Transition(kit.Enter).Children(func() {
				changed = choice.DayToggles(c, o.Days, &rule.Days, choice.DayTogglesOptions{
					Letters: i18n.Weekdays(loc, i18n.Narrow, true), Names: i18n.Weekdays(loc, i18n.Wide, true),
				}).Changed() || changed
			})
		}
		kit.Muted(c, summary(c, *rule, o)).ContentTransition(ui.ContentTransition{})
		if changed {
			root.SetChanged()
			kit.Log("Recurrence", "changed", "every", rule.Every, "unit", int(rule.Unit))
		}
	})
}

// summary writes a rule in words.
func summary(c *ui.Context, r Rule, o RecurrenceOptions) string {
	s := o.One[r.Unit]
	if r.Every > 1 {
		s = fmt.Sprintf(o.Many[r.Unit], r.Every)
	}
	if r.Unit != UnitWeek {
		return s
	}
	names, first := i18n.Weekdays(locale(c), i18n.Abbreviated, true), int(c.Region().FirstWeekday)
	var on []string
	for i := range 7 {
		if d := (first + i) % 7; r.Days[d] {
			on = append(on, names[d])
		}
	}
	if len(on) == 0 {
		return s
	}
	return fmt.Sprintf(o.On, s, strings.Join(on, ", "))
}
