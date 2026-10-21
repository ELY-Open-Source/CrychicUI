package input

import "github.com/egoist/mygo/ui"

// CommandInput opens / commands and inserts a choice.
func CommandInput(c *ui.Context, label string, value *string, commands []string) ui.Element {
	return triggerInput(c, label, value, commands, '/', "CommandInput")
}
