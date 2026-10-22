package form

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestAdvancedSearch(t *testing.T) {
	var rules []Rule
	searched, changed := 0, 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		root, search := AdvancedSearch(c, []string{"Title", "Author"}, &rules, AdvancedSearchOptions{})
		if root.Changed() {
			changed++
		}
		if search.Clicked() {
			searched++
		}
	}, 700, 400)
	click(t, tt, "Search")
	click(t, tt, "Add rule")
	typeInto(t, tt, "Value", "go")
	click(t, tt, "Add rule")
	settles(t, tt, "AdvancedSearch add")
	if len(rules) != 2 || rules[0] != (Rule{"Title", "contains", "go"}) {
		t.Fatalf("rules %+v", rules)
	}
	click(t, tt, "Remove rule")
	settles(t, tt, "AdvancedSearch remove")
	if len(rules) != 1 || rules[0].Value != "" {
		t.Fatalf("after removing the first: %+v", rules)
	}
	click(t, tt, "Search")
	if searched != 1 || changed < 4 {
		t.Errorf("searched %d, changed %d", searched, changed)
	}
	cues(t, before, sound.Tap, sound.Tap, sound.Tap, sound.Tap)
	panics(t, "no fields", func() { render(func(c *ui.Context) { AdvancedSearch(c, nil, &rules, AdvancedSearchOptions{}) }) })
}

func TestFilterPanel(t *testing.T) {
	var picked []Filter
	facets := []Facet{{"Color", []string{"Red", "Blue"}, []int{3, 5}}, {"Size", []string{"S"}, []int{2}}}
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) { FilterPanel(c, facets, &picked, FilterPanelOptions{}) }, 400, 400)
	find(t, tt, "5")
	if tt.HasText("Clear filters") {
		t.Error("clear shown with nothing picked")
	}
	click(t, tt, "Red")
	click(t, tt, "S")
	if len(picked) != 2 || picked[0] != (Filter{"Color", "Red"}) {
		t.Fatalf("picked %+v", picked)
	}
	click(t, tt, "Red")
	if len(picked) != 1 || picked[0] != (Filter{"Size", "S"}) {
		t.Fatalf("after unpicking: %+v", picked)
	}
	click(t, tt, "Clear filters")
	settles(t, tt, "FilterPanel")
	if len(picked) != 0 || tt.HasText("Clear filters") {
		t.Errorf("picked %+v after clear", picked)
	}
	cues(t, before, sound.Toggle, sound.Toggle, sound.Toggle, sound.Tap)
	panics(t, "counts", func() {
		render(func(c *ui.Context) { FilterPanel(c, []Facet{{"A", []string{"x"}, nil}}, &picked, FilterPanelOptions{}) })
	})
}

func TestSettingsForm(t *testing.T) {
	sound, theme := false, "Light"
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if SettingsForm(c, []SettingGroup{{"General", []Setting{
			{Title: "Sounds", Description: "Play cues", On: &sound},
			{Title: "Theme", Choice: &theme, Choices: []string{"Light", "Dark"}},
		}}}).Changed() {
			changed++
		}
	}, 500, 300)
	find(t, tt, "General")
	click(t, tt, "Sounds")
	if !sound || changed != 1 {
		t.Errorf("sound %v, changed %d", sound, changed)
	}
	panics(t, "both kinds", func() {
		render(func(c *ui.Context) {
			SettingsForm(c, []SettingGroup{{"G", []Setting{{Title: "X", On: &sound, Choice: &theme}}}})
		})
	})
}

func TestProfileForm(t *testing.T) {
	var v Profile
	saved := 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		if _, save := ProfileForm(c, &v, ProfileFormOptions{}); save.Clicked() {
			saved++
		}
	}, 500, 600)
	click(t, tt, "Save profile")
	click(t, tt, "Add link")
	typeInto(t, tt, "Link", "https://example.com")
	typeInto(t, tt, "Name", "Ada")
	click(t, tt, "Save profile")
	if saved != 1 || v.Name != "Ada" || len(v.Links) != 1 || v.Links[0] != "https://example.com" {
		t.Fatalf("saved %d, value %+v", saved, v)
	}
	click(t, tt, "Remove link")
	settles(t, tt, "ProfileForm")
	if len(v.Links) != 0 {
		t.Errorf("links %q", v.Links)
	}
	cues(t, before, sound.Tap, sound.Success, sound.Tap)
}

func TestAddressForm(t *testing.T) {
	var addr input.Address
	country, changed := "Canada", 0
	tt := ui.NewTester(func(c *ui.Context) {
		if AddressForm(c, &addr, &country, []string{"Canada", "Japan"}, AddressFormOptions{}).Changed() {
			changed++
		}
	}, 500, 400)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Type("1 Main St")
	tt.Frame()
	if addr.Street != "1 Main St" || changed == 0 {
		t.Errorf("street %q, changed %d", addr.Street, changed)
	}
	bad := "Mars"
	panics(t, "unknown country", func() {
		render(func(c *ui.Context) { AddressForm(c, &addr, &bad, []string{"Canada"}, AddressFormOptions{}) })
	})
}

func TestShippingForm(t *testing.T) {
	var v Shipping
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		ShippingForm(c, &v, []Speed{{"Standard", "Free, 5 days"}, {"Express", "$9, 1 day"}}, ShippingFormOptions{})
	}, 500, 600)
	click(t, tt, "Express")
	if v.Speed != "Express" {
		t.Errorf("speed %q", v.Speed)
	}
	cues(t, before, sound.Select)
	panics(t, "no speeds", func() { render(func(c *ui.Context) { ShippingForm(c, &v, nil, ShippingFormOptions{}) }) })
}

func TestAdvancedSearchEmptyOperators(t *testing.T) {
	var rules []Rule
	tt := ui.NewTester(func(c *ui.Context) {
		AdvancedSearch(c, []string{"Title"}, &rules, AdvancedSearchOptions{Operators: []string{}})
	}, 700, 300)
	click(t, tt, "Add rule")
	if len(rules) != 1 || rules[0].Op != "contains" {
		t.Errorf("rules %+v", rules)
	}
}
