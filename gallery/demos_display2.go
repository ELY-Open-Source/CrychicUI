package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/egoist/mygo/ui"
)

var displayEntries2 = []Entry{
	{164, "Feature", "display", "Icon box, title and text rising in.", demoFeature},
	{165, "Metric", "display", "Rolling value with a unit.", demoMetric},
	{166, "Stars", "display", "Read-only rating with halves.", demoStars},
	{167, "Price", "display", "New price rolls beside an old struck price.", demoPrice},
	{168, "Version", "display", "Version pill cross-fading on change.", demoVersion},
	{169, "Legend", "display", "Swatches with labels fading in.", demoLegend},
	{170, "Swatch", "display", "Color sample with a label and animated color.", demoSwatch},
	{171, "StepList", "display", "Done, current and next steps.", demoStepList},
	{172, "Checklist", "display", "Read-only checks growing in.", demoChecklist},
	{173, "Picture", "display", "Rounded image fading in with a Fit option.", demoPicture},
	{174, "Thumbnail", "display", "Small image with an icon fallback.", demoThumbnail},
	{175, "Figure", "display", "Image and caption fading in.", demoFigure},
	{176, "Locale", "display", "Language pill cross-fading.", demoLocale},
	{177, "Presence", "display", "Status dot and label with animated color.", demoPresence},
	{178, "OrSeparator", "display", "Label between two growing lines.", demoOrSeparator},
	{179, "Watermark", "display", "Faint tiled text over content.", demoWatermark},
	{180, "AccentCard", "display", "Card with a growing accent border.", demoAccentCard},
	{181, "Comparison", "display", "Before and after values with a rolling result.", demoComparison},
	{182, "Score", "display", "Number in a ring drawing to value.", demoScore},
	{183, "Level", "display", "Level number and bar to the next.", demoLevel},
	{184, "Streak", "display", "Flickering flame and rolling day count.", demoStreak},
	{185, "Achievement", "display", "Icon and name with an unlock shine.", demoAchievement},
	{186, "Counter", "display", "Big number counting toward a target.", demoCounter},
	{187, "LiveClock", "display", "Current time with ticking seconds.", demoLiveClock},
	{188, "CountdownTo", "display", "Days to seconds remaining with rolling digits.", demoCountdownTo},
	{189, "Weather", "display", "Condition icon and temperature.", demoWeather},
	{190, "IDCard", "display", "Avatar, name and fields rising on a card.", demoIDCard},
}

func demoFeature(c *ui.Context, log func(string)) {
	display.Feature(c, "sparkles", "Fast", "Motion and sound included.")
}

func demoMetric(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "latency", 12, log)
		display.Metric(c, float64(*n), "ms")
	})
}

func demoStars(c *ui.Context, log func(string)) {
	display.Stars(c, 3.5, display.StarsOptions{})
}

func demoPrice(c *ui.Context, log func(string)) {
	display.Price(c, 10, 12, "USD")
}

func demoVersion(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "version", 1, log)
		display.Version(c, fmt.Sprintf("1.%d.0", *n))
	})
}

func demoLegend(c *ui.Context, log func(string)) {
	t := c.Theme()
	display.Legend(c, []display.LegendItem{{Label: "Now", Color: t.Accent}, {Label: "Done", Color: t.Success}})
}

func demoSwatch(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	success := ui.Local(root, "success", func() bool { return false })
	root.Children(func() {
		if ui.Toggle(c, success, "Success").Changed() {
			log(fmt.Sprintf("success: %v", *success))
		}
		color := c.Theme().Accent
		if *success {
			color = c.Theme().Success
		}
		display.Swatch(c, color)
	})
}

func demoStepList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	current := ui.Local(root, "current", func() int { return 0 })
	steps := []string{"Draft", "Review", "Ship"}
	root.Children(func() {
		if ui.Button(c, "Next").Clicked() {
			*current = min(len(steps)-1, *current+1)
			log(fmt.Sprintf("current: %d", *current))
		}
		display.StepList(c, steps, *current)
	})
}

func demoChecklist(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	done := ui.Local(root, "done", func() bool { return false })
	root.Children(func() {
		if ui.Toggle(c, done, "Write").Changed() {
			log(fmt.Sprintf("done: %v", *done))
		}
		display.Checklist(c, []display.ChecklistItem{{Label: "Write", Done: *done}, {Label: "Review", Done: true}})
	})
}

func demoPicture(c *ui.Context, log func(string)) {
	display.Picture(c, sampleBitmap, display.PictureOptions{Fit: ui.Cover})
}

func demoThumbnail(c *ui.Context, log func(string)) {
	ui.Row(c).Gap(12).Children(func() {
		display.Thumbnail(c, sampleBitmap, "image")
		display.Thumbnail(c, nil, "image")
	})
}

func demoFigure(c *ui.Context, log func(string)) {
	display.Figure(c, sampleBitmap, "A generated square")
}

func demoLocale(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	chinese := ui.Local(root, "chinese", func() bool { return false })
	root.Children(func() {
		if ui.Toggle(c, chinese, "中文").Changed() {
			log(fmt.Sprintf("chinese: %v", *chinese))
		}
		label := "EN"
		if *chinese {
			label = "中文"
		}
		display.Locale(c, label)
	})
}

func demoPresence(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	status := ui.Local(root, "status", func() display.PresenceStatus { return display.PresenceOffline })
	labels := []string{"Offline", "Online", "Away", "Busy"}
	root.Children(func() {
		if ui.Button(c, "Cycle").Clicked() {
			*status = (*status + 1) % 4
			log("presence: " + labels[*status])
		}
		display.Presence(c, *status, labels[*status])
	})
}

func demoOrSeparator(c *ui.Context, log func(string)) {
	display.OrSeparator(c, "OR")
}

func demoWatermark(c *ui.Context, log func(string)) {
	display.Watermark(c, "DRAFT", func() { ui.Box(c).Size(220, 120).Background(c.Theme().Surface) })
}

func demoAccentCard(c *ui.Context, log func(string)) {
	display.AccentCard(c, func() { ui.Text(c, "An important card") })
}

func demoComparison(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "after", 12, log)
		display.Comparison(c, 10, float64(*n))
	})
}

func demoScore(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	value := ui.Local(root, "score", func() int { return 42 })
	root.Children(func() {
		if ui.Button(c, "+10").Clicked() {
			*value = (*value + 10) % 101
			log(fmt.Sprintf("score: %d", *value))
		}
		display.Score(c, float64(*value), display.ScoreOptions{})
	})
}

func demoLevel(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	progress := ui.Local(root, "progress", func() float64 { return 0.25 })
	root.Children(func() {
		if ui.Button(c, "Add progress").Clicked() {
			*progress += 0.25
			if *progress > 1 {
				*progress = 0.25
			}
			log(fmt.Sprintf("progress: %.2f", *progress))
		}
		display.Level(c, 2, *progress, display.LevelOptions{})
	})
}

func demoStreak(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "days", 4, log)
		display.Streak(c, *n)
	})
}

func demoAchievement(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	unlocked := ui.Local(root, "unlocked", func() bool { return false })
	root.Children(func() {
		if ui.Toggle(c, unlocked, "Unlocked").Changed() {
			log(fmt.Sprintf("unlocked: %v", *unlocked))
		}
		display.Achievement(c, "First ship", "trophy", *unlocked)
	})
}

func demoCounter(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	target := ui.Local(root, "target", func() int { return 100 })
	root.Children(func() {
		if ui.Button(c, "Add 25").Clicked() {
			*target += 25
			log(fmt.Sprintf("target: %d", *target))
		}
		display.Counter(c, *target)
	})
}

func demoLiveClock(c *ui.Context, log func(string)) {
	display.LiveClock(c)
}

func demoCountdownTo(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	target := ui.Local(root, "target", func() time.Time { return c.Now().Add(90 * time.Second) })
	root.Children(func() { display.CountdownTo(c, *target) })
}

func demoWeather(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	condition := ui.Local(root, "condition", func() display.WeatherCondition { return display.WeatherSunny })
	root.Children(func() {
		if ui.Button(c, "Cycle").Clicked() {
			*condition = (*condition + 1) % 4
			log(fmt.Sprintf("condition: %d", *condition))
		}
		display.Weather(c, *condition, 22.5, "°C")
	})
}

func demoIDCard(c *ui.Context, log func(string)) {
	display.IDCard(c, "Ada Lovelace", nil, []display.KeyValuePair{{Key: "Role", Value: "Engineer"}, {Key: "Team", Value: "Interface"}})
}
