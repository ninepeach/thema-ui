// Package bundle merges and minifies a theme's CSS and JS as listed in
// bundle.json. Used by cmd/build and by cmd/preview's dist mode.
package bundle

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	themaui "github.com/ninepeach/thema-ui"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
)

// Build returns the assets/dist files keyed by file name:
// <name>.css  <name>.min.css  <name>.js  <name>.min.js.
// version goes into the banner; "" leaves it out.
func Build(theme string, b themaui.Bundle, version string) (map[string][]byte, error) {
	banner := fmt.Sprintf("/*! %s */\n", b.Name)
	if version != "" {
		banner = fmt.Sprintf("/*! %s v%s */\n", b.Name, version)
	}
	assets := filepath.Join(theme, "assets")
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("application/javascript", js.Minify)

	out := map[string][]byte{}
	join := func(paths []string, kind string) ([]byte, error) {
		var buf bytes.Buffer
		for _, p := range paths {
			data, err := os.ReadFile(filepath.Join(assets, filepath.FromSlash(p)))
			if err != nil {
				return nil, fmt.Errorf("bundle.json %s: %w", kind, err)
			}
			if kind == "css" {
				fmt.Fprintf(&buf, "\n/* ---- %s ---- */\n", p)
			} else {
				fmt.Fprintf(&buf, "\n// ---- %s ----\n", p)
			}
			buf.Write(bytes.TrimSpace(data))
			buf.WriteString("\n")
		}
		return buf.Bytes(), nil
	}

	if len(b.CSS) > 0 {
		body, err := join(b.CSS, "css")
		if err != nil {
			return nil, err
		}
		min, err := m.Bytes("text/css", body)
		if err != nil {
			return nil, fmt.Errorf("minify css: %w", err)
		}
		out[b.Name+".css"] = append([]byte(banner), body...)
		out[b.Name+".min.css"] = append([]byte(banner), append(min, '\n')...)
	}
	if len(b.JS) > 0 {
		body, err := join(b.JS, "js")
		if err != nil {
			return nil, err
		}
		min, err := m.Bytes("application/javascript", body)
		if err != nil {
			return nil, fmt.Errorf("minify js: %w", err)
		}
		out[b.Name+".js"] = append([]byte(banner), body...)
		out[b.Name+".min.js"] = append([]byte(banner), append(min, '\n')...)
	}
	return out, nil
}

// Write stores files in dir, creating it if needed.
func Write(dir string, files map[string][]byte) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for name, data := range files {
		if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
			return err
		}
	}
	return nil
}
