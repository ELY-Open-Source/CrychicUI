package input

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SearchHistoryOptions tune a SearchHistory.
type SearchHistoryOptions struct {
	// Placeholder and Empty are caller copy.
	Placeholder, Empty string
}

// SearchHistory keeps recent searches below its input.
func SearchHistory(c *ui.Context, query *string, history *[]string, o SearchHistoryOptions) ui.Element {
	if query == nil || history == nil {
		panic("input: SearchHistory needs query and history")
	}
	if o.Placeholder == "" {
		o.Placeholder = "Search"
	}
	if o.Empty == "" {
		o.Empty = "No recent searches"
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	pending := ui.Local(root, "pending", func() bool { return false })
	return root.Children(func() {
		in := ui.TextInput(c, query).FillWidth().Placeholder(o.Placeholder).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		in.OnTextKey(func(k ui.TextKey) bool {
			if k.Key != ui.KeyEnter || k.Text() == "" {
				return false
			}
			*pending = true
			return true
		})
		if *pending {
			*pending = false
			*history = append([]string{*query}, slices.DeleteFunc(*history, func(s string) bool { return s == *query })...)
			root.SetChanged()
			sound.Play(sound.Pop)
			kit.Log("SearchHistory", "saved", "count", len(*history))
		}
		changed(root, in)
		if len(*history) == 0 {
			kit.Muted(c, o.Empty).FontSize(t.Rem(0.85))
			return
		}
		ui.Column(c).Gap(t.Space(1)).Children(func() {
			for i, item := range *history {
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
				kit.Arrive(row, ui.Pose{Fade: 1}, 0).Children(func() {
					choose := ui.ButtonBase(c).Label("Use "+item).Grow(1).Justify(ui.Start).Padding(t.Space(1), t.Space(2)).TextColor(t.Text)
					choose.Children(func() { ui.Text(c, item) })
					if kit.Clicked(choose, sound.Select) {
						*query = item
						root.SetChanged()
						kit.Log("SearchHistory", "selected", "index", i)
					}
					remove := ui.ButtonBase(c).Label("Remove " + item).Padding(t.Space(1)).TextColor(t.TextMuted)
					remove.Children(func() { ui.Icon(c, icon.Must("x")).Size(16, 16) })
					if kit.Clicked(remove, sound.Delete) {
						*history = append((*history)[:i], (*history)[i+1:]...)
						root.SetChanged()
						kit.Log("SearchHistory", "removed", "count", len(*history))
					}
				})
			}
		})
	})
}
