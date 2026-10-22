package data

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"
)

// JSONTreeOptions tune a JSONTree.
type JSONTreeOptions struct {
	// Expand names an arrow; it takes the key.
	Expand string
}

// pointer escapes a key as a JSON Pointer token.
var pointer = strings.NewReplacer("~", "~0", "/", "~1")

// JSONTree shows decoded JSON colored by type; open holds pointers.
func JSONTree(c *ui.Context, label string, value any, open *ui.Selection[string], o JSONTreeOptions) ui.Element {
	if label == "" || open == nil {
		panic(fmt.Sprintf("data: JSONTree %q needs open paths", label))
	}
	def(&o.Expand, "Expand %s")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Role(ui.RoleTree).Label(label)
	var node func(key string, v any, path string, depth int)
	node = func(key string, v any, path string, depth int) {
		var keys []string
		var at func(k string) any
		switch v := v.(type) {
		case map[string]any:
			for k := range v {
				keys = append(keys, k)
			}
			slices.Sort(keys)
			at = func(k string) any { return v[k] }
		case []any:
			for i := range v {
				keys = append(keys, strconv.Itoa(i))
			}
			at = func(k string) any { i, _ := strconv.Atoi(k); return v[i] }
		}
		c.Key(path)
		ui.Column(c).FillWidth().Children(func() {
			branch(c, root, "JSONTree", fmt.Sprintf(o.Expand, key), path, depth, at != nil, open, func() {
				ui.Text(c, key).Bold().SingleLine()
				shown, ink := scalar(t, v)
				if at != nil {
					shown, ink = fmt.Sprintf("{%d}", len(keys)), t.TextMuted
					if _, ok := v.([]any); ok {
						shown = fmt.Sprintf("[%d]", len(keys))
					}
				}
				ui.Text(c, shown).TextColor(ink).SingleLine()
			})
			if at != nil {
				kids(c, open.Has(path), func() {
					for _, k := range keys {
						node(k, at(k), path+"/"+pointer.Replace(k), depth+1)
					}
				})
			}
		})
	}
	root.Children(func() { node("$", value, "$", 0) })
	return root
}
