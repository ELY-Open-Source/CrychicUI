package gallery

import (
	"fmt"
	"image"
	stdcolor "image/color"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/egoist/mygo/ui"
)

var displayEntries = []Entry{
	{136, "Badge", "display", "Small semantic pill popping in.", demoBadge},
	{137, "DotBadge", "display", "Dot popping over a child's corner.", demoDotBadge},
	{138, "CountBadge", "display", "Rolling count with a 99+ cap.", demoCountBadge},
	{139, "StatusBadge", "display", "Status dot pulsing while live.", demoStatusBadge},
	{140, "Tag", "display", "Rounded tag with an optional icon.", demoTag},
	{141, "RemovableTag", "display", "Tag collapsing out when removed.", demoRemovableTag},
	{142, "Chip", "display", "Avatar and text pill lifting on hover.", demoChip},
	{143, "Avatar", "display", "Image or initials in a circle.", demoAvatar},
	{144, "AvatarGroup", "display", "Overlapping avatars spreading on hover.", demoAvatarGroup},
	{145, "AvatarStatus", "display", "Avatar with an animated presence dot.", demoAvatarStatus},
	{146, "Persona", "display", "Avatar, name and role fading in.", demoPersona},
	{147, "InfoCard", "display", "Title, body and footer card lifting on hover.", demoInfoCard},
	{148, "MediaCard", "display", "Card whose top image zooms on hover.", demoMediaCard},
	{149, "StatCard", "display", "Label, rolling value and delta card.", demoStatCard},
	{150, "Stat", "display", "Muted label over a rolling value.", demoStat},
	{151, "KPI", "display", "Big number with a trend arrow.", demoKPI},
	{152, "Delta", "display", "Signed change colored up or down.", demoDelta},
	{153, "ListItem", "display", "Icon, text and trailing slot with optional click.", demoListItem},
	{154, "DescriptionList", "display", "Key and value rows fading in.", demoDescriptionList},
	{155, "KeyValue", "display", "One key-value pair with cross-fading value.", demoKeyValue},
	{156, "Timeline", "display", "Vertical events staggering in.", demoTimeline},
	{157, "Activity", "display", "Actor, action, object and time fading in.", demoActivity},
	{158, "Ribbon", "display", "Corner ribbon sliding in.", demoRibbon},
	{159, "Indicator", "display", "Status dot animating to a new color.", demoIndicator},
	{160, "IconView", "display", "Lucide icon swapping with a content transition.", demoIconView},
	{161, "IconBox", "display", "Icon in a tinted square with animated color.", demoIconBox},
	{162, "Logo", "display", "Mark turning in beside a wordmark.", demoLogo},
	{163, "Placeholder", "display", "Dashed box with hover-animated dashes.", demoPlaceholder},
}

var sampleBitmap = ui.NewBitmap(func() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := range 32 {
		for x := range 32 {
			img.Set(x, y, stdcolor.RGBA{R: 80, G: 100, B: 220, A: 255})
		}
	}
	return img
}())

func demoBadge(c *ui.Context, log func(string)) {
	ui.Row(c).Gap(8).Children(func() {
		display.Badge(c, "Neutral", display.BadgeOptions{})
		display.Badge(c, "Accent", display.BadgeOptions{Variant: display.BadgeAccent})
		display.Badge(c, "Success", display.BadgeOptions{Variant: display.BadgeSuccess})
		display.Badge(c, "Warning", display.BadgeOptions{Variant: display.BadgeWarning})
		display.Badge(c, "Danger", display.BadgeOptions{Variant: display.BadgeDanger})
	})
}

func demoDotBadge(c *ui.Context, log func(string)) {
	display.DotBadge(c, display.DotBadgeOptions{}, func() { ui.Text(c, "Inbox") })
}

func demoCountBadge(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "count", 98, log)
		display.CountBadge(c, *n, func() { ui.Text(c, "Inbox") })
	})
}

func demoStatusBadge(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	live := ui.Local(root, "live", func() bool { return false })
	root.Children(func() {
		if ui.Toggle(c, live, "Live").Changed() {
			log(fmt.Sprintf("live: %v", *live))
		}
		display.StatusBadge(c, "Deploy", *live)
	})
}

func demoTag(c *ui.Context, log func(string)) {
	display.Tag(c, "Design", display.TagOptions{Icon: "tag"})
}

func demoRemovableTag(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	visible := ui.Local(root, "visible", func() bool { return true })
	root.Children(func() {
		if display.RemovableTag(c, "Docs", visible).Changed() {
			log("removed")
		}
	})
}

func demoChip(c *ui.Context, log func(string)) {
	display.Chip(c, "Ada", sampleBitmap)
}

func demoAvatar(c *ui.Context, log func(string)) {
	ui.Row(c).Gap(12).Children(func() {
		display.Avatar(c, "Crychic UI", nil)
		display.Avatar(c, "Ada", sampleBitmap)
	})
}

func demoAvatarGroup(c *ui.Context, log func(string)) {
	display.AvatarGroup(c, []string{"Ada Lovelace", "Grace Hopper", "Alan Turing", "Katherine Johnson"}, display.AvatarGroupOptions{})
}

func demoAvatarStatus(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	status := ui.Local(root, "status", func() display.PresenceStatus { return display.PresenceOffline })
	root.Children(func() {
		if ui.Button(c, "Cycle").Clicked() {
			*status = (*status + 1) % 4
			log(fmt.Sprintf("status: %d", *status))
		}
		display.AvatarStatus(c, "Ada", nil, *status)
	})
}

func demoPersona(c *ui.Context, log func(string)) {
	display.Persona(c, "Ada Lovelace", "Engineer", nil)
}

func demoInfoCard(c *ui.Context, log func(string)) {
	display.InfoCard(c, "Release", "Version two is ready.", func() { ui.Text(c, "Read the notes") })
}

func demoMediaCard(c *ui.Context, log func(string)) {
	display.MediaCard(c, "Photo", "A useful caption", sampleBitmap)
}

func demoStatCard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "revenue", 1234, log)
		display.StatCard(c, "Revenue", float64(*n)+0.5, 4.5)
	})
}

func demoStat(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "users", 12, log)
		display.Stat(c, "Users", float64(*n))
	})
}

func demoKPI(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "growth", 12000, log)
		display.KPI(c, "Growth", float64(*n), float64(*n-11900)/100)
	})
}

func demoDelta(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "delta", 8, log)
		display.Delta(c, float64(*n-10))
	})
}

func demoListItem(c *ui.Context, log func(string)) {
	if display.ListItem(c, "Inbox", "Three messages", display.ListItemOptions{Icon: "inbox", Clickable: true}, func() {
		ui.Text(c, "Now")
	}).Clicked() {
		log("open inbox")
	}
}

func demoDescriptionList(c *ui.Context, log func(string)) {
	display.DescriptionList(c, []display.KeyValuePair{{Key: "Name", Value: "Crychic"}, {Key: "Kind", Value: "UI"}})
}

func demoKeyValue(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	root.Children(func() {
		n := counter(c, root, "status", 1, log)
		display.KeyValue(c, "Status", fmt.Sprintf("Revision %d", *n))
	})
}

func demoTimeline(c *ui.Context, log func(string)) {
	display.Timeline(c, []string{"Created", "Reviewed", "Shipped"})
}

func demoActivity(c *ui.Context, log func(string)) {
	display.Activity(c, "Ada", "edited", "Spec", "now")
}

func demoRibbon(c *ui.Context, log func(string)) {
	display.Ribbon(c, "New", func() { ui.Box(c).Size(160, 100).Background(c.Theme().Surface) })
}

func demoIndicator(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	on := ui.Local(root, "on", func() bool { return false })
	root.Children(func() {
		if ui.Toggle(c, on, "Success").Changed() {
			log(fmt.Sprintf("on: %v", *on))
		}
		color := c.Theme().Accent
		if *on {
			color = c.Theme().Success
		}
		display.Indicator(c, color)
	})
}

func demoIconView(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	settings := ui.Local(root, "settings", func() bool { return false })
	root.Children(func() {
		if ui.Toggle(c, settings, "Settings").Changed() {
			log(fmt.Sprintf("settings: %v", *settings))
		}
		name := "house"
		if *settings {
			name = "settings"
		}
		display.IconView(c, name, display.IconViewOptions{})
	})
}

func demoIconBox(c *ui.Context, log func(string)) {
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
		display.IconBox(c, "bell", display.IconBoxOptions{Color: color})
	})
}

func demoLogo(c *ui.Context, log func(string)) {
	display.Logo(c, "sparkles", "CrychicUI")
}

func demoPlaceholder(c *ui.Context, log func(string)) {
	display.Placeholder(c, "Drop files", 240, 120)
}
