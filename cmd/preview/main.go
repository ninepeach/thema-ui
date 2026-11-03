// Command preview serves the thema-ui component gallery for debugging.
//
//	go run ./cmd/preview                 # from a thema-ui checkout: live edits
//	go run github.com/ninepeach/thema-ui/cmd/preview   # anywhere: embedded copy
//
// Pages: / home (pick example, page, language) · /{locale}/ website demo ·
// /{locale}/app/ app demo · /{locale}/components/ every component and state ·
// /frames/ several widths.
// Toolbar switches: example style (examples/*.css), assets (src files / dist
// built on the fly) and layers (base → +tokens → +style → +example, src only).
package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ninepeach/thema"
	themaui "github.com/ninepeach/thema-ui"
	"github.com/ninepeach/thema-ui/examples"
	"github.com/ninepeach/thema-ui/internal/bundle"
	"github.com/ninepeach/thema-ui/locale"
	"github.com/ninepeach/thema-ui/ui"
)

//go:embed all:theme
var overlayFS embed.FS

const themeID = "preview"

type server struct {
	base, overlay fs.FS // thema-ui theme, preview overlay
	examples      fs.FS // examples/*.css
	repo, dir     string
	mu            sync.Mutex
	views         *thema.Engine
	live          bool
}

func main() {
	addr := flag.String("addr", ":8090", "listen address")
	src := flag.String("src", ".", "thema-ui checkout for live editing (falls back to the embedded copy)")
	flag.Parse()

	s, err := newServer(*src)
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(s.repo)

	log.Printf("thema-ui preview on http://localhost%s/  (live=%v)", *addr, s.live)
	srv := &http.Server{Addr: *addr, Handler: s.handler(), ReadHeaderTimeout: 5 * time.Second}
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

// newServer prefers a thema-ui checkout at src (live editing) and falls
// back to the embedded copy.
func newServer(src string) (*server, error) {
	s := &server{base: themaui.FS(), examples: examples.FS}
	s.overlay, _ = fs.Sub(overlayFS, "theme")
	if info, err := os.Stat(filepath.Join(src, "theme", "assets", "src", "css", "common.css")); src != "" && err == nil && !info.IsDir() {
		s.base = os.DirFS(filepath.Join(src, "theme"))
		s.overlay = os.DirFS(filepath.Join(src, "cmd", "preview", "theme"))
		s.examples = os.DirFS(filepath.Join(src, "examples"))
		s.live = true
	}
	repo, err := os.MkdirTemp("", "thema-ui-preview-")
	if err != nil {
		return nil, err
	}
	s.repo, s.dir = repo, filepath.Join(repo, themeID)
	if err := s.build(); err != nil {
		return nil, err
	}
	if s.views, err = thema.New(repo, themeID); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *server) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /assets/"+themeID+"/", s.assets)
	mux.HandleFunc("GET /{$}", s.home)
	mux.HandleFunc("GET /open", s.open)
	mux.HandleFunc("GET /frames/", s.frames)
	mux.HandleFunc("GET /", s.page)
	return mux
}

// build assembles one theme directory the way a project would: the
// thema-ui theme copied in, the preview pages and d.* locale keys on top,
// examples under assets/examples/ and a freshly built assets/dist/.
func (s *server) build() error {
	if err := os.RemoveAll(s.dir); err != nil {
		return err
	}
	if err := copyFS(s.overlay, s.dir); err != nil {
		return err
	}
	if err := copyTheme(s.base, s.dir); err != nil {
		return err
	}
	if err := copyFS(s.examples, filepath.Join(s.dir, "assets", "examples")); err != nil {
		return err
	}
	b, err := themaui.ReadBundle(s.dir)
	if err != nil {
		return err
	}
	files, err := bundle.Build(s.dir, b, "")
	if err != nil {
		return err
	}
	return bundle.Write(filepath.Join(s.dir, "assets", "dist"), files)
}

// refresh rebuilds from disk on every request in live mode.
func (s *server) refresh(ctx context.Context) {
	if !s.live {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.build(); err != nil {
		log.Printf("build: %v", err)
		return
	}
	if _, err := s.views.Refresh(ctx); err != nil {
		log.Printf("refresh: %v", err)
	}
}

// ---------------------------------------------------------------------------
// settings: example, assets and layers, from ?query and remembered in
// cookies so the iframes on /frames/ follow the same choice.

type settings struct{ example, assets, layers string }

var (
	assetNames = []string{"src", "dist"}
	layerNames = []string{"base", "tokens", "style", "example"}
)

func readSettings(w http.ResponseWriter, r *http.Request) settings {
	pick := func(name, fallback string, allowed []string) string {
		value := r.URL.Query().Get(name)
		if value == "" {
			if c, err := r.Cookie("pv_" + name); err == nil {
				value = c.Value
			}
		}
		for _, a := range allowed {
			if a == value {
				if w != nil {
					http.SetCookie(w, &http.Cookie{Name: "pv_" + name, Value: value, Path: "/", SameSite: http.SameSiteLaxMode})
				}
				return value
			}
		}
		return fallback
	}
	return settings{
		example: pick("example", examples.Names[0], examples.Names),
		assets:  pick("assets", "src", assetNames),
		layers:  pick("layers", "example", layerNames),
	}
}

func (s *server) assets(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/assets/"+themeID+"/")
	w.Header().Set("Cache-Control", "no-store")
	http.ServeFile(w, r, filepath.Join(s.dir, "assets", filepath.FromSlash(name)))
}

// ---------------------------------------------------------------------------
// pages

type Page struct {
	ui.Document
	Toolbar *Toolbar
	Home    *Home
	Website *Website
	App     *App
	Kit     *Kit
	Frames  []Frame
}

type Toolbar struct {
	Examples, Pages, Assets, Layers []ToolLink
	FramesURL                       string
}

// Home is the start page: dropdowns for example, page and language.
type Home struct{ Examples, Pages, Langs []Option }

type Option struct {
	Value, Label string
	Selected     bool
}

type ToolLink struct {
	Label, URL string
	Current    bool
}

type Frame struct {
	Width int
	URL   string
}

var pages = []struct{ path, label string }{{"/", "website"}, {"/app/", "app"}, {"/components/", "components"}}

func (s *server) home(w http.ResponseWriter, r *http.Request) {
	s.refresh(r.Context())
	set := readSettings(w, r)
	page := s.newPage(locale.Default, "/", set)
	home := &Home{}
	for _, e := range examples.Names {
		home.Examples = append(home.Examples, Option{Value: e, Label: e, Selected: e == set.example})
	}
	for _, pg := range pages {
		home.Pages = append(home.Pages, Option{Value: pg.path, Label: pg.label})
	}
	for _, l := range locale.Languages {
		home.Langs = append(home.Langs, Option{Value: l.Code, Label: l.Name})
	}
	page.Home = home
	s.render(w, r, "pages/home", page)
}

// open sends the home page form to the chosen page with the example on.
func (s *server) open(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	code, pagePath, example := locale.Default, "/", examples.Names[0]
	for _, e := range examples.Names {
		if e == q.Get("example") {
			example = e
		}
	}
	for _, l := range locale.Languages {
		if l.Code == q.Get("lang") {
			code = l.Code
		}
	}
	for _, pg := range pages {
		if pg.path == q.Get("page") {
			pagePath = pg.path
		}
	}
	target := locale.Path(code, pagePath) + "?layers=example&example=" + example
	http.Redirect(w, r, target, http.StatusFound)
}

func (s *server) page(w http.ResponseWriter, r *http.Request) {
	code, pagePath, redirect := locale.Resolve(r.URL.Path)
	if redirect != "" {
		target := redirect
		if r.URL.RawQuery != "" {
			target += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, target, http.StatusFound)
		return
	}
	s.refresh(r.Context())
	set := readSettings(w, r)
	name, page, ok := s.buildPage(code, pagePath, set)
	if !ok {
		http.NotFound(w, r)
		return
	}
	s.render(w, r, name, page)
}

// buildPage returns the template and data for a preview page.
func (s *server) buildPage(code, pagePath string, set settings) (string, Page, bool) {
	page := s.newPage(code, pagePath, set)
	switch pagePath {
	case "/":
		page.Title = ui.K("d.title.website")
		page.Header, page.Hero, page.Footer = websiteChrome(code)
		page.Website = website(code)
		return "pages/website", page, true
	case "/app/":
		page.Title = ui.K("d.title.app")
		page.BodyClass = "ui-app"
		page.Hero = &ui.Hero{ShowBrand: true, Tagline: ui.K("d.app.tagline")}
		page.App = app()
		return "pages/app", page, true
	case "/components/":
		page.Title = ui.K("d.title.kit")
		page.Header, _, page.Footer = websiteChrome(code)
		page.Kit = kit(code)
		page.Breadcrumb = []ui.Link{{Label: ui.K("d.kit.title")}}
		return "pages/components", page, true
	}
	return "", page, false
}

func (s *server) frames(w http.ResponseWriter, r *http.Request) {
	s.refresh(r.Context())
	set := readSettings(w, r)
	target := r.URL.Query().Get("page")
	if target == "" || !strings.HasPrefix(target, "/") {
		target = "/en/"
	}
	code, pagePath, _ := locale.Resolve(target)
	page := s.newPage(code, pagePath, set)
	for _, width := range []int{340, 390, 768, 1280} {
		page.Frames = append(page.Frames, Frame{Width: width, URL: target})
	}
	s.render(w, r, "pages/frames", page)
}

func (s *server) render(w http.ResponseWriter, r *http.Request, name string, page Page) {
	var buf bytes.Buffer
	if err := s.views.Render(r.Context(), &buf, name, page, thema.WithLocale(page.Locale)); err != nil {
		log.Printf("render %s: %v", name, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = buf.WriteTo(w)
}

func (s *server) newPage(code, pagePath string, set settings) Page {
	css, js := s.bundleAssets(set)
	p := Page{Toolbar: &Toolbar{}, Document: ui.Document{
		Locale:         code,
		ScrollProgress: true,
		HomeURL:        locale.Path(code, "/"),
		Brand:          ui.Brand{Name: "Demo Lodge"},
		Lang:           locale.Current(code, pagePath),
		Langs:          locale.Options(code, pagePath),
		CSS:            css,
		JS:             js,
	}}
	link := func(label, query string, current bool) ToolLink {
		return ToolLink{Label: label, URL: locale.Path(code, pagePath) + "?" + query, Current: current}
	}
	for _, e := range examples.Names {
		p.Toolbar.Examples = append(p.Toolbar.Examples, link(e, "layers=example&example="+e, e == set.example && set.layers == "example"))
	}
	for _, pg := range pages {
		p.Toolbar.Pages = append(p.Toolbar.Pages, ToolLink{Label: pg.label, URL: locale.Path(code, pg.path), Current: pg.path == pagePath})
	}
	for _, a := range assetNames {
		p.Toolbar.Assets = append(p.Toolbar.Assets, link(a, "assets="+a, a == set.assets))
	}
	for _, l := range layerNames {
		p.Toolbar.Layers = append(p.Toolbar.Layers, link(l, "layers="+l, l == set.layers))
	}
	p.Toolbar.FramesURL = "/frames/?page=" + url.QueryEscape(locale.Path(code, pagePath))
	return p
}

// bundleAssets: the theme bundle (src files or dist build), then the
// preview toolbar (chrome.css) and the example style as project CSS.
// layers: base = default+common+pages, tokens = + tokens.css,
// style = + style.css, example = + examples/<name>.css.
// In dist mode tokens.css and style.css are in the bundle.
func (s *server) bundleAssets(set settings) (css, js []string) {
	b, err := themaui.ReadBundle(s.dir)
	if err != nil {
		log.Printf("bundle: %v", err)
	}
	all, js := b.Assets(set.assets == "src")
	for _, c := range all {
		if c == "src/css/tokens.css" && set.layers == "base" ||
			c == "src/css/style.css" && (set.layers == "base" || set.layers == "tokens") {
			continue
		}
		css = append(css, c)
	}
	css = append(css, "css/chrome.css")
	if set.layers == "example" {
		css = append(css, "examples/"+set.example+".css")
	}
	return css, js
}

// copyTheme copies the thema-ui theme over the overlay. Locale files are
// merged (one file per language in Thema): ui.* keys join the preview's d.* keys.
func copyTheme(src fs.FS, dst string) error {
	return fs.WalkDir(src, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || path.Base(p) == ".DS_Store" {
			return err
		}
		data, err := fs.ReadFile(src, p)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, filepath.FromSlash(p))
		if strings.HasPrefix(p, "locales/") {
			if data, err = mergeLocale(target, data); err != nil {
				return fmt.Errorf("%s: %w", p, err)
			}
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}

func mergeLocale(target string, base []byte) ([]byte, error) {
	overlay, err := os.ReadFile(target)
	if errors.Is(err, fs.ErrNotExist) {
		return base, nil
	}
	if err != nil {
		return nil, err
	}
	m := map[string]any{}
	for _, data := range [][]byte{overlay, base} {
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, err
		}
	}
	return json.MarshalIndent(m, "", "  ")
}

func copyFS(src fs.FS, dst string) error {
	return fs.WalkDir(src, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		target := filepath.Join(dst, filepath.FromSlash(p))
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := fs.ReadFile(src, p)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}
