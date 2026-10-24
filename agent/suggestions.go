package agent

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Suggestion is one prompt chip.
type Suggestion struct {
	Text, Icon string
}

// Suggestions are prompt chips; a tap fills the box.
func Suggestions(c *ui.Context, label string, prompts []Suggestion, draft *string) ui.Element {
	if label == "" || len(prompts) == 0 || draft == nil {
		panic(fmt.Sprintf("agent: Suggestions %q needs a label, prompts and a draft", label))
	}
	seen := map[string]bool{}
	for _, p := range prompts {
		if p.Text == "" || p.Icon == "" || seen[p.Text] {
			panic(fmt.Sprintf("agent: Suggestions %q prompts need unique text and an icon", label))
		}
		seen[p.Text] = true
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).FocusGroup(ui.Horizontal).Role(ui.RoleGroup).Label(label)
	return root.Children(func() {
		for i, p := range prompts {
			c.Key(p.Text)
			on := *draft == p.Text
			b := kit.Press(ui.ButtonBase(c).Label(p.Text).Row().Gap(t.Space(1.5)).AlignItems(ui.Center).
				Padding(t.Space(1.5), t.Space(3)).Radius(999).MaxWidthPercent(100).Checked(on))
			kit.Lift(b)
			b.Border(1, b.AnimateColor("ring", t.Border.Mix(t.Accent, unit(on || b.Hovered())), ui.Snappy))
			b.Background(b.AnimateColor("bg", t.Background.Mix(t.Accent.Alpha(0.08), unit(on)), ui.Snappy))
			kit.Arrive(b, ui.Pose{Y: t.Space(2), Fade: 1}, time.Duration(i)*40*time.Millisecond)
			b.Children(func() {
				ui.Icon(c, icon.Must(p.Icon)).Size(16, 16).TextColor(t.Accent)
				ui.Text(c, p.Text).FontSize(t.Rem(0.9)).MinWidth(0).SingleLine()
			})
			if b.Clicked() && !on {
				*draft = p.Text
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("Suggestions", "filled", "index", i)
			}
		}
	})
}
