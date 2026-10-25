package account

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Step is one getting-started task; Detail is optional.
type Step struct {
	ID, Title, Detail string
	Done              bool
}

// OnboardingOptions word the card; Count takes done and total.
type OnboardingOptions struct{ Title, Count, Finished string }

// Onboarding is a getting-started checklist; cue Toggle.
func Onboarding(c *ui.Context, steps *[]Step, o OnboardingOptions) ui.Element {
	if steps == nil || len(*steps) == 0 {
		panic("account: Onboarding needs steps")
	}
	ids := make([]string, len(*steps))
	done := 0
	for i, s := range *steps {
		if s.Title == "" {
			panic(fmt.Sprintf("account: Onboarding %q needs a title", s.ID))
		}
		ids[i] = s.ID
		if s.Done {
			done++
		}
	}
	unique("Onboarding", ids)
	def(&o.Title, "Get started")
	def(&o.Count, "%d of %d done")
	def(&o.Finished, "You're all set")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		ui.Column(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(4), t.Space(4), t.Space(2), t.Space(4)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
				tile(c, "rocket", t.Accent, t.Space(10))
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
					ui.Text(c, o.Title).FontSize(t.Rem(1.02)).Bold().Wrap().FillWidth()
					kit.Numeric(c, fmt.Sprintf(o.Count, done, len(*steps))).FontSize(t.Rem(0.78)).TextColor(t.TextMuted)
				})
			})
			ink := t.Accent
			if done == len(*steps) {
				ink = t.Success
			}
			meter(c, float32(done)/float32(len(*steps)), ink)
			if done == len(*steps) {
				c.Key("finished")
				note(c, "party-popper", o.Finished, t.Success)
			}
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			for i, s := range *steps {
				c.Key(s.ID)
				b := kit.Press(ui.ButtonBase(c).Label(s.Title).FillWidth().Padding(t.Space(2)).Radius(t.Radius).Checked(s.Done))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
				kit.Stagger(c, list, b, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Start).Children(func() {
						tick(c, s.Done, b.Hovered())
						ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
							title := ui.Text(c, s.Title).FontSize(t.Rem(0.85)).Bold().Wrap().FillWidth()
							title.TextColor(title.AnimateColor("ink", t.Text.Mix(t.TextMuted, unit(s.Done)), ui.Snappy))
							if s.Done {
								title.Strikethrough()
							}
							if s.Detail != "" {
								kit.Muted(c, s.Detail).FontSize(t.Rem(0.75)).Wrap().FillWidth()
							}
						})
					})
				})
				if b.Clicked() {
					(*steps)[i].Done = !s.Done
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("Onboarding", "toggled", "index", i, "done", !s.Done)
				}
			}
		})
	})
	return root
}
