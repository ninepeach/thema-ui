// Command build merges and minifies a theme's CSS and JS as listed in
// <theme>/bundle.json, writing <theme>/assets/dist/:
//
//	<name>.css  <name>.min.css  <name>.js  <name>.min.js
//
// A project runs it before deploying:
//
//	go run github.com/ninepeach/thema-ui/cmd/build@v0.1.0 -theme ./themes/myproject
//	go run github.com/ninepeach/thema-ui/cmd/build@v0.1.0 -theme ./themes/myproject -check
//
// Banner version: -version, else "version" in <theme>/theme.json, else none.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	themaui "github.com/ninepeach/thema-ui"
	"github.com/ninepeach/thema-ui/internal/bundle"
)

func main() {
	theme := flag.String("theme", "theme", "theme directory (contains bundle.json and assets/)")
	version := flag.String("version", "", "version for the banner (default: theme.json)")
	check := flag.Bool("check", false, "report whether assets/dist is up to date, write nothing")
	flag.Parse()

	b, err := themaui.ReadBundle(*theme)
	if err != nil {
		fail(err)
	}
	if *version == "" {
		*version = readVersion(*theme)
	}
	files, err := bundle.Build(*theme, b, *version)
	if err != nil {
		fail(err)
	}
	dist := filepath.Join(*theme, "assets", "dist")

	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)

	if *check {
		stale := 0
		for _, name := range names {
			got, err := os.ReadFile(filepath.Join(dist, name))
			if err != nil || !bytes.Equal(got, files[name]) {
				fmt.Println("stale assets/dist/" + name)
				stale++
			}
		}
		if stale > 0 {
			fmt.Printf("run: go run github.com/ninepeach/thema-ui/cmd/build -theme %s\n", *theme)
			os.Exit(1)
		}
		fmt.Printf("assets/dist up to date (%s)\n", b.Name)
		return
	}

	if err := bundle.Write(dist, files); err != nil {
		fail(err)
	}
	for _, name := range names {
		fmt.Printf("assets/dist/%-22s %7d bytes\n", name, len(files[name]))
	}
}

func readVersion(theme string) string {
	var manifest struct {
		Version string `json:"version"`
	}
	if data, err := os.ReadFile(filepath.Join(theme, "theme.json")); err == nil {
		if json.Unmarshal(data, &manifest) == nil {
			return manifest.Version
		}
	}
	return ""
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
