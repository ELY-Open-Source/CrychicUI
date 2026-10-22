package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Prop is one named value an Inspector lists.
type Prop struct{ Name, Value string }

// InspectorOptions tune an Inspector.
type InspectorOptions struct {
	// Close names the close button.
	Close string
}

// Inspector is a right panel of properties; it slides in.
func Inspector(c *ui.Context, title string, props []Prop, open *bool, o InspectorOptions) ui.Element {
	if title == "" || open == nil {
		panic(fmt.Sprintf("overlay: Inspector %q needs open", title))
	}
	for _, p := range props {
		if p.Name == "" {
			panic(fmt.Sprintf("overlay: Inspector %q props need names", title))
		}
	}
	def(&o.Close, "Close")
	t := c.Theme()
	root := ui.Box(c)
	from := t.Space(72)
	if c.Direction() == ui.RightToLeft {
		from = -from
	}
	if *open {
		ui.Overlay(c, func() {
			panel := ui.Column(c).Absolute().Top(0).Bottom(0).End(0).Width(t.Space(72)).Padding(t.Space(4)).Gap(t.Space(3)).
				Background(t.Background).Border(1, t.Border).Role(ui.RoleDialog).Label(title)
			kit.Arrive(panel, ui.Pose{X: from}, 0).Children(func() {
				if heading(c, title, o.Close).Clicked() {
					*open = false
					root.SetChanged()
				}
				for i, p := range props {
					c.Key(i)
					ui.Row(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(1.5), 0).Children(func() {
						kit.Muted(c, p.Name).Width(t.Space(24)).Shrink(0)
						ui.Text(c, p.Value).Grow(1)
					})
				}
			})
			if panel.OverlayShortcut(0, ui.KeyEscape) {
				*open = false
				root.SetChanged()
			}
		})
	}
	if opened, _ := moved(root, *open, "Inspector"); opened {
		sound.Play(sound.Open)
	}
	return root
}
