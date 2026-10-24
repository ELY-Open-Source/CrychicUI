package gallery

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/agent"
	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

var agentEntries = []Entry{
	{861, "PromptBox", "agent", "Large prompt input with model and attach; cue Send.", demoPromptBox},
	{862, "StreamText", "agent", "Text appearing token by token with a caret.", demoStreamText},
	{863, "Thinking", "agent", "Shimmering \"Thinking\" label with elapsed time.", demoThinking},
	{864, "Reasoning", "agent", "Collapsible reasoning panel; height springs; cue Toggle.", demoReasoning},
	{865, "ToolCall", "agent", "Tool name, arguments and status; status icon spins then checks.", demoToolCall},
	{866, "ToolResult", "agent", "Collapsible result body; cue Toggle.", demoToolResult},
	{867, "Plan", "agent", "Steps with pending, running, done marks; marks animate.", demoPlan},
	{868, "TaskProgress", "agent", "Current step and bar; bar fills.", demoTaskProgress},
	{869, "Citation", "agent", "Numbered chip opening a source; cue Open.", demoCitation},
	{870, "SourceCard", "agent", "Title, host and snippet; lifts on hover; cue Tap.", demoSourceCard},
	{871, "ModelPicker", "agent", "Models with speed and cost tags; cue Select.", demoModelPicker},
	{872, "TokenCount", "agent", "Tokens used of limit; digits roll.", demoTokenCount},
	{873, "ContextMeter", "agent", "Ring of context used; turns warning near full.", demoContextMeter},
	{874, "CostMeter", "agent", "Spend so far; digits roll.", demoCostMeter},
	{875, "Suggestions", "agent", "Prompt chips; tap fills the box; cue Tap.", demoSuggestions},
	{876, "PromptLibrary", "agent", "Saved prompts with search; cue Select.", demoPromptLibrary},
	{877, "Branches", "agent", "\"2 / 3\" switcher between answers; cue Tick.", demoBranches},
	{878, "AnswerActions", "agent", "Copy, retry, like, dislike; cue Pop.", demoAnswerActions},
	{879, "Artifact", "agent", "Card of a generated file with open; cue Open.", demoArtifact},
	{880, "Approval", "agent", "\"Allow this action?\" with Allow and Deny; cue Success, Close.", demoApproval},
}

// agentSources are pages a research answer cites.
func agentSources() []agent.Source {
	return []agent.Source{
		{Title: "Layout: rows, columns and grids", Host: "mygo.egoist.dev", URL: "https://mygo.egoist.dev/docs/ui", Snippet: "Every element lays out its children in a row, a column or a grid."},
		{Title: "Go 1.27 release notes", Host: "go.dev", URL: "https://go.dev/doc/go1.27", Snippet: "Iterators over maps now keep insertion order in the runtime."},
		{Title: "Designing calm interfaces", Host: "uxdesign.cc", URL: "https://uxdesign.cc/calm", Snippet: "Quiet defaults, loud exceptions: motion only where it explains."},
	}
}

func demoPromptBox(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	draft := textState(root, "draft", "")
	root.Children(func() {
		box, _, attach := agent.PromptBox(c, draft, "Claude Opus 5.5", func(s string) {
			log(fmt.Sprintf("sent %d runes", utf8.RuneCountInString(s)))
		}, agent.PromptBoxOptions{Placeholder: "Ask about the launch plan"})
		if box.Changed() {
			log(fmt.Sprintf("draft %d runes", utf8.RuneCountInString(*draft)))
		}
		if attach.Clicked() {
			log("attach")
		}
	})
}

func demoStreamText(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	restart := intState(root, "restart", 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(12).Children(func() {
			ui.Box(c).Size(32, 32).Shrink(0).Center().Radius(999).Background(c.Theme().Accent).TextColor(c.Theme().AccentText).
				Children(func() { ui.Text(c, "AI").FontSize(c.Theme().Rem(0.7)).Bold() })
			ui.Column(c).Grow(1).Key(*restart).Children(func() {
				agent.StreamText(c, "The deck is ready. It covers the launch, the risks and the next three milestones.", false)
			})
		})
		if button.Outline(c, "Restart").Clicked() {
			*restart++
			log("restarted")
		}
	})
}

func demoThinking(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Gap(12)
	since := timeState(root, "since", c)
	root.Children(func() {
		kit.Muted(c, "Plan a three-day trip to Kyoto in April")
		agent.Thinking(c, *since, agent.ThinkingOptions{})
	})
}

func demoReasoning(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(8)
	open := boolState(root, "open", false)
	root.Children(func() {
		body := "Ana, Kai and Mei share a free hour at noon on Friday. Kai has a flight at 4 pm, so lunch must end by 1:30. " +
			"The café on Pier 3 seats six and takes bookings, so I'll suggest it."
		if agent.Reasoning(c, body, open, agent.ReasoningOptions{Took: 14 * time.Second}).Changed() {
			log(fmt.Sprintf("open %v", *open))
		}
		ui.Text(c, "Let's meet at noon on Friday at the café on Pier 3.").Wrap()
	})
}

func demoToolCall(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460).Gap(8)
	done := boolState(root, "done", false)
	open := ui.Local(root, "open", func() []bool { return []bool{false, false, true} })
	calls := []struct{ name, args string }{
		{"web_search", `{"query": "Kyoto cherry blossom forecast 2027"}`},
		{"read_file", `{"path": "notes/trip.md"}`},
		{"write_file", `{"path": "plans/kyoto.md", "lines": 64}`},
	}
	root.Children(func() {
		for i, k := range calls {
			state := agent.StepDone
			if i == 2 && !*done {
				state = agent.StepRunning
			}
			c.Key(k.name)
			if agent.ToolCall(c, k.name, k.args, state, &(*open)[i], agent.ToolCallOptions{}).Changed() {
				log(fmt.Sprintf("%s open %v", k.name, (*open)[i]))
			}
		}
		label := "Finish write"
		if *done {
			label = "Run again"
		}
		if button.Outline(c, label).Clicked() {
			*done = !*done
			log(fmt.Sprintf("done %v", *done))
		}
	})
}

func demoToolResult(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	open := boolState(root, "open", true)
	body := "plans/kyoto.md        4.2 KB\nplans/osaka.md        2.9 KB\nnotes/trip.md         1.1 KB\nnotes/budget.csv      0.6 KB\nphotos/fushimi.jpg    2.4 MB\nphotos/arashiyama.jpg 3.1 MB"
	root.Children(func() {
		box, cp := agent.ToolResult(c, "list_files", body, open, agent.ToolResultOptions{})
		if box.Changed() {
			log(fmt.Sprintf("open %v", *open))
		}
		if cp.Clicked() {
			log("copied 6 lines")
		}
	})
}

func demoPlan(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	at := intState(root, "at", 1)
	titles := []struct{ title, detail string }{
		{"Read the brief", "Launch on 12 May, three markets"},
		{"Research competitors", "4 sources so far"},
		{"Draft the deck", "12 slides"},
		{"Review with Ana", "Friday 10:00"},
	}
	root.Children(func() {
		steps := make([]agent.PlanStep, len(titles))
		for i, s := range titles {
			steps[i] = agent.PlanStep{Title: s.title, Detail: s.detail}
			switch {
			case i < *at:
				steps[i].State = agent.StepDone
			case i == *at:
				steps[i].State = agent.StepRunning
			}
		}
		agent.Plan(c, "Launch deck", steps, agent.PlanOptions{})
		if button.Outline(c, "Next step").Clicked() {
			*at = min(*at+1, len(titles))
			log(fmt.Sprintf("step %d", *at))
		}
	})
}

func demoTaskProgress(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400).Gap(12)
	step := intState(root, "step", 2)
	steps := []string{"Gathering notes", "Outlining", "Drafting slides", "Adding charts", "Polishing"}
	root.Children(func() {
		agent.TaskProgress(c, "Launch deck", steps, *step, agent.TaskProgressOptions{})
		if button.Outline(c, "Advance").Clicked() {
			*step = min(*step+1, len(steps))
			log(fmt.Sprintf("step %d of %d", *step, len(steps)))
		}
	})
}

func demoCitation(c *ui.Context, log func(string)) {
	src := agentSources()
	ui.Column(c).Width(420).Gap(4).Children(func() {
		ui.Row(c).FillWidth().Wrap().Gap(4).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "Every element lays out its children in a row, a column or a grid")
			if agent.Citation(c, 1, src[0]).Changed() {
				log("opened 1")
			}
			ui.Text(c, "and Go now keeps map order")
			if agent.Citation(c, 2, src[1]).Changed() {
				log("opened 2")
			}
		})
	})
}

func demoSourceCard(c *ui.Context, log func(string)) {
	src := agentSources()
	ui.Column(c).Width(440).Gap(12).Children(func() {
		for i, s := range src {
			c.Key(s.URL)
			card := agent.SourceCard(c, i+1, s)
			kit.Arrive(card, ui.Pose{Y: 8, Fade: 1}, time.Duration(i)*60*time.Millisecond)
			if card.Changed() {
				log(fmt.Sprintf("opened %d", i+1))
			}
		}
	})
}

func demoModelPicker(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	selected := intState(root, "selected", 0)
	models := []agent.Model{
		{Name: "Opus 5.5", Description: "Deep reasoning for long tasks", Speed: "Steady", Cost: "$$$"},
		{Name: "Sonnet 5.5", Description: "Balanced for everyday work", Speed: "Fast", Cost: "$$"},
		{Name: "Haiku 5.5", Description: "Quick answers and drafts", Speed: "Fastest", Cost: "$"},
		{Name: "Local 8B", Description: "Runs on this machine, offline", Speed: "Slow", Cost: "Free"},
	}
	root.Children(func() {
		if agent.ModelPicker(c, "Models", models, selected).Changed() {
			log("selected " + models[*selected].Name)
		}
	})
}

func demoTokenCount(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320).Gap(12)
	used := intState(root, "used", 2848)
	root.Children(func() {
		agent.TokenCount(c, *used, 4096, agent.TokenCountOptions{})
		if button.Outline(c, "Add 512").Clicked() {
			*used += 512
			log(fmt.Sprintf("%d tokens", *used))
		}
	})
}

func demoContextMeter(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Gap(12)
	used := intState(root, "used", 96_000)
	root.Children(func() {
		agent.ContextMeter(c, *used, 200_000, agent.ContextMeterOptions{})
		if button.Outline(c, "Add 20k tokens").Clicked() {
			*used = min(*used+20_000, 200_000)
			log(fmt.Sprintf("%d tokens", *used))
		}
	})
}

func demoCostMeter(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320).Gap(12)
	amount := numberState(root, "amount", 31.6)
	root.Children(func() {
		agent.CostMeter(c, *amount, 50, "USD", agent.CostMeterOptions{})
		if button.Outline(c, "Run a batch").Clicked() {
			*amount += 4.8
			log(fmt.Sprintf("USD %.2f", *amount))
		}
	})
}

func demoSuggestions(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480).Gap(12)
	draft := textState(root, "draft", "")
	prompts := []agent.Suggestion{
		{Text: "Summarize this thread", Icon: "list"},
		{Text: "Draft a polite reply", Icon: "pen-line"},
		{Text: "List next steps", Icon: "list-checks"},
		{Text: "Explain the error", Icon: "bug"},
	}
	root.Children(func() {
		if agent.Suggestions(c, "Starters", prompts, draft).Changed() {
			log("fill: " + *draft)
		}
		box, _, attach := agent.PromptBox(c, draft, "Sonnet 5.5", func(s string) {
			log(fmt.Sprintf("sent %d runes", utf8.RuneCountInString(s)))
		}, agent.PromptBoxOptions{})
		if box.Changed() {
			log(fmt.Sprintf("draft %d runes", utf8.RuneCountInString(*draft)))
		}
		if attach.Clicked() {
			log("attach")
		}
	})
}

func demoPromptLibrary(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	query := textState(root, "query", "")
	chosen := intState(root, "chosen", -1)
	prompts := []agent.SavedPrompt{
		{Title: "Release note", Body: "Turn these commits into a release note for customers.", Tag: "Writing"},
		{Title: "Code review", Body: "Review this pull request for bugs and missing tests.", Tag: "Engineering"},
		{Title: "Sprint plan", Body: "Plan a two-week sprint from this backlog.", Tag: "Planning"},
		{Title: "Interview recap", Body: "Summarize the interview and list risks.", Tag: "Hiring"},
		{Title: "Bug triage", Body: "Sort these issues by severity and owner.", Tag: "Engineering"},
	}
	root.Children(func() {
		box, field := agent.PromptLibrary(c, "Library", prompts, query, chosen, agent.PromptLibraryOptions{})
		if field.Changed() {
			log(fmt.Sprintf("query %d runes", utf8.RuneCountInString(*query)))
		}
		if box.Changed() {
			log("selected " + prompts[*chosen].Title)
		}
	})
}

func demoBranches(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(8)
	current := intState(root, "current", 2)
	answers := []string{
		"Book the café on Pier 3 for noon on Friday.",
		"Order in at the office at noon; Kai can leave early.",
		"Move lunch to Thursday when everyone is free for two hours.",
	}
	root.Children(func() {
		kit.Slide(c, *current, func() { ui.Text(c, answers[*current-1]).Wrap() })
		_, back, forth := agent.Branches(c, current, len(answers), agent.BranchesOptions{})
		if back.Clicked() || forth.Clicked() {
			log(fmt.Sprintf("at %d of %d", *current, len(answers)))
		}
	})
}

func demoAnswerActions(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(8)
	vote := intState(root, "vote", 0)
	answer := "The deck is ready for review. I kept the risks slide and moved pricing to the appendix."
	root.Children(func() {
		ui.Text(c, answer).Wrap()
		box, cp, retry := agent.AnswerActions(c, answer, vote, agent.AnswerActionsOptions{})
		if box.Changed() {
			log(fmt.Sprintf("vote %d", *vote))
		}
		if cp.Clicked() {
			log(fmt.Sprintf("copied %d runes", utf8.RuneCountInString(answer)))
		}
		if retry.Clicked() {
			log("retry")
		}
	})
}

func demoArtifact(c *ui.Context, log func(string)) {
	ui.Column(c).Width(400).Gap(12).Children(func() {
		items := []struct {
			name string
			kind agent.ArtifactKind
			meta string
		}{{"launch-deck.html", agent.Web, "Web page · 18 KB"}, {"budget.csv", agent.Sheet, "Sheet · 42 rows"},
			{"kyoto-itinerary.md", agent.Document, "Document · 3 pages"}}
		for i, a := range items {
			c.Key(a.name)
			card, open := agent.Artifact(c, a.name, a.kind, a.meta, agent.ArtifactOptions{})
			kit.Arrive(card, ui.Pose{Y: 8, Fade: 1}, time.Duration(i)*60*time.Millisecond)
			if open.Clicked() {
				log("opened " + a.name)
			}
		}
	})
}

func demoApproval(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440).Gap(12)
	decision := ui.Local(root, "decision", func() agent.Decision { return agent.Undecided })
	root.Children(func() {
		_, allow, deny := agent.Approval(c, "shell", "git push --force origin main", "Rewrites 3 commits on the shared branch.", decision, agent.ApprovalOptions{})
		if allow.Clicked() {
			log("allowed")
		}
		if deny.Clicked() {
			log("denied")
		}
		if *decision != agent.Undecided && button.Outline(c, "Ask again").Clicked() {
			*decision = agent.Undecided
			log("asked again")
		}
	})
}
