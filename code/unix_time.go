// Code components for the gallery.
package code

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"strconv"
	"strings"
	"time"
)

type UnixTimeOptions struct {
	Epoch     string
	Date      string
	BadNumber string
	BadDate   string
}

// UnixTime converts between an epoch and a date; cue Tick.

func UnixTime(c *ui.Context, label string, epoch *int64, o UnixTimeOptions) ui.Element {
	const layout = "2006-01-02 15:04:05"
	if label == "" || epoch == nil {
		panic("code: UnixTime needs a label and a bound epoch")
	}
	def(&o.Epoch, "Unix seconds")
	def(&o.Date, "Date (UTC)")
	def(&o.BadNumber, "not a whole number")
	def(&o.BadDate, "use "+layout)
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	last := ui.Local(root, "last", func() int64 { return *epoch })
	epochText := ui.Local(root, "epochText", func() string { return strconv.FormatInt(*epoch, 10) })
	dateText := ui.Local(root, "dateText", func() string { return timeUnix(*epoch, layout) })
	bad := ui.Local(root, "bad", func() string { return "" })
	if *epoch != *last {
		*last = *epoch
		*epochText = strconv.FormatInt(*epoch, 10)
		*dateText = timeUnix(*epoch, layout)
		*bad = ""
	}
	root.Children(func() {
		ui.Row(c).Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			ef := ui.TextInput(c, epochText).Label(o.Epoch).Font("monospace").Width(t.Space(40))
			df := ui.TextInput(c, dateText).Label(o.Date).Font("monospace").Width(t.Space(48))
			if ef.Changed() {
				if v, err := strconv.ParseInt(strings.TrimSpace(*epochText), 10, 64); err != nil {
					*bad = o.BadNumber
					root.SetChanged()
				} else {
					*epoch, *last = v, v
					*dateText = timeUnix(v, layout)
					*bad = ""
					root.SetChanged()
					sound.Play(sound.Tick)
					kit.Log("UnixTime", "epoch")
				}
			}
			if df.Changed() {
				if tm, err := parseTime(layout, strings.TrimSpace(*dateText)); err != nil {
					*bad = o.BadDate
					root.SetChanged()
				} else {
					*epoch, *last = tm, tm
					*epochText = strconv.FormatInt(tm, 10)
					*bad = ""
					root.SetChanged()
					sound.Play(sound.Tick)
					kit.Log("UnixTime", "date")
				}
			}
		})
		if *bad != "" {
			ui.Text(c, *bad).FontSize(t.Rem(0.85)).TextColor(t.Danger)
		}
	})
	return root
}

// timeUnix formats epoch seconds in UTC.

func timeUnix(epoch int64, layout string) string {
	return time.Unix(epoch, 0).UTC().Format(layout)
}

// parseTime parses a UTC date to epoch seconds.

func parseTime(layout, s string) (int64, error) {
	tm, err := time.ParseInLocation(layout, s, time.UTC)
	if err != nil {
		return 0, err
	}
	return tm.Unix(), nil
}

// PackageCardOptions are the copy knobs of PackageCard.
