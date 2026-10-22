package gallery

import (
	"fmt"
	"image"
	"image/color"
	"strings"
	"sync"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/egoist/mygo/ui"
)

var choiceEntries2 = []Entry{
	{271, "Knob", "choice", "Dial dragged or keyed; its arc fills.", demoKnob},
	{272, "StarRating", "choice", "Stars with a hover preview; the chosen star bounces.", demoStarRating},
	{273, "EmojiRating", "choice", "Five faces; the chosen face grows.", demoEmojiRating},
	{274, "NPS", "choice", "0 to 10 scale colored by band.", demoNPS},
	{275, "Likert", "choice", "Disagree to agree scale.", demoLikert},
	{276, "Thumbs", "choice", "Up and down votes; the chosen thumb bounces.", demoThumbs},
	{277, "NumberStepper", "choice", "Minus and plus around a rolling value.", demoNumberStepper},
	{278, "Wheel", "choice", "Picker wheel that snaps and ticks per row.", demoWheel},
	{279, "YesNo", "choice", "Two joined buttons; the answer tints.", demoYesNo},
	{280, "ImageChoice", "choice", "Picture tiles; the chosen one rings.", demoImageChoice},
	{281, "ColorChoice", "choice", "Swatch radios; the chosen one rings.", demoColorChoice},
	{282, "SizeChoice", "choice", "Size boxes; unavailable ones struck.", demoSizeChoice},
	{283, "PlanChoice", "choice", "Plan radio cards with prices.", demoPlanChoice},
	{284, "ChoiceMatrix", "choice", "Rows of radios across shared columns.", demoChoiceMatrix},
	{285, "Ranking", "choice", "Rows dragged into order; rows slide.", demoRanking},
	{286, "IconToggle", "choice", "Pressable icon that swaps when toggled.", demoIconToggle},
	{287, "VisibilityToggle", "choice", "Eye that opens and closes.", demoVisibilityToggle},
	{288, "LockToggle", "choice", "Lock whose shackle lifts.", demoLockToggle},
	{289, "MuteToggle", "choice", "Speaker that mutes and unmutes.", demoMuteToggle},
	{290, "DayToggles", "choice", "Weekday letters toggled on and off.", demoDayToggles},
	{291, "SlotPicker", "choice", "Time slot chips; taken ones dim.", demoSlotPicker},
	{292, "OptionTiles", "choice", "Icon tiles choosing several; checks pop.", demoOptionTiles},
	{293, "TriCheck", "choice", "Check box whose mark morphs to a dash.", demoTriCheck},
	{294, "LabeledSwitch", "choice", "Switch whose inner text cross-fades.", demoLabeledSwitch},
	{295, "PickN", "choice", "Chips limited to n; extra choices shake.", demoPickN},
}

func intState(root ui.Element, key string, start int) *int {
	return ui.Local(root, key, func() int { return start })
}

func demoKnob(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	v := numberState(root, "v", 40)
	root.Children(func() {
		if choice.Knob(c, "Gain", v, 0, 100, choice.KnobOptions{Step: 5}).Changed() {
			log(fmt.Sprintf("gain: %.0f", *v))
		}
	})
}

func demoStarRating(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	v := intState(root, "v", 3)
	root.Children(func() {
		if choice.StarRating(c, "Rating", v, choice.StarRatingOptions{}).Changed() {
			log(fmt.Sprintf("stars: %d", *v))
		}
	})
}

func demoEmojiRating(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	v := intState(root, "v", 0)
	root.Children(func() {
		if choice.EmojiRating(c, "Mood", v, choice.EmojiRatingOptions{}).Changed() {
			log(fmt.Sprintf("mood: %d", *v))
		}
	})
}

func demoNPS(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	v := intState(root, "v", -1)
	root.Children(func() {
		if choice.NPS(c, "How likely are you to recommend us?", v, choice.NPSOptions{}).Changed() {
			log(fmt.Sprintf("score: %d", *v))
		}
	})
}

func demoLikert(c *ui.Context, log func(string)) {
	root := ui.Box(c).MaxWidth(560)
	v := intState(root, "v", -1)
	root.Children(func() {
		if choice.Likert(c, "The setup took little time.", v, choice.LikertOptions{}).Changed() {
			log(fmt.Sprintf("answer: %d", *v))
		}
	})
}

func demoThumbs(c *ui.Context, log func(string)) {
	root := ui.Row(c).Gap(12)
	v := intState(root, "v", 0)
	root.Children(func() {
		ui.Text(c, "Was this answer useful?")
		if choice.Thumbs(c, "Feedback", v, choice.ThumbsOptions{}).Changed() {
			log(fmt.Sprintf("vote: %d", *v))
		}
	})
}

func demoNumberStepper(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	v := numberState(root, "v", 2)
	root.Children(func() {
		if choice.NumberStepper(c, "Guests", v, 1, 12, 1, choice.NumberStepperOptions{}).Changed() {
			log(fmt.Sprintf("guests: %.0f", *v))
		}
	})
}

func demoWheel(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	v := intState(root, "v", 2)
	months := []string{"January", "February", "March", "April", "May", "June", "July", "August"}
	root.Children(func() {
		if choice.Wheel(c, "Month", months, v).Changed() {
			log("month: " + months[*v])
		}
	})
}

func demoYesNo(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	v := intState(root, "v", -1)
	root.Children(func() {
		if choice.YesNo(c, "Attending", v, choice.YesNoOptions{}).Changed() {
			log(fmt.Sprintf("answer: %d", *v))
		}
	})
}

// photos are three gradient pictures, made once.
var photos = sync.OnceValue(func() []*ui.Bitmap {
	var out []*ui.Bitmap
	for _, pair := range [][2]color.RGBA{{{246, 196, 120, 255}, {214, 120, 72, 255}}, {{120, 190, 230, 255}, {40, 90, 160, 255}}, {{150, 210, 140, 255}, {50, 120, 80, 255}}} {
		img := image.NewRGBA(image.Rect(0, 0, 96, 64))
		for y := range 64 {
			k := float64(y) / 63
			mix := func(a, b uint8) uint8 { return uint8(float64(a) + (float64(b)-float64(a))*k) }
			for x := range 96 {
				img.SetRGBA(x, y, color.RGBA{mix(pair[0].R, pair[1].R), mix(pair[0].G, pair[1].G), mix(pair[0].B, pair[1].B), 255})
			}
		}
		out = append(out, ui.NewBitmap(img))
	}
	return out
})

func demoImageChoice(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	v := intState(root, "v", 0)
	p := photos()
	items := []choice.ImageItem{{Label: "Dunes", Image: p[0]}, {Label: "Sea", Image: p[1]}, {Label: "Forest", Image: p[2]}}
	root.Children(func() {
		if choice.ImageChoice(c, "Cover", items, v).Changed() {
			log("cover: " + items[*v].Label)
		}
	})
}

func demoColorChoice(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	v := intState(root, "v", 0)
	t := c.Theme()
	swatches := []choice.Swatch{{Name: "Accent", Color: t.Accent}, {Name: "Success", Color: t.Success}, {Name: "Warning", Color: t.Warning}, {Name: "Danger", Color: t.Danger}, {Name: "Text", Color: t.Text}}
	root.Children(func() {
		if choice.ColorChoice(c, "Color", swatches, v).Changed() {
			log("color: " + swatches[*v].Name)
		}
	})
}

func demoSizeChoice(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	v := textState(root, "v", "M")
	root.Children(func() {
		if choice.SizeChoice(c, "Size", []string{"XS", "S", "M", "L", "XL"}, []string{"XS", "XL"}, v).Changed() {
			log("size: " + *v)
		}
	})
}

func demoPlanChoice(c *ui.Context, log func(string)) {
	root := ui.Box(c).MaxWidth(640)
	v := textState(root, "v", "Pro")
	plans := []choice.Plan{
		{Name: "Free", Price: "$0", Period: "/month", Note: "For trying things out."},
		{Name: "Pro", Price: "$12", Period: "/month", Note: "For one person shipping."},
		{Name: "Team", Price: "$40", Period: "/month", Note: "For five seats or more."},
	}
	root.Children(func() {
		if choice.PlanChoice(c, "Plan", plans, v).Changed() {
			log("plan: " + *v)
		}
	})
}

func demoChoiceMatrix(c *ui.Context, log func(string)) {
	root := ui.Box(c).MaxWidth(560)
	v := ui.Local(root, "v", func() []int { return []int{-1, -1, -1} })
	root.Children(func() {
		if choice.ChoiceMatrix(c, []string{"Speed", "Price", "Support"}, []string{"Poor", "Fair", "Good"}, v).Changed() {
			log(fmt.Sprintf("answers: %v", *v))
		}
	})
}

func demoRanking(c *ui.Context, log func(string)) {
	root := ui.Box(c).MaxWidth(360)
	v := listState(root, "v", "Speed", "Price", "Support", "Design")
	root.Children(func() {
		if choice.Ranking(c, "Priorities", v).Changed() {
			log("order: " + strings.Join(*v, ", "))
		}
	})
}

func demoIconToggle(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	v := boolState(root, "v", false)
	root.Children(func() {
		if choice.IconToggle(c, "Favorite", "heart", "heart-off", v).Changed() {
			log(fmt.Sprintf("favorite: %v", *v))
		}
	})
}

func demoVisibilityToggle(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	v := boolState(root, "v", false)
	root.Children(func() {
		if choice.VisibilityToggle(c, "Show password", v).Changed() {
			log(fmt.Sprintf("visible: %v", *v))
		}
	})
}

func demoLockToggle(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	v := boolState(root, "v", false)
	root.Children(func() {
		if choice.LockToggle(c, "Unlock layer", v).Changed() {
			log(fmt.Sprintf("unlocked: %v", *v))
		}
	})
}

func demoMuteToggle(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	v := boolState(root, "v", false)
	root.Children(func() {
		if choice.MuteToggle(c, "Mute", v).Changed() {
			log(fmt.Sprintf("muted: %v", *v))
		}
	})
}

func demoDayToggles(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	v := ui.Local(root, "v", func() [7]bool { return [7]bool{1: true, 3: true, 5: true} })
	root.Children(func() {
		if choice.DayToggles(c, "Repeat on", v, choice.DayTogglesOptions{}).Changed() {
			log(fmt.Sprintf("days: %v", *v))
		}
	})
}

func demoSlotPicker(c *ui.Context, log func(string)) {
	root := ui.Box(c).MaxWidth(420)
	v := textState(root, "v", "")
	slots := []string{"9:00", "9:30", "10:00", "10:30", "11:00", "13:00", "13:30", "14:00"}
	root.Children(func() {
		if choice.SlotPicker(c, "Time", slots, []string{"9:30", "11:00", "13:30"}, v).Changed() {
			log("slot: " + *v)
		}
	})
}

func demoOptionTiles(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	v := listState(root, "v", "Sleep")
	tiles := []choice.Tile{{Label: "Fitness", Icon: "dumbbell"}, {Label: "Sleep", Icon: "moon"}, {Label: "Focus", Icon: "target"}, {Label: "Food", Icon: "apple"}}
	root.Children(func() {
		if choice.OptionTiles(c, "Goals", tiles, v).Changed() {
			log("goals: " + strings.Join(*v, ", "))
		}
	})
}

func demoTriCheck(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	v := ui.Local(root, "v", func() choice.TriState { return choice.TriMixed })
	root.Children(func() {
		if choice.TriCheck(c, "Select all files", v).Changed() {
			log(fmt.Sprintf("state: %d", *v))
		}
	})
}

func demoLabeledSwitch(c *ui.Context, log func(string)) {
	root := ui.Box(c)
	v := boolState(root, "v", true)
	root.Children(func() {
		if choice.LabeledSwitch(c, "Autosave", v, choice.LabeledSwitchOptions{}).Changed() {
			log(fmt.Sprintf("autosave: %v", *v))
		}
	})
}

func demoPickN(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(8)
	v := listState(root, "v", "Go")
	root.Children(func() {
		ui.Text(c, "Pick up to three")
		o := choice.PickNOptions{Refused: func() { log("refused: three at most") }}
		if choice.PickN(c, "Languages", []string{"Go", "Rust", "Zig", "Swift", "Kotlin", "Elixir"}, 3, v, o).Changed() {
			log("picked: " + strings.Join(*v, ", "))
		}
	})
}
