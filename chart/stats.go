package chart

import (
	"math"
	"slices"

	mc "github.com/egoist/mygo/ui/chart"
)

// quartiles are the min, three quartiles and max.
func quartiles(vs []float64) [5]float64 {
	s := slices.Sorted(slices.Values(vs))
	at := func(p float64) float64 {
		f := p * float64(len(s)-1)
		i := int(f)
		if i+1 >= len(s) {
			return s[len(s)-1]
		}
		return s[i] + (s[i+1]-s[i])*(f-float64(i))
	}
	return [5]float64{s[0], at(0.25), at(0.5), at(0.75), s[len(s)-1]}
}

// density is a Gaussian estimate; bandwidth spans a step.
func density(vs, at []float64) []float64 {
	mean, sd := 0.0, 0.0
	for _, v := range vs {
		mean += v / float64(len(vs))
	}
	for _, v := range vs {
		sd += (v - mean) * (v - mean) / float64(len(vs))
	}
	step := (at[len(at)-1] - at[0]) / float64(len(at)-1)
	bw := max(1.06*math.Sqrt(sd)*math.Pow(float64(len(vs)), -0.2), step)
	out := make([]float64, len(at))
	for i, x := range at {
		for _, v := range vs {
			z := (x - v) / bw
			out[i] += math.Exp(-z * z / 2)
		}
	}
	return out
}

// squarify tiles the unit square by values, rows kept squarish.
func squarify(vs []float64) [][2]mc.Span {
	total := 0.0
	for _, v := range vs {
		total += max(v, 0)
	}
	out := make([][2]mc.Span, len(vs))
	x0, y0, w, h := 0.0, 0.0, 1.0, 1.0
	worst := func(row []float64, side float64) float64 {
		sum, hi, lo := 0.0, 0.0, math.Inf(1)
		for _, a := range row {
			sum, hi, lo = sum+a, max(hi, a), min(lo, a)
		}
		return max(side*side*hi/(sum*sum), sum*sum/(side*side*lo))
	}
	i := 0
	scale := 1 / total
	for i < len(vs) {
		side := min(w, h)
		row := []float64{max(vs[i], 1e-9) * scale}
		j := i + 1
		for j < len(vs) {
			next := append(slices.Clone(row), max(vs[j], 1e-9)*scale)
			if worst(next, side) > worst(row, side) {
				break
			}
			row, j = next, j+1
		}
		sum := 0.0
		for _, a := range row {
			sum += a
		}
		if w >= h {
			cw := sum / h
			cy := y0
			for k, a := range row {
				out[i+k] = [2]mc.Span{{Lo: float32(x0), Hi: float32(x0 + cw)}, {Lo: float32(cy), Hi: float32(cy + a/cw)}}
				cy += a / cw
			}
			x0, w = x0+cw, w-cw
		} else {
			ch := sum / w
			cx := x0
			for k, a := range row {
				out[i+k] = [2]mc.Span{{Lo: float32(cx), Hi: float32(cx + a/ch)}, {Lo: float32(y0 + h - ch), Hi: float32(y0 + h)}}
				cx += a / ch
			}
			h -= ch
		}
		i = j
	}
	return out
}
