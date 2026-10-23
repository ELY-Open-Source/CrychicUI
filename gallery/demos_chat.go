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

var chatEntries = []Entry{
	{816, "Bubble", "chat", "A message bubble; mine sits right in accent.", demoChatBubble},
	{817, "MessageList", "chat", "Bubbles grouped by sender; new ones slide up.", demoMessageList},
	{818, "Composer", "chat", "A text area with attach and send; Enter sends.", demoComposer},
	{819, "Typing", "chat", "Three dots bouncing in a bubble.", demoTyping},
	{820, "Receipt", "chat", "Sent, delivered and read checks cross-fade.", demoReceipt},
	{821, "Reactions", "chat", "Emoji counts; a tap toggles yours.", demoReactions},
	{822, "ReactionPicker", "chat", "A row of emoji that scales in.", demoReactionPicker},
	{823, "ReplyQuote", "chat", "A quoted message; a click jumps to it.", demoReplyQuote},
	{824, "Thread", "chat", "A reply count opening a side thread.", demoThread},
	{825, "Attachment", "chat", "A file chip with its size.", demoAttachment},
	{826, "ImageMessage", "chat", "An image bubble that opens large.", demoImageMessage},
	{827, "VoiceMessage", "chat", "Play, a waveform and the time left.", demoVoiceMessage},
	{828, "LinkPreview", "chat", "A title, host and image card.", demoLinkPreview},
	{829, "SystemMessage", "chat", "A centered muted event line.", demoSystemMessage},
	{830, "DayDivider", "chat", "A day label between two lines.", demoDayDivider},
	{831, "ConversationList", "chat", "Avatar, name, last line, time and unread.", demoConversationList},
	{832, "ChatSearch", "chat", "Search over messages; a hit jumps.", demoChatSearch},
	{833, "Pinned", "chat", "A pinned message bar with jump.", demoPinned},
	{834, "MessageMenu", "chat", "Reply, copy, edit, delete on hover.", demoMessageMenu},
	{835, "EditedMessage", "chat", "An edited mark whose tip shows the time.", demoEditedMessage},
	{836, "DeletedMessage", "chat", "A deleted placeholder that fades in.", demoDeletedMessage},
	{837, "Mention", "chat", "A highlighted @name in text.", demoMention},
	{838, "QuickReplies", "chat", "Reply chips; the tapped one sends.", demoQuickReplies},
}

// chatSample is a short talk ending at base.
func chatSample(base time.Time) []chat.Message {
	return []chat.Message{
		{Sender: "Ana", Text: "Morning! Did the build go through?", At: base.Add(-12 * time.Minute)},
		{Sender: "Ana", Text: "Lunch at noon?", At: base.Add(-11 * time.Minute)},
		{Sender: "me", Text: "Green on every runner. Noon works.", Mine: true, At: base.Add(-9 * time.Minute)},
		{Sender: "Kai", Text: "Can you review the deck before then?", At: base.Add(-3 * time.Minute)},
	}
}

func timeState(root ui.Element, key string, c *ui.Context) *time.Time {
	return ui.Local(root, key, func() time.Time { return c.Now() })
}

func demoChatBubble(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Gap(12)
	mine := boolState(root, "mine", false)
	base := timeState(root, "base", c)
	root.Children(func() {
		chat.Bubble(c, chat.Message{Sender: "Ana", Text: "The harbor shots came out well.", Mine: *mine, At: *base}, chat.BubbleOptions{ShowTime: true})
		if button.Outline(c, "Flip side").Clicked() {
			*mine = !*mine
			log(fmt.Sprintf("mine %v", *mine))
		}
	})
}

func demoMessageList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(12)
	base := timeState(root, "base", c)
	list := ui.Local(root, "list", func() []chat.Message { return chatSample(*base) })
	root.Children(func() {
		chat.MessageList(c, "Chat", *list, chat.MessageListOptions{Height: 320})
		ui.Row(c).Gap(8).Children(func() {
			if button.Outline(c, "Receive").Clicked() {
				*list = append(*list, chat.Message{Sender: "Ana", Text: fmt.Sprintf("Reply %d", len(*list)), At: c.Now()})
				sound.Play(sound.Notify)
				log(fmt.Sprintf("%d messages", len(*list)))
			}
			if button.Outline(c, "Send").Clicked() {
				*list = append(*list, chat.Message{Sender: "me", Text: fmt.Sprintf("Note %d", len(*list)), Mine: true, At: c.Now()})
				log(fmt.Sprintf("%d messages", len(*list)))
			}
		})
	})
}

func demoComposer(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	draft := ui.Local(root, "draft", func() string { return "" })
	root.Children(func() {
		root, _, attach := chat.Composer(c, draft, func(s string) { log(fmt.Sprintf("sent %d runes", utf8.RuneCountInString(s))) }, chat.ComposerOptions{})
		if root.Changed() {
			log(fmt.Sprintf("draft %d chars", len(*draft)))
		}
		if attach.Clicked() {
			log("attach")
		}
	})
}

func demoTyping(c *ui.Context, log func(string)) {
	ui.Column(c).Children(func() { chat.Typing(c, "Ana is typing") })
}

func demoReceipt(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	state := ui.Local(root, "state", func() chat.Delivery { return chat.Sent })
	root.Children(func() {
		chat.Receipt(c, *state, chat.ReceiptOptions{})
		if button.Outline(c, "Advance").Clicked() {
			*state = (*state + 1) % 3
			log(fmt.Sprintf("state %d", *state))
		}
	})
}

func demoReactions(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	items := ui.Local(root, "items", func() []chat.Reaction {
		return []chat.Reaction{{Emoji: "👍", Count: 3}, {Emoji: "🎉", Count: 1, Mine: true}, {Emoji: "👀", Count: 0}}
	})
	root.Children(func() {
		if chat.Reactions(c, "Reactions", items).Changed() {
			log(fmt.Sprintf("counts %d %d %d", (*items)[0].Count, (*items)[1].Count, (*items)[2].Count))
		}
	})
}

func demoReactionPicker(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	chosen := ui.Local(root, "chosen", func() string { return "" })
	root.Children(func() {
		if chat.ReactionPicker(c, "React", []string{"👍", "❤️", "😂", "😮", "😢"}, chosen).Changed() {
			log("picked " + *chosen)
		}
	})
}

func demoReplyQuote(c *ui.Context, log func(string)) {
	ui.Column(c).Width(320).Children(func() {
		if chat.ReplyQuote(c, "Ana", "Lunch at noon?", chat.ReplyQuoteOptions{}).Clicked() {
			log("jumped to the quote")
		}
	})
}

func demoThread(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := boolState(root, "open", false)
	root.Children(func() {
		if chat.Thread(c, 3, open, chat.ThreadOptions{}).Changed() {
			log(fmt.Sprintf("open %v", *open))
		}
	})
}

func demoAttachment(c *ui.Context, log func(string)) {
	ui.Column(c).Children(func() {
		if chat.Attachment(c, "report.pdf", 1234567).Clicked() {
			log("opened report.pdf")
		}
	})
}

func demoImageMessage(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := boolState(root, "open", false)
	root.Children(func() {
		if chat.ImageMessage(c, "Harbor photo", photos()[1], open, chat.ImageMessageOptions{}).Changed() {
			log(fmt.Sprintf("open %v", *open))
		}
	})
}

func demoVoiceMessage(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	playing := boolState(root, "playing", false)
	levels := []float64{0.2, 0.5, 0.9, 0.6, 0.3, 0.8, 1, 0.4, 0.7, 0.2, 0.5, 0.3}
	root.Children(func() {
		if chat.VoiceMessage(c, "Voice note", levels, 12*time.Second, playing, chat.VoiceMessageOptions{}).Changed() {
			log(fmt.Sprintf("playing %v", *playing))
		}
	})
}

func demoLinkPreview(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	withImage := boolState(root, "image", true)
	root.Children(func() {
		var img *ui.Bitmap
		if *withImage {
			img = photos()[0]
		}
		chat.LinkPreview(c, "MyGo: native UI for Go", "mygo.egoist.dev", img)
		if button.Outline(c, "Toggle image").Clicked() {
			*withImage = !*withImage
			log(fmt.Sprintf("image %v", *withImage))
		}
	})
}

func demoSystemMessage(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Children(func() { chat.SystemMessage(c, "Ana added Kai to the chat") })
}

func demoDayDivider(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Children(func() { chat.DayDivider(c, "Today") })
}

func demoConversationList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	selected := intState(root, "selected", 0)
	base := timeState(root, "base", c)
	p := photos()
	items := []chat.Conversation{
		{Name: "Ana Reyes", Last: "Lunch at noon?", At: base.Add(-11 * time.Minute), Unread: 2, Avatar: p[0]},
		{Name: "Kai Mori", Last: "Can you review the deck?", At: base.Add(-3 * time.Minute), Unread: 1},
		{Name: "Design team", Last: "Standup moved to 10", At: base.Add(-2 * time.Hour), Avatar: p[2]},
	}
	root.Children(func() {
		if chat.ConversationList(c, "Chats", items, selected).Changed() {
			log("selected " + items[*selected].Name)
		}
	})
}

func demoChatSearch(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	query := ui.Local(root, "query", func() string { return "" })
	jump := intState(root, "jump", -1)
	base := timeState(root, "base", c)
	list := chatSample(*base)
	root.Children(func() {
		box, field := chat.ChatSearch(c, "Search", list, query, jump, chat.ChatSearchOptions{})
		if field.Changed() {
			log(fmt.Sprintf("query %d chars", len(*query)))
		}
		if box.Changed() {
			log(fmt.Sprintf("jump %d", *jump))
		}
	})
}

func demoPinned(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Children(func() {
		if chat.Pinned(c, "Ana", "Standup moved to 10 on Fridays", chat.PinnedOptions{}).Clicked() {
			log("jumped to the pin")
		}
	})
}

func demoMessageMenu(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	chosen := intState(root, "chosen", -1)
	names := []string{"reply", "copy", "edit", "delete"}
	root.Children(func() {
		box := chat.MessageMenu(c, "Message", chosen, chat.MessageMenuOptions{}, func() {
			chat.Bubble(c, chat.Message{Sender: "Kai", Text: "Can you review the deck before then?", At: c.Now()}, chat.BubbleOptions{})
		})
		if box.Changed() {
			log("action " + names[*chosen])
		}
	})
}

func demoEditedMessage(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(4)
	base := timeState(root, "base", c)
	root.Children(func() {
		chat.Bubble(c, chat.Message{Sender: "me", Text: "Noon works, see you there.", Mine: true, At: *base}, chat.BubbleOptions{})
		chat.EditedMessage(c, *base, chat.EditedMessageOptions{}).AlignSelf(ui.End)
	})
}

func demoDeletedMessage(c *ui.Context, log func(string)) {
	ui.Column(c).Children(func() { chat.DeletedMessage(c, chat.DeletedMessageOptions{}) })
}

func demoMention(c *ui.Context, log func(string)) {
	ui.Row(c).Gap(4).AlignItems(ui.Center).Children(func() {
		ui.Text(c, "Ping")
		if chat.Mention(c, "ana").Clicked() {
			log("opened @ana")
		}
		ui.Text(c, "when the deck is ready.")
	})
}

func demoQuickReplies(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	chosen := intState(root, "chosen", -1)
	replies := []string{"Sounds good", "Running late", "Call me"}
	root.Children(func() {
		if chat.QuickReplies(c, "Quick replies", replies, chosen).Changed() {
			log("reply: " + replies[*chosen])
		}
	})
}
