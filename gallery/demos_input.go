package gallery

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/egoist/mygo/ui"
)

var inputEntries = []Entry{
	{191, "TextField", "input", "Label, input, help and a shaking error.", demoTextField},
	{192, "AutoTextArea", "input", "Text area growing with its lines.", demoAutoTextArea},
	{193, "PasswordField", "input", "Masked input with reveal and strength.", demoPasswordField},
	{194, "SearchInput", "input", "Search icon, clear button and Escape.", demoSearchInput},
	{195, "EmailField", "input", "Email validation when focus leaves.", demoEmailField},
	{196, "URLField", "input", "Fixed scheme prefix and host validation.", demoURLField},
	{197, "PhoneField", "input", "Grouped digits with a country prefix.", demoPhoneField},
	{198, "NumberField", "input", "Number arrows, keys and rolling digits.", demoNumberField},
	{199, "CurrencyField", "input", "Currency prefix and cents on blur.", demoCurrencyField},
	{200, "PercentField", "input", "Percent suffix clamped from zero to 100.", demoPercentField},
	{201, "UnitField", "input", "Number with a cross-fading unit select.", demoUnitField},
	{202, "OTPInput", "input", "Digit boxes with paste and completion.", demoOTPInput},
	{203, "PinInput", "input", "Masked digit boxes with popping dots.", demoPinInput},
	{204, "TagInput", "input", "Enter or comma adds a chip; Backspace removes.", demoTagInput},
	{205, "MentionInput", "input", "@ suggestions insert a mention.", demoMentionInput},
	{206, "CommandInput", "input", "/ commands insert a choice.", demoCommandInput},
	{207, "InlineEdit", "input", "Text edits inline; Enter saves, Escape cancels.", demoInlineEdit},
	{208, "EditableTitle", "input", "Heading-sized inline edit.", demoEditableTitle},
	{209, "Suggest", "input", "Styled suggestions filtering as you type.", demoSuggest},
	{210, "ComboField", "input", "Styled combobox with a scaling popup.", demoComboField},
	{211, "MultiCombo", "input", "Combobox adding removable chips.", demoMultiCombo},
	{212, "AddressField", "input", "Street, city, region and postal rows.", demoAddressField},
	{213, "NameField", "input", "Given and family name in one row.", demoNameField},
	{214, "CardNumberField", "input", "Card digits grouped with a brand icon.", demoCardNumberField},
	{215, "ExpiryField", "input", "MM/YY formatting and validation.", demoExpiryField},
	{216, "CVCField", "input", "Three or four masked digits.", demoCVCField},
	{217, "IBANField", "input", "IBAN grouped by four with checksum validation.", demoIBANField},
	{218, "SlugField", "input", "Slug derived from a title until locked or edited.", demoSlugField},
}

func textState(root ui.Element, key, start string) *string {
	return ui.Local(root, key, func() string { return start })
}

func numberState(root ui.Element, key string, start float64) *float64 {
	return ui.Local(root, key, func() float64 { return start })
}

func changeCount(root ui.Element) *int {
	return ui.Local(root, "changes", func() int { return 0 })
}

func demoTextField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.TextField(c, "Email", value, input.TextFieldOptions{Placeholder: "name@example.com", Help: "We never share it."}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoAutoTextArea(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.AutoTextArea(c, value, input.AutoTextAreaOptions{Placeholder: "Write a note"}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoPasswordField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.PasswordField(c, "Password", value, input.PasswordFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoSearchInput(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.SearchInput(c, value, input.SearchInputOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoEmailField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.EmailField(c, "Email", value, input.EmailFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoURLField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.URLField(c, "Site", value, input.URLFieldOptions{Placeholder: "example.com"}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoPhoneField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.PhoneField(c, "Phone", value, input.PhoneFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoNumberField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := numberState(root, "value", 4)
	root.Children(func() {
		if input.NumberField(c, "Copies", value, input.NumberFieldOptions{}).Changed() {
			log(fmt.Sprintf("copies: %.0f", *value))
		}
	})
}

func demoCurrencyField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := numberState(root, "value", 12.34)
	changes := changeCount(root)
	root.Children(func() {
		if input.CurrencyField(c, "Price", value, input.CurrencyFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoPercentField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := numberState(root, "value", 42)
	changes := changeCount(root)
	root.Children(func() {
		if input.PercentField(c, "Progress", value).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoUnitField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := numberState(root, "value", 2)
	unit := textState(root, "unit", "kg")
	root.Children(func() {
		if input.UnitField(c, "Weight", value, unit, []string{"g", "kg"}, input.UnitFieldOptions{}).Changed() {
			log(fmt.Sprintf("value: %.0f", *value))
		}
	})
}

func demoOTPInput(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.OTPInput(c, value, 4, input.OTPInputOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("digits: %d", len(*value)))
		}
	})
}

func demoPinInput(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.PinInput(c, value, 4, input.PinInputOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("digits: %d", len(*value)))
		}
	})
}

func demoTagInput(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	tags := ui.Local(root, "tags", func() []string { return []string{"ui"} })
	draft := textState(root, "draft", "")
	root.Children(func() {
		if input.TagInput(c, tags, draft, input.TagInputOptions{}).Changed() {
			log(fmt.Sprintf("tags: %d", len(*tags)))
		}
	})
}

func demoMentionInput(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.MentionInput(c, "Comment", value, []string{"Ada", "Grace"}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoCommandInput(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.CommandInput(c, "Command", value, []string{"delete", "duplicate"}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoInlineEdit(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "Draft title")
	root.Children(func() {
		if input.InlineEdit(c, value, input.InlineEditOptions{}).Changed() {
			log("saved")
		}
	})
}

func demoEditableTitle(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "Quarterly plan")
	root.Children(func() {
		if input.EditableTitle(c, value, input.EditableTitleOptions{}).Changed() {
			log("saved")
		}
	})
}

func demoSuggest(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	root.Children(func() {
		if input.Suggest(c, "Fruit", value, []string{"Apple", "Pear"}).Changed() {
			log("changed")
		}
	})
}

func demoComboField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	selected := textState(root, "selected", "One")
	root.Children(func() {
		if input.ComboField(c, "Choice", selected, []string{"One", "Two"}).Changed() {
			log("selected")
		}
	})
}

func demoMultiCombo(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	chosen := ui.Local(root, "chosen", func() []string { return nil })
	root.Children(func() {
		if input.MultiCombo(c, "Choose", chosen, []string{"One", "Two"}).Changed() {
			log(fmt.Sprintf("items: %d", len(*chosen)))
		}
	})
}

func demoAddressField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := ui.Local(root, "value", func() input.Address { return input.Address{} })
	changes := changeCount(root)
	root.Children(func() {
		if input.AddressField(c, value, input.AddressFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoNameField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	given := textState(root, "given", "")
	family := textState(root, "family", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.NameField(c, given, family, input.NameFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoCardNumberField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.CardNumberField(c, value, input.CardNumberFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoExpiryField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.ExpiryField(c, value, input.ExpiryFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoCVCField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.CVCField(c, value, input.CVCFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoIBANField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	value := textState(root, "value", "")
	changes := changeCount(root)
	root.Children(func() {
		if input.IBANField(c, value, input.IBANFieldOptions{}).Changed() {
			(*changes)++
			log(fmt.Sprintf("changes: %d", *changes))
		}
	})
}

func demoSlugField(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	title := "My First Post"
	slug := textState(root, "slug", "")
	locked := ui.Local(root, "locked", func() bool { return false })
	root.Children(func() {
		if input.SlugField(c, title, slug, locked, input.SlugFieldOptions{}).Changed() {
			log(fmt.Sprintf("locked: %v", *locked))
		}
	})
}
