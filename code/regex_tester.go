package code

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RegexTesterOptions are the copy knobs of RegexTester.
type RegexTesterOptions struct {
	Pattern, Text, Empty, Invalid string
	// Matches words the count; At a match's bounds, "%d" each.
	Matches, At string
	FormatError func(err error) string
}

// RegexTester marks and lists a pattern's matches in a text.
func RegexTester(c *ui.Context, label string, pattern, text *string, o RegexTesterOptions) ui.Element {
	if label == "" || pattern == nil || text == nil {
		panic(fmt.Sprintf("code: RegexTester %q needs a label and bound fields", label))
	}
	def(&o.Pattern, "Pattern")
	def(&o.Text, "Text")
	def(&o.Empty, "Type some text")
	def(&o.Invalid, "Invalid")
	def(&o.Matches, "%d matches")
	def(&o.At, "%d–%d")
	if o.FormatError == nil {
		o.FormatError = func(err error) string { return "bad pattern: " + err.Error() }
	}
	t := c.Theme()
	// scan reads the bound values; inputs may change them mid-build.
	scan := func() (matches [][]int, err error) {
		re, err := regexp.Compile(*pattern)
		if err != nil {
			return nil, err
		}
		return re.FindAllStringIndex(*text, -1), nil
	}
	matches, err := scan()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		kit.CardHead(c, "regex", label, func() {
			if err != nil {
				kit.Chip(c, o.Invalid, t.Danger)
				return
			}
			kit.Chip(c, fmt.Sprintf(o.Matches, len(matches)), t.Accent)
		})
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Children(func() {
			captioned(c, o.Pattern, func() {
				ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
					kit.Muted(c, "/").Font("monospace").FontSize(t.Rem(0.9))
					pf := ui.TextInput(c, pattern).Label(o.Pattern).Font("monospace").Grow(1)
					kit.Muted(c, "/g").Font("monospace").FontSize(t.Rem(0.9))
					if pf.Changed() {
						edited(root, "RegexTester")
					}
				})
			})
			captioned(c, o.Text, func() {
				if ui.TextArea(c, text).Label(o.Text).Font("monospace").Lines(2, 6).FillWidth().Changed() {
					edited(root, "RegexTester")
				}
			})
			matches, err := scan()
			if err != nil {
				ui.Text(c, o.FormatError(err)).FontSize(t.Rem(0.85)).TextColor(t.Danger).Wrap().
					Padding(t.Space(2)).Radius(t.Radius).Background(t.Danger.Alpha(0.08))
				return
			}
			regexPreview(c, *text, matches, o.Empty)
			if len(matches) > 0 {
				ui.Scroll(c).FillWidth().MaxHeight(t.Space(40)).Children(func() {
					ui.Column(c).FillWidth().Gap(t.Space(0.5)).Children(func() {
						for i, m := range matches {
							c.Key(i)
							row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(0.75), t.Space(1))
							kit.Stagger(c, root, row, i).Children(func() {
								kit.Numeric(c, "#"+strconv.Itoa(i+1)).FontSize(t.Rem(0.75)).TextColor(t.TextMuted).Width(t.Space(7)).Shrink(0)
								ui.Text(c, strings.ReplaceAll((*text)[m[0]:m[1]], "\n", "⏎")).Font("monospace").FontSize(t.Rem(0.8)).
									TextColor(t.Accent).Padding(0, t.Space(1.5)).Radius(t.Radius).Background(t.Accent.Alpha(0.1)).Shrink(1).MinWidth(0).SingleLine()
								ui.Box(c).Grow(1)
								kit.Numeric(c, fmt.Sprintf(o.At, m[0], m[1])).Font("monospace").FontSize(t.Rem(0.7)).TextColor(t.TextMuted).Shrink(0)
							})
						}
					})
				})
			}
		})
	})
	return root
}

// regexPreview shows text by line with the matches marked.
func regexPreview(c *ui.Context, text string, matches [][]int, empty string) {
	t := c.Theme()
	ui.Column(c).FillWidth().Padding(t.Space(2)).Background(t.Surface).Radius(t.Radius).Border(1, t.Border).Children(func() {
		if text == "" {
			kit.Muted(c, empty).FontSize(t.Rem(0.85))
		}
		for line, start := 0, 0; start <= len(text); line++ {
			end := strings.IndexByte(text[start:], '\n')
			if end < 0 {
				end = len(text)
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
				lo, hi := max(m[0], start), min(m[1], end)
				if lo > pos {
					spans = append(spans, ui.Span{Text: text[pos:lo]})
				}
				spans = append(spans, ui.Span{Text: text[lo:hi], Background: t.Accent.Alpha(0.25), Radius: t.Space(0.5)})
				pos = hi
			}
			if pos < end {
				spans = append(spans, ui.Span{Text: text[pos:end]})
			}
			if len(spans) == 0 {
				spans = append(spans, ui.Span{Text: " "})
			}
			ui.RichText(c, spans...).Font("monospace").FontSize(t.Rem(0.85)).LineHeight(1.5).FillWidth().Wrap()
			start = end + 1
		}
	})
}
