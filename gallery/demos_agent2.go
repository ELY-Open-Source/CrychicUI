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
	{881, "ProposedDiff", "agent", "A diff with Apply and Reject.", demoProposedDiff},
	{882, "AgentAvatar", "agent", "An avatar ring that pulses while working.", demoAgentAvatar},
	{883, "VoiceOrb", "agent", "An orb whose size follows a level.", demoVoiceOrb},
	{884, "Memories", "agent", "Remembered facts with delete.", demoMemories},
	{885, "AgentStatus", "agent", "A state dot and text.", demoAgentStatus},
	{886, "RunTimeline", "agent", "Steps of a run with duration bars.", demoRunTimeline},
	{887, "AgentBoard", "agent", "Agents as cards that move by status.", demoAgentBoard},
	{888, "ContextFiles", "agent", "Attached files as chips with remove.", demoContextFiles},
	{889, "SlashCommands", "agent", "The command list of a / prompt.", demoSlashCommands},
	{890, "SystemPrompt", "agent", "An editor with a token count.", demoSystemPrompt},
	{891, "Temperature", "agent", "A slider from Precise to Creative.", demoTemperature},
	{892, "EvalScore", "agent", "Scores per check with pass marks.", demoEvalScore},
	{893, "RateAnswer", "agent", "Thumbs, then a reason, then send.", demoRateAnswer},
	{894, "StreamCode", "agent", "A code block filling as it streams.", demoStreamCode},
	{895, "SearchSteps", "agent", "Search phases; the active one spins.", demoSearchSteps},
	{896, "ImageGenerating", "agent", "A sheen placeholder with a percent.", demoImageGenerating},
	{897, "AgentHandoff", "agent", "An arrow between agents with a note.", demoAgentHandoff},
	{898, "StopButton", "agent", "A square stop pulsing while running.", demoStopButton},
	{899, "PromptHistory", "agent", "Up and down recall earlier prompts.", demoPromptHistory},
	{900, "AgentWelcome", "agent", "A greeting with starters and capabilities.", demoAgentWelcome},
}

func demoProposedDiff(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() {
		_, apply, reject := agent.ProposedDiff(c, "deck.md", "@@ -1,2 +1,2 @@\n-Launch in June\n+Launch in July", agent.ProposedDiffOptions{})
		if apply.Clicked() {
			log("applied")
		}
		if reject.Clicked() {
			log("rejected")
		}
	})
}

func demoAgentAvatar(c *ui.Context, log func(string)) {
	ui.Column(c).Children(func() { agent.AgentAvatar(c, "Planner", nil, true) })
}

func demoVoiceOrb(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	level := numberState(root, "level", 0.3)
	root.Children(func() {
		agent.VoiceOrb(c, *level)
		if button.Outline(c, "Louder").Clicked() {
			*level = min(*level+0.5, 1)
			log(fmt.Sprintf("level %.1f", *level))
		}
	})
}

func demoMemories(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	facts := ui.Local(root, "facts", func() []string {
		return []string{"Likes green tea", "Works from Duluth", "Reviews on Fridays"}
	})
	root.Children(func() {
		if agent.Memories(c, "Memories", facts, agent.MemoriesOptions{}).Changed() {
			log(fmt.Sprintf("%d left", len(*facts)))
		}
	})
}

func demoAgentStatus(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	state := ui.Local(root, "state", func() agent.RunState { return agent.Working })
	root.Children(func() {
		agent.AgentStatus(c, *state, agent.AgentStatusOptions{})
		if button.Outline(c, "Finish").Clicked() {
			*state = agent.Done
			log(fmt.Sprintf("state %d", *state))
		}
	})
}

func demoRunTimeline(c *ui.Context, log func(string)) {
	steps := []agent.RunStep{{Name: "Search", Duration: 2 * time.Second}, {Name: "Read", Duration: time.Second}, {Name: "Write", Duration: 3 * time.Second}}
	ui.Column(c).Width(360).Children(func() { agent.RunTimeline(c, "Last run", steps) })
}

func demoAgentBoard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	agents := ui.Local(root, "agents", func() []agent.BoardAgent {
		return []agent.BoardAgent{{Name: "Alpha", Status: agent.Idle}, {Name: "Beta", Status: agent.Working}, {Name: "Gamma", Status: agent.Blocked}}
	})
	root.Children(func() {
		agent.AgentBoard(c, "Board", *agents, agent.AgentBoardOptions{})
		if button.Outline(c, "Unblock Gamma").Clicked() {
			(*agents)[2].Status = agent.Working
			log("gamma working")
		}
	})
}

func demoContextFiles(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	files := ui.Local(root, "files", func() []string { return []string{"deck.md", "brief.pdf", "notes.txt"} })
	root.Children(func() {
		if agent.ContextFiles(c, "Files", files, agent.ContextFilesOptions{}).Changed() {
			log(fmt.Sprintf("%d files left", len(*files)))
		}
	})
}

func demoSlashCommands(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	chosen := intState(root, "chosen", -1)
	commands := []string{"/deploy", "/review", "/test"}
	root.Children(func() {
		if agent.SlashCommands(c, "Commands", commands, chosen).Changed() {
			log("selected " + commands[*chosen])
		}
	})
}

func demoSystemPrompt(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	text := ui.Local(root, "text", func() string { return "You write launch decks that land." })
	root.Children(func() {
		if agent.SystemPrompt(c, "System prompt", text, agent.SystemPromptOptions{}).Changed() {
			log(fmt.Sprintf("edited %d runes", utf8.RuneCountInString(*text)))
		}
	})
}

func demoTemperature(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	value := numberState(root, "value", 0.5)
	root.Children(func() {
		if agent.Temperature(c, value, agent.TemperatureOptions{}).Changed() {
			log(fmt.Sprintf("value %.2f", *value))
		}
	})
}

func demoEvalScore(c *ui.Context, log func(string)) {
	checks := []agent.EvalCheck{{Name: "Accuracy", Score: 0.82, Pass: 0.7}, {Name: "Tone", Score: 0.55, Pass: 0.7}, {Name: "Speed", Score: 0.91, Pass: 0.8}}
	ui.Column(c).Width(380).Children(func() { agent.EvalScore(c, "Evals", checks) })
}

func demoRateAnswer(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	rating := intState(root, "rating", 0)
	reason := intState(root, "reason", -1)
	root.Children(func() {
		box, send := agent.RateAnswer(c, rating, []string{"Wrong facts", "Too long", "Off topic"}, reason, agent.RateAnswerOptions{})
		if box.Changed() {
			log(fmt.Sprintf("rating %d, reason %d", *rating, *reason))
		}
		if send.Clicked() {
			log("sent")
		}
	})
}

func demoStreamCode(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	restart := intState(root, "restart", 0)
	code := "package main\n\nfunc main() {\n\tprintln(\"deck\")\n\tprintln(\"ship\")\n}\n"
	root.Children(func() {
		c.Key(*restart)
		agent.StreamCode(c, "go", code, false)
		if button.Outline(c, "Restart").Clicked() {
			*restart++
			log("restarted")
		}
	})
}

func demoSearchSteps(c *ui.Context, log func(string)) {
	steps := []agent.SearchStep{{Label: "Searching", Done: true}, {Label: "Reading", Count: 3}, {Label: "Writing"}}
	ui.Column(c).Width(320).Children(func() { agent.SearchSteps(c, "Search", steps) })
}

func demoImageGenerating(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	progress := numberState(root, "progress", 0.4)
	root.Children(func() {
		agent.ImageGenerating(c, "Harbor at dusk", *progress)
		if button.Outline(c, "10% more").Clicked() {
			*progress = min(*progress+0.1, 1)
			log(fmt.Sprintf("%.0f%%", *progress*100))
		}
	})
}

func demoAgentHandoff(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	base := timeState(root, "base", c)
	root.Children(func() { agent.AgentHandoff(c, "Planner", "Writer", "Took the deck", *base) })
}

func demoStopButton(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	running := boolState(root, "running", true)
	root.Children(func() {
		if agent.StopButton(c, *running, agent.StopButtonOptions{}).Changed() {
			*running = false
			log("stopped")
		}
	})
}

func demoPromptHistory(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	at := intState(root, "at", 0)
	entries := []string{"Draft the deck", "Review the brief", "Plan the sprint"}
	root.Children(func() {
		if agent.PromptHistory(c, entries, at, agent.PromptHistoryOptions{}).Changed() {
			log(fmt.Sprintf("at %d", *at+1))
		}
	})
}

func demoAgentWelcome(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	chosen := intState(root, "chosen", -1)
	prompts := []string{"Plan the launch", "Draft the deck", "Review the brief"}
	root.Children(func() {
		if agent.AgentWelcome(c, "What are we shipping?", prompts, []string{"Search", "Code", "Files"}, chosen).Changed() {
			log("picked " + prompts[*chosen])
		}
	})
}
