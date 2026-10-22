package chart

import "github.com/egoist/mygo/ui"

// PieOptions tune a Pie.
type PieOptions struct {
	// Size is the pie's width; zero gives the default.
	Size float32
	// Format prints a value; "%g" when empty.
	Format string
}

// Pie slices a circle; the chosen slice lifts.
func Pie(c *ui.Context, label string, names []string, vs []float64, at *int, o PieOptions) ui.Element {
	return pies(c, label, names, vs, at, o.Size, 0, o.Format, "Pie", nil)
}
