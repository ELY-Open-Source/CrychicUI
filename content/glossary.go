package content

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Term is a word and its meaning.
type Term struct{ Word, Meaning string }

// GlossaryOptions title the glossary.
type GlossaryOptions struct{ Title string }

// Glossary lists terms under their first letter.
func Glossary(c *ui.Context, terms []Term, o GlossaryOptions) ui.Element {
	def(&o.Title, "Glossary")
	words := make([]string, len(terms))
	for i, tm := range terms {
		if strings.TrimSpace(tm.Meaning) == "" {
			panic(fmt.Sprintf("content: Glossary term %q needs a meaning", tm.Word))
		}
		words[i] = tm.Word
	}
	unique("Glossary", words)
	if len(terms) == 0 {
		panic(fmt.Sprintf("content: Glossary %q needs terms", o.Title))
	}
	sorted := slices.Clone(terms)
	// Letter first keeps groups in letter order.
	slices.SortStableFunc(sorted, func(a, b Term) int {
		if x, y := letter(a.Word), letter(b.Word); x != y {
			return strings.Compare(x, y)
		}
		return strings.Compare(strings.ToLower(a.Word), strings.ToLower(b.Word))
	})
	var letters []string
	for _, tm := range sorted {
		if l := letter(tm.Word); !slices.Contains(letters, l) {
			letters = append(letters, l)
		}
	}
	t := c.Theme()
	root := card(c).Gap(t.Space(4)).Padding(t.Space(5)).Role(ui.RoleList)
	root.Children(func() {
		head(c, "book-open", o.Title)
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(1)).Children(func() {
			for _, l := range letters {
				ui.Text(c, l).FontSize(t.Rem(0.72)).Bold().TextColor(t.Accent).Width(t.Space(6)).TextAlign(ui.Center).
					Padding(t.Space(0.5), 0).Radius(t.Radius).Background(t.Accent.Alpha(0.1))
			}
		})
		groups := ui.Column(c).FillWidth().Gap(t.Space(4))
		groups.Children(func() {
			for i, l := range letters {
				c.Key(strconv.Quote(l))
				g := ui.Row(c).FillWidth().Gap(t.Space(4)).AlignItems(ui.Start)
				kit.Stagger(c, groups, g, i).Children(func() {
					ui.Text(c, l).FontSize(t.Rem(1.3)).Bold().TextColor(t.Accent).Width(t.Space(7)).Shrink(0)
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(3)).Children(func() {
						for _, tm := range sorted {
							if letter(tm.Word) != l {
								continue
							}
							c.Key(strconv.Quote(tm.Word))
							ui.Column(c).FillWidth().Gap(t.Space(0.5)).Role(ui.RoleListItem).Children(func() {
								ui.Text(c, tm.Word).FontSize(t.Rem(0.88)).Bold().FillWidth().Wrap()
								kit.Muted(c, tm.Meaning).FontSize(t.Rem(0.8)).LineHeight(1.5).FillWidth().Wrap()
							})
						}
					})
				})
			}
		})
	})
	return root
}

// letter is a word's group: its first letter, or "#".
func letter(word string) string {
	if r, _ := utf8.DecodeRuneInString(strings.TrimSpace(word)); unicode.IsLetter(r) {
		return string(unicode.ToUpper(r))
	}
	return "#"
}
