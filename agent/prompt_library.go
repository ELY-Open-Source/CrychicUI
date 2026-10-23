package agent

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PromptLibraryOptions tune a PromptLibrary.
type PromptLibraryOptions struct {
	Placeholder, Empty string
}

// PromptLibrary searches saved prompts; a click selects one.
func PromptLibrary(c *ui.Context, label string, prompts []string, query *string, chosen *int, o PromptLibraryOptions) (root, field ui.Element) {
	if label == "" || len(prompts) == 0 || query == nil || chosen == nil || *chosen < -1 || *chosen >= len(prompts) {
		panic(fmt.Sprintf("agent: PromptLibrary %q needs a label, prompts, a query and a choice or -1", label))
	}
	for _, p := range prompts {
		if p == "" {
			panic(fmt.Sprintf("agent: PromptLibrary %q prompts must not be empty", label))
		}
	}
	def(&o.Placeholder, "Search prompts")
	def(&o.Empty, "No prompts match")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		field = input.SearchInput(c, query, input.SearchInputOptions{Placeholder: o.Placeholder})
		q := strings.ToLower(strings.TrimSpace(*query))
		var hits []int
		for i, p := range prompts {
			if q == "" || strings.Contains(strings.ToLower(p), q) {
				hits = append(hits, i)
			}
		}
		if len(hits) == 0 {
			kit.Arrive(kit.Muted(c, o.Empty).Padding(t.Space(2)), ui.Pose{Fade: 1}, 0)
			return
		}
		ui.Scroll(c).FillWidth().MaxHeight(t.Space(60)).Gap(t.Space(0.5)).FocusGroup(ui.Vertical).Children(func() {
			for _, i := range hits {
				c.Key(i)
				on := i == *chosen
				b := kit.Press(ui.ButtonBase(c).Label(prompts[i]).Row().FillWidth().
					Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).Checked(on))
				bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
				if on {
					bg = t.Accent.Alpha(0.1)
				}
				b.Background(b.AnimateColor(fmt.Sprint("bg.", i), bg, ui.Snappy)).Children(func() {
					ui.Text(c, prompts[i]).Grow(1).SingleLine()
				})
				if b.Clicked() && !on {
					*chosen = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("PromptLibrary", "selected", "index", i)
				}
			}
		})
	})
	return root, field
}
