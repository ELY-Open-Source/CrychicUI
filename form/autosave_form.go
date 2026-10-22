package form

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// AutosaveFormOptions tune an AutosaveForm.
type AutosaveFormOptions struct {
	// Delay is the pause before saving; one second when zero.
	Delay time.Duration
	// Saving and Saved are status copy.
	Saving, Saved string
}

// AutosaveForm reports Changed a pause after content bumps *edits.
func AutosaveForm(c *ui.Context, edits *uint64, o AutosaveFormOptions, content func()) ui.Element {
	if edits == nil || content == nil || o.Delay < 0 {
		panic(fmt.Sprintf("form: AutosaveForm needs edits, content and a delay of 0 or more, not %v", o.Delay))
	}
	if o.Delay == 0 {
		o.Delay = time.Second
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(4))
	seen := ui.Local(root, "seen", func() uint64 { return *edits })
	last := ui.Local(root, "last", func() time.Time { return time.Time{} })
	saving := ui.Local(root, "saving", func() bool { return false })
	return root.Children(func() {
		content()
		now := c.Now()
		if *edits != *seen {
			*seen, *last, *saving = *edits, now, true
		}
		if *saving {
			if left := o.Delay - now.Sub(*last); left > 0 {
				c.After(left)
			} else {
				*saving = false
				root.SetChanged()
				kit.Log("AutosaveForm", "saved", "edits", *edits)
			}
		}
		if !last.IsZero() {
			feedback.SaveStatus(c, *saving, feedback.SaveStatusOptions{Saving: o.Saving, Saved: o.Saved})
		}
	})
}
