// Package examples embeds the example project styles (lodge, forest, ink)
// for cmd/preview. Each is one project CSS file: copy it into a project's
// assets/css/ as a starting point.
package examples

import "embed"

//go:embed *.css
var FS embed.FS

// Names lists the examples in display order; the first is the default.
var Names = []string{"lodge", "forest", "ink"}
