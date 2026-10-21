package gallery

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

var typographyEntries2 = []Entry{
	{71, "ReadMore", "typography", "Clamped tail fading before a one-way expansion.", demoReadMore},
	{72, "Number", "typography", "Localized number with rolling digits.", demoNumber},
	{73, "Money", "typography", "Currency code and two decimals with rolling digits.", demoMoney},
	{74, "Percent", "typography", "Percent with sign coloring and rolling digits.", demoPercent},
	{75, "RelativeTime", "typography", "Relative time updating with the clock.", demoRelativeTime},
	{76, "Timestamp", "typography", "Localized date and time with a full tooltip.", demoTimestamp},
	{77, "Bytes", "typography", "Human byte size with rolling digits.", demoBytes},
	{78, "Duration", "typography", "Human duration with rolling digits.", demoDuration},
	{79, "BulletList", "typography", "Bulleted items staggering in.", demoBulletList},
	{80, "NumberedList", "typography", "Numbered items staggering in.", demoNumberedList},
	{81, "DefinitionList", "typography", "Terms and descriptions fading in.", demoDefinitionList},
	{82, "Prose", "typography", "Long-form paragraphs fading in.", demoProse},
	{83, "TextLink", "typography", "Link with an underline growing on hover.", demoTextLink},
	{84, "ExternalLink", "typography", "External link with a nudging arrow.", demoExternalLink},
	{85, "AnchorHeading", "typography", "Heading with a hover anchor to copy.", demoAnchorHeading},
	{86, "Footnote", "typography", "Superscript number opening a note.", demoFootnote},
	{87, "LoadingDots", "typography", "Loading label with cycling dots.", demoLoadingDots},
	{88, "DiffText", "typography", "Deleted and inserted text with changing fades.", demoDiffText},
	{89, "Monogram", "typography", "Name initials in a deterministic color circle.", demoMonogram},
	{90, "Ruby", "typography", "Annotation fading in over base text.", demoRuby},
}

func demoReadMore(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	open := ui.Local(root, "open", func() bool { return false })
	text := strings.Repeat("Long-form details wrap across lines and stay readable. ", 12)
	root.Children(func() {
		if typography.ReadMore(c, text, 3, open, typography.ReadMoreOptions{}).Changed() {
			log(fmt.Sprintf("open: %v", *open))
		}
	})
}

func demoNumber(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "number", 1234, log)
		typography.Number(c, float64(*n)+0.5)
	})
}

func demoMoney(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "amount", 1234, log)
		typography.Money(c, float64(*n)+0.5, "USD")
	})
}

func demoPercent(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "percent", 25, log)
		typography.Percent(c, float64(*n-50)/100)
	})
}

func demoRelativeTime(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	when := ui.Local(root, "when", func() time.Time { return c.Now().Add(-3 * time.Minute) })
	root.Children(func() { typography.RelativeTime(c, *when, typography.RelativeTimeOptions{}) })
}

func demoTimestamp(c *ui.Context, log func(string)) {
	when := time.Date(2026, time.January, 5, 13, 7, 0, 0, time.UTC)
	typography.Timestamp(c, when)
}

func demoBytes(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "bytes", 1500, log)
		typography.Bytes(c, int64(*n))
	})
}

func demoDuration(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "seconds", 64, log)
		typography.Duration(c, time.Duration(*n)*time.Second)
	})
}

func demoBulletList(c *ui.Context, log func(string)) {
	typography.BulletList(c, []string{"Read the requirement", "Build the component", "Run the tests"})
}

func demoNumberedList(c *ui.Context, log func(string)) {
	typography.NumberedList(c, []string{"Read the requirement", "Build the component", "Run the tests"})
}

func demoDefinitionList(c *ui.Context, log func(string)) {
	typography.DefinitionList(c, []typography.Definition{
		{Term: "Native", Description: "Drawn with the platform's own text and input."},
		{Term: "Motion", Description: "Short movement with a clear state change."},
	})
}

func demoProse(c *ui.Context, log func(string)) {
	typography.Prose(c, []string{
		"Components use the theme for color, spacing and type.",
		"Motion explains a change without getting in the way.",
	})
}

func demoTextLink(c *ui.Context, log func(string)) {
	if typography.TextLink(c, "MyGo documentation", "https://mygo.egoist.dev/docs/ui").Clicked() {
		log("open documentation")
	}
}

func demoExternalLink(c *ui.Context, log func(string)) {
	if typography.ExternalLink(c, "CrychicUI repository", "https://github.com/ELY-Open-Source/CrychicUI").Clicked() {
		log("open repository")
	}
}

func demoAnchorHeading(c *ui.Context, log func(string)) {
	_, copyAnchor := typography.AnchorHeading(c, "Getting started", "getting-started")
	if copyAnchor.Clicked() {
		log("anchor copied")
	}
}

func demoFootnote(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	open := ui.Local(root, "open", func() bool { return false })
	root.Children(func() {
		if typography.Footnote(c, 1, "This note adds the supporting detail.", open).Changed() {
			log(fmt.Sprintf("open: %v", *open))
		}
	})
}

func demoLoadingDots(c *ui.Context, log func(string)) {
	typography.LoadingDots(c, "Loading")
}

func demoDiffText(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "revision", 1, log)
		typography.DiffText(c, "old text", fmt.Sprintf("new text %d", *n))
	})
}

func demoMonogram(c *ui.Context, log func(string)) {
	typography.Monogram(c, "Crychic UI")
}

func demoRuby(c *ui.Context, log func(string)) {
	typography.Ruby(c, "東京", "とうきょう")
}
