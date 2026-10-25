package health

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Session is one workout; zero Kilometers or HeartRate hide them.
type Session struct {
	Kind, Icon        string
	When              time.Time
	Minutes, Calories int
	Kilometers        float64
	HeartRate         int
}

// WorkoutOptions word the stats; Units are min to bpm.
type WorkoutOptions struct {
	Date, Time, Distance, Energy, Heart, Pace string
	Units                                     [6]string
}

// Workout shows a finished workout with its stats.
func Workout(c *ui.Context, s Session, o WorkoutOptions) ui.Element {
	if strings.TrimSpace(s.Kind) == "" || s.Icon == "" || s.When.IsZero() || s.Minutes < 1 || s.Calories < 0 || s.HeartRate < 0 ||
		!(s.Kilometers == 0 || s.Kilometers >= 0.01 && s.Kilometers < 1e4) || s.Minutes > 24*60 {
		panic(fmt.Sprintf("health: Workout %q needs an icon, a time, 1 to 1440 minutes and stats of 0 or more", s.Kind))
	}
	def(&o.Date, "Mon, Jan 2 · 15:04")
	def(&o.Time, "Time")
	def(&o.Distance, "Distance")
	def(&o.Energy, "Energy")
	def(&o.Heart, "Avg heart rate")
	def(&o.Pace, "Pace")
	for i, u := range [...]string{"min", "h", "km", "/km", "kcal", "bpm"} {
		def(&o.Units[i], u)
	}
	t := c.Theme()
	type stat struct{ name, value, unit, mark string }
	stats := []stat{{o.Time, fmt.Sprint(s.Minutes), o.Units[0], "timer"}}
	if s.Minutes >= 60 {
		stats[0] = stat{o.Time, fmt.Sprintf("%d:%02d", s.Minutes/60, s.Minutes%60), o.Units[1], "timer"}
	}
	if s.Kilometers > 0 {
		pace := int(math.Round(float64(s.Minutes) * 60 / s.Kilometers))
		stats = append(stats, stat{o.Distance, fmt.Sprintf("%.2f", s.Kilometers), o.Units[2], "route"},
			stat{o.Pace, fmt.Sprintf("%d'%02d\"", pace/60, pace%60), o.Units[3], "gauge"})
	}
	stats = append(stats, stat{o.Energy, number(c, s.Calories), o.Units[4], "flame"})
	if s.HeartRate > 0 {
		stats = append(stats, stat{o.Heart, fmt.Sprint(s.HeartRate), o.Units[5], "heart"})
	}
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			kit.Arrive(tile(c, s.Icon, t.Success, t.Space(11)), ui.Pose{Scale: 0.6, Rotate: -15, Fade: 1}, 0)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, s.Kind).FontSize(t.Rem(0.95)).Bold().FillWidth().MinWidth(0).SingleLine()
				kit.Muted(c, s.When.Format(o.Date)).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
			})
		})
		grid := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(2))
		grid.Children(func() {
			for i, st := range stats {
				cell := ui.Column(c).MinWidth(0).Gap(t.Space(1)).Padding(t.Space(2.5)).Radius(t.Radius).Background(t.SurfaceHover.Alpha(0.6))
				kit.Stagger(c, grid, cell, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
						glyph(c, st.mark, 14)
						ui.Text(c, st.name).FontSize(t.Rem(0.68)).Grow(1).MinWidth(0).SingleLine()
					})
					ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.End).Children(func() {
						kit.Numeric(c, st.value).FontSize(t.Rem(1.1)).Bold().MinWidth(0).SingleLine().ContentTransition(ui.ContentTransition{Numeric: true})
						kit.Muted(c, st.unit).FontSize(t.Rem(0.68)).Shrink(0).MaxWidth(t.Space(12)).SingleLine()
					})
				})
			}
		})
	})
	return root
}
