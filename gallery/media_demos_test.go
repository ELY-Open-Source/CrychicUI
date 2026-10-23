package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestMediaDemosLog(t *testing.T) {
	press := func(name string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error { return tt.Click(name) }
	}
	keys := func(n int) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error {
			for range n {
				tt.Key(0, ui.KeyTab)
			}
			tt.Key(0, ui.KeyRight)
			return nil
		}
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"ImageView":      {press("Reload"), "reloading"},
		"Gallery":        {press("Photo 2"), "open: 1"},
		"Carousel":       {press("Next"), "slide 2"},
		"Slideshow":      {func(tt *ui.Tester) error { tt.Advance(4200 * time.Millisecond); return nil }, "slide 2"},
		"ThumbStrip":     {press("Photo 3"), "photo 3"},
		"ZoomImage":      {press("Zoom in"), "zoom 1.25"},
		"Compare":        {keys(1), "split 0.51"},
		"Cropper":        {keys(1), "crop 0.21"},
		"ImageFilters":   {press("Mono"), "filter 1"},
		"AudioPlayer":    {press("Play"), "playing true"},
		"Waveform":       {press("Play 10% more"), "played 40%"},
		"Spectrum":       {press("Play"), "playing true"},
		"PlayerControls": {press("Shuffle"), "playing false, shuffle true"},
		"SeekBar":        {keys(1), "seek 45s"},
		"Volume":         {press("Mute"), "volume 0"},
		"Playlist":       {press("Dune Road"), "playing track 2"},
		"NowPlaying":     {press("Skip 30 s"), "at 1m42s"},
		"MiniBar":        {press("Play"), "playing true"},
		"Lyrics":         {press("Next line"), "at 4s"},
		"PodcastCard":    {press("Play"), "playing true"},
		"Recorder":       {press("Record"), "state 1"},
		"Piano":          {press("C4"), "note 60"},
		"DrumPad":        {press("Snare"), "hit snare"},
		"Metronome":      {press("Play"), "96 BPM, playing true"},
		"Equalizer": {func(tt *ui.Tester) error {
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeyUp)
			return nil
		}, "gains [4 1"},
		"AlbumCard":    {press("Play"), "Blue Hour playing true"},
		"CoverFlow":    {keys(1), "cover 4"},
		"Stories":      {press("Stories"), "story 2"},
		"Filmstrip":    {press("Advance 2 s"), "at 6s"},
		"Scrubber":     {keys(1), "at 1s"},
		"Captions":     {press("Next line"), "at 3s"},
		"Chapters":     {press("Night shift"), "jumped to 3m0s"},
		"SpeedMenu":    {press("Speed: 1×"), "speed 1×, open true"},
		"QualityMenu":  {press("Quality: Auto"), "quality Auto, open true"},
		"Tuner":        {keys(1), "98.1 MHz"},
		"VideoFrame":   {press("Play"), "playing true"},
		"Turntable":    {press("Play"), "playing true"},
		"Cassette":     {press("Play"), "playing true"},
		"Visualizer":   {press("Play"), "playing true"},
		"MediaGrid":    {press("Ferns, Photo"), "picked Ferns"},
		"ImageStack":   {press("Photo stack"), "photo 2 on top"},
		"AvatarUpload": {keys(1), "zoom 1.5, pan 0.10"},
	}
	// Quiet demos only show; they have no events.
	quiet := map[string]bool{"PhotoGrid": true, "Collage": true, "Polaroids": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "media" {
			continue
		}
		seen++
		if quiet[e.Name] {
			continue
		}
		act, ok := acts[e.Name]
		if !ok {
			t.Errorf("%s has no logged action", e.Name)
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
		tt.Advance(time.Second)
		if err := act.do(tt); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Frame()
		found := false
		for _, l := range logs {
			found = found || strings.HasPrefix(l, act.want)
		}
		if !found {
			t.Errorf("%s logged %q, want a line starting %q", e.Name, logs, act.want)
		}
	}
	if seen != len(acts)+len(quiet) {
		t.Errorf("checked %d media demos, want %d", seen, len(acts))
	}
}
