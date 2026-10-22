package overlay

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

// Shortcut is an action and the keys that run it.
type Shortcut struct {
	Label string
	Keys  []string
}

// ShortcutGroup is a titled group of shortcuts.
type ShortcutGroup struct {
	Title     string
	Shortcuts []Shortcut
}

// ShortcutSheetOptions tune a ShortcutSheet.
type ShortcutSheetOptions struct {
	// Close names the close button.
	Close string
}

// ShortcutSheet lists shortcuts by group in a dialog.
func ShortcutSheet(c *ui.Context, title string, groups []ShortcutGroup, open *bool, o ShortcutSheetOptions) ui.Element {
	if title == "" || len(groups) == 0 || open == nil {
		panic(fmt.Sprintf("overlay: ShortcutSheet %q needs groups and open", title))
	}
	for _, g := range groups {
		if g.Title == "" || len(g.Shortcuts) == 0 {
			panic(fmt.Sprintf("overlay: ShortcutSheet %q groups need titles and shortcuts", title))
		}
		for _, s := range g.Shortcuts {
			if s.Label == "" || len(s.Keys) == 0 {
				panic(fmt.Sprintf("overlay: ShortcutSheet %q shortcuts need labels and keys", title))
			}
		}
	}
	def(&o.Close, "Close")
	t := c.Theme()
	root := ui.Box(c)
	if dialog(c, open, t.Space(150), func() {
		if heading(c, title, o.Close).Clicked() {
			*open = false
			root.SetChanged()
		}
		ui.Scroll(c).FillWidth().MaxHeight(t.Space(110)).Children(func() {
			ui.Row(c).FillWidth().Wrap().Gap(t.Space(6)).Children(func() {
				for g, grp := range groups {
					c.Key(g)
					col := ui.Column(c).Width(t.Space(60)).Gap(t.Space(2))
					kit.Arrive(col, ui.Pose{Y: t.Space(2), Fade: 1}, time.Duration(g)*40*time.Millisecond).Children(func() {
						kit.Muted(c, grp.Title).Bold().FontSize(t.Rem(0.75))
						for i, s := range grp.Shortcuts {
							c.Key(i)
							ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
								ui.Text(c, s.Label).Grow(1)
								ui.Row(c).Gap(t.Space(1)).Children(func() {
									for _, k := range s.Keys {
										typography.Kbd(c, k, false)
									}
								})
							})
						}
					})
				}
			})
		})
	}).Dismissed() {
		root.SetChanged()
	}
	switch opened, closed := moved(root, *open, "ShortcutSheet"); {
	case opened:
		sound.Play(sound.Open)
	case closed:
		sound.Play(sound.Close)
	}
	return root
}
