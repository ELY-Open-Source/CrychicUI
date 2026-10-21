package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/egoist/mygo/ui"
)

var buttonEntries2 = []Entry{
	{114, "DownloadButton", "button", "Progress fill and done check.", demoDownloadButton},
	{115, "UploadButton", "button", "Dropped path, progress fill and done check.", demoUploadButton},
	{116, "SubmitButton", "button", "Idle, loading, success and error states.", demoSubmitButton},
	{117, "RippleButton", "button", "Ripple growing from the press point.", demoRippleButton},
	{118, "MagneticButton", "button", "Button leaning toward the pointer.", demoMagneticButton},
	{119, "ShimmerButton", "button", "Sheen sweeping across in a loop.", demoShimmerButton},
	{120, "PulseButton", "button", "Pulsing ring until clicked.", demoPulseButton},
	{121, "ArrowButton", "button", "Arrow sliding right on hover.", demoArrowButton},
	{122, "ChipButton", "button", "Filter chip with a growing check.", demoChipButton},
	{123, "CloseButton", "button", "Cross turning ninety degrees on hover.", demoCloseButton},
	{124, "MoreButton", "button", "Three dots opening a menu.", demoMoreButton},
	{125, "DropdownButton", "button", "Label with a flipping chevron and menu.", demoDropdownButton},
	{126, "TextButton", "button", "Inline text action with animated color.", demoTextButton},
	{127, "ProviderButton", "button", "Sign-in button with a provider icon.", demoProviderButton},
	{128, "ThemeToggle", "button", "Sun and moon swapping with a turn.", demoThemeToggle},
	{129, "VoiceButton", "button", "Mic toggle with pulsing recording rings.", demoVoiceButton},
	{130, "PlayButton", "button", "Play and pause icons crossing.", demoPlayButton},
	{131, "ExpandButton", "button", "Chevron turning one hundred eighty degrees.", demoExpandButton},
	{132, "SortButton", "button", "Cycles none, ascending and descending.", demoSortButton},
	{133, "RefreshButton", "button", "Icon spinning while refreshing.", demoRefreshButton},
	{134, "ShortcutButton", "button", "Click or shortcut firing one action.", demoShortcutButton},
	{135, "SlideConfirm", "button", "Thumb sliding to confirm or springing back.", demoSlideConfirm},
}

func demoDownloadButton(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	progress := ui.Local(root, "progress", func() float64 { return 0 })
	until := ui.Local(root, "until", func() time.Time { return time.Time{} })
	if *progress > 0 && *progress < 1 {
		if !c.Now().Before(*until) {
			*progress = min(1, *progress+0.25)
			*until = c.Now().Add(100 * time.Millisecond)
			c.After(100 * time.Millisecond)
			log(fmt.Sprintf("progress: %.0f", *progress*100))
		} else {
			c.After(until.Sub(c.Now()))
		}
	}
	if button.DownloadButton(c, "Download", *progress).Clicked() {
		*progress = 0.25
		*until = c.Now().Add(100 * time.Millisecond)
		c.After(100 * time.Millisecond)
		log("start")
	}
}

func demoUploadButton(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	path := ui.Local(root, "path", func() string { return "" })
	progress := ui.Local(root, "progress", func() float64 { return 0 })
	until := ui.Local(root, "until", func() time.Time { return time.Time{} })
	if *progress > 0 && *progress < 1 {
		if !c.Now().Before(*until) {
			*progress = min(1, *progress+0.25)
			*until = c.Now().Add(100 * time.Millisecond)
			c.After(100 * time.Millisecond)
			log(fmt.Sprintf("progress: %.0f", *progress*100))
		} else {
			c.After(until.Sub(c.Now()))
		}
	}
	b := button.UploadButton(c, "Upload", path, *progress)
	if b.Changed() {
		log("files: 1")
	}
	if b.Clicked() {
		*progress = 0.25
		*until = c.Now().Add(100 * time.Millisecond)
		c.After(100 * time.Millisecond)
		log("start")
	}
}

func demoSubmitButton(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	state := ui.Local(root, "state", func() button.SubmitState { return button.SubmitIdle })
	until := ui.Local(root, "until", func() time.Time { return time.Time{} })
	if *state == button.SubmitLoading {
		if !c.Now().Before(*until) {
			*state = button.SubmitSuccess
			log("success")
		} else {
			c.After(until.Sub(c.Now()))
		}
	}
	if button.SubmitButton(c, "Submit", *state).Clicked() && *state != button.SubmitLoading {
		*state = button.SubmitLoading
		*until = c.Now().Add(500 * time.Millisecond)
		c.After(500 * time.Millisecond)
		log("loading")
	}
}

func demoRippleButton(c *ui.Context, log func(string)) {
	if button.RippleButton(c, "Ripple").Clicked() {
		log("ripple")
	}
}

func demoMagneticButton(c *ui.Context, log func(string)) {
	if button.MagneticButton(c, "Magnetic").Clicked() {
		log("magnetic")
	}
}

func demoShimmerButton(c *ui.Context, log func(string)) {
	if button.ShimmerButton(c, "Shimmer").Clicked() {
		log("shimmer")
	}
}

func demoPulseButton(c *ui.Context, log func(string)) {
	if button.PulseButton(c, "Pulse").Clicked() {
		log("pulse")
	}
}

func demoArrowButton(c *ui.Context, log func(string)) {
	if button.ArrowButton(c, "Continue").Clicked() {
		log("continue")
	}
}

func demoChipButton(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	on := ui.Local(root, "on", func() bool { return false })
	root.Children(func() {
		if button.ChipButton(c, "Active", on).Changed() {
			log(fmt.Sprintf("on: %v", *on))
		}
	})
}

func demoCloseButton(c *ui.Context, log func(string)) {
	if button.CloseButton(c, "Close").Clicked() {
		log("close")
	}
}

func demoMoreButton(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	selected := ui.Local(root, "selected", func() int { return -1 })
	targets := []string{"Rename", "Delete"}
	root.Children(func() {
		if button.MoreButton(c, "More", targets, selected).Changed() {
			log("selected: " + targets[*selected])
		}
	})
}

func demoDropdownButton(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	open := ui.Local(root, "open", func() bool { return false })
	selected := ui.Local(root, "selected", func() int { return -1 })
	root.Children(func() {
		if button.DropdownButton(c, "Assign", []string{"Email", "Call"}, open, selected).Changed() {
			log(fmt.Sprintf("open: %v, selected: %d", *open, *selected))
		}
	})
}

func demoTextButton(c *ui.Context, log func(string)) {
	if button.TextButton(c, "Learn more").Clicked() {
		log("learn more")
	}
}

func demoProviderButton(c *ui.Context, log func(string)) {
	if button.ProviderButton(c, "Continue with Provider", "globe").Clicked() {
		log("provider")
	}
}

func demoThemeToggle(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	dark := ui.Local(root, "dark", func() bool { return false })
	root.Children(func() {
		if button.ThemeToggle(c, dark, button.ThemeToggleOptions{}).Changed() {
			log(fmt.Sprintf("dark: %v", *dark))
		}
	})
}

func demoVoiceButton(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	recording := ui.Local(root, "recording", func() bool { return false })
	root.Children(func() {
		if button.VoiceButton(c, recording, button.VoiceButtonOptions{}).Changed() {
			log(fmt.Sprintf("recording: %v", *recording))
		}
	})
}

func demoPlayButton(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	playing := ui.Local(root, "playing", func() bool { return false })
	root.Children(func() {
		if button.PlayButton(c, playing, button.PlayButtonOptions{}).Changed() {
			log(fmt.Sprintf("playing: %v", *playing))
		}
	})
}

func demoExpandButton(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	expanded := ui.Local(root, "expanded", func() bool { return false })
	root.Children(func() {
		if button.ExpandButton(c, expanded, button.ExpandButtonOptions{}).Changed() {
			log(fmt.Sprintf("expanded: %v", *expanded))
		}
	})
}

func demoSortButton(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	direction := ui.Local(root, "direction", func() button.SortDirection { return button.SortNone })
	root.Children(func() {
		if button.SortButton(c, "Sort", direction).Changed() {
			log(fmt.Sprintf("direction: %d", *direction))
		}
	})
}

func demoRefreshButton(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	refreshing := ui.Local(root, "refreshing", func() bool { return false })
	until := ui.Local(root, "until", func() time.Time { return time.Time{} })
	if *refreshing {
		if !c.Now().Before(*until) {
			*refreshing = false
			log("ready")
		} else {
			c.After(until.Sub(c.Now()))
		}
	}
	if button.RefreshButton(c, "Refresh", *refreshing).Clicked() {
		*refreshing = true
		*until = c.Now().Add(time.Second)
		c.After(time.Second)
		log("refresh")
	}
}

func demoShortcutButton(c *ui.Context, log func(string)) {
	if button.ShortcutButton(c, "Command", ui.Cmd, ui.KeyK, []string{"Ctrl", "K"}, button.ShortcutButtonOptions{}).Changed() {
		log("shortcut")
	}
}

func demoSlideConfirm(c *ui.Context, log func(string)) {
	if button.SlideConfirm(c, "Slide to confirm", button.SlideConfirmOptions{}).Changed() {
		log("confirmed")
	}
}
