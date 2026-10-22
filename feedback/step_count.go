package feedback

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// StepCountOptions tune a StepCount.
type StepCountOptions struct {
	// Format writes step and total; "Step %d of %d".
	Format string
}

// StepCount is "Step 2 of 5" over a filling bar.
func StepCount(c *ui.Context, step, total int, o StepCountOptions) ui.Element {
	if total < 1 || step < 0 || step > total {
		panic("feedback: StepCount needs 0 <= step <= total, total >= 1")
	}
	if o.Format == "" {
		o.Format = "Step %d of %d"
	}
	text := fmt.Sprintf(o.Format, step, total)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1.5))
	return root.Children(func() {
		kit.Numeric(c, text).TextColor(t.TextMuted).FontSize(t.Rem(0.85))
		ProgressBar(c, text, float64(step)/float64(total))
	})
}
