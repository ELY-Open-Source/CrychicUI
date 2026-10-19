// Package archtest checks the project's import and file rules.
package archtest

import (
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"strconv"
	"strings"
	"unicode"
)

const module = "github.com/ELY-Open-Source/CrychicUI"

var tiers = map[string]int{"sound": 0, "icon": 0, "kit": 1, "gallery": 6, "cmd/gallery": 7}

func init() {
	for t, pkgs := range map[int]string{
		2: "layout typography button display",
		3: "input choice feedback motion fx",
		4: "form datetime color navigation overlay data chart media",
		5: "chat agent code files shop finance dashboard schedule project mail social account devtools game " +
			"learn health home marketing content desktop geo",
	} {
		for _, p := range strings.Fields(pkgs) {
			tiers[p] = t
		}
	}
}

func isCategory(dir string) bool { t, ok := tiers[dir]; return ok && t >= 2 && t <= 5 }

// checkImports enforces downward imports between module packages.
func checkImports(dir string, f *ast.File) []string {
	from, ok := tiers[dir]
	if !ok {
		return []string{fmt.Sprintf("%s: package not in a tier", dir)}
	}
	var bad []string
	for _, im := range f.Imports {
		p, _ := strconv.Unquote(im.Path.Value)
		rel, ok := strings.CutPrefix(p, module+"/")
		if !ok {
			continue
		}
		to, known := tiers[rel]
		if !known || to >= from {
			bad = append(bad, fmt.Sprintf("%s (tier %d) imports %s", dir, from, rel))
		}
	}
	return bad
}

// checkComments allows one-line comments of ten words at most.
func checkComments(fset *token.FileSet, f *ast.File) []string {
	var bad []string
	for _, g := range f.Comments {
		lines := 0
		for _, c := range g.List {
			if strings.HasPrefix(c.Text, "//go:") {
				continue
			}
			at := fset.Position(c.Pos())
			lines++
			body := strings.TrimPrefix(strings.TrimPrefix(c.Text, "//"), "/*")
			if strings.Contains(c.Text, "\n") {
				bad = append(bad, fmt.Sprintf("%s: comment spans lines", at))
			}
			if n := len(strings.Fields(strings.TrimSuffix(body, "*/"))); n > 10 {
				bad = append(bad, fmt.Sprintf("%s: comment has %d words", at, n))
			}
			if lines == 2 {
				bad = append(bad, fmt.Sprintf("%s: comment spans lines", at))
			}
		}
	}
	return bad
}

// checkComponents enforces one component per file and its signature.
func checkComponents(file string, f *ast.File) []string {
	var bad []string
	var names []string
	for _, d := range f.Decls {
		fn, ok := d.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || !fn.Name.IsExported() {
			continue
		}
		names = append(names, fn.Name.Name)
		if err := signature(fn); err != "" {
			bad = append(bad, fmt.Sprintf("%s: %s %s", file, fn.Name.Name, err))
		}
	}
	switch {
	case len(names) > 1:
		bad = append(bad, fmt.Sprintf("%s: %d components %v", file, len(names), names))
	case len(names) == 1 && path.Base(file) != snake(names[0])+".go":
		bad = append(bad, fmt.Sprintf("%s: holds %s, want %s.go", file, names[0], snake(names[0])))
	}
	return bad
}

func signature(fn *ast.FuncDecl) string {
	name, ps := fn.Name.Name, fn.Type.Params.List
	if len(ps) == 0 || typeString(ps[0].Type) != "*ui.Context" || len(ps[0].Names) != 1 ||
		ps[0].Names[0].Name != "c" {
		return "must take c *ui.Context first"
	}
	last := ps[len(ps)-1]
	contentLast := typeString(last.Type) == "func()" && len(last.Names) == 1
	for i, p := range ps {
		ts := typeString(p.Type)
		if _, ok := p.Type.(*ast.Ellipsis); ok {
			return "must not be variadic"
		}
		if strings.HasSuffix(ts, "Options") && (ts != name+"Options" || i != len(ps)-1 && !(i == len(ps)-2 && contentLast) ||
			len(p.Names) > 1) {
			return "options must be a final " + name + "Options"
		}
	}
	rs := fn.Type.Results
	if rs == nil || rs.NumFields() == 0 {
		return "must return ui.Element"
	}
	for _, r := range rs.List {
		if typeString(r.Type) != "ui.Element" {
			return "must return only ui.Element results"
		}
	}
	if rs.NumFields() > 1 && (len(rs.List[0].Names) == 0 || rs.List[0].Names[0].Name != "root") {
		return "several results must be named, root first"
	}
	return ""
}

func typeString(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeString(t.X)
	case *ast.SelectorExpr:
		return typeString(t.X) + "." + t.Sel.Name
	case *ast.FuncType:
		if t.Params.NumFields() == 0 && t.Results.NumFields() == 0 {
			return "func()"
		}
	}
	return fmt.Sprintf("%T", e)
}

// snake turns KPIGrid into kpi_grid.
func snake(s string) string {
	r := []rune(s)
	var b strings.Builder
	for i, c := range r {
		if i > 0 && unicode.IsUpper(c) {
			prev := r[i-1]
			endsAcronym := unicode.IsUpper(prev) && i+1 < len(r) && unicode.IsLower(r[i+1])
			if unicode.IsLower(prev) || unicode.IsDigit(prev) || endsAcronym {
				b.WriteByte('_')
			}
		}
		b.WriteRune(unicode.ToLower(c))
	}
	return b.String()
}
