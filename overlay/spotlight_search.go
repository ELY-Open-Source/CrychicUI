package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ResultGroup is a titled group of search results.
type ResultGroup struct {
	Title   string
	Results []string
}

// SpotlightSearchOptions tune a SpotlightSearch.
type SpotlightSearchOptions struct {
	// Empty shows when nothing matches.
	Empty string
}

// SpotlightSearch is a large centered search with grouped results.
func SpotlightSearch(c *ui.Context, label string, groups []ResultGroup, chosen *[2]int, open *bool, o SpotlightSearchOptions) ui.Element {
	if label == "" || len(groups) == 0 || chosen == nil || open == nil {
		panic(fmt.Sprintf("overlay: SpotlightSearch %q needs groups, a choice and open", label))
	}
	var flat [][2]int
	for g, grp := range groups {
		if grp.Title == "" {
			panic(fmt.Sprintf("overlay: SpotlightSearch %q groups need titles", label))
		}
		for i := range grp.Results {
			flat = append(flat, [2]int{g, i})
		}
	}
	def(&o.Empty, "No results")
	t := c.Theme()
	root := ui.Box(c)
	query := searchText(root, *open)
	if dialog(c, open, t.Space(150), func() {
		field := ui.SearchField(c, query).Label(label).FillWidth().FontSize(t.Rem(1.3))
		results(c, root, field, *query, len(flat), func(i int) string { return groups[flat[i][0]].Results[flat[i][1]] },
			func(i int) string { return groups[flat[i][0]].Title }, func(i int) {
				*chosen, *open = flat[i], false
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("SpotlightSearch", "chose", "group", flat[i][0], "index", flat[i][1])
			}, o.Empty)
	}).Dismissed() {
		root.SetChanged()
	}
	moved(root, *open, "SpotlightSearch")
	return root
}
