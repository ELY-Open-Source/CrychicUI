package form

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/input"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestContactForm(t *testing.T) {
	var v Contact
	sent := 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		if _, send := ContactForm(c, &v, ContactFormOptions{}); send.Clicked() {
			sent++
		}
	}, 400, 500)
	click(t, tt, "Send")
	typeInto(t, tt, "Name", "Ada")
	typeInto(t, tt, "Email", "ada")
	if tt.HasText("Enter a valid email address.") {
		t.Error("error shown while typing")
	}
	tt.Key(0, ui.KeyTab)
	tt.Frame()
	find(t, tt, "Enter a valid email address.")
	typeInto(t, tt, "Email", "@example.com")
	if tt.HasText("Enter a valid email address.") {
		t.Error("error stayed after the fix")
	}
	typeInto(t, tt, "Message", "Hello")
	click(t, tt, "Send")
	if sent != 1 || v != (Contact{"Ada", "ada@example.com", "Hello"}) {
		t.Fatalf("sent %d, value %+v", sent, v)
	}
	cues(t, before, sound.Send)
	panics(t, "nil value", func() { render(func(c *ui.Context) { ContactForm(c, nil, ContactFormOptions{}) }) })
}

func TestLoginBox(t *testing.T) {
	var v Login
	problem, signed := "", 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		if _, submit := LoginBox(c, &v, LoginBoxOptions{Error: problem}); submit.Clicked() {
			signed++
		}
	}, 400, 500)
	settles(t, tt, "LoginBox")
	click(t, tt, "Sign in")
	typeInto(t, tt, "Email", "ada@example.com")
	typeInto(t, tt, "Password", "secret")
	click(t, tt, "Remember me")
	click(t, tt, "Sign in")
	if signed != 1 || !v.Remember || v.Password != "secret" {
		t.Fatalf("signed %d, value %+v", signed, v)
	}
	if tt.HasText("secret") {
		t.Error("password shown in clear")
	}
	problem = "Wrong password."
	tt.Frame()
	tt.Frame()
	find(t, tt, "Wrong password.")
	cues(t, before, sound.Toggle, sound.Tap, sound.Error)
}

func TestFeedbackForm(t *testing.T) {
	var v Feedback
	sent, changed := 0, 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		root, send := FeedbackForm(c, &v, FeedbackFormOptions{})
		if root.Changed() {
			changed++
		}
		if send.Clicked() {
			sent++
		}
	}, 400, 400)
	click(t, tt, "Send feedback")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	typeInto(t, tt, "Comment", "Nice")
	click(t, tt, "Send feedback")
	if sent != 1 || v.Rating != 2 || v.Comment != "Nice" || changed == 0 {
		t.Fatalf("sent %d, changed %d, value %+v", sent, changed, v)
	}
	cues(t, before, sound.Select, sound.Select, sound.Send)
}

func TestNewsletterForm(t *testing.T) {
	var email string
	subscribed := 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		if _, sub := NewsletterForm(c, &email, NewsletterFormOptions{}); sub.Clicked() {
			subscribed++
		}
	}, 400, 120)
	typeInto(t, tt, "Email", "ada")
	click(t, tt, "Subscribe")
	if subscribed != 0 {
		t.Fatal("bad address subscribed")
	}
	typeInto(t, tt, "Email", "@example.com")
	click(t, tt, "Subscribe")
	settles(t, tt, "NewsletterForm")
	if subscribed != 1 || tt.HasText("Subscribe") {
		t.Errorf("subscribed %d, texts %q", subscribed, tt.Texts())
	}
	find(t, tt, "You're subscribed.")
	cues(t, before, sound.Success)
}

func TestPaymentForm(t *testing.T) {
	var v Payment
	paid := 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		if _, pay := PaymentForm(c, &v, PaymentFormOptions{}); pay.Clicked() {
			paid++
		}
	}, 600, 500)
	for _, s := range []string{"4242424242424242", "1230", "123"} {
		tt.Key(0, ui.KeyTab)
		tt.Type(s)
	}
	tt.Frame()
	click(t, tt, "Pay")
	if paid != 0 {
		t.Fatal("paid without a name")
	}
	typeInto(t, tt, "Name on card", "Ada Lovelace")
	before = baseline()
	click(t, tt, "Pay")
	if paid != 1 {
		t.Fatalf("paid %d, value %+v", paid, v)
	}
	cues(t, before, sound.Success)
}

func TestPayable(t *testing.T) {
	good := Payment{"4242 4242 4242 4242", "12/30", "123", "Ada"}
	if !payable(&good) {
		t.Fatal("good card refused")
	}
	for _, bad := range []Payment{
		{"4242 4242 4242", "12/30", "123", "Ada"},
		{"4242 4242 4242 4242", "13/30", "123", "Ada"},
		{"4242 4242 4242 4242", "12/30", "12a", "Ada"},
		{"4242 4242 4242 4242", "12/30", "12345", "Ada"},
		{"4242 4242 4242 4242", "12/30", "123", " "},
	} {
		if payable(&bad) {
			t.Errorf("accepted %+v", bad)
		}
	}
}

func TestPaymentFormCopy(t *testing.T) {
	var v Payment
	tt := ui.NewTester(func(c *ui.Context) {
		PaymentForm(c, &v, PaymentFormOptions{Card: input.CardNumberFieldOptions{Label: "Kortnummer"},
			Expiry: input.ExpiryFieldOptions{Label: "Utgång"}, CVC: input.CVCFieldOptions{Label: "Kod"}})
	}, 600, 500)
	for _, s := range []string{"Kortnummer", "Utgång", "Kod"} {
		find(t, tt, s)
	}
}
