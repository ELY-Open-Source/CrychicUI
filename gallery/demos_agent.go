package gallery

import (
	"fmt"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/agent"
	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/egoist/mygo/ui"
)

var agentEntries = []Entry{
	{861, "PromptBox", "agent", "A large prompt input with model and attach.", demoPromptBox},
	{862, "StreamText", "agent", "Text appearing token by token with a caret.", demoStreamText},
	{863, "Thinking", "agent", "A shimmering label with the elapsed time.", demoThinking},
	{864, "Reasoning", "agent", "A collapsible reasoning panel.", demoReasoning},
	{865, "ToolCall", "agent", "A tool name, arguments and status.", demoToolCall},
	{866, "ToolResult", "agent", "A collapsible result body.", demoToolResult},
	{867, "Plan", "agent", "Steps with pending, running and done marks.", demoPlan},
	{868, "TaskProgress", "agent", "The current step and a filling bar.", demoTaskProgress},
	{869, "Citation", "agent", "A numbered chip opening a source.", demoCitation},
	{870, "SourceCard", "agent", "A title, host and snippet; lifts on hover.", demoSourceCard},
	{871, "ModelPicker", "agent", "Models with speed and cost tags.", demoModelPicker},
	{872, "TokenCount", "agent", "Tokens used of a limit; digits roll.", demoTokenCount},
	{873, "ContextMeter", "agent", "A ring of context used; warns near full.", demoContextMeter},
	{874, "CostMeter", "agent", "The spend so far; digits roll.", demoCostMeter},
	{875, "Suggestions", "agent", "Prompt chips; a tap fills the box.", demoSuggestions},
	{876, "PromptLibrary", "agent", "Saved prompts with search.", demoPromptLibrary},
	{877, "Branches", "agent", "A position switcher between answers.", demoBranches},
	{878, "AnswerActions", "agent", "Copy, retry, like and dislike.", demoAnswerActions},
	{879, "Artifact", "agent", "A card of a generated file with open.", demoArtifact},
	{880, "Approval", "agent", "Allow or deny an action.", demoApproval},
}

func demoPromptBox(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	draft := ui.Local(root, "draft", func() string { return "" })
	root.Children(func() {
		box, _, attach := agent.PromptBox(c, draft, "GPT-6.1 Sol", func(s string) {
			log(fmt.Sprintf("sent %d runes", utf8.RuneCountInString(s)))
		}, agent.PromptBoxOptions{})
		if box.Changed() {
			log(fmt.Sprintf("draft %d chars", len(*draft)))
		}
		if attach.Clicked() {
			log("attach")
		}
	})
}

func demoStreamText(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	restart := intState(root, "restart", 0)
	root.Children(func() {
		c.Key(*restart)
		agent.StreamText(c, "The deck is ready. It covers the launch, the risks and the next three milestones.", false)
		if button.Outline(c, "Restart").Clicked() {
			*restart++
			log("restarted")
		}
	})
}

func demoThinking(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	since := timeState(root, "since", c)
	root.Children(func() { agent.Thinking(c, *since, agent.ThinkingOptions{}) })
}

func demoReasoning(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	open := boolState(root, "open", false)
	root.Children(func() {
		if agent.Reasoning(c, "The user asked about lunch; noon on Friday works for everyone.", open, agent.ReasoningOptions{}).Changed() {
			log(fmt.Sprintf("open %v", *open))
		}
	})
}

func demoToolCall(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Gap(12)
	running := boolState(root, "running", true)
	root.Children(func() {
		agent.ToolCall(c, "search", `query "launch deck"`, *running)
		if button.Outline(c, "Finish").Clicked() {
			*running = false
			log("done")
		}
	})
}

func demoToolResult(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	open := boolState(root, "open", true)
	root.Children(func() {
		if agent.ToolResult(c, "deck-v3.html\ndeck-notes.md\nlaunch-plan.pdf", open, agent.ToolResultOptions{}).Changed() {
			log(fmt.Sprintf("open %v", *open))
		}
	})
}

func demoPlan(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Gap(12)
	current := intState(root, "current", 0)
	steps := []string{"Read the brief", "Draft the deck", "Review together", "Ship it"}
	root.Children(func() {
		agent.Plan(c, "Plan", steps, min(*current, len(steps)-1))
		if button.Outline(c, "Next step").Clicked() {
			*current = min(*current+1, len(steps)-1)
			log(fmt.Sprintf("step %d", *current+1))
		}
	})
}

func demoTaskProgress(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Gap(12)
	step := intState(root, "step", 2)
	root.Children(func() {
		agent.TaskProgress(c, "Deck", *step, 5, agent.TaskProgressOptions{})
		if button.Outline(c, "Advance").Clicked() {
			*step = min(*step+1, 5)
			log(fmt.Sprintf("step %d of 5", *step))
		}
	})
}

func demoCitation(c *ui.Context, log func(string)) {
	ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
		ui.Text(c, "Per the docs")
		if agent.Citation(c, 1, "https://mygo.egoist.dev/docs/ui").Changed() {
			log("opened 1")
		}
		ui.Text(c, "lists follow the end.")
	})
}

func demoSourceCard(c *ui.Context, log func(string)) {
	ui.Column(c).Children(func() {
		if agent.SourceCard(c, "MyGo: native UI for Go", "mygo.egoist.dev", "One Go module for desktop and browser.").Changed() {
			log("opened the source")
		}
	})
}

func demoModelPicker(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	selected := intState(root, "selected", 0)
	models := []agent.Model{
		{Name: "GPT-6.1 Sol", Speed: "Fast", Cost: "$$$"},
		{Name: "Claude Fable", Speed: "Fast", Cost: "$$"},
		{Name: "Local 8B", Speed: "Slow", Cost: "Free"},
	}
	root.Children(func() {
		if agent.ModelPicker(c, "Models", models, selected).Changed() {
			log("selected " + models[*selected].Name)
		}
	})
}

func demoTokenCount(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	used := intState(root, "used", 1248)
	root.Children(func() {
		agent.TokenCount(c, *used, 4096, agent.TokenCountOptions{})
		if button.Outline(c, "Add 256").Clicked() {
			*used += 256
			log(fmt.Sprintf("%d tokens", *used))
		}
	})
}

func demoContextMeter(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	used := numberState(root, "used", 0.45)
	root.Children(func() {
		agent.ContextMeter(c, *used, agent.ContextMeterOptions{})
		if button.Outline(c, "Use 10% more").Clicked() {
			*used = min(*used+0.1, 1)
			log(fmt.Sprintf("%.0f%% used", *used*100))
		}
	})
}

func demoCostMeter(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	amount := numberState(root, "amount", 12.4)
	root.Children(func() {
		agent.CostMeter(c, *amount, "USD", agent.CostMeterOptions{})
		if button.Outline(c, "Add $0.40").Clicked() {
			*amount += 0.4
			log(fmt.Sprintf("USD %.2f", *amount))
		}
	})
}

func demoSuggestions(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	chosen := intState(root, "chosen", -1)
	prompts := []string{"Summarize the thread", "Draft a reply", "List next steps"}
	root.Children(func() {
		if agent.Suggestions(c, "Suggestions", prompts, chosen).Changed() {
			log("fill: " + prompts[*chosen])
		}
	})
}

func demoPromptLibrary(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	query := ui.Local(root, "query", func() string { return "" })
	chosen := intState(root, "chosen", -1)
	prompts := []string{"Deploy to production", "Draft a release note", "Review a pull request", "Plan the sprint"}
	root.Children(func() {
		box, field := agent.PromptLibrary(c, "Library", prompts, query, chosen, agent.PromptLibraryOptions{})
		if field.Changed() {
			log(fmt.Sprintf("query %d chars", len(*query)))
		}
		if box.Changed() {
			log("selected " + prompts[*chosen])
		}
	})
}

func demoBranches(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	current := ui.Local(root, "current", func() int { return 2 })
	root.Children(func() {
		_, back, forth := agent.Branches(c, current, 3, agent.BranchesOptions{})
		if back.Clicked() || forth.Clicked() {
			log(fmt.Sprintf("at %d of 3", *current))
		}
	})
}

func demoAnswerActions(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	vote := intState(root, "vote", 0)
	root.Children(func() {
		box, cp, retry := agent.AnswerActions(c, "The deck is ready for review.", vote, agent.AnswerActionsOptions{})
		if box.Changed() {
			log(fmt.Sprintf("vote %d", *vote))
		}
		if cp.Clicked() {
			log(fmt.Sprintf("copied %d runes", utf8.RuneCountInString("The deck is ready for review.")))
		}
		if retry.Clicked() {
			log("retry")
		}
	})
}

func demoArtifact(c *ui.Context, log func(string)) {
	ui.Column(c).Children(func() {
		_, open := agent.Artifact(c, "deck.html", "Web page", agent.ArtifactOptions{})
		if open.Clicked() {
			log("opened deck.html")
		}
	})
}

func demoApproval(c *ui.Context, log func(string)) {
	ui.Column(c).Width(380).Children(func() {
		_, allow, deny := agent.Approval(c, "git push --force origin main", agent.ApprovalOptions{})
		if allow.Clicked() {
			log("allowed")
		}
		if deny.Clicked() {
			log("denied")
		}
	})
}
