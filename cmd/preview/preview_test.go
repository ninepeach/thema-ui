package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/ninepeach/thema-ui/locale"
)

// Every preview page renders in every locale, layer and example — this
// executes every thema-ui template with real data.
func TestPreviewPagesRender(t *testing.T) {
	s, err := newServer("../..")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(s.repo)
	h := s.handler()

	for _, l := range locale.Languages {
		for _, pg := range pages {
			for _, q := range []string{"", "?layers=base", "?layers=tokens", "?layers=style", "?assets=dist", "?example=forest", "?example=ink"} {
				url := locale.Path(l.Code, pg.path) + q
				rec := httptest.NewRecorder()
				h.ServeHTTP(rec, httptest.NewRequest("GET", url, nil))
				if rec.Code != http.StatusOK {
					t.Fatalf("%s: %d %s", url, rec.Code, rec.Body.String())
				}
				if strings.Contains(rec.Body.String(), "d.") && strings.Contains(rec.Body.String(), ">d.") {
					t.Errorf("%s: untranslated key in output", url)
				}
			}
		}
	}
	get := func(url string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest("GET", url, nil))
		return rec
	}
	if rec := get("/"); rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `value="forest"`) {
		t.Fatalf("home: %d", rec.Code)
	}
	if rec := get("/open?example=ink&page=/app/&lang=ko"); rec.Header().Get("Location") != "/ko/app/?layers=example&example=ink" {
		t.Fatalf("open: %v", rec.Header())
	}
	if rec := get("/en/?example=ink"); !strings.Contains(rec.Body.String(), "examples/ink.css") {
		t.Fatal("example CSS not loaded")
	}
	if rec := get("/en/"); !strings.Contains(rec.Body.String(), `class="ui-scroll-progress" data-scroll-progress`) {
		t.Fatal("scroll progress bar not rendered")
	}
	for _, asset := range []string{"examples/forest.css", "dist/thema-ui.min.css", "dist/thema-ui.min.js"} {
		if rec := get("/assets/" + themeID + "/" + asset); rec.Code != http.StatusOK || rec.Body.Len() == 0 {
			t.Fatalf("%s: %d", asset, rec.Code)
		}
	}
	rec := get("/frames/?page=/ja/app/")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `src="/ja/app/"`) {
		t.Fatalf("frames: %d", rec.Code)
	}
}
