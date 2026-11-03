// Package ui holds the data for every thema-ui template. Each component
// under templates/components/<group>/<name>.html takes one struct from
// this package; pages compose components by filling structs.
//
// All visible text is a Text: either a locale key (translated in the
// template) or a literal value (already localised, e.g. from a database).
package ui

import "github.com/ninepeach/thema-ui/locale"

// ---------------------------------------------------------------------------
// Text

// Text is rendered by components/global/text:
//
//	Key == ""          → Value as-is
//	Key != ""          → {{ t Key "value" Value }}   (Value fills {{value}})
type Text struct{ Key, Value string }

// K is a translated text: K("home.title").
func K(key string) Text { return Text{Key: key} }

// S is a literal text: S("Room 201").
func S(value string) Text { return Text{Value: value} }

// KV is a translated text with {{value}}: KV("rules.after", "21:00").
func KV(key, value string) Text { return Text{Key: key, Value: value} }

// Set reports whether the text has anything to render.
func (t Text) Set() bool { return t.Key != "" || t.Value != "" }

// ---------------------------------------------------------------------------
// Document (template root)

// Document is the root of every page. Projects embed it in their own
// page struct; components that need page context (header, footer, hero,
// breadcrumb, language switcher) read it from the root.
type Document struct {
	Locale    string
	Title     Text
	Theme     string // <html data-theme>: selects the project palette
	BodyClass string
	// ScrollProgress shows components/global/scroll-progress (needs thema-ui.js).
	ScrollProgress bool
	HomeURL        string
	Brand          Brand
	Lang           locale.Option
	Langs          []locale.Option
	Header         *Header // nil: no site header
	Hero           *Hero   // nil: no hero
	Footer         *Footer // nil: no site footer
	Breadcrumb     []Link  // empty: no breadcrumb
	// CSS and JS are asset paths loaded in order, usually from
	// themaui.Bundle.Assets: source files in development, the minified
	// bundle in production. Empty → dist/thema-ui.min.css / .min.js.
	CSS []string
	JS  []string
}

// Brand → components/global/logo. Logo is a theme asset path ("" = text only).
type Brand struct{ Name, Logo string }

// ---------------------------------------------------------------------------
// Shared building blocks

type Link struct {
	Label    Text
	URL      string
	Icon     string
	Current  bool
	External bool
}

type Image struct {
	Src           string // theme asset path, or absolute URL when Remote
	Alt           Text
	Width, Height int
	Remote        bool
}

// Button → components/form/button. URL set → <a>, otherwise <button>.
type Button struct {
	Label   Text
	URL     string
	Type    string // submit | button (default button)
	Variant string // primary | secondary | ghost | danger (default secondary)
	Icon    string
	Name    string
	Value   string
	Block   bool // full width
	Disable bool
	Dialog  string // id of a dialog to open
	Close   bool   // closes the enclosing dialog
}

// Tone is shared by notice and badge: info | success | warning | danger | neutral.
type Tone string

// ---------------------------------------------------------------------------
// layout

// Header → components/layout/site-header (root dot).
type Header struct {
	Nav     []Link
	Actions []Button
	Overlay bool // transparent, on top of the hero image
}

// Footer → components/layout/site-footer (root dot).
type Footer struct {
	About   Text
	Columns []FooterColumn
	Bottom  Text // e.g. copyright
	Links   []Link
}

type FooterColumn struct {
	Title Text
	Links []Link
}

// Hero → components/layout/hero (root dot).
// ShowBrand puts logo + language switcher in the hero (app pages without
// a site header).
type Hero struct {
	Eyebrow, Title, Lead Text
	Tagline              Text
	Actions              []Button
	ShowBrand            bool
	Compact              bool
	Center               bool
}

// PageHeader → components/layout/page-header.
type PageHeader struct {
	Eyebrow, Title, Lead Text
}

// Section → components/layout/section-open / section-close.
type Section struct {
	ID                   string
	Eyebrow, Title, Lead Text
	Tint                 bool   // alternate background
	Width                string // "" (site width) | narrow
}

// ---------------------------------------------------------------------------
// navigation

type Tabs struct{ Items []Link }

type Steps struct{ Items []Step }

type Step struct {
	Label Text
	State string // done | current | todo
}

type Pagination struct {
	Prev, Next *Link
	Pages      []Link
}

// ---------------------------------------------------------------------------
// content

type Card struct {
	Image       *Image
	Eyebrow     Text
	Title, Body Text
	URL         string
	LinkLabel   Text // empty → ui.learn_more
	Badge       *Badge
}

type CardGrid struct {
	Cards   []Card
	Columns int // 2 or 3 at wide widths (default 3)
}

type ActionCard struct {
	Icon        string
	Title, Desc Text
	URL         string
	Primary     bool
}

type InfoStrip struct {
	Icon  string
	Title Text
	Items []InfoItem
}

type InfoItem struct{ Label, Value Text }

// ListCard rows are label/value rows, or text rows when Text is set.
type ListCard struct {
	ID, Icon   string
	Title, Sub Text
	Badge      Text
	Rows       []Row
}

type Row struct {
	Icon              string
	Label, Sub, Value Text
	Text              Text
}

type Features struct{ Items []Feature }

type Feature struct {
	Icon        string
	Title, Body Text
}

type Media struct {
	Image   Image
	Caption Text
	Ratio   string // CSS aspect-ratio, e.g. "16 / 9" (default 4 / 3)
}

type Accordion struct{ Items []AccordionItem }

type AccordionItem struct {
	Title, Body Text
	Open        bool
}

type CTA struct {
	Title, Body Text
	Actions     []Button
}

type EmptyState struct {
	Icon        string
	Title, Body Text
	Action      *Button
}

// CodeDisplay shows a value to read or copy (door code, booking number).
type CodeDisplay struct {
	ID    string
	Label Text
	Code  string
	Note  Text
	Copy  bool
}

// ---------------------------------------------------------------------------
// form

// Field → components/form/field.
// Type: text | email | tel | number | date | time | password | url | search | select | textarea.
type Field struct {
	ID, Name, Type     string
	Label, Hint, Error Text // Error set → invalid state
	Value, Placeholder string
	Autocomplete       string
	Inputmode          string
	Required, Disabled bool
	Options            []Option // select
}

type Option struct {
	Value    string
	Label    Text
	Selected bool
	Disabled bool
}

// Choice → components/form/choice: checkbox or radio group.
type Choice struct {
	Name, Type          string // checkbox | radio
	Legend, Hint, Error Text
	Required            bool
	Options             []Option // Selected = checked
}

// File → components/form/file.
type File struct {
	ID, Name           string
	Label, Hint, Error Text
	Accept, Capture    string // capture: environment | user (mobile camera)
	Multiple, Required bool
}

// ---------------------------------------------------------------------------
// feedback & overlay

type Notice struct {
	Tone        Tone
	Title, Body Text
	Action      *Button
}

type Badge struct {
	Label Text
	Tone  Tone
}

// Dialog → components/overlay/dialog. Open it with Button{Dialog: ID}.
type Dialog struct {
	ID          string
	Title, Body Text
	Actions     []Button // use Close: true for cancel
}
