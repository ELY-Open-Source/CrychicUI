package input

import (
	"fmt"
	"regexp"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RegexFieldOptions tune a RegexField.
type RegexFieldOptions struct {
	// Label, Error and Matches are caller copy.
	Label, Error, Matches string
}

// RegexField shows a match count for a pattern and sample.
func RegexField(c *ui.Context, pattern *string, sample string, o RegexFieldOptions) ui.Element {
	if pattern == nil {
		panic("input: RegexField needs a pattern")
	}
	if o.Label == "" {
		o.Label = "Pattern"
	}
	if o.Error == "" {
		o.Error = "Enter a valid regular expression."
	}
	if o.Matches == "" {
		o.Matches = "matches"
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		ui.Text(c, o.Label).FontWeight(500)
		in := ui.TextInput(c, pattern).FillWidth().Placeholder(o.Label).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).
			Background(t.Background).Border(1, t.Border)
		changed(root, in)
		re, err := regexp.Compile(*pattern)
		if err != nil {
			kit.Arrive(ui.Text(c, o.Error).TextColor(t.Danger).FontSize(t.Rem(0.85)), ui.Pose{Fade: 1}, 0)
			return
		}
		kit.Muted(c, fmt.Sprintf("%d %s", len(re.FindAllStringIndex(sample, -1)), o.Matches)).FontSize(t.Rem(0.85)).ContentTransition(ui.ContentTransition{})
	})
}
