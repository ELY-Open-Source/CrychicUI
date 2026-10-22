package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MentionMenuOptions tune a MentionMenu.
type MentionMenuOptions struct {
	// Empty shows when no one matches.
	Empty string
}

// MentionMenu pops people matching query under the caller's field.
func MentionMenu(c *ui.Context, label string, field ui.Element, query string, people []string, chosen *int, open *bool, o MentionMenuOptions) ui.Element {
	if label == "" || len(people) == 0 || chosen == nil || open == nil || *chosen < -1 || *chosen >= len(people) {
		panic(fmt.Sprintf("overlay: MentionMenu %q needs people, a choice or -1, and open", label))
	}
	for _, p := range people {
		if p == "" {
			panic(fmt.Sprintf("overlay: MentionMenu %q people need names", label))
		}
	}
	def(&o.Empty, "No one found")
	root := ui.Box(c)
	suggest(c, root, field, label, query, len(people), func(i int) string { return people[i] }, func(int) string { return "user" }, open,
		func(i int) {
			*chosen, *open = i, false
			root.SetChanged()
			sound.Play(sound.Select)
			kit.Log("MentionMenu", "chose", "index", i)
		}, o.Empty)
	moved(root, *open, "MentionMenu")
	return root
}
