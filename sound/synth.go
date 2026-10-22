package sound

import "math"

const rate = 48000

type wave uint8

const (
	sine wave = iota
	triangle
	noise
)

// note is one tone sweeping from f0 to f1 Hz.
type note struct {
	at, dur float64
	f0, f1  float64
	gain    float64
	wave    wave
}

var scores = [cueCount][]note{
	Tap:     {{0, 0.025, 1800, 1200, 0.5, sine}},
	Toggle:  {{0, 0.03, 880, 880, 0.4, sine}, {0.035, 0.03, 1320, 1320, 0.4, sine}},
	Select:  {{0, 0.05, 1046, 1046, 0.4, sine}},
	Tick:    {{0, 0.008, 3000, 3000, 0.3, sine}},
	Open:    {{0, 0.09, 500, 900, 0.4, sine}},
	Close:   {{0, 0.09, 900, 500, 0.4, sine}},
	Pop:     {{0, 0.04, 400, 1200, 0.5, sine}},
	Drop:    {{0, 0.08, 700, 200, 0.5, sine}},
	Success: {{0, 0.12, 523, 523, 0.3, sine}, {0.07, 0.12, 659, 659, 0.3, sine}, {0.14, 0.18, 784, 784, 0.3, sine}},
	Warning: {{0, 0.09, 660, 660, 0.35, triangle}, {0.13, 0.09, 660, 660, 0.35, triangle}},
	Error:   {{0, 0.08, 220, 180, 0.4, triangle}, {0.1, 0.08, 220, 180, 0.4, triangle}},
	Notify:  {{0, 0.25, 784, 784, 0.3, sine}, {0.1, 0.3, 1046, 1046, 0.3, sine}},
	Send:    {{0, 0.12, 300, 1500, 0.35, triangle}},
	Delete:  {{0, 0.14, 600, 150, 0.4, triangle}},
	Swipe:   {{0, 0.07, 1200, 400, 0.25, sine}},
	Note:    tone(440),
	Drum:    drumScores[Kick],
}

var drumScores = [pieceCount][]note{
	Kick:  {{0, 0.3, 140, 45, 0.9, sine}},
	Snare: {{0, 0.16, 220, 180, 0.35, triangle}, {0, 0.16, 1000, 1000, 0.45, noise}},
	Hat:   {{0, 0.06, 1000, 1000, 0.35, noise}},
	Clap:  {{0, 0.02, 1000, 1000, 0.5, noise}, {0.025, 0.02, 1000, 1000, 0.5, noise}, {0.05, 0.12, 1000, 1000, 0.45, noise}},
	Tom:   {{0, 0.25, 200, 110, 0.7, sine}},
}

// tone is a struck note with a soft overtone.
func tone(hz float64) []note {
	out := []note{{0, 0.45, hz, hz, 0.3, sine}}
	if 2*hz < rate/4 {
		out = append(out, note{0, 0.3, 2 * hz, 2 * hz, 0.08, triangle})
	}
	return out
}

// render turns a score into mono samples.
func render(score []note) []float32 {
	if len(score) == 0 {
		panic("sound: empty score")
	}
	end := 0.0
	for _, n := range score {
		end = max(end, n.at+n.dur)
	}
	out := make([]float32, int(end*rate)+1)
	for _, n := range score {
		start, frames := int(n.at*rate), int(n.dur*rate)
		phase := 0.0
		for i := range frames {
			t := float64(i) / rate
			f := n.f0 + (n.f1-n.f0)*t/n.dur
			phase += 2 * math.Pi * f / rate
			out[start+i] += float32(n.gain * envelope(t, n.dur) * shape(n.wave, phase))
		}
	}
	return out
}

// envelope: 3 ms attack, decay, 2 ms fade.
func envelope(t, dur float64) float64 {
	e := math.Exp(-5 * t / dur)
	if t < 0.003 {
		e *= t / 0.003
	}
	if left := dur - t; left < 0.002 {
		e *= left / 0.002
	}
	return e
}

func shape(w wave, phase float64) float64 {
	switch w {
	case triangle:
		return 2 / math.Pi * math.Asin(math.Sin(phase))
	case noise:
		// A hash of the phase: fixed, so tests stay stable.
		v := math.Sin(phase*12.9898) * 43758.5453
		return 2*(v-math.Floor(v)) - 1
	}
	return math.Sin(phase)
}
