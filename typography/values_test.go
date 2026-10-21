package typography

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestNumber(t *testing.T) {
	value := 1234.5
	tt := ui.NewTester(func(c *ui.Context) { Number(c, value) }, 300, 200)
	find(t, tt, "1,234.5")
	value = 12345.5
	tt.Frame()
	settles(t, tt, "Number")
	find(t, tt, "12,345.5")
	for _, value := range []float64{math.NaN(), math.Inf(1)} {
		panics(t, "Number", func() { ui.Render(func(c *ui.Context) { Number(c, value) }, 10, 10, 1) })
	}
}

func TestMoney(t *testing.T) {
	amount := 1234.5
	tt := ui.NewTester(func(c *ui.Context) { Money(c, amount, "USD") }, 300, 200)
	find(t, tt, "USD 1,234.50")
	amount = 5
	tt.Frame()
	settles(t, tt, "Money")
	find(t, tt, "USD 5.00")
	panics(t, "Money currency", func() { ui.Render(func(c *ui.Context) { Money(c, 1, "") }, 10, 10, 1) })
	panics(t, "Money NaN", func() { ui.Render(func(c *ui.Context) { Money(c, math.NaN(), "USD") }, 10, 10, 1) })
}

func TestPercent(t *testing.T) {
	value := 0.25
	tt := ui.NewTester(func(c *ui.Context) { Percent(c, value) }, 300, 200)
	find(t, tt, "25%")
	value = -0.5
	tt.Frame()
	settles(t, tt, "Percent")
	find(t, tt, "-50%")
	panics(t, "Percent NaN", func() { ui.Render(func(c *ui.Context) { Percent(c, math.NaN()) }, 10, 10, 1) })
}

func TestBytes(t *testing.T) {
	n := int64(1500)
	tt := ui.NewTester(func(c *ui.Context) { Bytes(c, n) }, 300, 200)
	find(t, tt, "1.5 KB")
	n = 2000000
	tt.Frame()
	settles(t, tt, "Bytes")
	find(t, tt, "2 MB")
	panics(t, "Bytes negative", func() { ui.Render(func(c *ui.Context) { Bytes(c, -1) }, 10, 10, 1) })
}

func TestDuration(t *testing.T) {
	d := time.Hour + 4*time.Minute
	tt := ui.NewTester(func(c *ui.Context) { Duration(c, d) }, 300, 200)
	find(t, tt, "1h 4m")
	d = 90 * time.Second
	tt.Frame()
	settles(t, tt, "Duration")
	find(t, tt, "1m 30s")
	panics(t, "Duration negative", func() { ui.Render(func(c *ui.Context) { Duration(c, -1) }, 10, 10, 1) })
}
