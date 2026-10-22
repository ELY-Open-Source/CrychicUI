package media

import (
	"bytes"
	"image/color"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func rgba(c ui.Color) color.RGBA { return color.RGBA{c.R, c.G, c.B, 255} }

func padded(view func(c *ui.Context)) func(c *ui.Context) {
	return func(c *ui.Context) { ui.Column(c).Padding(10).Children(func() { view(c) }) }
}

func TestSeekBar(t *testing.T) {
	pos := 30 * time.Second
	var acc ui.Color
	tt := ui.NewTester(padded(func(c *ui.Context) {
		acc = c.Theme().Accent
		SeekBar(c, "Seek", &pos, time.Minute, 2*time.Minute, SeekBarOptions{})
	}), 420, 200)
	settles(t, tt, "seek bar")
	r := find(t, tt, "Seek")
	if !slices.Contains(tt.Texts(), "0:30") || !slices.Contains(tt.Texts(), "2:00") {
		t.Errorf("times %q", tt.Texts())
	}
	if !at(tt, r.X+r.W*0.12, r.Y+r.H/2, rgba(acc)) || at(tt, r.X+r.W*0.6, r.Y+r.H/2, rgba(acc)) {
		t.Error("played part should end near a quarter")
	}
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	if pos != 35*time.Second || !slices.Contains(tt.Texts(), "0:35") {
		t.Errorf("pos %v after Right, want 35s", pos)
	}
	tt.Press(r.X+r.W*0.75, r.Y+r.H/2)
	tt.Frame()
	tt.Release(r.X+r.W*0.75, r.Y+r.H/2)
	tt.Frame()
	if d := pos - 90*time.Second; d < -6*time.Second || d > 6*time.Second {
		t.Errorf("pos %v after a press at three quarters", pos)
	}
	cues(t, before, sound.Tick, sound.Tick)
	panics(t, "pos past total", func() {
		p := 3 * time.Minute
		render(func(c *ui.Context) { SeekBar(c, "x", &p, 0, time.Minute, SeekBarOptions{}) })
	})
	panics(t, "no total", func() {
		p := time.Duration(0)
		render(func(c *ui.Context) { SeekBar(c, "x", &p, 0, 0, SeekBarOptions{}) })
	})
}

func TestAudioPlayer(t *testing.T) {
	clip := Clip{Notes: []float64{440, 0, 523.25}, Beat: 500 * time.Millisecond}
	pos, playing, changed := time.Duration(0), false, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if AudioPlayer(c, "Tune", clip, &pos, &playing, AudioPlayerOptions{}).Changed() {
			changed++
		}
	}), 500, 200)
	tt.Frame()
	before := baseline()
	click(t, tt, "Play")
	if !playing {
		t.Fatal("Play did not start")
	}
	tt.Advance(600 * time.Millisecond)
	if d := pos - 600*time.Millisecond; d < -40*time.Millisecond || d > 40*time.Millisecond {
		t.Errorf("pos %v after 600 ms", pos)
	}
	cues(t, before, sound.Toggle, sound.Note)
	tt.Advance(500 * time.Millisecond)
	cues(t, before, sound.Toggle, sound.Note, sound.Note)
	tt.Advance(time.Second)
	if playing || pos != 1500*time.Millisecond {
		t.Errorf("playing %v at %v, want stopped at the end", playing, pos)
	}
	settles(t, tt, "player")
	if changed < 2 {
		t.Errorf("changed %d times, want play and end", changed)
	}
	cues(t, before, sound.Toggle, sound.Note, sound.Note)
	panics(t, "bad note", func() {
		p, on := time.Duration(0), false
		render(func(c *ui.Context) {
			AudioPlayer(c, "x", Clip{Notes: []float64{5}, Beat: time.Second}, &p, &on, AudioPlayerOptions{})
		})
	})
	panics(t, "no beat", func() {
		p, on := time.Duration(0), false
		render(func(c *ui.Context) { AudioPlayer(c, "x", Clip{Notes: []float64{440}}, &p, &on, AudioPlayerOptions{}) })
	})
}

func TestAudioPlayerOutsidePause(t *testing.T) {
	clip := Clip{Notes: []float64{440, 0, 523.25, 0, 440, 0}, Beat: 500 * time.Millisecond}
	pos, playing := time.Duration(0), false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		AudioPlayer(c, "Tune", clip, &pos, &playing, AudioPlayerOptions{})
	}), 500, 200)
	tt.Frame()
	click(t, tt, "Play")
	tt.Advance(400 * time.Millisecond)
	playing = false
	tt.Frame()
	paused := pos
	tt.Advance(2 * time.Second)
	if pos != paused {
		t.Fatalf("pos moved from %v to %v while paused", paused, pos)
	}
	playing = true
	tt.Frame()
	tt.Advance(100 * time.Millisecond)
	if d := pos - paused - 100*time.Millisecond; d < -40*time.Millisecond || d > 40*time.Millisecond {
		t.Errorf("pos %v after resuming at %v for 100 ms", pos, paused)
	}
}

func TestWaveform(t *testing.T) {
	var acc ui.Color
	levels := []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		acc = c.Theme().Accent
		Waveform(c, "Clip", levels, 0.5, WaveformOptions{})
	}), 420, 200)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("bars do not grow in")
	}
	settles(t, tt, "waveform")
	r := find(t, tt, "Clip")
	if !at(tt, r.X+r.W*0.25, r.Y+r.H/2, rgba(acc)) || at(tt, r.X+r.W*0.85, r.Y+r.H/2, rgba(acc)) {
		t.Error("played half should be accent, the rest not")
	}
	panics(t, "level above 1", func() { render(func(c *ui.Context) { Waveform(c, "x", []float64{2}, 0, WaveformOptions{}) }) })
	panics(t, "played below 0", func() { render(func(c *ui.Context) { Waveform(c, "x", levels, -0.1, WaveformOptions{}) }) })
}

func TestSpectrum(t *testing.T) {
	playing := true
	tt := ui.NewTester(padded(func(c *ui.Context) { Spectrum(c, "Spec", playing, SpectrumOptions{}) }), 420, 200)
	tt.Advance(200 * time.Millisecond)
	a := bytes.Clone(tt.Image().Pix)
	tt.Advance(300 * time.Millisecond)
	if bytes.Equal(a, tt.Image().Pix) || !tt.FramePending() {
		t.Error("bars do not move while playing")
	}
	playing = false
	tt.Frame()
	settles(t, tt, "stopped spectrum")
	playing = true
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Frame()
	settles(t, tt, "spectrum under reduced motion")
}

func TestPlayerControls(t *testing.T) {
	playing, shuffle, rep := false, false, RepeatOff
	changed, prevs, nexts := 0, 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		root, prev, next := PlayerControls(c, "Player", &playing, &shuffle, &rep, PlayerControlsOptions{})
		if root.Changed() {
			changed++
		}
		if prev.Clicked() {
			prevs++
		}
		if next.Clicked() {
			nexts++
		}
	}, 500, 200)
	tt.Frame()
	before := baseline()
	click(t, tt, "Shuffle")
	if !shuffle {
		t.Error("Shuffle did not turn on")
	}
	for _, want := range []Repeat{RepeatAll, RepeatOne, RepeatOff} {
		click(t, tt, "Repeat")
		if rep != want {
			t.Errorf("repeat %d, want %d", rep, want)
		}
	}
	click(t, tt, "Previous")
	click(t, tt, "Next")
	click(t, tt, "Play")
	if !playing || prevs != 1 || nexts != 1 || changed != 5 {
		t.Errorf("playing %v, prev %d, next %d, changed %d", playing, prevs, nexts, changed)
	}
	cues(t, before, sound.Toggle, sound.Toggle, sound.Toggle, sound.Toggle, sound.Toggle, sound.Toggle, sound.Toggle)
	panics(t, "unknown repeat", func() {
		r := Repeat(7)
		render(func(c *ui.Context) { PlayerControls(c, "x", &playing, &shuffle, &r, PlayerControlsOptions{}) })
	})
}

func TestVolume(t *testing.T) {
	level := 60.0
	tt := ui.NewTester(padded(func(c *ui.Context) { Volume(c, "Volume", &level, VolumeOptions{}) }), 420, 200)
	tt.Frame()
	before := baseline()
	click(t, tt, "Mute")
	if level != 0 {
		t.Fatalf("level %v after Mute", level)
	}
	click(t, tt, "Unmute")
	if level != 60 {
		t.Fatalf("level %v after Unmute, want 60", level)
	}
	// The click left focus on the speaker button.
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	if level != 61 {
		t.Errorf("level %v after Right, want 61", level)
	}
	cues(t, before, sound.Toggle, sound.Toggle, sound.Tick)
}

func TestPlaylist(t *testing.T) {
	tracks := []Track{{Title: "First", Artist: "A", Length: 200 * time.Second}, {Title: "Second", Artist: "B", Length: 3 * time.Minute}}
	playing := -1
	tt := ui.NewTester(func(c *ui.Context) { Playlist(c, "Queue", tracks, &playing, PlaylistOptions{}) }, 500, 300)
	settles(t, tt, "playlist")
	before := baseline()
	click(t, tt, "Second")
	click(t, tt, "Second")
	if playing != 1 {
		t.Fatalf("playing %d, want 1", playing)
	}
	cues(t, before, sound.Select)
	tt.Advance(time.Second)
	if !tt.FramePending() || slices.Contains(tt.Texts(), "2") || !slices.Contains(tt.Texts(), "1") {
		t.Errorf("the playing row should show moving bars for its number: %q", tt.Texts())
	}
	panics(t, "no length", func() {
		p := -1
		render(func(c *ui.Context) { Playlist(c, "x", []Track{{Title: "a"}}, &p, PlaylistOptions{}) })
	})
}

func TestNowPlaying(t *testing.T) {
	track := Track{Title: "Blue Hour", Artist: "Kai", Length: 3 * time.Minute}
	tt := ui.NewTester(padded(func(c *ui.Context) { NowPlaying(c, track, time.Minute, NowPlayingOptions{}) }), 500, 300)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("NowPlaying does not arrive")
	}
	settles(t, tt, "now playing")
	for _, s := range []string{"Blue Hour", "Kai", "1:00", "3:00"} {
		if !slices.Contains(tt.Texts(), s) {
			t.Errorf("%q missing from %q", s, tt.Texts())
		}
	}
	panics(t, "pos past length", func() {
		render(func(c *ui.Context) { NowPlaying(c, track, 4*time.Minute, NowPlayingOptions{}) })
	})
}

func TestMiniBar(t *testing.T) {
	track := Track{Title: "Blue Hour", Artist: "Kai", Length: 4 * time.Minute}
	playing, changed := false, 0
	var acc ui.Color
	tt := ui.NewTester(padded(func(c *ui.Context) {
		acc = c.Theme().Accent
		if MiniBar(c, track, &playing, 2*time.Minute, MiniBarOptions{}).Changed() {
			changed++
		}
	}), 500, 200)
	settles(t, tt, "mini bar")
	r := find(t, tt, "Blue Hour")
	y := r.Y + r.H - 2
	if !at(tt, r.X+r.W*0.25, y, rgba(acc)) || at(tt, r.X+r.W*0.75, y, rgba(acc)) {
		t.Error("progress line should end at half")
	}
	before := baseline()
	click(t, tt, "Play")
	if !playing || changed != 1 {
		t.Errorf("playing %v, changed %d", playing, changed)
	}
	cues(t, before, sound.Toggle)
}

func TestLyrics(t *testing.T) {
	var lines []Line
	for i := range 12 {
		lines = append(lines, Line{At: time.Duration(i) * 2 * time.Second, Text: "Line " + string(rune('A'+i))})
	}
	pos := time.Duration(0)
	tt := ui.NewTester(padded(func(c *ui.Context) { Lyrics(c, "Song", lines, pos, LyricsOptions{}) }), 420, 300)
	settles(t, tt, "lyrics")
	view := find(t, tt, "Song")
	if r, ok := tt.Find("Line L"); ok && r.Y+r.H > view.Y && r.Y < view.Y+view.H && r.H > 1 {
		t.Fatalf("last line visible at the start: %v in %v", r, view)
	}
	pos = 23 * time.Second
	tt.Frame()
	settles(t, tt, "lyrics")
	if r := find(t, tt, "Line L"); r.Y < view.Y || r.Y+r.H > view.Y+view.H+1 {
		t.Errorf("lit line %v not scrolled into %v", r, view)
	}
	panics(t, "falling times", func() {
		render(func(c *ui.Context) {
			Lyrics(c, "x", []Line{{At: time.Second, Text: "a"}, {At: 0, Text: "b"}}, 0, LyricsOptions{})
		})
	})
}

func TestPodcastCard(t *testing.T) {
	ep := Track{Title: "Episode 9", Artist: "Night Talk", Length: 41*time.Minute + 40*time.Second}
	playing := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		PodcastCard(c, ep, &playing, PodcastCardOptions{Minutes: "%d Min."})
	}), 500, 300)
	settles(t, tt, "podcast card")
	if !slices.Contains(tt.Texts(), "42 Min.") {
		t.Errorf("length missing from %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Play")
	if !playing {
		t.Error("Play did not start")
	}
	cues(t, before, sound.Toggle)
}

func TestRecorder(t *testing.T) {
	state, changed := Idle, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Recorder(c, "Memo", &state, 0.6, RecorderOptions{}).Changed() {
			changed++
		}
	}), 500, 300)
	tt.Frame()
	before := baseline()
	shows := func(s string) {
		t.Helper()
		if !slices.Contains(tt.Texts(), s) {
			t.Errorf("want %q in %q", s, tt.Texts())
		}
	}
	click(t, tt, "Record")
	tt.Advance(2100 * time.Millisecond)
	shows("0:02")
	click(t, tt, "Pause")
	tt.Advance(3 * time.Second)
	shows("0:02")
	click(t, tt, "Resume")
	tt.Advance(1000 * time.Millisecond)
	shows("0:03")
	click(t, tt, "Stop")
	shows("0:00")
	_ = tt.Click("Stop")
	tt.Frame()
	if state != Idle || changed != 4 {
		t.Errorf("state %d, changed %d", state, changed)
	}
	cues(t, before, sound.Toggle, sound.Toggle, sound.Toggle, sound.Toggle)
	panics(t, "level above 1", func() {
		s := Idle
		render(func(c *ui.Context) { Recorder(c, "x", &s, 2, RecorderOptions{}) })
	})
}
