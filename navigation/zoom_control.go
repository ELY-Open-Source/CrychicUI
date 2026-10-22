package navigation

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/motion"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ZoomControlOptions tune a ZoomControl; zero values take defaults.
type ZoomControlOptions struct {
	Min, Max, Step, FitTo float64
	Out, In, Fit          string
}

// ZoomControl is minus, a rolling percent, plus and fit.
func ZoomControl(c *ui.Context, label string, zoom *float64, o ZoomControlOptions) ui.Element {
	if o.Min == 0 {
		o.Min = 0.25
	}
	if o.Max == 0 {
		o.Max = 4
	}
	if o.Step == 0 {
		o.Step = 0.25
	}
	if o.FitTo == 0 {
		o.FitTo = 1
	}
	for _, v := range []float64{o.Min, o.Max, o.Step, o.FitTo} {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			panic(fmt.Sprintf("navigation: ZoomControl %q options must be finite", label))
		}
	}
	if label == "" || zoom == nil || o.Min <= 0 || o.Min >= o.Max || o.Step <= 0 || o.FitTo < o.Min || o.FitTo > o.Max ||
		math.IsNaN(*zoom) || *zoom < o.Min || *zoom > o.Max {
		panic(fmt.Sprintf("navigation: ZoomControl %q needs a zoom inside a valid range", label))
	}
	def(&o.Out, "Zoom out")
	def(&o.In, "Zoom in")
	def(&o.Fit, "Fit")
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1)).Padding(t.Space(0.5)).AlignItems(ui.Center).Radius(t.Radius).
		Background(t.Surface).Border(1, t.Border).FocusGroup(ui.Horizontal).Label(label)
	z := *zoom
	set := func(v float64) {
		*zoom = v
		root.SetChanged()
		kit.Log("ZoomControl", "zoomed", "percent", math.Round(*zoom*100))
	}
	return root.Children(func() {
		if arrowButton(c, o.Out, "minus", z <= o.Min, sound.Tick) {
			set(max(o.Min, math.Round((z-o.Step)/o.Step)*o.Step))
		}
		ui.Row(c).Width(t.Space(14)).Justify(ui.Center).AlignItems(ui.Center).Children(func() {
			motion.RollingDigits(c, int64(math.Round(z*100))).FontSize(t.Rem(0.9))
			ui.Text(c, "%")
		})
		if arrowButton(c, o.In, "plus", z >= o.Max, sound.Tick) {
			set(min(o.Max, math.Round((z+o.Step)/o.Step)*o.Step))
		}
		if arrowButton(c, o.Fit, "scan", z == o.FitTo, sound.Tick) {
			set(o.FitTo)
		}
	})
}
