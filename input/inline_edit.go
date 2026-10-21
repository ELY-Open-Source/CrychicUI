package input

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// inlineEdit is the shared click-to-edit behavior.
func inlineEdit(c *ui.Context, value *string, size float32, bold bool, placeholder string) ui.Element {
	if value == nil {
		panic("input: inline edit needs a value")
	}
	t := c.Theme()
	root := ui.Box(c).FillWidth()
	editing := ui.Local(root, "editing", func() bool { return false })
	draft := ui.Local(root, "draft", func() string { return *value })
	saved := ui.Local(root, "saved", func() bool { return false })
	return root.Children(func() {
		if *saved {
			*saved = false
			root.SetChanged()
			sound.Play(sound.Success)
			kit.Log("InlineEdit", "saved", "count", len([]rune(*value)))
		}
		if !*editing {
			text := *value
			if text == "" {
				text = placeholder
			}
			button := ui.ButtonBase(c).Label("Edit text").FillWidth().Justify(ui.Start).Padding(t.Space(1), 0).TextColor(t.Text)
			kit.Press(button).Children(func() {
				e := ui.Text(c, text).FontSize(size)
				if bold {
					e.Bold()
				}
			})
			if button.Clicked() {
				*draft = *value
				*editing = true
			}
			return
		}
		in := ui.TextInput(c, draft).FillWidth().AutoFocus().Placeholder(placeholder).Padding(t.Space(1), 0)
		in.OnTextKey(func(k ui.TextKey) bool {
			switch k.Key {
			case ui.KeyEnter:
				*draft = k.Text()
				*value = *draft
				*editing = false
				*saved = true
				return true
			case ui.KeyEscape:
				*draft = *value
				*editing = false
				return true
			default:
				return false
			}
		})
	})
}

// InlineEditOptions tune an InlineEdit.
type InlineEditOptions struct {
	// Placeholder shows while the value is empty.
	Placeholder string
}

// InlineEdit turns text into an input and saves with Enter.
func InlineEdit(c *ui.Context, value *string, o InlineEditOptions) ui.Element {
	if o.Placeholder == "" {
		o.Placeholder = "Empty text"
	}
	return inlineEdit(c, value, c.Theme().FontSize, false, o.Placeholder)
}
