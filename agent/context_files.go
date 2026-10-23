package agent

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ContextFilesOptions tune the file chips.
type ContextFilesOptions struct {
	Remove string
}

// ContextFiles shows attached files as chips with remove.
func ContextFiles(c *ui.Context, label string, files *[]string, o ContextFilesOptions) ui.Element {
	if label == "" || files == nil {
		panic(fmt.Sprintf("agent: ContextFiles %q needs a label and files", label))
	}
	def(&o.Remove, "Remove")
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(1.5)).FocusGroup(ui.Horizontal).Role(ui.RoleGroup).Label(label)
	delAt := -1
	root.Children(func() {
		for i, f := range *files {
			if f == "" {
				panic(fmt.Sprintf("agent: ContextFiles %q files must not be empty", label))
			}
			c.Key(f)
			chip := ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2)).Radius(999).
				Background(t.Surface).Border(1, t.Border)
			kit.Arrive(chip, ui.Pose{Scale: 0.8, Fade: 1}, 0)
			var del ui.Element
			chip.Children(func() {
				ui.Icon(c, icon.Must("file-text")).Size(16, 16).TextColor(t.TextMuted)
				ui.Text(c, f).FontSize(t.Rem(0.85)).SingleLine()
				del = kit.Press(ui.ButtonBase(c).Label(o.Remove).Tooltip(o.Remove).Size(t.Space(5), t.Space(5)).Center().Radius(999))
				del.Background(del.AnimateColor("bg", t.Danger.Alpha(0.15*unit(del.Hovered())), ui.Snappy)).
					TextColor(t.TextMuted.Mix(t.Danger, unit(del.Hovered()))).
					Children(func() { ui.Icon(c, icon.Must("x")).Size(16, 16) })
			})
			if del.Clicked() {
				delAt = i
			}
		}
	})
	if delAt >= 0 {
		*files = append((*files)[:delAt], (*files)[delAt+1:]...)
		root.SetChanged()
		sound.Play(sound.Delete)
		kit.Log("ContextFiles", "removed", "left", len(*files))
	}
	return root
}
