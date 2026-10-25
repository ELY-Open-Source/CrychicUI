package gallery

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/mail"
	"github.com/egoist/mygo/ui"
)

func demoSignature(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	text := textState(root, "signature", "Maya Chen\nProduct lead, Atlas\nmaya.chen@atlas.app")
	root.Children(func() {
		if mail.Signature(c, text, mail.SignatureOptions{}).Changed() {
			log("runes " + fmt.Sprint(utf8.RuneCountInString(*text)))
		}
	})
}

func demoOutOfOffice(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	value := ui.Local(root, "away", func() mail.AutoReply {
		return mail.AutoReply{Enabled: true, Start: c.Now(), End: c.Now().AddDate(0, 0, 7),
			Message: "I'm away with little access to mail. For launch questions, write to Daniel Okafor."}
	})
	root.Children(func() {
		if mail.OutOfOffice(c, value, mail.OutOfOfficeOptions{}).Changed() {
			log("away " + fmt.Sprint(value.Enabled))
		}
	})
}

// demoPeople is the launch team's address book.
func demoPeople() []mail.Contact {
	return []mail.Contact{
		{ID: "maya", Name: "Maya Chen", Email: "maya.chen@atlas.app", Role: "Product lead"},
		{ID: "daniel", Name: "Daniel Okafor", Email: "daniel.okafor@atlas.app", Role: "Backend engineer"},
		{ID: "priya", Name: "Priya Nair", Email: "priya.nair@atlas.app", Phone: "+1 415 555 0142", Role: "Product designer"},
		{ID: "lucas", Name: "Lucas Martin", Email: "lucas.martin@atlas.app", Role: "Android engineer"},
		{ID: "sofia", Name: "Sofia Rossi", Email: "sofia.rossi@atlas.app", Role: "QA lead"},
		{ID: "david", Name: "David Kim", Email: "david.kim@atlas.app", Role: "Support"},
	}
}

func demoContacts(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	selected := textState(root, "contact", "maya")
	root.Children(func() {
		if mail.Contacts(c, demoPeople(), selected, mail.ContactsOptions{}).Changed() {
			log("contact " + *selected)
		}
	})
}

func demoContactCard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(340)
	action := textState(root, "action", "")
	root.Children(func() {
		if mail.ContactCard(c, demoPeople()[2], action, mail.ContactCardOptions{}).Changed() {
			log("action " + *action)
		}
	})
}

func demoCallScreen(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	ended := boolState(root, "ended", false)
	started := ui.Local(root, "started", func() time.Time { return c.Now().Add(-4*time.Minute - 12*time.Second) })
	root.Children(func() {
		if mail.CallScreen(c, "Daniel Okafor", *started, ended, mail.CallScreenOptions{}).Changed() {
			log("ended " + fmt.Sprint(*ended))
		}
	})
}

func demoIncomingCall(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	accepted, declined := boolState(root, "accepted", false), boolState(root, "declined", false)
	root.Children(func() {
		if mail.IncomingCall(c, "Priya Nair", accepted, declined, mail.IncomingCallOptions{}).Changed() {
			log(fmt.Sprintf("accepted %v declined %v", *accepted, *declined))
		}
	})
}

func demoCallControls(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	state := ui.Local(root, "call", func() mail.CallState { return mail.CallState{Video: true} })
	root.Children(func() {
		if mail.CallControls(c, state, mail.CallControlsOptions{}).Changed() {
			log(fmt.Sprintf("muted %v ended %v", state.Muted, state.Ended))
		}
	})
}

// demoCall is the launch stand-up's call.
func demoCall() []mail.Participant {
	return []mail.Participant{
		{ID: "maya", Name: "Maya Chen", Speaking: true}, {ID: "daniel", Name: "Daniel Okafor"},
		{ID: "priya", Name: "Priya Nair"}, {ID: "lucas", Name: "Lucas Martin", Muted: true},
	}
}

func demoVideoGrid(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() { mail.VideoGrid(c, demoCall()) })
}

func demoParticipants(c *ui.Context, log func(string)) {
	people := append(demoCall(), mail.Participant{ID: "sofia", Name: "Sofia Rossi", Muted: true})
	ui.Column(c).Width(380).Children(func() { mail.Participants(c, people, mail.ParticipantsOptions{}) })
}

func demoScreenShare(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	sharing := boolState(root, "sharing", true)
	root.Children(func() {
		if mail.ScreenShare(c, sharing, mail.ScreenShareOptions{Detail: "Atlas 2.0 roadmap · Keynote"}).Changed() {
			log("sharing " + fmt.Sprint(*sharing))
		}
	})
}

func demoVoicemail(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	playing := intState(root, "playing", -1)
	messages := []mail.VoiceMessage{
		{ID: "v1", Sender: "Daniel Okafor", Time: "10:24", Length: 38 * time.Second,
			Transcript: "Hey, the beta build is up. Give it a spin on your phone when you get a minute."},
		{ID: "v2", Sender: "David Kim", Time: "Yesterday", Length: 1*time.Minute + 12*time.Second,
			Transcript: "Support is ready for launch. We have answers for the top ten offline questions."},
		{ID: "v3", Sender: "Sofia Rossi", Time: "Mon", Length: 21 * time.Second,
			Transcript: "Quick one: two tests are flaky on Android 15. I quarantined them until Friday."},
	}
	root.Children(func() {
		if mail.Voicemail(c, messages, playing, mail.VoicemailOptions{}).Changed() {
			log("playing " + fmt.Sprint(*playing))
		}
	})
}

func demoMessageTemplates(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	selected := intState(root, "reply", -1)
	replies := []mail.SavedReply{
		{ID: "thanks", Title: "Thanks, on it", Body: "Thanks for the heads-up. I'm on it and will reply by end of day."},
		{ID: "later", Title: "Next week", Body: "Good idea. Let's pick this up after the launch next week."},
		{ID: "review", Title: "Ready for review", Body: "The change is ready for review. The test plan is in the description."},
		{ID: "away", Title: "Out until Monday", Body: "I'm out until Monday. For urgent launch issues, ask Daniel Okafor."},
	}
	root.Children(func() {
		if mail.MessageTemplates(c, replies, selected, mail.MessageTemplatesOptions{}).Changed() {
			log("reply " + fmt.Sprint(*selected))
		}
	})
}

func demoUnsubscribe(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	reason := textState(root, "reason", "")
	confirmed := boolState(root, "confirmed", false)
	reasons := []string{"I get too many emails", "The content isn't relevant to me", "I never signed up", "Other"}
	root.Children(func() {
		form, confirm := mail.Unsubscribe(c, reasons, reason, confirmed, mail.UnsubscribeOptions{})
		if confirm.Changed() {
			log("confirmed")
		} else if form.Changed() {
			log("reason picked")
		}
	})
}

func demoSpam(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	spam := boolState(root, "spam", true)
	root.Children(func() {
		if mail.Spam(c, spam, mail.SpamOptions{}).Changed() {
			log("spam " + fmt.Sprint(*spam))
		}
	})
}

func demoPriorityInbox(c *ui.Context, log func(string)) {
	m := demoMails()
	ui.Column(c).Width(440).Children(func() {
		mail.PriorityInbox(c, m[:2], m[3:], mail.PriorityInboxOptions{})
	})
}

func demoMailSearch(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	query := ui.Local(root, "query", func() mail.MailQuery { return mail.MailQuery{Text: "launch", From: "maya.chen@atlas.app"} })
	root.Children(func() {
		if mail.MailSearch(c, query, mail.MailSearchOptions{}).Changed() {
			log(fmt.Sprintf("runes %d attachments %v", utf8.RuneCountInString(query.Text), query.Attachments))
		}
	})
}

func demoUndoSend(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	visible := boolState(root, "visible", true)
	root.Children(func() {
		if mail.UndoSend(c, visible, mail.UndoSendOptions{}).Changed() {
			log("visible " + fmt.Sprint(*visible))
		}
	})
}

func demoReadingPane(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	selected := textState(root, "reading", "m1")
	root.Children(func() {
		if mail.ReadingPane(c, demoMails()[:4], selected).Changed() {
			log("reading " + *selected)
		}
	})
}
