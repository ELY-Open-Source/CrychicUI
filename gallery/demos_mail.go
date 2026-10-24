package gallery

import (
	"fmt"
	"time"

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

func demoMail() mail.Mail {
	return mail.Mail{ID: "m1", Sender: "Ana", Subject: "Release", Preview: "Notes", Body: "The release is ready.", Time: "10:00", Unread: true}
}

func demoInbox(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := textState(root, "selected", "")
	root.Children(func() {
		if mail.Inbox(c, []mail.Mail{demoMail()}, selected).Changed() {
			log("selected " + *selected)
		}
	})
}

func demoMailRow(c *ui.Context, log func(string)) {
	root := ui.Column(c).FillWidth()
	archived := boolState(root, "archived", false)
	root.Children(func() {
		if mail.MailRow(c, demoMail(), archived).Changed() {
			log("archived " + fmt.Sprint(*archived))
		}
	})
}

func demoComposeWindow(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	value := ui.Local(root, "value", func() mail.Compose {
		return mail.Compose{To: "ana@example.com", Subject: "Release", Body: "Ready"}
	})
	root.Children(func() {
		form, send := mail.ComposeWindow(c, value, mail.ComposeWindowOptions{})
		if send.Changed() {
			log("sent")
		} else if form.Changed() {
			log("edited " + fmt.Sprint(len([]rune(value.Body))))
		}
	})
}

func demoRecipients(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	to := ui.Local(root, "to", func() []string { return []string{"ana@example.com", "ben@example.com"} })
	cc := ui.Local(root, "cc", func() []string { return []string{} })
	bcc := ui.Local(root, "bcc", func() []string { return []string{} })
	root.Children(func() {
		if mail.Recipients(c, to, cc, bcc, mail.RecipientsOptions{}).Changed() {
			log("recipients " + fmt.Sprint(len(*to)+len(*cc)+len(*bcc)))
		}
	})
}

func demoFolders(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := textState(root, "folder", "inbox")
	folders := []mail.Folder{{ID: "inbox", Name: "Inbox", Count: 2}, {ID: "sent", Name: "Sent"}}
	root.Children(func() {
		if mail.Folders(c, folders, selected).Changed() {
			log("folder " + *selected)
		}
	})
}

func demoMailLabels(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := ui.Local(root, "labels", func() map[string]bool { return map[string]bool{} })
	labels := []mail.MailLabel{{ID: "work", Name: "Work", Color: c.Theme().Accent}}
	root.Children(func() {
		if mail.MailLabels(c, labels, selected).Changed() {
			log("labels " + fmt.Sprint(len(*selected)))
		}
	})
}

func demoMailThread(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := ui.Local(root, "open", func() map[string]bool { return map[string]bool{"m2": true} })
	messages := []mail.ThreadMessage{{ID: "m1", Sender: "Ana", Body: "First", Time: "09:00"}, {ID: "m2", Sender: "Ben", Body: "Last", Time: "10:00"}}
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
	mail.MailAttachments(c, []mail.Attachment{{Name: "notes.pdf", Size: 2048}})
}

func demoSnooze(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	at := ui.Local(root, "at", func() time.Time { return c.Now() })
	root.Children(func() {
		if mail.Snooze(c, at, mail.SnoozeOptions{}).Changed() {
			log("snoozed " + at.Format("2006-01-02"))
		}
	})
}

func demoSendLater(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	at := ui.Local(root, "at", func() time.Time { return c.Now() })
	root.Children(func() {
		if mail.SendLater(c, at, mail.SendLaterOptions{}).Changed() {
			log("scheduled " + at.Format("2006-01-02"))
		}
	})
}

func demoNotification(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	opened := boolState(root, "opened", false)
	notice := mail.Notice{ID: "n1", Text: "New reply", Time: "now"}
	root.Children(func() {
		if mail.Notification(c, notice, opened).Changed() {
			log("opened " + fmt.Sprint(*opened))
		}
	})
}

func demoNotificationSettings(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	settings := ui.Local(root, "settings", func() []mail.NotificationSetting {
		return []mail.NotificationSetting{{Event: "Replies", Email: true}}
	})
	root.Children(func() {
		if mail.NotificationSettings(c, settings, mail.NotificationSettingsOptions{}).Changed() {
			log(fmt.Sprintf("email %v push %v", (*settings)[0].Email, (*settings)[0].Push))
		}
	})
}

func demoSignature(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	text := textState(root, "signature", "Ana")
	root.Children(func() {
		if mail.Signature(c, text, mail.SignatureOptions{}).Changed() {
			log("runes " + fmt.Sprint(len([]rune(*text))))
		}
	})
}

func demoOutOfOffice(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	value := ui.Local(root, "away", func() mail.AutoReply {
		return mail.AutoReply{Start: c.Now(), End: c.Now().AddDate(0, 0, 1), Message: "Away"}
	})
	root.Children(func() {
		if mail.OutOfOffice(c, value, mail.OutOfOfficeOptions{}).Changed() {
			log("away " + fmt.Sprint(value.Enabled))
		}
	})
}

func demoContacts(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := textState(root, "contact", "")
	people := []mail.Contact{{ID: "ana", Name: "Ana", Email: "ana@example.com"}}
	root.Children(func() {
		if mail.Contacts(c, people, selected).Changed() {
			log("contact " + *selected)
		}
	})
}

func demoContactCard(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	action := textState(root, "action", "")
	person := mail.Contact{ID: "ana", Name: "Ana", Email: "ana@example.com", Phone: "123"}
	root.Children(func() {
		if mail.ContactCard(c, person, action, mail.ContactCardOptions{}).Changed() {
			log("action " + *action)
		}
	})
}

func demoCallScreen(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	ended := boolState(root, "ended", false)
	started := ui.Local(root, "started", func() time.Time { return c.Now() })
	root.Children(func() {
		if mail.CallScreen(c, "Ana", *started, ended, mail.CallScreenOptions{}).Changed() {
			log("ended " + fmt.Sprint(*ended))
		}
	})
}

func demoIncomingCall(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	accepted, declined := boolState(root, "accepted", false), boolState(root, "declined", false)
	root.Children(func() {
		if mail.IncomingCall(c, "Ana", accepted, declined, mail.IncomingCallOptions{}).Changed() {
			log(fmt.Sprintf("accepted %v declined %v", *accepted, *declined))
		}
	})
}

func demoCallControls(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	state := ui.Local(root, "call", func() mail.CallState { return mail.CallState{Video: true} })
	root.Children(func() {
		if mail.CallControls(c, state, mail.CallControlsOptions{}).Changed() {
			log(fmt.Sprintf("muted %v ended %v", state.Muted, state.Ended))
		}
	})
}

func demoVideoGrid(c *ui.Context, log func(string)) {
	mail.VideoGrid(c, []mail.Participant{{ID: "ana", Name: "Ana", Speaking: true}})
}

func demoParticipants(c *ui.Context, log func(string)) {
	mail.Participants(c, []mail.Participant{{ID: "ana", Name: "Ana", Muted: true}})
}

func demoScreenShare(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	sharing := boolState(root, "sharing", true)
	root.Children(func() {
		if mail.ScreenShare(c, sharing, mail.ScreenShareOptions{}).Changed() {
			log("sharing " + fmt.Sprint(*sharing))
		}
	})
}

func demoVoicemail(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	playing := intState(root, "playing", -1)
	messages := []mail.VoiceMessage{{ID: "v1", Sender: "Ana", Time: "1m", Transcript: "Call me"}}
	root.Children(func() {
		if mail.Voicemail(c, messages, playing, mail.VoicemailOptions{}).Changed() {
			log("playing " + fmt.Sprint(*playing))
		}
	})
}

func demoMessageTemplates(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := intState(root, "reply", -1)
	replies := []mail.SavedReply{{ID: "r1", Title: "Thanks", Body: "Thank you"}}
	root.Children(func() {
		if mail.MessageTemplates(c, replies, selected).Changed() {
			log("reply " + fmt.Sprint(*selected))
		}
	})
}

func demoUnsubscribe(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	reason := textState(root, "reason", "")
	confirmed := boolState(root, "confirmed", false)
	root.Children(func() {
		form, confirm := mail.Unsubscribe(c, []string{"Too many"}, reason, confirmed, mail.UnsubscribeOptions{})
		if confirm.Changed() {
			log("confirmed")
		} else if form.Changed() {
			log("reason " + *reason)
		}
	})
}

func demoSpam(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	spam := boolState(root, "spam", true)
	root.Children(func() {
		if mail.Spam(c, spam, mail.SpamOptions{}).Changed() {
			log("spam " + fmt.Sprint(*spam))
		}
	})
}

func demoPriorityInbox(c *ui.Context, log func(string)) {
	other := demoMail()
	other.ID, other.Unread = "m2", false
	mail.PriorityInbox(c, []mail.Mail{demoMail()}, []mail.Mail{other}, mail.PriorityInboxOptions{})
}

func demoMailSearch(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	query := ui.Local(root, "query", func() mail.MailQuery { return mail.MailQuery{} })
	root.Children(func() {
		if mail.MailSearch(c, query, mail.MailSearchOptions{}).Changed() {
			log(fmt.Sprintf("runes %d attachments %v", len([]rune(query.Text)), query.Attachments))
		}
	})
}

func demoUndoSend(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	visible := boolState(root, "visible", true)
	root.Children(func() {
		if mail.UndoSend(c, visible, mail.UndoSendOptions{}).Changed() {
			log("visible " + fmt.Sprint(*visible))
		}
	})
}

func demoReadingPane(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	selected := textState(root, "reading", "m1")
	second := demoMail()
	second.ID, second.Subject, second.Body = "m2", "Second", "Second body"
	root.Children(func() {
		if mail.ReadingPane(c, []mail.Mail{demoMail(), second}, selected).Changed() {
			log("reading " + *selected)
		}
	})
}
