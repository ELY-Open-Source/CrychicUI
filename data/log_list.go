package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Level is a log line's severity.
type Level uint8

// The four levels.
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// LogLine is a leveled, timed message.
type LogLine struct {
	Level      Level
	Time, Text string
}

// LogListOptions tune a LogList.
type LogListOptions struct {
	// Height of the list; zero gives the default.
	Height float32
	// Levels name the four levels, debug first.
	Levels [4]string
}

// LogList follows leveled lines; new ones fade in.
func LogList(c *ui.Context, label string, lines []LogLine, o LogListOptions) ui.Element {
	if label == "" {
		panic("data: LogList needs a label")
	}
	for _, l := range lines {
		if l.Level > LevelError {
			panic(fmt.Sprintf("data: LogList %q has an unknown level", label))
		}
	}
	for i, n := range [4]string{"DEBUG", "INFO", "WARN", "ERROR"} {
		def(&o.Levels[i], n)
	}
	t := c.Theme()
	ink := [4]ui.Color{t.TextMuted, t.Accent, t.Warning, t.Danger}
	root := frame(c, label).Height(height(c, o.Height))
	list := ui.Local(root, "list", func() ui.ListState { return ui.ListState{FollowEnd: true, FollowAnimation: &ui.Smooth} })
	seen := ui.Local(root, "seen", func() int { return 0 })
	from := *seen
	*seen = len(lines)
	root.Children(func() {
		ui.List(c, list, len(lines), func(i int) {
			l := lines[i]
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(1), t.Space(3)).FontFeatures("tnum")
			if *ui.Local(row, "fresh", func() bool { return i >= from }) {
				kit.Arrive(row, ui.Pose{Fade: 1}, 0)
			}
			row.Children(func() {
				kit.Muted(c, l.Time).Shrink(0).SingleLine()
				ui.Text(c, o.Levels[l.Level]).Bold().TextColor(ink[l.Level]).Width(t.Space(12)).Shrink(0).SingleLine()
				ui.Text(c, l.Text).Grow(1)
			})
		}).Label(label).Grow(1)
	})
	if len(lines) > from && from > 0 {
		kit.Log("LogList", "appended", "count", len(lines)-from)
	}
	return root
}
