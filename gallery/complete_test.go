package gallery

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCatalogComplete(t *testing.T) {
	es := Entries()
	if len(es) != 1500 {
		t.Errorf("entries = %d, want 1500", len(es))
	}
	for i, e := range es {
		if e.ID != i+1 {
			t.Fatalf("entry %d (%s) has ID %d", i+1, e.Name, e.ID)
		}
	}
	shown := map[string]int{}
	for _, e := range es {
		shown[e.Category+"/"+e.Name]++
	}
	built, tested := map[string]bool{}, map[string]bool{}
	for _, k := range Categories {
		files, err := os.ReadDir(filepath.Join("..", k.Key))
		if err != nil {
			t.Fatalf("category %s: %v", k.Key, err)
		}
		for _, f := range files {
			if !strings.HasSuffix(f.Name(), ".go") {
				continue
			}
			src, err := parser.ParseFile(token.NewFileSet(), filepath.Join("..", k.Key, f.Name()), nil, parser.SkipObjectResolution)
			if err != nil {
				t.Fatal(err)
			}
			if strings.HasSuffix(f.Name(), "_test.go") {
				ast.Inspect(src, func(n ast.Node) bool {
					if id, ok := n.(*ast.Ident); ok {
						tested[k.Key+"/"+id.Name] = true
					}
					return true
				})
				continue
			}
			for _, d := range src.Decls {
				if fn, ok := d.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Name.IsExported() {
					built[k.Key+"/"+fn.Name.Name] = true
				}
			}
		}
	}
	for name := range built {
		if shown[name] != 1 {
			t.Errorf("%s has %d gallery entries, want 1", name, shown[name])
		}
		if !tested[name] {
			t.Errorf("%s is never named in its package's tests", name)
		}
	}
	for name := range shown {
		if !built[name] {
			t.Errorf("entry %s names no exported component", name)
		}
	}
	if len(built) != 1500 {
		t.Errorf("exported components = %d, want 1500", len(built))
	}
}

func TestShownPage(t *testing.T) {
	for path, want := range map[string]string{"/": "/", "/c/layout": "/c/layout", "/c/game/dice": "/c/game/dice",
		"/c/not-a-category": "", "/c/game/no-such": ""} {
		a, tt := open(path)
		tt.Frame()
		if got := a.Snapshot().Shown; got != want {
			t.Errorf("%s shows %q, want %q", path, got, want)
		}
	}
}
