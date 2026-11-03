package main

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	themaui "github.com/ninepeach/thema-ui"
	"github.com/ninepeach/thema-ui/internal/bundle"
)

// A project adds its own CSS to the bundle and gets one file.
func TestProjectBundle(t *testing.T) {
	theme := copyTheme(t)
	must(t, os.MkdirAll(filepath.Join(theme, "assets", "css"), 0o755))
	must(t, os.WriteFile(filepath.Join(theme, "assets", "css", "proj.css"), []byte(".proj-only { color: red }"), 0o644))
	must(t, os.WriteFile(filepath.Join(theme, "theme.json"), []byte(`{"name":"proj","version":"2.3.4","thema":"0.1"}`), 0o644))
	b, _ := themaui.ReadBundle(theme)
	b.Name = "app"
	b.CSS = append(b.CSS, "css/proj.css")
	files, err := bundle.Build(theme, b, readVersion(theme))
	if err != nil {
		t.Fatal(err)
	}
	min := string(files["app.min.css"])
	if !strings.HasPrefix(min, "/*! app v2.3.4 */") || !strings.Contains(min, ".proj-only{color:red}") || !strings.Contains(min, ".ui-button") {
		t.Fatalf("unexpected bundle: %.200s", min)
	}
	if css, js := b.Assets(false); css[0] != "dist/app.min.css" || js[0] != "dist/app.min.js" {
		t.Fatalf("Assets(false) = %v %v", css, js)
	}
}

// build writes dist; -check passes afterwards and fails after an edit.
func TestBuildAndCheck(t *testing.T) {
	theme := copyTheme(t)
	run := func(args ...string) error {
		cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	must(t, run("-theme", theme))
	must(t, run("-theme", theme, "-check"))
	css := filepath.Join(theme, "assets", "src", "css", "style.css")
	must(t, os.WriteFile(css, []byte(".changed{}"), 0o644))
	if run("-theme", theme, "-check") == nil {
		t.Fatal("-check passed after an edit")
	}
}

// copyTheme does what a project does: cp -R theme.
func copyTheme(t *testing.T) string {
	t.Helper()
	dst := filepath.Join(t.TempDir(), "proj")
	must(t, copyFS(themaui.FS(), dst))
	return dst
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
