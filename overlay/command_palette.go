package overlay

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CommandPaletteOptions tune a CommandPalette.
type CommandPaletteOptions struct {
	// Empty shows when no command matches.
	Empty string
}

// CommandPalette opens on Cmd+K and filters commands as typed.
func CommandPalette(c *ui.Context, label string, commands []MenuItem, chosen *int, open *bool, o CommandPaletteOptions) ui.Element {
	items("CommandPalette", label, commands, chosen)
	if open == nil {
		panic("overlay: CommandPalette needs open")
	}
	def(&o.Empty, "No matching commands")
	t := c.Theme()
	root := ui.Box(c)
	query := searchText(root, *open)
	if c.Shortcut(ui.Cmd, ui.KeyK) && !*open {
		*open = true
		root.SetChanged()
	}
	if dialog(c, open, t.Space(120), func() {
		field := ui.SearchField(c, query).Label(label).FillWidth()
		results(c, root, field, *query, len(commands), func(i int) string { return commands[i].Label },
			func(int) string { return "" }, nil, func(i int) {
				*chosen, *open = i, false
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("CommandPalette", "chose", "index", i)
			}, o.Empty)
	}).Dismissed() {
		root.SetChanged()
	}
	if opened, _ := moved(root, *open, "CommandPalette"); opened {
		sound.Play(sound.Open)
	}
	return root
}
