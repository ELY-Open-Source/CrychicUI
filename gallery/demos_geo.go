package gallery

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/geo"
	"github.com/egoist/mygo/ui"
)

var geoEntries = []Entry{
	{ID: 1486, Name: "MapCanvas", Category: "geo", Summary: "Drawn map that pans and zooms; cue Tick.", Demo: demoGeoMapCanvas},
	{ID: 1487, Name: "Pin", Category: "geo", Summary: "Pin that drops in; cue Drop.", Demo: demoGeoPin},
	{ID: 1488, Name: "Cluster", Category: "geo", Summary: "Pin cluster with count.", Demo: demoGeoCluster},
	{ID: 1489, Name: "LocationSearch", Category: "geo", Summary: "Places search; cue Select.", Demo: demoGeoLocationSearch},
	{ID: 1490, Name: "PlaceCard", Category: "geo", Summary: "Place with rating and hours.", Demo: demoGeoPlaceCard},
	{ID: 1491, Name: "Directions", Category: "geo", Summary: "Route steps.", Demo: demoGeoDirections},
	{ID: 1492, Name: "Distance", Category: "geo", Summary: "Distance badge.", Demo: demoGeoDistance},
	{ID: 1493, Name: "Compass", Category: "geo", Summary: "Needle turning to a heading.", Demo: demoGeoCompass},
	{ID: 1494, Name: "Coordinates", Category: "geo", Summary: "Latitude and longitude fields.", Demo: demoGeoCoordinates},
	{ID: 1495, Name: "RegionSelect", Category: "geo", Summary: "Country and region; cue Select.", Demo: demoGeoRegionSelect},
	{ID: 1496, Name: "StoreLocator", Category: "geo", Summary: "Stores list and map.", Demo: demoGeoStoreLocator},
	{ID: 1497, Name: "DeliveryTracker", Category: "geo", Summary: "Courier moving along a route.", Demo: demoGeoDeliveryTracker},
	{ID: 1498, Name: "Geofence", Category: "geo", Summary: "Radius circle resized; cue Tick.", Demo: demoGeoGeofence},
	{ID: 1499, Name: "Elevation", Category: "geo", Summary: "Elevation profile.", Demo: demoGeoElevation},
	{ID: 1500, Name: "AddressLookup", Category: "geo", Summary: "Address suggestions; cue Select.", Demo: demoGeoAddressLookup},
}

// geoStudio is the Riverbend Clay Co. studio in the Mission.
var geoStudio = geo.Point{Lat: 37.7599, Lng: -122.4148}

func demoGeoMapCanvas(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	view := ui.Local(root, "view", func() geo.View { return geo.View{Center: geoStudio, Zoom: 14} })
	root.Children(func() {
		if geo.MapCanvas(c, view, []geo.MapPin{
			{Name: "Riverbend Clay Co.", At: geoStudio},
			{Name: "Mission Clay Supply", At: geo.Point{Lat: 37.7650, Lng: -122.4190}},
			{Name: "Kiln & Wheel Repair", At: geo.Point{Lat: 37.7545, Lng: -122.4060}},
		}, geo.MapCanvasOptions{}).Changed() {
			log(fmt.Sprintf("center %.4f, %.4f, zoom %g", view.Center.Lat, view.Center.Lng, view.Zoom))
		}
	})
}

func demoGeoPin(c *ui.Context, log func(string)) {
	ui.Row(c).Width(360).Gap(16).Justify(ui.Center).Children(func() {
		for i, name := range []string{"Studio", "Clay supply", "Kiln repair"} {
			if geo.Pin(c, name, geo.PinOptions{Tint: [3]int{0, 1, 3}[i]}).Clicked() {
				log("dropped " + name)
			}
		}
	})
}

func demoGeoCluster(c *ui.Context, log func(string)) {
	ui.Row(c).Width(360).Gap(16).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
		for _, n := range []int{4, 27, 312, 1840} {
			geo.Cluster(c, n, geo.ClusterOptions{})
		}
	})
}

func geoPlaces() []geo.Place {
	return []geo.Place{
		{Name: "Riverbend Clay Co.", Address: "2125 Bryant St", Icon: "flame", At: geoStudio},
		{Name: "Mission Clay Supply", Address: "2414 Mission St", Icon: "package", At: geo.Point{Lat: 37.7650, Lng: -122.4190}},
		{Name: "Fogline Ceramics", Address: "2900 18th St", Icon: "store", At: geo.Point{Lat: 37.7620, Lng: -122.4110}},
		{Name: "Kiln & Wheel Repair", Address: "1450 Folsom St", Icon: "wrench", At: geo.Point{Lat: 37.7545, Lng: -122.4060}},
		{Name: "Dolores Park Café", Address: "501 Dolores St", Icon: "coffee", At: geo.Point{Lat: 37.7610, Lng: -122.4260}},
	}
}

func demoGeoLocationSearch(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	query, picked := textState(root, "query", ""), textState(root, "picked", "")
	root.Children(func() {
		before := *picked
		from := geoStudio
		if geo.LocationSearch(c, query, geoPlaces(), picked, geo.LocationSearchOptions{From: &from}).Changed() {
			if *picked != before {
				log("picked " + *picked)
			} else {
				log(fmt.Sprintf("query %d runes", len([]rune(*query))))
			}
		}
	})
}

func demoGeoPlaceCard(c *ui.Context, log func(string)) {
	ui.Column(c).Width(380).Children(func() {
		geo.PlaceCard(c, geo.Place{
			Name: "Dolores Park Café", Address: "501 Dolores St, San Francisco", Icon: "coffee", Hours: "7:00–19:00",
			At: geo.Point{Lat: 37.7610, Lng: -122.4260}, Rating: 4.6, Reviews: 1284, Open: true, Photo: productShots()[1],
		}, geo.PlaceCardOptions{})
	})
}

func demoGeoDirections(c *ui.Context, log func(string)) {
	ui.Column(c).Width(400).Children(func() {
		geo.Directions(c, []geo.Step{
			{Text: "Head north on Bryant St", Turn: "arrow-up", Meters: 350, Minutes: 1},
			{Text: "Turn left onto 20th St", Turn: "corner-up-left", Meters: 600, Minutes: 2},
			{Text: "Turn right onto Mission St", Turn: "corner-up-right", Meters: 900, Minutes: 3},
			{Text: "Mission Clay Supply is on the left", Turn: "arrow-up-left", Meters: 150, Minutes: 1},
		}, geo.DirectionsOptions{Arrive: "Arrive at Mission Clay Supply"})
	})
}

func demoGeoDistance(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Gap(10).AlignItems(ui.Start).Children(func() {
		geo.Distance(c, 850, geo.DistanceOptions{Icon: "footprints", Minutes: 11})
		geo.Distance(c, 4200, geo.DistanceOptions{Icon: "bike", Minutes: 16})
		geo.Distance(c, 12400, geo.DistanceOptions{Icon: "car", Minutes: 68})
		geo.Distance(c, 4200, geo.DistanceOptions{Icon: "car", Imperial: true})
	})
}

func demoGeoCompass(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320).Gap(12).AlignItems(ui.Center)
	heading := intState(root, "heading", 292)
	root.Children(func() {
		geo.Compass(c, float64(*heading), geo.CompassOptions{})
		if button.Outline(c, "Turn 90°").Clicked() {
			*heading = (*heading + 90) % 360
			log(fmt.Sprintf("heading %d", *heading))
		}
	})
}

func demoGeoCoordinates(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	at := ui.Local(root, "at", func() geo.Point { return geoStudio })
	root.Children(func() {
		r, copied := geo.Coordinates(c, at, geo.CoordinatesOptions{})
		if copied.Clicked() {
			log("copied")
		}
		if r.Changed() {
			log("coordinates edited")
		}
	})
}

func demoGeoRegionSelect(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400).Height(300)
	country, region := textState(root, "country", "United States"), textState(root, "region", "California")
	root.Children(func() {
		was := *country + "/" + *region
		if geo.RegionSelect(c, []geo.Country{
			{Name: "United States", Regions: []string{"California", "Oregon", "Washington", "New York"}},
			{Name: "Canada", Regions: []string{"British Columbia", "Ontario", "Quebec"}},
			{Name: "Japan", Regions: []string{"Tokyo", "Kyoto", "Osaka"}},
			{Name: "Mexico", Regions: []string{"Jalisco", "Oaxaca", "Yucatán"}},
		}, country, region, geo.RegionSelectOptions{}).Changed() {
			if *country+"/"+*region == was {
				log("menu open")
			} else {
				log(fmt.Sprintf("country %s, region %s", *country, *region))
			}
		}
	})
}

func demoGeoStoreLocator(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	picked := textState(root, "picked", "Riverbend Mission")
	root.Children(func() {
		from := geo.Point{Lat: 37.7700, Lng: -122.4300}
		if geo.StoreLocator(c, []geo.Store{
			{Name: "Riverbend Mission", Address: "2125 Bryant St", At: geoStudio, Open: true},
			{Name: "Riverbend Dogpatch", Address: "901 Minnesota St", At: geo.Point{Lat: 37.7576, Lng: -122.3880}, Open: true},
			{Name: "Riverbend North Beach", Address: "1600 Stockton St", At: geo.Point{Lat: 37.8000, Lng: -122.4100}},
			{Name: "Riverbend Sunset", Address: "1220 Irving St", At: geo.Point{Lat: 37.7640, Lng: -122.4700}, Open: true},
		}, picked, geo.StoreLocatorOptions{From: &from}).Changed() {
			log("picked " + *picked)
		}
	})
}

func demoGeoDeliveryTracker(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440).Gap(10).AlignItems(ui.Start)
	step := intState(root, "step", 3)
	root.Children(func() {
		progress := float64(*step) / 10
		geo.DeliveryTracker(c, []geo.Point{
			{Lat: 37.7650, Lng: -122.4190}, {Lat: 37.7650, Lng: -122.4150}, {Lat: 37.7625, Lng: -122.4150},
			{Lat: 37.7625, Lng: -122.4120}, {Lat: 37.7599, Lng: -122.4120}, geoStudio,
		}, progress, geo.DeliveryTrackerOptions{Courier: "Lucas Martin", Minutes: int(math.Round((1 - progress) * 14))})
		if button.Outline(c, "Advance").Disabled(*step >= 10).Clicked() {
			*step = min(10, *step+3)
			log(fmt.Sprintf("progress %d%%", *step*10))
		}
	})
}

func demoGeoGeofence(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	radius := ui.Local(root, "radius", func() float64 { return 500 })
	root.Children(func() {
		if geo.Geofence(c, geoStudio, radius, geo.GeofenceOptions{Label: "Studio alert radius"}).Changed() {
			log(fmt.Sprintf("radius %.0f m", *radius))
		}
	})
}

func demoGeoElevation(c *ui.Context, log func(string)) {
	samples := make([]float64, 41)
	for i := range samples {
		x := float64(i) / 40
		samples[i] = math.Round(60 + 210*math.Sin(x*math.Pi) + 30*math.Sin(x*9))
	}
	ui.Column(c).Width(440).Children(func() {
		geo.Elevation(c, samples, 5400, geo.ElevationOptions{Title: "Twin Peaks loop"})
	})
}

func demoGeoAddressLookup(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	query, chosen := textState(root, "query", "21"), textState(root, "chosen", "")
	root.Children(func() {
		before := *chosen
		if geo.AddressLookup(c, query, []geo.Address{
			{Line: "2125 Bryant St", City: "San Francisco, CA 94110"},
			{Line: "2140 Folsom St", City: "San Francisco, CA 94110"},
			{Line: "501 Dolores St", City: "San Francisco, CA 94114"},
			{Line: "2100 Mission St", City: "San Francisco, CA 94110"},
			{Line: "1450 Folsom St", City: "San Francisco, CA 94103"},
		}, chosen, geo.AddressLookupOptions{}).Changed() {
			if *chosen != before && *chosen != "" {
				log("chose " + *chosen)
			} else {
				log(fmt.Sprintf("query %d runes", len([]rune(*query))))
			}
		}
	})
}
