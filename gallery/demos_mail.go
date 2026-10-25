package gallery

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/mail"
	"github.com/egoist/mygo/ui"
)

var mailEntries = []Entry{
	{ID: 1161, Name: "Inbox", Category: "mail", Summary: "Mail rows with unread bold; cue Select.", Demo: demoInbox},
	{ID: 1162, Name: "MailRow", Category: "mail", Summary: "Sender, subject, preview, time; swipe archives; cue Swipe.", Demo: demoMailRow},
	{ID: 1163, Name: "ComposeWindow", Category: "mail", Summary: "To, subject, body, send; cue Send.", Demo: demoComposeWindow},
	{ID: 1164, Name: "Recipients", Category: "mail", Summary: "To, Cc, Bcc chips; cue Pop.", Demo: demoRecipients},
	{ID: 1165, Name: "Folders", Category: "mail", Summary: "Inbox, sent, drafts with counts; cue Select.", Demo: demoFolders},
	{ID: 1166, Name: "MailLabels", Category: "mail", Summary: "Colored labels; cue Toggle.", Demo: demoMailLabels},
	{ID: 1167, Name: "MailThread", Category: "mail", Summary: "Messages collapsed but the last; cue Toggle.", Demo: demoMailThread},
	{ID: 1168, Name: "MailAttachments", Category: "mail", Summary: "Attachment cards with size.", Demo: demoMailAttachments},
	{ID: 1169, Name: "Snooze", Category: "mail", Summary: "Later today, tomorrow, pick; cue Select.", Demo: demoSnooze},
	{ID: 1170, Name: "SendLater", Category: "mail", Summary: "Scheduled send choices; cue Select.", Demo: demoSendLater},
	{ID: 1171, Name: "Notification", Category: "mail", Summary: "Notification row with icon and time; cue Tap.", Demo: demoNotification},
	{ID: 1172, Name: "NotificationSettings", Category: "mail", Summary: "Channels by event with toggles; cue Toggle.", Demo: demoNotificationSettings},
	{ID: 1173, Name: "Signature", Category: "mail", Summary: "Signature editor with preview.", Demo: demoSignature},
	{ID: 1174, Name: "OutOfOffice", Category: "mail", Summary: "Toggle, dates and message; cue Toggle.", Demo: demoOutOfOffice},
	{ID: 1175, Name: "Contacts", Category: "mail", Summary: "Contacts list with letters; cue Select.", Demo: demoContacts},
	{ID: 1176, Name: "ContactCard", Category: "mail", Summary: "Avatar, name, fields, actions.", Demo: demoContactCard},
	{ID: 1177, Name: "CallScreen", Category: "mail", Summary: "Caller, timer, end; cue Close.", Demo: demoCallScreen},
	{ID: 1178, Name: "IncomingCall", Category: "mail", Summary: "Accept and decline; rings pulse; cue Notify.", Demo: demoIncomingCall},
	{ID: 1179, Name: "CallControls", Category: "mail", Summary: "Mute, video, share, end; cue Toggle.", Demo: demoCallControls},
	{ID: 1180, Name: "VideoGrid", Category: "mail", Summary: "Participant tiles; speaker ring.", Demo: demoVideoGrid},
	{ID: 1181, Name: "Participants", Category: "mail", Summary: "List with mute marks.", Demo: demoParticipants},
	{ID: 1182, Name: "ScreenShare", Category: "mail", Summary: "\"You are sharing\" bar with stop; cue Close.", Demo: demoScreenShare},
	{ID: 1183, Name: "Voicemail", Category: "mail", Summary: "Messages with play and transcript; cue Toggle.", Demo: demoVoicemail},
	{ID: 1184, Name: "MessageTemplates", Category: "mail", Summary: "Saved replies; cue Select.", Demo: demoMessageTemplates},
	{ID: 1185, Name: "Unsubscribe", Category: "mail", Summary: "Reason and confirm; cue Success.", Demo: demoUnsubscribe},
	{ID: 1186, Name: "Spam", Category: "mail", Summary: "Spam banner with not spam; cue Tap.", Demo: demoSpam},
	{ID: 1187, Name: "PriorityInbox", Category: "mail", Summary: "Important and other sections.", Demo: demoPriorityInbox},
	{ID: 1188, Name: "MailSearch", Category: "mail", Summary: "Search with from and has chips.", Demo: demoMailSearch},
	{ID: 1189, Name: "UndoSend", Category: "mail", Summary: "Sent toast with undo; cue Close.", Demo: demoUndoSend},
	{ID: 1190, Name: "ReadingPane", Category: "mail", Summary: "Selected mail beside the list.", Demo: demoReadingPane},
}

// demoMails is one product team's inbox during a launch week.
func demoMails() []mail.Mail {
	return []mail.Mail{
		{ID: "m1", Sender: "Maya Chen", Subject: "Atlas 2.0 launch checklist", Time: "9:42", Unread: true, Attached: true,
			Preview: "Three items left before Thursday: screenshots, notes, support brief.",
			Body:    "Three items left before Thursday: App Store screenshots, release notes and the support brief. I attached the checklist; tick yours off when done."},
		{ID: "m2", Sender: "Daniel Okafor", Subject: "Crash fix is in review", Time: "9:10", Unread: true,
			Preview: "The deep link guard is ready. Can you look before lunch?",
			Body:    "The deep link guard is ready, with a test for cold starts. Can you look before lunch? It blocks the beta build."},
		{ID: "m3", Sender: "Priya Nair", Subject: "Onboarding copy, final pass", Time: "Yesterday", Attached: true,
			Preview: "I tightened the second screen and cut the jargon.",
			Body:    "I tightened the second screen and cut the jargon. The new copy fits on small phones without wrapping."},
		{ID: "m4", Sender: "Lucas Martin", Subject: "Android widget mockups", Time: "Yesterday",
			Preview: "Two options attached. I lean toward the compact one.",
			Body:    "Two options attached. I lean toward the compact one: it shows the next trip and its weather in one line."},
		{ID: "m5", Sender: "Sofia Rossi", Subject: "Regression run: 214 passed", Time: "Mon",
			Preview: "All green on iOS 18 and Android 15.",
			Body:    "All green on iOS 18 and Android 15. Two flaky tests are quarantined until Friday."},
	}
}

func demoInbox(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	selected := textState(root, "selected", "")
	root.Children(func() {
		if mail.Inbox(c, demoMails(), selected, mail.InboxOptions{}).Changed() {
			log("selected " + *selected)
		}
	})
}

func demoMailRow(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	archived := boolState(root, "archived", false)
	root.Children(func() {
		if !*archived {
			c.Key("mail-row")
			if mail.MailRow(c, demoMails()[1], archived).Changed() {
				log("archived " + fmt.Sprint(*archived))
			}
			return
		}
		if button.Outline(c, "Undo archive").Clicked() {
			*archived = false
			log("restored")
		}
	})
}

func demoComposeWindow(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	value := ui.Local(root, "value", func() mail.Compose {
		return mail.Compose{To: "daniel.okafor@atlas.app", Subject: "Re: Crash fix is in review",
			Body: "Looks good to me. Ship it after the regression run, then tag the beta build."}
	})
	root.Children(func() {
		form, send := mail.ComposeWindow(c, value, mail.ComposeWindowOptions{})
		if send.Changed() {
			log("sent")
		} else if form.Changed() {
			log("edited " + fmt.Sprint(utf8.RuneCountInString(value.Body)))
		}
	})
}

func demoRecipients(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	to := ui.Local(root, "to", func() []string { return []string{"maya.chen@atlas.app", "daniel.okafor@atlas.app"} })
	cc := ui.Local(root, "cc", func() []string { return []string{"priya.nair@atlas.app"} })
	bcc := ui.Local(root, "bcc", func() []string { return []string{} })
	root.Children(func() {
		if mail.Recipients(c, to, cc, bcc, mail.RecipientsOptions{}).Changed() {
			log("recipients " + fmt.Sprint(len(*to)+len(*cc)+len(*bcc)))
		}
	})
}

func demoFolders(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320)
	selected := textState(root, "folder", "inbox")
	folders := []mail.Folder{
		{ID: "inbox", Name: "Inbox", Icon: "inbox", Count: 12},
		{ID: "starred", Name: "Starred", Icon: "star", Count: 3},
		{ID: "sent", Name: "Sent", Icon: "send"},
		{ID: "drafts", Name: "Drafts", Icon: "file-pen", Count: 2},
		{ID: "archive", Name: "Archive", Icon: "archive"},
		{ID: "trash", Name: "Trash", Icon: "trash-2"},
	}
	root.Children(func() {
		if mail.Folders(c, folders, selected).Changed() {
			log("folder " + *selected)
		}
	})
}

func demoMailLabels(c *ui.Context, log func(string)) {
	t := c.Theme()
	root := ui.Column(c).Width(400)
	selected := ui.Local(root, "labels", func() map[string]bool { return map[string]bool{"launch": true} })
	labels := []mail.MailLabel{
		{ID: "launch", Name: "Launch", Color: t.Accent}, {ID: "design", Name: "Design", Color: t.Success},
		{ID: "bugs", Name: "Bugs", Color: t.Danger}, {ID: "follow", Name: "Follow up", Color: t.Warning},
	}
	root.Children(func() {
		if mail.MailLabels(c, labels, selected).Changed() {
			on := 0
			for _, v := range *selected {
				if v {
					on++
				}
			}
			log("labels " + fmt.Sprint(on))
		}
	})
}

func demoMailThread(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	open := ui.Local(root, "open", func() map[string]bool { return map[string]bool{} })
	messages := []mail.ThreadMessage{
		{ID: "t1", Sender: "Daniel Okafor", Time: "09:02", Body: "Found the crash: the map router runs before the trip store loads on a cold start."},
		{ID: "t2", Sender: "Maya Chen", Time: "09:15", Body: "Nice catch. Can we guard it and add a cold-start test so it stays fixed?"},
		{ID: "t3", Sender: "Priya Nair", Time: "09:40", Body: "I checked the fix on my phone. Shared links open straight to the map now."},
	}
	root.Children(func() {
		if mail.MailThread(c, messages, open).Changed() {
			count := 0
			for _, shown := range *open {
				if shown {
					count++
				}
			}
			log("open " + fmt.Sprint(count))
		}
	})
}

func demoMailAttachments(c *ui.Context, log func(string)) {
	files := []mail.Attachment{
		{Name: "launch-checklist.pdf", Size: 248 << 10}, {Name: "store-screenshots.zip", Size: 18<<20 + 400<<10},
		{Name: "widget-mockup.png", Size: 1200 << 10}, {Name: "launch-budget.xlsx", Size: 42 << 10},
	}
	ui.Column(c).Width(440).Children(func() { mail.MailAttachments(c, files, mail.MailAttachmentsOptions{}) })
}

func demoSnooze(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	at := ui.Local(root, "at", func() time.Time { return c.Now() })
	root.Children(func() {
		if mail.Snooze(c, at, mail.SnoozeOptions{}).Changed() {
			log("snoozed " + at.Format("2006-01-02"))
		}
	})
}

func demoSendLater(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	at := ui.Local(root, "at", func() time.Time { return c.Now() })
	root.Children(func() {
		if mail.SendLater(c, at, mail.SendLaterOptions{}).Changed() {
			log("scheduled " + at.Format("2006-01-02"))
		}
	})
}

func demoNotification(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Gap(8)
	first, second := boolState(root, "first", false), boolState(root, "second", true)
	reply := mail.Notice{ID: "n1", Icon: "reply", Time: "2m ago", Text: "Daniel Okafor replied to Crash fix is in review",
		Detail: "The fix and its test are in review now."}
	shared := mail.Notice{ID: "n2", Icon: "file-text", Time: "1h ago", Text: "Sofia Rossi shared Regression run: 214 passed"}
	root.Children(func() {
		a := mail.Notification(c, reply, first).Changed()
		b := mail.Notification(c, shared, second).Changed()
		if a || b {
			opened := 0
			for _, v := range []bool{*first, *second} {
				if v {
					opened++
				}
			}
			log("opened " + fmt.Sprint(opened))
		}
	})
}

func demoNotificationSettings(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	settings := ui.Local(root, "settings", func() []mail.NotificationSetting {
		return []mail.NotificationSetting{{Event: "Replies", Email: true, Push: true}, {Event: "Mentions", Push: true},
			{Event: "Weekly digest", Email: true}, {Event: "Security alerts", Email: true, Push: true}}
	})
	root.Children(func() {
		if mail.NotificationSettings(c, settings, mail.NotificationSettingsOptions{}).Changed() {
			on := 0
			for _, s := range *settings {
				for _, v := range []bool{s.Email, s.Push} {
					if v {
						on++
					}
				}
			}
			log("channels " + fmt.Sprint(on))
		}
	})
}
