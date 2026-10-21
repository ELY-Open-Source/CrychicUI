package input

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/markdown"
)

// MarkdownEditorOptions tune a MarkdownEditor.
type MarkdownEditorOptions struct {
	// Placeholder is caller copy.
	Placeholder string
}

// MarkdownEditor edits Markdown beside a cross-fading preview.
func MarkdownEditor(c *ui.Context, source *string, o MarkdownEditorOptions) ui.Element {
	if source == nil {
		panic("input: MarkdownEditor needs a source")
	}
	if o.Placeholder == "" {
		o.Placeholder = "Write Markdown"
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Start)
	last := ui.Local(root, "last", func() string { return *source })
	previous := ui.Local(root, "previous", func() string { return *source })
	return root.Children(func() {
		editor := ui.TextArea(c, source).Grow(1).Basis(0).MinHeight(t.Rem(10)).Placeholder(o.Placeholder).Padding(t.Space(3)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		if editor.Changed() {
			root.SetChanged()
			sound.Play(sound.Tick)
		}
		preview := ui.Box(c).Grow(1).Basis(0).MinHeight(t.Rem(10)).Padding(t.Space(3)).Radius(t.Radius).Background(t.Surface)
		if *source != *last {
			*previous = *last
			*last = *source
			preview.Jump("fade", 0)
		}
		p := preview.AnimateTo("fade", 1, ui.Animation{Duration: 160 * time.Millisecond, Ease: ui.EaseOut})
		preview.Children(func() {
			current := ui.Box(c).Fill().Opacity(p)
			current.Children(func() { markdown.View(c, *source, nil) })
			if p < 1 {
				old := ui.Box(c).Absolute().Fill().Opacity(1 - p)
				old.Children(func() { markdown.View(c, *previous, nil) })
			}
		})
	})
}
