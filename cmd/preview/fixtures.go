package main

import (
	"github.com/ninepeach/thema-ui/locale"
	"github.com/ninepeach/thema-ui/ui"
)

// k shortens demo keys: k("hero.title") → ui.K("d.hero.title").
func k(key string) ui.Text { return ui.K("d." + key) }

// ---------------------------------------------------------------------------
// website demo

type Website struct {
	Intro, Rooms, About, FAQ, Closing ui.Section
	Features                          ui.Features
	Cards                             ui.CardGrid
	Media                             ui.Media
	Hours                             ui.ListCard
	Accordion                         ui.Accordion
	CTA                               ui.CTA
}

func websiteChrome(code string) (*ui.Header, *ui.Hero, *ui.Footer) {
	p := func(path string) string { return locale.Path(code, path) }
	nav := []ui.Link{
		{Label: k("nav.rooms"), URL: p("/") + "#rooms"},
		{Label: k("nav.dining"), URL: p("/") + "#about"},
		{Label: k("nav.access"), URL: p("/app/")},
		{Label: k("nav.faq"), URL: p("/") + "#faq"},
	}
	header := &ui.Header{Nav: nav, Actions: []ui.Button{{Label: k("nav.book"), URL: p("/") + "#book", Variant: "primary"}}}
	hero := &ui.Hero{
		Eyebrow: k("hero.eyebrow"),
		Title:   k("hero.title"),
		Lead:    k("hero.lead"),
		Actions: []ui.Button{
			{Label: k("hero.primary"), URL: "#book", Variant: "primary", Icon: "calendar"},
			{Label: k("hero.secondary"), URL: "#rooms", Variant: "secondary"},
		},
	}
	footer := &ui.Footer{
		About: k("footer.about"),
		Columns: []ui.FooterColumn{
			{Title: k("footer.stay"), Links: nav[:3]},
			{Title: k("footer.help"), Links: []ui.Link{
				{Label: k("nav.faq"), URL: p("/") + "#faq"},
				{Label: k("footer.contact"), URL: "mailto:hello@example.com", Icon: "mail"},
				{Label: ui.S("thema"), URL: "https://github.com/ninepeach/thema", External: true},
			}},
		},
		Bottom: ui.S("© 2026 Demo Lodge"),
		Links:  []ui.Link{{Label: k("footer.privacy"), URL: "#privacy"}},
	}
	return header, hero, footer
}

func website(code string) *Website {
	img := func(name, alt string) *ui.Image {
		return &ui.Image{Src: "images/" + name, Alt: k(alt), Width: 1200, Height: 800}
	}
	return &Website{
		Intro: ui.Section{ID: "intro", Eyebrow: k("intro.eyebrow"), Title: k("intro.title")},
		Features: ui.Features{Items: []ui.Feature{
			{Icon: "wifi", Title: k("f.wifi"), Body: k("f.wifi_body")},
			{Icon: "car", Title: k("f.parking"), Body: k("f.parking_body")},
			{Icon: "utensils", Title: k("f.breakfast"), Body: k("f.breakfast_body")},
			{Icon: "map-pin", Title: k("f.access"), Body: k("f.access_body")},
		}},
		Rooms: ui.Section{ID: "rooms", Eyebrow: k("rooms.eyebrow"), Title: k("rooms.title"), Lead: k("rooms.lead"), Tint: true},
		Cards: ui.CardGrid{Cards: []ui.Card{
			{Image: img("room.svg", "room.twin"), Title: k("room.twin"), Body: k("room.twin_body"), URL: "#twin", Badge: &ui.Badge{Label: k("badge.popular"), Tone: "info"}},
			{Image: img("lake.svg", "room.family"), Title: k("room.family"), Body: k("room.family_body"), URL: "#family"},
			{Image: img("peak.svg", "room.tatami"), Title: k("room.tatami"), Body: k("room.tatami_body"), URL: "#tatami"},
		}},
		About: ui.Section{ID: "about", Title: k("about.title")},
		Media: ui.Media{Image: *img("hero.svg", "about.caption"), Caption: k("about.caption"), Ratio: "3x2"},
		Hours: rules(),
		FAQ:   ui.Section{ID: "faq", Title: k("faq.title"), Width: "narrow"},
		Accordion: ui.Accordion{Items: []ui.AccordionItem{
			{Title: k("faq.q1"), Body: k("faq.a1"), Open: true},
			{Title: k("faq.q2"), Body: k("faq.a2")},
			{Title: k("faq.q3"), Body: k("faq.a3")},
		}},
		Closing: ui.Section{ID: "book"},
		CTA: ui.CTA{Title: k("cta.title"), Body: k("cta.body"), Actions: []ui.Button{
			{Label: k("hero.primary"), URL: "#book", Variant: "secondary", Icon: "calendar"},
			{Label: k("cta.contact"), URL: "mailto:hello@example.com", Variant: "ghost", Icon: "mail"},
		}},
	}
}

func rules() ui.ListCard {
	return ui.ListCard{ID: "house-rules", Icon: "info", Title: k("rules.title"), Rows: []ui.Row{
		{Icon: "bed", Label: k("rules.checkin"), Value: ui.S("16:00 – 21:00")},
		{Icon: "key", Label: k("rules.self"), Value: ui.KV("d.rules.after", "21:00")},
		{Icon: "logout", Label: k("rules.checkout"), Value: ui.KV("d.rules.before", "10:00")},
		{Icon: "utensils", Label: k("rules.breakfast"), Value: ui.S("07:00 – 09:00")},
	}}
}

// ---------------------------------------------------------------------------
// app demo

type App struct {
	Actions []ui.ActionCard
	WiFi    ui.InfoStrip
	Rules   ui.ListCard
	Quiet   ui.ListCard
	Notice  ui.Notice
}

func app() *App {
	return &App{
		Actions: []ui.ActionCard{
			{Icon: "clipboard", Title: k("app.checkin"), Desc: k("app.checkin_desc"), URL: "#checkin", Primary: true},
			{Icon: "logout", Title: k("app.checkout"), Desc: k("app.checkout_desc"), URL: "#checkout"},
		},
		WiFi: ui.InfoStrip{Icon: "wifi", Title: k("wifi.title"), Items: []ui.InfoItem{
			{Label: k("wifi.ssid"), Value: ui.S("DEMO")},
			{Label: k("wifi.password"), Value: k("wifi.none")},
		}},
		Rules: rules(),
		Quiet: ui.ListCard{ID: "quiet", Icon: "moon", Title: k("quiet.title"), Badge: ui.S("22:30 – 07:00"), Rows: []ui.Row{
			{Icon: "volume-off", Text: k("quiet.n1")},
			{Icon: "footsteps", Text: k("quiet.n2")},
			{Icon: "hair-dryer", Text: k("quiet.n3")},
		}},
		Notice: ui.Notice{Tone: "info", Title: k("app.notice_title"), Body: k("app.notice_body")},
	}
}

// ---------------------------------------------------------------------------
// component kit

type Kit struct {
	Header               ui.PageHeader
	Buttons              []ui.Button
	Loading, BlockButton ui.Button
	Fields               []ui.Field
	Choices              []ui.Choice
	Files                []ui.File
	Submit               ui.Button
	Notices              []ui.Notice
	Badges               []ui.Badge
	Steps                []ui.Steps
	Tabs                 ui.Tabs
	Pagination           ui.Pagination
	Back                 ui.Link
	Code, CodeLong       ui.CodeDisplay
	Empty                ui.EmptyState
	Cards                ui.CardGrid
	Minimal              ui.ListCard
	LongStrip            ui.InfoStrip
	LongActions          []ui.ActionCard
	DialogOpen           ui.Button
	Dialog               ui.Dialog
	Icons                []string
}

func kit(code string) *Kit {
	countries := []ui.Option{
		{Value: "", Label: ui.S("—")},
		{Value: "JP", Label: k("form.c_jp")},
		{Value: "HK", Label: k("form.c_hk"), Selected: true},
		{Value: "TW", Label: k("form.c_tw")},
		{Value: "XX", Label: k("form.c_other")},
	}
	page := func(n string, current bool) ui.Link {
		return ui.Link{Label: ui.S(n), URL: "?page=" + n, Current: current}
	}
	return &Kit{
		Header: ui.PageHeader{Eyebrow: ui.S("thema-ui"), Title: k("kit.title"), Lead: k("kit.lead")},
		Buttons: []ui.Button{
			{Label: k("btn.primary"), Variant: "primary"},
			{Label: k("btn.secondary")},
			{Label: k("btn.ghost"), Variant: "ghost"},
			{Label: k("btn.danger"), Variant: "danger"},
			{Label: k("btn.icon"), Icon: "upload", URL: "#buttons"},
			{Label: k("btn.disabled"), Variant: "primary", Disable: true},
		},
		Loading:     ui.Button{Label: k("btn.loading")},
		BlockButton: ui.Button{Label: k("btn.block"), Variant: "primary", Block: true, Icon: "arrow-right"},
		Fields: []ui.Field{
			{ID: "f-name", Name: "name", Label: k("form.name"), Hint: k("form.name_hint"), Required: true, Autocomplete: "name"},
			{ID: "f-email", Name: "email", Type: "email", Label: k("form.email"), Value: "guest@", Error: k("form.email_error"), Required: true},
			{ID: "f-phone", Name: "phone", Type: "tel", Label: k("form.phone"), Placeholder: "+81 90 1234 5678", Inputmode: "tel"},
			{ID: "f-country", Name: "country", Type: "select", Label: k("form.country"), Options: countries},
			{ID: "f-arrival", Name: "arrival", Type: "date", Label: k("form.arrival"), Value: "2026-12-24"},
			{ID: "f-booking", Name: "booking", Label: k("form.booking"), Value: "KM-24816", Disabled: true},
			{ID: "f-note", Name: "note", Type: "textarea", Label: k("form.note")},
		},
		Choices: []ui.Choice{
			{Name: "purpose", Type: "radio", Legend: k("form.purpose"), Required: true, Options: []ui.Option{
				{Value: "leisure", Label: k("form.p_leisure"), Selected: true},
				{Value: "business", Label: k("form.p_business")},
			}},
			{Name: "extras", Type: "checkbox", Legend: k("form.extras"), Options: []ui.Option{
				{Value: "breakfast", Label: k("form.e_breakfast"), Selected: true},
				{Value: "parking", Label: k("form.e_parking")},
				{Value: "shuttle", Label: k("form.e_shuttle"), Disabled: true},
			}},
			{Name: "agree", Type: "checkbox", Error: k("form.agree_error"), Required: true, Options: []ui.Option{
				{Value: "yes", Label: k("form.agree")},
			}},
		},
		Files: []ui.File{
			{ID: "f-passport", Name: "passport", Label: k("form.passport"), Hint: k("form.passport_hint"), Accept: "image/*", Capture: "environment", Required: true},
		},
		Submit: ui.Button{Label: k("form.submit"), Type: "submit", Variant: "primary", Block: true},
		Notices: []ui.Notice{
			{Tone: "info", Title: k("n.info_t"), Body: k("n.info_b")},
			{Tone: "success", Title: k("n.success_t")},
			{Tone: "warning", Title: k("n.warning_t")},
			{Tone: "danger", Title: k("n.danger_t"), Body: k("n.danger_b"), Action: &ui.Button{Label: k("n.retry"), Variant: "secondary"}},
		},
		Badges: []ui.Badge{
			{Label: k("b.neutral")}, {Label: k("b.info"), Tone: "info"}, {Label: k("b.success"), Tone: "success"},
			{Label: k("b.warning"), Tone: "warning"}, {Label: k("b.danger"), Tone: "danger"},
		},
		Steps: []ui.Steps{
			{Items: []ui.Step{{Label: k("s.find"), State: "done"}, {Label: k("s.register"), State: "current"}, {Label: k("s.done")}}},
			{Items: []ui.Step{{Label: k("s.find"), State: "done"}, {Label: k("s.register"), State: "done"}, {Label: k("s.done"), State: "done"}}},
		},
		Tabs: ui.Tabs{Items: []ui.Link{
			{Label: k("t.overview"), URL: "#navigation", Current: true},
			{Label: k("nav.rooms"), URL: "#navigation"},
			{Label: k("t.location"), URL: "#navigation", Icon: "map-pin"},
			{Label: k("nav.faq"), URL: "#navigation"},
		}},
		Pagination: ui.Pagination{
			Prev:  &ui.Link{URL: "?page=1"},
			Next:  &ui.Link{URL: "?page=3"},
			Pages: []ui.Link{page("1", false), page("2", true), page("3", false), page("4", false)},
		},
		Back:     ui.Link{Label: k("back"), URL: locale.Path(code, "/")},
		Code:     ui.CodeDisplay{ID: "door-code", Label: k("code.label"), Code: "4821#", Note: k("code.note"), Copy: true},
		CodeLong: ui.CodeDisplay{ID: "booking-no", Label: k("code.booking"), Code: "BK-2026-1224-00381", Copy: true},
		Empty:    ui.EmptyState{Icon: "calendar", Title: k("empty.title"), Body: k("empty.body"), Action: &ui.Button{Label: k("empty.action"), Variant: "primary"}},
		Cards: ui.CardGrid{Columns: 2, Cards: []ui.Card{
			{Eyebrow: ui.S("Text only"), Title: k("room.twin"), Body: k("room.twin_body")},
			{Title: k("long.a1"), Body: k("long.a1d"), URL: "#content", LinkLabel: k("hero.primary"), Badge: &ui.Badge{Label: k("b.warning"), Tone: "warning"}},
		}},
		Minimal: ui.ListCard{Title: k("min.title"), Rows: []ui.Row{
			{Label: ui.S("Label"), Value: ui.S("Value")},
			{Label: ui.S("Label with a much longer description"), Value: ui.S("¥12,800")},
			{Text: ui.S("A text row without icon.")},
		}},
		LongStrip: ui.InfoStrip{Icon: "car", Title: k("long.title"), Items: []ui.InfoItem{
			{Label: k("long.lot"), Value: ui.S("North lot, spaces 12–18")},
			{Label: k("long.gate"), Value: ui.S("0420-7731")},
		}},
		LongActions: []ui.ActionCard{
			{Icon: "utensils", Title: k("long.a1"), Desc: k("long.a1d"), URL: "#list", Primary: true},
			{Icon: "bed", Title: k("long.a2"), Desc: k("long.a2d"), URL: "#list"},
		},
		DialogOpen: ui.Button{Label: k("dlg.open"), Icon: "logout", Dialog: "checkout-dialog", Variant: "danger"},
		Dialog: ui.Dialog{ID: "checkout-dialog", Title: k("dlg.title"), Body: k("dlg.body"), Actions: []ui.Button{
			{Label: k("dlg.cancel"), Close: true},
			{Label: k("dlg.confirm"), Variant: "danger", Close: true},
		}},
		Icons: []string{
			"chevron-right", "chevron-left", "chevron-down", "arrow-right", "external", "menu", "x", "home", "globe",
			"check", "info", "alert", "copy", "mail", "phone", "map-pin", "calendar", "clock", "user", "upload",
			"clipboard", "logout", "wifi", "bed", "key", "utensils", "moon", "volume-off", "footsteps", "hair-dryer",
			"car", "mountain",
		},
	}
}
