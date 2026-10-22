// Package sound plays the cues components make.
package sound

import (
	"fmt"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/egoist/mygo/audio"
)

// Cue names one sound.
type Cue uint8

const (
	Tap Cue = iota
	Toggle
	Select
	Tick
	Open
	Close
	Pop
	Drop
	Success
	Warning
	Error
	Notify
	Send
	Delete
	Swipe
	Note
	Drum
	cueCount
)

var names = [cueCount]string{"tap", "toggle", "select", "tick", "open", "close", "pop", "drop", "success",
	"warning", "error", "notify", "send", "delete", "swipe", "note", "drum"}

// Piece names one synthesized drum.
type Piece uint8

const (
	Kick Piece = iota
	Snare
	Hat
	Clap
	Tom
	pieceCount
)

var pieceNames = [pieceCount]string{"kick", "snare", "hat", "clap", "tom"}

func (p Piece) String() string {
	if p >= pieceCount {
		return fmt.Sprintf("piece(%d)", uint8(p))
	}
	return pieceNames[p]
}

// Pieces lists every drum piece in order.
func Pieces() []Piece {
	out := make([]Piece, pieceCount)
	for i := range out {
		out[i] = Piece(i)
	}
	return out
}

func (c Cue) String() string {
	if c >= cueCount {
		return fmt.Sprintf("cue(%d)", uint8(c))
	}
	return names[c]
}

// Cues lists every cue in order.
func Cues() []Cue {
	out := make([]Cue, cueCount)
	for i := range out {
		out[i] = Cue(i)
	}
	return out
}

const recentCap = 32

var (
	mu sync.Mutex
	// Test binaries start muted; tests read Recent.
	enabled  = !testing.Testing()
	volume   = 0.5
	recent   []Cue
	sounds   [cueCount]*audio.Sound
	drums    [pieceCount]*audio.Sound
	watching sync.Once
)

// Play records the cue and plays it when enabled.
func Play(c Cue) {
	if c >= cueCount {
		panic("sound: unknown " + c.String())
	}
	emit(c, &sounds[c], func() []float32 { return render(scores[c]) }, "cue", c.String())
}

// PlayNote records a Note and plays hz when enabled.
func PlayNote(hz float64) {
	if !(hz >= 20 && hz <= 20000) {
		panic(fmt.Sprintf("sound: note %v Hz outside 20–20000", hz))
	}
	emit(Note, nil, func() []float32 { return render(tone(hz)) }, "cue", "note", "hz", hz)
}

// PlayDrum records a Drum and plays the piece when enabled.
func PlayDrum(p Piece) {
	if p >= pieceCount {
		panic("sound: unknown " + p.String())
	}
	emit(Drum, &drums[p], func() []float32 { return render(drumScores[p]) }, "cue", "drum", "piece", p.String())
}

// emit records c and plays its sound when enabled.
func emit(c Cue, cache **audio.Sound, fresh func() []float32, attrs ...any) {
	mu.Lock()
	recent = append(recent, c)
	if len(recent) > recentCap {
		recent = recent[len(recent)-recentCap:]
	}
	on, v := enabled, volume
	var s *audio.Sound
	if on {
		if cache != nil {
			if *cache == nil {
				*cache = audio.NewSound(rate, 1, fresh())
			}
			s = *cache
		} else {
			s = audio.NewSound(rate, 1, fresh())
		}
	}
	mu.Unlock()
	slog.Debug("crychic: sound", append(attrs, "played", on)...)
	if !on {
		return
	}
	s.NewPlayer().SetVolume(v).Play()
	watching.Do(func() { go watch() })
}

// watch logs each change of the output's error.
func watch() {
	var last error
	for range time.Tick(time.Second) {
		last = noteErr(last, audio.Err())
	}
}

func noteErr(last, err error) error {
	switch {
	case err != nil && (last == nil || err.Error() != last.Error()):
		slog.Warn("crychic: sound output unavailable", "err", err)
	case err == nil && last != nil:
		slog.Info("crychic: sound output back")
	}
	return err
}

// SetEnabled turns every cue on or off.
func SetEnabled(on bool) {
	mu.Lock()
	enabled = on
	mu.Unlock()
	slog.Debug("crychic: sound enabled", "on", on)
}

// Enabled reports whether cues play.
func Enabled() bool {
	mu.Lock()
	defer mu.Unlock()
	return enabled
}

// SetVolume sets the gain of cues, from 0 to 1.
func SetVolume(v float64) {
	if !(v >= 0 && v <= 1) {
		panic(fmt.Sprintf("sound: volume %v outside 0–1", v))
	}
	mu.Lock()
	volume = v
	mu.Unlock()
}

// Volume returns the gain of cues.
func Volume() float64 {
	mu.Lock()
	defer mu.Unlock()
	return volume
}

// Recent returns the last cues played or muted, oldest first.
func Recent() []Cue {
	mu.Lock()
	defer mu.Unlock()
	return append([]Cue(nil), recent...)
}
