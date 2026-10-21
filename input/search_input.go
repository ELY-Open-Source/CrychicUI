package input

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SearchInputOptions tune a SearchInput.
type SearchInputOptions struct {
	// Placeholder is caller copy.
	Placeholder string
}

// SearchInput has a search icon, clear button and Escape clearing.
func SearchInput(c *ui.Context, value *string, o SearchInputOptions) ui.Element {
	if value == nil {
		panic("input: SearchInput needs a value")
	}
	if o.Placeholder == "" {
		o.Placeholder = "Search"
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(t.Space(2)).Padding(0, t.Space(2)).Radius(t.Radius).
		Background(t.Background).Border(1, t.Border)
	cleared := ui.Local(root, "cleared", func() bool { return false })
	return root.Children(func() {
		ui.Icon(c, icon.Must("search")).Size(16, 16).TextColor(t.TextMuted)
		in := ui.TextInput(c, value).Grow(1).Placeholder(o.Placeholder).Padding(t.Space(2), 0)
		in.OnTextKey(func(k ui.TextKey) bool {
			if k.Key != ui.KeyEscape || k.Text() == "" {
				return false
			}
			k.SetText("")
			*cleared = true
			return true
		})
		if *cleared {
			*cleared = false
			root.SetChanged()
			sound.Play(sound.Close)
			kit.Log("SearchInput", "cleared")
		}
		changed(root, in)
		if *value != "" {
			clear := ui.ButtonBase(c).Key("clear").Label("Clear search").Padding(t.Space(1)).TextColor(t.TextMuted)
			kit.Arrive(kit.Press(clear), ui.Pose{Scale: 0.5, Fade: 1}, 0).Children(func() {
				ui.Icon(c, icon.Must("x")).Size(16, 16)
			})
			if kit.Clicked(clear, sound.Close) {
				*value = ""
				root.SetChanged()
				kit.Log("SearchInput", "cleared")
			}
		}
	})
}
