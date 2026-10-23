package finance

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PayeesOptions are the payee labels.
type PayeesOptions struct{ Placeholder, Empty string }

// Payees filters saved payees; cue Select.
func Payees(c *ui.Context, names []string, query *string, selected *int, o PayeesOptions) ui.Element {
	if len(names) == 0 || query == nil || selected == nil || *selected < -1 || *selected >= len(names) {
		panic("finance: Payees needs names and selection")
	}
	def(&o.Placeholder, "Search payees")
	def(&o.Empty, "No payees match")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		field := input.SearchInput(c, query, input.SearchInputOptions{Placeholder: o.Placeholder})
		if field.Changed() {
			root.SetChanged()
		}
		q := strings.ToLower(*query)
		hits := 0
		for i, name := range names {
			if name == "" {
				panic("finance: Payees name must not be empty")
			}
			if q != "" && !strings.Contains(strings.ToLower(name), q) {
				continue
			}
			hits++
			c.Key(i)
			b := kit.Press(ui.ButtonBase(c).Label(name).Row().FillWidth().Padding(t.Space(2)).Radius(t.Radius).
				Border(1, t.Border).Checked(i == *selected))
			b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(i == *selected)), ui.Snappy))
			b.Children(func() { ui.Text(c, name).SingleLine() })
			if b.Clicked() && *selected != i {
				*selected = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Payees", "selected", "index", i)
			}
		}
		if hits == 0 {
			kit.Muted(c, o.Empty).FontSize(t.Rem(0.8))
		}
	})
	return root
}
