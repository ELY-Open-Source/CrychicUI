package gallery

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/egoist/mygo/ui"
)

var feedbackEntries = []Entry{
	{296, "Alert", "feedback", "Tinted message box with an icon; slides in.", demoAlert},
	{297, "DismissibleAlert", "feedback", "Alert with a close button; collapses out.", demoDismissibleAlert},
	{298, "Banner", "feedback", "Full-width bar with action and close; slides down.", demoBanner},
	{299, "Callout", "feedback", "Tinted box whose left border grows in.", demoCallout},
	{300, "InlineMessage", "feedback", "Icon and text in a line; fades in.", demoInlineMessage},
	{301, "Toaster", "feedback", "Styled toasts for c.AddToast; they slide and stack.", demoToaster},
	{302, "ToastCard", "feedback", "One toast card with action and close.", demoToastCard},
	{303, "Snackbar", "feedback", "Dark bar with one action; slides up.", demoSnackbar},
	{304, "Notice", "feedback", "Centered message on a surface; scales in.", demoNotice},
	{305, "ProgressBar", "feedback", "Bar whose fill eases to the value.", demoProgressBar},
	{306, "IndeterminateBar", "feedback", "Segment sliding along a track.", demoIndeterminateBar},
	{307, "StripedProgress", "feedback", "Bar with moving stripes.", demoStripedProgress},
	{308, "SegmentProgress", "feedback", "Segments filling in order.", demoSegmentProgress},
	{309, "ProgressRing", "feedback", "Circle whose arc draws; the value rolls.", demoProgressRing},
	{310, "Gauge", "feedback", "Half circle whose needle springs.", demoGauge},
	{311, "Spin", "feedback", "Arc spinner in three theme sizes.", demoSpin},
	{312, "DotsLoader", "feedback", "Three dots bouncing in turn.", demoDotsLoader},
	{313, "PulseLoader", "feedback", "Circle sending rings outward.", demoPulseLoader},
	{314, "BarsLoader", "feedback", "Equalizer bars rising and falling.", demoBarsLoader},
	{315, "OrbitLoader", "feedback", "Dot orbiting a ring.", demoOrbitLoader},
	{316, "Skeleton", "feedback", "Placeholder block with a moving sheen.", demoSkeleton},
	{317, "SkeletonText", "feedback", "Sheen lines, the last one shorter.", demoSkeletonText},
	{318, "SkeletonCard", "feedback", "Image block, title and lines of sheen.", demoSkeletonCard},
	{319, "SkeletonList", "feedback", "Rows of avatar and lines of sheen.", demoSkeletonList},
	{320, "EmptyState", "feedback", "Icon, title, text and action; rises in.", demoEmptyState},
	{321, "ErrorState", "feedback", "Shaking error icon, message and Retry.", demoErrorState},
	{322, "SuccessState", "feedback", "Check that draws, title and action.", demoSuccessState},
	{323, "NotFound", "feedback", "404 state with a home action; rises in.", demoNotFound},
}

// progressControls steps a demo's value with two buttons.
func progressControls(c *ui.Context, v *float64, log func(string)) {
	ui.Row(c).Gap(8).Children(func() {
		if ui.Button(c, "Less").Disabled(*v <= 0).Clicked() {
			*v = max(0, *v-0.25)
			log(fmt.Sprintf("value: %.2f", *v))
		}
		if ui.Button(c, "More").Disabled(*v >= 1).Clicked() {
			*v = min(1, *v+0.25)
			log(fmt.Sprintf("value: %.2f", *v))
		}
	})
}

func demoAlert(c *ui.Context, log func(string)) {
	ui.Column(c).Gap(8).MaxWidth(480).Children(func() {
		feedback.Alert(c, feedback.ToneInfo, "Heads up", "The beta ends on Friday.")
		feedback.Alert(c, feedback.ToneSuccess, "Deployed", "Version 2.4 is live.")
		feedback.Alert(c, feedback.ToneWarning, "Disk almost full", "Free some space soon.")
		feedback.Alert(c, feedback.ToneDanger, "Build failed", "Three tests did not pass.")
	})
}

func demoDismissibleAlert(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(8).MaxWidth(480)
	open := boolState(root, "open", true)
	root.Children(func() {
		if feedback.DismissibleAlert(c, feedback.ToneWarning, "Trial ends soon", "Add a card to keep your data.", open, feedback.DismissibleAlertOptions{}).Changed() {
			log("dismissed")
		}
		if !*open && ui.Button(c, "Show again").Clicked() {
			*open = true
			log("shown")
		}
	})
}

func demoBanner(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(8)
	open := boolState(root, "open", true)
	root.Children(func() {
		bar, act := feedback.Banner(c, "A new version is ready.", "Update", open, feedback.BannerOptions{})
		if act.Clicked() {
			log("update")
		}
		if bar.Changed() {
			log("closed")
		}
		if !*open && ui.Button(c, "Show again").Clicked() {
			*open = true
			log("shown")
		}
	})
}

func demoCallout(c *ui.Context, log func(string)) {
	ui.Column(c).MaxWidth(480).Children(func() {
		feedback.Callout(c, feedback.ToneInfo, "Tip", func() {
			ui.Text(c, "Press Ctrl+K anywhere to search every component.")
		})
	})
}

func demoInlineMessage(c *ui.Context, log func(string)) {
	ui.Column(c).Gap(8).Children(func() {
		feedback.InlineMessage(c, feedback.ToneSuccess, "Changes saved")
		feedback.InlineMessage(c, feedback.ToneDanger, "Name is taken")
	})
}

func demoToaster(c *ui.Context, log func(string)) {
	root := ui.Row(c).Gap(8)
	n := intState(root, "n", 0)
	root.Children(func() {
		for _, kind := range []string{"info", "success", "error"} {
			if ui.Button(c, "Toast "+kind).Clicked() {
				*n++
				id := fmt.Sprint(*n)
				c.AddToast(ui.Toast{ID: id, Title: "Toast " + id, Description: "A " + kind + " message.", Type: kind,
					Action: "Undo", OnAction: func() { log("undo " + id) }, OnClose: func() { log("closed " + id) }})
				log("added " + id)
			}
		}
	})
	feedback.Toaster(c, feedback.ToasterOptions{})
}

func demoToastCard(c *ui.Context, log func(string)) {
	_, act, dismiss := feedback.ToastCard(c, feedback.ToneSuccess, "Invite sent", "Ada will get an email.", feedback.ToastCardOptions{Action: "Copy link"})
	if act.Clicked() {
		log("copy link")
	}
	if dismiss.Clicked() {
		log("close")
	}
}

func demoSnackbar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).Width(420)
	open := boolState(root, "open", true)
	root.Children(func() {
		if ui.Button(c, "Archive").Clicked() {
			*open = true
			log("archived")
		}
		if _, act := feedback.Snackbar(c, "Conversation archived", "Undo", open); act.Clicked() {
			log("undo")
		}
	})
}

func demoNotice(c *ui.Context, log func(string)) {
	feedback.Notice(c, "inbox", "All caught up", "Nothing new since this morning.")
}

func demoProgressBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).Width(320)
	v := numberState(root, "v", 0.5)
	root.Children(func() {
		feedback.ProgressBar(c, "Upload", *v)
		progressControls(c, v, log)
	})
}

func demoIndeterminateBar(c *ui.Context, log func(string)) {
	ui.Column(c).Width(320).Children(func() { feedback.IndeterminateBar(c, "Connecting") })
}

func demoStripedProgress(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).Width(320)
	v := numberState(root, "v", 0.5)
	root.Children(func() {
		feedback.StripedProgress(c, "Sync", *v)
		progressControls(c, v, log)
	})
}

func demoSegmentProgress(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).Width(320)
	v := numberState(root, "v", 0.5)
	root.Children(func() {
		feedback.SegmentProgress(c, "Steps", 4, int(*v*4))
		progressControls(c, v, log)
	})
}

func demoProgressRing(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	v := numberState(root, "v", 0.75)
	root.Children(func() {
		feedback.ProgressRing(c, "Storage", *v, feedback.ProgressRingOptions{})
		progressControls(c, v, log)
	})
}

func demoGauge(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	v := numberState(root, "v", 0.5)
	root.Children(func() {
		feedback.Gauge(c, "CPU", *v*100, 0, 100, feedback.GaugeOptions{Format: "%.0f%%"})
		progressControls(c, v, log)
	})
}

func demoSpin(c *ui.Context, log func(string)) {
	ui.Row(c).Gap(16).Children(func() {
		feedback.Spin(c, "Loading small", feedback.SpinOptions{Size: feedback.SpinSmall})
		feedback.Spin(c, "Loading", feedback.SpinOptions{})
		feedback.Spin(c, "Loading large", feedback.SpinOptions{Size: feedback.SpinLarge})
	})
}

func demoDotsLoader(c *ui.Context, log func(string)) { feedback.DotsLoader(c, "Loading") }

func demoPulseLoader(c *ui.Context, log func(string)) { feedback.PulseLoader(c, "Loading") }

func demoBarsLoader(c *ui.Context, log func(string)) { feedback.BarsLoader(c, "Loading") }

func demoOrbitLoader(c *ui.Context, log func(string)) { feedback.OrbitLoader(c, "Loading") }

func demoSkeleton(c *ui.Context, log func(string)) {
	ui.Column(c).Gap(8).Width(320).Children(func() {
		feedback.Skeleton(c)
		feedback.Skeleton(c).Height(80)
	})
}

func demoSkeletonText(c *ui.Context, log func(string)) {
	ui.Column(c).Width(320).Children(func() { feedback.SkeletonText(c, 4) })
}

func demoSkeletonCard(c *ui.Context, log func(string)) { feedback.SkeletonCard(c) }

func demoSkeletonList(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Children(func() { feedback.SkeletonList(c, 3) })
}

func demoEmptyState(c *ui.Context, log func(string)) {
	if _, act := feedback.EmptyState(c, "inbox", "No messages", "Mail you get shows up here.", "Compose"); act.Clicked() {
		log("compose")
	}
}

func demoErrorState(c *ui.Context, log func(string)) {
	if _, retry := feedback.ErrorState(c, "The server did not answer in time.", feedback.ErrorStateOptions{}); retry.Clicked() {
		log("retry")
	}
}

func demoSuccessState(c *ui.Context, log func(string)) {
	if _, act := feedback.SuccessState(c, "Payment complete", "A receipt is on its way.", "View order"); act.Clicked() {
		log("view order")
	}
}

func demoNotFound(c *ui.Context, log func(string)) {
	if _, home := feedback.NotFound(c, feedback.NotFoundOptions{}); home.Clicked() {
		log("home")
	}
}
