// Package project holds task management components.
package project

import (
	"fmt"
	"hash/fnv"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Task is one project task; Assignee and Due are optional.
type Task struct {
	ID, Title, Assignee, Due string
	Labels                   []string
	Done                     bool
}

// PriorityLevel is a task priority.
type PriorityLevel uint8

const (
	Urgent PriorityLevel = iota
	High
	Medium
	Low
)

// TaskStatus is a task state.
type TaskStatus uint8

const (
	Todo TaskStatus = iota
	Doing
	Done
)

func def(s *string, d string) {
	if *s == "" {
		*s = d
	}
}

func unit(v bool) float32 {
	if v {
		return 1
	}
	return 0
}

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

// percent rounds a share; 100 only when full.
func percent(v float64) int {
	if v < 1 {
		return min(99, int(math.Round(v*100)))
	}
	return int(math.Round(v * 100))
}

func validTask(component string, task Task) {
	if task.ID == "" || task.Title == "" {
		panic(fmt.Sprintf("project: %s has a task %q without ID or title", component, task.ID))
	}
}

func priorityText(level PriorityLevel, urgent, high, medium, low string) string {
	switch level {
	case Urgent:
		return urgent
	case High:
		return high
	case Medium:
		return medium
	default:
		return low
	}
}

func statusText(status TaskStatus, todo, doing, done string) string {
	switch status {
	case Todo:
		return todo
	case Doing:
		return doing
	default:
		return done
	}
}

// priorityLook is a level's icon and ink.
func priorityLook(t *ui.Theme, level PriorityLevel) (string, ui.Color) {
	switch level {
	case Urgent:
		return "flame", t.Danger
	case High:
		return "chevrons-up", t.Warning
	case Medium:
		return "equal", t.Accent
	}
	return "chevrons-down", t.TextMuted
}

// statusLook is a status's icon and ink.
func statusLook(t *ui.Theme, status TaskStatus) (string, ui.Color) {
	switch status {
	case Doing:
		return "circle-dot", t.Accent
	case Done:
		return "circle-check", t.Success
	}
	return "circle-dashed", t.TextMuted
}

// labelInk colors a label by its name, the same everywhere.
func labelInk(t *ui.Theme, name string) ui.Color {
	h := fnv.New32a()
	h.Write([]byte(name))
	return []ui.Color{t.Accent, t.Success, t.Warning, t.Danger}[h.Sum32()%4]
}

func progress(c *ui.Context, value float64, color ui.Color) ui.Element {
	if value < 0 || value > 1 || !finite(value) {
		panic("project: progress must be in [0,1]")
	}
	t := c.Theme()
	bar := ui.Box(c).FillWidth().Height(t.Space(1.5)).Radius(999).Clip().Background(t.SurfaceHover)
	w := bar.AnimateTo("progress", float32(value), ui.Smooth)
	tint := bar.AnimateColor("ink", color, ui.Snappy)
	bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * w, H: r.H}, tint, 999)
	}))
	return bar
}

// card is the shared arriving card surface.
func card(c *ui.Context) ui.Element {
	t := c.Theme()
	return kit.Arrive(kit.Card(c).FillWidth().Role(ui.RoleGroup), ui.Pose{Y: t.Space(2), Fade: 1}, 0)
}

// glyph is a small icon in the current text color.
func glyph(c *ui.Context, name string, size float32) ui.Element {
	return ui.Icon(c, icon.Must(name)).Size(size, size).Shrink(0)
}

// pill is a choice chip tinted by ink while on.
func pill(c *ui.Context, label string, on bool, ink ui.Color, content func()) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(1), t.Space(2.5)).Radius(999).Shrink(0).Checked(on))
	b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(ink.Alpha(0.14), unit(on)), ui.Snappy))
	b.Border(1, b.AnimateColor("line", t.Border.Mix(ink, unit(on)), ui.Snappy))
	b.TextColor(b.AnimateColor("ink", t.Text.Mix(ink, unit(on)), ui.Snappy))
	return b.Children(func() { ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(content) })
}

// tile is a tinted rounded icon square.
func tile(c *ui.Context, name string, ink ui.Color) ui.Element {
	t := c.Theme()
	return ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(t.Radius).Background(ink.Alpha(0.12)).
		TextColor(ink).Children(func() { glyph(c, name, 16) })
}

// tag is a small rounded label in ink.
func tag(c *ui.Context, text string, ink ui.Color) ui.Element {
	return fitTag(c, text, ink, math.MaxFloat32)
}

// fitTag is a tag at most room wide.
func fitTag(c *ui.Context, text string, ink ui.Color, room float32) ui.Element {
	t := c.Theme()
	return ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.25), t.Space(1.5)).Radius(999).
		Background(ink.Alpha(0.12)).TextColor(ink).Children(func() {
		ui.Box(c).Size(t.Space(1.5), t.Space(1.5)).Shrink(0).Radius(999).Background(ink)
		ui.Text(c, text).FontSize(t.Rem(0.65)).Bold().Shrink(0).MaxWidth(min(t.Space(40), max(0, room-t.Space(5.5)))).SingleLine()
	})
}

// unique panics when names repeat or are empty.
func unique(component string, names []string) {
	seen := map[string]bool{}
	for _, name := range names {
		if name == "" || seen[name] {
			panic(fmt.Sprintf("project: %s has an empty or repeated %q", component, name))
		}
		seen[name] = true
	}
}
