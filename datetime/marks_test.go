package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// tally counts pixels near col in r, or everywhere.
func tally(tt *ui.Tester, r ui.Rect, col ui.Color) int {
	img := tt.Image()
	b := img.Bounds()
	if r == (ui.Rect{}) {
		r = ui.Rect{W: float32(b.Dx()), H: float32(b.Dy())}
	}
	n := 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			p := img.RGBAAt(x, y)
			if absInt(int(p.R)-int(col.R)) < 12 && absInt(int(p.G)-int(col.G)) < 12 && absInt(int(p.B)-int(col.B)) < 12 {
				n++
			}
		}
	}
	return n
}

func absInt(n int) int { return max(n, -n) }

func TestTimeAgo(t *testing.T) {
	var at time.Time
	o := TimeAgoOptions{}
	tt := ui.NewTester(func(c *ui.Context) {
		if !at.IsZero() {
			TimeAgo(c, at, o)
		}
	}, 300, 60)
	at = tt.Now().Add(-90 * time.Second)
	tt.Frame()
	if !tt.HasText("1 min ago") {
		t.Fatalf("90 s: %q", tt.Texts())
	}
	tt.Advance(31 * time.Second)
	if !tt.HasText("2 min ago") {
		t.Errorf("121 s: %q", tt.Texts())
	}
	for _, c := range []struct {
		ago  time.Duration
		want string
	}{{3 * time.Hour, "3 h ago"}, {50 * time.Hour, "2 d ago"}, {-time.Hour, "just now"}} {
		at = tt.Now().Add(-c.ago)
		tt.Frame()
		if !tt.HasText(c.want) {
			t.Errorf("%v ago: %q, want %q", c.ago, tt.Texts(), c.want)
		}
	}
	at = tt.Now().AddDate(0, 0, -30)
	tt.Frame()
	if tt.HasText("30 d ago") {
		t.Error("a month ago reads as days")
	}
	o.Minutes, at = "vor %d Min.", tt.Now().Add(-5*time.Minute)
	tt.Frame()
	if !tt.HasText("vor 5 Min.") {
		t.Errorf("copy: %q", tt.Texts())
	}
	settles(t, tt, "TimeAgo")
	panics(t, "zero time", func() { render(func(c *ui.Context) { TimeAgo(c, time.Time{}, TimeAgoOptions{}) }) })
}

func TestDatePill(t *testing.T) {
	day := time.Date(2026, 3, 4, 9, 0, 0, 0, time.Local)
	tt := ui.NewTester(func(c *ui.Context) { DatePill(c, day) }, 300, 60)
	if !tt.HasText("Wed") || !tt.HasText("4") {
		t.Fatalf("texts %q", tt.Texts())
	}
	day = day.AddDate(0, 0, 1)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("change did not cross-fade")
	}
	settles(t, tt, "DatePill")
	if !tt.HasText("Thu") || !tt.HasText("5") {
		t.Errorf("after change %q", tt.Texts())
	}
}

func TestDueChipTint(t *testing.T) {
	th := ui.LightTheme()
	var due time.Time
	tt := ui.NewTester(func(c *ui.Context) {
		if !due.IsZero() {
			DueChip(c, due, DueChipOptions{})
		}
	}, 300, 60)
	due = tt.Now().AddDate(0, 0, 10)
	tt.Frame()
	tt.Advance(time.Second)
	if tally(tt, ui.Rect{}, th.Warning) > 0 || tally(tt, ui.Rect{}, th.Danger) > 0 {
		t.Error("a far date is tinted")
	}
	due = tt.Now().AddDate(0, 0, 1)
	tt.Frame()
	tt.Advance(time.Second)
	if tally(tt, ui.Rect{}, th.Warning) < 50 {
		t.Error("tomorrow is not warning")
	}
	tt.Advance(49 * time.Hour)
	if tally(tt, ui.Rect{}, th.Danger) < 50 {
		t.Error("the day after did not turn overdue on its own")
	}
	panics(t, "negative Soon", func() { render(func(c *ui.Context) { DueChip(c, time.Now(), DueChipOptions{Soon: -1}) }) })
}

func TestDeadlineBadge(t *testing.T) {
	var due time.Time
	tt := ui.NewTester(func(c *ui.Context) {
		if !due.IsZero() {
			DeadlineBadge(c, due, DeadlineBadgeOptions{})
		}
	}, 300, 60)
	due = tt.Now().AddDate(0, 0, 3)
	tt.Frame()
	if !tt.HasText("In 3 d") {
		t.Errorf("ahead: %q", tt.Texts())
	}
	settles(t, tt, "badge ahead")
	due = tt.Now()
	tt.Frame()
	if !tt.HasText("Due today") {
		t.Errorf("today: %q", tt.Texts())
	}
	due = tt.Now().AddDate(0, 0, -1)
	tt.Frame()
	tt.Advance(2 * time.Second)
	if !tt.HasText("Overdue") || !tt.FramePending() {
		t.Errorf("overdue: %q, pulsing %v", tt.Texts(), tt.FramePending())
	}
}

func TestDateBadgeFlips(t *testing.T) {
	day := time.Date(2026, 3, 4, 0, 0, 0, 0, time.Local)
	tt := ui.NewTester(func(c *ui.Context) { DateBadge(c, day) }, 200, 120)
	settles(t, tt, "DateBadge")
	if !tt.HasText("Mar") || !tt.HasText("4") {
		t.Fatalf("texts %q", tt.Texts())
	}
	r := find(t, tt, "Mar 4, 2026")
	rest := tally(tt, r, ui.LightTheme().Accent)
	day = day.AddDate(0, 1, 1)
	tt.Frame()
	tt.Advance(80 * time.Millisecond)
	if mid := tally(tt, r, ui.LightTheme().Accent); mid >= rest {
		t.Errorf("leaf did not turn: header %d of %d pixels", mid, rest)
	}
	settles(t, tt, "flip")
	if !tt.HasText("Apr") || !tt.HasText("5") || tt.HasText("Mar") {
		t.Errorf("after flip %q", tt.Texts())
	}
}

func TestAgeTurnsOnBirthday(t *testing.T) {
	var birth time.Time
	tt := ui.NewTester(func(c *ui.Context) {
		if !birth.IsZero() {
			Age(c, birth, AgeOptions{})
		}
	}, 300, 80)
	next := tomorrow(tt.Now())
	birth = time.Date(next.Year()-30, next.Month(), next.Day(), 0, 0, 0, 0, time.Local)
	tt.Frame()
	if !tt.HasText("29") || !tt.HasText("years") {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt.Advance(next.Sub(tt.Now()) + time.Second)
	if !tt.HasText("30") {
		t.Errorf("birthday passed: %q", tt.Texts())
	}
	panics(t, "future birth", func() { render(func(c *ui.Context) { Age(c, time.Now().AddDate(1, 0, 0), AgeOptions{}) }) })
}

func TestMoonPhaseLit(t *testing.T) {
	th := ui.LightTheme()
	for _, c := range []struct {
		days        float64
		name        string
		left, right bool
	}{{0, "New moon", false, false}, {synodic / 4, "First quarter", false, true},
		{synodic / 2, "Full moon", true, true}, {synodic * 3 / 4, "Last quarter", true, false}} {
		day := newMoon.Add(time.Duration(c.days * 24 * float64(time.Hour)))
		tt := ui.NewTester(func(ctx *ui.Context) { MoonPhase(ctx, day, MoonPhaseOptions{}) }, 200, 200)
		settles(t, tt, c.name)
		r := find(t, tt, c.name)
		side := func(dx float32) bool {
			return tally(tt, ui.Rect{X: r.X + r.W/2 + dx - 2, Y: r.Y + r.H/2 - 2, W: 4, H: 4}, th.Warning) > 8
		}
		if side(-r.W/4) != c.left || side(r.W/4) != c.right {
			t.Errorf("%s: left lit %v, right lit %v", c.name, side(-r.W/4), side(r.W/4))
		}
	}
}

func TestSunArcSun(t *testing.T) {
	th := ui.LightTheme()
	var rise, set time.Time
	var name string
	tt := ui.NewTester(func(c *ui.Context) {
		if !rise.IsZero() {
			at := func(x time.Time) string { return i18n.FormatDate(locale(c), x, i18n.DateOptions{Skeleton: "jm"}) }
			name = at(rise) + " – " + at(set)
			SunArc(c, rise, set)
		}
	}, 400, 200)
	rise, set = tt.Now().Add(-time.Hour), tt.Now().Add(time.Hour)
	tt.Frame()
	r := find(t, tt, name)
	pad, h := th.Space(3), th.Space(24)
	rad := min(r.W/2-pad, h-2*pad)
	cx, cy := r.X+r.W/2, r.Y+h-pad
	if tally(tt, ui.Rect{X: cx - 2, Y: cy - rad - 2, W: 4, H: 4}, th.Warning) < 8 {
		t.Error("midday sun is not at the top")
	}
	rise, set = tt.Now().Add(time.Hour), tt.Now().Add(2*time.Hour)
	tt.Frame()
	settles(t, tt, "SunArc")
	if tally(tt, ui.Rect{X: cx - rad - 2, Y: cy - 2, W: 4, H: 4}, th.TextMuted) < 8 {
		t.Error("before sunrise the sun is not resting at the left")
	}
	panics(t, "sunset first", func() { render(func(c *ui.Context) { SunArc(c, time.Now(), time.Now()) }) })
}

func TestDayProgress(t *testing.T) {
	var want string
	tt := ui.NewTester(func(c *ui.Context) {
		now := c.Now()
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		pct := int64(now.Sub(start) * 100 / tomorrow(now).Sub(start))
		want = fmt.Sprintf("%s of the day", i18n.FormatNumber(locale(c), float64(pct)/100, i18n.NumberOptions{Percent: true}))
		DayProgress(c, "Today", DayProgressOptions{})
	}, 400, 100)
	if !tt.HasText(want) {
		t.Fatalf("texts %q, want %q", tt.Texts(), want)
	}
	first := want
	tt.Advance(15 * time.Minute)
	if want == first || !tt.HasText(want) {
		t.Errorf("15 min on: %q, want %q", tt.Texts(), want)
	}
}

func TestFiscalYearFills(t *testing.T) {
	accent := ui.LightTheme().Accent
	var start time.Time
	tt := ui.NewTester(func(c *ui.Context) {
		if !start.IsZero() {
			FiscalYear(c, "FY", start, FiscalYearOptions{})
		}
	}, 400, 100)
	lit := func(s time.Time) int {
		start = s
		tt.Frame()
		tt.Advance(2 * time.Second)
		return tally(tt, ui.Rect{}, accent)
	}
	month := monthOf(tt.Now())
	done, part, none := lit(month.AddDate(-2, 0, 0)), lit(month.AddDate(0, -4, 0)), lit(month.AddDate(0, 1, 0))
	if none != 0 || part <= done/4 || part >= done/2 {
		t.Errorf("accent pixels: ended %d, four months in %d, not begun %d", done, part, none)
	}
	panics(t, "mid-month start", func() {
		render(func(c *ui.Context) { FiscalYear(c, "FY", month.AddDate(0, 0, 3), FiscalYearOptions{}) })
	})
}

func TestCalendarHeat(t *testing.T) {
	counts := make([]int, 30)
	counts[10] = 4
	from := time.Date(2026, 3, 1, 0, 0, 0, 0, time.Local)
	lead, loc := 0, ""
	tt := ui.NewTester(func(c *ui.Context) {
		lead, loc = (int(from.Weekday())-int(c.Region().FirstWeekday)+7)%7, locale(c)
		CalendarHeat(c, "Commits", from, counts, CalendarHeatOptions{})
	}, 400, 200)
	r := find(t, tt, "Commits")
	tt.Advance(50 * time.Millisecond)
	if !tt.FramePending() {
		t.Error("cells did not sweep in")
	}
	settles(t, tt, "CalendarHeat")
	th := ui.LightTheme()
	cell := th.Space(2.5) + th.Space(0.5)
	at := func(i int) (float32, float32) {
		return r.X + float32((lead+i)/7)*cell + 4, r.Y + float32((lead+i)%7)*cell + 4
	}
	x, y := at(10)
	if tally(tt, ui.Rect{X: x - 1, Y: y - 1, W: 2, H: 2}, th.Accent) == 0 {
		t.Error("the busiest day is not full accent")
	}
	tt.Move(x, y)
	tt.Frame()
	want := fmt.Sprintf("4 on %s", i18n.FormatDate(loc, from.AddDate(0, 0, 10), i18n.DateOptions{Skeleton: "yMMMd"}))
	if !tt.HasText(want) {
		t.Errorf("hover caption %q, want %q", tt.Texts(), want)
	}
	panics(t, "negative count", func() {
		render(func(c *ui.Context) { CalendarHeat(c, "x", from, []int{-1}, CalendarHeatOptions{}) })
	})
}
