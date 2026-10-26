package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/content"
	"github.com/egoist/mygo/ui"
)

var contentEntries = []Entry{
	{ID: 1431, Name: "Article", Category: "content", Summary: "Title, byline, body.", Demo: demoContentArticle},
	{ID: 1432, Name: "DocTOC", Category: "content", Summary: "Page outline with scroll marker; cue Select.", Demo: demoContentDocTOC},
	{ID: 1433, Name: "DocSidebar", Category: "content", Summary: "Doc sections tree; cue Select.", Demo: demoContentDocSidebar},
	{ID: 1434, Name: "Admonition", Category: "content", Summary: "Note, tip, warning blocks.", Demo: demoContentAdmonition},
	{ID: 1435, Name: "DocSteps", Category: "content", Summary: "Numbered steps with lines.", Demo: demoContentDocSteps},
	{ID: 1436, Name: "DocCodeTabs", Category: "content", Summary: "Language tabs for code; cue Select.", Demo: demoContentDocCodeTabs},
	{ID: 1437, Name: "APIRow", Category: "content", Summary: "Method, path, description.", Demo: demoContentAPIRow},
	{ID: 1438, Name: "PropTable", Category: "content", Summary: "Props with types and defaults.", Demo: demoContentPropTable},
	{ID: 1439, Name: "Markdown", Category: "content", Summary: "Styled markdown view.", Demo: demoContentMarkdown},
	{ID: 1440, Name: "Footnotes", Category: "content", Summary: "Notes list with back links; cue Tap.", Demo: demoContentFootnotes},
	{ID: 1441, Name: "Glossary", Category: "content", Summary: "Terms by letter.", Demo: demoContentGlossary},
	{ID: 1442, Name: "DocSearch", Category: "content", Summary: "Search with grouped results; cue Select.", Demo: demoContentDocSearch},
	{ID: 1443, Name: "PrevNext", Category: "content", Summary: "Previous and next pages; cue Tap.", Demo: demoContentPrevNext},
}

// docDay is a date in October 2026.
func docDay(d int) time.Time { return time.Date(2026, 10, d, 9, 0, 0, 0, time.UTC) }

func demoContentArticle(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		content.Article(c, "Firing schedules for cone 6 glazes", []string{
			"A slow cool after peak temperature grows crystals in matte glazes and keeps glossy ones clear.",
			"Start with the default schedule, log what comes out of the kiln, and adjust the hold in ten-minute steps.",
		}, content.ArticleOptions{Author: "Priya Nair", Role: "Ceramicist", Tag: "Guides", Cover: productShots()[5], Date: docDay(14), Minutes: 7})
	})
}

func demoContentDocTOC(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320)
	active := intState(root, "active", 0)
	root.Children(func() {
		r := content.DocTOC(c, []content.Heading{
			{Title: "Overview"}, {Title: "Install the CLI"}, {Title: "Authenticate", Depth: 1},
			{Title: "Book a kiln"}, {Title: "Webhooks"}, {Title: "Retries", Depth: 1},
		}, active, content.DocTOCOptions{})
		if r.Changed() {
			log(fmt.Sprintf("active %d", *active))
		}
	})
}

func demoContentDocSidebar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320)
	current := textState(root, "current", "Quickstart")
	root.Children(func() {
		before := *current
		r := content.DocSidebar(c, []content.DocSection{
			{Title: "Getting started", Icon: "rocket", Pages: []string{"Introduction", "Quickstart", "Studio setup"}},
			{Title: "Kilns", Icon: "flame", Pages: []string{"Book a kiln", "Firing logs", "Shelf maps"}},
			{Title: "API", Icon: "code", Pages: []string{"Authentication", "Webhooks", "Rate limits"}},
		}, current)
		if r.Changed() {
			if *current == before {
				log("section toggled")
			} else {
				log("page " + *current)
			}
		}
	})
}

func demoContentAdmonition(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Gap(12).Children(func() {
		content.Admonition(c, content.Note, "", "Bookings lock 24 hours before a firing starts.", content.AdmonitionOptions{})
		content.Admonition(c, content.Tip, "Load by cone", "Group pieces by cone so one firing serves the whole shelf.", content.AdmonitionOptions{})
		content.Admonition(c, content.Warning, "", "Unwaxed feet fuse to the shelf. Check every piece before loading.", content.AdmonitionOptions{})
		content.Admonition(c, content.Danger, "Never open a hot kiln", "Wait until the controller reads below 150°F.", content.AdmonitionOptions{})
	})
}

func demoContentDocSteps(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		content.DocSteps(c, []content.DocStep{
			{Title: "Install the CLI", Text: "brew install glaze, or download it from the releases page."},
			{Title: "Sign in", Text: "glaze login opens your browser and stores a token."},
			{Title: "Pick a studio", Text: "glaze use riverbend sets the studio for later commands."},
			{Title: "Book your first firing", Text: "glaze book kiln-2 --cone 6 --date friday"},
		})
	})
}

func demoContentDocCodeTabs(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	picked := intState(root, "picked", 0)
	root.Children(func() {
		r, copy := content.DocCodeTabs(c, []content.CodeSample{
			{Name: "cURL", Lang: "bash", Code: "curl -X POST https://api.glaze.studio/v1/kilns/kiln-2/bookings \\\n  -H \"Authorization: Bearer $GLAZE_TOKEN\" \\\n  -d cone=6 -d date=2026-10-30"},
			{Name: "Go", Lang: "go", Code: "b, err := client.Bookings.Create(ctx, glaze.Booking{\n\tKiln: \"kiln-2\",\n\tCone: \"6\",\n})\nif err != nil {\n\treturn err\n}"},
			{Name: "JavaScript", Lang: "javascript", Code: "const booking = await glaze.bookings.create({\n  kiln: \"kiln-2\",\n  cone: \"6\",\n});"},
		}, picked, content.DocCodeTabsOptions{})
		if r.Changed() {
			log(fmt.Sprintf("tab %d", *picked))
		}
		if copy.Clicked() {
			log("copied")
		}
	})
}

func demoContentAPIRow(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Gap(8).Children(func() {
		content.APIRow(c, "GET", "/v1/kilns", "List the studio's kilns and their next free slot.", content.APIRowOptions{})
		content.APIRow(c, "POST", "/v1/kilns/{id}/bookings", "Book shelf space in a firing.", content.APIRowOptions{})
		content.APIRow(c, "PATCH", "/v1/bookings/{id}", "Change a booking's cone or shelf count.", content.APIRowOptions{})
		content.APIRow(c, "DELETE", "/v1/bookings/{id}", "Cancel a booking before it locks.", content.APIRowOptions{})
		content.APIRow(c, "GET", "/v1/shelves", "Use /v1/kilns/{id}/shelves.", content.APIRowOptions{Deprecated: true})
	})
}

func demoContentPropTable(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		content.PropTable(c, []content.Prop{
			{Name: "kiln", Type: "string", Text: "The kiln to book, as kiln-2.", Required: true},
			{Name: "cone", Type: "string", Default: "\"6\"", Text: "Target cone of the firing."},
			{Name: "date", Type: "time.Time", Text: "The firing's start; it locks a day before.", Required: true},
			{Name: "shelves", Type: "int", Default: "1", Text: "Shelves to hold, 1 to 4."},
			{Name: "notify", Type: "bool", Default: "true", Text: "Text the member when the kiln cools."},
		}, content.PropTableOptions{Title: "BookingRequest"})
	})
}

func demoContentMarkdown(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		content.Markdown(c, "## Glaze CLI\n\nBook kilns and read firing logs from your terminal.\n\n"+
			"- `glaze book` holds shelf space\n- `glaze logs` shows the last firing\n- `glaze members` lists who is waiting\n\n"+
			"> Bookings lock **24 hours** before a firing.\n\n```bash\nglaze book kiln-2 --cone 6\n```\n",
			content.MarkdownOptions{Title: "README.md"})
	})
}

func demoContentFootnotes(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	back := intState(root, "back", -1)
	root.Children(func() {
		r := content.Footnotes(c, []string{
			"Cone 6 is about 2232°F at a heating rate of 108°F an hour.",
			"Oxidation firings leave iron speckles brown; reduction turns them blue-grey.",
			"Most studios hold for 10 minutes at peak; see Firing logs.",
		}, back, content.FootnotesOptions{})
		if r.Changed() {
			log(fmt.Sprintf("back %d", *back))
		}
	})
}

func demoContentGlossary(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		content.Glossary(c, []content.Term{
			{Word: "Greenware", Meaning: "Clay that is shaped but not yet fired."},
			{Word: "Bisque", Meaning: "Ware after its first, low firing; porous and ready for glaze."},
			{Word: "Cone", Meaning: "A ceramic pyramid that bends at a set heat; names firing temperatures."},
			{Word: "Bone dry", Meaning: "Clay with no water left to lose; ready for the bisque firing."},
			{Word: "Wedging", Meaning: "Kneading clay to drive out air before throwing."},
		}, content.GlossaryOptions{})
	})
}

func demoContentDocSearch(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	query, picked := textState(root, "query", ""), textState(root, "picked", "")
	root.Children(func() {
		before := *picked
		r := content.DocSearch(c, query, []content.DocPage{
			{Section: "Guides", Title: "Book a kiln", Text: "Hold shelf space in a firing."},
			{Section: "Guides", Title: "Firing logs", Text: "Read temperatures and holds from past firings."},
			{Section: "Guides", Title: "Shelf maps", Text: "See which member's work sits where."},
			{Section: "API", Title: "Webhooks", Text: "Get a POST when a firing starts or cools."},
			{Section: "API", Title: "Rate limits", Text: "1,000 requests an hour per token."},
		}, picked, content.DocSearchOptions{})
		if r.Changed() {
			if *picked != before {
				log("picked " + *picked)
			} else {
				log(fmt.Sprintf("query %d runes", len([]rune(*query))))
			}
		}
	})
}

func demoContentPrevNext(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	current := intState(root, "current", 1)
	pages := []string{"Quickstart", "Studio setup", "Book a kiln", "Firing logs"}
	root.Children(func() {
		if content.PrevNext(c, pages, current, content.PrevNextOptions{}).Changed() {
			log("page " + pages[*current])
		}
	})
}
