package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/egoist/mygo/ui"
)

var feedbackEntries2 = []Entry{
	{324, "OfflineState", "feedback", "Pulsing offline icon and Retry.", demoOfflineState},
	{325, "Maintenance", "feedback", "Rocking tools icon and a message.", demoMaintenance},
	{326, "LoadingVeil", "feedback", "Dims content under a spinner; fades in and out.", demoLoadingVeil},
	{327, "Busy", "feedback", "Inline spinner and text cross-fading to done.", demoBusy},
	{328, "UploadProgress", "feedback", "File rows with bars; rows come and go.", demoUploadProgress},
	{329, "StepCount", "feedback", "\"Step 2 of 5\" over a filling bar.", demoStepCount},
	{330, "TimeLeftBar", "feedback", "Bar draining over a duration.", demoTimeLeftBar},
	{331, "RetryIn", "feedback", "Counts down, then enables Retry.", demoRetryIn},
	{332, "ErrorList", "feedback", "Field errors; a click focuses the field.", demoErrorList},
	{333, "FieldError", "feedback", "Error text that shakes in.", demoFieldError},
	{334, "Hint", "feedback", "Muted help text with an info icon.", demoHint},
	{335, "Attention", "feedback", "Ring pulsing around a child.", demoAttention},
	{336, "Beacon", "feedback", "Pulsing hotspot; a click opens a tip.", demoBeacon},
	{337, "SaveStatus", "feedback", "Saving spinner cross-fading to a saved check.", demoSaveStatus},
	{338, "SyncStatus", "feedback", "Synced, syncing and failed; syncing spins.", demoSyncStatus},
	{339, "ConnectionStatus", "feedback", "Online, reconnecting and offline dot.", demoConnectionStatus},
	{340, "LiveDot", "feedback", "Red dot pulsing beside LIVE.", demoLiveDot},
	{341, "QuotaBar", "feedback", "Used of a limit; warning past 80%.", demoQuotaBar},
	{342, "StorageBar", "feedback", "Usage bar of growing parts with a legend.", demoStorageBar},
	{343, "Strength", "feedback", "Password strength in four bars and a word.", demoStrength},
	{344, "LevelMeter", "feedback", "Level bars whose peak falls back.", demoLevelMeter},
	{345, "Battery", "feedback", "Battery filling to a level; charging pulses.", demoBattery},
	{346, "Signal", "feedback", "Four bars by strength; bars rise.", demoSignal},
	{347, "ServiceStatus", "feedback", "Service rows with status dots and uptime.", demoServiceStatus},
	{348, "AnnouncementBar", "feedback", "Messages sliding up in turn.", demoAnnouncementBar},
	{349, "UndoToast", "feedback", "Toast with Undo and a countdown ring.", demoUndoToast},
	{350, "UnlockToast", "feedback", "Achievement toast with a passing shine.", demoUnlockToast},
}

func demoOfflineState(c *ui.Context, log func(string)) {
	if _, retry := feedback.OfflineState(c, feedback.OfflineStateOptions{}); retry.Clicked() {
		log("retry")
	}
}

func demoMaintenance(c *ui.Context, log func(string)) {
	feedback.Maintenance(c, feedback.MaintenanceOptions{})
}

func demoLoadingVeil(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).Width(320)
	until := ui.Local(root, "until", func() time.Time { return time.Time{} })
	loading := c.Now().Before(*until)
	if loading {
		c.After(until.Sub(c.Now()))
	}
	root.Children(func() {
		if ui.Button(c, "Reload").Disabled(loading).Clicked() {
			*until = c.Now().Add(2 * time.Second)
			log("reloading")
		}
		feedback.LoadingVeil(c, loading, feedback.LoadingVeilOptions{}, func() {
			ui.Column(c).Gap(6).Padding(16).Background(c.Theme().Surface).Children(func() {
				ui.Text(c, "Quarterly report").Bold()
				ui.Text(c, "Revenue rose eight percent.")
			})
		})
	})
}

func demoBusy(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	busy := boolState(root, "busy", false)
	root.Children(func() {
		if ui.Button(c, "Toggle").Clicked() {
			*busy = !*busy
			log(fmt.Sprintf("busy: %v", *busy))
		}
		feedback.Busy(c, *busy, "Saving draft", "Draft saved")
	})
}

func demoUploadProgress(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).Width(360)
	ups := ui.Local(root, "ups", func() []feedback.Upload {
		return []feedback.Upload{{Name: "cover.png", Progress: 0.4}, {Name: "notes.pdf", Progress: 1}}
	})
	next := intState(root, "next", 3)
	root.Children(func() {
		ui.Row(c).Gap(8).Children(func() {
			if ui.Button(c, "Advance").Clicked() {
				for i := range *ups {
					(*ups)[i].Progress = min(1, (*ups)[i].Progress+0.2)
				}
				log("advanced")
			}
			if ui.Button(c, "Add").Clicked() {
				*ups = append(*ups, feedback.Upload{Name: fmt.Sprintf("file-%d.txt", *next)})
				*next++
				log("added")
			}
			if ui.Button(c, "Clear done").Clicked() {
				var kept []feedback.Upload
				for _, u := range *ups {
					if u.Progress < 1 {
						kept = append(kept, u)
					}
				}
				*ups = kept
				log("cleared")
			}
		})
		feedback.UploadProgress(c, *ups)
	})
}

func demoStepCount(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).Width(320)
	step := intState(root, "step", 2)
	root.Children(func() {
		feedback.StepCount(c, *step, 5, feedback.StepCountOptions{})
		if ui.Button(c, "Next").Disabled(*step >= 5).Clicked() {
			*step++
			log(fmt.Sprintf("step: %d", *step))
		}
	})
}

func demoTimeLeftBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).Width(320)
	ends := ui.Local(root, "ends", func() time.Time { return time.Time{} })
	root.Children(func() {
		if ui.Button(c, "Start 5 s").Clicked() {
			*ends = c.Now().Add(5 * time.Second)
			log("started")
		}
		feedback.TimeLeftBar(c, "Offer", *ends, 5*time.Second)
	})
}

func demoRetryIn(c *ui.Context, log func(string)) {
	root := ui.Row(c).Gap(8)
	at := ui.Local(root, "at", func() time.Time { return time.Time{} })
	root.Children(func() {
		if ui.Button(c, "Fail").Clicked() {
			*at = c.Now().Add(5 * time.Second)
			log("failed")
		}
		if feedback.RetryIn(c, *at, feedback.RetryInOptions{}).Clicked() {
			log("retry")
		}
	})
}

func demoErrorList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(8).Width(320)
	name, email := ui.Local(root, "name", func() *ui.Handle { return &ui.Handle{} }), ui.Local(root, "email", func() *ui.Handle { return &ui.Handle{} })
	a, b := textState(root, "a", ""), textState(root, "b", "ada@")
	root.Children(func() {
		na := ui.TextInput(c.Key("name"), a).Placeholder("Name").Bind(*name)
		if na.Changed() {
			log(fmt.Sprintf("name length: %d", len(*a)))
		}
		em := ui.TextInput(c.Key("email"), b).Placeholder("Email").Bind(*email)
		if em.Changed() {
			log(fmt.Sprintf("email length: %d", len(*b)))
		}
		errs := []feedback.FormError{{Message: "Name is required", Field: *name}, {Message: "Email is invalid", Field: *email}}
		feedback.ErrorList(c, errs, feedback.ErrorListOptions{})
		focus := ui.Local(root, "focus", func() string { return "" })
		now := ""
		if na.Focused() {
			now = "name"
		} else if em.Focused() {
			now = "email"
		}
		if now != *focus {
			*focus = now
			if now != "" {
				log("focused: " + now)
			}
		}
	})
}

func demoFieldError(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(6).Width(320)
	v := textState(root, "v", "ab")
	root.Children(func() {
		if ui.TextInput(c, v).Placeholder("Username").Changed() {
			log(fmt.Sprintf("length: %d", len(*v)))
		}
		msg := ""
		if len(*v) < 3 {
			msg = "Use three characters or more."
		}
		feedback.FieldError(c, msg)
	})
}

func demoHint(c *ui.Context, log func(string)) {
	feedback.Hint(c, "Passwords need eight characters or more.")
}

func demoAttention(c *ui.Context, log func(string)) {
	root := ui.Box(c).Padding(12)
	on := boolState(root, "on", true)
	root.Children(func() {
		feedback.Attention(c, *on, func() {
			if ui.Button(c, "Try the new editor").Clicked() {
				*on = false
				log("clicked")
			}
		})
	})
}

func demoBeacon(c *ui.Context, log func(string)) {
	ui.Row(c).Gap(12).Padding(12).Children(func() {
		ui.Text(c, "Export")
		if feedback.Beacon(c, "What is new", "Charts now export to PNG and SVG.").Clicked() {
			log("tip toggled")
		}
	})
}

func demoSaveStatus(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	until := ui.Local(root, "until", func() time.Time { return time.Time{} })
	saving := c.Now().Before(*until)
	if saving {
		c.After(until.Sub(c.Now()))
	}
	root.Children(func() {
		if ui.Button(c, "Edit").Clicked() {
			*until = c.Now().Add(1500 * time.Millisecond)
			log("saving")
		}
		feedback.SaveStatus(c, saving, feedback.SaveStatusOptions{})
	})
}

func demoSyncStatus(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	s := ui.Local(root, "s", func() feedback.SyncState { return feedback.Synced })
	root.Children(func() {
		ui.Row(c).Gap(8).Children(func() {
			for i, name := range []string{"Synced", "Syncing", "Failed"} {
				if ui.Button(c, name).Clicked() {
					*s = feedback.SyncState(i)
					log("state: " + name)
				}
			}
		})
		feedback.SyncStatus(c, *s, feedback.SyncStatusOptions{})
	})
}

func demoConnectionStatus(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	s := ui.Local(root, "s", func() feedback.Connection { return feedback.Online })
	root.Children(func() {
		ui.Row(c).Gap(8).Children(func() {
			for i, name := range []string{"Online", "Reconnecting", "Offline"} {
				c.Key(name)
				if ui.Button(c, "Set "+name).Clicked() {
					*s = feedback.Connection(i)
					log("state: " + name)
				}
			}
		})
		feedback.ConnectionStatus(c, *s, feedback.ConnectionStatusOptions{})
	})
}

func demoLiveDot(c *ui.Context, log func(string)) { feedback.LiveDot(c, feedback.LiveDotOptions{}) }

func demoQuotaBar(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12).Width(320)
	used := numberState(root, "used", 6)
	root.Children(func() {
		feedback.QuotaBar(c, "Seats", *used, 10, feedback.QuotaBarOptions{})
		if ui.Button(c, "Add seat").Disabled(*used >= 10).Clicked() {
			*used++
			log(fmt.Sprintf("seats: %.0f", *used))
		}
	})
}

func demoStorageBar(c *ui.Context, log func(string)) {
	t := c.Theme()
	parts := []feedback.StoragePart{{Name: "Photos", Size: 48, Color: t.Accent}, {Name: "Apps", Size: 22, Color: t.Success}, {Name: "Mail", Size: 9, Color: t.Warning}}
	ui.Column(c).Width(360).Children(func() {
		feedback.StorageBar(c, "Disk", parts, 128, feedback.StorageBarOptions{})
	})
}

func demoStrength(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(8).Width(320)
	pw := textState(root, "pw", "hunter")
	root.Children(func() {
		if ui.TextInput(c, pw).Password().Placeholder("Password").Changed() {
			log(fmt.Sprintf("length: %d", len(*pw)))
		}
		feedback.Strength(c, "Password strength", min(4, len(*pw)/3), feedback.StrengthOptions{})
	})
}

func demoLevelMeter(c *ui.Context, log func(string)) {
	root := ui.Row(c).Gap(16)
	v := numberState(root, "v", 0.6)
	root.Children(func() {
		feedback.LevelMeter(c, "Input", *v)
		if ui.Slider(c, v, 0, 1).Label("Level").Width(200).Changed() {
			log(fmt.Sprintf("level: %.2f", *v))
		}
	})
}

func demoBattery(c *ui.Context, log func(string)) {
	root := ui.Row(c).Gap(16)
	v, charging := numberState(root, "v", 0.15), boolState(root, "charging", false)
	root.Children(func() {
		feedback.Battery(c, "Battery", *v, *charging)
		if ui.Checkbox(c, charging, "Charging").Changed() {
			log(fmt.Sprintf("charging: %v", *charging))
		}
		if ui.Slider(c, v, 0, 1).Label("Level").Width(160).Changed() {
			log(fmt.Sprintf("level: %.2f", *v))
		}
	})
}

func demoSignal(c *ui.Context, log func(string)) {
	root := ui.Row(c).Gap(16)
	v := intState(root, "v", 3)
	root.Children(func() {
		feedback.Signal(c, "Wi-Fi", *v)
		if ui.Button(c, "Weaker").Disabled(*v == 0).Clicked() {
			*v--
			log(fmt.Sprintf("bars: %d", *v))
		}
		if ui.Button(c, "Stronger").Disabled(*v == 4).Clicked() {
			*v++
			log(fmt.Sprintf("bars: %d", *v))
		}
	})
}

func demoServiceStatus(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() {
		feedback.ServiceStatus(c, []feedback.Service{
			{Name: "API", State: feedback.Operational, Uptime: 0.9998},
			{Name: "Dashboard", State: feedback.Operational, Uptime: 0.9991},
			{Name: "Search", State: feedback.Degraded, Uptime: 0.973},
			{Name: "Webhooks", State: feedback.Down, Uptime: 0.912},
		}, feedback.ServiceStatusOptions{})
	})
}

func demoAnnouncementBar(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() {
		if feedback.AnnouncementBar(c, []string{"Free shipping over $50", "New colors just landed", "Members save ten percent"}, feedback.AnnouncementBarOptions{Every: 3 * time.Second}).Changed() {
			log("rotated")
		}
	})
}

func demoUndoToast(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	open := boolState(root, "open", false)
	root.Children(func() {
		if ui.Button(c, "Delete message").Clicked() {
			*open = true
			log("deleted")
		}
		toast, undo := feedback.UndoToast(c, "Message deleted", open, feedback.UndoToastOptions{})
		if undo.Clicked() {
			log("undo")
		} else if toast.Changed() {
			log("closed")
		}
	})
}

func demoUnlockToast(c *ui.Context, log func(string)) {
	root := ui.Column(c).Gap(12)
	open := boolState(root, "open", false)
	root.Children(func() {
		if ui.Button(c, "Finish lesson").Clicked() {
			*open = true
			log("unlocked")
		}
		feedback.UnlockToast(c, "Seven-day streak", "You studied every day this week.", open)
	})
}
