package code

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RegexTesterOptions are the copy knobs of RegexTester.
type RegexTesterOptions struct {
	Pattern     string
	Text        string
	Empty       string
	Matches     string
	FormatError func(err error) string
}

// RegexTester marks the matches of a pattern in a text.
func RegexTester(c *ui.Context, label string, pattern, text *string, o RegexTesterOptions) ui.Element {
	if label == "" || pattern == nil || text == nil {
		panic("code: RegexTester needs a label and bound fields")
	}
	def(&o.Pattern, "Pattern")
	def(&o.Text, "Text")
	def(&o.Empty, "Type some text")
	def(&o.Matches, "%d matches")
	if o.FormatError == nil {
		o.FormatError = func(err error) string { return "bad pattern: " + err.Error() }
	}
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Width(t.Space(90)).Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		pf := ui.TextInput(c, pattern).Label(o.Pattern).Font("monospace").FillWidth()
		tf := ui.TextArea(c, text).Label(o.Text).Font("monospace").Lines(2, 6).FillWidth()
		if pf.Changed() || tf.Changed() {
			root.SetChanged()
			sound.Play(sound.Tick)
			kit.Log("RegexTester", "edited")
		}
		re, err := regexp.Compile(*pattern)
		if err != nil {
			ui.Text(c, o.FormatError(err)).FontSize(t.Rem(0.85)).TextColor(t.Danger).Wrap()
			return
		}
		matches := re.FindAllStringIndex(*text, -1)
		box := ui.Column(c).FillWidth().Padding(t.Space(2)).Background(t.Surface).Radius(t.Radius)
		box.Children(func() {
			if *text == "" {
				kit.Muted(c, o.Empty).FontSize(t.Rem(0.85))
			}
			for line, start := 0, 0; start <= len(*text); line++ {
				end := strings.IndexByte((*text)[start:], '\n')
				if end < 0 {
					end = len(*text)
				} else {
					end += start
				}
				c.Key(line)
				spans := make([]ui.Span, 0, len(matches)*2+1)
				pos := start
				for _, m := range matches {
					if m[1] <= start {
						continue
					}
					if m[0] >= end {
						break
					}
					lo, hi := m[0], m[1]
					if lo < start {
						lo = start
					}
					if hi > end {
						hi = end
					}
					if lo > pos {
						spans = append(spans, ui.Span{Text: (*text)[pos:lo]})
					}
					spans = append(spans, ui.Span{Text: (*text)[lo:hi], Background: t.Accent.Alpha(0.3), Radius: t.Space(0.5)})
					pos = hi
				}
				if pos < end {
					spans = append(spans, ui.Span{Text: (*text)[pos:end]})
				}
				if len(spans) == 0 {
					spans = append(spans, ui.Span{Text: " "})
				}
				ui.RichText(c, spans...).Font("monospace").FontSize(t.Rem(0.85)).FillWidth().Wrap()
				start = end + 1
			}
		})
		kit.Numeric(c, fmt.Sprintf(o.Matches, len(matches))).FontSize(t.Rem(0.8)).TextColor(t.TextMuted)
	})
	return root
}
