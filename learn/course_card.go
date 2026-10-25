package learn

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Course is a course with lessons done of total.
type Course struct {
	Title, Author, Level string
	Lessons, Done        int
	Minutes              int
	Cover                ui.ImageSource
}

// CourseCardOptions word the card; Lessons and Time take counts.
type CourseCardOptions struct{ Lessons, Lesson, Time, Short, Progress, Complete string }

// CourseCard shows a course with its cover and progress.
func CourseCard(c *ui.Context, k Course, o CourseCardOptions) ui.Element {
	if strings.TrimSpace(k.Title) == "" || strings.TrimSpace(k.Author) == "" || k.Lessons < 1 || k.Done < 0 || k.Done > k.Lessons || k.Minutes < 0 {
		panic(fmt.Sprintf("learn: CourseCard %q needs an author, lessons and 0 <= done <= lessons", k.Title))
	}
	def(&o.Lessons, "%d lessons")
	def(&o.Lesson, "1 lesson")
	def(&o.Time, "%dh %02dm")
	def(&o.Short, "%d min")
	def(&o.Progress, "%d%% complete")
	def(&o.Complete, "Completed")
	t := c.Theme()
	done := k.Done == k.Lessons
	root := card(c).Clip()
	root.Children(func() {
		band := ui.Box(c).FillWidth().Height(t.Space(28)).Clip().Background(t.Accent.Alpha(0.16))
		band.Children(func() {
			if k.Cover != nil {
				ui.Image(c, k.Cover).Fill().Fit(ui.Cover)
			} else {
				ui.Box(c).Fill().Center().TextColor(t.Accent).Children(func() { glyph(c, "graduation-cap", 40) })
			}
			if k.Level != "" {
				ui.Box(c).Absolute().Left(t.Space(3)).Top(t.Space(3)).Children(func() { kit.Chip(c, k.Level, t.Accent) })
			}
		})
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3.5)).Children(func() {
			ui.Text(c, k.Title).FontSize(t.Rem(0.95)).Bold().FillWidth().Wrap()
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				display.Avatar(c, k.Author, nil).Size(t.Space(6), t.Space(6))
				ui.Text(c, k.Author).FontSize(t.Rem(0.76)).Bold().MinWidth(0).SingleLine()
			})
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
				glyph(c, "book-open", 14)
				count := fmt.Sprintf(o.Lessons, k.Lessons)
				if k.Lessons == 1 {
					count = o.Lesson
				}
				ui.Text(c, count).FontSize(t.Rem(0.72)).Shrink(0)
				dot(c)
				glyph(c, "clock", 14)
				span := fmt.Sprintf(o.Short, k.Minutes)
				if k.Minutes >= 60 {
					span = fmt.Sprintf(o.Time, k.Minutes/60, k.Minutes%60)
				}
				ui.Text(c, span).FontSize(t.Rem(0.72)).MinWidth(0).SingleLine()
			})
			ink := t.Accent
			if done {
				ink = t.Success
			}
			meter(c, float32(float64(k.Done)/float64(k.Lessons)), ink)
			ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
				c.Key(done)
				if done {
					ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).TextColor(t.Success).Children(func() {
						glyph(c, "circle-check", 14)
						ui.Text(c, o.Complete).FontSize(t.Rem(0.74)).Bold().Shrink(0)
					})
				} else {
					kit.Muted(c, fmt.Sprintf(o.Progress, percent(k.Done, k.Lessons))).FontSize(t.Rem(0.74)).Bold().Shrink(0)
				}
			})
		})
	})
	return root
}
