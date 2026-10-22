package gallery

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/motion"
	"github.com/egoist/mygo/ui"
)

var motionEntries2 = []Entry{
	{379, "Rings", "motion", "Rings expand and fade from a center.", demoRings},
	{380, "Morph", "motion", "A dot springs into a card and back.", demoMorph},
	{381, "SharedMove", "motion", "A row flies into its detail and back.", demoSharedMove},
	{382, "Shuffle", "motion", "Rows slide to a new order.", demoShuffle},
	{383, "ArriveList", "motion", "New rows slide in; the rest make room.", demoArriveList},
	{384, "Presence", "motion", "Content comes and goes with motion.", demoMotionPresence},
	{385, "Expand", "motion", "Height follows the content.", demoExpand},
	{386, "CrossFade", "motion", "Two contents fade into each other.", demoCrossFade},
	{387, "CardStack", "motion", "Drag the top card away to see the next.", demoCardStack},
	{388, "SwipeActions", "motion", "Drag a row aside to reveal actions.", demoSwipeActions},
	{389, "Throwable", "motion", "Throw the tile; it lands in a slot.", demoThrowable},
	{390, "PullRefresh", "motion", "Pull down past the line to refresh.", demoPullRefresh},
	{391, "Magnet", "motion", "Content leans toward the pointer.", demoMagnet},
	{392, "Follower", "motion", "A dot trails the pointer.", demoFollower},
	{393, "ScrollProgress", "motion", "A bar tracks the scroll.", demoScrollProgress},
	{394, "ScrollFade", "motion", "A header fades as the page scrolls.", demoScrollFade},
	{395, "ScrollReveal", "motion", "Rows rise in as they scroll into view.", demoScrollReveal},
	{396, "Sequence", "motion", "Title, body and action arrive in turn.", demoSequence},
	{397, "KeyPath", "motion", "A dot runs a path of keyframes.", demoKeyPath},
	{398, "ArcCounter", "motion", "An arc draws to a value and counts.", demoArcCounter},
	{399, "Sweep", "motion", "A marker sweeps under a phrase.", demoSweep},
	{400, "UnderlineGrow", "motion", "An underline grows on hover.", demoUnderlineGrow},
	{401, "Wipe", "motion", "Text revealed by a moving clip.", demoWipe},
	{402, "IconSwap", "motion", "An icon swaps with a transition.", demoIconSwap},
	{403, "DrawCheck", "motion", "A check draws itself in.", demoDrawCheck},
	{404, "TextShimmer", "motion", "A sheen passes over text.", demoTextShimmer},
	{405, "Magnify", "motion", "Icons swell near the pointer.", demoMagnify},
}

func demoRings(c *ui.Context, log func(string)) { motion.Rings(c, motion.RingsOptions{}) }

func demoMorph(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	open := boolState(root, "open", false)
	root.Children(func() {
		if motion.Morph(c, "Expand", open, func() { ui.Text(c, "Now a card.") }).Changed() {
			log(fmt.Sprintf("open: %v", *open))
		}
	})
}

func demoSharedMove(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320)
	open := intState(root, "open", -1)
	root.Children(func() {
		if motion.SharedMove(c, []string{"Inbox", "Drafts", "Sent"}, open, func(i int) {
			ui.Text(c, fmt.Sprintf("%d messages.", 3*i+4))
		}, motion.SharedMoveOptions{}).Changed() {
			log(fmt.Sprintf("open: %d", *open))
		}
	})
}

func demoShuffle(c *ui.Context, log func(string)) {
	if motion.Shuffle(c, []string{"Ada", "Lin", "Kai", "Mo"}, motion.ShuffleOptions{}).Changed() {
		log("shuffled")
	}
}

func demoArriveList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).Width(280)
	n := intState(root, "n", 2)
	root.Children(func() {
		if ui.Button(c, "New message").Clicked() {
			*n++
			log(fmt.Sprintf("messages %d", *n))
		}
		var items []string
		for i := *n; i > 0; i-- {
			items = append(items, fmt.Sprintf("Message %d", i))
		}
		motion.ArriveList(c, items)
	})
}

func demoMotionPresence(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).AlignItems(ui.Start)
	shown := boolState(root, "shown", true)
	root.Children(func() {
		if ui.Button(c, "Toggle").Clicked() {
			*shown = !*shown
			log(fmt.Sprintf("shown: %v", *shown))
		}
		motion.Presence(c, *shown, func() { tile(c, "Here now", 44) })
	})
}

func demoExpand(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).AlignItems(ui.Start).Width(280)
	more := boolState(root, "more", false)
	root.Children(func() {
		if ui.Button(c, "More").Clicked() {
			*more = !*more
			log(fmt.Sprintf("more: %v", *more))
		}
		motion.Expand(c, func() {
			ui.Text(c, "Every part moves.")
			if *more {
				ui.Text(c, "Height follows what the box holds, on a spring, with nothing to measure.")
			}
		})
	})
}

func demoCrossFade(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).AlignItems(ui.Start)
	b := boolState(root, "b", false)
	root.Children(func() {
		if ui.Button(c, "Swap").Clicked() {
			*b = !*b
			log(fmt.Sprintf("second: %v", *b))
		}
		motion.CrossFade(c, *b, func() { tile(c, "Day", 44) }, func() { tile(c, "Night", 44) })
	})
}

func demoCardStack(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).AlignItems(ui.Start).Width(260)
	top := intState(root, "top", 0)
	root.Children(func() {
		if motion.CardStack(c, 4, top, func(i int) { ui.Text(c, fmt.Sprintf("Card %d of 4", i+1)) }, motion.CardStackOptions{}).Changed() {
			log(fmt.Sprintf("top %d", *top))
		}
		if ui.Button(c, "Restack").Clicked() {
			*top = 0
			log("restacked")
		}
	})
}

func demoSwipeActions(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(340)
	chosen := intState(root, "chosen", -1)
	root.Children(func() {
		if motion.SwipeActions(c, []string{"Archive", "Delete"}, chosen, func() {
			ui.Row(c).Padding(14, 16).Children(func() { ui.Text(c, "Swipe this row left") })
		}).Changed() {
			log(fmt.Sprintf("action %d", *chosen))
		}
	})
}

func demoThrowable(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	slot := intState(root, "slot", 0)
	root.Children(func() {
		if motion.Throwable(c, "Go", slot, motion.ThrowableOptions{}).Changed() {
			log(fmt.Sprintf("slot %d", *slot))
		}
	})
}

func demoPullRefresh(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(300)
	busy := boolState(root, "busy", false)
	until := ui.Local(root, "until", func() time.Time { return time.Time{} })
	if *busy {
		if left := until.Sub(c.Now()); left > 0 {
			c.After(left)
		} else {
			*busy = false
			log("refreshed")
		}
	}
	root.Children(func() {
		if motion.PullRefresh(c, busy, motion.PullRefreshOptions{}, func() {
			ui.Column(c).Gap(8).Padding(12).Children(func() {
				for _, s := range []string{"Morning run", "Lunch with Lin", "Ship the build"} {
					ui.Text(c, s)
				}
			})
		}).Changed() {
			*until = c.Now().Add(1500 * time.Millisecond)
			log("refreshing")
		}
	})
}

func demoMagnet(c *ui.Context, log func(string)) {
	m := motion.Magnet(c, func() { tile(c, "Pull me", 44) })
	hover(m, log, "near", "away")
}

func demoFollower(c *ui.Context, log func(string)) {
	f := motion.Follower(c, motion.FollowerOptions{Height: 160}).Width(320)
	hover(f, log, "following", "resting")
}

// ends logs the scroll reaching its end or its top.
func ends(root ui.Element, st *ui.ScrollState, log func(string)) {
	at := ui.Local(root, "demo.ends", func() int { return 0 })
	now := 1
	switch {
	case st.Y <= 0:
		now = 0
	case st.MaxY > 0 && st.Y >= st.MaxY:
		now = 2
	}
	if now != *at {
		*at = now
		switch now {
		case 0:
			log("at the top")
		case 2:
			log("at the end")
		}
	}
}

func demoScrollProgress(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(8).Width(300)
	st := ui.Local(root, "st", func() *ui.ScrollState { return &ui.ScrollState{} })
	ends(root, *st, log)
	root.Children(func() {
		motion.ScrollProgress(c, "Reading", *st)
		ui.Scroll(c).TrackScroll(*st).Height(140).Gap(8).Children(func() {
			for i := range 12 {
				ui.Text(c, fmt.Sprintf("Paragraph %d", i+1))
			}
		})
	})
}

func demoScrollFade(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(300)
	st := ui.Local(root, "st", func() *ui.ScrollState { return &ui.ScrollState{} })
	ends(root, *st, log)
	root.Children(func() {
		ui.Scroll(c).TrackScroll(*st).Height(160).Gap(8).Children(func() {
			motion.ScrollFade(c, *st, motion.ScrollFadeOptions{}, func() { ui.Text(c, "Header").Bold().FontSize(24) })
			for i := range 12 {
				ui.Text(c, fmt.Sprintf("Line %d", i+1))
			}
		})
	})
}

func demoScrollReveal(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(300)
	st := ui.Local(root, "st", func() *ui.ScrollState { return &ui.ScrollState{} })
	ends(root, *st, log)
	root.Children(func() {
		ui.Scroll(c).TrackScroll(*st).Height(180).Children(func() {
			motion.ScrollReveal(c, 12, func(i int) { tile(c, fmt.Sprintf("Row %d", i+1), 40) })
		})
	})
}

func demoSequence(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).AlignItems(ui.Start).Width(320)
	run := ui.Local(root, "run", func() uint64 { return 0 })
	root.Children(func() {
		if ui.Button(c, "Replay").Clicked() {
			*run++
			log(fmt.Sprintf("run %d", *run))
		}
		if _, action := motion.Sequence(c, *run, "Welcome", "Three short steps to start.", motion.SequenceOptions{}); action.Clicked() {
			log("continue")
		}
	})
}

func demoKeyPath(c *ui.Context, log func(string)) { motion.KeyPath(c) }

func demoArcCounter(c *ui.Context, log func(string)) {
	root := ui.Row(c).Gap(16).AlignItems(ui.Center)
	v := numberState(root, "v", 72)
	root.Children(func() {
		motion.ArcCounter(c, "Score", *v, 100, motion.ArcCounterOptions{})
		if ui.Button(c, "New score").Clicked() {
			*v = math.Mod(*v+37, 101)
			log(fmt.Sprintf("score %.0f", *v))
		}
	})
}

func demoSweep(c *ui.Context, log func(string)) {
	trigger(c, log, "Mark", func(run uint64) { motion.Sweep(c, run, "The part that matters") })
}

func demoUnderlineGrow(c *ui.Context, log func(string)) {
	u := motion.UnderlineGrow(c, "Read the story")
	hover(u, log, "underline", "plain")
}

func demoWipe(c *ui.Context, log func(string)) {
	trigger(c, log, "Wipe again", func(run uint64) { motion.Wipe(c, run, "Revealed, left to right.") })
}

func demoIconSwap(c *ui.Context, log func(string)) {
	root := ui.Row(c).Gap(16).AlignItems(ui.Center)
	day := boolState(root, "day", true)
	root.Children(func() {
		name := "moon"
		if *day {
			name = "sun"
		}
		motion.IconSwap(c, name)
		if ui.Button(c, "Swap").Clicked() {
			*day = !*day
			log(fmt.Sprintf("day: %v", *day))
		}
	})
}

func demoDrawCheck(c *ui.Context, log func(string)) {
	trigger(c, log, "Save", func(run uint64) { motion.DrawCheck(c, "Saved", run) })
}

func demoTextShimmer(c *ui.Context, log func(string)) { motion.TextShimmer(c, "Loading your feed…") }

func demoMagnify(c *ui.Context, log func(string)) {
	root := ui.Column(c)
	picked := intState(root, "picked", -1)
	items := []motion.MagnifyItem{{Icon: "house", Label: "Home"}, {Icon: "mail", Label: "Mail"}, {Icon: "camera", Label: "Camera"},
		{Icon: "music", Label: "Music"}, {Icon: "settings", Label: "Settings"}}
	root.Children(func() {
		if motion.Magnify(c, items, picked).Changed() {
			log("picked " + items[*picked].Label)
		}
	})
}
