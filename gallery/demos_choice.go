package gallery

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/egoist/mygo/ui"
)

var choiceEntries = []Entry{
	{246, "Check", "choice", "Check box whose mark strokes in.", demoCheck},
	{247, "CheckCard", "choice", "Selectable card; border and check animate.", demoCheckCard},
	{248, "CheckGroup", "choice", "Checks under a parent that goes mixed.", demoCheckGroup},
	{249, "Toggle", "choice", "Switch whose knob springs across.", demoToggle},
	{250, "ToggleRow", "choice", "Setting row: title, description and switch.", demoToggleRow},
	{251, "Option", "choice", "Radio whose dot springs in.", demoOption},
	{252, "OptionCard", "choice", "Radio card with an animated border.", demoOptionCard},
	{253, "OptionGroup", "choice", "Radios in a row or column, one Tab stop.", demoOptionGroup},
	{254, "Segments", "choice", "Segmented control whose thumb slides.", demoSegments},
	{255, "IconSegments", "choice", "Icon-only segments with tooltips.", demoIconSegments},
	{256, "ToggleSet", "choice", "Joined toggles choosing several.", demoToggleSet},
	{257, "ChipChoice", "choice", "One choice among pill chips.", demoChipChoice},
	{258, "FilterChips", "choice", "Chips choosing several, with a rolling count.", demoFilterChips},
	{259, "Dropdown", "choice", "Select whose popup scales in; chosen row checked.", demoDropdown},
	{260, "MultiDropdown", "choice", "Select of checks with a count in its trigger.", demoMultiDropdown},
	{261, "GroupedDropdown", "choice", "Select with section titles.", demoGroupedDropdown},
	{262, "SearchDropdown", "choice", "Select whose popup filters as you type.", demoSearchDropdown},
	{263, "TreeDropdown", "choice", "Select of nested items that expand.", demoTreeDropdown},
	{264, "Cascade", "choice", "Columns that drill down; columns slide in.", demoCascade},
	{265, "Listbox", "choice", "Inline list choosing one; highlight slides.", demoListbox},
	{266, "Transfer", "choice", "Two lists; clicked items fly across.", demoTransfer},
	{267, "ValueSlider", "choice", "Slider with a value bubble on the knob.", demoValueSlider},
	{268, "RangeSlider", "choice", "Two knobs, values rolling above.", demoRangeSlider},
	{269, "TickSlider", "choice", "Stepped slider with ticks and labels.", demoTickSlider},
	{270, "VerticalSlider", "choice", "Slider rising from the bottom.", demoVerticalSlider},
}

func boolState(root ui.Element, key string, start bool) *bool {
	return ui.Local(root, key, func() bool { return start })
}

func listState(root ui.Element, key string, start ...string) *[]string {
	return ui.Local(root, key, func() []string { return start })
}

func demoCheck(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(8)
	a, b := boolState(root, "a", true), boolState(root, "b", false)
	root.Children(func() {
		if choice.Check(c, "Sync across devices", a).Changed() {
			log(fmt.Sprintf("sync: %v", *a))
		}
		if choice.Check(c, "Show hidden files", b).Changed() {
			log(fmt.Sprintf("hidden: %v", *b))
		}
	})
}

func demoCheckCard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(8).MaxWidth(360)
	on := boolState(root, "on", false)
	root.Children(func() {
		if choice.CheckCard(c, "Nightly backups", "Copies your files at 2 AM.", on).Changed() {
			log(fmt.Sprintf("backups: %v", *on))
		}
	})
}

func demoCheckGroup(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	sel := listState(root, "sel", "Mail")
	root.Children(func() {
		if choice.CheckGroup(c, []string{"Mail", "Calendar", "Messages"}, sel, choice.CheckGroupOptions{}).Changed() {
			log(fmt.Sprintf("selected: %d", len(*sel)))
		}
	})
}

func demoToggle(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	on := boolState(root, "on", true)
	root.Children(func() {
		if choice.Toggle(c, "Wi-Fi", on).Changed() {
			log(fmt.Sprintf("wi-fi: %v", *on))
		}
	})
}

func demoToggleRow(c *ui.Context, log func(string)) {
	root := ui.Column(c).MaxWidth(420).Dividers(1, c.Theme().Border)
	cues, dark := boolState(root, "sound", true), boolState(root, "dark", false)
	root.Children(func() {
		if choice.ToggleRow(c, "Sound cues", "Play a cue on each action.", cues).Changed() {
			log(fmt.Sprintf("sound: %v", *cues))
		}
		if choice.ToggleRow(c, "Dim at night", "Lower the brightness after sunset.", dark).Changed() {
			log(fmt.Sprintf("dim: %v", *dark))
		}
	})
}

func demoOption(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(8)
	size := textState(root, "size", "Small")
	root.Children(func() {
		for _, s := range []string{"Small", "Large"} {
			c.Key(s)
			if choice.Option(c, s, s, size).Changed() {
				log("size: " + *size)
			}
		}
	})
}

func demoOptionCard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(8).MaxWidth(360)
	plan := textState(root, "plan", "Pro")
	root.Children(func() {
		for _, p := range [][2]string{{"Pro", "For one person."}, {"Team", "For five or more."}} {
			c.Key(p[0])
			if choice.OptionCard(c, p[0], p[1], p[0], plan).Changed() {
				log("plan: " + *plan)
			}
		}
	})
}

func demoOptionGroup(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	align := textState(root, "align", "Left")
	root.Children(func() {
		if choice.OptionGroup(c, "Alignment", []string{"Left", "Center", "Right"}, align, choice.OptionGroupOptions{Row: true}).Changed() {
			log("align: " + *align)
		}
	})
}

func demoSegments(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	view := ui.Local(root, "view", func() int { return 0 })
	root.Children(func() {
		if choice.Segments(c, "Range", []string{"Day", "Week", "Month"}, view).Changed() {
			log(fmt.Sprintf("view: %d", *view))
		}
	})
}

func demoIconSegments(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	view := ui.Local(root, "view", func() int { return 0 })
	items := []choice.SegmentItem{{Label: "List", Icon: "list"}, {Label: "Grid", Icon: "layout-grid"}, {Label: "Columns", Icon: "columns-3"}}
	root.Children(func() {
		if choice.IconSegments(c, "Layout", items, view).Changed() {
			log("layout: " + items[*view].Label)
		}
	})
}

func demoToggleSet(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	sel := listState(root, "sel", "Bold")
	root.Children(func() {
		if choice.ToggleSet(c, "Style", []string{"Bold", "Italic", "Underline"}, sel).Changed() {
			log("style: " + strings.Join(*sel, ", "))
		}
	})
}

func demoChipChoice(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	size := textState(root, "size", "M")
	root.Children(func() {
		if choice.ChipChoice(c, "Size", []string{"XS", "S", "M", "L", "XL"}, size).Changed() {
			log("size: " + *size)
		}
	})
}

func demoFilterChips(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	sel := listState(root, "sel")
	root.Children(func() {
		if choice.FilterChips(c, "Filters", []string{"Open", "Assigned to me", "Has PR", "Stale"}, sel).Changed() {
			log(fmt.Sprintf("filters: %d", len(*sel)))
		}
	})
}

func demoDropdown(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	plan := textState(root, "plan", "Free")
	root.Children(func() {
		if choice.Dropdown(c, "Plan", []string{"Free", "Pro", "Team", "Enterprise"}, plan).Changed() {
			log("plan: " + *plan)
		}
	})
}

func demoMultiDropdown(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	sel := listState(root, "sel", "Bug")
	root.Children(func() {
		if choice.MultiDropdown(c, "Labels", []string{"Bug", "Docs", "Feature", "UI"}, sel).Changed() {
			log(fmt.Sprintf("labels: %d", len(*sel)))
		}
	})
}

func demoGroupedDropdown(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	city := textState(root, "city", "Paris")
	groups := []choice.DropdownGroup{{Title: "Europe", Options: []string{"Paris", "Rome", "Berlin"}}, {Title: "Asia", Options: []string{"Tokyo", "Seoul"}}}
	root.Children(func() {
		if choice.GroupedDropdown(c, "City", groups, city).Changed() {
			log("city: " + *city)
		}
	})
}

func demoSearchDropdown(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	lang := textState(root, "lang", "Go")
	root.Children(func() {
		options := []string{"Go", "Rust", "Ruby", "Python", "TypeScript", "Swift", "Kotlin"}
		if choice.SearchDropdown(c, "Language", options, lang, choice.SearchDropdownOptions{}).Changed() {
			log("language: " + *lang)
		}
	})
}

func demoTreeDropdown(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	item := textState(root, "item", "Guides")
	tree := []choice.TreeItem{
		{Label: "Docs", Children: []choice.TreeItem{{Label: "Guides"}, {Label: "API", Children: []choice.TreeItem{{Label: "REST"}, {Label: "GraphQL"}}}}},
		{Label: "Blog"},
	}
	root.Children(func() {
		if choice.TreeDropdown(c, "Section", tree, item, choice.TreeDropdownOptions{}).Changed() {
			log("section: " + *item)
		}
	})
}

func demoCascade(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	path := listState(root, "path")
	tree := []choice.TreeItem{
		{Label: "France", Children: []choice.TreeItem{{Label: "Paris"}, {Label: "Lyon"}}},
		{Label: "Japan", Children: []choice.TreeItem{{Label: "Tokyo"}, {Label: "Osaka"}}},
	}
	root.Children(func() {
		if choice.Cascade(c, "Place", tree, path).Changed() {
			log("place: " + strings.Join(*path, " / "))
		}
	})
}

func demoListbox(c *ui.Context, log func(string)) {
	root := ui.Box(c).Width(220)
	fruit := textState(root, "fruit", "Apple")
	root.Children(func() {
		if choice.Listbox(c, "Fruit", []string{"Apple", "Pear", "Plum", "Cherry"}, fruit).Changed() {
			log("fruit: " + *fruit)
		}
	})
}

func demoTransfer(c *ui.Context, log func(string)) {
	root := ui.Box(c).MaxWidth(520)
	left, right := listState(root, "left", "Alpha", "Beta", "Delta"), listState(root, "right", "Gamma")
	root.Children(func() {
		if choice.Transfer(c, left, right, choice.TransferOptions{}).Changed() {
			log(fmt.Sprintf("chosen: %d", len(*right)))
		}
	})
}

func demoValueSlider(c *ui.Context, log func(string)) {
	root := ui.Box(c).Width(300).PaddingY(24)
	v := numberState(root, "v", 40)
	root.Children(func() {
		if choice.ValueSlider(c, "Volume", v, 0, 100, choice.ValueSliderOptions{}).Changed() {
			log(fmt.Sprintf("volume: %.0f", *v))
		}
	})
}

func demoRangeSlider(c *ui.Context, log func(string)) {
	root := ui.Box(c).Width(300)
	low, high := numberState(root, "low", 200), numberState(root, "high", 800)
	root.Children(func() {
		if choice.RangeSlider(c, "Price", low, high, 0, 1000, 50, choice.RangeSliderOptions{Format: "$%.0f"}).Changed() {
			log(fmt.Sprintf("price: %.0f–%.0f", *low, *high))
		}
	})
}

func demoTickSlider(c *ui.Context, log func(string)) {
	root := ui.Box(c).Width(320)
	v := numberState(root, "v", 2)
	root.Children(func() {
		o := choice.TickSliderOptions{Labels: []string{"XS", "S", "M", "L", "XL"}}
		if choice.TickSlider(c, "Size", v, 0, 4, 1, o).Changed() {
			log("size: " + o.Labels[int(*v)])
		}
	})
}

func demoVerticalSlider(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	v := numberState(root, "v", 30)
	root.Children(func() {
		if choice.VerticalSlider(c, "Level", v, 0, 100, choice.VerticalSliderOptions{}).Changed() {
			log(fmt.Sprintf("level: %.0f", *v))
		}
	})
}
