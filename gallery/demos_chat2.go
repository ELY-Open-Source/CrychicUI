package gallery

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/chat"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

var chatEntries2 = []Entry{
	{839, "PollMessage", "chat", "Options with vote bars; bars grow.", demoPollMessage},
	{840, "LocationMessage", "chat", "A map tile with a pin and an address.", demoLocationMessage},
	{841, "ContactMessage", "chat", "A contact card in a bubble with add.", demoContactMessage},
	{842, "Sticker", "chat", "A large image without a bubble; pops in.", demoSticker},
	{843, "EmojiBurst", "chat", "A big emoji message that bounces.", demoEmojiBurst},
	{844, "ChannelList", "chat", "Channels with unread dots.", demoChannelList},
	{845, "Members", "chat", "Online and offline groups.", demoMembers},
	{846, "CallBanner", "chat", "A call in progress bar with join.", demoCallBanner},
	{847, "JumpToLatest", "chat", "A pill with a count after scrolling up.", demoJumpToLatest},
	{848, "MessageStatus", "chat", "A sending clock and a failed retry.", demoMessageStatus},
	{849, "Translate", "chat", "Toggle original and translation; cross-fades.", demoTranslate},
	{850, "CodeMessage", "chat", "A code block bubble with copy.", demoCodeMessage},
	{851, "MarkdownMessage", "chat", "A bubble rendered with ui/markdown.", demoMarkdownMessage},
	{852, "Forwarded", "chat", "A forwarded-from header in a bubble.", demoForwarded},
	{853, "ChatEmpty", "chat", "An empty chat with a wave.", demoChatEmpty},
	{854, "OnlineBar", "chat", "Avatars of who is online; new ones pop.", demoOnlineBar},
	{855, "HoldToRecord", "chat", "Hold to record, slide away to cancel.", demoHoldToRecord},
	{856, "MessageTime", "chat", "A time that shows on hover.", demoMessageTime},
	{857, "Handoff", "chat", "An agent-joined event with an avatar.", demoHandoff},
	{858, "ChatHeader", "chat", "Avatar, name, status and call buttons.", demoChatHeader},
	{859, "UnreadDivider", "chat", "A new-messages line that fades.", demoUnreadDivider},
	{860, "ChatWindow", "chat", "A header, list and composer together.", demoChatWindow},
}

func demoPollMessage(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	votes := ui.Local(root, "votes", func() []int { return []int{3, 2, 1} })
	chosen := intState(root, "chosen", -1)
	root.Children(func() {
		if chat.PollMessage(c, "Lunch on Friday?", []string{"Pizza", "Sushi", "Salad"}, votes, chosen, chat.PollMessageOptions{}).Changed() {
			log(fmt.Sprintf("votes %v", *votes))
		}
	})
}

func demoLocationMessage(c *ui.Context, log func(string)) {
	ui.Column(c).Children(func() { chat.LocationMessage(c, "471 Harbor Road, Duluth") })
}

func demoContactMessage(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	added := boolState(root, "added", false)
	root.Children(func() {
		if chat.ContactMessage(c, "Ana Reyes", photos()[0], added, chat.ContactMessageOptions{}).Changed() {
			log(fmt.Sprintf("added %v", *added))
		}
	})
}

func demoSticker(c *ui.Context, log func(string)) {
	ui.Column(c).Children(func() { chat.Sticker(c, "Sunset sticker", photos()[0], chat.StickerOptions{}) })
}

func demoEmojiBurst(c *ui.Context, log func(string)) {
	ui.Column(c).Children(func() { chat.EmojiBurst(c, "🎉") })
}

func demoChannelList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(280)
	selected := intState(root, "selected", 0)
	channels := []chat.Channel{{Name: "general", Unread: 2}, {Name: "random"}, {Name: "design", Unread: 5}}
	root.Children(func() {
		if chat.ChannelList(c, "Channels", channels, selected).Changed() {
			log("selected " + channels[*selected].Name)
		}
	})
}

func demoMembers(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(280)
	chosen := intState(root, "chosen", -1)
	names := []string{"Ana Reyes", "Kai Mori", "Bo Tan"}
	root.Children(func() {
		if chat.Members(c, "Members", names[:2], names[2:], chosen, chat.MembersOptions{}).Changed() {
			log("picked " + names[*chosen])
		}
	})
}

func demoCallBanner(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	joined := boolState(root, "joined", false)
	root.Children(func() {
		if chat.CallBanner(c, "Call in progress · 12:04", joined, chat.CallBannerOptions{}).Changed() {
			log(fmt.Sprintf("joined %v", *joined))
		}
	})
}

func demoJumpToLatest(c *ui.Context, log func(string)) {
	ui.Column(c).Children(func() {
		if chat.JumpToLatest(c, 5, chat.JumpToLatestOptions{}).Changed() {
			log("jumped to latest")
		}
	})
}

func demoMessageStatus(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	state := ui.Local(root, "state", func() chat.SendState { return chat.Failed })
	root.Children(func() {
		box := chat.MessageStatus(c, *state, func() {
			*state = chat.Sending
			log("retried")
		}, chat.MessageStatusOptions{})
		if box.Changed() {
			log(fmt.Sprintf("state %d", *state))
		}
		if button.Outline(c, "Fail").Clicked() {
			*state = chat.Failed
			sound.Play(sound.Error)
			log("failed")
		}
	})
}

func demoTranslate(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	shown := boolState(root, "shown", false)
	root.Children(func() {
		if chat.Translate(c, "The harbor shots came out well.", "Les photos du port sont réussies.", shown, chat.TranslateOptions{}).Changed() {
			log(fmt.Sprintf("translated %v", *shown))
		}
	})
}

func demoCodeMessage(c *ui.Context, log func(string)) {
	code := "fmt.Println(\"harbor\")"
	ui.Column(c).Children(func() {
		if _, cp := chat.CodeMessage(c, "go", code, chat.CodeMessageOptions{}); cp.Clicked() {
			log(fmt.Sprintf("copied %d runes", utf8.RuneCountInString(code)))
		}
	})
}

func demoMarkdownMessage(c *ui.Context, log func(string)) {
	m := chat.Message{Sender: "Ana", Text: "Ship it **today**:\n\n- run the tests\n- tag `v0.6`", At: time.Now()}
	ui.Column(c).Width(320).Children(func() { chat.MarkdownMessage(c, m) })
}

func demoForwarded(c *ui.Context, log func(string)) {
	m := chat.Message{Sender: "Kai", Text: "Standup moved to 10 on Fridays", At: time.Now()}
	ui.Column(c).Width(320).Children(func() { chat.Forwarded(c, "Kai Mori", m, chat.ForwardedOptions{}) })
}

func demoChatEmpty(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Children(func() {
		if _, act := chat.ChatEmpty(c, "No messages yet", "Say something nice to Ana.", "Say hi"); act.Clicked() {
			log("waved hello")
		}
	})
}

func demoOnlineBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	names := ui.Local(root, "names", func() []string { return []string{"Ana Reyes"} })
	root.Children(func() {
		chat.OnlineBar(c, "Online now", *names)
		if button.Outline(c, "Kai joins").Clicked() {
			*names = append(*names, "Kai Mori")
			log(fmt.Sprintf("%d online", len(*names)))
		}
	})
}

func demoHoldToRecord(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	recording := boolState(root, "recording", false)
	root.Children(func() {
		box := chat.HoldToRecord(c, "Voice note", recording, func() { log("cancelled") }, chat.HoldToRecordOptions{})
		if box.Changed() {
			log(fmt.Sprintf("recording %v", *recording))
		}
	})
}

func demoMessageTime(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320)
	base := timeState(root, "base", c)
	root.Children(func() {
		chat.MessageTime(c, *base, func() {
			chat.Bubble(c, chat.Message{Sender: "Ana", Text: "Lunch at noon?", At: *base}, chat.BubbleOptions{})
		})
	})
}

func demoHandoff(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	base := timeState(root, "base", c)
	root.Children(func() { chat.Handoff(c, "Ana Reyes", photos()[0], *base, chat.HandoffOptions{}) })
}

func demoChatHeader(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Children(func() {
		_, voice, video := chat.ChatHeader(c, "Ana Reyes", "Online", photos()[0], chat.ChatHeaderOptions{})
		if voice.Clicked() {
			log("voice call")
		}
		if video.Clicked() {
			log("video call")
		}
	})
}

func demoUnreadDivider(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Children(func() { chat.UnreadDivider(c, "New messages") })
}

func demoChatWindow(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	base := timeState(root, "base", c)
	list := chatSample(*base)
	root.Children(func() {
		_, _, voice, video, attach := chat.ChatWindow(c, "Ana Reyes", "Online", list, func(s string) {
			log(fmt.Sprintf("sent %d runes", utf8.RuneCountInString(s)))
		}, chat.ChatWindowOptions{Height: 300, Avatar: photos()[0]})
		if voice.Clicked() {
			log("voice call")
		}
		if video.Clicked() {
			log("video call")
		}
		if attach.Clicked() {
			log("attach")
		}
	})
}
