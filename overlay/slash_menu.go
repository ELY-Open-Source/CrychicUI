package overlay

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SlashMenuOptions tune a SlashMenu.
type SlashMenuOptions struct {
	// Empty shows when no command matches.
	Empty string
}

// SlashMenu pops icon commands under the caller's field.
func SlashMenu(c *ui.Context, label string, field ui.Element, query string, commands []MenuItem, chosen *int, open *bool, o SlashMenuOptions) ui.Element {
	items("SlashMenu", label, commands, chosen)
	if open == nil {
		panic("overlay: SlashMenu needs open")
	}
	def(&o.Empty, "No matching commands")
	root := ui.Box(c)
	suggest(c, root, field, label, query, len(commands), func(i int) string { return commands[i].Label },
		func(i int) string { return commands[i].Icon }, open, func(i int) {
			*chosen, *open = i, false
			root.SetChanged()
			sound.Play(sound.Select)
			kit.Log("SlashMenu", "chose", "index", i)
		}, o.Empty)
	moved(root, *open, "SlashMenu")
	return root
}
