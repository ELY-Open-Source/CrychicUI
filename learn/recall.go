package learn

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RecallOptions word the prompt and the four grades.
type RecallOptions struct {
	Prompt string
	Grades [4]string
}

// Recall grades how well a card was remembered; cue Select.
func Recall(c *ui.Context, due [4]string, grade *int, o RecallOptions) ui.Element {
	if grade == nil || *grade < -1 || *grade > 3 {
		panic("learn: Recall needs a grade of -1 to 3")
	}
	for _, d := range due {
		if strings.TrimSpace(d) == "" {
			panic(fmt.Sprintf("learn: Recall needs a next-review time for each grade, got %q", due))
		}
	}
	def(&o.Prompt, "How well did you remember?")
	for i, g := range [...]string{"Again", "Hard", "Good", "Easy"} {
		def(&o.Grades[i], g)
	}
	t := c.Theme()
	inks := [...]ui.Color{t.Danger, t.Warning, t.Accent, t.Success}
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		kit.Muted(c, o.Prompt).FontSize(t.Rem(0.82)).FillWidth().TextAlign(ui.Center).Wrap()
		row := ui.Row(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleRadioGroup)
		row.Children(func() {
			for i, name := range o.Grades {
				ink, on := inks[i], *grade == i
				c.Key(i)
				b := kit.Press(ui.ButtonBase(c).Label(name).Role(ui.RoleRadio).Checked(on).Grow(1).Basis(0).MinWidth(0).
					Column().Center().Gap(t.Space(0.5)).Padding(t.Space(2.5), t.Space(1)).Radius(t.Radius))
				b.Background(b.AnimateColor("bg", ink.Alpha(0.1).Mix(ink.Alpha(0.2), unit(b.Hovered())).Mix(ink, unit(on)), ui.Snappy))
				b.TextColor(b.AnimateColor("ink", ink.Mix(t.AccentText, unit(on)), ui.Snappy))
				kit.Stagger(c, row, b, i).Children(func() {
					ui.Text(c, name).FontSize(t.Rem(0.84)).Bold().FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine()
					ui.Text(c, due[i]).FontSize(t.Rem(0.66)).FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine().
						TextColor(t.TextMuted.Mix(t.AccentText, unit(on)))
				})
				if b.Clicked() && !on {
					*grade = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Recall", "graded", "grade", i)
				}
			}
		})
	})
	return root
}
