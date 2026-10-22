package navigation

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CommanderOptions tune a Commander.
type CommanderOptions struct {
	// Empty shows when no command matches.
	Empty string
}

// Commander is a searchable command list; the highlight slides.
func Commander(c *ui.Context, label string, commands []string, chosen *int, o CommanderOptions) ui.Element {
	if label == "" || len(commands) == 0 || chosen == nil || *chosen < -1 || *chosen >= len(commands) {
		panic(fmt.Sprintf("navigation: Commander %q needs commands and a choice or -1", label))
	}
	def(&o.Empty, "No matches")
	t := c.Theme()
	root := ui.Column(c).Width(t.Space(80)).Gap(t.Space(1)).Padding(t.Space(2)).Radius(t.Radius).Background(t.Background).
		Border(1, t.Border).Shadow(0, 2, 8, 0, kit.Shade(0.12))
	query := ui.Local(root, "query", func() string { return "" })
	hi := ui.Local(root, "hi", func() int { return 0 })
	pointed := ui.Local(root, "pointed", func() int { return -1 })
	id := root.ID()
	choose := func(i int) {
		*chosen = i
		root.SetChanged()
		sound.Play(sound.Select)
		kit.Log("Commander", "chose", "index", i)
	}
	return root.Children(func() {
		field := ui.SearchField(c, query).Label(label).FillWidth()
		if field.Changed() {
			*hi = 0
		}
		q := strings.ToLower(strings.TrimSpace(*query))
		var matches []int
		for i, cmd := range commands {
			if len(matches) < 8 && strings.Contains(strings.ToLower(cmd), q) {
				matches = append(matches, i)
			}
		}
		*hi = max(0, min(*hi, len(matches)-1))
		switch {
		case field.Shortcut(0, ui.KeyDown):
			*hi = min(*hi+1, len(matches)-1)
		case field.Shortcut(0, ui.KeyUp):
			*hi = max(*hi-1, 0)
		case field.Submitted() && len(matches) > 0:
			choose(matches[*hi])
		}
		if len(matches) == 0 {
			kit.Muted(c, o.Empty).Padding(t.Space(2)).Key("empty").Transition(kit.Enter)
		}
		// Rows first: the highlight settles before any is drawn.
		rows := make([]ui.Element, len(matches))
		now := -1
		for n, i := range matches {
			c.Key(i)
			rows[n] = ui.ButtonBase(c).Label(commands[i]).FillWidth().Radius(t.Radius)
			// Only a pointer arriving takes the highlight.
			if rows[n].Hovered() {
				if now = i; *pointed != i {
					*hi = n
				}
			}
		}
		*pointed = now
		for n, i := range matches {
			b := rows[n]
			b.Children(func() {
				stack(c, n == *hi, func() ui.Element {
					return ui.Box(c).Radius(t.Radius).Background(t.Accent.Alpha(0.12)).Transition(glide(id))
				}, func() ui.Element {
					return ui.Box(c).FillWidth().Padding(t.Space(1.5), t.Space(2.5)).Children(func() { ui.Text(c, commands[i]).SingleLine() })
				})
			})
			if b.Clicked() {
				choose(i)
			}
		}
	})
}
