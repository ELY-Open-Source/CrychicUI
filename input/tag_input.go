package input

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TagInputOptions tune a TagInput.
type TagInputOptions struct {
	// Placeholder is caller copy.
	Placeholder string
}

// TagInput adds chips with Enter or comma.
func TagInput(c *ui.Context, tags *[]string, draft *string, o TagInputOptions) ui.Element {
	if tags == nil || draft == nil {
		panic("input: TagInput needs tags and draft")
	}
	if o.Placeholder == "" {
		o.Placeholder = "Add tag"
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	pending := ui.Local(root, "pending", func() int { return 0 })
	pendingText := ui.Local(root, "pendingText", func() string { return "" })
	add := func(text string) {
		tag := strings.TrimSpace(text)
		if tag == "" {
			return
		}
		*tags = append(*tags, tag)
		*draft = ""
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("TagInput", "added", "count", len(*tags))
	}
	remove := func() {
		if len(*tags) == 0 {
			return
		}
		*tags = (*tags)[:len(*tags)-1]
		root.SetChanged()
		sound.Play(sound.Delete)
		kit.Log("TagInput", "removed", "count", len(*tags))
	}
	return root.Children(func() {
		if len(*tags) > 0 {
			ui.Row(c).Wrap().Gap(t.Space(2)).Children(func() {
				for i, tag := range *tags {
					chip := ui.Row(c).Key(i).Gap(t.Space(1)).Padding(t.Space(1), t.Space(2)).Radius(999).
						Background(t.Accent.Alpha(0.12)).TextColor(t.Accent)
					kit.Arrive(chip, ui.Pose{Scale: 0.7, Fade: 1}, 0).Children(func() { ui.Text(c, tag) })
				}
			})
		}
		c.Key("input")
		in := ui.TextInput(c, draft).FillWidth().Placeholder(o.Placeholder).Padding(t.Space(2), t.Space(3)).
			Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		in.OnTextKey(func(k ui.TextKey) bool {
			switch {
			case k.Key == ui.KeyEnter:
				*pendingText = k.Text()
				k.SetText("")
				*pending = 1
				return true
			case k.Key == ui.KeyBackspace && k.Text() == "":
				*pending = 2
				return true
			default:
				return false
			}
		})
		if *pending == 1 {
			add(*pendingText)
			*pendingText = ""
		} else if *pending == 2 {
			remove()
		}
		*pending = 0
		if changed(root, in) && strings.HasSuffix(*draft, ",") {
			add(strings.TrimSuffix(*draft, ","))
		}
	})
}
