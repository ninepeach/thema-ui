package themaui

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"testing"

	"github.com/ninepeach/thema"
)

// A copied theme plus a theme.json compiles as a Thema theme on its own.
func TestCopiedThemeCompiles(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "demo")
	must(t, copyFS(FS(), dst))
	must(t, os.WriteFile(filepath.Join(dst, "theme.json"), []byte(`{"name":"demo","version":"0.1.0","thema":"0.1"}`), 0o644))
	if _, err := thema.New(filepath.Dir(dst), filepath.Base(dst)); err != nil {
		t.Fatal(err)
	}
}

func TestBundleAssets(t *testing.T) {
	b, err := ReadBundle("theme")
	if err != nil {
		t.Fatal(err)
	}
	if css, js := b.Assets(true); len(css) != len(b.CSS) || css[0] != "src/css/default.css" || js[0] != "src/js/thema-ui.js" {
		t.Fatalf("Assets(true) = %v %v", css, js)
	}
	if css, js := b.Assets(false); len(css) != 1 || css[0] != "dist/thema-ui.min.css" || js[0] != "dist/thema-ui.min.js" {
		t.Fatalf("Assets(false) = %v %v", css, js)
	}
}

// Each CSS layer works on its own: default + common (+ pages) render with
// common's neutral defaults, so every variable they read must be defined in
// common.css. tokens.css and the examples only override known variables.
func TestCSSVariables(t *testing.T) {
	src := "theme/assets/src/css/"
	common := cssVars(t, src+"common.css")
	known := union(common.defs, cssVars(t, src+"style.css").defs)
	for _, layer := range []struct {
		file string
		defs map[string]bool
	}{
		{"common.css", common.defs},
		{"pages.css", common.defs},
		{"style.css", common.defs},
		{"tokens.css", known},
	} {
		v := cssVars(t, src+layer.file)
		missing(t, layer.file+" reads", v.uses, union(layer.defs, v.defs))
	}
	missing(t, "tokens.css defines", cssVars(t, src+"tokens.css").defs, known)

	examples, _ := filepath.Glob("examples/*.css")
	for _, file := range examples {
		v := cssVars(t, file)
		missing(t, file+" defines", v.defs, union(known, v.uses))
	}
}

type vars struct{ uses, defs map[string]bool }

var (
	useRE = regexp.MustCompile(`var\((--[a-z0-9-]+)\)`) // var(--x, fallback) is self-contained
	defRE = regexp.MustCompile(`(?:^|[;{\s])(--[a-z0-9-]+)\s*:`)
)

func cssVars(t *testing.T, file string) vars {
	t.Helper()
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	v := vars{map[string]bool{}, map[string]bool{}}
	for _, m := range useRE.FindAllStringSubmatch(string(data), -1) {
		v.uses[m[1]] = true
	}
	for _, m := range defRE.FindAllStringSubmatch(string(data), -1) {
		v.defs[m[1]] = true
	}
	return v
}

func union(sets ...map[string]bool) map[string]bool {
	out := map[string]bool{}
	for _, s := range sets {
		for k := range s {
			out[k] = true
		}
	}
	return out
}

func missing(t *testing.T, what string, names, defined map[string]bool) {
	t.Helper()
	var out []string
	for n := range names {
		if !defined[n] {
			out = append(out, n)
		}
	}
	sort.Strings(out)
	if len(out) > 0 {
		t.Errorf("%s undefined variables: %v", what, out)
	}
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
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
