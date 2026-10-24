package account

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Providers chooses a sign-in provider; cue Tap.
func Providers(c *ui.Context, providers []string, chosen *string) ui.Element {
	if len(providers) == 0 || chosen == nil {
		panic("account: Providers needs providers and selection")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, provider := range providers {
			if strings.TrimSpace(provider) == "" {
				panic("account: provider must not be empty")
			}
			c.Key(provider)
			b := button(c, provider).FillWidth().Checked(*chosen == provider)
			b.Children(func() { ui.Text(c, provider) })
			if b.Clicked() {
				*chosen = provider
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("Providers", "selected")
			}
		}
	})
	return root
}
