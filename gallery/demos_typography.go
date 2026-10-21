package gallery

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

var typographyEntries = []Entry{
	{51, "Display", "typography", "Hero text at 3 rem, fading up.", demoDisplay},
	{52, "Heading", "typography", "Six leveled heading sizes and roles.", demoHeading},
	{53, "Title", "typography", "Bold title with cross-fading changes.", demoTitle},
	{54, "Subtitle", "typography", "Muted subtitle with cross-fading changes.", demoSubtitle},
	{55, "Lead", "typography", "Large introductory text fading in.", demoLead},
	{56, "Paragraph", "typography", "Wrapping body text with cross-fading changes.", demoParagraph},
	{57, "Caption", "typography", "Small muted text with cross-fading changes.", demoCaption},
	{58, "Overline", "typography", "Uppercase spaced label with cross-fading changes.", demoOverline},
	{59, "Eyebrow", "typography", "Accent label sliding above a heading.", demoEyebrow},
	{60, "Label", "typography", "Medium label with an optional required mark.", demoLabel},
	{61, "InlineCode", "typography", "Monospace text on a surface pill.", demoInlineCode},
	{62, "Kbd", "typography", "Key cap pressing when its key matches.", demoKbd},
	{63, "KeyCombo", "typography", "Key caps pressing in sequence on a match.", demoKeyCombo},
	{64, "Blockquote", "typography", "Quote beside a growing accent bar.", demoBlockquote},
	{65, "PullQuote", "typography", "Large centered quote fading up.", demoPullQuote},
	{66, "Highlight", "typography", "Text with a marker sweeping behind it.", demoHighlight},
	{67, "Mark", "typography", "Search-match highlight fading in.", demoMark},
	{68, "Abbr", "typography", "Dotted abbreviation with a full-term tooltip.", demoAbbr},
	{69, "Truncate", "typography", "One ellipsized line with a full-text tooltip.", demoTruncate},
	{70, "Clamp", "typography", "Line clamp with a springy More and Less toggle.", demoClamp},
}

func demoDisplay(c *ui.Context, log func(string)) {
	typography.Display(c, "Build native interfaces")
}

func demoHeading(c *ui.Context, log func(string)) {
	ui.Column(c).Gap(12).Children(func() {
		typography.Heading(c, "Second-level heading", 2)
		typography.Heading(c, "Fourth-level heading", 4)
	})
}

func demoTitle(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "revisions", 1, log)
		typography.Title(c, fmt.Sprintf("Release %d", *n))
	})
}

func demoSubtitle(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "subtitles", 1, log)
		typography.Subtitle(c, fmt.Sprintf("Subtitle revision %d", *n))
	})
}

func demoLead(c *ui.Context, log func(string)) {
	typography.Lead(c, "A clear introduction to the component library.")
}

func demoParagraph(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "paragraphs", 1, log)
		typography.Paragraph(c, fmt.Sprintf("Paragraph revision %d wraps when its room runs out.", *n))
	})
}

func demoCaption(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "captions", 1, log)
		typography.Caption(c, fmt.Sprintf("Caption revision %d", *n))
	})
}

func demoOverline(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "overlines", 1, log)
		typography.Overline(c, fmt.Sprintf("Collection %d", *n))
	})
}

func demoEyebrow(c *ui.Context, log func(string)) {
	ui.Column(c).Gap(8).Children(func() {
		typography.Eyebrow(c, "COMPONENTS")
		typography.Heading(c, "Typography", 2)
	})
}

func demoLabel(c *ui.Context, log func(string)) {
	ui.Column(c).Gap(12).Children(func() {
		typography.Label(c, "Display name", typography.LabelOptions{})
		typography.Label(c, "Email address", typography.LabelOptions{Required: true})
	})
}

func demoInlineCode(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "commands", 1, log)
		typography.InlineCode(c, fmt.Sprintf("go test ./... # %d", *n))
	})
}

func demoKbd(c *ui.Context, log func(string)) {
	matched := c.Shortcut(0, ui.KeyS)
	if matched {
		log("S pressed")
	}
	ui.Row(c).Gap(12).Children(func() {
		ui.Text(c, "Press S")
		typography.Kbd(c, "S", matched)
	})
}

func demoKeyCombo(c *ui.Context, log func(string)) {
	matched := c.Shortcut(ui.Cmd, ui.KeyK)
	if matched {
		log("combo matched")
	}
	ui.Row(c).Gap(12).Children(func() {
		ui.Text(c, "Press the modifier and K")
		typography.KeyCombo(c, []string{"Ctrl / ⌘", "K"}, matched)
	})
}

func demoBlockquote(c *ui.Context, log func(string)) {
	typography.Blockquote(c, "The best interface is the one that stays out of the way.")
}

func demoPullQuote(c *ui.Context, log func(string)) {
	typography.PullQuote(c, "Make it simple. Make it work.", "The Crychic team")
}

func demoHighlight(c *ui.Context, log func(string)) {
	typography.Highlight(c, "Remember this sentence")
}

func demoMark(c *ui.Context, log func(string)) {
	typography.Mark(c, "search match")
}

func demoAbbr(c *ui.Context, log func(string)) {
	typography.Abbr(c, "API", "Application Programming Interface")
}

func demoTruncate(c *ui.Context, log func(string)) {
	ui.Column(c).Width(220).Children(func() {
		typography.Truncate(c, "A very long component title that cannot fit on one line")
	})
}

func demoClamp(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	open := ui.Local(root, "open", func() bool { return false })
	text := strings.Repeat("Long-form details wrap across lines and stay readable. ", 12)
	root.Children(func() {
		if typography.Clamp(c, text, 3, open, typography.ClampOptions{}).Changed() {
			log(fmt.Sprintf("open: %v", *open))
		}
	})
}
