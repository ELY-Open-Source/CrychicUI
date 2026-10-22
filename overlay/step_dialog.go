package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// StepDialogOptions tune a StepDialog.
type StepDialogOptions struct {
	Back, Next, Done string
	// Count formats the step and total.
	Count string
}

// StepDialog pages through steps; page builds step i.
func StepDialog(c *ui.Context, title string, steps []string, step *int, open *bool, page func(i int), o StepDialogOptions) ui.Element {
	if title == "" || len(steps) == 0 || step == nil || open == nil || page == nil || !valid(*step, len(steps)) {
		panic(fmt.Sprintf("overlay: StepDialog %q needs steps, a step, open and a page", title))
	}
	for _, s := range steps {
		if s == "" {
			panic(fmt.Sprintf("overlay: StepDialog %q steps need names", title))
		}
	}
	def(&o.Back, "Back")
	def(&o.Next, "Next")
	def(&o.Done, "Done")
	def(&o.Count, "Step %d of %d")
	t := c.Theme()
	root := ui.Box(c)
	shown, last := *step, len(steps)-1
	turn := func(to int) {
		*step = to
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("StepDialog", "step", "index", to)
	}
	if dialog(c, open, t.Space(110), func() {
		ui.Text(c, title).Bold().FontSize(t.Rem(1.1))
		ui.Row(c).FillWidth().Gap(t.Space(1)).Children(func() {
			for i := range steps {
				c.Key(i)
				bar := ui.Box(c).Grow(1).Height(t.Space(1)).Radius(999)
				bar.Background(bar.AnimateColor("bg", t.Border.Mix(t.Accent, unit(i <= shown)), ui.Snappy))
			}
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
			kit.Numeric(c, fmt.Sprintf(o.Count, shown+1, len(steps))).TextColor(t.TextMuted).FontSize(t.Rem(0.8))
			ui.Text(c, steps[shown]).Bold().FontSize(t.Rem(0.8)).ContentTransition(ui.ContentTransition{})
		})
		kit.Slide(c, shown, func() {
			ui.Column(c).FillWidth().Gap(t.Space(2)).Children(func() { page(shown) })
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).Justify(ui.End).Children(func() {
			if quiet(c.Key("back"), o.Back).Disabled(shown == 0).Clicked() {
				turn(shown - 1)
			}
			if shown < last {
				if primary(c.Key("next"), o.Next, t.Accent, t.AccentText).Clicked() {
					turn(shown + 1)
				}
				return
			}
			if primary(c.Key("done"), o.Done, t.Accent, t.AccentText).Clicked() {
				*open = false
				root.SetChanged()
				sound.Play(sound.Tick)
				kit.Log("StepDialog", "done")
			}
		})
	}).Dismissed() {
		root.SetChanged()
	}
	moved(root, *open, "StepDialog")
	return root
}
