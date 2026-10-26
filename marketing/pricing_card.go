package marketing

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PricingCardOptions word the card; Current disables its button.
type PricingCardOptions struct {
	PlanWords
	Yearly, Current bool
}

// PricingCard is one plan with its features; cue Tap.
func PricingCard(c *ui.Context, p Plan, o PricingCardOptions) (root, choose ui.Element) {
	validPlan("PricingCard", p)
	o.fill()
	root, choose = planCard(c, p, o.Yearly, o.Current, o.PlanWords)
	if choose.Clicked() && !o.Current {
		sound.Play(sound.Tap)
		kit.Log("PricingCard", "chosen", "yearly", o.Yearly)
	}
	return root, choose
}
