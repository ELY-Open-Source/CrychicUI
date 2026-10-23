package gallery

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/code"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

var codeEntries2 = []Entry{
	{ID: 921, Name: "Breadcrumb", Category: "code", Summary: "File and symbol path; cue Tap.", Demo: demoCodeBreadcrumb},
	{ID: 922, Name: "SymbolSearch", Category: "code", Summary: "Symbols filtered as typed; cue Select.", Demo: demoCodeSymbolSearch},
	{ID: 923, Name: "RegexTester", Category: "code", Summary: "Pattern, text, matches marked.", Demo: demoRegexTester},
	{ID: 924, Name: "JSONFormat", Category: "code", Summary: "Formats or minifies JSON; cue Success, Error.", Demo: demoJSONFormat},
	{ID: 925, Name: "Base64Tool", Category: "code", Summary: "Encode and decode panes; cue Tick.", Demo: demoBase64Tool},
	{ID: 926, Name: "HashTool", Category: "code", Summary: "SHA-256 of text; digest cross-fades.", Demo: demoHashTool},
	{ID: 927, Name: "UnixTime", Category: "code", Summary: "Epoch and date converter; cue Tick.", Demo: demoUnixTime},
	{ID: 928, Name: "Commit", Category: "code", Summary: "Hash, message, author and time row; cue Pop on copy.", Demo: demoCodeCommit},
	{ID: 929, Name: "BranchPicker", Category: "code", Summary: "Branches with search and current mark; cue Select.", Demo: demoBranchPicker},
	{ID: 930, Name: "PullRequest", Category: "code", Summary: "Title, status, checks and reviewers.", Demo: demoPullRequest},
	{ID: 931, Name: "Checks", Category: "code", Summary: "CI checks with spinners and results.", Demo: demoCodeChecks},
	{ID: 932, Name: "Blame", Category: "code", Summary: "Lines with author and age gutter.", Demo: demoBlame},
	{ID: 933, Name: "CommitGraph", Category: "code", Summary: "Lanes and dots of commits.", Demo: demoCommitGraph},
	{ID: 934, Name: "Language", Category: "code", Summary: "Language bar of shares; segments grow.", Demo: demoCodeLanguage},
	{ID: 935, Name: "KeyTable", Category: "code", Summary: "Keyboard shortcuts in key caps by group.", Demo: demoKeyTable},
	{ID: 936, Name: "Console", Category: "code", Summary: "Log, warn, error rows with filter; cue Toggle.", Demo: demoCodeConsole},
	{ID: 937, Name: "REPL", Category: "code", Summary: "Input evaluating arithmetic with history; cue Send.", Demo: demoREPL},
	{ID: 938, Name: "PackageCard", Category: "code", Summary: "Name, version, downloads, install copy; cue Pop.", Demo: demoPackageCard},
	{ID: 939, Name: "Env", Category: "code", Summary: "KEY=value rows with masked values; cue Toggle.", Demo: demoEnv},
	{ID: 940, Name: "Playground", Category: "code", Summary: "Code editor and output pane; run; cue Tap.", Demo: demoPlayground},
}

func demoCodeBreadcrumb(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	at := intState(root, "at", 3)
	root.Children(func() {
		if code.Breadcrumb(c, []string{"cmd", "gallery", "demos_code.go", "demoCodeBlock"}, at).Changed() {
			log("at " + []string{"cmd", "gallery", "demos_code.go", "demoCodeBlock"}[*at])
		}
	})
}

func demoCodeSymbolSearch(c *ui.Context, log func(string)) {
	symbols := []code.Symbol{
		{Kind: code.SymFunc, Name: "runServer", Line: 12},
		{Kind: code.SymType, Name: "Server", Line: 30},
		{Kind: code.SymVar, Name: "port", Line: 5},
		{Kind: code.SymFunc, Name: "handleRequest", Line: 44},
	}
	root := ui.Column(c)
	query, chosen := textState(root, "query", ""), intState(root, "chosen", -1)
	root.Children(func() {
		r, field := code.SymbolSearch(c, "Symbols", symbols, query, chosen, code.SymbolSearchOptions{})
		if r.Changed() {
			log("selected " + fmt.Sprint(*chosen))
		}
		if field.Changed() {
			log("query " + fmt.Sprint(utf8.RuneCountInString(*query)) + " chars")
		}
	})
}

func demoRegexTester(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	pattern := textState(root, "pattern", "[a-z]+")
	text := textState(root, "text", "ship the 20 MiB build")
	root.Children(func() {
		if code.RegexTester(c, "Regex", pattern, text, code.RegexTesterOptions{}).Changed() {
			log("edited")
		}
	})
}

func demoJSONFormat(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	text := textState(root, "text", `{"name":"crychic","count":940}`)
	root.Children(func() {
		r, _, _ := code.JSONFormat(c, "JSON", text, code.JSONFormatOptions{})
		if r.Changed() {
			log("action on " + fmt.Sprint(utf8.RuneCountInString(*text)) + " runes")
		}
	})
}

func demoBase64Tool(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	text := textState(root, "text", "hello world")
	encode := boolState(root, "encode", true)
	root.Children(func() {
		if code.Base64Tool(c, "Base64", text, encode, code.Base64ToolOptions{}).Changed() {
			log("encode " + fmt.Sprint(*encode))
		}
	})
}

func demoHashTool(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	text := textState(root, "text", "crychic")
	root.Children(func() {
		if code.HashTool(c, "SHA-256", text).Changed() {
			log("edited")
		}
	})
}

func demoUnixTime(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	epoch := int64State(root, "epoch", 1700000000)
	root.Children(func() {
		if code.UnixTime(c, "Unix time", epoch, code.UnixTimeOptions{}).Changed() {
			log("epoch " + fmt.Sprint(*epoch))
		}
	})
}

func demoCodeCommit(c *ui.Context, log func(string)) {
	_, copy := code.Commit(c, "1849e62abcdef", "Add code components 0901-0920", "Mei", time.Now().Add(-2*time.Hour), code.CommitOptions{})
	if copy.Clicked() {
		log("copied hash")
	}
}

func demoBranchPicker(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	current, query := intState(root, "current", 0), textState(root, "query", "")
	root.Children(func() {
		r, _ := code.BranchPicker(c, "Branches", []string{"main", "dev", "feature/review", "release/1.4"}, current, query, code.BranchPickerOptions{})
		if r.Changed() {
			log("selected " + []string{"main", "dev", "feature/review", "release/1.4"}[*current])
		}
	})
}

func demoPullRequest(c *ui.Context, log func(string)) {
	checks := []code.Check{{Name: "build", State: code.CheckPass}, {Name: "test", State: code.CheckPending}, {Name: "lint", State: code.CheckPass}}
	code.PullRequest(c, "Add the code category", code.PROpen, checks, []string{"Ana Reyes", "Kai Mori"}, code.PullRequestOptions{})
}

func demoCodeChecks(c *ui.Context, log func(string)) {
	code.Checks(c, "CI", []code.Check{
		{Name: "build", State: code.CheckPass},
		{Name: "test", State: code.CheckPass},
		{Name: "deploy", State: code.CheckPending},
	})
}

func demoBlame(c *ui.Context, log func(string)) {
	code.Blame(c, "Blame", []code.BlameLine{
		{Author: "Mei", Text: "package main", At: time.Now().Add(-96 * time.Hour)},
		{Author: "Kai", Text: "", At: time.Now().Add(-48 * time.Hour)},
		{Author: "Mei", Text: "func main() {", At: time.Now().Add(-3 * time.Hour)},
		{Author: "Ana", Text: "\trun()", At: time.Now().Add(-30 * time.Minute)},
	}, code.BlameOptions{})
}

func demoCommitGraph(c *ui.Context, log func(string)) {
	code.CommitGraph(c, "Graph", []code.GraphCommit{
		{Hash: "aaaaaaa1", Message: "Merge the review fixes", Lane: 0},
		{Hash: "bbbbbbb2", Message: "Polish the minimap", Lane: 1},
		{Hash: "ccccccc3", Message: "Add the code category", Lane: 0},
		{Hash: "ddddddd4", Message: "Start the branch", Lane: 0},
	})
}

func demoCodeLanguage(c *ui.Context, log func(string)) {
	code.Language(c, "Languages", []code.LangShare{{Name: "Go", Share: 0.72}, {Name: "Shell", Share: 0.18}, {Name: "JS", Share: 0.1}})
}

func demoKeyTable(c *ui.Context, log func(string)) {
	code.KeyTable(c, "Shortcuts", []code.KeyGroup{
		{Name: "Editor", Shortcuts: []code.Shortcut{
			{Keys: "Ctrl+S", Action: "Save the file"},
			{Keys: "Ctrl+Shift+P", Action: "Open the palette"},
		}},
		{Name: "View", Shortcuts: []code.Shortcut{
			{Keys: "Ctrl+B", Action: "Toggle the sidebar"},
			{Keys: "Ctrl+`", Action: "Toggle the terminal"},
		}},
	})
}

func demoCodeConsole(c *ui.Context, log func(string)) {
	entries := []code.LogEntry{
		{Level: code.LevelLog, Text: "gallery booted"},
		{Level: code.LevelLog, Text: "wasm ready"},
		{Level: code.LevelWarn, Text: "bundle over 30 MB"},
		{Level: code.LevelError, Text: "chunk 7 failed to load"},
	}
	root := ui.Column(c)
	filter := intState(root, "filter", 7)
	root.Children(func() {
		if code.Console(c, "Console", entries, filter, code.ConsoleOptions{}).Changed() {
			log("filter " + fmt.Sprint(*filter))
		}
	})
}

func demoREPL(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	history := replHistoryState(root, "history")
	draft := textState(root, "draft", "")
	root.Children(func() {
		r, _ := code.REPL(c, "REPL", history, draft, code.REPLOptions{})
		if r.Changed() {
			log(fmt.Sprint(len(*history)) + " lines")
		}
	})
}

func demoPackageCard(c *ui.Context, log func(string)) {
	_, copy := code.PackageCard(c, "github.com/egoist/mygo", "v0.5.5", 12400, code.PackageCardOptions{})
	if copy.Clicked() {
		log("copied install")
	}
}

func demoEnv(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	shown := boolState(root, "shown", false)
	vars := []code.EnvVar{
		{Key: "API_KEY", Value: "sk-live-redacted"},
		{Key: "REGION", Value: "us-east-1"},
		{Key: "MAX_RETRIES", Value: "3"},
	}
	root.Children(func() {
		if code.Env(c, "Environment", vars, shown, code.EnvOptions{}).Changed() {
			log("values " + map[bool]string{true: "shown", false: "hidden"}[*shown])
		}
	})
}

func demoPlayground(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	st := ui.Local(root, "editor", func() *editor.State {
		buf := editor.NewBuffer("println(\"hello\")\nprintln(2 + 3)\n")
		buf.SetLanguage(editor.LanguageFor("main.go"))
		return &editor.State{Buffer: buf}
	})
	output := textState(root, "output", "")
	root.Children(func() {
		r, _ := code.Playground(c, *st, output, func(src string) {
			*output = "hello\n5"
		}, code.PlaygroundOptions{})
		if r.Changed() {
			log("ran " + fmt.Sprint(utf8.RuneCountInString((*st).Buffer.String())) + " runes")
		}
	})
}

func replHistoryState(root ui.Element, key string) *[]code.REPLLine {
	return ui.Local(root, key, func() []code.REPLLine {
		return []code.REPLLine{{In: "1 + 2", Out: "3"}}
	})
}

func int64State(root ui.Element, key string, initial int64) *int64 {
	return ui.Local(root, key, func() int64 { return initial })
}
