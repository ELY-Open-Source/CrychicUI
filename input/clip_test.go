package input

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestSignatureStrokeStaysOnPad(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { SignaturePad(c, SignaturePadOptions{}) }, 400, 300)
	pad := ui.LightTheme().Rem(20)
	tt.Press(pad-40, 100)
	for i := 1; i <= 8; i++ {
		tt.Advance(16 * time.Millisecond)
		tt.Move(pad-40+float32(i)*15, 100)
	}
	tt.Release(pad+80, 100)
	tt.Advance(time.Second)
	bg := ui.LightTheme().Background
	for x := int(pad) + 6; x < int(pad)+70; x++ {
		if p := tt.Image().RGBAAt(x, 100); p.R != bg.R || p.G != bg.G || p.B != bg.B {
			t.Fatalf("stroke ink at %d,100 outside the pad: %v", x, p)
		}
	}
}
