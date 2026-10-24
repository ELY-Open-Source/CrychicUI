package project

import (
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// QuickAddOptions are the quick add labels.
type QuickAddOptions struct{ Placeholder, Add string }

// QuickAdd enters a task with parsed due; cue Pop.
func QuickAdd(c *ui.Context, title *string, due *time.Time, created *bool, o QuickAddOptions) (root, add ui.Element) {
	if title == nil || due == nil || created == nil {
		panic("project: QuickAdd needs title, due, and state")
	}
	def(&o.Placeholder, "Task and due")
	def(&o.Add, "Add")
	t := c.Theme()
	root = ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup)
	pending := ui.Local(root, "pending", func() bool { return false })
	root.Children(func() {
		in := ui.TextInput(c, title).Label(o.Placeholder).Placeholder(o.Placeholder).Grow(1)
		in.OnTextKey(func(k ui.TextKey) bool {
			if k.Key != ui.KeyEnter || strings.TrimSpace(k.Text()) == "" {
				return false
			}
			*pending = true
			return true
		})
		if in.Changed() {
			root.SetChanged()
		}
		add = kit.Press(ui.ButtonBase(c).Label(o.Add).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).
			Background(t.Accent).TextColor(t.AccentText).Disabled(strings.TrimSpace(*title) == ""))
		add.Children(func() { ui.Text(c, o.Add).SingleLine() })
	})
	if *pending || add.Clicked() {
		*pending = false
		*due = parseQuickDue(c.Now(), *title)
		*created = true
		add.SetChanged()
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("QuickAdd", "added", "hasDue", !due.IsZero())
	}
	return root, add
}

func parseQuickDue(now time.Time, title string) time.Time {
	text := strings.ToLower(strings.TrimSpace(title))
	for _, field := range strings.Fields(text) {
		switch field {
		case "today":
			return now
		case "tomorrow":
			return now.AddDate(0, 0, 1)
		default:
			if parsed, err := time.ParseInLocation("2006-01-02", field, now.Location()); err == nil {
				return parsed
			}
		}
	}
	return time.Time{}
}
