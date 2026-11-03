// Package themaui embeds the thema-ui base theme and reads its bundle.json.
//
// thema-ui is a starter kit: a project copies theme/, ui/ and locale/ with
// cp -R and owns the copy from then on. This package keeps only what the
// preview and projects need: FS, ReadBundle and Bundle.Assets (bundle.go,
// small enough to copy into a project too).
package themaui

import (
	"embed"
	"io/fs"
)

//go:embed all:theme
var embedded embed.FS

// FS returns the embedded theme (rooted at theme/).
func FS() fs.FS {
	sub, err := fs.Sub(embedded, "theme")
	if err != nil {
		panic(err)
	}
	return sub
}
