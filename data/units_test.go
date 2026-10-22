package data

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestMove(t *testing.T) {
	items := []string{"a", "b", "c", "d"}
	for _, tc := range []struct {
		rows []int
		to   int
		want []string
	}{
		{[]int{1}, 0, []string{"b", "a", "c", "d"}},
		{[]int{0}, 4, []string{"b", "c", "d", "a"}},
		{[]int{0, 2}, 4, []string{"b", "d", "a", "c"}},
		{[]int{3}, 1, []string{"a", "d", "b", "c"}},
	} {
		if got := move(items, tc.rows, tc.to); !slices.Equal(got, tc.want) {
			t.Errorf("move %v to %d: %v, want %v", tc.rows, tc.to, got, tc.want)
		}
	}
}

func TestLetters(t *testing.T) {
	for i, want := range map[int]string{0: "A", 25: "Z", 26: "AA", 27: "AB", 701: "ZZ", 702: "AAA"} {
		if got := letters(i); got != want {
			t.Errorf("letters(%d) = %q, want %q", i, got, want)
		}
	}
}

func TestMarked(t *testing.T) {
	th := ui.LightTheme()
	spans := marked(th, "Go is not GOLANG", "go")
	var hits []string
	for _, s := range spans {
		if s.Background.A > 0 {
			hits = append(hits, s.Text)
		}
	}
	if !slices.Equal(hits, []string{"Go", "GO"}) {
		t.Errorf("hits %q", hits)
	}
	if len(marked(th, "İstanbul straße", "STRASSE")) == 0 {
		t.Error("no spans for unicode text")
	}
}
