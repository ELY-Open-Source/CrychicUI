// Package gallery shows every component live.
package gallery

import (
	"embed"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"reflect"
	"runtime"
	"strings"
	"unicode"

	"github.com/egoist/mygo/ui"
)

// Entry is one component's page.
type Entry struct {
	ID       int
	Name     string
	Category string
	Summary  string
	Demo     func(c *ui.Context, log func(string))
}

// Category is a component package shown in the sidebar.
type Category struct {
	Key, Name, Icon string
}

var Categories = []Category{
	{"layout", "Layout", "layout-grid"},
	{"typography", "Typography", "type"},
	{"button", "Button", "mouse-pointer-click"},
	{"display", "Display", "id-card"},
	{"input", "Input", "text-cursor-input"},
	{"choice", "Choice", "list-checks"},
	{"feedback", "Feedback", "message-square-warning"},
	{"motion", "Motion", "move"},
	{"fx", "Effects", "sparkles"},
	{"form", "Form", "clipboard-list"},
	{"datetime", "Date & Time", "calendar-clock"},
	{"color", "Color", "palette"},
	{"navigation", "Navigation", "compass"},
	{"overlay", "Overlay", "layers"},
	{"data", "Data", "table"},
	{"chart", "Chart", "chart-column"},
	{"media", "Media", "image"},
	{"chat", "Chat", "message-circle"},
	{"agent", "Agent", "bot"},
	{"code", "Code", "code"},
	{"files", "Files", "folder"},
	{"shop", "Shop", "shopping-cart"},
	{"finance", "Finance", "wallet"},
	{"dashboard", "Dashboard", "layout-dashboard"},
	{"schedule", "Schedule", "calendar-days"},
	{"project", "Project", "kanban"},
	{"mail", "Mail", "mail"},
	{"social", "Social", "users"},
	{"account", "Account", "user-round"},
	{"devtools", "Devtools", "terminal"},
	{"game", "Game", "gamepad-2"},
	{"learn", "Learn", "graduation-cap"},
	{"health", "Health", "heart-pulse"},
	{"home", "Home", "house"},
	{"marketing", "Marketing", "megaphone"},
	{"content", "Content", "book-open"},
	{"desktop", "Desktop", "monitor"},
	{"geo", "Geo", "map-pin"},
}

// Entries lists every registered component in ID order.
func Entries() []Entry { return entries }

func category(key string) (Category, bool) {
	for _, k := range Categories {
		if k.Key == key {
			return k, true
		}
	}
	return Category{}, false
}

func inCategory(key string) []Entry {
	var out []Entry
	for _, e := range entries {
		if e.Category == key {
			out = append(out, e)
		}
	}
	return out
}

// Slug turns KPIGrid into kpi-grid.
func (e Entry) Slug() string {
	r := []rune(e.Name)
	var b strings.Builder
	for i, c := range r {
		if i > 0 && unicode.IsUpper(c) {
			prev := r[i-1]
			endsAcronym := unicode.IsUpper(prev) && i+1 < len(r) && unicode.IsLower(r[i+1])
			if unicode.IsLower(prev) || unicode.IsDigit(prev) || endsAcronym {
				b.WriteByte('-')
			}
		}
		b.WriteRune(unicode.ToLower(c))
	}
	return b.String()
}

func (e Entry) Path() string { return "/c/" + e.Category + "/" + e.Slug() }

//go:embed demos*.go
var demoFiles embed.FS

var codes = sources(demoFiles)

// Code is the source of the entry's demo function.
func (e Entry) Code() string {
	name := funcName(e.Demo)
	src, ok := codes[name]
	if !ok {
		panic("gallery: no source for " + name)
	}
	return src
}

func funcName(f any) string {
	n := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	return n[strings.LastIndex(n, ".")+1:]
}

// sources indexes every top-level function by name.
func sources(fsys fs.FS) map[string]string {
	names, err := fs.Glob(fsys, "*.go")
	if err != nil {
		panic(err)
	}
	out := map[string]string{}
	for _, f := range names {
		b, err := fs.ReadFile(fsys, f)
		if err != nil {
			panic(err)
		}
		// A Windows checkout may bring CRLF.
		src := strings.ReplaceAll(string(b), "\r\n", "\n")
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, f, src, 0)
		if err != nil {
			panic(err)
		}
		for _, d := range file.Decls {
			if fn, ok := d.(*ast.FuncDecl); ok && fn.Recv == nil {
				out[fn.Name.Name] = src[fset.Position(fn.Pos()).Offset:fset.Position(fn.End()).Offset]
			}
		}
	}
	return out
}
