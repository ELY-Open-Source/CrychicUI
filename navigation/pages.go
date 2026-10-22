package navigation

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
)

// Page is a destination's name and router path.
type Page struct{ Name, Path string }

// routed checks pages and finds the deepest one r shows.
func routed(component, label string, r *ui.Router, pages []Page) int {
	if label == "" || r == nil || len(pages) == 0 {
		panic(fmt.Sprintf("navigation: %s %q needs a router and pages", component, label))
	}
	at, path := -1, r.Path()
	for i, p := range pages {
		if p.Name == "" || !strings.HasPrefix(p.Path, "/") {
			panic(fmt.Sprintf("navigation: %s %q pages need a name and a path from /", component, label))
		}
		for _, q := range pages[:i] {
			if q.Path == p.Path {
				panic(fmt.Sprintf("navigation: %s %q lists %q twice", component, label, p.Path))
			}
		}
		inside := path == p.Path || strings.HasPrefix(path, strings.TrimSuffix(p.Path, "/")+"/")
		if inside && (at < 0 || len(p.Path) > len(pages[at].Path)) {
			at = i
		}
	}
	if at < 0 {
		panic(fmt.Sprintf("navigation: %s %q shows %q, under no page", component, label, path))
	}
	return at
}
