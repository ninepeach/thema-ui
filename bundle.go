package themaui

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Bundle is <theme>/bundle.json: the CSS and JS files, relative to
// <theme>/assets, that cmd/build merges and minifies into
// assets/dist/<name>.min.css and assets/dist/<name>.min.js.
// A project adds its own files (e.g. "css/kurmon.css") to ship one file.
type Bundle struct {
	Name string   `json:"name"`
	CSS  []string `json:"css"`
	JS   []string `json:"js"`
}

// ReadBundle loads <themeDir>/bundle.json.
func ReadBundle(themeDir string) (Bundle, error) {
	var b Bundle
	file := filepath.Join(themeDir, "bundle.json")
	data, err := os.ReadFile(file)
	if errors.Is(err, fs.ErrNotExist) {
		return b, fmt.Errorf("%s not found: -theme must point at a thema-ui theme directory "+
			"(in the thema-ui repo: theme; in a project: e.g. themes/myproject)", file)
	}
	if err != nil {
		return b, err
	}
	if err := json.Unmarshal(data, &b); err != nil {
		return b, fmt.Errorf("bundle.json: %w", err)
	}
	if b.Name == "" {
		return b, fmt.Errorf("bundle.json: name is empty")
	}
	return b, nil
}

// Assets returns the asset paths a page should load: the source files one
// by one (source == true, for development) or the minified bundle that
// cmd/build writes to assets/dist.
// Put the result in ui.Document.CSS / ui.Document.JS.
func (b Bundle) Assets(source bool) (css, js []string) {
	if source {
		return append([]string(nil), b.CSS...), append([]string(nil), b.JS...)
	}
	if len(b.CSS) > 0 {
		css = []string{"dist/" + b.Name + ".min.css"}
	}
	if len(b.JS) > 0 {
		js = []string{"dist/" + b.Name + ".min.js"}
	}
	return css, js
}
