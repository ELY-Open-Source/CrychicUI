package archtest

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

const root = "../.."

var skipDirs = map[string]bool{".git": true, "node_modules": true, "dist": true, "icon/svg": true}

// Generated lock files are exempt.
var exempt = map[string]bool{"go.sum": true, "package-lock.json": true}

func walk(t *testing.T, visit func(rel string, src []byte)) {
	t.Helper()
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, p)
		rel = filepath.ToSlash(rel)
		if d.IsDir() {
			if skipDirs[rel] || skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		src, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		visit(rel, src)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileLength(t *testing.T) {
	walk(t, func(rel string, src []byte) {
		if exempt[filepath.Base(rel)] || !utf8.Valid(src) || bytes.IndexByte(src, 0) >= 0 {
			return
		}
		n := bytes.Count(src, []byte("\n"))
		if len(src) > 0 && src[len(src)-1] != '\n' {
			n++
		}
		if n >= 500 {
			t.Errorf("%s: %d lines, want under 500", rel, n)
		}
	})
}

func TestGoRules(t *testing.T) {
	fset, files := token.NewFileSet(), 0
	walk(t, func(rel string, src []byte) {
		if filepath.Ext(rel) != ".go" {
			return
		}
		files++
		f, err := parser.ParseFile(fset, rel, src, parser.ParseComments)
		if err != nil {
			t.Error(err)
			return
		}
		bad := checkComments(fset, f)
		if dir := filepath.ToSlash(filepath.Dir(rel)); dir != "internal/archtest" {
			bad = append(bad, checkImports(dir, f)...)
			if isCategory(dir) && !strings.HasSuffix(rel, "_test.go") {
				bad = append(bad, checkComponents(rel, f)...)
			}
		}
		for _, b := range bad {
			t.Error(b)
		}
	})
	if files == 0 {
		t.Fatal("no Go files found")
	}
}

func parse(t *testing.T, src string) (*token.FileSet, *ast.File) {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "x.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	return fset, f
}

func TestRulesCatch(t *testing.T) {
	const head = "package button\nimport \"github.com/egoist/mygo/ui\"\n"
	cases := []struct {
		name, file, src string
		rule            func(*token.FileSet, *ast.File, string) []string
		want            int
	}{
		{"same tier", "", "package button\nimport _ \"" + module + "/display\"\n", imports("button"), 1},
		{"upward", "", "package button\nimport _ \"" + module + "/form\"\n", imports("button"), 1},
		{"unknown pkg", "", "package x\n", imports("x"), 1},
		{"downward", "", "package button\nimport _ \"" + module + "/kit\"\n", imports("button"), 0},
		{"audio outside sound", "", "package media\nimport _ \"github.com/egoist/mygo/audio\"\n", imports("media"), 1},
		{"audio in sound", "", "package sound\nimport _ \"github.com/egoist/mygo/audio\"\n", imports("sound"), 0},
		{"long comment", "", "package x\n// one two three four five six seven eight nine ten eleven\n", comments, 1},
		{"two lines", "", "package x\n// one\n// two\n", comments, 1},
		{"block", "", "package x\n/* one\ntwo */\n", comments, 1},
		{"directive", "", "package x\n// Files.\n//go:generate echo\n", comments, 0},
		{"good", "button/icon_button.go", head + "func IconButton(c *ui.Context, o IconButtonOptions) ui.Element\n" +
			"type IconButtonOptions struct{}\n", components, 0},
		{"parts", "button/split.go", head + "func Split(c *ui.Context) (root, menu ui.Element)\n", components, 0},
		{"struct result", "button/split.go", head + "func Split(c *ui.Context) SplitResult\ntype SplitResult struct{}\n",
			components, 1},
		{"unnamed parts", "button/split.go", head + "func Split(c *ui.Context) (ui.Element, ui.Element)\n",
			components, 1},
		{"root not first", "button/split.go", head + "func Split(c *ui.Context) (menu, root ui.Element)\n",
			components, 1},
		{"two", "button/a.go", head + "func A(c *ui.Context) ui.Element\nfunc B(c *ui.Context) ui.Element\n",
			components, 1},
		{"file name", "button/b.go", head + "func KPIGrid(c *ui.Context) ui.Element\n", components, 1},
		{"no context", "button/a.go", head + "func A() ui.Element\n", components, 1},
		{"variadic", "button/a.go", head + "func A(c *ui.Context, o ...AOptions) ui.Element\n", components, 1},
		{"options not last", "button/a.go", head + "func A(c *ui.Context, o AOptions, n int) ui.Element\n",
			components, 1},
		{"returns", "button/a.go", head + "func A(c *ui.Context) int\n", components, 1},
		{"ctx name", "button/a.go", head + "func A(ctx *ui.Context) ui.Element\n", components, 1},
		{"content after options", "button/a.go", head + "func A(c *ui.Context, o AOptions, content func()) ui.Element\n",
			components, 0},
		{"two after options", "button/a.go", head + "func A(c *ui.Context, o AOptions, a, b func()) ui.Element\n",
			components, 1},
		{"grouped options", "button/a.go", head + "func A(c *ui.Context, x, y AOptions) ui.Element\n",
			components, 1},
	}
	for _, tc := range cases {
		fset, f := parse(t, tc.src)
		if got := tc.rule(fset, f, tc.file); len(got) != tc.want {
			t.Errorf("%s: %d findings %q, want %d", tc.name, len(got), got, tc.want)
		}
	}
}

func imports(dir string) func(*token.FileSet, *ast.File, string) []string {
	return func(_ *token.FileSet, f *ast.File, _ string) []string { return checkImports(dir, f) }
}

func comments(fset *token.FileSet, f *ast.File, _ string) []string { return checkComments(fset, f) }

func components(_ *token.FileSet, f *ast.File, file string) []string { return checkComponents(file, f) }

func TestSnake(t *testing.T) {
	for in, want := range map[string]string{"Stack": "stack", "KPIGrid": "kpi_grid", "Base64Tool": "base64_tool",
		"VDivider": "v_divider", "IDCard": "id_card", "PiPBox": "pi_p_box", "KPI": "kpi"} {
		if got := snake(in); got != want {
			t.Errorf("snake(%s) = %s, want %s", in, got, want)
		}
	}
}
