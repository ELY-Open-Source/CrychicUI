package gallery

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/overlay"
	"github.com/egoist/mygo/ui"
)

var overlayEntries = []Entry{
	{611, "Modal", "overlay", "A dialog scaling in over a dim.", demoModal},
	{612, "ConfirmDialog", "overlay", "Title, text, Cancel and Confirm.", demoConfirmDialog},
	{613, "DangerDialog", "overlay", "A destructive confirm.", demoDangerDialog},
	{614, "PromptDialog", "overlay", "A dialog with one input.", demoPromptDialog},
	{615, "Drawer", "overlay", "A side panel sliding in.", demoDrawer},
	{616, "BottomSheet", "overlay", "A sheet dragged between heights.", demoBottomSheet},
	{617, "TopSheet", "overlay", "A panel dropping from the top.", demoTopSheet},
	{618, "PopoverCard", "overlay", "A card with an arrow.", demoPopoverCard},
	{619, "HoverCard", "overlay", "A card while hovering.", demoHoverCard},
	{620, "RichTooltip", "overlay", "A tooltip with a title.", demoRichTooltip},
	{621, "ContextMenu", "overlay", "A styled right-click menu.", demoContextMenu},
	{622, "Menu", "overlay", "A menu with icons and keys.", demoMenu},
	{623, "CommandPalette", "overlay", "Cmd+K commands.", demoCommandPalette},
	{624, "SpotlightSearch", "overlay", "Search with result groups.", demoSpotlightSearch},
	{625, "Lightbox", "overlay", "Images over a dark backdrop.", demoLightbox},
	{626, "FullscreenModal", "overlay", "A modal filling the window.", demoFullscreenModal},
	{627, "Inspector", "overlay", "A panel of properties.", demoInspector},
	{628, "Tour", "overlay", "Steps spotlighting anchors.", demoTour},
	{629, "Coachmark", "overlay", "One tip at an anchor.", demoCoachmark},
	{630, "WhatsNew", "overlay", "Pages of new features.", demoWhatsNew},
	{631, "NotificationPanel", "overlay", "A bell with grouped notices.", demoNotificationPanel},
	{632, "QuickView", "overlay", "A preview over a list item.", demoQuickView},
	{633, "PreviewPane", "overlay", "The selected item, cross-fading.", demoPreviewPane},
}

var overlayItems = []overlay.MenuItem{{Label: "Rename", Icon: "pencil", Keys: "F2"}, {Label: "Duplicate", Icon: "copy", Keys: "Ctrl D"},
	{Label: "Share", Icon: "share-2"}, {Label: "Delete", Icon: "trash-2", Keys: "Del"}}

// launch is a demo button that opens an overlay.
func launch(c *ui.Context, label string, open *bool, log func(string)) ui.Element {
	b := button.Outline(c, label)
	if b.Clicked() {
		*open = true
		log("opened")
	}
	return b
}

// flag is a demo bool kept across builds.
func flag(root ui.Element, key string) *bool {
	return ui.Local(root, key, func() bool { return false })
}

func demoModal(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := flag(root, "open")
	root.Children(func() {
		launch(c, "Open settings", open, log)
		if overlay.Modal(c, "Settings", open, overlay.ModalOptions{}, func() { kit.Muted(c, "Notifications, theme and language.") }).Changed() {
			log("closed")
		}
	})
}

func demoConfirmDialog(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := flag(root, "open")
	root.Children(func() {
		launch(c, "Publish", open, log)
		dialog, confirm := overlay.ConfirmDialog(c, "Publish post?", "Everyone will be able to read it.", open, overlay.ConfirmDialogOptions{})
		if dialog.Changed() {
			log("canceled")
		}
		if confirm.Clicked() {
			log("confirmed")
		}
	})
}

func demoDangerDialog(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := flag(root, "open")
	root.Children(func() {
		launch(c, "Delete project", open, log)
		dialog, confirm := overlay.DangerDialog(c, "Delete project?", "Its files and history go for good.", open, overlay.DangerDialogOptions{})
		if dialog.Changed() {
			log("canceled")
		}
		if confirm.Clicked() {
			log("deleted")
		}
	})
}

func demoPromptDialog(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(8)
	open := flag(root, "open")
	name := ui.Local(root, "name", func() string { return "Untitled" })
	root.Children(func() {
		launch(c, "Rename", open, log)
		kit.Muted(c, "Name: "+*name)
		before := *name
		if dialog, _ := overlay.PromptDialog(c, "Rename file", name, open, overlay.PromptDialogOptions{}); dialog.Changed() {
			if *name != before {
				log(fmt.Sprintf("renamed, %d characters", len(*name)))
			} else {
				log("closed")
			}
		}
	})
}

func demoDrawer(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := flag(root, "open")
	root.Children(func() {
		launch(c, "Filters", open, log)
		if overlay.Drawer(c, "Filters", open, overlay.DrawerOptions{}, func() { kit.Muted(c, "Price, size and color.") }).Changed() {
			log("closed")
		}
	})
}

func demoBottomSheet(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := flag(root, "open")
	snap := ui.Local(root, "snap", func() int { return 0 })
	root.Children(func() {
		launch(c, "Details", open, log)
		if overlay.BottomSheet(c, "Details", open, snap, overlay.BottomSheetOptions{}, func() { kit.Muted(c, "Drag the sheet up or down.") }).Changed() {
			log(fmt.Sprintf("open %v, height %d", *open, *snap))
		}
	})
}

func demoTopSheet(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := flag(root, "open")
	root.Children(func() {
		launch(c, "Announcements", open, log)
		if overlay.TopSheet(c, "Announcements", open, overlay.TopSheetOptions{}, func() { kit.Muted(c, "Maintenance on Sunday at 02:00.") }).Changed() {
			log("closed")
		}
	})
}

func demoPopoverCard(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := flag(root, "open")
	root.Children(func() {
		a := button.Outline(c, "About this plan")
		if a.Clicked() {
			*open = !*open
			log(fmt.Sprintf("open %v", *open))
		}
		if overlay.PopoverCard(c, a, open, func() { kit.Muted(c, "Five seats, 100 GB and priority support.") }).Changed() {
			log("dismissed")
		}
	})
}

func demoHoverCard(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	root.Children(func() {
		a := ui.Text(c, "@ada").TextColor(c.Theme().Accent).Label("Ada")
		overlay.HoverCard(c, "Ada Lovelace", a, func() {
			ui.Text(c, "Ada Lovelace").Bold()
			kit.Muted(c, "Wrote the first program.")
		})
	})
}

func demoRichTooltip(c *ui.Context, log func(string)) {
	a := button.Outline(c, "Share")
	overlay.RichTooltip(c, a, "Share", "Send a link to this page.")
	if a.Clicked() {
		log("shared")
	}
}

func demoContextMenu(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	chosen := ui.Local(root, "chosen", func() int { return -1 })
	open := flag(root, "open")
	root.Children(func() {
		target := ui.ButtonBase(c).Label("report.pdf").Padding(16).Radius(8).Border(1, c.Theme().Border).
			Children(func() { ui.Text(c, "Right-click report.pdf") })
		if overlay.ContextMenu(c, "File actions", target, overlayItems, chosen, open).Changed() {
			log(fmt.Sprintf("open %v, chose %d", *open, *chosen))
		}
	})
}

func demoMenu(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	chosen := ui.Local(root, "chosen", func() int { return -1 })
	open := flag(root, "open")
	root.Children(func() {
		if overlay.Menu(c, "Actions", overlayItems, chosen, open).Changed() {
			log(fmt.Sprintf("open %v, chose %d", *open, *chosen))
		}
	})
}

func demoCommandPalette(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(8)
	chosen := ui.Local(root, "chosen", func() int { return -1 })
	open := flag(root, "open")
	root.Children(func() {
		launch(c, "Commands", open, log)
		kit.Muted(c, "Or press Ctrl+K.")
		if overlay.CommandPalette(c, "Commands", overlayItems, chosen, open, overlay.CommandPaletteOptions{}).Changed() {
			log(fmt.Sprintf("open %v, chose %d", *open, *chosen))
		}
	})
}

func demoSpotlightSearch(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	chosen := ui.Local(root, "chosen", func() [2]int { return [2]int{-1, -1} })
	open := flag(root, "open")
	groups := []overlay.ResultGroup{{Title: "Files", Results: []string{"budget.xlsx", "roadmap.md", "logo.svg"}},
		{Title: "People", Results: []string{"Ada Lovelace", "Alan Turing"}}, {Title: "Settings", Results: []string{"Appearance", "Accounts"}}}
	root.Children(func() {
		launch(c, "Search", open, log)
		if overlay.SpotlightSearch(c, "Search", groups, chosen, open, overlay.SpotlightSearchOptions{}).Changed() {
			log(fmt.Sprintf("open %v, chose %v", *open, *chosen))
		}
	})
}

func demoLightbox(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	at := ui.Local(root, "at", func() int { return 0 })
	open := flag(root, "open")
	root.Children(func() {
		launch(c, "View photos", open, log)
		if overlay.Lightbox(c, "Photos", photos(), []string{"Sunset", "Lake", "Meadow"}, at, open, overlay.LightboxOptions{}).Changed() {
			log(fmt.Sprintf("open %v, photo %d", *open, *at+1))
		}
	})
}

func demoFullscreenModal(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := flag(root, "open")
	root.Children(func() {
		launch(c, "Open editor", open, log)
		if overlay.FullscreenModal(c, "Editor", open, overlay.FullscreenModalOptions{}, func() { kit.Muted(c, "A full-window canvas.") }).Changed() {
			log("closed")
		}
	})
}

func demoInspector(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := flag(root, "open")
	props := []overlay.Prop{{Name: "Width", Value: "320"}, {Name: "Height", Value: "240"}, {Name: "Fill", Value: "Accent"}, {Name: "Radius", Value: "8"}}
	root.Children(func() {
		launch(c, "Inspect layer", open, log)
		if overlay.Inspector(c, "Layer", props, open, overlay.InspectorOptions{}).Changed() {
			log("closed")
		}
	})
}

func demoTour(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(16)
	open := flag(root, "open")
	step := ui.Local(root, "step", func() int { return 0 })
	steps := []overlay.TourStep{{Title: "Search", Text: "Find anything here."}, {Title: "Create", Text: "Start a new document."}}
	root.Children(func() {
		if launch(c, "Start tour", open, log).Clicked() {
			*step = 0
		}
		ui.Row(c).Gap(24).Children(func() {
			search := button.IconButton(c, "Search", "search")
			create := button.Solid(c, "New document")
			if overlay.Tour(c, "Welcome tour", []ui.Element{search, create}, steps, step, open, overlay.TourOptions{}).Changed() {
				log(fmt.Sprintf("open %v, step %d", *open, *step+1))
			}
		})
	})
}

func demoCoachmark(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(16)
	open := flag(root, "open")
	root.Children(func() {
		launch(c, "Show tip", open, log)
		a := button.Outline(c, "Export")
		if overlay.Coachmark(c, a, "New: export", "Save any page as a PDF.", open, overlay.CoachmarkOptions{}).Changed() {
			log("dismissed")
		}
	})
}

func demoWhatsNew(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := flag(root, "open")
	page := ui.Local(root, "page", func() int { return 0 })
	features := []overlay.Feature{{Icon: "palette", Title: "Themes", Text: "Pick light, dark or your own colors."},
		{Icon: "zap", Title: "Faster", Text: "Pages open twice as fast."}, {Icon: "bell", Title: "Alerts", Text: "Hear about what matters."}}
	root.Children(func() {
		if launch(c, "What's new", open, log).Clicked() {
			*page = 0
		}
		if overlay.WhatsNew(c, "What's new", features, page, open, overlay.WhatsNewOptions{}).Changed() {
			log(fmt.Sprintf("open %v, page %d", *open, *page+1))
		}
	})
}

func demoNotificationPanel(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := flag(root, "open")
	groups := []overlay.NoticeGroup{{Title: "Today", Notices: []overlay.Notice{{Title: "Build passed", Text: "main is green", Unread: true},
		{Title: "New comment", Text: "Ada replied to your post", Unread: true}}}, {Title: "Earlier", Notices: []overlay.Notice{{Title: "Deployed", Text: "Version 1.2 is live"}}}}
	root.Children(func() {
		if overlay.NotificationPanel(c, "Notifications", groups, open, overlay.NotificationPanelOptions{}).Changed() {
			log(fmt.Sprintf("open %v", *open))
		}
	})
}

func demoQuickView(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320).Gap(4)
	open := flag(root, "open")
	at := ui.Local(root, "at", func() int { return 0 })
	root.Children(func() {
		for i, name := range []string{"Invoice 1042", "Invoice 1043"} {
			c.Key(i)
			row := ui.ButtonBase(c).Label(name).FillWidth().Padding(10).Radius(8).Border(1, c.Theme().Border).Children(func() { ui.Text(c, name) })
			if row.Clicked() {
				*at, *open = i, true
				log("previewed " + name)
			}
			if i == *at && overlay.QuickView(c, row, name, open, overlay.QuickViewOptions{}, func() { kit.Muted(c, "Total 1,240.00, due in 14 days") }).Changed() {
				log("closed")
			}
		}
	})
}
func demoPreviewPane(c *ui.Context, log func(string)) {
	root := ui.Row(c).Width(460).Gap(12).AlignItems(ui.Start)
	sel := ui.Local(root, "sel", func() int { return -1 })
	names := []string{"Design brief", "Budget", "Timeline"}
	root.Children(func() {
		ui.Column(c).Width(140).Gap(4).Children(func() {
			for i, n := range names {
				c.Key(i)
				if button.Ghost(c, n).Clicked() {
					*sel = i
					log("selected " + n)
				}
			}
		})
		overlay.PreviewPane(c, "Preview", *sel, overlay.PreviewPaneOptions{}, func() {
			ui.Text(c, names[*sel]).Bold()
			kit.Muted(c, "Last edited today.")
		}).Grow(1)
	})
}
