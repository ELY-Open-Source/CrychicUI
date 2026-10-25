package account

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ProvidersOptions word the buttons; Continue takes the provider.
type ProvidersOptions struct{ Continue string }

// Providers chooses a sign-in provider; cue Tap.
func Providers(c *ui.Context, providers []string, chosen *string, o ProvidersOptions) ui.Element {
	if len(providers) == 0 || chosen == nil {
		panic("account: Providers needs providers and selection")
	}
	unique("Providers", providers)
	def(&o.Continue, "Continue with %s")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, provider := range providers {
			c.Key(provider)
			b := providerButton(c, provider, o.Continue, *chosen == provider)
			kit.Stagger(c, root, b, i)
			if b.Clicked() {
				*chosen = provider
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("Providers", "selected", "index", i)
			}
		}
	})
	return root
}

// providerButton is a full-width provider row with a letter tile.
func providerButton(c *ui.Context, provider, format string, on bool) ui.Element {
	t := c.Theme()
	inks := []ui.Color{t.Accent, t.Danger, t.Success, t.Warning}
	sum := 0
	for _, r := range provider {
		sum += int(r)
	}
	ink := inks[sum%len(inks)]
	b := kit.Press(ui.ButtonBase(c).Label(provider).FillWidth().Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Checked(on))
	b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent.Alpha(0.08), unit(on)), ui.Snappy))
	b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
	return b.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
			letter := strings.ToUpper(string([]rune(provider)[:1]))
			ui.Box(c).Size(t.Space(6), t.Space(6)).Shrink(0).Center().Radius(999).Background(ink.Alpha(0.14)).TextColor(ink).
				Children(func() { ui.Text(c, letter).FontSize(t.Rem(0.75)).Bold().Shrink(0) })
			ui.Text(c, fmt.Sprintf(format, provider)).FontSize(t.Rem(0.85)).Bold().Grow(1).MinWidth(0).SingleLine()
			if on {
				c.Key("check")
				kit.Arrive(glyph(c, "check", 16).TextColor(t.Accent), ui.Pose{Scale: 0.4, Fade: 1}, 0)
			}
		})
	})
}
