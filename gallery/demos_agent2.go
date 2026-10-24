package gallery

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/agent"
	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/egoist/mygo/ui"
)

var agentEntries2 = []Entry{
	{881, "ProposedDiff", "agent", "Diff with Apply and Reject; cue Success, Close.", demoProposedDiff},
	{882, "AgentAvatar", "agent", "Avatar ring that pulses while working.", demoAgentAvatar},
	{883, "VoiceOrb", "agent", "Orb whose size follows a level.", demoVoiceOrb},
	{884, "Memories", "agent", "Remembered facts with delete; cue Delete.", demoMemories},
	{885, "AgentStatus", "agent", "Idle, working, blocked, done dot and text.", demoAgentStatus},
	{886, "RunTimeline", "agent", "Steps of a run with durations; bars grow.", demoRunTimeline},
	{887, "AgentBoard", "agent", "Agents as cards with status; cards move by status.", demoAgentBoard},
	{888, "ContextFiles", "agent", "Attached files as chips with remove; cue Delete.", demoContextFiles},
	{889, "SlashCommands", "agent", "Command list for \"/\" in the prompt; cue Select.", demoSlashCommands},
	{890, "SystemPrompt", "agent", "Editor with a token count.", demoSystemPrompt},
	{891, "Temperature", "agent", "Slider with \"Precise\" to \"Creative\" labels; cue Tick.", demoTemperature},
	{892, "EvalScore", "agent", "Scores per check with pass marks.", demoEvalScore},
	{893, "RateAnswer", "agent", "Thumbs then a reason; cue Send.", demoRateAnswer},
	{894, "StreamCode", "agent", "Code block filling as it streams.", demoStreamCode},
	{895, "SearchSteps", "agent", "\"Searching, Reading 3, Writing\" steps.", demoSearchSteps},
	{896, "ImageGenerating", "agent", "Sheen placeholder with a percent.", demoImageGenerating},
	{897, "AgentHandoff", "agent", "Arrow from one agent to another with a note.", demoAgentHandoff},
	{898, "StopButton", "agent", "Square stop that pulses while running; cue Close.", demoStopButton},
	{899, "PromptHistory", "agent", "Up and down recall earlier prompts; cue Tick.", demoPromptHistory},
	{900, "AgentWelcome", "agent", "Greeting, starter prompts and capabilities; rises in.", demoAgentWelcome},
}

func demoProposedDiff(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480).Gap(12)
	decision := ui.Local(root, "decision", func() agent.Decision { return agent.Undecided })
	diff := "@@ -12,6 +12,7 @@ func Price(p Plan) Money {\n \tbase := p.Seats * p.Rate\n-\tif p.Annual {\n-\t\treturn base * 10\n+\tif p.Annual && p.Seats >= 5 {\n+\t\treturn base * 10 * 0.9\n \t}\n+\t// monthly plans keep list price\n \treturn base * 12"
	root.Children(func() {
		_, apply, reject := agent.ProposedDiff(c, "billing/price.go", diff, decision, agent.ProposedDiffOptions{})
		if apply.Clicked() {
			log("applied")
		}
		if reject.Clicked() {
			log("rejected")
		}
		if *decision != agent.Undecided && button.Outline(c, "Show again").Clicked() {
			*decision = agent.Undecided
			log("reset")
		}
	})
}

func demoAgentAvatar(c *ui.Context, log func(string)) {
	ui.Row(c).Width(360).Gap(24).AlignItems(ui.Center).Children(func() {
		for i, a := range []struct {
			name    string
			working bool
		}{{"Planner Bot", true}, {"Researcher", true}, {"Writer", false}} {
			c.Key(i)
			ui.Column(c).Gap(6).AlignItems(ui.Center).Children(func() {
				agent.AgentAvatar(c, a.name, nil, a.working)
				ui.Text(c, a.name).FontSize(c.Theme().Rem(0.8))
			})
		}
	})
}

func demoVoiceOrb(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320).Gap(16).AlignItems(ui.Center)
	level := numberState(root, "level", 0.3)
	root.Children(func() {
		agent.VoiceOrb(c, *level, agent.VoiceOrbOptions{})
		ui.Text(c, "Listening…").FontSize(c.Theme().Rem(0.9))
		ui.Row(c).Gap(8).Children(func() {
			if button.Outline(c, "Quieter").Clicked() {
				*level = max(*level-0.25, 0)
				log(fmt.Sprintf("level %.2f", *level))
			}
			if button.Outline(c, "Louder").Clicked() {
				*level = min(*level+0.25, 1)
				log(fmt.Sprintf("level %.2f", *level))
			}
		})
	})
}

func demoMemories(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	facts := ui.Local(root, "facts", func() []agent.Memory {
		return []agent.Memory{
			{Fact: "Prefers metric units and 24-hour time", Saved: "Saved 3 days ago"},
			{Fact: "Works at Lumen as a product designer", Saved: "Saved last week"},
			{Fact: "Vegetarian; allergic to peanuts", Saved: "Saved in May"},
			{Fact: "Writes release notes in British English", Saved: "Saved in April"},
		}
	})
	root.Children(func() {
		if agent.Memories(c, "Memory", facts, agent.MemoriesOptions{}).Changed() {
			log(fmt.Sprintf("%d left", len(*facts)))
		}
	})
}

func demoAgentStatus(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Gap(12)
	state := ui.Local(root, "state", func() agent.RunState { return agent.Working })
	root.Children(func() {
		ui.Row(c).Gap(8).Wrap().Children(func() {
			for _, s := range []agent.RunState{agent.Idle, agent.Working, agent.Blocked, agent.Done} {
				c.Key(s)
				agent.AgentStatus(c, s, agent.AgentStatusOptions{})
			}
		})
		ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "Deploy agent")
			agent.AgentStatus(c, *state, agent.AgentStatusOptions{})
		})
		if button.Outline(c, "Next state").Clicked() {
			*state = (*state + 1) % 4
			log(fmt.Sprintf("state %d", *state))
		}
	})
}

func demoRunTimeline(c *ui.Context, log func(string)) {
	steps := []agent.RunStep{
		{Name: "Plan the answer", Duration: 2 * time.Second, State: agent.StepDone},
		{Name: "Search the web", Duration: 9 * time.Second, State: agent.StepDone},
		{Name: "Read 6 pages", Duration: 14 * time.Second, State: agent.StepDone},
		{Name: "Run the numbers", Duration: 4 * time.Second, State: agent.StepFailed},
		{Name: "Write the report", Duration: 7 * time.Second, State: agent.StepRunning},
	}
	ui.Column(c).Width(480).Children(func() { agent.RunTimeline(c, "Run #1842", steps, agent.RunTimelineOptions{}) })
}

func demoAgentBoard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520).Gap(12)
	agents := ui.Local(root, "agents", func() []agent.BoardAgent {
		return []agent.BoardAgent{
			{Name: "Alpha", Task: "Drafting the launch email", Status: agent.Working},
			{Name: "Beta", Task: "Summarizing 40 tickets", Status: agent.Working},
			{Name: "Gamma", Task: "Needs a database password", Status: agent.Blocked},
			{Name: "Delta", Task: "Waiting for a task", Status: agent.Idle},
			{Name: "Epsilon", Task: "Shipped the changelog", Status: agent.Done},
		}
	})
	root.Children(func() {
		agent.AgentBoard(c, "Agents", *agents, agent.AgentBoardOptions{})
		ui.Row(c).Gap(8).Children(func() {
			if button.Outline(c, "Unblock Gamma").Clicked() {
				(*agents)[2].Status = agent.Working
				log("gamma working")
			}
			if button.Outline(c, "Finish Alpha").Clicked() {
				(*agents)[0].Status = agent.Done
				log("alpha done")
			}
		})
	})
}

func demoContextFiles(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	files := ui.Local(root, "files", func() []agent.ContextFile {
		return []agent.ContextFile{
			{Name: "brief.pdf", Size: 248_000}, {Name: "price.go", Size: 4_120},
			{Name: "q3-revenue.csv", Size: 18_400}, {Name: "hero.png", Size: 1_240_000},
		}
	})
	root.Children(func() {
		if agent.ContextFiles(c, "Attached files", files, agent.ContextFilesOptions{}).Changed() {
			log(fmt.Sprintf("%d files left", len(*files)))
		}
	})
}

func demoSlashCommands(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400).Gap(8)
	query := textState(root, "query", "/")
	chosen := textState(root, "chosen", "")
	cmds := []agent.Command{
		{Name: "/deploy", Description: "Ship the branch to production", Icon: "rocket"},
		{Name: "/review", Description: "Review the open pull request", Icon: "eye"},
		{Name: "/summarize", Description: "Summarize this conversation", Icon: "list"},
		{Name: "/translate", Description: "Translate the last answer", Icon: "languages"},
		{Name: "/image", Description: "Generate an image", Icon: "image"},
	}
	root.Children(func() {
		if ui.TextInput(c, query).Label("Command").FillWidth().Font("monospace").Changed() {
			log(fmt.Sprintf("query %d runes", utf8.RuneCountInString(*query)))
		}
		if agent.SlashCommands(c, "Commands", cmds, *query, chosen, agent.SlashCommandsOptions{}).Changed() {
			log("selected " + *chosen)
		}
	})
}

func demoSystemPrompt(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	text := textState(root, "text", "You are Lumen's support agent. Answer in two short paragraphs, cite the help center and never promise refunds.")
	root.Children(func() {
		if agent.SystemPrompt(c, "System prompt", text, agent.SystemPromptOptions{}).Changed() {
			log(fmt.Sprintf("edited %d runes", utf8.RuneCountInString(*text)))
		}
	})
}

func demoTemperature(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	value := numberState(root, "value", 0.5)
	root.Children(func() {
		if agent.Temperature(c, value, agent.TemperatureOptions{}).Changed() {
			log(fmt.Sprintf("value %.2f", *value))
		}
	})
}

func demoEvalScore(c *ui.Context, log func(string)) {
	checks := []agent.EvalCheck{
		{Name: "Factual accuracy", Score: 0.94, Pass: 0.85},
		{Name: "Cites sources", Score: 0.88, Pass: 0.8},
		{Name: "Tone matches brand", Score: 0.62, Pass: 0.7},
		{Name: "No unsafe advice", Score: 1, Pass: 0.99},
	}
	ui.Column(c).Width(480).Children(func() { agent.EvalScore(c, "Support eval · v14", checks, agent.EvalScoreOptions{}) })
}

func demoRateAnswer(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	rating, reason := intState(root, "rating", 0), intState(root, "reason", -1)
	sent := boolState(root, "sent", false)
	reasons := []string{"Inaccurate", "Too long", "Missed the question", "Unsafe"}
	root.Children(func() {
		box, send := agent.RateAnswer(c, rating, reasons, reason, sent, agent.RateAnswerOptions{})
		if send.Clicked() && *sent {
			log(fmt.Sprintf("rating %d, reason %d", *rating, *reason))
		} else if box.Changed() {
			log(fmt.Sprintf("rating %d", *rating))
		}
		if *sent && button.Outline(c, "Rate again").Clicked() {
			*rating, *reason, *sent = 0, -1, false
			log("reset")
		}
	})
}

func demoStreamCode(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480).Gap(12)
	restart := intState(root, "restart", 0)
	code := "func Price(p Plan) Money {\n\tbase := p.Seats * p.Rate\n\tif p.Annual {\n\t\treturn base * 10\n\t}\n\treturn base * 12\n}\n"
	root.Children(func() {
		ui.Column(c).FillWidth().Key(*restart).Children(func() {
			if _, cp := agent.StreamCode(c, "Go", code, false, agent.StreamCodeOptions{}); cp.Clicked() {
				log("copied 7 lines")
			}
		})
		if button.Outline(c, "Restart").Clicked() {
			*restart++
			log("restarted")
		}
	})
}

func demoSearchSteps(c *ui.Context, log func(string)) {
	steps := []agent.SearchStep{
		{Label: "Searching", Count: 24, State: agent.StepDone},
		{Label: "Reading", Count: 3, Detail: "of 8", State: agent.StepRunning},
		{Label: "Writing", State: agent.StepPending},
	}
	ui.Column(c).Width(520).Children(func() { agent.SearchSteps(c, "Research", steps) })
}

func demoImageGenerating(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400).Gap(12)
	progress := numberState(root, "progress", 0.4)
	root.Children(func() {
		agent.ImageGenerating(c, "A lighthouse at dusk, film grain, 35mm", *progress)
		if button.Outline(c, "10% more").Clicked() {
			*progress = min(*progress+0.1, 1)
			log(fmt.Sprintf("%.0f%%", *progress*100))
		}
	})
}

func demoAgentHandoff(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	at := timeState(root, "at", c)
	root.Children(func() {
		agent.AgentHandoff(c, "Researcher", "Writer", "Gathered 8 sources on EU pricing; draft the comparison table next.", *at)
	})
}

func demoStopButton(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380).Gap(12)
	running := boolState(root, "running", true)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(12).AlignItems(ui.Center).Children(func() {
			label := "Writing the report…"
			if !*running {
				label = "Stopped"
			}
			ui.Text(c, label).Grow(1)
			if agent.StopButton(c, *running, agent.StopButtonOptions{}).Changed() {
				*running = false
				log("stopped")
			}
		})
		if !*running && button.Outline(c, "Run again").Clicked() {
			*running = true
			log("running")
		}
	})
}

func demoPromptHistory(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	at := intState(root, "at", 0)
	entries := []string{
		"Summarize the Q3 revenue sheet",
		"Draft a reply to Kai about the deck",
		"Translate the release note to Japanese",
		"List risks for the May launch",
		"Plan a three-day trip to Kyoto",
	}
	root.Children(func() {
		if agent.PromptHistory(c, entries, at, agent.PromptHistoryOptions{}).Changed() {
			log(fmt.Sprintf("at %d", *at+1))
		}
	})
}

func demoAgentWelcome(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	chosen := intState(root, "chosen", -1)
	prompts := []agent.Suggestion{
		{Text: "Draft the deck", Icon: "presentation"},
		{Text: "Plan a trip", Icon: "map"},
		{Text: "Explain a chart", Icon: "chart-line"},
	}
	caps := []agent.Capability{
		{Icon: "search", Title: "Research", Text: "Reads and cites the web"},
		{Icon: "code", Title: "Code", Text: "Writes and runs Go"},
		{Icon: "image", Title: "Images", Text: "Draws from a prompt"},
	}
	root.Children(func() {
		if agent.AgentWelcome(c, "Good morning, Ana", prompts, caps, chosen, agent.AgentWelcomeOptions{Subtitle: "What should we work on today?"}).Changed() {
			log("picked " + prompts[*chosen].Text)
		}
	})
}
