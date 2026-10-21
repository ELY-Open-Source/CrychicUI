// Package icon embeds the Lucide icons, parsed on first use.
package icon

import (
	"embed"
	"slices"
	"strings"
	"sync"

	"github.com/egoist/mygo/ui"
)

//go:embed svg/*.svg
var files embed.FS

var (
	names       = list()
	mu          sync.Mutex
	cache       = map[string]*ui.SVG{}
	filledCache = map[string]*ui.SVG{}
)

func list() []string {
	entries, err := files.ReadDir("svg")
	if err != nil {
		panic(err)
	}
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = strings.TrimSuffix(e.Name(), ".svg")
	}
	slices.Sort(out)
	return out
}

// Must returns the icon named name; it panics when missing.
func Must(name string) *ui.SVG {
	mu.Lock()
	defer mu.Unlock()
	if s, ok := cache[name]; ok {
		return s
	}
	if _, ok := slices.BinarySearch(names, name); !ok {
		panic("icon: no icon " + name)
	}
	data, err := files.ReadFile("svg/" + name + ".svg")
	if err != nil {
		panic(err)
	}
	s := ui.MustParseSVG(data)
	cache[name] = s
	return s
}

// Filled returns a filled Lucide icon by name.
func Filled(name string) *ui.SVG {
	mu.Lock()
	defer mu.Unlock()
	if s, ok := filledCache[name]; ok {
		return s
	}
	if _, ok := slices.BinarySearch(names, name); !ok {
		panic("icon: no icon " + name)
	}
	data, err := files.ReadFile("svg/" + name + ".svg")
	if err != nil {
		panic(err)
	}
	data = []byte(strings.Replace(string(data), `fill="none"`, `fill="currentColor"`, 1))
	s := ui.MustParseSVG(data)
	filledCache[name] = s
	return s
}

// Names lists every icon name in order.
func Names() []string {
	return slices.Clone(names)
}
