package gallery

import (
	"fmt"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/code"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

var codeEntries = []Entry{
	{901, "CodeBlock", "code", "Highlighted code with a language and copy.", demoCodeBlock},
	{902, "CodeEditor", "code", "MyGo's code editor with a status bar.", demoCodeEditor},
	{903, "DiffView", "code", "MyGo's diff viewer, split or unified.", demoDiffView},
	{904, "Terminal", "code", "Output lines and a blinking prompt.", demoTerminal},
	{905, "CommandLine", "code", "A prompt with history and output rows.", demoCommandLine},
	{906, "CodeTabs", "code", "Files as tabs over a code block.", demoCodeTabs},
	{907, "FileHeader", "code", "A path, language and actions bar.", demoFileHeader},
	{908, "Snippet", "code", "Short code with a title and copy.", demoSnippet},
	{909, "LineHighlight", "code", "Code with marked lines pulsing once.", demoLineHighlight},
	{910, "Annotated", "code", "Code with numbered notes beside lines.", demoAnnotatedCode},
	{911, "StackTrace", "code", "Frames with app frames bold; opens.", demoStackTrace},
	{912, "ErrorOverlay", "code", "A compile error panel sliding down.", demoErrorOverlay},
	{913, "Gutter", "code", "Line numbers with breakpoints toggled.", demoGutterCode},
	{914, "DebugBar", "code", "Continue, step over, in, out and stop.", demoDebugBar},
	{915, "Variables", "code", "A tree of names, types and values.", demoVariables},
	{916, "Watches", "code", "Expressions and values; add and remove.", demoWatches},
	{917, "CallStack", "code", "Frames with the current one marked.", demoCallStack},
	{918, "Problems", "code", "Errors and warnings by file.", demoProblems},
	{919, "Outline", "code", "Symbols of a file by kind.", demoOutlineCode},
	{920, "Minimap", "code", "A narrow file view with a viewport box.", demoMinimap},
}

const demoGoSource = "package main\n\nfunc main() {\n\tprintln(\"ship the deck\")\n}\n"

func demoCodeBlock(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() {
		_, cp := code.CodeBlock(c, "go", demoGoSource, code.CodeBlockOptions{})
		if cp.Clicked() {
			log(fmt.Sprintf("copied %d runes", utf8.RuneCountInString(demoGoSource)))
		}
	})
}

func demoCodeEditor(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	st := ui.Local(root, "editor", func() *editor.State {
		buf := editor.NewBuffer(demoGoSource)
		buf.SetLanguage(editor.LanguageFor("main.go"))
		return &editor.State{Buffer: buf}
	})
	root.Children(func() {
		if code.CodeEditor(c, *st, code.CodeEditorOptions{Height: 260}).Changed() {
			log("edited")
		}
	})
}

func demoDiffView(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	st := ui.Local(root, "diff", func() *editor.DiffState {
		oldBuf := editor.NewBuffer("package main\n\nfunc main() {\n\tprintln(\"deck\")\n}\n")
		newBuf := editor.NewBuffer(demoGoSource)
		oldBuf.SetLanguage(editor.LanguageFor("main.go"))
		newBuf.SetLanguage(editor.LanguageFor("main.go"))
		return &editor.DiffState{Old: oldBuf, New: newBuf}
	})
	split := boolState(root, "split", true)
	root.Children(func() {
		if code.DiffView(c, *st, split, code.DiffViewOptions{Height: 240}).Changed() {
			if *split {
				log("layout split")
			} else {
				log("layout unified")
			}
		}
	})
}

func demoTerminal(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	lines := ui.Local(root, "lines", func() []string { return []string{"building crychic…", "done in 1.2s"} })
	root.Children(func() {
		code.Terminal(c, "Shell", *lines, code.TerminalOptions{Height: 200})
		if button.Outline(c, "New line").Clicked() {
			*lines = append(*lines, fmt.Sprintf("line %d ok", len(*lines)-1))
			log(fmt.Sprintf("%d lines", len(*lines)))
		}
	})
}

func demoCommandLine(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	history := ui.Local(root, "history", func() []code.Command {
		return []code.Command{{Text: "ls", Output: "deck.md  notes.md"}}
	})
	draft := ui.Local(root, "draft", func() string { return "" })
	root.Children(func() {
		code.CommandLine(c, "Shell", *history, draft, func(s string) {
			*history = append(*history, code.Command{Text: s, Output: "ok"})
			log(fmt.Sprintf("ran %d runes", utf8.RuneCountInString(s)))
		}, code.CommandLineOptions{})
	})
}

func demoCodeTabs(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	current := intState(root, "current", 0)
	files := []code.CodeFile{
		{Name: "main.go", Lang: "go", Code: demoGoSource},
		{Name: "util.go", Lang: "go", Code: "package main\n\nfunc shout(s string) {\n\tprintln(s + \"!\")\n}\n"},
	}
	root.Children(func() {
		if code.CodeTabs(c, files, current).Changed() {
			log("selected " + files[*current].Name)
		}
	})
}

func demoFileHeader(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() {
		_, cp := code.FileHeader(c, "cmd/gallery/main.go", "Go", code.FileHeaderOptions{})
		if cp.Clicked() {
			log("copied the path")
		}
	})
}

func demoSnippet(c *ui.Context, log func(string)) {
	code2 := "http.ListenAndServe(\":8080\", nil)"
	ui.Column(c).Width(420).Children(func() {
		_, cp := code.Snippet(c, "Serve the gallery", "go", code2, code.SnippetOptions{})
		if cp.Clicked() {
			log(fmt.Sprintf("copied %d runes", utf8.RuneCountInString(code2)))
		}
	})
}

func demoLineHighlight(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() {
		code.LineHighlight(c, "go", demoGoSource, []int{4})
	})
}

func demoAnnotatedCode(c *ui.Context, log func(string)) {
	notes := []code.Note{{Line: 3, Text: "The entry point."}, {Line: 4, Text: "The whole program."}}
	ui.Column(c).Width(420).Children(func() { code.Annotated(c, "go", demoGoSource, notes) })
}

func demoStackTrace(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	open := boolState(root, "open", false)
	list := []code.StackFrame{
		{Func: "main.main", File: "main.go", Line: 12, App: true},
		{Func: "http.serve", File: "server.go", Line: 99},
		{Func: "runtime.goexit", File: "asm.s", Line: 1},
	}
	root.Children(func() {
		if code.StackTrace(c, list, open, code.StackTraceOptions{}).Changed() {
			log(fmt.Sprintf("open %v", *open))
		}
	})
}

func demoErrorOverlay(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	open := boolState(root, "open", false)
	root.Children(func() {
		box := code.ErrorOverlay(c, "build failed", "main.go:12: undefined: deck", open, code.ErrorOverlayOptions{}, func() {
			ui.Text(c, "The editor sits under the panel.")
		})
		if box.Changed() {
			log(fmt.Sprintf("open %v", *open))
		}
		if button.Outline(c, "Trigger").Clicked() {
			*open = true
			log("open true")
		}
	})
}

func demoGutterCode(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	breakpoints := ui.Local(root, "breakpoints", func() []int { return nil })
	root.Children(func() {
		if code.Gutter(c, 9, breakpoints).Changed() {
			log(fmt.Sprintf("%d breakpoints", len(*breakpoints)))
		}
	})
}

func demoDebugBar(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	running := boolState(root, "running", false)
	names := []string{"continue", "over", "into", "out", "stop"}
	root.Children(func() {
		_, cont, over, into, out, stop := code.DebugBar(c, *running, code.DebugBarOptions{})
		for i, b := range []ui.Element{cont, over, into, out, stop} {
			if b.Clicked() {
				log(names[i])
			}
		}
	})
}

func demoVariables(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	vars := []code.Variable{
		{Name: "cfg", Type: "Config", Value: "{…}", Children: []code.Variable{
			{Name: "port", Type: "int", Value: "8080"},
			{Name: "debug", Type: "bool", Value: "true"},
		}},
		{Name: "deck", Type: "string", Value: "\"launch\""},
	}
	root.Children(func() {
		if code.Variables(c, "Variables", vars).Changed() {
			log("toggled")
		}
	})
}

func demoWatches(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	watches := ui.Local(root, "watches", func() []code.Watch {
		return []code.Watch{{Expr: "len(deck)", Value: "12"}}
	})
	root.Children(func() {
		if code.Watches(c, "Watches", watches, code.WatchesOptions{}).Changed() {
			log(fmt.Sprintf("%d watches", len(*watches)))
		}
	})
}

func demoCallStack(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	current := intState(root, "current", 0)
	list := []code.StackFrame{
		{Func: "main.main", File: "main.go", Line: 12, App: true},
		{Func: "http.serve", File: "server.go", Line: 99},
		{Func: "runtime.goexit", File: "asm.s", Line: 1},
	}
	root.Children(func() {
		if code.CallStack(c, list, current).Changed() {
			log(fmt.Sprintf("selected %d", *current))
		}
	})
}

func demoProblems(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	selected := intState(root, "selected", -1)
	problems := []code.Problem{
		{Severity: code.ProbError, File: "main.go", Line: 12, Message: "undefined: deck"},
		{Severity: code.ProbWarning, File: "main.go", Line: 30, Message: "unused variable"},
		{Severity: code.ProbWarning, File: "util.go", Line: 4, Message: "missing return"},
	}
	root.Children(func() {
		if code.Problems(c, "Problems", problems, selected).Changed() {
			log("selected " + problems[*selected].Message)
		}
	})
}

func demoOutlineCode(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	selected := intState(root, "selected", -1)
	symbols := []code.Symbol{
		{Kind: code.SymType, Name: "Config", Line: 3},
		{Kind: code.SymFunc, Name: "Serve", Line: 20},
		{Kind: code.SymVar, Name: "cache", Line: 41},
		{Kind: code.SymConst, Name: "version", Line: 55},
	}
	root.Children(func() {
		if code.Outline(c, "Outline", symbols, selected).Changed() {
			log("selected " + symbols[*selected].Name)
		}
	})
}

func demoMinimap(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	at := intState(root, "at", 0)
	root.Children(func() {
		if code.Minimap(c, "Editor map", 200, 20, at, code.MinimapOptions{}).Changed() {
			log(fmt.Sprintf("at %d", *at))
		}
	})
}
