package learn

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Teacher is an instructor; Rating runs 0 to 5.
type Teacher struct {
	Name, Role, Bio   string
	Rating            float64
	Students, Courses int
	Photo             *ui.Bitmap
}

// InstructorOptions word the card's button and stats.
type InstructorOptions struct{ Follow, Following, Rating, Students, Courses string }

// Instructor shows a teacher, stats and follow; cue Toggle.
func Instructor(c *ui.Context, k Teacher, followed *bool, o InstructorOptions) ui.Element {
	if followed == nil || strings.TrimSpace(k.Name) == "" || !(k.Rating >= 0 && k.Rating <= 5) || k.Students < 0 || k.Courses < 0 {
		panic(fmt.Sprintf("learn: Instructor %q needs a rating of 0 to 5, counts of at least 0 and a state", k.Name))
	}
	def(&o.Follow, "Follow")
	def(&o.Following, "Following")
	def(&o.Rating, "Rating")
	def(&o.Students, "Students")
	def(&o.Courses, "Courses")
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			av := display.Avatar(c, k.Name, k.Photo).Size(t.Space(15), t.Space(15))
			kit.Arrive(av, ui.Pose{Scale: 0.6, Fade: 1}, 0)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
				ui.Text(c, k.Name).FontSize(t.Rem(1)).Bold().FillWidth().MinWidth(0).SingleLine()
				if k.Role != "" {
					kit.Muted(c, k.Role).FontSize(t.Rem(0.76)).FillWidth().MinWidth(0).SingleLine()
				}
				rating(c, k.Rating, t.Space(3.5))
			})
		})
		if k.Bio != "" {
			ui.Text(c, k.Bio).FontSize(t.Rem(0.8)).TextColor(t.Text.Mix(t.TextMuted, 0.3)).FillWidth().Wrap()
		}
		stats := ui.Row(c).FillWidth().Gap(t.Space(2))
		stats.Children(func() {
			cells := [3][2]string{{fmt.Sprintf("%.1f", k.Rating), o.Rating}, {number(c, k.Students), o.Students}, {number(c, k.Courses), o.Courses}}
			for i, cell := range cells {
				col := ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(0.5)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).
					Background(t.SurfaceHover.Alpha(0.6))
				kit.Stagger(c, stats, col, i).Children(func() {
					kit.Numeric(c, cell[0]).FontSize(t.Rem(0.95)).Bold().FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine()
					kit.Muted(c, cell[1]).FontSize(t.Rem(0.66)).FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine()
				})
			}
		})
		word, mark := o.Follow, "user-plus"
		if *followed {
			word, mark = o.Following, "user-check"
		}
		c.Key("follow")
		b := kit.Press(ui.ButtonBase(c).Label(word).FillWidth().Padding(t.Space(2.5), 0).Radius(999).Checked(*followed))
		b.Background(b.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(b.Hovered())).Mix(t.SurfaceHover, unit(*followed)), ui.Snappy))
		b.TextColor(b.AnimateColor("ink", t.AccentText.Mix(t.Text, unit(*followed)), ui.Snappy))
		b.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				c.Key(*followed)
				kit.Arrive(glyph(c, mark, 16), ui.Pose{Scale: 0.5, Fade: 1}, 0)
				ui.Text(c, word).FontSize(t.Rem(0.85)).Bold().Shrink(0)
			})
		})
		if b.Clicked() {
			*followed = !*followed
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("Instructor", "followed", "on", *followed)
		}
	})
	return root
}
