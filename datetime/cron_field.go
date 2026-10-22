package datetime

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CronFieldOptions tune a CronField.
type CronFieldOptions struct {
	// Names label the five fields, minute first.
	Names [5]string
	// Every summarizes all stars; Invalid formats a bad field.
	Every, Invalid string
}

// cronRange bounds each cron field.
var cronRange = [5][2]int{{0, 59}, {0, 23}, {1, 31}, {1, 12}, {0, 6}}

// CronField edits five cron fields and summarizes them.
func CronField(c *ui.Context, label string, fields *[5]string, o CronFieldOptions) ui.Element {
	if label == "" || fields == nil {
		panic("datetime: CronField needs a label and fields")
	}
	if o.Names == ([5]string{}) {
		o.Names = [5]string{"Minute", "Hour", "Day", "Month", "Weekday"}
	}
	def(&o.Every, "Every minute")
	def(&o.Invalid, "Check the %s field")
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(2)).Label(label)
	return kit.Arrive(root, ui.Pose{Y: 4, Fade: 1}, 0).Children(func() {
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			for i := range fields {
				c.Key(i)
				ui.Column(c).Gap(t.Space(1)).Children(func() {
					kit.Muted(c, o.Names[i]).FontSize(t.Rem(0.75))
					if kit.Changed(ui.TextInput(c, &fields[i]).Label(o.Names[i]).Width(t.Space(16)), sound.Tick) {
						root.SetChanged()
						kit.Log("CronField", "edited", "field", i)
					}
				})
			}
		})
		var parts []string
		bad := -1
		for i, f := range fields {
			if !cronValid(f, cronRange[i][0], cronRange[i][1]) {
				bad = i
				break
			}
			if f != "*" {
				parts = append(parts, o.Names[i]+" "+f)
			}
		}
		text, ink := o.Every, t.TextMuted
		switch {
		case bad >= 0:
			text, ink = fmt.Sprintf(o.Invalid, o.Names[bad]), t.Danger
		case len(parts) > 0:
			text = strings.Join(parts, ", ")
		}
		line := ui.Text(c, text).FontSize(t.Rem(0.85)).ContentTransition(ui.ContentTransition{})
		line.TextColor(line.AnimateColor("ink", ink, ui.Snappy))
	})
}

// cronValid reports whether f is a cron field in range.
func cronValid(f string, lo, hi int) bool {
	for _, part := range strings.Split(f, ",") {
		base, step, stepped := strings.Cut(part, "/")
		if n, err := strconv.Atoi(step); stepped && (err != nil || n < 1) {
			return false
		}
		if base == "*" {
			continue
		}
		a, b, ranged := strings.Cut(base, "-")
		x, err := strconv.Atoi(a)
		y := x
		if err == nil && ranged {
			y, err = strconv.Atoi(b)
		}
		if err != nil || x < lo || y > hi || x > y {
			return false
		}
	}
	return true
}
