package marketing

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Consent is a cookie choice; it counts once Decided.
type Consent struct{ Decided, Analytics, Marketing bool }

// CookieBannerOptions word the banner; Open takes the settings state.
type CookieBannerOptions struct {
	Open                                        *bool
	Text, Accept, Settings, Save, Saved, Change string
	Necessary, Analytics, Marketing, Always     string
}

// CookieBanner asks for cookie consent; cue Tap.
func CookieBanner(c *ui.Context, consent *Consent, o CookieBannerOptions) ui.Element {
	def(&o.Text, "We use cookies to keep you signed in and, with your consent, to measure and improve the site.")
	def(&o.Accept, "Accept all")
	def(&o.Settings, "Settings")
	def(&o.Save, "Save choices")
	def(&o.Saved, "Cookie choices saved")
	def(&o.Change, "Change")
	def(&o.Necessary, "Necessary")
	def(&o.Analytics, "Analytics")
	def(&o.Marketing, "Marketing")
	def(&o.Always, "Always on")
	if consent == nil {
		panic("marketing: CookieBanner needs a consent")
	}
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	open := o.Open
	if open == nil {
		open = ui.Local(root, "open", func() bool { return false })
	}
	decide := func(a, m bool, how string) {
		consent.Analytics, consent.Marketing, consent.Decided = a, m, true
		*open = false
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("CookieBanner", how, "analytics", a, "marketing", m)
	}
	root.Children(func() {
		if consent.Decided {
			c.Key("saved")
			row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center)
			kit.Arrive(row, ui.Pose{Y: t.Space(1.5), Fade: 1}, 0).Children(func() {
				glyph(c, "shield-check", 18).TextColor(t.Success)
				ui.Text(c, o.Saved).FontSize(t.Rem(0.84)).Bold().Grow(1).MinWidth(0).Wrap()
				c.Key("change")
				change := action(c, o.Change, "", false, t.Space(40))
				if change.Clicked() {
					consent.Decided, *open = false, true
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("CookieBanner", "reopened")
				}
			})
			return
		}
		c.Key("ask")
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			tile(c, "cookie", t.Accent, t.Space(10))
			kit.Muted(c, o.Text).FontSize(t.Rem(0.82)).Grow(1).MinWidth(0).Wrap()
		})
		c.Key("settings")
		layout.Collapse(c, *open, func() {
			ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(1), 0).Children(func() {
				for i, name := range [3]string{o.Necessary, o.Analytics, o.Marketing} {
					on := i == 0 || [2]bool{consent.Analytics, consent.Marketing}[i-1]
					b := ui.ButtonBase(c).Label(name).Checked(on).Disabled(i == 0).FillWidth().Row().Gap(t.Space(3)).AlignItems(ui.Center).
						Padding(t.Space(2), t.Space(2.5)).Radius(t.Radius)
					b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
					b.Children(func() {
						ui.Text(c, name).FontSize(t.Rem(0.84)).Bold().Grow(1).MinWidth(0).SingleLine()
						if i == 0 {
							kit.Muted(c, o.Always).FontSize(t.Rem(0.72)).Shrink(0)
						}
						track(c, on, b.Hovered())
					})
					if b.Clicked() && i > 0 {
						flag := &consent.Analytics
						if i == 2 {
							flag = &consent.Marketing
						}
						*flag = !*flag
						root.SetChanged()
						sound.Play(sound.Toggle)
						kit.Log("CookieBanner", "toggled", "index", i, "on", *flag)
					}
				}
			})
		})
		c.Key("actions")
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).Justify(ui.End).Children(func() {
			if *open {
				c.Key("save")
				if action(c, o.Save, "", false, t.Space(60)).Clicked() {
					decide(consent.Analytics, consent.Marketing, "saved")
				}
			} else {
				c.Key("open")
				if action(c, o.Settings, "settings-2", false, t.Space(60)).Clicked() {
					*open = true
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("CookieBanner", "settings")
				}
			}
			c.Key("accept")
			if action(c, o.Accept, "", true, t.Space(60)).Clicked() {
				decide(true, true, "accepted")
			}
		})
	})
	return root
}
