package agent

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Suggestions are prompt chips; a tap fills the box.
func Suggestions(c *ui.Context, label string, prompts []string, chosen *int) ui.Element {
	if label == "" || len(prompts) == 0 || chosen == nil || *chosen < -1 || *chosen >= len(prompts) {
		panic(fmt.Sprintf("agent: Suggestions %q needs a label, prompts and a choice or -1", label))
	}
	for _, p := range prompts {
		if p == "" {
			panic(fmt.Sprintf("agent: Suggestions %q prompts must not be empty", label))
		}
	}
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(1.5)).FocusGroup(ui.Horizontal).Role(ui.RoleGroup).Label(label)
	return root.Children(func() {
		for i, p := range prompts {
			c.Key(i)
			b := kit.Press(ui.ButtonBase(c).Label(p).Padding(t.Space(1.5), t.Space(3)).Radius(999).Border(1, t.Border))
			b.Background(b.AnimateColor("bg", t.Surface.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
			kit.Arrive(b, ui.Pose{Y: t.Space(2), Fade: 1}, time.Duration(i)*40*time.Millisecond)
			b.Children(func() { ui.Text(c, p).SingleLine() })
			if b.Clicked() {
				*chosen = i
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("Suggestions", "picked", "index", i)
			}
		}
	})
}
