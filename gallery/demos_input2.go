package gallery

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/richtext"
)

var inputEntries2 = []Entry{
	{219, "CountedField", "input", "Character counter warning near its limit.", demoCountedField},
	{220, "AdornedField", "input", "Prefix and suffix slots with focus tint.", demoAdornedField},
	{221, "ClearableField", "input", "Clear button appears with text.", demoClearableField},
	{222, "FloatingLabel", "input", "Label floats above on focus or value.", demoFloatingLabel},
	{223, "UnderlineField", "input", "Underline grows from center on focus.", demoUnderlineField},
	{224, "FilledField", "input", "Filled surface field with animated color.", demoFilledField},
	{225, "JSONField", "input", "JSON text area with validation.", demoJSONField},
	{226, "RegexField", "input", "Pattern with a match count.", demoRegexField},
	{227, "KeyValueField", "input", "Key and value rows with add and remove.", demoKeyValueField},
	{228, "ListField", "input", "String rows with add and remove.", demoListField},
	{229, "FilterSearch", "input", "Search with toggleable filter chips.", demoFilterSearch},
	{230, "SearchHistory", "input", "Recent searches below the input.", demoSearchHistory},
	{231, "DebouncedSearch", "input", "Reports a query after a pause.", demoDebouncedSearch},
	{232, "VerifyField", "input", "Send code counts down before resend.", demoVerifyField},
	{233, "SignaturePad", "input", "Pointer strokes with a clear action.", demoSignaturePad},
	{234, "RichEditor", "input", "Rich text editor with a styled toolbar.", demoRichEditor},
	{235, "MarkdownEditor", "input", "Edit and preview panes with cross-fades.", demoMarkdownEditor},
	{236, "NoteInput", "input", "Borderless large text area.", demoNoteInput},
	{237, "QuantityInput", "input", "Minus, rolling value and plus.", demoQuantityInput},
	{238, "IPField", "input", "Four octets with range validation.", demoIPField},
	{239, "PortField", "input", "Port from 1 through 65535.", demoPortField},
	{240, "MaskedField", "input", "Digits filling a mask.", demoMaskedField},
	{241, "TemplateField", "input", "Variable tokens inserted from a menu.", demoTemplateField},
	{242, "AutosizeInput", "input", "Single-line input following its text width.", demoAutosizeInput},
	{243, "InputGroup", "input", "Attached inputs and an action sharing borders.", demoInputGroup},
	{244, "FileNameField", "input", "Name editing with a fixed extension.", demoFileNameField},
	{245, "ConfirmPassword", "input", "Two masked fields with a match mark.", demoConfirmPassword},
}

func demoCountedField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.CountedField(c, "Bio", value, input.CountedFieldOptions{Max: 40}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoAdornedField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.AdornedField(c, "Amount", value, func() { ui.Text(c, "$") }, func() { ui.Text(c, "USD") }, input.AdornedFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoClearableField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.ClearableField(c, "Name", value, input.ClearableFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoFloatingLabel(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.FloatingLabel(c, "Email", value, input.FloatingLabelOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoUnderlineField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.UnderlineField(c, "Name", value, input.UnderlineFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoFilledField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.FilledField(c, "Name", value, input.FilledFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoJSONField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", `{"name":"Crychic"}`)
	changes := changeCount(root)
	root.Children(func() {
		if input.JSONField(c, value, input.JSONFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoRegexField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	pattern := textState(root, "pattern", "a+")
	changes := changeCount(root)
	root.Children(func() {
		if input.RegexField(c, pattern, "aa ba", input.RegexFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoKeyValueField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	rows := ui.Local(root, "rows", func() []input.KeyValueRow { return []input.KeyValueRow{{Key: "host", Value: "localhost"}} })
	root.Children(func() {
		if input.KeyValueField(c, rows, input.KeyValueFieldOptions{}).Changed() {
			log(fmt.Sprintf("rows: %d", len(*rows)))
		}
	})
}

func demoListField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	items := ui.Local(root, "items", func() []string { return []string{"First"} })
	root.Children(func() {
		if input.ListField(c, "Items", items, input.ListFieldOptions{}).Changed() {
			log(fmt.Sprintf("items: %d", len(*items)))
		}
	})
}

func demoFilterSearch(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	query := textState(root, "query", "")
	selected := ui.Local(root, "selected", func() []string { return nil })
	root.Children(func() {
		if input.FilterSearch(c, query, []string{"UI", "Docs"}, selected, input.FilterSearchOptions{}).Changed() {
			log(fmt.Sprintf("selected: %d", len(*selected)))
		}
	})
}

func demoSearchHistory(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	query := textState(root, "query", "")
	history := ui.Local(root, "history", func() []string { return nil })
	root.Children(func() {
		if input.SearchHistory(c, query, history, input.SearchHistoryOptions{}).Changed() {
			log(fmt.Sprintf("history: %d", len(*history)))
		}
	})
}

func demoDebouncedSearch(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	query := textState(root, "query", "")
	reported := textState(root, "reported", "")
	root.Children(func() {
		if input.DebouncedSearch(c, "Search", query, reported, input.DebouncedSearchOptions{}).Changed() {
			log("reported")
		}
	})
}

func demoVerifyField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	code := textState(root, "code", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.VerifyField(c, "Code", code, input.VerifyFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoSignaturePad(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	changes := changeCount(root)
	root.Children(func() {
		if input.SignaturePad(c, input.SignaturePadOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("strokes: %d", *changes))
		}
	})
}

func demoRichEditor(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	doc := ui.Local(root, "doc", func() *richtext.Document { return richtext.NewDocument(richtext.NewText("Hello")) })
	changes := changeCount(root)
	root.Children(func() {
		if input.RichEditor(c, *doc).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoMarkdownEditor(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	source := textState(root, "source", "# Hello")
	changes := changeCount(root)
	root.Children(func() {
		if input.MarkdownEditor(c, source, input.MarkdownEditorOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoNoteInput(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.NoteInput(c, value, "Write a note").Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoQuantityInput(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := ui.Local(root, "value", func() int { return 1 })
	root.Children(func() {
		if input.QuantityInput(c, "Items", value, input.QuantityInputOptions{}).Changed() {
			log(fmt.Sprintf("value: %d", *value))
		}
	})
}

func demoIPField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.IPField(c, value, input.IPFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoPortField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := ui.Local(root, "value", func() int { return 8080 })
	root.Children(func() {
		if input.PortField(c, "Port", value, input.PortFieldOptions{}).Changed() {
			log(fmt.Sprintf("value: %d", *value))
		}
	})
}

func demoMaskedField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.MaskedField(c, "Serial", value, "###-##", input.MaskedFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoTemplateField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "Hello ")
	changes := changeCount(root)
	root.Children(func() {
		if input.TemplateField(c, "Template", value, []string{"name", "date"}, input.TemplateFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoAutosizeInput(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.AutosizeInput(c, value, input.AutosizeInputOptions{Placeholder: "Name"}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoInputGroup(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	host := textState(root, "host", "")
	port := textState(root, "port", "")
	changes := changeCount(root)
	root.Children(func() {
		view, action := input.InputGroup(c, []string{"Host", "Port"}, []*string{host, port}, "Connect")
		if view.Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
		if action.Clicked() {
			log("connect")
		}
	})
}

func demoFileNameField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	name := textState(root, "name", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.FileNameField(c, "File", name, "pdf", input.FileNameFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoConfirmPassword(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	password := textState(root, "password", "")
	confirm := textState(root, "confirm", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.ConfirmPassword(c, password, confirm, input.ConfirmPasswordOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}
