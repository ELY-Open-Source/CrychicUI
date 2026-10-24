package dashboard

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RangeBarOptions word the bar; Day and Days take the length.
type RangeBarOptions struct{ Compare, Previous, Next, Format, Day, Days string }

// RangeBar steps a range by its length; cue Select.
func RangeBar(c *ui.Context, start, end *time.Time, compare *bool, o RangeBarOptions) ui.Element {
	if start == nil || end == nil || start.After(*end) || compare == nil {
		panic("dashboard: RangeBar needs a start before its end and a compare flag")
	}
	def(&o.Compare, "Compare")
	def(&o.Previous, "Previous")
	def(&o.Next, "Next")
	def(&o.Format, "Jan 2, 2006")
	def(&o.Day, "%d day")
	def(&o.Days, "%d days")
	// Calendar days from Unix seconds: Sub saturates past 292 years.
	day := func(d time.Time) int64 {
		return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC).Unix() / 86400
	}
	span := int(day(*end) - day(*start))
	step := max(1, span)
	t := c.Theme()
	root := kit.Card(c).Row().FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(1.5)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for i, action := range []struct {
			name, glyph string
			days        int
		}{{o.Previous, "chevron-left", -step}, {o.Next, "chevron-right", step}} {
			c.Key(i)
			b := kit.Press(ui.ButtonBase(c).Label(action.name).Tooltip(action.name).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999))
			b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy))
			b.Children(func() { ui.Icon(c, icon.Must(action.glyph)).Size(16, 16) })
			if b.Clicked() {
				*start, *end = start.AddDate(0, 0, action.days), end.AddDate(0, 0, action.days)
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("RangeBar", "shifted", "days", action.days)
			}
		}
		ui.Row(c).Grow(1).MinWidth(0).Gap(t.Space(2)).AlignItems(ui.Center).Padding(0, t.Space(1)).Children(func() {
			ui.Icon(c, icon.Must("calendar")).Size(14, 14).Shrink(0).TextColor(t.TextMuted)
			ui.Text(c, start.Format(o.Format)+" – "+end.Format(o.Format)).FontSize(t.Rem(0.85)).Bold().Grow(1).MinWidth(0).SingleLine().
				ContentTransition(ui.ContentTransition{})
			length := o.Days
			if span == 1 {
				length = o.Day
			}
			kit.Chip(c, fmt.Sprintf(length, span), t.TextMuted)
		})
		b := kit.Press(ui.ButtonBase(c).Label(o.Compare).Shrink(0).Radius(999).Border(1, t.Border).Checked(*compare))
		b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(t.Accent.Alpha(0.14), unit(*compare)), ui.Snappy)).
			TextColor(b.AnimateColor("ink", t.Text.Mix(t.Accent, unit(*compare)), ui.Snappy))
		b.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2.5)).Children(func() {
				glyph := "git-compare"
				if *compare {
					glyph = "check"
				}
				ui.Icon(c, icon.Must(glyph)).Size(14, 14).ContentTransition(ui.ContentTransition{})
				ui.Text(c, o.Compare).FontSize(t.Rem(0.8)).Bold().Shrink(0)
			})
		})
		if b.Clicked() {
			*compare = !*compare
			root.SetChanged()
			sound.Play(sound.Select)
			kit.Log("RangeBar", "compare", "value", *compare)
		}
	})
	return root
}
