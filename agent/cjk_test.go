package agent

import (
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"
)

func TestStreamTextUnspacedScripts(t *testing.T) {
	for _, text := range []string{
		"今天天气很好，我们去公园散步，看到了很多鸫在树上唱歌。",
		"きょうはいい天気ですね。コーヒーを飲みましょう。",
	} {
		tt := ui.NewTester(padded(func(c *ui.Context) { StreamText(c, text, false) }), 600, 200)
		tt.Frame()
		shown, steps, biggest := 0, 0, 0
		for range 400 {
			tt.Advance(16 * time.Millisecond)
			n := 0
			for _, x := range tt.Texts() {
				if strings.HasPrefix(text, x) {
					n = max(n, utf8.RuneCountInString(x))
				}
			}
			if n != shown {
				steps++
				biggest = max(biggest, n-shown)
				shown = n
			}
		}
		if shown != utf8.RuneCountInString(text) || steps < shown/2 || biggest > 2 {
			t.Errorf("%q: %d steps, largest jump %d, shown %d", text, steps, biggest, shown)
		}
	}
}

func TestPrefixTokens(t *testing.T) {
	for _, c := range []struct {
		s    string
		n    int
		want string
	}{
		{"one two three", 2, "one two"},
		{"中文字", 2, "中文"},
		{"Go 语言很快", 2, "Go 语"},
		{"Go语言", 1, "Go"},
		{"好，走", 2, "好，"},
		{"中文", 0, ""},
		{"中文", 9, "中文"},
	} {
		if got := prefixTokens(c.s, c.n); got != c.want {
			t.Errorf("prefixTokens(%q, %d) = %q, want %q", c.s, c.n, got, c.want)
		}
	}
}
