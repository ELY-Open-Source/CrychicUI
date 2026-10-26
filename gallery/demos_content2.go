package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/content"
	"github.com/egoist/mygo/ui"
)

var contentEntries2 = []Entry{
	{ID: 1444, Name: "EditLink", Category: "content", Summary: "Edit this page link; cue Tap.", Demo: demoContentEditLink},
	{ID: 1445, Name: "WasHelpful", Category: "content", Summary: "Yes and no feedback; cue Pop.", Demo: demoContentWasHelpful},
	{ID: 1446, Name: "DocCrumbs", Category: "content", Summary: "Doc path.", Demo: demoContentDocCrumbs},
	{ID: 1447, Name: "VersionSelect", Category: "content", Summary: "Docs version select; cue Select.", Demo: demoContentVersionSelect},
	{ID: 1448, Name: "LinkCard", Category: "content", Summary: "Link card with description; cue Tap.", Demo: demoContentLinkCard},
	{ID: 1449, Name: "CaptionedImage", Category: "content", Summary: "Image with caption.", Demo: demoContentCaptionedImage},
	{ID: 1450, Name: "Embed", Category: "content", Summary: "Embedded box with title.", Demo: demoContentEmbed},
	{ID: 1451, Name: "ReadingTime", Category: "content", Summary: "Minutes to read.", Demo: demoContentReadingTime},
	{ID: 1452, Name: "StickyNote", Category: "content", Summary: "Note with tilt.", Demo: demoContentStickyNote},
	{ID: 1453, Name: "Quote", Category: "content", Summary: "Quote with source.", Demo: demoContentQuote},
	{ID: 1454, Name: "StoryTimeline", Category: "content", Summary: "Dated story beats.", Demo: demoContentStoryTimeline},
	{ID: 1455, Name: "PrintBar", Category: "content", Summary: "Print and share actions; cue Tap.", Demo: demoContentPrintBar},
}

func demoContentEditLink(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		_, edit := content.EditLink(c, "https://github.com/glaze-studio/docs/edit/main/kilns/booking.md", docDay(20), content.EditLinkOptions{})
		if edit.Clicked() {
			log("edit opened")
		}
	})
}

func demoContentWasHelpful(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	vote := ui.Local(root, "vote", func() content.Helpful { return content.Unrated })
	root.Children(func() {
		if content.WasHelpful(c, vote, content.WasHelpfulOptions{}).Changed() {
			log(fmt.Sprintf("vote %d", *vote))
		}
	})
}

func demoContentDocCrumbs(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		content.DocCrumbs(c, []string{"Docs", "API reference", "Webhooks", "Retries"})
	})
}

func demoContentVersionSelect(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Height(300)
	current := textState(root, "current", "v3.2")
	root.Children(func() {
		before := *current
		r := content.VersionSelect(c, []content.DocVersion{
			{Name: "v3.2", Note: "Current release", Latest: true},
			{Name: "v3.1", Note: "Supported until March"},
			{Name: "v3.0"},
			{Name: "v2.9", Note: "Legacy API"},
		}, current, content.VersionSelectOptions{})
		if r.Changed() {
			if *current == before {
				log("menu open")
			} else {
				log("version " + *current)
			}
		}
	})
}

func demoContentLinkCard(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Gap(10).Children(func() {
		if content.LinkCard(c, "Glaze API reference", "Every endpoint, with examples in cURL, Go and JavaScript.",
			"https://docs.glaze.studio/api", content.LinkCardOptions{Icon: "book-open"}).Clicked() {
			log("opened api")
		}
		if content.LinkCard(c, "Community forum", "Ask other studios how they run firings and classes.",
			"https://forum.glaze.studio", content.LinkCardOptions{Icon: "message-square"}).Clicked() {
			log("opened forum")
		}
	})
}

func demoContentCaptionedImage(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() {
		content.CaptionedImage(c, productShots()[1], "Speckled stoneware mug after a cone 6 oxidation firing.",
			content.CaptionedImageOptions{Credit: "Photo: Daniel Okafor"})
	})
}

func demoContentEmbed(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		content.Embed(c, "Kiln 2 · last firing", content.EmbedOptions{Icon: "flame", Host: "glaze.studio"}, func() {
			t := c.Theme()
			ui.Column(c).FillWidth().Gap(t.Space(2)).Children(func() {
				for _, row := range [][2]string{{"Ramp", "108°F per hour"}, {"Peak", "2232°F, cone 6"}, {"Hold", "10 minutes"}, {"Cooled", "Fri 6:40 am"}} {
					ui.Row(c).FillWidth().Children(func() {
						ui.Text(c, row[0]).FontSize(t.Rem(0.82)).TextColor(t.TextMuted).Grow(1)
						ui.Text(c, row[1]).FontSize(t.Rem(0.82)).Bold()
					})
				}
			})
		})
	})
}

func demoContentReadingTime(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Gap(10).Children(func() {
		content.ReadingTime(c, 1610, content.ReadingTimeOptions{})
		content.ReadingTime(c, 1610, content.ReadingTimeOptions{Progress: 0.4})
		content.ReadingTime(c, 420, content.ReadingTimeOptions{Progress: 0.9})
	})
}

func demoContentStickyNote(c *ui.Context, log func(string)) {
	ui.Row(c).Width(520).Gap(8).Children(func() {
		content.StickyNote(c, "Wax every foot ring before glazing. Kiln 2 shelves are new.", content.StickyNoteOptions{Author: "Maya Chen"})
		content.StickyNote(c, "Friday's cone 6 firing is full. Next one is Tuesday.", content.StickyNoteOptions{Author: "Daniel Okafor", Tint: 1, Tilt: 3})
	})
}

func demoContentQuote(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		content.Quote(c, "We stopped double-booking the kiln the week we switched. Members love seeing when their pieces are ready.",
			"Priya Nair", content.QuoteOptions{Role: "Owner, Riverbend Clay Co.", Cite: "Glaze customer stories"})
	})
}

func demoContentStoryTimeline(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		d := func(m time.Month, day int) time.Time { return time.Date(2026, m, day, 9, 0, 0, 0, time.UTC) }
		content.StoryTimeline(c, []content.Beat{
			{Date: d(9, 2), Title: "Thrown on the wheel", Text: "Maya threw twelve mugs from speckled stoneware.", Icon: "hand"},
			{Date: d(9, 9), Title: "Trimmed and handled", Text: "Leather hard, so the handles pulled clean.", Icon: "scissors"},
			{Date: d(9, 30), Title: "Bisque fired", Text: "Cone 04 in kiln 1, shared with the Tuesday class.", Icon: "flame"},
			{Date: d(10, 23), Title: "Glazed and fired", Text: "Cone 6 oxidation; the speckles came through.", Icon: "sparkles"},
		}, content.StoryTimelineOptions{Date: "Jan 2"})
	})
}

func demoContentPrintBar(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		_, print, share := content.PrintBar(c, "https://docs.glaze.studio/kilns/booking", content.PrintBarOptions{})
		if print.Clicked() {
			log("print")
		}
		if share.Clicked() {
			log("link copied")
		}
	})
}
