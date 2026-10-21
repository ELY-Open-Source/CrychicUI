package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/egoist/mygo/ui"
)

var buttonEntries = []Entry{
	{91, "Solid", "button", "Accent button with animated hover and press.", demoSolid},
	{92, "Soft", "button", "Accent-tinted button with animated hover and press.", demoSoft},
	{93, "Outline", "button", "Bordered button with animated border color.", demoOutline},
	{94, "Ghost", "button", "Transparent button gaining a hover surface.", demoGhost},
	{95, "Destructive", "button", "Danger button with a delete cue.", demoDestructive},
	{96, "IconButton", "button", "Square icon button with an accessible label.", demoIconButton},
	{97, "IconLabel", "button", "Icon and text with a nudging icon.", demoIconLabel},
	{98, "Pill", "button", "Fully rounded solid button.", demoPill},
	{99, "Fab", "button", "Floating round action bouncing in.", demoFab},
	{100, "SpeedDial", "button", "Fab fanning actions out with a stagger.", demoSpeedDial},
	{101, "SplitButton", "button", "Main action beside a chevron menu.", demoSplitButton},
	{102, "ButtonGroup", "button", "Joined buttons with a tinted selection.", demoButtonGroup},
	{103, "ToggleButton", "button", "Pressed state with an animated tint.", demoToggleButton},
	{104, "LoadingButton", "button", "Spinner and disabled state while busy.", demoLoadingButton},
	{105, "ConfirmButton", "button", "Two-step confirmation with a timeout.", demoConfirmButton},
	{106, "HoldButton", "button", "Hold fill completing an action.", demoHoldButton},
	{107, "CopyButton", "button", "Clipboard action with a temporary check.", demoCopyButton},
	{108, "ShareButton", "button", "Menu of share targets.", demoShareButton},
	{109, "LikeButton", "button", "Heart filling and bursting as the count rolls.", demoLikeButton},
	{110, "StarButton", "button", "Star filling with a spin.", demoStarButton},
	{111, "FollowButton", "button", "Follow and Following with a width transition.", demoFollowButton},
	{112, "CounterButton", "button", "Label with a rolling count badge.", demoCounterButton},
	{113, "UndoButton", "button", "Countdown ring with an undo action.", demoUndoButton},
}

func demoSolid(c *ui.Context, log func(string)) {
	if button.Solid(c, "Continue").Clicked() {
		log("continue")
	}
}

func demoSoft(c *ui.Context, log func(string)) {
	if button.Soft(c, "Continue").Clicked() {
		log("continue")
	}
}

func demoOutline(c *ui.Context, log func(string)) {
	if button.Outline(c, "Continue").Clicked() {
		log("continue")
	}
}

func demoGhost(c *ui.Context, log func(string)) {
	if button.Ghost(c, "Continue").Clicked() {
		log("continue")
	}
}

func demoDestructive(c *ui.Context, log func(string)) {
	if button.Destructive(c, "Delete").Clicked() {
		log("delete")
	}
}

func demoIconButton(c *ui.Context, log func(string)) {
	if button.IconButton(c, "Edit", "pencil").Clicked() {
		log("edit")
	}
}

func demoIconLabel(c *ui.Context, log func(string)) {
	if button.IconLabel(c, "Edit", "pencil").Clicked() {
		log("edit")
	}
}

func demoPill(c *ui.Context, log func(string)) {
	if button.Pill(c, "Continue").Clicked() {
		log("continue")
	}
}

func demoFab(c *ui.Context, log func(string)) {
	if button.Fab(c, "Add", "plus").Clicked() {
		log("add")
	}
}

func demoSpeedDial(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	open := ui.Local(root, "open", func() bool { return false })
	selected := ui.Local(root, "selected", func() int { return -1 })
	actions := []button.SpeedDialAction{{Label: "Edit", Icon: "pencil"}, {Label: "Share", Icon: "share-2"}}
	root.Children(func() {
		if button.SpeedDial(c, open, actions, selected, button.SpeedDialOptions{}).Changed() {
			log(fmt.Sprintf("open: %v, selected: %d", *open, *selected))
		}
	})
}

func demoSplitButton(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	selected := ui.Local(root, "selected", func() int { return -1 })
	root.Children(func() {
		view, main := button.SplitButton(c, "Publish", []string{"Draft", "Schedule"}, selected)
		if view.Changed() {
			log(fmt.Sprintf("selected: %d", *selected))
		}
		if main.Clicked() {
			log("publish")
		}
	})
}

func demoButtonGroup(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	selected := ui.Local(root, "selected", func() int { return 0 })
	items := []string{"Day", "Week", "Month"}
	root.Children(func() {
		if button.ButtonGroup(c, items, selected).Changed() {
			log("selected: " + items[*selected])
		}
	})
}

func demoToggleButton(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	on := ui.Local(root, "on", func() bool { return false })
	root.Children(func() {
		if button.ToggleButton(c, "Bold", on).Changed() {
			log(fmt.Sprintf("on: %v", *on))
		}
	})
}

func demoLoadingButton(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	busy := ui.Local(root, "busy", func() bool { return false })
	until := ui.Local(root, "until", func() time.Time { return time.Time{} })
	if *busy && !c.Now().Before(*until) {
		*busy = false
		log("ready")
	}
	if button.LoadingButton(c, "Save", *busy).Clicked() {
		*busy = true
		*until = c.Now().Add(2 * time.Second)
		c.After(2 * time.Second)
		log("saving")
	}
}

func demoConfirmButton(c *ui.Context, log func(string)) {
	b := button.ConfirmButton(c, "Delete", button.ConfirmButtonOptions{})
	if b.Clicked() {
		log("click")
	}
	if b.Changed() {
		log("confirmed")
	}
}

func demoHoldButton(c *ui.Context, log func(string)) {
	if button.HoldButton(c, "Hold to confirm", button.HoldButtonOptions{}).Changed() {
		log("completed")
	}
}

func demoCopyButton(c *ui.Context, log func(string)) {
	if button.CopyButton(c, "Copy", "CrychicUI").Clicked() {
		log("copied")
	}
}

func demoShareButton(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	selected := ui.Local(root, "selected", func() int { return -1 })
	targets := []string{"Email", "Copy link"}
	root.Children(func() {
		if button.ShareButton(c, "Share", targets, selected).Changed() {
			log("selected: " + targets[*selected])
		}
	})
}

func demoLikeButton(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	liked := ui.Local(root, "liked", func() bool { return false })
	likes := ui.Local(root, "likes", func() int { return 4 })
	root.Children(func() {
		if button.LikeButton(c, liked, likes, button.LikeButtonOptions{}).Changed() {
			log(fmt.Sprintf("likes: %d", *likes))
		}
	})
}

func demoStarButton(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	on := ui.Local(root, "on", func() bool { return false })
	root.Children(func() {
		if button.StarButton(c, on, button.StarButtonOptions{}).Changed() {
			log(fmt.Sprintf("on: %v", *on))
		}
	})
}

func demoFollowButton(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	following := ui.Local(root, "following", func() bool { return false })
	root.Children(func() {
		if button.FollowButton(c, following, button.FollowButtonOptions{}).Changed() {
			log(fmt.Sprintf("following: %v", *following))
		}
	})
}

func demoCounterButton(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	n := ui.Local(root, "count", func() int { return 2 })
	root.Children(func() {
		if button.CounterButton(c, "Notifications", n).Changed() {
			log(fmt.Sprintf("count: %d", *n))
		}
	})
}

func demoUndoButton(c *ui.Context, log func(string)) {
	b := button.UndoButton(c, "Undo", button.UndoButtonOptions{})
	if b.Clicked() {
		log("undo")
	}
	if b.Changed() {
		log("commit")
	}
}
