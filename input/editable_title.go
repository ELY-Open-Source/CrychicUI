package input

import "github.com/egoist/mygo/ui"

// EditableTitleOptions tune an EditableTitle.
type EditableTitleOptions struct {
	// Placeholder shows while the title is empty.
	Placeholder string
}

// EditableTitle is a heading-sized inline edit.
func EditableTitle(c *ui.Context, value *string, o EditableTitleOptions) ui.Element {
	if o.Placeholder == "" {
		o.Placeholder = "Untitled"
	}
	return inlineEdit(c, value, c.Theme().Rem(1.5), true, o.Placeholder)
}
