package gallery

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/overlay"
	"github.com/egoist/mygo/ui"
)

var overlayEntries2 = []Entry{
	{634, "EmojiPicker", "overlay", "Emoji by group, with search.", demoEmojiPicker},
	{635, "DropOverlay", "overlay", "A veil while files drag over.", demoDropOverlay},
	{636, "BusyVeil", "overlay", "Blocks content with a message.", demoBusyVeil},
	{637, "LockOverlay", "overlay", "A locked panel over blurred content.", demoLockOverlay},
	{638, "Paywall", "overlay", "An upgrade card over blurred content.", demoPaywall},
	{639, "ShareSheet", "overlay", "Share targets and a link to copy.", demoShareSheet},
	{640, "ActionSheet", "overlay", "Bottom actions with Cancel.", demoActionSheet},
	{641, "MiniPlayer", "overlay", "A floating player bar.", demoMiniPlayer},
	{642, "FloatingPanel", "overlay", "A panel dragged anywhere.", demoFloatingPanel},
	{643, "PiPBox", "overlay", "A box thrown to a corner.", demoPiPBox},
	{644, "CalloutPopover", "overlay", "A popover with a toned edge.", demoCalloutPopover},
	{645, "InfoPopover", "overlay", "An info button with text.", demoInfoPopover},
	{646, "ConfirmPopover", "overlay", "Are you sure? by its anchor.", demoConfirmPopover},
	{647, "SelectMenu", "overlay", "A searchable menu of choices.", demoSelectMenu},
	{648, "MentionMenu", "overlay", "People matching @ text.", demoMentionMenu},
	{649, "SlashMenu", "overlay", "Commands matching / text.", demoSlashMenu},
	{650, "ShortcutSheet", "overlay", "Shortcuts listed by group.", demoShortcutSheet},
	{651, "StepDialog", "overlay", "A dialog paging through steps.", demoStepDialog},
	{652, "RenameDialog", "overlay", "A name field, stem selected.", demoRenameDialog},
	{653, "SessionTimeout", "overlay", "A countdown with Stay.", demoSessionTimeout},
	{654, "FeedbackWidget", "overlay", "A corner button with a form.", demoFeedbackWidget},
	{655, "SheetStack", "overlay", "Sheets stacked, lower ones back.", demoSheetStack},
}

var demoEmoji = []overlay.EmojiGroup{
	{Title: "Smileys", Emoji: []overlay.Emoji{{Char: "😀", Name: "grinning"}, {Char: "😂", Name: "joy"}, {Char: "😍", Name: "heart eyes"}, {Char: "🤔", Name: "thinking"}}},
	{Title: "Animals", Emoji: []overlay.Emoji{{Char: "🐶", Name: "dog"}, {Char: "🐱", Name: "cat"}, {Char: "🦊", Name: "fox"}, {Char: "🐼", Name: "panda"}}},
	{Title: "Food", Emoji: []overlay.Emoji{{Char: "🍎", Name: "apple"}, {Char: "🍕", Name: "pizza"}, {Char: "🍣", Name: "sushi"}, {Char: "☕", Name: "coffee"}}},
}

func demoEmojiPicker(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	chosen := ui.Local(root, "chosen", func() string { return "" })
	root.Children(func() {
		if overlay.EmojiPicker(c, "Emoji", demoEmoji, chosen, overlay.EmojiPickerOptions{}).Changed() {
			log("chose " + *chosen)
		}
	})
}

func demoDropOverlay(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	files := ui.Local(root, "files", func() []string { return nil })
	root.Children(func() {
		if overlay.DropOverlay(c, files, overlay.DropOverlayOptions{}, func() {
			ui.Column(c).FillWidth().Height(140).Center().Radius(c.Theme().Radius).Border(1, c.Theme().Border).Children(func() {
				kit.Muted(c, "Drag files from your desktop here")
			})
		}).Changed() {
			log(fmt.Sprint("dropped ", len(*files), " files"))
		}
	})
}

// veiled is sample content for the veils to cover.
func veiled(c *ui.Context, title string) {
	t := c.Theme()
	ui.Column(c).FillWidth().Padding(t.Space(4)).Gap(t.Space(2)).Radius(t.Radius).Border(1, t.Border).Children(func() {
		ui.Text(c, title).Bold()
		kit.Muted(c, "Quarterly revenue grew 12% on the year.")
		button.Outline(c, "Open report")
	})
}

func demoBusyVeil(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Gap(12)
	busy := flag(root, "busy")
	root.Children(func() {
		word := "Start saving"
		if *busy {
			word = "Finish"
		}
		if button.Outline(c, word).Clicked() {
			*busy = !*busy
			log(fmt.Sprint("busy ", *busy))
		}
		overlay.BusyVeil(c, *busy, "Saving…", func() { veiled(c, "Report") })
	})
}

func demoLockOverlay(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360).Gap(12)
	locked := ui.Local(root, "locked", func() bool { return true })
	root.Children(func() {
		if !*locked && button.Outline(c, "Lock").Clicked() {
			*locked = true
			log("locked")
		}
		if overlay.LockOverlay(c, "Vault is locked", locked, overlay.LockOverlayOptions{Text: "Unlock to read your notes."},
			func() { veiled(c, "Notes") }).Changed() {
			log("unlocked")
		}
	})
}

func demoPaywall(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Children(func() {
		if _, up := overlay.Paywall(c, "Pro charts", "Upgrade to see every quarter.", overlay.PaywallOptions{},
			func() { veiled(c, "Charts") }); up.Clicked() {
			log("upgrade")
		}
	})
}
func demoShareSheet(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := flag(root, "open")
	chosen := ui.Local(root, "chosen", func() int { return -1 })
	targets := []overlay.ShareTarget{{Label: "Mail", Icon: "mail"}, {Label: "Messages", Icon: "message-circle"},
		{Label: "Notes", Icon: "notebook-pen"}, {Label: "Drive", Icon: "hard-drive"}}
	root.Children(func() {
		launch(c, "Share", open, log)
		if overlay.ShareSheet(c, "Share report", "https://crychic.dev/r/42", targets, chosen, open, overlay.ShareSheetOptions{}).Changed() {
			log(fmt.Sprintf("target %d, open %v", *chosen, *open))
		}
	})
}

func demoActionSheet(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := flag(root, "open")
	chosen := ui.Local(root, "chosen", func() int { return -1 })
	root.Children(func() {
		launch(c, "Message actions", open, log)
		if overlay.ActionSheet(c, "Message", overlayItems, chosen, open, overlay.ActionSheetOptions{}).Changed() {
			log(fmt.Sprintf("chose %d, open %v", *chosen, *open))
		}
	})
}

func demoMiniPlayer(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	playing := flag(root, "playing")
	root.Children(func() {
		if overlay.MiniPlayer(c, "Haruhikage", "CRYCHIC", playing, 0.4, overlay.MiniPlayerOptions{}).Changed() {
			log(fmt.Sprint("playing ", *playing))
		}
	})
}

func demoFloatingPanel(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	at := ui.Local(root, "at", func() ui.Point { return ui.Point{X: 16, Y: 16} })
	root.Children(func() {
		if overlay.FloatingPanel(c, "Layers", at, overlay.FloatingPanelOptions{Height: 260}, func() {
			kit.Muted(c, "Drag me anywhere.")
		}).Changed() {
			log(fmt.Sprintf("moved to %.0f, %.0f", at.X, at.Y))
		}
	})
}

func demoPiPBox(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	corner := ui.Local(root, "corner", func() int { return 3 })
	root.Children(func() {
		if overlay.PiPBox(c, "Call", corner, overlay.PiPBoxOptions{Height: 260}, func() {
			ui.Column(c).FillWidth().FillHeight().Center().Background(c.Theme().Accent.Alpha(0.12)).Children(func() {
				ui.Text(c, "Ada").Bold()
			})
		}).Changed() {
			log(fmt.Sprint("corner ", *corner))
		}
	})
}

func demoCalloutPopover(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := flag(root, "open")
	root.Children(func() {
		a := button.Outline(c, "Before you publish")
		if a.Clicked() {
			*open = !*open
			log(fmt.Sprint("open ", *open))
		}
		if overlay.CalloutPopover(c, a, open, feedback.ToneWarning, "Check the links", "Two links point to drafts.").Changed() {
			log("dismissed")
		}
	})
}

func demoInfoPopover(c *ui.Context, log func(string)) {
	root := ui.Row(c).Gap(8).AlignItems(ui.Center)
	open := flag(root, "open")
	root.Children(func() {
		ui.Text(c, "Price per seat")
		if overlay.InfoPopover(c, "About pricing", "Billed monthly; taxes added at checkout.", open).Changed() {
			log(fmt.Sprint("open ", *open))
		}
	})
}

func demoConfirmPopover(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := flag(root, "open")
	root.Children(func() {
		a := button.Outline(c, "Delete draft")
		if a.Clicked() {
			*open = !*open
			log(fmt.Sprint("open ", *open))
		}
		r, yes := overlay.ConfirmPopover(c, a, open, overlay.ConfirmPopoverOptions{Question: "Delete this draft?"})
		if yes.Clicked() {
			log("deleted")
		} else if r.Changed() {
			log("kept")
		}
	})
}

func demoSelectMenu(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := flag(root, "open")
	chosen := ui.Local(root, "chosen", func() int { return -1 })
	fruits := []string{"Apple", "Banana", "Cherry", "Grape", "Lemon", "Mango", "Peach", "Pear"}
	root.Children(func() {
		if overlay.SelectMenu(c, "Fruit", fruits, chosen, open, overlay.SelectMenuOptions{}).Changed() {
			log(fmt.Sprintf("chose %d, open %v", *chosen, *open))
		}
	})
}

// typedMenu opens a menu when its field gains mark.
func typedMenu(c *ui.Context, label, mark string, menu func(field ui.Element, text *string, open *bool, done func())) {
	root := ui.Column(c).Width(320)
	text := ui.Local(root, "text", func() string { return "" })
	open := flag(root, "open")
	jump := flag(root, "jump")
	root.Children(func() {
		field := ui.TextInput(c, text).Label(label).Placeholder("Type " + mark).FillWidth()
		// The field holds the new text only from this build.
		if *jump {
			field.SetTextSelection(utf8.RuneCountInString(*text), utf8.RuneCountInString(*text))
			*jump = false
		}
		word := (*text)[strings.LastIndex(*text, " ")+1:]
		if field.Changed() && strings.HasPrefix(word, mark) {
			*open = true
		}
		menu(field, text, open, func() {
			*jump = true
			c.Invalidate()
		})
	})
}

func demoMentionMenu(c *ui.Context, log func(string)) {
	people := []string{"Ada Lovelace", "Alan Turing", "Grace Hopper", "Linus Torvalds"}
	typedMenu(c, "Message", "@", func(field ui.Element, text *string, open *bool, done func()) {
		chosen := ui.Local(field, "chosen", func() int { return -1 })
		at := strings.LastIndex(*text, "@")
		if overlay.MentionMenu(c, "People", field, (*text)[at+1:], people, chosen, open, overlay.MentionMenuOptions{}).Changed() {
			if !*open && *chosen >= 0 && at >= 0 {
				*text = (*text)[:at] + "@" + people[*chosen] + " "
				*chosen = -1
				done()
				log("mentioned")
			} else {
				log(fmt.Sprint("open ", *open))
			}
		}
	})
}

func demoSlashMenu(c *ui.Context, log func(string)) {
	blocks := []overlay.MenuItem{{Label: "Heading", Icon: "heading"}, {Label: "List", Icon: "list"},
		{Label: "Image", Icon: "image"}, {Label: "Table", Icon: "table"}, {Label: "Quote", Icon: "quote"}}
	typedMenu(c, "Note", "/", func(field ui.Element, text *string, open *bool, done func()) {
		chosen := ui.Local(field, "chosen", func() int { return -1 })
		at := strings.LastIndex(*text, "/")
		if overlay.SlashMenu(c, "Blocks", field, (*text)[at+1:], blocks, chosen, open, overlay.SlashMenuOptions{}).Changed() {
			if !*open && *chosen >= 0 && at >= 0 {
				*text = (*text)[:at]
				done()
				log("inserted " + blocks[*chosen].Label)
				*chosen = -1
			} else {
				log(fmt.Sprint("open ", *open))
			}
		}
	})
}
func demoShortcutSheet(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := flag(root, "open")
	groups := []overlay.ShortcutGroup{
		{Title: "General", Shortcuts: []overlay.Shortcut{{Label: "Command palette", Keys: []string{"Ctrl", "K"}}, {Label: "Settings", Keys: []string{"Ctrl", ","}}}},
		{Title: "Editing", Shortcuts: []overlay.Shortcut{{Label: "Undo", Keys: []string{"Ctrl", "Z"}}, {Label: "Redo", Keys: []string{"Ctrl", "Shift", "Z"}}}},
		{Title: "View", Shortcuts: []overlay.Shortcut{{Label: "Zoom in", Keys: []string{"Ctrl", "+"}}, {Label: "Zoom out", Keys: []string{"Ctrl", "-"}}}},
	}
	root.Children(func() {
		launch(c, "Keyboard shortcuts", open, log)
		if overlay.ShortcutSheet(c, "Keyboard shortcuts", groups, open, overlay.ShortcutSheetOptions{}).Changed() {
			log("closed")
		}
	})
}

func demoStepDialog(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := flag(root, "open")
	step := ui.Local(root, "step", func() int { return 0 })
	steps := []string{"Account", "Plan", "Payment"}
	root.Children(func() {
		if launch(c, "Start setup", open, log).Clicked() {
			*step = 0
		}
		if overlay.StepDialog(c, "Set up workspace", steps, step, open, func(i int) {
			kit.Muted(c, fmt.Sprintf("Fill in the %s details.", strings.ToLower(steps[i])))
		}, overlay.StepDialogOptions{}).Changed() {
			log(fmt.Sprintf("step %d, open %v", *step, *open))
		}
	})
}

func demoRenameDialog(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(8)
	open := flag(root, "open")
	name := ui.Local(root, "name", func() string { return "report.pdf" })
	root.Children(func() {
		ui.Text(c, *name)
		launch(c, "Rename", open, log)
		before := *name
		if r, _ := overlay.RenameDialog(c, "Rename file", name, open, overlay.RenameDialogOptions{}); r.Changed() {
			if *name != before {
				log(fmt.Sprint("renamed, ", len(*name), " characters"))
			} else {
				log("closed")
			}
		}
	})
}

func demoSessionTimeout(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := flag(root, "open")
	deadline := ui.Local(root, "deadline", func() time.Time { return c.Now().Add(10 * time.Second) })
	root.Children(func() {
		if launch(c, "Simulate idle", open, log).Clicked() {
			*deadline = c.Now().Add(10 * time.Second)
		}
		r, stay, leave := overlay.SessionTimeout(c, *deadline, open, overlay.SessionTimeoutOptions{})
		switch {
		case stay.Clicked():
			log("stayed")
		case leave.Clicked():
			log("signed out")
		case r.Changed():
			log(fmt.Sprint("open ", *open))
		}
	})
}

func demoFeedbackWidget(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420).Height(300).Justify(ui.End)
	open := flag(root, "open")
	text := ui.Local(root, "text", func() string { return "" })
	root.Children(func() {
		r, send := overlay.FeedbackWidget(c, text, open, overlay.FeedbackWidgetOptions{})
		if send.Clicked() {
			log(fmt.Sprint("sent ", len(*text), " characters"))
			*text = ""
		} else if r.Changed() {
			log(fmt.Sprint("open ", *open, ", ", len(*text), " characters"))
		}
	})
}

func demoSheetStack(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	depth := ui.Local(root, "depth", func() int { return 0 })
	titles := []string{"Settings", "Account", "Password"}
	root.Children(func() {
		if button.Outline(c, "Open settings").Clicked() {
			*depth = 1
			log("depth 1")
		}
		if overlay.SheetStack(c, titles, depth, func(i int) {
			kit.Muted(c, titles[i]+" options go here.")
			if i+1 < len(titles) && button.Solid(c, "Open "+titles[i+1]).Clicked() {
				*depth = i + 2
				log(fmt.Sprint("depth ", *depth))
			}
		}, overlay.SheetStackOptions{}).Changed() {
			log(fmt.Sprint("depth ", *depth))
		}
	})
}
