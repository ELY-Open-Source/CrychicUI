package learn

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Week is one syllabus week; done weeks lead.
type Week struct {
	Title  string
	Topics []string
	Done   bool
}

// SyllabusOptions word the list; Week takes a number.
type SyllabusOptions struct{ Title, Week, Topics, Topic, Now string }

// Syllabus folds weeks open to their topics; cues Open, Close.
func Syllabus(c *ui.Context, weeks []Week, open *int, o SyllabusOptions) ui.Element {
	def(&o.Title, "Syllabus")
	if open == nil || len(weeks) == 0 || *open < -1 || *open >= len(weeks) {
		panic(fmt.Sprintf("learn: Syllabus %q needs weeks and an open week in range", o.Title))
	}
	titles := make([]string, len(weeks))
	current := len(weeks)
	for i, w := range weeks {
		if len(w.Topics) == 0 {
			panic(fmt.Sprintf("learn: Syllabus week %q needs topics", w.Title))
		}
		unique("Syllabus", w.Topics)
		if w.Done && current < i {
			panic(fmt.Sprintf("learn: Syllabus week %q is done after an open one", w.Title))
		}
		if !w.Done && current == len(weeks) {
			current = i
		}
		titles[i] = w.Title
	}
	unique("Syllabus", titles)
	def(&o.Week, "Week %d")
	def(&o.Topics, "%d topics")
	def(&o.Topic, "1 topic")
	def(&o.Now, "This week")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "calendar-range", o.Title, nil)
		list := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2))
		list.Children(func() {
			for i, w := range weeks {
				on := *open == i
				ink := t.TextMuted
				switch {
				case w.Done:
					ink = t.Success
				case i == current:
					ink = t.Accent
				}
				c.Key(strconv.Quote(w.Title))
				item := ui.Column(c).FillWidth().Radius(t.Radius)
				item.Background(item.AnimateColor("bg", t.Background.Mix(t.SurfaceHover.Alpha(0.6), unit(on)), ui.Snappy))
				kit.Stagger(c, list, item, i).Children(func() {
					b := kit.Press(ui.ButtonBase(c).Label(w.Title).FillWidth().Padding(t.Space(2.5)).Radius(t.Radius).AlignItems(ui.Stretch).
						Expanded(on))
					b.Background(b.AnimateColor("bg", t.Background.Alpha(0).Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
					b.Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
							badge := ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999).Background(ink.Alpha(0.14)).TextColor(ink)
							badge.Children(func() {
								c.Key(w.Done)
								if w.Done {
									glyph(c, "check", 16)
								} else {
									ui.Text(c, fmt.Sprint(i+1)).FontSize(t.Rem(0.8)).Bold().Shrink(0)
								}
							})
							ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
								kit.Muted(c, fmt.Sprintf(o.Week, i+1)).FontSize(t.Rem(0.66)).Bold().FillWidth().MinWidth(0).SingleLine()
								ui.Text(c, w.Title).FontSize(t.Rem(0.86)).Bold().FillWidth().MinWidth(0).SingleLine()
							})
							if i == current {
								kit.Chip(c, o.Now, t.Accent)
							}
							count := fmt.Sprintf(o.Topics, len(w.Topics))
							if len(w.Topics) == 1 {
								count = o.Topic
							}
							kit.Muted(c, count).FontSize(t.Rem(0.7)).Shrink(0)
							c.Key("chevron")
							chev := glyph(c, "chevron-down", 16).TextColor(t.TextMuted)
							chev.Rotate(chev.AnimateTo("turn", 180*unit(on), ui.Snappy))
						})
					})
					if b.Clicked() {
						cue := sound.Open
						*open = i
						if on {
							*open, cue = -1, sound.Close
						}
						root.SetChanged()
						sound.Play(cue)
						kit.Log("Syllabus", "toggled", "week", i, "open", !on)
					}
					layout.Collapse(c, on, func() {
						ui.Column(c).FillWidth().Gap(t.Space(1.5)).Padding(0, t.Space(3), t.Space(3), t.Space(13)).Children(func() {
							for _, topic := range w.Topics {
								ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
									ui.Box(c).Size(t.Space(1.5), t.Space(1.5)).Shrink(0).Radius(999).Background(ink)
									ui.Text(c, topic).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).Wrap()
								})
							}
						})
					})
				})
			}
		})
	})
	return root
}
