// Package project holds task management components.
package project

import (
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Task is one project task.
type Task struct {
	ID, Title, Assignee, Due string
	Labels                   []string
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

func progress(c *ui.Context, value float64, color ui.Color) ui.Element {
	if value < 0 || value > 1 || !finite(value) {
		panic("project: progress must be in [0,1]")
	}
	t := c.Theme()
	bar := ui.Box(c).FillWidth().Height(t.Space(1)).Radius(999).Clip().Background(t.SurfaceHover)
	w := bar.AnimateTo("progress", float32(value), ui.Smooth)
	bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * w, H: r.H}, color, 999)
	}))
	return bar
}
