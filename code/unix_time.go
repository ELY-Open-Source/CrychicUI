package code

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// UnixTimeOptions are the copy knobs of UnixTime.
type UnixTimeOptions struct {
	Epoch, Date, BadNumber, BadDate, Now string
	// FormatLong words the date; FormatSince its distance from now.
	FormatLong  func(t time.Time) string
	FormatSince func(now, at time.Time) string
}

// UnixTime converts between an epoch and a date; cue Tick.
func UnixTime(c *ui.Context, label string, epoch *int64, o UnixTimeOptions) (root, now ui.Element) {
	const layout = "2006-01-02 15:04:05"
	if label == "" || epoch == nil {
		panic(fmt.Sprintf("code: UnixTime %q needs a label and a bound epoch", label))
	}
	def(&o.Epoch, "Unix seconds")
	def(&o.Date, "Date (UTC)")
	def(&o.BadNumber, "not a whole number")
	def(&o.BadDate, "use "+layout)
	def(&o.Now, "Now")
	if o.FormatLong == nil {
		o.FormatLong = func(t time.Time) string { return t.Format("Monday, 2 January 2006") }
	}
	if o.FormatSince == nil {
		o.FormatSince = sinceText
	}
	t := c.Theme()
	root = kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(label)
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
	set := func(v int64, event string) {
		*epoch, *last = v, v
		*bad = ""
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("UnixTime", event)
	}
	root.Children(func() {
		kit.CardHead(c, "clock", label, func() { now = action(c, o.Now, "timer", false) })
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(3)).Children(func() {
				ui.Column(c).Grow(1).MinWidth(0).Children(func() {
					captioned(c, o.Epoch, func() {
						if ui.TextInput(c, epochText).Label(o.Epoch).Font("monospace").FillWidth().Changed() {
							if v, err := strconv.ParseInt(strings.TrimSpace(*epochText), 10, 64); err != nil {
								*bad = o.BadNumber
								root.SetChanged()
							} else {
								*dateText = timeUnix(v, layout)
								set(v, "epoch")
							}
						}
					})
				})
				ui.Column(c).Grow(1).MinWidth(0).Children(func() {
					captioned(c, o.Date, func() {
						if ui.TextInput(c, dateText).Label(o.Date).Font("monospace").FillWidth().Changed() {
							if v, err := parseTime(layout, strings.TrimSpace(*dateText)); err != nil {
								*bad = o.BadDate
								root.SetChanged()
							} else {
								*epochText = strconv.FormatInt(v, 10)
								set(v, "date")
							}
						}
					})
				})
			})
			if *bad != "" {
				ui.Row(c).Key("bad").Gap(t.Space(1.5)).AlignItems(ui.Center).Transition(kit.Enter).Children(func() {
					ui.Icon(c, icon.Must("circle-x")).Size(14, 14).TextColor(t.Danger)
					ui.Text(c, *bad).FontSize(t.Rem(0.8)).TextColor(t.Danger)
				})
				return
			}
			at := time.Unix(*epoch, 0).UTC()
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).Background(t.Surface).Children(func() {
				ui.Icon(c, icon.Must("calendar")).Size(14, 14).Shrink(0).TextColor(t.Accent)
				ui.Text(c, o.FormatLong(at)).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine().ContentTransition(ui.ContentTransition{})
				kit.Muted(c, o.FormatSince(c.Now(), at)).FontSize(t.Rem(0.8)).Shrink(0)
			})
		})
	})
	if now.Clicked() {
		v := c.Now().Unix()
		*epochText, *dateText = strconv.FormatInt(v, 10), timeUnix(v, layout)
		set(v, "now")
	}
	return root, now
}

// sinceText words at from now: calendar years, else age.
func sinceText(now, at time.Time) string {
	lo, hi, ahead := at, now, at.After(now)
	if ahead {
		lo, hi = now, at
	}
	text := ""
	if y := hi.Year() - lo.Year(); y > 0 {
		if hi.AddDate(-y, 0, 0).Before(lo) {
			y--
		}
		if y > 0 {
			text = strconv.Itoa(y) + "y"
		}
	}
	if text == "" {
		d := hi.Sub(lo)
		if d < time.Minute {
			return "now"
		}
		text = age(d)
	}
	if ahead {
		return "in " + text
	}
	return text + " ago"
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
