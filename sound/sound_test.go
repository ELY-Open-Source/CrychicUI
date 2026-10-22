package sound

import (
	"bytes"
	"errors"
	"log/slog"
	"math"
	"slices"
	"strings"
	"testing"
)

func resetRecent(t *testing.T) {
	mu.Lock()
	recent = nil
	mu.Unlock()
	t.Cleanup(func() {
		mu.Lock()
		recent = nil
		mu.Unlock()
	})
}

func TestTestBinaryStartsMuted(t *testing.T) {
	if Enabled() {
		t.Fatal("sound enabled in a test binary")
	}
}

func TestCueNames(t *testing.T) {
	seen := map[string]bool{}
	for _, c := range Cues() {
		s := c.String()
		if seen[s] {
			t.Errorf("duplicate name %q", s)
		}
		seen[s] = true
	}
	if len(seen) != 17 {
		t.Errorf("cues = %d, want 17", len(seen))
	}
	if got := Cue(99).String(); got != "cue(99)" {
		t.Errorf("unknown cue name = %q", got)
	}
}

// bounded checks a sound's length, peak and clean ends.
func bounded(t *testing.T, name string, s []float32) {
	t.Helper()
	if d := float64(len(s)) / rate; d < 0.005 || d > 0.5 {
		t.Errorf("%s lasts %.3f s", name, d)
	}
	peak := 0.0
	for _, v := range s {
		peak = max(peak, math.Abs(float64(v)))
	}
	if peak < 0.05 || peak > 1 {
		t.Errorf("%s peak %.3f", name, peak)
	}
	if a, z := math.Abs(float64(s[0])), math.Abs(float64(s[len(s)-1])); a > 0.01 || z > 0.01 {
		t.Errorf("%s clicks: first %.3f last %.3f", name, a, z)
	}
}

func TestSynthBounds(t *testing.T) {
	for _, c := range Cues() {
		bounded(t, c.String(), render(scores[c]))
	}
}

func TestCuesDiffer(t *testing.T) {
	cues := Cues()
	for i, a := range cues {
		for _, b := range cues[i+1:] {
			if slices.Equal(render(scores[a]), render(scores[b])) {
				t.Errorf("%v and %v sound the same", a, b)
			}
		}
	}
}

func TestRecentKeepsLast32(t *testing.T) {
	resetRecent(t)
	var played []Cue
	for i := range 40 {
		c := Cue(i % int(cueCount))
		Play(c)
		played = append(played, c)
	}
	if got, want := Recent(), played[8:]; !slices.Equal(got, want) {
		t.Errorf("recent = %v, want %v", got, want)
	}
}

func TestRecentIsACopy(t *testing.T) {
	resetRecent(t)
	Play(Tap)
	r := Recent()
	r[0] = Error
	if Recent()[0] != Tap {
		t.Error("Recent exposed its buffer")
	}
}

func TestSetEnabled(t *testing.T) {
	t.Cleanup(func() { SetEnabled(false) })
	SetEnabled(true)
	if !Enabled() {
		t.Error("SetEnabled(true) ignored")
	}
	SetEnabled(false)
	if Enabled() {
		t.Error("SetEnabled(false) ignored")
	}
}

func TestSetVolume(t *testing.T) {
	old := Volume()
	t.Cleanup(func() { SetVolume(old) })
	SetVolume(0.8)
	if Volume() != 0.8 {
		t.Errorf("volume = %v", Volume())
	}
	for _, v := range []float64{-0.1, 1.1, math.NaN()} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("SetVolume(%v) did not panic", v)
				}
			}()
			SetVolume(v)
		}()
	}
}

func TestPlayPanicsOnUnknownCue(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Play(cueCount) did not panic")
		}
	}()
	Play(cueCount)
}

func TestNoteErrLogsChanges(t *testing.T) {
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(old) })
	down, other := errors.New("no device"), errors.New("device lost")
	steps := []struct {
		err  error
		want string
	}{
		{nil, ""},
		{down, "level=WARN msg=\"crychic: sound output unavailable\" err=\"no device\""},
		{down, ""},
		{other, "level=WARN msg=\"crychic: sound output unavailable\" err=\"device lost\""},
		{nil, "level=INFO msg=\"crychic: sound output back\""},
		{nil, ""},
	}
	var last error
	for i, s := range steps {
		buf.Reset()
		last = noteErr(last, s.err)
		got := strings.TrimSpace(buf.String())
		if s.want == "" && got != "" || s.want != "" && !strings.Contains(got, s.want) {
			t.Errorf("step %d logged %q, want %q", i, got, s.want)
		}
		if last != s.err {
			t.Errorf("step %d kept %v, want %v", i, last, s.err)
		}
	}
}
