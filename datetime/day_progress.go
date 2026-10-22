package datetime

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// DayProgressOptions tune a DayProgress.
type DayProgressOptions struct {
	// Format writes the elapsed percent; "%s of the day".
	Format string
}

// DayProgress fills a bar with the day elapsed.
func DayProgress(c *ui.Context, label string, o DayProgressOptions) ui.Element {
	if label == "" {
		panic("datetime: DayProgress needs a label")
	}
	def(&o.Format, "%s of the day")
	t := c.Theme()
	now := c.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	length := tomorrow(now).Sub(start)
	pct := int64(now.Sub(start) * 100 / length)
	c.After(start.Add((length*time.Duration(pct+1) + 99) / 100).Sub(now))
	text := fmt.Sprintf(o.Format, i18n.FormatNumber(locale(c), float64(pct)/100, i18n.NumberOptions{Percent: true}))
	root := ui.Column(c).FillWidth().Gap(t.Space(1.5))
	return root.Children(func() {
		ui.Row(c).FillWidth().Justify(ui.SpaceBetween).Children(func() {
			ui.Text(c, label).Bold()
			kit.Muted(c, text).ContentTransition(ui.ContentTransition{Numeric: true})
		})
		feedback.ProgressBar(c, label, float64(now.Sub(start))/float64(length))
	})
}
